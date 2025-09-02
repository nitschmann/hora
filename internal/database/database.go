package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nitschmann/hora/internal/database/migrations"
	"github.com/nitschmann/hora/internal/model"
	"github.com/nitschmann/hora/internal/repository"
)

// SQLiteDatabase implements the Database interface using SQLite
type SQLiteDatabase struct {
	db                  *sql.DB
	projectRepository   repository.Project
	timeEntryRepository repository.TimeEntry
	pauseRepository     repository.Pause
}

// NewSQLiteDatabase creates a new SQLite database connection
func NewSQLiteDatabase() (*SQLiteDatabase, error) {
	dbPath, err := GetDatabasePath()
	if err != nil {
		return nil, fmt.Errorf("failed to get database path: %w", err)
	}

	// Open database connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable foreign key constraints
	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	// Run migrations
	ctx := context.Background()
	if err := migrations.RunMigrations(ctx, db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	// Create repositories
	projectRepo := repository.NewProject(db)
	timeEntryRepo := repository.NewTimeEntry(db)
	pauseRepo := repository.NewPause(db)

	return &SQLiteDatabase{
		db:                  db,
		projectRepository:   projectRepo,
		timeEntryRepository: timeEntryRepo,
		pauseRepository:     pauseRepo,
	}, nil
}

// GetDB returns the underlying database connection
func (d *SQLiteDatabase) GetDB() *sql.DB {
	return d.db
}

// GetProjectRepository returns the project repository
func (d *SQLiteDatabase) GetProjectRepository() repository.Project {
	return d.projectRepository
}

// GetTimeEntryRepository returns the time entry repository
func (d *SQLiteDatabase) GetTimeEntryRepository() repository.TimeEntry {
	return d.timeEntryRepository
}

// GetPauseRepository returns the pause repository
func (d *SQLiteDatabase) GetPauseRepository() repository.Pause {
	return d.pauseRepository
}

// Project management methods

// GetOrCreateProject retrieves a project by name, or creates it if it doesn't exist
func (d *SQLiteDatabase) GetOrCreateProject(ctx context.Context, name string) (*model.Project, error) {
	return d.projectRepository.GetOrCreate(ctx, name)
}

// GetProject retrieves a project by its ID
func (d *SQLiteDatabase) GetProject(ctx context.Context, id int) (*model.Project, error) {
	return d.projectRepository.GetByID(ctx, id)
}

// GetProjectByName retrieves a project by its name
func (d *SQLiteDatabase) GetProjectByName(ctx context.Context, name string) (*model.Project, error) {
	return d.projectRepository.GetByName(ctx, name)
}

// GetAllProjects retrieves all projects with their last tracked time
func (d *SQLiteDatabase) GetAllProjects(ctx context.Context) ([]model.Project, error) {
	return d.projectRepository.GetAll(ctx)
}

// GetProjectByIDOrName retrieves a project by ID (if numeric) or name
func (d *SQLiteDatabase) GetProjectByIDOrName(ctx context.Context, idOrName string) (*model.Project, error) {
	return d.projectRepository.GetByIDOrName(ctx, idOrName)
}

// RemoveProject removes a project and all its associated time entries
func (d *SQLiteDatabase) RemoveProject(ctx context.Context, name string) error {
	// Get project first to get its ID
	project, err := d.projectRepository.GetByName(ctx, name)
	if err != nil {
		return fmt.Errorf("project not found: %w", err)
	}

	// Delete time entries first (foreign key constraint will handle this, but being explicit)
	if err := d.timeEntryRepository.DeleteByProject(ctx, project.ID); err != nil {
		return fmt.Errorf("failed to delete time entries: %w", err)
	}

	// Delete project
	return d.projectRepository.Delete(ctx, name)
}

// RemoveProjectByIDOrName removes a project by ID (if numeric) or name and all its associated time entries
func (d *SQLiteDatabase) RemoveProjectByIDOrName(ctx context.Context, idOrName string) error {
	// Get project first to get its ID
	project, err := d.projectRepository.GetByIDOrName(ctx, idOrName)
	if err != nil {
		return fmt.Errorf("project not found: %w", err)
	}

	// Delete time entries first (foreign key constraint will handle this, but being explicit)
	if err := d.timeEntryRepository.DeleteByProject(ctx, project.ID); err != nil {
		return fmt.Errorf("failed to delete time entries: %w", err)
	}

	// Delete project by ID
	return d.projectRepository.DeleteByID(ctx, project.ID)
}

// Time tracking methods

// StartTracking starts tracking time for a project
func (d *SQLiteDatabase) StartTracking(ctx context.Context, project string) error {
	// Check for active entry
	activeEntry, err := d.timeEntryRepository.GetActive(ctx)
	if err == nil && activeEntry != nil {
		return fmt.Errorf("a time tracking session is already active for project '%s'", activeEntry.Project.Name)
	}

	// Get or create project
	proj, err := d.projectRepository.GetOrCreate(ctx, project)
	if err != nil {
		return fmt.Errorf("failed to get or create project: %w", err)
	}

	// Create new time entry
	_, err = d.timeEntryRepository.Create(ctx, proj.ID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to create time entry: %w", err)
	}

	return nil
}

// StartTrackingForce starts tracking time for a project, stopping any active session first
func (d *SQLiteDatabase) StartTrackingForce(ctx context.Context, project string) error {
	// Stop all active entries
	if err := d.timeEntryRepository.StopAllActive(ctx); err != nil {
		return fmt.Errorf("failed to stop active entries: %w", err)
	}

	// Get or create project
	proj, err := d.projectRepository.GetOrCreate(ctx, project)
	if err != nil {
		return fmt.Errorf("failed to get or create project: %w", err)
	}

	// Create new time entry
	_, err = d.timeEntryRepository.Create(ctx, proj.ID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to create time entry: %w", err)
	}

	return nil
}

// StopTracking stops the currently active time tracking session
func (d *SQLiteDatabase) StopTracking(ctx context.Context) (*model.TimeEntry, error) {
	// Get active entry
	activeEntry, err := d.timeEntryRepository.GetActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("no active time tracking session found: %w", err)
	}

	// End any active pause first
	activePause, err := d.pauseRepository.GetActivePause(ctx, activeEntry.ID)
	if err == nil {
		// There's an active pause, end it
		now := time.Now()
		pauseDuration := now.Sub(activePause.PauseStart)
		if err := d.pauseRepository.EndPause(ctx, activePause.ID, now, pauseDuration); err != nil {
			return nil, fmt.Errorf("failed to end active pause: %w", err)
		}
	}

	// Calculate total duration minus pause time
	endTime := time.Now()
	totalDuration := endTime.Sub(activeEntry.StartTime)

	// Get all pauses for this time entry to calculate total pause time
	pauses, err := d.pauseRepository.GetByTimeEntry(ctx, activeEntry.ID)
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
	if err := d.timeEntryRepository.UpdateEndTime(ctx, activeEntry.ID, endTime, workDuration); err != nil {
		return nil, fmt.Errorf("failed to update time entry: %w", err)
	}

	// Get updated entry
	updatedEntry, err := d.timeEntryRepository.GetByID(ctx, activeEntry.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated time entry: %w", err)
	}

	return updatedEntry, nil
}

