package stream

import (
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stream",
		Short: "Call the stream method",
		RunE:  action,
		Args:  cobra.NoArgs,
	}
	return cmd
}

func action(cmd *cobra.Command, args []string) error {
	return nil
}
