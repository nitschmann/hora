package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nitschmann/hora/internal/backgroundtracker"
)

func NewStopCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop the current time tracking session",
		Long:  `Stop the currently active time tracking session and display the duration.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			entry, err := timeService.StopTracking(ctx)
			if err != nil {
				fmt.Println("Failed to stop time tracking")
				return err
			}

			durationStr := timeService.FormatDuration(*entry.Duration)

			fmt.Printf("Stopped tracking time for project: %s\n", entry.Project.Name)
			fmt.Printf("Duration: %s\n", durationStr)

			if backgroundtracker.IsRunning() {
				err = backgroundtracker.Stop()
				if err != nil {
					fmt.Printf("Warning: Failed to stop background tracker: %v\n", err)
					return err
				}
			}

			return nil
		},
	}

	return cmd
}
