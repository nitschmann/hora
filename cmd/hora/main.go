package main

import (
	"context"
	"fmt"
	"os"

	"github.com/nitschmann/hora/internal/cmd"
)

// Version is the current version of the application
var Version string

func main() {
	ctx := context.Background()
	rootCmd := cmd.NewRootCmd()
	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		fmt.Printf("Error: %+v\n", err)
		os.Exit(1)
	}
}
