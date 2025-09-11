package backgroundtracker

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

var pidFile = filepath.Join(os.TempDir(), "hora-backgroundtracker.pid")

func Daemonize() {
	if runtime.GOOS != "darwin" {
		// Daemonization is only implemented for macOS for now
		return
	}

	if os.Getenv("IS_DAEMON") != "1" {
		binaryPath, err := exec.LookPath(os.Args[0])
		if err != nil {
			Logger().Error("Failed to find binary path", "err", err, "binary", os.Args[0])
			os.Exit(1)
		}

		Logger().Info("Starting daemonization", "binary", binaryPath, "args", os.Args)

		attr := &syscall.ProcAttr{
			Files: []uintptr{0, 1, 2},
			Env:   append(os.Environ(), "IS_DAEMON=1"),
		}
		pid, err := syscall.ForkExec(binaryPath, os.Args, attr)
		if err != nil {
			Logger().Error("Failed to daemonize", "err", err, "binary", binaryPath)
			os.Exit(1)
		}

		Logger().Info("Daemon started", "pid", pid)

		// Parent: write PID file
		if err := os.WriteFile(pidFile, []byte(strconv.Itoa(pid)), 0644); err != nil {
			Logger().Error("Failed to write PID file", "err", err, "pidFile", pidFile)
			os.Exit(1)
		}

		os.Exit(0)
	} else {
		Logger().Info("Running as daemon", "pid", os.Getpid())
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

	// Wait until the daemon removes its PID file (graceful shutdown)
	const maxWait = 2 // seconds
	for i := 0; i < maxWait*10; i++ {
		if _, err := os.Stat(pidFile); os.IsNotExist(err) {
			Logger().Info("Daemon stopped cleanly", "pid", pid)
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}

	_ = os.Remove(pidFile)
	Logger().Warn("Daemon did not clean up PID file, removed manually", "pid", pid)
	return nil
}
