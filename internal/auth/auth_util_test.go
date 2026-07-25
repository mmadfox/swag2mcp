package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mmadfox/swag2mcp/internal/httpclient"
)

func TestSetAuthHeader_EmptyValue(t *testing.T) {
	t.Parallel()

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com", nil)
	setAuthHeader(req, nil, "Authorization", "")
	if req.Header.Get("Authorization") != "" {
		t.Error("header should not be set for empty value")
	}
}

func TestSetAuthHeader_WithInfo(t *testing.T) {
	t.Parallel()

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com", nil)
	var info Info
	setAuthHeader(req, &info, "Authorization", "Bearer token")
	if req.Header.Get("Authorization") != "Bearer token" {
		t.Errorf("Authorization = %q, want %q", req.Header.Get("Authorization"), "Bearer token")
	}
	if info.Headers["Authorization"] != "Bearer token" {
		t.Errorf("info.Headers[Authorization] = %q, want %q", info.Headers["Authorization"], "Bearer token")
	}
}

func TestSetAuthHeader_NilInfo(t *testing.T) {
	t.Parallel()

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com", nil)
	setAuthHeader(req, nil, "X-Key", "value")
	if req.Header.Get("X-Key") != "value" {
		t.Errorf("X-Key = %q, want %q", req.Header.Get("X-Key"), "value")
	}
}

func TestSetAuthQuery_EmptyValue(t *testing.T) {
	t.Parallel()

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	setAuthQuery(req, nil, "key", "")
	if req.URL.Query().Get("key") != "" {
		t.Error("query param should not be set for empty value")
	}
}

func TestSetAuthQuery_WithInfo(t *testing.T) {
	t.Parallel()

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	var info Info
	setAuthQuery(req, &info, "api_key", "secret")
	if req.URL.Query().Get("api_key") != "secret" {
		t.Errorf("api_key = %q, want %q", req.URL.Query().Get("api_key"), "secret")
	}
	if info.QueryParams["api_key"] != "secret" {
		t.Errorf("info.QueryParams[api_key] = %q, want %q", info.QueryParams["api_key"], "secret")
	}
}

func TestSetAuthQuery_NilInfo(t *testing.T) {
	t.Parallel()

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	setAuthQuery(req, nil, "key", "val")
	if req.URL.Query().Get("key") != "val" {
		t.Errorf("key = %q, want %q", req.URL.Query().Get("key"), "val")
	}
}

func TestTransport_RoundTrip(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client := &BearerTokenAuthClient{Token: "test-token"}
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}

	transport := &Transport{
		Base: http.DefaultTransport,
		Auth: client,
	}

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, srv.URL, nil)
	resp, err := transport.RoundTrip(req)
	if err != nil {
		t.Fatalf("RoundTrip() = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	if req.Header.Get("Authorization") != "Bearer test-token" {
		t.Errorf("Authorization = %q, want %q", req.Header.Get("Authorization"), "Bearer test-token")
	}
}

func TestTransport_RoundTrip_Error(t *testing.T) {
	t.Parallel()

	client := &BearerTokenAuthClient{Token: "test-token"}
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}

	transport := &Transport{
		Base: http.DefaultTransport,
		Auth: client,
	}

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://nonexistent.example.com", nil)
	_, err := transport.RoundTrip(req)
	if err == nil {
		t.Fatal("expected error for nonexistent host")
	}
}

func TestNewHTTPClient(t *testing.T) {
	t.Parallel()

	client := &BearerTokenAuthClient{Token: "test-token"}
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}

	httpClient := NewHTTPClient(client)
	if httpClient == nil {
		t.Fatal("NewHTTPClient() returned nil")
	}

	transport, ok := httpClient.Transport.(*Transport)
	if !ok {
		t.Fatalf("Transport type = %T, want *Transport", httpClient.Transport)
	}
	if transport.Auth != client {
		t.Error("Auth mismatch")
	}
}

func TestNoAuthClient_New(t *testing.T) {
	t.Parallel()

	client := NewNoAuthClient()
	if err := client.New(); err != nil {
		t.Errorf("New() = %v", err)
	}
}

func TestNoAuthClient_Type(t *testing.T) {
	t.Parallel()

	client := NewNoAuthClient()
	if client.Type() != NoAuth {
		t.Errorf("Type() = %q, want %q", client.Type(), NoAuth)
	}
}

func TestNoAuthClient_Validate(t *testing.T) {
	t.Parallel()

	client := NewNoAuthClient()
	if err := client.Validate(); err != nil {
		t.Errorf("Validate() = %v", err)
	}
}

