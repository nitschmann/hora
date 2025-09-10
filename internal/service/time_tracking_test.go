package service

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/nitschmann/hora/internal/model"
	"github.com/nitschmann/hora/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProjectRepo struct {
	mock.Mock
}

func (m *MockProjectRepo) Create(ctx context.Context, name string) (*model.Project, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(*model.Project), args.Error(1)
}

func (m *MockProjectRepo) GetByID(ctx context.Context, id int) (*model.Project, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Project), args.Error(1)
}

func (m *MockProjectRepo) GetByName(ctx context.Context, name string) (*model.Project, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(*model.Project), args.Error(1)
}

func (m *MockProjectRepo) GetOrCreate(ctx context.Context, name string) (*model.Project, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(*model.Project), args.Error(1)
}

func (m *MockProjectRepo) GetAll(ctx context.Context) ([]model.Project, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Project), args.Error(1)
}

func (m *MockProjectRepo) Delete(ctx context.Context, name string) error {
	args := m.Called(ctx, name)
	return args.Error(0)
}

func (m *MockProjectRepo) DeleteByID(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProjectRepo) GetByIDOrName(ctx context.Context, idOrName string) (*model.Project, error) {
	args := m.Called(ctx, idOrName)
	return args.Get(0).(*model.Project), args.Error(1)
}

type MockTimeEntryRepo struct {
	mock.Mock
}

func (m *MockTimeEntryRepo) Create(ctx context.Context, projectID int, startTime time.Time, category *string) (*model.TimeEntry, error) {
	args := m.Called(ctx, projectID, startTime, category)
	return args.Get(0).(*model.TimeEntry), args.Error(1)
}

func (m *MockTimeEntryRepo) GetByID(ctx context.Context, id int) (*model.TimeEntry, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.TimeEntry), args.Error(1)
}

func (m *MockTimeEntryRepo) GetActive(ctx context.Context) (*model.TimeEntry, error) {
	args := m.Called(ctx)
	return args.Get(0).(*model.TimeEntry), args.Error(1)
}

func (m *MockTimeEntryRepo) UpdateEndTime(ctx context.Context, id int, endTime time.Time, duration time.Duration) error {
	args := m.Called(ctx, id, endTime, duration)
	return args.Error(0)
}

func (m *MockTimeEntryRepo) StopAllActive(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockTimeEntryRepo) GetByProject(ctx context.Context, projectID int, limit int, sortOrder string) ([]model.TimeEntry, error) {
	args := m.Called(ctx, projectID, limit, sortOrder)
	return args.Get(0).([]model.TimeEntry), args.Error(1)
}

func (m *MockTimeEntryRepo) GetByProjectName(ctx context.Context, projectName string, limit int, sortOrder string) ([]model.TimeEntry, error) {
	args := m.Called(ctx, projectName, limit, sortOrder)
	return args.Get(0).([]model.TimeEntry), args.Error(1)
}

func (m *MockTimeEntryRepo) GetByProjectWithPauses(ctx context.Context, projectID int, limit int, sortOrder string, since *time.Time) ([]repository.TimeEntryWithPauses, error) {
	args := m.Called(ctx, projectID, limit, sortOrder, since)
	return args.Get(0).([]repository.TimeEntryWithPauses), args.Error(1)
}

func (m *MockTimeEntryRepo) GetByProjectNameWithPauses(ctx context.Context, projectName string, limit int, sortOrder string, since *time.Time) ([]repository.TimeEntryWithPauses, error) {
	args := m.Called(ctx, projectName, limit, sortOrder, since)
	return args.Get(0).([]repository.TimeEntryWithPauses), args.Error(1)
}

func (m *MockTimeEntryRepo) GetAll(ctx context.Context, limit int) ([]model.TimeEntry, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]model.TimeEntry), args.Error(1)
}

func (m *MockTimeEntryRepo) GetAllWithPauses(ctx context.Context, limit int, sortOrder string, since *time.Time) ([]repository.TimeEntryWithPauses, error) {
	args := m.Called(ctx, limit, sortOrder, since)
	return args.Get(0).([]repository.TimeEntryWithPauses), args.Error(1)
}

