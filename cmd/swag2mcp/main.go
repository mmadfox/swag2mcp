package main

import (
	"fmt"
	"os"

	"github.com/mmadfox/swag2mcp/internal/commands"
)

func main() {
	if err := commands.NewRootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
