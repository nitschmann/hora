package service

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Logger provides structured logging for the agent
type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	logFile     *os.File
}

// NewLogger creates a new logger instance
func NewLogger() (*Logger, error) {
	// Create logs directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	logsDir := filepath.Join(homeDir, ".hora", "logs")
	err = os.MkdirAll(logsDir, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Create log file with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logFileName := fmt.Sprintf("hora-agent-%s.log", timestamp)
	logFilePath := filepath.Join(logsDir, logFileName)

	// Open log file
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	// Create multi-writer to write to both file and stdout
	multiWriter := io.MultiWriter(logFile, os.Stdout)

	// Create loggers
	infoLogger := log.New(multiWriter, "[INFO] ", log.LstdFlags|log.Lshortfile)
	errorLogger := log.New(multiWriter, "[ERROR] ", log.LstdFlags|log.Lshortfile)

	return &Logger{
		infoLogger:  infoLogger,
		errorLogger: errorLogger,
		logFile:     logFile,
	}, nil
}

// Info logs an info message
func (l *Logger) Info(format string, v ...interface{}) {
	l.infoLogger.Printf(format, v...)
}

// Error logs an error message
func (l *Logger) Error(format string, v ...interface{}) {
	l.errorLogger.Printf(format, v...)
}

// Close closes the log file
func (l *Logger) Close() error {
	if l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
}

// GetLogFilePath returns the path to the current log file
func (l *Logger) GetLogFilePath() string {
	if l.logFile != nil {
		return l.logFile.Name()
	}
	return ""
}

// GetLatestLogFile returns the path to the most recent log file
func GetLatestLogFile() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	logsDir := filepath.Join(homeDir, ".hora", "logs")

	// Check if logs directory exists
	if _, err := os.Stat(logsDir); os.IsNotExist(err) {
		return "", fmt.Errorf("logs directory does not exist")
	}

	// Read directory to find the latest log file
	entries, err := os.ReadDir(logsDir)
	if err != nil {
		return "", fmt.Errorf("failed to read logs directory: %w", err)
	}

	var latestFile string
	var latestTime time.Time

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".log" {
			info, err := entry.Info()
			if err != nil {
				continue
			}

			if info.ModTime().After(latestTime) {
				latestTime = info.ModTime()
				latestFile = filepath.Join(logsDir, entry.Name())
			}
		}
	}

	if latestFile == "" {
		return "", fmt.Errorf("no log files found")
	}

	return latestFile, nil
}
