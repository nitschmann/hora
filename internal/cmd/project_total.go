package cmd

import (
	"fmt"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func NewProjectTotalCmd() *cobra.Command {
	var since string

	cmd := &cobra.Command{
		Use:   "total [PROJECT_ID_OR_NAME]",
		Short: "Show total tracked time for a project",
		Long:  `Show the total tracked time for a specific project, including all time entries and accounting for pauses. You can specify either the project ID (numeric) or name.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			projectIDOrName := args[0]

			// Parse since date if provided
			var sinceTime *time.Time
			if since != "" {
				parsedTime, err := time.Parse("2006-01-02", since)
				if err != nil {
					return fmt.Errorf("invalid date format for --since flag. Use YYYY-MM-DD format: %w", err)
				}
				sinceTime = &parsedTime
			}

			// First check if the project exists
			project, err := timeService.GetProjectByIDOrName(ctx, projectIDOrName)
			if err != nil {
				return fmt.Errorf("failed to get project: %w", err)
			}

			// Get total time for the project
			totalTime, err := timeService.GetTotalTimeForProject(ctx, projectIDOrName, sinceTime)
			if err != nil {
				return fmt.Errorf("failed to get total time: %w", err)
			}

			table := tablewriter.NewTable(cmd.OutOrStdout())
			table.Header("Project", "Total Time", "Since")
			formatedSince := ""
			if sinceTime != nil {
				formatedSince = sinceTime.Format("2006-01-02")
			}

			table.Append([]string{
				project.Name,
				timeService.FormatDuration(totalTime),
				formatedSince,
			})

			table.Render()

			return nil
		},
	}

	cmd.Flags().StringVar(&since, "since", "", "Only include time since this date (YYYY-MM-DD format)")

	return cmd
}
