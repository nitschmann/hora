package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/doug-martin/goqu/v9"

	"github.com/nitschmann/hora/internal/model"
)

const timeEntryTable = "time_entries"

// TimeEntryWithPauses represents a time entry with its pause information
type TimeEntryWithPauses struct {
	model.TimeEntry
	PauseCount int           `json:"pause_count"`
	PauseTime  time.Duration `json:"pause_time"`
}

// TimeEntry defines the interface for time entry data operations
type TimeEntry interface {
	// Create creates a new time entry
	Create(ctx context.Context, projectID int, startTime time.Time, category *string) (*model.TimeEntry, error)
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
	// GetTotalTimeByProjectIDOrName retrieves the total tracked time for a project by ID (if numeric) or name
	GetTotalTimeByProjectIDOrName(ctx context.Context, projectIDOrName string, since *time.Time) (time.Duration, error)
	// GetAllWithPauses retrieves all time entries with pause information across all projects
	GetAllWithPauses(ctx context.Context, limit int, sortOrder string, since *time.Time) ([]TimeEntryWithPauses, error)
	// GetAllWithPausesByCategory retrieves all time entries with pause information across all projects filtered by category
	GetAllWithPausesByCategory(ctx context.Context, limit int, sortOrder string, since *time.Time, category *string) ([]TimeEntryWithPauses, error)
	// GetAll retrieves all time entries with a limit
	GetAll(ctx context.Context, limit int) ([]model.TimeEntry, error)
	// DeleteByProject deletes all time entries for a specific project
	DeleteByProject(ctx context.Context, projectID int) error
	// DeleteAll deletes all time entries
	DeleteAll(ctx context.Context) error
}

// timeEntry implements TimeEntry using SQLite
type timeEntry struct {
	db *sql.DB
}

// NewTimeEntry creates a new time entry repository
func NewTimeEntry(db *sql.DB) TimeEntry {
	return &timeEntry{db: db}
}

// Create creates a new time entry
func (r *timeEntry) Create(ctx context.Context, projectID int, startTime time.Time, category *string) (*model.TimeEntry, error) {
	record := goqu.Record{
		"project_id": projectID,
		"start_time": startTime,
	}

	if category != nil {
		record["category"] = *category
	}

	query, args, err := goqu.Insert(timeEntryTable).Rows(record).ToSQL()
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

	return r.GetByID(ctx, int(id))
}

