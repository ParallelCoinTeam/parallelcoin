package rununit

import (
	uberatomic "go.uber.org/atomic"
	
	"github.com/p9c/pod/pkg/comm/stdconn/worker"
	"github.com/p9c/pod/pkg/util/logi"
	"github.com/p9c/pod/pkg/util/logi/pipe/consume"
)

// RunUnit handles correctly starting and stopping child processes that have StdConn pipe logging enabled, allowing
// custom hooks to run on start and stop,
type RunUnit struct {
	running, shuttingDown uberatomic.Bool
	commandChan           chan bool
	worker                *worker.Worker
	quit                  chan struct{}
}

// New creates and starts a new rununit. run and stop functions are executed after starting and stopping. logger
// receives log entries and processes them (such as logging them).
func New(
	run, stop func(), logger func(ent *logi.Entry) (err error), pkgFilter func(pkg string) (out bool),
	args ...string,
) (r *RunUnit) {
	r = &RunUnit{
		commandChan: make(chan bool),
		quit:        make(chan struct{}),
	}
	r.running.Store(false)
	r.shuttingDown.Store(false)
	go func() {
	out:
		for {
			select {
			case cmd := <-r.commandChan:
				switch cmd {
				case true:
					Debug("run called for", args)
					if r.running.Load() {
						Debug("already running", args)
						continue
					}
					r.worker = consume.Log(r.quit, logger, pkgFilter, args...)
					// Debug(r.worker)
					consume.Start(r.worker)
					r.running.Store(true)
					run()
					Debug(r.running.Load())
				case false:
					Debug(r.running.Load())
					Debug("stop called for", args)
					if !r.running.Load() {
						Debug("wasn't running", args)
						continue
					}
					consume.Kill(r.worker)
					r.running.Store(false)
					stop()
				}
			case <-r.quit:
				Debug("quitting on run unit quit channel", args, r.running.Load())
				if r.running.Load() {
					Debug("wasn't running", args)
					// continue
				}
				consume.Kill(r.worker)
				r.running.Store(false)
				stop()
				// r.commandChan <- false
				break out
			}
		}
	}()
	return
}

// Running returns whether the unit is running
func (r *RunUnit) Running() bool {
	return r.running.Load()
}

// Run signals the run unit to start
func (r *RunUnit) Start() {
	r.commandChan <- true
}

// Stop signals the run unit to stop
func (r *RunUnit) Stop() {
	r.commandChan <- false
}

// Shutdown terminates the run unit
func (r *RunUnit) Shutdown() {
	// debug.PrintStack()
	if !r.shuttingDown.Load() && r.running.Load() {
		r.shuttingDown.Store(true)
		close(r.quit)
	}
}
