package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewProjectRemoveCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:     "remove [project-name]",
		Aliases: []string{"rm"},
		Short:   "Remove a project and all its time entries",
		Long:    `Remove a project and all its associated time entries. This action cannot be undone.`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectName := args[0]

			if !force {
				fmt.Printf("This will delete project '%s' and ALL its time entries. Are you sure? (y/N): ", projectName)
				var response string
				fmt.Scanln(&response)
				if response != "y" && response != "Y" {
					fmt.Println("Operation cancelled.")
					return nil
				}
			}

			err := timeService.RemoveProject(projectName)
			if err != nil {
				return fmt.Errorf("failed to remove project: %w", err)
			}

			fmt.Printf("Project '%s' and all its time entries have been removed.\n", projectName)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Skip confirmation prompt")

	return cmd
}
