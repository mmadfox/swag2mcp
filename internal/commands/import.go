/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package commands

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/httpclient"
	"github.com/mmadfox/swag2mcp/internal/service"
	"github.com/mmadfox/swag2mcp/internal/workspace"
)

func formatImportErr(err error) string {
	llmErr, ok := errors.AsType[*service.LLMError](err)
	if ok {
		return fmt.Sprintf("Error: %s\n  %s", llmErr.Code, llmErr.Message)
	}
	return fmt.Sprintf("Error: %s", err)
}

func newImportCmd() *cobra.Command {
	opts := struct {
		Specs   []string
		FromZip string
		Force   bool
	}{}

	var specFlag string

	cmd := &cobra.Command{
		Use:   "import [path] [source] [name]",
		Short: "Import spec files or restore workspace from backup",
		Long: `Import spec files into the workspace specs/ directory for local use.

After downloading, manually add the spec to swag2mcp.yaml and set location to specs/<filename>.

Single import — download a spec file and save it to specs/:
  swag2mcp import https://example.com/spec.yaml example-api.yaml
  swag2mcp import /path/to/workspace https://example.com/spec.yaml example-api.yaml
  swag2mcp import ./local-spec.yaml example-api.yaml

Single import (name derived from URL):
  swag2mcp import https://example.com/specs/petstore.yaml

Single import with overwrite:
  swag2mcp import --force https://example.com/spec.yaml example-api.yaml

Bulk import (--spec) — download all collection spec files from the config
and update their locations to local paths in specs/:
  swag2mcp import --spec                (all specs)
  swag2mcp import --spec meteo           (specific spec)
  swag2mcp import --spec meteo,github     (multiple specs)
  swag2mcp import /path/to/workspace --spec meteo

Restore from backup (--from-zip or .zip file):
  swag2mcp import --from-zip /path/to/backup.zip
  swag2mcp import /path/to/workspace /path/to/backup.zip`,
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			specFlagSet := cmd.Flags().Changed("spec")
			if specFlagSet {
				if specFlag == noOptDefValMarker {
					if len(args) > 0 {
						last := args[len(args)-1]
						if !strings.Contains(last, "/") && !strings.HasPrefix(last, ".") {
							opts.Specs = strings.Split(last, ",")
							args = args[:len(args)-1]
						} else {
							opts.Specs = nil
						}
					} else {
						opts.Specs = nil
					}
				} else {
					opts.Specs = strings.Split(specFlag, ",")
				}
			}
			parsed := parseImportArgs(args, opts.FromZip, specFlagSet)
			return runImport(parsed, opts.Specs, opts.Force, cmd)
		},
	}

	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false,
		"Overwrite existing spec files without error")
	cmd.Flags().StringVarP(&specFlag, "spec", "s", "",
		"Download collection spec files from the config (use without value for all specs, or specify domains like --spec meteo,github)")
	cmd.Flags().Lookup("spec").NoOptDefVal = noOptDefValMarker
	cmd.Flags().StringVar(&opts.FromZip, "from-zip", "",
		"Restore workspace from a swag2mcp backup ZIP archive")
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	return cmd
}

const noOptDefValMarker = "*"

type importMode int

const (
	importModeSingle importMode = iota
	importModeBulk
	importModeZip
)

const (
	importArgsMinForFull = 3
	importArgsSourceName = 2
)

type importArgs struct {
	mode     importMode
	basePath string
	source   string
	name     string
	zipFile  string
}

func parseImportArgs(args []string, fromZip string, specFlagSet bool) importArgs {
	if fromZip != "" {
		basePath := ""
		if len(args) > 0 {
			basePath = args[0]
		}
		return importArgs{mode: importModeZip, basePath: basePath, zipFile: fromZip}
	}

	if specFlagSet {
		basePath := ""
		if len(args) > 0 {
			basePath = args[0]
		}
		return importArgs{mode: importModeBulk, basePath: basePath}
	}

	return parseSingleOrZipArgs(args)
}

func parseSingleOrZipArgs(args []string) importArgs {
	l := len(args)
	if l == 0 {
		return importArgs{mode: importModeSingle}
	}

	last := args[l-1]
	if isZipFile(last) {
		basePath := ""
		if l > 1 {
			basePath = args[0]
		}
		return importArgs{mode: importModeZip, basePath: basePath, zipFile: last}
	}

	if l >= importArgsMinForFull {
		return importArgs{mode: importModeSingle, basePath: args[0], source: args[1], name: args[2]}
	}
	if l == importArgsSourceName {
		src, name := args[0], args[1]
		if isURL(name) {
			return importArgs{mode: importModeSingle, basePath: src, source: name}
		}
		return importArgs{mode: importModeSingle, source: src, name: name}
	}
	return importArgs{mode: importModeSingle, source: args[0]}
}

