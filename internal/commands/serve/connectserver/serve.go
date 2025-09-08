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

	"connectrpc.com/connect"
	"github.com/agentio/echo-go/genproto/echopb"
	"github.com/agentio/echo-go/genproto/echopb/echopbconnect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var verbose bool

func Run(port int, socket string, _verbose bool) error {
	verbose = _verbose
	mux := http.NewServeMux()
	mux.Handle(echopbconnect.NewEchoHandler(&echoServer{}))
	var socketListener net.Listener
	var err error
	if socket != "" {
		socketListener, err = net.Listen("unix", socket)
		if verbose {
			log.Printf("serving on %s", socket)
		}
	} else {
		socketListener, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
		if verbose {
			log.Printf("serving on %d", port)
		}
	}
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
	if verbose {
		log.Printf("Get received: %s", req.Msg.Text)
	}
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
	if verbose {
		log.Printf("Expand received: %s", req.Msg.Text)
	}
	parts := strings.Split(req.Msg.Text, " ")
	for i, part := range parts {
		if err := stream.Send(&echopb.EchoResponse{
			Text: fmt.Sprintf("Go echo expand (%d): %s", i, part),
		}); err != nil {
			return err
		}
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
			break
		}
		request := stream.Msg()
		if verbose {
			log.Printf("Collect received: %s", request.Text)
		}
		parts = append(parts, request.Text)
	}
	return connect.NewResponse(
		&echopb.EchoResponse{
			Text: fmt.Sprintf("Go echo collect: %s", strings.Join(parts, " ")),
		},
	), nil
}

// Streams back messages as they are received in an input stream.
func (s echoServer) Update(
	ctx context.Context,
	stream *connect.BidiStream[echopb.EchoRequest, echopb.EchoResponse],
) error {
	count := 0
	for {
		request, err := stream.Receive()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		if verbose {
			log.Printf("Stream received: %s", request.Text)
		}
		count++
		if err := stream.Send(&echopb.EchoResponse{
			Text: fmt.Sprintf("Go echo update (%d): %s", count, request.Text),
		}); err != nil {
			return err
		}
	}
}
