/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package spec

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseV3_InvalidJSON(t *testing.T) {
	t.Parallel()

	_, err := parseV3([]byte(`{"openapi":"3.0.0","paths":"not-an-object"}`))
	require.Error(t, err, "expected error for invalid OpenAPI document")
}

func TestOpenapi3OpToOp_NilValue(t *testing.T) {
	t.Parallel()

	op := openapi3OpToOp(&openapi3.Operation{
		OperationID: "testOp",
		Parameters: []*openapi3.ParameterRef{
			nil,
			{Value: nil},
		},
	})
	require.NotNil(t, op, "op is nil")
	assert.Empty(t, op.Parameters)
}

func TestOpenapi3OpToOp_NilDescription(t *testing.T) {
	t.Parallel()

	responses := openapi3.NewResponses()
	responses.Set("200", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: nil,
		},
	})

	op := openapi3OpToOp(&openapi3.Operation{
		OperationID: "testOp",
		Responses:   responses,
	})
	require.NotNil(t, op, "op is nil")
	resp, ok := op.Responses["200"]
	require.True(t, ok, "200 response not found")
	assert.Empty(t, resp.Description)
}

func TestSchemaRefToSchema_Nil(t *testing.T) {
	t.Parallel()

	s := schemaRefToSchema(nil)
	require.Nil(t, s, "expected nil")
}

func TestSchemaRefToSchema_MultipleTypes(t *testing.T) {
	t.Parallel()

	s := schemaRefToSchema(&openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type: &openapi3.Types{"string", "null"},
		},
	})
	require.NotNil(t, s, "schema is nil")
	assert.Equal(t, "string", s.Type)
}

func TestSchemaRefToSchema_Items(t *testing.T) {
	t.Parallel()

	s := schemaRefToSchema(&openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Items: &openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: &openapi3.Types{"string"},
				},
			},
		},
	})
	require.NotNil(t, s, "schema is nil")
	require.NotNil(t, s.Items, "Items is nil")
	assert.Equal(t, "string", s.Items.Type)
}

func TestSchemaRefToSchema_Composition(t *testing.T) {
	t.Parallel()

	s := schemaRefToSchema(&openapi3.SchemaRef{
		Value: &openapi3.Schema{
			OneOf: []*openapi3.SchemaRef{
				{Value: &openapi3.Schema{Type: &openapi3.Types{"string"}}},
			},
			AnyOf: []*openapi3.SchemaRef{
				{Value: &openapi3.Schema{Type: &openapi3.Types{"integer"}}},
			},
			AllOf: []*openapi3.SchemaRef{
				{Value: &openapi3.Schema{Type: &openapi3.Types{"number"}}},
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

func TestOpenapi3PathItemToOps_MergesPathLevelParams(t *testing.T) {
	t.Parallel()

	item := &openapi3.PathItem{
		Parameters: []*openapi3.ParameterRef{
			{
				Value: &openapi3.Parameter{
					Name:     "id",
					In:       "path",
					Required: true,
					Schema:   &openapi3.SchemaRef{Value: &openapi3.Schema{Type: &openapi3.Types{"string"}}},
				},
			},
		},
		Get: &openapi3.Operation{
			OperationID: "getItem",
			Parameters: []*openapi3.ParameterRef{
				{
					Value: &openapi3.Parameter{
						Name:   "limit",
						In:     "query",
						Schema: &openapi3.SchemaRef{Value: &openapi3.Schema{Type: &openapi3.Types{"integer"}}},
					},
				},
			},
		},
		Post: &openapi3.Operation{
			OperationID: "createItem",
		},
	}

	items := openapi3PathItemToOps("/items/{id}", item)
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

func TestOpenapi3PathItemToOps_OpLevelOverridesPathLevel(t *testing.T) {
	t.Parallel()

	item := &openapi3.PathItem{
		Parameters: []*openapi3.ParameterRef{
			{
				Value: &openapi3.Parameter{
					Name:     "id",
					In:       "path",
					Required: true,
					Schema:   &openapi3.SchemaRef{Value: &openapi3.Schema{Type: &openapi3.Types{"string"}}},
				},
			},
		},
		Get: &openapi3.Operation{
			OperationID: "getItem",
			Parameters: []*openapi3.ParameterRef{
				{
					Value: &openapi3.Parameter{
						Name:     "id",
						In:       "path",
						Required: true,
						Schema: &openapi3.SchemaRef{
							Value: &openapi3.Schema{Type: &openapi3.Types{"integer"}},
						}, // override type
					},
				},
			},
		},
	}

	items := openapi3PathItemToOps("/items/{id}", item)
	require.Len(t, items, 1)

	getOp := items[0]
	require.Len(t, getOp.Operation.Parameters, 1)
	require.NotNil(t, getOp.Operation.Parameters[0].Schema)
	assert.Equal(t, "integer", getOp.Operation.Parameters[0].Schema.Type)
}

func TestOpenapi3OpToOp_Security(t *testing.T) {
	t.Parallel()

	sec := openapi3.SecurityRequirements{
		{"bearerAuth": {}},
	}
	op := openapi3OpToOp(&openapi3.Operation{
		OperationID: "secureOp",
		Security:    &sec,
	})

	require.NotNil(t, op)
	require.Len(t, op.Security, 1)
	assert.Contains(t, op.Security[0], "bearerAuth")
}

func TestOpenapi3OpToOp_Security_Nil(t *testing.T) {
	t.Parallel()

	op := openapi3OpToOp(&openapi3.Operation{
		OperationID: "publicOp",
	})

	require.NotNil(t, op)
	assert.Empty(t, op.Security)
}

func TestOpenapi3OpToOp_Security_Empty(t *testing.T) {
	t.Parallel()

	sec := openapi3.SecurityRequirements{}
	op := openapi3OpToOp(&openapi3.Operation{
		OperationID: "noSecurityOp",
		Security:    &sec,
	})

	require.NotNil(t, op)
	assert.Empty(t, op.Security)
}
