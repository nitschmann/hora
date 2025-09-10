package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Long:  `Display the current version of hora.`,
		Run: func(cmd *cobra.Command, args []string) {
			if Version == "" {
				Version = "dev"
			}

			fmt.Printf("%s\n", Version)
		},
	}

	return cmd
}
