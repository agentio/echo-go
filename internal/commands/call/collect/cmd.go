package collect

import (
	"fmt"
	"log"
	"time"

	"github.com/agentio/echo-go/genproto/echopb"
	"github.com/agentio/echo-go/internal/connection"
	"github.com/agentio/echo-go/internal/track"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var message string
	var address string
	var useTLS bool
	var count int
	var stack string
	var n int
	cmd := &cobra.Command{
		Use:   "collect",
		Short: "Call the collect method",
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
				defer track.Measure(time.Now(), "collect", n)
				for j := 0; j < n; j++ {
					stream, err := client.Collect(cmd.Context())
					if err != nil {
						return err
					}
					for i := 0; i < count; i++ {
						if err := stream.Send(&echopb.EchoRequest{
							Text: fmt.Sprintf("%s %d", message, i),
						}); err != nil {
							return err
						}
					}
					response, err := stream.CloseAndRecv()
					if n == 1 {
						log.Printf("Received: %s", response.Text)
					}
				}
				return nil
			case "connect", "connect-grpc", "connect-grpc-web":
				client, err := connection.NewConnectEchoClient(address, useTLS, stack)
				if err != nil {
					return nil
				}
				defer track.Measure(time.Now(), "collect", n)
				for j := 0; j < n; j++ {
					stream := client.Collect(cmd.Context())
					for i := 0; i < count; i++ {
						if err := stream.Send(&echopb.EchoRequest{
							Text: fmt.Sprintf("%s %d", message, i),
						}); err != nil {
							return err
						}
					}
					response, err := stream.CloseAndReceive()
					if err != nil {
						return err
					}
					if n == 1 {
						log.Printf("Received: %s", response.Msg.Text)
					}
				}
				return nil
			default:
				return fmt.Errorf("unsupported stack: %s", stack)
			}
		},
	}
	cmd.Flags().StringVar(&message, "message", "hello", "message")
	cmd.Flags().StringVar(&address, "address", "", "address of the echo server to use")
	cmd.Flags().BoolVar(&useTLS, "tls", false, "use tls for connections")
	cmd.Flags().IntVar(&count, "count", 3, "number of messages to send")
	cmd.Flags().StringVar(&stack, "stack", "grpc", "stack to use to connect")
	cmd.Flags().IntVarP(&n, "number", "n", 1, "number of times to call the method")
	return cmd
}
