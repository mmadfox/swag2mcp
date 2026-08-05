/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package mockserver

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/MicahParks/jwkset"
	"github.com/golang-jwt/jwt/v5"
)

const (
	authNonceBytes      = 8
	authNonceTTL        = 5 * time.Minute
	authTokenLength     = 32
	authShutdownTimeout = 5
	authTokenExpiresIn  = 3600
	authExpiresInKey    = "expires_in"
	grantTypePassword   = "password"
	// jwtKeySize is the RSA key size in bits for the JWT mock server.
	// Used by newAuthMockServer when generating the signing key.
	jwtKeySize = 2048
	// jwtTokenExpiry is the token lifetime in seconds for the JWT mock server.
	// Used by handleJWTToken when building the token response.
	jwtTokenExpiry = 3600
	// accessTokenKey is the JSON key for the access token in OAuth2 and JWT
	// token responses. Used by handleOAuth2CC, handleOAuth2Password, and
	// handleJWTToken.
	accessTokenKey = "access_token"
	// tokenTypeKey is the JSON key for the token type in OAuth2 and JWT
	// token responses. Used by handleOAuth2CC, handleOAuth2Password, and
	// handleJWTToken.
	tokenTypeKey = "token_type"
	// tokenTypeBearer is the Bearer token type value used in OAuth2 and JWT
	// token responses. Used by handleOAuth2CC, handleOAuth2Password, and
	// handleJWTToken.
	tokenTypeBearer = "Bearer"
)

type authServerType string

const (
	authServerOAuth2 authServerType = "oauth2"
	authServerDigest authServerType = "digest"
	authServerHMAC   authServerType = "hmac"
	authServerJWT    authServerType = "jwt"
)

type authMockServer struct {
	serverType authServerType
	addr       string
	server     *http.Server
	logger     *slog.Logger
	tlsConfig  *tls.Config
	mu         sync.Mutex
	nonce      string
	opaque     string
	nonceTime  time.Time
	jwtKey     *rsa.PrivateKey
	jwkSet     jwkset.JWKSMarshal
}

// newAuthMockServer creates a new auth mock server of the given type.
func newAuthMockServer(
	serverType authServerType,
	addr string,
	tlsConfig *tls.Config,
	logger *slog.Logger,
) *authMockServer {
	if logger == nil {
		logger = slog.New(slog.DiscardHandler)
	}
	s := &authMockServer{
		serverType: serverType,
		addr:       addr,
		tlsConfig:  tlsConfig,
		logger:     logger,
	}
	if serverType == authServerJWT {
		key, err := rsa.GenerateKey(rand.Reader, jwtKeySize)
		if err != nil {
			logger.ErrorContext(context.Background(), "failed to generate JWT key", "error", err)
		} else {
			s.jwtKey = key
			s.jwkSet = jwksMarshalFromKey(key)
		}
	}
	return s
}

// start begins listening for HTTP requests on the configured address.
func (m *authMockServer) start(ctx context.Context) {
	mux := http.NewServeMux()

	switch m.serverType {
	case authServerOAuth2:
		mux.HandleFunc("/token", m.handleOAuth2)
	case authServerDigest:
		mux.HandleFunc("/", m.handleDigest)
	case authServerHMAC:
		mux.HandleFunc("/", m.handleHMAC)
	case authServerJWT:
		mux.HandleFunc("/token", m.handleJWTToken)
		mux.HandleFunc("/.well-known/jwks.json", m.handleJWKS)
	default:
		m.logger.ErrorContext(ctx, "unsupported auth server type",
			"type", m.serverType,
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
				"type", m.serverType,
				"addr", m.addr,
				"error", serveError,
			)
		}
	}()

	m.logger.InfoContext(ctx, "auth mock server started",
		"type", m.serverType,
		"addr", m.addr,
	)
}

// shutdown gracefully stops the auth mock server.
func (m *authMockServer) shutdown() {
	if m.server != nil {
		shutdownContext, shutdownCancel := context.WithTimeout(
			context.Background(),
			authShutdownTimeout*time.Second,
		)
		defer shutdownCancel()
		if err := m.server.Shutdown(shutdownContext); err != nil {
			m.logger.Warn("mock auth server shutdown error", "error", err)
		}
	}
}

