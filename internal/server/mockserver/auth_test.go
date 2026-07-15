package mockserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAuthMockServer_handleDigest_NoChallenge(t *testing.T) {
	t.Parallel()

	server := newAuthMockServer(authServerDigest, "127.0.0.1:0", nil, nil)
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)

	server.handleDigest(responseRecorder, request)

	if responseRecorder.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", responseRecorder.Code)
	}
	wwwAuth := responseRecorder.Header().Get("WWW-Authenticate")
	if !strings.HasPrefix(wwwAuth, "Digest ") {
		t.Errorf("expected WWW-Authenticate to start with 'Digest ', got %q", wwwAuth)
	}
}

func TestAuthMockServer_handleDigest_ValidResponse(t *testing.T) {
	t.Parallel()

	server := newAuthMockServer(authServerDigest, "127.0.0.1:0", nil, nil)
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("Authorization", `Digest username="test", realm="test", nonce="abc", uri="/", response="def"`)

	server.handleDigest(responseRecorder, request)

	if responseRecorder.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", responseRecorder.Code)
	}
}

func TestAuthMockServer_handleOAuth2_InvalidMethod(t *testing.T) {
	t.Parallel()

	server := newAuthMockServer(authServerOAuth2, "127.0.0.1:0", nil, nil)
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/token", nil)

	server.handleOAuth2(responseRecorder, request)

	if responseRecorder.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", responseRecorder.Code)
	}
}

func TestAuthMockServer_handleOAuth2_CC_ValidRequest(t *testing.T) {
	t.Parallel()

	server := newAuthMockServer(authServerOAuth2, "127.0.0.1:0", nil, nil)
	responseRecorder := httptest.NewRecorder()
	body := strings.NewReader("grant_type=client_credentials&client_id=test-client")
	request := httptest.NewRequest(http.MethodPost, "/token", body)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	server.handleOAuth2(responseRecorder, request)

	if responseRecorder.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", responseRecorder.Code)
	}

	var response map[string]any
	if err := json.NewDecoder(responseRecorder.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if response["access_token"] == "" {
		t.Error("expected access_token to be non-empty")
	}
}

func TestAuthMockServer_handleOAuth2_Password_ValidRequest(t *testing.T) {
	t.Parallel()

	server := newAuthMockServer(authServerOAuth2, "127.0.0.1:0", nil, nil)
	responseRecorder := httptest.NewRecorder()
	body := strings.NewReader("grant_type=password&username=alice&password=secret")
	request := httptest.NewRequest(http.MethodPost, "/token", body)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	server.handleOAuth2(responseRecorder, request)

	if responseRecorder.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", responseRecorder.Code)
	}

	var response map[string]any
	if err := json.NewDecoder(responseRecorder.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if response["access_token"] == "" {
		t.Error("expected access_token to be non-empty")
	}
}

func TestAuthMockServer_handleOAuth2_InvalidGrantType(t *testing.T) {
	t.Parallel()

	server := newAuthMockServer(authServerOAuth2, "127.0.0.1:0", nil, nil)
	responseRecorder := httptest.NewRecorder()
	body := strings.NewReader("grant_type=invalid_grant")
	request := httptest.NewRequest(http.MethodPost, "/token", body)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	server.handleOAuth2(responseRecorder, request)

	if responseRecorder.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", responseRecorder.Code)
	}
}

func TestParseDigestAuthorization(t *testing.T) {
	t.Parallel()

	server := &authMockServer{}
	authorization := `Digest username="test", realm="example", nonce="abc123", uri="/", response="def456", opaque="xyz789", algorithm=MD5, qop=auth`

	params := server.parseDigestAuthorization(authorization)

	expected := map[string]string{
		"username":  "test",
		"realm":     "example",
		"nonce":     "abc123",
		"uri":       "/",
		"response":  "def456",
		"opaque":    "xyz789",
		"algorithm": "MD5",
		"qop":       "auth",
	}

	for key, expectedValue := range expected {
		if params[key] != expectedValue {
			t.Errorf("expected %q=%q, got %q", key, expectedValue, params[key])
		}
	}
}

func TestAuthMockServer_handleHMAC_Valid(t *testing.T) {
	t.Parallel()

	server := newAuthMockServer(authServerHMAC, "127.0.0.1:0", nil, nil)
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api?symbol=BTCUSDT&signature=abc123&timestamp=1770736694138", nil)
	request.Header.Set("X-MBX-APIKEY", "test-api-key")

	server.handleHMAC(responseRecorder, request)

	if responseRecorder.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", responseRecorder.Code)
	}
}

func TestAuthMockServer_handleHMAC_MissingAPIKey(t *testing.T) {
	t.Parallel()

	server := newAuthMockServer(authServerHMAC, "127.0.0.1:0", nil, nil)
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api?signature=abc&timestamp=123", nil)

	server.handleHMAC(responseRecorder, request)

	if responseRecorder.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", responseRecorder.Code)
	}
}

func TestAuthMockServer_handleHMAC_MissingSignature(t *testing.T) {
	t.Parallel()

	server := newAuthMockServer(authServerHMAC, "127.0.0.1:0", nil, nil)
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api?timestamp=123", nil)
	request.Header.Set("X-MBX-APIKEY", "key")

	server.handleHMAC(responseRecorder, request)

	if responseRecorder.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", responseRecorder.Code)
	}
}

func TestAuthMockServer_handleHMAC_MissingTimestamp(t *testing.T) {
	t.Parallel()

	server := newAuthMockServer(authServerHMAC, "127.0.0.1:0", nil, nil)
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api?signature=abc", nil)
	request.Header.Set("X-MBX-APIKEY", "key")

	server.handleHMAC(responseRecorder, request)

	if responseRecorder.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", responseRecorder.Code)
	}
}
