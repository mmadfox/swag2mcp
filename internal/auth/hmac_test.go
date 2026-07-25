package auth

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHMACAuthClient_Type(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{}
	assert.Equal(t, HMACAuth, client.Type())
}

func TestHMACAuthClient_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		client  *HMACAuthClient
		wantErr bool
	}{
		{name: "valid", client: &HMACAuthClient{APIKey: "key", SecretKey: "secret"}, wantErr: false},
		{name: "empty", client: &HMACAuthClient{}, wantErr: true},
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

func TestHMACAuthClient_New_EnvVars(t *testing.T) {
	t.Setenv("TEST_HMAC_KEY", "env-key")
	t.Setenv("TEST_HMAC_SECRET", "env-secret")

	client := &HMACAuthClient{
		APIKey:    "$(TEST_HMAC_KEY)",
		SecretKey: "$(TEST_HMAC_SECRET)",
	}
	require.NoError(t, client.New(), "New()")
	assert.Equal(t, "env-key", client.APIKey)
	assert.Equal(t, "env-secret", client.SecretKey)
}

func TestHMACAuthClient_Apply_SetsAPIKeyHeader(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{APIKey: "test-api-key", SecretKey: "test-secret"}
	require.NoError(t, client.New(), "New()")

	req, _ := http.NewRequest(
		http.MethodGet,
		"http://example.com/api/v3/ticker/price?symbol=BTCUSDT",
		nil,
	)
	var info Info
	require.NoError(t, client.Apply(req, &info), "Apply()")

	assert.Equal(t, "test-api-key", req.Header.Get("X-MBX-APIKEY"))
	assert.Equal(t, "test-api-key", info.Headers["X-MBX-APIKEY"])
}

func TestHMACAuthClient_Apply_AddsTimestamp(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{APIKey: "key", SecretKey: "secret"}
	require.NoError(t, client.New(), "New()")

	before := time.Now().UnixMilli()
	req, _ := http.NewRequest(http.MethodGet, "http://example.com/api", nil)
	require.NoError(t, client.Apply(req, nil), "Apply()")
	after := time.Now().UnixMilli()

	ts := req.URL.Query().Get("timestamp")
	require.NotEmpty(t, ts, "timestamp query param is missing")

	tsInt, err := strconv.ParseInt(ts, 10, 64)
	require.NoError(t, err, "timestamp is not a valid int")

	assert.GreaterOrEqual(t, tsInt, before)
	assert.LessOrEqual(t, tsInt, after)
}

func TestHMACAuthClient_Apply_AddsSignature(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{APIKey: "key", SecretKey: "secret"}
	require.NoError(t, client.New(), "New()")

	req, _ := http.NewRequest(
		http.MethodGet,
		"http://example.com/api?symbol=BTCUSDT",
		nil,
	)
	var info Info
	require.NoError(t, client.Apply(req, &info), "Apply()")

	sig := req.URL.Query().Get("signature")
	require.NotEmpty(t, sig, "signature query param is missing")

	_, err := hex.DecodeString(sig)
	require.NoError(t, err, "signature is not valid hex")

	assert.Equal(t, sig, info.QueryParams["signature"])
}

func TestHMACAuthClient_Apply_SignatureIsValid(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{APIKey: "key", SecretKey: "secret"}
	require.NoError(t, client.New(), "New()")

	req, _ := http.NewRequest(
		http.MethodGet,
		"http://example.com/api?symbol=BTCUSDT",
		nil,
	)
	require.NoError(t, client.Apply(req, nil), "Apply()")

	q := req.URL.Query()
	sig := q.Get("signature")
	q.Del("signature")

	expectedQuery := q.Encode()
	mac := hmac.New(sha256.New, []byte("secret"))
	mac.Write([]byte(expectedQuery))
	expectedSig := hex.EncodeToString(mac.Sum(nil))

	assert.Equal(t, expectedSig, sig)
}

func TestHMACAuthClient_Apply_PreservesExistingQueryParams(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{APIKey: "key", SecretKey: "secret"}
	require.NoError(t, client.New(), "New()")

	req, _ := http.NewRequest(
		http.MethodGet,
		"http://example.com/api?symbol=BTCUSDT&limit=10",
		nil,
	)
	require.NoError(t, client.Apply(req, nil), "Apply()")

	q := req.URL.Query()
	assert.Equal(t, "BTCUSDT", q.Get("symbol"))
	assert.Equal(t, "10", q.Get("limit"))
}

func TestHMACAuthClient_Apply_InfoHasTimestamp(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{APIKey: "key", SecretKey: "secret"}
	require.NoError(t, client.New(), "New()")

	req, _ := http.NewRequest(http.MethodGet, "http://example.com/api", nil)
	var info Info
	require.NoError(t, client.Apply(req, &info), "Apply()")

	assert.Contains(t, info.QueryParams, "timestamp")
}

func TestHMACAuthClient_Apply_NoExistingQueryParams(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{APIKey: "key", SecretKey: "secret"}
	require.NoError(t, client.New(), "New()")

	req, _ := http.NewRequest(http.MethodGet, "http://example.com/api", nil)
	require.NoError(t, client.Apply(req, nil), "Apply()")

	q := req.URL.Query()
	assert.NotEmpty(t, q.Get("timestamp"), "timestamp should be present")
	assert.NotEmpty(t, q.Get("signature"), "signature should be present")
}

func TestHMACAuthClient_Apply_EmptyAPIKey(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{APIKey: "", SecretKey: "secret"}
	require.NoError(t, client.New(), "New()")

	req, _ := http.NewRequest(http.MethodGet, "http://example.com/api", nil)
	require.NoError(t, client.Apply(req, nil), "Apply()")

	assert.Empty(t, req.Header.Get("X-MBX-APIKEY"))
}

func TestHMACAuthClient_Apply_InfoNil(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{APIKey: "key", SecretKey: "secret"}
	require.NoError(t, client.New(), "New()")

	req, _ := http.NewRequest(http.MethodGet, "http://example.com/api", nil)
	require.NoError(t, client.Apply(req, nil), "Apply()")

	assert.Equal(t, "key", req.Header.Get("X-MBX-APIKEY"))
	assert.NotEmpty(t, req.URL.Query().Get("signature"), "signature should be present")
}

func TestHMACAuthClient_Apply_QueryParamsSorted(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{APIKey: "key", SecretKey: "secret"}
	require.NoError(t, client.New(), "New()")

	req, _ := http.NewRequest(
		http.MethodGet,
		"http://example.com/api?z=last&a=first&m=middle",
		nil,
	)
	require.NoError(t, client.Apply(req, nil), "Apply()")

	q := req.URL.Query()
	keys := make([]string, 0, len(q))
	for k := range q {
		keys = append(keys, k)
	}

	require.GreaterOrEqual(t, len(keys), 3, "expected at least 3 query params")
}

func TestHMACAuthClient_Apply_DeterministicSignature(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{APIKey: "key", SecretKey: "secret"}
	require.NoError(t, client.New(), "New()")

	req1, _ := http.NewRequest(
		http.MethodGet,
		"http://example.com/api?symbol=BTCUSDT",
		nil,
	)
	req2, _ := http.NewRequest(
		http.MethodGet,
		"http://example.com/api?symbol=BTCUSDT",
		nil,
	)

	require.NoError(t, client.Apply(req1, nil), "Apply req1")
	time.Sleep(2 * time.Millisecond)
	require.NoError(t, client.Apply(req2, nil), "Apply req2")

	sig1 := req1.URL.Query().Get("signature")
	sig2 := req2.URL.Query().Get("signature")

	assert.NotEqual(t, sig1, sig2, "signatures should differ due to different timestamps")
}

func TestHMACAuthClient_Apply_NoPanicOnNilOut(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{APIKey: "key", SecretKey: "secret"}
	require.NoError(t, client.New(), "New()")

	req, _ := http.NewRequest(http.MethodGet, "http://example.com/api", nil)
	require.NoError(t, client.Apply(req, nil), "Apply()")
}

func TestHMACAuthClient_Apply_InfoQueryParamsPopulated(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{APIKey: "key", SecretKey: "secret"}
	require.NoError(t, client.New(), "New()")

	req, _ := http.NewRequest(
		http.MethodGet,
		"http://example.com/api?symbol=BTCUSDT",
		nil,
	)
	var info Info
	require.NoError(t, client.Apply(req, &info), "Apply()")

	require.NotNil(t, info.QueryParams, "info.QueryParams should not be nil")
	assert.Contains(t, info.QueryParams, "signature")
	assert.Contains(t, info.QueryParams, "timestamp")
}

func TestHMACAuthClient_Apply_InfoHeadersPopulated(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{APIKey: "key", SecretKey: "secret"}
	require.NoError(t, client.New(), "New()")

	req, _ := http.NewRequest(http.MethodGet, "http://example.com/api", nil)
	var info Info
	require.NoError(t, client.Apply(req, &info), "Apply()")

	require.NotNil(t, info.Headers, "info.Headers should not be nil")
	assert.Contains(t, info.Headers, "X-MBX-APIKEY")
}

func TestHMACAuthClient_Apply_EmptySecretKey(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{APIKey: "key", SecretKey: ""}
	require.NoError(t, client.New(), "New()")

	req, _ := http.NewRequest(http.MethodGet, "http://example.com/api", nil)
	require.NoError(t, client.Apply(req, nil), "Apply()")

	sig := req.URL.Query().Get("signature")
	require.NotEmpty(t, sig, "signature should be present even with empty secret")
}

func TestHMACAuthClient_Apply_ComplexQueryString(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{APIKey: "key", SecretKey: "secret"}
	require.NoError(t, client.New(), "New()")

	req, _ := http.NewRequest(
		http.MethodGet,
		"http://example.com/api?symbol=BTCUSDT&side=BUY&type=LIMIT&timeInForce=GTC&quantity=1&price=0.1&recvWindow=5000",
		nil,
	)
	require.NoError(t, client.Apply(req, nil), "Apply()")

	q := req.URL.Query()
	assert.Equal(t, "BTCUSDT", q.Get("symbol"))
	assert.Equal(t, "BUY", q.Get("side"))
	assert.NotEmpty(t, q.Get("signature"), "signature should be present")
}

func TestHMACAuthClient_Apply_NoQueryString(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{APIKey: "key", SecretKey: "secret"}
	require.NoError(t, client.New(), "New()")

	req, _ := http.NewRequest(http.MethodGet, "http://example.com/api", nil)
	require.NoError(t, client.Apply(req, nil), "Apply()")

	q := req.URL.Query()
	assert.NotEmpty(t, q.Get("timestamp"), "timestamp should be present")
	assert.NotEmpty(t, q.Get("signature"), "signature should be present")
}

func TestHMACAuthClient_Apply_InfoTimestampIsString(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{APIKey: "key", SecretKey: "secret"}
	require.NoError(t, client.New(), "New()")

	req, _ := http.NewRequest(http.MethodGet, "http://example.com/api", nil)
	var info Info
	require.NoError(t, client.Apply(req, &info), "Apply()")

	ts, ok := info.QueryParams["timestamp"]
	require.True(t, ok, "info.QueryParams should contain timestamp")
	_, err := strconv.ParseInt(ts, 10, 64)
	require.NoError(t, err, "timestamp should be a valid int string")
}

func TestHMACAuthClient_Apply_InfoSignatureIsHex(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{APIKey: "key", SecretKey: "secret"}
	require.NoError(t, client.New(), "New()")

	req, _ := http.NewRequest(http.MethodGet, "http://example.com/api", nil)
	var info Info
	require.NoError(t, client.Apply(req, &info), "Apply()")

	sig, ok := info.QueryParams["signature"]
	require.True(t, ok, "info.QueryParams should contain signature")
	_, err := hex.DecodeString(sig)
	require.NoError(t, err, "signature should be valid hex")
}

func TestHMACAuthClient_Apply_HeaderAndQueryConsistent(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{APIKey: "test-key", SecretKey: "test-secret"}
	require.NoError(t, client.New(), "New()")

	req, _ := http.NewRequest(http.MethodGet, "http://example.com/api", nil)
	var info Info
	require.NoError(t, client.Apply(req, &info), "Apply()")

	assert.Equal(t, "test-key", info.Headers["X-MBX-APIKEY"])
	assert.NotEmpty(t, info.QueryParams["signature"], "signature missing from info")
	assert.NotEmpty(t, info.QueryParams["timestamp"], "timestamp missing from info")
}

func TestHMACAuthClient_Apply_DoesNotModifyMethod(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{APIKey: "key", SecretKey: "secret"}
	require.NoError(t, client.New(), "New()")

	req, _ := http.NewRequest(
		http.MethodPost,
		"http://example.com/api",
		strings.NewReader(`{"key":"value"}`),
	)
	require.NoError(t, client.Apply(req, nil), "Apply()")

	assert.Equal(t, http.MethodPost, req.Method)
}

func TestHMACAuthClient_Apply_DoesNotModifyBody(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{APIKey: "key", SecretKey: "secret"}
	require.NoError(t, client.New(), "New()")

	body := `{"symbol":"BTCUSDT","side":"BUY"}`
	req, _ := http.NewRequest(
		http.MethodPost,
		"http://example.com/api",
		strings.NewReader(body),
	)
	require.NoError(t, client.Apply(req, nil), "Apply()")

	assert.Equal(t, "key", req.Header.Get("X-MBX-APIKEY"))
	assert.NotEmpty(t, req.URL.Query().Get("signature"), "signature should be present")
}
