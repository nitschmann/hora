package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func resetViper() {
	viper.Reset()
}

func TestLoad_WithValidConfigFile(t *testing.T) {
	defer resetViper()
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	configContent := `database_dir: "/tmp/test"
debug: true
list_limit: 100
list_order: "asc"
use_background_tracker: false
web_ui_port: 9090
background_tracker_auto_stop: true
background_tracker_auto_stop_after: 45`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	cfg, usedPath, err := Load(configPath)
	require.NoError(t, err)
	assert.Equal(t, configPath, usedPath)
	assert.Equal(t, "/tmp/test", cfg.DatabaseDir)
	assert.True(t, cfg.Debug)
	assert.Equal(t, 100, cfg.ListLimit)
	assert.Equal(t, "asc", cfg.ListOrder)
	assert.False(t, cfg.UseBackgroundTracker)
	assert.Equal(t, 9090, cfg.WebUIPort)
	assert.True(t, cfg.BackgroundTrackerAutoStop)
	assert.Equal(t, 45, cfg.BackgroundTrackerAutoStopAfter)
}

func TestLoad_WithInvalidConfigFile(t *testing.T) {
	defer resetViper()
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	configContent := `database_dir: "/tmp/test"
debug: true
list_limit: 0
list_order: "invalid"
use_background_tracker: false`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	cfg, usedPath, err := Load(configPath)
	assert.Error(t, err)
	assert.Nil(t, cfg)
	assert.Equal(t, "", usedPath)
}

func TestLoad_WithNonExistentFile(t *testing.T) {
	defer resetViper()
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "nonexistent.yaml")

	cfg, usedPath, err := Load(configPath)
	assert.Error(t, err)
	assert.Nil(t, cfg)
	assert.Equal(t, "", usedPath)
}

func TestLoad_WithEmptyConfigFile(t *testing.T) {
	defer resetViper()
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	err := os.WriteFile(configPath, []byte(""), 0644)
	require.NoError(t, err)

	cfg, usedPath, err := Load(configPath)
	require.NoError(t, err)
	assert.Equal(t, configPath, usedPath)
	assert.NotEmpty(t, cfg.DatabaseDir)
	assert.False(t, cfg.Debug)
	assert.Equal(t, 50, cfg.ListLimit)
	assert.Equal(t, "desc", cfg.ListOrder)
	assert.True(t, cfg.UseBackgroundTracker)
	assert.Equal(t, 8080, cfg.WebUIPort)
	assert.False(t, cfg.BackgroundTrackerAutoStop)
	assert.Equal(t, 120, cfg.BackgroundTrackerAutoStopAfter)
}

func TestLoad_WithDefaultsOnly(t *testing.T) {
	defer resetViper()
	tempDir := t.TempDir()
	originalLookupPaths := LookupPaths
	LookupPaths = []string{tempDir}
	defer func() { LookupPaths = originalLookupPaths }()

	cfg, usedPath, err := Load("")
	require.NoError(t, err)
	assert.Equal(t, "", usedPath)
	assert.NotEmpty(t, cfg.DatabaseDir)
	assert.False(t, cfg.Debug)
	assert.Equal(t, 50, cfg.ListLimit)
	assert.Equal(t, "desc", cfg.ListOrder)
	assert.True(t, cfg.UseBackgroundTracker)
	assert.Equal(t, 8080, cfg.WebUIPort)
	assert.False(t, cfg.BackgroundTrackerAutoStop)
	assert.Equal(t, 120, cfg.BackgroundTrackerAutoStopAfter)
}

func TestCreateDefault_WithEmptyDirectory(t *testing.T) {
	defer resetViper()
	tempDir := t.TempDir()

	configPath, err := CreateDefault(tempDir, false)
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(tempDir, "config.yaml"), configPath)

	_, err = os.Stat(configPath)
	require.NoError(t, err)

	cfg, _, err := Load(configPath)
	require.NoError(t, err)
	assert.NotEmpty(t, cfg.DatabaseDir)
	assert.False(t, cfg.Debug)
	assert.Equal(t, 50, cfg.ListLimit)
	assert.Equal(t, "desc", cfg.ListOrder)
	assert.True(t, cfg.UseBackgroundTracker)
	assert.Equal(t, 8080, cfg.WebUIPort)
	assert.False(t, cfg.BackgroundTrackerAutoStop)
	assert.Equal(t, 120, cfg.BackgroundTrackerAutoStopAfter)
}

