package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/nitschmann/hora/internal/config"
)

var conf *config.Config

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:               "hora",
		Short:             "hora is a simple time tracking CLI tool",
		Long:              `hora is a simple command-line time tracking tool. Track your project time with ease.`,
		SilenceUsage:      true,
		SilenceErrors:     true,
		DisableAutoGenTag: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error

			configFile, err := cmd.Flags().GetString("config")
			if err != nil {
				return fmt.Errorf("failed to get config flag: %w", err)
			}

			conf, _, err = config.Load(configFile)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			err = initDatabaseConnectionAndService()
			if err != nil {
				return fmt.Errorf("failed to initialize database and service in daemon: %w", err)
			}

			return nil
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			if os.Getenv("IS_DAEMON") == "1" {
				return nil
			}

			if dbConn != nil {
				return dbConn.Close()
			}

			return nil
		},
	}

	rootCmd.PersistentFlags().StringP("config", "c", "", "Path to configuration file")

	rootCmd.AddCommand(NewDeleteAllCmd())
	rootCmd.AddCommand(NewContinueCmd())
	rootCmd.AddCommand(NewStartCmd())
	rootCmd.AddCommand(NewStopCmd())
	rootCmd.AddCommand(NewPauseCmd())
	rootCmd.AddCommand(NewStatusCmd())
	rootCmd.AddCommand(NewTimesCmd())
	rootCmd.AddCommand(NewExportCmd())
	rootCmd.AddCommand(NewProjectCmd())

	return rootCmd
}
