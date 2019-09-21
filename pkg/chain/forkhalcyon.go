package blockchain

import (
	"fmt"
	"math/big"
	"time"

	"github.com/parallelcointeam/parallelcoin/pkg/chain/fork"
)

// calcNextRequiredDifficultyHalcyon calculates the required difficulty for the
// block after the passed previous block node based on the difficulty retarget
// rules. This function differs from the exported  CalcNextRequiredDifficulty
// in that the exported version uses the current best chain as the previous
// block node while this function accepts any block node.
func (b *BlockChain) CalcNextRequiredDifficultyHalcyon(lastNode *blockNode,
	newBlockTime time.Time, algoname string, l bool) (newTargetBits uint32,
	err error) {

	TRACE("algoname", algoname)
	nH := lastNode.height + 1
	// INFO{nH}

	algo := fork.GetAlgoVer(algoname, nH)
	algoName := fork.GetAlgoName(algo, nH)
	newTargetBits = fork.GetMinBits(algoName, nH)
	if lastNode == nil {
		TRACE("lastnode was nil", newTargetBits)
		return newTargetBits, nil
	}
	prevNode := lastNode.GetLastWithAlgo(algo)
	if prevNode == nil {
		TRACE("prevnode was nil", newTargetBits)
		return newTargetBits, nil
	}
	firstNode := prevNode
	for i := int64(0); firstNode != nil &&
		i < fork.GetAveragingInterval(nH)-1; i++ {
		firstNode = firstNode.RelativeAncestor(1)
		firstNode = firstNode.GetLastWithAlgo(algo)
	}
	if firstNode == nil {
		TRACE("firstnode was nil", newTargetBits)
		return newTargetBits, nil
	}
	actualTimespan := prevNode.timestamp - firstNode.timestamp
	adjustedTimespan := actualTimespan
	if actualTimespan < b.params.MinActualTimespan {
		adjustedTimespan = b.params.MinActualTimespan
	} else if actualTimespan > b.params.MaxActualTimespan {
		adjustedTimespan = b.params.MaxActualTimespan
	}
	TRACE("from bits", newTargetBits)
	newTarget := fork.CompactToBig(prevNode.bits)
	TRACE("to big", newTarget)
	bigAdjustedTimespan := big.NewInt(adjustedTimespan)
	newTarget = newTarget.Mul(bigAdjustedTimespan, newTarget)
	TRACE("multiplied", newTarget, bigAdjustedTimespan)
	newTarget = newTarget.Div(newTarget, big.NewInt(b.params.AveragingTargetTimespan))
	TRACE("divided", newTarget)
	if newTarget.Cmp(fork.CompactToBig(newTargetBits)) > 0 {
		TRACE("fell under", newTarget)
	}
	newTarget.Set(fork.CompactToBig(newTargetBits))
	TRACE("newTarget", newTarget)
	newTargetBits = BigToCompact(newTarget)
	TRACE("divided", newTargetBits)
	DEBUGC(func() string {
		return fmt.Sprintf("difficulty retarget at block height %d, "+
			"old %08x new %08x", lastNode.height+1, prevNode.bits,
			newTargetBits)
	})
	TRACEC(func() string {
		return fmt.Sprintf(
			"actual timespan %v, adjusted timespan %v, target timespan %v",
			// "\nOld %064x\nNew %064x",
			actualTimespan,
			adjustedTimespan,
			b.params.AveragingTargetTimespan,
			// oldTarget,
			// fork.CompactToBig(newTargetBits),
		)
	})
	TRACEF("newtarget bits %8x %s", newTargetBits)
	return BigToCompact(newTarget), nil

}
