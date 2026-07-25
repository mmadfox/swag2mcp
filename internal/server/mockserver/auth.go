package mockserver

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	authNonceBytes      = 8
	authNonceTTL        = 5 * time.Minute
	authTokenLength     = 32
	authShutdownTimeout = 5
	authTokenExpiresIn  = 3600
	authExpiresInKey    = "expires_in"
	grantTypePassword   = "password"
)

// authMockServer simulates an authentication server for a specific auth type.
// It runs on its own address and responds to requests as a real auth provider would.
type authMockServer struct {
	specDomain string
	authType   string
	addr       string
	server     *http.Server
	logger     *slog.Logger
	tlsConfig  *tls.Config
	mu         sync.Mutex
	nonce      string
	opaque     string
	nonceTime  time.Time
}

// newAuthMockServer creates a new auth mock server for the given domain and auth type.
func newAuthMockServer(
	domain string,
	authType string,
	addr string,
	tlsConfig *tls.Config,
	logger *slog.Logger,
) *authMockServer {
	return &authMockServer{
		specDomain: domain,
		authType:   authType,
		addr:       addr,
		tlsConfig:  tlsConfig,
		logger:     logger,
	}
}

// start begins listening for HTTP requests on the configured port.
// It registers the appropriate handler based on the auth type.
func (m *authMockServer) start(ctx context.Context) {
	mux := http.NewServeMux()

	switch m.authType {
	case "basic":
		mux.HandleFunc("/", m.handleBasic)
	case "bearer":
		mux.HandleFunc("/", m.handleBearer)
	case "digest":
		mux.HandleFunc("/", m.handleDigest)
	case "oauth2-cc":
		mux.HandleFunc("/token", m.handleOAuth2CC)
	case "oauth2-pwd":
		mux.HandleFunc("/token", m.handleOAuth2Password)
	case "api-key":
		mux.HandleFunc("/", m.handleAPIKey)
	case "script":
		mux.HandleFunc("/token", m.handleScript)
	default:
		m.logger.ErrorContext(ctx, "unsupported auth type",
			"type", m.authType,
		)
		return
	}

	address := extractHostPort(m.addr)
	if !strings.Contains(address, ":") {
		address = ":" + address
	}
	m.server = &http.Server{
		Addr:              address,
		Handler:           mux,
		ReadHeaderTimeout: authShutdownTimeout * time.Second,
	}

	if m.tlsConfig != nil {
		m.server.TLSConfig = m.tlsConfig
	}

	go func() {
		serveError := m.server.ListenAndServe()
		if serveError != nil && serveError != http.ErrServerClosed {
			m.logger.ErrorContext(ctx, "auth mock server error",
				"type", m.authType,
				"addr", m.addr,
				"error", serveError,
			)
		}
	}()

	m.logger.InfoContext(ctx, "auth mock server started",
		"type", m.authType,
		"addr", m.addr,
	)
}

// shutdown gracefully stops the HTTP server.
func (m *authMockServer) shutdown() {
	if m.server != nil {
		shutdownContext, shutdownCancel := context.WithTimeout(
			context.Background(),
			authShutdownTimeout*time.Second,
		)
		defer shutdownCancel()
		_ = m.server.Shutdown(shutdownContext)
	}
}

