package main

import (
	"bytes"
	"io"
	"log"
	"testing"
	"time"

	"github.com/agentio/echo-go/internal/commands"
)

func TestSocketGrpcServiceGrpcClient(t *testing.T) {
	test_service(t,
		[]string{"serve", "grpc", "--socket", "@echotest"},
		[]string{"--address", "unix:@echotest", "--stack", "grpc"},
	)
}

func TestSocketConnectServiceGrpcClient(t *testing.T) {
	test_service(t,
		[]string{"serve", "connect", "--socket", "@echotestconnect"},
		[]string{"--address", "unix:@echotestconnect", "--stack", "grpc"},
	)
}

func TestLocalGrpcServiceGrpcClient(t *testing.T) {
	const port = "19876"
	test_service(t,
		[]string{"serve", "grpc", "--port", port},
		[]string{"--address", "localhost:" + port, "--stack", "grpc"},
	)
}

func TestLocalGrpcServiceConnectGrpcClient(t *testing.T) {
	const port = "19875"
	test_service(t,
		[]string{"serve", "grpc", "--port", port},
		[]string{"--address", "localhost:" + port, "--stack", "connect-grpc"},
	)
}

func TestLocalConnectServiceGrpcClient(t *testing.T) {
	const port = "19874"
	test_service(t,
		[]string{"serve", "connect", "--port", port},
		[]string{"--address", "localhost:" + port, "--stack", "grpc"},
	)
}

func TestLocalConnectServiceConnectClient(t *testing.T) {
	const port = "19873"
	test_service(t,
		[]string{"serve", "connect", "--port", port},
		[]string{"--address", "localhost:" + port, "--stack", "connect"},
	)
}

func TestLocalConnectServiceConnectGrpcClient(t *testing.T) {
	const port = "19872"
	test_service(t,
		[]string{"serve", "connect", "--port", port},
		[]string{"--address", "localhost:" + port, "--stack", "connect-grpc"},
	)
}

func TestLocalConnectServiceConnectGrpcWebClient(t *testing.T) {
	const port = "19871"
	test_service(t,
		[]string{"serve", "connect", "--port", port},
		[]string{"--address", "localhost:" + port, "--stack", "connect-grpc-web"},
	)
}

func test_service(t *testing.T, serverArgs, clientArgs []string) {
	go func() {
		serveCmd := commands.Cmd()
		serveCmd.SetArgs(serverArgs)
		err := serveCmd.Execute()
		if err != nil {
			log.Printf("failed to read output from buffer: %v", err)
		}
	}()
	time.Sleep(100 * time.Millisecond)
	tests := []struct {
		Args     []string
		Expected string
	}{
		{
			Args:     []string{"call", "get"},
			Expected: expected_get,
		},
		{
			Args:     []string{"call", "collect"},
			Expected: expected_collect,
		},
		{
			Args:     []string{"call", "expand"},
			Expected: expected_expand,
		},
		{
			Args:     []string{"call", "stream"},
			Expected: expected_stream,
		},
	}
	for _, test := range tests {
		cmd := commands.Cmd()
		buffer := new(bytes.Buffer)
		cmd.SetOut(buffer)
		cmd.SetArgs(append(test.Args, clientArgs...))
		err := cmd.Execute()
		if err != nil {
			t.Errorf("%s", err)
		}
		out, err := io.ReadAll(buffer)
		if err != nil {
			t.Fatalf("failed to read output: %v", err)
		}
		if string(out) != test.Expected {
			t.Errorf("expected %q, got %q", test.Expected, string(out))
		}
	}
}

const expected_get = `{"text":"Go echo get: hello"}
`
const expected_collect = `{"text":"Go echo collect: hello 0 hello 1 hello 2"}
`
const expected_expand = `{"text":"Go echo expand (0): 1"}
{"text":"Go echo expand (1): 2"}
{"text":"Go echo expand (2): 3"}
`
const expected_stream = `{"text":"Go echo stream (1): hello 0"}
{"text":"Go echo stream (2): hello 1"}
{"text":"Go echo stream (3): hello 2"}
{"text":"Go echo stream (4): hello 3"}
{"text":"Go echo stream (5): hello 4"}
{"text":"Go echo stream (6): hello 5"}
`
