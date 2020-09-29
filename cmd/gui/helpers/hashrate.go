package helpers

import (
	"container/ring"
	"time"

	"github.com/p9c/pod/cmd/kopach/control/hashrate"
)

// GetHashrate returns the exponential weighted moving average of the total hashrate and
// a simple moving average for each block version. To decode the caller needs to use
// fork.GetAlgoName(height, version) which returns the string name of the algorithm/block version
// fork.List[<current hard fork id>].Algos[<block version number>].VersionInterval tells the number of
// seconds for this block version interval and the coinbase payment is scaled according to this ratio,
// which is computed by
func GetHashrate(hrb *ring.Ring) (hr float64, hrp map[int32]float64) {
	var hashTotal int
	hashPerVersion := make(map[int32]int)
	hrp = make(map[int32]float64)
	var firstHashTime, lastHashTime time.Time
	var started bool
	hrb.Do(func(entry interface{}) {
		e, ok := entry.(hashrate.Hashrate)
		// Debug("iterating hashrate buffer", entry)
		if ok {
			Debug("got entry in hashrate buffer")
			if !started {
				started = true
				firstHashTime = e.Time
			}
			hashTotal += e.Count
			hashPerVersion[e.Version] += e.Count
			lastHashTime = e.Time
		}
	})
	Debug(hashTotal, hashPerVersion, firstHashTime, lastHashTime)
	hashDuration := lastHashTime.Sub(firstHashTime)
	hr = float64(hashDuration) / float64(hashTotal)
	for i := range hashPerVersion {
		hrp[i] = float64(hashDuration) / float64(hashPerVersion[i])
	}
	return
}
