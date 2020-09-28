package connmgr

import (
	"fmt"
	"github.com/stalker-loki/app/slog"
	mrand "math/rand"
	"net"
	"strconv"
	"time"

	"github.com/p9c/pod/pkg/chain/config/netparams"
	"github.com/p9c/pod/pkg/chain/wire"
)

const (
	// These constants are used by the DNS seed code to pick a random last
	// seen time.
	secondsIn3Days int32 = 24 * 60 * 60 * 3
	secondsIn4Days int32 = 24 * 60 * 60 * 4
)

// OnSeed is the signature of the callback function which is invoked when DNS
// seeding is successful
type OnSeed func(addrs []*wire.NetAddress)

// LookupFunc is the signature of the DNS lookup function.
type LookupFunc func(string) ([]net.IP, error)

// SeedFromDNS uses DNS seeding to populate the address manager with peers.
func SeedFromDNS(chainParams *netparams.Params, reqServices wire.ServiceFlag,
	lookupFn LookupFunc, seedFn OnSeed) {
	for _, dnsSeed := range chainParams.DNSSeeds {
		var host string
		if !dnsSeed.HasFiltering || reqServices == wire.SFNodeNetwork {
			host = dnsSeed.Host
		} else {
			host = fmt.Sprintf("x%x.%s", uint64(reqServices), dnsSeed.Host)
		}
		go func(host string) {
			randSource := mrand.New(mrand.NewSource(time.Now().UnixNano()))
			seedPeers, err := lookupFn(host)
			if err != nil {
				slog.Error("DNS routeable failed on seed %s: %v", host, err)
				return
			}
			numPeers := len(seedPeers)
			slog.Debugf("%d addresses found from DNS seed %s", numPeers, host)
			if numPeers == 0 {
				return
			}
			addresses := make([]*wire.NetAddress, len(seedPeers))
			// if this errors then we have *real* problems
			intPort, _ := strconv.Atoi(chainParams.DefaultPort)
			for i, peer := range seedPeers {
				addresses[i] = wire.NewNetAddressTimestamp(
					// bitcoind seeds with addresses from a time randomly
					// selected between 3 and 7 days ago.
					time.Now().Add(-1*time.Second*time.Duration(secondsIn3Days+
						randSource.Int31n(secondsIn4Days))),
					0, peer, uint16(intPort))
			}
			seedFn(addresses)
		}(host)
	}
}
