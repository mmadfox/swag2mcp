package auth

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mmadfox/swag2mcp/internal/httpclient"
)

func TestSetAuthHeader_EmptyValue(t *testing.T) {
	t.Parallel()

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com", nil)
	setAuthHeader(req, nil, headerAuthorization, "")
	assert.Empty(t, req.Header.Get(headerAuthorization), "header should not be set for empty value")
}

func TestSetAuthHeader_WithInfo(t *testing.T) {
	t.Parallel()

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com", nil)
	var info Info
	setAuthHeader(req, &info, headerAuthorization, "Bearer token")
	assert.Equal(t, "Bearer token", req.Header.Get(headerAuthorization))
	assert.Equal(t, "Bearer token", info.Headers[headerAuthorization])
}

func TestSetAuthHeader_NilInfo(t *testing.T) {
	t.Parallel()

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com", nil)
	setAuthHeader(req, nil, "X-Key", "value")
	assert.Equal(t, "value", req.Header.Get("X-Key"))
}

func TestSetAuthQuery_EmptyValue(t *testing.T) {
	t.Parallel()

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	setAuthQuery(req, nil, "key", "")
	assert.Empty(t, req.URL.Query().Get("key"), "query param should not be set for empty value")
}

func TestSetAuthQuery_WithInfo(t *testing.T) {
	t.Parallel()

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	var info Info
	setAuthQuery(req, &info, "api_key", "secret")
	assert.Equal(t, "secret", req.URL.Query().Get("api_key"))
	assert.Equal(t, "secret", info.QueryParams["api_key"])
}

func TestSetAuthQuery_NilInfo(t *testing.T) {
	t.Parallel()

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	setAuthQuery(req, nil, "key", "val")
	assert.Equal(t, "val", req.URL.Query().Get("key"))
}

func TestTransport_RoundTrip(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client := &BearerTokenAuthClient{Token: "test-token"}
	require.NoError(t, client.New(), "New()")

	transport := &Transport{
		Base: http.DefaultTransport,
		Auth: client,
	}

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, srv.URL, nil)
	resp, err := transport.RoundTrip(req)
	require.NoError(t, err, "RoundTrip()")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Bearer test-token", req.Header.Get(headerAuthorization))
}

func TestTransport_RoundTrip_Error(t *testing.T) {
	t.Parallel()

	client := &BearerTokenAuthClient{Token: "test-token"}
	require.NoError(t, client.New(), "New()")

	transport := &Transport{
		Base: http.DefaultTransport,
		Auth: client,
	}

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://nonexistent.example.com", nil)
	_, err := transport.RoundTrip(req)
	require.Error(t, err, "expected error for nonexistent host")
}

func TestNewHTTPClient(t *testing.T) {
	t.Parallel()

	client := &BearerTokenAuthClient{Token: "test-token"}
	require.NoError(t, client.New(), "New()")

	httpClient := newHTTPClient(client)
	require.NotNil(t, httpClient, "newHTTPClient() returned nil")

	transport, ok := httpClient.Transport.(*Transport)
	require.True(t, ok, "Transport type should be *Transport")
	assert.Equal(t, client, transport.Auth)
}

func TestNoAuthClient_New(t *testing.T) {
	t.Parallel()

	client := NewNoAuthClient()
	require.NoError(t, client.New())
}

func TestNoAuthClient_Type(t *testing.T) {
	t.Parallel()

	client := NewNoAuthClient()
	assert.Equal(t, NoAuth, client.Type())
}

func TestNoAuthClient_Validate(t *testing.T) {
	t.Parallel()

	client := NewNoAuthClient()
	require.NoError(t, client.Validate())
}

func TestBearerTokenAuthClient_Validate(t *testing.T) {
	t.Parallel()

	client := &BearerTokenAuthClient{Token: "valid-token"}
	require.NoError(t, client.Validate())
}

func TestBearerTokenAuthClient_Validate_EmptyToken(t *testing.T) {
	t.Parallel()

	client := &BearerTokenAuthClient{Token: ""}
	require.Error(t, client.Validate())
}

func TestBasicAuthClient_Type(t *testing.T) {
	t.Parallel()

	client := &BasicAuthClient{}
	assert.Equal(t, BasicAuth, client.Type())
}

func TestBearerTokenAuthClient_Type(t *testing.T) {
	t.Parallel()

	client := &BearerTokenAuthClient{}
	assert.Equal(t, BearerTokenAuth, client.Type())
}

func TestDigestAuthClient_SetMockBaseURL(t *testing.T) {
	t.Parallel()

	client := &DigestAuthClient{}
	client.SetMockBaseURL("http://localhost:9091/")
	assert.Equal(t, "http://localhost:9091/", client.MockBaseURL)
}

