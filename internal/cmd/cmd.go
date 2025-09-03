package cmd

import (
	"github.com/agentio/echo-go/internal/cmd/call"
	"github.com/agentio/echo-go/internal/cmd/serve"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "echo-go",
		Short: "Echo client and server",
	}
	cmd.AddCommand(serve.Cmd())
	cmd.AddCommand(call.Cmd())
	return cmd
}
