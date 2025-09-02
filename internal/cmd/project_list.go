package cmd

import (
	"fmt"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func NewProjectListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all projects",
		Long:    `List all projects in your time tracking system.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			projects, err := timeService.GetProjects(ctx)
			if err != nil {
				return fmt.Errorf("failed to get projects: %w", err)
			}

			if len(projects) == 0 {
				fmt.Println("No projects found.")
				return nil
			}

			table := tablewriter.NewTable(cmd.OutOrStdout())
			table.Header("ID", "Name", "Created", "Last Tracked")

			for _, project := range projects {
				createdStr := project.CreatedAt.Format("2006-01-02 15:04")

				var lastTrackedStr string
				if project.LastTrackedAt != nil {
					lastTrackedStr = project.LastTrackedAt.Format("2006-01-02 15:04")
				} else {
					lastTrackedStr = "Never"
				}

				table.Append([]string{
					fmt.Sprintf("%d", project.ID),
					project.Name,
					createdStr,
					lastTrackedStr,
				})
			}

			table.Render()

			return nil
		},
	}

	return cmd
}
