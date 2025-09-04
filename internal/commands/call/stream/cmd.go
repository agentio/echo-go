package stream

import (
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/agentio/echo-go/genproto/echopb"
	"github.com/agentio/echo-go/internal/connection"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var message string
	var address string
	var useTLS bool
	var count int
	var stack string
	cmd := &cobra.Command{
		Use:   "stream",
		Short: "Call the stream method",
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
				stream, err := client.Stream(cmd.Context())
				if err != nil {
					return err
				}
				waitc := make(chan struct{})
				go func() {
					for {
						in, err := stream.Recv()
						if err == io.EOF {
							close(waitc)
							return
						}
						if err != nil {
							log.Fatalf("Failed to receive an echo : %v", err)
						}
						log.Printf("Received: %s", in.Text)
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
				return nil
			case "connect", "connect-grpc", "connect-grpc-web":
				client, err := connection.NewConnectEchoClient(address, useTLS, stack)
				if err != nil {
					return nil
				}
				stream := client.Stream(cmd.Context())
				waitc := make(chan struct{})
				go func() {
					for {
						in, err := stream.Receive()
						if errors.Is(err, io.EOF) {
							close(waitc)
							return
						}
						if err != nil {
							log.Fatalf("Failed to receive an echo : %v", err)
						}
						log.Printf("Received: %s", in.Text)
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
	return cmd
}
