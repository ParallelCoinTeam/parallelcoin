// +build !windows

package old

import (
	"os/exec"
	"syscall"
)

func pause(s *State, c *exec.Cmd) {
	if err := c.Process.Signal(syscall.SIGSTOP); !Check(err) {
		s.Config.Pausing = true
		Debug("paused")
	}
}
func resume(s *State, c *exec.Cmd) {
	if err := c.Process.Signal(syscall.SIGCONT); !Check(err) {
		s.Config.Pausing = false
		Debug("resumed")
	}
}
