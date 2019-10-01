package app

import (
	"sync"

	"github.com/urfave/cli"

	"github.com/p9c/pod/cmd/kopach"
	"github.com/p9c/pod/pkg/conte"
)

func kopachHandle(cx *conte.Xt) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		var wg sync.WaitGroup
		Configure(cx)
		quit := make(chan struct{})
		kopach.Main(cx, quit, &wg)
		wg.Wait()
		return
	}
}
