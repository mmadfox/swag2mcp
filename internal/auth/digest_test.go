package auth

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

		authHeader := r.Header.Get(headerAuthorization)
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

func TestDigestAuthClient_Apply(t *testing.T) {
	t.Parallel()

	t.Run("successful digest auth with challenge-response", func(t *testing.T) {
		t.Parallel()

		ds := newDigestTestServer(t, "digest-user", "digest-pass")

		client := &DigestAuthClient{
			Username: "digest-user",
			Password: "digest-pass",
		}
		require.NoError(t, client.New(), "New()")

		req, _ := http.NewRequest(http.MethodGet, ds.srv.URL+"/api/resource", nil)
		var info Info
		require.NoError(t, client.Apply(req, &info), "Apply()")

		auth := req.Header.Get(headerAuthorization)
		require.True(t, strings.HasPrefix(auth, "Digest "), "Authorization should have Digest prefix")

		assert.Contains(t, auth, `username="digest-user"`)
		assert.Contains(t, auth, `realm="test-realm"`)
		assert.Contains(t, auth, `nonce="test-nonce-abc123"`)
		assert.Contains(t, auth, `opaque="test-opaque-xyz789"`)
		assert.Contains(t, auth, `response="`)
		assert.Contains(t, auth, `qop=auth`)
		assert.Contains(t, auth, `nc=`)
		assert.Contains(t, auth, `cnonce="`)

		assert.Equal(t, auth, info.Headers[headerAuthorization])

		assert.Equal(t, int32(1), ds.reqCount.Load(), "expected 1 request (401 challenge)")
	})

	t.Run("caches challenge and reuses on second Apply", func(t *testing.T) {
		t.Parallel()

		ds := newDigestTestServer(t, "user", "pass")

		client := &DigestAuthClient{
			Username: "user",
			Password: "pass",
		}
		require.NoError(t, client.New(), "New()")

		req1, _ := http.NewRequest(http.MethodGet, ds.srv.URL+"/api", nil)
		require.NoError(t, client.Apply(req1, nil), "Apply #1")

		req2, _ := http.NewRequest(http.MethodGet, ds.srv.URL+"/api", nil)
		require.NoError(t, client.Apply(req2, nil), "Apply #2")

		// First Apply: 2 requests (401 + authed). Second Apply: 0 requests (cached).
		assert.Equal(t, int32(1), ds.reqCount.Load(), "expected 1 total request (cached)")

		auth := req2.Header.Get(headerAuthorization)
		assert.True(t, strings.HasPrefix(auth, "Digest "), "Authorization should have Digest prefix")
	})

	t.Run("refetches challenge after nonce TTL expires", func(t *testing.T) {
		t.Parallel()

		ds := newDigestTestServer(t, "user", "pass")

		client := &DigestAuthClient{
			Username: "user",
			Password: "pass",
		}
		require.NoError(t, client.New(), "New()")

		req1, _ := http.NewRequest(http.MethodGet, ds.srv.URL+"/api", nil)
		require.NoError(t, client.Apply(req1, nil), "Apply #1")

		// Override cachedAt to force expiration
		client.mu.Lock()
		client.cachedAt = time.Now().Add(-10 * time.Minute)
		client.mu.Unlock()

		req2, _ := http.NewRequest(http.MethodGet, ds.srv.URL+"/api", nil)
		require.NoError(t, client.Apply(req2, nil), "Apply #2")

		// Each Apply makes 2 requests (401 + authed) = 4 total
		assert.Equal(t, int32(2), ds.reqCount.Load(), "expected 2 total requests (1 per Apply)")
	})

	t.Run("increments nonce count on each Apply", func(t *testing.T) {
		t.Parallel()

		ds := newDigestTestServer(t, "user", "pass")

		client := &DigestAuthClient{
			Username: "user",
			Password: "pass",
		}
		require.NoError(t, client.New(), "New()")

		req1, _ := http.NewRequest(http.MethodGet, ds.srv.URL+"/api", nil)
		require.NoError(t, client.Apply(req1, nil), "Apply #1")

		req2, _ := http.NewRequest(http.MethodGet, ds.srv.URL+"/api", nil)
		require.NoError(t, client.Apply(req2, nil), "Apply #2")

		auth1 := req1.Header.Get(headerAuthorization)
		auth2 := req2.Header.Get(headerAuthorization)

		nc1 := extractNCDigest(t, auth1)
		nc2 := extractNCDigest(t, auth2)

		require.NotEmpty(t, nc1, "nc not found in first auth header")
		require.NotEmpty(t, nc2, "nc not found in second auth header")
		assert.NotEqual(t, nc1, nc2, "nc should increment")
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
		require.NoError(t, client.New(), "New()")

		req, _ := http.NewRequest(http.MethodGet, srv.URL+"/api", nil)
		err := client.Apply(req, nil)
		require.Error(t, err, "expected error for non-401 response")
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
		require.NoError(t, client.New(), "New()")

		req, _ := http.NewRequest(http.MethodGet, srv.URL+"/api", nil)
		err := client.Apply(req, nil)
		require.Error(t, err, "expected error for missing WWW-Authenticate")
	})
}

