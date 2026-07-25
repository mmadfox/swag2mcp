package auth

import (
	"context"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/mmadfox/swag2mcp/internal/httpclient"
)

const (
	digestNonceTTL   = 5 * time.Minute
	digestNonceBytes = 8
	digestAlgoMD5    = "MD5"
	digestQopAuth    = "auth"
	digestQopAuthInt = "auth-int"
)

// DigestAuthClient holds credentials for HTTP Digest authentication.
type DigestAuthClient struct {
	Username string `yaml:"username" validate:"required"`
	Password string `yaml:"password" validate:"required"`

	MockBaseURL string `yaml:"-"`

	mu         sync.Mutex
	realm      string
	nonce      string
	opaque     string
	qop        string
	algorithm  string
	nonceCount int
	cnonce     string
	cachedAt   time.Time
}

type digestChallenge struct {
	realm     string
	nonce     string
	opaque    string
	qop       string
	algorithm string
}

// New resolves environment variables in Username and Password and returns nil.
func (c *DigestAuthClient) New() error {
	c.Username = resolveEnv(c.Username)
	c.Password = resolveEnv(c.Password)
	return nil
}

// Type returns the authentication type for HTTP Digest auth.
func (c *DigestAuthClient) Type() Type {
	return DigestAuth
}

// Apply performs Digest authentication by fetching a challenge, computing the response, and setting the Authorization header.
func (c *DigestAuthClient) Apply(req *http.Request, out *Info) error {
	challenge, nc, cnonce, ok := c.readChallenge()
	if !ok {
		var err error
		challenge, err = c.fetchChallenge(req)
		if err != nil {
			return fmt.Errorf("digest: %w", err)
		}
		nc = 0
		cnonce = c.generateCnonce()
		c.writeChallenge(challenge, cnonce)
	}

	nc++
	auth := c.buildDigest(req.Method, req.URL.RequestURI(), challenge, nc, cnonce)
	setAuthHeader(req, out, headerAuthorization, auth)

	c.updateNonceCount(nc)
	return nil
}

func (c *DigestAuthClient) readChallenge() (digestChallenge, int, string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	ok := c.realm != "" && c.nonce != "" && time.Since(c.cachedAt) < digestNonceTTL
	return digestChallenge{
		realm: c.realm, nonce: c.nonce, opaque: c.opaque,
		qop: c.qop, algorithm: c.algorithm,
	}, c.nonceCount, c.cnonce, ok
}

func (c *DigestAuthClient) writeChallenge(ch digestChallenge, cnonce string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.realm = ch.realm
	c.nonce = ch.nonce
	c.opaque = ch.opaque
	c.qop = ch.qop
	c.algorithm = ch.algorithm
	c.nonceCount = 0
	c.cnonce = cnonce
	c.cachedAt = time.Now()
}

func (c *DigestAuthClient) updateNonceCount(nc int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.nonceCount = nc
}

// SetMockBaseURL sets a mock base URL used for fetching the Digest challenge during testing.
func (c *DigestAuthClient) SetMockBaseURL(url string) {
	c.MockBaseURL = url
}

func (c *DigestAuthClient) fetchChallenge(req *http.Request) (digestChallenge, error) {
	challengeURL := req.URL.String()
	if c.MockBaseURL != "" {
		challengeURL = c.MockBaseURL
	}

	fakeReq, err := http.NewRequestWithContext(context.Background(), req.Method, challengeURL, nil)
	if err != nil {
		return digestChallenge{}, fmt.Errorf("create challenge request: %w", err)
	}

	cli, err := httpclient.NewDefault()
	if err != nil {
		return digestChallenge{}, fmt.Errorf("create http client: %w", err)
	}
	resp, err := cli.Do(fakeReq)
	if err != nil {
		return digestChallenge{}, fmt.Errorf("challenge request: %w", err)
	}
	_ = resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		return digestChallenge{}, fmt.Errorf("expected 401 for digest challenge, got %d", resp.StatusCode)
	}

	authHeader := resp.Header.Get("WWW-Authenticate")
	if !strings.HasPrefix(authHeader, "Digest ") {
		return digestChallenge{}, fmt.Errorf("expected Digest challenge, got: %s", authHeader)
	}

	return parseDigestChallenge(authHeader), nil
}

func parseDigestChallenge(header string) digestChallenge {
	c := digestChallenge{algorithm: digestAlgoMD5}
	rest := header[len("Digest "):]

	for part := range strings.SplitSeq(rest, ",") {
		part = strings.TrimSpace(part)
		key, val, found := strings.Cut(part, "=")
		if !found {
			continue
		}
		key = strings.TrimSpace(key)
		val = strings.Trim(strings.TrimSpace(val), "\"")
		switch strings.ToLower(key) {
		case "realm":
			c.realm = val
		case "nonce":
			c.nonce = val
		case "opaque":
			c.opaque = val
		case "qop":
			c.qop = val
		case "algorithm":
			c.algorithm = val
		}
	}
	return c
}

func (c *DigestAuthClient) buildDigest(method, uri string, ch digestChallenge, nc int, cnonce string) string {
	ha1Input := fmt.Sprintf("%s:%s:%s", c.Username, ch.realm, c.Password)
	ha1 := md5hex(ha1Input)

	ha2Input := fmt.Sprintf("%s:%s", method, uri)
	ha2 := md5hex(ha2Input)

	ncStr := fmt.Sprintf("%08x", nc)

	var response string
	if ch.qop == digestQopAuth || ch.qop == digestQopAuthInt {
		respInput := fmt.Sprintf("%s:%s:%s:%s:%s:%s", ha1, ch.nonce, ncStr, cnonce, ch.qop, ha2)
		response = md5hex(respInput)
	} else {
		respInput := fmt.Sprintf("%s:%s:%s", ha1, ch.nonce, ha2)
		response = md5hex(respInput)
	}

	var b strings.Builder
	fmt.Fprintf(&b, `Digest username="%s", realm="%s", nonce="%s", uri="%s", response="%s"`,
		c.Username, ch.realm, ch.nonce, uri, response)

	if ch.algorithm != "" {
		fmt.Fprintf(&b, `, algorithm="%s"`, ch.algorithm)
	}
	if ch.opaque != "" {
		fmt.Fprintf(&b, `, opaque="%s"`, ch.opaque)
	}
	if ch.qop != "" {
		fmt.Fprintf(&b, `, qop=%s, nc=%s, cnonce="%s"`, ch.qop, ncStr, cnonce)
	}

	return b.String()
}

func (c *DigestAuthClient) generateCnonce() string {
	b := make([]byte, digestNonceBytes)
	if _, err := rand.Read(b); err != nil {
		slog.Default().Warn("failed to generate cnonce", "error", err)
	}
	return hex.EncodeToString(b)
}

func md5hex(s string) string {
	h := md5.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}

// Validate checks that the Username and Password fields are present and valid.
func (c *DigestAuthClient) Validate() error {
	return authValidator.Struct(c)
}
