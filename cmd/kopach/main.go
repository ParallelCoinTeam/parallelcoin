package kopach

import (
	"context"
	"github.com/p9c/pod/pkg/chain/mining"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/controller/transport"
	"github.com/p9c/pod/pkg/log"
	"go.uber.org/atomic"
	"net"
	"sync"
)

type Worker struct {
	active                 *atomic.Bool
	blockTemplateGenerator *mining.BlkTmplGenerator
	conn                   *transport.Connection
	ctx                    context.Context
	cx                     *conte.Xt
	mx                     *sync.Mutex
	receiveChan            chan []byte
	sendAddresses          []*net.UDPAddr
}

func Main(cx *conte.Xt, quit chan struct{}) {
	log.DEBUG("miner controller starting")
out:
	for {
		select {
		case <-quit:
			break out
		}
	}
}
