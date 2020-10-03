package worker

import (
	"io"
	"os"
	"os/exec"
	"runtime"
	"syscall"

	"github.com/p9c/pod/pkg/comm/stdconn"
)

type Worker struct {
	cmd     *exec.Cmd
	args    []string
	StdConn stdconn.StdConn
}

// Spawn starts up an arbitrary executable file with given arguments and
// attaches a connection to its stdin/stdout
func Spawn(args ...string) (w *Worker, err error) {
	if runtime.GOOS == "windows" {
		args = append([]string{"cmd.exe", "/C", "start"}, args...)
	}
	w = &Worker{
		cmd:  exec.Command(args[0], args[1:]...),
		args: args,
	}
	// w.cmd.Stderr = os.Stderr
	var cmdOut io.ReadCloser
	if cmdOut, err = w.cmd.StdoutPipe(); Check(err) {
		return
	}
	var cmdIn io.WriteCloser
	if cmdIn, err = w.cmd.StdinPipe(); Check(err) {
		return
	}
	w.StdConn = stdconn.New(cmdOut, cmdIn, make(chan struct{}))
	if err = w.cmd.Start(); Check(err) {
	}
	return
}

func (w *Worker) Wait() (err error) {
	return w.cmd.Wait()
}

func (w *Worker) Interrupt() (err error) {
	if err = w.cmd.Process.Signal(syscall.SIGINT); !Check(err) {
		Debug("interrupted")
	}
	if err = w.cmd.Process.Release(); !Check(err) {
		Debug("released")
	}
	return
}

// Kill forces the child process to shut down without cleanup
func (w *Worker) Kill() (err error) {
	if err = w.cmd.Process.Signal(syscall.SIGKILL); !Check(err) {
		Debug("killed")
	}
	return
}

// Stop signals the worker to shut down cleanly.
//
// Note that the worker must have handlers for os.Signal messages.
//
// It is possible and neater to put a quit method in the IPC API and use the quit channel built into the StdConn
func (w *Worker) Stop() (err error) {
	return w.cmd.Process.Signal(os.Interrupt)
}
