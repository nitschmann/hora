package main

import (
	"context"
	"log"
	"os"

	"github.com/nitschmann/hora/internal/cmd"
)

// Version is the current version of the application
var Version string

func main() {
	ctx := context.Background()

	rootCmd := cmd.NewRootCmd()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
