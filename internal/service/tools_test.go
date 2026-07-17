package service

import (
	"fmt"
	"strings"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/model"
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

func TestLoadInstructionFromEmbed_Success(t *testing.T) {
	t.Parallel()

	content, err := loadInstructionFromEmbed()
	if err != nil {
		t.Fatalf("loadInstructionFromEmbed() = %v", err)
	}
	if content == "" {
		t.Error("instruction content is empty")
	}
}

func TestLoadToolFromEmbed_Success(t *testing.T) {
	t.Parallel()

	tool, err := loadToolFromEmbed("spec_list.md")
	if err != nil {
		t.Fatalf("loadToolFromEmbed() = %v", err)
	}
	if tool.Name != SpecList {
		t.Errorf("Name = %q, want %q", tool.Name, SpecList)
	}
	if tool.Description == "" {
		t.Error("Description is empty")
	}
}

func TestLoadToolFromEmbed_NotFound(t *testing.T) {
	t.Parallel()

	_, err := loadToolFromEmbed("nonexistent.md")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestLoadToolFromEmbed_AllToolsHaveName(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	defs, err := svc.MakeToolDefinitions()
	if err != nil {
		t.Fatalf("MakeToolDefinitions() = %v", err)
	}

	for _, tool := range defs.Tools {
		if tool.Name == "" {
			t.Error("found tool with empty name")
		}
		if tool.Description == "" {
			t.Errorf("tool %q has empty description", tool.Name)
		}
	}
}

func TestMakeAvaliablesSpecs_NoSpecs(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	result := svc.makeAvailableSpecs()
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

func TestMakeAvaliablesSpecs_WithSpecs(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	seedTestData(t, svc, t.Name())

	result := svc.makeAvailableSpecs()
	if result == "" {
		t.Fatal("expected non-empty result")
	}
	if !strings.Contains(result, "# Available specs") {
		t.Error("missing 'Available specs' header")
	}
	if !strings.Contains(result, "specID:") {
		t.Error("missing specID")
	}
	if !strings.Contains(result, "domain:") {
		t.Error("missing domain")
	}
}

func TestMakeAvaliablesSpecs_WithLLMInstruction(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	specInfo, _, _, _ := seedTestData(t, svc, t.Name())
	specInfo.LLMInstruction = "line1\nline2\nline3"

	result := svc.makeAvailableSpecs()
	if !strings.Contains(result, "instruction:") {
		t.Error("missing instruction field")
	}
}

func TestMakeAvaliablesSpecs_CollectionsBySpecError(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	specInfo, _, _, _ := seedTestData(t, svc, t.Name())

	svc.index.RemoveCollectionsBySpec(specInfo.ID)

	result := svc.makeAvailableSpecs()
	if !strings.Contains(result, "No available collections") {
		t.Error("expected 'No available collections' message")
	}
}

func TestMakeAvaliablesSpecs_MoreThan10Collections(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	specInfo, _, _, _ := seedTestData(t, svc, t.Name())

	// Add 11 collections to the spec
	for i := range 11 {
		collID := fmt.Sprintf("coll-%d", i)
		svc.index.AddCollection(&model.Collection{
			ID:     collID,
			SpecID: specInfo.ID,
			Title:  fmt.Sprintf("Collection %d", i),
		})
	}

	result := svc.makeAvailableSpecs()
	if !strings.Contains(result, "more than 10 collections") {
		t.Error("expected 'more than 10 collections' message")
	}
}

func TestMakeAvaliablesSpecs_CollectionLLMInstruction(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	_, collectionInfo, _, _ := seedTestData(t, svc, t.Name())
	collectionInfo.LLMInstruction = "custom instruction for this collection"

	result := svc.makeAvailableSpecs()
	if !strings.Contains(result, "custom instruction for this collection") {
		t.Error("expected collection LLMInstruction in output")
	}
}
