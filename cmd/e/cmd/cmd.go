package cmd

import (
	"github.com/agentio/echo/cmd/e/cmd/collect"
	"github.com/agentio/echo/cmd/e/cmd/expand"
	"github.com/agentio/echo/cmd/e/cmd/get"
	"github.com/agentio/echo/cmd/e/cmd/serve"
	"github.com/agentio/echo/cmd/e/cmd/stream"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "e",
		Short: "Echo client and server",
	}
	cmd.AddCommand(serve.Cmd())
	cmd.AddCommand(get.Cmd())
	cmd.AddCommand(expand.Cmd())
	cmd.AddCommand(collect.Cmd())
	cmd.AddCommand(stream.Cmd())
	return cmd
}
