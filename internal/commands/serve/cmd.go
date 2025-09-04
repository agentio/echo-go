package serve

import (
	"errors"

	"github.com/agentio/echo-go/internal/commands/serve/connectserver"
	"github.com/agentio/echo-go/internal/commands/serve/grpcserver"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var port int
	var socket string
	var verbose bool
	cmd := &cobra.Command{
		Use:   "serve [PROTOCOL]",
		Short: "Run the echo server",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			protocol := ""
			if len(args) > 0 {
				protocol = args[0]
			}
			switch protocol {
			case "", "grpc":
				return grpcserver.Run(port, socket, verbose)
			case "connect":
				return connectserver.Run(port, socket, verbose)
			default:
				return errors.New("please specify 'grpc' or 'connect'")
			}
		},
	}
	cmd.Flags().IntVar(&port, "port", 8080, "server port")
	cmd.Flags().StringVar(&socket, "socket", "", "server socket")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose")
	return cmd
}
