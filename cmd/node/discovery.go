package node

import (
	"strings"
	"time"

	"github.com/grandcat/zeroconf"

	"github.com/parallelcointeam/parallelcoin/cmd/node/rpc"
	"github.com/parallelcointeam/parallelcoin/pkg/conte"
	"github.com/parallelcointeam/parallelcoin/pkg/discovery"
	"github.com/parallelcointeam/parallelcoin/pkg/log"
	"github.com/parallelcointeam/parallelcoin/pkg/peer/connmgr"
)

// DiscoverPeers connects lan peers in the same group.
// This allows faster initial sync for new peers and instant creation
// of single machine or lan testnets and running more than one isolated
// network of the same kind on a LAN also without configuration.
// A function is
func DiscoverPeers(cx *conte.Xt) (cancel func()) {
	quit := make(chan struct{})
	cancel = func() { close(quit) }
	serviceName := discovery.GetParallelcoinServiceName(cx.ActiveNet)
	cancelSearch, resultsChan, err := discovery.AsyncZeroConfSearch(
		serviceName, *cx.Config.Group)
	if err != nil {
		log.ERROR("error running zeroconf search ", err)
		return func() {}
	}
	ticker := time.NewTicker(time.Second * 10)
	var zcPeers []*zeroconf.ServiceEntry
	go func() {
	out:
		for {
		selectOut:
			select {
			case <-ticker.C:
				// every 10 seconds we clear this - it is just to stop the
				// multiple rebroadcasts from each peer that is well stopped
				// by 10 seconds
				zcPeers = []*zeroconf.ServiceEntry{}
			case r := <-resultsChan:
				for i := range zcPeers {
					if r.Instance == zcPeers[i].Instance {
						break selectOut
					}
				}
				for i := range r.Text {
					split := strings.Split(r.Text[i], "=")
					if split[0] == "node" {
						nodeAddress, err := rpc.AddrStringToNetAddr(cx.Config,
							cx.StateCfg, split[1])
						if err != nil {
							continue
						}
						cx.RealNode.ConnManager.Connect(&connmgr.ConnReq{
							Addr:      nodeAddress,
							Permanent: true,
						})
						zcPeers = append(zcPeers, r)
						log.WARN("connecting to zeroconf peer ", nodeAddress)
					}
				}
			case <-quit:
				cancelSearch()
				break out
			}
		}
	}()
	// log.TRACE("started up discovery loop")
	return cancel
}
