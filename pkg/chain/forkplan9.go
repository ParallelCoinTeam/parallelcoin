package blockchain

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/VividCortex/ewma"

	"github.com/p9c/pod/pkg/chain/fork"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/util/cl"
)

// CalcNextRequiredDifficultyPlan9 calculates the required difficulty for the
// block after the passed previous block node based on the difficulty retarget
// rules. This function differs from the exported  CalcNextRequiredDifficulty
// in that the exported version uses the current best chain as the previous
// block node while this function accepts any block node.
func (b *BlockChain) CalcNextRequiredDifficultyPlan9(lastNode *blockNode,
	newBlockTime time.Time, algoname string, l bool) (newTargetBits uint32,
	adjustment float64, err error) {
	log <- cl.Trace{"algoname ", algoname, cl.Ine()}
	const max float64 = 10000
	const maxA, minA = max, 1 / max
	const minAvSamples = 9
	// square := func(f float64) float64 {
	// 	return f * f
	// }
	nH := lastNode.height + 1
	if lastNode == nil {
		return fork.FirstPowLimitBits, 1, nil
	}
	// At activation difficulty resets
	if b.params.Net == wire.MainNet {
		if fork.List[1].ActivationHeight == nH {
			if l {
				log <- cl.Debug{"on plan 9 hardfork", cl.Ine()}
			}
			return fork.FirstPowLimitBits, 1, nil
		}
	}
	if b.params.Net == wire.TestNet3 {
		if fork.List[1].TestnetStart == nH {
			if l {
				log <- cl.Debug{"on plan 9 hardfork", cl.Ine()}
			}
			return fork.FirstPowLimitBits, 1, nil
		}
	}
	algoVer := fork.GetAlgoVer(algoname, nH)
	newTargetBits = fork.SecondPowLimitBits
	log <- cl.Tracef{"newTarget %08x %s %d %s", newTargetBits, algoname, algoVer, cl.Ine()}
	last := lastNode
	// find the most recent block of the same algo
	//
	if last.version != algoVer {
		ln := last.RelativeAncestor(1)
		if ln == nil {
			return fork.FirstPowLimitBits, 1, nil
		}
		if ln.version == algoVer {
			last = ln
		} else {
			for ln != nil && ln.version != algoVer {
				ln = ln.RelativeAncestor(1)
				// if it found nothing, return baseline
				//
				if ln == nil {
					if l {
						log <- cl.Debug{"before first ", algoname, cl.Ine()}
					}
					return fork.FirstPowLimitBits, 1, nil
				}
				// ignore the first block as its time is not a normal timestamp
				//
				if ln.height < 1 {
					return fork.SecondPowLimitBits, 1, nil
				}
				last = ln
			}
		}
	}
	ttpb := float64(fork.List[1].TargetTimePerBlock)
	startHeight := fork.List[1].ActivationHeight
	if b.params.Net == wire.TestNet3 {
		startHeight = fork.List[1].TestnetStart
	}
	f, _ := b.BlockByHeight(startHeight)
	fh := f.MsgBlock().Header.BlockHash()
	first := b.Index.LookupNode(&fh)
	// time from lastNode timestamp until start
	//
	allTime := float64(lastNode.timestamp - first.timestamp)
	allBlocks := float64(lastNode.height - first.height)
	if allBlocks == 0 {
		allBlocks = 1
	}
	allTimeAv := allTime / allBlocks
	allTimeDiv := float64(1)
	if allTimeAv > 0 {
		allTimeDiv = allTimeAv / ttpb
	}
	allTimeDiv *= allTimeDiv * allTimeDiv * allTimeDiv * allTimeDiv * allTimeDiv * allTimeDiv
	// collect timestamps of same algo of equal number as avinterval
	algDiv := float64(1)
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
			if l {
				log <- cl.Trace{"algs", algIntervals, cl.Ine()}
			}
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
	} else {
		// if there is no intervals this algo needs some love
		// return fork.FirstPowLimitBits, 1, nil
	}
	tspb := ttpb * float64(len(fork.List[1].Algos))
	since := float64(lastNode.timestamp - last.timestamp)
	// ratio of seconds since to target seconds per block times the
	// all time divergence ensures the change scales with the divergence
	// from the target, and favours algos that are later
	timeSinceAlgo := (since / tspb) / 5 // * (since / tspb) // * allTimeDiv
	oneHour := 60 * 60 / fork.List[1].TargetTimePerBlock
	oneDay := oneHour * 24
	qHour := 60 * 60 / fork.List[1].TargetTimePerBlock / 4
	dayBlock := lastNode.RelativeAncestor(oneDay)
	dayDiv := allTimeDiv
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
				if l {
					log <- cl.Trace{"da", dayIntervals, cl.Ine()}
				}
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
	hourDiv := allTimeDiv
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
				if l {
					log <- cl.Trace{"hr", hourIntervals, cl.Ine()}
				}
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
	qhourDiv := allTimeDiv
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
				if l {
					log <- cl.Trace{"qh", qhourIntervals, cl.Ine()}
				}
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
	adjustment = (allTimeDiv + algDiv + dayDiv + hourDiv + qhourDiv + timeSinceAlgo) / 6
	log <- cl.Trace{"adjustment %3.4f %08x %s", adjustment, last.bits, cl.Ine()}
	bigAdjustment := big.NewFloat(adjustment)
	bigOldTarget := big.NewFloat(1.0).SetInt(fork.CompactToBig(last.bits))
	bigNewTargetFloat := big.NewFloat(1.0).Mul(bigAdjustment, bigOldTarget)
	newTarget, _ := bigNewTargetFloat.Int(nil)
	if newTarget == nil {
		log <- cl.Info{"newTarget is nil ", cl.Ine()}
		return newTargetBits, 1, nil
	}
	if newTarget.Cmp(&fork.FirstPowLimit) < 0 {
		newTargetBits = BigToCompact(newTarget)
		log <- cl.Tracef{"newTarget %064x %08x %s", newTarget, newTargetBits, cl.Ine()}
	}
	if l {
		an := fork.List[1].AlgoVers[algoVer]
		pad := 14 - len(an)
		if pad > 0 {
			an += strings.Repeat(" ", pad)
		}
		log <- cl.Debugf{
			"%d %s %s %s %s %s %s %s %s %s %s",
			lastNode.height + 1,
			an,
			RightJustify(fmt.Sprintf("%3.2f", allTimeAv), 5),
			RightJustify(fmt.Sprintf("%3.2fa", allTimeDiv*ttpb), 7),
			RightJustify(fmt.Sprintf("%3.2fd", dayDiv*ttpb), 7),
			RightJustify(fmt.Sprintf("%3.2fh", hourDiv*ttpb), 7),
			RightJustify(fmt.Sprintf("%3.2fq", qhourDiv*ttpb), 7),
			RightJustify(fmt.Sprintf("%3.2fA", algDiv*ttpb), 7),
			RightJustify(fmt.Sprintf("%3.0f %3.3fD",
				since-ttpb*float64(len(fork.List[1].Algos)), timeSinceAlgo*ttpb), 13),
			RightJustify(fmt.Sprintf("%4.4fx", 1/adjustment), 11),
			cl.Ine()}
	}
	return newTargetBits, adjustment, nil
}
