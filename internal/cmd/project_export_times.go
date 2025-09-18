package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

func NewProjectExportTimesCmd() *cobra.Command {
	var (
		since  string
		sort   string
		output string
		limit  int
	)

	cmd := &cobra.Command{
		Use:   "export-times [project]",
		Short: "Export project time entries to CSV",
		Long:  "Export time entries for a specific project to a CSV file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			projectName := args[0]

			var sinceTime *time.Time
			if since != "" {
				parsed, err := time.Parse("2006-01-02", since)
				if err != nil {
					return fmt.Errorf("invalid date format: %w", err)
				}
				sinceTime = &parsed
			}

			entries, err := timeService.GetEntriesForProjectWithPauses(ctx, projectName, limit, sort, sinceTime)
			if err != nil {
				return fmt.Errorf("failed to get project time entries: %w", err)
			}

			filename, err := exportTimesToCSV(entries, output, projectName)
			if err != nil {
				return fmt.Errorf("failed to export CSV: %w", err)
			}

			fmt.Printf("Exported %d time entries for project %q to %s\n", len(entries), projectName, filename)
			return nil
		},
	}

	addListCommandCommonFlags(cmd, &limit, &since, &sort)

	cmd.Flags().StringVarP(&output, "output", "o", "", "Output file path (default: TIMESTAMP_PROJECT_times.csv)")

	return cmd
}
