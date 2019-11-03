package netparams

import (
	"sync"
	
	chaincfg `github.com/p9c/pod/pkg/chain/config`
)

// Params is used to group parameters for various networks such as the main network and test networks.
type Params struct {
	sync.Mutex
	*chaincfg.Params
	RPCClientPort       string
	WalletRPCServerPort string
}

// MainNetParams contains parameters specific running btcwallet and pod on the main network (wire.MainNet).
var MainNetParams = Params{
	Params:              &chaincfg.MainNetParams,
	RPCClientPort:       "11048",
	WalletRPCServerPort: "11046",
}

// SimNetParams contains parameters specific to the simulation test network (wire.SimNet).
var SimNetParams = Params{
	Params:              &chaincfg.SimNetParams,
	RPCClientPort:       "41048",
	WalletRPCServerPort: "41046",
}

// TestNet3Params contains parameters specific running btcwallet and pod on the test network (version 3) (wire.TestNet3).
var TestNet3Params = Params{
	Params:              &chaincfg.TestNet3Params,
	RPCClientPort:       "21048",
	WalletRPCServerPort: "21046",
}

// RegressionTestParams contains parameters specific to the simulation test network (wire.SimNet).
var RegressionTestParams = Params{
	Params:              &chaincfg.RegressionTestParams,
	RPCClientPort:       "31048",
	WalletRPCServerPort: "31046",
}
