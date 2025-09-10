package repository

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/doug-martin/goqu/v9"

	"github.com/nitschmann/hora/internal/model"
)

const projectTable = "projects"

// Project defines the interface for project data operations
type Project interface {
	// Create creates a new project
	Create(ctx context.Context, name string) (*model.Project, error)
	// GetByID retrieves a project by its ID
	GetByID(ctx context.Context, id int) (*model.Project, error)
	// GetByName retrieves a project by its name
	GetByName(ctx context.Context, name string) (*model.Project, error)
	// GetOrCreate retrieves a project by name, or creates it if it doesn't exist
	GetOrCreate(ctx context.Context, name string) (*model.Project, error)
	// GetAll retrieves all projects with their last tracked time
	GetAll(ctx context.Context) ([]model.Project, error)
	// Delete deletes a project by name
	Delete(ctx context.Context, name string) error
	// DeleteByID deletes a project by ID
	DeleteByID(ctx context.Context, id int) error
	// GetByIDOrName retrieves a project by ID (if numeric) or name
	GetByIDOrName(ctx context.Context, idOrName string) (*model.Project, error)
}

type project struct {
	db *sql.DB
}

// NewProject creates a new project repository
func NewProject(db *sql.DB) Project {
	return &project{db: db}
}

// Create creates a new project
func (r *project) Create(ctx context.Context, name string) (*model.Project, error) {
	query, args, err := goqu.Insert(projectTable).Rows(goqu.Record{
		"name": name,
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

	// Get the created project
	return r.GetByID(ctx, int(id))
}

// GetByID retrieves a project by its ID
func (r *project) GetByID(ctx context.Context, id int) (*model.Project, error) {
	query, args, err := goqu.From(projectTable).
		Select("id", "name", "created_at").
		Where(goqu.C("id").Eq(id)).
		ToSQL()
	if err != nil {
		return nil, err
	}

	var project model.Project
	err = r.db.QueryRowContext(ctx, query, args...).Scan(
		&project.ID,
		&project.Name,
		&project.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &project, nil
}

// GetByName retrieves a project by its name
func (r *project) GetByName(ctx context.Context, name string) (*model.Project, error) {
	query, args, err := goqu.From(projectTable).
		Select("id", "name", "created_at").
		Where(goqu.C("name").Eq(name)).
		ToSQL()
	if err != nil {
		return nil, err
	}

	var project model.Project
	err = r.db.QueryRowContext(ctx, query, args...).Scan(
		&project.ID,
		&project.Name,
		&project.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &project, nil
}

// GetOrCreate retrieves a project by name, or creates it if it doesn't exist
func (r *project) GetOrCreate(ctx context.Context, name string) (*model.Project, error) {
	// Try to get existing project first
	project, err := r.GetByName(ctx, name)
	if err == nil {
		return project, nil
	}

	// If not found, create it
	return r.Create(ctx, name)
}

// GetAll retrieves all projects with their last tracked time
func (r *project) GetAll(ctx context.Context) ([]model.Project, error) {
	query, args, err := goqu.From(projectTable).
		LeftJoin(goqu.T("time_entries"), goqu.On(goqu.I("projects.id").Eq(goqu.I("time_entries.project_id")))).
		Select(
			goqu.I("projects.id"),
			goqu.I("projects.name"),
			goqu.I("projects.created_at"),
			goqu.MAX(goqu.I("time_entries.end_time")).As("last_tracked_at"),
		).
		GroupBy(goqu.I("projects.id"), goqu.I("projects.name"), goqu.I("projects.created_at")).
		Order(goqu.I("projects.name").Asc()).
		ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []model.Project
	for rows.Next() {
		var project model.Project
		var lastTrackedAtStr *string

		err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.CreatedAt,
			&lastTrackedAtStr,
		)
		if err != nil {
			return nil, err
		}

		// Parse last tracked time if present
		if lastTrackedAtStr != nil {
			// Try different time formats
			formats := []string{
				time.RFC3339Nano,
				time.RFC3339,
				"2006-01-02 15:04:05.999999999-07:00",
				"2006-01-02 15:04:05",
			}

			for _, format := range formats {
				if lastTrackedAt, err := time.Parse(format, *lastTrackedAtStr); err == nil {
					project.LastTrackedAt = &lastTrackedAt
					break
				}
			}
		}

		projects = append(projects, project)
	}

	return projects, rows.Err()
}

// Delete deletes a project by name
func (r *project) Delete(ctx context.Context, name string) error {
	query, args, err := goqu.Delete(projectTable).
		Where(goqu.C("name").Eq(name)).
		ToSQL()
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

// DeleteByID deletes a project by ID
func (r *project) DeleteByID(ctx context.Context, id int) error {
	query, args, err := goqu.Delete(projectTable).
		Where(goqu.C("id").Eq(id)).
		ToSQL()
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

// GetByIDOrName retrieves a project by ID (if numeric) or name
func (r *project) GetByIDOrName(ctx context.Context, idOrName string) (*model.Project, error) {
	if id, err := strconv.Atoi(idOrName); err == nil {
		return r.GetByID(ctx, id)
	}

	return r.GetByName(ctx, idOrName)
}
