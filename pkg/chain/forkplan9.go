package blockchain

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/p9c/pod/pkg/chain/fork"
	"github.com/p9c/pod/pkg/log"
)

// CalcNextRequiredDifficultyPlan9 calculates the required difficulty for the
// block after the passed previous block node based on the difficulty retarget
// rules. This function differs from the exported  CalcNextRequiredDifficulty
// in that the exported version uses the current best chain as the previous
// block node while this function accepts any block node.
func (b *BlockChain) CalcNextRequiredDifficultyPlan9(workerNumber uint32,
	lastNode *BlockNode, newBlockTime time.Time, algoname string,
	l bool) (newTargetBits uint32, adjustment float64, err error) {
	//log.TRACE("algoname ", algoname)
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
			return fmt.Sprintf("wrkr: %s hght: %d %08x %s %s %s %s %s %s %s"+
				" %s %s %08x",
				RightJustify(fmt.Sprint(workerNumber), 3),
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