func (m *MockTimeEntryRepo) GetAllWithPausesByCategory(ctx context.Context, limit int, sortOrder string, since *time.Time, category *string) ([]repository.TimeEntryWithPauses, error) {
	args := m.Called(ctx, limit, sortOrder, since, category)
	return args.Get(0).([]repository.TimeEntryWithPauses), args.Error(1)
}

func (m *MockTimeEntryRepo) GetTotalTimeByProject(ctx context.Context, projectID int, since *time.Time) (time.Duration, error) {
	args := m.Called(ctx, projectID, since)
	return args.Get(0).(time.Duration), args.Error(1)
}

func (m *MockTimeEntryRepo) GetTotalTimeByProjectName(ctx context.Context, projectName string, since *time.Time) (time.Duration, error) {
	args := m.Called(ctx, projectName, since)
	return args.Get(0).(time.Duration), args.Error(1)
}

func (m *MockTimeEntryRepo) DeleteByProject(ctx context.Context, projectID int) error {
	args := m.Called(ctx, projectID)
	return args.Error(0)
}

func (m *MockTimeEntryRepo) DeleteAll(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockTimeEntryRepo) GetByProjectIDOrName(ctx context.Context, projectIDOrName string, limit int, sortOrder string) ([]model.TimeEntry, error) {
	args := m.Called(ctx, projectIDOrName, limit, sortOrder)
	return args.Get(0).([]model.TimeEntry), args.Error(1)
}

func (m *MockTimeEntryRepo) GetByProjectIDOrNameWithPauses(ctx context.Context, projectIDOrName string, limit int, sortOrder string, since *time.Time) ([]repository.TimeEntryWithPauses, error) {
	args := m.Called(ctx, projectIDOrName, limit, sortOrder, since)
	return args.Get(0).([]repository.TimeEntryWithPauses), args.Error(1)
}

func (m *MockTimeEntryRepo) GetTotalTimeByProjectIDOrName(ctx context.Context, projectIDOrName string, since *time.Time) (time.Duration, error) {
	args := m.Called(ctx, projectIDOrName, since)
	return args.Get(0).(time.Duration), args.Error(1)
}

type MockPauseRepo struct {
	mock.Mock
}

func (m *MockPauseRepo) Create(ctx context.Context, timeEntryID int, pauseStart time.Time) (*model.Pause, error) {
	args := m.Called(ctx, timeEntryID, pauseStart)
	return args.Get(0).(*model.Pause), args.Error(1)
}

func (m *MockPauseRepo) GetByID(ctx context.Context, id int) (*model.Pause, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Pause), args.Error(1)
}

func (m *MockPauseRepo) GetActivePause(ctx context.Context, timeEntryID int) (*model.Pause, error) {
	args := m.Called(ctx, timeEntryID)
	return args.Get(0).(*model.Pause), args.Error(1)
}

func (m *MockPauseRepo) EndPause(ctx context.Context, id int, pauseEnd time.Time, duration time.Duration) error {
	args := m.Called(ctx, id, pauseEnd, duration)
	return args.Error(0)
}

func (m *MockPauseRepo) GetByTimeEntry(ctx context.Context, timeEntryID int) ([]model.Pause, error) {
	args := m.Called(ctx, timeEntryID)
	return args.Get(0).([]model.Pause), args.Error(1)
}

func (m *MockPauseRepo) DeleteByTimeEntry(ctx context.Context, timeEntryID int) error {
	args := m.Called(ctx, timeEntryID)
	return args.Error(0)
}

