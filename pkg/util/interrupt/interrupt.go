package interrupt

import (
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	
	"github.com/kardianos/osext"
	
	"github.com/p9c/pod/app/appdata"
	"github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/pkg/log"
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
		log.DEBUG("interrupt handlers finished")
		if Restart {
			log.DEBUG("restarting")
			file, err := osext.Executable()
			if err != nil {
				log.ERROR(err)
				return
			}
			err = syscall.Exec(file, os.Args, os.Environ())
			if err != nil {
				log.FATAL(err)
			}
			// return
			os.Exit(1)
		} else {
			// return
			os.Exit(1)
		}
	}
	for {
		select {
		case sig := <-Chan:
			log.Printf("\r>>> received signal (%s)\n", sig)
			log.DEBUG("received interrupt signal")
			requested = true
			invokeCallbacks()
			return
		case <-ShutdownRequestChan:
			log.WARN("received shutdown request - shutting down...")
			requested = true
			invokeCallbacks()
			return
		case handler := <-AddHandlerChan:
			log.DEBUG("adding handler")
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

// RequestRestart sets the reset flag and requests a restart
func RequestRestart() {
	Restart = true
	Request()
}

// Requested returns true if an interrupt has been requested
func Requested() bool {
	return requested
}

// cleanAndExpandPath expands environement variables and leading ~ in the passed path, cleans the result, and returns it.
func cleanAndExpandPath(path string) string {
	// Expand initial ~ to OS specific home directory.
	if strings.HasPrefix(path, "~") {
		appHomeDir := appdata.Dir("gencerts", false)
		homeDir := filepath.Dir(appHomeDir)
		path = strings.Replace(path, "~", homeDir, 1)
	}
	if !apputil.FileExists(path) {
		wd, err := os.Getwd()
		if err != nil {
			log.ERROR("can't get working dir:", err)
		}
		path = filepath.Join(wd, path)
	}
	// NOTE: The os.ExpandEnv doesn't work with Windows-style %VARIABLE%, but they variables can still be expanded via POSIX-style $VARIABLE.
	return filepath.Clean(os.ExpandEnv(path))
}
