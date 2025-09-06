package main

import (
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
		log.Printf("%s", err)
	}()

	time.Sleep(100 * time.Millisecond)
	{
		getCmd := commands.Cmd()
		getCmd.SetArgs([]string{"call", "get", "--address", "unix:@echo", "-n", "1000"})
		err := getCmd.Execute()
		if err != nil {
			t.Errorf("%s", err)
		}
	}
	{
		collectCmd := commands.Cmd()
		collectCmd.SetArgs([]string{"call", "collect", "--address", "unix:@echo", "-n", "1000"})
		err := collectCmd.Execute()
		if err != nil {
			t.Errorf("%s", err)
		}
	}
	{
		expandCmd := commands.Cmd()
		expandCmd.SetArgs([]string{"call", "expand", "--address", "unix:@echo", "-n", "1000"})
		err := expandCmd.Execute()
		if err != nil {
			t.Errorf("%s", err)
		}
	}
	{
		streamCmd := commands.Cmd()
		streamCmd.SetArgs([]string{"call", "stream", "--address", "unix:@echo", "-n", "1000"})
		err := streamCmd.Execute()
		if err != nil {
			t.Errorf("%s", err)
		}
	}
}
