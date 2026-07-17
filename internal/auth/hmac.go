package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"
)

const (
	hmacTimestampKey = "timestamp"
	hmacSignatureKey = "signature"
	hmacAPIKeyHeader = "X-MBX-APIKEY" //nolint:gosec // This is a header name, not a credential.
)

// HMACAuthClient holds credentials for HMAC-SHA256 signature authentication
// (Binance-style). It signs requests by computing an HMAC-SHA256 of the
// query string (including a timestamp) and appends the signature as a
// query parameter. The API key is sent via the X-MBX-APIKEY header.
type HMACAuthClient struct {
	APIKey    string `yaml:"api_key"    validate:"required"`
	SecretKey string `yaml:"secret_key" validate:"required"`
}

// New resolves environment variable references in the API key and secret key.
func (c *HMACAuthClient) New() error {
	c.APIKey = resolveEnv(c.APIKey)
	c.SecretKey = resolveEnv(c.SecretKey)
	return nil
}

// Type returns the authentication method type (hmac).
func (c *HMACAuthClient) Type() Type {
	return HMACAuth
}

// Apply signs the request by setting the X-MBX-APIKEY header, adding a
// millisecond-precision timestamp to the query string, and appending an
// HMAC-SHA256 signature of the sorted query parameters.
func (c *HMACAuthClient) Apply(req *http.Request, out *Info) error {
	setAuthHeader(req, out, hmacAPIKeyHeader, c.APIKey)

	q := req.URL.Query()
	q.Set(hmacTimestampKey, strconv.FormatInt(time.Now().UnixMilli(), 10))

	queryString := q.Encode()
	mac := hmac.New(sha256.New, []byte(c.SecretKey))
	mac.Write([]byte(queryString))
	signature := hex.EncodeToString(mac.Sum(nil))

	q.Set(hmacSignatureKey, signature)
	req.URL.RawQuery = q.Encode()

	if out == nil {
		return nil
	}
	if out.QueryParams == nil {
		out.QueryParams = make(map[string]string)
	}
	out.QueryParams[hmacSignatureKey] = signature
	out.QueryParams[hmacTimestampKey] = q.Get(hmacTimestampKey)

	return nil
}

// Validate checks that the API key and secret key are not empty.
func (c *HMACAuthClient) Validate() error {
	return authValidator.Struct(c)
}
