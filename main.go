package main

import (
	"os"

	"github.com/agentio/echo-go/internal/commands"
)

func main() {
	if err := commands.Cmd().Execute(); err != nil {
		os.Exit(1)
	}
}
