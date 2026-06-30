package main

import (
	"github.com/mmadfox/swag2mcp/internal/commands"
)

func main() {
	if err := commands.NewRootCmd().Execute(); err != nil {
		panic(err)
	}
}
