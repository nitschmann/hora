package model

import "time"

// Pause represents a pause in time tracking
type Pause struct {
	ID          int            `json:"id" db:"id"`
	TimeEntryID int            `json:"time_entry_id" db:"time_entry_id"`
	PauseStart  time.Time      `json:"pause_start" db:"pause_start"`
	PauseEnd    *time.Time     `json:"pause_end,omitempty" db:"pause_end"`
	Duration    *time.Duration `json:"duration,omitempty" db:"duration"`
	CreatedAt   time.Time      `json:"created_at" db:"created_at"`
}
