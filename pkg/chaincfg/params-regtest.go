package chaincfg

import (
	"github.com/p9c/matrjoska/pkg/wire"
)

// RegressionTestParams defines the network parameters for the regression test Bitcoin network. Not to be confused with
// the test Bitcoin network (version 3), this network is sometimes simply called "testnet".
var RegressionTestParams = Params{
	Name:        "regtest",
	Net:         wire.TestNet,
	DefaultPort: "31047",
	DNSSeeds:    []DNSSeed{},
	// Chain parameters
	GenesisBlock:     &regTestGenesisBlock,
	GenesisHash:      &regTestGenesisHash,
	PowLimit:         regressionPowLimit,
	PowLimitBits:     0x207fffff,
	CoinbaseMaturity: 100,
	// BIP0034Height:            100000000, // Not active - Permit ver 1 blocks
	// BIP0065Height:            100000000, // Used by regression tests
	// BIP0066Height:            100000000, // Used by regression tests
	SubsidyReductionInterval: 150,
	TargetTimespan:           30000, // 14 days
	TargetTimePerBlock:       300,   // 5 minutes
	RetargetAdjustmentFactor: 2,     // 50% less, 200% more
	ReduceMinDifficulty:      true,
	MinDiffReductionTime:     300 * 2,
	GenerateSupported:        true,
	// Checkpoints ordered from oldest to newest.
	Checkpoints: nil,
	// Consensus rule change deployments.
	//
	// The miner confirmation window is defined as:
	//   target proof of work timespan / target proof of work spacing
	RuleChangeActivationThreshold: 108, // 75%  of MinerConfirmationWindow
	MinerConfirmationWindow:       144,
	// Deployments: [DefinedDeployments]ConsensusDeployment{
	// 	DeploymentTestDummy: {
	// 		BitNumber:  28,
	// 		StartTime:  0,             // Always available for vote
	// 		ExpireTime: math.MaxInt64, // Never expires
	// 	},
	// 	DeploymentCSV: {
	// 		BitNumber:  0,
	// 		StartTime:  0,             // Always available for vote
	// 		ExpireTime: math.MaxInt64, // Never expires
	// 	},
	// 	DeploymentSegwit: {
	// 		BitNumber:  1,
	// 		StartTime:  0,             // Always available for vote
	// 		ExpireTime: math.MaxInt64, // Never expires.
	// 	},
	// },
	// Mempool parameters
	RelayNonStdTxs: true,
	// Human-readable part for Bech32 encoded segwit addresses, as defined in
	// // BIP 173.
	// Bech32HRPSegwit: "bcrt", // always bcrt for reg test net
	// Address encoding magics
	PubKeyHashAddrID: 0x00,
	ScriptHashAddrID: 0x05,
	PrivateKeyID:     0x80,
	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x04, 0x35, 0x83, 0x94}, // starts with tprv
	HDPublicKeyID:  [4]byte{0x04, 0x35, 0x87, 0xcf}, // starts with tpub
	// BIP44 coin type used in the hierarchical deterministic path for address generation.
	HDCoinType: 1,
	// Parallelcoin specific difficulty adjustment parameters
	Interval:                Interval,
	AveragingInterval:       10, // Extend to target timespan to adjust better to hashpower (30000/300=100) post hardfork
	AveragingTargetTimespan: AveragingTargetTimespan,
	MaxAdjustDown:           MaxAdjustDown,
	MaxAdjustUp:             MaxAdjustUp,
	TargetTimespanAdjDown:   AveragingTargetTimespan * (Interval + MaxAdjustDown) / Interval,
	MinActualTimespan:       AveragingTargetTimespan * (Interval - MaxAdjustUp) / Interval,
	MaxActualTimespan:       AveragingTargetTimespan * (Interval + MaxAdjustDown) / Interval,
	ScryptPowLimit:          &scryptPowLimit,
	ScryptPowLimitBits:      ScryptPowLimitBits,
	RPCClientPort:           "31048",
	WalletRPCServerPort:     "31046",
}
