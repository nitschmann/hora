package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nitschmann/hora/internal/config"
)

func NewConfigInitCmd() *cobra.Command {
	var (
		directory string
		force     bool
	)

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize configuration",
		Long:  `Initialize the configuration for hora with default values.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			usedConfigFilepath, err := config.CreateDefault(directory, force)
			if err != nil {
				return err
			}

			fmt.Printf("Config file created at %s\n", usedConfigFilepath)

			return nil
		},
	}

	cmd.Flags().StringVarP(&directory, "directory", "d", "", fmt.Sprintf("Directory to create the config file in (default: %q)", config.LookupPaths[0]))
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force re-initialization even if config file exists")

	return cmd
}
