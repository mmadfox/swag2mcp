package mockserver

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	authNonceBytes      = 8
	authNonceTTL        = 5 * time.Minute
	authTokenLength     = 32
	authShutdownTimeout = 5
	authTokenExpiresIn  = 3600
	authExpiresInKey    = "expires_in"
	grantTypePassword   = "password"
)

type authServerType string

const (
	authServerOAuth2 authServerType = "oauth2"
	authServerDigest authServerType = "digest"
)

type authMockServer struct {
	serverType authServerType
	addr       string
	server     *http.Server
	logger     *slog.Logger
	tlsConfig  *tls.Config
	mu         sync.Mutex
	nonce      string
	opaque     string
	nonceTime  time.Time
}

func newAuthMockServer(
	serverType authServerType,
	addr string,
	tlsConfig *tls.Config,
	logger *slog.Logger,
) *authMockServer {
	if logger == nil {
		logger = slog.New(slog.DiscardHandler)
	}
	return &authMockServer{
		serverType: serverType,
		addr:       addr,
		tlsConfig:  tlsConfig,
		logger:     logger,
	}
}

func (m *authMockServer) start(ctx context.Context) {
	mux := http.NewServeMux()

	switch m.serverType {
	case authServerOAuth2:
		mux.HandleFunc("/token", m.handleOAuth2)
	case authServerDigest:
		mux.HandleFunc("/", m.handleDigest)
	default:
		m.logger.ErrorContext(ctx, "unsupported auth server type",
			"type", m.serverType,
		)
		return
	}

	address := extractHostPort(m.addr)
	if !strings.Contains(address, ":") {
		address = ":" + address
	}
	m.server = &http.Server{
		Addr:              address,
		Handler:           mux,
		ReadHeaderTimeout: authShutdownTimeout * time.Second,
	}

	if m.tlsConfig != nil {
		m.server.TLSConfig = m.tlsConfig
	}

	go func() {
		serveError := m.server.ListenAndServe()
		if serveError != nil && serveError != http.ErrServerClosed {
			m.logger.ErrorContext(ctx, "auth mock server error",
				"type", m.serverType,
				"addr", m.addr,
				"error", serveError,
			)
		}
	}()

	m.logger.InfoContext(ctx, "auth mock server started",
		"type", m.serverType,
		"addr", m.addr,
	)
}

func (m *authMockServer) shutdown() {
	if m.server != nil {
		shutdownContext, shutdownCancel := context.WithTimeout(
			context.Background(),
			authShutdownTimeout*time.Second,
		)
		defer shutdownCancel()
		_ = m.server.Shutdown(shutdownContext)
	}
}

