package main

import (
	"fmt"
	"net"
	"time"
	
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/transport"
)

func main() {
	log.L.SetLevel("trace", true)
	ready := make(chan struct{})
	var listener *net.UDPConn
	var err error
	transport.Address = "224.0.0.1:1234"
	go func() {
		listener, err = transport.Listen(8192,
			func(addr *net.UDPAddr, count int, data []byte) {
				log.INFOF("%s [%d] received '%s'", addr, count, string(data[:count]))
			})
		ready <- struct{}{}
	}()
	<-ready
	bc, err := transport.NewBroadcaster(8192)
	if err != nil {
		log.ERROR(err)
	}
	var n int
	for i := 0; i < 10; i++ {
		text := fmt.Sprintf("this is a test %d", i)
		n, err = bc.Write([]byte(text))
		if err != nil {
			log.ERROR(err)
		}
		log.INFOF("%s sent %d '%s'", bc.RemoteAddr(), n, text)
	}
	err = bc.Close()
	if err != nil {
		log.ERROR(err)
	}
	time.Sleep(time.Second * 1)
	err = listener.Close()
	if err != nil {
		log.ERROR(err)
	}
	time.Sleep(time.Second * 1)
}
