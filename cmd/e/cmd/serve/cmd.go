package serve

import (
	"github.com/agentio/echo/cmd/e/cmd/serve/grpcserver"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Run the echo server",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			port := 8080
			return grpcserver.Run(port)
		},
	}
	return cmd
}
