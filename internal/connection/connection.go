package connection

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"github.com/agentio/echo-go/genproto/echopb/echopbconnect"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func NewGRPCConnection(address string, useTLS bool) (*grpc.ClientConn, error) {
	if useTLS {
		if address == "" {
			address = "localhost:443"
		}
		return grpc.NewClient(address, grpc.WithTransportCredentials(credentials.NewTLS(nil)))
	} else {
		if address == "" {
			address = "localhost:8080"
		}
		return grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
}

func NewConnectEchoClient(address string, useTLS bool, stack string) (echopbconnect.EchoClient, error) {
	if address == "" {
		address = "localhost"
	}
	var url string
	var httpClient *http.Client
	if useTLS {
		url = "https://" + address
		httpClient = http.DefaultClient
	} else {
		url = "http://" + address
		httpClient = &http.Client{
			Transport: &http2.Transport{
				AllowHTTP: true,
				DialTLSContext: func(ctx context.Context, network string, addr string, cfg *tls.Config) (net.Conn, error) {
					return net.DialTimeout(network, addr, 5*time.Second)
				},
			},
		}
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
		httpClient,
		url,
		options...,
	), nil
}
