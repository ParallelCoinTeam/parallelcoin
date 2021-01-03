// Package fork handles tracking the hard fork status and is used to determine which consensus rules apply on a block
package fork

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"math/rand"
	"sort"
	"time"
)

const (
	Scrypt  = "scrypt"
	SHA256d = "sha256d"
)

// AlgoParams are the identifying block version number and their minimum target bits
type AlgoParams struct {
	Version         int32
	MinBits         uint32
	AlgoID          uint32
	VersionInterval int
}

// HardForks is the details related to a hard fork, number, name and activation height
type HardForks struct {
	Number             uint32
	ActivationHeight   int32
	Name               string
	Algos              map[string]AlgoParams
	AlgoVers           map[int32]string
	TargetTimePerBlock int32
	AveragingInterval  int32
	TestnetStart       int32
}

func init() {
	Trace("running fork data init")
	for i := range P9AlgosNumeric {
		List[1].AlgoVers[i] = fmt.Sprintf("Div%d", P9AlgosNumeric[i].VersionInterval)
	}
	for i, v := range P9AlgoVers {
		List[1].Algos[v] = P9AlgosNumeric[i]
	}
	AlgoSlices = append(AlgoSlices, AlgoSpecs{})
	for i := range Algos {
		AlgoSlices[0] = append(AlgoSlices[0], AlgoSpec{
			List[0].Algos[i].Version,
			i,
		})
	}
	AlgoSlices = append(AlgoSlices, AlgoSpecs{})
	for i := range P9Algos {
		AlgoSlices[1] = append(AlgoSlices[1], AlgoSpec{
			List[1].Algos[i].Version,
			i,
		})
	}
	sort.Sort(AlgoSlices[0])
	sort.Sort(AlgoSlices[1])
	Debug(P9AlgoVers)
	baseVersionName := AlgoSlices[1][0].Name
	baseVersionInterval := float64(P9Algos[baseVersionName].VersionInterval)
	Debug(baseVersionName, baseVersionInterval)
	P9Average = 0
	for _, i := range AlgoSlices[1] {
		vi := float64(P9Algos[i.Name].VersionInterval)
		p9a := baseVersionInterval / vi
		P9Average += p9a
		// Tracef("P9Average %4.4f %4.4f %d %4.4f", p9a, P9Average, IntervalBase, vi)
	}
	Debug(P9Average)
	P9Average = baseVersionInterval / P9Average
	Debug(P9Average)
}

type AlgoSpec struct {
	Version int32
	Name    string
}
type AlgoSpecs []AlgoSpec

func (a AlgoSpecs) Len() int {
	return len(a)
}

func (a AlgoSpecs) Less(i, j int) bool {
	return a[i].Version < a[j].Version
}

func (a AlgoSpecs) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

