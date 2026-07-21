package workspace

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"sync"
	"time"
)

// Workspace manages the workspace directory and its standard subdirectories.
type Workspace struct {
	mu   sync.RWMutex
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
		return &Workspace{root: filepath.Join(home, DefaultRootName)}, nil
	}
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("resolve path: %w", err)
	}
	return &Workspace{root: absRoot}, nil
}

// NewFromBase creates a Workspace rooted at the given base directory.
// If base is empty, it defaults to ~/.swag2mcp.
// If base is provided, it is used as the workspace root directly.
func NewFromBase(base string) (*Workspace, error) {
	if base == "" {
		return New("")
	}
	abs, err := filepath.Abs(base)
	if err != nil {
		return nil, fmt.Errorf("resolve path: %w", err)
	}
	return &Workspace{root: abs}, nil
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
	w.mu.RLock()
	defer w.mu.RUnlock()
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
	w.mu.RLock()
	defer w.mu.RUnlock()
	return ConfigPathIn(w.root)
}

// IsEmpty checks whether the workspace directory is empty or does not exist.
// Returns true if the directory does not exist, exists but is empty,
// or contains only swag2mcp.yaml (from a previous init).
func (w *Workspace) IsEmpty() (bool, error) {
	entries, err := os.ReadDir(w.root)
	if os.IsNotExist(err) {
		return true, nil
	}
	if err != nil {
		return false, fmt.Errorf("read directory %q: %w", w.root, err)
	}
	for _, entry := range entries {
		if entry.Name() == "swag2mcp.yaml" {
			continue
		}
		return false, nil
	}
	return true, nil
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

// CleanOldResponses removes response files older than maxAge from the responses directory.
func (w *Workspace) CleanOldResponses(maxAge time.Duration) error {
	dir := w.ResponsesDir()
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("read responses dir: %w", err)
	}

	now := time.Now()
	var errs []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, statErr := entry.Info()
		if statErr != nil {
			errs = append(errs, statErr.Error())
			continue
		}
		if now.Sub(info.ModTime()) > maxAge {
			p := filepath.Join(dir, entry.Name())
			if rmErr := os.Remove(p); rmErr != nil {
				errs = append(errs, rmErr.Error())
			}
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to remove some old response files: %s", strings.Join(errs, "; "))
	}
	return nil
}

// AuthScriptPath returns the path to the auth script for the given domain.
func (w *Workspace) AuthScriptPath(domain string) string {
	ext := ".sh"
	if runtime.GOOS == osWindows {
		ext = ".bat"
	}
	return filepath.Join(w.AuthScriptsDir(), domain+ext)
}

// EnsureAuthScript creates an auth script stub for the given domain if it does not exist.
func (w *Workspace) EnsureAuthScript(domain string) error {
	if err := os.MkdirAll(w.AuthScriptsDir(), 0750); err != nil {
		return fmt.Errorf("create auth_scripts dir: %w", err)
	}

	scriptPath := w.AuthScriptPath(domain)
	if _, err := os.Stat(scriptPath); err == nil {
		return nil
	}

	var content string
	if runtime.GOOS == osWindows {
		content = `@echo off
echo {"token": "your-token-here", "expires_in": 3600}
`
	} else {
		content = `#!/bin/sh
echo '{"token": "your-token-here", "expires_in": 3600}'
`
	}

	if err := os.WriteFile(scriptPath, []byte(content), 0600); err != nil {
		return fmt.Errorf("write auth script %s: %w", scriptPath, err)
	}

	return nil
}

// DownloadSpec downloads a spec file from a URL or reads it from a local path.
// Returns the raw file data.
func (w *Workspace) DownloadSpec(ctx context.Context, source string) ([]byte, error) {
	if source == "" {
		return nil, errors.New("source is empty")
	}

	isURL := strings.HasPrefix(source, "https://") || strings.HasPrefix(source, "http://")
	isFileURL := strings.HasPrefix(source, "file://")

	switch {
	case isFileURL:
		return w.downloadFromFileURL(source)
	case isURL:
		return w.downloadFromHTTP(ctx, source)
	default:
		return w.downloadFromLocalPath(source)
	}
}