func (m *MockPauseRepo) DeleteAll(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestTimeTracking_StartTracking(t *testing.T) {
	ctx := context.Background()
	mockProjectRepo := &MockProjectRepo{}
	mockTimeEntryRepo := &MockTimeEntryRepo{}
	mockPauseRepo := &MockPauseRepo{}

	service := &timeTracking{
		projectRepo:   mockProjectRepo,
		timeEntryRepo: mockTimeEntryRepo,
		pauseRepo:     mockPauseRepo,
	}

	project := &model.Project{ID: 1, Name: "Test Project"}
	timeEntry := &model.TimeEntry{ID: 1, ProjectID: 1, StartTime: time.Now()}

	mockProjectRepo.On("GetOrCreate", ctx, "Test Project").Return(project, nil)
	mockTimeEntryRepo.On("GetActive", ctx).Return((*model.TimeEntry)(nil), nil)
	mockTimeEntryRepo.On("Create", ctx, 1, mock.AnythingOfType("time.Time"), (*string)(nil)).Return(timeEntry, nil)

	err := service.StartTracking(ctx, "Test Project", false, nil)

	assert.NoError(t, err)
	mockProjectRepo.AssertExpectations(t)
	mockTimeEntryRepo.AssertExpectations(t)
}

func TestTimeTracking_StartTracking_Force(t *testing.T) {
	ctx := context.Background()
	mockProjectRepo := &MockProjectRepo{}
	mockTimeEntryRepo := &MockTimeEntryRepo{}
	mockPauseRepo := &MockPauseRepo{}

	service := &timeTracking{
		projectRepo:   mockProjectRepo,
		timeEntryRepo: mockTimeEntryRepo,
		pauseRepo:     mockPauseRepo,
	}

	project := &model.Project{ID: 1, Name: "Test Project"}
	newEntry := &model.TimeEntry{ID: 2, ProjectID: 1, StartTime: time.Now()}

	mockProjectRepo.On("GetOrCreate", ctx, "Test Project").Return(project, nil)
	mockTimeEntryRepo.On("StopAllActive", ctx).Return(nil)
	mockTimeEntryRepo.On("Create", ctx, 1, mock.AnythingOfType("time.Time"), (*string)(nil)).Return(newEntry, nil)

	err := service.StartTracking(ctx, "Test Project", true, nil)

	assert.NoError(t, err)
	mockProjectRepo.AssertExpectations(t)
	mockTimeEntryRepo.AssertExpectations(t)
}

func TestTimeTracking_StartTracking_AlreadyActive(t *testing.T) {
	ctx := context.Background()
	mockProjectRepo := &MockProjectRepo{}
	mockTimeEntryRepo := &MockTimeEntryRepo{}
	mockPauseRepo := &MockPauseRepo{}

	service := &timeTracking{
		projectRepo:   mockProjectRepo,
		timeEntryRepo: mockTimeEntryRepo,
		pauseRepo:     mockPauseRepo,
	}

	project := &model.Project{ID: 1, Name: "Existing Project"}
	activeEntry := &model.TimeEntry{ID: 1, ProjectID: 1, StartTime: time.Now(), Project: project}

	mockTimeEntryRepo.On("GetActive", ctx).Return(activeEntry, nil)

	err := service.StartTracking(ctx, "Test Project", false, nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already active")
	mockTimeEntryRepo.AssertExpectations(t)
}

func TestTimeTracking_StopTracking(t *testing.T) {
	ctx := context.Background()
	mockProjectRepo := &MockProjectRepo{}
	mockTimeEntryRepo := &MockTimeEntryRepo{}
	mockPauseRepo := &MockPauseRepo{}

	service := &timeTracking{
		projectRepo:   mockProjectRepo,
		timeEntryRepo: mockTimeEntryRepo,
		pauseRepo:     mockPauseRepo,
	}

	activeEntry := &model.TimeEntry{ID: 1, ProjectID: 1, StartTime: time.Now().Add(-time.Hour)}

	mockTimeEntryRepo.On("GetActive", ctx).Return(activeEntry, nil)
	mockPauseRepo.On("GetActivePause", ctx, 1).Return((*model.Pause)(nil), sql.ErrNoRows)
	mockPauseRepo.On("GetByTimeEntry", ctx, 1).Return([]model.Pause{}, nil)
	mockTimeEntryRepo.On("UpdateEndTime", ctx, 1, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Duration")).Return(nil)
	mockTimeEntryRepo.On("GetByID", ctx, 1).Return(activeEntry, nil)

	result, err := service.StopTracking(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockTimeEntryRepo.AssertExpectations(t)
}

func TestTimeTracking_StopTracking_NoActive(t *testing.T) {
	ctx := context.Background()
	mockProjectRepo := &MockProjectRepo{}
	mockTimeEntryRepo := &MockTimeEntryRepo{}
	mockPauseRepo := &MockPauseRepo{}

	service := &timeTracking{
		projectRepo:   mockProjectRepo,
		timeEntryRepo: mockTimeEntryRepo,
		pauseRepo:     mockPauseRepo,
	}

	mockTimeEntryRepo.On("GetActive", ctx).Return((*model.TimeEntry)(nil), sql.ErrNoRows)

	result, err := service.StopTracking(ctx)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no active")
	mockTimeEntryRepo.AssertExpectations(t)
}

func TestTimeTracking_GetActiveEntry(t *testing.T) {
	ctx := context.Background()
	mockProjectRepo := &MockProjectRepo{}
	mockTimeEntryRepo := &MockTimeEntryRepo{}
	mockPauseRepo := &MockPauseRepo{}

	service := &timeTracking{
		projectRepo:   mockProjectRepo,
		timeEntryRepo: mockTimeEntryRepo,
		pauseRepo:     mockPauseRepo,
	}

	activeEntry := &model.TimeEntry{ID: 1, ProjectID: 1, StartTime: time.Now()}

	mockTimeEntryRepo.On("GetActive", ctx).Return(activeEntry, nil)

	result, err := service.GetActiveEntry(ctx)

	assert.NoError(t, err)
	assert.Equal(t, activeEntry, result)
	mockTimeEntryRepo.AssertExpectations(t)
}

func TestTimeTracking_GetProjects(t *testing.T) {
	ctx := context.Background()
	mockProjectRepo := &MockProjectRepo{}
	mockTimeEntryRepo := &MockTimeEntryRepo{}
	mockPauseRepo := &MockPauseRepo{}

	service := &timeTracking{
		projectRepo:   mockProjectRepo,
		timeEntryRepo: mockTimeEntryRepo,
		pauseRepo:     mockPauseRepo,
	}

	projects := []model.Project{
		{ID: 1, Name: "Project 1"},
		{ID: 2, Name: "Project 2"},
	}

	mockProjectRepo.On("GetAll", ctx).Return(projects, nil)

	result, err := service.GetProjects(ctx)

	assert.NoError(t, err)
	assert.Equal(t, projects, result)
	mockProjectRepo.AssertExpectations(t)
}

func TestTimeTracking_GetOrCreateProject(t *testing.T) {
	ctx := context.Background()
	mockProjectRepo := &MockProjectRepo{}
	mockTimeEntryRepo := &MockTimeEntryRepo{}
	mockPauseRepo := &MockPauseRepo{}

	service := &timeTracking{
		projectRepo:   mockProjectRepo,
		timeEntryRepo: mockTimeEntryRepo,
		pauseRepo:     mockPauseRepo,
	}

	project := &model.Project{ID: 1, Name: "Test Project"}

	mockProjectRepo.On("GetOrCreate", ctx, "Test Project").Return(project, nil)

	result, err := service.GetOrCreateProject(ctx, "Test Project")

	assert.NoError(t, err)
	assert.Equal(t, project, result)
	mockProjectRepo.AssertExpectations(t)
}

func TestTimeTracking_ClearAllData(t *testing.T) {
	ctx := context.Background()
	mockProjectRepo := &MockProjectRepo{}
	mockTimeEntryRepo := &MockTimeEntryRepo{}
	mockPauseRepo := &MockPauseRepo{}

	service := &timeTracking{
		projectRepo:   mockProjectRepo,
		timeEntryRepo: mockTimeEntryRepo,
		pauseRepo:     mockPauseRepo,
	}

	mockTimeEntryRepo.On("DeleteAll", ctx).Return(nil)
	mockPauseRepo.On("DeleteAll", ctx).Return(nil)

	err := service.ClearAllData(ctx)

	assert.NoError(t, err)
	mockTimeEntryRepo.AssertExpectations(t)
	mockPauseRepo.AssertExpectations(t)
}

func TestTimeTracking_FormatDuration(t *testing.T) {
	service := &timeTracking{}

	tests := []struct {
		duration time.Duration
		expected string
	}{
		{time.Hour, "01:00:00"},
		{time.Hour + 30*time.Minute, "01:30:00"},
		{45*time.Minute + 30*time.Second, "00:45:30"},
		{30 * time.Second, "00:00:30"},
	}

	for _, test := range tests {
		result := service.FormatDuration(test.duration)
		assert.Equal(t, test.expected, result)
	}
}