func (m *authMockServer) handleOAuth2(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		m.logger.WarnContext(request.Context(), "oauth2 mock: method not allowed",
			"method", request.Method,
		)
		http.Error(responseWriter, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	parseError := request.ParseForm()
	if parseError != nil {
		m.logger.WarnContext(request.Context(), "oauth2 mock: failed to parse form",
			"error", parseError,
		)
		http.Error(responseWriter, "Bad Request", http.StatusBadRequest)
		return
	}

	grantType := request.FormValue("grant_type")

	m.logger.InfoContext(request.Context(), "oauth2 mock: token requested",
		"grant_type", grantType,
		"client_id", request.FormValue("client_id"),
		"username", request.FormValue("username"),
	)

	switch grantType {
	case "client_credentials":
		m.handleOAuth2CC(responseWriter, request)
	case grantTypePassword:
		m.handleOAuth2Password(responseWriter, request)
	default:
		m.logger.WarnContext(request.Context(), "oauth2 mock: invalid grant_type",
			"grant_type", grantType,
		)
		http.Error(responseWriter, "Invalid grant_type", http.StatusBadRequest)
	}
}

func (m *authMockServer) handleOAuth2CC(responseWriter http.ResponseWriter, request *http.Request) {
	clientID := request.FormValue("client_id")
	if clientID == "" {
		m.logger.WarnContext(request.Context(), "oauth2 mock: missing client_id")
		http.Error(responseWriter, "Missing client_id", http.StatusBadRequest)
		return
	}

	token := generateRandomToken()
	m.logger.InfoContext(request.Context(), "oauth2 mock: token issued",
		"grant_type", "client_credentials",
		"client_id", clientID,
		"token_prefix", token[:8],
	)

	tokenResponse := map[string]any{
		"access_token":   token,
		"token_type":     "Bearer",
		authExpiresInKey: authTokenExpiresIn,
		"scope":          request.FormValue("scope"),
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(responseWriter).Encode(tokenResponse)
}

func (m *authMockServer) handleOAuth2Password(responseWriter http.ResponseWriter, request *http.Request) {
	username := request.FormValue("username")
	password := request.FormValue("password")
	if username == "" || password == "" {
		m.logger.WarnContext(request.Context(), "oauth2 mock: missing username or password",
			"username", username,
		)
		http.Error(responseWriter, "Missing username or password", http.StatusBadRequest)
		return
	}

	token := generateRandomToken()
	m.logger.InfoContext(request.Context(), "oauth2 mock: token issued",
		"grant_type", "password",
		"username", username,
		"token_prefix", token[:8],
	)

	tokenResponse := map[string]any{
		"access_token":   token,
		"token_type":     "Bearer",
		authExpiresInKey: authTokenExpiresIn,
		"scope":          request.FormValue("scope"),
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(responseWriter).Encode(tokenResponse)
}

func (m *authMockServer) handleDigest(responseWriter http.ResponseWriter, request *http.Request) {
	authorization := request.Header.Get("Authorization")

	if !strings.HasPrefix(authorization, "Digest ") {
		m.logger.InfoContext(request.Context(), "digest mock: sending challenge")
		m.generateDigestChallenge(responseWriter)
		return
	}

	digestParams := m.parseDigestAuthorization(authorization)
	response := digestParams["response"]
	if response == "" {
		m.logger.InfoContext(request.Context(), "digest mock: empty response, sending challenge")
		m.generateDigestChallenge(responseWriter)
		return
	}

	m.logger.InfoContext(request.Context(), "digest mock: authentication successful",
		"username", digestParams["username"],
	)

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	_, _ = responseWriter.Write([]byte(`{"status":"authenticated","method":"digest"}`))
}

func (m *authMockServer) generateDigestChallenge(responseWriter http.ResponseWriter) {
	m.mu.Lock()
	now := time.Now()
	if m.nonce == "" || now.Sub(m.nonceTime) > authNonceTTL {
		nonceBytes := make([]byte, authNonceBytes)
		_, _ = rand.Read(nonceBytes)
		m.nonce = hex.EncodeToString(nonceBytes)
		m.opaque = hex.EncodeToString(nonceBytes)
		m.nonceTime = now
	}
	nonce := m.nonce
	opaque := m.opaque
	m.mu.Unlock()

	challenge := fmt.Sprintf(
		"Digest realm=\"swag2mcp-mock\", nonce=\"%s\", opaque=\"%s\", algorithm=MD5, qop=\"auth\"",
		nonce, opaque,
	)

	m.logger.Info("digest mock: challenge sent",
		"nonce_prefix", nonce[:8],
	)

	responseWriter.Header().Set("WWW-Authenticate", challenge)
	http.Error(responseWriter, "Unauthorized", http.StatusUnauthorized)
}

func (m *authMockServer) parseDigestAuthorization(authorization string) map[string]string {
	parameters := make(map[string]string)
	headerValue := strings.TrimPrefix(authorization, "Digest ")

	for part := range strings.SplitSeq(headerValue, ",") {
		part = strings.TrimSpace(part)
		key, value, found := strings.Cut(part, "=")
		if !found {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.Trim(strings.TrimSpace(value), "\"")
		parameters[key] = value
	}

	return parameters
}

func generateRandomToken() string {
	tokenBytes := make([]byte, authTokenLength)
	_, _ = rand.Read(tokenBytes)
	return hex.EncodeToString(tokenBytes)
}
