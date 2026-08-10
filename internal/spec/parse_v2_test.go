/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package spec

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse_swaggerHost(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile(filepath.Join("testdata", "valid_v20_swagger.yaml"))
	require.NoError(t, err)

	doc, err := Parse(data)
	require.NoError(t, err)

	require.NotEmpty(t, doc.Servers, "expected at least 1 server from swagger host")

	assert.Equal(t, "https://api.example.com/v1", doc.Servers[0].URL)
}

func TestParse_operationMetadata(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile(filepath.Join("testdata", "valid_v20_swagger.yaml"))
	require.NoError(t, err)

	doc, err := Parse(data)
	require.NoError(t, err)

	var found bool
	for _, pi := range doc.PathItems {
		if pi.Path == "/users" && pi.Method == http.MethodGet {
			found = true
			op := pi.Operation
			assert.Equal(t, "Список пользователей", op.Summary)
			require.NotEmpty(t, op.Parameters, "expected parameters")
			assert.Equal(t, "limit", op.Parameters[0].Name)
			assert.Equal(t, "query", op.Parameters[0].In)
			break
		}
	}
	require.True(t, found, "GET /users not found in parsed doc")
}

func TestParse_requestBody(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile(filepath.Join("testdata", "valid_v20_swagger.yaml"))
	require.NoError(t, err)

	doc, err := Parse(data)
	require.NoError(t, err)

	var found bool
	for _, pi := range doc.PathItems {
		if pi.Path == "/users" && pi.Method == http.MethodPost {
			found = true
			op := pi.Operation
			require.NotNil(t, op.RequestBody, "expected request body")
			assert.True(t, op.RequestBody.Required, "expected required request body")
			require.NotNil(t, op.RequestBody.Content, "expected request body content")

			for _, p := range op.Parameters {
				assert.NotEqual(t, "body", p.In, "body param should not appear in parameters list")
			}
			break
		}
	}
	require.True(t, found, "POST /users not found")
}

func TestParse_responses(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile(filepath.Join("testdata", "valid_v20_swagger.yaml"))
	require.NoError(t, err)

	doc, err := Parse(data)
	require.NoError(t, err)

	var found bool
	for _, pi := range doc.PathItems {
		if pi.Path == "/users" && pi.Method == http.MethodGet {
			found = true
			op := pi.Operation
			require.NotEmpty(t, op.Responses, "expected responses")
			resp, ok := op.Responses["200"]
			require.True(t, ok, "expected 200 response")
			assert.Equal(t, "OK", resp.Description)
			break
		}
	}
	require.True(t, found, "GET /users not found")
}

func TestParse_swaggerFileUpload(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile(filepath.Join("testdata", "valid_v20_swagger.yaml"))
	require.NoError(t, err)

	doc, err := Parse(data)
	require.NoError(t, err)

	var found bool
	for _, pi := range doc.PathItems {
		if pi.Path == "/files/upload" && pi.Method == http.MethodPost {
			found = true
			op := pi.Operation
			require.NotEmpty(t, op.Parameters, "expected parameters")
			assert.Equal(t, "file", op.Parameters[0].Name)
			assert.Equal(t, "formData", op.Parameters[0].In)
			break
		}
	}
	require.True(t, found, "POST /files/upload not found")
}

func TestParseV2_InvalidJSON(t *testing.T) {
	t.Parallel()

	_, err := parseV2([]byte("{invalid}"))
	require.Error(t, err, "expected error for invalid JSON")
}

func TestParseV2_EmptyHost(t *testing.T) {
	t.Parallel()

	jsonData := []byte(`{"swagger":"2.0","info":{"title":"Test","version":"1.0"},"paths":{}}`)
	result, err := parseV2(jsonData)
	require.NoError(t, err, "parseV2() failed")
	assert.Empty(t, result.Servers)
}

func TestSwaggerSchemaToSchema_Nil(t *testing.T) {
	t.Parallel()

	s := swaggerSchemaToSchema(nil)
	require.Nil(t, s, "expected nil")
}

