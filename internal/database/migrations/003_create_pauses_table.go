package migrations

import (
	"context"
	"database/sql"
)

func init() {
	up := func(ctx context.Context, tx *sql.Tx) error {
		// Create pauses table
		query := `
		CREATE TABLE IF NOT EXISTS pauses (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			time_entry_id INTEGER NOT NULL,
			pause_start DATETIME NOT NULL,
			pause_end DATETIME,
			duration INTEGER, -- Duration in seconds
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (time_entry_id) REFERENCES time_entries(id) ON DELETE CASCADE
		);
		CREATE INDEX IF NOT EXISTS idx_pauses_time_entry_id ON pauses(time_entry_id);
		CREATE INDEX IF NOT EXISTS idx_pauses_pause_start ON pauses(pause_start);
		`

		_, err := tx.ExecContext(ctx, query)
		return err
	}

	down := func(ctx context.Context, tx *sql.Tx) error {
		query := `DROP TABLE IF EXISTS pauses;`
		_, err := tx.ExecContext(ctx, query)
		return err
	}

	AddMigration("003_create_pauses_table", up, down)
}
