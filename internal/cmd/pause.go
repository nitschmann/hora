package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func NewPauseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pause",
		Short: "Pause the currently active time tracking session",
		Long:  `Pause the currently active time tracking session. You can resume it later with the 'continue' command.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			err := timeService.PauseTracking(ctx)
			if err != nil {
				return fmt.Errorf("failed to pause tracking: %w", err)
			}

			fmt.Println("Time tracking paused.")
			return nil
		},
	}

	return cmd
}
