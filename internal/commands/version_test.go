package commands

import (
	"strings"
	"testing"
)

func TestNewVersionCmd(t *testing.T) {
	cmd := newVersionCmd()
	if cmd == nil {
		t.Fatal("newVersionCmd() returned nil")
	}
	if cmd.Use != "version" {
		t.Errorf("Use = %q, want %q", cmd.Use, "version")
	}
}

func TestNewVersionCmd_Output(t *testing.T) {
	Version = "v1.0.0"
	cmd := newVersionCmd()

	buf := new(strings.Builder)
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() = %v", err)
	}

	output := buf.String()
	if output != "swag2mcp v1.0.0\n" {
		t.Errorf("output = %q, want %q", output, "swag2mcp v1.0.0\n")
	}
}
