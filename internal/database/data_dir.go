package database

import (
	"os"
	"path/filepath"
	"runtime"
)

// GetDataDir returns the appropriate data directory for the current operating system
func GetDataDir() (string, error) {
	var dataDir string

	switch runtime.GOOS {
	case "windows":
		// Use AppData\Local on Windows
		appData := os.Getenv("LOCALAPPDATA")
		if appData == "" {
			appData = os.Getenv("APPDATA")
		}
		if appData == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			appData = filepath.Join(homeDir, "AppData", "Local")
		}
		dataDir = filepath.Join(appData, "Hora")
	case "darwin":
		// Use ~/Library/Application Support on macOS
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		dataDir = filepath.Join(homeDir, "Library", "Application Support", "Hora")
	default:
		// Use ~/.local/share on Linux and other Unix-like systems
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		dataDir = filepath.Join(homeDir, ".local", "share", "hora")
	}

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return "", err
	}

	return dataDir, nil
}

// GetDatabasePath returns the full path to the SQLite database file
func GetDatabasePath() (string, error) {
	dataDir, err := GetDataDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dataDir, "hora.db"), nil
}