func TestSwaggerSchemaToSchema_Ref(t *testing.T) {
	t.Parallel()

	s := swaggerSchemaToSchema(&spec.Schema{
		SchemaProps: spec.SchemaProps{
			Ref: spec.MustCreateRef("#/definitions/Pet"),
		},
	})
	require.NotNil(t, s, "schema is nil")
	assert.Equal(t, "#/definitions/Pet", s.Ref)
}

func TestSwaggerSchemaToSchema_Items(t *testing.T) {
	t.Parallel()

	s := swaggerSchemaToSchema(&spec.Schema{
		SchemaProps: spec.SchemaProps{
			Items: &spec.SchemaOrArray{
				Schema: &spec.Schema{
					SchemaProps: spec.SchemaProps{
						Type: spec.StringOrArray{"string"},
					},
				},
			},
		},
	})
	require.NotNil(t, s, "schema is nil")
	require.NotNil(t, s.Items, "Items is nil")
	assert.Equal(t, "string", s.Items.Type)
}

func TestSwaggerSchemaToSchema_Properties(t *testing.T) {
	t.Parallel()

	s := swaggerSchemaToSchema(&spec.Schema{
		SchemaProps: spec.SchemaProps{
			Properties: map[string]spec.Schema{
				"name": {SchemaProps: spec.SchemaProps{Type: spec.StringOrArray{"string"}}},
			},
		},
	})
	require.NotNil(t, s, "schema is nil")
	require.Len(t, s.Properties, 1)
	assert.Equal(t, "string", s.Properties["name"].Type)
}

func TestSwaggerSchemaToSchema_OneOfAnyOfAllOf(t *testing.T) {
	t.Parallel()

	s := swaggerSchemaToSchema(&spec.Schema{
		SchemaProps: spec.SchemaProps{
			OneOf: []spec.Schema{
				{SchemaProps: spec.SchemaProps{Type: spec.StringOrArray{"string"}}},
			},
			AnyOf: []spec.Schema{
				{SchemaProps: spec.SchemaProps{Type: spec.StringOrArray{"integer"}}},
			},
			AllOf: []spec.Schema{
				{SchemaProps: spec.SchemaProps{Type: spec.StringOrArray{"number"}}},
			},
		},
	})
	require.NotNil(t, s, "schema is nil")
	require.Len(t, s.OneOf, 1)
	assert.Equal(t, "string", s.OneOf[0].Type)
	require.Len(t, s.AnyOf, 1)
	assert.Equal(t, "integer", s.AnyOf[0].Type)
	require.Len(t, s.AllOf, 1)
	assert.Equal(t, "number", s.AllOf[0].Type)
}

func TestSwaggerOpToOp_DefaultResponse(t *testing.T) {
	t.Parallel()

	op := swaggerOpToOp(&spec.Operation{
		OperationProps: spec.OperationProps{
			ID: "testOp",
			Responses: &spec.Responses{
				ResponsesProps: spec.ResponsesProps{
					Default: &spec.Response{
						ResponseProps: spec.ResponseProps{
							Description: "Default response",
						},
					},
				},
			},
		},
	})
	require.NotNil(t, op, "op is nil")
	resp, ok := op.Responses["default"]
	require.True(t, ok, "default response not found")
	assert.Equal(t, "Default response", resp.Description)
}

func TestPathItemToOps_MergesPathLevelParams(t *testing.T) {
	t.Parallel()

	item := spec.PathItem{
		PathItemProps: spec.PathItemProps{
			Parameters: []spec.Parameter{
				{
					ParamProps: spec.ParamProps{
						Name:     "id",
						In:       "path",
						Required: true,
					},
					SimpleSchema: spec.SimpleSchema{
						Type: "string",
					},
				},
			},
			Get: &spec.Operation{
				OperationProps: spec.OperationProps{
					ID: "getItem",
					Parameters: []spec.Parameter{
						{
							ParamProps: spec.ParamProps{
								Name: "limit",
								In:   "query",
							},
							SimpleSchema: spec.SimpleSchema{
								Type: "integer",
							},
						},
					},
				},
			},
			Post: &spec.Operation{
				OperationProps: spec.OperationProps{
					ID: "createItem",
				},
			},
		},
	}

	items := pathItemToOps("/items/{id}", item)
	require.Len(t, items, 2)

	// GET should have both path-level (id) and operation-level (limit) params
	getOp := items[0]
	require.Equal(t, "getItem", getOp.Operation.ID)
	require.Len(t, getOp.Operation.Parameters, 2)
	assert.Equal(t, "id", getOp.Operation.Parameters[0].Name)
	assert.Equal(t, "path", getOp.Operation.Parameters[0].In)
	assert.True(t, getOp.Operation.Parameters[0].Required)
	assert.Equal(t, "limit", getOp.Operation.Parameters[1].Name)
	assert.Equal(t, "query", getOp.Operation.Parameters[1].In)

	// POST should have only path-level params
	postOp := items[1]
	require.Equal(t, "createItem", postOp.Operation.ID)
	require.Len(t, postOp.Operation.Parameters, 1)
	assert.Equal(t, "id", postOp.Operation.Parameters[0].Name)
	assert.Equal(t, "path", postOp.Operation.Parameters[0].In)
}

