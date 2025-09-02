package cmd

import (
	"context"
	"fmt"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func NewProjectTimesCmd() *cobra.Command {
	var limit int
	var sort string

	cmd := &cobra.Command{
		Use:     "times [ID_OR_NAME]",
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

			ctx := context.Background()

			// First check if the project exists
			project, err := timeService.GetProjectByIDOrName(ctx, projectIDOrName)
			if err != nil {
				return fmt.Errorf("failed to get project: %w", err)
			}

			// Get time entries for the project
			entries, err := timeService.GetEntriesForProjectByIDOrName(ctx, projectIDOrName, limit, sort)
			if err != nil {
				return fmt.Errorf("failed to get time entries: %w", err)
			}

			if len(entries) == 0 {
				fmt.Printf("No time entries found for project '%s'.\n", project.Name)
				return nil
			}

			// Create table
			table := tablewriter.NewTable(cmd.OutOrStdout())
			table.Header("Start Time", "End Time", "Duration")

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

				table.Append([]string{
					startStr,
					endStr,
					durationStr,
				})
			}

			// Render table
			table.Render()

			return nil
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 50, "Maximum number of entries to show")
	cmd.Flags().StringVarP(&sort, "sort", "s", "desc", "Sort order: 'asc' (oldest first) or 'desc' (newest first)")

	return cmd
}
