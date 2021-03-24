package state

import (
	"github.com/p9c/pod/pkg/amt"
	"github.com/p9c/pod/pkg/btcaddr"
	"net"
	"time"
	
	"github.com/p9c/pod/pkg/chaincfg"
)

// Config stores current state of the node
type Config struct {
	Lookup              func(string) ([]net.IP, error)
	Oniondial           func(string, string, time.Duration) (net.Conn, error)
	Dial                func(string, string, time.Duration) (net.Conn, error)
	AddedCheckpoints    []chaincfg.Checkpoint
	ActiveMiningAddrs   []btcaddr.Address
	ActiveMinerKey      []byte
	ActiveMinRelayTxFee amt.Amount
	ActiveWhitelists    []*net.IPNet
	DropAddrIndex       bool
	DropTxIndex         bool
	DropCfIndex         bool
	Save                bool
	// Miner               *worker.Worker
}