func TestBearerTokenAuthClient_Validate(t *testing.T) {
	t.Parallel()

	client := &BearerTokenAuthClient{Token: "valid-token"}
	if err := client.Validate(); err != nil {
		t.Errorf("Validate() = %v", err)
	}
}

func TestBearerTokenAuthClient_Validate_EmptyToken(t *testing.T) {
	t.Parallel()

	client := &BearerTokenAuthClient{Token: ""}
	if err := client.Validate(); err == nil {
		t.Error("expected validation error for empty token")
	}
}

func TestBasicAuthClient_Validate(t *testing.T) {
	t.Parallel()

	client := &BasicAuthClient{Username: "u", Password: "p"}
	if err := client.Validate(); err != nil {
		t.Errorf("Validate() = %v", err)
	}
}

func TestBasicAuthClient_Validate_Empty(t *testing.T) {
	t.Parallel()

	client := &BasicAuthClient{}
	if err := client.Validate(); err == nil {
		t.Error("expected validation error for empty basic auth")
	}
}

func TestDigestAuthClient_New(t *testing.T) {
	t.Parallel()

	client := &DigestAuthClient{Username: "u", Password: "p"}
	if err := client.New(); err != nil {
		t.Errorf("New() = %v", err)
	}
}

func TestDigestAuthClient_Validate(t *testing.T) {
	t.Parallel()

	client := &DigestAuthClient{Username: "u", Password: "p"}
	if err := client.Validate(); err != nil {
		t.Errorf("Validate() = %v", err)
	}
}

func TestDigestAuthClient_Validate_Empty(t *testing.T) {
	t.Parallel()

	client := &DigestAuthClient{}
	if err := client.Validate(); err == nil {
		t.Error("expected validation error for empty digest auth")
	}
}

func TestDigestAuthClient_Type(t *testing.T) {
	t.Parallel()

	client := &DigestAuthClient{}
	if client.Type() != DigestAuth {
		t.Errorf("Type() = %q, want %q", client.Type(), DigestAuth)
	}
}

func TestParseDigestChallenge(t *testing.T) {
	t.Parallel()

	header := `Digest realm="test-realm", nonce="test-nonce", opaque="test-opaque", qop="auth", algorithm="MD5"`
	ch := parseDigestChallenge(header)
	if ch.realm != "test-realm" {
		t.Errorf("realm = %q, want %q", ch.realm, "test-realm")
	}
	if ch.nonce != "test-nonce" {
		t.Errorf("nonce = %q, want %q", ch.nonce, "test-nonce")
	}
	if ch.opaque != "test-opaque" {
		t.Errorf("opaque = %q, want %q", ch.opaque, "test-opaque")
	}
	if ch.qop != "auth" {
		t.Errorf("qop = %q, want %q", ch.qop, "auth")
	}
	if ch.algorithm != "MD5" {
		t.Errorf("algorithm = %q, want %q", ch.algorithm, "MD5")
	}
}

func TestParseDigestChallenge_Defaults(t *testing.T) {
	t.Parallel()

	ch := parseDigestChallenge(`Digest realm="r", nonce="n"`)
	if ch.algorithm != "MD5" {
		t.Errorf("algorithm = %q, want %q", ch.algorithm, "MD5")
	}
}

func TestParseDigestChallenge_Empty(t *testing.T) {
	t.Parallel()

	ch := parseDigestChallenge(`Digest `)
	if ch.realm != "" {
		t.Errorf("realm = %q, want empty", ch.realm)
	}
}

func TestMd5hex(t *testing.T) {
	t.Parallel()

	result := md5hex("hello")
	if len(result) != 32 {
		t.Errorf("len = %d, want 32", len(result))
	}
}

func TestGenerateCnonce(t *testing.T) {
	t.Parallel()

	client := &DigestAuthClient{}
	c1 := client.generateCnonce()
	c2 := client.generateCnonce()
	if len(c1) != 16 {
		t.Errorf("len = %d, want 16", len(c1))
	}
	if c1 == c2 {
		t.Error("cnonces should be unique")
	}
}

