package interrupt

import (
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	
	"github.com/kardianos/osext"
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
	DataDir      string
)

// Receiver listens for interrupt signals, registers interrupt callbacks, and responds to custom shutdown signals as
// required
func Listener() {
	var interruptCallbacks []func()
	invokeCallbacks := func() {
		Debug("running interrupt callbacks")
		// run handlers in LIFO order.
		for i := range interruptCallbacks {
			idx := len(interruptCallbacks) - 1 - i
			interruptCallbacks[idx]()
		}
		close(HandlersDone)
		Debug("interrupt handlers finished")
		if Restart {
			file, err := osext.Executable()
			if err != nil {
				Error(err)
				return
			}
			Debug("restarting")
			if runtime.GOOS != "windows" {
				err = syscall.Exec(file, os.Args, os.Environ())
				if err != nil {
					Fatal(err)
				}
			} else {
				Debug("doing windows restart")

				// procAttr := new(os.ProcAttr)
				// procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
				// os.StartProcess(os.Args[0], os.Args[1:], procAttr)

				var s []string
				// s = []string{"cmd.exe", "/C", "start"}
				s = append(s, os.Args[0])
				// s = append(s, "--delaystart")
				s = append(s, os.Args[1:]...)
				cmd := exec.Command(s[0], s[1:]...)
				Debug("windows restart done")
				if err = cmd.Start(); Check(err) {
				}
				// // select{}
				// os.Exit(0)
			}
		}
		// time.Sleep(time.Second * 3)
		// os.Exit(1)
	}
	defer Debug("interrupt listener terminated")
	for {
		select {
		case sig := <-Chan:
			// L.Printf("\r>>> received signal (%s)\n", sig)
			Debug("received interrupt signal", sig)
			requested = true
			invokeCallbacks()
			// pprof.Lookup("goroutine").WriteTo(os.Stderr, 2)
			return
		case <-ShutdownRequestChan:
			Warn("received shutdown request - shutting down...")
			requested = true
			invokeCallbacks()
			return
		case handler := <-AddHandlerChan:
			Debug("adding handler")
			interruptCallbacks = append(interruptCallbacks, handler)
		}
	}
}

// AddHandler adds a handler to call when a SIGINT (Ctrl+C) is received.
func AddHandler(handler func()) {
	// Create the channel and start the main interrupt handler which invokes all other callbacks and exits if not
	// already done.
	if Chan == nil {
		Chan = make(chan os.Signal, 1)
		signal.Notify(Chan, Signals...)
		go Listener()
	}
	AddHandlerChan <- handler
}

// Request programmatically requests a shutdown
func Request() {
	Debug("interrupt requested")
	ShutdownRequestChan <- struct{}{}
	// var ok bool
	// select {
	// case _, ok = <-ShutdownRequestChan:
	// default:
	// }
	// Debug("shutdownrequestchan", ok)
	// if ok {
	// 	close(ShutdownRequestChan)
	// }
}

// RequestRestart sets the reset flag and requests a restart
func RequestRestart() {
	Restart = true
	Debug("requesting restart")
	Request()
}

// Requested returns true if an interrupt has been requested
func Requested() bool {
	return requested
}
