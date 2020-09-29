// Package monitor is a log viewer and filter and configuration interface
//
// +build !headless,!windows

package pkg

import (
	"os/exec"
	"syscall"

	"github.com/p9c/pkg/app/slog"
)

func pause(s *State, c *exec.Cmd) {
	if err := c.Process.Signal(syscall.SIGSTOP); !slog.Check(err) {
		s.Config.Pausing = true
		slog.Debug("paused")
	}
}
func resume(s *State, c *exec.Cmd) {
	if err := c.Process.Signal(syscall.SIGCONT); !slog.Check(err) {
		s.Config.Pausing = false
		slog.Debug("resumed")
	}
}
