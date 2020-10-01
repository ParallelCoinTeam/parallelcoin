// +build !windows,!plan9

package rename

import (
	"os"
)

// Atomic provides an atomic file rename. newpath is replaced if it already exists.
func Atomic(oldpath, newpath string) error {
	return os.Rename(oldpath, newpath)
}
