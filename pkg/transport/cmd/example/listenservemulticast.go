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
	if c, err := transport.NewBroadcastChannel(nil, "cipher",
		1234, 8192, transport.Handlers{
			TestMagic: func(ctx interface{}, src *net.UDPAddr, dst string,
				b []byte) (err error) {
				log.INFOF("%s <- %s [%d] '%s'", src.String(), dst, len(b), string(b))
				return
			},
		},
	); log.Check(err) {
		panic(err)
	} else {
		var n int
		loop.To(10, func(i int) {
			text := []byte(fmt.Sprintf("this is a test %d", i))
			if err = c.SendMany(TestMagicB, transport.GetShards(text)); log.Check(err) {
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
