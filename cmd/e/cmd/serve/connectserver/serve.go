package connectserver

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"connectrpc.com/connect"
	"github.com/agentio/echo/genproto/echopb"
	"github.com/agentio/echo/genproto/echopb/echopbconnect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// main.go (continued)
func Run(port int) error {
	mux := http.NewServeMux()
	mux.Handle(echopbconnect.NewEchoHandler(&echoServer{}))
	log.Printf("serving on %d", port)
	var socketListener net.Listener
	socketListener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	server := http.Server{
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}
	return server.Serve(socketListener)
}

type echoServer struct {
}

// Immediately returns an echo of a request.
func (s echoServer) Get(ctx context.Context, req *connect.Request[echopb.EchoRequest]) (*connect.Response[echopb.EchoResponse], error) {
	fmt.Printf("Get received: %s\n", req.Msg.Text)
	return connect.NewResponse(&echopb.EchoResponse{
		Text: "Go echo get: " + req.Msg.Text,
	}), nil
}

// Splits a request into words and returns each word in a stream of messages.
func (s echoServer) Expand(
	ctx context.Context,
	req *connect.Request[echopb.EchoRequest],
	stream *connect.ServerStream[echopb.EchoResponse],
) error {
	fmt.Printf("Expand received: %s\n", req.Msg.Text)
	parts := strings.Split(req.Msg.Text, " ")
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
func (s echoServer) Collect(
	ctx context.Context,
	stream *connect.ClientStream[echopb.EchoRequest],
) (*connect.Response[echopb.EchoResponse], error) {
	parts := []string{}
	for {
		running := stream.Receive()
		if !running {
			return connect.NewResponse(
				&echopb.EchoResponse{
					Text: fmt.Sprintf("Go echo collect: %s", strings.Join(parts, " ")),
				},
			), nil
		}
		request := stream.Msg()
		fmt.Printf("Collect received: %s\n", request.Text)
		parts = append(parts, request.Text)
	}
}

// Streams back messages as they are received in an input stream.
func (s echoServer) Stream(
	ctx context.Context,
	stream *connect.BidiStream[echopb.EchoRequest, echopb.EchoResponse],
) error {
	count := 0
	for {
		request, err := stream.Receive()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		fmt.Printf("Stream received: %s\n", request.Text)
		count++
		if err := stream.Send(
			&echopb.EchoResponse{
				Text: fmt.Sprintf("Go echo stream (%d): %s", count, request.Text),
			},
		); err != nil {
			return err
		}
	}
}
