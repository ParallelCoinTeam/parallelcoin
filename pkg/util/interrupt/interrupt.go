package interrupt

import (
	"fmt"
	
	uberatomic "go.uber.org/atomic"
	
	qu "github.com/p9c/pod/pkg/util/quit"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	
	"github.com/kardianos/osext"
)

type HandlerWithSource struct {
	Source string
	Fn     func()
}

var (
	Restart   bool // = true
	requested uberatomic.Bool
	// Chan is used to receive SIGINT (Ctrl+C) signals.
	Chan chan os.Signal
	// Signals is the list of signals that cause the interrupt
	Signals = []os.Signal{os.Interrupt}
	// ShutdownRequestChan is a channel that can receive shutdown requests
	ShutdownRequestChan = qu.T()
	// AddHandlerChan is used to add an interrupt handler to the list of
	// handlers to be invoked on SIGINT (Ctrl+C) signals.
	AddHandlerChan = make(chan HandlerWithSource)
	// HandlersDone is closed after all interrupt handlers run the first time
	// an interrupt is signaled.
	HandlersDone = make(qu.C)
	DataDir      string
)

var interruptCallbacks []func()
var interruptCallbackSources []string

// Receiver listens for interrupt signals, registers interrupt callbacks, and responds to custom shutdown signals as
// required
func Listener() {
	invokeCallbacks := func() {
		Debug("running interrupt callbacks", len(interruptCallbacks), interruptCallbackSources)
		// run handlers in LIFO order.
		for i := range interruptCallbacks {
			idx := len(interruptCallbacks) - 1 - i
			Debug("running callback", idx, interruptCallbackSources[idx])
			interruptCallbacks[idx]()
		}
		Debug("interrupt handlers finished")
		HandlersDone.Q()
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
		// close(HandlersDone)
	}
out:
	for {
		select {
		case sig := <-Chan:
			// if !requested {
			// 	L.Printf("\r>>> received signal (%s)\n", sig)
			Debug("received interrupt signal", sig)
			requested.Store(true)
			invokeCallbacks()
			// pprof.Lookup("goroutine").WriteTo(os.Stderr, 2)
			// }
			break out
		case <-ShutdownRequestChan:
			// if !requested {
			Warn("received shutdown request - shutting down...")
			requested.Store(true)
			invokeCallbacks()
			break out
			// }
		case handler := <-AddHandlerChan:
			// if !requested {
			// Debug("adding handler")
			interruptCallbacks = append(interruptCallbacks, handler.Fn)
			interruptCallbackSources = append(interruptCallbackSources, handler.Source)
			// }
		case <-HandlersDone:
			break out
		}
	}
}

// AddHandler adds a handler to call when a SIGINT (Ctrl+C) is received.
func AddHandler(handler func()) {
	// Create the channel and start the main interrupt handler which invokes all other callbacks and exits if not
	// already done.
	_, loc, line, _ := runtime.Caller(1)
	msg := fmt.Sprintf("%s:%d", loc, line)
	Debug("handler added by:", msg)
	if Chan == nil {
		Chan = make(chan os.Signal)
		signal.Notify(Chan, Signals...)
		go Listener()
	}
	AddHandlerChan <- HandlerWithSource{
		msg, handler,
	}
}

// Request programmatically requests a shutdown
func Request() {
	_, f, l, _ := runtime.Caller(1)
	Debugf("interrupt requested %s:%d %v", f, l, requested)
	if requested.Load() {
		Debug("requested again")
		return
	}
	requested.Store(true)
	ShutdownRequestChan.Q()
	// qu.PrintChanState()
	var ok bool
	select {
	case _, ok = <-ShutdownRequestChan:
	default:
	}
	Debug("shutdownrequestchan", ok)
	if ok {
		close(ShutdownRequestChan)
	}
}

// GoroutineDump returns a string with the current goroutine dump in order to show what's going on in case of timeout.
func GoroutineDump() string {
	buf := make([]byte, 1<<18)
	n := runtime.Stack(buf, true)
	return string(buf[:n])
}

// RequestRestart sets the reset flag and requests a restart
func RequestRestart() {
	Restart = true
	Debug("requesting restart")
	Request()
}

// Requested returns true if an interrupt has been requested
func Requested() bool {
	return requested.Load()
}