func TestPathItemToOps_OpLevelOverridesPathLevel(t *testing.T) {
	t.Parallel()

	item := spec.PathItem{
		PathItemProps: spec.PathItemProps{
			Parameters: []spec.Parameter{
				{
					ParamProps: spec.ParamProps{
						Name:     "id",
						In:       "path",
						Required: true,
					},
					SimpleSchema: spec.SimpleSchema{
						Type: "string",
					},
				},
			},
			Get: &spec.Operation{
				OperationProps: spec.OperationProps{
					ID: "getItem",
					Parameters: []spec.Parameter{
						{
							ParamProps: spec.ParamProps{
								Name:     "id",
								In:       "path",
								Required: true,
							},
							SimpleSchema: spec.SimpleSchema{
								Type: "integer",
							},
						},
					},
				},
			},
		},
	}

	items := pathItemToOps("/items/{id}", item)
	require.Len(t, items, 1)

	getOp := items[0]
	require.Len(t, getOp.Operation.Parameters, 1)
	require.NotNil(t, getOp.Operation.Parameters[0].Schema)
	assert.Equal(t, "integer", getOp.Operation.Parameters[0].Schema.Type)
}

func TestSwaggerParamsToParams_Full(t *testing.T) {
	t.Parallel()

	params := []spec.Parameter{
		{
			ParamProps: spec.ParamProps{
				Name:        "id",
				In:          "path",
				Description: "Item ID",
				Required:    true,
			},
			SimpleSchema: spec.SimpleSchema{
				Type: "string",
			},
		},
		{
			ParamProps: spec.ParamProps{
				Name:        "limit",
				In:          "query",
				Description: "Page limit",
			},
			SimpleSchema: spec.SimpleSchema{
				Type:    "integer",
				Default: 10,
			},
		},
	}

	result := swaggerParamsToParams(params)
	require.Len(t, result, 2)

	assert.Equal(t, "id", result[0].Name)
	assert.Equal(t, "path", result[0].In)
	assert.True(t, result[0].Required)
	require.NotNil(t, result[0].Schema)
	assert.Equal(t, "string", result[0].Schema.Type)

	assert.Equal(t, "limit", result[1].Name)
	assert.Equal(t, "query", result[1].In)
	require.NotNil(t, result[1].Schema)
	assert.Equal(t, "integer", result[1].Schema.Type)
	assert.Equal(t, 10, result[1].Schema.Default)
}

func TestSwaggerParamsToParams_NilSchema(t *testing.T) {
	t.Parallel()

	params := []spec.Parameter{
		{
			ParamProps: spec.ParamProps{
				Name: "x-api-key",
				In:   "header",
			},
			// No SimpleSchema set - this tests the nil schema handling
		},
	}

	result := swaggerParamsToParams(params)
	require.Len(t, result, 1)
	assert.Equal(t, "x-api-key", result[0].Name)
	assert.Equal(t, "header", result[0].In)
	assert.Nil(t, result[0].Schema)
}

