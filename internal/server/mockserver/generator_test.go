package mockserver

import (
	"testing"

	"github.com/mmadfox/swag2mcp/internal/spec"
)

func TestGenerateFromSchema_NilSchema(t *testing.T) {
	t.Parallel()

	result := GenerateFromSchema(nil)
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestGenerateFromSchema_String(t *testing.T) {
	t.Parallel()

	schema := &spec.Schema{Type: "string"}
	result := GenerateFromSchema(schema)
	_, ok := result.(string)
	if !ok {
		t.Errorf("expected string, got %T", result)
	}
}

func TestGenerateFromSchema_Integer(t *testing.T) {
	t.Parallel()

	schema := &spec.Schema{Type: "integer"}
	result := GenerateFromSchema(schema)
	_, ok := result.(int)
	if !ok {
		t.Errorf("expected int, got %T", result)
	}
}

func TestGenerateFromSchema_Number(t *testing.T) {
	t.Parallel()

	schema := &spec.Schema{Type: "number"}
	result := GenerateFromSchema(schema)
	_, ok := result.(float64)
	if !ok {
		t.Errorf("expected float64, got %T", result)
	}
}

func TestGenerateFromSchema_Boolean(t *testing.T) {
	t.Parallel()

	schema := &spec.Schema{Type: "boolean"}
	result := GenerateFromSchema(schema)
	_, ok := result.(bool)
	if !ok {
		t.Errorf("expected bool, got %T", result)
	}
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
	if !ok {
		t.Fatalf("expected []any, got %T", result)
	}
	if len(array) < 1 || len(array) > 3 {
		t.Errorf("expected array length 1-3, got %d", len(array))
	}
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
	if !ok {
		t.Fatalf("expected map[string]any, got %T", result)
	}
	if _, exists := obj["name"]; !exists {
		t.Error("expected 'name' key in result")
	}
	if _, exists := obj["age"]; !exists {
		t.Error("expected 'age' key in result")
	}
}

func TestGenerateFromSchema_Enum(t *testing.T) {
	t.Parallel()

	schema := &spec.Schema{
		Type: "string",
		Enum: []any{"red", "green", "blue"},
	}
	result := GenerateFromSchema(schema)
	value, ok := result.(string)
	if !ok {
		t.Fatalf("expected string, got %T", result)
	}
	valid := value == "red" || value == "green" || value == "blue"
	if !valid {
		t.Errorf("expected one of [red, green, blue], got %q", value)
	}
}

func TestGenerateFromSchema_Default(t *testing.T) {
	t.Parallel()

	schema := &spec.Schema{
		Type:    "string",
		Default: "default-value",
	}
	result := GenerateFromSchema(schema)
	if result != "default-value" {
		t.Errorf("expected %q, got %v", "default-value", result)
	}
}

func TestGenerateFromSchema_Example(t *testing.T) {
	t.Parallel()

	schema := &spec.Schema{
		Type:    "string",
		Example: "example-value",
	}
	result := GenerateFromSchema(schema)
	if result != "example-value" {
		t.Errorf("expected %q, got %v", "example-value", result)
	}
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
	if !ok {
		t.Fatalf("expected map[string]any, got %T", result)
	}
	if _, exists := obj["name"]; !exists {
		t.Error("expected 'name' key from allOf merge")
	}
	if _, exists := obj["age"]; !exists {
		t.Error("expected 'age' key from allOf merge")
	}
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
			if !ok {
				t.Errorf("expected string for format %q, got %T", tt.format, result)
			}
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
		if !ok {
			t.Errorf("expected int32, got %T", result)
		}
	})

	t.Run("int64", func(t *testing.T) {
		t.Parallel()
		schema := &spec.Schema{Type: "integer", Format: "int64"}
		result := GenerateFromSchema(schema)
		_, ok := result.(int64)
		if !ok {
			t.Errorf("expected int64, got %T", result)
		}
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
	if !ok {
		t.Fatalf("expected map[string]any, got %T", result)
	}
	if _, exists := obj["readOnlyField"]; exists {
		t.Error("readOnly field should be skipped")
	}
	if _, exists := obj["writeOnlyField"]; exists {
		t.Error("writeOnly field should be skipped")
	}
	if _, exists := obj["normalField"]; !exists {
		t.Error("normal field should be present")
	}
}

func TestGenerateRandomToken(t *testing.T) {
	t.Parallel()

	token := generateRandomToken()
	if len(token) != authTokenLength*2 {
		t.Errorf("expected token length %d, got %d", authTokenLength*2, len(token))
	}
}
