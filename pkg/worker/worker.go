package worker

import (
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/stdconn"
	"os"
	"os/exec"
)

type Worker struct {
	*exec.Cmd
	args    []string
	StdConn stdconn.StdConn
}

func SpawnWorker(args ...string) (w *Worker) {
	w = &Worker{
		Cmd:  exec.Command(args[0], args[1:]...),
		args: args,
	}
	w.Stderr = os.Stdout
	cmdOut, err := w.StdoutPipe()
	if err != nil {
		log.ERROR(err)
		return
	}
	cmdIn, err := w.StdinPipe()
	if err != nil {
		log.ERROR(err)
		return
	}
	w.StdConn = stdconn.New(cmdOut, cmdIn, make(chan struct{}))
	err = w.Start()
	if err != nil {
		log.ERROR(err)
		return nil
	} else {
		return
	}
}

func (w *Worker) KillWorker() (err error) {
	err = w.StdConn.Close()
	if err != nil {
		log.ERROR(err)
	}
	err = w.Wait()
	if err != nil {
		log.ERROR(err)
	}
	return
}