var (
	AlgoSlices []AlgoSpecs
	// AlgoVers is the lookup for pre hardfork
	//
	AlgoVers = map[int32]string{
		2:   SHA256d,
		514: Scrypt,
	}
	// Algos are the specifications identifying the algorithm used in the
	// block proof
	Algos = map[string]AlgoParams{
		AlgoVers[2]: {
			Version: 2,
			MinBits: MainPowLimitBits,
		},
		AlgoVers[514]: {
			Version: 514,
			MinBits: MainPowLimitBits,
			AlgoID:  1,
		},
	}
	// FirstPowLimit is
	FirstPowLimit = func() big.Int {
		mplb, _ := hex.DecodeString(
			"0fffff0000000000000000000000000000000000000000000000000000000000")
		return *big.NewInt(0).SetBytes(mplb)
	}()
	// FirstPowLimitBits is
	FirstPowLimitBits = BigToCompact(&FirstPowLimit)
	// IsTestnet is set at startup here to be accessible to all other libraries
	IsTestnet bool
	// List is the list of existing hard forks and when they activate
	List = []HardForks{
		{
			Number:             0,
			Name:               "Halcyon days",
			ActivationHeight:   0,
			Algos:              Algos,
			AlgoVers:           AlgoVers,
			TargetTimePerBlock: 300,
			AveragingInterval:  10, // 50 minutes
			TestnetStart:       0,
		},
		{
			Number:             1,
			Name:               "Plan 9 from Crypto Space",
			ActivationHeight:   2500000,
			Algos:              P9Algos,
			AlgoVers:           P9AlgoVers,
			TargetTimePerBlock: 36,
			AveragingInterval:  3600,
			TestnetStart:       1,
		},
	}
	// P9AlgoVers is the lookup for after 1st hardfork
	P9AlgoVers = make(map[int32]string)
	
	P9PrimeSequence = []int{2, 3, 5, 7, 11, 13, 17, 19, 23}
	IntervalDivisor = 1
	IntervalBase    = 101
	// P9Algos is the algorithm specifications after the hard fork
	P9Algos        = make(map[string]AlgoParams)
	P9AlgosNumeric = map[int32]AlgoParams{
		5:  {5, FirstPowLimitBits, 0, IntervalBase * P9PrimeSequence[0] / IntervalDivisor},  // 2
		6:  {6, FirstPowLimitBits, 1, IntervalBase * P9PrimeSequence[1] / IntervalDivisor},  // 3
		7:  {7, FirstPowLimitBits, 2, IntervalBase * P9PrimeSequence[2] / IntervalDivisor},  // 5
		8:  {8, FirstPowLimitBits, 3, IntervalBase * P9PrimeSequence[3] / IntervalDivisor},  // 7
		9:  {9, FirstPowLimitBits, 4, IntervalBase * P9PrimeSequence[4] / IntervalDivisor},  // 11
		10: {10, FirstPowLimitBits, 5, IntervalBase * P9PrimeSequence[5] / IntervalDivisor}, // 13
		11: {11, FirstPowLimitBits, 7, IntervalBase * P9PrimeSequence[7] / IntervalDivisor}, // 17
		12: {12, FirstPowLimitBits, 6, IntervalBase * P9PrimeSequence[6] / IntervalDivisor}, // 19
		13: {13, FirstPowLimitBits, 8, IntervalBase * P9PrimeSequence[8] / IntervalDivisor}, // 23
	}
	
	P9Average float64
	
	// SecondPowLimit is
	SecondPowLimit = func() big.Int {
		mplb, _ := hex.DecodeString(
			// "01f1f1f1f1f1f1f1f1f1f1f1f1f1f1f1f1f1f1f1f1f1f1f1f1f1f1f1f1f1f1f1")
			"0fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
		return *big.NewInt(0).SetBytes(mplb)
	}()
	SecondPowLimitBits = BigToCompact(&SecondPowLimit)
	MainPowLimit       = func() big.Int {
		mplb, _ := hex.DecodeString(
			"00000fffff000000000000000000000000000000000000000000000000000000")
		return *big.NewInt(0).SetBytes(mplb)
	}()
	MainPowLimitBits = BigToCompact(&MainPowLimit)
)

// GetAlgoID returns the 'algo_id' which in pre-hardfork is not the same as the block version number, but is afterwards
func GetAlgoID(algoname string, height int32) uint32 {
	if GetCurrent(height) > 1 {
		return P9Algos[algoname].AlgoID
	}
	return Algos[algoname].AlgoID
}

// GetAlgoName returns the string identifier of an algorithm depending on
// hard fork activation status
func GetAlgoName(algoVer int32, height int32) (name string) {
	hf := GetCurrent(height)
	var ok bool
	name, ok = List[hf].AlgoVers[algoVer]
	if hf < 1 && !ok {
		name = SHA256d
	}
	// INFO("GetAlgoName", algoVer, height, name}
	return
}

// GetRandomVersion returns a random version relevant to the current hard fork state and height
func GetRandomVersion(height int32) int32 {
	rand.Seed(time.Now().UnixNano())
	return int32(rand.Intn(len(List[GetCurrent(height)].Algos)) + 5)
}

// GetAlgoVer returns the version number for a given algorithm (by string name) at a given height. If "random" is given,
// a random number is taken from the system secure random source (for randomised cpu mining)
func GetAlgoVer(name string, height int32) (version int32) {
	hf := GetCurrent(height)
	n := AlgoSlices[hf][0].Name
	// DEBUG("GetAlgoVer", name, height, hf, n)
	if _, ok := List[hf].Algos[name]; ok {
		n = name
	}
	version = List[hf].Algos[n].Version
	return
}

// GetAveragingInterval returns the active block interval target based on hard fork status
func GetAveragingInterval(height int32) (r int32) {
	r = List[GetCurrent(height)].AveragingInterval
	return
}

// GetCurrent returns the hardfork number code
func GetCurrent(height int32) (curr int) {
	// Trace("istestnet", IsTestnet)
	if IsTestnet {
		for i := range List {
			if height >= List[i].TestnetStart {
				curr = i
			}
		}
	} else {
		for i := range List {
			if height >= List[i].ActivationHeight {
				curr = i
			}
		}
	}
	return
}

// GetMinBits returns the minimum diff bits based on height and testnet
func GetMinBits(algoname string, height int32) (mb uint32) {
	curr := GetCurrent(height)
	// Trace("GetMinBits", algoname, height, curr, List[curr].Algos)
	mb = List[curr].Algos[algoname].MinBits
	// TraceF("minbits %08x, %d", mb, mb)
	return
}

// GetMinDiff returns the minimum difficulty in uint256 form
func GetMinDiff(algoname string, height int32) (md *big.Int) {
	// Trace("GetMinDiff", algoname)
	minbits := GetMinBits(algoname, height)
	// TraceF("mindiff minbits %08x", minbits)
	return CompactToBig(minbits)
}

// GetTargetTimePerBlock returns the active block interval target based on hard fork status
func GetTargetTimePerBlock(height int32) (r int64) {
	r = int64(List[GetCurrent(height)].TargetTimePerBlock)
	return
}
