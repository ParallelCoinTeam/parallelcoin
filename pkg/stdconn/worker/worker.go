package worker

import (
	"os"
	"os/exec"

	"github.com/p9c/pod/pkg/stdconn"
)

type Worker struct {
	cmd *exec.Cmd
	args    []string
	StdConn stdconn.StdConn
}

// Spawn starts up an arbitrary executable file with given arguments and
// attaches a connection to its stdin/stdout
func Spawn(args ...string) (w *Worker) {
	w = &Worker{
		cmd:  exec.Command(args[0], args[1:]...),
		args: args,
	}
	//w.Stderr = os.Stderr
	cmdOut, err := w.cmd.StdoutPipe()
	if err != nil {
		L.Error(err)
		return
	}
	cmdIn, err := w.cmd.StdinPipe()
	if err != nil {
		L.Error(err)
		return
	}
	w.StdConn = stdconn.New(cmdOut, cmdIn, make(chan struct{}))
	err = w.cmd.Start()
	if err != nil {
		L.Error(err)
		return nil
	} else {
		return
	}
}

// Kill forces the child process to shut down without cleanup
func (w *Worker) Kill() (err error) {
	return w.cmd.Process.Kill()
}

// Stop signals the worker to shut down cleanly.
// Note that the worker must have handlers for os.Signal messages.
// It is possible and neater to put a quit method in the IPC API and use the
// quit channel built into the StdConn
func (w *Worker) Stop() (err error) {
	return w.cmd.Process.Signal(os.Interrupt)
}
