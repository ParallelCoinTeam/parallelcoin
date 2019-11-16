package blockchain

import (
	"github.com/VividCortex/ewma"
	"github.com/p9c/pod/pkg/chain/fork"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/log"
)

func (b *BlockChain) GetCommonP9Averages(lastNode *BlockNode,
	nH int32) (allTimeAv, allTimeDiv, qhourDiv, hourDiv, dayDiv float64) {
	const minAvSamples = 2
	allTimeAv, allTimeDiv, qhourDiv, hourDiv, dayDiv = 1.0, 1.0, 1.0, 1.0, 1.0
	ttpb := float64(fork.List[1].TargetTimePerBlock)
	startHeight := fork.List[1].ActivationHeight
	if b.params.Net == wire.TestNet3 {
		startHeight = fork.List[1].TestnetStart
	}
	if nH == startHeight {
		log.DEBUG("on hard fork", nH, startHeight)
		return
	}
	//log.DEBUG("startHeight", startHeight)
	f, _ := b.BlockByHeight(startHeight)
	if f != nil {
		fh := f.MsgBlock().Header.BlockHash()
		first := b.Index.LookupNode(&fh)
		allTime := float64(lastNode.timestamp - first.timestamp)
		allBlocks := float64(lastNode.height - first.height)
		// time from lastNode timestamp until start
		if allBlocks == 0 {
			allBlocks = 1
		}
		allTimeAv = allTime / allBlocks
		allTimeDiv = float64(1)
		if allTimeAv > 0 {
			allTimeDiv = allTimeAv / ttpb
		}
		allTimeDiv *= allTimeDiv * allTimeDiv * allTimeDiv * allTimeDiv
	}

	oneHour := 60 * 60 / fork.List[1].TargetTimePerBlock
	oneDay := oneHour * 24
	qHour := 60 * 60 / fork.List[1].TargetTimePerBlock / 4
	dayBlock := lastNode.RelativeAncestor(oneDay)
	dayDiv = allTimeDiv
	if dayBlock != nil {
		// collect timestamps within averaging interval
		dayStamps := []int64{lastNode.timestamp}
		for ln := lastNode; ln != nil && ln.height > startHeight+1 &&
			len(dayStamps) <= int(fork.List[1].AveragingInterval); {
			ln = ln.RelativeAncestor(oneDay)
			if ln == nil {
				break
			}
			dayStamps = append(dayStamps, ln.timestamp)
		}
		if len(dayStamps) > minAvSamples {
			intervals := float64(0)
			// calculate intervals
			dayIntervals := []int64{}
			for i := range dayStamps {
				if i > 0 {
					r := dayStamps[i-1] - dayStamps[i]
					intervals++
					dayIntervals = append(dayIntervals, r)
				}
			}
			if intervals >= minAvSamples {
				// calculate exponential weighted moving average from intervals
				dw := ewma.NewMovingAverage()
				for _, x := range dayIntervals {
					dw.Add(float64(x))
				}
				dayDiv = capP9Adjustment(dw.Value() / ttpb / float64(oneDay))
			}
		}
	}
	hourBlock := lastNode.RelativeAncestor(oneHour)
	hourDiv = allTimeDiv
	if hourBlock != nil {
		// collect timestamps within averaging interval
		hourStamps := []int64{lastNode.timestamp}
		for ln := lastNode; ln.height > startHeight+1 &&
			len(hourStamps) <= int(fork.List[1].AveragingInterval); {
			ln = ln.RelativeAncestor(oneHour)
			if ln == nil {
				break
			}
			hourStamps = append(hourStamps, ln.timestamp)
		}
		if len(hourStamps) > minAvSamples {
			intervals := float64(0)
			// calculate intervals
			hourIntervals := []int64{}
			for i := range hourStamps {
				if i > 0 {
					r := hourStamps[i-1] - hourStamps[i]
					intervals++
					hourIntervals = append(hourIntervals, r)
				}
			}
			if intervals >= minAvSamples {
				// calculate exponential weighted moving average from intervals
				hw := ewma.NewMovingAverage()
				for _, x := range hourIntervals {
					hw.Add(float64(x))
				}
				hourDiv = capP9Adjustment(hw.Value() / ttpb / float64(oneHour))
			}
		}
	}
	qhourBlock := lastNode.RelativeAncestor(qHour)
	qhourDiv = allTimeDiv
	if qhourBlock != nil {
		// collect timestamps within averaging interval
		qhourStamps := []int64{lastNode.timestamp}
		for ln := lastNode; ln != nil && ln.height > startHeight+1 &&
			len(qhourStamps) <= int(fork.List[1].AveragingInterval); {
			ln = ln.RelativeAncestor(qHour)
			if ln == nil {
				break
			}
			qhourStamps = append(qhourStamps, ln.timestamp)
		}
		if len(qhourStamps) > minAvSamples {
			intervals := float64(0)
			// calculate intervals
			qhourIntervals := []int64{}
			for i := range qhourStamps {
				if i > 0 {
					r := qhourStamps[i-1] - qhourStamps[i]
					intervals++
					qhourIntervals = append(qhourIntervals, r)
				}
			}
			if intervals >= minAvSamples {
				// calculate exponential weighted moving average from intervals
				qhw := ewma.NewMovingAverage()
				for _, x := range qhourIntervals {
					qhw.Add(float64(x))
				}
				qhourDiv = capP9Adjustment(qhw.Value() / ttpb / float64(qHour))
			}
		}
	}
	return
}

