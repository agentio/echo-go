package grpcserver

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"github.com/agentio/echo-go/genproto/echopb"
	"google.golang.org/grpc"
)

func Run(port int, socket string) error {
	var lis net.Listener
	var err error
	if socket != "" {
		lis, err = net.Listen("unix", socket)
		log.Printf("serving on %s", socket)
	} else {
		lis, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
		log.Printf("serving on %d", port)
	}
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	echopb.RegisterEchoServer(grpcServer, &echoServer{})
	return grpcServer.Serve(lis)
}

type echoServer struct {
	echopb.UnimplementedEchoServer
}

// Immediately returns an echo of a request.
func (s *echoServer) Get(ctx context.Context, request *echopb.EchoRequest) (*echopb.EchoResponse, error) {
	fmt.Printf("Get received: %s\n", request.Text)
	return &echopb.EchoResponse{
		Text: "Go echo get: " + request.Text,
	}, nil
}

// Splits a request into words and returns each word in a stream of messages.
func (s *echoServer) Expand(request *echopb.EchoRequest, stream echopb.Echo_ExpandServer) error {
	fmt.Printf("Expand received: %s\n", request.Text)
	parts := strings.Split(request.Text, " ")
	for i, part := range parts {
		if err := stream.Send(&echopb.EchoResponse{
			Text: fmt.Sprintf("Go echo expand (%d): %s", i, part),
		}); err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

// Collects a stream of messages and returns them concatenated when the caller closes.
func (s *echoServer) Collect(stream echopb.Echo_CollectServer) error {
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
	if err := stream.SendAndClose(&echopb.EchoResponse{
		Text: fmt.Sprintf("Go echo collect: %s", strings.Join(parts, " ")),
	}); err != nil {
		return err
	}
	return nil
}

// Streams back messages as they are received in an input stream.
func (s *echoServer) Stream(stream echopb.Echo_StreamServer) error {
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
		count++
		if err := stream.Send(&echopb.EchoResponse{
			Text: fmt.Sprintf("Go echo stream (%d): %s", count, request.Text),
		}); err != nil {
			return err
		}
	}
}
