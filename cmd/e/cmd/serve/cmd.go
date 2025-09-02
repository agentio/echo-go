package serve

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

	pb "github.com/agentio/echo/genproto/echopb"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Run the echo server",
		RunE:  action,
		Args:  cobra.NoArgs,
	}
	return cmd
}

func action(cmd *cobra.Command, args []string) error {
	var err error
	var lis net.Listener
	var grpcServer *grpc.Server

	lis, err = net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer = grpc.NewServer()
	var echoServer EchoServer

	pb.RegisterEchoServer(grpcServer, &echoServer)
	return grpcServer.Serve(lis)
}

type EchoServer struct {
	pb.UnimplementedEchoServer
}

// requests are immediately returned, no inbound or outbound streaming
func (s *EchoServer) Get(ctx context.Context, request *pb.EchoRequest) (*pb.EchoResponse, error) {
	fmt.Printf("Get received: %s\n", request.Text)
	response := &pb.EchoResponse{}
	response.Text = "Go echo get: " + request.Text
	return response, nil
}

// requests stream in and are immediately streamed out
func (s *EchoServer) Stream(stream pb.Echo_StreamServer) error {
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
		response := &pb.EchoResponse{}
		response.Text = fmt.Sprintf("Go echo stream (%d): %s", count, request.Text)
		count++
		if err := stream.Send(response); err != nil {
			return err
		}
	}
}

// requests stream in, are appended together, and are returned in a single response when the input is closed
func (s *EchoServer) Collect(stream pb.Echo_CollectServer) error {
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
	response := &pb.EchoResponse{}
	response.Text = fmt.Sprintf("Go echo collect: %s", strings.Join(parts, " "))
	if err := stream.SendAndClose(response); err != nil {
		return err
	}
	return nil
}

// a single request is accepted and split into parts which are individually returned with a time delay
func (s *EchoServer) Expand(request *pb.EchoRequest, stream pb.Echo_ExpandServer) error {
	fmt.Printf("Expand received: %s\n", request.Text)
	parts := strings.Split(request.Text, " ")
	for i, part := range parts {
		response := &pb.EchoResponse{}
		response.Text = fmt.Sprintf("Go echo expand (%d): %s", i, part)
		if err := stream.Send(response); err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}
