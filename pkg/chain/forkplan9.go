package blockchain

import (
	"fmt"
	"math/big"
	"strings"
	
	"github.com/VividCortex/ewma"
	
	"github.com/p9c/pod/pkg/chain/fork"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/log"
)

// CalcNextRequiredDifficultyPlan9 implements the Parallel Prime Difficulty
// Adjustment.
// From the first 9 primes 2, 3, 5, 7, 11, 13, 17, 19, 23
// using these values to pick the primes in sequence of these numbers:
// 3, 5, 11, 17, 31, 41, 59, 67, 83 being the interval in seconds
// This sequence has an effective cumulative product of 37066.310186611 years
// Thus it is effectively random as a whole though each parallel interval/version
// is cyclic.
func (b *BlockChain) CalcNextRequiredDifficultyPlan9(workerNumber uint32,
	lastNode *BlockNode, algoname string, l bool) (newTargetBits uint32,
	adjustment float64, err error) {
	newTargetBits = fork.SecondPowLimitBits
	if lastNode == nil {
		log.TRACE("lastNode is nil")
	}
	algoInterval := fork.P9Algos[algoname].VersionInterval
	algoVer := fork.P9Algos[algoname].Version
	_ = algoInterval
	adjustment = 1
	startHeight := fork.List[1].ActivationHeight
	if b.params.Net == wire.TestNet3 {
		startHeight = fork.List[1].TestnetStart
	}
	var last *BlockNode
	var found bool
	algStamps := []uint64{uint64(lastNode.timestamp)}
	for ln := lastNode; ln != nil && ln.height > startHeight &&
		len(algStamps) <= int(fork.List[1].AveragingInterval); ln = ln.
		RelativeAncestor(1) {
		if ln.version == algoVer && ln.height > startHeight {
			algStamps = append(algStamps, uint64(ln.timestamp))
			if !found {
				found = true
				last = ln
			}
		}
	}
	const minAvSamples = 2
	ttpb := float64(fork.List[1].Algos[algoname].VersionInterval)
	if len(algStamps) > minAvSamples {
		intervals := float64(0)
		// calculate intervals
		algIntervals := []uint64{}
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
			adjustment = capP9Adjustment(awi.Value() / ttpb / float64(len(fork.
				P9Algos)))
		}
	}
	if last == nil {
		log.TRACE("last was nil")
		last = new(BlockNode)
		last.bits = fork.SecondPowLimitBits
		last.version = algoVer
	}
	unsquared := adjustment
	// adjustment *= adjustment
	bigAdjustment := big.NewFloat(adjustment)
	bigOldTarget := big.NewFloat(1.0).SetInt(fork.CompactToBig(last.bits))
	bigNewTargetFloat := big.NewFloat(1.0).Mul(bigAdjustment, bigOldTarget)
	newTarget, _ := bigNewTargetFloat.Int(nil)
	if newTarget == nil {
		log.INFO("newTarget is nil ")
		return
	}
	if newTarget.Cmp(&fork.FirstPowLimit) < 0 {
		newTargetBits = BigToCompact(newTarget)
		// log.TRACEF("newTarget %064x %08x", newTarget, newTargetBits)
	}
	// if l {
	an := fork.List[1].AlgoVers[algoVer]
	pad := 9 - len(an)
	if pad > 0 {
		an += strings.Repeat(" ", pad)
	}
	if l && workerNumber == 0 {
		log.DEBUGC(func() string {
			return fmt.Sprintf("%08x %s %s %08x av %s, %4.0f interval",
				// RightJustify(fmt.Sprint(workerNumber), 3),
				// RightJustify(fmt.Sprint(last.height+1), 9),
				last.bits,
				an,
				RightJustify(fmt.Sprintf("%4.4fx", 1/adjustment), 11),
				newTargetBits,
				RightJustify(fmt.Sprintf("%4.4f", unsquared*ttpb), 11),
				ttpb, // fork.List[1].Algos[algoname].VersionInterval,
			)
		})
	}
	// }
	return
}

// CalcNextRequiredDifficultyPlan9 calculates the required difficulty for the
// block after the passed previous block node based on the difficulty retarget
// rules. This function differs from the exported  CalcNextRequiredDifficulty
// in that the exported version uses the current best chain as the previous
// block node while this function accepts any block node.
func (b *BlockChain) CalcNextRequiredDifficultyPlan9old(lastNode *BlockNode, algoname string, l bool) (newTargetBits uint32, adjustment float64, err error) {
	nH := lastNode.height + 1
	newTargetBits = fork.SecondPowLimitBits
	adjustment = 1.0
	if lastNode == nil || b.IsP9HardFork(nH) {
		return
	}
	allTimeAv, allTimeDiv, qhourDiv, hourDiv,
	dayDiv := b.GetCommonP9Averages(lastNode, nH)
	algoVer := fork.GetAlgoVer(algoname, nH)
	since, ttpb, timeSinceAlgo, startHeight, last := b.GetP9Since(lastNode,
		algoVer)
	if last == nil {
		return
	}
	algDiv := b.GetP9AlgoDiv(allTimeDiv, last, startHeight, algoVer, ttpb)
	adjustment = (allTimeDiv + algDiv + dayDiv + hourDiv + qhourDiv +
		timeSinceAlgo) / 6
	bigAdjustment := big.NewFloat(adjustment)
	bigOldTarget := big.NewFloat(1.0).SetInt(fork.CompactToBig(last.bits))
	bigNewTargetFloat := big.NewFloat(1.0).Mul(bigAdjustment, bigOldTarget)
	newTarget, _ := bigNewTargetFloat.Int(nil)
	if newTarget == nil {
		log.INFO("newTarget is nil ")
		return
	}
	if newTarget.Cmp(&fork.FirstPowLimit) < 0 {
		newTargetBits = BigToCompact(newTarget)
		log.TRACEF("newTarget %064x %08x", newTarget, newTargetBits)
	}
	if l {
		an := fork.List[1].AlgoVers[algoVer]
		pad := 9 - len(an)
		if pad > 0 {
			an += strings.Repeat(" ", pad)
		}
		log.DEBUGC(func() string {
			return fmt.Sprintf("hght: %d %08x %s %s %s %s %s %s %s"+
				" %s %s %08x",
				lastNode.height+1,
				last.bits,
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
				newTargetBits,
			)
		})
	}
	return
}