// handleOAuth2 handles OAuth2 token requests (client_credentials and password grants).
func (m *authMockServer) handleOAuth2(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		m.logger.WarnContext(request.Context(), "oauth2 mock: method not allowed",
			"method", request.Method,
		)
		http.Error(responseWriter, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	params, parseError := parseOAuth2Body(request)
	if parseError != nil {
		m.logger.WarnContext(request.Context(), "oauth2 mock: failed to parse body",
			"error", parseError,
		)
		http.Error(responseWriter, "Bad Request", http.StatusBadRequest)
		return
	}

	grantType := params["grant_type"]

	m.logger.InfoContext(request.Context(), "oauth2 mock: token requested",
		"grant_type", grantType,
		"client_id", params["client_id"],
		"username", params["username"],
	)

	switch grantType {
	case "client_credentials":
		m.handleOAuth2CC(request.Context(), responseWriter, params)
	case grantTypePassword:
		m.handleOAuth2Password(request.Context(), responseWriter, params)
	default:
		m.logger.WarnContext(request.Context(), "oauth2 mock: invalid grant_type",
			"grant_type", grantType,
		)
		http.Error(responseWriter, "Invalid grant_type", http.StatusBadRequest)
	}
}

// parseOAuth2Body parses the request body as form-urlencoded or JSON based on Content-Type.
func parseOAuth2Body(r *http.Request) (map[string]string, error) {
	ct := r.Header.Get("Content-Type")
	if strings.HasPrefix(ct, "application/json") {
		var body map[string]string
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			return nil, err
		}
		return body, nil
	}
	if err := r.ParseForm(); err != nil {
		return nil, err
	}
	params := make(map[string]string, len(r.Form))
	for k, v := range r.Form {
		params[k] = v[0]
	}
	return params, nil
}

