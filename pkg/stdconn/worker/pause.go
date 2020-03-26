// +build !windows

package worker

import (
	"syscall"
)

func (w *Worker) Pause() (err error) {
	if err = w.cmd.Process.Signal(syscall.SIGSTOP); !L.Check(err) {
		L.Debug("paused")
	}
	return
}
func (w *Worker) Resume() (err error) {
	if err = w.cmd.Process.Signal(syscall.SIGCONT); !L.Check(err) {
		L.Debug("resumed")
	}
	return
}
