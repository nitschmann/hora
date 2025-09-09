package service

import (
	"context"
	"fmt"
	"time"

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
	GetAllEntriesWithPauses(ctx context.Context, limit int, sortOrder string, since *time.Time) ([]repository.TimeEntryWithPauses, error)
	GetTotalTimeForProject(ctx context.Context, projectIDOrName string, since *time.Time) (time.Duration, error)
	ClearAllData(ctx context.Context) error
	GetProjects(ctx context.Context) ([]model.Project, error)
	GetOrCreateProject(ctx context.Context, name string) (*model.Project, error)
	GetProjectByIDOrName(ctx context.Context, idOrName string) (*model.Project, error)
	RemoveProject(ctx context.Context, idOrName string) error
	PauseTracking(ctx context.Context) error
	ContinueTracking(ctx context.Context) error
	FormatDuration(duration time.Duration) string
}

// timeTracking implements the TimeTracking interface
type timeTracking struct {
	projectRepo   repository.Project
	timeEntryRepo repository.TimeEntry
	pauseRepo     repository.Pause
}

// NewTimeTracking creates a new time tracking service
func NewTimeTracking(projectRepo repository.Project, timeEntryRepo repository.TimeEntry, pauseRepo repository.Pause) TimeTracking {
	return &timeTracking{
		projectRepo:   projectRepo,
		timeEntryRepo: timeEntryRepo,
		pauseRepo:     pauseRepo,
	}
}

// StartTracking starts a new time tracking session for the given project
func (s *timeTracking) StartTracking(ctx context.Context, projectName string, force bool) error {
	// Check for active entry if not forcing
	if !force {
		activeEntry, err := s.timeEntryRepo.GetActive(ctx)
		if err == nil && activeEntry != nil {
			return fmt.Errorf("a time tracking session is already active for project '%s'", activeEntry.Project.Name)
		}
	} else {
		// Stop all active entries when forcing
		if err := s.timeEntryRepo.StopAllActive(ctx); err != nil {
			return fmt.Errorf("failed to stop active entries: %w", err)
		}
	}

	// Get or create project
	proj, err := s.projectRepo.GetOrCreate(ctx, projectName)
	if err != nil {
		return fmt.Errorf("failed to get or create project: %w", err)
	}

	// Create new time entry
	_, err = s.timeEntryRepo.Create(ctx, proj.ID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to create time entry: %w", err)
	}

	return nil
}

// StopTracking stops the current active time tracking session
func (s *timeTracking) StopTracking(ctx context.Context) (*model.TimeEntry, error) {
	// Get active entry
	activeEntry, err := s.timeEntryRepo.GetActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("no active time tracking session found: %w", err)
	}

	// End any active pause first
	activePause, err := s.pauseRepo.GetActivePause(ctx, activeEntry.ID)
	if err == nil {
		// There's an active pause, end it
		now := time.Now()
		pauseDuration := now.Sub(activePause.PauseStart)
		if err := s.pauseRepo.EndPause(ctx, activePause.ID, now, pauseDuration); err != nil {
			return nil, fmt.Errorf("failed to end active pause: %w", err)
		}
	}

	// Calculate total duration minus pause time
	endTime := time.Now()
	totalDuration := endTime.Sub(activeEntry.StartTime)

	// Get all pauses for this time entry to calculate total pause time
	pauses, err := s.pauseRepo.GetByTimeEntry(ctx, activeEntry.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pauses: %w", err)
	}

	var totalPauseTime time.Duration
	for _, pause := range pauses {
		if pause.Duration != nil {
			totalPauseTime += *pause.Duration
		}
	}

	// Calculate actual work duration (total time minus pause time)
	workDuration := totalDuration - totalPauseTime

	// Update the entry
	if err := s.timeEntryRepo.UpdateEndTime(ctx, activeEntry.ID, endTime, workDuration); err != nil {
		return nil, fmt.Errorf("failed to update time entry: %w", err)
	}

	// Get updated entry
	updatedEntry, err := s.timeEntryRepo.GetByID(ctx, activeEntry.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated time entry: %w", err)
	}

	return updatedEntry, nil
}

// GetActiveEntry returns the currently active time tracking entry, if any
func (s *timeTracking) GetActiveEntry(ctx context.Context) (*model.TimeEntry, error) {
	return s.timeEntryRepo.GetActive(ctx)
}

// GetEntries returns a list of time entries, limited by the given count
func (s *timeTracking) GetEntries(ctx context.Context, limit int) ([]model.TimeEntry, error) {
	return s.timeEntryRepo.GetAll(ctx, limit)
}

