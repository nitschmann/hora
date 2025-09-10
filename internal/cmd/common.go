package cmd

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/nitschmann/hora/internal/database"
	"github.com/nitschmann/hora/internal/repository"
	"github.com/nitschmann/hora/internal/service"
)

var (
	dbConn      *database.Connection
	timeService service.TimeTracking
)

// formatTimeInLocal formats a time value in the local timezone
func formatTimeInLocal(t time.Time) string {
	return t.Local().Format("2006-01-02 15:04:05")
}

// formatTimeInLocalShort formats a time value in the local timezone with shorter format
func formatTimeInLocalShort(t time.Time) string {
	return t.Local().Format("2006-01-02 15:04")
}

// formatDateInLocal formats a date in the local timezone
func formatDateInLocal(t time.Time) string {
	return t.Local().Format("2006-01-02")
}

// validateCategory validates that a category contains only alphanumeric characters, underscores, and hyphens
func validateCategory(category string) error {
	if category == "" {
		return nil
	}

	// Check for common shell special characters that might cause issues
	if strings.ContainsAny(category, "!$`\\") {
		return fmt.Errorf("category contains shell special characters (!$`\\) that may cause issues. Use only alphanumeric characters, underscores (_), and hyphens (-)")
	}

	matched, err := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, category)
	if err != nil {
		return fmt.Errorf("failed to validate category: %w", err)
	}

	if !matched {
		return fmt.Errorf("category must contain only alphanumeric characters, underscores (_), and hyphens (-)")
	}

	return nil
}

func initDatabaseConnectionAndService() error {
	var err error
	dbConn, err = database.NewConnection(conf)
	if err != nil {
		return fmt.Errorf("failed to initialize database in daemon: %w", err)
	}

	// repositories for service
	projectRepo := repository.NewProject(dbConn.GetDB())
	timeEntryRepo := repository.NewTimeEntry(dbConn.GetDB())
	pauseRepo := repository.NewPause(dbConn.GetDB())

	timeService = service.NewTimeTracking(projectRepo, timeEntryRepo, pauseRepo)

	return err
}
