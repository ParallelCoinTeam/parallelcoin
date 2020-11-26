// Package multicast provides a UDP multicast connection with an in-process multicast interface for sending and receiving.
//
// In order to allow processes on the same machine (windows) to receive the messages this code enables multicast
// loopback. It is up to the consuming library to discard messages it sends. This is only necessary because the
// net standard library disables loopback by default though on windows this takes effect whereas on unix platforms
// it does not.
//
// This code was derived from the information found here:
// https://stackoverflow.com/questions/43109552/how-to-set-ip-multicast-loop-on-multicast-udpconn-in-golang

package multicast

import (
	"net"

	"golang.org/x/net/ipv4"

	"github.com/p9c/pod/pkg/util/routeable"
)

func Conn(port int) (conn *net.UDPConn, err error) {
	var ipv4Addr = &net.UDPAddr{IP: net.IPv4(224, 0, 0, 1), Port: port}
	if conn, err = net.ListenUDP("udp4", ipv4Addr); Check(err) {
		return
	}

	pc := ipv4.NewPacketConn(conn)
	// var ifaces []net.Interface
	var iface *net.Interface
	// if ifaces, err = net.Interfaces(); Check(err) {
	// }
	// // This grabs the first physical interface with multicast that is up. Note that this should filter out
	// // VPN connections which would normally be selected first but don't actually have a multicast connection
	// // to the local area network.
	// for i := range ifaces {
	// 	if ifaces[i].Flags&net.FlagMulticast != 0 &&
	// 		ifaces[i].Flags&net.FlagUp != 0 &&
	// 		ifaces[i].HardwareAddr != nil {
	// 		iface = ifaces[i]
	// 		break
	// 	}
	// }
	iface = routeable.GetInterface()[0]
	if err = pc.JoinGroup(iface, &net.UDPAddr{IP: net.IPv4(224, 0, 0, 1)}); Check(err) {
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
