/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package service

import (
	"testing"

	"github.com/mmadfox/swag2mcp/internal/auth"
	"github.com/mmadfox/swag2mcp/internal/spec"
	"github.com/stretchr/testify/require"
)

func TestValidateParameters_unknown(t *testing.T) {
	t.Parallel()

	op := &spec.Operation{
		Parameters: []*spec.Parameter{
			{Name: "id", In: "path"},
		},
	}
	err := validateParameters(op, map[string]any{"unknown": "val"}, nil)
	require.NoError(t, err, "unknown parameters should be silently ignored")
}

func TestValidateParameters_missingRequired(t *testing.T) {
	t.Parallel()

	op := &spec.Operation{
		Parameters: []*spec.Parameter{
			{Name: "id", In: "path", Required: true},
		},
	}
	err := validateParameters(op, nil, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing required")
}

func TestValidateParameters_valid(t *testing.T) {
	t.Parallel()

	op := &spec.Operation{
		Parameters: []*spec.Parameter{
			{Name: "id", In: "path", Required: true},
		},
	}
	err := validateParameters(op, map[string]any{"id": "123"}, nil)
	require.NoError(t, err)
}

func TestValidateParameters_authCoversRequired(t *testing.T) {
	t.Parallel()

	op := &spec.Operation{
		Parameters: []*spec.Parameter{
			{Name: "api_key", In: "query", Required: true},
			{Name: "ip_address", In: "query", Required: false},
		},
	}
	err := validateParameters(op, nil, map[string]struct{}{"api_key": {}})
	require.NoError(t, err, "api_key is covered by auth and must not be required")
}

func TestValidateParameters_authCoversOnlySome(t *testing.T) {
	t.Parallel()

	op := &spec.Operation{
		Parameters: []*spec.Parameter{
			{Name: "api_key", In: "query", Required: true},
			{Name: "ip_address", In: "query", Required: true},
		},
	}
	err := validateParameters(op, nil, map[string]struct{}{"api_key": {}})
	require.Error(t, err)
	require.Contains(t, err.Error(), "ip_address")
	require.NotContains(t, err.Error(), "api_key")
}

func TestValidateParameters_authCoversNonRequired(t *testing.T) {
	t.Parallel()

	op := &spec.Operation{
		Parameters: []*spec.Parameter{
			{Name: "api_key", In: "query", Required: false},
		},
	}
	err := validateParameters(op, nil, map[string]struct{}{"api_key": {}})
	require.NoError(t, err)
}

func TestAuthQueryParamNames_nil(t *testing.T) {
	t.Parallel()

	require.Nil(t, authQueryParamNames(nil))
}

func TestAuthQueryParamNames_apiKeyQuery(t *testing.T) {
	t.Parallel()

	a := &auth.APIKeyAuthClient{Key: "api_key", In: "query"}
	require.Equal(t, map[string]struct{}{"api_key": {}}, authQueryParamNames(a))
}

func TestAuthQueryParamNames_apiKeyHeader(t *testing.T) {
	t.Parallel()

	a := &auth.APIKeyAuthClient{Key: "api_key", In: "header"}
	require.Nil(t, authQueryParamNames(a))
}

func TestAuthQueryParamNames_hmac(t *testing.T) {
	t.Parallel()

	a := &auth.HMACAuthClient{Options: &auth.HMACOptions{RecvWindow: 5000}}
	got := authQueryParamNames(a)
	require.Equal(t, map[string]struct{}{
		"timestamp":  {},
		"signature":  {},
		"recvWindow": {},
	}, got)
}

func TestAuthQueryParamNames_hmacNoRecvWindow(t *testing.T) {
	t.Parallel()

	a := &auth.HMACAuthClient{}
	got := authQueryParamNames(a)
	require.Equal(t, map[string]struct{}{
		"timestamp": {},
		"signature": {},
	}, got)
}

func TestAuthQueryParamNames_headerOnly(t *testing.T) {
	t.Parallel()

	for _, a := range []auth.Authenticator{
		&auth.BasicAuthClient{},
		&auth.BearerTokenAuthClient{},
		&auth.DigestAuthClient{},
		&auth.OAuth2ClientCredentialsAuthClient{},
		&auth.OAuth2PasswordAuthClient{},
		&auth.ScriptAuthClient{},
		auth.NewNoAuthClient(),
	} {
		require.Nil(t, authQueryParamNames(a), "auth type %s must not inject query params", a.Type())
	}
}

func TestValidateRequestBody_notRequired(t *testing.T) {
	t.Parallel()

	op := &spec.Operation{
		RequestBody: &spec.RequestBody{Required: false},
	}
	err := validateRequestBody(op, nil)
	require.NoError(t, err)
}

func TestValidateRequestBody_requiredMissing(t *testing.T) {
	t.Parallel()

	op := &spec.Operation{
		RequestBody: &spec.RequestBody{Required: true},
	}
	err := validateRequestBody(op, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "request body is required")
}

func TestValidateRequestBody_unknownField(t *testing.T) {
	t.Parallel()

	op := &spec.Operation{
		RequestBody: &spec.RequestBody{
			Required: true,
			Content: map[string]*spec.MediaType{
				"application/json": {
					Schema: &spec.Schema{
						Type:       "object",
						Properties: map[string]*spec.Schema{"name": {Type: "string"}},
					},
				},
			},
		},
	}
	err := validateRequestBody(op, map[string]any{"unknown": "val"})
	require.NoError(t, err, "unknown fields in request body should be silently ignored")
}

func TestValidateRequestBody_valid(t *testing.T) {
	t.Parallel()

	op := &spec.Operation{
		RequestBody: &spec.RequestBody{
			Required: true,
			Content: map[string]*spec.MediaType{
				"application/json": {
					Schema: &spec.Schema{
						Type:       "object",
						Properties: map[string]*spec.Schema{"name": {Type: "string"}},
					},
				},
			},
		},
	}
	err := validateRequestBody(op, map[string]any{"name": "test"})
	require.NoError(t, err)
}

func TestSchemaForContentType_nil(t *testing.T) {
	t.Parallel()

	require.Nil(t, schemaForContentType(nil))
}

func TestSchemaForContentType_noJSON(t *testing.T) {
	t.Parallel()

	ct := map[string]*spec.MediaType{
		"text/plain": {Schema: &spec.Schema{Type: "string"}},
	}
	require.Nil(t, schemaForContentType(ct))
}

func TestSchemaForContentType_found(t *testing.T) {
	t.Parallel()

	ct := map[string]*spec.MediaType{
		"application/json": {Schema: &spec.Schema{Type: "object"}},
	}
	require.NotNil(t, schemaForContentType(ct))
}

func TestValidateArraySchema_valid(t *testing.T) {
	t.Parallel()

	sc := &spec.Schema{
		Type: "array",
		Items: &spec.Schema{
			Type:       "object",
			Properties: map[string]*spec.Schema{"name": {Type: "string"}},
		},
	}
	err := validateSchemaValue(sc, []any{map[string]any{"name": "test"}}, "$")
	require.NoError(t, err)
}

func TestValidateArraySchema_missingField(t *testing.T) {
	t.Parallel()

	sc := &spec.Schema{
		Type: "array",
		Items: &spec.Schema{
			Type:       "object",
			Required:   []string{"name"},
			Properties: map[string]*spec.Schema{"name": {Type: "string"}},
		},
	}
	err := validateSchemaValue(sc, []any{map[string]any{}}, "$")
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing required field")
}

func TestValidateArraySchema_notArray(t *testing.T) {
	t.Parallel()

	sc := &spec.Schema{
		Type: "array",
		Items: &spec.Schema{
			Type: "string",
		},
	}
	err := validateSchemaValue(sc, "not-an-array", "$")
	require.NoError(t, err)
}

func TestValidateSchemaValue_nilSchema(t *testing.T) {
	t.Parallel()

	err := validateSchemaValue(nil, "anything", "$")
	require.NoError(t, err)
}

func TestValidateSchemaValue_unknownType(t *testing.T) {
	t.Parallel()

	sc := &spec.Schema{Type: "string"}
	err := validateSchemaValue(sc, "hello", "$")
	require.NoError(t, err)
}

func TestValidateObjectSchema_notObject(t *testing.T) {
	t.Parallel()

	sc := &spec.Schema{
		Type:       "object",
		Properties: map[string]*spec.Schema{"name": {Type: "string"}},
	}
	err := validateSchemaValue(sc, "not-an-object", "$")
	require.NoError(t, err)
}

func TestValidateObjectSchema_unknownField(t *testing.T) {
	t.Parallel()

	sc := &spec.Schema{
		Type:       "object",
		Properties: map[string]*spec.Schema{"name": {Type: "string"}},
	}
	err := validateSchemaValue(sc, map[string]any{"unknown": "val"}, "$")
	require.NoError(t, err, "unknown fields in object schema should be silently ignored")
}

func TestValidateRequestBody_nilContent(t *testing.T) {
	t.Parallel()

	op := &spec.Operation{
		RequestBody: &spec.RequestBody{
			Required: true,
			Content:  nil,
		},
	}
	err := validateRequestBody(op, map[string]any{"name": "test"})
	require.NoError(t, err)
}

func TestValidateRequestBody_noJSONContent(t *testing.T) {
	t.Parallel()

	op := &spec.Operation{
		RequestBody: &spec.RequestBody{
			Required: true,
			Content: map[string]*spec.MediaType{
				"text/plain": {Schema: &spec.Schema{Type: "string"}},
			},
		},
	}
	err := validateRequestBody(op, map[string]any{"name": "test"})
	require.NoError(t, err)
}

func TestValidateRequestBody_graphQLStyle(t *testing.T) {
	t.Parallel()

	op := &spec.Operation{
		RequestBody: &spec.RequestBody{
			Required: true,
			Content: map[string]*spec.MediaType{
				"application/json": {
					Schema: &spec.Schema{
						Type:        "object",
						Description: "GraphQL query",
					},
				},
			},
		},
	}
	err := validateRequestBody(op, map[string]any{"query": "{ users { id } }"})
	require.NoError(t, err)
}

func TestValidateRequestBody_graphQLStyleWithVariables(t *testing.T) {
	t.Parallel()

	op := &spec.Operation{
		RequestBody: &spec.RequestBody{
			Required: true,
			Content: map[string]*spec.MediaType{
				"application/json": {
					Schema: &spec.Schema{
						Type:        "object",
						Description: "GraphQL query",
					},
				},
			},
		},
	}
	err := validateRequestBody(op, map[string]any{
		"query":     "query($id: ID!) { user(id: $id) { name } }",
		"variables": map[string]any{"id": "123"},
	})
	require.NoError(t, err)
}

func TestValidateObjectSchema_noProperties(t *testing.T) {
	t.Parallel()

	sc := &spec.Schema{
		Type: "object",
	}
	err := validateSchemaValue(sc, map[string]any{"query": "test", "variables": map[string]any{"id": "1"}}, "$")
	require.NoError(t, err)
}

func TestValidateObjectSchema_noPropertiesWithRequired(t *testing.T) {
	t.Parallel()

	sc := &spec.Schema{
		Type:     "object",
		Required: []string{"query"},
	}
	err := validateSchemaValue(sc, map[string]any{"query": "test"}, "$")
	require.NoError(t, err)
}

func TestValidateObjectSchema_noPropertiesMissingRequired(t *testing.T) {
	t.Parallel()

	sc := &spec.Schema{
		Type:     "object",
		Required: []string{"query"},
	}
	err := validateSchemaValue(sc, map[string]any{"variables": "test"}, "$")
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing required field")
}
