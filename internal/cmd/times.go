package cmd

import (
	"fmt"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/nitschmann/hora/internal/repository"
)

func NewTimesCmd() *cobra.Command {
	var (
		category string
		since    string
		sort     string
	)

	cmd := &cobra.Command{
		Use:   "times",
		Short: "List all time entries across all projects",
		Long:  `List all time entries across all projects, showing start time, end time, duration, and project information.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			// Get the limit from command flags
			limit, err := cmd.Flags().GetInt("limit")
			if err != nil {
				return fmt.Errorf("failed to get limit flag: %w", err)
			}

			if sort != "asc" && sort != "desc" {
				return fmt.Errorf("sort order must be 'asc' or 'desc', got: %s", sort)
			}

			var sinceTime *time.Time
			if since != "" {
				parsedTime, err := time.Parse("2006-01-02", since)
				if err != nil {
					return fmt.Errorf("invalid date format for --since flag. Use YYYY-MM-DD format: %w", err)
				}
				sinceTime = &parsedTime
			}

			// Validate category if provided
			var categoryPtr *string
			if category != "" {
				if err := validateCategory(category); err != nil {
					return fmt.Errorf("invalid category: %w", err)
				}
				categoryPtr = &category
			}

			// Get all time entries across all projects
			var entries []repository.TimeEntryWithPauses
			if categoryPtr != nil {
				entries, err = timeService.GetAllEntriesWithPausesByCategory(ctx, limit, sort, sinceTime, categoryPtr)
			} else {
				entries, err = timeService.GetAllEntriesWithPauses(ctx, limit, sort, sinceTime)
			}
			if err != nil {
				return fmt.Errorf("failed to get time entries: %w", err)
			}

			if len(entries) == 0 {
				fmt.Println("No time entries found.")
				return nil
			}

			// Create table
			table := tablewriter.NewTable(cmd.OutOrStdout())
			table.Header("Start Time", "End Time", "Project", "Category", "Duration", "Pauses", "Pause Time", "Effective Work Time")

			// Add rows
			for _, entry := range entries {
				startStr := formatTimeInLocal(entry.StartTime)

				var endStr, durationStr string
				if entry.EndTime != nil {
					endStr = formatTimeInLocal(*entry.EndTime)
					if entry.Duration != nil {
						totalDuration := *entry.Duration + entry.PauseTime
						durationStr = timeService.FormatDuration(totalDuration)
					} else {
						durationStr = "Unknown"
					}
				} else {
					endStr = "Active"
					durationStr = "In progress"
				}

				pauseCountStr := fmt.Sprintf("%d", entry.PauseCount)
				pauseTimeStr := timeService.FormatDuration(entry.PauseTime)
				if entry.PauseCount == 0 {
					pauseTimeStr = "-"
				}

				// Calculate and format effective work time (duration - pause time)
				var effectiveWorkTimeStr string
				if entry.EndTime != nil && entry.Duration != nil {
					effectiveWorkTimeStr = timeService.FormatDuration(*entry.Duration)
				} else {
					effectiveWorkTimeStr = "In progress"
				}

				// Format category
				var categoryStr string
				if entry.Category != nil {
					categoryStr = *entry.Category
				} else {
					categoryStr = "-"
				}

				table.Append([]string{
					startStr,
					endStr,
					entry.Project.Name,
					categoryStr,
					durationStr,
					pauseCountStr,
					pauseTimeStr,
					effectiveWorkTimeStr,
				})
			}

			table.Render()

			return nil
		},
	}

	cmd.Flags().StringVar(&category, "category", "", "Filter by category (avoid shell special characters like ! $ ` \\")
	cmd.Flags().IntP("limit", "l", 50, "Maximum number of entries to show")
	cmd.Flags().StringVar(&since, "since", "", "Only show entries since this date (YYYY-MM-DD format)")
	cmd.Flags().StringVarP(&sort, "sort", "s", "desc", "Sort order: 'asc' (oldest first) or 'desc' (newest first)")

	return cmd
}
