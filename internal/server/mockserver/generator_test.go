package mockserver

import (
	"testing"

	"github.com/mmadfox/swag2mcp/internal/spec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateFromSchema_NilSchema(t *testing.T) {
	t.Parallel()

	result := GenerateFromSchema(nil)
	assert.Nil(t, result, "expected nil")
}

func TestGenerateFromSchema_String(t *testing.T) {
	t.Parallel()

	schema := &spec.Schema{Type: "string"}
	result := GenerateFromSchema(schema)
	_, ok := result.(string)
	assert.True(t, ok, "expected string, got %T", result)
}

func TestGenerateFromSchema_Integer(t *testing.T) {
	t.Parallel()

	schema := &spec.Schema{Type: "integer"}
	result := GenerateFromSchema(schema)
	_, ok := result.(int)
	assert.True(t, ok, "expected int, got %T", result)
}

func TestGenerateFromSchema_Number(t *testing.T) {
	t.Parallel()

	schema := &spec.Schema{Type: "number"}
	result := GenerateFromSchema(schema)
	_, ok := result.(float64)
	assert.True(t, ok, "expected float64, got %T", result)
}

func TestGenerateFromSchema_Boolean(t *testing.T) {
	t.Parallel()

	schema := &spec.Schema{Type: "boolean"}
	result := GenerateFromSchema(schema)
	_, ok := result.(bool)
	assert.True(t, ok, "expected bool, got %T", result)
}

func TestGenerateFromSchema_Array(t *testing.T) {
	t.Parallel()

	schema := &spec.Schema{
		Type: "array",
		Items: &spec.Schema{
			Type: "string",
		},
	}
	result := GenerateFromSchema(schema)
	array, ok := result.([]any)
	require.True(t, ok, "expected []any, got %T", result)
	assert.GreaterOrEqual(t, len(array), 1, "expected array length 1-3")
	assert.LessOrEqual(t, len(array), 3, "expected array length 1-3")
}

func TestGenerateFromSchema_Object(t *testing.T) {
	t.Parallel()

	schema := &spec.Schema{
		Type: "object",
		Properties: map[string]*spec.Schema{
			"name": {Type: "string"},
			"age":  {Type: "integer"},
		},
	}
	result := GenerateFromSchema(schema)
	obj, ok := result.(map[string]any)
	require.True(t, ok, "expected map[string]any, got %T", result)
	assert.Contains(t, obj, "name", "expected 'name' key in result")
	assert.Contains(t, obj, "age", "expected 'age' key in result")
}

func TestGenerateFromSchema_Enum(t *testing.T) {
	t.Parallel()

	schema := &spec.Schema{
		Type: "string",
		Enum: []any{"red", "green", "blue"},
	}
	result := GenerateFromSchema(schema)
	value, ok := result.(string)
	require.True(t, ok, "expected string, got %T", result)
	valid := value == "red" || value == "green" || value == "blue"
	assert.True(t, valid, "expected one of [red, green, blue], got %q", value)
}

func TestGenerateFromSchema_Default(t *testing.T) {
	t.Parallel()

	schema := &spec.Schema{
		Type:    "string",
		Default: "default-value",
	}
	result := GenerateFromSchema(schema)
	assert.Equal(t, "default-value", result, "expected default value")
}

func TestGenerateFromSchema_Example(t *testing.T) {
	t.Parallel()

	schema := &spec.Schema{
		Type:    "string",
		Example: "example-value",
	}
	result := GenerateFromSchema(schema)
	assert.Equal(t, "example-value", result, "expected example value")
}

func TestGenerateFromSchema_AllOf(t *testing.T) {
	t.Parallel()

	schema := &spec.Schema{
		AllOf: []*spec.Schema{
			{
				Type: "object",
				Properties: map[string]*spec.Schema{
					"name": {Type: "string"},
				},
			},
			{
				Type: "object",
				Properties: map[string]*spec.Schema{
					"age": {Type: "integer"},
				},
			},
		},
	}
	result := GenerateFromSchema(schema)
	obj, ok := result.(map[string]any)
	require.True(t, ok, "expected map[string]any, got %T", result)
	assert.Contains(t, obj, "name", "expected 'name' key from allOf merge")
	assert.Contains(t, obj, "age", "expected 'age' key from allOf merge")
}

func TestGenerateFromSchema_OneOf(t *testing.T) {
	t.Parallel()

	schema := &spec.Schema{
		OneOf: []*spec.Schema{
			{Type: "string"},
			{Type: "integer"},
		},
	}
	result := GenerateFromSchema(schema)
	switch result.(type) {
	case string, int:
	default:
		t.Errorf("expected string or int, got %T", result)
	}
}

func TestGenerateFromSchema_AnyOf(t *testing.T) {
	t.Parallel()

	schema := &spec.Schema{
		AnyOf: []*spec.Schema{
			{Type: "string"},
			{Type: "integer"},
		},
	}
	result := GenerateFromSchema(schema)
	switch result.(type) {
	case string, int:
	default:
		t.Errorf("expected string or int, got %T", result)
	}
}

func TestGenerateFromSchema_StringFormat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		format string
	}{
		{"date"},
		{"date-time"},
		{"email"},
		{"uri"},
		{"uuid"},
		{"ipv4"},
		{"hostname"},
		{"phone"},
		{"byte"},
		{"binary"},
		{"password"},
	}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			t.Parallel()
			schema := &spec.Schema{Type: "string", Format: tt.format}
			result := GenerateFromSchema(schema)
			_, ok := result.(string)
			assert.True(t, ok, "expected string for format %q, got %T", tt.format, result)
		})
	}
}

func TestGenerateFromSchema_IntegerFormat(t *testing.T) {
	t.Parallel()

	t.Run("int32", func(t *testing.T) {
		t.Parallel()
		schema := &spec.Schema{Type: "integer", Format: "int32"}
		result := GenerateFromSchema(schema)
		_, ok := result.(int32)
		assert.True(t, ok, "expected int32, got %T", result)
	})

	t.Run("int64", func(t *testing.T) {
		t.Parallel()
		schema := &spec.Schema{Type: "integer", Format: "int64"}
		result := GenerateFromSchema(schema)
		_, ok := result.(int64)
		assert.True(t, ok, "expected int64, got %T", result)
	})
}

func TestGenerateFromSchema_ReadOnlySkipped(t *testing.T) {
	t.Parallel()

	schema := &spec.Schema{
		Type: "object",
		Properties: map[string]*spec.Schema{
			"readOnlyField":  {Type: "string", ReadOnly: true},
			"writeOnlyField": {Type: "string", WriteOnly: true},
			"normalField":    {Type: "string"},
		},
	}
	result := GenerateFromSchema(schema)
	obj, ok := result.(map[string]any)
	require.True(t, ok, "expected map[string]any, got %T", result)
	assert.NotContains(t, obj, "readOnlyField", "readOnly field should be skipped")
	assert.NotContains(t, obj, "writeOnlyField", "writeOnly field should be skipped")
	assert.Contains(t, obj, "normalField", "normal field should be present")
}

func TestGenerateRandomToken(t *testing.T) {
	t.Parallel()

	token := generateRandomToken()
	assert.Len(t, token, authTokenLength*2, "expected token length %d", authTokenLength*2)
}
