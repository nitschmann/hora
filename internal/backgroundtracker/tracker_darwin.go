//go:build darwin

package backgroundtracker

/*
#cgo CFLAGS: -x objective-c -framework Foundation
#cgo LDFLAGS: -framework Foundation
#import <Foundation/Foundation.h>

extern void onScreenLocked();
extern void onScreenUnlocked();

static void startLockEventListenerHora(void) {
    NSDistributedNotificationCenter *center = [NSDistributedNotificationCenter defaultCenter];
    [center addObserverForName:@"com.apple.screenIsLocked"
                        object:nil
                         queue:[NSOperationQueue mainQueue]
                    usingBlock:^(NSNotification *note) {
        onScreenLocked();
    }];

    [center addObserverForName:@"com.apple.screenIsUnlocked"
                        object:nil
                         queue:[NSOperationQueue mainQueue]
                    usingBlock:^(NSNotification *note) {
        onScreenUnlocked();
    }];

    [[NSRunLoop mainRunLoop] run];
}
*/
import "C"

import (
	"context"
	"strings"

	"github.com/nitschmann/hora/internal/service"
)

// Global variable to hold the time tracking service
var timeService service.TimeTracking

// SetTimeTrackingService sets the time tracking service for screen lock integration
func SetTimeTrackingService(ts service.TimeTracking) {
	timeService = ts
}

//export onScreenLocked
func onScreenLocked() {
	Logger().Info("Screen locked - attempting to pause time tracking")

	if timeService != nil {
		ctx := context.Background()
		err := timeService.PauseTracking(ctx)
		if err != nil {
			Logger().Error("Failed to pause time tracking on screen lock", "error", err)
		} else {
			Logger().Info("Time tracking paused due to screen lock")
		}
	} else {
		Logger().Warn("Time tracking service not available for screen lock pause")
	}
}

//export onScreenUnlocked
func onScreenUnlocked() {
	Logger().Info("Screen unlocked - attempting to resume time tracking")

	if timeService != nil {
		ctx := context.Background()
		err := timeService.ContinueTracking(ctx)
		if err != nil {
			// If there's no active pause, this is not an error
			if !strings.Contains(err.Error(), "no active pause") {
				Logger().Error("Failed to resume time tracking on screen unlock", "error", err)
			} else {
				Logger().Info("Screen unlocked but no active pause to continue")
			}
		} else {
			Logger().Info("Time tracking resumed due to screen unlock")
		}
	} else {
		Logger().Warn("Time tracking service not available for screen unlock resume")
	}
}

// Start begins listening for screen lock events
func Start() {
	Logger().Info("Starting screen lock detection...")
	C.startLockEventListenerHora()
}

