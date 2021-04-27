package chaincfg

import (
	"encoding/hex"
	"errors"
	"math/big"
	"time"
	
	"github.com/p9c/matrjoska/pkg/chainhash"
	"github.com/p9c/matrjoska/pkg/wire"
)

var (
	// ErrDuplicateNet describes an error where the parameters for a Bitcoin network could not be set due to the network
	// already being a standard network or previously-registered into this package.
	ErrDuplicateNet = errors.New("duplicate Bitcoin network")
	// ErrUnknownHDKeyID describes an error where the provided id which is intended to identify the network for a
	// hierarchical deterministic private extended key is not registered.
	ErrUnknownHDKeyID    = errors.New("unknown hd private extended key bytes")
	registeredNets       = make(map[wire.BitcoinNet]struct{})
	pubKeyHashAddrIDs    = make(map[byte]struct{})
	scriptHashAddrIDs    = make(map[byte]struct{})
	bech32SegwitPrefixes = make(map[string]struct{})
	hdPrivToPubKeyIDs    = make(map[[4]byte][]byte)
	// allOnes is 32 bytes of 0xff, the maximum target
	allOnes = func() big.Int {
		b := big.NewInt(1)
		t := make([]byte, 32)
		for i := range t {
			t[i] = ^byte(0)
		}
		b.SetBytes(t)
		return *b
	}()
	// mainPowLimit is the highest proof of work value a Parallelcoin block can have for the main network.
	mainPowLimit = func() big.Int {
		mplb, _ := hex.DecodeString("00000fffff000000000000000000000000000000000000000000000000000000")
		return *big.NewInt(0).SetBytes(mplb) // AllOnes.Rsh(&AllOnes, 0)
	}()
	// MainPowLimit is the pre-hardfork pow limit
	MainPowLimit = mainPowLimit
	// MainPowLimitBits is the bits version of the above
	MainPowLimitBits = BigToCompact(&MainPowLimit)
	scryptPowLimit   = func() big.Int {
		mplb, _ := hex.DecodeString("0000fffff0000000000000000000000000000000000000000000000000000000")
		return *big.NewInt(0).SetBytes(mplb) // AllOnes.Rsh(&AllOnes, 0)
	}()
	// ScryptPowLimit is the pre-hardfork maximum hash for Scrypt algorithm
	ScryptPowLimit = scryptPowLimit
	// ScryptPowLimitBits is the bits version of the above
	ScryptPowLimitBits = BigToCompact(&scryptPowLimit)
	// regressionPowLbimit is the highest proof of work value a Bitcoin block can have for the regression test network.
	// It is the value 2^255 - 1, all ones, 256 bits.
	regressionPowLimit = &allOnes
	// testnetBits = ScryptPowLimitBits testNet3PowLimit = ScryptPowLimit simNetPowLimit is the highest proof of work
	// value a Bitcoin block can have for the simulation test network. It is the value 2^255 - 1, all ones, 256 bits.
	simNetPowLimit = &allOnes
	// Interval is the number of blocks in the averaging window
	Interval int64 = 100
	// MaxAdjustDown is the percentage hard limit for downwards difficulty adjustment (ie 90%)
	MaxAdjustDown int64 = 10
	// MaxAdjustUp is the percentage hard limit for upwards (ie 120%)
	MaxAdjustUp int64 = 20
	// TargetTimePerBlock is the pre hardfork target time for blocks
	TargetTimePerBlock int64 = 300
	// AveragingInterval is the number of blocks to average (per algorithm)
	AveragingInterval int64 = 10
	// AveragingTargetTimespan is how many seconds for the averaging target interval
	AveragingTargetTimespan = TargetTimePerBlock * AveragingInterval
	// TargetTimespan is the base for adjustment
	TargetTimespan = Interval * TargetTimePerBlock
	// TestnetInterval is the number of blocks in the averaging window
	TestnetInterval int64 = 100
	// TestnetMaxAdjustDown is the percentage hard limit for downwards difficulty adjustment (ie 90%)
	TestnetMaxAdjustDown int64 = 10
	// TestnetMaxAdjustUp is the percentage hard limit for upwards (ie 120%)
	TestnetMaxAdjustUp int64 = 20
	// TestnetTargetTimePerBlock is the pre hardfork target time for blocks
	TestnetTargetTimePerBlock int64 = 9
	// TestnetAveragingInterval is the number of blocks to average (per algorithm)
	TestnetAveragingInterval int64 = 1600
	// TestnetAveragingTargetTimespan is how many seconds for the averaging target interval
	TestnetAveragingTargetTimespan = TestnetTargetTimePerBlock * TestnetAveragingInterval
	// TestnetTargetTimespan is the base for adjustment
	TestnetTargetTimespan = TestnetInterval * TestnetTargetTimePerBlock
)

// Checkpoint identifies a known good point in the block chain. Using checkpoints allows a few optimizations for old
// blocks during initial download and also prevents forks from old blocks. Each checkpoint is selected based upon
// several factors. See the documentation for blockchain.IsCheckpointCandidate for details on the selection criteria.
type Checkpoint struct {
	Height int32
	Hash   *chainhash.Hash
}

// DNSSeed identifies a DNS seed.
type DNSSeed struct {
	// Host defines the hostname of the seed.
	Host string
	// HasFiltering defines whether the seed supports filtering by service flags (wire.ServiceFlag).
	HasFiltering bool
}

