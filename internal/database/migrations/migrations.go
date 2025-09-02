package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
)

// Migration represents a database migration
type Migration struct {
	Name string
	Up   func(ctx context.Context, tx *sql.Tx) error
	Down func(ctx context.Context, tx *sql.Tx) error
}

var migrations = make(map[string]*Migration)

// AddMigration registers a new migration
func AddMigration(name string, up, down func(ctx context.Context, tx *sql.Tx) error) {
	migrations[name] = &Migration{
		Name: name,
		Up:   up,
		Down: down,
	}
}

// RunMigrations runs all pending migrations
func RunMigrations(ctx context.Context, db *sql.DB) error {
	// Create migrations table if it doesn't exist
	if err := createMigrationsTable(ctx, db); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get applied migrations
	applied, err := getAppliedMigrations(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Get all available migrations
	var available []string
	for name := range migrations {
		available = append(available, name)
	}
	sort.Strings(available)

	// Run pending migrations
	for _, name := range available {
		if !applied[name] {
			if err := runMigration(ctx, db, name, migrations[name]); err != nil {
				return fmt.Errorf("failed to run migration %s: %w", name, err)
			}
		}
	}

	return nil
}

// createMigrationsTable creates the migrations tracking table
func createMigrationsTable(ctx context.Context, db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS migrations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.ExecContext(ctx, query)
	return err
}

// getAppliedMigrations returns a map of applied migration names
func getAppliedMigrations(ctx context.Context, db *sql.DB) (map[string]bool, error) {
	query := `SELECT name FROM migrations;`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		applied[name] = true
	}

	return applied, rows.Err()
}

// runMigration runs a single migration
func runMigration(ctx context.Context, db *sql.DB, name string, migration *Migration) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Run the migration
	if err := migration.Up(ctx, tx); err != nil {
		return err
	}

	// Record the migration as applied
	query := `INSERT INTO migrations (name) VALUES (?);`
	if _, err := tx.ExecContext(ctx, query, name); err != nil {
		return err
	}

	return tx.Commit()
}
