package cmd

import (
	"fmt"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func NewProjectTimesCmd() *cobra.Command {
	var (
		limit int
		sort  string
		since string
	)

	cmd := &cobra.Command{
		Use:     "times [PROJECT_ID_OR_NAME]",
		Aliases: []string{"t"},
		Short:   "List time entries for a specific project",
		Long:    `List all time entries for a specific project, showing start time, end time, and duration. You can specify either the project ID (numeric) or name.`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			projectIDOrName := args[0]

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

			project, err := timeService.GetProjectByIDOrName(ctx, projectIDOrName)
			if err != nil {
				return fmt.Errorf("failed to get project: %w", err)
			}

			entries, err := timeService.GetEntriesForProjectWithPauses(ctx, projectIDOrName, limit, sort, sinceTime)
			if err != nil {
				return fmt.Errorf("failed to get time entries: %w", err)
			}

			if len(entries) == 0 {
				fmt.Printf("No time entries found for project '%s'.\n", project.Name)
				return nil
			}

			table := tablewriter.NewTable(cmd.OutOrStdout())
			table.Header("Start Time", "End Time", "Category", "Duration", "Pauses", "Pause Time", "Effective Work Time")

			// Add rows
			for _, entry := range entries {
				startStr := formatTimeInLocal(entry.StartTime)

				var endStr string
				var durationStr string

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

				// Calculate and format net working time (duration - pause time)
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

	addListCommandCommonFlags(cmd, &limit, &since, &sort)

	return cmd
}
