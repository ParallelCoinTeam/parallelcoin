package rununit

import (
	uberatomic "go.uber.org/atomic"

	"github.com/p9c/log"
	"github.com/p9c/matrjoska/pkg/pipe"

	"github.com/p9c/interrupt"
	"github.com/p9c/qu"
)

// RunUnit handles correctly starting and stopping child processes that have StdConn pipe logging enabled, allowing
// custom hooks to run on start and stop,
type RunUnit struct {
	name                  string
	args                  []string
	running, shuttingDown uberatomic.Bool
	commandChan           chan bool
	worker                *pipe.Worker
	quit                  qu.C
}

// New creates and starts a new rununit. run and stop functions are executed after starting and stopping. logger
// receives log entries and processes them (such as logging them).
func New(
	name string, run, stop func(), logger func(ent *log.Entry) (e error),
	pkgFilter func(pkg string) (out bool), quit qu.C, args ...string,
) (r *RunUnit) {
	r = &RunUnit{
		name:        name,
		args:        args,
		commandChan: make(chan bool),
		quit:        qu.T(),
	}
	r.running.Store(false)
	r.shuttingDown.Store(false)
	go func() {
		D.Ln("run unit command loop", args)
		var e error
	out:
		for {
			select {
			case cmd := <-r.commandChan:
				switch cmd {
				case true:
					D.Ln(r.running.Load(), "run called for", args)
					if r.running.Load() {
						D.Ln("already running", args)
						continue
					}
					if r.worker != nil {
						if e = r.worker.Kill(); E.Chk(e) {
						}
					}
					// quit from rununit's quit, which closes after the main quit triggers stopping in the watcher loop
					r.worker = pipe.LogConsume(r.quit, logger, pkgFilter, args...)
					// D.S(r.worker)
					pipe.Start(r.worker)
					r.running.Store(true)
					run()
					// D.Ln(r.running.Load())
				case false:
					running := r.running.Load()
					D.Ln("stop called for", args, running)
					if !running {
						D.Ln("wasn't running", args)
						continue
					}
					pipe.Kill(r.worker)
					// var e error
					// if e = r.worker.Wait(); E.Chk(e) {
					// }
					r.running.Store(false)
					stop()
					D.Ln(args, "after stop", r.running.Load())
				}
				break
			case <-r.quit.Wait():
				D.Ln("runner stopped for", args)
				break out
			}
		}
	}()
	// when the main quit signal is triggered, stop the run unit cleanly
	go func() {
	out:
		select {
		case <-quit.Wait():
			D.Ln("runner quit trigger called", args)
			running := r.running.Load()
			if !running {
				D.Ln("wasn't running", args)
				break out
			}
			// r.quit.Q()
			pipe.Kill(r.worker)
			var e error
			if e = r.worker.Wait(); E.Chk(e) {
			}
			r.running.Store(false)
			stop()
			D.Ln(args, "after stop", r.running.Load())
		}
	}()
	interrupt.AddHandler(
		func() {
			quit.Q()
		},
	)
	return
}

// Running returns whether the unit is running
func (r *RunUnit) Running() bool {
	return r.running.Load()
}

// Start signals the run unit to start
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
	if !r.shuttingDown.Load() {
		r.shuttingDown.Store(true)
		r.quit.Q()
	}
}

// ShuttingDown returns true if the server is shuting down
func (r *RunUnit) ShuttingDown() bool {
	return r.shuttingDown.Load()
}
