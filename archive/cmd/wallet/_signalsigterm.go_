// +build darwin dragonfly freebsd linux netbsd openbsd solaris
package wallet

import (
	"os"
	"syscall"
)

func init(	) {

	signals = []os.Signal{os.Interrupt, syscall.SIGTERM}
}
