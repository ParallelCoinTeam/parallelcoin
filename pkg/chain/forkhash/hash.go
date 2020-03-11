package forkhash

import (
	"math/big"
	
	"github.com/bitbandi/go-x11"
	
	"github.com/p9c/pod/pkg/chain/fork"
	log "github.com/p9c/logi"
	
	skein "github.com/enceve/crypto/skein/skein256"
	gost "github.com/programmer10110/gostreebog"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/scrypt"
	"golang.org/x/crypto/sha3"
	"lukechampine.com/blake3"
	
	chainhash "github.com/p9c/pod/pkg/chain/hash"
)

// HashReps allows the number of multiplication/division cycles to be
// repeated before the final hash,
// on release for mainnet this is probably set to 9 or so to raise the
// difficulty to a reasonable level for the hard fork.
// at 5 repetitions (first plus repeats, thus 4),
// an example block header produces a number around 48kb in byte size and
// ~119000 decimal digits, which is then finally hashed down to 32 bytes
var HashReps = 2

// Argon2i takes bytes, generates a Blake3 hash as salt, generates an argon2i key
func Argon2i(bytes []byte) []byte {
	return argon2.IDKey(reverse(bytes), bytes, 1, 4*1024, 1, 32)
}

// Blake2b takes bytes and returns a blake2b 256 bit hash
func Blake2b(bytes []byte) []byte {
	b := blake2b.Sum256(bytes)
	return b[:]
}

// X11 takes bytes and returns a X11 256 bit hash
func X11(bytes []byte) (out []byte) {
	hf := x11.New()
	out = make([]byte, 32)
	hf.Hash(bytes, out)
	// log.L.Debugf("x11 %x", out)
	return
	// return cryptonight.Sum(bytes, 2)
}

func reverse(b []byte) []byte {
	out := make([]byte, len(b))
	for i := range b {
		out[i] = b[len(b)-1-i]
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
//
// This hash generates very large random numbers from an 80 byte block using
// a procedure involving two squares of recombined spliced halves,
// multiplied together and then divided by
// the original block, reversed, repeated 4 more times,
// of over 48kb to represent the product, which is hashed afterwards.
// ( see example in fork/scratch/divhash.go )
//
// This would be around half of the available level 1 cache of a ryzen 5 1600,
// likely distributed as 6 using the 64kb and 3 threads using the 32kb
// smaller ones. Here is a block diagram of Zen architecture:
// https://en.wikichip.org/w/images/thumb/0/02/zen_block_diagram.svg/1178px-zen_block_diagram.svg.png
// This is one, with Zen the cores are independently partitioned.
// But it shows that it has one divider per core. Thus for a 6 core,
// 6 would be the right number to use with it as other numbers will lead to
// contention and memory copies.
// Probably its ability to branch twice per cycle will be a big boost for its
// performance in this task.
//
// Long division units are expensive and slow,
// and make a perfect application specific proof of work because a
// substantial part of the cost as proportional to the relative surface area
// of circuitry it is substantially more than 10% of the total logic on a
// CPU. There is low chances of putting these units into one package with
// half half IDIV and IMUL units and enough cache for each one,
// would be economic or accessible to most chip manufacturers at a scale and
// clock that beats the CPU price.
//
// Most GPUs still only have 32 bit integer divide units because the type of
// mathematics done by GPUs is mainly based on multiplication,
// addition and subtraction, specifically, with matrixes,
// which are per-bit equivalent to big (128-512 bit) addition,
// built for walking graphs and generating directional or particle effects,
// under gravity, and the like.
// Video is very parallelisable so generally speaking GPU's main bulk of
// processing capability does not help here,
// caches holding a fraction of the number at a time as it is computed,
// and only 32 bits wide at a time for the special purpose dividers,
// that are relatively swamped by stream processors in big grids.
//
// The cheaper arithmetic units can be programmed to also perform the
// calculations but they are going to be funny letter log differences to the
// point it adds up to less than 10% better due to complexity of the code and
// scheduling it.
//
// Long story short,
// this hash function should be the end of big margin ASIC mining,
// and a lot of R&D funds going to improving smaller fabs for spewing out
// such processors.
func DivHash(hf func([]byte) []byte, blockbytes []byte, howmany int) []byte {
	blocklen := len(blockbytes)
	rfb := reverse(blockbytes[:blocklen/2])
	firsthalf := append(blockbytes, rfb...)
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
	sqfh := fh.Mul(fh, fh)
	sqsh := sh.Mul(sh, sh)
	sqsq := fh.Mul(sqfh, sqsh)
	divd := sqsq.Div(sqsq, bl)
	divdb := divd.Bytes()
	dlen := len(divdb)
	ddd := make([]byte, dlen)
	copy(ddd, reverse(divdb))
	// this allows us run this operation an arbitrary number of times
	// log.Printf("%d bytes %x\n", len(ddd), ddd)
	if howmany > 0 {
		return DivHash(hf, append(ddd, reverse(ddd)...), howmany-1)
	}
	// return X11(hf(ddd))
	return hf(ddd)
}

// Hash computes the hash of bytes using the named hash
func Hash(bytes []byte, name string, height int32) (out chainhash.Hash) {
	hR := HashReps
	if fork.IsTestnet {
		switch {
		case height == 1:
			hR = 0
			// case height < 10:
			//	hR = 6
		}
	}
	switch name {
	case fork.Scrypt:
		if fork.GetCurrent(height) > 0 {
			_ = out.SetBytes(DivHash(Scrypt, bytes, hR))
		} else {
			_ = out.SetBytes(Scrypt(bytes))
		}
	case fork.SHA256d:
		if fork.GetCurrent(height) > 0 {
			_ = out.SetBytes(DivHash(chainhash.DoubleHashB, bytes, hR))
		} else {
			_ = out.SetBytes(chainhash.DoubleHashB(
				bytes))
		}
	default:
		_ = out.SetBytes(DivHash(Blake3, bytes, hR))
	}
	return
}

// Keccak takes bytes and returns a keccak (sha-3) 256 bit hash
func Keccak(bytes []byte) []byte {
	sum := sha3.Sum256(bytes)
	return sum[:]
}

// Blake3 takes bytes and returns a lyra2rev2 256 bit hash
func Blake3(bytes []byte) []byte {
	b := blake3.Sum256(bytes)
	return b[:]
}

// Scrypt takes bytes and returns a scrypt 256 bit hash
func Scrypt(bytes []byte) []byte {
	b := bytes
	c := make([]byte, len(b))
	copy(c, b)
	dk, err := scrypt.Key(c, c, 1024, 1, 1, 32)
	if err != nil {
		log.L.Error(err)
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
