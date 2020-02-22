// +build windows

package transport

import (
	"syscall"
	
	"github.com/p9c/pod/pkg/log"
)

func reusePort(network, address string, conn syscall.RawConn) error {
	return conn.Control(func(descriptor uintptr) {
		err := syscall.SetsockoptInt(syscall.Handle(descriptor), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
		if err != nil {
			log.ERROR(err)
		}
	})
}
