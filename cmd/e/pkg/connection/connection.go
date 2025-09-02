package connection

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func NewConnection(address string, useTLS bool) (*grpc.ClientConn, error) {
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
