package blockchain

import (
	"fmt"
	bits2 "github.com/p9c/matrjoska/pkg/bits"
	"github.com/p9c/matrjoska/pkg/fork"
	"math/big"
	"strings"
	
	"github.com/VividCortex/ewma"
	
	"github.com/p9c/matrjoska/pkg/wire"
)

// GetAlgStamps ...
func GetAlgStamps(algoName string, startHeight int32, lastNode *BlockNode) (last *BlockNode,
	found bool, algStamps []int64, version int32,
) {
	
	version = fork.P9Algos[algoName].Version
	for ln := lastNode; ln != nil && ln.height > startHeight &&
		len(algStamps) <= int(fork.List[1].AveragingInterval); ln = ln.
		RelativeAncestor(1) {
		if ln.version == version && ln.height > startHeight {
			algStamps = append(algStamps, ln.timestamp)
			if !found {
				found = true
				last = ln
			}
		}
	}
	// D.Ln(algStamps)
	// reverse order of stamps
	for i := 0; i < len(algStamps)/2; i++ {
		algStamps[i], algStamps[len(algStamps)-i-1] = algStamps[len(
			algStamps,
		)-i-1], algStamps[i]
	}
	// D.Ln(algStamps)
	return
}

func GetAllStamps(startHeight int32, lastNode *BlockNode) (allStamps []int64) {
	
	for ln := lastNode; ln != nil && ln.height > startHeight &&
		len(allStamps) <= int(fork.List[1].AveragingInterval); ln = ln.RelativeAncestor(1) {
		allStamps = append(allStamps, ln.timestamp)
	}
	// D.Ln(allStamps)
	// reverse order of stamps
	for i := 0; i < len(allStamps)/2; i++ {
		allStamps[i], allStamps[len(allStamps)-i-1] =
			allStamps[len(allStamps)-i-1], allStamps[i]
	}
	// D.Ln(allStamps)
	return
}

func GetAll(allStamps []int64) (allAv, allAdj float64) {
	allAdj = 1
	allAv = fork.P9Average
	// calculate intervals
	allIntervals := make([]float64, len(allStamps)-1)
	for i := range allStamps {
		if i > 0 {
			r := allStamps[i] - allStamps[i-1]
			allIntervals[i-1] = float64(r)
		}
	}
	// D.Ln(allStamps)
	// calculate exponential weighted moving average from intervals
	aewma := ewma.NewMovingAverage()
	for _, x := range allIntervals {
		aewma.Add(x)
	}
	allAv = aewma.Value()
	// W.Ln(allAv)
	if allAv != 0 {
		allAdj = allAv / fork.P9Average
	}
	return
}

func GetAlg(algStamps []int64, targetTimePerBlock float64) (algAv, algAdj float64) {
	// calculate intervals
	algIntervals := make([]int64, len(algStamps)-1)
	for i := range algStamps {
		if i > 0 {
			r := algStamps[i] - algStamps[i-1]
			algIntervals[i-1] = r
		}
	}
	// D.Ln(algStamps)
	// calculate exponential weighted moving average from intervals
	gewma := ewma.NewMovingAverage()
	for _, x := range algIntervals {
		gewma.Add(float64(x))
	}
	algAv = gewma.Value()
	if algAv != 0 {
		algAdj = algAv / targetTimePerBlock
	}
	return
}

