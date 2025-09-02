package migrations

import (
	"context"
	"database/sql"
)

func init() {
	up := func(ctx context.Context, tx *sql.Tx) error {
		// Create time_entries table
		query := `
		CREATE TABLE IF NOT EXISTS time_entries (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			project_id INTEGER NOT NULL,
			start_time DATETIME NOT NULL,
			end_time DATETIME,
			duration INTEGER, -- Duration in seconds
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
		);`

		_, err := tx.ExecContext(ctx, query)
		if err != nil {
			return err
		}

		// Create index on project_id for better performance
		indexQuery := `CREATE INDEX IF NOT EXISTS idx_time_entries_project_id ON time_entries(project_id);`
		_, err = tx.ExecContext(ctx, indexQuery)
		if err != nil {
			return err
		}

		// Create index on start_time for sorting
		indexQuery2 := `CREATE INDEX IF NOT EXISTS idx_time_entries_start_time ON time_entries(start_time);`
		_, err = tx.ExecContext(ctx, indexQuery2)
		return err
	}

	down := func(ctx context.Context, tx *sql.Tx) error {
		// Drop indexes first
		_, err := tx.ExecContext(ctx, `DROP INDEX IF EXISTS idx_time_entries_start_time;`)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, `DROP INDEX IF EXISTS idx_time_entries_project_id;`)
		if err != nil {
			return err
		}

		// Drop time_entries table
		query := `DROP TABLE IF EXISTS time_entries;`
		_, err = tx.ExecContext(ctx, query)
		return err
	}

	// Register the migration
	AddMigration("002_create_time_entries_table", up, down)
}
