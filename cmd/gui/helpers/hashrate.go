package helpers

import (
	"time"
	
	"github.com/VividCortex/ewma"
	
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/controller/hashrate"
)

// GetHashrate returns the exponential weighted moving average of the total hashrate and
// a simple moving average for each block version. To decode the caller needs to use
// fork.GetAlgoName(height, version) which returns the string name of the algorithm/block version
// fork.List[<current hard fork id>].Algos[<block version number>].VersionInterval tells the number of
// seconds for this block version interval and the coinbase payment is scaled according to this ratio,
// which is computed by
func GetHashrate(cx *conte.Xt) (hr float64, hrp map[int32]float64) {
	var hashTotal int
	hashPerVersion := make(map[int32]int)
	hrp = make(map[int32]float64)
	var firstHashTime, lastHashTime time.Time
	var started bool
	ma := ewma.NewMovingAverage()
	cx.HRBMutex.Lock()
	cx.HashrateBuffer.Do(func(entry interface{}) {
		e := entry.(hashrate.Hashrate)
		if !started {
			started = true
			firstHashTime = e.Time
		}
		hashTotal += e.Count
		hashPerVersion[e.Version] += e.Count
		lastHashTime = e.Time
		ma.Add(float64(e.Count))
	})
	cx.HRBMutex.Unlock()
	hashDuration := lastHashTime.Sub(firstHashTime)
	for i, v := range hashPerVersion {
		hrp[i] = float64(v) / float64(hashDuration)
	}
	return
}