func TestOAuth2ClientCredentialsAuthClient_SetTokenURL(t *testing.T) {
	t.Parallel()

	client := &OAuth2ClientCredentialsAuthClient{}
	client.SetTokenURL("http://localhost:9090/token")
	assert.Equal(t, "http://localhost:9090/token", client.TokenURL)
}

func TestOAuth2PasswordAuthClient_SetTokenURL(t *testing.T) {
	t.Parallel()

	client := &OAuth2PasswordAuthClient{}
	client.SetTokenURL("http://localhost:9090/token")
	assert.Equal(t, "http://localhost:9090/token", client.TokenURL)
}

func TestBasicAuthClient_Validate(t *testing.T) {
	t.Parallel()

	client := &BasicAuthClient{Username: "u", Password: "p"}
	require.NoError(t, client.Validate())
}

func TestBasicAuthClient_Validate_Empty(t *testing.T) {
	t.Parallel()

	client := &BasicAuthClient{}
	require.Error(t, client.Validate())
}

func TestDigestAuthClient_New(t *testing.T) {
	t.Parallel()

	client := &DigestAuthClient{Username: "u", Password: "p"}
	require.NoError(t, client.New())
}

func TestDigestAuthClient_Validate(t *testing.T) {
	t.Parallel()

	client := &DigestAuthClient{Username: "u", Password: "p"}
	require.NoError(t, client.Validate())
}

func TestDigestAuthClient_Validate_Empty(t *testing.T) {
	t.Parallel()

	client := &DigestAuthClient{}
	require.Error(t, client.Validate())
}

func TestDigestAuthClient_Type(t *testing.T) {
	t.Parallel()

	client := &DigestAuthClient{}
	assert.Equal(t, DigestAuth, client.Type())
}

func TestParseDigestChallenge(t *testing.T) {
	t.Parallel()

	header := `Digest realm="test-realm", nonce="test-nonce", opaque="test-opaque", qop="auth", algorithm="MD5"`
	ch := parseDigestChallenge(header)
	assert.Equal(t, "test-realm", ch.realm)
	assert.Equal(t, "test-nonce", ch.nonce)
	assert.Equal(t, "test-opaque", ch.opaque)
	assert.Equal(t, "auth", ch.qop)
	assert.Equal(t, "MD5", ch.algorithm)
}

func TestParseDigestChallenge_Defaults(t *testing.T) {
	t.Parallel()

	ch := parseDigestChallenge(`Digest realm="r", nonce="n"`)
	assert.Equal(t, "MD5", ch.algorithm)
}

func TestParseDigestChallenge_Empty(t *testing.T) {
	t.Parallel()

	ch := parseDigestChallenge(`Digest `)
	assert.Empty(t, ch.realm)
}

func TestMd5hex(t *testing.T) {
	t.Parallel()

	result := md5hex("hello")
	assert.Len(t, result, 32)
}

func TestGenerateCnonce(t *testing.T) {
	t.Parallel()

	client := &DigestAuthClient{}
	c1 := client.generateCnonce()
	c2 := client.generateCnonce()
	assert.Len(t, c1, 16)
	assert.NotEqual(t, c1, c2, "cnonces should be unique")
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
	assert.True(t, strings.HasPrefix(auth, "Digest "), "expected Digest prefix")
	assert.Contains(t, auth, `username="user"`)
	assert.Contains(t, auth, `realm="test-realm"`)
	assert.Contains(t, auth, `nonce="test-nonce"`)
	assert.Contains(t, auth, `response="`)
	assert.Contains(t, auth, `opaque="test-opaque"`)
	assert.Contains(t, auth, `qop=auth`)
	assert.Contains(t, auth, `nc=00000001`)
	assert.Contains(t, auth, `cnonce="cnonce123"`)
}

func TestBuildDigest_NoQop(t *testing.T) {
	t.Parallel()

	client := &DigestAuthClient{Username: "user", Password: "pass"}
	ch := digestChallenge{
		realm: "test-realm",
		nonce: "test-nonce",
	}

	auth := client.buildDigest("GET", "/api", ch, 1, "cnonce123")
	assert.Contains(t, auth, `response="`)
	assert.NotContains(t, auth, "qop=")
}

func TestAPIKeyAuthClient_New(t *testing.T) {
	t.Parallel()

	client := &APIKeyAuthClient{Key: "X-Key", Value: "val", In: "header"}
	require.NoError(t, client.New())
}

func TestAPIKeyAuthClient_Validate(t *testing.T) {
	t.Parallel()

	client := &APIKeyAuthClient{Key: "X-Key", Value: "val", In: "header"}
	require.NoError(t, client.Validate())
}

