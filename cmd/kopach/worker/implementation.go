package worker

import (
	"net"
	"os"

	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/sem"
	"github.com/p9c/pod/pkg/stdconn"
)

type handleFunc func(blk *wire.MsgBlock, sem sem.T)

type Worker struct {
	sem sem.T
	conn net.Conn
	Quit    chan struct{}
	handler handleFunc
}

// NewWithConnAndSemaphore is exposed to enable use an actual network
// connection while retaining the same RPC API to allow a worker to be
// configured to run on a bare metal system with a different launcher main
func NewWithConnAndSemaphore(
	handler handleFunc,
	conn net.Conn,
	s sem.T,
	quit chan struct{},
) *Worker {
	log.DEBUG("creating new Worker")
	return &Worker{
		sem:       s,
		Quit:    quit,
		conn:    conn,
		handler: handler,
	}
}

// New initialises the state for a worker,
// loading the work function handler that runs a round of processing between
// checking quit signal and work semaphore
func New(handler handleFunc, s sem.T) (w *Worker, conn net.Conn) {
	quit := make(chan struct{})
	conn = stdconn.New(os.Stdin, os.Stdout, quit)
	return NewWithConnAndSemaphore(
		handler,
		conn,
		s,
		quit), conn
}

// NewJob is a delivery of a new job for the worker, this starts a miner thread
func (w *Worker) NewJob(blk *wire.MsgBlock, reply *bool) (err error) {
	log.DEBUG("received new job")
	*reply = true
	// previous thread loses its semaphore when a new job arrives
	w.sem.Acquire()
	go func() {
	out:
		for {
			// mine!
			w.handler(blk, w.sem)
			select {
			case <-w.sem.Release():
				// yield when w.Pause() sends w.Acquire()
				log.DEBUG("pausing work")
				break out
			case <-w.Quit:
				// quit when w.Stop() is called
				log.DEBUG("worker stopping on quit message")
				break out
			default:
			}
		}
		log.DEBUG("finished job")
	}()
	return
}

// Pause signals the worker to stop working,
// releases its semaphore and the worker is then idle
func (w *Worker) Pause(_ int, reply *bool) (err error) {
	w.sem.Acquire()
	*reply = true
	return
}

// Stop signals the worker to quit
func (w *Worker) Stop(_ int, reply *bool) (err error) {
	close(w.Quit)
	*reply = true
	return
}
