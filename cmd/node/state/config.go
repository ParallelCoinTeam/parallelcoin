package state

import (
	"net"
	"time"

	chaincfg "git.parallelcoin.io/dev/pod/pkg/chain/config"
	"git.parallelcoin.io/dev/pod/pkg/discovery"
	"git.parallelcoin.io/dev/pod/pkg/util"
)

// Config stores current state of the node
type Config struct {
	Lookup              func(string) ([]net.IP, error)
	Oniondial           func(string, string, time.Duration) (net.Conn, error)
	Dial                func(string, string, time.Duration) (net.Conn, error)
	AddedCheckpoints    []chaincfg.Checkpoint
	ActiveMiningAddrs   []util.Address
	ActiveMinerKey      []byte
	ActiveMinRelayTxFee util.Amount
	ActiveWhitelists    []*net.IPNet
	DiscoveryUpdate     discovery.RequestFunc
	RouteableAddress    string
	DropAddrIndex       bool
	DropTxIndex         bool
	DropCfIndex         bool
	Save                bool
}
