package mockserver

import (
	"context"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/config"
)

func TestNewMockServer_NoServers(t *testing.T) {
	t.Parallel()

	server := New(Options{
		Config: &config.Config{},
	})
	err := server.Start(context.Background())
	if err == nil {
		t.Error("expected error when mock_enabled is false")
	}
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
	if err == nil {
		t.Error("expected error when mock_enabled is false")
	}
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
			if got != tt.want {
				t.Errorf("extractHostPort(%q) = %q, want %q", tt.addr, got, tt.want)
			}
		})
	}
}
