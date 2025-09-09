package cmd

import (
	"fmt"
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