func TestDigestAuthClient_New(t *testing.T) {
	t.Parallel()

	client := &DigestAuthClient{Username: "u", Password: "p"}
	require.NoError(t, client.New())
}

func TestDigestAuthClient_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		client  *DigestAuthClient
		wantErr bool
	}{
		{name: "valid", client: &DigestAuthClient{Username: "u", Password: "p"}, wantErr: false},
		{name: "empty", client: &DigestAuthClient{}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.client.Validate()
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestDigestAuthClient_Type(t *testing.T) {
	t.Parallel()

	client := &DigestAuthClient{}
	assert.Equal(t, DigestAuth, client.Type())
}

func TestDigestAuthClient_SetMockBaseURL(t *testing.T) {
	t.Parallel()

	client := &DigestAuthClient{}
	client.SetMockBaseURL("http://localhost:9091/")
	assert.Equal(t, "http://localhost:9091/", client.MockBaseURL)
}

func TestDigestAuthClient_fetchChallenge_Non401(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()

	client := &DigestAuthClient{Username: "u", Password: "p"}
	req, _ := http.NewRequest(http.MethodGet, srv.URL, nil)
	_, err := client.fetchChallenge(context.Background(), req)
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
	_, err := client.fetchChallenge(context.Background(), req)
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
	_, err := client.fetchChallenge(context.Background(), req)
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

func TestParseDigestChallenge(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		header     string
		wantRealm  string
		wantNonce  string
		wantOpaque string
		wantQop    string
		wantAlgo   string
	}{
		{
			name:       "full challenge",
			header:     `Digest realm="test-realm", nonce="test-nonce", opaque="test-opaque", qop="auth", algorithm="MD5"`,
			wantRealm:  "test-realm",
			wantNonce:  "test-nonce",
			wantOpaque: "test-opaque",
			wantQop:    "auth",
			wantAlgo:   "MD5",
		},
		{
			name:      "defaults algorithm to MD5",
			header:    `Digest realm="r", nonce="n"`,
			wantRealm: "r",
			wantNonce: "n",
			wantAlgo:  "MD5",
		},
		{
			name:     "empty challenge",
			header:   `Digest `,
			wantAlgo: "MD5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ch := parseDigestChallenge(tt.header)
			assert.Equal(t, tt.wantRealm, ch.realm)
			assert.Equal(t, tt.wantNonce, ch.nonce)
			assert.Equal(t, tt.wantOpaque, ch.opaque)
			assert.Equal(t, tt.wantQop, ch.qop)
			assert.Equal(t, tt.wantAlgo, ch.algorithm)
		})
	}
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

	tests := []struct {
		name   string
		ch     digestChallenge
		hasQop bool
	}{
		{
			name: "with qop",
			ch: digestChallenge{
				realm: "test-realm", nonce: "test-nonce",
				opaque: "test-opaque", qop: "auth", algorithm: "MD5",
			},
			hasQop: true,
		},
		{
			name: "without qop",
			ch: digestChallenge{
				realm: "test-realm", nonce: "test-nonce",
			},
			hasQop: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			auth := client.buildDigest("GET", "/api", tt.ch, 1, "cnonce123")
			assert.True(t, strings.HasPrefix(auth, "Digest "), "expected Digest prefix")
			assert.Contains(t, auth, `username="user"`)
			assert.Contains(t, auth, `realm="test-realm"`)
			assert.Contains(t, auth, `nonce="test-nonce"`)
			assert.Contains(t, auth, `response="`)
			if tt.hasQop {
				assert.Contains(t, auth, `opaque="test-opaque"`)
				assert.Contains(t, auth, `qop=auth`)
				assert.Contains(t, auth, `nc=00000001`)
				assert.Contains(t, auth, `cnonce="cnonce123"`)
			} else {
				assert.NotContains(t, auth, "qop=")
			}
		})
	}
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
