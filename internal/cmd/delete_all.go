package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewDeleteAllCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "delete-all",
		Short: "Delete all time tracking data",
		Long:  `Delete all time tracking data including all time entries and projects. This action cannot be undone.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !force {
				fmt.Print("This will delete ALL time tracking data. Are you sure? (y/N): ")
				var response string
				fmt.Scanln(&response)
				if response != "y" && response != "Y" {
					fmt.Println("Operation cancelled.")
					return nil
				}
			}

			ctx := cmd.Context()
			err := timeService.ClearAllData(ctx)
			if err != nil {
				return fmt.Errorf("failed to delete all data: %w", err)
			}

			fmt.Println("All time tracking data has been deleted.")
			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Skip confirmation prompt")

	return cmd
}
