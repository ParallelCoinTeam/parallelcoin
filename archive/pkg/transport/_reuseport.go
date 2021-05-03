// +build !windows

package transport

import (
	"syscall"
)

func reusePort(network, address string, conn syscall.RawConn) (e error) {
	return conn.Control(func(descriptor uintptr) {
		e := syscall.SetsockoptInt(int(descriptor), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
		if e != nil {
		}
	},
	)
}
