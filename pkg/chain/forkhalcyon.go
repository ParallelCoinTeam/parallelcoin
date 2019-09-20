package blockchain

import (
	"fmt"
	"math/big"
	"time"

	"git.parallelcoin.io/dev/pod/pkg/chain/fork"
	"git.parallelcoin.io/dev/pod/pkg/util/cl"
)

// calcNextRequiredDifficultyHalcyon calculates the required difficulty for the
// block after the passed previous block node based on the difficulty retarget
// rules. This function differs from the exported  CalcNextRequiredDifficulty
// in that the exported version uses the current best chain as the previous
// block node while this function accepts any block node.
func (b *BlockChain) CalcNextRequiredDifficultyHalcyon(lastNode *blockNode,
	newBlockTime time.Time, algoname string, l bool) (newTargetBits uint32,
	err error) {

	log <- cl.Trace{"algoname", algoname}
	nH := lastNode.height + 1
	// log <- cl.Info{nH}

	algo := fork.GetAlgoVer(algoname, nH)
	algoName := fork.GetAlgoName(algo, nH)
	newTargetBits = fork.GetMinBits(algoName, nH)
	if lastNode == nil {
		log <- cl.Trace{"lastnode was nil", newTargetBits}
		return newTargetBits, nil
	}
	prevNode := lastNode.GetLastWithAlgo(algo)
	if prevNode == nil {
		log <- cl.Trace{"prevnode was nil", newTargetBits}
		return newTargetBits, nil
	}
	firstNode := prevNode
	for i := int64(0); firstNode != nil &&
		i < fork.GetAveragingInterval(nH)-1; i++ {
		firstNode = firstNode.RelativeAncestor(1)
		firstNode = firstNode.GetLastWithAlgo(algo)
	}
	if firstNode == nil {
		log <- cl.Trace{"firstnode was nil", newTargetBits}
		return newTargetBits, nil
	}
	actualTimespan := prevNode.timestamp - firstNode.timestamp
	adjustedTimespan := actualTimespan
	if actualTimespan < b.params.MinActualTimespan {
		adjustedTimespan = b.params.MinActualTimespan
	} else if actualTimespan > b.params.MaxActualTimespan {
		adjustedTimespan = b.params.MaxActualTimespan
	}
	log <- cl.Trace{"from bits", newTargetBits}
	newTarget := fork.CompactToBig(prevNode.bits)
	log <- cl.Trace{"to big", newTarget}
	bigAdjustedTimespan := big.NewInt(adjustedTimespan)
	newTarget = newTarget.Mul(bigAdjustedTimespan, newTarget)
	log <- cl.Trace{"multiplied", newTarget, bigAdjustedTimespan}
	newTarget = newTarget.Div(newTarget, big.NewInt(b.params.AveragingTargetTimespan))
	log <- cl.Trace{"divided", newTarget}
	if newTarget.Cmp(fork.CompactToBig(newTargetBits)) > 0 {
		log <- cl.Trace{"fell under", newTarget}
		newTarget.Set(fork.CompactToBig(newTargetBits))
	}
	log <- cl.Trace{"newTarget", newTarget}
	newTargetBits = BigToCompact(newTarget)
	log <- cl.Trace{"divided", newTargetBits}
	log <- cl.Debugc(func() string {
		return fmt.Sprintf(
			"difficulty retarget at block height %d, old %08x new %08x %s",
			lastNode.height+1, prevNode.bits, newTargetBits, cl.Ine())
	})
	log <- cl.Tracec(func() string {
		return fmt.Sprintf(
			"actual timespan %v, adjusted timespan %v, target timespan %v",
			// +					"\nOld %064x\nNew %064x",
			actualTimespan,
			adjustedTimespan,
			b.params.AveragingTargetTimespan,
			// oldTarget,
			// fork.CompactToBig(newTargetBits),
		)
	})
	log <- cl.Tracef{"newtarget bits %8x %s", newTargetBits, cl.Ine()}
	return BigToCompact(newTarget), nil

}
