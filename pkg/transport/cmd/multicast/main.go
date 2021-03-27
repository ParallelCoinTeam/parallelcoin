package main

import (
	"fmt"
	"net"
	"time"

	"golang.org/x/net/ipv4"
)

var ipv4Addr = &net.UDPAddr{IP: net.IPv4(224, 0, 0, 1), Port: 1234}

func main() {
	conn, e := net.ListenUDP("udp4", ipv4Addr)
	if e != nil  {
		fmt.Printf("ListenUDP error %v\n", e)
		return
	}

	pc := ipv4.NewPacketConn(conn)
	var ifaces []net.Interface
	var iface net.Interface
	if ifaces, e = net.Interfaces(); E.Chk(e) {
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
	if e = pc.JoinGroup(&iface, &net.UDPAddr{IP: net.IPv4(224, 0, 0, 1)}); E.Chk(e) {
		return
	}
	// test
	var loop bool
	if loop, e = pc.MulticastLoopback(); e ==  nil {
		fmt.Printf("MulticastLoopback status:%v\n", loop)
		if !loop {
			if e = pc.SetMulticastLoopback(true); E.Chk(e) {
				fmt.Printf("SetMulticastLoopback error:%v\n", e)
			}
		}
	}
	go func() {
		for {
			if _, e = conn.WriteTo([]byte("hello"), ipv4Addr); E.Chk(e) {
				fmt.Printf("Write failed, %v\n", e)
			}
			time.Sleep(time.Second)
		}
	}()

	buf := make([]byte, 1024)
	for {
		if n, addr, e := conn.ReadFrom(buf); E.Chk(e) {
			fmt.Printf("error %v", e)
		} else {
			fmt.Printf("recv %s from %v\n", string(buf[:n]), addr)
		}
	}

	// return
}

//
// func main() {
// 	var ifs []net.Interface
// 	var e error
// 	if ifs, e = net.Interfaces(); E.Chk(e) {
// 	}
// 	D.S(ifs)
// 	var addrs []net.Addr
// 	var addr net.Addr
// 	for i := range ifs {
// 		if ifs[i].Flags&net.FlagUp != 0 && ifs[i].Flags&net.FlagMulticast != 0 {
// 			if addrs, e = ifs[i].MulticastAddrs(); E.Chk(e) {
// 			}
// 			for j := range addrs {
// 				if addrs[j].String() == Multicast {
// 					addr = addrs[j]
// 					break
// 				}
// 			}
// 		}
// 	}
// 	D.S(addr)
// }
