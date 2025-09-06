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
		serveCmd.SetArgs([]string{"serve", "grpc", "--socket", "@echo"})
		err := serveCmd.Execute()
		if err != nil {
			log.Printf("failed to read output from buffer: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)
	{
		getCmd := commands.Cmd()
		buffer := new(bytes.Buffer)
		getCmd.SetOut(buffer)
		getCmd.SetArgs([]string{"call", "get", "--address", "unix:@echo"})
		err := getCmd.Execute()
		if err != nil {
			t.Errorf("%s", err)
		}
		out, err := io.ReadAll(buffer)
		if err != nil {
			t.Fatalf("failed to read output from buffer: %v", err)
		}
		expectedOutput := `{"text":"Go echo get: hello"}
`
		if string(out) != expectedOutput {
			t.Errorf("expected output %q, got %q", expectedOutput, string(out))
		}
	}
	{
		collectCmd := commands.Cmd()
		buffer := new(bytes.Buffer)
		collectCmd.SetOut(buffer)
		collectCmd.SetArgs([]string{"call", "collect", "--address", "unix:@echo"})
		err := collectCmd.Execute()
		if err != nil {
			t.Errorf("%s", err)
		}
		out, err := io.ReadAll(buffer)
		if err != nil {
			t.Fatalf("failed to read output from buffer: %v", err)
		}
		expectedOutput := `{"text":"Go echo collect: hello 0 hello 1 hello 2"}
`
		if string(out) != expectedOutput {
			t.Errorf("expected output %q, got %q", expectedOutput, string(out))
		}
	}
	{
		expandCmd := commands.Cmd()
		buffer := new(bytes.Buffer)
		expandCmd.SetOut(buffer)
		expandCmd.SetArgs([]string{"call", "expand", "--address", "unix:@echo"})
		err := expandCmd.Execute()
		if err != nil {
			t.Errorf("%s", err)
		}
		out, err := io.ReadAll(buffer)
		if err != nil {
			t.Fatalf("failed to read output from buffer: %v", err)
		}
		expectedOutput := `{"text":"Go echo expand (0): 1"}
{"text":"Go echo expand (1): 2"}
{"text":"Go echo expand (2): 3"}
`
		if string(out) != expectedOutput {
			t.Errorf("expected output %q, got %q", expectedOutput, string(out))
		}
	}
	{
		streamCmd := commands.Cmd()
		buffer := new(bytes.Buffer)
		streamCmd.SetOut(buffer)
		streamCmd.SetArgs([]string{"call", "stream", "--address", "unix:@echo"})
		err := streamCmd.Execute()
		if err != nil {
			t.Errorf("%s", err)
		}
		out, err := io.ReadAll(buffer)
		if err != nil {
			t.Fatalf("failed to read output from buffer: %v", err)
		}
		expectedOutput := `{"text":"Go echo stream (1): hello 0"}
{"text":"Go echo stream (2): hello 1"}
{"text":"Go echo stream (3): hello 2"}
{"text":"Go echo stream (4): hello 3"}
{"text":"Go echo stream (5): hello 4"}
{"text":"Go echo stream (6): hello 5"}
`
		if string(out) != expectedOutput {
			t.Errorf("expected output %q, got %q", expectedOutput, string(out))
		}
	}
}
