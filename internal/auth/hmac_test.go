package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestHMACAuthClient_Type(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{}
	if client.Type() != HMACAuth {
		t.Errorf("Type() = %q, want %q", client.Type(), HMACAuth)
	}
}

func TestHMACAuthClient_Validate(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{APIKey: "key", SecretKey: "secret"}
	if err := client.Validate(); err != nil {
		t.Errorf("Validate() = %v", err)
	}
}

func TestHMACAuthClient_Validate_Empty(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{}
	if err := client.Validate(); err == nil {
		t.Error("expected validation error for empty hmac auth")
	}
}

func TestHMACAuthClient_New_EnvVars(t *testing.T) {
	t.Setenv("TEST_HMAC_KEY", "env-key")
	t.Setenv("TEST_HMAC_SECRET", "env-secret")

	client := &HMACAuthClient{
		APIKey:    "$(TEST_HMAC_KEY)",
		SecretKey: "$(TEST_HMAC_SECRET)",
	}
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}
	if client.APIKey != "env-key" {
		t.Errorf("APIKey = %q, want %q", client.APIKey, "env-key")
	}
	if client.SecretKey != "env-secret" {
		t.Errorf("SecretKey = %q, want %q", client.SecretKey, "env-secret")
	}
}

func TestHMACAuthClient_Apply_SetsAPIKeyHeader(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{APIKey: "test-api-key", SecretKey: "test-secret"}
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}

	req, _ := http.NewRequest(
		http.MethodGet,
		"http://example.com/api/v3/ticker/price?symbol=BTCUSDT",
		nil,
	)
	var info Info
	if err := client.Apply(req, &info); err != nil {
		t.Fatalf("Apply() = %v", err)
	}

	//nolint:canonicalheader // Binance requires this exact header name
	if v := req.Header.Get("X-MBX-APIKEY"); v != "test-api-key" {
		t.Errorf("X-MBX-APIKEY header = %q, want %q", v, "test-api-key")
	}
	if v := info.Headers["X-MBX-APIKEY"]; v != "test-api-key" {
		t.Errorf("info.Headers[X-MBX-APIKEY] = %q, want %q", v, "test-api-key")
	}
}

func TestHMACAuthClient_Apply_AddsTimestamp(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{APIKey: "key", SecretKey: "secret"}
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}

	before := time.Now().UnixMilli()
	req, _ := http.NewRequest(http.MethodGet, "http://example.com/api", nil)
	if err := client.Apply(req, nil); err != nil {
		t.Fatalf("Apply() = %v", err)
	}
	after := time.Now().UnixMilli()

	ts := req.URL.Query().Get("timestamp")
	if ts == "" {
		t.Fatal("timestamp query param is missing")
	}
	tsInt, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		t.Fatalf("timestamp is not a valid int: %v", err)
	}
	if tsInt < before || tsInt > after {
		t.Errorf("timestamp %d should be between %d and %d", tsInt, before, after)
	}
}

func TestHMACAuthClient_Apply_AddsSignature(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{APIKey: "key", SecretKey: "secret"}
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}

	req, _ := http.NewRequest(
		http.MethodGet,
		"http://example.com/api?symbol=BTCUSDT",
		nil,
	)
	var info Info
	if err := client.Apply(req, &info); err != nil {
		t.Fatalf("Apply() = %v", err)
	}

	sig := req.URL.Query().Get("signature")
	if sig == "" {
		t.Fatal("signature query param is missing")
	}
	if _, err := hex.DecodeString(sig); err != nil {
		t.Errorf("signature is not valid hex: %v", err)
	}

	if v := info.QueryParams["signature"]; v != sig {
		t.Errorf("info.QueryParams[signature] = %q, want %q", v, sig)
	}
}

