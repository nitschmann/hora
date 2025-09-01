package service

import (
	"fmt"
	"time"

	"github.com/nitschmann/hora/internal/database"
	"github.com/nitschmann/hora/internal/model"
)

// TimeTracking defines the interface for time tracking operations
type TimeTracking interface {
	StartTracking(projectName string) error
	StartTrackingForce(projectName string) error
	StopTracking() (*model.TimeEntry, error)
	GetActiveEntry() (*model.TimeEntry, error)
	GetEntries(limit int) ([]model.TimeEntry, error)
	GetEntriesForProject(projectName string, limit int, sortOrder string) ([]model.TimeEntry, error)
	ClearAllData() error
	GetProjects() ([]model.Project, error)
	GetOrCreateProject(name string) (*model.Project, error)
	RemoveProject(name string) error
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
func (s *timeTrackingService) StartTracking(projectName string) error {
	return s.db.StartTracking(projectName)
}

// StartTrackingForce starts a new time tracking session, stopping any existing session first
func (s *timeTrackingService) StartTrackingForce(projectName string) error {
	return s.db.StartTrackingForce(projectName)
}

// StopTracking stops the current active time tracking session
func (s *timeTrackingService) StopTracking() (*model.TimeEntry, error) {
	return s.db.StopTracking()
}

// GetActiveEntry returns the currently active time tracking entry, if any
func (s *timeTrackingService) GetActiveEntry() (*model.TimeEntry, error) {
	return s.db.GetActiveEntry()
}

// GetEntries returns a list of time entries, limited by the given count
func (s *timeTrackingService) GetEntries(limit int) ([]model.TimeEntry, error) {
	return s.db.GetEntries(limit)
}

// GetEntriesForProject returns time entries for a specific project
func (s *timeTrackingService) GetEntriesForProject(projectName string, limit int, sortOrder string) ([]model.TimeEntry, error) {
	return s.db.GetEntriesForProject(projectName, limit, sortOrder)
}

// ClearAllData removes all time entries and projects from the database
func (s *timeTrackingService) ClearAllData() error {
	return s.db.ClearAllEntries()
}

// GetProjects returns all projects in the system
func (s *timeTrackingService) GetProjects() ([]model.Project, error) {
	return s.db.GetAllProjects()
}

// GetOrCreateProject gets an existing project or creates a new one
func (s *timeTrackingService) GetOrCreateProject(name string) (*model.Project, error) {
	return s.db.GetOrCreateProject(name)
}

// RemoveProject removes a project and all its time entries
func (s *timeTrackingService) RemoveProject(name string) error {
	return s.db.RemoveProject(name)
}

// FormatDuration formats a duration into HH:MM:SS format
func (s *timeTrackingService) FormatDuration(duration time.Duration) string {
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}
