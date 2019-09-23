package blockchain

import (
	"math/big"
	"time"

	"github.com/p9c/pod/pkg/chain/fork"
	"github.com/p9c/pod/pkg/log"
)

// calcNextRequiredDifficultyHalcyon calculates the required difficulty for the
// block after the passed previous block node based on the difficulty retarget
// rules. This function differs from the exported  CalcNextRequiredDifficulty
// in that the exported version uses the current best chain as the previous
// block node while this function accepts any block node.
func (b *BlockChain) CalcNextRequiredDifficultyHalcyon(lastNode *blockNode,
	newBlockTime time.Time, algoname string, l bool) (newTargetBits uint32,
	err error) {

	// log.WARN("next required diff for halcyon", algoname)
	nH := lastNode.height + 1
	// INFO{nH}

	algo := fork.GetAlgoVer(algoname, nH)
	algoName := fork.GetAlgoName(algo, nH)
	newTargetBits = fork.GetMinBits(algoName, nH)
	if lastNode == nil {
		log.WARN("lastnode was nil", newTargetBits)
		return newTargetBits, nil
	}
	prevNode := lastNode.GetLastWithAlgo(algo)
	if prevNode == nil {
		log.WARN("prevnode was nil", newTargetBits)
		return newTargetBits, nil
	}
	newTargetBits = prevNode.bits
	firstNode := prevNode
	for i := int64(0); firstNode != nil &&
		i < fork.GetAveragingInterval(nH)-1; i++ {
		// firstNode = firstNode.RelativeAncestor(1)
		firstNode = firstNode.GetLastWithAlgo(algo)
	}
	if firstNode == nil {
		log.WARN("firstnode was nil", newTargetBits)
		return newTargetBits, nil
	}
	actualTimespan := prevNode.timestamp - firstNode.timestamp
	adjustedTimespan := actualTimespan
	if actualTimespan < b.params.MinActualTimespan {
		adjustedTimespan = b.params.MinActualTimespan
	} else if actualTimespan > b.params.MaxActualTimespan {
		adjustedTimespan = b.params.MaxActualTimespan
	}
	// log.WARNF("from bits %08x", newTargetBits)
	newTarget := fork.CompactToBig(prevNode.bits)
	// log.WARNF("to big %064x", newTarget)
	bigAdjustedTimespan := big.NewInt(adjustedTimespan)
	newTarget = newTarget.Mul(bigAdjustedTimespan, newTarget)
	// log.WARNF("multiplied %064x, %d", newTarget, bigAdjustedTimespan)
	newTarget = newTarget.Div(newTarget, big.NewInt(b.params.AveragingTargetTimespan))
	// log.WARNF("divided %064x", newTarget)
	if newTarget.Cmp(fork.CompactToBig(newTargetBits)) > 0 {
		// log.WARNF("fell under %064x", newTarget)
	}
	// newTarget.Set(fork.CompactToBig(newTargetBits))
	// log.WARNF("newTarget %064x", newTarget)
	newTargetBits = BigToCompact(newTarget)
	// log.WARNF("divided %08x", newTargetBits)
	// log.DEBUGC(func() string {
	// 	return fmt.Sprintf("difficulty retarget at block height %d, "+
	// 		"old %08x new %08x", lastNode.height+1, prevNode.bits,
	// 		newTargetBits) +
	// 		fmt.Sprintf(
	// 			"\nactual timespan %v, adjusted timespan %v, " +
	// 				"target timespan %v",
	// 			// "\nOld %064x\nNew %064x",
	// 			actualTimespan,
	// 			adjustedTimespan,
	// 			b.params.AveragingTargetTimespan,
	// 			// oldTarget,
	// 			// fork.CompactToBig(newTargetBits),
	// 		)
	// })
	log.TRACEF("newtarget bits %8x %s", newTargetBits)
	return BigToCompact(newTarget), nil
}
