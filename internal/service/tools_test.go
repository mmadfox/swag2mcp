package service

import (
	"testing"
)

func TestMakeToolDefinitions_IncludesAuthByDefault(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	defs, err := svc.MakeToolDefinitions()
	if err != nil {
		t.Fatalf("MakeToolDefinitions() = %v", err)
	}

	found := false
	for _, tool := range defs.Tools {
		if tool.Name == Auth {
			found = true
			break
		}
	}
	if !found {
		t.Error("auth tool not found in definitions, expected it to be present")
	}
}

func TestMakeToolDefinitions_ExcludesAuthWhenDisabled(t *testing.T) {
	t.Parallel()

	svc := newTestService(t, WithDisableLLMAuth(true))
	defs, err := svc.MakeToolDefinitions()
	if err != nil {
		t.Fatalf("MakeToolDefinitions() = %v", err)
	}

	for _, tool := range defs.Tools {
		if tool.Name == Auth {
			t.Error("auth tool found in definitions, expected it to be excluded")
		}
	}
}

func TestMakeToolDefinitions_InstructionNotEmpty(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	defs, err := svc.MakeToolDefinitions()
	if err != nil {
		t.Fatalf("MakeToolDefinitions() = %v", err)
	}
	if defs.Instruction == "" {
		t.Error("Instruction is empty")
	}
}

func TestMakeToolDefinitions_HasTools(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	defs, err := svc.MakeToolDefinitions()
	if err != nil {
		t.Fatalf("MakeToolDefinitions() = %v", err)
	}
	if len(defs.Tools) == 0 {
		t.Error("Tools is empty")
	}
}
