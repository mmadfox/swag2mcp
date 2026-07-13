package commands

import (
	"strings"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/workspace"
)

func TestRunClean_EmptyWorkspace(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := workspace.New(tmpDir)
	if err := ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

	var buf strings.Builder
	err := runClean(tmpDir, &buf)
	if err != nil {
		t.Fatalf("runClean() = %v", err)
	}
	if !strings.Contains(buf.String(), "Removed contents") {
		t.Errorf("output = %q, want success message", buf.String())
	}
}

func TestRunClean_InvalidPath(t *testing.T) {
	var buf strings.Builder
	err := runClean("/nonexistent/path", &buf)
	if err != nil {
		t.Fatalf("runClean() = %v", err)
	}
	if !strings.Contains(buf.String(), "Removed contents") {
		t.Errorf("output = %q, want success message", buf.String())
	}
}
