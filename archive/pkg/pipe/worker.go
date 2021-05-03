package pipe

import (
	"io"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	
	"github.com/p9c/qu"
)

type Worker struct {
	Cmd  *exec.Cmd
	Args []string
	// Stderr  io.WriteCloser
	// StdPipe io.ReadCloser
	StdConn *StdConn
}

// Spawn starts up an arbitrary executable file with given arguments and
// attaches a connection to its stdin/stdout
func Spawn(quit qu.C, args ...string) (w *Worker, e error) {
	w = &Worker{
		Cmd:  exec.Command(args[0], args[1:]...),
		Args: args,
	}
	var cmdOut io.ReadCloser
	if cmdOut, e = w.Cmd.StdoutPipe(); E.Chk(e) {
		return
	}
	var cmdIn io.WriteCloser
	if cmdIn, e = w.Cmd.StdinPipe(); E.Chk(e) {
		return
	}
	w.StdConn = New(cmdOut, cmdIn, quit)
	w.Cmd.Stderr = os.Stderr
	if e = w.Cmd.Start(); E.Chk(e) {
	}
	return
}

// Wait for the process to finish running
func (w *Worker) Wait() (e error) {
	return w.Cmd.Wait()
}

// Interrupt the child process.
// This invokes kill on windows because windows doesn't have an interrupt
// signal.
func (w *Worker) Interrupt() (e error) {
	if runtime.GOOS == "windows" {
		if e = w.Cmd.Process.Kill(); E.Chk(e) {
		}
		return
	}
	if e = w.Cmd.Process.Signal(syscall.SIGINT); !E.Chk(e) {
		D.Ln("interrupted")
	}
	return
}

// Kill forces the child process to shut down without cleanup
func (w *Worker) Kill() (e error) {
	if e = w.Cmd.Process.Kill(); !E.Chk(e) {
		D.Ln("killed")
	}
	return
}

// Stop signals the worker to shut down cleanly.
//
// Note that the worker must have handlers for os.Signal messages.
//
// It is possible and neater to put a quit method in the IPC API and use the quit channel built into the StdConn
func (w *Worker) Stop() (e error) {
	return w.Cmd.Process.Signal(os.Interrupt)
}
