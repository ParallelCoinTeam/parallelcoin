// +build !windows

package monitor

import (
	"os/exec"
	"syscall"
)

func pause(s *State, c *exec.Cmd) {
	if err := c.Process.Signal(syscall.SIGSTOP); !L.Check(err) {
		s.Config.Pausing.Store(true)
		L.Debug("paused")
	}
}
func resume(s *State, c *exec.Cmd) {
	if err := c.Process.Signal(syscall.SIGCONT); !L.Check(err) {
		s.Config.Pausing.Store(false)
		L.Debug("resumed")
	}
}
