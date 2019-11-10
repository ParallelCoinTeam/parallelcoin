package blockchain

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/p9c/pod/pkg/chain/fork"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/log"
)

// CalcNextRequiredDifficultyPlan9 calculates the required difficulty for the
// block after the passed previous block node based on the difficulty retarget
// rules. This function differs from the exported  CalcNextRequiredDifficulty
// in that the exported version uses the current best chain as the previous
// block node while this function accepts any block node.
func (b *BlockChain) CalcNextRequiredDifficultyPlan9(
	workerNumber uint32, lastNode *blockNode,
	newBlockTime time.Time, algoname string, l bool) (newTargetBits uint32,
	adjustment float64, err error) {
	log.TRACE("algoname ", algoname)
	const max float64 = 65536
	const maxA, minA = max, 1 / max
	const minAvSamples = 9
	nH := lastNode.height + 1
	if lastNode == nil {
		return fork.SecondPowLimitBits, 1, nil
	}
	// At activation difficulty resets
	if b.params.Net == wire.MainNet {
		if fork.List[1].ActivationHeight == nH {
			if l {
				log.DEBUG("on plan 9 hardfork")
			}
			return fork.SecondPowLimitBits, 1, nil
		}
	}
	if b.params.Net == wire.TestNet3 {
		if fork.List[1].TestnetStart == nH {
			if l {
				log.DEBUG("wrkr:", workerNumber, "on plan 9 hardfork", algoname)
			}
			return fork.SecondPowLimitBits, 1, nil
		}
	}
	algoVer := fork.GetAlgoVer(algoname, nH)
	newTargetBits = fork.SecondPowLimitBits
	log.TRACEF("newTarget %08x %s %d", newTargetBits, algoname, algoVer)
	last := lastNode
	// find the most recent block of the same algo
	ln := last
	for ln.version != algoVer {
		ln = ln.RelativeAncestor(1)
		// if it found nothing, return baseline
		if ln == nil {
			return fork.SecondPowLimitBits, 1, nil
		}
		last = ln
	}
	since := float64(lastNode.timestamp - last.timestamp)
	ttpb := float64(fork.List[1].TargetTimePerBlock)
	tspb := ttpb * float64(len(fork.List[1].Algos))
	timeSinceAlgo := (since / tspb) / 5
	startHeight := fork.List[1].ActivationHeight
	if b.params.Net == wire.TestNet3 {
		startHeight = fork.List[1].TestnetStart
	}
	allTimeAv, allTimeDiv, qhourDiv, hourDiv,
	dayDiv := b.GetCommonP9Averages(lastNode, newBlockTime, nH, algoname)

	// ratio of seconds since to target seconds per block times the
	// all time divergence ensures the change scales with the divergence
	// from the target, and favours algos that are later

	algDiv := b.GetP9AlgoDiv(allTimeDiv, last, startHeight, algoVer, ttpb)

	adjustment = (allTimeDiv + algDiv + dayDiv + hourDiv + qhourDiv + timeSinceAlgo) / 6
	if adjustment > maxA {
		adjustment = maxA
	}
	if adjustment < minA {
		adjustment = minA
	}
	log.TRACEF("adjustment %3.4f %08x", adjustment, last.bits)
	bigAdjustment := big.NewFloat(adjustment)
	bigOldTarget := big.NewFloat(1.0).SetInt(fork.CompactToBig(last.bits))
	bigNewTargetFloat := big.NewFloat(1.0).Mul(bigAdjustment, bigOldTarget)
	newTarget, _ := bigNewTargetFloat.Int(nil)
	if newTarget == nil {
		log.INFO("newTarget is nil ")
		return newTargetBits, 1, nil
	}
	if newTarget.Cmp(&fork.FirstPowLimit) < 0 {
		newTargetBits = BigToCompact(newTarget)
		log.TRACEF("newTarget %064x %08x", newTarget, newTargetBits)
	}
	if l {
		an := fork.List[1].AlgoVers[algoVer]
		pad := 14 - len(an)
		if pad > 0 {
			an += strings.Repeat(" ", pad)
		}
		log.DEBUGC(func() string {
			return fmt.Sprintf("wrkr: %s hght: %s %08x %s %s %s %s %s %s %s"+
				" %s %s %08x",
				RightJustify(fmt.Sprint(workerNumber), 3),
				RightJustify(fmt.Sprint(lastNode.height+1), 8),
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
	return newTargetBits, adjustment, nil
}
