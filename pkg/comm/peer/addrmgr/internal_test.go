package addrmgr

import (
	"time"

	"github.com/p9c/pod/pkg/chain/wire"
)

func TstKnownAddressIsBad(ka *KnownAddress) bool {
	return ka.isBad()
}
func TstKnownAddressChance(ka *KnownAddress) float64 {
	return ka.chance()
}
func TstNewKnownAddress(na *wire.NetAddress, attempts int, lastAttempt, lastSuccess time.Time, tried bool, refs int) *KnownAddress {
	return &KnownAddress{na: na, attempts: attempts, lastAttempt: lastAttempt, lastSuccess: lastSuccess, tried: tried, refs: refs}
}