func TestAPIKeyAuthClient_Validate_Empty(t *testing.T) {
	t.Parallel()

	client := &APIKeyAuthClient{}
	require.Error(t, client.Validate())
}

func TestAPIKeyAuthClient_Type(t *testing.T) {
	t.Parallel()

	client := &APIKeyAuthClient{}
	assert.Equal(t, APIKeyAuth, client.Type())
}

func TestScriptAuthClient_New_PathSeparators(t *testing.T) {
	t.Parallel()

	client := &ScriptAuthClient{Domain: "my/domain"}
	err := client.New()
	require.Error(t, err, "expected error for domain with path separator")
}

func TestScriptAuthClient_New_Backslash(t *testing.T) {
	t.Parallel()

	client := &ScriptAuthClient{Domain: "my\\domain"}
	err := client.New()
	require.Error(t, err, "expected error for domain with backslash")
}

func TestScriptAuthClient_Type(t *testing.T) {
	t.Parallel()

	client := &ScriptAuthClient{}
	assert.Equal(t, ScriptAuth, client.Type())
}

func TestScriptAuthClient_SetWorkspaceDir(t *testing.T) {
	t.Parallel()

	client := &ScriptAuthClient{Domain: "test"}
	client.SetWorkspaceDir("/custom/workspace")
	assert.Equal(t, "/custom/workspace", client.workspaceDir)
}

func TestScriptAuthClient_ScriptPath(t *testing.T) {
	t.Parallel()

	client := &ScriptAuthClient{Domain: "test", workspaceDir: "/ws"}
	path := client.scriptPath()
	expected := filepath.Join("/ws", "auth_scripts", "test.sh")
	assert.Equal(t, expected, path)
}

func TestScriptAuthClient_New_EnvVars(t *testing.T) {
	t.Setenv("TEST_SCRIPT_DOMAIN", "env-domain")

	client := &ScriptAuthClient{Domain: "$(TEST_SCRIPT_DOMAIN)"}
	require.NoError(t, client.New(), "New()")
	assert.Equal(t, "env-domain", client.Domain)
}

func TestScriptAuthClient_New_TrimsSpace(t *testing.T) {
	client := &ScriptAuthClient{Domain: "  my-domain  "}
	require.NoError(t, client.New(), "New()")
	assert.Equal(t, "my-domain", client.Domain)
}

func TestOAuth2ClientCredentialsAuthClient_Validate(t *testing.T) {
	t.Parallel()

	client := &OAuth2ClientCredentialsAuthClient{
		ClientID: "cid", ClientSecret: "cs", TokenURL: "https://example.com/token",
	}
	require.NoError(t, client.Validate())
}

func TestOAuth2ClientCredentialsAuthClient_Validate_Empty(t *testing.T) {
	t.Parallel()

	client := &OAuth2ClientCredentialsAuthClient{}
	require.Error(t, client.Validate())
}

func TestOAuth2ClientCredentialsAuthClient_Type(t *testing.T) {
	t.Parallel()

	client := &OAuth2ClientCredentialsAuthClient{}
	assert.Equal(t, OAuth2ClientCredentials, client.Type())
}

func TestOAuth2ClientCredentialsAuthClient_New_EnvVars(t *testing.T) {
	t.Setenv("TEST_CC_ID", "env-cid")
	t.Setenv("TEST_CC_SECRET", "env-cs")

	client := &OAuth2ClientCredentialsAuthClient{
		ClientID: "$(TEST_CC_ID)", ClientSecret: "$(TEST_CC_SECRET)",
		TokenURL: "https://example.com/token",
	}
	require.NoError(t, client.New(), "New()")
	assert.Equal(t, "env-cid", client.ClientID)
	assert.Equal(t, "env-cs", client.ClientSecret)
}

func TestOAuth2PasswordAuthClient_Validate(t *testing.T) {
	t.Parallel()

	client := &OAuth2PasswordAuthClient{
		Username: "u", Password: "p", ClientID: "cid",
		TokenURL: "https://example.com/token",
	}
	require.NoError(t, client.Validate())
}

func TestOAuth2PasswordAuthClient_Validate_Empty(t *testing.T) {
	t.Parallel()

	client := &OAuth2PasswordAuthClient{}
	require.Error(t, client.Validate())
}

func TestOAuth2PasswordAuthClient_Type(t *testing.T) {
	t.Parallel()

	client := &OAuth2PasswordAuthClient{}
	assert.Equal(t, OAuth2Password, client.Type())
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
	require.NoError(t, client.New(), "New()")
	assert.Equal(t, "env-user", client.Username)
	assert.Equal(t, "env-pass", client.Password)
}

