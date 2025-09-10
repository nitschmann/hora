package repository

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/nitschmann/hora/internal/config"
	"github.com/nitschmann/hora/internal/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	// Create a temporary directory for the test database
	tempDir, err := os.MkdirTemp("", "hora-test-*")
	require.NoError(t, err)

	// Clean up the temp directory after the test
	t.Cleanup(func() {
		os.RemoveAll(tempDir)
	})

	conf := &config.Config{
		DatabaseDir: tempDir,
	}

	conn, err := database.NewConnection(conf)
	require.NoError(t, err)

	return conn.GetDB()
}

func TestProjectIntegration(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewProject(db)
	ctx := context.Background()

	// Test Create
	project, err := repo.Create(ctx, "Test Project Basic")
	require.NoError(t, err)
	assert.Equal(t, "Test Project Basic", project.Name)
	assert.NotZero(t, project.ID)

	// Test GetByID
	found, err := repo.GetByID(ctx, project.ID)
	require.NoError(t, err)
	assert.Equal(t, project.ID, found.ID)
	assert.Equal(t, "Test Project Basic", found.Name)

	// Test GetByName
	foundByName, err := repo.GetByName(ctx, "Test Project Basic")
	require.NoError(t, err)
	assert.Equal(t, project.ID, foundByName.ID)

	// Test GetAll
	projects, err := repo.GetAll(ctx)
	require.NoError(t, err)
	assert.Len(t, projects, 1)
	assert.Equal(t, "Test Project Basic", projects[0].Name)

	// Test Delete
	err = repo.Delete(ctx, "Test Project Basic")
	require.NoError(t, err)

	// Verify deletion
	_, err = repo.GetByID(ctx, project.ID)
	assert.Error(t, err)
}

func TestTimeEntryIntegration(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	projectRepo := NewProject(db)
	timeEntryRepo := NewTimeEntry(db)
	ctx := context.Background()

	// Create a project first
	project, err := projectRepo.Create(ctx, "Test Project Time")
	require.NoError(t, err)

	// Test Create time entry
	timeEntry, err := timeEntryRepo.Create(ctx, project.ID, time.Now(), nil)
	require.NoError(t, err)
	assert.Equal(t, project.ID, timeEntry.ProjectID)
	assert.NotZero(t, timeEntry.ID)

	// Test GetByID
	found, err := timeEntryRepo.GetByID(ctx, timeEntry.ID)
	require.NoError(t, err)
	assert.Equal(t, timeEntry.ID, found.ID)

	// Test GetActive
	active, err := timeEntryRepo.GetActive(ctx)
	require.NoError(t, err)
	assert.Equal(t, timeEntry.ID, active.ID)

	// Test UpdateEndTime
	err = timeEntryRepo.UpdateEndTime(ctx, timeEntry.ID, time.Now(), time.Hour)
	require.NoError(t, err)

	// Test GetByProject
	entries, err := timeEntryRepo.GetByProject(ctx, project.ID, 10, "desc")
	require.NoError(t, err)
	assert.Len(t, entries, 1)
}

func TestPauseIntegration(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	projectRepo := NewProject(db)
	timeEntryRepo := NewTimeEntry(db)
	pauseRepo := NewPause(db)
	ctx := context.Background()

	// Create a project and time entry first
	project, err := projectRepo.Create(ctx, "Test Project Pause")
	require.NoError(t, err)

	timeEntry, err := timeEntryRepo.Create(ctx, project.ID, time.Now(), nil)
	require.NoError(t, err)

	// Test Create pause
	pause, err := pauseRepo.Create(ctx, timeEntry.ID, time.Now())
	require.NoError(t, err)
	assert.Equal(t, timeEntry.ID, pause.TimeEntryID)
	assert.NotZero(t, pause.ID)

	// Test GetByID
	found, err := pauseRepo.GetByID(ctx, pause.ID)
	require.NoError(t, err)
	assert.Equal(t, pause.ID, found.ID)

	// Test GetActivePause
	active, err := pauseRepo.GetActivePause(ctx, timeEntry.ID)
	require.NoError(t, err)
	assert.Equal(t, pause.ID, active.ID)

	// Test EndPause
	err = pauseRepo.EndPause(ctx, pause.ID, time.Now(), time.Hour)
	require.NoError(t, err)

	// Test GetByTimeEntry
	pauses, err := pauseRepo.GetByTimeEntry(ctx, timeEntry.ID)
	require.NoError(t, err)
	assert.Len(t, pauses, 1)
}
