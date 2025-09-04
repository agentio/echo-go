package commands

import (
	"github.com/agentio/echo-go/internal/commands/call"
	"github.com/agentio/echo-go/internal/commands/serve"
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