// // ConsensusDeployment defines details related to a specific consensus rule change that is voted in. This is part of
// // BIP0009.
// type ConsensusDeployment struct {
// 	// BitNumber defines the specific bit number within the block version this particular soft-fork deployment refers
// 	// to.
// 	BitNumber uint8
// 	// StartTime is the median block time after which voting on the deployment starts.
// 	StartTime uint64
// 	// ExpireTime is the median block time after which the attempted deployment expires.
// 	ExpireTime uint64
// }

// // Constants that define the deployment offset in the deployments field of the
// // parameters for each deployment. This is useful to be able to get the details
// // of a specific deployment by name.
// const (
// 	// DeploymentTestDummy defines the rule change deployment ID for testing
// 	// purposes.
// 	DeploymentTestDummy = iota
// 	// DeploymentCSV defines the rule change deployment ID for the CSV soft-fork
// 	// package. The CSV package includes the deployment of BIPS 68, 112, and 113.
// 	DeploymentCSV
// 	// DeploymentSegwit defines the rule change deployment ID for the Segregated
// 	// Witness (segwit) soft-fork package. The segwit package includes the
// 	// deployment of BIPS 141, 142, 144, 145, 147 and 173.
// 	DeploymentSegwit
// 	// DefinedDeployments is the number of currently defined deployments.
// 	//
// 	// NOTE: DefinedDeployments must always come last since it is used to determine
// 	// how many defined deployments there currently are.
// 	DefinedDeployments
// )

// Params defines a Bitcoin network by its parameters. These parameters may be
// used by Bitcoin applications to differentiate networks as well as addresses
// and keys for one network from those intended for use on another network.
type Params struct {
	// Name defines a human-readable identifier for the network.
	Name string
	// Net defines the magic bytes used to identify the network.
	Net wire.BitcoinNet
	// DefaultPort defines the default peer-to-peer port for the network.
	DefaultPort         string
	RPCClientPort       string
	WalletRPCServerPort string
	// DNSSeeds defines a list of DNS seeds for the network that are used
	// as one method to discover peers.
	DNSSeeds []DNSSeed
	// GenesisBlock defines the first block of the chain.
	GenesisBlock *wire.Block
	// GenesisHash is the starting block hash.
	GenesisHash *chainhash.Hash
	// PowLimit defines the highest allowed proof of work value for a // as a uint256.
	PowLimit *big.Int
	// PowLimitBits defines the highest allowed proof of work value for a block in compact form.
	PowLimitBits uint32
	// // These fields define the block heights at which the specified softfork BIP became active.
	// BIP0034Height int32
	// BIP0065Height int32
	// BIP0066Height int32
	// CoinbaseMaturity is the number of blocks required before newly mined coins (coinbase transactions) can be spent.
	CoinbaseMaturity uint16
	// SubsidyReductionInterval is the interval of blocks before the subsidy is reduced.
	SubsidyReductionInterval int32
	// TargetTimespan is the desired amount of time that should elapse before the block difficulty requirement is
	// examined to determine how it should be changed in order to maintain the desired block generation rate.
	TargetTimespan int64
	// TargetTimePerBlock is the desired amount of time to generate each block. Same as TargetSpacing in legacy client.
	TargetTimePerBlock int64
	// RetargetAdjustmentFactor is the adjustment factor used to limit the minimum and maximum amount of adjustment that
	// can occur between difficulty retargets.
	RetargetAdjustmentFactor int64
	// ReduceMinDifficulty defines whether the network should reduce the minimum required difficulty after a long enough
	// period of time has passed without finding a block. This is really only useful for test networks and should not be
	// set on a main network.
	ReduceMinDifficulty bool
	// MinDiffReductionTime is the amount of time after which the minimum required difficulty should be reduced when a
	// block hasn't been found. NOTE: This only applies if ReduceMinDifficulty is true.
	MinDiffReductionTime time.Duration
	// GenerateSupported specifies whether or not CPU mining is allowed.
	GenerateSupported bool
	// Checkpoints ordered from oldest to newest.
	Checkpoints []Checkpoint
	// These fields are related to voting on consensus rule changes as defined by BIP0009.
	//
	// RuleChangeActivationThreshold is the number of blocks in a threshold state retarget window for which a positive
	// vote for a rule change must be cast in order to lock in a rule change. It should typically be 95% for the main
	// network and 75% for test networks.
	RuleChangeActivationThreshold uint32
	// MinerConfirmationWindow is the number of blocks in each threshold state retarget window.
	MinerConfirmationWindow uint32
	// // Deployments define the specific consensus rule changes to be voted on.
	// Deployments [DefinedDeployments]ConsensusDeployment
	// Mempool parameters
	RelayNonStdTxs bool
	// // Human-readable part for Bech32 encoded segwit addresses, as defined in BIP 173.
	// Bech32HRPSegwit string
	// Address encoding magics
	PubKeyHashAddrID byte // First byte of a P2PKH address
	ScriptHashAddrID byte // First byte of a P2SH address
	PrivateKeyID     byte // First byte of a WIF private key
	// WitnessPubKeyHashAddrID byte // First byte of a P2WPKH address
	// WitnessScriptHashAddrID byte // First byte of a P2WSH address
	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID [4]byte
	HDPublicKeyID  [4]byte
	// BIP44 coin type used in the hierarchical deterministic path for address generation.
	HDCoinType uint32
	// Parallelcoin specific difficulty adjustment parameters
	Interval                int64
	AveragingInterval       int64
	AveragingTargetTimespan int64
	MaxAdjustDown           int64
	MaxAdjustUp             int64
	TargetTimespanAdjDown   int64
	MinActualTimespan       int64
	MaxActualTimespan       int64
	// PowLimit defines the highest allowed proof of work value for a scrypt block as a uint256.
	ScryptPowLimit     *big.Int
	ScryptPowLimitBits uint32
}