func TestSwaggerOpToOp_FiltersBodyParam(t *testing.T) {
	t.Parallel()

	op := swaggerOpToOp(&spec.Operation{
		OperationProps: spec.OperationProps{
			ID: "createUser",
			Parameters: []spec.Parameter{
				{
					ParamProps: spec.ParamProps{
						Name:     "userId",
						In:       "path",
						Required: true,
					},
					SimpleSchema: spec.SimpleSchema{Type: "string"},
				},
				{
					ParamProps: spec.ParamProps{
						Name:     "body",
						In:       "body",
						Required: true,
					},
					SimpleSchema: spec.SimpleSchema{Type: "object"},
				},
			},
		},
	})

	require.NotNil(t, op.RequestBody, "expected request body")
	assert.True(t, op.RequestBody.Required)

	for _, p := range op.Parameters {
		assert.NotEqual(t, "body", p.In, "body param should not appear in parameters")
	}
	assert.Len(t, op.Parameters, 1, "expected only non-body params")
	assert.Equal(t, "userId", op.Parameters[0].Name)
}

func TestSwaggerParamsToParams_FiltersBodyParam(t *testing.T) {
	t.Parallel()

	params := []spec.Parameter{
		{
			ParamProps: spec.ParamProps{
				Name: "id",
				In:   "path",
			},
			SimpleSchema: spec.SimpleSchema{Type: "string"},
		},
		{
			ParamProps: spec.ParamProps{
				Name:     "body",
				In:       "body",
				Required: true,
			},
		},
	}

	result := swaggerParamsToParams(params)
	for _, p := range result {
		assert.NotEqual(t, "body", p.In, "body param should not appear in result")
	}
	assert.Len(t, result, 1, "expected only non-body params")
	assert.Equal(t, "id", result[0].Name)
}

func TestSwaggerOpToOp_Security(t *testing.T) {
	t.Parallel()

	op := swaggerOpToOp(&spec.Operation{
		OperationProps: spec.OperationProps{
			ID: "secureOp",
			Security: []map[string][]string{
				{"ApiKeyAuth": {}},
			},
		},
	})

	require.NotNil(t, op)
	require.Len(t, op.Security, 1)
	assert.Contains(t, op.Security[0], "ApiKeyAuth")
}

func TestSwaggerOpToOp_Security_Empty(t *testing.T) {
	t.Parallel()

	op := swaggerOpToOp(&spec.Operation{
		OperationProps: spec.OperationProps{
			ID: "publicOp",
		},
	})

	require.NotNil(t, op)
	assert.Empty(t, op.Security)
}

func TestParseV2_TopLevelSecurityPropagates(t *testing.T) {
	t.Parallel()

	data, err := os.ReadFile(filepath.Join("testdata", "valid_v20_top_level_security.json"))
	require.NoError(t, err)

	doc, err := Parse(data)
	require.NoError(t, err)
	require.Len(t, doc.PathItems, 1)
	op := doc.PathItems[0].Operation
	require.NotNil(t, op)
	require.Len(t, op.Security, 1, "top-level security must propagate to operation")
	assert.Contains(t, op.Security[0], "UserSecurity")
}

func TestParseV2_TopLevelSecurityEmptyNotPropagated(t *testing.T) {
	t.Parallel()

	data, err := os.ReadFile(filepath.Join("testdata", "valid_v20_empty_security.json"))
	require.NoError(t, err)

	doc, err := Parse(data)
	require.NoError(t, err)
	require.Len(t, doc.PathItems, 1)
	assert.Empty(t, doc.PathItems[0].Operation.Security, "empty top-level security must not propagate")
}

func TestParseV2_OperationSecurityOverridesTopLevel(t *testing.T) {
	t.Parallel()

	data, err := os.ReadFile(filepath.Join("testdata", "valid_v20_operation_security_override.json"))
	require.NoError(t, err)

	doc, err := Parse(data)
	require.NoError(t, err)
	require.Len(t, doc.PathItems, 2)

	var publicOp, secureOp *Operation
	for _, pi := range doc.PathItems {
		switch pi.Path {
		case "/public":
			publicOp = pi.Operation
		case "/secure":
			secureOp = pi.Operation
		}
	}
	require.NotNil(t, publicOp)
	require.NotNil(t, secureOp)
	assert.Empty(t, publicOp.Security, "operation-level empty security must win")
	assert.Len(t, secureOp.Security, 1, "operation without security inherits top-level")
}
