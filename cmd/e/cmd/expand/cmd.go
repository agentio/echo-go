package expand

import (
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "expand",
		Short: "Call the expand method",
		RunE:  action,
		Args:  cobra.NoArgs,
	}
	return cmd
}

func action(cmd *cobra.Command, args []string) error {
	return nil
}
