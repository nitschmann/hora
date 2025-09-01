package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewClearCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "clear",
		Short: "Clear all time tracking data",
		Long:  `Clear all time tracking data including all time entries and projects. This action cannot be undone.`,
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

			err := timeService.ClearAllData()
			if err != nil {
				return fmt.Errorf("failed to clear data: %w", err)
			}

			fmt.Println("All time tracking data has been cleared.")
			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Skip confirmation prompt")

	return cmd
}
