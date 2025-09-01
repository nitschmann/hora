package database

import (
	"fmt"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nitschmann/hora/internal/model"
)

// SQLiteDatabase implements the Database interface using SQLite
type SQLiteDatabase struct {
	db *sqlx.DB
}

// NewSQLiteDatabase creates a new SQLite database connection
func NewSQLiteDatabase() (Database, error) {
	dbPath, err := GetDatabasePath()
	if err != nil {
		return nil, fmt.Errorf("failed to get database path: %w", err)
	}

	db, err := sqlx.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	database := &SQLiteDatabase{db: db}

	// Initialize database schema
	if err := database.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return database, nil
}

func (d *SQLiteDatabase) initSchema() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS projects (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS time_entries (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			project_id INTEGER NOT NULL,
			start_time DATETIME NOT NULL,
			end_time DATETIME,
			duration INTEGER,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (project_id) REFERENCES projects (id)
		);`,
	}

	for _, query := range queries {
		if _, err := d.db.Exec(query); err != nil {
			return err
		}
	}

	return nil
}

// Project management methods
func (d *SQLiteDatabase) GetOrCreateProject(name string) (*model.Project, error) {
	// Try to get existing project first
	project, err := d.GetProjectByName(name)
	if err == nil && project != nil {
		return project, nil
	}

	// Create new project if it doesn't exist
	ds := goqu.Insert("projects").Rows(goqu.Record{"name": name})
	query, args, err := ds.ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build insert query: %w", err)
	}

	result, err := d.db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get project ID: %w", err)
	}

	return &model.Project{
		ID:        int(id),
		Name:      name,
		CreatedAt: time.Now(),
	}, nil
}

func (d *SQLiteDatabase) GetProject(id int) (*model.Project, error) {
	ds := goqu.Select("id", "name", "created_at").From("projects").Where(goqu.C("id").Eq(id))
	query, args, err := ds.ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	var project model.Project
	err = d.db.Get(&project, query, args...)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, err
	}
	return &project, nil
}

func (d *SQLiteDatabase) GetProjectByName(name string) (*model.Project, error) {
	ds := goqu.Select("id", "name", "created_at").From("projects").Where(goqu.C("name").Eq(name))
	query, args, err := ds.ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	var project model.Project
	err = d.db.Get(&project, query, args...)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, err
	}
	return &project, nil
}

func (d *SQLiteDatabase) GetAllProjects() ([]model.Project, error) {
	// Use raw SQL to get the last tracked time for each project
	query := `
		SELECT p.id, p.name, p.created_at, MAX(te.end_time) as last_tracked_at
		FROM projects p
		LEFT JOIN time_entries te ON p.id = te.project_id
		GROUP BY p.id, p.name, p.created_at
		ORDER BY p.name ASC`

	var results []struct {
		ID            int       `db:"id"`
		Name          string    `db:"name"`
		CreatedAt     time.Time `db:"created_at"`
		LastTrackedAt *string   `db:"last_tracked_at"`
	}

	err := d.db.Select(&results, query)
	if err != nil {
		return nil, err
	}

	var projects []model.Project
	for _, result := range results {
		var lastTrackedAt *time.Time
		if result.LastTrackedAt != nil && *result.LastTrackedAt != "" {
			// Try parsing with RFC3339Nano first (includes microseconds)
			if parsedTime, err := time.Parse(time.RFC3339Nano, *result.LastTrackedAt); err == nil {
				lastTrackedAt = &parsedTime
			} else if parsedTime, err := time.Parse(time.RFC3339, *result.LastTrackedAt); err == nil {
				// Fallback to RFC3339
				lastTrackedAt = &parsedTime
			}
		}

		projects = append(projects, model.Project{
			ID:            result.ID,
			Name:          result.Name,
			CreatedAt:     result.CreatedAt,
			LastTrackedAt: lastTrackedAt,
		})
	}

	return projects, nil
}

func (d *SQLiteDatabase) RemoveProject(name string) error {
	// First, get the project to check if it exists
	project, err := d.GetProjectByName(name)
	if err != nil {
		return fmt.Errorf("failed to get project: %w", err)
	}
	if project == nil {
		return fmt.Errorf("project '%s' not found", name)
	}

	// Delete all time entries for this project first (due to foreign key constraint)
	ds := goqu.Delete("time_entries").Where(goqu.C("project_id").Eq(project.ID))
	query, args, err := ds.ToSQL()
	if err != nil {
		return fmt.Errorf("failed to build delete query: %w", err)
	}
	_, err = d.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete time entries for project: %w", err)
	}

	// Delete the project
	ds = goqu.Delete("projects").Where(goqu.C("id").Eq(project.ID))
	query, args, err = ds.ToSQL()
	if err != nil {
		return fmt.Errorf("failed to build delete query: %w", err)
	}
	_, err = d.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	return nil
}

func (d *SQLiteDatabase) StartTracking(project string) error {
	// Check if there's already an active session
	activeEntry, err := d.GetActiveEntry()
	if err != nil {
		return fmt.Errorf("failed to check for active entry: %w", err)
	}

	if activeEntry != nil {
		return fmt.Errorf("already tracking time for project: %s", activeEntry.Project.Name)
	}

	// Get or create the project
	proj, err := d.GetOrCreateProject(project)
	if err != nil {
		return fmt.Errorf("failed to get or create project: %w", err)
	}

	ds := goqu.Insert("time_entries").Rows(goqu.Record{
		"project_id": proj.ID,
		"start_time": time.Now(),
	})
	query, args, err := ds.ToSQL()
	if err != nil {
		return fmt.Errorf("failed to build insert query: %w", err)
	}

	_, err = d.db.Exec(query, args...)
	return err
}

func (d *SQLiteDatabase) StartTrackingForce(project string) error {
	// Stop ALL active sessions first
	updateDs := goqu.Update("time_entries").
		Set(goqu.Record{
			"end_time": time.Now(),
		}).
		Where(goqu.C("end_time").IsNull())

	query, args, err := updateDs.ToSQL()
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	_, err = d.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to stop active sessions: %w", err)
	}

	// Get or create the project
	proj, err := d.GetOrCreateProject(project)
	if err != nil {
		return fmt.Errorf("failed to get or create project: %w", err)
	}

	insertDs := goqu.Insert("time_entries").Rows(goqu.Record{
		"project_id": proj.ID,
		"start_time": time.Now(),
	})
	query, args, err = insertDs.ToSQL()
	if err != nil {
		return fmt.Errorf("failed to build insert query: %w", err)
	}

	_, err = d.db.Exec(query, args...)
	return err
}

func (d *SQLiteDatabase) StopTracking() (*model.TimeEntry, error) {
	activeEntry, err := d.GetActiveEntry()
	if err != nil {
		return nil, fmt.Errorf("failed to get active entry: %w", err)
	}

	if activeEntry == nil {
		return nil, fmt.Errorf("no active time tracking session found")
	}

	endTime := time.Now()
	duration := endTime.Sub(activeEntry.StartTime)

	ds := goqu.Update("time_entries").
		Set(goqu.Record{
			"end_time": endTime,
			"duration": int64(duration),
		}).
		Where(goqu.C("id").Eq(activeEntry.ID))

	query, args, err := ds.ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build update query: %w", err)
	}

	_, err = d.db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to stop tracking: %w", err)
	}

	// Return the updated entry
	activeEntry.EndTime = &endTime
	activeEntry.Duration = &duration

	return activeEntry, nil
}

func (d *SQLiteDatabase) GetActiveEntry() (*model.TimeEntry, error) {
	query := `
		SELECT te.id, te.project_id, te.start_time, te.created_at, p.name
		FROM time_entries te
		JOIN projects p ON te.project_id = p.id
		WHERE te.end_time IS NULL 
		ORDER BY te.start_time DESC 
		LIMIT 1`

	var result struct {
		ID          int       `db:"id"`
		ProjectID   int       `db:"project_id"`
		StartTime   time.Time `db:"start_time"`
		CreatedAt   time.Time `db:"created_at"`
		ProjectName string    `db:"name"`
	}

	err := d.db.Get(&result, query)
	if err != nil {
		// Check if it's a "no rows" error
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, err
	}

	entry := &model.TimeEntry{
		ID:        result.ID,
		ProjectID: result.ProjectID,
		StartTime: result.StartTime,
		CreatedAt: result.CreatedAt,
		Project: &model.Project{
			ID:   result.ProjectID,
			Name: result.ProjectName,
		},
	}

	return entry, nil
}

func (d *SQLiteDatabase) GetEntries(limit int) ([]model.TimeEntry, error) {
	ds := goqu.Select(
		"te.id", "te.project_id", "te.start_time", "te.end_time", "te.duration", "te.created_at",
		goqu.C("p.id").As("project_id2"), goqu.C("p.name").As("project_name"), goqu.C("p.created_at").As("project_created_at"),
	).
		From(goqu.T("time_entries").As("te")).
		Join(goqu.T("projects").As("p"), goqu.On(goqu.C("te.project_id").Eq(goqu.C("p.id")))).
		Order(goqu.C("te.start_time").Desc()).
		Limit(uint(limit))

	query, args, err := ds.ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	var results []struct {
		ID               int        `db:"id"`
		ProjectID        int        `db:"project_id"`
		StartTime        time.Time  `db:"start_time"`
		EndTime          *time.Time `db:"end_time"`
		Duration         *int64     `db:"duration"`
		CreatedAt        time.Time  `db:"created_at"`
		ProjectID2       int        `db:"project_id2"`
		ProjectName      string     `db:"project_name"`
		ProjectCreatedAt time.Time  `db:"project_created_at"`
	}

	err = d.db.Select(&results, query, args...)
	if err != nil {
		return nil, err
	}

	var entries []model.TimeEntry
	for _, result := range results {
		entry := model.TimeEntry{
			ID:        result.ID,
			ProjectID: result.ProjectID,
			StartTime: result.StartTime,
			EndTime:   result.EndTime,
			CreatedAt: result.CreatedAt,
			Project: &model.Project{
				ID:        result.ProjectID2,
				Name:      result.ProjectName,
				CreatedAt: result.ProjectCreatedAt,
			},
		}

		if result.Duration != nil {
			duration := time.Duration(*result.Duration)
			entry.Duration = &duration
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

func (d *SQLiteDatabase) GetEntriesForProject(projectName string, limit int, sortOrder string) ([]model.TimeEntry, error) {
	// Validate sort order
	var orderClause string
	switch sortOrder {
	case "asc":
		orderClause = "ORDER BY te.start_time ASC"
	case "desc":
		orderClause = "ORDER BY te.start_time DESC"
	default:
		orderClause = "ORDER BY te.start_time DESC" // default to desc
	}

	// Use raw SQL to avoid goqu issues
	query := fmt.Sprintf(`
		SELECT te.id, te.project_id, te.start_time, te.end_time, te.duration, te.created_at,
		       p.id as project_id2, p.name as project_name, p.created_at as project_created_at
		FROM time_entries te
		JOIN projects p ON te.project_id = p.id
		WHERE p.name = ?
		%s
		LIMIT ?`, orderClause)

	var results []struct {
		ID               int        `db:"id"`
		ProjectID        int        `db:"project_id"`
		StartTime        time.Time  `db:"start_time"`
		EndTime          *time.Time `db:"end_time"`
		Duration         *int64     `db:"duration"`
		CreatedAt        time.Time  `db:"created_at"`
		ProjectID2       int        `db:"project_id2"`
		ProjectName      string     `db:"project_name"`
		ProjectCreatedAt time.Time  `db:"project_created_at"`
	}

	err := d.db.Select(&results, query, projectName, limit)
	if err != nil {
		return nil, err
	}

	var entries []model.TimeEntry
	for _, result := range results {
		entry := model.TimeEntry{
			ID:        result.ID,
			ProjectID: result.ProjectID,
			StartTime: result.StartTime,
			EndTime:   result.EndTime,
			CreatedAt: result.CreatedAt,
			Project: &model.Project{
				ID:        result.ProjectID2,
				Name:      result.ProjectName,
				CreatedAt: result.ProjectCreatedAt,
			},
		}

		if result.Duration != nil {
			duration := time.Duration(*result.Duration)
			entry.Duration = &duration
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

func (d *SQLiteDatabase) ClearAllEntries() error {
	// Delete all time entries first (due to foreign key constraint)
	ds := goqu.Delete("time_entries")
	query, args, err := ds.ToSQL()
	if err != nil {
		return fmt.Errorf("failed to build delete query: %w", err)
	}
	_, err = d.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to clear time entries: %w", err)
	}

	// Delete all projects
	ds = goqu.Delete("projects")
	query, args, err = ds.ToSQL()
	if err != nil {
		return fmt.Errorf("failed to build delete query: %w", err)
	}
	_, err = d.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to clear projects: %w", err)
	}

	return nil
}

func (d *SQLiteDatabase) Close() error {
	return d.db.Close()
}
