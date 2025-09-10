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
		since    string
		sort     string
		output   string
	)

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export time entries to CSV",
		Long:  "Export all time entries across projects to a CSV file",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			limit, err := cmd.Flags().GetInt("limit")
			if err != nil {
				return fmt.Errorf("failed to get limit flag: %w", err)
			}

			var sinceTime *time.Time
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

	cmd.Flags().StringVar(&category, "category", "", "Filter by category")
	cmd.Flags().IntP("limit", "l", 50, "Maximum number of entries to show")
	cmd.Flags().StringVar(&since, "since", "", "Only show entries since this date (YYYY-MM-DD format)")
	cmd.Flags().StringVarP(&sort, "sort", "s", "desc", "Sort order: 'asc' (oldest first) or 'desc' (newest first)")
	cmd.Flags().StringVarP(&output, "output", "o", "", "Output file path (default: TIMESTAMP_times.csv)")

	return cmd
}
