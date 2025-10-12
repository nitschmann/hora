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
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/nitschmann/hora/internal/config"
	"github.com/nitschmann/hora/internal/service"
)

var (
	cfg              config.Config
	stopCheckChannel = make(chan struct{})
)

//export onScreenLocked
func onScreenLocked() {
	Logger().Info("Screen locked - attempting to pause time tracking")

	if timeService != nil {
		ctx := context.Background()
		err := timeService.PauseTracking(ctx)
		if err != nil {
			Logger().Error("Failed to pause time tracking on screen lock", "error", err)
		} else {
			activeEntry, err := timeService.GetActiveEntry(ctx)
			if err != nil {
				Logger().Error("Failed to get active entry after pausing", "error", err)
			}

			Logger().Info(
				"Time tracking paused due to screen lock",
				"project", activeEntry.Project.Name,
			)

			// start monitoring pause duration
			go monitorPauseDuration(ctx)
		}
	} else {
		// very unlikely case - maybe even panic?
		Logger().Warn("Time tracking service not available for screen lock pause")
	}
}

//export onScreenUnlocked
func onScreenUnlocked() {
	Logger().Info("Screen unlocked - attempting to resume time tracking")

	// stop pause monitoring if running
	close(stopCheckChannel)
	stopCheckChannel = make(chan struct{})

	if timeService != nil {
		ctx := context.Background()
		err := timeService.ContinueTracking(ctx)
		if err != nil {
			if !strings.Contains(err.Error(), "no active pause") {
				Logger().Error("Failed to resume time tracking on screen unlock", "error", err)
			} else {
				Logger().Info("Screen unlocked but no active pause to continue")
			}
		} else {
			activeEntry, err := timeService.GetActiveEntry(ctx)
			if err != nil {
				Logger().Error("Failed to get active entry after resuming", "error", err)
			}

			Logger().Info(
				"Time tracking resumed due to screen unlock",
				"project", activeEntry.Project.Name,
			)
		}
	} else {
		// very unlikely case - maybe even panic?
		Logger().Warn("Time tracking service not available for screen unlock resume")
	}
}

func monitorPauseDuration(ctx context.Context) {
	if !cfg.BackgroundTrackerAutoStop {
		Logger().Info("Auto-stop on long pause is disabled, not monitoring pause duration")
		return
	}

	pauseLimit := time.Duration(cfg.BackgroundTrackerAutoStopAfter) * time.Minute
	start := time.Now()
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			elapsed := time.Since(start)
			Logger().Info("Monitoring pause duration...", "elapsed", elapsed.String(), "limit", pauseLimit.String())
			if elapsed >= pauseLimit {
				timeEntry, err := timeService.StopTracking(ctx)
				if err != nil {
					Logger().Error("Failed to stop tracking after long pause", "error", err)
				}

				Logger().Info(
					"Tracking session stopped due to long pause",
					"project", timeEntry.Project.Name,
				)

				Stop()
				os.Exit(0)
			}
		case <-stopCheckChannel:
			Logger().Info("Pause duration monitoring stopped (seesion resumed)")
			return
		}
	}
}

// Start begins listening for screen lock events
func Start(conf *config.Config, timeService service.TimeTracking) {
	SetTimeTrackingService(timeService)
	if conf != nil {
		// dereference to avoid potential nil pointer dereference later
		cfg = *conf
	}

	Logger().Info("Starting screen lock detection...")

	// Handle SIGTERM / SIGINT for graceful shutdown
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
		s := <-c
		Logger().Info("Received shutdown signal, cleaning up...", "signal", s)

		if timeService != nil {
			ctx := context.Background()
			activeEntry, err := timeService.GetActiveEntry(ctx)

			if activeEntry != nil && err == nil {
				entry, err := timeService.StopTracking(ctx)
				if err != nil {
					Logger().ErrorContext(
						ctx,
						"Failed to stop tracking during shutdown",
						"error", err,
					)
				} else {
					Logger().InfoContext(
						ctx,
						"Tracking session stopped due to shutdown",
						"project", entry.Project.Name,
					)
				}
			} else {
				Logger().Info("No active tracking session to stop during shutdown")
			}
		}

		_ = os.Remove(pidFile)

		Logger().Info("Daemon shut down cleanly")
		os.Exit(0)
	}()

	C.startLockEventListenerHora()
}
