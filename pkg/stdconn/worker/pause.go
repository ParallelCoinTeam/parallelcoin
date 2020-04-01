// +build !windows

package worker

import (
	"syscall"
)

func (w *Worker) Pause() (err error) {
	if err = w.cmd.Process.Signal(syscall.SIGSTOP); !Check(err) {
		Debug("paused")
	}
	return
}
func (w *Worker) Resume() (err error) {
	if err = w.cmd.Process.Signal(syscall.SIGCONT); !Check(err) {
		Debug("resumed")
	}
	return
}
