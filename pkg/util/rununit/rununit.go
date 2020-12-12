package rununit

import "sync"

type RunUnit struct {
	sync.Mutex
	name        string
	running     bool
	commandChan chan string
	quit        chan struct{}
}

func New(run, stop func(), args ...string) (out *RunUnit) {
	out = &RunUnit{commandChan: make(chan string), name: args[0]}
	go func() {
	out:
		for {
			select {
			case cmd := <-out.commandChan:
				switch cmd {
				case "run":
					Debug("run called for", args[0])

					run()
					out.running = true
				case "stop":
					Debug("stop called for", args[0])

					stop()
					out.running = false
				}
			case <-out.quit:
				out.commandChan <- "stop"
				break out
			}
		}
	}()
	return
}
//
// func (r *RunUnit) Chan() chan<- string {
// 	return r.commandChan
// }

func (r *RunUnit) Running() bool {
	return r.running
}

func (r *RunUnit) Run() {
	r.commandChan <- "run"
	r.Lock()
	defer r.Unlock()
}

func (r *RunUnit) Stop() {
	r.commandChan <- "stop"
	r.Lock()
	defer r.Unlock()
	r.running = false
}

func (r *RunUnit) Shutdown() {
	close(r.quit)
}
