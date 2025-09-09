package backgroundtracker

import "github.com/nitschmann/hora/internal/service"

// global variable to hold the time tracking service
var timeService service.TimeTracking

// SetTimeTrackingService sets the time tracking service for screen lock integration
func SetTimeTrackingService(ts service.TimeTracking) {
	timeService = ts
}
