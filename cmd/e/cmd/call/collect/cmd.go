package collect

import (
	"context"
	"fmt"
	"log"

	"github.com/agentio/echo/cmd/e/pkg/connection"
	"github.com/agentio/echo/genproto/echopb"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var message string
	var address string
	var useTLS bool
	var count int
	cmd := &cobra.Command{
		Use:   "collect",
		Short: "Call the collect method",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			conn, err := connection.NewConnection(address, useTLS)
			if err != nil {
				return err
			}
			defer conn.Close()
			c := echopb.NewEchoClient(conn)
			stream, err := c.Collect(context.Background())
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
			log.Printf("Received: %s", response.Text)
			return nil
		},
	}
	cmd.Flags().StringVar(&message, "message", "hello", "message")
	cmd.Flags().StringVar(&address, "address", "", "address of the echo server to use")
	cmd.Flags().BoolVar(&useTLS, "tls", false, "use tls for connections")
	cmd.Flags().IntVar(&count, "count", 3, "number of messages to send")
	return cmd
}
