package auth

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// digestTestServer simulates a real HTTP server that challenges with Digest auth.
type digestTestServer struct {
	srv       *httptest.Server
	realm     string
	nonce     string
	opaque    string
	qop       string
	algorithm string
	username  string
	password  string
	reqCount  atomic.Int32
}

func newDigestTestServer(t *testing.T, username, password string) *digestTestServer {
	t.Helper()

	ds := &digestTestServer{
		realm:     "test-realm",
		nonce:     "test-nonce-abc123",
		opaque:    "test-opaque-xyz789",
		qop:       "auth",
		algorithm: "MD5",
		username:  username,
		password:  password,
	}

	ds.srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ds.reqCount.Add(1)

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.Header().Set("WWW-Authenticate",
				fmt.Sprintf(`Digest realm="%s", nonce="%s", opaque="%s", qop="%s", algorithm="%s"`,
					ds.realm, ds.nonce, ds.opaque, ds.qop, ds.algorithm))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(authHeader, "Digest ") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		params := parseDigestParams(authHeader[7:])

		if params["username"] != ds.username {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if params["realm"] != ds.realm {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if params["nonce"] != ds.nonce {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if params["opaque"] != ds.opaque {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		expectedResponse := ds.computeResponse(params, r.Method, r.URL.RequestURI())
		if params["response"] != expectedResponse {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	t.Cleanup(ds.srv.Close)

	return ds
}

func parseDigestParams(auth string) map[string]string {
	params := make(map[string]string)
	for part := range strings.SplitSeq(auth, ",") {
		part = strings.TrimSpace(part)
		key, val, found := strings.Cut(part, "=")
		if !found {
			continue
		}
		key = strings.TrimSpace(key)
		val = strings.Trim(strings.TrimSpace(val), "\"")
		params[key] = val
	}
	return params
}

func (ds *digestTestServer) computeResponse(params map[string]string, method, uri string) string {
	ha1Input := fmt.Sprintf("%s:%s:%s", params["username"], params["realm"], ds.password)
	ha1 := md5hex(ha1Input)

	ha2Input := fmt.Sprintf("%s:%s", method, uri)
	ha2 := md5hex(ha2Input)

	respInput := fmt.Sprintf("%s:%s:%s:%s:%s:%s",
		ha1, params["nonce"], params["nc"], params["cnonce"], params["qop"], ha2)
	return md5hex(respInput)
}

//nolint:gocognit,gocyclo,cyclop // test table with many sub-tests
func TestDigestAuthClient_Apply(t *testing.T) {
	t.Parallel()

	t.Run("successful digest auth with challenge-response", func(t *testing.T) {
		t.Parallel()

		ds := newDigestTestServer(t, "digest-user", "digest-pass")

		client := &DigestAuthClient{
			Username: "digest-user",
			Password: "digest-pass",
		}
		if err := client.New(); err != nil {
			t.Fatalf("New() = %v", err)
		}

		req, _ := http.NewRequest(http.MethodGet, ds.srv.URL+"/api/resource", nil)
		var info Info
		if err := client.Apply(req, &info); err != nil {
			t.Fatalf("Apply() = %v", err)
		}

		auth := req.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Digest ") {
			t.Fatalf("Authorization = %q, want Digest prefix", auth)
		}
		if !strings.Contains(auth, `username="digest-user"`) {
			t.Errorf("Authorization missing username: %s", auth)
		}
		if !strings.Contains(auth, `realm="test-realm"`) {
			t.Errorf("Authorization missing realm: %s", auth)
		}
		if !strings.Contains(auth, `nonce="test-nonce-abc123"`) {
			t.Errorf("Authorization missing nonce: %s", auth)
		}
		if !strings.Contains(auth, `opaque="test-opaque-xyz789"`) {
			t.Errorf("Authorization missing opaque: %s", auth)
		}
		if !strings.Contains(auth, `response="`) {
			t.Errorf("Authorization missing response: %s", auth)
		}
		if !strings.Contains(auth, `qop=auth`) {
			t.Errorf("Authorization missing qop: %s", auth)
		}
		if !strings.Contains(auth, `nc=`) {
			t.Errorf("Authorization missing nc: %s", auth)
		}
		if !strings.Contains(auth, `cnonce="`) {
			t.Errorf("Authorization missing cnonce: %s", auth)
		}

		if v := info.Headers["Authorization"]; v != auth {
			t.Errorf("info.Headers[Authorization] = %q, want %q", v, auth)
		}

		if ds.reqCount.Load() != 1 {
			t.Errorf("expected 1 request (401 challenge), got %d", ds.reqCount.Load())
		}
	})

	t.Run("caches challenge and reuses on second Apply", func(t *testing.T) {
		t.Parallel()

		ds := newDigestTestServer(t, "user", "pass")

		client := &DigestAuthClient{
			Username: "user",
			Password: "pass",
		}
		if err := client.New(); err != nil {
			t.Fatalf("New() = %v", err)
		}

		req1, _ := http.NewRequest(http.MethodGet, ds.srv.URL+"/api", nil)
		if err := client.Apply(req1, nil); err != nil {
			t.Fatalf("Apply #1 = %v", err)
		}

		req2, _ := http.NewRequest(http.MethodGet, ds.srv.URL+"/api", nil)
		if err := client.Apply(req2, nil); err != nil {
			t.Fatalf("Apply #2 = %v", err)
		}

		// First Apply: 2 requests (401 + authed). Second Apply: 0 requests (cached).
		if ds.reqCount.Load() != 1 {
			t.Errorf("expected 1 total request (cached), got %d", ds.reqCount.Load())
		}

		auth := req2.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Digest ") {
			t.Fatalf("Authorization = %q, want Digest prefix", auth)
		}
	})

	t.Run("refetches challenge after nonce TTL expires", func(t *testing.T) {
		t.Parallel()

		ds := newDigestTestServer(t, "user", "pass")

		client := &DigestAuthClient{
			Username: "user",
			Password: "pass",
		}
		if err := client.New(); err != nil {
			t.Fatalf("New() = %v", err)
		}

		req1, _ := http.NewRequest(http.MethodGet, ds.srv.URL+"/api", nil)
		if err := client.Apply(req1, nil); err != nil {
			t.Fatalf("Apply #1 = %v", err)
		}

		// Override cachedAt to force expiration
		client.mu.Lock()
		client.cachedAt = time.Now().Add(-10 * time.Minute)
		client.mu.Unlock()

		req2, _ := http.NewRequest(http.MethodGet, ds.srv.URL+"/api", nil)
		if err := client.Apply(req2, nil); err != nil {
			t.Fatalf("Apply #2 = %v", err)
		}

		// Each Apply makes 2 requests (401 + authed) = 4 total
		if ds.reqCount.Load() != 2 {
			t.Errorf("expected 2 total requests (1 per Apply), got %d", ds.reqCount.Load())
		}
	})

	t.Run("increments nonce count on each Apply", func(t *testing.T) {
		t.Parallel()

		ds := newDigestTestServer(t, "user", "pass")

		client := &DigestAuthClient{
			Username: "user",
			Password: "pass",
		}
		if err := client.New(); err != nil {
			t.Fatalf("New() = %v", err)
		}

		req1, _ := http.NewRequest(http.MethodGet, ds.srv.URL+"/api", nil)
		if err := client.Apply(req1, nil); err != nil {
			t.Fatalf("Apply #1 = %v", err)
		}

		req2, _ := http.NewRequest(http.MethodGet, ds.srv.URL+"/api", nil)
		if err := client.Apply(req2, nil); err != nil {
			t.Fatalf("Apply #2 = %v", err)
		}

		auth1 := req1.Header.Get("Authorization")
		auth2 := req2.Header.Get("Authorization")

		nc1 := extractNCDigest(t, auth1)
		nc2 := extractNCDigest(t, auth2)

		if nc1 == "" {
			t.Fatal("nc not found in first auth header")
		}
		if nc2 == "" {
			t.Fatal("nc not found in second auth header")
		}
		if nc1 == nc2 {
			t.Errorf("nc should increment, got %s both times", nc1)
		}
	})

	t.Run("returns error on non-401 challenge response", func(t *testing.T) {
		t.Parallel()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusForbidden)
		}))
		t.Cleanup(srv.Close)

		client := &DigestAuthClient{
			Username: "u",
			Password: "p",
		}
		if err := client.New(); err != nil {
			t.Fatalf("New() = %v", err)
		}

		req, _ := http.NewRequest(http.MethodGet, srv.URL+"/api", nil)
		if err := client.Apply(req, nil); err == nil {
			t.Fatal("expected error for non-401 response, got nil")
		}
	})

	t.Run("returns error on missing WWW-Authenticate header", func(t *testing.T) {
		t.Parallel()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		t.Cleanup(srv.Close)

		client := &DigestAuthClient{
			Username: "u",
			Password: "p",
		}
		if err := client.New(); err != nil {
			t.Fatalf("New() = %v", err)
		}

		req, _ := http.NewRequest(http.MethodGet, srv.URL+"/api", nil)
		if err := client.Apply(req, nil); err == nil {
			t.Fatal("expected error for missing WWW-Authenticate, got nil")
		}
	})
}

func extractNCDigest(t *testing.T, auth string) string {
	t.Helper()
	for part := range strings.SplitSeq(auth, ",") {
		part = strings.TrimSpace(part)
		key, val, found := strings.Cut(part, "=")
		if !found {
			continue
		}
		if strings.TrimSpace(key) == "nc" {
			return strings.Trim(strings.TrimSpace(val), "\"")
		}
	}
	return ""
}
