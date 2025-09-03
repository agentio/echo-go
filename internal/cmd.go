package internal

import (
	"github.com/agentio/echo-go/internal/call"
	"github.com/agentio/echo-go/internal/serve"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "e",
		Short: "Echo client and server",
	}
	cmd.AddCommand(serve.Cmd())
	cmd.AddCommand(call.Cmd())
	return cmd
}
