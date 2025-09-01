package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewStartCmd() *cobra.Command {
	var project string
	var force bool

	cmd := &cobra.Command{
		Use:     "start [project]",
		Aliases: []string{"s"},
		Short:   "Start tracking time for a project",
		Long:    `Start tracking time for a specific project. If no project name is provided, it will prompt for one. Use --force to stop any existing session and start a new one.`,
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				project = args[0]
			}

			if project == "" {
				fmt.Print("Enter project name: ")
				fmt.Scanln(&project)
				if project == "" {
					return fmt.Errorf("project name cannot be empty")
				}
			}

			var err error
			if force {
				err = timeService.StartTrackingForce(project)
				if err != nil {
					return err
				}
				fmt.Printf("Started tracking time for project: %s (stopped previous session)\n", project)
			} else {
				err = timeService.StartTracking(project)
				if err != nil {
					return err
				}
				fmt.Printf("Started tracking time for project: %s\n", project)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&project, "project", "p", "", "Project name to track time for")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Stop any existing session and start a new one")

	return cmd
}
