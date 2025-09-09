package migrations

import (
	"context"
	"database/sql"
)

func init() {
	up := func(ctx context.Context, tx *sql.Tx) error {
		// Add category column to time_entries table
		query := `ALTER TABLE time_entries ADD COLUMN category VARCHAR(50);`
		_, err := tx.ExecContext(ctx, query)
		if err != nil {
			return err
		}

		// Create index on category for better performance
		indexQuery := `CREATE INDEX IF NOT EXISTS idx_time_entries_category ON time_entries(category);`
		_, err = tx.ExecContext(ctx, indexQuery)
		return err
	}

	down := func(ctx context.Context, tx *sql.Tx) error {
		// Drop the category index first
		_, err := tx.ExecContext(ctx, `DROP INDEX IF EXISTS idx_time_entries_category;`)
		if err != nil {
			return err
		}

		// SQLite doesn't support DROP COLUMN directly, so we need to recreate the table
		// This is a simplified approach - in production you might want to use a more sophisticated migration
		query := `
		CREATE TABLE time_entries_new (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			project_id INTEGER NOT NULL,
			start_time DATETIME NOT NULL,
			end_time DATETIME,
			duration INTEGER, -- Duration in seconds
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
		);
		
		INSERT INTO time_entries_new (id, project_id, start_time, end_time, duration, created_at)
		SELECT id, project_id, start_time, end_time, duration, created_at FROM time_entries;
		
		DROP TABLE time_entries;
		ALTER TABLE time_entries_new RENAME TO time_entries;
		
		-- Recreate indexes
		CREATE INDEX IF NOT EXISTS idx_time_entries_project_id ON time_entries(project_id);
		CREATE INDEX IF NOT EXISTS idx_time_entries_start_time ON time_entries(start_time);
		`
		_, err = tx.ExecContext(ctx, query)
		return err
	}

	// Register the migration
	AddMigration("004_add_category_to_time_entries", up, down)
}
