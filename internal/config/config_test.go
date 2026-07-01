package config

import "testing"

func TestTest(t *testing.T) {
	t.Parallel()
	conf := &Config{
		WorkspaceDir: "/tmp/swag2mcp-test",
		Specs: []Spec{
			{
				Domain:   "DZZ-12",
				LLMTitle: "DZZ-12 API - comprehensive integration platform for enterprise resource planning and customer relationship management",
				BaseURL:  "https://api.dzz-12.example.com",
				Collections: []Collection{
					{
						Location: "https://api.dzz-12.example.com/openapi.json",
					},
				},
			},
		},
	}
	f := NewFilter(nil)
	if err := conf.Validate(f); err != nil {
		t.Fatal(err)
	}
}
