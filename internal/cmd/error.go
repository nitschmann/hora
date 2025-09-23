package cmd

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
)

// mapCmdError maps internal errors to user-friendly cmd errors
func mapCmdError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("no matching record(s) found")
	}

	return err
}

// mapCmdErrorAndExit maps an error and exits the program with appropriate exit code
func mapCmdErrorAndExit(err error) {
	cmdErr := mapCmdError(err)
	if cmdErr == nil {
		os.Exit(0)
	}

	fmt.Printf("An unexpected error occurred:\n%v\n", cmdErr.Error())

	os.Exit(1)
}
