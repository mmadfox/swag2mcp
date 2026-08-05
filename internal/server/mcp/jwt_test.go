/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package mcp

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
	"time"

	"github.com/MicahParks/jwkset"
	"github.com/golang-jwt/jwt/v5"
	"github.com/modelcontextprotocol/go-sdk/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewJWTVerifier_UnsupportedType(t *testing.T) {
	t.Parallel()

	verifier := newJWTVerifier(&JWTConfig{Type: "unknown"})
	_, err := verifier(context.Background(), "token", nil)
	assert.ErrorContains(t, err, "unsupported JWT auth type")
}

func TestNewJWTVerifier_JWKS_InvalidURL(t *testing.T) {
	t.Parallel()

	verifier := newJWTVerifier(&JWTConfig{
		Type:    "jwks",
		JWKSURL: "http://127.0.0.1:1/nonexistent",
	})
	_, err := verifier(context.Background(), "token", nil)
	assert.Error(t, err)
}

func TestNewJWTVerifier_JWKS_Valid(t *testing.T) {
	t.Parallel()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	jwksServer := newJWKSServer(t, &key.PublicKey)
	defer jwksServer.Close()

	verifier := newJWTVerifier(&JWTConfig{
		Type:     "jwks",
		JWKSURL:  jwksServer.URL + "/.well-known/jwks.json",
		Issuer:   "test-issuer",
		Audience: "test-audience",
	})

	token := signJWT(t, key, jwt.MapClaims{
		"iss": "test-issuer",
		"aud": []string{"test-audience"},
		"sub": "user-123",
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	info, err := verifier(context.Background(), token, nil)
	require.NoError(t, err)
	assert.Equal(t, "user-123", info.UserID)
	assert.False(t, info.Expiration.IsZero())
}

func TestNewJWTVerifier_JWKS_InvalidIssuer(t *testing.T) {
	t.Parallel()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	jwksServer := newJWKSServer(t, &key.PublicKey)
	defer jwksServer.Close()

	verifier := newJWTVerifier(&JWTConfig{
		Type:    "jwks",
		JWKSURL: jwksServer.URL + "/.well-known/jwks.json",
		Issuer:  "expected-issuer",
	})

	token := signJWT(t, key, jwt.MapClaims{
		"iss": "wrong-issuer",
		"sub": "user-123",
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	_, err = verifier(context.Background(), token, nil)
	assert.ErrorIs(t, err, auth.ErrInvalidToken)
}

func TestNewJWTVerifier_JWKS_InvalidAudience(t *testing.T) {
	t.Parallel()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	jwksServer := newJWKSServer(t, &key.PublicKey)
	defer jwksServer.Close()

	verifier := newJWTVerifier(&JWTConfig{
		Type:     "jwks",
		JWKSURL:  jwksServer.URL + "/.well-known/jwks.json",
		Audience: "expected-audience",
	})

	token := signJWT(t, key, jwt.MapClaims{
		"iss": "test-issuer",
		"aud": []string{"wrong-audience"},
		"sub": "user-123",
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	_, err = verifier(context.Background(), token, nil)
	assert.ErrorIs(t, err, auth.ErrInvalidToken)
}

func TestNewJWTVerifier_JWKS_ExpiredToken(t *testing.T) {
	t.Parallel()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	jwksServer := newJWKSServer(t, &key.PublicKey)
	defer jwksServer.Close()

	verifier := newJWTVerifier(&JWTConfig{
		Type:    "jwks",
		JWKSURL: jwksServer.URL + "/.well-known/jwks.json",
	})

	token := signJWT(t, key, jwt.MapClaims{
		"sub": "user-123",
		"exp": time.Now().Add(-time.Hour).Unix(),
	})

	_, err = verifier(context.Background(), token, nil)
	assert.Error(t, err)
}

func TestNewJWTVerifier_OIDC(t *testing.T) {
	t.Parallel()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	var serverURL string
	mux := http.NewServeMux()
	mux.HandleFunc("/.well-known/openid-configuration", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"issuer":   serverURL,
			"jwks_uri": serverURL + "/.well-known/jwks.json",
		})
	})
	mux.HandleFunc("/.well-known/jwks.json", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jwksMarshalFromKey(t, &key.PublicKey))
	})

	oidcServer := httptest.NewServer(mux)
	defer oidcServer.Close()
	serverURL = oidcServer.URL

	verifier := newJWTVerifier(&JWTConfig{
		Type:   "oidc",
		Issuer: oidcServer.URL,
	})

	token := signJWT(t, key, jwt.MapClaims{
		"iss": oidcServer.URL,
		"sub": "user-oidc",
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	info, err := verifier(context.Background(), token, nil)
	require.NoError(t, err)
	assert.Equal(t, "user-oidc", info.UserID)
}

func TestNewJWTVerifier_Introspection_Active(t *testing.T) {
	t.Parallel()

	introServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
		user, pass, ok := r.BasicAuth()
		assert.True(t, ok)
		assert.Equal(t, "my-client", user)
		assert.Equal(t, "my-secret", pass)

		json.NewEncoder(w).Encode(introspectionResponse{
			Active: true,
			Sub:    "user-intro",
			Exp:    time.Now().Add(time.Hour).Unix(),
		})
	}))
	defer introServer.Close()

	verifier := newJWTVerifier(&JWTConfig{
		Type:             "introspection",
		IntrospectionURL: introServer.URL,
		ClientID:         "my-client",
		ClientSecret:     "my-secret",
	})

	info, err := verifier(context.Background(), "some-token", nil)
	require.NoError(t, err)
	assert.Equal(t, "user-intro", info.UserID)
	assert.False(t, info.Expiration.IsZero())
}

func TestNewJWTVerifier_Introspection_Inactive(t *testing.T) {
	t.Parallel()

	introServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(introspectionResponse{Active: false})
	}))
	defer introServer.Close()

	verifier := newJWTVerifier(&JWTConfig{
		Type:             "introspection",
		IntrospectionURL: introServer.URL,
		ClientID:         "my-client",
		ClientSecret:     "my-secret",
	})

	_, err := verifier(context.Background(), "some-token", nil)
	assert.ErrorIs(t, err, auth.ErrInvalidToken)
}

func TestContains(t *testing.T) {
	t.Parallel()

	assert.True(t, slices.Contains([]string{"a", "b", "c"}, "b"))
	assert.False(t, slices.Contains([]string{"a", "b", "c"}, "d"))
	assert.False(t, slices.Contains([]string{}, "a"))
}

// helpers

func signJWT(t *testing.T, key *rsa.PrivateKey, claims jwt.Claims) string {
	t.Helper()

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := token.SignedString(key)
	require.NoError(t, err)
	return signed
}

func newJWKSServer(t *testing.T, pubKey *rsa.PublicKey) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()
	mux.HandleFunc("/.well-known/jwks.json", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jwksMarshalFromKey(t, pubKey))
	})

	return httptest.NewServer(mux)
}

func jwksMarshalFromKey(t *testing.T, pubKey *rsa.PublicKey) jwkset.JWKSMarshal {
	t.Helper()

	jwk, err := jwkset.NewJWKFromKey(pubKey, jwkset.JWKOptions{
		Metadata: jwkset.JWKMetadataOptions{
			KID: "test-kid",
			ALG: jwkset.AlgRS256,
		},
	})
	require.NoError(t, err)

	storage := jwkset.NewMemoryStorage()
	err = storage.KeyWrite(context.Background(), jwk)
	require.NoError(t, err)

	marshal, err := storage.Marshal(context.Background())
	require.NoError(t, err)

	return marshal
}