func TestBuildDigest(t *testing.T) {
	t.Parallel()

	client := &DigestAuthClient{Username: "user", Password: "pass"}
	ch := digestChallenge{
		realm:     "test-realm",
		nonce:     "test-nonce",
		opaque:    "test-opaque",
		qop:       "auth",
		algorithm: "MD5",
	}

	auth := client.buildDigest("GET", "/api", ch, 1, "cnonce123")
	if !strings.HasPrefix(auth, "Digest ") {
		t.Errorf("expected Digest prefix, got %q", auth)
	}
	if !strings.Contains(auth, `username="user"`) {
		t.Errorf("missing username: %s", auth)
	}
	if !strings.Contains(auth, `realm="test-realm"`) {
		t.Errorf("missing realm: %s", auth)
	}
	if !strings.Contains(auth, `nonce="test-nonce"`) {
		t.Errorf("missing nonce: %s", auth)
	}
	if !strings.Contains(auth, `response="`) {
		t.Errorf("missing response: %s", auth)
	}
	if !strings.Contains(auth, `opaque="test-opaque"`) {
		t.Errorf("missing opaque: %s", auth)
	}
	if !strings.Contains(auth, `qop=auth`) {
		t.Errorf("missing qop: %s", auth)
	}
	if !strings.Contains(auth, `nc=00000001`) {
		t.Errorf("missing nc: %s", auth)
	}
	if !strings.Contains(auth, `cnonce="cnonce123"`) {
		t.Errorf("missing cnonce: %s", auth)
	}
}

func TestBuildDigest_NoQop(t *testing.T) {
	t.Parallel()

	client := &DigestAuthClient{Username: "user", Password: "pass"}
	ch := digestChallenge{
		realm: "test-realm",
		nonce: "test-nonce",
	}

	auth := client.buildDigest("GET", "/api", ch, 1, "cnonce123")
	if !strings.Contains(auth, `response="`) {
		t.Errorf("missing response: %s", auth)
	}
	if strings.Contains(auth, "qop=") {
		t.Errorf("unexpected qop: %s", auth)
	}
}

func TestAPIKeyAuthClient_New(t *testing.T) {
	t.Parallel()

	client := &APIKeyAuthClient{Key: "X-Key", Value: "val", In: "header"}
	if err := client.New(); err != nil {
		t.Errorf("New() = %v", err)
	}
}

func TestAPIKeyAuthClient_Validate(t *testing.T) {
	t.Parallel()

	client := &APIKeyAuthClient{Key: "X-Key", Value: "val", In: "header"}
	if err := client.Validate(); err != nil {
		t.Errorf("Validate() = %v", err)
	}
}

func TestAPIKeyAuthClient_Validate_Empty(t *testing.T) {
	t.Parallel()

	client := &APIKeyAuthClient{}
	if err := client.Validate(); err == nil {
		t.Error("expected validation error for empty api key")
	}
}

func TestAPIKeyAuthClient_Type(t *testing.T) {
	t.Parallel()

	client := &APIKeyAuthClient{}
	if client.Type() != APIKeyAuth {
		t.Errorf("Type() = %q, want %q", client.Type(), APIKeyAuth)
	}
}

func TestScriptAuthClient_New_PathSeparators(t *testing.T) {
	t.Parallel()

	client := &ScriptAuthClient{Domain: "my/domain"}
	err := client.New()
	if err == nil {
		t.Fatal("expected error for domain with path separator")
	}
}

func TestScriptAuthClient_New_Backslash(t *testing.T) {
	t.Parallel()

	client := &ScriptAuthClient{Domain: "my\\domain"}
	err := client.New()
	if err == nil {
		t.Fatal("expected error for domain with backslash")
	}
}

func TestScriptAuthClient_Type(t *testing.T) {
	t.Parallel()

	client := &ScriptAuthClient{}
	if client.Type() != ScriptAuth {
		t.Errorf("Type() = %q, want %q", client.Type(), ScriptAuth)
	}
}

func TestScriptAuthClient_SetWorkspaceDir(t *testing.T) {
	t.Parallel()

	client := &ScriptAuthClient{Domain: "test"}
	client.SetWorkspaceDir("/custom/workspace")
	if client.workspaceDir != "/custom/workspace" {
		t.Errorf("workspaceDir = %q, want %q", client.workspaceDir, "/custom/workspace")
	}
}

func TestScriptAuthClient_ScriptPath(t *testing.T) {
	t.Parallel()

	client := &ScriptAuthClient{Domain: "test", workspaceDir: "/ws"}
	path := client.scriptPath()
	expected := filepath.Join("/ws", "auth_scripts", "test.sh")
	if path != expected {
		t.Errorf("scriptPath() = %q, want %q", path, expected)
	}
}

