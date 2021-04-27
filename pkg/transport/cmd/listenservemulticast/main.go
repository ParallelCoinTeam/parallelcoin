package main

import (
	"fmt"
	"github.com/p9c/log"
	"net"
	"time"
	
	"github.com/p9c/qu"
	
	"github.com/p9c/matrjoska/pkg/transport"
	"github.com/p9c/matrjoska/pkg/util/loop"
)

const (
	TestMagic = "TEST"
)

var (
	TestMagicB = []byte(TestMagic)
)

func main() {
	log.SetLogLevel("trace")
	D.Ln("starting test")
	quit := qu.T()
	var c *transport.Channel
	var e error
	if c, e = transport.NewBroadcastChannel(
		"test", nil, []byte("cipher"),
		1234, 8192, transport.Handlers{
			TestMagic: func(
				ctx interface{}, src net.Addr, dst string,
				b []byte,
			) (e error) {
				I.F("%s <- %s [%d] '%s'", src.String(), dst, len(b), string(b))
				return
			},
		},
		quit,
	); E.Chk(e) {
		panic(e)
	}
	time.Sleep(time.Second)
	var n int
	loop.To(
		10, func(i int) {
			text := []byte(fmt.Sprintf("this is a test %d", i))
			I.F("%s -> %s [%d] '%s'", c.Sender.LocalAddr(), c.Sender.RemoteAddr(), n-4, text)
			if e = c.SendMany(TestMagicB, transport.GetShards(text)); E.Chk(e) {
			} else {
			}
		},
	)
	time.Sleep(time.Second * 5)
	if e = c.Close(); !E.Chk(e) {
		time.Sleep(time.Second * 1)
	}
	quit.Q()
}
