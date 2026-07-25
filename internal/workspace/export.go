package workspace

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	// MetaFileName is the marker file inside a swag2mcp backup ZIP.
	MetaFileName = "swag2mcp.meta"

	// MetaType is the type value in the meta file.
	MetaType = "swag2mcp-backup"

	zipExt = ".zip"
)

// Meta holds metadata about a swag2mcp backup archive.
type Meta struct {
	Type    string `json:"type"`
	Version string `json:"version"`
	Created string `json:"created"`
}

// CreateExportDir creates a temporary directory for building an export archive.
func (w *Workspace) CreateExportDir() (string, error) {
	dir, err := os.MkdirTemp("", "swag2mcp-export-*")
	if err != nil {
		return "", fmt.Errorf("create temp dir: %w", err)
	}
	return dir, nil
}

// WriteSpecToExport writes spec data to the export directory under specs/.
func WriteSpecToExport(exportDir, name string, data []byte) error {
	specsDir := filepath.Join(exportDir, DirSpecs)
	if err := os.MkdirAll(specsDir, 0750); err != nil {
		return fmt.Errorf("create specs dir in export: %w", err)
	}
	path := filepath.Join(specsDir, name)
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("write spec to export: %w", err)
	}
	return nil
}

// CopyAuthScriptsToExport copies all auth scripts from the workspace to the export directory.
func (w *Workspace) CopyAuthScriptsToExport(exportDir string) error {
	entries, err := os.ReadDir(w.AuthScriptsDir())
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("read auth_scripts dir: %w", err)
	}

	authDir := filepath.Join(exportDir, DirAuthScripts)
	if err := os.MkdirAll(authDir, 0750); err != nil {
		return fmt.Errorf("create auth_scripts dir in export: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		src := filepath.Join(w.AuthScriptsDir(), entry.Name())
		dst := filepath.Join(authDir, entry.Name())
		data, err := os.ReadFile(src)
		if err != nil {
			return fmt.Errorf("read auth script %s: %w", entry.Name(), err)
		}
		if err := os.WriteFile(filepath.Clean(dst), data, 0600); err != nil {
			return fmt.Errorf("write auth script %s: %w", entry.Name(), err)
		}
	}
	return nil
}

// CreateMetaFile creates the swag2mcp.meta marker file in the export directory.
func CreateMetaFile(exportDir, version string) error {
	meta := Meta{
		Type:    MetaType,
		Version: version,
		Created: time.Now().UTC().Format(time.RFC3339),
	}
	data, err := json.Marshal(meta)
	if err != nil {
		return fmt.Errorf("marshal meta: %w", err)
	}
	path := filepath.Join(exportDir, MetaFileName)
	if err := os.WriteFile(filepath.Clean(path), data, 0600); err != nil {
		return fmt.Errorf("write meta file: %w", err)
	}
	return nil
}

// CreateZip creates a ZIP archive from the source directory at the output path.
func CreateZip(sourceDir, outputPath string) error {
	outputPath = ensureZipExt(outputPath)

	f, err := os.Create(filepath.Clean(outputPath))
	if err != nil {
		return fmt.Errorf("create zip file: %w", err)
	}
	defer f.Close()

	zw := zip.NewWriter(f)
	defer zw.Close()

	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == sourceDir {
			return nil
		}

		rel, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return fmt.Errorf("relative path: %w", err)
		}

		if info.IsDir() {
			_, err := zw.Create(rel + "/")
			return err
		}

		w, err := zw.Create(rel)
		if err != nil {
			return fmt.Errorf("create zip entry %s: %w", rel, err)
		}

		data, err := os.ReadFile(filepath.Clean(path))
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}

		if _, err := w.Write(data); err != nil {
			return fmt.Errorf("write %s to zip: %w", rel, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("walk source dir: %w", err)
	}

	return nil
}

// ValidateZip checks whether the given ZIP file is a valid swag2mcp backup archive.
// It reads the swag2mcp.meta file inside the archive without extracting.
func ValidateZip(path string) (bool, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return false, fmt.Errorf("open zip: %w", err)
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name == MetaFileName {
			rc, err := f.Open()
			if err != nil {
				return false, fmt.Errorf("open meta file in zip: %w", err)
			}
			defer rc.Close()

			data, err := io.ReadAll(rc)
			if err != nil {
				return false, fmt.Errorf("read meta file: %w", err)
			}

			var meta Meta
			if err := json.Unmarshal(data, &meta); err != nil {
				return false, nil
			}
			return meta.Type == MetaType, nil
		}
	}
	return false, nil
}