// GetEntriesForProject returns time entries for a project by ID (if numeric) or name
func (s *timeTracking) GetEntriesForProject(ctx context.Context, projectIDOrName string, limit int, sortOrder string) ([]model.TimeEntry, error) {
	return s.timeEntryRepo.GetByProjectIDOrName(ctx, projectIDOrName, limit, sortOrder)
}

// GetEntriesForProjectWithPauses returns time entries with pause information for a project by ID (if numeric) or name
func (s *timeTracking) GetEntriesForProjectWithPauses(ctx context.Context, projectIDOrName string, limit int, sortOrder string, since *time.Time) ([]repository.TimeEntryWithPauses, error) {
	return s.timeEntryRepo.GetByProjectIDOrNameWithPauses(ctx, projectIDOrName, limit, sortOrder, since)
}

// GetAllEntriesWithPauses returns all time entries with pause information across all projects
func (s *timeTracking) GetAllEntriesWithPauses(ctx context.Context, limit int, sortOrder string, since *time.Time) ([]repository.TimeEntryWithPauses, error) {
	return s.timeEntryRepo.GetAllWithPauses(ctx, limit, sortOrder, since)
}

// GetTotalTimeForProject returns the total tracked time for a project by ID (if numeric) or name
func (s *timeTracking) GetTotalTimeForProject(ctx context.Context, projectIDOrName string, since *time.Time) (time.Duration, error) {
	return s.timeEntryRepo.GetTotalTimeByProjectIDOrName(ctx, projectIDOrName, since)
}

// ClearAllData removes all time entries and projects from the database
func (s *timeTracking) ClearAllData(ctx context.Context) error {
	// Delete all pauses first
	if err := s.pauseRepo.DeleteAll(ctx); err != nil {
		return fmt.Errorf("failed to delete pauses: %w", err)
	}

	// Delete all time entries
	if err := s.timeEntryRepo.DeleteAll(ctx); err != nil {
		return fmt.Errorf("failed to delete time entries: %w", err)
	}

	return nil
}

// GetProjects returns all projects in the system
func (s *timeTracking) GetProjects(ctx context.Context) ([]model.Project, error) {
	return s.projectRepo.GetAll(ctx)
}

// GetOrCreateProject gets an existing project or creates a new one
func (s *timeTracking) GetOrCreateProject(ctx context.Context, name string) (*model.Project, error) {
	return s.projectRepo.GetOrCreate(ctx, name)
}

// GetProjectByIDOrName gets a project by ID (if numeric) or name
func (s *timeTracking) GetProjectByIDOrName(ctx context.Context, idOrName string) (*model.Project, error) {
	return s.projectRepo.GetByIDOrName(ctx, idOrName)
}

// RemoveProject removes a project by ID (if numeric) or name and all its time entries
func (s *timeTracking) RemoveProject(ctx context.Context, idOrName string) error {
	// Get project first to get its ID
	project, err := s.projectRepo.GetByIDOrName(ctx, idOrName)
	if err != nil {
		return fmt.Errorf("project not found: %w", err)
	}

	// Delete time entries first
	if err := s.timeEntryRepo.DeleteByProject(ctx, project.ID); err != nil {
		return fmt.Errorf("failed to delete time entries: %w", err)
	}

	// Delete project by ID
	return s.projectRepo.DeleteByID(ctx, project.ID)
}

// PauseTracking pauses the currently active time tracking session
func (s *timeTracking) PauseTracking(ctx context.Context) error {
	// Get the active time entry
	activeEntry, err := s.timeEntryRepo.GetActive(ctx)
	if err != nil {
		return fmt.Errorf("no active time tracking session found: %w", err)
	}

	// Check if there's already an active pause
	_, err = s.pauseRepo.GetActivePause(ctx, activeEntry.ID)
	if err == nil {
		return fmt.Errorf("time tracking is already paused")
	}

	// Create a new pause
	_, err = s.pauseRepo.Create(ctx, activeEntry.ID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to create pause: %w", err)
	}

	return nil
}

// ContinueTracking continues the currently paused time tracking session
func (s *timeTracking) ContinueTracking(ctx context.Context) error {
	// Get the active time entry
	activeEntry, err := s.timeEntryRepo.GetActive(ctx)
	if err != nil {
		return fmt.Errorf("no active time tracking session found: %w", err)
	}

	// Get the active pause
	activePause, err := s.pauseRepo.GetActivePause(ctx, activeEntry.ID)
	if err != nil {
		return fmt.Errorf("no active pause found: %w", err)
	}

	// End the pause
	now := time.Now()
	duration := now.Sub(activePause.PauseStart)
	err = s.pauseRepo.EndPause(ctx, activePause.ID, now, duration)
	if err != nil {
		return fmt.Errorf("failed to end pause: %w", err)
	}

	return nil
}

// FormatDuration formats a duration into HH:MM:SS format
func (s *timeTracking) FormatDuration(duration time.Duration) string {
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}
