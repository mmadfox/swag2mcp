package commands

import (
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "swag2mcp",
		Short:        "swag2mcp - MCP server for OpenAPI/Swagger specifications",
		Long:         `swag2mcp provides LLM agents with tools to work with Swagger/OpenAPI 3+ specifications.`,
		SilenceUsage: true,
	}

	cmd.AddCommand(newInitCmd())
	cmd.AddCommand(newAddCmd())
	cmd.AddCommand(newDeleteCmd())
	cmd.AddCommand(newLsCmd())
	cmd.AddCommand(newRunCmd())
	cmd.AddCommand(newValidateCmd())
	cmd.AddCommand(newCleanCmd())
	cmd.AddCommand(newUpdateCmd())
	cmd.AddCommand(newMCPCmd(Version))
	cmd.AddCommand(newVersionCmd())
	cmd.AddCommand(newInfoCmd())
	cmd.AddCommand(newImportCmd())
	cmd.AddCommand(newExportCmd())

	return cmd
}
