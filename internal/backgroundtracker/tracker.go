//go:build !darwin

package backgroundtracker

import "fmt"

func Start() {
	// On non-macOS systems, just log or noop
	fmt.Println("Screen lock tracking is only supported on macOS")
}
