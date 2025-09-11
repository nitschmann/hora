package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nitschmann/hora/internal/config"
	"github.com/nitschmann/hora/internal/database/migrations"
)

const databaseFileName = "hora.db"

type Connection struct {
	db *sql.DB
}

// NewConnection creates a new database connection
func NewConnection(conf *config.Config) (*Connection, error) {
	dbPath, err := setUpSQLiteDatabaseFile(conf.DatabaseDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get database path: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	ctx := context.Background()
	if err := migrations.RunMigrations(ctx, db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return &Connection{db: db}, nil
}

// GetDB returns the underlying database connection
func (c *Connection) GetDB() *sql.DB {
	return c.db
}

func (c *Connection) Close() error {
	return c.db.Close()
}

func setUpSQLiteDatabaseFile(databaseDir string) (string, error) {
	if err := os.MkdirAll(databaseDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create database directory %s: %w", databaseDir, err)
	}

	_ = os.Remove(filepath.Join(databaseDir, databaseFileName+"-wal"))
	_ = os.Remove(filepath.Join(databaseDir, databaseFileName+"-shm"))

	return filepath.Join(databaseDir, databaseFileName), nil
}
