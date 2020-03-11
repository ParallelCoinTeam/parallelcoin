package state

import (
	"net"
	"time"

	chaincfg "github.com/p9c/chaincfg"
	"github.com/p9c/util"
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
	DropAddrIndex       bool
	DropTxIndex         bool
	DropCfIndex         bool
	Save                bool
}
