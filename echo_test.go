package main

import (
	"bytes"
	"io"
	"log"
	"testing"
	"time"

	"github.com/agentio/echo-go/internal/commands"
)

func TestServices(t *testing.T) {
	go func() {
		serveCmd := commands.Cmd()
		serveCmd.SetArgs([]string{"serve", "grpc", "--socket", "@echotest"})
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
			Args: []string{"call", "get", "--address", "unix:@echotest"},
			Expected: `{"text":"Go echo get: hello"}
`,
		},
		{
			Args: []string{"call", "collect", "--address", "unix:@echotest"},
			Expected: `{"text":"Go echo collect: hello 0 hello 1 hello 2"}
`,
		},
		{
			Args: []string{"call", "expand", "--address", "unix:@echotest"},
			Expected: `{"text":"Go echo expand (0): 1"}
{"text":"Go echo expand (1): 2"}
{"text":"Go echo expand (2): 3"}
`,
		},
		{
			Args: []string{"call", "stream", "--address", "unix:@echotest"},
			Expected: `{"text":"Go echo stream (1): hello 0"}
{"text":"Go echo stream (2): hello 1"}
{"text":"Go echo stream (3): hello 2"}
{"text":"Go echo stream (4): hello 3"}
{"text":"Go echo stream (5): hello 4"}
{"text":"Go echo stream (6): hello 5"}
`,
		},
	}
	for _, test := range tests {
		cmd := commands.Cmd()
		buffer := new(bytes.Buffer)
		cmd.SetOut(buffer)
		cmd.SetArgs(test.Args)
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
