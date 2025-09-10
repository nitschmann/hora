package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	"github.com/nitschmann/hora/internal/backgroundtracker"
)

func NewLogsCmd() *cobra.Command {
	var follow bool

	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Display background (daemon) tracker logs",
		Long:  `Display logs from the background (daemon) tracker. Use --follow to tail the logs in real-time.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if runtime.GOOS != "darwin" {
				return fmt.Errorf("logs command is only supported on macOS")
			}

			logPath, err := backgroundtracker.GetLogPath()
			if err != nil {
				return fmt.Errorf("failed to get log path: %w", err)
			}

			if _, err := os.Stat(logPath); os.IsNotExist(err) {
				return fmt.Errorf("log file does not exist: %s", logPath)
			}

			fmt.Printf("logs from: %s\n\n", logPath)

			var tailCmd *exec.Cmd
			if follow {
				fmt.Println("Press Ctrl+C to stop")
				fmt.Println(strings.Repeat("-", 50))
				tailCmd = exec.Command("tail", "-f", logPath)
			} else {
				tailCmd = exec.Command("tail", logPath)
			}

			tailCmd.Stdout = os.Stdout
			tailCmd.Stderr = os.Stderr

			return tailCmd.Run()
		},
	}

	cmd.Flags().BoolVarP(&follow, "follow", "f", false, "Follow log output (like tail -f)")

	return cmd
}