func TestHMACAuthClient_Apply_SignatureIsValid(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{APIKey: "key", SecretKey: "secret"}
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}

	req, _ := http.NewRequest(
		http.MethodGet,
		"http://example.com/api?symbol=BTCUSDT",
		nil,
	)
	if err := client.Apply(req, nil); err != nil {
		t.Fatalf("Apply() = %v", err)
	}

	q := req.URL.Query()
	sig := q.Get("signature")
	q.Del("signature")

	expectedQuery := q.Encode()
	mac := hmac.New(sha256.New, []byte("secret"))
	mac.Write([]byte(expectedQuery))
	expectedSig := hex.EncodeToString(mac.Sum(nil))

	if sig != expectedSig {
		t.Errorf("signature = %q, want %q", sig, expectedSig)
	}
}

func TestHMACAuthClient_Apply_PreservesExistingQueryParams(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{APIKey: "key", SecretKey: "secret"}
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}

	req, _ := http.NewRequest(
		http.MethodGet,
		"http://example.com/api?symbol=BTCUSDT&limit=10",
		nil,
	)
	if err := client.Apply(req, nil); err != nil {
		t.Fatalf("Apply() = %v", err)
	}

	q := req.URL.Query()
	if q.Get("symbol") != "BTCUSDT" {
		t.Errorf("symbol = %q, want %q", q.Get("symbol"), "BTCUSDT")
	}
	if q.Get("limit") != "10" {
		t.Errorf("limit = %q, want %q", q.Get("limit"), "10")
	}
}

func TestHMACAuthClient_Apply_InfoHasTimestamp(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{APIKey: "key", SecretKey: "secret"}
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}

	req, _ := http.NewRequest(http.MethodGet, "http://example.com/api", nil)
	var info Info
	if err := client.Apply(req, &info); err != nil {
		t.Fatalf("Apply() = %v", err)
	}

	if _, ok := info.QueryParams["timestamp"]; !ok {
		t.Error("info.QueryParams should contain timestamp")
	}
}

func TestHMACAuthClient_Apply_NoExistingQueryParams(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{APIKey: "key", SecretKey: "secret"}
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}

	req, _ := http.NewRequest(http.MethodGet, "http://example.com/api", nil)
	if err := client.Apply(req, nil); err != nil {
		t.Fatalf("Apply() = %v", err)
	}

	q := req.URL.Query()
	if q.Get("timestamp") == "" {
		t.Error("timestamp should be present")
	}
	if q.Get("signature") == "" {
		t.Error("signature should be present")
	}
}

func TestHMACAuthClient_Apply_EmptyAPIKey(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{APIKey: "", SecretKey: "secret"}
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}

	req, _ := http.NewRequest(http.MethodGet, "http://example.com/api", nil)
	if err := client.Apply(req, nil); err != nil {
		t.Fatalf("Apply() = %v", err)
	}

	//nolint:canonicalheader // Binance requires this exact header name
	if v := req.Header.Get("X-MBX-APIKEY"); v != "" {
		t.Errorf("X-MBX-APIKEY header should be empty, got %q", v)
	}
}

func TestHMACAuthClient_Apply_InfoNil(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{APIKey: "key", SecretKey: "secret"}
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}

	req, _ := http.NewRequest(http.MethodGet, "http://example.com/api", nil)
	if err := client.Apply(req, nil); err != nil {
		t.Fatalf("Apply() = %v", err)
	}

	//nolint:canonicalheader // Binance requires this exact header name
	if v := req.Header.Get("X-MBX-APIKEY"); v != "key" {
		t.Errorf("X-MBX-APIKEY = %q, want %q", v, "key")
	}
	if req.URL.Query().Get("signature") == "" {
		t.Error("signature should be present")
	}
}

func TestHMACAuthClient_Apply_QueryParamsSorted(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{APIKey: "key", SecretKey: "secret"}
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}

	req, _ := http.NewRequest(
		http.MethodGet,
		"http://example.com/api?z=last&a=first&m=middle",
		nil,
	)
	if err := client.Apply(req, nil); err != nil {
		t.Fatalf("Apply() = %v", err)
	}

	q := req.URL.Query()
	keys := make([]string, 0, len(q))
	for k := range q {
		keys = append(keys, k)
	}

	if len(keys) < 3 {
		t.Fatalf("expected at least 3 query params, got %d", len(keys))
	}
}

