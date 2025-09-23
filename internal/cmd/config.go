package cmd

import "github.com/spf13/cobra"

func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration",
		Long:  `Manage the configuration hora`,
	}

	cmd.AddCommand(NewConfigInitCmd())

	return cmd
}
