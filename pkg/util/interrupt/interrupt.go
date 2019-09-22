package interrupt

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/parallelcointeam/parallelcoin/pkg/log"
)

var (
	requested bool
	// Chan is used to receive SIGINT (Ctrl+C) signals.
	Chan chan os.Signal
	// Signals is the list of signals that cause the interrupt
	Signals = []os.Signal{os.Interrupt, syscall.SIGTERM}
	// ShutdownRequestChan is a channel that can receive shutdown requests
	ShutdownRequestChan = make(chan struct{})
	// AddHandlerChan is used to add an interrupt handler to the list of
	// handlers to be invoked on SIGINT (Ctrl+C) signals.
	AddHandlerChan = make(chan func())
	// HandlersDone is closed after all interrupt handlers run the first time
	// an interrupt is signaled.
	HandlersDone = make(chan struct{})
)

// Listener listens for interrupt signals, registers interrupt callbacks, and
// responds to custom shutdown signals as required
func Listener() {
	var interruptCallbacks []func()
	invokeCallbacks := func() {
		// run handlers in LIFO order.
		for i := range interruptCallbacks {
			idx := len(interruptCallbacks) - 1 - i
			interruptCallbacks[idx]()
		}
		close(HandlersDone)
	}
	for {
		select {
		case sig := <-Chan:
			log.WARNF("received signal (%s) - shutting down", sig)
			_ = sig
			requested = true
			invokeCallbacks()
			return
		case <-ShutdownRequestChan:
			log.WARN("received shutdown request - shutting down...")
			requested = true
			invokeCallbacks()
			return
		case handler := <-AddHandlerChan:
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
	close(ShutdownRequestChan)
}

// Requested returns true if an interrupt has been requested
func Requested() bool {
	return requested
}
