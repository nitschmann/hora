package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/nitschmann/hora/internal/database"
	"github.com/nitschmann/hora/internal/repository"
	"github.com/nitschmann/hora/internal/service"
	"github.com/spf13/cobra"
)

var (
	dbConn      *database.Connection
	timeService service.TimeTracking
)

// addListCommandCommonFlags adds common flags for lists to the given cobra command
func addListCommandCommonFlags(cmd *cobra.Command,
	limitVar *int,
	sinceVar *string,
	sortVar *string,
) {
	cmd.Flags().IntVarP(limitVar, "limit", "l", conf.ListLimit, "Maximum number of entries to show")
	cmd.Flags().StringVar(sinceVar, "since", "", "Only show entries since this date (YYYY-MM-DD format)")
	cmd.Flags().StringVar(sortVar, "sort", conf.ListOrder, "Sort order: 'asc' (oldest first) or 'desc' (newest first)")
}

// exportTimesToCSV exports the given time entries to a CSV file with the specified filename
func exportTimesToCSV(entries []repository.TimeEntryWithPauses, filename string, project string) (string, error) {
	if filename == "" {
		timestamp := time.Now().Format("20060102150405")
		filename = fmt.Sprintf("%s_times.csv", timestamp)
	}

	file, err := os.Create(filename)
	if err != nil {
		return filename, err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{
		"Start Time",
		"End Time",
		"Project",
		"Category",
		"Duration",
		"Pauses",
		"Pause Time",
		"Effective Work Time",
	}
	if err := writer.Write(header); err != nil {
		return filename, err
	}

	for _, entry := range entries {
		startTime := formatTimeInLocal(entry.StartTime)
		endTime := "-"
		if entry.EndTime != nil {
			endTime = formatTimeInLocal(*entry.EndTime)
		}

		rowProject := project
		if rowProject == "" {
			rowProject = entry.Project.Name
		}

		category := "-"
		if entry.Category != nil {
			category = *entry.Category
		}

		duration := "-"
		if entry.Duration != nil {
			duration = formatDuration(*entry.Duration)
		}

		pauseTime := "-"
		if entry.PauseTime > 0 {
			pauseTime = formatDuration(entry.PauseTime)
		}

		effectiveTime := "-"
		if entry.Duration != nil {
			effective := *entry.Duration - entry.PauseTime
			if effective > 0 {
				effectiveTime = formatDuration(effective)
			}
		}

		record := []string{
			startTime,
			endTime,
			rowProject,
			category,
			duration,
			strconv.Itoa(entry.PauseCount),
			pauseTime,
			effectiveTime,
		}

		if err := writer.Write(record); err != nil {
			return filename, err
		}
	}

	return filename, nil
}

// formatTimeInLocal formats a time value in the local timezone
func formatTimeInLocal(t time.Time) string {
	return t.Local().Format("2006-01-02 15:04:05")
}

// formatTimeInLocalShort formats a time value in the local timezone with shorter format
func formatTimeInLocalShort(t time.Time) string {
	return t.Local().Format("2006-01-02 15:04")
}

// formatDateInLocal formats a date in the local timezone
func formatDateInLocal(t time.Time) string {
	return t.Local().Format("2006-01-02")
}

// formatDuration formats a duration as HH:MM:SS
func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func initDatabaseConnectionAndService() error {
	var err error
	dbConn, err = database.NewConnection(conf)
	if err != nil {
		return fmt.Errorf("failed to initialize database in daemon: %w", err)
	}

	// repositories for service
	projectRepo := repository.NewProject(dbConn.GetDB())
	timeEntryRepo := repository.NewTimeEntry(dbConn.GetDB())
	pauseRepo := repository.NewPause(dbConn.GetDB())

	timeService = service.NewTimeTracking(projectRepo, timeEntryRepo, pauseRepo)

	return err
}

// validateCategory validates that a category contains only alphanumeric characters, underscores, and hyphens
func validateCategory(category string) error {
	if category == "" {
		return nil
	}

	// Check for common shell special characters that might cause issues
	if strings.ContainsAny(category, "!$`\\") {
		return fmt.Errorf("category contains shell special characters (!$`\\) that may cause issues. Use only alphanumeric characters, underscores (_), and hyphens (-)")
	}

	matched, err := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, category)
	if err != nil {
		return fmt.Errorf("failed to validate category: %w", err)
	}

	if !matched {
		return fmt.Errorf("category must contain only alphanumeric characters, underscores (_), and hyphens (-)")
	}

	return nil
}
