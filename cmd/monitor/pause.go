// +build !windows

package monitor

import (
	"github.com/stalker-loki/app/slog"
	"os/exec"
	"syscall"
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