func (b *BlockChain) GetP9AlgoDiv(allTimeDiv float64, last *BlockNode,
	startHeight int32, algoVer int32, ttpb float64) (algDiv float64) {
	const minAvSamples = 9
	// collect timestamps of same algo of equal number as avinterval
	algDiv = allTimeDiv
	algStamps := []int64{last.timestamp}
	for ln := last; ln != nil && ln.height > startHeight+1 &&
		len(algStamps) <= int(fork.List[1].AveragingInterval); {
		ln = ln.RelativeAncestor(1)
		if ln.version == algoVer {
			algStamps = append(algStamps, ln.timestamp)
		}
	}
	if len(algStamps) > minAvSamples {
		intervals := float64(0)
		// calculate intervals
		algIntervals := []int64{}
		for i := range algStamps {
			if i > 0 {
				r := algStamps[i-1] - algStamps[i]
				intervals++
				algIntervals = append(algIntervals, r)
			}
		}
		if intervals >= minAvSamples {
			// calculate exponential weighted moving average from intervals
			awi := ewma.NewMovingAverage()
			for _, x := range algIntervals {
				awi.Add(float64(x))
			}
			algDiv = capP9Adjustment(awi.Value() / ttpb / float64(len(fork.
				P9Algos)))
		}
	}
	return
}

func (b *BlockChain) GetP9Since(lastNode *BlockNode, algoVer int32) (since,
	ttpb, timeSinceAlgo float64, startHeight int32, last *BlockNode) {
	last = lastNode
	// find the most recent block of the same algo
	ln := last
	for ln.version != algoVer {
		ln = ln.RelativeAncestor(1)
		// if it found nothing, return baseline
		if ln == nil {
			last = nil
			return
		}
		last = ln
	}
	since = float64(lastNode.timestamp - last.timestamp)
	ttpb = float64(fork.List[1].TargetTimePerBlock)
	tspb := ttpb * float64(len(fork.List[1].Algos))
	// ratio of seconds since to target seconds per block times the
	// all time divergence ensures the change scales with the divergence
	// from the target, and favours algos that are later
	timeSinceAlgo = (since / tspb) / 5
	startHeight = fork.List[1].ActivationHeight
	if b.params.Net == wire.TestNet3 {
		startHeight = fork.List[1].TestnetStart
	}
	return
}

func (b *BlockChain) IsP9HardFork(nH int32) bool {
	// At activation difficulty resets
	if b.params.Net == wire.MainNet {
		if fork.List[1].ActivationHeight == nH {
			return true
		}
	}
	if b.params.Net == wire.TestNet3 {
		if fork.List[1].TestnetStart == nH {
			return true
		}
	}
	return false
}

func capP9Adjustment(adjustment float64) float64 {
	const max float64 = 65536
	const maxA, minA = max, 1 / max
	if adjustment > maxA {
		adjustment = maxA
	}
	if adjustment < minA {
		adjustment = minA
	}
	return adjustment
}
