package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewProjectRemoveCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:     "remove [PROJECT_ID_OR_NAME]",
		Aliases: []string{"rm"},
		Short:   "Remove a project and all its time entries",
		Long:    `Remove a project and all its associated time entries. This action cannot be undone. You can specify either the project ID (numeric) or name.`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectIDOrName := args[0]

			if !force {
				fmt.Printf("This will delete project '%s' and ALL its time entries. Are you sure? (y/N): ", projectIDOrName)
				var response string
				fmt.Scanln(&response)
				if response != "y" && response != "Y" {
					fmt.Println("Operation cancelled.")
					return nil
				}
			}

			ctx := cmd.Context()
			err := timeService.RemoveProjectByIDOrName(ctx, projectIDOrName)
			if err != nil {
				return fmt.Errorf("failed to remove project: %w", err)
			}

			fmt.Printf("Project '%s' and all its time entries have been removed.\n", projectIDOrName)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Skip confirmation prompt")

	return cmd
}
