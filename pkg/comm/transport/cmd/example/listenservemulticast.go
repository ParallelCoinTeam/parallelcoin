package main

import (
	"fmt"
	"github.com/p9c/pkg/app/slog"
	"net"
	"time"

	"github.com/p9c/pod/pkg/comm/transport"
)

const (
	TestMagic = "TEST"
)

var (
	TestMagicB = []byte(TestMagic)
)

func main() {
	quit := make(chan struct{})
	if c, err := transport.NewBroadcastChannel("test", nil, "cipher",
		1234, 8192, transport.Handlers{
			TestMagic: func(ctx interface{}, src net.Addr, dst string,
				b []byte) (err error) {
				slog.Infof("%s <- %s [%d] '%s'", src.String(), dst, len(b), string(b))
				return
			},
		},
		quit,
	); slog.Check(err) {
		panic(err)
	} else {
		var n int
		for i := 0; i < 10; i++ {
			text := []byte(fmt.Sprintf("this is a test %d", i))
			if err = c.SendMany(TestMagicB, transport.GetShards(text)); slog.Check(err) {
			} else {
				slog.Infof("%s -> %s [%d] '%s'",
					c.Sender.LocalAddr(), c.Sender.RemoteAddr(), n-4, text)
			}
		}
		close(quit)
		if err = c.Close(); !slog.Check(err) {
			time.Sleep(time.Second * 1)
		}
	}
}
