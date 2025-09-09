package model

import "time"

// TimeEntry represents a time tracking entry
type TimeEntry struct {
	ID        int            `json:"id" db:"id"`
	ProjectID int            `json:"project_id" db:"project_id"`
	Project   *Project       `json:"project,omitempty" db:"-"`
	StartTime time.Time      `json:"start_time" db:"start_time"`
	EndTime   *time.Time     `json:"end_time,omitempty" db:"end_time"`
	Duration  *time.Duration `json:"duration,omitempty" db:"duration"`
	Category  *string        `json:"category,omitempty" db:"category"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
}
