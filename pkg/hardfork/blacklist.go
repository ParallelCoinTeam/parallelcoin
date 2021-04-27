package hardfork

import (
	"github.com/p9c/matrjoska/pkg/btcaddr"
)

// Blacklist is a list of addresses that have been suspended
var Blacklist = []btcaddr.Address{
	// Cryptopia liquidation wallet
	// Addr("8JEEhaMxJf4dZh5rvVCVSA7JKeYBvy8fir", mn),
}
