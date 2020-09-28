// +build !windows

package worker

import (
	"github.com/stalker-loki/app/slog"
	"syscall"
)

func (w *Worker) Pause() (err error) {
	if err = w.cmd.Process.Signal(syscall.SIGSTOP); !slog.Check(err) {
		slog.Debug("paused")
	}
	return
}
func (w *Worker) Resume() (err error) {
	if err = w.cmd.Process.Signal(syscall.SIGCONT); !slog.Check(err) {
		slog.Debug("resumed")
	}
	return
}
