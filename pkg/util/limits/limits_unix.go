// +build !windows,!plan9

package limits

import (
	"fmt"
	"syscall"
)

const (
	fileLimitWant = 32768
	fileLimitMin  = 1024
)

// SetLimits raises some process limits to values which allow pod and associated utilities to run.
func SetLimits() (e error) {
	var rLimit syscall.Rlimit
	if e = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); E.Chk(e) {
		return
	}
	if rLimit.Cur > fileLimitWant {
		return
	}
	if rLimit.Max < fileLimitMin {
		e = fmt.Errorf(
			"need at least %v file descriptors",
			fileLimitMin,
		)
		return
	}
	if rLimit.Max < fileLimitWant {
		rLimit.Cur = rLimit.Max
	} else {
		rLimit.Cur = fileLimitWant
	}
	e = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if e != nil {
		E.Ln(e)
		// try min value
		rLimit.Cur = fileLimitMin
		if e = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); E.Chk(e) {
			return
		}
	}
	return
}
