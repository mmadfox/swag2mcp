package main

import (
	"fmt"
	"os"

	"github.com/mmadfox/swag2mcp/internal/commands"
)

func main() {
	if err := runMain(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func runMain() error {
	return commands.NewRootCmd().Execute()
}
