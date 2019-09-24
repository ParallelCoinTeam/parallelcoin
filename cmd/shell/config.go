package shell

import (
	"github.com/parallelcointeam/parallelcoin/cmd/node"
	"github.com/parallelcointeam/parallelcoin/cmd/walletmain"
	"github.com/parallelcointeam/parallelcoin/pkg/chain/config/netparams"
)

// Config is the combined app and logging configuration data
type Config struct {
	ConfigFile      string
	DataDir         string
	AppDataDir      string
	Node            *node.Config
	Wallet          *walletmain.Config
	Levels          map[string]string
	nodeActiveNet   *netparams.Params
	walletActiveNet *netparams.Params
}

// GetNodeActiveNet returns the activenet netparams
func (r *Config) GetNodeActiveNet() *netparams.Params {
	return r.nodeActiveNet
}

// GetWalletActiveNet returns the activenet netparams
func (r *Config) GetWalletActiveNet() *netparams.Params {
	return r.walletActiveNet
}

// SetNodeActiveNet returns the activenet netparams
func (r *Config) SetNodeActiveNet(in *netparams.Params) {
	r.nodeActiveNet = in
}

// SetWalletActiveNet returns the activenet netparams
func (r *Config) SetWalletActiveNet(in *netparams.Params) {
	r.walletActiveNet = in
}
