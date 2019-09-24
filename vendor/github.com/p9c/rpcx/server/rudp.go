package server

import (
	"net"

	"git.parallelcoin.io/dev/rudp"
)

func init() {
	makeListeners["rudp"] = rudpMakeListener
}

func rudpMakeListener(s *Server, address string) (ln net.Listener, err error) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}

	l, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	return rudp.NewListener(l), nil
}
