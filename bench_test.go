package main

import (
	"fmt"
	"io"
	"log"
	"testing"
	"time"

	"github.com/agentio/echo-go/genproto/echopb"
	"github.com/agentio/echo-go/internal/commands"
	"github.com/agentio/echo-go/internal/connection"
)

func runServerOnSocket(socket string) {
	serveCmd := commands.Cmd()
	serveCmd.SetArgs([]string{"serve", "grpc", "--socket", socket})
	if err := serveCmd.Execute(); err != nil {
		log.Printf("%s", err)
	}
}

func runServerOnPort(port string) {
	serveCmd := commands.Cmd()
	serveCmd.SetArgs([]string{"serve", "grpc", "--port", port})
	if err := serveCmd.Execute(); err != nil {
		log.Printf("%s", err)
	}
}

func pause() {
	time.Sleep(100 * time.Millisecond)
}

func measureGet(b *testing.B, client echopb.EchoClient) {
	message := "hello"
	for b.Loop() {
		_, err := client.Get(b.Context(), &echopb.EchoRequest{Text: message})
		if err != nil {
			b.FailNow()
		}
	}
}

func measureExpand(b *testing.B, client echopb.EchoClient) {
	message := "1 2 3"
	for b.Loop() {
		stream, err := client.Expand(b.Context(), &echopb.EchoRequest{Text: message})
		if err != nil {
			b.FailNow()
		}
		for {
			_, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				b.FailNow()
			}
		}
	}
}

func measureCollect(b *testing.B, client echopb.EchoClient) {
	message := "hello"
	for b.Loop() {
		stream, err := client.Expand(b.Context(), &echopb.EchoRequest{Text: message})
		if err != nil {
			b.FailNow()
		}
		for {
			_, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				b.FailNow()
			}
		}
	}
}

func measureUpdate(b *testing.B, client echopb.EchoClient) {
	message := "hello"
	count := 6
	for b.Loop() {
		stream, err := client.Update(b.Context())
		if err != nil {
			b.FailNow()
		}
		waitc := make(chan struct{})
		go func() {
			for {
				_, err := stream.Recv()
				if err == io.EOF {
					close(waitc)
					return
				}
				if err != nil {
					log.Printf("%s", err)
				}
			}
		}()
		for i := 0; i < count; i++ {
			if err := stream.Send(&echopb.EchoRequest{
				Text: fmt.Sprintf("%s %d", message, i),
			}); err != nil {
				b.FailNow()
			}
		}
		stream.CloseSend()
		<-waitc
	}
}

func BenchmarkEchoGetSocket(b *testing.B) {
	socket := "@echoget"
	go runServerOnSocket(socket)
	pause()
	conn, err := connection.NewGRPCConnection("unix:"+socket, false)
	if err != nil {
		b.FailNow()
	}
	defer conn.Close()
	measureGet(b, echopb.NewEchoClient(conn))
}

func BenchmarkEchoExpandSocket(b *testing.B) {
	socket := "@echoexpand"
	go runServerOnSocket(socket)
	pause()
	conn, err := connection.NewGRPCConnection("unix:"+socket, false)
	if err != nil {
		b.FailNow()
	}
	defer conn.Close()
	measureExpand(b, echopb.NewEchoClient(conn))
}

func BenchmarkEchoCollectSocket(b *testing.B) {
	socket := "@echocollect"
	go runServerOnSocket(socket)
	pause()
	conn, err := connection.NewGRPCConnection("unix:"+socket, false)
	if err != nil {
		b.FailNow()
	}
	defer conn.Close()
	measureCollect(b, echopb.NewEchoClient(conn))
}

func BenchmarkEchoUpdateSocket(b *testing.B) {
	socket := "@echoupdate"
	go runServerOnSocket(socket)
	pause()
	conn, err := connection.NewGRPCConnection("unix:"+socket, false)
	if err != nil {
		b.FailNow()
	}
	defer conn.Close()
	measureUpdate(b, echopb.NewEchoClient(conn))
}

func BenchmarkEchoGetPort(b *testing.B) {
	port := "21001"
	go runServerOnPort(port)
	pause()
	conn, err := connection.NewGRPCConnection("localhost:"+port, false)
	if err != nil {
		b.FailNow()
	}
	defer conn.Close()
	measureGet(b, echopb.NewEchoClient(conn))
}

func BenchmarkEchoExpandPort(b *testing.B) {
	port := "21002"
	go runServerOnPort(port)
	pause()
	conn, err := connection.NewGRPCConnection("localhost:"+port, false)
	if err != nil {
		b.FailNow()
	}
	defer conn.Close()
	measureExpand(b, echopb.NewEchoClient(conn))
}

func BenchmarkEchoCollectPort(b *testing.B) {
	port := "21003"
	go runServerOnPort(port)
	pause()
	conn, err := connection.NewGRPCConnection("localhost:"+port, false)
	if err != nil {
		b.FailNow()
	}
	defer conn.Close()
	measureCollect(b, echopb.NewEchoClient(conn))
}

func BenchmarkEchoUpdatePort(b *testing.B) {
	port := "21004"
	go runServerOnPort(port)
	pause()
	conn, err := connection.NewGRPCConnection("localhost:"+port, false)
	if err != nil {
		b.FailNow()
	}
	defer conn.Close()
	measureUpdate(b, echopb.NewEchoClient(conn))
}
