package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nitschmann/hora/internal/model"
)

// Pause defines the interface for pause data operations
type Pause interface {
	// Create creates a new pause
	Create(ctx context.Context, timeEntryID int, pauseStart time.Time) (*model.Pause, error)

	// GetByID retrieves a pause by its ID
	GetByID(ctx context.Context, id int) (*model.Pause, error)

	// GetActivePause retrieves the currently active pause for a time entry
	GetActivePause(ctx context.Context, timeEntryID int) (*model.Pause, error)

	// EndPause ends a pause by setting its end time and duration
	EndPause(ctx context.Context, id int, pauseEnd time.Time, duration time.Duration) error

	// GetByTimeEntry retrieves all pauses for a specific time entry
	GetByTimeEntry(ctx context.Context, timeEntryID int) ([]model.Pause, error)

	// DeleteByTimeEntry deletes all pauses for a specific time entry
	DeleteByTimeEntry(ctx context.Context, timeEntryID int) error

	// DeleteAll deletes all pauses
	DeleteAll(ctx context.Context) error
}

// pause implements Pause using SQLite
type pause struct {
	db *sql.DB
}

// NewPause creates a new pause repository
func NewPause(db *sql.DB) Pause {
	return &pause{db: db}
}

// Create creates a new pause
func (r *pause) Create(ctx context.Context, timeEntryID int, pauseStart time.Time) (*model.Pause, error) {
	query, args, err := goqu.Insert("pauses").Rows(goqu.Record{
		"time_entry_id": timeEntryID,
		"pause_start":   pauseStart,
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

	// Get the created pause
	return r.GetByID(ctx, int(id))
}

// GetByID retrieves a pause by its ID
func (r *pause) GetByID(ctx context.Context, id int) (*model.Pause, error) {
	query, args, err := goqu.From("pauses").
		Select(goqu.Star()).
		Where(goqu.C("id").Eq(id)).
		ToSQL()
	if err != nil {
		return nil, err
	}

	var pause model.Pause
	var pauseEnd *time.Time
	var duration *int64

	err = r.db.QueryRowContext(ctx, query, args...).Scan(
		&pause.ID,
		&pause.TimeEntryID,
		&pause.PauseStart,
		&pauseEnd,
		&duration,
		&pause.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	pause.PauseEnd = pauseEnd
	if duration != nil {
		d := time.Duration(*duration) * time.Second
		pause.Duration = &d
	}

	return &pause, nil
}

// GetActivePause retrieves the currently active pause for a time entry
func (r *pause) GetActivePause(ctx context.Context, timeEntryID int) (*model.Pause, error) {
	query, args, err := goqu.From("pauses").
		Select(goqu.Star()).
		Where(goqu.C("time_entry_id").Eq(timeEntryID), goqu.C("pause_end").IsNull()).
		Order(goqu.C("pause_start").Desc()).
		Limit(1).
		ToSQL()
	if err != nil {
		return nil, err
	}

	var pause model.Pause
	var pauseEnd *time.Time
	var duration *int64

	err = r.db.QueryRowContext(ctx, query, args...).Scan(
		&pause.ID,
		&pause.TimeEntryID,
		&pause.PauseStart,
		&pauseEnd,
		&duration,
		&pause.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	pause.PauseEnd = pauseEnd
	if duration != nil {
		d := time.Duration(*duration) * time.Second
		pause.Duration = &d
	}

	return &pause, nil
}

// EndPause ends a pause by setting its end time and duration
func (r *pause) EndPause(ctx context.Context, id int, pauseEnd time.Time, duration time.Duration) error {
	durationSeconds := int64(duration.Seconds())
	query, args, err := goqu.Update("pauses").
		Set(goqu.Record{
			"pause_end": pauseEnd,
			"duration":  durationSeconds,
		}).
		Where(goqu.C("id").Eq(id)).
		ToSQL()
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

// GetByTimeEntry retrieves all pauses for a specific time entry
func (r *pause) GetByTimeEntry(ctx context.Context, timeEntryID int) ([]model.Pause, error) {
	query, args, err := goqu.From("pauses").
		Select(goqu.Star()).
		Where(goqu.C("time_entry_id").Eq(timeEntryID)).
		Order(goqu.C("pause_start").Asc()).
		ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanPauses(rows)
}

// DeleteByTimeEntry deletes all pauses for a specific time entry
func (r *pause) DeleteByTimeEntry(ctx context.Context, timeEntryID int) error {
	query, args, err := goqu.Delete("pauses").
		Where(goqu.C("time_entry_id").Eq(timeEntryID)).
		ToSQL()
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

// DeleteAll deletes all pauses
func (r *pause) DeleteAll(ctx context.Context) error {
	query, args, err := goqu.Delete("pauses").ToSQL()
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

// scanPauses is a helper method to scan pauses from rows
func (r *pause) scanPauses(rows *sql.Rows) ([]model.Pause, error) {
	var pauses []model.Pause

	for rows.Next() {
		var pause model.Pause
		var pauseEnd *time.Time
		var duration *int64

		err := rows.Scan(
			&pause.ID,
			&pause.TimeEntryID,
			&pause.PauseStart,
			&pauseEnd,
			&duration,
			&pause.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		pause.PauseEnd = pauseEnd
		if duration != nil {
			d := time.Duration(*duration) * time.Second
			pause.Duration = &d
		}

		pauses = append(pauses, pause)
	}

	return pauses, rows.Err()
}
