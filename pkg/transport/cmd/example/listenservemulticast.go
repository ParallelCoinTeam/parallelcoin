package main

import (
	"fmt"
	"net"
	"time"

	"github.com/p9c/pod/pkg/log"
)

const UDP4MulticastAddress = "224.0.0.1:11049"

func main() {
	ready := make(chan struct{})
	go func() {
		log.INFO("listening")
		ready <- struct{}{}
		Listen(UDP4MulticastAddress, func(addr *net.UDPAddr, count int, data []byte) {
			log.INFO(addr, count, string(data[:count]))
		})
	}()
	bc, err := NewBroadcaster(UDP4MulticastAddress)
	if err != nil {
		log.ERROR(err)
	}
	<-ready
	for i := 0; i < 10; i++ {
		log.INFO("sending", i)
		// var n int
		_, err = bc.Write([]byte(fmt.Sprintf("this is a test %d", i)))
		// log.INFO(n, err)
		if err != nil {
			log.ERROR(err)
		}
	}
	time.Sleep(time.Second * 1)
}

const (
	maxDatagramSize = 8192
)

// NewBroadcaster creates a new UDP multicast connection on which to broadcast
func NewBroadcaster(address string) (*net.UDPConn, error) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// Listen binds to the UDP address and port given and writes packets received
// from that address to a buffer which is passed to a hander
func Listen(address string, handler func(*net.UDPAddr, int, []byte)) {
	// Parse the string address
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.ERROR(err)
	}

	// Open up a connection
	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		log.ERROR(err)
	}

	conn.SetReadBuffer(maxDatagramSize)

	// Loop forever reading from the socket
	for {
		buffer := make([]byte, maxDatagramSize)
		numBytes, src, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.ERROR("ReadFromUDP failed:", err)
		}

		handler(src, numBytes, buffer)
	}
}
