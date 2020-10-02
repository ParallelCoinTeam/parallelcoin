// Package multicast provides a UDP multicast connection with an in-process multicast interface for sending and receiving.
//
// In order to allow processes on the same machine (windows) to receive the messages this code enables multicast
// loopback. It is up to the consuming library to discard messages it sends. This is only necessary because the
// net standard library disables loopback by default though on windows this takes effect whereas on unix platforms
// it does not.

package multicast

import (
	"net"

	"golang.org/x/net/ipv4"
)

func GetMulticastConn(port int) (conn *net.UDPConn, err error) {
	var ipv4Addr = &net.UDPAddr{IP: net.IPv4(224, 0, 0, 1), Port: port}
	conn, err = net.ListenUDP("udp4", ipv4Addr)
	if err != nil {
		Errorf("ListenUDP error %v\n", err)
		return
	}

	pc := ipv4.NewPacketConn(conn)
	var ifaces []net.Interface
	var iface net.Interface
	if ifaces, err = net.Interfaces(); Check(err) {
	}
	// This grabs the first physical interface with multicast that is up
	for i := range ifaces {
		if ifaces[i].Flags&net.FlagMulticast != 0 &&
			ifaces[i].Flags&net.FlagUp != 0 &&
			ifaces[i].HardwareAddr != nil {
			iface = ifaces[i]
			break
		}
	}
	if err = pc.JoinGroup(&iface, &net.UDPAddr{IP: net.IPv4(224, 0, 0, 1)}); Check(err) {
		return
	}
	// test
	if loop, err := pc.MulticastLoopback(); err == nil {
		Debugf("MulticastLoopback status:%v\n", loop)
		if !loop {
			if err := pc.SetMulticastLoopback(true); err != nil {
				Errorf("SetMulticastLoopback error:%v\n", err)
			}
		}
	}

	return
}