//go:build !darwin

package backgroundtracker

func Start() {
	// On non-macOS systems, just log or noop
	return
}
