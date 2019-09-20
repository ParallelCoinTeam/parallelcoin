package chaincfg

import "git.parallelcoin.io/dev/pod/pkg/chain/wire"

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
	},
	// Chain parameters
	GenesisBlock:             &genesisBlock,
	GenesisHash:              &genesisHash,
	PowLimit:                 &mainPowLimit,
	PowLimitBits:             MainPowLimitBits, //0x1e0fffff,
	BIP0034Height:            100000000,        // Reserved for future change
	BIP0065Height:            100000000,
	BIP0066Height:            100000000,
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
	Deployments: [DefinedDeployments]ConsensusDeployment{
		DeploymentTestDummy: {
			BitNumber:  28,
			StartTime:  1199145601, // January 1, 2008 UTC
			ExpireTime: 1230767999, // December 31, 2008 UTC
		},
		DeploymentCSV: {
			BitNumber:  0,
			StartTime:  1462060800, // May 1st, 2016
			ExpireTime: 1493596800, // May 1st, 2017
		},
		DeploymentSegwit: {
			BitNumber:  1,
			StartTime:  1479168000, // November 15, 2016 UTC
			ExpireTime: 1510704000, // November 15, 2017 UTC.
		},
	},
	// Mempool parameters
	RelayNonStdTxs: false,
	// Human-readable part for Bech32 encoded segwit addresses, as defined in
	// BIP 173.
	Bech32HRPSegwit: "pc", // always bc for main net
	// Address encoding magics
	PubKeyHashAddrID:        83,  // 0x00, // starts with 1
	ScriptHashAddrID:        9,   // 0x05, // starts with 3
	PrivateKeyID:            178, // 0x80, // starts with 5 (uncompressed) or K (compressed)
	WitnessPubKeyHashAddrID: 84,  // 0x06, // starts with p2
	WitnessScriptHashAddrID: 19,  // 0x0A, // starts with 7Xh
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
	MinActualTimespan:  2400,
	MaxActualTimespan:  3300,
	ScryptPowLimit:     &scryptPowLimit,
	ScryptPowLimitBits: ScryptPowLimitBits,
}
