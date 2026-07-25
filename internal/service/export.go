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

// ExportRequest is the request for the Export method.
type ExportRequest struct {
	OutputPath string
	SpecFilter []string
}

// ExportResponse holds the result of an export operation.
type ExportResponse struct {
	OutputPath string `json:"outputPath"`
	FileCount  int    `json:"fileCount"`
}

// Export creates a portable ZIP backup of the workspace.
func (s *Service) Export(ctx context.Context, req ExportRequest) (ExportResponse, error) {
	cfg, loadErr := s.loadExportConfig()
	if loadErr != nil {
		return ExportResponse{}, loadErr
	}

	exportDir, dirErr := s.ws.CreateExportDir()
	if dirErr != nil {
		return ExportResponse{}, NewInvokeError(
			"Failed to create temporary export directory.",
			dirErr,
		)
	}
	defer os.RemoveAll(exportDir)

	if err := workspace.CreateEmptyDirsInExport(exportDir); err != nil {
		return ExportResponse{}, NewInvokeError(
			"Failed to create export directory structure.",
			err,
		)
	}

	fileCount, exportErr := s.exportCollections(ctx, cfg, req.SpecFilter, exportDir)
	if exportErr != nil {
		return ExportResponse{}, exportErr
	}

	if fileCount == 0 {
		return ExportResponse{}, NewNotFoundError(
			"No collections found to export. Ensure the workspace has specs with valid collections.",
			errors.New("no collections to export"),
		)
	}

	if err := s.finalizeExport(cfg, exportDir); err != nil {
		return ExportResponse{}, err
	}

	outputPath := req.OutputPath
	if outputPath == "" {
		outputPath = workspace.DefaultExportName()
	}

	if err := workspace.CreateZip(exportDir, outputPath); err != nil {
		return ExportResponse{}, NewInvokeError(
			fmt.Sprintf("Failed to create ZIP archive at %q.", outputPath),
			err,
		)
	}

	return ExportResponse{
		OutputPath: outputPath,
		FileCount:  fileCount,
	}, nil
}

func (s *Service) loadExportConfig() (*config.Config, error) {
	cfgPath := s.ws.ConfigPath()
	if s.ws.ConfigNotExists() {
		return nil, NewNotFoundError(
			"No configuration found in the workspace. Run 'swag2mcp init' first.",
			fmt.Errorf("config not found at %s", cfgPath),
		)
	}
	cfg, loadErr := config.Load(cfgPath)
	if loadErr != nil {
		return nil, NewInvokeError(
			fmt.Sprintf("Failed to load configuration from %q.", cfgPath),
			loadErr,
		)
	}
	return cfg, nil
}

func (s *Service) exportCollections(ctx context.Context, cfg *config.Config, specFilter []string, exportDir string) (int, error) {
	filter := makeFilter(specFilter)
	locationMap := make(map[string]string)
	fileCount := 0

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
			if _, already := locationMap[coll.Location]; already {
				continue
			}

			data, dlErr := s.ws.DownloadSpec(ctx, coll.Location)
			if dlErr != nil {
				return 0, NewInvokeError(
					fmt.Sprintf("Failed to download spec from %q for collection %q in spec %q. "+
						"Check that the location is correct and accessible.",
						coll.Location, coll.Title, spec.Domain),
					dlErr,
				)
			}

			name := specFileName(spec.Domain, coll.Title, coll.Location)
			if writeErr := workspace.WriteSpecToExport(exportDir, name, data); writeErr != nil {
				return 0, NewInvokeError(
					fmt.Sprintf("Failed to write spec %q to export.", name),
					writeErr,
				)
			}

			locationMap[coll.Location] = filepath.Join("specs", name)
			fileCount++
		}
	}

	s.updateConfigLocations(cfg, specFilter, locationMap)

	return fileCount, nil
}

func (s *Service) updateConfigLocations(cfg *config.Config, specFilter []string, locationMap map[string]string) {
	filter := makeFilter(specFilter)
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
			if newLoc, ok := locationMap[coll.Location]; ok {
				coll.Location = newLoc
			}
		}
	}
}

func (s *Service) finalizeExport(cfg *config.Config, exportDir string) error {
	if err := config.Save(cfg, filepath.Join(exportDir, "swag2mcp.yaml")); err != nil {
		return NewInvokeError(
			"Failed to save updated configuration to export.",
			err,
		)
	}
	if err := s.ws.CopyAuthScriptsToExport(exportDir); err != nil {
		return NewInvokeError(
			"Failed to copy auth scripts to export.",
			err,
		)
	}
	if err := workspace.CreateMetaFile(exportDir, s.version); err != nil {
		return NewInvokeError(
			"Failed to create backup metadata.",
			err,
		)
	}
	return nil
}
