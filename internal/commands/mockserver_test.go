package commands

import (
	"context"
	"testing"
)

func TestRunMockServer_NoConfig(t *testing.T) {
	tmpDir := t.TempDir()
	opts := &mockServerCmdOptions{}
	err := runMockServer(tmpDir, opts, context.Background())
	if err == nil {
		t.Fatal("runMockServer() expected error, got nil")
	}
}
