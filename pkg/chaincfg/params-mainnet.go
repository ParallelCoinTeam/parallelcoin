package chaincfg

import (
	"github.com/p9c/matrjoska/pkg/wire"
)

// MainNetParams defines the network parameters for the main Bitcoin network.
var MainNetParams = Params{
	Name:        "mainnet",
	Net:         wire.MainNet,
	DefaultPort: "11047",
	DNSSeeds: []DNSSeed{
		{"seed1.parallelcoin.io", true},
		{"seed2.parallelcoin.io", true},
		{"seed3.parallelcoin.io", true},
		{"seed4.parallelcoin.io", true},
		{"seed5.parallelcoin.io", true},
		{"185.69.55.35", true},
		{"46.28.107.182", true},
		{"91.206.16.214", true},
		{"157.161.128.62", true},
		{"85.15.179.171", true},
		{"103.254.148.9", true},
		{"144.217.73.92", true},
		{"165.227.110.22", true},
		{"194.135.88.119", true},
		{"73.164.170.207", true},
		{"76.176.77.120", true},
		{"89.40.12.55", true},
		{"coins.prohashing.com:6245", true},
		//
	},
	// Chain parameters
	GenesisBlock: &genesisBlock,
	GenesisHash:  &genesisHash,
	PowLimit:     &mainPowLimit,
	PowLimitBits: MainPowLimitBits, // 0x1e0fffff,
	// BIP0034Height:            math.MaxInt32,        // Reserved for future change
	// BIP0065Height:            math.MaxInt32,
	// BIP0066Height:            math.MaxInt32,
	CoinbaseMaturity:         100,
	SubsidyReductionInterval: 250000,
	TargetTimespan:           TargetTimespan,
	TargetTimePerBlock:       TargetTimePerBlock,
	RetargetAdjustmentFactor: 2, // 50% less, 200% more (not used in parallelcoin)
	ReduceMinDifficulty:      false,
	MinDiffReductionTime:     0,
	GenerateSupported:        true,
	// Checkpoints ordered from oldest to newest.
	Checkpoints: []Checkpoint{
		// {, newHashFromStr("")},
		// {200069, newHashFromStr("000000000000044e641986c8ee672460e853a11b352869cb8a4a8ba0b3f3e6dc")},
	},
	// Consensus rule change deployments.
	//
	// The miner confirmation window is defined as:
	//   target proof of work timespan / target proof of work spacing
	RuleChangeActivationThreshold: 1916, // 95% of MinerConfirmationWindow
	MinerConfirmationWindow:       2016, //
	// Deployments: [DefinedDeployments]ConsensusDeployment{
	// 	DeploymentTestDummy: {
	// 		BitNumber:  28,
	// 		StartTime:  math.MaxUint64, // January 1, 2008 UTC
	// 		ExpireTime: math.MaxUint64, // December 31, 2008 UTC
	// 	},
	// 	DeploymentCSV: {
	// 		BitNumber:  0,
	// 		StartTime:  math.MaxUint64, // May 1st, 2016
	// 		ExpireTime: math.MaxUint64, // May 1st, 2017
	// 	},
	// 	DeploymentSegwit: {
	// 		BitNumber:  1,
	// 		StartTime:  math.MaxUint64, // November 15, 2016 UTC
	// 		ExpireTime: math.MaxUint64, // November 15, 2017 UTC.
	// 	},
	// },
	// Mempool parameters
	RelayNonStdTxs: false,
	// // Human-readable part for Bech32 encoded segwit addresses, as defined in
	// // BIP 173.
	// Bech32HRPSegwit: "p9", // always bc for main net
	// Address encoding magics
	PubKeyHashAddrID: 83,  // 0x00, // starts with 1
	ScriptHashAddrID: 9,   // 0x05, // starts with 3
	PrivateKeyID:     178, // 0x80, // starts with 5 (uncompressed) or K (compressed)
	// WitnessPubKeyHashAddrID: 84,  // 0x06, // starts with p2
	// WitnessScriptHashAddrID: 19,  // 0x0A, // starts with 7Xh
	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x04, 0x88, 0xad, 0xe4}, // starts with xprv
	HDPublicKeyID:  [4]byte{0x04, 0x88, 0xb2, 0x1e}, // starts with xpub
	// BIP44 coin type used in the hierarchical deterministic path for address generation.
	HDCoinType: 0,
	// Parallelcoin specific difficulty adjustment parameters
	Interval:                Interval,
	AveragingInterval:       AveragingInterval, // Extend to target timespan to adjust better to hashpower (30000/300=100) post hardfork
	AveragingTargetTimespan: AveragingTargetTimespan,
	MaxAdjustDown:           MaxAdjustDown,
	MaxAdjustUp:             MaxAdjustUp,
	TargetTimespanAdjDown: AveragingTargetTimespan *
		(Interval + MaxAdjustDown) / Interval,
	MinActualTimespan:   2400,
	MaxActualTimespan:   3300,
	ScryptPowLimit:      &scryptPowLimit,
	ScryptPowLimitBits:  ScryptPowLimitBits,
	RPCClientPort:       "11048",
	WalletRPCServerPort: "11046",
	
}
