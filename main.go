package main

import (
	"os"

	"github.com/agentio/echo-go/internal/cmd"
)

func main() {
	if err := cmd.Cmd().Execute(); err != nil {
		os.Exit(1)
	}
}
