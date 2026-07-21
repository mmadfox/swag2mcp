package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

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

	if req.Source == "" && len(req.SpecFilter) == 0 {
		return ImportResponse{}, NewValidationError(
			"Import requires a source URL with filename, a spec filter, or a zip backup path.",
			errors.New("no import source specified"),
		)
	}

	if req.Source != "" && req.Name == "" {
		return ImportResponse{}, NewValidationError(
			"Single import requires both a source URL and a filename.",
			errors.New("name is required when source is provided"),
		)
	}

	if len(req.SpecFilter) > 0 {
		return ims.importSpecs(ctx, req)
	}
	return ims.importSingle(ctx, req)
}

func (ims *importService) importFromZip(_ context.Context, req ImportRequest) (ImportResponse, error) {
	if !workspace.IsSwag2mcpZip(req.ZipSource) {
		return ImportResponse{}, NewValidationError(
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
		return ImportResponse{}, NewInvokeError(
			fmt.Sprintf("Failed to download spec from %q.", req.Source),
			err,
		)
	}

	path, err := ims.ws.SaveSpec(req.Name, data)
	if err != nil {
		return ImportResponse{}, NewValidationError(
			fmt.Sprintf("Failed to save spec as %q. The filename may already exist.", req.Name),
			err,
		)
	}

	return ImportResponse{
		Files: []ImportedFile{
			{
				Source:    req.Source,
				Name:      req.Name,
				SavedPath: path,
			},
		},
	}, nil
}

func (ims *importService) importSpecs(ctx context.Context, req ImportRequest) (ImportResponse, error) {
	if req.ConfFilePath == "" {
		return ImportResponse{}, NewValidationError(
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

			data, err := ims.ws.DownloadSpec(ctx, coll.Location)
			if err != nil {
				return ImportResponse{}, NewInvokeError(
					fmt.Sprintf("Failed to download spec for collection %q.", coll.Title),
					err,
				)
			}

			name := specFileName(spec.Domain, coll.Title, coll.Location)
			sp, err := ims.ws.SaveSpec(name, data)
			if err != nil {
				return ImportResponse{}, NewValidationError(
					fmt.Sprintf("Failed to save spec as %q. The filename may already exist.", name),
					err,
				)
			}

			coll.Location = filepath.Join("specs", name)
			updated = true

			imported = append(imported, ImportedFile{
				Source:    coll.Location,
				Name:      name,
				SavedPath: sp,
			})
		}
	}

	if !updated {
		return ImportResponse{}, NewNotFoundError(
			fmt.Sprintf("No matching specs found for filter %v.", req.SpecFilter),
			fmt.Errorf("no specs matched filter %v", req.SpecFilter),
		)
	}

	if err := config.Save(cfg, req.ConfFilePath); err != nil {
		return ImportResponse{}, NewConfigError(
			fmt.Sprintf("Failed to save updated configuration to %q.", req.ConfFilePath),
			err,
		)
	}

	return ImportResponse{Files: imported}, nil
}
