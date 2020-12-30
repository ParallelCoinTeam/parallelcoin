package worker

import (
	qu "github.com/p9c/pod/pkg/util/quit"
	"io"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	
	"github.com/p9c/pod/pkg/comm/stdconn"
)

type Worker struct {
	Cmd  *exec.Cmd
	Args []string
	// Stderr  io.WriteCloser
	// StdPipe io.ReadCloser
	StdConn *stdconn.StdConn
}

// Spawn starts up an arbitrary executable file with given arguments and
// attaches a connection to its stdin/stdout
func Spawn(quit qu.C, args ...string) (w *Worker, err error) {
	// if runtime.GOOS == "windows" {
	// 	args = append([]string{"Cmd.exe", "/C", "start"}, args...)
	// }
	// args = apputil.PrependForWindows(args)
	// var pipeReader, pipeWriter *os.File
	// if pipeReader, pipeWriter, err = os.Pipe(); Check(err) {
	// }
	w = &Worker{
		Cmd:  exec.Command(args[0], args[1:]...),
		Args: args,
		// Stderr:  pipeWriter,
		// StdPipe: pipeReader,
	}
	// w.Cmd.Stderr = pipeWriter
	var cmdOut io.ReadCloser
	if cmdOut, err = w.Cmd.StdoutPipe(); Check(err) {
		return
	}
	var cmdIn io.WriteCloser
	if cmdIn, err = w.Cmd.StdinPipe(); Check(err) {
		return
	}
	w.StdConn = stdconn.New(cmdOut, cmdIn, quit)
	// w.Cmd.Stderr = os.Stderr
	if err = w.Cmd.Start(); Check(err) {
	}
	// data := make([]byte, 8192)
	// go func() {
	// out:
	// 	for {
	// 		select {
	// 		case <-quit:
	// 			Debug("passed quit chan closed", args)
	// 			break out
	// 		default:
	// 		}
	// 		var n int
	// 		if n, err = w.StdPipe.Read(data); Check(err) {
	// 		}
	// 		// if !onBackup {
	// 		if n > 0 {
	// 			if n, err = os.Stderr.Write(append([]byte("PIPED:\n"), data[:n]...)); Check(err) {
	// 			}
	// 		}
	// 	}
	// }()
	return
}

func (w *Worker) Wait() (err error) {
	return w.Cmd.Wait()
}

func (w *Worker) Interrupt() (err error) {
	if runtime.GOOS == "windows" {
		if err = w.Cmd.Process.Kill(); Check(err) {
		}
		return
	}
	if err = w.Cmd.Process.Signal(syscall.SIGINT); !Check(err) {
		Debug("interrupted")
	}
	// if err = w.Cmd.Process.Release(); !Check(err) {
	//	Debug("released")
	// }
	return
}

// Kill forces the child process to shut down without cleanup
func (w *Worker) Kill() (err error) {
	if err = w.Cmd.Process.Kill(); !Check(err) {
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
	return w.Cmd.Process.Signal(os.Interrupt)
}
