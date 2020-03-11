package interrupt

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kardianos/osext"

	log "github.com/p9c/logi"
)

var (
	Restart   bool // = true
	requested bool
	// Chan is used to receive SIGINT (Ctrl+C) signals.
	Chan chan os.Signal
	// Signals is the list of signals that cause the interrupt
	Signals = []os.Signal{os.Interrupt}
	// ShutdownRequestChan is a channel that can receive shutdown requests
	ShutdownRequestChan = make(chan struct{})
	// AddHandlerChan is used to add an interrupt handler to the list of
	// handlers to be invoked on SIGINT (Ctrl+C) signals.
	AddHandlerChan = make(chan func())
	// HandlersDone is closed after all interrupt handlers run the first time
	// an interrupt is signaled.
	HandlersDone = make(chan struct{})
)

// Receiver listens for interrupt signals, registers interrupt callbacks, and
// responds to custom shutdown signals as required
func Listener() {
	var interruptCallbacks []func()
	invokeCallbacks := func() {
		log.L.Debug("running interrupt callbacks")
		// run handlers in LIFO order.
		for i := range interruptCallbacks {
			idx := len(interruptCallbacks) - 1 - i
			interruptCallbacks[idx]()
		}
		close(HandlersDone)
		log.L.Debug("interrupt handlers finished")
		if Restart {
			log.L.Debug("restarting")
			file, err := osext.Executable()
			if err != nil {
				log.L.Error(err)
				return
			}
			err = syscall.Exec(file, os.Args, os.Environ())
			if err != nil {
				log.L.Fatal(err)
			}
			// return
		}
		time.Sleep(time.Second)
		os.Exit(1)
	}
	for {
		select {
		case sig := <-Chan:
			// log.L.Printf("\r>>> received signal (%s)\n", sig)
			log.L.Debug("received interrupt signal", sig)
			requested = true
			invokeCallbacks()
			return
		case <-ShutdownRequestChan:
			log.L.Warn("received shutdown request - shutting down...")
			requested = true
			invokeCallbacks()
			return
		case handler := <-AddHandlerChan:
			log.L.Debug("adding handler")
			interruptCallbacks = append(interruptCallbacks, handler)
		}
	}
}

// AddHandler adds a handler to call when a SIGINT (Ctrl+C) is received.
func AddHandler(handler func()) {
	// Create the channel and start the main interrupt handler which invokes all
	// other callbacks and exits if not already done.
	if Chan == nil {
		Chan = make(chan os.Signal, 1)
		signal.Notify(Chan, Signals...)
		go Listener()
	}
	AddHandlerChan <- handler
}

// Request programatically requests a shutdown
func Request() {
	log.L.Debug("interrupt requested")
	ShutdownRequestChan <- struct{}{}
	// var ok bool
	// select {
	// case _, ok = <-ShutdownRequestChan:
	// default:
	// }
	// log.L.Debug("shutdownrequestchan", ok)
	// if ok {
	// 	close(ShutdownRequestChan)
	// }
}

// RequestRestart sets the reset flag and requests a restart
func RequestRestart() {
	Restart = true
	log.L.Debug("requesting restart")
	Request()
}

// Requested returns true if an interrupt has been requested
func Requested() bool {
	return requested
}
