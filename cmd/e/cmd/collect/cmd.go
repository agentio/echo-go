package collect

import (
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collect",
		Short: "Call the collect method",
		RunE:  action,
		Args:  cobra.NoArgs,
	}
	return cmd
}

func action(cmd *cobra.Command, args []string) error {
	return nil
}
