package config

import "testing"

func TestTest(t *testing.T) {
	conf := &Config{
		Specs: []Spec{
			{
				Domain: "DZZ-12",
				Collections: []Collection{
					{},
				},
			},
		},
	}
	f := NewFilter(nil)
	if err := conf.Validate(f); err != nil {
		t.Fatal(err)
	}
}