// CalcNextRequiredDifficultyPlan9 returns the consensus difficulty adjustment by processing recent past blocks
func (b *BlockChain) CalcNextRequiredDifficultyPlan9(lastNodeP *BlockNode, algoName string,
	l bool,
) (newTargetBits uint32, adjustment float64, e error) {
	lastNode := lastNodeP
	
	algoVer := fork.GetAlgoVer(algoName, lastNode.height+1)
	ttpb := float64(fork.List[1].Algos[algoName].VersionInterval)
	newTargetBits = fork.SecondPowLimitBits
	const minAvSamples = 3
	adjustment = 1
	var algAdj, allAdj, algAv, allAv float64 = 1, 1, ttpb, fork.P9Average
	if lastNode == nil {
		D.Ln("lastNode is nil")
	}
	// algoInterval := fork.P9Algos[algoname].VersionInterval
	startHeight := fork.List[1].ActivationHeight
	if b.params.Net == wire.TestNet3 {
		startHeight = fork.List[1].TestnetStart
	}
	allStamps := GetAllStamps(startHeight, lastNode)
	last, _, algStamps, algoVer := GetAlgStamps(algoName, startHeight, lastNode)
	if len(allStamps) > minAvSamples {
		allAv, allAdj = GetAll(allStamps)
	}
	if len(algStamps) > minAvSamples {
		algAv, algAdj = GetAlg(algStamps, ttpb)
	}
	bits := fork.SecondPowLimitBits
	if last != nil {
		bits = last.bits
	}
	// D.Ln(
	//	"allAv", allAv,
	//	"fork.P9Average", fork.P9Average,
	//	"allAv/fork.P9Average", allAv/fork.P9Average,
	//	"algAv", algAv,
	//	"algAdj", algAdj,
	//	"allAdj", allAdj)
	
	adjustment = algAdj * allAdj
	
	// thisInterval := fork.P9AlgosNumeric[algoVer].VersionInterval
	// baseInterval := fork.P9AlgosNumeric[5].VersionInterval
	// timeFactor := float64(baseInterval) / float64(thisInterval)
	// // if adjustment < 1 {
	// // if the difficulty is adjusting upwards, accelerate it in proportion with block interval's ratio
	// adjustment *= timeFactor
	// // }
	
	// adjustment *= adjustment
	bigAdjustment := big.NewFloat(adjustment)
	bigOldTarget := big.NewFloat(1.0).SetInt(bits2.CompactToBig(bits))
	bigNewTargetFloat := big.NewFloat(1.0).Mul(bigAdjustment, bigOldTarget)
	newTarget, _ := bigNewTargetFloat.Int(nil)
	if newTarget == nil {
		I.Ln("newTarget is nil ")
		return
	}
	if newTarget.Cmp(&fork.FirstPowLimit) < 0 {
		newTargetBits = bits2.BigToCompact(newTarget)
		// Tracef("newTarget %064x %08x", newTarget, newTargetBits)
	}
	// if l {
	// if lastNode.version == algoVer {
	I.Ln(func() string {
		an := fork.List[1].AlgoVers[algoVer]
		pad := 8 - len(an)
		if pad > 0 {
			an += strings.Repeat(" ", pad)
		}
		factor := 1 / adjustment
		symbol := ">"
		if factor < 1 {
			factor = adjustment
			symbol = "<"
		}
		if factor == 1 {
			symbol = "-"
		}
		isNewest := ""
		if lastNode.version == algoVer {
			isNewest = "*"
		}
		return fmt.Sprintf("%s %s av %s/%2.2f %s %s %08x %08x%s",
			an,
			RightJustify(fmt.Sprintf("%4.2f", algAv), 8),
			RightJustify(fmt.Sprintf("%4.2f", allAv), 7),
			fork.P9Average,
			RightJustify(fmt.Sprintf("%4.2f", factor), 7),
			symbol,
			bits,
			newTargetBits,
			isNewest,
		)
	}(),
	)
	// }
	// }
	return
}

// CalcNextRequiredDifficultyPlan9old calculates the required difficulty for the block after the passed previous block node
// based on the difficulty retarget rules. This function differs from the exported CalcNextRequiredDifficulty in that
// the exported version uses the current best chain as the previous block node while this function accepts any block
// node.
func (b *BlockChain) CalcNextRequiredDifficultyPlan9old(lastNode *BlockNode, algoName string, l bool,
) (newTargetBits uint32, adjustment float64, e error) {
	
	nH := lastNode.height + 1
	newTargetBits = fork.SecondPowLimitBits
	adjustment = 1.0
	if lastNode == nil || b.IsP9HardFork(nH) {
		return
	}
	allTimeAv, allTimeDiv, qhourDiv, hourDiv,
	dayDiv := b.GetCommonP9Averages(lastNode, nH)
	algoVer := fork.GetAlgoVer(algoName, nH)
	since, ttpb, timeSinceAlgo, startHeight, last := b.GetP9Since(lastNode, algoVer)
	if last == nil {
		return
	}
	algDiv := b.GetP9AlgoDiv(allTimeDiv, last, startHeight, algoVer, ttpb)
	adjustment = (allTimeDiv + algDiv + dayDiv + hourDiv + qhourDiv +
		timeSinceAlgo) / 6
	bigAdjustment := big.NewFloat(adjustment)
	bigOldTarget := big.NewFloat(1.0).SetInt(bits2.CompactToBig(last.bits))
	bigNewTargetFloat := big.NewFloat(1.0).Mul(bigAdjustment, bigOldTarget)
	newTarget, _ := bigNewTargetFloat.Int(nil)
	if newTarget == nil {
		I.Ln("newTarget is nil ")
		return
	}
	if newTarget.Cmp(&fork.FirstPowLimit) < 0 {
		newTargetBits = bits2.BigToCompact(newTarget)
		T.F("newTarget %064x %08x", newTarget, newTargetBits)
	}
	if l {
		an := fork.List[1].AlgoVers[algoVer]
		pad := 9 - len(an)
		if pad > 0 {
			an += strings.Repeat(" ", pad)
		}
		D.C(func() string {
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
					since-ttpb*float64(len(fork.List[1].Algos)), timeSinceAlgo*ttpb,
				), 13,
				),
				RightJustify(fmt.Sprintf("%4.4fx", 1/adjustment), 11),
				newTargetBits,
			)
		},
		)
	}
	return
}