// handleBasic validates HTTP Basic authentication.
// It returns 401 with a WWW-Authenticate header if no credentials are provided,
// or 200 with a success response if valid Basic credentials are present.
func (m *authMockServer) handleBasic(responseWriter http.ResponseWriter, request *http.Request) {
	authorization := request.Header.Get("Authorization")

	if !strings.HasPrefix(authorization, "Basic ") {
		responseWriter.Header().Set("WWW-Authenticate", "Basic realm=\"swag2mcp-mock\"")
		http.Error(responseWriter, "Unauthorized", http.StatusUnauthorized)
		return
	}

	payload := strings.TrimPrefix(authorization, "Basic ")
	decodedBytes, decodeError := base64.StdEncoding.DecodeString(payload)
	if decodeError != nil {
		http.Error(responseWriter, "Bad Request", http.StatusBadRequest)
		return
	}

	credentials := string(decodedBytes)
	_, _, hasColon := strings.Cut(credentials, ":")
	if !hasColon {
		http.Error(responseWriter, "Bad Request", http.StatusBadRequest)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	_, _ = responseWriter.Write([]byte(`{"status":"authenticated","method":"basic"}`))
}

// handleBearer validates HTTP Bearer token authentication.
// It returns 401 if no Bearer token is provided, or 200 with a success response.
func (m *authMockServer) handleBearer(responseWriter http.ResponseWriter, request *http.Request) {
	authorization := request.Header.Get("Authorization")

	if !strings.HasPrefix(authorization, "Bearer ") {
		http.Error(responseWriter, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token := strings.TrimPrefix(authorization, "Bearer ")
	if token == "" {
		http.Error(responseWriter, "Unauthorized", http.StatusUnauthorized)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	_, _ = responseWriter.Write([]byte(`{"status":"authenticated","method":"bearer"}`))
}

// handleDigest validates HTTP Digest authentication.
// It returns 401 with a Digest challenge on first request, then validates
// the digest response on subsequent requests.
func (m *authMockServer) handleDigest(responseWriter http.ResponseWriter, request *http.Request) {
	authorization := request.Header.Get("Authorization")

	if !strings.HasPrefix(authorization, "Digest ") {
		m.generateDigestChallenge(responseWriter)
		return
	}

	digestParams := m.parseDigestAuthorization(authorization)
	response := digestParams["response"]
	if response == "" {
		m.generateDigestChallenge(responseWriter)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	_, _ = responseWriter.Write([]byte(`{"status":"authenticated","method":"digest"}`))
}

// generateDigestChallenge sends a 401 response with a Digest WWW-Authenticate header.
// The nonce and opaque values are cached and rotated after authNonceTTL.
func (m *authMockServer) generateDigestChallenge(responseWriter http.ResponseWriter) {
	m.mu.Lock()
	now := time.Now()
	if m.nonce == "" || now.Sub(m.nonceTime) > authNonceTTL {
		nonceBytes := make([]byte, authNonceBytes)
		_, _ = rand.Read(nonceBytes)
		m.nonce = hex.EncodeToString(nonceBytes)
		m.opaque = hex.EncodeToString(nonceBytes)
		m.nonceTime = now
	}
	nonce := m.nonce
	opaque := m.opaque
	m.mu.Unlock()

	challenge := fmt.Sprintf(
		"Digest realm=\"swag2mcp-mock\", nonce=\"%s\", opaque=\"%s\", algorithm=MD5, qop=\"auth\"",
		nonce, opaque,
	)

	responseWriter.Header().Set("WWW-Authenticate", challenge)
	http.Error(responseWriter, "Unauthorized", http.StatusUnauthorized)
}

// parseDigestAuthorization extracts Digest authentication parameters from the
// Authorization header value.
func (m *authMockServer) parseDigestAuthorization(authorization string) map[string]string {
	parameters := make(map[string]string)
	headerValue := strings.TrimPrefix(authorization, "Digest ")

	for part := range strings.SplitSeq(headerValue, ",") {
		part = strings.TrimSpace(part)
		key, value, found := strings.Cut(part, "=")
		if !found {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.Trim(strings.TrimSpace(value), "\"")
		parameters[key] = value
	}

	return parameters
}

// handleOAuth2CC simulates the OAuth2 Client Credentials grant flow.
// It accepts POST requests to /token with grant_type=client_credentials
// and returns an access token.
func (m *authMockServer) handleOAuth2CC(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(responseWriter, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	parseError := request.ParseForm()
	if parseError != nil {
		http.Error(responseWriter, "Bad Request", http.StatusBadRequest)
		return
	}

	grantType := request.FormValue("grant_type")
	if grantType != "client_credentials" {
		http.Error(responseWriter, "Invalid grant_type", http.StatusBadRequest)
		return
	}

	clientID := request.FormValue("client_id")
	if clientID == "" {
		http.Error(responseWriter, "Missing client_id", http.StatusBadRequest)
		return
	}

	tokenResponse := map[string]any{
		"access_token":   generateRandomToken(),
		"token_type":     "Bearer",
		authExpiresInKey: authTokenExpiresIn,
		"scope":          request.FormValue("scope"),
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(responseWriter).Encode(tokenResponse)
}

// handleOAuth2Password simulates the OAuth2 Resource Owner Password grant flow.
// It accepts POST requests to /token with grant_type=password and returns
// an access token.
func (m *authMockServer) handleOAuth2Password(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(responseWriter, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	parseError := request.ParseForm()
	if parseError != nil {
		http.Error(responseWriter, "Bad Request", http.StatusBadRequest)
		return
	}

	grantType := request.FormValue("grant_type")
	if grantType != grantTypePassword {
		http.Error(responseWriter, "Invalid grant_type", http.StatusBadRequest)
		return
	}

	username := request.FormValue("username")
	password := request.FormValue("password")
	if username == "" || password == "" {
		http.Error(responseWriter, "Missing username or password", http.StatusBadRequest)
		return
	}

	tokenResponse := map[string]any{
		"access_token":   generateRandomToken(),
		"token_type":     "Bearer",
		authExpiresInKey: authTokenExpiresIn,
		"scope":          request.FormValue("scope"),
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(responseWriter).Encode(tokenResponse)
}

// handleAPIKey validates API key authentication.
// It checks for the key in the X-Api-Key header or the api_key query parameter.
func (m *authMockServer) handleAPIKey(responseWriter http.ResponseWriter, request *http.Request) {
	apiKeyHeader := request.Header.Get("X-Api-Key")
	apiKeyQuery := request.URL.Query().Get("api_key")

	if apiKeyHeader == "" && apiKeyQuery == "" {
		http.Error(responseWriter, "Unauthorized", http.StatusUnauthorized)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	_, _ = responseWriter.Write([]byte(`{"status":"authenticated","method":"api-key"}`))
}

// handleScript simulates a script-based authentication endpoint.
// It returns a token response as if a shell script was executed.
func (m *authMockServer) handleScript(responseWriter http.ResponseWriter, _ *http.Request) {
	tokenResponse := map[string]any{
		"token":          generateRandomToken(),
		authExpiresInKey: authTokenExpiresIn,
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(responseWriter).Encode(tokenResponse)
}

// generateRandomToken returns a random hex-encoded token string.
func generateRandomToken() string {
	tokenBytes := make([]byte, authTokenLength)
	_, _ = rand.Read(tokenBytes)
	return hex.EncodeToString(tokenBytes)
}
