package repository

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nitschmann/hora/internal/model"
)

// TimeEntryWithPauses represents a time entry with its pause information
type TimeEntryWithPauses struct {
	model.TimeEntry
	PauseCount int           `json:"pause_count"`
	PauseTime  time.Duration `json:"pause_time"`
}

// TimeEntry defines the interface for time entry data operations
type TimeEntry interface {
	// Create creates a new time entry
	Create(ctx context.Context, projectID int, startTime time.Time) (*model.TimeEntry, error)

	// GetByID retrieves a time entry by its ID
	GetByID(ctx context.Context, id int) (*model.TimeEntry, error)

	// GetActive retrieves the currently active time entry
	GetActive(ctx context.Context) (*model.TimeEntry, error)

	// UpdateEndTime updates the end time and duration of a time entry
	UpdateEndTime(ctx context.Context, id int, endTime time.Time, duration time.Duration) error

	// StopAllActive stops all active time entries by setting their end time
	StopAllActive(ctx context.Context) error

	// GetByProject retrieves time entries for a specific project
	GetByProject(ctx context.Context, projectID int, limit int, sortOrder string) ([]model.TimeEntry, error)

	// GetByProjectName retrieves time entries for a project by name
	GetByProjectName(ctx context.Context, projectName string, limit int, sortOrder string) ([]model.TimeEntry, error)

	// GetByProjectIDOrName retrieves time entries for a project by ID (if numeric) or name
	GetByProjectIDOrName(ctx context.Context, projectIDOrName string, limit int, sortOrder string) ([]model.TimeEntry, error)

	// GetByProjectIDOrNameWithPauses retrieves time entries with pause information for a project by ID (if numeric) or name
	GetByProjectIDOrNameWithPauses(ctx context.Context, projectIDOrName string, limit int, sortOrder string, since *time.Time) ([]TimeEntryWithPauses, error)

	// GetAll retrieves all time entries with a limit
	GetAll(ctx context.Context, limit int) ([]model.TimeEntry, error)

	// DeleteByProject deletes all time entries for a specific project
	DeleteByProject(ctx context.Context, projectID int) error

	// DeleteAll deletes all time entries
	DeleteAll(ctx context.Context) error
}

// TimeEntryImpl implements TimeEntry using SQLite
type TimeEntryImpl struct {
	db *sql.DB
}

// NewTimeEntry creates a new time entry repository
func NewTimeEntry(db *sql.DB) TimeEntry {
	return &TimeEntryImpl{db: db}
}

