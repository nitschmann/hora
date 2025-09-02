package cmd

import (
	"fmt"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func NewProjectTimesCmd() *cobra.Command {
	var limit int
	var sort string
	var since string

	cmd := &cobra.Command{
		Use:     "times [PROJECT_ID_OR_NAME]",
		Aliases: []string{"t"},
		Short:   "List time entries for a specific project",
		Long:    `List all time entries for a specific project, showing start time, end time, and duration. You can specify either the project ID (numeric) or name.`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectIDOrName := args[0]

			// Validate sort order
			if sort != "asc" && sort != "desc" {
				return fmt.Errorf("sort order must be 'asc' or 'desc', got: %s", sort)
			}

			// Parse since date if provided
			var sinceTime *time.Time
			if since != "" {
				parsedTime, err := time.Parse("2006-01-02", since)
				if err != nil {
					return fmt.Errorf("invalid date format for --since flag. Use YYYY-MM-DD format: %w", err)
				}
				sinceTime = &parsedTime
			}

			ctx := cmd.Context()

			// First check if the project exists
			project, err := timeService.GetProjectByIDOrName(ctx, projectIDOrName)
			if err != nil {
				return fmt.Errorf("failed to get project: %w", err)
			}

			// Get time entries with pause information for the project
			entries, err := timeService.GetEntriesForProjectWithPauses(ctx, projectIDOrName, limit, sort, sinceTime)
			if err != nil {
				return fmt.Errorf("failed to get time entries: %w", err)
			}

			if len(entries) == 0 {
				fmt.Printf("No time entries found for project '%s'.\n", project.Name)
				return nil
			}

			// Create table
			table := tablewriter.NewTable(cmd.OutOrStdout())
			table.Header("Start Time", "End Time", "Duration", "Pauses", "Pause Time")

			// Add rows
			for _, entry := range entries {
				startStr := entry.StartTime.Format("2006-01-02 15:04:05")

				var endStr string
				var durationStr string

				if entry.EndTime != nil {
					endStr = entry.EndTime.Format("2006-01-02 15:04:05")
					if entry.Duration != nil {
						durationStr = timeService.FormatDuration(*entry.Duration)
					} else {
						durationStr = "Unknown"
					}
				} else {
					endStr = "Active"
					durationStr = "In progress"
				}

				// Format pause information
				pauseCountStr := fmt.Sprintf("%d", entry.PauseCount)
				pauseTimeStr := timeService.FormatDuration(entry.PauseTime)
				if entry.PauseCount == 0 {
					pauseTimeStr = "-"
				}

				table.Append([]string{
					startStr,
					endStr,
					durationStr,
					pauseCountStr,
					pauseTimeStr,
				})
			}

			// Render table
			table.Render()

			return nil
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 50, "Maximum number of entries to show")
	cmd.Flags().StringVarP(&sort, "sort", "s", "desc", "Sort order: 'asc' (oldest first) or 'desc' (newest first)")
	cmd.Flags().StringVar(&since, "since", "", "Only show entries since this date (YYYY-MM-DD format)")

	return cmd
}
