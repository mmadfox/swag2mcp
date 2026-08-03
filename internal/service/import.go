/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/workspace"
)

// ImportRequest is the request for the Import method.
type ImportRequest struct {
	Source       string
	Name         string
	SpecFilter   []string
	ConfFilePath string
	ZipSource    string
	Force        bool
}

// ImportResponse holds the result of an import operation.
type ImportResponse struct {
	Files []ImportedFile `json:"files"`
}

// ImportedFile represents a single imported spec file.
type ImportedFile struct {
	Source    string `json:"source"`
	Name      string `json:"name"`
	SavedPath string `json:"savedPath"`
	Skipped   bool   `json:"skipped,omitempty"`
}

type importService struct {
	ws WorkspaceOps
}

func newImportService(ws WorkspaceOps) *importService {
	return &importService{ws: ws}
}

// Import imports spec files into the workspace specs/ directory.
func (ims *importService) Import(ctx context.Context, req ImportRequest) (ImportResponse, error) {
	if req.ZipSource != "" {
		return ims.importFromZip(ctx, req)
	}

	if req.Source == "" {
		return ims.importSpecs(ctx, req)
	}

	if req.Name == "" {
		derived := specFileNameBase(req.Source)
		if derived == "" || derived == defaultSpecName {
			return ImportResponse{}, NewImportSourceError(
				errors.New("filename not found in URL"),
			)
		}
		req.Name = derived
	}

	return ims.importSingle(ctx, req)
}

func (ims *importService) importFromZip(_ context.Context, req ImportRequest) (ImportResponse, error) {
	if !workspace.IsSwag2mcpZip(req.ZipSource) {
		return ImportResponse{}, NewImportError(
			fmt.Sprintf("File %q is not a valid swag2mcp backup archive.", req.ZipSource),
			fmt.Errorf("invalid swag2mcp zip: %s", req.ZipSource),
		)
	}

	extractDir, dirErr := os.MkdirTemp("", "swag2mcp-restore-*")
	if dirErr != nil {
		return ImportResponse{}, NewWorkspaceError(
			"Failed to create temporary directory for extraction.",
			dirErr,
		)
	}
	defer os.RemoveAll(extractDir)

	if extractErr := workspace.ExtractZip(req.ZipSource, extractDir); extractErr != nil {
		return ImportResponse{}, NewWorkspaceError(
			fmt.Sprintf("Failed to extract archive %q.", req.ZipSource),
			extractErr,
		)
	}

	if initErr := ims.ws.Init(); initErr != nil {
		return ImportResponse{}, NewWorkspaceError(
			"Failed to initialize workspace directories.",
			initErr,
		)
	}

	if copyErr := ims.ws.CopySpecsToWorkspace(extractDir); copyErr != nil {
		return ImportResponse{}, NewWorkspaceError(
			"Failed to copy spec files from backup to workspace.",
			copyErr,
		)
	}

	if copyErr := ims.ws.CopyAuthScriptsToWorkspace(extractDir); copyErr != nil {
		return ImportResponse{}, NewWorkspaceError(
			"Failed to copy auth scripts from backup to workspace.",
			copyErr,
		)
	}

	cfgData, cfgReadErr := workspace.ReadConfigFromExport(extractDir)
	if cfgReadErr != nil {
		return ImportResponse{}, NewConfigError(
			"Failed to read configuration from backup.",
			cfgReadErr,
		)
	}

	cfgPath := ims.ws.ConfigPath()
	if writeErr := os.WriteFile(cfgPath, cfgData, 0600); writeErr != nil {
		return ImportResponse{}, NewConfigError(
			fmt.Sprintf("Failed to write configuration to %q.", cfgPath),
			writeErr,
		)
	}

	specs, listErr := ims.ws.ListSpecs()
	if listErr != nil {
		return ImportResponse{}, NewWorkspaceError(
			"Failed to list imported spec files.",
			listErr,
		)
	}

	files := make([]ImportedFile, 0, len(specs))
	for _, name := range specs {
		files = append(files, ImportedFile{
			Source:    req.ZipSource,
			Name:      name,
			SavedPath: ims.ws.SpecPath(name),
		})
	}

	return ImportResponse{
		Files: files,
	}, nil
}

func (ims *importService) importSingle(ctx context.Context, req ImportRequest) (ImportResponse, error) {
	data, err := ims.ws.DownloadSpec(ctx, req.Source)
	if err != nil {
		return ImportResponse{}, NewImportError(
			fmt.Sprintf("Failed to download spec from %q.", req.Source),
			err,
		)
	}

	name := req.Name
	pathPart := pathPartFromLocation(req.Source)
	if pathPart != "" {
		ext := filepath.Ext(name)
		base := strings.TrimSuffix(name, ext)
		name = base + "-" + pathPart + ext
	}

	if filepath.Ext(name) == "" {
		if ext := extFromLocation(req.Source); ext != "" {
			name += ext
		}
	}

	var path string
	if req.Force {
		path, err = ims.ws.SaveOrUpdateSpec(name, data)
	} else {
		path, err = ims.ws.SaveSpec(name, data)
	}
	if err != nil {
		return ImportResponse{}, NewImportError(
			fmt.Sprintf("Failed to save spec as %q. The filename may already exist.", name),
			err,
		)
	}

	return ImportResponse{
		Files: []ImportedFile{
			{
				Source:    req.Source,
				Name:      name,
				SavedPath: path,
			},
		},
	}, nil
}

