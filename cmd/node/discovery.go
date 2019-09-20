package node

import (
	"strings"
	"time"

	"github.com/grandcat/zeroconf"

	"git.parallelcoin.io/dev/pod/cmd/node/rpc"
	"git.parallelcoin.io/dev/pod/pkg/conte"
	"git.parallelcoin.io/dev/pod/pkg/discovery"
	"git.parallelcoin.io/dev/pod/pkg/peer/connmgr"
	"git.parallelcoin.io/dev/pod/pkg/util/cl"
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
		log <- cl.Error{"error running zeroconf search ", err, cl.Ine()}
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
						log <- cl.Warn{"connecting to zeroconf peer ",
							nodeAddress, cl.Ine()}
					}
				}
			case <-quit:
				cancelSearch()
				break out
			}
		}
	}()
	log <- cl.Warn{"started up discovery loop", cl.Ine()}
	return func() {}
}
