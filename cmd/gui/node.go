package gui

import (
	"github.com/p9c/pod/cmd/node"
	"github.com/p9c/pod/cmd/node/rpc"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/util/interrupt"
	"os"
	"sync"
	"sync/atomic"

)

func DuOSnode(cx *conte.Xt) error {
	nodeChan := make(chan *rpc.Server)
	cx.NodeKill = make(chan struct{})
	cx.Node = &atomic.Value{}
	cx.Node.Store(false)
	var err error
	var wg sync.WaitGroup
	if !*cx.Config.NodeOff {
		go func() {
			log.INFO(cx.Language.RenderText("goApp_STARTINGNODE"))
			//utils.GetBiosMessage(view, cx.Language.RenderText("goApp_STARTINGNODE"))
			err = node.Main(cx, nil, cx.NodeKill, nodeChan, &wg)
			if err != nil {
				log.INFO("error running node:", err)
				os.Exit(1)
			}
		}()
		log.DEBUG("waiting for nodeChan")
		cx.RPCServer = <-nodeChan
		log.DEBUG("nodeChan sent")
		cx.Node.Store(true)
	}
	interrupt.AddHandler(func() {
		log.WARN("interrupt received, " +
			"shutting down shell modules")
		close(cx.NodeKill)
	})
	return err
}



