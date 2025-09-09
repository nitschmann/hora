package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/nitschmann/hora/internal/backgroundtracker"
	"github.com/nitschmann/hora/internal/database"
	"github.com/nitschmann/hora/internal/repository"
	"github.com/nitschmann/hora/internal/service"
)

func NewStartCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:     "start [project]",
		Aliases: []string{"s"},
		Short:   "Start tracking time for a project",
		Long:    `Start tracking time for a specific project. Use -f to stop any existing session and start a new one.`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			project := args[0]

			if conf.UseBackgroundTracker {
				if backgroundtracker.IsRunning() {
					activeEntry, err := timeService.GetActiveEntry(ctx)
					if err == nil && activeEntry != nil {
						if !force {
							return fmt.Errorf(
								"a time tracking session is already active for project '%s'. Use --force to stop it and start a new one",
								activeEntry.Project.Name,
							)
						}
						if err := backgroundtracker.Stop(); err != nil {
							return fmt.Errorf("failed to stop existing daemon: %w", err)
						}
					}
				}

				// Close parent database connection before forking
				if dbConn != nil {
					dbConn.Close()
					dbConn = nil
				}

				backgroundtracker.Daemonize()

				// Parent exits here, daemon continues
				if os.Getenv("IS_DAEMON") != "1" {
					fmt.Printf("Started tracking time for project: %s\n", project)
					return nil
				}

				// Daemon process: reinitialize database connection
				var err error
				dbConn, err = database.NewConnection(conf)
				if err != nil {
					return fmt.Errorf("failed to initialize database in daemon: %w", err)
				}

				// Create repositories for daemon
				projectRepo := repository.NewProject(dbConn.GetDB())
				timeEntryRepo := repository.NewTimeEntry(dbConn.GetDB())
				pauseRepo := repository.NewPause(dbConn.GetDB())

				timeService = service.NewTimeTracking(projectRepo, timeEntryRepo, pauseRepo)
			}

			// --- only daemon process reaches this point ---
			err := timeService.StartTracking(ctx, project, force)
			if err != nil {
				return err
			}

			fmt.Printf("Started tracking time for project: %s\n", project)

			if conf.UseBackgroundTracker {
				backgroundtracker.SetTimeTrackingService(timeService)
				backgroundtracker.Start()
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Stop any existing tracking session and start a new one")

	return cmd
}
