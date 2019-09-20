package state

import (
	"net"
	"time"

	chaincfg "github.com/parallelcointeam/parallelcoin/pkg/chain/config"
	"github.com/parallelcointeam/parallelcoin/pkg/discovery"
	"github.com/parallelcointeam/parallelcoin/pkg/util"
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