// GetActiveEntry retrieves the currently active time entry
func (d *SQLiteDatabase) GetActiveEntry(ctx context.Context) (*model.TimeEntry, error) {
	return d.timeEntryRepository.GetActive(ctx)
}

// GetEntries retrieves all time entries with a limit
func (d *SQLiteDatabase) GetEntries(ctx context.Context, limit int) ([]model.TimeEntry, error) {
	return d.timeEntryRepository.GetAll(ctx, limit)
}

// GetEntriesForProject retrieves time entries for a specific project
func (d *SQLiteDatabase) GetEntriesForProject(ctx context.Context, projectName string, limit int, sortOrder string) ([]model.TimeEntry, error) {
	return d.timeEntryRepository.GetByProjectName(ctx, projectName, limit, sortOrder)
}

// GetEntriesForProjectByIDOrName retrieves time entries for a project by ID (if numeric) or name
func (d *SQLiteDatabase) GetEntriesForProjectByIDOrName(ctx context.Context, projectIDOrName string, limit int, sortOrder string) ([]model.TimeEntry, error) {
	return d.timeEntryRepository.GetByProjectIDOrName(ctx, projectIDOrName, limit, sortOrder)
}

// Data management methods

// ClearAllEntries removes all time entries and projects from the database
func (d *SQLiteDatabase) ClearAllEntries(ctx context.Context) error {
	// Delete all pauses first
	if err := d.pauseRepository.DeleteAll(ctx); err != nil {
		return fmt.Errorf("failed to delete pauses: %w", err)
	}

	// Delete all time entries (foreign key constraints will handle remaining pauses)
	if err := d.timeEntryRepository.DeleteAll(ctx); err != nil {
		return fmt.Errorf("failed to delete time entries: %w", err)
	}

	// Note: We don't delete projects as they might be referenced elsewhere
	// If you want to delete projects too, you would need to add a DeleteAll method to ProjectRepository

	return nil
}

// Pause management methods

// PauseTracking pauses the currently active time tracking session
func (d *SQLiteDatabase) PauseTracking(ctx context.Context) error {
	// Get the active time entry
	activeEntry, err := d.timeEntryRepository.GetActive(ctx)
	if err != nil {
		return fmt.Errorf("no active time tracking session found: %w", err)
	}

	// Check if there's already an active pause
	_, err = d.pauseRepository.GetActivePause(ctx, activeEntry.ID)
	if err == nil {
		return fmt.Errorf("time tracking is already paused")
	}

	// Create a new pause
	_, err = d.pauseRepository.Create(ctx, activeEntry.ID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to create pause: %w", err)
	}

	return nil
}

// ContinueTracking continues the currently paused time tracking session
func (d *SQLiteDatabase) ContinueTracking(ctx context.Context) error {
	// Get the active time entry
	activeEntry, err := d.timeEntryRepository.GetActive(ctx)
	if err != nil {
		return fmt.Errorf("no active time tracking session found: %w", err)
	}

	// Get the active pause
	activePause, err := d.pauseRepository.GetActivePause(ctx, activeEntry.ID)
	if err != nil {
		return fmt.Errorf("no active pause found: %w", err)
	}

	// End the pause
	now := time.Now()
	duration := now.Sub(activePause.PauseStart)
	err = d.pauseRepository.EndPause(ctx, activePause.ID, now, duration)
	if err != nil {
		return fmt.Errorf("failed to end pause: %w", err)
	}

	return nil
}

// Close closes the database connection
func (d *SQLiteDatabase) Close() error {
	return d.db.Close()
}
