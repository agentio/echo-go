package get

import (
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Call the get method",
		RunE:  action,
		Args:  cobra.NoArgs,
	}
	return cmd
}

func action(cmd *cobra.Command, args []string) error {
	return nil
}
