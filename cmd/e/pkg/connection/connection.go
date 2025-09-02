package connection

import (
	"net/http"

	"connectrpc.com/connect"
	"github.com/agentio/echo/genproto/echopb/echopbconnect"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func NewGRPCConnection(address string, useTLS bool) (*grpc.ClientConn, error) {
	if !useTLS {
		if address == "" {
			address = "localhost:8080"
		}
		return grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		if address == "" {
			address = "localhost:443"
		}
		return grpc.NewClient(address, grpc.WithTransportCredentials(credentials.NewTLS(nil)))
	}
}

func NewConnectEchoClient(address string, useTLS bool, stack string) (echopbconnect.EchoClient, error) {
	if address == "" {
		address = "localhost"
	}
	var url string
	if useTLS {
		url = "https://" + address
	} else {
		url = "http://" + address
	}
	options := []connect.ClientOption{}
	switch stack {
	case "connect-grpc":
		options = append(options, connect.WithGRPC())
	case "connect-grpc-web":
		options = append(options, connect.WithGRPCWeb())
	default:
	}
	return echopbconnect.NewEchoClient(
		http.DefaultClient,
		url,
		options...,
	), nil
}
