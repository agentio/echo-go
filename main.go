package main

import (
	"os"

	"github.com/agentio/echo-go/internal"
)

func main() {
	if err := internal.Cmd().Execute(); err != nil {
		os.Exit(1)
	}
}