// downloadFromFileURL reads a file from a file:// URL.
func (w *Workspace) downloadFromFileURL(rawURL string) ([]byte, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid file URL: %w", err)
	}
	path, err := url.PathUnescape(u.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to unescape URL path: %w", err)
	}
	if runtime.GOOS == osWindows && len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}
	data, err := os.ReadFile(filepath.FromSlash(path))
	if err != nil {
		return nil, fmt.Errorf("read file %q: %w", path, err)
	}
	return data, nil
}

// downloadFromHTTP downloads a spec file from an HTTP(S) URL.
func (w *Workspace) downloadFromHTTP(ctx context.Context, source string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, source, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("download %q: %w", source, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected HTTP status %d for %q", resp.StatusCode, source)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}
	return data, nil
}

// downloadFromLocalPath reads a spec file from a local filesystem path.
func (w *Workspace) downloadFromLocalPath(source string) ([]byte, error) {
	path := source
	if !filepath.IsAbs(path) {
		abs, err := filepath.Abs(path)
		if err != nil {
			return nil, fmt.Errorf("resolve path: %w", err)
		}
		path = abs
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file %q: %w", path, err)
	}
	return data, nil
}

// SpecPath returns the full path to a spec file in the specs/ directory.
func (w *Workspace) SpecPath(name string) string {
	return filepath.Join(w.SpecsDir(), name)
}

// ListSpecs returns the filenames of all files in the specs/ directory.
func (w *Workspace) ListSpecs() ([]string, error) {
	entries, err := os.ReadDir(w.SpecsDir())
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read specs dir: %w", err)
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		names = append(names, entry.Name())
	}
	return names, nil
}

// specNameExists checks if a filename already exists in the specs/ directory.
func (w *Workspace) specNameExists(name string) (bool, error) {
	names, err := w.ListSpecs()
	if err != nil {
		return false, err
	}
	return slices.Contains(names, name), nil
}

// SaveSpec saves spec data to the specs/ directory with the given name.
// Returns an error if a file with that name already exists.
func (w *Workspace) SaveSpec(name string, data []byte) (string, error) {
	if name == "" {
		return "", errors.New("name is empty")
	}
	if len(data) == 0 {
		return "", errors.New("data is empty")
	}

	if err := os.MkdirAll(w.SpecsDir(), 0750); err != nil {
		return "", fmt.Errorf("create specs dir: %w", err)
	}

	exists, existsErr := w.specNameExists(name)
	if existsErr != nil {
		return "", existsErr
	}
	if exists {
		return "", fmt.Errorf("spec file %q already exists in %s", name, w.SpecsDir())
	}

	path := w.SpecPath(name)
	if writeErr := os.WriteFile(filepath.Clean(path), data, 0600); writeErr != nil {
		return "", fmt.Errorf("write spec file %q: %w", path, writeErr)
	}
	return path, nil
}

// RemoveOrphanAuthScripts removes auth script files for domains not in the active list.
func (w *Workspace) RemoveOrphanAuthScripts(activeDomains []string) error {
	active := make(map[string]bool, len(activeDomains))
	for _, d := range activeDomains {
		active[d] = true
	}

	entries, err := os.ReadDir(w.AuthScriptsDir())
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("read auth_scripts dir: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		domain := strings.TrimSuffix(name, ".sh")
		domain = strings.TrimSuffix(domain, ".bat")
		if domain == name {
			continue
		}
		if !active[domain] {
			p := filepath.Join(w.AuthScriptsDir(), name)
			if rErr := os.Remove(p); rErr != nil {
				return fmt.Errorf("remove orphan auth script %s: %w", name, rErr)
			}
		}
	}

	return nil
}
