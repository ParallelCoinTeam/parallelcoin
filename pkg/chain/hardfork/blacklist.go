package hardfork

import "git.parallelcoin.io/dev/pod/pkg/util"

// Blacklist is a list of addresses that have been suspended
var Blacklist = []util.Address{
	// Cryptopia liquidation wallet
	Addr("8JEEhaMxJf4dZh5rvVCVSA7JKeYBvy8fir", mn),
}
