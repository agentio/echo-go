package update

import (
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/agentio/echo-go/genproto/echopb"
	"github.com/agentio/echo-go/internal/connection"
	"github.com/agentio/echo-go/internal/track"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

func Cmd() *cobra.Command {
	var message string
	var address string
	var useTLS bool
	var count int
	var stack string
	var n int
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Call the update method",
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
				defer track.Measure(time.Now(), "stream", n, cmd.OutOrStdout())
				for j := 0; j < n; j++ {
					stream, err := client.Update(cmd.Context())
					if err != nil {
						return err
					}
					waitc := make(chan struct{})
					go func() {
						for {
							response, err := stream.Recv()
							if err == io.EOF {
								close(waitc)
								return
							}
							if err != nil {
								log.Fatalf("Failed to receive an echo : %v", err)
							}
							if n == 1 {
								body, err := protojson.Marshal(response)
								if err != nil {
									close(waitc)
									return
								}
								_, _ = cmd.OutOrStdout().Write(body)
								_, _ = cmd.OutOrStdout().Write([]byte("\n"))
							}
						}
					}()
					for i := 0; i < count; i++ {
						if err := stream.Send(&echopb.EchoRequest{
							Text: fmt.Sprintf("%s %d", message, i),
						}); err != nil {
							return err
						}
					}
					stream.CloseSend()
					<-waitc
				}
				return nil
			case "connect", "connect-grpc", "connect-grpc-web":
				client, err := connection.NewConnectEchoClient(address, useTLS, stack)
				if err != nil {
					return nil
				}
				defer track.Measure(time.Now(), "stream", n, cmd.OutOrStdout())
				for j := 0; j < n; j++ {
					stream := client.Update(cmd.Context())
					waitc := make(chan struct{})
					go func() {
						for {
							response, err := stream.Receive()
							if errors.Is(err, io.EOF) {
								close(waitc)
								return
							}
							if err != nil {
								log.Fatalf("Failed to receive an echo : %v", err)
							}
							if n == 1 {
								body, err := protojson.Marshal(response)
								if err != nil {
									close(waitc)
									return
								}
								_, _ = cmd.OutOrStdout().Write(body)
								_, _ = cmd.OutOrStdout().Write([]byte("\n"))
							}
						}
					}()
					for i := 0; i < count; i++ {
						if err := stream.Send(&echopb.EchoRequest{
							Text: fmt.Sprintf("%s %d", message, i),
						}); err != nil {
							return err
						}
					}
					stream.CloseRequest()
					<-waitc
				}
				return nil
			default:
				return fmt.Errorf("unsupported stack: %s", stack)
			}
		},
	}
	cmd.Flags().StringVarP(&message, "message", "m", "hello", "message")
	cmd.Flags().StringVarP(&address, "address", "a", "", "address of the echo server to use")
	cmd.Flags().BoolVar(&useTLS, "tls", false, "use tls for connections")
	cmd.Flags().IntVarP(&count, "count", "c", 6, "number of messages to send")
	cmd.Flags().StringVar(&stack, "stack", "grpc", "stack to use to connect")
	cmd.Flags().IntVarP(&n, "number", "n", 1, "number of times to call the method")
	return cmd
}
