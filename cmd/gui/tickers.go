package gui

import (
	"github.com/p9c/pod/pkg/conte"
	"time"
)

func DuOSgatherer(cx *conte.Xt) {
	rcv.cx = cx
	go func() {
		for _ = range time.NewTicker(time.Second * 1).C {
			rcv.GetDuOSbalance()
			rcv.GetDuOStransactions(0, 10, "all")
			rcv.GetDuOSblockCount()
			rcv.GetDuOSnetworkLastBlock()
			rcv.GetDuOSconnectionCount()
			rcv.GetDuOStatus()

		}
	}()

}