func TestScriptAuthClient_New_EnvVars(t *testing.T) {
	t.Setenv("TEST_SCRIPT_DOMAIN", "env-domain")

	client := &ScriptAuthClient{Domain: "$(TEST_SCRIPT_DOMAIN)"}
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}
	if client.Domain != "env-domain" {
		t.Errorf("Domain = %q, want %q", client.Domain, "env-domain")
	}
}

func TestScriptAuthClient_New_TrimsSpace(t *testing.T) {
	client := &ScriptAuthClient{Domain: "  my-domain  "}
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}
	if client.Domain != "my-domain" {
		t.Errorf("Domain = %q, want %q", client.Domain, "my-domain")
	}
}

func TestOAuth2ClientCredentialsAuthClient_Validate(t *testing.T) {
	t.Parallel()

	client := &OAuth2ClientCredentialsAuthClient{
		ClientID: "cid", ClientSecret: "cs", TokenURL: "https://example.com/token",
	}
	if err := client.Validate(); err != nil {
		t.Errorf("Validate() = %v", err)
	}
}

func TestOAuth2ClientCredentialsAuthClient_Validate_Empty(t *testing.T) {
	t.Parallel()

	client := &OAuth2ClientCredentialsAuthClient{}
	if err := client.Validate(); err == nil {
		t.Error("expected validation error for empty oauth2-cc")
	}
}

func TestOAuth2ClientCredentialsAuthClient_Type(t *testing.T) {
	t.Parallel()

	client := &OAuth2ClientCredentialsAuthClient{}
	if client.Type() != OAuth2ClientCredentials {
		t.Errorf("Type() = %q, want %q", client.Type(), OAuth2ClientCredentials)
	}
}

func TestOAuth2ClientCredentialsAuthClient_New_EnvVars(t *testing.T) {
	t.Setenv("TEST_CC_ID", "env-cid")
	t.Setenv("TEST_CC_SECRET", "env-cs")

	client := &OAuth2ClientCredentialsAuthClient{
		ClientID: "$(TEST_CC_ID)", ClientSecret: "$(TEST_CC_SECRET)",
		TokenURL: "https://example.com/token",
	}
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}
	if client.ClientID != "env-cid" {
		t.Errorf("ClientID = %q, want %q", client.ClientID, "env-cid")
	}
	if client.ClientSecret != "env-cs" {
		t.Errorf("ClientSecret = %q, want %q", client.ClientSecret, "env-cs")
	}
}

func TestOAuth2PasswordAuthClient_Validate(t *testing.T) {
	t.Parallel()

	client := &OAuth2PasswordAuthClient{
		Username: "u", Password: "p", ClientID: "cid",
		TokenURL: "https://example.com/token",
	}
	if err := client.Validate(); err != nil {
		t.Errorf("Validate() = %v", err)
	}
}

func TestOAuth2PasswordAuthClient_Validate_Empty(t *testing.T) {
	t.Parallel()

	client := &OAuth2PasswordAuthClient{}
	if err := client.Validate(); err == nil {
		t.Error("expected validation error for empty oauth2-pwd")
	}
}

func TestOAuth2PasswordAuthClient_Type(t *testing.T) {
	t.Parallel()

	client := &OAuth2PasswordAuthClient{}
	if client.Type() != OAuth2Password {
		t.Errorf("Type() = %q, want %q", client.Type(), OAuth2Password)
	}
}

func TestOAuth2PasswordAuthClient_New_EnvVars(t *testing.T) {
	t.Setenv("TEST_PWD_USER", "env-user")
	t.Setenv("TEST_PWD_PASS", "env-pass")
	t.Setenv("TEST_PWD_CID", "env-cid")
	t.Setenv("TEST_PWD_CS", "env-cs")

	client := &OAuth2PasswordAuthClient{
		Username: "$(TEST_PWD_USER)", Password: "$(TEST_PWD_PASS)",
		ClientID: "$(TEST_PWD_CID)", ClientSecret: "$(TEST_PWD_CS)",
		TokenURL: "https://example.com/token",
	}
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}
	if client.Username != "env-user" {
		t.Errorf("Username = %q, want %q", client.Username, "env-user")
	}
	if client.Password != "env-pass" {
		t.Errorf("Password = %q, want %q", client.Password, "env-pass")
	}
}

func TestDefaultHTTPClient(t *testing.T) {
	t.Parallel()

	cli, err := httpclient.NewDefault()
	if err != nil {
		t.Fatalf("NewDefault() = %v", err)
	}
	if cli == nil {
		t.Fatal("NewDefault() returned nil")
	}
	if cli.Timeout == 0 {
		t.Error("Timeout should be set")
	}
}

