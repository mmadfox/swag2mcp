package id

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDomain(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
			result := Domain(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCollection(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
			result := Collection(tt.domain, tt.spec)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTag(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		domain     string
		collection string
		tag        string
		expected   string
	}{
		{
			name:       "simple tag",
			domain:     "petstore",
			collection: "v1",
			tag:        "pets",
			expected:   "0c612ec95cacabbf36b22d967b1dff00",
		},
		{
			name:       "tag with spaces",
			domain:     "petstore",
			collection: "v1",
			tag:        "pet store",
			expected:   "5dc29dee1b2f7ffe9f41c91c3ff47753",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := Tag(tt.domain, tt.collection, tt.tag)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMethod(t *testing.T) {
	t.Parallel()
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
			expected:   "cb838b8e4a0a09492b88733a5996546d",
		},
		{
			name:       "method with path parameters",
			domain:     "petstore",
			collection: "v1",
			method:     "GET",
			path:       "/pets/{id}",
			opID:       "getPet",
			expected:   "e80ec318753224589986c6dd3b834e11",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := Method(tt.domain, tt.collection, "", tt.method, tt.path, tt.opID)
			assert.Equal(t, tt.expected, result)
		})
	}
}
