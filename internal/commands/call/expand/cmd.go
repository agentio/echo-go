package expand

import (
	"fmt"
	"io"
	"log"
	"time"

	"connectrpc.com/connect"
	"github.com/agentio/echo-go/genproto/echopb"
	"github.com/agentio/echo-go/internal/connection"
	"github.com/agentio/echo-go/internal/track"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var message string
	var address string
	var useTLS bool
	var stack string
	var n int
	cmd := &cobra.Command{
		Use:   "expand",
		Short: "Call the expand method",
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
				defer track.Measure(time.Now(), "expand", n)
				for j := 0; j < n; j++ {
					stream, err := client.Expand(cmd.Context(), &echopb.EchoRequest{Text: message})
					if err != nil {
						return err
					}
					for {
						in, err := stream.Recv()
						if err == io.EOF {
							break
						}
						if err != nil {
							return err
						}
						if n == 1 {
							log.Printf("Received: %s", in.Text)
						}
					}
				}
				return nil
			case "connect", "connect-grpc", "connect-grpc-web":
				client, err := connection.NewConnectEchoClient(address, useTLS, stack)
				if err != nil {
					return nil
				}
				defer track.Measure(time.Now(), "expand", n)
				for j := 0; j < n; j++ {
					stream, err := client.Expand(cmd.Context(), connect.NewRequest(&echopb.EchoRequest{Text: message}))
					if err != nil {
						return err
					}
					for {
						running := stream.Receive()
						if !running {
							break
						}
						in := stream.Msg()
						if n == 1 {
							log.Printf("Received: %s", in.Text)
						}
					}
				}
				return nil
			default:
				return fmt.Errorf("unsupported stack: %s", stack)
			}
		},
	}
	cmd.Flags().StringVarP(&message, "message", "m", "1 2 3", "message")
	cmd.Flags().StringVarP(&address, "address", "a", "", "address of the echo server to use")
	cmd.Flags().BoolVar(&useTLS, "tls", false, "use tls for connections")
	cmd.Flags().StringVar(&stack, "stack", "grpc", "stack to use to connect")
	cmd.Flags().IntVarP(&n, "number", "n", 1, "number of times to call the method")
	return cmd
}
