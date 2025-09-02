package migrations

import (
	"context"
	"database/sql"
)

func init() {
	up := func(ctx context.Context, tx *sql.Tx) error {
		// Create projects table
		query := `
		CREATE TABLE IF NOT EXISTS projects (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`

		_, err := tx.ExecContext(ctx, query)
		return err
	}

	down := func(ctx context.Context, tx *sql.Tx) error {
		// Drop projects table
		query := `DROP TABLE IF EXISTS projects;`
		_, err := tx.ExecContext(ctx, query)
		return err
	}

	// Register the migration
	AddMigration("001_create_projects_table", up, down)
}