// ExtractZip extracts a ZIP archive to the destination directory.
func ExtractZip(path, destDir string) error {
	r, err := zip.OpenReader(path)
	if err != nil {
		return fmt.Errorf("open zip: %w", err)
	}
	defer r.Close()

	destDir = filepath.Clean(destDir)

	for _, f := range r.File {
		fpath := filepath.Join(destDir, f.Name)

		rel, err := filepath.Rel(destDir, fpath)
		if err != nil || strings.HasPrefix(rel, "..") {
			return fmt.Errorf("zip slip detected: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, 0750); err != nil {
				return fmt.Errorf("create dir %s: %w", fpath, err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), 0750); err != nil {
			return fmt.Errorf("create parent dir for %s: %w", fpath, err)
		}

		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf("open zip entry %s: %w", f.Name, err)
		}

		data, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return fmt.Errorf("read zip entry %s: %w", f.Name, err)
		}

		if err := os.WriteFile(filepath.Clean(fpath), data, 0600); err != nil {
			return fmt.Errorf("write %s: %w", fpath, err)
		}
	}
	return nil
}

// ensureZipExt adds .zip extension if not present.
func ensureZipExt(path string) string {
	if filepath.Ext(path) != zipExt {
		return path + zipExt
	}
	return path
}

// DefaultExportName returns the default ZIP filename with timestamp.
func DefaultExportName() string {
	return fmt.Sprintf("swag2mcp-backup-%s.zip", time.Now().UTC().Format("2006-01-02-150405"))
}

// IsSwag2mcpZip checks if a file is a swag2mcp backup ZIP without opening it.
// Returns true if the file has .zip extension and contains swag2mcp.meta.
func IsSwag2mcpZip(path string) bool {
	if filepath.Ext(path) != ".zip" {
		return false
	}
	valid, err := ValidateZip(path)
	return err == nil && valid
}

// ReadMetaFromZip reads the meta file from a swag2mcp backup ZIP.
func ReadMetaFromZip(path string) (*Meta, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, fmt.Errorf("open zip: %w", err)
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name == MetaFileName {
			rc, err := f.Open()
			if err != nil {
				return nil, fmt.Errorf("open meta file: %w", err)
			}
			defer rc.Close()

			data, err := io.ReadAll(rc)
			if err != nil {
				return nil, fmt.Errorf("read meta file: %w", err)
			}

			var meta Meta
			if err := json.Unmarshal(data, &meta); err != nil {
				return nil, fmt.Errorf("unmarshal meta: %w", err)
			}
			return &meta, nil
		}
	}
	return nil, errors.New("meta file not found in zip")
}

// CreateEmptyDirsInExport creates empty cache/ and responses/ directories in the export.
func CreateEmptyDirsInExport(exportDir string) error {
	for _, dir := range []string{DirCache, DirResponses} {
		path := filepath.Join(exportDir, dir)
		if err := os.MkdirAll(path, 0750); err != nil {
			return fmt.Errorf("create %s dir in export: %w", dir, err)
		}
	}
	return nil
}

// CopyConfigToExport copies the config file to the export directory.
func (w *Workspace) copyConfigToExport(exportDir string) error {
	data, err := os.ReadFile(w.ConfigPath())
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}
	dst := filepath.Join(exportDir, "swag2mcp.yaml")
	if err := os.WriteFile(filepath.Clean(dst), data, 0600); err != nil {
		return fmt.Errorf("write config to export: %w", err)
	}
	return nil
}

// ReadConfigFromExport reads the config from an exported workspace directory.
func ReadConfigFromExport(exportDir string) ([]byte, error) {
	path := filepath.Join(exportDir, "swag2mcp.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config from export: %w", err)
	}
	return data, nil
}

// CopySpecsToWorkspace copies all spec files from the export specs/ to the workspace specs/.
func (w *Workspace) CopySpecsToWorkspace(exportDir string) error {
	srcDir := filepath.Join(exportDir, DirSpecs)
	entries, err := os.ReadDir(srcDir)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("read export specs dir: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		src := filepath.Join(srcDir, entry.Name())
		data, err := os.ReadFile(src)
		if err != nil {
			return fmt.Errorf("read spec %s: %w", entry.Name(), err)
		}
		if _, err := w.SaveSpec(entry.Name(), data); err != nil {
			return fmt.Errorf("save spec %s: %w", entry.Name(), err)
		}
	}
	return nil
}

// CopyAuthScriptsToWorkspace copies all auth scripts from the export to the workspace.
func (w *Workspace) CopyAuthScriptsToWorkspace(exportDir string) error {
	srcDir := filepath.Join(exportDir, DirAuthScripts)
	entries, err := os.ReadDir(srcDir)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("read export auth_scripts dir: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		src := filepath.Join(srcDir, entry.Name())
		data, err := os.ReadFile(src)
		if err != nil {
			return fmt.Errorf("read auth script %s: %w", entry.Name(), err)
		}
		dst := filepath.Join(w.AuthScriptsDir(), entry.Name())
		if err := os.WriteFile(filepath.Clean(dst), data, 0600); err != nil {
			return fmt.Errorf("write auth script %s: %w", entry.Name(), err)
		}
	}
	return nil
}