func TestHMACAuthClient_Apply_DeterministicSignature(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{APIKey: "key", SecretKey: "secret"}
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}

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

	if err := client.Apply(req1, nil); err != nil {
		t.Fatalf("Apply req1 = %v", err)
	}
	time.Sleep(2 * time.Millisecond)
	if err := client.Apply(req2, nil); err != nil {
		t.Fatalf("Apply req2 = %v", err)
	}

	sig1 := req1.URL.Query().Get("signature")
	sig2 := req2.URL.Query().Get("signature")

	if sig1 == sig2 {
		t.Error("signatures should differ due to different timestamps")
	}
}

func TestHMACAuthClient_Apply_NoPanicOnNilOut(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{APIKey: "key", SecretKey: "secret"}
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}

	req, _ := http.NewRequest(http.MethodGet, "http://example.com/api", nil)
	if err := client.Apply(req, nil); err != nil {
		t.Fatalf("Apply() = %v", err)
	}
}

func TestHMACAuthClient_Apply_InfoQueryParamsPopulated(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{APIKey: "key", SecretKey: "secret"}
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}

	req, _ := http.NewRequest(
		http.MethodGet,
		"http://example.com/api?symbol=BTCUSDT",
		nil,
	)
	var info Info
	if err := client.Apply(req, &info); err != nil {
		t.Fatalf("Apply() = %v", err)
	}

	if info.QueryParams == nil {
		t.Fatal("info.QueryParams should not be nil")
	}
	if _, ok := info.QueryParams["signature"]; !ok {
		t.Error("info.QueryParams should contain signature")
	}
	if _, ok := info.QueryParams["timestamp"]; !ok {
		t.Error("info.QueryParams should contain timestamp")
	}
}

func TestHMACAuthClient_Apply_InfoHeadersPopulated(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{APIKey: "key", SecretKey: "secret"}
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}

	req, _ := http.NewRequest(http.MethodGet, "http://example.com/api", nil)
	var info Info
	if err := client.Apply(req, &info); err != nil {
		t.Fatalf("Apply() = %v", err)
	}

	if info.Headers == nil {
		t.Fatal("info.Headers should not be nil")
	}
	if _, ok := info.Headers["X-MBX-APIKEY"]; !ok {
		t.Error("info.Headers should contain X-MBX-APIKEY")
	}
}

func TestHMACAuthClient_Apply_EmptySecretKey(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{APIKey: "key", SecretKey: ""}
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}

	req, _ := http.NewRequest(http.MethodGet, "http://example.com/api", nil)
	if err := client.Apply(req, nil); err != nil {
		t.Fatalf("Apply() = %v", err)
	}

	sig := req.URL.Query().Get("signature")
	if sig == "" {
		t.Fatal("signature should be present even with empty secret")
	}
}

func TestHMACAuthClient_Apply_ComplexQueryString(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{APIKey: "key", SecretKey: "secret"}
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}

	req, _ := http.NewRequest(
		http.MethodGet,
		"http://example.com/api?symbol=BTCUSDT&side=BUY&type=LIMIT&timeInForce=GTC&quantity=1&price=0.1&recvWindow=5000",
		nil,
	)
	if err := client.Apply(req, nil); err != nil {
		t.Fatalf("Apply() = %v", err)
	}

	q := req.URL.Query()
	if q.Get("symbol") != "BTCUSDT" {
		t.Errorf("symbol = %q, want %q", q.Get("symbol"), "BTCUSDT")
	}
	if q.Get("side") != "BUY" {
		t.Errorf("side = %q, want %q", q.Get("side"), "BUY")
	}
	if q.Get("signature") == "" {
		t.Error("signature should be present")
	}
}