func TestCreateDefault_WithForceOverwrite(t *testing.T) {
	defer resetViper()
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	configContent := `database_dir: "/tmp/old"
debug: true
list_limit: 25
list_order: "asc"
use_background_tracker: false`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	newConfigPath, err := CreateDefault(tempDir, true)
	require.NoError(t, err)
	assert.Equal(t, configPath, newConfigPath)

	cfg, _, err := Load(configPath)
	require.NoError(t, err)
	assert.NotEqual(t, "/tmp/old", cfg.DatabaseDir)
	assert.False(t, cfg.Debug)
	assert.Equal(t, 50, cfg.ListLimit)
	assert.Equal(t, "desc", cfg.ListOrder)
	assert.True(t, cfg.UseBackgroundTracker)
	assert.Equal(t, 8080, cfg.WebUIPort)
	assert.False(t, cfg.BackgroundTrackerAutoStop)
	assert.Equal(t, 120, cfg.BackgroundTrackerAutoStopAfter)
}

func TestCreateDefault_WithoutForceOverwrite(t *testing.T) {
	defer resetViper()
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	configContent := `database_dir: "/tmp/old"
debug: true
list_limit: 25
list_order: "asc"
use_background_tracker: false`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	_, err = CreateDefault(tempDir, false)
	assert.Error(t, err)
}

func TestCreateDefault_WithEmptyDirectoryString(t *testing.T) {
	defer resetViper()
	tempDir := t.TempDir()
	originalLookupPaths := LookupPaths
	LookupPaths = []string{tempDir}
	defer func() { LookupPaths = originalLookupPaths }()

	configPath, err := CreateDefault("", false)
	require.NoError(t, err)
	assert.Contains(t, configPath, "config.yaml")

	_, err = os.Stat(configPath)
	require.NoError(t, err)
}

func TestExpandPath_WithTilde(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	require.NoError(t, err)

	expanded, err := expandPath("~")
	require.NoError(t, err)
	assert.Equal(t, homeDir, expanded)

	expanded, err = expandPath("~/.hora")
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(homeDir, ".hora"), expanded)
}

func TestExpandPath_WithoutTilde(t *testing.T) {
	expanded, err := expandPath("/tmp/test")
	require.NoError(t, err)
	assert.Equal(t, "/tmp/test", expanded)

	expanded, err = expandPath("./test")
	require.NoError(t, err)
	assert.Equal(t, "./test", expanded)
}

func TestGetDefaultDatabaseDir(t *testing.T) {
	dbDir, err := getDefaultDatabaseDir()
	require.NoError(t, err)
	assert.NotEmpty(t, dbDir)
	assert.Contains(t, dbDir, "hora")
}

func TestValidateConfig_WithValidConfig(t *testing.T) {
	cfg := &Config{
		DatabaseDir:                    "/tmp/test",
		Debug:                          true,
		ListLimit:                      100,
		ListOrder:                      "asc",
		UseBackgroundTracker:           false,
		WebUIPort:                      3000,
		BackgroundTrackerAutoStop:      true,
		BackgroundTrackerAutoStopAfter: 60,
	}

	err := validateConfig(cfg)
	assert.NoError(t, err)
}

func TestValidateConfig_WithInvalidBackgroundAutoStopAfter(t *testing.T) {
	cfg := &Config{
		DatabaseDir:                    "/tmp/test",
		Debug:                          true,
		ListLimit:                      50,
		ListOrder:                      "asc",
		UseBackgroundTracker:           true,
		WebUIPort:                      8080,
		BackgroundTrackerAutoStop:      true,
		BackgroundTrackerAutoStopAfter: 0,
	}

	err := validateConfig(cfg)
	assert.Error(t, err)
}

func TestValidateConfig_WithInvalidListLimit(t *testing.T) {
	cfg := &Config{
		DatabaseDir:          "/tmp/test",
		Debug:                true,
		ListLimit:            0,
		ListOrder:            "asc",
		UseBackgroundTracker: false,
	}

	err := validateConfig(cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation error")
}

func TestValidateConfig_WithInvalidListOrder(t *testing.T) {
	cfg := &Config{
		DatabaseDir:          "/tmp/test",
		Debug:                true,
		ListLimit:            50,
		ListOrder:            "invalid",
		UseBackgroundTracker: false,
	}

	err := validateConfig(cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation error")
}

func TestValidateConfig_WithMultipleValidationErrors(t *testing.T) {
	cfg := &Config{
		DatabaseDir:          "/tmp/test",
		Debug:                true,
		ListLimit:            0,
		ListOrder:            "invalid",
		UseBackgroundTracker: false,
	}

	err := validateConfig(cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation errors")
}
