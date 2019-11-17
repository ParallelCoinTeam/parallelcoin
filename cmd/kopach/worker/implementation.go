package worker

import (
	"fmt"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/sem"
	"github.com/p9c/pod/pkg/stdconn"
	"net"
	"os"
)

type Worker struct {
	sem.T
	net.Conn
	Quit    chan struct{}
	handler func(blk *wire.MsgBlock, sem sem.T)
}

// NewWithConnAndSemaphore is exposed to enable use an actual network
// connection while retaining the same RPC API to allow a worker to be
// configured to run on a bare metal system with a different launcher main
func NewWithConnAndSemaphore(
	handler func(blk *wire.MsgBlock, s sem.T),
	conn net.Conn,
	s sem.T,
	quit chan struct{},
) *Worker {
	printlnE("creating new Worker")
	return &Worker{
		T:       s,
		Quit:    quit,
		Conn:    conn,
		handler: handler,
	}
}

// New initialises the state for a worker,
// loading the work function handler that runs a round of processing between
// checking quit signal and work semaphore
func New(handler func(blk *wire.MsgBlock, s sem.T), s sem.T) *Worker {
	quit := make(chan struct{})
	return NewWithConnAndSemaphore(
		handler,
		stdconn.New(os.Stdin, os.Stdout, quit),
		s,
		quit)
}

// NewJob is a delivery of a new job for the worker, this starts a miner thread
func (w *Worker) NewJob(blk *wire.MsgBlock, reply *bool) (err error) {
	printlnE("received new job")
	*reply = true
	// previous thread loses its semaphore when a new job arrives
	w.Acquire()
	go func() {
	out:
		for {
			// mine!
			w.handler(blk, w.T)
			select {
			case <-w.Release():
				// yield when w.Pause() sends w.Acquire()
				printlnE("pausing work")
				break out
			case <-w.Quit:
				// quit when w.Stop() is called
				printlnE("worker stopping on quit message")
				break out
			default:
			}
		}
		printlnE("finished job")
	}()
	return
}

// Pause signals the worker to stop working,
// releases its semaphore and the worker is then idle
func (w *Worker) Pause(_ int, reply *bool) (err error) {
	w.Acquire()
	*reply = true
	return
}

// Stop signals the worker to quit
func (w *Worker) Stop(_ int, reply *bool) (err error) {
	close(w.Quit)
	*reply = true
	return
}

func printlnE(a ...interface{}) {
	out := append([]interface{}{"[Worker]"}, a...)
	_, _ = fmt.Fprintln(os.Stderr, out...)
}