// Create creates a new time entry
func (r *TimeEntryImpl) Create(ctx context.Context, projectID int, startTime time.Time) (*model.TimeEntry, error) {
	query, args, err := goqu.Insert("time_entries").Rows(goqu.Record{
		"project_id": projectID,
		"start_time": startTime,
	}).ToSQL()
	if err != nil {
		return nil, err
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Get the created time entry
	return r.GetByID(ctx, int(id))
}

// GetByID retrieves a time entry by its ID
func (r *TimeEntryImpl) GetByID(ctx context.Context, id int) (*model.TimeEntry, error) {
	query := `
		SELECT te.id, te.project_id, te.start_time, te.end_time, te.duration, te.created_at,
		       p.id as project_id2, p.name as project_name, p.created_at as project_created_at
		FROM time_entries te
		JOIN projects p ON te.project_id = p.id
		WHERE te.id = ?;`

	var entry model.TimeEntry
	var endTime *time.Time
	var duration *int64
	var project model.Project

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&entry.ID,
		&entry.ProjectID,
		&entry.StartTime,
		&endTime,
		&duration,
		&entry.CreatedAt,
		&project.ID,
		&project.Name,
		&project.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	entry.EndTime = endTime
	if duration != nil {
		d := time.Duration(*duration) * time.Second
		entry.Duration = &d
	}
	entry.Project = &project

	return &entry, nil
}

// GetActive retrieves the currently active time entry
func (r *TimeEntryImpl) GetActive(ctx context.Context) (*model.TimeEntry, error) {
	query := `
		SELECT te.id, te.project_id, te.start_time, te.end_time, te.duration, te.created_at,
		       p.id as project_id2, p.name as project_name, p.created_at as project_created_at
		FROM time_entries te
		JOIN projects p ON te.project_id = p.id
		WHERE te.end_time IS NULL
		ORDER BY te.start_time DESC
		LIMIT 1;`

	var entry model.TimeEntry
	var endTime *time.Time
	var duration *int64
	var project model.Project

	err := r.db.QueryRowContext(ctx, query).Scan(
		&entry.ID,
		&entry.ProjectID,
		&entry.StartTime,
		&endTime,
		&duration,
		&entry.CreatedAt,
		&project.ID,
		&project.Name,
		&project.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	entry.EndTime = endTime
	if duration != nil {
		d := time.Duration(*duration) * time.Second
		entry.Duration = &d
	}
	entry.Project = &project

	return &entry, nil
}

// UpdateEndTime updates the end time and duration of a time entry
func (r *TimeEntryImpl) UpdateEndTime(ctx context.Context, id int, endTime time.Time, duration time.Duration) error {
	durationSeconds := int64(duration.Seconds())
	query, args, err := goqu.Update("time_entries").
		Set(goqu.Record{
			"end_time": endTime,
			"duration": durationSeconds,
		}).
		Where(goqu.C("id").Eq(id)).
		ToSQL()
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

// StopAllActive stops all active time entries by setting their end time
func (r *TimeEntryImpl) StopAllActive(ctx context.Context) error {
	now := time.Now()
	query, args, err := goqu.Update("time_entries").
		Set(goqu.Record{
			"end_time": now,
			"duration": 0, // Duration will be calculated later
		}).
		Where(goqu.C("end_time").IsNull()).
		ToSQL()
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

// GetByProject retrieves time entries for a specific project
func (r *TimeEntryImpl) GetByProject(ctx context.Context, projectID int, limit int, sortOrder string) ([]model.TimeEntry, error) {
	var orderClause string
	switch sortOrder {
	case "asc":
		orderClause = "ORDER BY te.start_time ASC"
	case "desc":
		orderClause = "ORDER BY te.start_time DESC"
	default:
		orderClause = "ORDER BY te.start_time DESC"
	}

	query := `
		SELECT te.id, te.project_id, te.start_time, te.end_time, te.duration, te.created_at,
		       p.id as project_id2, p.name as project_name, p.created_at as project_created_at
		FROM time_entries te
		JOIN projects p ON te.project_id = p.id
		WHERE te.project_id = ?
		` + orderClause + `
		LIMIT ?;`

	rows, err := r.db.QueryContext(ctx, query, projectID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTimeEntries(rows)
}

// GetByProjectName retrieves time entries for a project by name
func (r *TimeEntryImpl) GetByProjectName(ctx context.Context, projectName string, limit int, sortOrder string) ([]model.TimeEntry, error) {
	var orderClause string
	switch sortOrder {
	case "asc":
		orderClause = "ORDER BY te.start_time ASC"
	case "desc":
		orderClause = "ORDER BY te.start_time DESC"
	default:
		orderClause = "ORDER BY te.start_time DESC"
	}

	query := `
		SELECT te.id, te.project_id, te.start_time, te.end_time, te.duration, te.created_at,
		       p.id as project_id2, p.name as project_name, p.created_at as project_created_at
		FROM time_entries te
		JOIN projects p ON te.project_id = p.id
		WHERE p.name = ?
		` + orderClause + `
		LIMIT ?;`

	rows, err := r.db.QueryContext(ctx, query, projectName, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTimeEntries(rows)
}

// GetAll retrieves all time entries with a limit
func (r *TimeEntryImpl) GetAll(ctx context.Context, limit int) ([]model.TimeEntry, error) {
	query, args, err := goqu.From("time_entries").
		Join(goqu.T("projects"), goqu.On(goqu.C("time_entries.project_id").Eq(goqu.C("projects.id")))).
		Select(
			"time_entries.id", "time_entries.project_id", "time_entries.start_time", "time_entries.end_time", "time_entries.duration", "time_entries.created_at",
			goqu.C("projects.id").As("project_id2"), goqu.C("projects.name").As("project_name"), goqu.C("projects.created_at").As("project_created_at"),
		).
		Order(goqu.C("time_entries.start_time").Desc()).
		Limit(uint(limit)).
		ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTimeEntries(rows)
}

// GetByProjectIDOrName retrieves time entries for a project by ID (if numeric) or name
func (r *TimeEntryImpl) GetByProjectIDOrName(ctx context.Context, projectIDOrName string, limit int, sortOrder string) ([]model.TimeEntry, error) {
	// Try to parse as integer first
	if projectID, err := strconv.Atoi(projectIDOrName); err == nil {
		// It's a numeric ID
		return r.GetByProject(ctx, projectID, limit, sortOrder)
	}

	// It's a name
	return r.GetByProjectName(ctx, projectIDOrName, limit, sortOrder)
}

// GetByProjectIDOrNameWithPauses retrieves time entries with pause information for a project by ID (if numeric) or name
func (r *TimeEntryImpl) GetByProjectIDOrNameWithPauses(ctx context.Context, projectIDOrName string, limit int, sortOrder string, since *time.Time) ([]TimeEntryWithPauses, error) {
	// Try to parse as integer first
	if projectID, err := strconv.Atoi(projectIDOrName); err == nil {
		// It's a numeric ID
		return r.GetByProjectWithPauses(ctx, projectID, limit, sortOrder, since)
	}

	// It's a name
	return r.GetByProjectNameWithPauses(ctx, projectIDOrName, limit, sortOrder, since)
}

// GetByProjectWithPauses retrieves time entries with pause information for a specific project
func (r *TimeEntryImpl) GetByProjectWithPauses(ctx context.Context, projectID int, limit int, sortOrder string, since *time.Time) ([]TimeEntryWithPauses, error) {
	var orderClause string
	switch sortOrder {
	case "asc":
		orderClause = "ORDER BY te.start_time ASC"
	case "desc":
		orderClause = "ORDER BY te.start_time DESC"
	default:
		orderClause = "ORDER BY te.start_time DESC"
	}

	query := `
		SELECT te.id, te.project_id, te.start_time, te.end_time, te.duration, te.created_at,
		       p.id as project_id2, p.name as project_name, p.created_at as project_created_at,
		       COALESCE(pause_stats.pause_count, 0) as pause_count,
		       COALESCE(pause_stats.total_pause_time, 0) as total_pause_time
		FROM time_entries te
		JOIN projects p ON te.project_id = p.id
		LEFT JOIN (
			SELECT time_entry_id, 
			       COUNT(*) as pause_count,
			       SUM(COALESCE(duration, 0)) as total_pause_time
			FROM pauses 
			WHERE pause_end IS NOT NULL
			GROUP BY time_entry_id
		) pause_stats ON te.id = pause_stats.time_entry_id
		WHERE te.project_id = ?`

	args := []interface{}{projectID}

	if since != nil {
		query += ` AND te.start_time >= ?`
		args = append(args, *since)
	}

	query += ` ` + orderClause + ` LIMIT ?;`
	args = append(args, limit)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTimeEntriesWithPauses(rows)
}

// GetByProjectNameWithPauses retrieves time entries with pause information for a project by name
func (r *TimeEntryImpl) GetByProjectNameWithPauses(ctx context.Context, projectName string, limit int, sortOrder string, since *time.Time) ([]TimeEntryWithPauses, error) {
	var orderClause string
	switch sortOrder {
	case "asc":
		orderClause = "ORDER BY te.start_time ASC"
	case "desc":
		orderClause = "ORDER BY te.start_time DESC"
	default:
		orderClause = "ORDER BY te.start_time DESC"
	}

	query := `
		SELECT te.id, te.project_id, te.start_time, te.end_time, te.duration, te.created_at,
		       p.id as project_id2, p.name as project_name, p.created_at as project_created_at,
		       COALESCE(pause_stats.pause_count, 0) as pause_count,
		       COALESCE(pause_stats.total_pause_time, 0) as total_pause_time
		FROM time_entries te
		JOIN projects p ON te.project_id = p.id
		LEFT JOIN (
			SELECT time_entry_id, 
			       COUNT(*) as pause_count,
			       SUM(COALESCE(duration, 0)) as total_pause_time
			FROM pauses 
			WHERE pause_end IS NOT NULL
			GROUP BY time_entry_id
		) pause_stats ON te.id = pause_stats.time_entry_id
		WHERE p.name = ?`

	args := []interface{}{projectName}

	if since != nil {
		query += ` AND te.start_time >= ?`
		args = append(args, *since)
	}

	query += ` ` + orderClause + ` LIMIT ?;`
	args = append(args, limit)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTimeEntriesWithPauses(rows)
}

// scanTimeEntriesWithPauses scans rows into TimeEntryWithPauses structs
func (r *TimeEntryImpl) scanTimeEntriesWithPauses(rows *sql.Rows) ([]TimeEntryWithPauses, error) {
	var entries []TimeEntryWithPauses

	for rows.Next() {
		var entry TimeEntryWithPauses
		var projectID2 int
		var projectName string
		var projectCreatedAt time.Time
		var pauseCount int
		var totalPauseTimeSeconds int64

		err := rows.Scan(
			&entry.ID,
			&entry.ProjectID,
			&entry.StartTime,
			&entry.EndTime,
			&entry.Duration,
			&entry.CreatedAt,
			&projectID2,
			&projectName,
			&projectCreatedAt,
			&pauseCount,
			&totalPauseTimeSeconds,
		)
		if err != nil {
			return nil, err
		}

		// Set project information
		entry.Project = &model.Project{
			ID:        projectID2,
			Name:      projectName,
			CreatedAt: projectCreatedAt,
		}

		// Set pause information
		entry.PauseCount = pauseCount
		entry.PauseTime = time.Duration(totalPauseTimeSeconds) * time.Second

		entries = append(entries, entry)
	}

	return entries, rows.Err()
}

// DeleteByProject deletes all time entries for a specific project
func (r *TimeEntryImpl) DeleteByProject(ctx context.Context, projectID int) error {
	query, args, err := goqu.Delete("time_entries").
		Where(goqu.C("project_id").Eq(projectID)).
		ToSQL()
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

// DeleteAll deletes all time entries
func (r *TimeEntryImpl) DeleteAll(ctx context.Context) error {
	query, args, err := goqu.Delete("time_entries").ToSQL()
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

// scanTimeEntries is a helper method to scan time entries from rows
func (r *TimeEntryImpl) scanTimeEntries(rows *sql.Rows) ([]model.TimeEntry, error) {
	var entries []model.TimeEntry

	for rows.Next() {
		var entry model.TimeEntry
		var endTime *time.Time
		var duration *int64
		var project model.Project

		err := rows.Scan(
			&entry.ID,
			&entry.ProjectID,
			&entry.StartTime,
			&endTime,
			&duration,
			&entry.CreatedAt,
			&project.ID,
			&project.Name,
			&project.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		entry.EndTime = endTime
		if duration != nil {
			d := time.Duration(*duration) * time.Second
			entry.Duration = &d
		}
		entry.Project = &project

		entries = append(entries, entry)
	}

	return entries, rows.Err()
}
