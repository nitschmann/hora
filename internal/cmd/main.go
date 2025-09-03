package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nitschmann/hora/internal/database"
	"github.com/nitschmann/hora/internal/service"
)

// package variables for database and service
var (
	db          database.Database
	timeService service.TimeTracking
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "hora",
		Short:         "Hora is a simple time tracking CLI tool",
		Long:          `Hora is a simple command-line time tracking tool. Track your project time with ease.`,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error

			db, err = database.NewDatabase()
			if err != nil {
				return fmt.Errorf("failed to initialize database: %w", err)
			}

			timeService = service.NewTimeTracking(db)

			return nil
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			// if db != nil {
			// 	db.Close()
			// }
		},
	}

	rootCmd.AddCommand(NewStartCmd())
	rootCmd.AddCommand(NewStopCmd())
	rootCmd.AddCommand(NewPauseCmd())
	rootCmd.AddCommand(NewContinueCmd())
	rootCmd.AddCommand(NewStatusCmd())
	rootCmd.AddCommand(NewTimesCmd())
	rootCmd.AddCommand(NewDeleteAllCmd())
	rootCmd.AddCommand(NewProjectCmd())
	rootCmd.AddCommand(NewLogsCmd())

	return rootCmd
}
