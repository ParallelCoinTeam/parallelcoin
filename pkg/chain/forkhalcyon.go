package blockchain

import (
	"fmt"
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
func (b *BlockChain) CalcNextRequiredDifficultyHalcyon(
	workerNumber uint32, lastNode *blockNode,
	newBlockTime time.Time, algoname string, l bool) (newTargetBits uint32,
	err error) {
	if workerNumber != 0 {
		l = false
	}
	if l {
		log.TRACE("CalcNextRequiredDifficultyHalcyon", workerNumber, algoname, lastNode.height)
	}
	nH := lastNode.height + 1
	if l {
		log.TRACE("on pre-hardfork")
	}
	if lastNode == nil {
		if l {
			log.DEBUG("lastNode is nil")
		}
		return newTargetBits, nil
	}
	algo := fork.GetAlgoVer(algoname, nH)
	algoName := fork.GetAlgoName(algo, nH)
	newTargetBits = fork.GetMinBits(algoName, nH)
	if l {
		log.TRACEF("last %d %d %8x", lastNode.height, lastNode.version, lastNode.bits)
	}
	prevNode := lastNode.GetLastWithAlgo(algo)
	if prevNode == nil {
		if l {
			log.DEBUG("prevNode is nil")
		}
		return newTargetBits, nil
	}
	if l {
		log.DEBUG("prev with algo", algoname, prevNode.height, prevNode.hash)
	}
	firstNode := prevNode
	for i := int64(0); firstNode != nil &&
		i < fork.GetAveragingInterval(nH)-1; i++ {
		if l {
			log.TRACEF("%d: prev %d %d %8x",
				i, firstNode.height, firstNode.version, firstNode.bits)
		}
		firstNode = firstNode.RelativeAncestor(1)
		firstNode = firstNode.GetLastWithAlgo(algo)
	}
	if firstNode == nil {
		return newTargetBits, nil
	}
	//log.DEBUGF("9: first %d %d %8x",		firstNode.height, firstNode.version, firstNode.bits)
	actualTimespan := prevNode.timestamp - firstNode.timestamp
	adjustedTimespan := actualTimespan
	if l {
		log.TRACEF("actual %d", actualTimespan)
	}
	if actualTimespan < b.params.MinActualTimespan {
		adjustedTimespan = b.params.MinActualTimespan
	} else if actualTimespan > b.params.MaxActualTimespan {
		adjustedTimespan = b.params.MaxActualTimespan
	}
	if l {
		log.TRACEF("adjusted %d", adjustedTimespan)
	}
	oldTarget := CompactToBig(prevNode.bits)
	newTarget := new(big.Int).
		Mul(oldTarget, big.NewInt(adjustedTimespan))
	newTarget = newTarget.
		Div(newTarget, big.NewInt(b.params.AveragingTargetTimespan))
	if newTarget.Cmp(CompactToBig(newTargetBits)) > 0 {
		newTarget.Set(CompactToBig(newTargetBits))
	}
	newTargetBits = BigToCompact(newTarget)
	if l {
		log.DEBUGF(
			"difficulty retarget at block height %d, old %08x new %08x",
			lastNode.height+1,
			prevNode.bits,
			newTargetBits,
		)
	}
	if l {
		log.TRACEC(func() string {
			return fmt.Sprintf(
				"actual timespan %v, adjusted timespan %v, target timespan %v"+
					"\nOld %064x\nNew %064x",
				actualTimespan,
				adjustedTimespan,
				b.params.AveragingTargetTimespan,
				oldTarget,
				CompactToBig(newTargetBits),
			)
		})
	}
	return newTargetBits, nil
}

//
//
// 		// log.WARN("next required diff for halcyon", algoname)
// 	nH := lastNode.height + 1
// 	// INFO{nH}
//
// 	algo := fork.GetAlgoVer(algoname, nH)
// 	algoName := fork.GetAlgoName(algo, nH)
// 	newTargetBits = fork.GetMinBits(algoName, nH)
// 	if lastNode == nil {
// 		log.WARN("lastnode was nil", newTargetBits)
// 		return newTargetBits, nil
// 	}
// 	prevNode := lastNode.GetLastWithAlgo(algo)
// 	if prevNode == nil {
// 		log.WARN("prevnode was nil", newTargetBits)
// 		return newTargetBits, nil
// 	}
// 	newTargetBits = prevNode.bits
// 	firstNode := prevNode
// 	for i := int64(0); firstNode != nil &&
// 		i < fork.GetAveragingInterval(nH)-1; i++ {
// 		// firstNode = firstNode.RelativeAncestor(1)
// 		firstNode = firstNode.GetLastWithAlgo(algo)
// 	}
// 	if firstNode == nil {
// 		log.WARN("firstnode was nil", newTargetBits)
// 		return newTargetBits, nil
// 	}
// 	actualTimespan := prevNode.timestamp - firstNode.timestamp
// 	adjustedTimespan := actualTimespan
// 	if actualTimespan < b.params.MinActualTimespan {
// 		adjustedTimespan = b.params.MinActualTimespan
// 	} else if actualTimespan > b.params.MaxActualTimespan {
// 		adjustedTimespan = b.params.MaxActualTimespan
// 	}
// 	// log.WARNF("from bits %08x", newTargetBits)
// 	newTarget := fork.CompactToBig(prevNode.bits)
// 	// log.WARNF("to big %064x", newTarget)
// 	bigAdjustedTimespan := big.NewInt(adjustedTimespan)
// 	newTarget = newTarget.Mul(bigAdjustedTimespan, newTarget)
// 	// log.WARNF("multiplied %064x, %d", newTarget, bigAdjustedTimespan)
// 	newTarget = newTarget.Div(newTarget, big.NewInt(b.params.AveragingTargetTimespan))
// 	// log.WARNF("divided %064x", newTarget)
// 	if newTarget.Cmp(fork.CompactToBig(newTargetBits)) > 0 {
// 		// log.WARNF("fell under %064x", newTarget)
// 	}
// 	// newTarget.Set(fork.CompactToBig(newTargetBits))
// 	// log.WARNF("newTarget %064x", newTarget)
// 	newTargetBits = BigToCompact(newTarget)
// 	// log.WARNF("divided %08x", newTargetBits)
// 	// log.DEBUGC(func() string {
// 	// 	return fmt.Sprintf("difficulty retarget at block height %d, "+
// 	// 		"old %08x new %08x", lastNode.height+1, prevNode.bits,
// 	// 		newTargetBits) +
// 	// 		fmt.Sprintf(
// 	// 			"\nactual timespan %v, adjusted timespan %v, " +
// 	// 				"target timespan %v",
// 	// 			// "\nOld %064x\nNew %064x",
// 	// 			actualTimespan,
// 	// 			adjustedTimespan,
// 	// 			b.params.AveragingTargetTimespan,
// 	// 			// oldTarget,
// 	// 			// fork.CompactToBig(newTargetBits),
// 	// 		)
// 	// })
// 	log.TRACEF("newtarget bits %8x %s", newTargetBits)
// 	return BigToCompact(newTarget), nil
// }
