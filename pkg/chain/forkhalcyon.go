package blockchain

import (
	"fmt"
	"math/big"

	"github.com/p9c/pod/pkg/chain/fork"
)

// calcNextRequiredDifficultyHalcyon calculates the required difficulty for the
// block after the passed previous block node based on the difficulty retarget
// rules. This function differs from the exported  CalcNextRequiredDifficulty
// in that the exported version uses the current best chain as the previous
// block node while this function accepts any block node.
func (b *BlockChain) CalcNextRequiredDifficultyHalcyon(workerNumber uint32, lastNode *BlockNode, algoname string, l bool) (newTargetBits uint32, err error) {
	if workerNumber != 0 {
		l = false
	}
	nH := lastNode.height + 1
	if lastNode == nil {
		if l {
			L.Debug("lastNode is nil")
		}
		return newTargetBits, nil
	}
	// this sanitises invalid block versions according to legacy consensus quirks
	algo := fork.GetAlgoVer(algoname, nH)
	algoName := fork.GetAlgoName(algo, nH)
	newTargetBits = fork.GetMinBits(algoName, nH)
	prevNode := lastNode.GetLastWithAlgo(algo)
	if prevNode == nil {
		if l {
			L.Debug("prevNode is nil")
		}
		return newTargetBits, nil
	}
	firstNode := prevNode
	for i := int32(0); firstNode != nil &&
		i < fork.GetAveragingInterval(nH)-1; i++ {
		firstNode = firstNode.RelativeAncestor(1)
		firstNode = firstNode.GetLastWithAlgo(algo)
	}
	if firstNode == nil {
		return newTargetBits, nil
	}
	actualTimespan := prevNode.timestamp - firstNode.timestamp
	adjustedTimespan := actualTimespan
	if l {
		L.Tracef("actual %d", actualTimespan)
	}
	if actualTimespan < b.params.MinActualTimespan {
		adjustedTimespan = b.params.MinActualTimespan
	} else if actualTimespan > b.params.MaxActualTimespan {
		adjustedTimespan = b.params.MaxActualTimespan
	}
	if l {
		L.Tracef("adjusted %d", adjustedTimespan)
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
		L.Debugf(
			"difficulty retarget at block height %d, old %08x new %08x",
			lastNode.height+1,
			prevNode.bits,
			newTargetBits,
		)
	}
	if l {
		L.Tracec(func() string {
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
