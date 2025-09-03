package main

import (
	"context"
	"log"
	"os"

	"github.com/nitschmann/hora/internal/backgroundtracker"
	"github.com/nitschmann/hora/internal/cmd"
)

// Version is the current version of the application
var Version string

func main() {
	ctx := context.Background()

	rootCmd := cmd.NewRootCmd()
	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		// check if background tracker is running and stop it
		if backgroundtracker.IsRunning() {
			err := backgroundtracker.Stop()
			if err != nil {
				log.Printf("Warning: Failed to stop background tracker: %v\n", err)
			}
		}

		log.Fatal(err)
		os.Exit(1)
	}
}
