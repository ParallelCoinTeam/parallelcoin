package kopach

import (
	"context"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/controller"
	"github.com/p9c/pod/pkg/log"
	"net"
	"sync"
)

func Main(cx *conte.Xt, quit chan struct{}, wg *sync.WaitGroup) {
	log.DEBUG("kopach miner starting")
	wg.Add(1)
	//port := controller.GetPort(*cx.Config.Controller)
	//controller.MCAddresses[0].Port = int(port.(*controller.Port).Get())
	//controller.MCAddresses[1].Port = int(port.(*controller.Port).Get())
	var cancel context.CancelFunc
	var err error
	for _, j := range controller.MCAddresses {
		i := j
		cancel, err = controller.Listen(i, func(a *net.UDPAddr, n int,
			b []byte) {
			log.DEBUG(i, a)
			log.SPEW(b[:n])
		})
		if err != nil {
			continue
		}
		if cancel != nil {
			log.DEBUG("listener started", i.IP, i.Port, i.Zone, i.String(),
				i.Network())
			break
		}
	}
	select {
	case <-quit:
		log.DEBUG("kopach miner shutting down")
		cancel()
		break
	}
	wg.Done()
}
