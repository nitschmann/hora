package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/nitschmann/hora/internal/service"
)

func NewLogsCmd() *cobra.Command {
	var lines int

	cmd := &cobra.Command{
		Use:   "logs",
		Short: "View time tracking logs",
		Long:  `View the time tracking logs to see screen lock detection events and other activities.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get the latest log file
			logFile, err := service.GetLatestLogFile()
			if err != nil {
				return fmt.Errorf("failed to get log file: %w", err)
			}

			// Read and display the log file
			content, err := os.ReadFile(logFile)
			if err != nil {
				return fmt.Errorf("failed to read log file: %w", err)
			}

			logLines := strings.Split(string(content), "\n")

			// Show last N lines if specified
			if lines > 0 && len(logLines) > lines {
				logLines = logLines[len(logLines)-lines:]
			}

			// Display the logs
			for _, line := range logLines {
				if line != "" {
					fmt.Println(line)
				}
			}

			return nil
		},
	}

	cmd.Flags().IntVarP(&lines, "lines", "n", 50, "Number of lines to show from the end of the log")

	return cmd
}
