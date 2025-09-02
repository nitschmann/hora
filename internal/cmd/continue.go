package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewContinueCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "continue",
		Short: "Continue the currently paused time tracking session",
		Long:  `Continue the currently paused time tracking session. This will end the current pause and resume tracking.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			err := timeService.ContinueTracking(ctx)
			if err != nil {
				return fmt.Errorf("failed to continue tracking: %w", err)
			}

			fmt.Println("Time tracking resumed.")
			return nil
		},
	}

	return cmd
}
