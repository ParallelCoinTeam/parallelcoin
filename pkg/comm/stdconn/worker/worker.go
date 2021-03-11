package worker

import (
	"io"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	
	qu "github.com/p9c/pod/pkg/util/qu"
	
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
func Spawn(quit qu.C, args ...string) (w *Worker, e error) {
	// if runtime.GOOS == "windows" {
	// 	args = append([]string{"Cmd.exe", "/C", "start"}, args...)
	// }
	// args = apputil.PrependForWindows(args)
	// var pipeReader, pipeWriter *os.File
	// if pipeReader, pipeWriter, e = os.Pipe(); err.Chk(e) {
	// }
	w = &Worker{
		Cmd:  exec.Command(args[0], args[1:]...),
		Args: args,
		// Stderr:  pipeWriter,
		// StdPipe: pipeReader,
	}
	// w.Cmd.Stderr = pipeWriter
	var cmdOut io.ReadCloser
	if cmdOut, e = w.Cmd.StdoutPipe(); err.Chk(e) {
		return
	}
	var cmdIn io.WriteCloser
	if cmdIn, e = w.Cmd.StdinPipe(); err.Chk(e) {
		return
	}
	w.StdConn = stdconn.New(cmdOut, cmdIn, quit)
	w.Cmd.Stderr = os.Stderr
	if e = w.Cmd.Start(); err.Chk(e) {
	}
	// data := make([]byte, 8192)
	// go func() {
	// out:
	// 	for {
	// 		select {
	// 		case <-quit:
	// 			dbg.Ln("passed quit chan closed", args)
	// 			break out
	// 		default:
	// 		}
	// 		var n int
	// 		if n, e = w.StdPipe.Read(data); err.Chk(e) {
	// 		}
	// 		// if !onBackup {
	// 		if n > 0 {
	// 			if n, e = os.Stderr.Write(append([]byte("PIPED:\n"), data[:n]...)); err.Chk(e) {
	// 			}
	// 		}
	// 	}
	// }()
	return
}

func (w *Worker) Wait() (e error) {
	return w.Cmd.Wait()
}

func (w *Worker) Interrupt() (e error) {
	if runtime.GOOS == "windows" {
		if e = w.Cmd.Process.Kill(); err.Chk(e) {
		}
		return
	}
	if e = w.Cmd.Process.Signal(syscall.SIGINT); !err.Chk(e) {
		dbg.Ln("interrupted")
	}
	// if e = w.Cmd.Process.Release(); !err.Chk(e) {
	//	dbg.Ln("released")
	// }
	return
}

// Kill forces the child process to shut down without cleanup
func (w *Worker) Kill() (e error) {
	if e = w.Cmd.Process.Kill(); !err.Chk(e) {
		dbg.Ln("killed")
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
