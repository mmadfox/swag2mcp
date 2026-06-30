package id

import (
	"testing"
)

func TestDomain(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple domain",
			input:    "petstore",
			expected: "a748d94c1b42369ef4df9c7dbc53639a",
		},
		{
			name:     "domain with dots",
			input:    "api.example.com",
			expected: "0aa7c02afb2118bf6c103c79876c7808",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Domain(tt.input)
			if result != tt.expected {
				t.Errorf("Domain(%s) = %s; want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCollection(t *testing.T) {
	tests := []struct {
		name     string
		domain   string
		spec     string
		expected string
	}{
		{
			name:     "simple collection",
			domain:   "petstore",
			spec:     "v1",
			expected: "19cbda7523af05362f044eff965f73f9",
		},
		{
			name:     "collection with dots",
			domain:   "api.example.com",
			spec:     "swagger.json",
			expected: "bfd7fdd0bba63fdf56e301f1a5347a06",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Collection(tt.domain, tt.spec)
			if result != tt.expected {
				t.Errorf("Collection(%s, %s) = %s; want %s", tt.domain, tt.spec, result, tt.expected)
			}
		})
	}
}

func TestMethod(t *testing.T) {
	tests := []struct {
		name       string
		domain     string
		collection string
		method     string
		path       string
		opID       string
		expected   string
	}{
		{
			name:       "simple method",
			domain:     "petstore",
			collection: "v1",
			method:     "GET",
			path:       "/pets",
			opID:       "listPets",
			expected:   "f6a7a32b8166a3c5f36ec26311b347c1",
		},
		{
			name:       "method with path parameters",
			domain:     "petstore",
			collection: "v1",
			method:     "GET",
			path:       "/pets/{id}",
			opID:       "getPet",
			expected:   "789115df2f8fd369d1c03abdc1e5fefa",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Method(tt.domain, tt.collection, "", tt.method, tt.path, tt.opID)
			if result != tt.expected {
				t.Errorf("Method(%s, %s, %s, %s, %s) = %s; want %s", tt.domain, tt.collection, tt.method, tt.path, tt.opID, result, tt.expected)
			}
		})
	}
}