// handleOAuth2CC handles client_credentials grant type requests.
func (m *authMockServer) handleOAuth2CC(
	ctx context.Context,
	responseWriter http.ResponseWriter,
	params map[string]string,
) {
	clientID := params["client_id"]
	if clientID == "" {
		m.logger.WarnContext(ctx, "oauth2 mock: missing client_id")
		http.Error(responseWriter, "Missing client_id", http.StatusBadRequest)
		return
	}

	token := generateRandomToken()
	m.logger.InfoContext(ctx, "oauth2 mock: token issued",
		"grant_type", "client_credentials",
		"client_id", clientID,
		"token_prefix", token[:8],
	)

	tokenResponse := map[string]any{
		accessTokenKey:   token,
		tokenTypeKey:     tokenTypeBearer,
		authExpiresInKey: authTokenExpiresIn,
		"scope":          params["scope"],
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(responseWriter).Encode(tokenResponse); err != nil {
		m.logger.WarnContext(ctx, "failed to encode oauth2 cc token response", "error", err)
	}
}

// handleOAuth2Password handles password grant type requests.
func (m *authMockServer) handleOAuth2Password(
	ctx context.Context,
	responseWriter http.ResponseWriter,
	params map[string]string,
) {
	username := params["username"]
	password := params["password"]
	if username == "" || password == "" {
		m.logger.WarnContext(ctx, "oauth2 mock: missing username or password",
			"username", username,
		)
		http.Error(responseWriter, "Missing username or password", http.StatusBadRequest)
		return
	}

	token := generateRandomToken()
	m.logger.InfoContext(ctx, "oauth2 mock: token issued",
		"grant_type", "password",
		"username", username,
		"token_prefix", token[:8],
	)

	tokenResponse := map[string]any{
		accessTokenKey:   token,
		tokenTypeKey:     tokenTypeBearer,
		authExpiresInKey: authTokenExpiresIn,
		"scope":          params["scope"],
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(responseWriter).Encode(tokenResponse); err != nil {
		m.logger.WarnContext(ctx, "failed to encode oauth2 password token response", "error", err)
	}
}

// handleDigest handles Digest authentication requests, sending a challenge or validating credentials.
func (m *authMockServer) handleDigest(responseWriter http.ResponseWriter, request *http.Request) {
	authorization := request.Header.Get("Authorization")

	if !strings.HasPrefix(authorization, "Digest ") {
		m.logger.InfoContext(request.Context(), "digest mock: sending challenge")
		m.generateDigestChallenge(responseWriter)
		return
	}

	digestParams := m.parseDigestAuthorization(authorization)
	response := digestParams["response"]
	if response == "" {
		m.logger.InfoContext(request.Context(), "digest mock: empty response, sending challenge")
		m.generateDigestChallenge(responseWriter)
		return
	}

	m.logger.InfoContext(request.Context(), "digest mock: authentication successful",
		"username", digestParams["username"],
	)

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	if _, err := responseWriter.Write([]byte(`{"status":"authenticated","method":"digest"}`)); err != nil {
		m.logger.WarnContext(request.Context(), "failed to write digest response", "error", err)
	}
}

// handleHMAC handles HMAC authentication requests, validating API key, signature, and timestamp.
func (m *authMockServer) handleHMAC(responseWriter http.ResponseWriter, request *http.Request) {
	apiKey := request.Header.Get("X-MBX-APIKEY")
	if apiKey == "" {
		m.logger.WarnContext(request.Context(), "hmac mock: missing X-MBX-APIKEY header")
		http.Error(responseWriter, "Missing X-MBX-APIKEY header", http.StatusUnauthorized)
		return
	}

	signature := request.URL.Query().Get("signature")
	if signature == "" {
		m.logger.WarnContext(request.Context(), "hmac mock: missing signature query param")
		http.Error(responseWriter, "Missing signature query param", http.StatusUnauthorized)
		return
	}

	timestamp := request.URL.Query().Get("timestamp")
	if timestamp == "" {
		m.logger.WarnContext(request.Context(), "hmac mock: missing timestamp query param")
		http.Error(responseWriter, "Missing timestamp query param", http.StatusUnauthorized)
		return
	}

	m.logger.InfoContext(request.Context(), "hmac mock: authentication successful",
		"api_key_prefix", apiKey[:min(len(apiKey), authNonceBytes)],
	)

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	if _, err := responseWriter.Write([]byte(`{"status":"authenticated","method":"hmac"}`)); err != nil {
		m.logger.WarnContext(request.Context(), "failed to write hmac response", "error", err)
	}
}

// generateDigestChallenge sends a WWW-Authenticate Digest challenge to the client.
func (m *authMockServer) generateDigestChallenge(responseWriter http.ResponseWriter) {
	m.mu.Lock()
	now := time.Now()
	if m.nonce == "" || now.Sub(m.nonceTime) > authNonceTTL {
		nonceBytes := make([]byte, authNonceBytes)
		if _, err := rand.Read(nonceBytes); err != nil {
			m.logger.Warn("failed to generate nonce", "error", err)
		}
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

	m.logger.Info("digest mock: challenge sent",
		"nonce_prefix", nonce[:8],
	)

	responseWriter.Header().Set("WWW-Authenticate", challenge)
	http.Error(responseWriter, "Unauthorized", http.StatusUnauthorized)
}

// parseDigestAuthorization parses a Digest authorization header into key-value pairs.
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

func generateRandomToken() string {
	tokenBytes := make([]byte, authTokenLength)
	if _, err := rand.Read(tokenBytes); err != nil {
		slog.Default().Warn("failed to generate random token", "error", err)
	}
	return hex.EncodeToString(tokenBytes)
}

// handleJWTToken generates a signed JWT access token.
func (m *authMockServer) handleJWTToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if m.jwtKey == nil {
		http.Error(w, "JWT key not initialized", http.StatusInternalServerError)
		return
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"iss": fmt.Sprintf("http://%s", m.addr),
		"sub": "mock-user",
		"aud": []string{"swag2mcp"},
		"exp": now.Add(time.Hour).Unix(),
		"iat": now.Unix(),
		"jti": generateRandomToken()[:16],
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = "mock-kid"

	signed, err := token.SignedString(m.jwtKey)
	if err != nil {
		m.logger.ErrorContext(r.Context(), "failed to sign JWT", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]any{
		accessTokenKey: signed,
		tokenTypeKey:   tokenTypeBearer,
		"expires_in":   jwtTokenExpiry,
	}); err != nil {
		m.logger.ErrorContext(r.Context(), "failed to encode JWT token response", "error", err)
	}
}

// handleJWKS serves the JWKS public key set.
func (m *authMockServer) handleJWKS(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(m.jwkSet); err != nil {
		slog.Default().WarnContext(context.Background(), "failed to encode JWKS response", "error", err)
	}
}

// jwksMarshalFromKey creates a JWKSMarshal from an RSA public key.
func jwksMarshalFromKey(key *rsa.PrivateKey) jwkset.JWKSMarshal {
	jwk, err := jwkset.NewJWKFromKey(&key.PublicKey, jwkset.JWKOptions{
		Metadata: jwkset.JWKMetadataOptions{
			KID: "mock-kid",
			ALG: jwkset.AlgRS256,
		},
	})
	if err != nil {
		slog.Default().Error("failed to create JWK", "error", err)
		return jwkset.JWKSMarshal{}
	}

	storage := jwkset.NewMemoryStorage()
	if err := storage.KeyWrite(context.Background(), jwk); err != nil {
		slog.Default().Error("failed to write JWK to storage", "error", err)
		return jwkset.JWKSMarshal{}
	}

	marshal, err := storage.Marshal(context.Background())
	if err != nil {
		slog.Default().Error("failed to marshal JWK set", "error", err)
		return jwkset.JWKSMarshal{}
	}

	return marshal
}
