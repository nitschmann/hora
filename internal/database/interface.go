package database

import "github.com/nitschmann/hora/internal/model"

// Database defines the interface for time tracking database operations
type Database interface {
	// Project management
	GetOrCreateProject(name string) (*model.Project, error)
	GetProject(id int) (*model.Project, error)
	GetProjectByName(name string) (*model.Project, error)
	GetAllProjects() ([]model.Project, error)
	RemoveProject(name string) error

	// Time tracking
	StartTracking(project string) error
	StartTrackingForce(project string) error
	StopTracking() (*model.TimeEntry, error)
	GetActiveEntry() (*model.TimeEntry, error)
	GetEntries(limit int) ([]model.TimeEntry, error)
	GetEntriesForProject(projectName string, limit int, sortOrder string) ([]model.TimeEntry, error)

	// Data management
	ClearAllEntries() error

	// Connection management
	Close() error
}
