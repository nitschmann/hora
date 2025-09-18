package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/nitschmann/hora/internal/repository"
)

func NewExportCmd() *cobra.Command {
	var (
		category string
		limit    int
		since    string
		sort     string
		output   string
	)

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export time entries to CSV",
		Long:  "Export all time entries across projects to a CSV file",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err       error
				sinceTime *time.Time
			)

			ctx := context.Background()

			if since != "" {
				parsed, err := time.Parse("2006-01-02", since)
				if err != nil {
					return fmt.Errorf("invalid date format: %w", err)
				}
				sinceTime = &parsed
			}

			var categoryPtr *string
			if category != "" {
				if err := validateCategory(category); err != nil {
					return err
				}
				categoryPtr = &category
			}

			var entries []repository.TimeEntryWithPauses
			if categoryPtr != nil {
				entries, err = timeService.GetAllEntriesWithPausesByCategory(ctx, limit, sort, sinceTime, categoryPtr)
			} else {
				entries, err = timeService.GetAllEntriesWithPauses(ctx, limit, sort, sinceTime)
			}
			if err != nil {
				return fmt.Errorf("failed to get time entries: %w", err)
			}

			filename, err := exportTimesToCSV(entries, output, "")
			if err != nil {
				return fmt.Errorf("failed to export CSV: %w", err)
			}

			fmt.Printf("Exported %d time entries to %s\n", len(entries), filename)

			return nil
		},
	}

	addListCommandCommonFlags(cmd, &limit, &since, &sort)

	cmd.Flags().StringVar(&category, "category", "", "Filter by category")
	cmd.Flags().StringVarP(&output, "output", "o", "", "Output file path (default: TIMESTAMP_times.csv)")

	return cmd
}
