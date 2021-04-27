package chaincfg

import (
	"github.com/p9c/matrjoska/pkg/fork"
	"github.com/p9c/matrjoska/pkg/wire"
)

// TestNet3Params defines the network parameters for the test Bitcoin network (version 3). Not to be confused with the
// regression test network, this network is sometimes simply called "testnet".
var TestNet3Params = Params{
	Name:        "testnet",
	Net:         wire.TestNet3,
	DefaultPort: "21047",
	DNSSeeds: []DNSSeed{
		// {"seed3.parallelcoin.io", true},
	},
	// Chain parameters
	GenesisBlock: &testNet3GenesisBlock,
	GenesisHash:  &testNet3GenesisHash,
	PowLimit:     &fork.SecondPowLimit,    // fork&testNet3PowLimit,
	PowLimitBits: fork.SecondPowLimitBits, // testnetBits,
	// BIP0034Height:            math.MaxInt32,                       // 0000000023b3a96d3484e5abb3755c413e7d41500f8e2a5c3f0dd01299cd8ef8
	// BIP0065Height:            math.MaxInt32,                       // 00000000007f6655f22f98e72ed80d8b06dc761d5da09df0fa1dc4be4f861eb6
	// BIP0066Height:            math.MaxInt32,                       // 000000002104c8c45e99a8853285a3b592602a3ccde2b832481da85e9e4ba182
	CoinbaseMaturity:         9,
	SubsidyReductionInterval: 250000,
	TargetTimespan:           TestnetTargetTimespan,
	TargetTimePerBlock:       TestnetTargetTimePerBlock,
	RetargetAdjustmentFactor: 2,
	ReduceMinDifficulty:      false,
	MinDiffReductionTime:     0, // time.Minute * 10, // TargetTimePerBlock * 2
	GenerateSupported:        true,
	// Checkpoints ordered from oldest to newest.
	Checkpoints: []Checkpoint{
		// {546, newHashFromStr("000000002a936ca763904c3c35fce2f3556c559c0214345d31b1bcebf76acb70")},
	},
	// Consensus rule change deployments.
	//
	// The miner confirmation window is defined as:
	//   target proof of work timespan / target proof of work spacing
	RuleChangeActivationThreshold: 2, // 75% of MinerConfirmationWindow
	MinerConfirmationWindow:       2016,
	// Deployments: [DefinedDeployments]ConsensusDeployment{
	// 	DeploymentTestDummy: {
	// 		BitNumber:  28,
	// 		StartTime:  math.MaxInt64, // January 1, 2008 UTC
	// 		ExpireTime: math.MaxInt64, // December 31, 2008 UTC
	// 	},
	// 	DeploymentCSV: {
	// 		BitNumber:  29,
	// 		StartTime:  math.MaxInt64, // March 1st, 2016
	// 		ExpireTime: math.MaxInt64, // May 1st, 2017
	// 	},
	// 	DeploymentSegwit: {
	// 		BitNumber:  29,
	// 		StartTime:  math.MaxInt64, // always available
	// 		ExpireTime: math.MaxInt64, // never expires
	// 	},
	// },
	// Mempool parameters
	RelayNonStdTxs: true,
	// // Human-readable part for Bech32 encoded segwit addresses, as defined in BIP 173.
	// Bech32HRPSegwit: "t9", // always tb for test net
	// Address encoding magics
	PubKeyHashAddrID: 18,  // starts with m or n
	ScriptHashAddrID: 188, // starts with 2
	// WitnessPubKeyHashAddrID: 0x03, // starts with QW
	// WitnessScriptHashAddrID: 0x28, // starts with T7n
	PrivateKeyID: 239, // starts with 9 (uncompressed) or c (compressed)
	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x04, 0x35, 0x83, 0x94}, // starts with tprv
	HDPublicKeyID:  [4]byte{0x04, 0x35, 0x87, 0xcf}, // starts with tpub
	// BIP44 coin type used in the hierarchical deterministic path for address generation.
	HDCoinType: 1,
	// Parallelcoin specific difficulty adjustment parameters
	Interval:                TestnetInterval,
	AveragingInterval:       TestnetAveragingInterval, // Extend to target timespan to adjust better to hashpower (30000/300=100) post hardforkTestnet
	AveragingTargetTimespan: TestnetAveragingTargetTimespan,
	MaxAdjustDown:           TestnetMaxAdjustDown,
	MaxAdjustUp:             TestnetMaxAdjustUp,
	TargetTimespanAdjDown:   TestnetAveragingTargetTimespan * (TestnetInterval + TestnetMaxAdjustDown) / TestnetInterval,
	MinActualTimespan:       TestnetAveragingTargetTimespan * (TestnetInterval - TestnetMaxAdjustUp) / TestnetInterval,
	MaxActualTimespan:       TestnetAveragingTargetTimespan * (TestnetInterval + TestnetMaxAdjustDown) / TestnetInterval,
	ScryptPowLimit:          &scryptPowLimit,
	ScryptPowLimitBits:      ScryptPowLimitBits,
	RPCClientPort:           "21048",
	WalletRPCServerPort:     "21046",
}
