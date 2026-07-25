package mockserver

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"context"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMockServer_NoServers(t *testing.T) {
	t.Parallel()

	server := New(Options{
		Config: &config.Config{},
	})
	err := server.Start(context.Background())
	require.Error(t, err, "expected error when mock_enabled is false")
}

func TestNewMockServer_MockDisabled(t *testing.T) {
	t.Parallel()

	server := New(Options{
		Config: &config.Config{
			MockEnabled: false,
			Specs: []config.Spec{
				{
					Domain:   "test-api",
					LLMTitle: "Test API v1",
					BaseURL:  "https://api.example.com",
					Collections: []config.Collection{
						{
							LLMTitle:    "Main Collection",
							Location:    "https://example.com/spec.yaml",
							BaseMockURL: "localhost:8080",
						},
					},
				},
			},
		},
	})
	err := server.Start(context.Background())
	require.Error(t, err, "expected error when mock_enabled is false")
}

func TestExtractHostPort(t *testing.T) {
	t.Parallel()

	tests := []struct {
		addr string
		want string
	}{
		{"localhost:8080", "localhost:8080"},
		{"127.0.0.1:9000/v1/smev", "127.0.0.1:9000"},
		{"localhost:8080/api/v1", "localhost:8080"},
		{"0.0.0.0:3000/path/to/service", "0.0.0.0:3000"},
		{"127.0.0.1:80", "127.0.0.1:80"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.addr, func(t *testing.T) {
			t.Parallel()
			got := extractHostPort(tt.addr)
			assert.Equal(t, tt.want, got, "extractHostPort(%q)", tt.addr)
		})
	}
}
