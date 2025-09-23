package apperror

import (
	"fmt"

	"github.com/spf13/viper"
)

type configError struct {
	ErrorMessage  string
	OriginalError error
}

// NewConfigError creates a new config error wrapping the original error and adding a custom message
func NewConfigError(err error, message string) *configError {
	if message == "" {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			vErr := err.(viper.ConfigFileNotFoundError)
			message = vErr.Error()
		default:
			message = fmt.Sprintf("unexpected config error: %v", err)
		}
	}

	return &configError{
		OriginalError: err,
		ErrorMessage:  message,
	}
}

// NewConfigErrorWithoutMessage creates a new config error wrapping the original error without adding a custom message
func NewConfigErrorWithoutMessage(err error) *configError {
	return NewConfigError(err, "")
}

func (e *configError) Error() string {
	return e.ErrorMessage
}
