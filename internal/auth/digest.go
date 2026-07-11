package auth

import (
	"context"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
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

func (c *DigestAuthClient) New() error {
	c.Username = resolveEnv(c.Username)
	c.Password = resolveEnv(c.Password)
	return nil
}

func (c *DigestAuthClient) Type() Type {
	return DigestAuth
}

func (c *DigestAuthClient) Apply(req *http.Request, out *Info) error {
	c.mu.Lock()
	challenge := digestChallenge{
		realm:     c.realm,
		nonce:     c.nonce,
		opaque:    c.opaque,
		qop:       c.qop,
		algorithm: c.algorithm,
	}
	nc := c.nonceCount
	cnonce := c.cnonce
	hasChallenge := c.realm != "" && c.nonce != "" && time.Since(c.cachedAt) < digestNonceTTL
	c.mu.Unlock()

	if !hasChallenge {
		newChallenge, err := c.fetchChallenge(req)
		if err != nil {
			return fmt.Errorf("digest: %w", err)
		}
		challenge = newChallenge
		nc = 0
		cnonce = c.generateCnonce()

		c.mu.Lock()
		c.realm = challenge.realm
		c.nonce = challenge.nonce
		c.opaque = challenge.opaque
		c.qop = challenge.qop
		c.algorithm = challenge.algorithm
		c.nonceCount = 0
		c.cnonce = cnonce
		c.cachedAt = time.Now()
		c.mu.Unlock()
	}

	nc++
	auth := c.buildDigest(req.Method, req.URL.RequestURI(), challenge, nc, cnonce)
	setAuthHeader(req, out, "Authorization", auth)

	c.mu.Lock()
	c.nonceCount = nc
	c.mu.Unlock()

	return nil
}

func (c *DigestAuthClient) SetMockBaseURL(url string) {
	c.MockBaseURL = url
}

func (c *DigestAuthClient) fetchChallenge(req *http.Request) (digestChallenge, error) {
	challengeURL := req.URL.String()
	if c.MockBaseURL != "" {
		challengeURL = c.MockBaseURL
	}

	fakeReq, reqErr := http.NewRequestWithContext(context.Background(), req.Method, challengeURL, nil)
	if reqErr != nil {
		return digestChallenge{}, fmt.Errorf("create challenge request: %w", reqErr)
	}

	cli, cliErr := httpclient.NewDefault()
	if cliErr != nil {
		return digestChallenge{}, fmt.Errorf("create http client: %w", cliErr)
	}
	resp, doErr := cli.Do(fakeReq)
	if doErr != nil {
		return digestChallenge{}, fmt.Errorf("challenge request: %w", doErr)
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
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func md5hex(s string) string {
	h := md5.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}

func (c *DigestAuthClient) Validate() error {
	return authValidator.Struct(c)
}
