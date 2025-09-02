package service

import (
	"context"
	"fmt"
	"time"

	"github.com/nitschmann/hora/internal/database"
	"github.com/nitschmann/hora/internal/model"
	"github.com/nitschmann/hora/internal/repository"
)

// TimeTracking defines the interface for time tracking operations
type TimeTracking interface {
	StartTracking(ctx context.Context, projectName string, force bool) error
	StopTracking(ctx context.Context) (*model.TimeEntry, error)
	GetActiveEntry(ctx context.Context) (*model.TimeEntry, error)
	GetEntries(ctx context.Context, limit int) ([]model.TimeEntry, error)
	GetEntriesForProject(ctx context.Context, projectIDOrName string, limit int, sortOrder string) ([]model.TimeEntry, error)
	GetEntriesForProjectWithPauses(ctx context.Context, projectIDOrName string, limit int, sortOrder string, since *time.Time) ([]repository.TimeEntryWithPauses, error)
	ClearAllData(ctx context.Context) error
	GetProjects(ctx context.Context) ([]model.Project, error)
	GetOrCreateProject(ctx context.Context, name string) (*model.Project, error)
	GetProjectByIDOrName(ctx context.Context, idOrName string) (*model.Project, error)
	RemoveProject(ctx context.Context, name string) error
	RemoveProjectByIDOrName(ctx context.Context, idOrName string) error
	PauseTracking(ctx context.Context) error
	ContinueTracking(ctx context.Context) error
	FormatDuration(duration time.Duration) string
}

// timeTracking implements the TimeTracking interface
type timeTracking struct {
	db database.Database
}

// NewTimeTracking creates a new time tracking service
func NewTimeTracking(db database.Database) TimeTracking {
	return &timeTracking{db: db}
}

// StartTracking starts a new time tracking session for the given project
func (s *timeTracking) StartTracking(ctx context.Context, projectName string, force bool) error {
	return s.db.StartTracking(ctx, projectName, force)
}

// StopTracking stops the current active time tracking session
func (s *timeTracking) StopTracking(ctx context.Context) (*model.TimeEntry, error) {
	return s.db.StopTracking(ctx)
}

// GetActiveEntry returns the currently active time tracking entry, if any
func (s *timeTracking) GetActiveEntry(ctx context.Context) (*model.TimeEntry, error) {
	return s.db.GetActiveEntry(ctx)
}

// GetEntries returns a list of time entries, limited by the given count
func (s *timeTracking) GetEntries(ctx context.Context, limit int) ([]model.TimeEntry, error) {
	return s.db.GetEntries(ctx, limit)
}

// GetEntriesForProject returns time entries for a project by ID (if numeric) or name
func (s *timeTracking) GetEntriesForProject(ctx context.Context, projectIDOrName string, limit int, sortOrder string) ([]model.TimeEntry, error) {
	return s.db.GetEntriesForProject(ctx, projectIDOrName, limit, sortOrder)
}

// GetEntriesForProjectWithPauses returns time entries with pause information for a project by ID (if numeric) or name
func (s *timeTracking) GetEntriesForProjectWithPauses(ctx context.Context, projectIDOrName string, limit int, sortOrder string, since *time.Time) ([]repository.TimeEntryWithPauses, error) {
	return s.db.GetEntriesForProjectWithPauses(ctx, projectIDOrName, limit, sortOrder, since)
}

// ClearAllData removes all time entries and projects from the database
func (s *timeTracking) ClearAllData(ctx context.Context) error {
	return s.db.ClearAllEntries(ctx)
}

// GetProjects returns all projects in the system
func (s *timeTracking) GetProjects(ctx context.Context) ([]model.Project, error) {
	return s.db.GetAllProjects(ctx)
}

// GetOrCreateProject gets an existing project or creates a new one
func (s *timeTracking) GetOrCreateProject(ctx context.Context, name string) (*model.Project, error) {
	return s.db.GetOrCreateProject(ctx, name)
}

// GetProjectByIDOrName gets a project by ID (if numeric) or name
func (s *timeTracking) GetProjectByIDOrName(ctx context.Context, idOrName string) (*model.Project, error) {
	return s.db.GetProjectByIDOrName(ctx, idOrName)
}

// RemoveProject removes a project and all its time entries
func (s *timeTracking) RemoveProject(ctx context.Context, name string) error {
	return s.db.RemoveProject(ctx, name)
}

// RemoveProjectByIDOrName removes a project by ID (if numeric) or name and all its time entries
func (s *timeTracking) RemoveProjectByIDOrName(ctx context.Context, idOrName string) error {
	return s.db.RemoveProjectByIDOrName(ctx, idOrName)
}

// PauseTracking pauses the currently active time tracking session
func (s *timeTracking) PauseTracking(ctx context.Context) error {
	return s.db.PauseTracking(ctx)
}

// ContinueTracking continues the currently paused time tracking session
func (s *timeTracking) ContinueTracking(ctx context.Context) error {
	return s.db.ContinueTracking(ctx)
}

// FormatDuration formats a duration into HH:MM:SS format
func (s *timeTracking) FormatDuration(duration time.Duration) string {
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}
