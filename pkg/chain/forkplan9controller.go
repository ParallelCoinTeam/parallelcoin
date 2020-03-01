package blockchain

import (
	"sort"
	
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

type Algo struct {
	Name   string
	Params fork.AlgoParams
}

type AlgoList []Algo

func (al AlgoList) Len() int {
	return len(al)
}

func (al AlgoList) Less(i, j int) bool {
	return al[i].Params.Version < al[j].Params.Version
}

func (al AlgoList) Swap(i, j int) {
	al[i], al[j] = al[j], al[i]
}

// CalcNextRequiredDifficultyPlan9Controller returns all of the algorithm
// difficulty targets for sending out with the other pieces required to
// construct a block, as these numbers are generated from block timestamps
func (b *BlockChain) CalcNextRequiredDifficultyPlan9Controller(
	lastNode *BlockNode) (newTargetBits map[int32]uint32, err error) {
	nH := lastNode.height + 1
	currFork := fork.GetCurrent(nH)
	nTB := make(map[int32]uint32)
	switch currFork {
	case 0:
		for i := range fork.List[0].Algos {
			v := fork.List[0].Algos[i].Version
			nTB[v], err = b.CalcNextRequiredDifficultyHalcyon(0, lastNode, i, true)
		}
		return nTB, nil
	case 1:
		if b.DifficultyHeight.Load() != nH {
			b.DifficultyHeight.Store(nH)
			currFork := fork.GetCurrent(nH)
			algos := make(AlgoList, len(fork.List[currFork].Algos))
			var counter int
			for i := range fork.List[1].Algos {
				algos[counter] = Algo{
					Name:   i,
					Params: fork.List[currFork].Algos[i],
				}
				counter++
			}
			sort.Sort(algos)
			for _, v := range algos {
				nTB[v.Params.Version], _, err = b.CalcNextRequiredDifficultyPlan9(0, lastNode, v.Name, true)
			}
			newTargetBits = nTB
			// log.SPEW(newTargetBits)
		} else {
			newTargetBits = b.DifficultyBits.Load().(map[int32]uint32)
		}
		return
	}
	log.TRACE("should not fall through here")
	return
}
