package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewStopCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop the current time tracking session",
		Long:  `Stop the currently active time tracking session and display the duration.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			entry, err := timeService.StopTracking()
			if err != nil {
				return err
			}

			durationStr := timeService.FormatDuration(*entry.Duration)

			fmt.Printf("Stopped tracking time for project: %s\n", entry.Project.Name)
			fmt.Printf("Duration: %s\n", durationStr)

			return nil
		},
	}

	return cmd
}
