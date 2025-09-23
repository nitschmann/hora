package apperror

import "fmt"

type databaseError struct {
	OriginalError error
}

// NewDatabaseError creates a new database error wrapping the original error
func NewDatabaseError(err error) error {
	return &databaseError{
		OriginalError: err,
	}
}

// Error implements the error interface for databaseError
func (e *databaseError) Error() string {
	return fmt.Sprintf("unexpected database error: %v", e.OriginalError)
}