func TestInfo_HeadersNil(t *testing.T) {
	t.Parallel()

	var info Info
	if info.Headers != nil {
		t.Error("Headers should be nil initially")
	}
	if info.QueryParams != nil {
		t.Error("QueryParams should be nil initially")
	}
}

func TestDigestAuthClient_fetchChallenge_Non401(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()

	client := &DigestAuthClient{Username: "u", Password: "p"}
	req, _ := http.NewRequest(http.MethodGet, srv.URL, nil)
	_, err := client.fetchChallenge(req)
	if err == nil {
		t.Fatal("expected error for non-401 response")
	}
}

func TestDigestAuthClient_fetchChallenge_MissingWWWAuth(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	client := &DigestAuthClient{Username: "u", Password: "p"}
	req, _ := http.NewRequest(http.MethodGet, srv.URL, nil)
	_, err := client.fetchChallenge(req)
	if err == nil {
		t.Fatal("expected error for missing WWW-Authenticate")
	}
}

func TestDigestAuthClient_fetchChallenge_NonDigest(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("WWW-Authenticate", "Basic realm=test")
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	client := &DigestAuthClient{Username: "u", Password: "p"}
	req, _ := http.NewRequest(http.MethodGet, srv.URL, nil)
	_, err := client.fetchChallenge(req)
	if err == nil {
		t.Fatal("expected error for non-Digest challenge")
	}
}

func TestDigestAuthClient_Apply_EnvVars(t *testing.T) {
	t.Setenv("TEST_DIGEST_USER", "env-user")
	t.Setenv("TEST_DIGEST_PASS", "env-pass")

	ds := newDigestTestServer(t, "env-user", "env-pass")

	client := &DigestAuthClient{
		Username: "$(TEST_DIGEST_USER)",
		Password: "$(TEST_DIGEST_PASS)",
	}
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}

	req, _ := http.NewRequest(http.MethodGet, ds.srv.URL+"/api", nil)
	if err := client.Apply(req, nil); err != nil {
		t.Fatalf("Apply() = %v", err)
	}

	auth := req.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Digest ") {
		t.Errorf("Authorization = %q, want Digest prefix", auth)
	}
}

func TestScriptAuthClient_Apply_RefetchesAfterExpiry(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeScript(t, dir, `echo '{"token": "script-token", "expires_in": 1}'`)

	client := &ScriptAuthClient{
		Domain:       "testdomain",
		workspaceDir: dir,
	}
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}

	req1, _ := newGetRequest()
	if err := client.Apply(req1, nil); err != nil {
		t.Fatalf("Apply #1 = %v", err)
	}

	client.mu.Lock()
	client.expiresAt = time.Now().Add(-time.Second)
	client.mu.Unlock()

	req2, _ := newGetRequest()
	if err := client.Apply(req2, nil); err != nil {
		t.Fatalf("Apply #2 = %v", err)
	}

	if v := req2.Header.Get("Authorization"); v != "Bearer script-token" {
		t.Errorf("Authorization = %q, want %q", v, "Bearer script-token")
	}
}

func TestScriptAuthClient_Apply_NoWorkspaceDir(t *testing.T) {
	t.Parallel()

	client := &ScriptAuthClient{Domain: "testdomain"}
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}

	req, _ := newGetRequest()
	err := client.Apply(req, nil)
	if err == nil {
		t.Fatal("expected error for missing workspace dir")
	}
}

func TestScriptAuthClient_Apply_EnvVars(t *testing.T) {
	t.Setenv("TEST_SCRIPT_DOMAIN", "testdomain")

	dir := t.TempDir()
	writeScript(t, dir, `echo '{"token": "env-script-token", "expires_in": 3600}'`)

	client := &ScriptAuthClient{
		Domain:       "$(TEST_SCRIPT_DOMAIN)",
		workspaceDir: dir,
	}
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}

	req, _ := newGetRequest()
	if err := client.Apply(req, nil); err != nil {
		t.Fatalf("Apply() = %v", err)
	}

	if v := req.Header.Get("Authorization"); v != "Bearer env-script-token" {
		t.Errorf("Authorization = %q, want %q", v, "Bearer env-script-token")
	}
}

