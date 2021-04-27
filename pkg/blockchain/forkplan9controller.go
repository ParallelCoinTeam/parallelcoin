package blockchain

import (
	"github.com/p9c/matrjoska/pkg/chainhash"
	"github.com/p9c/matrjoska/pkg/fork"
	"sort"
)

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

type Diffs map[int32]uint32

type Merkles map[int32]*chainhash.Hash

// CalcNextRequiredDifficultyPlan9Controller returns all of the algorithm difficulty targets for sending out with the
// other pieces required to construct a block, as these numbers are generated from block timestamps
func (b *BlockChain) CalcNextRequiredDifficultyPlan9Controller(lastNode *BlockNode) (
	diffs Diffs, e error,
) {
	nH := lastNode.height + 1
	currFork := fork.GetCurrent(nH)
	nTB := make(Diffs)
	switch currFork {
	case 0:
		for i := range fork.List[0].Algos {
			v := fork.List[currFork].Algos[i].Version
			nTB[v], e = b.CalcNextRequiredDifficultyHalcyon(lastNode, i, true)
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
				nTB[v.Params.Version], _, e = b.CalcNextRequiredDifficultyPlan9(lastNode, v.Name, true)
			}
			diffs = nTB
			// Traces(diffs)
		} else {
			diffs = b.DifficultyBits.Load().(Diffs)
		}
		return
	}
	F.Ln("should not fall through here")
	return
}
