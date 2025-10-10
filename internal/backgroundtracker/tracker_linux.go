//go:build !darwin

package backgroundtracker

import (
	"github.com/nitschmann/hora/internal/config"
	"github.com/nitschmann/hora/internal/service"
)

func Start(_ *config.Config, _ service.TimeTracking) {
	return
}
