// +build !windows,!plan9

package limits

import (
	"fmt"
	"github.com/p9c/pkg/app/slog"
	"syscall"
)

const (
	fileLimitWant = 32768
	fileLimitMin  = 1024
)

// SetLimits raises some process limits to values which allow pod and associated utilities to run.
func SetLimits() (err error) {
	var rLimit syscall.Rlimit
	if err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); slog.Check(err) {
		return
	}
	if rLimit.Cur > fileLimitWant {
		return
	}
	if rLimit.Max < fileLimitMin {
		err = fmt.Errorf("need at least %v file descriptors", fileLimitMin)
		slog.Debug(err)
		return
	}
	if rLimit.Max < fileLimitWant {
		rLimit.Cur = rLimit.Max
	} else {
		rLimit.Cur = fileLimitWant
	}
	if err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); slog.Check(err) {
		// try min value
		rLimit.Cur = fileLimitMin
		if err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); slog.Check(err) {
			return
		}
	}
	return
}
