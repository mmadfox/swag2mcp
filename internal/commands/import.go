package commands

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/httpclient"
	"github.com/mmadfox/swag2mcp/internal/service"
	"github.com/mmadfox/swag2mcp/internal/workspace"
)

func newImportCmd() *cobra.Command {
	opts := struct {
		Specs   []string
		FromZip string
	}{}

	cmd := &cobra.Command{
		Use:   "import [path] [source] [name]",
		Short: "Import a spec file into the workspace",
		Long: `Import a spec file into the workspace specs/ directory.

Single import (requires source and name):
  swag2mcp import https://example.com/spec.yaml myspec
  swag2mcp import /path/to/workspace https://example.com/spec.yaml myspec
  swag2mcp import ./local-spec.yaml myspec

Bulk import (requires --spec flag):
  swag2mcp import --spec meteo
  swag2mcp import /path/to/workspace --spec meteo,store

Restore from backup (--from-zip flag or .zip file as source):
  swag2mcp import --from-zip /path/to/backup.zip
  swag2mcp import /path/to/workspace /path/to/backup.zip

The --spec flag imports all collections from the matching specs in the config,
saves them to specs/, and updates the config with the new locations.
The --from-zip flag or a .zip source restores a full workspace from a swag2mcp backup archive.`,
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			parsed := parseImportArgs(args, opts.Specs, opts.FromZip)
			return runImport(parsed, opts.Specs, cmd)
		},
	}

	cmd.Flags().StringSliceVarP(&opts.Specs, "spec", "s", nil,
		"Import all collections from specified specs (comma-separated)")
	cmd.Flags().StringVar(&opts.FromZip, "from-zip", "",
		"Restore workspace from a swag2mcp backup ZIP archive")
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	return cmd
}

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

func parseImportArgs(args []string, specs []string, fromZip string) importArgs {
	if fromZip != "" {
		basePath := ""
		if len(args) > 0 {
			basePath = args[0]
		}
		return importArgs{mode: importModeZip, basePath: basePath, zipFile: fromZip}
	}

	if len(specs) > 0 {
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
		return importArgs{mode: importModeSingle, source: args[0], name: args[1]}
	}
	return importArgs{mode: importModeSingle, source: args[0]}
}

func isZipFile(path string) bool {
	return filepath.Ext(path) == ".zip"
}

func runImport(parsed importArgs, specs []string, cmd *cobra.Command) error {
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
			return importErr
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
			return importErr
		}

		cmd.Printf("✅ Imported %d spec files:\n", len(resp.Files))
		for _, f := range resp.Files {
			cmd.Printf("   • %s → %s\n", f.Source, f.SavedPath)
		}
		return nil

	case importModeSingle:
		if parsed.source == "" || parsed.name == "" {
			return errors.New("import requires a source and name (single import), --spec flag (bulk import), or --from-zip (restore from backup)\n\n" +
				"Single import:\n" +
				"  swag2mcp import <source> <name>\n" +
				"  swag2mcp import /path/to/workspace <source> <name>\n\n" +
				"Bulk import:\n" +
				"  swag2mcp import --spec meteo\n" +
				"  swag2mcp import /path/to/workspace --spec meteo,store\n\n" +
				"Restore from backup:\n" +
				"  swag2mcp import --from-zip /path/to/backup.zip\n" +
				"  swag2mcp import /path/to/workspace /path/to/backup.zip")
		}

		if initErr := ws.Init(); initErr != nil {
			return fmt.Errorf("workspace init: %w", initErr)
		}

		resp, importErr := svc.Import(cmd.Context(), service.ImportRequest{
			Source: parsed.source,
			Name:   parsed.name,
		})
		if importErr != nil {
			return importErr
		}

		cmd.Printf("✅ Imported %s → %s\n", resp.Files[0].Source, resp.Files[0].SavedPath)
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
