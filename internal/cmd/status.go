package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

func NewStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show the currently active time tracking session",
		Long:  `Show information about the currently active time tracking session, including project name, start time, and current duration.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			activeEntry, err := timeService.GetActiveEntry(ctx)
			if err != nil {
				return fmt.Errorf("failed to get active entry: %w", err)
			}

			if activeEntry == nil {
				fmt.Println("No active time tracking session.")
				return nil
			}

			// Calculate current duration
			currentDuration := time.Since(activeEntry.StartTime)
			durationStr := timeService.FormatDuration(currentDuration)

			fmt.Printf("Active session:\n\n")
			fmt.Printf("Project: %s\n", activeEntry.Project.Name)
			fmt.Printf("Started: %s\n", activeEntry.StartTime.Format("2006-01-02 15:04:05"))
			fmt.Printf("Duration: %s\n", durationStr)

			return nil
		},
	}

	return cmd
}
