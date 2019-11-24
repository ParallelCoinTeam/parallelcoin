package blockchain

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/p9c/pod/pkg/chain/fork"
	"github.com/p9c/pod/pkg/log"
)

func secondPowLimitBits(currFork int) (out *map[int32]uint32) {
	aV := fork.List[currFork].AlgoVers
	o := make(map[int32]uint32, len(aV))
	for i := range aV {
		o[i] = fork.SecondPowLimitBits
	}
	return &o
}

// CalcNextRequiredDifficultyPlan9Controller returns all of the algorithm
// difficulty targets for sending out with the other pieces required to
// construct a block, as these numbers are generated from block timestamps
func (b *BlockChain) CalcNextRequiredDifficultyPlan9Controller(
	lastNode *BlockNode) (newTargetBits *map[int32]uint32, err error) {
	nH := lastNode.height + 1
	currFork := fork.GetCurrent(nH)
	nTB := make(map[int32]uint32)
	newTargetBits = &nTB
	if currFork == 0 {
		for i := range fork.List[0].Algos {
			v := fork.List[0].Algos[i].Version
			nTB[v], err = b.CalcNextRequiredDifficultyHalcyon(0,
				lastNode, time.Now(), i, true)
		}
		return
	}
	lnh := lastNode.Header()
	hD := &lnh
	newTargetBits = secondPowLimitBits(currFork)
	if lastNode == nil || b.IsP9HardFork(nH) {
		return
	}
	log.TRACEC(func() string {
		return fmt.Sprint("calculating difficulty targets to attach to"+
			" block ", hD.BlockHashWithAlgos(lastNode.height), lastNode.height)
	})
	tn := time.Now()
	defer log.TRACEC(func() string {
		return fmt.Sprint(time.Now().Sub(tn), " to calculate all diffs")
	})
	// here we only need to do this once
	allTimeAv, allTimeDiv, qhourDiv, hourDiv, dayDiv := b.
		GetCommonP9Averages(lastNode, nH)
	for aV := range fork.List[currFork].AlgoVers {
		// TODO: merge this with the single algorithm one
		since, ttpb, timeSinceAlgo, startHeight, last := b.GetP9Since(lastNode,
			aV)
		if last == nil {
			log.DEBUG("last is nil")
			return
		}
		algDiv := b.GetP9AlgoDiv(allTimeDiv, last, startHeight, aV, ttpb)
		adjustment := (allTimeDiv + algDiv + dayDiv + hourDiv + qhourDiv +
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
			(*newTargetBits)[aV] = BigToCompact(newTarget)
		}
		an := fork.List[1].AlgoVers[aV]
		pad := 9 - len(an)
		if pad > 0 {
			an += strings.Repeat(" ", pad)
		}
		log.DEBUGC(func() string {
			return fmt.Sprintf("hght: %d %08x %s %s %s %s %s %s %s"+
				" %s %s %08x",
				nH,
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
				(*newTargetBits)[aV],
			)
		})
	}
	return
}
