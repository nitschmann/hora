package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/nitschmann/hora/internal/model"
	"github.com/nitschmann/hora/internal/repository"
)

// Database defines the interface for time tracking database operations
type Database interface {
	// Project management
	GetOrCreateProject(ctx context.Context, name string) (*model.Project, error)
	GetProject(ctx context.Context, id int) (*model.Project, error)
	GetProjectByName(ctx context.Context, name string) (*model.Project, error)
	GetProjectByIDOrName(ctx context.Context, idOrName string) (*model.Project, error)
	GetAllProjects(ctx context.Context) ([]model.Project, error)
	RemoveProject(ctx context.Context, name string) error
	RemoveProjectByIDOrName(ctx context.Context, idOrName string) error

	// Time tracking
	StartTracking(ctx context.Context, project string, force bool) error
	StopTracking(ctx context.Context) (*model.TimeEntry, error)
	GetActiveEntry(ctx context.Context) (*model.TimeEntry, error)
	GetEntries(ctx context.Context, limit int) ([]model.TimeEntry, error)
	GetEntriesForProject(ctx context.Context, projectIDOrName string, limit int, sortOrder string) ([]model.TimeEntry, error)
	GetEntriesForProjectWithPauses(ctx context.Context, projectIDOrName string, limit int, sortOrder string, since *time.Time) ([]repository.TimeEntryWithPauses, error)
	GetAllEntriesWithPauses(ctx context.Context, limit int, sortOrder string, since *time.Time) ([]repository.TimeEntryWithPauses, error)
	GetTotalTimeForProject(ctx context.Context, projectIDOrName string, since *time.Time) (time.Duration, error)

	// Data management
	ClearAllEntries(ctx context.Context) error

	// Pause management
	PauseTracking(ctx context.Context) error
	ContinueTracking(ctx context.Context) error

	// Connection management
	Close() error

	// Get underlying database connection for repositories
	GetDB() *sql.DB

	// Get repositories
	GetProjectRepository() repository.Project
	GetTimeEntryRepository() repository.TimeEntry
	GetPauseRepository() repository.Pause
}