func TestDefaultHTTPClient(t *testing.T) {
	t.Parallel()

	cli, err := httpclient.NewDefault()
	require.NoError(t, err, "NewDefault()")
	require.NotNil(t, cli, "NewDefault() returned nil")
	assert.NotZero(t, cli.Timeout, "Timeout should be set")
}

func TestInfo_HeadersNil(t *testing.T) {
	t.Parallel()

	var info Info
	assert.Nil(t, info.Headers, "Headers should be nil initially")
	assert.Nil(t, info.QueryParams, "QueryParams should be nil initially")
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
	require.Error(t, err, "expected error for non-401 response")
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
	require.Error(t, err, "expected error for missing WWW-Authenticate")
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
	require.Error(t, err, "expected error for non-Digest challenge")
}

func TestDigestAuthClient_Apply_EnvVars(t *testing.T) {
	t.Setenv("TEST_DIGEST_USER", "env-user")
	t.Setenv("TEST_DIGEST_PASS", "env-pass")

	ds := newDigestTestServer(t, "env-user", "env-pass")

	client := &DigestAuthClient{
		Username: "$(TEST_DIGEST_USER)",
		Password: "$(TEST_DIGEST_PASS)",
	}
	require.NoError(t, client.New(), "New()")

	req, _ := http.NewRequest(http.MethodGet, ds.srv.URL+"/api", nil)
	require.NoError(t, client.Apply(req, nil), "Apply()")

	auth := req.Header.Get(headerAuthorization)
	assert.True(t, strings.HasPrefix(auth, "Digest "), "Authorization should have Digest prefix")
}

func TestScriptAuthClient_Apply_RefetchesAfterExpiry(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeScript(t, dir, `echo '{"token": "script-token", "expires_in": 1}'`)

	client := &ScriptAuthClient{
		Domain:       "testdomain",
		workspaceDir: dir,
	}
	require.NoError(t, client.New(), "New()")

	req1, _ := newGetRequest()
	require.NoError(t, client.Apply(req1, nil), "Apply #1")

	client.mu.Lock()
	client.expiresAt = time.Now().Add(-time.Second)
	client.mu.Unlock()

	req2, _ := newGetRequest()
	require.NoError(t, client.Apply(req2, nil), "Apply #2")

	assert.Equal(t, "Bearer script-token", req2.Header.Get(headerAuthorization))
}

func TestScriptAuthClient_Apply_NoWorkspaceDir(t *testing.T) {
	t.Parallel()

	client := &ScriptAuthClient{Domain: "testdomain"}
	require.NoError(t, client.New(), "New()")

	req, _ := newGetRequest()
	err := client.Apply(req, nil)
	require.Error(t, err, "expected error for missing workspace dir")
}

func TestScriptAuthClient_Apply_EnvVars(t *testing.T) {
	t.Setenv("TEST_SCRIPT_DOMAIN", "testdomain")

	dir := t.TempDir()
	writeScript(t, dir, `echo '{"token": "env-script-token", "expires_in": 3600}'`)

	client := &ScriptAuthClient{
		Domain:       "$(TEST_SCRIPT_DOMAIN)",
		workspaceDir: dir,
	}
	require.NoError(t, client.New(), "New()")

	req, _ := newGetRequest()
	require.NoError(t, client.Apply(req, nil), "Apply()")

	assert.Equal(t, "Bearer env-script-token", req.Header.Get(headerAuthorization))
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
	require.NoError(t, client.New(), "New()")

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	require.NoError(t, client.Apply(req, nil), "Apply()")
	assert.Equal(t, "Bearer scoped-token", req.Header.Get(headerAuthorization))
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
	require.NoError(t, client.New(), "New()")

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	require.NoError(t, client.Apply(req, nil), "Apply()")
	assert.Equal(t, "Bearer token", req.Header.Get(headerAuthorization))
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
	require.NoError(t, client.New(), "New()")

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	require.NoError(t, client.Apply(req, nil), "Apply()")
	assert.Equal(t, "Bearer token", req.Header.Get(headerAuthorization))
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
	require.NoError(t, client.New(), "New()")

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	require.NoError(t, client.Apply(req, nil), "Apply()")
	assert.Equal(t, "Bearer scoped-token", req.Header.Get(headerAuthorization))
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
	require.NoError(t, client.New(), "New()")

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	err := client.Apply(req, nil)
	require.Error(t, err, "expected error for empty access_token")
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
	require.NoError(t, client.New(), "New()")

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	err := client.Apply(req, nil)
	require.Error(t, err, "expected error for non-200")
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
	require.NoError(t, client.New(), "New()")

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	err := client.Apply(req, nil)
	require.Error(t, err, "expected error for empty access_token")
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
	require.NoError(t, client.New(), "New()")

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	err := client.Apply(req, nil)
	require.Error(t, err, "expected error for non-200")
}
