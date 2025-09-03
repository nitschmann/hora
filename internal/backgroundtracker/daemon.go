package backgroundtracker

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"syscall"
)

var pidFile = filepath.Join(os.TempDir(), "hora-backgroundtracker.pid")

func Daemonize() {
	if runtime.GOOS != "darwin" {
		// Daemonization is only implemented for macOS for now
		return
	}

	if os.Getenv("IS_DAEMON") != "1" {
		attr := &syscall.ProcAttr{
			Files: []uintptr{0, 1, 2},
			Env:   append(os.Environ(), "IS_DAEMON=1"),
		}
		pid, err := syscall.ForkExec(os.Args[0], os.Args, attr)
		if err != nil {
			Logger().Error("Failed to daemonize", "err", err)
			os.Exit(1)
		}

		// Parent: write PID file
		if err := os.WriteFile(pidFile, []byte(strconv.Itoa(pid)), 0644); err != nil {
			Logger().Error("Failed to write PID file", "err", err)
			os.Exit(1)
		}

		// TODO: Remove this logging to stdout once stable
		fmt.Println("Background daemon started with PID", pid)

		os.Exit(0)
	}
}

// Stop stops the background tracker daemon if it's running
func Stop() error {
	if runtime.GOOS != "darwin" {
		Logger().Error("Stop is only supported on macOS")
		return nil
	}

	return stopDaemon(pidFile)
}

// IsRunning checks if the background tracker daemon is currently running
func IsRunning() bool {
	// only supported on macOS for now
	if runtime.GOOS != "darwin" {
		return false
	}

	data, err := os.ReadFile(pidFile)
	if err != nil {
		return false
	}

	pid, err := strconv.Atoi(string(data))
	if err != nil {
		return false
	}

	return syscall.Kill(pid, 0) == nil
}

func stopDaemon(pidFile string) error {
	data, err := os.ReadFile(pidFile)
	if err != nil {
		return fmt.Errorf("could not read pid file: %w", err)
	}

	pid, err := strconv.Atoi(string(data))
	if err != nil {
		return fmt.Errorf("invalid pid file: %w", err)
	}

	if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
		return fmt.Errorf("failed to kill process %d: %w", pid, err)
	}

	_ = os.Remove(pidFile)
	Logger().Info("Daemon stopped", "pid", pid)
	return nil
}
