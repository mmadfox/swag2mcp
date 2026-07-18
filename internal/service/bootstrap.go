package service

import (
	"context"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mmadfox/swag2mcp/internal/auth"
	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/httpclient"
	"github.com/mmadfox/swag2mcp/internal/id"
	"github.com/mmadfox/swag2mcp/internal/model"
	"github.com/mmadfox/swag2mcp/internal/spec"
	"github.com/mmadfox/swag2mcp/internal/workspace"
)

// BootstrapRequest is the request for the Bootstrap method.
type BootstrapRequest struct {
	ConfFilePath string
	Tags         []string
}

// Bootstrap loads the configuration, initializes the workspace, creates the
// global HTTP client, and indexes all specs, collections, tags, and endpoints.
func (s *Service) Bootstrap(ctx context.Context, request BootstrapRequest) error {
	s.startedAt = time.Now()

	cfg, err := s.loadConfiguration(request.ConfFilePath, request.Tags)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	s.config = cfg
	s.disableRateLimiter.Store(cfg.DisableRateLimiter)

	if err := s.initializeWorkspace(filepath.Dir(request.ConfFilePath)); err != nil {
		return err
	}

	httpCfg := buildGlobalHTTPConfig(cfg.HTTPClient)
	if httpCfg.Randomize {
		httpclient.RandomizeConfig(&httpCfg)
	}

	client, err := httpclient.New(httpCfg)
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %w", err)
	}
	s.httpClient = client
	s.httpClientConfig = httpCfg
	s.maxResponseSize = resolveMaxResponseSize(httpCfg.MaxResponseSize)
	s.globalHeaders = httpCfg.Headers
	s.globalUserAgent = httpCfg.UserAgent
	s.globalCookies = httpCfg.Cookies

	httpclient.SetGlobalConfig(httpCfg)

	filter := config.NewFilter(request.Tags)

	for sc := range cfg.Iterate(filter) {
		if err := s.processSpec(ctx, sc, cfg.MockEnabled, cfg.MockAuth); err != nil {
			return err
		}
	}

	s.buildSnapshot()

	return nil
}

func buildGlobalHTTPConfig(global *config.GlobalHTTPClientConfig) httpclient.Config {
	if global == nil {
		return httpclient.Config{
			UserAgent: defaultUserAgent,
		}
	}

	cfg := httpclient.Config{
		Randomize:       global.Randomize,
		UserAgent:       global.UserAgent,
		Timeout:         global.Timeout,
		FollowRedirects: global.FollowRedirects,
		MaxRedirects:    global.MaxRedirects,
		MaxResponseSize: global.MaxResponseSize,
	}

	if cfg.UserAgent == "" && !cfg.Randomize {
		cfg.UserAgent = defaultUserAgent
	}
	if global.Headers != nil {
		cfg.Headers = make(map[string]string, len(global.Headers))
		maps.Copy(cfg.Headers, global.Headers)
	}
	if len(global.Cookies) > 0 {
		cfg.Cookies = make([]httpclient.Cookie, len(global.Cookies))
		for i, cookie := range global.Cookies {
			cfg.Cookies[i] = httpclient.Cookie{
				Name:     cookie.Name,
				Value:    cookie.Value,
				Domain:   cookie.Domain,
				Path:     cookie.Path,
				Secure:   cookie.Secure,
				HTTPOnly: cookie.HTTPOnly,
			}
		}
	}
	if global.Proxy != nil {
		cfg.Proxy = &httpclient.ProxyConfig{
			URL:      global.Proxy.URL,
			Username: global.Proxy.Username,
			Password: global.Proxy.Password,
			Bypass:   append([]string{}, global.Proxy.Bypass...),
		}
	}
	return cfg
}

