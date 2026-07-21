package service

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
