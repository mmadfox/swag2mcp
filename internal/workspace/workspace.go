package workspace

import (
	"fmt"
	"os"
	"path/filepath"
)

// Workspace manages the workspace directory and its standard subdirectories.
type Workspace struct {
	root string
}

// New creates a Workspace rooted at the given directory.
// If root is empty, it defaults to ~/.swag2mcp.
func New(root string) (*Workspace, error) {
	if root == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("cannot determine home directory: %w", err)
		}
		root = filepath.Join(home, DefaultRootName)
	}
	return &Workspace{root: root}, nil
}

// Init creates the workspace root and all standard subdirectories.
func (w *Workspace) Init() error {
	dirs := []string{
		w.root,
		w.CacheDir(),
		w.SpecsDir(),
		w.ResponsesDir(),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0750); err != nil {
			return fmt.Errorf("failed to create directory %q: %w", dir, err)
		}
	}
	return nil
}

// Root returns the absolute path to the workspace root directory.
func (w *Workspace) Root() string {
	return w.root
}

// Sub returns the path to a named subdirectory inside the workspace.
func (w *Workspace) Sub(name string) string {
	return filepath.Join(w.root, name)
}

// CacheDir returns the path to the cache subdirectory.
func (w *Workspace) CacheDir() string {
	return w.Sub(DirCache)
}

// SpecsDir returns the path to the specs subdirectory.
func (w *Workspace) SpecsDir() string {
	return w.Sub(DirSpecs)
}

// ResponsesDir returns the path to the responses subdirectory.
func (w *Workspace) ResponsesDir() string {
	return w.Sub(DirResponses)
}
