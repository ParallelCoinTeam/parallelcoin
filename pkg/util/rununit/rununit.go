package rununit

import (
	uberatomic "go.uber.org/atomic"

	"github.com/p9c/pod/pkg/comm/stdconn/worker"
	"github.com/p9c/pod/pkg/util/logi"
	"github.com/p9c/pod/pkg/util/logi/consume"
)

// RunUnit handles correctly starting and stopping child processes that have StdConn pipe logging enabled, allowing
// custom hooks to run on start and stop,
type RunUnit struct {
	running     uberatomic.Bool
	commandChan chan bool
	worker      *worker.Worker
	quit        chan struct{}
}

// New creates and starts a new rununit. run and stop functions are executed after starting and stopping. logger
// receives log entries and processes them (such as logging them).
func New(run, stop func(), logger func(ent *logi.Entry) (err error), pkgFilter func(pkg string) (out bool),
	args ...string) (r *RunUnit) {
	r = &RunUnit{
		commandChan: make(chan bool),
		quit:        make(chan struct{}),
	}
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
					consume.Start(r.worker)
					run()
					r.running.Store(true)
				case false:
					Debug("stop called for", args)
					if r.running.Load() {
						Debug("wasn't running", args)
						continue
					}
					consume.Kill(r.worker)
					stop()
					r.running.Store(false)
				}
			case <-r.quit:
				r.commandChan <- false
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
func (r *RunUnit) Run() {
	r.commandChan <- true
}

// Stop signals the run unit to stop
func (r *RunUnit) Stop() {
	r.commandChan <- false
}

// Shutdown terminates the run unit
func (r *RunUnit) Shutdown() {
	close(r.quit)
}
