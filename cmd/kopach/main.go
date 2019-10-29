package main

import (
	"crypto/cipher"
	"github.com/p9c/pod/pkg/controller/broadcast"
	"github.com/p9c/pod/pkg/controller/gcm"
	"github.com/p9c/pod/pkg/log"
	"net"
	"sync"
	"time"

	"github.com/p9c/pod/pkg/conte"
)

func Main(cx *conte.Xt, quit chan struct{}, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		log.WARN("starting kopach standalone miner worker")
		m := newMsgHandle(*cx.Config.MinerPass)
	out:
		for {
			cancel := broadcast.Listen(broadcast.DefaultAddress, m.msgHandler)
			select {
			case <-quit:
				log.DEBUG("quitting on killswitch")
				cancel()
				break out
			}
		}
		wg.Done()
	}()
}