func (s *Service) loadConfiguration(configFilepath string, tags []string) (*config.Config, error) {
	cfg, err := config.Load(configFilepath)
	if err != nil {
		return nil, err
	}

	filter := config.NewFilter(tags)
	if err := cfg.Validate(filter); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (s *Service) initializeWorkspace(wsDir string) error {
	if wsDir != "" && wsDir != s.ws.Root() {
		ws, err := workspace.New(wsDir)
		if err != nil {
			return fmt.Errorf("failed to create workspace: %w", err)
		}
		s.ws = ws
	}

	if err := s.ws.Init(); err != nil {
		return fmt.Errorf("failed to init workspace: %w", err)
	}

	s.cache.SetWorkspaceDir(s.ws.Root())
	return nil
}

func (s *Service) processSpec(ctx context.Context, sc *config.Spec, mockEnabled bool, ma *config.MockAuthConfig) error {
	sp, err := s.buildSpecInfo(sc, mockEnabled, ma)
	if err != nil {
		return err
	}

	tags := make(map[string]*model.Tag)
	colls := make(map[string]*model.Collection)
	eps := make(map[string]*model.Endpoint)

	for i := range sc.Collections {
		cc := &sc.Collections[i]
		if cc.Disable {
			continue
		}

		coll, err := s.processCollection(
			ctx, sp, sc, cc,
			tags, eps,
		)
		if err != nil {
			return err
		}

		colls[coll.ID] = coll
		sp.Stats.Collections++
	}

	sp.Stats.Tags = len(tags)
	sp.Stats.Methods = len(eps)

	return s.indexSpec(sp, colls, tags, eps)
}

func (s *Service) processCollection(
	ctx context.Context,
	sp *model.Spec,
	sc *config.Spec,
	cc *config.Collection,
	tags map[string]*model.Tag,
	eps map[string]*model.Endpoint,
) (*model.Collection, error) {
	coll := &model.Collection{
		ID:             id.Collection(sp.ID, cc.Location),
		SpecID:         sp.ID,
		LLMTitle:       cc.LLMTitle,
		LLMInstruction: cc.LLMInstruction,
		BaseURL:        cc.BaseURL,
		BaseMockURL:    cc.BaseMockURL,
		HTTPClient: mergeHTTPClientConfig(
			sc.HTTPClient,
			cc.HTTPClient,
		),
	}

	doc, err := s.parseSpecDocument(ctx, cc.Location)
	if err != nil {
		return nil, err
	}

	applySpecMetadata(coll, doc)

	for _, pi := range doc.PathItems {
		op := pi.Operation
		if op == nil {
			continue
		}

		tagName := resolveTagName(op.Tags)
		tagID := id.Tag(sp.ID, coll.ID, tagName)

		tag, exists := tags[tagID]
		if !exists {
			coll.Stats.Tags++
			tag = &model.Tag{
				ID:           tagID,
				SpecID:       sp.ID,
				CollectionID: coll.ID,
				Name:         tagName,
			}
			tags[tagID] = tag
		}

		coll.Stats.Methods++
		tag.Stats.Methods++

		ep := model.Endpoint{
			ID: id.Method(
				sp.ID,
				coll.ID,
				tagID,
				pi.Method,
				pi.Path,
				op.ID,
			),
			SpecID:       sp.ID,
			CollectionID: coll.ID,
			TagID:        tagID,
			Tag:          tagName,
			Name:         pi.Method,
			Path:         pi.Path,
			Operation:    op,
		}
		eps[ep.ID] = &ep
	}

	return coll, nil
}

func (s *Service) parseSpecDocument(ctx context.Context, location string) (*spec.Doc, error) {
	lp := location

	if s.cache != nil {
		var err error
		lp, err = s.cache.Resolve(ctx, location)
		if err != nil {
			return nil, fmt.Errorf("collection %q: %w", location, err)
		}
	}

	data, err := os.ReadFile(lp)
	if err != nil {
		return nil, fmt.Errorf("collection %q: read file: %w", location, err)
	}

	doc, err := spec.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("collection %q: parse spec: %w", location, err)
	}

	return doc, nil
}

func (s *Service) buildSpecInfo(sc *config.Spec, mockEnabled bool, ma *config.MockAuthConfig) (*model.Spec, error) {
	sp := &model.Spec{
		ID:             id.Domain(sc.Domain),
		Domain:         sc.Domain,
		LLMTitle:       sc.LLMTitle,
		LLMInstruction: sc.LLMInstruction,
		BaseURL:        sc.BaseURL,
		Auth:           sc.Auth.Client,
	}

	if sc.HTTPClient != nil {
		sp.HTTPClient = &model.HTTPClientConfig{
			Headers: sc.HTTPClient.Headers,
			Cookies: convertCookies(sc.HTTPClient.Cookies),
		}
	}

	if sc.Auth.Client != nil {
		if err := sp.InitAuthenticator(); err != nil {
			return nil, fmt.Errorf(
				"spec %s, failed to initialize authenticator: %w",
				sc.Domain, err,
			)
		}

		if scriptClient, isScript := sc.Auth.Client.(*auth.ScriptAuthClient); isScript {
			scriptClient.SetWorkspaceDir(s.ws.Root())
		}

		if mockEnabled {
			applyMockAuthURLs(sc.Auth.Client, ma)
		}
	}

	return sp, nil
}

