package rcd

import (
	"os"
	"sync/atomic"
	
	"github.com/p9c/pod/cmd/node"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/util/interrupt"
)

func (r *RcVar) DuoUInode() error {
	r.cx.NodeKill = make(chan struct{})
	r.cx.Node = &atomic.Value{}
	r.cx.Node.Store(false)
	var err error
	if !*r.cx.Config.NodeOff {
		go func() {
			log.INFO(r.cx.Language.RenderText("goApp_STARTINGNODE"))
			// utils.GetBiosMessage(view, cx.Language.RenderText("goApp_STARTINGNODE"))
			err = node.Main(r.cx, nil)
			if err != nil {
				log.INFO("error running node:", err)
				os.Exit(1)
			}
		}()
		
	}
	interrupt.AddHandler(func() {
		log.WARN("interrupt received, " +
			"shutting down node")
		close(r.cx.NodeKill)
	})
	return err
}