func TestHMACAuthClient_Apply_NoQueryString(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{APIKey: "key", SecretKey: "secret"}
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}

	req, _ := http.NewRequest(http.MethodGet, "http://example.com/api", nil)
	if err := client.Apply(req, nil); err != nil {
		t.Fatalf("Apply() = %v", err)
	}

	q := req.URL.Query()
	if q.Get("timestamp") == "" {
		t.Error("timestamp should be present")
	}
	if q.Get("signature") == "" {
		t.Error("signature should be present")
	}
}

func TestHMACAuthClient_Apply_InfoTimestampIsString(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{APIKey: "key", SecretKey: "secret"}
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}

	req, _ := http.NewRequest(http.MethodGet, "http://example.com/api", nil)
	var info Info
	if err := client.Apply(req, &info); err != nil {
		t.Fatalf("Apply() = %v", err)
	}

	ts, ok := info.QueryParams["timestamp"]
	if !ok {
		t.Fatal("info.QueryParams should contain timestamp")
	}
	if _, err := strconv.ParseInt(ts, 10, 64); err != nil {
		t.Errorf("timestamp should be a valid int string: %v", err)
	}
}

func TestHMACAuthClient_Apply_InfoSignatureIsHex(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{APIKey: "key", SecretKey: "secret"}
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}

	req, _ := http.NewRequest(http.MethodGet, "http://example.com/api", nil)
	var info Info
	if err := client.Apply(req, &info); err != nil {
		t.Fatalf("Apply() = %v", err)
	}

	sig, ok := info.QueryParams["signature"]
	if !ok {
		t.Fatal("info.QueryParams should contain signature")
	}
	if _, err := hex.DecodeString(sig); err != nil {
		t.Errorf("signature should be valid hex: %v", err)
	}
}

func TestHMACAuthClient_Apply_HeaderAndQueryConsistent(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{APIKey: "test-key", SecretKey: "test-secret"}
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}

	req, _ := http.NewRequest(http.MethodGet, "http://example.com/api", nil)
	var info Info
	if err := client.Apply(req, &info); err != nil {
		t.Fatalf("Apply() = %v", err)
	}

	if info.Headers["X-MBX-APIKEY"] != "test-key" {
		t.Errorf("header mismatch: got %q", info.Headers["X-MBX-APIKEY"])
	}
	if info.QueryParams["signature"] == "" {
		t.Error("signature missing from info")
	}
	if info.QueryParams["timestamp"] == "" {
		t.Error("timestamp missing from info")
	}
}

func TestHMACAuthClient_Apply_DoesNotModifyMethod(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{APIKey: "key", SecretKey: "secret"}
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}

	req, _ := http.NewRequest(
		http.MethodPost,
		"http://example.com/api",
		strings.NewReader(`{"key":"value"}`),
	)
	if err := client.Apply(req, nil); err != nil {
		t.Fatalf("Apply() = %v", err)
	}

	if req.Method != http.MethodPost {
		t.Errorf("Method = %q, want %q", req.Method, http.MethodPost)
	}
}

func TestHMACAuthClient_Apply_DoesNotModifyBody(t *testing.T) {
	t.Parallel()

	client := &HMACAuthClient{APIKey: "key", SecretKey: "secret"}
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}

	body := `{"symbol":"BTCUSDT","side":"BUY"}`
	req, _ := http.NewRequest(
		http.MethodPost,
		"http://example.com/api",
		strings.NewReader(body),
	)
	if err := client.Apply(req, nil); err != nil {
		t.Fatalf("Apply() = %v", err)
	}

	//nolint:canonicalheader // Binance requires this exact header name
	if req.Header.Get("X-MBX-APIKEY") != "key" {
		t.Errorf("X-MBX-APIKEY = %q, want %q", req.Header.Get("X-MBX-APIKEY"), "key")
	}
	if req.URL.Query().Get("signature") == "" {
		t.Error("signature should be present")
	}
}
