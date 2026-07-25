package tui

import (
	"os"
	"testing"
)

// withStdinPipe replaces [os.Stdin] with a pipe for testing.
func withStdinPipe(t *testing.T, input string) {
	t.Helper()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Pipe() = %v", err)
	}
	t.Cleanup(func() { r.Close(); w.Close() })

	oldStdin := os.Stdin
	os.Stdin = r                              //nolint:reassign // test helper temporarily replaces os.Stdin to mock user input
	t.Cleanup(func() { os.Stdin = oldStdin }) //nolint:reassign // restore original stdin

	w.WriteString(input)
	w.Close()
}

func TestPromptSelection_Valid(t *testing.T) {
	withStdinPipe(t, "2\n")

	got, err := promptSelection("test", 3)
	if err != nil {
		t.Fatalf("promptSelection() = %v", err)
	}
	if got != 2 {
		t.Errorf("got %d, want 2", got)
	}
}

func TestPromptSelection_Invalid(t *testing.T) {
	withStdinPipe(t, "0\n")

	_, err := promptSelection("test", 3)
	if err == nil {
		t.Fatal("promptSelection() expected error for 0, got nil")
	}
}

func TestPromptSelection_OutOfRange(t *testing.T) {
	withStdinPipe(t, "99\n")

	_, err := promptSelection("test", 3)
	if err == nil {
		t.Fatal("promptSelection() expected error for out of range, got nil")
	}
}

func TestConfirmAction_Yes(t *testing.T) {
	withStdinPipe(t, "y\n")

	if !confirmAction("test?") {
		t.Error("confirmAction() = false, want true")
	}
}

func TestConfirmAction_No(t *testing.T) {
	withStdinPipe(t, "n\n")

	if confirmAction("test?") {
		t.Error("confirmAction() = true, want false")
	}
}
