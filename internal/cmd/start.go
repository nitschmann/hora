package cmd

import (
	"fmt"
	"os"

	"github.com/nitschmann/hora/internal/backgroundtracker"
	"github.com/spf13/cobra"
)

// TimeTrackingScreenLockHandler implements the ScreenLockHandler interface
type TimeTrackingScreenLockHandler struct{}

func (h *TimeTrackingScreenLockHandler) OnScreenLocked() {
	fmt.Println("🔒 Screen locked - pausing time tracking")
	// You can add logic here to pause time tracking
	// For example: timeService.PauseTracking(context.Background())
}

func (h *TimeTrackingScreenLockHandler) OnScreenUnlocked() {
	fmt.Println("🔓 Screen unlocked - resuming time tracking")
	// You can add logic here to resume time tracking
	// For example: timeService.ContinueTracking(context.Background())
}

func NewStartCmd() *cobra.Command {
	var project string
	var force bool

	cmd := &cobra.Command{
		Use:     "start [project]",
		Aliases: []string{"s"},
		Short:   "Start tracking time for a project",
		Long:    `Start tracking time for a specific project. If no project name is provided, it will prompt for one. Use --force to stop any existing session and start a new one.`,
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if len(args) > 0 {
				project = args[0]
			}

			if project == "" {
				fmt.Print("Enter project name: ")
				fmt.Scanln(&project)
				if project == "" {
					return fmt.Errorf("project name cannot be empty")
				}
			}
			// fork into daemon
			backgroundtracker.Daemonize()

			// Parent exits here, daemon continues
			if os.Getenv("IS_DAEMON") != "1" {
				fmt.Printf("Started tracking time for project: %s\n", project)
				return nil
			}

			err := timeService.StartTracking(ctx, project, force)
			if err != nil {
				return err
			}

			fmt.Printf("Started tracking time for project: %s\n", project)

			// Set up the time tracking service for screen lock integration
			backgroundtracker.SetTimeTrackingService(timeService)
			backgroundtracker.Start()

			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Stop any existing tracking session and start a new one")

	return cmd
}
