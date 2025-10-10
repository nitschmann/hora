package cmd

import (
	"github.com/nitschmann/hora/internal/ui"
	"github.com/spf13/cobra"
)

func NewUICommand() *cobra.Command {
	var port int

	cmd := &cobra.Command{
		Use:   "ui",
		Short: "Start the web UI",
		Long:  `Start the web UI for time tracking with interactive charts and analytics.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			server := ui.NewServer(timeService)
			return server.Start(ctx, port)
		},
	}

	cmd.Flags().IntVarP(&port, "port", "p", conf.WebUIPort, "Port to run the web UI on")

	return cmd
}
