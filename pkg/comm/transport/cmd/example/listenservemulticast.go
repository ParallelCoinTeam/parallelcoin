package main

import (
	"fmt"
	"net"
	"time"
	
	qu "github.com/p9c/pod/pkg/util/quit"
	
	"github.com/p9c/pod/pkg/comm/transport"
	log "github.com/p9c/pod/pkg/util/logi"
	"github.com/p9c/pod/pkg/util/loop"
)

const (
	TestMagic = "TEST"
)

var (
	TestMagicB = []byte(TestMagic)
)

func main() {
	log.L.SetLevel("trace", true, "pod")
	dbg.Ln("starting test")
	quit := qu.T()
	var c *transport.Channel
	var e error
	if c, e = transport.NewBroadcastChannel("test", nil, "cipher",
		1234, 8192, transport.Handlers{
			TestMagic: func(ctx interface{}, src net.Addr, dst string,
				b []byte) (e error) {
				inf.F("%s <- %s [%d] '%s'", src.String(), dst, len(b), string(b))
				return
			},
		},
		quit,
	); dbg.Chk(e) {
		panic(err)
	}
	time.Sleep(time.Second)
	var n int
	loop.To(10, func(i int) {
		text := []byte(fmt.Sprintf("this is a test %d", i))
		inf.F("%s -> %s [%d] '%s'", c.Sender.LocalAddr(), c.Sender.RemoteAddr(), n-4, text)
		if e = c.SendMany(TestMagicB, transport.GetShards(text)); dbg.Chk(e) {
		} else {
		}
	})
	time.Sleep(time.Second * 5)
	if e = c.Close(); !dbg.Chk(e) {
		time.Sleep(time.Second * 1)
	}
	quit.Q()
}
