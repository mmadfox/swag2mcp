package service

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"testing"

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
	err := validateParameters(op, map[string]any{"unknown": "val"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "unknown parameter")
}

func TestValidateParameters_missingRequired(t *testing.T) {
	t.Parallel()

	op := &spec.Operation{
		Parameters: []*spec.Parameter{
			{Name: "id", In: "path", Required: true},
		},
	}
	err := validateParameters(op, nil)
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
	err := validateParameters(op, map[string]any{"id": "123"})
	require.NoError(t, err)
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
	require.Error(t, err)
	require.Contains(t, err.Error(), "unknown field")
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
	require.Error(t, err)
	require.Contains(t, err.Error(), "unknown field")
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
