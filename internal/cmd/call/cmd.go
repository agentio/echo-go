package call

import (
	"github.com/agentio/echo-go/internal/cmd/call/collect"
	"github.com/agentio/echo-go/internal/cmd/call/expand"
	"github.com/agentio/echo-go/internal/cmd/call/get"
	"github.com/agentio/echo-go/internal/cmd/call/stream"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "call",
		Short: "Call echo methods",
	}
	cmd.AddCommand(get.Cmd())
	cmd.AddCommand(expand.Cmd())
	cmd.AddCommand(collect.Cmd())
	cmd.AddCommand(stream.Cmd())
	return cmd
}
