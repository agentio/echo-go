package serve

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"github.com/agentio/echo/genproto/echopb"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Run the echo server",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			lis, err := net.Listen("tcp", ":8080")
			if err != nil {
				log.Fatalf("failed to listen: %v", err)
			}
			grpcServer := grpc.NewServer()
			var echoServer EchoServer
			echopb.RegisterEchoServer(grpcServer, &echoServer)
			return grpcServer.Serve(lis)
		},
	}
	return cmd
}

type EchoServer struct {
	echopb.UnimplementedEchoServer
}

// requests are immediately returned, no inbound or outbound streaming
func (s *EchoServer) Get(ctx context.Context, request *echopb.EchoRequest) (*echopb.EchoResponse, error) {
	fmt.Printf("Get received: %s\n", request.Text)
	response := &echopb.EchoResponse{}
	response.Text = "Go echo get: " + request.Text
	return response, nil
}

// requests stream in and are immediately streamed out
func (s *EchoServer) Stream(stream echopb.Echo_StreamServer) error {
	count := 0
	for {
		request, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		fmt.Printf("Stream received: %s\n", request.Text)
		response := &echopb.EchoResponse{}
		response.Text = fmt.Sprintf("Go echo stream (%d): %s", count, request.Text)
		count++
		if err := stream.Send(response); err != nil {
			return err
		}
	}
}

// requests stream in, are appended together, and are returned in a single response when the input is closed
func (s *EchoServer) Collect(stream echopb.Echo_CollectServer) error {
	parts := []string{}
	for {
		request, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		fmt.Printf("Collect received: %s\n", request.Text)
		parts = append(parts, request.Text)
	}
	response := &echopb.EchoResponse{}
	response.Text = fmt.Sprintf("Go echo collect: %s", strings.Join(parts, " "))
	if err := stream.SendAndClose(response); err != nil {
		return err
	}
	return nil
}

// a single request is accepted and split into parts which are individually returned with a time delay
func (s *EchoServer) Expand(request *echopb.EchoRequest, stream echopb.Echo_ExpandServer) error {
	fmt.Printf("Expand received: %s\n", request.Text)
	parts := strings.Split(request.Text, " ")
	for i, part := range parts {
		response := &echopb.EchoResponse{}
		response.Text = fmt.Sprintf("Go echo expand (%d): %s", i, part)
		if err := stream.Send(response); err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}
