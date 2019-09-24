package fork

import (
	"crypto/sha256"
	"math/big"

	"git.parallelcoin.io/dev/cryptonight"
	"github.com/bitgoin/lyra2rev2"
	"github.com/dchest/blake256"
	skein "github.com/enceve/crypto/skein/skein256"
	gost "github.com/programmer10110/gostreebog"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/blake2s"
	"golang.org/x/crypto/scrypt"
	"golang.org/x/crypto/sha3"

	chainhash "github.com/p9c/pod/pkg/chain/hash"
)

const HashReps = 1

// Argon2i takes bytes, generates a Lyra2REv2 hash as salt, generates an argon2i key
func Argon2i(bytes []byte) []byte {
	return argon2.IDKey(bytes, bytes, 1, 64*1024, 1, 32)
}

// Blake14lr takes bytes and returns a blake14lr 256 bit hash
func Blake14lr(bytes []byte) []byte {
	a := blake256.New()
	_, _ = a.Write(bytes)
	return a.Sum(nil)
}

// Blake2b takes bytes and returns a blake2b 256 bit hash
func Blake2b(bytes []byte) []byte {
	b := blake2b.Sum256(bytes)
	return b[:]
}

// Blake2s takes bytes and returns a blake2s 256 bit hash
func Blake2s(bytes []byte) []byte {
	b := blake2s.Sum256(bytes)
	return b[:]
}

// Cryptonight7v2 takes bytes and returns a cryptonight 7 v2 256 bit hash
func Cryptonight7v2(bytes []byte) []byte {
	return cryptonight.Sum(bytes, 2)
}

func reverse(b []byte) []byte {
	out := make([]byte, len(b))
	for i := range b {
		out[i] = b[len(b)-1]
	}
	return out
}

// DivHash first runs an arbitrary big number calculation involving a very
// large integer, and hashes the result. In this way, this hash requires
// both one of 9 arbitrary hash functions plus a big number long division
// operation and three multiplication operations, unlikely to be satisfied
// on anything other than CPU and GPU, with contrary advantages on each -
// GPU division is 32 bits wide operations, CPU is 64, but GPU hashes about
// equal to a CPU to varying degrees of memory hardness (and CPU cache size
// then improves CPU performance at some hashes)
func DivHash(hf func([]byte) []byte, blockbytes []byte, howmany int) []byte {
	blocklen := len(blockbytes)
	firsthalf := append(blockbytes, reverse(blockbytes[:blocklen/2])...)
	fhc := make([]byte, len(firsthalf))
	copy(fhc, firsthalf)
	secondhalf := append(blockbytes, reverse(blockbytes[blocklen/2:])...)
	shc := make([]byte, len(secondhalf))
	copy(shc, secondhalf)
	bbb := make([]byte, len(blockbytes))
	copy(bbb, blockbytes)
	bl := big.NewInt(0).SetBytes(bbb)
	fh := big.NewInt(0).SetBytes(fhc)
	sh := big.NewInt(0).SetBytes(shc)
	// sqfh := fh.Mul(fh, fh)
	// sqsh := sh.Mul(sh, sh)
	// sqsq := fh.Mul(sqfh, sqsh)
	sqsq := fh.Mul(fh, sh)
	divd := sqsq.Div(sqsq, bl)
	divdb := divd.Bytes()
	// pad out to 512 bit wide chunks to handle all hash functions
	// note that this alters the cut boundary in the mangling phase
	// above when hash is repeated
	dlen := len(divdb)
	if dlen%64 != 0 {
		dlen = (dlen/64 + 1) * 64
	}
	ddd := make([]byte, dlen)
	copy(ddd, divdb)
	// this allows us run this operation an arbitrary number of times
	if howmany > 0 {
		return DivHash(hf, ddd, howmany-1)
	}
	// fmt.Printf("%x\n", ddd)
	return hf(ddd)
}

// Hash computes the hash of bytes using the named hash
func Hash(bytes []byte, name string, height int32) (out chainhash.Hash) {
	// if IsTestnet && height < 10 {
	// 	// log <- cl.Warn{"hash", name, height}
	// 	time.Sleep(time.Second / 20)
	// }
	// INFO("hash", name, height}
	switch name {
	case "blake2b":
		_ = out.SetBytes(DivHash(Blake2b, bytes, HashReps))
	case "argon2i":
		_ = out.SetBytes(DivHash(Argon2i, bytes, HashReps))
	case "cryptonight7v2":
		_ = out.SetBytes(DivHash(Cryptonight7v2, bytes, HashReps))
	case "lyra2rev2":
		_ = out.SetBytes(DivHash(Lyra2REv2, bytes, HashReps))
	case "scrypt":
		if GetCurrent(height) > 0 {
			_ = out.SetBytes(DivHash(Scrypt, bytes, HashReps))
		} else {
			_ = out.SetBytes(Scrypt(bytes))
		}
	case "sha256d": // sha256d
		if GetCurrent(height) > 0 {
			_ = out.SetBytes(DivHash(chainhash.DoubleHashB, bytes, HashReps))
		} else {
			_ = out.SetBytes(chainhash.DoubleHashB(bytes))
		}
	case "stribog":
		_ = out.SetBytes(DivHash(Stribog, bytes, HashReps))
	case "skein":
		_ = out.SetBytes(DivHash(Skein, bytes, HashReps))
	case "keccak":
		_ = out.SetBytes(DivHash(Keccak, bytes, HashReps))
	}
	return
}

// Keccak takes bytes and returns a keccak (sha-3) 256 bit hash
func Keccak(bytes []byte) []byte {
	sum := sha3.Sum256(bytes)
	return sum[:]
}

// Lyra2REv2 takes bytes and returns a lyra2rev2 256 bit hash
func Lyra2REv2(bytes []byte) []byte {
	bytes, _ = lyra2rev2.Sum(bytes)
	return bytes
}

// SHA256D takes bytes and returns a double SHA256 hash
func SHA256D(bytes []byte) []byte {
	h := sha256.Sum256(bytes)
	h = sha256.Sum256(h[:])
	return h[:]
}

// Scrypt takes bytes and returns a scrypt 256 bit hash
func Scrypt(bytes []byte) []byte {
	b := bytes
	c := make([]byte, len(b))
	copy(c, b)
	dk, err := scrypt.Key(c, c, 1024, 1, 1, 32)
	if err != nil {
		return make([]byte, 32)
	}
	o := make([]byte, 32)
	for i := range dk {
		o[i] = dk[len(dk)-1-i]
	}
	copy(o, dk)
	return o
}

// Skein takes bytes and returns a skein 256 bit hash
func Skein(bytes []byte) []byte {
	var out [32]byte
	skein.Sum256(&out, bytes, nil)
	return out[:]
}

// Stribog takes bytes and returns a double GOST Stribog 256 bit hash
func Stribog(bytes []byte) []byte {
	return gost.Hash(bytes, "256")
}
