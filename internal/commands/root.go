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
	cmd.AddCommand(newMCPCmd("v0.1.0"))

	return cmd
}