func (ims *importService) importSpecs(ctx context.Context, req ImportRequest) (ImportResponse, error) {
	if req.ConfFilePath == "" {
		return ImportResponse{}, NewImportError(
			"Configuration file path is required for bulk import.",
			errors.New("config file path is empty"),
		)
	}

	cfg, err := config.Load(req.ConfFilePath)
	if err != nil {
		return ImportResponse{}, NewConfigError(
			fmt.Sprintf("Failed to load configuration from %q.", req.ConfFilePath),
			err,
		)
	}

	filter := makeFilter(req.SpecFilter)
	var imported []ImportedFile
	updated := false

	if len(req.SpecFilter) > 0 {
		active := make(map[string]struct{}, len(cfg.Specs))
		for i := range cfg.Specs {
			if !cfg.Specs[i].Disable {
				active[cfg.Specs[i].Domain] = struct{}{}
			}
		}
		for _, d := range req.SpecFilter {
			d = strings.TrimSpace(d)
			if _, ok := active[d]; !ok {
				return ImportResponse{}, NewImportSpecNotFoundError(d)
			}
		}
	}

	for i := range cfg.Specs {
		spec := &cfg.Specs[i]
		if spec.Disable {
			continue
		}
		if !filter.match(spec.Domain) {
			continue
		}

		for j := range spec.Collections {
			coll := &spec.Collections[j]
			if coll.Disable {
				continue
			}
			var importErr error
			imported, updated, importErr = ims.importCollection(ctx, spec.Domain, coll, req.ConfFilePath, imported, updated, j)
			if importErr != nil {
				return ImportResponse{}, importErr
			}
		}
	}

	if !updated {
		if len(req.SpecFilter) == 0 {
			return ImportResponse{}, NewImportNoMatchError("no active specs found in config")
		}
		return ImportResponse{}, NewImportNoMatchError(fmt.Sprintf("%v", req.SpecFilter))
	}

	if err := config.Save(cfg, req.ConfFilePath); err != nil {
		return ImportResponse{}, NewConfigError(
			fmt.Sprintf("Failed to save updated configuration to %q.", req.ConfFilePath),
			err,
		)
	}

	return ImportResponse{Files: imported}, nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (ims *importService) importCollection(
	ctx context.Context,
	domain string,
	coll *config.Collection,
	confFilePath string,
	imported []ImportedFile,
	updated bool,
	collIndex int,
) ([]ImportedFile, bool, error) {
	if isLocalSpecPath(coll.Location, ims.ws.SpecsDir(), filepath.Dir(confFilePath)) {
		absPath := coll.Location
		if !filepath.IsAbs(absPath) {
			absPath = filepath.Join(filepath.Dir(confFilePath), absPath)
		}
		if !fileExists(absPath) {
			return imported, updated, NewImportError(
				fmt.Sprintf("Spec file %q not found. The location points to a "+
					"local spec file that does not exist. Make sure the file "+
					"exists or update the location to a remote URL.", coll.Location),
				fmt.Errorf("spec file not found: %s", coll.Location),
			)
		}
		imported = append(imported, ImportedFile{
			Source:    coll.Location,
			Name:      filepath.Base(coll.Location),
			SavedPath: coll.Location,
			Skipped:   true,
		})
		return imported, true, nil
	}

	name := specFileName(domain, collTitle(coll, collIndex), coll.Location, pathPartFromLocation(coll.Location))
	specPath := ims.ws.SpecPath(name)
	exists := fileExists(specPath)

	data, err := ims.ws.DownloadSpec(ctx, coll.Location)
	if err != nil {
		if exists {
			coll.Location = filepath.Join("specs", name)
			imported = append(imported, ImportedFile{
				Source:    coll.Location,
				Name:      name,
				SavedPath: specPath,
				Skipped:   true,
			})
			return imported, true, nil
		}
		return imported, updated, NewImportError(
			fmt.Sprintf("Failed to download spec for collection %q.", coll.Title),
			err,
		)
	}

	sp, err := ims.ws.SaveOrUpdateSpec(name, data)
	if err != nil {
		return imported, updated, NewImportError(
			fmt.Sprintf("Failed to save spec as %q.", name),
			err,
		)
	}

	coll.Location = filepath.Join("specs", name)
	imported = append(imported, ImportedFile{
		Source:    coll.Location,
		Name:      name,
		SavedPath: sp,
	})
	return imported, true, nil
}

func collTitle(coll *config.Collection, index int) string {
	if coll.LLMTitle != "" {
		return coll.LLMTitle
	}
	if coll.Title != "" {
		return coll.Title
	}
	return fmt.Sprintf("#%d", index+1)
}

// isLocalSpecPath checks whether the location points to a file inside the
// workspace specs/ directory. If so, the collection is already imported
// and should be skipped during bulk import.
func isLocalSpecPath(location, specsDir, workspaceRoot string) bool {
	if strings.HasPrefix(location, "http://") || strings.HasPrefix(location, "https://") {
		return false
	}
	absLoc := location
	if !filepath.IsAbs(location) {
		absLoc = filepath.Join(workspaceRoot, location)
	}
	absLoc, err := filepath.Abs(absLoc)
	if err != nil {
		return false
	}
	absSpecs, err := filepath.Abs(specsDir)
	if err != nil {
		return false
	}
	rel, err := filepath.Rel(absSpecs, absLoc)
	if err != nil {
		return false
	}
	return !strings.HasPrefix(rel, "..")
}
