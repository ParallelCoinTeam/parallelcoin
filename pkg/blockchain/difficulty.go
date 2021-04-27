package blockchain

import (
	"encoding/hex"
	bits2 "github.com/p9c/matrjoska/pkg/bits"
	"github.com/p9c/matrjoska/pkg/fork"
	"math/big"
	"strings"
	"time"
	
	"github.com/p9c/matrjoska/pkg/chainhash"
)

var (
	// ScryptPowLimit is
	ScryptPowLimit = scryptPowLimit
	// ScryptPowLimitBits is
	ScryptPowLimitBits = bits2.BigToCompact(&scryptPowLimit)
	// bigOne is 1 represented as a big.Int. It is defined here to avoid the overhead of creating it multiple times.
	bigOne = big.NewInt(1)
	// oneLsh256 is 1 shifted left 256 bits. It is defined here to avoid the overhead of creating it multiple times.
	oneLsh256      = new(big.Int).Lsh(bigOne, 256)
	scryptPowLimit = func() big.Int {
		mplb, _ := hex.DecodeString(
			"000000039fcaa04ac30b6384471f337748ef5c87c7aeffce5e51770ce6283137,",
		)
		return *big.NewInt(0).SetBytes(mplb) // AllOnes.Rsh(&AllOnes, 0)
	}()
)

// CalcNextRequiredDifficulty calculates the required difficulty for the block after the end of the current best chain
// based on the difficulty retarget rules. This function is safe for concurrent access.
func (b *BlockChain) CalcNextRequiredDifficulty(algo string) (difficulty uint32, e error) {
	b.ChainLock.Lock()
	difficulty, e = b.CalcNextRequiredDifficultyFromNode(
		b.BestChain.
			Tip(), algo, false,
	)
	// F.Ln("CalcNextRequiredDifficulty", difficulty)
	b.ChainLock.Unlock()
	return
}

// calcEasiestDifficulty calculates the easiest possible difficulty that a block can have given starting difficulty bits
// and a duration.
//
// It is mainly used to verify that claimed proof of work by a block is sane as compared to a known good checkpoint.
func (b *BlockChain) calcEasiestDifficulty(bits uint32, duration time.Duration) uint32 {
	// Convert types used in the calculations below.
	durationVal := int64(duration / time.Second)
	adjustmentFactor := big.NewInt(b.params.RetargetAdjustmentFactor)
	// Since easier difficulty equates to higher numbers, the easiest difficulty for a given duration is the largest
	// value possible given the number of retargets for the duration and starting difficulty multiplied by the max
	// adjustment factor.
	newTarget := bits2.CompactToBig(bits)
	for durationVal > 0 && newTarget.Cmp(b.params.PowLimit) < 0 {
		newTarget.Mul(newTarget, adjustmentFactor)
		durationVal -= b.maxRetargetTimespan
	}
	// Limit new value to the proof of work limit.
	if newTarget.Cmp(b.params.PowLimit) > 0 {
		newTarget.Set(b.params.PowLimit)
	}
	return bits2.BigToCompact(newTarget)
}

// CalcNextRequiredDifficultyFromNode calculates the required difficulty for the block after the passed previous block node
// based on the difficulty retarget rules.
//
// This function differs from the exported CalcNextRequiredDifficulty in that the exported version uses the current best
// chain as the previous block node while this function accepts any block node.
func (b *BlockChain) CalcNextRequiredDifficultyFromNode(lastNode *BlockNode, algoname string, l bool,) (
	newTargetBits uint32,
	e error,
) {
	nH := lastNode.height + 1
	cF := fork.GetCurrent(nH)
	newTargetBits = fork.GetMinBits(algoname, nH)
	// Tracef("CalcNextRequiredDifficultyFromNode %08x", newTargetBits)
	switch cF {
	// Legacy difficulty adjustment
	case 0:
		// F.Ln("before hardfork")
		return b.CalcNextRequiredDifficultyHalcyon(lastNode, algoname, l)
	// Plan 9 from Crypto Space
	case 1:
		bits, ok := lastNode.Diffs.Load().(Diffs)
		if bits == nil || !ok {
			lastNode.Diffs.Store(make(Diffs))
		}
		version := fork.GetAlgoVer(algoname, lastNode.height+1)
		if bits[version] == 0 {
			bits, e = b.CalcNextRequiredDifficultyPlan9Controller(lastNode)
			if e != nil {
				E.Ln(e)
				return
			}
			// D.Ln(bits, reflect.TypeOf(bits))
			b.DifficultyBits.Store(bits)
			// D.F("got difficulty %d %08x %+v", version, (*b.DifficultyBits)[version], *bits)
		}
		newTargetBits = bits[version]
		return
	}
	return
}

// RightJustify takes a string and right justifies it by a width or crops it
func RightJustify(s string, w int) string {
	sw := len(s)
	diff := w - sw
	if diff > 0 {
		s = strings.Repeat(" ", diff) + s
	} else if diff < 0 {
		s = s[:w]
	}
	return s
}

// CalcWork calculates a work value from difficulty bits.
// Bitcoin increases the difficulty for generating a block by decreasing the
// value which the generated hash must be less than.
// This difficulty target is stored in each block header using a compact
// representation as described in the documentation for CompactToBig.
// The main chain is selected by choosing the chain that has the most proof
// of work (highest difficulty).
// Since a lower target difficulty value equates to higher actual difficulty,
// the work value which will be accumulated must be the inverse of the
// difficulty.  Also,
// in order to avoid potential division by zero and really small floating
// point numbers, the result adds 1 to the denominator and multiplies the
// numerator by 2^256.
func CalcWork(bits uint32, height int32, algover int32) *big.Int {
	// Return a work value of zero if the passed difficulty bits represent a negative number. Note this should not
	// happen in practice with valid blocks, but an invalid block could trigger it.
	difficultyNum := bits2.CompactToBig(bits)
	// To make the difficulty values correlate to number of hash operations, multiply this difficulty base by the
	// nanoseconds/hash figures in the fork algorithms list
	if difficultyNum.Sign() <= 0 {
		return big.NewInt(0)
	}
	denominator := new(big.Int).Add(difficultyNum, bigOne)
	r := new(big.Int).Div(oneLsh256, denominator)
	return r
}

// HashToBig converts a chainhash.Hash into a big. Int that can be used to perform math comparisons.
func HashToBig(hash *chainhash.Hash) *big.Int {
	// A Hash is in little-endian, but the big package wants the bytes in big-endian, so reverse them.
	buf := *hash
	blen := len(buf)
	for i := 0; i < blen/2; i++ {
		buf[i], buf[blen-1-i] = buf[blen-1-i], buf[i]
	}
	// buf := hash.CloneBytes()
	return new(big.Int).SetBytes(buf[:])
}
