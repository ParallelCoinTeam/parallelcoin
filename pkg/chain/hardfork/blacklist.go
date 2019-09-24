package hardfork

import "github.com/parallelcointeam/parallelcoin/pkg/util"

// Blacklist is a list of addresses that have been suspended
var Blacklist = []util.Address{
	// Cryptopia liquidation wallet
	Addr("8JEEhaMxJf4dZh5rvVCVSA7JKeYBvy8fir", mn),
}
