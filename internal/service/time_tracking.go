package service

import (
	"context"
	"fmt"
	"time"

	"github.com/nitschmann/hora/internal/database"
	"github.com/nitschmann/hora/internal/model"
)

// TimeTracking defines the interface for time tracking operations
type TimeTracking interface {
	StartTracking(ctx context.Context, projectName string) error
	StartTrackingForce(ctx context.Context, projectName string) error
	StopTracking(ctx context.Context) (*model.TimeEntry, error)
	GetActiveEntry(ctx context.Context) (*model.TimeEntry, error)
	GetEntries(ctx context.Context, limit int) ([]model.TimeEntry, error)
	GetEntriesForProject(ctx context.Context, projectName string, limit int, sortOrder string) ([]model.TimeEntry, error)
	GetEntriesForProjectByIDOrName(ctx context.Context, projectIDOrName string, limit int, sortOrder string) ([]model.TimeEntry, error)
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

// timeTrackingService implements the TimeTracking interface
type timeTrackingService struct {
	db database.Database
}

// NewTimeTracking creates a new time tracking service
func NewTimeTracking(db database.Database) TimeTracking {
	return &timeTrackingService{db: db}
}

// StartTracking starts a new time tracking session for the given project
func (s *timeTrackingService) StartTracking(ctx context.Context, projectName string) error {
	return s.db.StartTracking(ctx, projectName)
}

// StartTrackingForce starts a new time tracking session, stopping any existing session first
func (s *timeTrackingService) StartTrackingForce(ctx context.Context, projectName string) error {
	return s.db.StartTrackingForce(ctx, projectName)
}

// StopTracking stops the current active time tracking session
func (s *timeTrackingService) StopTracking(ctx context.Context) (*model.TimeEntry, error) {
	return s.db.StopTracking(ctx)
}

// GetActiveEntry returns the currently active time tracking entry, if any
func (s *timeTrackingService) GetActiveEntry(ctx context.Context) (*model.TimeEntry, error) {
	return s.db.GetActiveEntry(ctx)
}

// GetEntries returns a list of time entries, limited by the given count
func (s *timeTrackingService) GetEntries(ctx context.Context, limit int) ([]model.TimeEntry, error) {
	return s.db.GetEntries(ctx, limit)
}

// GetEntriesForProject returns time entries for a specific project
func (s *timeTrackingService) GetEntriesForProject(ctx context.Context, projectName string, limit int, sortOrder string) ([]model.TimeEntry, error) {
	return s.db.GetEntriesForProject(ctx, projectName, limit, sortOrder)
}

// GetEntriesForProjectByIDOrName returns time entries for a project by ID (if numeric) or name
func (s *timeTrackingService) GetEntriesForProjectByIDOrName(ctx context.Context, projectIDOrName string, limit int, sortOrder string) ([]model.TimeEntry, error) {
	return s.db.GetEntriesForProjectByIDOrName(ctx, projectIDOrName, limit, sortOrder)
}

// ClearAllData removes all time entries and projects from the database
func (s *timeTrackingService) ClearAllData(ctx context.Context) error {
	return s.db.ClearAllEntries(ctx)
}

// GetProjects returns all projects in the system
func (s *timeTrackingService) GetProjects(ctx context.Context) ([]model.Project, error) {
	return s.db.GetAllProjects(ctx)
}

// GetOrCreateProject gets an existing project or creates a new one
func (s *timeTrackingService) GetOrCreateProject(ctx context.Context, name string) (*model.Project, error) {
	return s.db.GetOrCreateProject(ctx, name)
}

// GetProjectByIDOrName gets a project by ID (if numeric) or name
func (s *timeTrackingService) GetProjectByIDOrName(ctx context.Context, idOrName string) (*model.Project, error) {
	return s.db.GetProjectByIDOrName(ctx, idOrName)
}

// RemoveProject removes a project and all its time entries
func (s *timeTrackingService) RemoveProject(ctx context.Context, name string) error {
	return s.db.RemoveProject(ctx, name)
}

// RemoveProjectByIDOrName removes a project by ID (if numeric) or name and all its time entries
func (s *timeTrackingService) RemoveProjectByIDOrName(ctx context.Context, idOrName string) error {
	return s.db.RemoveProjectByIDOrName(ctx, idOrName)
}

// PauseTracking pauses the currently active time tracking session
func (s *timeTrackingService) PauseTracking(ctx context.Context) error {
	return s.db.PauseTracking(ctx)
}

// ContinueTracking continues the currently paused time tracking session
func (s *timeTrackingService) ContinueTracking(ctx context.Context) error {
	return s.db.ContinueTracking(ctx)
}

// FormatDuration formats a duration into HH:MM:SS format
func (s *timeTrackingService) FormatDuration(duration time.Duration) string {
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}