func TestOAuth2ClientCredentialsAuthClient_Apply_Scopes(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		if r.Form.Get("scope") != "read write" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		resp := oauth2TokenResponse{AccessToken: "scoped-token", TokenType: "Bearer", ExpiresIn: 3600}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	t.Cleanup(srv.Close)

	client := &OAuth2ClientCredentialsAuthClient{
		ClientID: "c", ClientSecret: "s", TokenURL: srv.URL + "/token",
		Scopes: []string{"read", "write"},
	}
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	if err := client.Apply(req, nil); err != nil {
		t.Fatalf("Apply() = %v", err)
	}
	if v := req.Header.Get("Authorization"); v != "Bearer scoped-token" {
		t.Errorf("Authorization = %q, want %q", v, "Bearer scoped-token")
	}
}

func TestOAuth2ClientCredentialsAuthClient_Apply_DefaultExpiry(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		resp := oauth2TokenResponse{AccessToken: "token", TokenType: "Bearer", ExpiresIn: 0}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	t.Cleanup(srv.Close)

	client := &OAuth2ClientCredentialsAuthClient{
		ClientID: "c", ClientSecret: "s", TokenURL: srv.URL + "/token",
	}
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	if err := client.Apply(req, nil); err != nil {
		t.Fatalf("Apply() = %v", err)
	}
	if v := req.Header.Get("Authorization"); v != "Bearer token" {
		t.Errorf("Authorization = %q, want %q", v, "Bearer token")
	}
}

func TestOAuth2PasswordAuthClient_Apply_DefaultExpiry(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		resp := oauth2TokenResponse{AccessToken: "token", TokenType: "Bearer", ExpiresIn: 0}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	t.Cleanup(srv.Close)

	client := &OAuth2PasswordAuthClient{
		Username: "u", Password: "p", ClientID: "c",
		TokenURL: srv.URL + "/token",
	}
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	if err := client.Apply(req, nil); err != nil {
		t.Fatalf("Apply() = %v", err)
	}
	if v := req.Header.Get("Authorization"); v != "Bearer token" {
		t.Errorf("Authorization = %q, want %q", v, "Bearer token")
	}
}

func TestOAuth2PasswordAuthClient_Apply_Scopes(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		if r.Form.Get("scope") != "openid profile" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		resp := oauth2TokenResponse{AccessToken: "scoped-token", TokenType: "Bearer", ExpiresIn: 3600}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	t.Cleanup(srv.Close)

	client := &OAuth2PasswordAuthClient{
		Username: "u", Password: "p", ClientID: "c",
		TokenURL: srv.URL + "/token",
		Scopes:   []string{"openid", "profile"},
	}
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	if err := client.Apply(req, nil); err != nil {
		t.Fatalf("Apply() = %v", err)
	}
	if v := req.Header.Get("Authorization"); v != "Bearer scoped-token" {
		t.Errorf("Authorization = %q, want %q", v, "Bearer scoped-token")
	}
}

func TestOAuth2PasswordAuthClient_Apply_EmptyAccessToken(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		resp := oauth2TokenResponse{AccessToken: "", TokenType: "Bearer", ExpiresIn: 3600}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	t.Cleanup(srv.Close)

	client := &OAuth2PasswordAuthClient{
		Username: "u", Password: "p", ClientID: "c",
		TokenURL: srv.URL + "/token",
	}
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	if err := client.Apply(req, nil); err == nil {
		t.Fatal("expected error for empty access_token")
	}
}

func TestOAuth2PasswordAuthClient_Apply_Non200(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	t.Cleanup(srv.Close)

	client := &OAuth2PasswordAuthClient{
		Username: "u", Password: "p", ClientID: "c",
		TokenURL: srv.URL + "/token",
	}
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	if err := client.Apply(req, nil); err == nil {
		t.Fatal("expected error for non-200")
	}
}

func TestOAuth2ClientCredentialsAuthClient_Apply_EmptyAccessToken(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		resp := oauth2TokenResponse{AccessToken: "", TokenType: "Bearer", ExpiresIn: 3600}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	t.Cleanup(srv.Close)

	client := &OAuth2ClientCredentialsAuthClient{
		ClientID: "c", ClientSecret: "s", TokenURL: srv.URL + "/token",
	}
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	if err := client.Apply(req, nil); err == nil {
		t.Fatal("expected error for empty access_token")
	}
}

func TestOAuth2ClientCredentialsAuthClient_Apply_Non200(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	t.Cleanup(srv.Close)

	client := &OAuth2ClientCredentialsAuthClient{
		ClientID: "c", ClientSecret: "s", TokenURL: srv.URL + "/token",
	}
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	if err := client.Apply(req, nil); err == nil {
		t.Fatal("expected error for non-200")
	}
}