// GetByID retrieves a time entry by its ID
func (r *timeEntry) GetByID(ctx context.Context, id int) (*model.TimeEntry, error) {
	query, args, err := goqu.From(goqu.T(timeEntryTable).As("te")).
		Select(
			goqu.I("te.id"),
			goqu.I("te.project_id"),
			goqu.I("te.start_time"),
			goqu.I("te.end_time"),
			goqu.I("te.duration"),
			goqu.I("te.category"),
			goqu.I("te.created_at"),
			goqu.I("p.id").As("project_id2"),
			goqu.I("p.name").As("project_name"),
			goqu.I("p.created_at").As("project_created_at"),
		).
		Join(goqu.T(projectTable).As("p"), goqu.On(goqu.I("te.project_id").Eq(goqu.I("p.id")))).
		Where(goqu.I("te.id").Eq(id)).
		ToSQL()
	if err != nil {
		return nil, err
	}

	var entry model.TimeEntry
	var endTime *time.Time
	var duration *int64
	var category *string
	var project model.Project

	err = r.db.QueryRowContext(ctx, query, args...).Scan(
		&entry.ID,
		&entry.ProjectID,
		&entry.StartTime,
		&endTime,
		&duration,
		&category,
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
	entry.Category = category
	entry.Project = &project

	return &entry, nil
}

// GetActive retrieves the currently active time entry
func (r *timeEntry) GetActive(ctx context.Context) (*model.TimeEntry, error) {
	query, args, err := goqu.From(goqu.T(timeEntryTable).As("te")).
		Select(
			goqu.I("te.id"),
			goqu.I("te.project_id"),
			goqu.I("te.start_time"),
			goqu.I("te.end_time"),
			goqu.I("te.duration"),
			goqu.I("te.category"),
			goqu.I("te.created_at"),
			goqu.I("p.id").As("project_id2"),
			goqu.I("p.name").As("project_name"),
			goqu.I("p.created_at").As("project_created_at"),
		).
		Join(goqu.T(projectTable).As("p"), goqu.On(goqu.I("te.project_id").Eq(goqu.I("p.id")))).
		Where(goqu.I("te.end_time").IsNull()).
		Order(goqu.I("te.start_time").Desc()).
		Limit(1).
		ToSQL()
	if err != nil {
		return nil, err
	}

	var entry model.TimeEntry
	var endTime *time.Time
	var duration *int64
	var category *string
	var project model.Project

	err = r.db.QueryRowContext(ctx, query, args...).Scan(
		&entry.ID,
		&entry.ProjectID,
		&entry.StartTime,
		&endTime,
		&duration,
		&category,
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
	entry.Category = category
	entry.Project = &project

	return &entry, nil
}

// UpdateEndTime updates the end time and duration of a time entry
func (r *timeEntry) UpdateEndTime(ctx context.Context, id int, endTime time.Time, duration time.Duration) error {
	durationSeconds := int64(duration.Seconds())
	query, args, err := goqu.Update(timeEntryTable).
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
func (r *timeEntry) StopAllActive(ctx context.Context) error {
	now := time.Now()
	query, args, err := goqu.Update(timeEntryTable).
		Set(goqu.Record{
			"end_time": now,
			"duration": 0, // duration will be calculated later
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
func (r *timeEntry) GetByProject(ctx context.Context, projectID int, limit int, sortOrder string) ([]model.TimeEntry, error) {
	var orderDirection string
	switch sortOrder {
	case "asc":
		orderDirection = "ASC"
	case "desc":
		orderDirection = "DESC"
	default:
		orderDirection = "DESC"
	}

	queryBuilder := goqu.From(goqu.T(timeEntryTable).As("te")).
		Select(
			goqu.I("te.id"),
			goqu.I("te.project_id"),
			goqu.I("te.start_time"),
			goqu.I("te.end_time"),
			goqu.I("te.duration"),
			goqu.I("te.category"),
			goqu.I("te.created_at"),
			goqu.I("p.id").As("project_id2"),
			goqu.I("p.name").As("project_name"),
			goqu.I("p.created_at").As("project_created_at"),
		).
		Join(goqu.T(projectTable).As("p"), goqu.On(goqu.I("te.project_id").Eq(goqu.I("p.id")))).
		Where(goqu.I("te.project_id").Eq(projectID))

	if orderDirection == "ASC" {
		queryBuilder = queryBuilder.Order(goqu.I("te.start_time").Asc())
	} else {
		queryBuilder = queryBuilder.Order(goqu.I("te.start_time").Desc())
	}

	query, args, err := queryBuilder.Limit(uint(limit)).ToSQL()
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

// GetByProjectName retrieves time entries for a project by name
func (r *timeEntry) GetByProjectName(ctx context.Context, projectName string, limit int, sortOrder string) ([]model.TimeEntry, error) {
	var orderDirection string
	switch sortOrder {
	case "asc":
		orderDirection = "ASC"
	case "desc":
		orderDirection = "DESC"
	default:
		orderDirection = "DESC"
	}

	queryBuilder := goqu.From(goqu.T(timeEntryTable).As("te")).
		Select(
			goqu.I("te.id"),
			goqu.I("te.project_id"),
			goqu.I("te.start_time"),
			goqu.I("te.end_time"),
			goqu.I("te.duration"),
			goqu.I("te.category"),
			goqu.I("te.created_at"),
			goqu.I("p.id").As("project_id2"),
			goqu.I("p.name").As("project_name"),
			goqu.I("p.created_at").As("project_created_at"),
		).
		Join(goqu.T(projectTable).As("p"), goqu.On(goqu.I("te.project_id").Eq(goqu.I("p.id")))).
		Where(goqu.I("p.name").Eq(projectName))

	if orderDirection == "ASC" {
		queryBuilder = queryBuilder.Order(goqu.I("te.start_time").Asc())
	} else {
		queryBuilder = queryBuilder.Order(goqu.I("te.start_time").Desc())
	}

	query, args, err := queryBuilder.Limit(uint(limit)).ToSQL()
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

// GetAll retrieves all time entries with a limit
func (r *timeEntry) GetAll(ctx context.Context, limit int) ([]model.TimeEntry, error) {
	query, args, err := goqu.From(timeEntryTable).
		Join(goqu.T(projectTable), goqu.On(goqu.C("time_entries.project_id").Eq(goqu.C("projects.id")))).
		Select(
			"time_entries.id", "time_entries.project_id", "time_entries.start_time", "time_entries.end_time", "time_entries.duration", "time_entries.category", "time_entries.created_at",
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
func (r *timeEntry) GetByProjectIDOrName(ctx context.Context, projectIDOrName string, limit int, sortOrder string) ([]model.TimeEntry, error) {
	if projectID, err := strconv.Atoi(projectIDOrName); err == nil {
		return r.GetByProject(ctx, projectID, limit, sortOrder)
	}

	return r.GetByProjectName(ctx, projectIDOrName, limit, sortOrder)
}

// GetByProjectIDOrNameWithPauses retrieves time entries with pause information for a project by ID (if numeric) or name
func (r *timeEntry) GetByProjectIDOrNameWithPauses(ctx context.Context, projectIDOrName string, limit int, sortOrder string, since *time.Time) ([]TimeEntryWithPauses, error) {
	if projectID, err := strconv.Atoi(projectIDOrName); err == nil {
		return r.GetByProjectWithPauses(ctx, projectID, limit, sortOrder, since)
	}

	return r.GetByProjectNameWithPauses(ctx, projectIDOrName, limit, sortOrder, since)
}

// GetByProjectWithPauses retrieves time entries with pause information for a specific project
func (r *timeEntry) GetByProjectWithPauses(ctx context.Context, projectID int, limit int, sortOrder string, since *time.Time) ([]TimeEntryWithPauses, error) {
	var orderDirection string
	switch sortOrder {
	case "asc":
		orderDirection = "ASC"
	case "desc":
		orderDirection = "DESC"
	default:
		orderDirection = "DESC"
	}

	// Build the subquery for pause stats
	pauseStatsSubquery := goqu.From(pauseTable).
		Select(
			goqu.I("time_entry_id"),
			goqu.COUNT("*").As("pause_count"),
			goqu.SUM(goqu.COALESCE(goqu.I("duration"), 0)).As("total_pause_time"),
		).
		Where(goqu.I("pause_end").IsNotNull()).
		GroupBy(goqu.I("time_entry_id"))

	// main query
	queryBuilder := goqu.From(goqu.T(timeEntryTable).As("te")).
		Select(
			goqu.I("te.id"),
			goqu.I("te.project_id"),
			goqu.I("te.start_time"),
			goqu.I("te.end_time"),
			goqu.I("te.duration"),
			goqu.I("te.category"),
			goqu.I("te.created_at"),
			goqu.I("p.id").As("project_id2"),
			goqu.I("p.name").As("project_name"),
			goqu.I("p.created_at").As("project_created_at"),
			goqu.COALESCE(goqu.I("pause_stats.pause_count"), 0).As("pause_count"),
			goqu.COALESCE(goqu.I("pause_stats.total_pause_time"), 0).As("total_pause_time"),
		).
		Join(goqu.T(projectTable).As("p"), goqu.On(goqu.I("te.project_id").Eq(goqu.I("p.id")))).
		LeftJoin(pauseStatsSubquery.As("pause_stats"), goqu.On(goqu.I("te.id").Eq(goqu.I("pause_stats.time_entry_id")))).
		Where(goqu.I("te.project_id").Eq(projectID))

	if since != nil {
		queryBuilder = queryBuilder.Where(goqu.I("te.start_time").Gte(*since))
	}

	if orderDirection == "ASC" {
		queryBuilder = queryBuilder.Order(goqu.I("te.start_time").Asc())
	} else {
		queryBuilder = queryBuilder.Order(goqu.I("te.start_time").Desc())
	}

	query, args, err := queryBuilder.Limit(uint(limit)).ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTimeEntriesWithPauses(rows)
}

// GetByProjectNameWithPauses retrieves time entries with pause information for a project by name
func (r *timeEntry) GetByProjectNameWithPauses(ctx context.Context, projectName string, limit int, sortOrder string, since *time.Time) ([]TimeEntryWithPauses, error) {
	var orderDirection string
	switch sortOrder {
	case "asc":
		orderDirection = "ASC"
	case "desc":
		orderDirection = "DESC"
	default:
		orderDirection = "DESC"
	}

	// Build the subquery for pause stats
	pauseStatsSubquery := goqu.From(pauseTable).
		Select(
			goqu.I("time_entry_id"),
			goqu.COUNT("*").As("pause_count"),
			goqu.SUM(goqu.COALESCE(goqu.I("duration"), 0)).As("total_pause_time"),
		).
		Where(goqu.I("pause_end").IsNotNull()).
		GroupBy(goqu.I("time_entry_id"))

	// main query
	queryBuilder := goqu.From(goqu.T(timeEntryTable).As("te")).
		Select(
			goqu.I("te.id"),
			goqu.I("te.project_id"),
			goqu.I("te.start_time"),
			goqu.I("te.end_time"),
			goqu.I("te.duration"),
			goqu.I("te.category"),
			goqu.I("te.created_at"),
			goqu.I("p.id").As("project_id2"),
			goqu.I("p.name").As("project_name"),
			goqu.I("p.created_at").As("project_created_at"),
			goqu.COALESCE(goqu.I("pause_stats.pause_count"), 0).As("pause_count"),
			goqu.COALESCE(goqu.I("pause_stats.total_pause_time"), 0).As("total_pause_time"),
		).
		Join(goqu.T(projectTable).As("p"), goqu.On(goqu.I("te.project_id").Eq(goqu.I("p.id")))).
		LeftJoin(pauseStatsSubquery.As("pause_stats"), goqu.On(goqu.I("te.id").Eq(goqu.I("pause_stats.time_entry_id")))).
		Where(goqu.I("p.name").Eq(projectName))

	if since != nil {
		queryBuilder = queryBuilder.Where(goqu.I("te.start_time").Gte(*since))
	}

	if orderDirection == "ASC" {
		queryBuilder = queryBuilder.Order(goqu.I("te.start_time").Asc())
	} else {
		queryBuilder = queryBuilder.Order(goqu.I("te.start_time").Desc())
	}

	query, args, err := queryBuilder.Limit(uint(limit)).ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTimeEntriesWithPauses(rows)
}

// scanTimeEntriesWithPauses scans rows into TimeEntryWithPauses structs
func (r *timeEntry) scanTimeEntriesWithPauses(rows *sql.Rows) ([]TimeEntryWithPauses, error) {
	var entries []TimeEntryWithPauses

	for rows.Next() {
		var entry TimeEntryWithPauses
		var projectID2 int
		var projectName string
		var projectCreatedAt time.Time
		var pauseCount int
		var totalPauseTimeSeconds int64
		var durationSeconds *int64
		var category *string

		err := rows.Scan(
			&entry.ID,
			&entry.ProjectID,
			&entry.StartTime,
			&entry.EndTime,
			&durationSeconds,
			&category,
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

		if durationSeconds != nil {
			duration := time.Duration(*durationSeconds) * time.Second
			entry.Duration = &duration
		}

		entry.Category = category

		entry.Project = &model.Project{
			ID:        projectID2,
			Name:      projectName,
			CreatedAt: projectCreatedAt,
		}

		entry.PauseCount = pauseCount
		entry.PauseTime = time.Duration(totalPauseTimeSeconds) * time.Second

		entries = append(entries, entry)
	}

	return entries, rows.Err()
}

// DeleteByProject deletes all time entries for a specific project
func (r *timeEntry) DeleteByProject(ctx context.Context, projectID int) error {
	query, args, err := goqu.Delete(timeEntryTable).
		Where(goqu.C("project_id").Eq(projectID)).
		ToSQL()
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

// DeleteAll deletes all time entries
func (r *timeEntry) DeleteAll(ctx context.Context) error {
	query, args, err := goqu.Delete(timeEntryTable).ToSQL()
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

// scanTimeEntries is a helper method to scan time entries from rows
func (r *timeEntry) scanTimeEntries(rows *sql.Rows) ([]model.TimeEntry, error) {
	var entries []model.TimeEntry

	for rows.Next() {
		var entry model.TimeEntry
		var endTime *time.Time
		var duration *int64
		var category *string
		var project model.Project

		err := rows.Scan(
			&entry.ID,
			&entry.ProjectID,
			&entry.StartTime,
			&endTime,
			&duration,
			&category,
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
		entry.Category = category
		entry.Project = &project

		entries = append(entries, entry)
	}

	return entries, rows.Err()
}

// GetTotalTimeByProjectIDOrName retrieves the total tracked time for a project by ID (if numeric) or name
func (r *timeEntry) GetTotalTimeByProjectIDOrName(ctx context.Context, projectIDOrName string, since *time.Time) (time.Duration, error) {
	if projectID, err := strconv.Atoi(projectIDOrName); err == nil {
		return r.GetTotalTimeByProject(ctx, projectID, since)
	}

	return r.GetTotalTimeByProjectName(ctx, projectIDOrName, since)
}

// GetTotalTimeByProject retrieves the total tracked time for a specific project
func (r *timeEntry) GetTotalTimeByProject(ctx context.Context, projectID int, since *time.Time) (time.Duration, error) {
	pauseStatsSubquery := goqu.From(pauseTable).
		Select(
			goqu.I("time_entry_id"),
			goqu.SUM(goqu.COALESCE(goqu.I("duration"), 0)).As("total_pause_time"),
		).
		Where(goqu.I("pause_end").IsNotNull()).
		GroupBy(goqu.I("time_entry_id"))

	// main query
	queryBuilder := goqu.From(goqu.T(timeEntryTable).As("te")).
		Select(
			goqu.COALESCE(
				goqu.SUM(goqu.L("te.duration + COALESCE(pause_stats.total_pause_time, 0)")),
				0,
			).As("total_time"),
		).
		LeftJoin(pauseStatsSubquery.As("pause_stats"), goqu.On(goqu.I("te.id").Eq(goqu.I("pause_stats.time_entry_id")))).
		Where(
			goqu.I("te.project_id").Eq(projectID),
			goqu.I("te.end_time").IsNotNull(),
		)

	if since != nil {
		queryBuilder = queryBuilder.Where(goqu.I("te.start_time").Gte(*since))
	}

	query, args, err := queryBuilder.ToSQL()
	if err != nil {
		return 0, err
	}

	var totalTimeSeconds int64
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&totalTimeSeconds)
	if err != nil {
		return 0, err
	}

	return time.Duration(totalTimeSeconds) * time.Second, nil
}

// GetTotalTimeByProjectName retrieves the total tracked time for a project by name
func (r *timeEntry) GetTotalTimeByProjectName(ctx context.Context, projectName string, since *time.Time) (time.Duration, error) {
	pauseStatsSubquery := goqu.From(pauseTable).
		Select(
			goqu.I("time_entry_id"),
			goqu.SUM(goqu.COALESCE(goqu.I("duration"), 0)).As("total_pause_time"),
		).
		Where(goqu.I("pause_end").IsNotNull()).
		GroupBy(goqu.I("time_entry_id"))

	queryBuilder := goqu.From(goqu.T(timeEntryTable).As("te")).
		Select(
			goqu.COALESCE(
				goqu.SUM(goqu.L("te.duration + COALESCE(pause_stats.total_pause_time, 0)")),
				0,
			).As("total_time"),
		).
		Join(goqu.T(projectTable).As("p"), goqu.On(goqu.I("te.project_id").Eq(goqu.I("p.id")))).
		LeftJoin(pauseStatsSubquery.As("pause_stats"), goqu.On(goqu.I("te.id").Eq(goqu.I("pause_stats.time_entry_id")))).
		Where(
			goqu.I("p.name").Eq(projectName),
			goqu.I("te.end_time").IsNotNull(),
		)

	if since != nil {
		queryBuilder = queryBuilder.Where(goqu.I("te.start_time").Gte(*since))
	}

	query, args, err := queryBuilder.ToSQL()
	if err != nil {
		return 0, err
	}

	var totalTimeSeconds int64
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&totalTimeSeconds)
	if err != nil {
		return 0, err
	}

	return time.Duration(totalTimeSeconds) * time.Second, nil
}

// GetAllWithPauses retrieves all time entries with pause information across all projects
func (r *timeEntry) GetAllWithPauses(ctx context.Context, limit int, sortOrder string, since *time.Time) ([]TimeEntryWithPauses, error) {
	var orderClause string
	switch sortOrder {
	case "asc":
		orderClause = "ORDER BY te.start_time ASC"
	case "desc":
		orderClause = "ORDER BY te.start_time DESC"
	default:
		orderClause = "ORDER BY te.start_time DESC"
	}

	query := fmt.Sprintf(`
		SELECT te.id, te.project_id, te.start_time, te.end_time, te.duration, te.category, te.created_at,
		       p.id, p.name, p.created_at,
		       COALESCE(pause_stats.pause_count, 0) as pause_count,
		       COALESCE(pause_stats.total_pause_time, 0) as total_pause_time
		FROM %s te
		JOIN %s p ON te.project_id = p.id
		LEFT JOIN (
			SELECT time_entry_id,
			       COUNT(*) as pause_count,
			       SUM(COALESCE(duration, 0)) as total_pause_time
			FROM %s
			WHERE pause_end IS NOT NULL
			GROUP BY time_entry_id
		) pause_stats ON te.id = pause_stats.time_entry_id`, timeEntryTable, projectTable, pauseTable)

	args := []interface{}{}

	if since != nil {
		query += ` WHERE te.start_time >= ?`
		args = append(args, *since)
	}

	query += ` ` + orderClause

	// limit
	query += ` LIMIT ?`
	args = append(args, limit)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTimeEntriesWithPauses(rows)
}

// GetAllWithPausesByCategory retrieves all time entries with pause information across all projects filtered by category
func (r *timeEntry) GetAllWithPausesByCategory(ctx context.Context, limit int, sortOrder string, since *time.Time, category *string) ([]TimeEntryWithPauses, error) {
	var orderClause string
	switch sortOrder {
	case "asc":
		orderClause = "ORDER BY te.start_time ASC"
	case "desc":
		orderClause = "ORDER BY te.start_time DESC"
	default:
		orderClause = "ORDER BY te.start_time DESC"
	}

	query := fmt.Sprintf(`
		SELECT te.id, te.project_id, te.start_time, te.end_time, te.duration, te.category, te.created_at,
		       p.id, p.name, p.created_at,
		       COALESCE(pause_stats.pause_count, 0) as pause_count,
		       COALESCE(pause_stats.total_pause_time, 0) as total_pause_time
		FROM %s te
		JOIN %s p ON te.project_id = p.id
		LEFT JOIN (
			SELECT time_entry_id,
			       COUNT(*) as pause_count,
			       SUM(COALESCE(duration, 0)) as total_pause_time
			FROM %s
			WHERE pause_end IS NOT NULL
			GROUP BY time_entry_id
		) pause_stats ON te.id = pause_stats.time_entry_id`, timeEntryTable, projectTable, pauseTable)

	args := []interface{}{}
	whereClauses := []string{}

	if since != nil {
		whereClauses = append(whereClauses, "te.start_time >= ?")
		args = append(args, *since)
	}

	if category != nil {
		whereClauses = append(whereClauses, "te.category = ?")
		args = append(args, *category)
	}

	if len(whereClauses) > 0 {
		query += " WHERE " + whereClauses[0]
		for i := 1; i < len(whereClauses); i++ {
			query += " AND " + whereClauses[i]
		}
	}

	// sorting
	query += ` ` + orderClause

	// limit
	query += ` LIMIT ?`
	args = append(args, limit)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTimeEntriesWithPauses(rows)
}
