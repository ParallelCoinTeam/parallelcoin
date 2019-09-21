package util

const (
	// SatoshiPerBitcent is the number of satoshi in one bitcoin cent.
	SatoshiPerBitcent Amount = 1e6
	// SatoshiPerBitcoin is the number of satoshi in one bitcoin (1 DUO).
	SatoshiPerBitcoin Amount = 1e8
	// MaxSatoshi is the maximum transaction amount allowed in satoshi.
	MaxSatoshi = 21e6 * SatoshiPerBitcoin
)
