package workspace

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
	entries, readErr := os.ReadDir(w.AuthScriptsDir())
	if os.IsNotExist(readErr) {
		return nil
	}
	if readErr != nil {
		return fmt.Errorf("read auth_scripts dir: %w", readErr)
	}

	authDir := filepath.Join(exportDir, DirAuthScripts)
	if mkdirErr := os.MkdirAll(authDir, 0750); mkdirErr != nil {
		return fmt.Errorf("create auth_scripts dir in export: %w", mkdirErr)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		src := filepath.Join(w.AuthScriptsDir(), entry.Name())
		dst := filepath.Join(authDir, entry.Name())
		data, readFileErr := os.ReadFile(src)
		if readFileErr != nil {
			return fmt.Errorf("read auth script %s: %w", entry.Name(), readFileErr)
		}
		if writeErr := os.WriteFile(filepath.Clean(dst), data, 0600); writeErr != nil {
			return fmt.Errorf("write auth script %s: %w", entry.Name(), writeErr)
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
	data, marshalErr := json.Marshal(meta)
	if marshalErr != nil {
		return fmt.Errorf("marshal meta: %w", marshalErr)
	}
	path := filepath.Join(exportDir, MetaFileName)
	if writeErr := os.WriteFile(filepath.Clean(path), data, 0600); writeErr != nil {
		return fmt.Errorf("write meta file: %w", writeErr)
	}
	return nil
}

// CreateZip creates a ZIP archive from the source directory at the output path.
func CreateZip(sourceDir, outputPath string) error {
	outputPath = ensureZipExt(outputPath)

	f, createErr := os.Create(filepath.Clean(outputPath))
	if createErr != nil {
		return fmt.Errorf("create zip file: %w", createErr)
	}
	defer f.Close()

	zw := zip.NewWriter(f)
	defer zw.Close()

	walkErr := filepath.Walk(sourceDir, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if path == sourceDir {
			return nil
		}

		rel, relErr := filepath.Rel(sourceDir, path)
		if relErr != nil {
			return fmt.Errorf("relative path: %w", relErr)
		}

		if info.IsDir() {
			_, zipErr := zw.Create(rel + "/")
			return zipErr
		}

		w, zipErr := zw.Create(rel)
		if zipErr != nil {
			return fmt.Errorf("create zip entry %s: %w", rel, zipErr)
		}

		data, readErr := os.ReadFile(filepath.Clean(path))
		if readErr != nil {
			return fmt.Errorf("read %s: %w", path, readErr)
		}

		if _, writeErr := w.Write(data); writeErr != nil {
			return fmt.Errorf("write %s to zip: %w", rel, writeErr)
		}
		return nil
	})
	if walkErr != nil {
		return fmt.Errorf("walk source dir: %w", walkErr)
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
			rc, openErr := f.Open()
			if openErr != nil {
				return false, fmt.Errorf("open meta file in zip: %w", openErr)
			}
			defer rc.Close()

			data, readErr := io.ReadAll(rc)
			if readErr != nil {
				return false, fmt.Errorf("read meta file: %w", readErr)
			}

			var meta Meta
			if unmarshalErr := json.Unmarshal(data, &meta); unmarshalErr != nil {
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

	for _, f := range r.File {
		fpath := filepath.Join(destDir, f.Name)

		if !strings.HasPrefix(filepath.Clean(fpath), filepath.Clean(destDir)+string(filepath.Separator)) {
			return fmt.Errorf("illegal file path in zip: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			if mkdirErr := os.MkdirAll(fpath, 0750); mkdirErr != nil {
				return fmt.Errorf("create dir %s: %w", fpath, mkdirErr)
			}
			continue
		}

		if mkdirErr := os.MkdirAll(filepath.Dir(fpath), 0750); mkdirErr != nil {
			return fmt.Errorf("create parent dir for %s: %w", fpath, mkdirErr)
		}

		rc, openErr := f.Open()
		if openErr != nil {
			return fmt.Errorf("open zip entry %s: %w", f.Name, openErr)
		}

		data, readErr := io.ReadAll(rc)
		rc.Close()
		if readErr != nil {
			return fmt.Errorf("read zip entry %s: %w", f.Name, readErr)
		}

		if writeErr := os.WriteFile(filepath.Clean(fpath), data, 0600); writeErr != nil {
			return fmt.Errorf("write %s: %w", fpath, writeErr)
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
			rc, openErr := f.Open()
			if openErr != nil {
				return nil, fmt.Errorf("open meta file: %w", openErr)
			}
			defer rc.Close()

			data, readErr := io.ReadAll(rc)
			if readErr != nil {
				return nil, fmt.Errorf("read meta file: %w", readErr)
			}

			var meta Meta
			if unmarshalErr := json.Unmarshal(data, &meta); unmarshalErr != nil {
				return nil, fmt.Errorf("unmarshal meta: %w", unmarshalErr)
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
func (w *Workspace) CopyConfigToExport(exportDir string) error {
	data, err := os.ReadFile(w.ConfigPath())
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}
	dst := filepath.Join(exportDir, "swag2mcp.yaml")
	if writeErr := os.WriteFile(filepath.Clean(dst), data, 0600); writeErr != nil {
		return fmt.Errorf("write config to export: %w", writeErr)
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
		data, readErr := os.ReadFile(src)
		if readErr != nil {
			return fmt.Errorf("read spec %s: %w", entry.Name(), readErr)
		}
		if _, saveErr := w.SaveSpec(entry.Name(), data); saveErr != nil {
			return fmt.Errorf("save spec %s: %w", entry.Name(), saveErr)
		}
	}
	return nil
}

// CopyAuthScriptsToWorkspace copies all auth scripts from the export to the workspace.
func (w *Workspace) CopyAuthScriptsToWorkspace(exportDir string) error {
	srcDir := filepath.Join(exportDir, DirAuthScripts)
	entries, readErr := os.ReadDir(srcDir)
	if os.IsNotExist(readErr) {
		return nil
	}
	if readErr != nil {
		return fmt.Errorf("read export auth_scripts dir: %w", readErr)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		src := filepath.Join(srcDir, entry.Name())
		data, readFileErr := os.ReadFile(src)
		if readFileErr != nil {
			return fmt.Errorf("read auth script %s: %w", entry.Name(), readFileErr)
		}
		dst := filepath.Join(w.AuthScriptsDir(), entry.Name())
		if writeErr := os.WriteFile(filepath.Clean(dst), data, 0600); writeErr != nil {
			return fmt.Errorf("write auth script %s: %w", entry.Name(), writeErr)
		}
	}
	return nil
}
