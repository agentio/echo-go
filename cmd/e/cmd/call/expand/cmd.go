package expand

import (
	"context"
	"io"
	"log"

	"github.com/agentio/echo/cmd/e/pkg/connection"
	"github.com/agentio/echo/genproto/echopb"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var message string
	var address string
	var useTLS bool
	cmd := &cobra.Command{
		Use:   "expand",
		Short: "Call the expand method",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			conn, err := connection.NewGRPCConnection(address, useTLS)
			if err != nil {
				return err
			}
			defer conn.Close()
			c := echopb.NewEchoClient(conn)
			stream, err := c.Expand(context.Background(), &echopb.EchoRequest{Text: message})
			if err != nil {
				return err
			}
			for {
				in, err := stream.Recv()
				if err == io.EOF {
					// read done.
					return nil
				}
				if err != nil {
					return err
				}
				log.Printf("Received: %s", in.Text)
			}
		},
	}
	cmd.Flags().StringVar(&message, "message", "hello", "message")
	cmd.Flags().StringVar(&address, "address", "", "address of the echo server to use")
	cmd.Flags().BoolVar(&useTLS, "tls", false, "use tls for connections")
	return cmd
}
