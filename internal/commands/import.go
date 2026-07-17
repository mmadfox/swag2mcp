package commands

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/mmadfox/swag2mcp/internal/service"
	"github.com/mmadfox/swag2mcp/internal/workspace"
)

const (
	importMaxArgs = 3
	importTwoArgs = 2
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
  swag2mcp import --spec petstore
  swag2mcp import /path/to/workspace --spec petstore,store

Restore from backup (--from-zip flag or .zip file as source):
  swag2mcp import --from-zip /path/to/backup.zip
  swag2mcp import /path/to/workspace /path/to/backup.zip

The --spec flag imports all collections from the matching specs in the config,
saves them to specs/, and updates the config with the new locations.
The --from-zip flag or a .zip source restores a full workspace from a swag2mcp backup archive.`,
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			basePath, source, name, zipSource := parseImportArgs(args, opts.Specs, opts.FromZip)
			return runImport(basePath, source, name, zipSource, opts.Specs, cmd)
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

func parseImportArgs(args []string, specs []string, fromZip string) (string, string, string, string) {
	if fromZip != "" {
		if len(args) > 0 {
			return args[0], "", "", fromZip
		}
		return "", "", "", fromZip
	}

	if len(specs) > 0 {
		if len(args) > 0 {
			return args[0], "", "", ""
		}
		return "", "", "", ""
	}

	l := len(args)
	if l >= importMaxArgs {
		if isZipFile(args[l-1]) {
			return args[0], "", "", args[l-1]
		}
		return args[0], args[1], args[2], ""
	}
	if l == importTwoArgs {
		if isZipFile(args[1]) {
			return args[0], "", "", args[1]
		}
		return "", args[0], args[1], ""
	}
	if l == 1 {
		if isZipFile(args[0]) {
			return "", "", "", args[0]
		}
		return "", args[0], "", ""
	}
	return "", "", "", ""
}

func isZipFile(path string) bool {
	return filepath.Ext(path) == ".zip"
}

func runImport(basePath, source, name, zipSource string, specs []string, cmd *cobra.Command) error {
	ws, wsErr := workspace.NewFromBase(basePath)
	if wsErr != nil {
		return fmt.Errorf("workspace: %w", wsErr)
	}

	svc, svcErr := service.New(service.WithWorkspace(ws))
	if svcErr != nil {
		return fmt.Errorf("service: %w", svcErr)
	}

	if zipSource != "" {
		if initErr := ws.Init(); initErr != nil {
			return fmt.Errorf("workspace init: %w", initErr)
		}

		_, importErr := svc.Import(cmd.Context(), service.ImportRequest{
			ZipSource: zipSource,
		})
		if importErr != nil {
			return importErr
		}

		cmd.Println("✅ Restored successfully!")
		return nil
	}

	if len(specs) > 0 {
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
	}

	if source == "" || name == "" {
		return errors.New("import requires a source and name (single import), --spec flag (bulk import), or --from-zip (restore from backup)\n\n" +
			"Single import:\n" +
			"  swag2mcp import <source> <name>\n" +
			"  swag2mcp import /path/to/workspace <source> <name>\n\n" +
			"Bulk import:\n" +
			"  swag2mcp import --spec petstore\n" +
			"  swag2mcp import /path/to/workspace --spec petstore,store\n\n" +
			"Restore from backup:\n" +
			"  swag2mcp import --from-zip /path/to/backup.zip\n" +
			"  swag2mcp import /path/to/workspace /path/to/backup.zip")
	}

	if initErr := ws.Init(); initErr != nil {
		return fmt.Errorf("workspace init: %w", initErr)
	}

	resp, importErr := svc.Import(cmd.Context(), service.ImportRequest{
		Source: source,
		Name:   name,
	})
	if importErr != nil {
		return importErr
	}

	cmd.Printf("✅ Imported %s → %s\n", resp.Files[0].Source, resp.Files[0].SavedPath)
	return nil
}
