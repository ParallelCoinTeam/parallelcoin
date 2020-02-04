package blockchain

import (
	"github.com/p9c/pod/pkg/chain/fork"
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
			nTB[v], err = b.CalcNextRequiredDifficultyHalcyon(0, lastNode, i, true)
		}
		return &nTB, nil
	}
	for i := range fork.List[1].Algos {
		v := fork.List[1].Algos[i].Version
		nTB[v], _, err = b.CalcNextRequiredDifficultyPlan9(0, lastNode, i, true)
	}
	newTargetBits = &nTB
	return
}
