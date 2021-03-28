// +build !windows

package worker

import (
	"syscall"
)

// Pause sends a signal to the worker process to stop
func (w *Worker) Pause() (e error) {
	if e = w.Cmd.Process.Signal(syscall.SIGSTOP); !E.Chk(e) {
		D.Ln("paused")
	}
	return
}

// Continue sends a signal to a worker process to resume work
func (w *Worker) Continue() (e error) {
	if e = w.Cmd.Process.Signal(syscall.SIGCONT); !E.Chk(e) {
		D.Ln("resumed")
	}
	return
}
