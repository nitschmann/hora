package backgroundtracker

import (
	"fmt"
	"log/slog"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"sync"
)

const logFilename = "hora-backgroundtracker.log"

var (
	logger *slog.Logger
	once   sync.Once
)

func initLogger() {
	if runtime.GOOS == "darwin" {
		usr, err := user.Current()
		if err != nil {
			panic("Failed to get current user: " + err.Error())
		}

		logDir := filepath.Join(usr.HomeDir, "Library", "Logs")
		err = os.MkdirAll(logDir, 0755)
		if err != nil {
			panic("Failed to create log directory: " + err.Error())
		}

		logPath := filepath.Join(logDir, logFilename)
		f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			panic("Failed to open log file: " + err.Error())
		}

		handler := slog.NewTextHandler(f, &slog.HandlerOptions{Level: slog.LevelInfo})
		logger = slog.New(handler)
		return
	}

	// fallback to stderr/text logging on other systems
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger = slog.New(handler)
}

// Logger returns the singleton logger instance for the background tracker
func Logger() *slog.Logger {
	once.Do(initLogger)
	return logger
}

// GetLogPath returns the path to the used log file
func GetLogPath() (string, error) {
	if runtime.GOOS != "darwin" {
		return "", fmt.Errorf("log path is only available on macOS")
	}

	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("failed to get current user: %w", err)
	}

	logDir := filepath.Join(usr.HomeDir, "Library", "Logs")
	return filepath.Join(logDir, logFilename), nil
}
