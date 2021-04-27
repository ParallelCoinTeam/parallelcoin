// Package multicast provides a UDP multicast connection with an in-process multicast interface for sending and
// receiving.
//
// In order to allow processes on the same machine (windows) to receive the messages this code enables multicast
// loopback. It is up to the consuming library to discard messages it sends. This is only necessary because the net
// standard library disables loopback by default though on windows this takes effect whereas on unix platforms it does
// not.
//
// This code was derived from the information found here:
// https://stackoverflow.com/questions/43109552/how-to-set-ip-multicast-loop-on-multicast-udpconn-in-golang

package multicast

import (
	"net"
	
	"golang.org/x/net/ipv4"
	
	"github.com/p9c/matrjoska/pkg/util/routeable"
)

func Conn(port int) (conn *net.UDPConn, e error) {
	var ipv4Addr = &net.UDPAddr{IP: net.IPv4(224, 0, 0, 1), Port: port}
	if conn, e = net.ListenUDP("udp4", ipv4Addr); E.Chk(e) {
		return
	}
	D.Ln("listening on", conn.LocalAddr(), "to", conn.RemoteAddr())
	pc := ipv4.NewPacketConn(conn)
	// var ifaces []net.Interface
	var iface *net.Interface
	// if ifaces, e = net.Interfaces(); E.Chk(e) {
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
	ifcs, _ := routeable.GetAllInterfacesAndAddresses()
	for _, ifc := range ifcs {
		iface = ifc
		if e = pc.JoinGroup(iface, &net.UDPAddr{IP: net.IPv4(224, 0, 0, 1)}); E.Chk(e) {
			return
		}
		// test
		var loop bool
		if loop, e = pc.MulticastLoopback(); e == nil {
			D.F("MulticastLoopback status:%v\n", loop)
			if !loop {
				if e = pc.SetMulticastLoopback(true); e != nil {
					E.F("SetMulticastLoopback error:%v\n", e)
				}
			}
			D.Ln("Multicast Loopback enabled on", ifc.Name)
		}
	}
	return
}
