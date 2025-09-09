package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/nitschmann/hora/internal/backgroundtracker"
)

func NewStartCmd() *cobra.Command {
	var (
		force                 bool
		skipBackgroundTracker bool
	)

	cmd := &cobra.Command{
		Use:     "start [project]",
		Aliases: []string{"s"},
		Short:   "Start tracking time for a project",
		Long:    `Start tracking time for a specific project. Use -f to stop any existing session and start a new one.`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			project := args[0]
			useBackgroundTracker := conf.UseBackgroundTracker && !skipBackgroundTracker

			if useBackgroundTracker {
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
					_ = dbConn.Close()
					dbConn = nil
					timeService = nil
				}

				backgroundtracker.Daemonize()

				// Parent exits here, daemon continues
				if os.Getenv("IS_DAEMON") != "1" {
					fmt.Printf("Started tracking time for project: %s\n", project)
					return nil
				}

				err := initDatabaseConnectionAndService()
				if err != nil {
					return fmt.Errorf("failed to initialize database and service in daemon: %w", err)
				}
			}

			// --- only daemon process reaches this point ---
			err := timeService.StartTracking(ctx, project, force)
			if err != nil {
				return err
			}

			fmt.Printf("Started tracking time for project: %s\n", project)

			if useBackgroundTracker {
				backgroundtracker.SetTimeTrackingService(timeService)
				backgroundtracker.Start()
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Stop any existing tracking session and start a new one")
	cmd.Flags().BoolVar(&skipBackgroundTracker, "skip-background-tracker", false, "Skip starting the background tracker")

	return cmd
}