const (
	defaultUserAgent      = "swag2mcp-global/1.0"
	defaultMockOAuth2Port = 9090
	defaultMockDigestPort = 9091
)

func applyMockAuthURLs(client auth.Authenticator, mockAuth *config.MockAuthConfig) {
	oauth2Port := defaultMockOAuth2Port
	digestPort := defaultMockDigestPort
	if mockAuth != nil {
		if mockAuth.OAuth2Port > 0 {
			oauth2Port = mockAuth.OAuth2Port
		}
		if mockAuth.DigestPort > 0 {
			digestPort = mockAuth.DigestPort
		}
	}
	if setter, ok := client.(auth.TokenURLSetter); ok {
		setter.SetTokenURL(fmt.Sprintf("http://127.0.0.1:%d/token", oauth2Port))
	}
	if setter, ok := client.(auth.MockBaseURLSetter); ok {
		setter.SetMockBaseURL(fmt.Sprintf("http://127.0.0.1:%d/", digestPort))
	}
}

func (s *Service) indexSpec(
	sp *model.Spec,
	colls map[string]*model.Collection,
	tags map[string]*model.Tag,
	eps map[string]*model.Endpoint,
) error {
	collections := make([]*model.Collection, 0, len(colls))
	for _, c := range colls {
		collections = append(collections, c)
	}

	ts := make([]*model.Tag, 0, len(tags))
	for _, t := range tags {
		ts = append(ts, t)
	}

	endpoints := make([]*model.Endpoint, 0, len(eps))
	for _, e := range eps {
		endpoints = append(endpoints, e)
	}

	if err := s.index.EnsureIndex(sp, collections, ts, endpoints); err != nil {
		return fmt.Errorf("failed to ensure index: %w", err)
	}

	return nil
}

func resolveTagName(tags []string) string {
	if len(tags) > 0 {
		return strings.Join(tags, ",")
	}
	return "default"
}

func applySpecMetadata(collection *model.Collection, specDocument *spec.Doc) {
	if len(collection.LLMTitle) == 0 && len(specDocument.Title) > 0 {
		collection.LLMTitle = specDocument.Title
	}
	if len(collection.LLMInstruction) == 0 && len(specDocument.Description) > 0 {
		collection.LLMInstruction = specDocument.Description
	}
	collection.Title = specDocument.Title
}

func convertCookies(cookies []config.Cookie) []httpclient.Cookie {
	if len(cookies) == 0 {
		return nil
	}

	result := make([]httpclient.Cookie, len(cookies))
	for index, cookie := range cookies {
		result[index] = httpclient.Cookie{
			Name:     cookie.Name,
			Value:    cookie.Value,
			Domain:   cookie.Domain,
			Path:     cookie.Path,
			Secure:   cookie.Secure,
			HTTPOnly: cookie.HTTPOnly,
		}
	}

	return result
}

// mergeHTTPClientConfig merges per-request HTTP configs with cascade:
// spec → collection. Collection overrides spec.
func mergeHTTPClientConfig(
	spec, collection *config.HTTPClientConfig,
) *model.HTTPClientConfig {
	result := &model.HTTPClientConfig{}

	levels := []*config.HTTPClientConfig{spec, collection}

	for _, level := range levels {
		if level == nil {
			continue
		}
		mergeHeaders(result, level)
		mergeCookies(result, level)
	}

	return result
}

func mergeHeaders(result *model.HTTPClientConfig, level *config.HTTPClientConfig) {
	if result.Headers == nil && len(level.Headers) > 0 {
		result.Headers = make(map[string]string, len(level.Headers))
		maps.Copy(result.Headers, level.Headers)
	}
}

func mergeCookies(result *model.HTTPClientConfig, level *config.HTTPClientConfig) {
	if len(result.Cookies) > 0 || len(level.Cookies) == 0 {
		return
	}

	result.Cookies = make([]httpclient.Cookie, len(level.Cookies))
	for index, cookie := range level.Cookies {
		result.Cookies[index] = httpclient.Cookie{
			Name:     cookie.Name,
			Value:    cookie.Value,
			Domain:   cookie.Domain,
			Path:     cookie.Path,
			Secure:   cookie.Secure,
			HTTPOnly: cookie.HTTPOnly,
		}
	}
}
