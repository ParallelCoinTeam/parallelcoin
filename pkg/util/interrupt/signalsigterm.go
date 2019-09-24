// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package interrupt

import (
	"os"
	"syscall"
)

func init() {
	InterruptSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}
}
