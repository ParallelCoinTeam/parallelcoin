// Bitcoin fork genesis block generator, based on https://bitcointalk.org/index.php?topic=181981.0 hosted at https://pastebin.com/nhuuV7y9
package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	log "github.com/p9c/pod/pkg/logi"
	"os"
	"runtime"
	"strconv"
	"time"

	chainhash "github.com/p9c/pod/pkg/chain/hash"
)

type transaction struct {
	merkleHash     []byte // 32 bytes long
	serializedData []byte
	version        uint32
	numInputs      byte
	prevOutput     []byte // 32 bytes long
	prevoutIndex   uint32
	scriptSig      []byte
	sequence       uint32
	numOutputs     byte
	outValue       uint64
	pubkeyScript   []byte
	locktime       uint32
}

const coin uint64 = 10000000

var (
	opChecksig byte = 172
	startNonce uint32
	unixtime   uint32
	maxNonce   = ^uint32(0)
)

// This function reverses the bytes in a byte array
func byteswap(buf []byte) {
	length := len(buf)
	for i := 0; i < length/2; i++ {
		buf[i], buf[length-i-1] = buf[length-i-1], buf[i]
	}
}
func initTransaction() (t transaction) {
	t.version = 1
	t.numInputs = 1
	t.numOutputs = 1
	t.locktime = 0
	t.prevoutIndex = 0xffffffff
	t.sequence = 0xfffffff
	t.outValue = coin
	t.prevOutput = make([]byte, 32)
	return
}
func main() {
	args := os.Args
	if len(args) != 4 {
		log.Println("Bitcoin fork genesis block generator")
		log.Println("Usage:")
		log.Println("    ", args[0], "<pubkey> <timestamp> <nBits>")
		log.Println("Example:")
		log.Println("    ", args[0], 
			"04678afdb0fe5548271967f1a67130b7105cd6a828e03909a67962e0ea1f61deb649f6bc3f4cef38c4f35504e51ec112de5c384df7ba0b8d578a4c702b6bf11d5f \"The Times 03/Jan/2009 Chancellor on brink of second bailout for banks\" 486604799")
		log.Println("\nIf you execute this without parameters another one in" +
			" the source code will be generated, using a random public key")
		args = []string{
			os.Args[0],
			"04678afdb0fe5548271967f1a67130b7105cd6a828e03909a67962e0ea1f61deb649f6bc3f4cef38c4f35504e51ec112de5c384df7ba0b8d578a4c702b6bf11d5f",
			"“All rational action is in the first place individual action. Only the individual thinks. Only the individual reasons. Only the individual acts.” - Ludwig von Mises, Socialism: An Economic and Sociological Analysis",
			"486604799",
		}
	}
	var pubkey []byte
	if args[1] == "" {
		pubkey = make([]byte, 65)
		n, err := rand.Read(pubkey)
		if err != nil {
			log.L.Error(err)
			os.Exit(1)
		}
		if n != 65 {
			log.L.Error("For some reason did not get 65 random bytes")
			os.Exit(1)
		}
		log.Printf("\nGenerated random public key:\n0x%x\n", pubkey)
	} else {
		if len(args[1]) != 130 {
			log.L.Error("Invalid public key length. Should be 130 hex digits,")
			os.Exit(1)
		}
		var err error
		pubkey, err = hex.DecodeString(args[1])
		if err != nil {
			log.L.Error(err)
			log.Println("Public key had invalid characters")
		}
	}
	timestamp := args[2]
	if len(timestamp) > 254 || len(timestamp) < 1 {
		log.L.Error("Timestamp was either longer than 254 characters or zero" +
			" length")
		os.Exit(1)
	}
	tx := initTransaction()
	nbits, err := strconv.ParseInt(args[3], 10, 32)
	if err != nil {
		log.L.Error("nBits was not a decimal number or exceeded the precision of 32 bits")
		os.Exit(0)
	}
	nBits := uint32(nbits)
	tx.pubkeyScript = joinBytes([]byte{0x41}, pubkey, []byte{opChecksig})
	switch {
	case nBits <= 255:
		tx.scriptSig = append([]byte{1}, byte(nBits))
	case nBits <= 65535:
		tx.scriptSig = joinBytes([]byte{2, byte(nBits)}, []byte{byte(nBits >> 8)})
	case nBits <= 16777215:
		tx.scriptSig = append([]byte{3}, byte(nBits))
		for i := uint(1); i < 3; i++ {
			tx.scriptSig = append(tx.scriptSig, byte(nBits>>(8*i)))
		}
	default:
		tx.scriptSig = append([]byte{4}, byte(nBits))
		for i := uint(1); i < 4; i++ {
			tx.scriptSig = append(tx.scriptSig, byte(nBits>>(8*i)))
		}
	}
	tx.scriptSig = joinBytes([]byte{0x01, 0x04, byte(len(timestamp))}, []byte(timestamp))
	tx.serializedData = joinBytes(
		uint32tobytes(tx.version),
		[]byte{tx.numInputs},
		tx.prevOutput,
		uint32tobytes(tx.prevoutIndex),
		[]byte{byte(len(tx.scriptSig))},
		tx.scriptSig,
		uint32tobytes(tx.sequence),
		[]byte{tx.numOutputs},
		uint64tobytes(tx.outValue),
		[]byte{byte(len(tx.pubkeyScript))},
		tx.pubkeyScript,
		uint32tobytes(tx.locktime))
	// hash1 := sha256.Sum256(tx.serializedData)
	// hash2 := sha256.Sum256(hash1[:])
	hash1 := chainhash.HashB(tx.serializedData)
	hash2 := sha256.Sum256(hash1[:])
	tx.merkleHash = hash2[:]
	merkleHash := hex.EncodeToString(tx.merkleHash)
	byteswap(tx.merkleHash)
	merkleHashSwapped := hex.EncodeToString(tx.merkleHash)
	byteswap(tx.merkleHash)
	txScriptSig := hex.EncodeToString(tx.scriptSig)
	pubScriptSig := hex.EncodeToString(tx.pubkeyScript)
	log.Printf("\nCoinbase:\n0x%s\n\nPubKeyScript:\n0x%s\n\nMerkle Hash:\n0x%s\n\nByteswapped:\n0x%s\n", txScriptSig, pubScriptSig, merkleHash, merkleHashSwapped)
	unixtime := uint32(time.Now().Unix())
	var blockversion uint32 = 4
	blockHeader := joinBytes(uint32tobytes(blockversion), make([]byte, 32), tx.merkleHash,
		uint32tobytes(uint32(unixtime)), // byte 68 - 71
		uint32tobytes(uint32(nBits)),
		uint32tobytes(startNonce)) // byte 76 - 79
	bytes := nBits >> 24
	// bytes := 31
	body := nBits << 8 >> 8
	var bits uint32
	for bits < 24 {
		bits++
		if body<<bits == 0 {
			break
		}
	}
	bits = 32 - bits
	if bits < 31 {
		bytes = bytes - bits/8
		bits = bits % 8
	}
	log.L.Info("\nSearching for nonce/unixtime combination that satisfies " +
		"minimum target %d with %d threads on %d cores...\nPlease wait... ",
		nBits, runtime.GOMAXPROCS(-1), runtime.NumCPU())
	start := time.Now()
	for i := 0; i < runtime.NumCPU(); i++ {
		go findNonce(blockHeader, bytes, bits, start)
		time.Sleep(time.Second)
	}
	time.Sleep(time.Hour)
}
func findNonce(b []byte, bytes, bits uint32, start time.Time) []byte {
	blockHeader := append([]byte(nil), b...)
	unixtime = uint32(time.Now().Unix())
	blockHeader[68] = byte(unixtime)
	blockHeader[69] = byte(unixtime >> 8)
	blockHeader[70] = byte(unixtime >> 16)
	blockHeader[71] = byte(unixtime >> 24)
	for {
		blockhash1 := sha256.Sum256(blockHeader)
		blockhash2 := sha256.Sum256(blockhash1[:])
		if undertarget(blockhash2[bytes:], bits) {
			byteswap(blockhash2[:])
			log.Printf("Block found!\n\nHash:\n0x%x\n\nNonce:\n%d\n\nUnix time:\n%d\n", blockhash2, startNonce, unixtime)
			log.Printf("\nBlock header encoded in hex:\n0x%x\n", blockHeader)
			log.Println("\nTime for nonce search:", time.Since(start))
			os.Exit(1)
		}
		startNonce++
		if startNonce < maxNonce {
			blockHeader[76] = byte(startNonce)
			blockHeader[77] = byte(startNonce >> 8)
			blockHeader[78] = byte(startNonce >> 16)
			blockHeader[79] = byte(startNonce >> 24)
		} else {
			startNonce = 0
			unixtime = uint32(time.Now().Unix())
			blockHeader[68] = byte(unixtime)
			blockHeader[69] = byte(unixtime >> 8)
			blockHeader[70] = byte(unixtime >> 16)
			blockHeader[71] = byte(unixtime >> 24)
		}
	}
}
func joinBytes(segment ...[]byte) (joined []byte) {
	joined = make([]byte, 0)
	for i := range segment {
		joined = append(joined, segment[i]...)
	}
	return
}
func undertarget(hash []byte, bits uint32) bool {
	// for i:=len(hash)-1; i>0; i-- { hash[i]=0 }
	// log.Println(hash)
	for i := len(hash) - 1; i > 0; i-- {
		// log.Println(hash[i])
		if hash[i] != 0 {
			return false
		}
	}
	// hash[0] = 0
	for i := bits; i > 0; i-- {
		if hash[0]<<i != 0 {
			return false
		}
	}
	return true
}
func uint32tobytes(u uint32) []byte {
	b := make([]byte, 4)
	b[0] = byte(u)
	for i := uint(1); i < 4; i++ {
		b[i] = byte(u >> (i * 8))
	}
	return b
}

// func bytestouint32(// 	b []byte) uint32 {
// 	if len(b) > 4 {
// 		return 0
// 	}
// 	var u uint32
// 	for i := uint32(3); i > 0; i-- {
// 		u += uint32(b[i]) << (i * 8)
// 	}
// 	return u
// }
func uint64tobytes(u uint64) []byte {
	b := make([]byte, 8)
	b[0] = byte(u)
	for i := uint(1); i < 8; i++ {
		b[i] = byte(u >> (i * 8))
	}
	return b
}
