package service

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/workspace"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// ImportRequest is the request for the Import method.
type ImportRequest struct {
	Source       string
	Name         string
	SpecFilter   []string
	ConfFilePath string
	ZipSource    string // path to a swag2mcp backup ZIP to restore
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

// Import imports spec files into the workspace specs/ directory.
//
// When SpecFilter is empty, it imports a single spec from Source and saves it as Name.
// When SpecFilter is set, it reads the config, downloads all collections from matching specs,
// saves them to specs/, and updates the config with the new locations.
// When ZipSource is set, it restores a full workspace from a swag2mcp backup ZIP.
func (s *Service) Import(ctx context.Context, req ImportRequest) (ImportResponse, error) {
	if req.ZipSource != "" {
		return s.importFromZip(ctx, req)
	}

	if req.Source == "" && len(req.SpecFilter) == 0 {
		return ImportResponse{}, NewValidationError(
			"Import requires either a source and name (single import), a spec filter (bulk import), "+
				"or a zip source (restore from backup). "+
				"For single import: provide the spec URL/path and a unique filename. "+
				"For bulk import: use --spec flag with one or more spec domain names. "+
				"For restore: use --from-zip flag with a swag2mcp backup archive.",
			errors.New("no import source specified"),
		)
	}

	if req.Source != "" && req.Name == "" {
		return ImportResponse{}, NewValidationError(
			"Single import requires both a source and name. "+
				"Example: swag2mcp import https://example.com/spec.yaml myspec",
			errors.New("name is required when source is provided"),
		)
	}

	if len(req.SpecFilter) > 0 {
		return s.importSpecs(ctx, req)
	}
	return s.importSingle(ctx, req)
}

func (s *Service) importFromZip(_ context.Context, req ImportRequest) (ImportResponse, error) {
	if !workspace.IsSwag2mcpZip(req.ZipSource) {
		return ImportResponse{}, NewValidationError(
			fmt.Sprintf("The file %q is not a valid swag2mcp backup archive. "+
				"Use 'swag2mcp export' to create a backup, then import it here.", req.ZipSource),
			fmt.Errorf("invalid swag2mcp zip: %s", req.ZipSource),
		)
	}

	extractDir, dirErr := os.MkdirTemp("", "swag2mcp-restore-*")
	if dirErr != nil {
		return ImportResponse{}, NewInvokeError(
			"Failed to create temporary directory for extraction.",
			dirErr,
		)
	}
	defer os.RemoveAll(extractDir)

	if extractErr := workspace.ExtractZip(req.ZipSource, extractDir); extractErr != nil {
		return ImportResponse{}, NewInvokeError(
			fmt.Sprintf("Failed to extract archive %q.", req.ZipSource),
			extractErr,
		)
	}

	if initErr := s.ws.Init(); initErr != nil {
		return ImportResponse{}, NewInvokeError(
			"Failed to initialize workspace directories.",
			initErr,
		)
	}

	if copyErr := s.ws.CopySpecsToWorkspace(extractDir); copyErr != nil {
		return ImportResponse{}, NewInvokeError(
			"Failed to copy spec files from backup to workspace.",
			copyErr,
		)
	}

	if copyErr := s.ws.CopyAuthScriptsToWorkspace(extractDir); copyErr != nil {
		return ImportResponse{}, NewInvokeError(
			"Failed to copy auth scripts from backup to workspace.",
			copyErr,
		)
	}

	cfgData, cfgReadErr := workspace.ReadConfigFromExport(extractDir)
	if cfgReadErr != nil {
		return ImportResponse{}, NewInvokeError(
			"Failed to read configuration from backup.",
			cfgReadErr,
		)
	}

	cfgPath := s.ws.ConfigPath()
	if writeErr := os.WriteFile(cfgPath, cfgData, 0600); writeErr != nil {
		return ImportResponse{}, NewInvokeError(
			fmt.Sprintf("Failed to write configuration to %q.", cfgPath),
			writeErr,
		)
	}

	specs, listErr := s.ws.ListSpecs()
	if listErr != nil {
		return ImportResponse{}, NewInvokeError(
			"Failed to list imported spec files.",
			listErr,
		)
	}

	files := make([]ImportedFile, 0, len(specs))
	for _, name := range specs {
		files = append(files, ImportedFile{
			Source:    req.ZipSource,
			Name:      name,
			SavedPath: s.ws.SpecPath(name),
		})
	}

	return ImportResponse{
		Files: files,
	}, nil
}

func (s *Service) importSingle(ctx context.Context, req ImportRequest) (ImportResponse, error) {
	data, err := s.ws.DownloadSpec(ctx, req.Source)
	if err != nil {
		return ImportResponse{}, NewInvokeError(
			fmt.Sprintf("Failed to download spec from %q. Check that the URL or file path is correct and accessible.", req.Source),
			err,
		)
	}

	path, err := s.ws.SaveSpec(req.Name, data)
	if err != nil {
		return ImportResponse{}, NewValidationError(
			fmt.Sprintf("Failed to save spec as %q. The filename may already exist in the specs/ directory. "+
				"Use a different name or remove the existing file first.", req.Name),
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

func (s *Service) importSpecs(ctx context.Context, req ImportRequest) (ImportResponse, error) {
	if req.ConfFilePath == "" {
		return ImportResponse{}, NewValidationError(
			"Configuration file path is required when using --spec filter. "+
				"Run the command from within a workspace directory or provide the workspace path.",
			errors.New("config file path is empty"),
		)
	}

	cfg, err := config.Load(req.ConfFilePath)
	if err != nil {
		return ImportResponse{}, NewInvokeError(
			fmt.Sprintf("Failed to load configuration from %q. Ensure the config file exists and is valid YAML.", req.ConfFilePath),
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

			data, err := s.ws.DownloadSpec(ctx, coll.Location)
			if err != nil {
				return ImportResponse{}, NewInvokeError(
					fmt.Sprintf("Failed to download spec from %q for collection %q in spec %q. "+
						"Check that the location is correct and accessible.",
						coll.Location, coll.Title, spec.Domain),
					err,
				)
			}

			name := specFileName(spec.Domain, coll.Title, coll.Location)
			sp, err := s.ws.SaveSpec(name, data)
			if err != nil {
				return ImportResponse{}, NewValidationError(
					fmt.Sprintf("Failed to save spec as %q for collection %q in spec %q. "+
						"The filename may already exist in the specs/ directory.",
						name, coll.Title, spec.Domain),
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
			fmt.Sprintf("No matching specs found for filter %v. Use spec_list to see available specs.", req.SpecFilter),
			fmt.Errorf("no specs matched filter %v", req.SpecFilter),
		)
	}

	if err := config.Save(cfg, req.ConfFilePath); err != nil {
		return ImportResponse{}, NewInvokeError(
			fmt.Sprintf("Failed to save updated configuration to %q.", req.ConfFilePath),
			err,
		)
	}

	return ImportResponse{Files: imported}, nil
}

type specFilter struct {
	domains map[string]struct{}
}

func makeFilter(domains []string) *specFilter {
	f := &specFilter{domains: make(map[string]struct{}, len(domains))}
	for _, d := range domains {
		f.domains[strings.TrimSpace(d)] = struct{}{}
	}
	return f
}

func (f *specFilter) match(domain string) bool {
	if len(f.domains) == 0 {
		return true
	}
	_, ok := f.domains[domain]
	return ok
}

func specFileName(domain, title, location string) string {
	base := title
	if base == "" {
		base = specFileNameBase(location)
	}

	ext := filepath.Ext(base)
	if ext == "" {
		ext = ".yaml"
	}
	base = strings.TrimSuffix(base, ext)
	base = strings.TrimSuffix(base, ".yml")

	sanitized := strings.ToLower(base)
	sanitized = strings.NewReplacer(
		" ", "-",
		"_", "-",
		".", "-",
	).Replace(sanitized)
	sanitized = removeDiacritics(sanitized)

	if sanitized == domain {
		return fmt.Sprintf("%s%s", domain, ext)
	}

	return fmt.Sprintf("%s-%s%s", domain, sanitized, ext)
}

func removeDiacritics(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.Predicate(unicode.IsMark)), norm.NFC)
	result, _, _ := transform.String(t, s)
	return result
}

const defaultSpecName = "spec"

func specFileNameBase(location string) string {
	if strings.HasPrefix(location, "http://") || strings.HasPrefix(location, "https://") {
		u, err := url.Parse(location)
		if err == nil && u.Path != "" && u.Path != "/" {
			return filepath.Base(u.Path)
		}
		return defaultSpecName
	}
	return filepath.Base(location)
}
