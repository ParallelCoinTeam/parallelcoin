package main

import (
	"fmt"
	"net"
	"time"

	"golang.org/x/net/ipv4"
)

var ipv4Addr = &net.UDPAddr{IP: net.IPv4(224, 0, 0, 1), Port: 1234}

func main() {
	conn, err := net.ListenUDP("udp4", ipv4Addr)
	if err != nil {
		fmt.Printf("ListenUDP error %v\n", err)
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
		fmt.Printf("MulticastLoopback status:%v\n", loop)
		if !loop {
			if err := pc.SetMulticastLoopback(true); err != nil {
				fmt.Printf("SetMulticastLoopback error:%v\n", err)
			}
		}
	}
	go func() {
		for {
			if _, err := conn.WriteTo([]byte("hello"), ipv4Addr); err != nil {
				fmt.Printf("Write failed, %v\n", err)
			}
			time.Sleep(time.Second)
		}
	}()

	buf := make([]byte, 1024)
	for {
		if n, addr, err := conn.ReadFrom(buf); err != nil {
			fmt.Printf("error %v", err)
		} else {
			fmt.Printf("recv %s from %v\n", string(buf[:n]), addr)
		}
	}

	// return
}

//
// func main() {
// 	var ifs []net.Interface
// 	var err error
// 	if ifs, err = net.Interfaces(); Check(err) {
// 	}
// 	Debugs(ifs)
// 	var addrs []net.Addr
// 	var addr net.Addr
// 	for i := range ifs {
// 		if ifs[i].Flags&net.FlagUp != 0 && ifs[i].Flags&net.FlagMulticast != 0 {
// 			if addrs, err = ifs[i].MulticastAddrs(); Check(err) {
// 			}
// 			for j := range addrs {
// 				if addrs[j].String() == Multicast {
// 					addr = addrs[j]
// 					break
// 				}
// 			}
// 		}
// 	}
// 	Debugs(addr)
// }
