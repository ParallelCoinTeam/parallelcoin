package main

import (
	"fmt"
	"net"
	"time"

	log "github.com/p9c/pod/pkg/logi"
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
	L.SetLevel("trace", true, "pod")
	quit := make(chan struct{})
	if c, err := transport.NewBroadcastChannel("test", nil, "cipher",
		1234, 8192, transport.Handlers{
			TestMagic: func(ctx interface{}, src net.Addr, dst string,
				b []byte) (err error) {
				L.Infof("%s <- %s [%d] '%s'", src.String(), dst, len(b), string(b))
				return
			},
		},
		quit,
	); L.Check(err) {
		panic(err)
	} else {
		var n int
		loop.To(10, func(i int) {
			text := []byte(fmt.Sprintf("this is a test %d", i))
			if err = c.SendMany(TestMagicB, transport.GetShards(text)); L.Check(err) {
			} else {
				L.Infof("%s -> %s [%d] '%s'",
					c.Sender.LocalAddr(), c.Sender.RemoteAddr(), n-4, text)
			}
		})
		close(quit)
		if err = c.Close(); !L.Check(err) {
			time.Sleep(time.Second * 1)
		}
	}
}
