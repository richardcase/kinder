package internal

import (
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kinder",
		Short: "A CLI to do stuff with kind",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(newKubeconfigCommand())
	cmd.AddCommand(newDeleteCommand())

	return cmd
}