func isURL(s string) bool {
	return strings.HasPrefix(s, "https://") || strings.HasPrefix(s, "http://")
}

func isZipFile(path string) bool {
	return filepath.Ext(path) == ".zip"
}

var specExts = []string{".yaml", ".yml", ".json", ".swagger", ".postman"} //nolint:gochecknoglobals // Allowed spec file extensions.

func isValidSpecLocation(location string) bool {
	lower := strings.ToLower(location)
	for _, ext := range specExts {
		if strings.HasSuffix(lower, ext) {
			return true
		}
	}
	return false
}

func runImport(parsed importArgs, specs []string, force bool, cmd *cobra.Command) error {
	ws, wsErr := workspace.NewFromBase(parsed.basePath)
	if wsErr != nil {
		return fmt.Errorf("workspace: %w", wsErr)
	}

	svc, svcErr := service.New(service.WithWorkspace(ws))
	if svcErr != nil {
		return fmt.Errorf("service: %w", svcErr)
	}

	setupGlobalHTTPClient(ws.ConfigPath())

	switch parsed.mode {
	case importModeZip:
		if initErr := ws.Init(); initErr != nil {
			return fmt.Errorf("workspace init: %w", initErr)
		}

		_, importErr := svc.Import(cmd.Context(), service.ImportRequest{
			ZipSource: parsed.zipFile,
		})
		if importErr != nil {
			return errors.New(formatImportErr(importErr))
		}

		cmd.Println("✅ Restored successfully!")
		return nil

	case importModeBulk:
		cfgPath := ws.ConfigPath()
		if ws.ConfigNotExists() {
			return fmt.Errorf("configuration not found at %s\n  Run 'swag2mcp init' first or provide a workspace path with a valid config", cfgPath)
		}

		resp, importErr := svc.Import(cmd.Context(), service.ImportRequest{
			SpecFilter:   specs,
			ConfFilePath: cfgPath,
		})
		if importErr != nil {
			return errors.New(formatImportErr(importErr))
		}

		cmd.Printf("✅ Imported %d spec files:\n", len(resp.Files))
		for _, f := range resp.Files {
			cmd.Printf("   • %s → %s\n", f.Source, f.SavedPath)
		}
		return nil

	case importModeSingle:
		if parsed.source == "" {
			return errors.New("import requires a source, --spec flag (bulk import), or --from-zip (restore from backup)\n\n" +
				"Single import — download a spec file and save it to specs/:\n" +
				"  swag2mcp import <source> <name>\n" +
				"  swag2mcp import <source>  (name derived from URL)\n" +
				"  swag2mcp import /path/to/workspace <source> <name>\n\n" +
				"Bulk import (--spec) — download all collection spec files from the config\n" +
				"and update their locations to local paths in specs/:\n" +
				"  swag2mcp import --spec                (all specs)\n" +
				"  swag2mcp import --spec meteo           (specific spec)\n" +
				"  swag2mcp import /path/to/workspace --spec meteo,github\n\n" +
				"Restore from backup:\n" +
				"  swag2mcp import --from-zip /path/to/backup.zip\n" +
				"  swag2mcp import /path/to/workspace /path/to/backup.zip")
		}

		if initErr := ws.Init(); initErr != nil {
			return fmt.Errorf("workspace init: %w", initErr)
		}

		if !isValidSpecLocation(parsed.source) {
			return errors.New("source must be a spec file with one of these extensions: " +
				".yaml, .yml, .json, .swagger, .postman")
		}

		resp, importErr := svc.Import(cmd.Context(), service.ImportRequest{
			Source: parsed.source,
			Name:   parsed.name,
			Force:  force,
		})
		if importErr != nil {
			return errors.New(formatImportErr(importErr))
		}

		relPath := resp.Files[0].SavedPath
		if r, err := filepath.Rel(ws.Root(), resp.Files[0].SavedPath); err == nil {
			relPath = r
		}
		cmd.Printf("✅ Imported to %s\n", ws.Root())
		cmd.Printf("   %s\n\n", relPath)
		cmd.Printf("   Add to swag2mcp.yaml:\n")
		cmd.Printf("     specs:\n")
		cmd.Printf("       - domain: <your-domain>\n")
		cmd.Printf("         collections:\n")
		cmd.Printf("           - location: %s\n", relPath)
		return nil
	}

	return nil
}

// setupGlobalHTTPClient loads the config and sets the global HTTP client config
// so that httpclient.NewDefault() returns a properly configured client.
// If the config file does not exist, the default client with a 30s timeout is used.
func setupGlobalHTTPClient(configPath string) {
	cfg, err := config.Load(configPath)
	if err != nil {
		return
	}

	httpCfg := service.BuildGlobalHTTPConfig(cfg.HTTPClient)
	httpclient.SetGlobalConfig(httpCfg)
}
