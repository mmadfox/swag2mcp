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
// If root is a relative path, it is resolved to an absolute path.
func New(root string) (*Workspace, error) {
	if root == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("cannot determine home directory: %w", err)
		}
		root = filepath.Join(home, DefaultRootName)
	} else {
		absRoot, err := filepath.Abs(root)
		if err != nil {
			return nil, fmt.Errorf("resolve path: %w", err)
		}
		root = absRoot
	}
	return &Workspace{root: root}, nil
}

// NewFromBase creates a Workspace rooted at base/.swag2mcp.
// If base is empty, it defaults to ~/.swag2mcp.
func NewFromBase(base string) (*Workspace, error) {
	if base == "" {
		return New("")
	}
	abs, err := filepath.Abs(base)
	if err != nil {
		return nil, fmt.Errorf("resolve path: %w", err)
	}
	return New(filepath.Join(abs, DefaultRootName))
}

// Init creates the workspace root and all standard subdirectories.
func (w *Workspace) Init() error {
	dirs := []string{
		w.root,
		w.CacheDir(),
		w.SpecsDir(),
		w.ResponsesDir(),
		w.AuthScriptsDir(),
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

// DefaultRoot returns the default workspace root path (~/.swag2mcp).
func DefaultRoot() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return DefaultRootName
	}
	return filepath.Join(home, DefaultRootName)
}

// DefaultConfigPath returns the default config file path (~/.swag2mcp/swag2mcp.yaml).
func DefaultConfigPath() string {
	return filepath.Join(DefaultRoot(), "swag2mcp.yaml")
}

// ConfigPathIn returns the config path inside a given workspace directory.
func ConfigPathIn(workspaceDir string) string {
	return filepath.Join(workspaceDir, "swag2mcp.yaml")
}

// ConfigPath returns the config file path inside this workspace.
func (w *Workspace) ConfigPath() string {
	return ConfigPathIn(w.root)
}

// ConfigExists checks whether the config file exists in this workspace.
func (w *Workspace) ConfigExists() bool {
	_, err := os.Stat(w.ConfigPath())
	return err == nil
}

// ConfigNotExists checks whether the config file does NOT exist in this workspace.
func (w *Workspace) ConfigNotExists() bool {
	return !w.ConfigExists()
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

// AuthScriptsDir returns the path to the auth scripts subdirectory.
func (w *Workspace) AuthScriptsDir() string {
	return w.Sub(DirAuthScripts)
}

// Clean removes all contents of cache/ and responses/ directories
// without removing the directories themselves.
func (w *Workspace) Clean() error {
	dirs := []string{
		w.CacheDir(),
		w.ResponsesDir(),
	}
	for _, dir := range dirs {
		if err := removeContents(dir); err != nil {
			return fmt.Errorf("clean %s: %w", filepath.Base(dir), err)
		}
	}
	return nil
}

// removeContents removes all files and subdirectories inside dir,
// but keeps the directory itself.
func removeContents(dir string) error {
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	for _, entry := range entries {
		p := filepath.Join(dir, entry.Name())
		if rErr := os.RemoveAll(p); rErr != nil {
			return rErr
		}
	}
	return nil
}
