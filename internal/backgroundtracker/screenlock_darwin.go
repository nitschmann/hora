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

//export onScreenLocked
func onScreenLocked() {
	Logger().Info("Screen locked")
}

//export onScreenUnlocked
func onScreenUnlocked() {
	Logger().Info("Screen unlocked")
	// log.Println("Screen unlocked (no handler set)")
}

// Start begins listening for screen lock events
func Start() {
	C.startLockEventListenerHora()
}
