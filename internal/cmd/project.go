package cmd

import (
	"github.com/spf13/cobra"
)

func NewProjectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "project",
		Aliases: []string{"p"},
		Short:   "Manage projects",
		Long:    `Manage projects in your time tracking system.`,
	}

	cmd.AddCommand(NewProjectListCmd())
	cmd.AddCommand(NewProjectRemoveCmd())
	cmd.AddCommand(NewProjectTimesCmd())
	cmd.AddCommand(NewProjectTotalCmd())

	return cmd
}
