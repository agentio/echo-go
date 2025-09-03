package get

import (
	"fmt"
	"log"

	"connectrpc.com/connect"
	"github.com/agentio/echo/cmd/e/pkg/connection"
	"github.com/agentio/echo/genproto/echopb"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var message string
	var address string
	var useTLS bool
	var stack string
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Call the get method",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			switch stack {
			case "grpc":
				conn, err := connection.NewGRPCConnection(address, useTLS)
				if err != nil {
					return err
				}
				defer conn.Close()
				client := echopb.NewEchoClient(conn)
				response, err := client.Get(cmd.Context(), &echopb.EchoRequest{Text: message})
				if err != nil {
					return err
				}
				log.Printf("Received: %s", response.Text)
				return nil
			case "connect", "connect-grpc", "connect-grpc-web":
				client, err := connection.NewConnectEchoClient(address, useTLS, stack)
				response, err := client.Get(cmd.Context(), connect.NewRequest(&echopb.EchoRequest{Text: message}))
				if err != nil {
					return err
				}
				log.Printf("Received: %s", response.Msg.Text)
				return nil
			default:
				return fmt.Errorf("unsupported stack: %s", stack)
			}
		},
	}
	cmd.Flags().StringVar(&message, "message", "hello", "message")
	cmd.Flags().StringVar(&address, "address", "", "address of the echo server to use")
	cmd.Flags().BoolVar(&useTLS, "tls", false, "use tls for connections")
	cmd.Flags().StringVar(&stack, "stack", "grpc", "stack to use to connect")
	return cmd
}
