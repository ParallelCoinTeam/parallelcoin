package blockchain

import (
	"github.com/VividCortex/ewma"
	"github.com/p9c/pod/pkg/chain/fork"
	"github.com/p9c/pod/pkg/chain/wire"
	"time"
)

func (b *BlockChain) GetCommonP9Averages(lastNode *blockNode,
	newBlockTime time.Time, nH int32, algoname string) (allTimeAv,
	allTimeDiv, qhourDiv, hourDiv, dayDiv float64) {
	const max float64 = 65536
	const maxA, minA = max, 1 / max
	const minAvSamples = 9
	ttpb := float64(fork.List[1].TargetTimePerBlock)
	startHeight := fork.List[1].ActivationHeight
	if b.params.Net == wire.TestNet3 {
		startHeight = fork.List[1].TestnetStart
	}
	f, _ := b.BlockByHeight(startHeight)
	fh := f.MsgBlock().Header.BlockHash()
	first := b.Index.LookupNode(&fh)
	// time from lastNode timestamp until start
	allTime := float64(lastNode.timestamp - first.timestamp)
	allBlocks := float64(lastNode.height - first.height)
	if allBlocks == 0 {
		allBlocks = 1
	}
	allTimeAv = allTime / allBlocks
	allTimeDiv = float64(1)
	if allTimeAv > 0 {
		allTimeDiv = allTimeAv / ttpb
	}
	allTimeDiv *= allTimeDiv * allTimeDiv * allTimeDiv * allTimeDiv * allTimeDiv * allTimeDiv

	oneHour := 60 * 60 / fork.List[1].TargetTimePerBlock
	oneDay := oneHour * 24
	qHour := 60 * 60 / fork.List[1].TargetTimePerBlock / 4
	dayBlock := lastNode.RelativeAncestor(oneDay)
	dayDiv = allTimeDiv
	if dayBlock != nil {
		// collect timestamps within averaging interval
		dayStamps := []int64{lastNode.timestamp}
		for ln := lastNode; ln != nil && ln.height > startHeight &&
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
			if intervals > minAvSamples {
				// calculate exponential weighted moving average from intervals
				dw := ewma.NewMovingAverage()
				for _, x := range dayIntervals {
					dw.Add(float64(x))
				}
				dayDiv = dw.Value() / ttpb / float64(oneDay)
				if dayDiv < minA {
					dayDiv = minA
				}
				if dayDiv > maxA {
					dayDiv = maxA
				}
			}
		}
	}
	hourBlock := lastNode.RelativeAncestor(oneHour)
	hourDiv = allTimeDiv
	if hourBlock != nil {
		// collect timestamps within averaging interval
		hourStamps := []int64{lastNode.timestamp}
		for ln := lastNode; ln.height > startHeight &&
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
			if intervals > minAvSamples {
				// calculate exponential weighted moving average from intervals
				hw := ewma.NewMovingAverage()
				for _, x := range hourIntervals {
					hw.Add(float64(x))
				}
				hourDiv = hw.Value() / ttpb / float64(oneHour)
				if hourDiv < minA {
					hourDiv = minA
				}
				if hourDiv > maxA {
					hourDiv = maxA
				}
			}
		}
	}
	qhourBlock := lastNode.RelativeAncestor(qHour)
	qhourDiv = allTimeDiv
	if qhourBlock != nil {
		// collect timestamps within averaging interval
		qhourStamps := []int64{lastNode.timestamp}
		for ln := lastNode; ln != nil && ln.height > startHeight &&
			len(qhourStamps) <= int(fork.List[1].AveragingInterval); {
			ln = ln.RelativeAncestor(qHour)
			if ln == nil {
				break
			}
			qhourStamps = append(qhourStamps, ln.timestamp)
		}
		if len(qhourStamps) > 2 {
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
			if intervals > 1 {
				// calculate exponential weighted moving average from intervals
				qhw := ewma.NewMovingAverage()
				for _, x := range qhourIntervals {
					qhw.Add(float64(x))
				}
				qhourDiv = qhw.Value() / ttpb / float64(qHour)
				if qhourDiv < minA {
					qhourDiv = minA
				}
				if qhourDiv > maxA {
					qhourDiv = maxA
				}
			}
		}
	}
	return
}

func (b *BlockChain) GetP9AlgoDiv(allTimeDiv float64, last *blockNode,
	startHeight int32, algoVer int32, ttpb float64) (algDiv float64) {
	const max float64 = 65536
	const maxA, minA = max, 1 / max
	const minAvSamples = 9
	// collect timestamps of same algo of equal number as avinterval
	algDiv = allTimeDiv
	algStamps := []int64{last.timestamp}
	for ln := last; ln != nil && ln.height > startHeight &&
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
		if intervals > minAvSamples {
			// calculate exponential weighted moving average from intervals
			awi := ewma.NewMovingAverage()
			for _, x := range algIntervals {
				awi.Add(float64(x))
			}
			algDiv = awi.Value() / ttpb / float64(len(fork.P9Algos))
			if algDiv < minA {
				algDiv = minA
			}
			if algDiv > maxA {
				algDiv = maxA
			}
		}
	}
	return
}
