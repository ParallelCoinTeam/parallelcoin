package main

import (
	"fmt"
	"net"
	"time"
	
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/loop"
	"github.com/p9c/pod/pkg/transport"
)

const (
	TestMagic = "TEST"
)

var (
	TestMagicB = []byte(TestMagic)
)

func main() {
	log.L.SetLevel("trace", true)
	if c, err := transport.NewBroadcastChannel("cipher",
		1234, 8192, transport.Handlers{
			TestMagic: func(src *net.UDPAddr, dst string, count int, data []byte) (err error) {
				log.INFOF("%s <- %s [%d] '%s'", src.String(), dst, count, string(data[:count]))
				return
			},
		},
	); log.Check(err) {
		panic(err)
	} else {
		var n int
		loop.To(10, func(i int) {
			text := []byte(fmt.Sprintf("this is a test %d", i))
			if n, err = c.Send(TestMagicB, text); log.Check(err) {
			} else {
				log.INFOF("%s -> %s [%d] '%s'",
					c.Sender.LocalAddr(), c.Sender.RemoteAddr(), n-4, text)
			}
		})
		time.Sleep(time.Second)
		if err = c.Close(); !log.Check(err) {
			time.Sleep(time.Second * 1)
		}
	}
}
