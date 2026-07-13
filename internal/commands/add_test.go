package commands

import (
	"testing"
)

func TestRunAddSpec_Example(t *testing.T) {
	err := runAddSpec("", "", true)
	if err != nil {
		t.Fatalf("runAddSpec(example=true) = %v", err)
	}
}

func TestRunAddCollection_Example(t *testing.T) {
	err := runAddCollection("", "", true)
	if err != nil {
		t.Fatalf("runAddCollection(example=true) = %v", err)
	}
}
