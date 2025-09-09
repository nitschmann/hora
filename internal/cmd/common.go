package cmd

import (
	"fmt"

	"github.com/nitschmann/hora/internal/database"
	"github.com/nitschmann/hora/internal/repository"
	"github.com/nitschmann/hora/internal/service"
)

var (
	dbConn      *database.Connection
	timeService service.TimeTracking
)

func initDatabaseConnectionAndService() error {
	var err error
	dbConn, err = database.NewConnection(conf)
	if err != nil {
		return fmt.Errorf("failed to initialize database in daemon: %w", err)
	}

	// Create repositories for daemon
	projectRepo := repository.NewProject(dbConn.GetDB())
	timeEntryRepo := repository.NewTimeEntry(dbConn.GetDB())
	pauseRepo := repository.NewPause(dbConn.GetDB())

	timeService = service.NewTimeTracking(projectRepo, timeEntryRepo, pauseRepo)

	return err
}
