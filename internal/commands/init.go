package commands

import (
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/mmadfox/swag2mcp/internal/auth"
	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/tui"
	"github.com/mmadfox/swag2mcp/internal/workspace"
)

func newInitCmd() *cobra.Command {
	opts := struct {
		Interactive bool
		Force       bool
	}{}

	cmd := &cobra.Command{
		Use:   "init [path]",
		Short: "Initialize workspace and configuration",
		Long: `Initialize workspace and configuration.

  swag2mcp init              — create ~/.swag2mcp/swag2mcp.yaml
  swag2mcp init ./           — create ./swag2mcp.yaml
  swag2mcp init path/to      — create path/to/swag2mcp.yaml
  swag2mcp init -i           — interactive wizard
  swag2mcp init -f           — force overwrite existing configuration`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			basePath := ""
			if len(args) > 0 {
				basePath = args[0]
			}
			return runInit(basePath, opts.Interactive, opts.Force, cmd)
		},
	}

	cmd.Flags().BoolVarP(&opts.Interactive, "interactive", "i", false, "Run interactive wizard")
	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "Overwrite existing configuration")
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	return cmd
}

func runInit(basePath string, interactive, force bool, cmd *cobra.Command) error {
	if interactive {
		cfgPath, wsDir, _, err := tui.Run()
		if err != nil {
			return fmt.Errorf("init wizard: %w", err)
		}

		ws, wsErr := workspace.New(wsDir)
		if wsErr == nil {
			if cfg, loadErr := config.Load(cfgPath); loadErr == nil {
				ensureAuthScripts(cfg, ws)
			}
		}

		return nil
	}

	var (
		workspaceDir string
		configPath   string
	)

	if basePath == "" {
		ws, wsErr := workspace.New("")
		if wsErr != nil {
			return fmt.Errorf("workspace: %w", wsErr)
		}
		workspaceDir = ws.Root()
		configPath = ws.ConfigPath()
	} else {
		absBase, err := filepath.Abs(basePath)
		if err != nil {
			return fmt.Errorf("resolve path: %w", err)
		}
		workspaceDir = absBase
		configPath = filepath.Join(absBase, "swag2mcp.yaml")
	}

	if !force {
		ws, wsErr := workspace.New(workspaceDir)
		if wsErr != nil {
			return fmt.Errorf("workspace: %w", wsErr)
		}
		empty, emptyErr := ws.IsEmpty()
		if emptyErr != nil {
			return fmt.Errorf("check directory: %w", emptyErr)
		}
		if !empty {
			return fmt.Errorf("directory %q is not empty\n  Use --force to initialize in a non-empty directory", workspaceDir)
		}
	}

	if err := tui.Setup(configPath, workspaceDir); err != nil {
		return fmt.Errorf("init: %w", err)
	}

	ws, wsErr := workspace.New(workspaceDir)
	if wsErr == nil {
		if cfg, loadErr := config.Load(configPath); loadErr == nil {
			ensureAuthScripts(cfg, ws)
		}
	}

	cmd.Printf("✅ Configuration written to %s\n", configPath)
	cmd.Printf("✅ Workspace initialized at %s\n", workspaceDir)
	cmd.Printf("\nNext step: edit %s or run 'swag2mcp ls' to list configured specs\n", filepath.Base(configPath))
	return nil
}

func ensureAuthScripts(cfg *config.Config, ws *workspace.Workspace) {
	for spec := range cfg.Iterate(nil) {
		if spec.Auth.Client != nil && spec.Auth.Client.Type() == auth.ScriptAuth {
			if err := ws.EnsureAuthScript(spec.Domain); err != nil {
				slog.Default().Warn("failed to ensure auth script", "domain", spec.Domain, "error", err)
			}
		}
	}
}
