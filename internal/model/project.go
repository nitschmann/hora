package model

import "time"

// Project represents a project in the system
type Project struct {
	ID            int        `json:"id" db:"id"`
	Name          string     `json:"name" db:"name"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	LastTrackedAt *time.Time `json:"last_tracked_at,omitempty" db:"last_tracked_at"`
}
