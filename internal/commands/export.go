package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/mmadfox/swag2mcp/internal/service"
	"github.com/mmadfox/swag2mcp/internal/workspace"
)

const exportMaxArgs = 2

func newExportCmd() *cobra.Command {
	opts := struct {
		Specs []string
	}{}

	cmd := &cobra.Command{
		Use:   "export [path] [output]",
		Short: "Export workspace as a portable ZIP backup",
		Long: `Export the workspace as a portable ZIP backup archive.

The archive contains all spec files, configuration, and auth scripts,
ready to be imported on another machine or restored later.

Examples:
  swag2mcp export
  swag2mcp export /path/to/workspace
  swag2mcp export /path/to/workspace /path/to/backup.zip
  swag2mcp export --spec petstore
  swag2mcp export --spec petstore,store`,
		Args: cobra.MaximumNArgs(exportMaxArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			basePath, outputPath := parseExportArgs(args)
			return runExport(basePath, outputPath, opts.Specs, cmd)
		},
	}

	cmd.Flags().StringSliceVarP(&opts.Specs, "spec", "s", nil,
		"Export only specified specs (comma-separated)")
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	return cmd
}

func parseExportArgs(args []string) (string, string) {
	switch len(args) {
	case exportMaxArgs:
		return args[0], args[1]
	case exportMaxArgs - 1:
		return args[0], ""
	default:
		return "", ""
	}
}

func runExport(basePath, outputPath string, specs []string, cmd *cobra.Command) error {
	ws, wsErr := workspace.NewFromBase(basePath)
	if wsErr != nil {
		return fmt.Errorf("workspace: %w", wsErr)
	}

	svc, svcErr := service.New(service.WithWorkspace(ws))
	if svcErr != nil {
		return fmt.Errorf("service: %w", svcErr)
	}

	if outputPath == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("get current directory: %w", err)
		}
		outputPath = filepath.Join(cwd, workspace.DefaultExportName())
	}

	absOutput, absErr := filepath.Abs(outputPath)
	if absErr != nil {
		return fmt.Errorf("resolve output path: %w", absErr)
	}

	resp, exportErr := svc.Export(cmd.Context(), service.ExportRequest{
		OutputPath: absOutput,
		SpecFilter: specs,
	})
	if exportErr != nil {
		return exportErr
	}

	cmd.Printf("✅ Exported %d spec files to %s\n", resp.FileCount, resp.OutputPath)
	return nil
}
