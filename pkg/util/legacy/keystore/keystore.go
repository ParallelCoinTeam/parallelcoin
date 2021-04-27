package keystore

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/p9c/matrjoska/pkg/btcaddr"
	"github.com/p9c/matrjoska/pkg/chaincfg"
	"io"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"sync"
	"time"
	
	"golang.org/x/crypto/ripemd160"
	
	"github.com/p9c/matrjoska/pkg/chainhash"
	ec "github.com/p9c/matrjoska/pkg/ecc"
	"github.com/p9c/matrjoska/pkg/txscript"
	"github.com/p9c/matrjoska/pkg/util"
	"github.com/p9c/matrjoska/pkg/util/legacy/rename"
	"github.com/p9c/matrjoska/pkg/wire"
)

// A bunch of constants
const (
	Filename = "wallet.bin"
	// Length in bytes of KDF output.
	kdfOutputBytes = 32
	// Maximum length in bytes of a comment that can have a size represented as a uint16.
	maxCommentLen         = (1 << 16) - 1
	defaultKdfComputeTime = 0.25
	defaultKdfMaxMem      = 32 * 1024 * 1024
)

// Possible errors when dealing with key stores.
var (
	ErrAddressNotFound  = errors.New("address not found")
	ErrAlreadyEncrypted = errors.New("private key is already encrypted")
	ErrChecksumMismatch = errors.New("checksum mismatch")
	ErrDuplicate        = errors.New("duplicate key or address")
	ErrMalformedEntry   = errors.New("malformed entry")
	ErrWatchingOnly     = errors.New("keystore is watching-only")
	ErrLocked           = errors.New("keystore is locked")
	ErrWrongPassphrase  = errors.New("wrong passphrase")
)
var fileID = [8]byte{0xba, 'W', 'A', 'L', 'L', 'E', 'T', 0x00}

type entryHeader byte

const (
	addrCommentHeader entryHeader = 1 << iota
	txCommentHeader
	deletedHeader
	scriptHeader
	addrHeader entryHeader = 0
)

// We want to use binaryRead and binaryWrite instead of binary.Read and binary.Write because those from the binary
// package do not return the number of bytes actually written or read. We need to return this value to correctly support
// the io.ReaderFrom and io.WriterTo interfaces.
func binaryRead(r io.Reader, order binary.ByteOrder, data interface{}) (n int64, e error) {
	var read int
	buf := make([]byte, binary.Size(data))
	if read, e = io.ReadFull(r, buf); E.Chk(e) {
		return int64(read), e
	}
	return int64(read), binary.Read(bytes.NewBuffer(buf), order, data)
}

// See comment for binaryRead().
func binaryWrite(w io.Writer, order binary.ByteOrder, data interface{}) (n int64, e error) {
	buf := bytes.Buffer{}
	if e = binary.Write(&buf, order, data); E.Chk(e) {
		return 0, e
	}
	written, e := w.Write(buf.Bytes())
	return int64(written), e
}

// pubkeyFromPrivkey creates an encoded pubkey based on a 32-byte privkey. The returned pubkey is 33 bytes if
// compressed, or 65 bytes if uncompressed.
func pubkeyFromPrivkey(privkey []byte, compress bool) (pubkey []byte) {
	_, pk := ec.PrivKeyFromBytes(ec.S256(), privkey)
	if compress {
		return pk.SerializeCompressed()
	}
	return pk.SerializeUncompressed()
}
func keyOneIter(passphrase, salt []byte, memReqts uint64) []byte {
	saltedpass := append(passphrase, salt...)
	lutbl := make([]byte, memReqts)
	// Seed for lookup table
	seed := sha512.Sum512(saltedpass)
	copy(lutbl[:sha512.Size], seed[:])
	for nByte := 0; nByte < (int(memReqts) - sha512.Size); nByte += sha512.Size {
		hash := sha512.Sum512(lutbl[nByte : nByte+sha512.Size])
		copy(lutbl[nByte+sha512.Size:nByte+2*sha512.Size], hash[:])
	}
	x := lutbl[cap(lutbl)-sha512.Size:]
	seqCt := uint32(memReqts / sha512.Size)
	nLookups := seqCt / 2
	for i := uint32(0); i < nLookups; i++ {
		// Armory ignores endianness here.  We assume LE.
		newIdx := binary.LittleEndian.Uint32(x[cap(x)-4:]) % seqCt
		// Index of hash result at newIdx
		vIdx := newIdx * sha512.Size
		v := lutbl[vIdx : vIdx+sha512.Size]
		// XOR hash x with hash v
		for j := 0; j < sha512.Size; j++ {
			x[j] ^= v[j]
		}
		// Save new hash to x
		hash := sha512.Sum512(x)
		copy(x, hash[:])
	}
	return x[:kdfOutputBytes]
}

// kdf implements the key derivation function used by Armory based on the ROMix algorithm described in Colin Percival's
// paper "Stronger Key Derivation via Sequential Memory-Hard Functions" (http://www.tarsnap.com/scrypt/scrypt.pdf).
func kdf(passphrase []byte, params *kdfParameters) []byte {
	masterKey := passphrase
	for i := uint32(0); i < params.nIter; i++ {
		masterKey = keyOneIter(masterKey, params.salt[:], params.mem)
	}
	return masterKey
}
func pad(size int, b []byte) []byte {
	// Prevent a possible panic if the input exceeds the expected size.
	if len(b) > size {
		size = len(b)
	}
	p := make([]byte, size)
	copy(p[size-len(b):], b)
	return p
}

// chainedPrivKey deterministically generates a new private key using a previous address and chaincode. privkey and
// chaincode must be 32 bytes long, and pubkey may either be 33 or 65 bytes.
func chainedPrivKey(privkey, pubkey, chaincode []byte) ([]byte, error) {
	if len(privkey) != 32 {
		return nil, fmt.Errorf(
			"invalid privkey length %d (must be 32)",
			len(privkey),
		)
	}
	if len(chaincode) != 32 {
		return nil, fmt.Errorf(
			"invalid chaincode length %d (must be 32)",
			len(chaincode),
		)
	}
	switch n := len(pubkey); n {
	case ec.PubKeyBytesLenUncompressed, ec.PubKeyBytesLenCompressed:
		// Correct length
	default:
		return nil, fmt.Errorf("invalid pubkey length %d", n)
	}
	xorbytes := make([]byte, 32)
	chainMod := chainhash.DoubleHashB(pubkey)
	for i := range xorbytes {
		xorbytes[i] = chainMod[i] ^ chaincode[i]
	}
	chainXor := new(big.Int).SetBytes(xorbytes)
	privint := new(big.Int).SetBytes(privkey)
	t := new(big.Int).Mul(chainXor, privint)
	b := t.Mod(t, ec.S256().N).Bytes()
	return pad(32, b), nil
}

// chainedPubKey deterministically generates a new public key using a previous public key and chaincode. pubkey must be
// 33 or 65 bytes, and chaincode must be 32 bytes long.
func chainedPubKey(pubkey, chaincode []byte) ([]byte, error) {
	var compressed bool
	switch n := len(pubkey); n {
	case ec.PubKeyBytesLenUncompressed:
		compressed = false
	case ec.PubKeyBytesLenCompressed:
		compressed = true
	default:
		// Incorrect serialized pubkey length
		return nil, fmt.Errorf("invalid pubkey length %d", n)
	}
	if len(chaincode) != 32 {
		return nil, fmt.Errorf(
			"invalid chaincode length %d (must be 32)",
			len(chaincode),
		)
	}
	xorbytes := make([]byte, 32)
	chainMod := chainhash.DoubleHashB(pubkey)
	for i := range xorbytes {
		xorbytes[i] = chainMod[i] ^ chaincode[i]
	}
	oldPk, e := ec.ParsePubKey(pubkey, ec.S256())
	if e != nil {
		E.Ln(e)
		return nil, e
	}
	newX, newY := ec.S256().ScalarMult(oldPk.X, oldPk.Y, xorbytes)
	// if e != nil  {
	// 	return nil, e
	// }
	newPk := &ec.PublicKey{
		Curve: ec.S256(),
		X:     newX,
		Y:     newY,
	}
	if compressed {
		return newPk.SerializeCompressed(), nil
	}
	return newPk.SerializeUncompressed(), nil
}

type ksVersion struct {
	major         byte
	minor         byte
	bugfix        byte
	autoincrement byte
}

// Enforce that ksVersion satisfies the io.ReaderFrom and io.WriterTo interfaces.
var _ io.ReaderFrom = &ksVersion{}
var _ io.WriterTo = &ksVersion{}

// readerFromVersion is an io.ReaderFrom and io.WriterTo that can specify any particular key store file format for
// reading depending on the key store file ksVersion.
type readerFromVersion interface {
	readFromVersion(ksVersion, io.Reader) (int64, error)
	io.WriterTo
}

func (v ksVersion) String() string {
	str := fmt.Sprintf("%d.%d", v.major, v.minor)
	if v.bugfix != 0x00 || v.autoincrement != 0x00 {
		str += fmt.Sprintf(".%d", v.bugfix)
	}
	if v.autoincrement != 0x00 {
		str += fmt.Sprintf(".%d", v.autoincrement)
	}
	return str
}
func (v ksVersion) Uint32() uint32 {
	return uint32(v.major)<<6 | uint32(v.minor)<<4 | uint32(v.bugfix)<<2 | uint32(v.autoincrement)
}
func (v *ksVersion) ReadFrom(r io.Reader) (int64, error) {
	// Read 4 bytes for the ksVersion.
	var versBytes [4]byte
	n, e := io.ReadFull(r, versBytes[:])
	if e != nil {
		E.Ln(e)
		return int64(n), e
	}
	v.major = versBytes[0]
	v.minor = versBytes[1]
	v.bugfix = versBytes[2]
	v.autoincrement = versBytes[3]
	return int64(n), nil
}
func (v *ksVersion) WriteTo(w io.Writer) (int64, error) {
	// Write 4 bytes for the ksVersion.
	versBytes := []byte{
		v.major,
		v.minor,
		v.bugfix,
		v.autoincrement,
	}
	n, e := w.Write(versBytes)
	return int64(n), e
}

// LT returns whether v is an earlier ksVersion than v2.
func (v ksVersion) LT(v2 ksVersion) bool {
	switch {
	case v.major < v2.major:
		return true
	case v.minor < v2.minor:
		return true
	case v.bugfix < v2.bugfix:
		return true
	case v.autoincrement < v2.autoincrement:
		return true
	default:
		return false
	}
}

// EQ returns whether v2 is an equal ksVersion to v.
func (v ksVersion) EQ(v2 ksVersion) bool {
	switch {
	case v.major != v2.major:
		return false
	case v.minor != v2.minor:
		return false
	case v.bugfix != v2.bugfix:
		return false
	case v.autoincrement != v2.autoincrement:
		return false
	default:
		return true
	}
}

// GT returns whether v is a later ksVersion than v2.
func (v ksVersion) GT(v2 ksVersion) bool {
	switch {
	case v.major > v2.major:
		return true
	case v.minor > v2.minor:
		return true
	case v.bugfix > v2.bugfix:
		return true
	case v.autoincrement > v2.autoincrement:
		return true
	default:
		return false
	}
}

// Various versions.
var (
	// VersArmory is the latest ksVersion used by Armory.
	VersArmory = ksVersion{1, 35, 0, 0}
	// Vers20LastBlocks is the ksVersion where key store files now hold the 20 most recently seen block hashes.
	Vers20LastBlocks = ksVersion{1, 36, 0, 0}
	// VersUnsetNeedsPrivkeyFlag is the bugfix ksVersion where the createPrivKeyNextUnlock address flag is correctly unset
	// after creating and encrypting its private key after unlock.
	//
	// Otherwise, re-creating private keys will occur too early in the address chain and fail due to encrypting an
	// already encrypted address.
	//
	// Key store versions at or before this ksVersion include a special case to allow the duplicate encrypt.
	VersUnsetNeedsPrivkeyFlag = ksVersion{1, 36, 1, 0}
	// VersCurrent is the current key store file ksVersion.
	VersCurrent = VersUnsetNeedsPrivkeyFlag
)

type varEntries struct {
	store   *Store
	entries []io.WriterTo
}

func (v *varEntries) WriteTo(w io.Writer) (n int64, e error) {
	ss := v.entries
	var written int64
	for _, s := range ss {
		var e error
		if written, e = s.WriteTo(w); E.Chk(e) {
			return n + written, e
		}
		n += written
	}
	return n, nil
}
func (v *varEntries) ReadFrom(r io.Reader) (n int64, e error) {
	var read int64
	// Remove any previous entries.
	v.entries = nil
	wts := v.entries
	// Keep reading entries until an EOF is reached.
	for {
		var header entryHeader
		if read, e = binaryRead(r, binary.LittleEndian, &header); E.Chk(e) {
			// EOF here is not an error.
			if e == io.EOF {
				return n + read, nil
			}
			return n + read, e
		}
		n += read
		var wt io.WriterTo
		switch header {
		case addrHeader:
			var entry addrEntry
			entry.addr.store = v.store
			if read, e = entry.ReadFrom(r); E.Chk(e) {
				return n + read, e
			}
			n += read
			wt = &entry
		case scriptHeader:
			var entry scriptEntry
			entry.script.store = v.store
			if read, e = entry.ReadFrom(r); E.Chk(e) {
				return n + read, e
			}
			n += read
			wt = &entry
		default:
			return n, fmt.Errorf("unknown entry header: %d", uint8(header))
		}
		if wt != nil {
			wts = append(wts, wt)
			v.entries = wts
		}
	}
}

// Key stores use a custom network parameters type so it can be an io.ReaderFrom.
//
// Due to the way and order that key stores are currently serialized and how address reading requires the key store's
// network parameters, setting and failing on unknown key store networks must happen on the read itself and not after
// the fact.
//
// This is admittedly a hack, but with a bip32 keystore on the horizon I'm not too motivated to clean this up.
type netParams chaincfg.Params

func (net *netParams) ReadFrom(r io.Reader) (int64, error) {
	var buf [4]byte
	uint32Bytes := buf[:4]
	n, e := io.ReadFull(r, uint32Bytes)
	n64 := int64(n)
	if e != nil {
		E.Ln(e)
		return n64, e
	}
	switch wire.BitcoinNet(binary.LittleEndian.Uint32(uint32Bytes)) {
	case wire.MainNet:
		*net = netParams(chaincfg.MainNetParams)
	case wire.TestNet3:
		*net = netParams(chaincfg.TestNet3Params)
	case wire.SimNet:
		*net = netParams(chaincfg.SimNetParams)
	default:
		return n64, errors.New("unknown network")
	}
	return n64, nil
}
func (net *netParams) WriteTo(w io.Writer) (int64, error) {
	var buf [4]byte
	uint32Bytes := buf[:4]
	binary.LittleEndian.PutUint32(uint32Bytes, uint32(net.Net))
	n, e := w.Write(uint32Bytes)
	n64 := int64(n)
	return n64, e
}

// Stringified byte slices for use as map lookup keys.
type addressKey string

// type transactionHashKey string
// type comment []byte

func getAddressKey(addr btcaddr.Address) addressKey {
	return addressKey(addr.ScriptAddress())
}

// Store represents an key store in memory. It implements the io.ReaderFrom and io.WriterTo interfaces to read from and
// write to any type of byte streams, including files.
type Store struct {
	// TODO: Use atomic operations for dirty so the reader lock doesn't need to be grabbed.
	dirty        bool
	path         string
	dir          string
	file         string
	mtx          sync.RWMutex
	vers         ksVersion
	net          *netParams
	flags        walletFlags
	createDate   int64
	name         [32]byte
	desc         [256]byte
	highestUsed  int64
	kdfParams    kdfParameters
	keyGenerator btcAddress
	// These are non-standard and fit in the extra 1024 bytes between the root address and the appended entries.
	recent  recentBlocks
	addrMap map[addressKey]walletAddress
	// The rest of the fields in this struct are not serialized.
	passphrase       []byte
	secret           []byte
	chainIdxMap      map[int64]btcaddr.Address
	importedAddrs    []walletAddress
	lastChainIdx     int64
	missingKeysStart int64
}

// New creates and initializes a new Store. name's and desc's byte length must not exceed 32 and 256 bytes,
// respectively. All address private keys are encrypted with passphrase. The key store is returned locked.
func New(
	dir string, desc string, passphrase []byte, net *chaincfg.Params,
	createdAt *BlockStamp,
) (s *Store, e error) {
	// Chk txsizes of inputs.
	if len(desc) > 256 {
		return nil, errors.New("desc exceeds 256 byte maximum size")
	}
	// Randomly-generate rootkey and chaincode.
	rootkey := make([]byte, 32)
	if _, e = rand.Read(rootkey); E.Chk(e) {
		return nil, e
	}
	chaincode := make([]byte, 32)
	if _, e = rand.Read(chaincode); E.Chk(e) {
		return nil, e
	}
	// Compute AES key and encrypt root address.
	kdfp, e := computeKdfParameters(defaultKdfComputeTime, defaultKdfMaxMem)
	if e != nil {
		E.Ln(e)
		return nil, e
	}
	aeskey := kdf(passphrase, kdfp)
	// Create and fill key store.
	s = &Store{
		path: filepath.Join(dir, Filename),
		dir:  dir,
		file: Filename,
		vers: VersCurrent,
		net:  (*netParams)(net),
		flags: walletFlags{
			useEncryption: true,
			watchingOnly:  false,
		},
		createDate:  time.Now().Unix(),
		highestUsed: rootKeyChainIdx,
		kdfParams:   *kdfp,
		recent: recentBlocks{
			lastHeight: createdAt.Height,
			hashes: []*chainhash.Hash{
				createdAt.Hash,
			},
		},
		addrMap:          make(map[addressKey]walletAddress),
		chainIdxMap:      make(map[int64]btcaddr.Address),
		lastChainIdx:     rootKeyChainIdx,
		missingKeysStart: rootKeyChainIdx,
		secret:           aeskey,
	}
	copy(s.desc[:], desc)
	// Create new root address from key and chaincode.
	root, e := newRootBtcAddress(
		s, rootkey, nil, chaincode,
		createdAt,
	)
	if e != nil {
		return nil, e
	}
	// Verify root address keypairs.
	if e := root.verifyKeypairs(); E.Chk(e) {
		return nil, e
	}
	if e := root.encrypt(aeskey); E.Chk(e) {
		return nil, e
	}
	s.keyGenerator = *root
	// Add root address to maps.
	rootAddr := s.keyGenerator.Address()
	s.addrMap[getAddressKey(rootAddr)] = &s.keyGenerator
	s.chainIdxMap[rootKeyChainIdx] = rootAddr
	// key store must be returned locked.
	if e := s.Lock(); E.Chk(e) {
		return nil, e
	}
	return s, nil
}

// ReadFrom reads data from a io.Reader and saves it to a key store, returning the number of bytes read and any errors
// encountered.
func (s *Store) ReadFrom(r io.Reader) (n int64, e error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	var read int64
	s.net = &netParams{}
	s.addrMap = make(map[addressKey]walletAddress)
	s.chainIdxMap = make(map[int64]btcaddr.Address)
	var id [8]byte
	appendedEntries := varEntries{store: s}
	s.keyGenerator.store = s
	// Iterate through each entry needing to be read. If data implements io.ReaderFrom, use its ReadFrom func.
	// Otherwise, data is a pointer to a fixed sized value.
	datas := []interface{}{
		&id,
		&s.vers,
		s.net,
		&s.flags,
		make([]byte, 6), // Hash for Armory unique ID
		&s.createDate,
		&s.name,
		&s.desc,
		&s.highestUsed,
		&s.kdfParams,
		make([]byte, 256),
		&s.keyGenerator,
		newUnusedSpace(1024, &s.recent),
		&appendedEntries,
	}
	for _, data := range datas {
		var e error
		switch d := data.(type) {
		case readerFromVersion:
			read, e = d.readFromVersion(s.vers, r)
		case io.ReaderFrom:
			read, e = d.ReadFrom(r)
		default:
			read, e = binaryRead(r, binary.LittleEndian, d)
		}
		n += read
		if e != nil {
			return n, e
		}
	}
	if id != fileID {
		return n, errors.New("unknown file ID")
	}
	// Add root address to address map.
	rootAddr := s.keyGenerator.Address()
	s.addrMap[getAddressKey(rootAddr)] = &s.keyGenerator
	s.chainIdxMap[rootKeyChainIdx] = rootAddr
	s.lastChainIdx = rootKeyChainIdx
	// Fill unserializied fields.
	wts := appendedEntries.entries
	for _, wt := range wts {
		switch e := wt.(type) {
		case *addrEntry:
			addr := e.addr.Address()
			s.addrMap[getAddressKey(addr)] = &e.addr
			if e.addr.Imported() {
				s.importedAddrs = append(s.importedAddrs, &e.addr)
			} else {
				s.chainIdxMap[e.addr.chainIndex] = addr
				if s.lastChainIdx < e.addr.chainIndex {
					s.lastChainIdx = e.addr.chainIndex
				}
			}
			// If the private keys have not been created yet, mark the
			// earliest so all can be created on next key store unlock.
			if e.addr.flags.createPrivKeyNextUnlock {
				switch {
				case s.missingKeysStart == rootKeyChainIdx:
					fallthrough
				case e.addr.chainIndex < s.missingKeysStart:
					s.missingKeysStart = e.addr.chainIndex
				}
			}
		case *scriptEntry:
			addr := e.script.Address()
			s.addrMap[getAddressKey(addr)] = &e.script
			// script are always imported.
			s.importedAddrs = append(s.importedAddrs, &e.script)
		default:
			return n, errors.New("unknown appended entry")
		}
	}
	return n, nil
}

// WriteTo serializes a key store and writes it to a io.Writer, returning the number of bytes written and any errors
// encountered.
func (s *Store) WriteTo(w io.Writer) (n int64, e error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.writeTo(w)
}
func (s *Store) writeTo(w io.Writer) (n int64, e error) {
	var wts []io.WriterTo
	var chainedAddrs = make([]io.WriterTo, len(s.chainIdxMap)-1)
	var importedAddrs []io.WriterTo
	for _, wAddr := range s.addrMap {
		switch btcAddr := wAddr.(type) {
		case *btcAddress:
			e := &addrEntry{
				addr: *btcAddr,
			}
			copy(e.pubKeyHash160[:], btcAddr.AddrHash())
			if btcAddr.Imported() {
				// No order for imported addresses.
				importedAddrs = append(importedAddrs, e)
			} else if btcAddr.chainIndex >= 0 {
				// Chained addresses are sorted.  This is
				// kind of nice but probably isn't necessary.
				chainedAddrs[btcAddr.chainIndex] = e
			}
		case *scriptAddress:
			e := &scriptEntry{
				script: *btcAddr,
			}
			copy(e.scriptHash160[:], btcAddr.AddrHash())
			// scripts are always imported
			importedAddrs = append(importedAddrs, e)
		}
	}
	wts = append(chainedAddrs, importedAddrs...)
	appendedEntries := varEntries{store: s, entries: wts}
	// Iterate through each entry needing to be written. If data implements io.WriterTo, use its WriteTo func.
	// Otherwise, data is a pointer to a fixed size value.
	datas := []interface{}{
		&fileID,
		&VersCurrent,
		s.net,
		&s.flags,
		make([]byte, 6), // Hash for Armory unique ID
		&s.createDate,
		&s.name,
		&s.desc,
		&s.highestUsed,
		&s.kdfParams,
		make([]byte, 256),
		&s.keyGenerator,
		newUnusedSpace(1024, &s.recent),
		&appendedEntries,
	}
	var written int64
	for _, data := range datas {
		if s, ok := data.(io.WriterTo); ok {
			written, e = s.WriteTo(w)
		} else {
			written, e = binaryWrite(w, binary.LittleEndian, data)
		}
		n += written
		if e != nil {
			return n, e
		}
	}
	return n, nil
}

// TODO: set this automatically.

// MarkDirty is
func (s *Store) MarkDirty() {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.dirty = true
}

// WriteIfDirty is
func (s *Store) WriteIfDirty() (e error) {
	s.mtx.RLock()
	if !s.dirty {
		s.mtx.RUnlock()
		return nil
	}
	// TempFile creates the file 0600, so no need to chmod it.
	fi, e := ioutil.TempFile(s.dir, s.file)
	if e != nil {
		s.mtx.RUnlock()
		return e
	}
	fiPath := fi.Name()
	_, e = s.writeTo(fi)
	if e != nil {
		s.mtx.RUnlock()
		if e = fi.Close(); E.Chk(e) {
		}
		return e
	}
	e = fi.Sync()
	if e != nil {
		s.mtx.RUnlock()
		if e = fi.Close(); E.Chk(e) {
		}
		return e
	}
	if e = fi.Close(); E.Chk(e) {
	}
	e = rename.Atomic(fiPath, s.path)
	s.mtx.RUnlock()
	if e == nil {
		s.mtx.Lock()
		s.dirty = false
		s.mtx.Unlock()
	}
	return e
}

// OpenDir opens a new key store from the specified directory.
//
// If the file does not exist, the error from the os package will be returned, and can be checked with os.IsNotExist to
// differentiate missing file errors from others (including deserialization).
func OpenDir(dir string) (*Store, error) {
	path := filepath.Join(dir, Filename)
	fi, e := os.OpenFile(path, os.O_RDONLY, 0)
	if e != nil {
		return nil, e
	}
	defer func() {
		if e = fi.Close(); E.Chk(e) {
		}
	}()
	store := new(Store)
	_, e = store.ReadFrom(fi)
	if e != nil {
		return nil, e
	}
	store.path = path
	store.dir = dir
	store.file = Filename
	return store, nil
}

// Unlock derives an AES key from passphrase and key store's KDF parameters and unlocks the root key of the key store.
// If the unlock was successful, the key store's secret key is saved, allowing the decryption of any encrypted private
// key. Any addresses created while the key store was locked without private keys are created at this time.
func (s *Store) Unlock(passphrase []byte) (e error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if s.flags.watchingOnly {
		return ErrWatchingOnly
	}
	// Derive key from KDF parameters and passphrase.
	key := kdf(passphrase, &s.kdfParams)
	// Unlock root address with derived key.
	if _, e = s.keyGenerator.unlock(key); E.Chk(e) {
		return e
	}
	// If unlock was successful, save the passphrase and aes key.
	s.passphrase = passphrase
	s.secret = key
	return s.createMissingPrivateKeys()
}

// Lock performs a best try effort to remove and zero all secret keys associated with the key store.
func (s *Store) Lock() (e error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if s.flags.watchingOnly {
		return ErrWatchingOnly
	}
	// Remove clear text passphrase from key store.
	if s.isLocked() {
		e = ErrLocked
	} else {
		zero(s.passphrase)
		s.passphrase = nil
		zero(s.secret)
		s.secret = nil
	}
	// Remove clear text private keys from all address entries.
	for _, addr := range s.addrMap {
		if baddr, ok := addr.(*btcAddress); ok {
			_ = baddr.lock()
		}
	}
	return e
}

// ChangePassphrase creates a new AES key from a new passphrase and re-encrypts all encrypted private keys with the new
// key.
func (s *Store) ChangePassphrase(new []byte) (e error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if s.flags.watchingOnly {
		return ErrWatchingOnly
	}
	if s.isLocked() {
		return ErrLocked
	}
	oldkey := s.secret
	newkey := kdf(new, &s.kdfParams)
	for _, wa := range s.addrMap {
		// Only btcAddresses curently have private keys.
		a, ok := wa.(*btcAddress)
		if !ok {
			continue
		}
		if e := a.changeEncryptionKey(oldkey, newkey); E.Chk(e) {
			return e
		}
	}
	// zero old secrets.
	zero(s.passphrase)
	zero(s.secret)
	// Save new secrets.
	s.passphrase = new
	s.secret = newkey
	return nil
}
func zero(b []byte) {
	for i := range b {
		b[i] = 0
	}
}

// IsLocked returns whether a key store is unlocked (in which case the key is saved in memory), or locked.
func (s *Store) IsLocked() bool {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.isLocked()
}
func (s *Store) isLocked() bool {
	return len(s.secret) != 32
}

// NextChainedAddress attempts to get the next chained address. If the key store is unlocked, the next pubkey and
// private key of the address chain are derived.
//
// If the key store is locke, only the next pubkey is derived, and the private key will be generated on next unlock.
func (s *Store) NextChainedAddress(bs *BlockStamp) (btcaddr.Address, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	return s.nextChainedAddress(bs)
}
func (s *Store) nextChainedAddress(bs *BlockStamp) (btcaddr.Address, error) {
	addr, e := s.nextChainedBtcAddress(bs)
	if e != nil {
		return nil, e
	}
	return addr.Address(), nil
}

// ChangeAddress returns the next chained address from the key store, marking the address for a change transaction
// output.
func (s *Store) ChangeAddress(bs *BlockStamp) (btcaddr.Address, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	addr, e := s.nextChainedBtcAddress(bs)
	if e != nil {
		return nil, e
	}
	addr.flags.change = true
	// Create and return payment address for address hash.
	return addr.Address(), nil
}
func (s *Store) nextChainedBtcAddress(bs *BlockStamp) (*btcAddress, error) {
	// Attempt to get address hash of next chained address.
	nextAPKH, ok := s.chainIdxMap[s.highestUsed+1]
	if !ok {
		if s.isLocked() {
			// Chain pubkeys.
			if e := s.extendLocked(bs); E.Chk(e) {
				return nil, e
			}
		} else {
			// Chain private and pubkeys.
			if e := s.extendUnlocked(bs); E.Chk(e) {
				return nil, e
			}
		}
		// Should be added to the internal maps, try lookup again.
		nextAPKH, ok = s.chainIdxMap[s.highestUsed+1]
		if !ok {
			return nil, errors.New("chain index map inproperly updated")
		}
	}
	// Look up address.
	addr, ok := s.addrMap[getAddressKey(nextAPKH)]
	if !ok {
		return nil, errors.New("cannot find generated address")
	}
	btcAddr, ok := addr.(*btcAddress)
	if !ok {
		return nil, errors.New("found non-pubkey chained address")
	}
	s.highestUsed++
	return btcAddr, nil
}

// LastChainedAddress returns the most recently requested chained address from calling NextChainedAddress, or the root
// address if no chained addresses have been requested.
func (s *Store) LastChainedAddress() btcaddr.Address {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.chainIdxMap[s.highestUsed]
}

// extendUnlocked grows address chain for an unlocked keystore.
func (s *Store) extendUnlocked(bs *BlockStamp) (e error) {
	// Get last chained address. New chained addresses will be chained off of this address's chaincode and private key.
	a := s.chainIdxMap[s.lastChainIdx]
	waddr, ok := s.addrMap[getAddressKey(a)]
	if !ok {
		return errors.New("expected last chained address not found")
	}
	if s.isLocked() {
		return ErrLocked
	}
	lastAddr, ok := waddr.(*btcAddress)
	if !ok {
		return errors.New("found non-pubkey chained address")
	}
	privkey, e := lastAddr.unlock(s.secret)
	if e != nil {
		return e
	}
	cc := lastAddr.chaincode[:]
	privkey, e = chainedPrivKey(privkey, lastAddr.pubKeyBytes(), cc)
	if e != nil {
		return e
	}
	newAddr, e := newBtcAddress(s, privkey, nil, bs, true)
	if e != nil {
		return e
	}
	if e = newAddr.verifyKeypairs(); E.Chk(e) {
		return e
	}
	if e = newAddr.encrypt(s.secret); E.Chk(e) {
		return e
	}
	a = newAddr.Address()
	s.addrMap[getAddressKey(a)] = newAddr
	newAddr.chainIndex = lastAddr.chainIndex + 1
	s.chainIdxMap[newAddr.chainIndex] = a
	s.lastChainIdx++
	copy(newAddr.chaincode[:], cc)
	return nil
}

// extendLocked creates one new address without a private key (allowing for extending the address chain from a locked
// key store) chained from the last used chained address and adds the address to the key store's internal bookkeeping
// structures.
func (s *Store) extendLocked(bs *BlockStamp) (e error) {
	a := s.chainIdxMap[s.lastChainIdx]
	waddr, ok := s.addrMap[getAddressKey(a)]
	if !ok {
		return errors.New("expected last chained address not found")
	}
	addr, ok := waddr.(*btcAddress)
	if !ok {
		return errors.New("found non-pubkey chained address")
	}
	cc := addr.chaincode[:]
	nextPubkey, e := chainedPubKey(addr.pubKeyBytes(), cc)
	if e != nil {
		return e
	}
	newaddr, e := newBtcAddressWithoutPrivkey(s, nextPubkey, nil, bs)
	if e != nil {
		return e
	}
	a = newaddr.Address()
	s.addrMap[getAddressKey(a)] = newaddr
	newaddr.chainIndex = addr.chainIndex + 1
	s.chainIdxMap[newaddr.chainIndex] = a
	s.lastChainIdx++
	copy(newaddr.chaincode[:], cc)
	if s.missingKeysStart == rootKeyChainIdx {
		s.missingKeysStart = newaddr.chainIndex
	}
	return nil
}
func (s *Store) createMissingPrivateKeys() (e error) {
	idx := s.missingKeysStart
	if idx == rootKeyChainIdx {
		return nil
	}
	// Lookup previous address.
	apkh, ok := s.chainIdxMap[idx-1]
	if !ok {
		return errors.New("missing previous chained address")
	}
	prevWAddr := s.addrMap[getAddressKey(apkh)]
	if s.isLocked() {
		return ErrLocked
	}
	prevAddr, ok := prevWAddr.(*btcAddress)
	if !ok {
		return errors.New("found non-pubkey chained address")
	}
	prevPrivKey, e := prevAddr.unlock(s.secret)
	if e != nil {
		return e
	}
	for i := idx; ; i++ {
		// Get the next private key for the ith address in the address chain.
		ithPrivKey, e := chainedPrivKey(
			prevPrivKey,
			prevAddr.pubKeyBytes(), prevAddr.chaincode[:],
		)
		if e != nil {
			return e
		}
		// Get the address with the missing private key, set, and encrypt.
		apkh, ok := s.chainIdxMap[i]
		if !ok {
			// Finished.
			break
		}
		waddr := s.addrMap[getAddressKey(apkh)]
		addr, ok := waddr.(*btcAddress)
		if !ok {
			return errors.New("found non-pubkey chained address")
		}
		addr.privKeyCT = ithPrivKey
		if e := addr.encrypt(s.secret); E.Chk(e) {
			// Avoid bug: see comment for VersUnsetNeedsPrivkeyFlag.
			if e != ErrAlreadyEncrypted || s.vers.LT(VersUnsetNeedsPrivkeyFlag) {
				return e
			}
		}
		addr.flags.createPrivKeyNextUnlock = false
		// Set previous address and private key for next iteration.
		prevAddr = addr
		prevPrivKey = ithPrivKey
	}
	s.missingKeysStart = rootKeyChainIdx
	return nil
}

// Address returns an walletAddress structure for an address in a key store. This address may be typecast into other
// interfaces (like PubKeyAddress and ScriptAddress) if specific information e.g. keys is required.
func (s *Store) Address(a btcaddr.Address) (WalletAddress, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	// Look up address by address hash.
	btcaddr, ok := s.addrMap[getAddressKey(a)]
	if !ok {
		return nil, ErrAddressNotFound
	}
	return btcaddr, nil
}

// Net returns the bitcoin network parameters for this key store.
func (s *Store) Net() *chaincfg.Params {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.netParams()
}
func (s *Store) netParams() *chaincfg.Params {
	return (*chaincfg.Params)(s.net)
}

// SetSyncStatus sets the sync status for a single key store address. This may error if the address is not found in the
// key store.
//
// When marking an address as unsynced, only the type Unsynced matters. The value is ignored.
func (s *Store) SetSyncStatus(a btcaddr.Address, ss SyncStatus) (e error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	wa, ok := s.addrMap[getAddressKey(a)]
	if !ok {
		return ErrAddressNotFound
	}
	wa.setSyncStatus(ss)
	return nil
}

// SetSyncedWith marks already synced addresses in the key store to be in sync with the recently-seen block described by
// the blockstamp.
//
// Unsynced addresses are unaffected by this method and must be marked as in sync with MarkAddressSynced or
// MarkAllSynced to be considered in sync with bs.
//
// If bs is nil, the entire key store is marked unsynced.
func (s *Store) SetSyncedWith(bs *BlockStamp) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if bs == nil {
		s.recent.hashes = s.recent.hashes[:0]
		s.recent.lastHeight = s.keyGenerator.firstBlock
		s.keyGenerator.setSyncStatus(Unsynced(s.keyGenerator.firstBlock))
		return
	}
	// Chk if we're trying to rollback the last seen history.
	//
	// If so, and this bs is already saved, remove anything after and return. Otherwise, remove previous hashes.
	if bs.Height < s.recent.lastHeight {
		maybeIdx := len(s.recent.hashes) - 1 - int(s.recent.lastHeight-bs.Height)
		if maybeIdx >= 0 && maybeIdx < len(s.recent.hashes) &&
			*s.recent.hashes[maybeIdx] == *bs.Hash {
			s.recent.lastHeight = bs.Height
			// subslice out the removed hashes.
			s.recent.hashes = s.recent.hashes[:maybeIdx]
			return
		}
		s.recent.hashes = nil
	}
	if bs.Height != s.recent.lastHeight+1 {
		s.recent.hashes = nil
	}
	s.recent.lastHeight = bs.Height
	if len(s.recent.hashes) == 20 {
		// Make room for the most recent hash.
		copy(s.recent.hashes, s.recent.hashes[1:])
		// Set new block in the last position.
		s.recent.hashes[19] = bs.Hash
	} else {
		s.recent.hashes = append(s.recent.hashes, bs.Hash)
	}
}

// SyncedTo returns details about the block that a wallet is marked at least synced through. The height is the height
// that rescans should start at when syncing a wallet back to the best chain.
//
// NOTE: If the hash of the synced block is not known, hash will be nil, and must be obtained from elsewhere. This must
// be explicitly checked before dereferencing the pointer.
func (s *Store) SyncedTo() (hash *chainhash.Hash, height int32) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	switch h, ok := s.keyGenerator.SyncStatus().(PartialSync); {
	case ok && int32(h) > s.recent.lastHeight:
		height = int32(h)
	default:
		height = s.recent.lastHeight
		if n := len(s.recent.hashes); n != 0 {
			hash = s.recent.hashes[n-1]
		}
	}
	for _, a := range s.addrMap {
		var syncHeight int32
		switch e := a.SyncStatus().(type) {
		case Unsynced:
			syncHeight = int32(e)
		case PartialSync:
			syncHeight = int32(e)
		case FullSync:
			continue
		}
		if syncHeight < height {
			height = syncHeight
			hash = nil
			// Can't go lower than 0.
			if height == 0 {
				return
			}
		}
	}
	return
}

// NewIterateRecentBlocks returns an iterator for recently-seen blocks. The iterator starts at the most recently-added
// block, and Prev should be used to access earlier blocks.
func (s *Store) NewIterateRecentBlocks() *BlockIterator {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.recent.iter(s)
}

// ImportPrivateKey imports a WIF private key into the keystore. The imported address is created using either a
// compressed or uncompressed serialized public key, depending on the CompressPubKey bool of the WIF.
func (s *Store) ImportPrivateKey(wif *util.WIF, bs *BlockStamp) (btcaddr.Address, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if s.flags.watchingOnly {
		return nil, ErrWatchingOnly
	}
	// First, must check that the key being imported will not result in a duplicate address.
	pkh := btcaddr.Hash160(wif.SerializePubKey())
	if _, ok := s.addrMap[addressKey(pkh)]; ok {
		return nil, ErrDuplicate
	}
	// The key store must be unlocked to encrypt the imported private key.
	if s.isLocked() {
		return nil, ErrLocked
	}
	// Create new address with this private key.
	privKey := wif.PrivKey.Serialize()
	btcaddr, e := newBtcAddress(s, privKey, nil, bs, wif.CompressPubKey)
	if e != nil {
		return nil, e
	}
	btcaddr.chainIndex = importedKeyChainIdx
	// Mark as unsynced if import height is below currently-synced height.
	if len(s.recent.hashes) != 0 && bs.Height < s.recent.lastHeight {
		btcaddr.flags.unsynced = true
	}
	// Encrypt imported address with the derived AES key.
	if e = btcaddr.encrypt(s.secret); E.Chk(e) {
		return nil, e
	}
	addr := btcaddr.Address()
	// Add address to key store's bookkeeping structures. Adding to the map will result in the imported address being
	// serialized on the next WriteTo call.
	s.addrMap[getAddressKey(addr)] = btcaddr
	s.importedAddrs = append(s.importedAddrs, btcaddr)
	// Create and return address.
	return addr, nil
}

// ImportScript creates a new scriptAddress with a user-provided script and adds it to the key store.
func (s *Store) ImportScript(script []byte, bs *BlockStamp) (btcaddr.Address, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if s.flags.watchingOnly {
		return nil, ErrWatchingOnly
	}
	if _, ok := s.addrMap[addressKey(btcaddr.Hash160(script))]; ok {
		return nil, ErrDuplicate
	}
	// Create new address with this private key.
	scriptaddr, e := newScriptAddress(s, script, bs)
	if e != nil {
		return nil, e
	}
	// Mark as unsynced if import height is below currently-synced height.
	if len(s.recent.hashes) != 0 && bs.Height < s.recent.lastHeight {
		scriptaddr.flags.unsynced = true
	}
	// Add address to key store's bookkeeping structures. Adding to the map will result in the imported address being
	// serialized on the next WriteTo call.
	addr := scriptaddr.Address()
	s.addrMap[getAddressKey(addr)] = scriptaddr
	s.importedAddrs = append(s.importedAddrs, scriptaddr)
	// Create and return address.
	return addr, nil
}

// CreateDate returns the Unix time of the key store creation time. This is used to compare the key store creation time
// against block headers and set a better minimum block height of where to being rescans.
func (s *Store) CreateDate() int64 {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.createDate
}

// ExportWatchingWallet creates and returns a new key store with the same addresses in w, but as a watching-only key
// store without any private keys.
//
// New addresses created by the watching key store will match the new addresses created the original key store (thanks
// to public key address chaining), but will be missing the associated private keys.
func (s *Store) ExportWatchingWallet() (*Store, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	// Don't continue if key store is already watching-only.
	if s.flags.watchingOnly {
		return nil, ErrWatchingOnly
	}
	// Copy members of w into a new key store, but mark as watching-only and do not include any private keys.
	ws := &Store{
		vers: s.vers,
		net:  s.net,
		flags: walletFlags{
			useEncryption: false,
			watchingOnly:  true,
		},
		name:        s.name,
		desc:        s.desc,
		createDate:  s.createDate,
		highestUsed: s.highestUsed,
		recent: recentBlocks{
			lastHeight: s.recent.lastHeight,
		},
		addrMap: make(map[addressKey]walletAddress),
		// todo oga make me a list
		chainIdxMap:  make(map[int64]btcaddr.Address),
		lastChainIdx: s.lastChainIdx,
	}
	kgwc := s.keyGenerator.watchingCopy(ws)
	ws.keyGenerator = *(kgwc.(*btcAddress))
	if len(s.recent.hashes) != 0 {
		ws.recent.hashes = make([]*chainhash.Hash, 0, len(s.recent.hashes))
		for _, hash := range s.recent.hashes {
			hashCpy := *hash
			ws.recent.hashes = append(ws.recent.hashes, &hashCpy)
		}
	}
	for apkh, addr := range s.addrMap {
		if !addr.Imported() {
			// Must be a btcAddress if !imported.
			btcAddr := addr.(*btcAddress)
			ws.chainIdxMap[btcAddr.chainIndex] =
				addr.Address()
		}
		apkhCopy := apkh
		ws.addrMap[apkhCopy] = addr.watchingCopy(ws)
	}
	if len(s.importedAddrs) != 0 {
		ws.importedAddrs = make(
			[]walletAddress, 0,
			len(s.importedAddrs),
		)
		for _, addr := range s.importedAddrs {
			ws.importedAddrs = append(ws.importedAddrs, addr.watchingCopy(ws))
		}
	}
	return ws, nil
}

// SyncStatus is the interface type for all sync variants.
type SyncStatus interface {
	ImplementsSyncStatus()
}
type (
	// Unsynced is a type representing an unsynced address. When this is returned by a key store method, the value is
	// the recorded first seen block height.
	Unsynced int32
	// PartialSync is a type representing a partially synced address (for example, due to the result of a
	// partially-completed rescan).
	PartialSync int32
	// FullSync is a type representing an address that is in sync with the recently seen blocks.
	FullSync struct{}
)

// ImplementsSyncStatus is implemented to make Unsynced a SyncStatus.
func (u Unsynced) ImplementsSyncStatus() {
}

// ImplementsSyncStatus is implemented to make PartialSync a SyncStatus.
func (p PartialSync) ImplementsSyncStatus() {
}

// ImplementsSyncStatus is implemented to make FullSync a SyncStatus.
func (f FullSync) ImplementsSyncStatus() {
}

// WalletAddress is an interface that provides acces to information regarding an address managed by a key store.
// Concrete implementations of this type may provide further fields to provide information specific to that type of
// address.
type WalletAddress interface {
	// Address returns a util.Address for the backing address.
	Address() btcaddr.Address
	// AddrHash returns the key or script hash related to the address
	AddrHash() string
	// FirstBlock returns the first block an address could be in.
	FirstBlock() int32
	// Imported returns true if the backing address was imported instead
	// of being part of an address chain.
	Imported() bool
	// Change returns true if the backing address was created for a change output
	// of a transaction.
	Change() bool
	// Compressed returns true if the backing address is compressed.
	Compressed() bool
	// SyncStatus returns the current synced state of an address.
	SyncStatus() SyncStatus
}

// SortedActiveAddresses returns all key store addresses that have been requested to be generated. These do not include
// unused addresses in the key pool. Use this when ordered addresses are needed. Otherwise, ActiveAddresses is
// preferred.
func (s *Store) SortedActiveAddresses() []WalletAddress {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	addrs := make(
		[]WalletAddress, 0,
		s.highestUsed+int64(len(s.importedAddrs))+1,
	)
	for i := int64(rootKeyChainIdx); i <= s.highestUsed; i++ {
		a := s.chainIdxMap[i]
		info, ok := s.addrMap[getAddressKey(a)]
		if ok {
			addrs = append(addrs, info)
		}
	}
	for _, addr := range s.importedAddrs {
		addrs = append(addrs, addr)
	}
	return addrs
}

// ActiveAddresses returns a map between active payment addresses and their full info. These do not include unused
// addresses in the key pool. If addresses must be sorted, use SortedActiveAddresses.
func (s *Store) ActiveAddresses() map[btcaddr.Address]WalletAddress {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	addrs := make(map[btcaddr.Address]WalletAddress)
	for i := int64(rootKeyChainIdx); i <= s.highestUsed; i++ {
		a := s.chainIdxMap[i]
		addr := s.addrMap[getAddressKey(a)]
		addrs[addr.Address()] = addr
	}
	for _, addr := range s.importedAddrs {
		addrs[addr.Address()] = addr
	}
	return addrs
}

// ExtendActiveAddresses gets or creates the next n addresses from the address chain and marks each as active. This is
// used to recover deterministic (not imported) addresses from a key store backup, or to keep the active addresses in
// sync between an encrypted key store with private keys and an exported watching key store without.
//
// A slice is returned with the util.Address of each new address. The blockchain must be rescanned for these addresses.
func (s *Store) ExtendActiveAddresses(n int) ([]btcaddr.Address, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	last := s.addrMap[getAddressKey(s.chainIdxMap[s.highestUsed])]
	bs := &BlockStamp{Height: last.FirstBlock()}
	addrs := make([]btcaddr.Address, n)
	for i := 0; i < n; i++ {
		addr, e := s.nextChainedAddress(bs)
		if e != nil {
			return nil, e
		}
		addrs[i] = addr
	}
	return addrs, nil
}

type walletFlags struct {
	useEncryption bool
	watchingOnly  bool
}

func (wf *walletFlags) ReadFrom(r io.Reader) (int64, error) {
	var b [8]byte
	n, e := io.ReadFull(r, b[:])
	if e != nil {
		return int64(n), e
	}
	wf.useEncryption = b[0]&(1<<0) != 0
	wf.watchingOnly = b[0]&(1<<1) != 0
	return int64(n), nil
}
func (wf *walletFlags) WriteTo(w io.Writer) (int64, error) {
	var b [8]byte
	if wf.useEncryption {
		b[0] |= 1 << 0
	}
	if wf.watchingOnly {
		b[0] |= 1 << 1
	}
	n, e := w.Write(b[:])
	return int64(n), e
}

type addrFlags struct {
	hasPrivKey              bool
	hasPubKey               bool
	encrypted               bool
	createPrivKeyNextUnlock bool
	compressed              bool
	change                  bool
	unsynced                bool
	partialSync             bool
}

func (af *addrFlags) ReadFrom(r io.Reader) (int64, error) {
	var b [8]byte
	n, e := io.ReadFull(r, b[:])
	if e != nil {
		return int64(n), e
	}
	af.hasPrivKey = b[0]&(1<<0) != 0
	af.hasPubKey = b[0]&(1<<1) != 0
	af.encrypted = b[0]&(1<<2) != 0
	af.createPrivKeyNextUnlock = b[0]&(1<<3) != 0
	af.compressed = b[0]&(1<<4) != 0
	af.change = b[0]&(1<<5) != 0
	af.unsynced = b[0]&(1<<6) != 0
	af.partialSync = b[0]&(1<<7) != 0
	// Currently (at least until watching-only key stores are implemented) btcwallet shall refuse to open any
	// unencrypted addresses. This check only makes sense if there is a private key to encrypt, which there may not be
	// if the keypool was extended from just the last public key and no private keys were written.
	if af.hasPrivKey && !af.encrypted {
		return int64(n), errors.New("private key is unencrypted")
	}
	return int64(n), nil
}
func (af *addrFlags) WriteTo(w io.Writer) (int64, error) {
	var b [8]byte
	if af.hasPrivKey {
		b[0] |= 1 << 0
	}
	if af.hasPubKey {
		b[0] |= 1 << 1
	}
	if af.hasPrivKey && !af.encrypted {
		// We only support encrypted privkeys.
		return 0, errors.New("address must be encrypted")
	}
	if af.encrypted {
		b[0] |= 1 << 2
	}
	if af.createPrivKeyNextUnlock {
		b[0] |= 1 << 3
	}
	if af.compressed {
		b[0] |= 1 << 4
	}
	if af.change {
		b[0] |= 1 << 5
	}
	if af.unsynced {
		b[0] |= 1 << 6
	}
	if af.partialSync {
		b[0] |= 1 << 7
	}
	n, e := w.Write(b[:])
	return int64(n), e
}

// recentBlocks holds at most the last 20 seen block hashes as well as the block height of the most recently seen block.
type recentBlocks struct {
	hashes     []*chainhash.Hash
	lastHeight int32
}

func (rb *recentBlocks) readFromVersion(v ksVersion, r io.Reader) (int64, error) {
	if !v.LT(Vers20LastBlocks) {
		// Use current ksVersion.
		return rb.ReadFrom(r)
	}
	// Old file versions only saved the most recently seen block height and hash, not the last 20.
	var read int64
	// Read height.
	var heightBytes [4]byte // 4 bytes for a int32
	n, e := io.ReadFull(r, heightBytes[:])
	read += int64(n)
	if e != nil {
		return read, e
	}
	rb.lastHeight = int32(binary.LittleEndian.Uint32(heightBytes[:]))
	// If height is -1, the last synced block is unknown, so don't try to read a block hash.
	if rb.lastHeight == -1 {
		rb.hashes = nil
		return read, nil
	}
	// Read block hash.
	var syncedBlockHash chainhash.Hash
	n, e = io.ReadFull(r, syncedBlockHash[:])
	read += int64(n)
	if e != nil {
		return read, e
	}
	rb.hashes = []*chainhash.Hash{
		&syncedBlockHash,
	}
	return read, nil
}
func (rb *recentBlocks) ReadFrom(r io.Reader) (int64, error) {
	var read int64
	// Read number of saved blocks.  This should not exceed 20.
	var nBlockBytes [4]byte // 4 bytes for a uint32
	n, e := io.ReadFull(r, nBlockBytes[:])
	read += int64(n)
	if e != nil {
		return read, e
	}
	nBlocks := binary.LittleEndian.Uint32(nBlockBytes[:])
	if nBlocks > 20 {
		return read, errors.New("number of last seen blocks exceeds maximum of 20")
	}
	// Read most recently seen block height.
	var heightBytes [4]byte // 4 bytes for a int32
	n, e = io.ReadFull(r, heightBytes[:])
	read += int64(n)
	if e != nil {
		return read, e
	}
	height := int32(binary.LittleEndian.Uint32(heightBytes[:]))
	// height should not be -1 (or any other negative number) since at this point we should be reading in at least one
	// known block.
	if height < 0 {
		return read, errors.New("expected a block but specified height is negative")
	}
	// Set last seen height.
	rb.lastHeight = height
	// Read nBlocks block hashes. Merkles are expected to be in order of oldest to newest, but there's no way to check
	// that here.
	rb.hashes = make([]*chainhash.Hash, 0, nBlocks)
	for i := uint32(0); i < nBlocks; i++ {
		var blockHash chainhash.Hash
		n, e := io.ReadFull(r, blockHash[:])
		read += int64(n)
		if e != nil {
			return read, e
		}
		rb.hashes = append(rb.hashes, &blockHash)
	}
	return read, nil
}
func (rb *recentBlocks) WriteTo(w io.Writer) (int64, error) {
	var written int64
	// Write number of saved blocks.  This should not exceed 20.
	nBlocks := uint32(len(rb.hashes))
	if nBlocks > 20 {
		return written, errors.New("number of last seen blocks exceeds maximum of 20")
	}
	if nBlocks != 0 && rb.lastHeight < 0 {
		return written, errors.New("number of block hashes is positive, but height is negative")
	}
	var nBlockBytes [4]byte // 4 bytes for a uint32
	binary.LittleEndian.PutUint32(nBlockBytes[:], nBlocks)
	n, e := w.Write(nBlockBytes[:])
	written += int64(n)
	if e != nil {
		return written, e
	}
	// Write most recently seen block height.
	var heightBytes [4]byte // 4 bytes for a int32
	binary.LittleEndian.PutUint32(heightBytes[:], uint32(rb.lastHeight))
	n, e = w.Write(heightBytes[:])
	written += int64(n)
	if e != nil {
		return written, e
	}
	// Write block hashes.
	for _, hash := range rb.hashes {
		n, e := w.Write(hash[:])
		written += int64(n)
		if e != nil {
			return written, e
		}
	}
	return written, nil
}

// BlockIterator allows for the forwards and backwards iteration of recently seen blocks.
type BlockIterator struct {
	storeMtx *sync.RWMutex
	height   int32
	index    int
	rb       *recentBlocks
}

func (rb *recentBlocks) iter(s *Store) *BlockIterator {
	if rb.lastHeight == -1 || len(rb.hashes) == 0 {
		return nil
	}
	return &BlockIterator{
		storeMtx: &s.mtx,
		height:   rb.lastHeight,
		index:    len(rb.hashes) - 1,
		rb:       rb,
	}
}

// Next is
func (it *BlockIterator) Next() bool {
	it.storeMtx.RLock()
	defer it.storeMtx.RUnlock()
	if it.index+1 >= len(it.rb.hashes) {
		return false
	}
	it.index++
	return true
}

// Prev is
func (it *BlockIterator) Prev() bool {
	it.storeMtx.RLock()
	defer it.storeMtx.RUnlock()
	if it.index-1 < 0 {
		return false
	}
	it.index--
	return true
}

// BlockStamp is
func (it *BlockIterator) BlockStamp() BlockStamp {
	it.storeMtx.RLock()
	defer it.storeMtx.RUnlock()
	return BlockStamp{
		Height: it.rb.lastHeight - int32(len(it.rb.hashes)-1-it.index),
		Hash:   it.rb.hashes[it.index],
	}
}

// unusedSpace is a wrapper type to read or write one or more types that btcwallet fits into an unused space left by
// Armory's key store file format.
type unusedSpace struct {
	nBytes int // number of unused bytes that armory left.
	rfvs   []readerFromVersion
}

func newUnusedSpace(nBytes int, rfvs ...readerFromVersion) *unusedSpace {
	return &unusedSpace{
		nBytes: nBytes,
		rfvs:   rfvs,
	}
}
func (u *unusedSpace) readFromVersion(v ksVersion, r io.Reader) (int64, error) {
	var read int64
	for _, rfv := range u.rfvs {
		n, e := rfv.readFromVersion(v, r)
		if e != nil {
			return read + n, e
		}
		read += n
		if read > int64(u.nBytes) {
			return read, errors.New("read too much from armory's unused space")
		}
	}
	// Read rest of actually unused bytes.
	unused := make([]byte, u.nBytes-int(read))
	n, e := io.ReadFull(r, unused)
	return read + int64(n), e
}
func (u *unusedSpace) WriteTo(w io.Writer) (int64, error) {
	var written int64
	for _, wt := range u.rfvs {
		n, e := wt.WriteTo(w)
		if e != nil {
			return written + n, e
		}
		written += n
		if written > int64(u.nBytes) {
			return written, errors.New("wrote too much to armory's unused space")
		}
	}
	// Write rest of actually unused bytes.
	unused := make([]byte, u.nBytes-int(written))
	n, e := w.Write(unused)
	return written + int64(n), e
}

// walletAddress is the internal interface used to abstracted around the different address types.
type walletAddress interface {
	io.ReaderFrom
	io.WriterTo
	WalletAddress
	watchingCopy(*Store) walletAddress
	setSyncStatus(SyncStatus)
}
type btcAddress struct {
	store             *Store
	address           btcaddr.Address
	flags             addrFlags
	chaincode         [32]byte
	chainIndex        int64
	chainDepth        int64 // unused
	initVector        [16]byte
	privKey           [32]byte
	pubKey            *ec.PublicKey
	firstSeen         int64
	lastSeen          int64
	firstBlock        int32
	partialSyncHeight int32  // This is reappropriated from armory's `lastBlock` field.
	privKeyCT         []byte // non-nil if unlocked.
}

const (
	// Root address has a chain index of -1. Each subsequent chained address increments the index.
	rootKeyChainIdx = -1
	// Imported private keys are not part of the chain, and have a special index of -2.
	importedKeyChainIdx = -2
)
const (
	pubkeyCompressed   byte = 0x2
	pubkeyUncompressed byte = 0x4
)

type publicKey []byte

func (k *publicKey) ReadFrom(r io.Reader) (n int64, e error) {
	var read int64
	var format byte
	read, e = binaryRead(r, binary.LittleEndian, &format)
	if e != nil {
		return n + read, e
	}
	n += read
	// Remove the oddness from the format
	noodd := format
	noodd &= ^byte(0x1)
	var s []byte
	switch noodd {
	case pubkeyUncompressed:
		// Read the remaining 64 bytes.
		s = make([]byte, 64)
	case pubkeyCompressed:
		// Read the remaining 32 bytes.
		s = make([]byte, 32)
	default:
		return n, errors.New("unrecognized pubkey format")
	}
	read, e = binaryRead(r, binary.LittleEndian, &s)
	if e != nil {
		return n + read, e
	}
	n += read
	*k = append([]byte{format}, s...)
	return
}
func (k *publicKey) WriteTo(w io.Writer) (n int64, e error) {
	return binaryWrite(w, binary.LittleEndian, []byte(*k))
}

// PubKeyAddress implements WalletAddress and additionally provides the pubkey for a pubkey-based address.
type PubKeyAddress interface {
	WalletAddress
	// PubKey returns the public key associated with the address.
	PubKey() *ec.PublicKey
	// ExportPubKey returns the public key associated with the address serialised as a hex encoded string.
	ExportPubKey() string
	// PrivKey returns the private key for the address. It can fail if the key store is watching only, the key store is
	// locked, or the address doesn't have any keys.
	PrivKey() (*ec.PrivateKey, error)
	// ExportPrivKey exports the WIF private key.
	ExportPrivKey() (*util.WIF, error)
}

// newBtcAddress initializes and returns a new address. privkey must be 32 bytes. iv must be 16 bytes, or nil (in which
// case it is randomly generated).
func newBtcAddress(wallet *Store, privkey, iv []byte, bs *BlockStamp, compressed bool) (addr *btcAddress, e error) {
	if len(privkey) != 32 {
		return nil, errors.New("private key is not 32 bytes")
	}
	addr, e = newBtcAddressWithoutPrivkey(
		wallet,
		pubkeyFromPrivkey(privkey, compressed), iv, bs,
	)
	if e != nil {
		return nil, e
	}
	addr.flags.createPrivKeyNextUnlock = false
	addr.flags.hasPrivKey = true
	addr.privKeyCT = privkey
	return addr, nil
}

// newBtcAddressWithoutPrivkey initializes and returns a new address with an unknown (at the time) private key that must
// be found later. pubkey must be 33 or 65 bytes, and iv must be 16 bytes or empty (in which case it is randomly
// generated).
func newBtcAddressWithoutPrivkey(s *Store, pubkey, iv []byte, bs *BlockStamp) (addr *btcAddress, e error) {
	var compressed bool
	switch n := len(pubkey); n {
	case ec.PubKeyBytesLenCompressed:
		compressed = true
	case ec.PubKeyBytesLenUncompressed:
		compressed = false
	default:
		return nil, fmt.Errorf("invalid pubkey length %d", n)
	}
	if len(iv) == 0 {
		iv = make([]byte, 16)
		if _, e = rand.Read(iv); E.Chk(e) {
			return nil, e
		}
	} else if len(iv) != 16 {
		return nil, errors.New("init vector must be nil or 16 bytes large")
	}
	pk, e := ec.ParsePubKey(pubkey, ec.S256())
	if e != nil {
		return nil, e
	}
	address, e := btcaddr.NewPubKeyHash(btcaddr.Hash160(pubkey), s.netParams())
	if e != nil {
		return nil, e
	}
	addr = &btcAddress{
		flags: addrFlags{
			hasPrivKey:              false,
			hasPubKey:               true,
			encrypted:               false,
			createPrivKeyNextUnlock: true,
			compressed:              compressed,
			change:                  false,
			unsynced:                false,
		},
		store:      s,
		address:    address,
		firstSeen:  time.Now().Unix(),
		firstBlock: bs.Height,
		pubKey:     pk,
	}
	copy(addr.initVector[:], iv)
	return addr, nil
}

// newRootBtcAddress generates a new address, also setting the chaincode and chain index to represent this address as a
// root address.
func newRootBtcAddress(
	s *Store, privKey, iv, chaincode []byte,
	bs *BlockStamp,
) (addr *btcAddress, e error) {
	if len(chaincode) != 32 {
		return nil, errors.New("chaincode is not 32 bytes")
	}
	// Create new btcAddress with provided inputs. This will always use a compressed pubkey.
	addr, e = newBtcAddress(s, privKey, iv, bs, true)
	if e != nil {
		return nil, e
	}
	copy(addr.chaincode[:], chaincode)
	addr.chainIndex = rootKeyChainIdx
	return addr, e
}

// verifyKeypairs creates a signature using the parsed private key and verifies the signature with the parsed public
// key. If either of these steps fail, the keypair generation failed and any funds sent to this address will be
// unspendable. This step requires an unencrypted or unlocked btcAddress.
func (a *btcAddress) verifyKeypairs() (e error) {
	if len(a.privKeyCT) != 32 {
		return errors.New("private key unavailable")
	}
	privKey := &ec.PrivateKey{
		PublicKey: *a.pubKey.ToECDSA(),
		D:         new(big.Int).SetBytes(a.privKeyCT),
	}
	data := "String to sign."
	sig, e := privKey.Sign([]byte(data))
	if e != nil {
		return e
	}
	ok := sig.Verify([]byte(data), privKey.PubKey())
	if !ok {
		return errors.New("pubkey verification failed")
	}
	return nil
}

// ReadFrom reads an encrypted address from an io.Reader.
func (a *btcAddress) ReadFrom(r io.Reader) (n int64, e error) {
	var read int64
	// Checksums
	var chkPubKeyHash uint32
	var chkChaincode uint32
	var chkInitVector uint32
	var chkPrivKey uint32
	var chkPubKey uint32
	var pubKeyHash [ripemd160.Size]byte
	var pubKey publicKey
	// Read serialized key store into addr fields and checksums.
	datas := []interface{}{
		&pubKeyHash,
		&chkPubKeyHash,
		make([]byte, 4), // ksVersion
		&a.flags,
		&a.chaincode,
		&chkChaincode,
		&a.chainIndex,
		&a.chainDepth,
		&a.initVector,
		&chkInitVector,
		&a.privKey,
		&chkPrivKey,
		&pubKey,
		&chkPubKey,
		&a.firstSeen,
		&a.lastSeen,
		&a.firstBlock,
		&a.partialSyncHeight,
	}
	for _, data := range datas {
		if rf, ok := data.(io.ReaderFrom); ok {
			read, e = rf.ReadFrom(r)
		} else {
			read, e = binaryRead(r, binary.LittleEndian, data)
		}
		if e != nil {
			return n + read, e
		}
		n += read
	}
	// Verify checksums, correct errors where possible.
	checks := []struct {
		data []byte
		chk  uint32
	}{
		{pubKeyHash[:], chkPubKeyHash},
		{a.chaincode[:], chkChaincode},
		{a.initVector[:], chkInitVector},
		{a.privKey[:], chkPrivKey},
		{pubKey, chkPubKey},
	}
	for i := range checks {
		if e = verifyAndFix(checks[i].data, checks[i].chk); E.Chk(e) {
			return n, e
		}
	}
	if !a.flags.hasPubKey {
		return n, errors.New("read in an address without a public key")
	}
	pk, e := ec.ParsePubKey(pubKey, ec.S256())
	if e != nil {
		return n, e
	}
	a.pubKey = pk
	addr, e := btcaddr.NewPubKeyHash(pubKeyHash[:], a.store.netParams())
	if e != nil {
		return n, e
	}
	a.address = addr
	return n, nil
}
func (a *btcAddress) WriteTo(w io.Writer) (n int64, e error) {
	var written int64
	pubKey := a.pubKeyBytes()
	hash := a.address.ScriptAddress()
	datas := []interface{}{
		&hash,
		walletHash(hash),
		make([]byte, 4), // ksVersion
		&a.flags,
		&a.chaincode,
		walletHash(a.chaincode[:]),
		&a.chainIndex,
		&a.chainDepth,
		&a.initVector,
		walletHash(a.initVector[:]),
		&a.privKey,
		walletHash(a.privKey[:]),
		&pubKey,
		walletHash(pubKey),
		&a.firstSeen,
		&a.lastSeen,
		&a.firstBlock,
		&a.partialSyncHeight,
	}
	for _, data := range datas {
		if wt, ok := data.(io.WriterTo); ok {
			written, e = wt.WriteTo(w)
		} else {
			written, e = binaryWrite(w, binary.LittleEndian, data)
		}
		if e != nil {
			return n + written, e
		}
		n += written
	}
	return n, nil
}

// encrypt attempts to encrypt an address's clear text private key, failing if the address is already encrypted or if
// the private key is not 32 bytes. If successful, the encryption flag is set.
func (a *btcAddress) encrypt(key []byte) (e error) {
	if a.flags.encrypted {
		return ErrAlreadyEncrypted
	}
	if len(a.privKeyCT) != 32 {
		return errors.New("invalid clear text private key")
	}
	aesBlockEncrypter, e := aes.NewCipher(key)
	if e != nil {
		return e
	}
	aesEncrypter := cipher.NewCFBEncrypter(aesBlockEncrypter, a.initVector[:])
	aesEncrypter.XORKeyStream(a.privKey[:], a.privKeyCT)
	a.flags.hasPrivKey = true
	a.flags.encrypted = true
	return nil
}

// lock removes the reference this address holds to its clear text private key. This function fails if the address is
// not encrypted.
func (a *btcAddress) lock() (e error) {
	if !a.flags.encrypted {
		return errors.New("unable to lock unencrypted address")
	}
	zero(a.privKeyCT)
	a.privKeyCT = nil
	return nil
}

// unlock decrypts and stores a pointer to an address's private key, failing if the address is not encrypted, or the
// provided key is incorrect.
//
// The returned clear text private key will always be a copy that may be safely used by the caller without worrying
// about it being zeroed during an address lock.
func (a *btcAddress) unlock(key []byte) (privKeyCT []byte, e error) {
	if !a.flags.encrypted {
		return nil, errors.New("unable to unlock unencrypted address")
	}
	// Decrypt private key with AES key.
	aesBlockDecryptor, e := aes.NewCipher(key)
	if e != nil {
		return nil, e
	}
	aesDecryptor := cipher.NewCFBDecrypter(aesBlockDecryptor, a.initVector[:])
	privkey := make([]byte, 32)
	aesDecryptor.XORKeyStream(privkey, a.privKey[:])
	// If secret is already saved, simply compare the bytes.
	if len(a.privKeyCT) == 32 {
		if !bytes.Equal(a.privKeyCT, privkey) {
			return nil, ErrWrongPassphrase
		}
		privKeyCT := make([]byte, 32)
		copy(privKeyCT, a.privKeyCT)
		return privKeyCT, nil
	}
	x, y := ec.S256().ScalarBaseMult(privkey)
	if x.Cmp(a.pubKey.X) != 0 || y.Cmp(a.pubKey.Y) != 0 {
		return nil, ErrWrongPassphrase
	}
	privkeyCopy := make([]byte, 32)
	copy(privkeyCopy, privkey)
	a.privKeyCT = privkey
	return privkeyCopy, nil
}

// changeEncryptionKey re-encrypts the private keys for an address with a new AES encryption key. oldkey must be the old
// AES encryption key and is used to decrypt the private key.
func (a *btcAddress) changeEncryptionKey(oldkey, newkey []byte) (e error) {
	// Address must have a private key and be encrypted to continue.
	if !a.flags.hasPrivKey {
		return errors.New("no private key")
	}
	if !a.flags.encrypted {
		return errors.New("address is not encrypted")
	}
	privKeyCT, e := a.unlock(oldkey)
	if e != nil {
		return e
	}
	aesBlockEncrypter, e := aes.NewCipher(newkey)
	if e != nil {
		return e
	}
	newIV := make([]byte, len(a.initVector))
	if _, e = rand.Read(newIV); E.Chk(e) {
		return e
	}
	copy(a.initVector[:], newIV)
	aesEncrypter := cipher.NewCFBEncrypter(aesBlockEncrypter, a.initVector[:])
	aesEncrypter.XORKeyStream(a.privKey[:], privKeyCT)
	return nil
}

// Address returns the pub key address, implementing AddressInfo.
func (a *btcAddress) Address() btcaddr.Address {
	return a.address
}

// AddrHash returns the pub key hash, implementing WalletAddress.
func (a *btcAddress) AddrHash() string {
	return string(a.address.ScriptAddress())
}

// FirstBlock returns the first block the address is seen in, implementing AddressInfo.
func (a *btcAddress) FirstBlock() int32 {
	return a.firstBlock
}

// Imported returns the pub if the address was imported, or a chained address, implementing AddressInfo.
func (a *btcAddress) Imported() bool {
	return a.chainIndex == importedKeyChainIdx
}

// Change returns true if the address was created as a change address, implementing AddressInfo.
func (a *btcAddress) Change() bool {
	return a.flags.change
}

// Compressed returns true if the address backing key is compressed, implementing AddressInfo.
func (a *btcAddress) Compressed() bool {
	return a.flags.compressed
}

// SyncStatus returns a SyncStatus type for how the address is currently synced. For an Unsynced type, the value is the
// recorded first seen block height of the address.
func (a *btcAddress) SyncStatus() SyncStatus {
	switch {
	case a.flags.unsynced && !a.flags.partialSync:
		return Unsynced(a.firstBlock)
	case a.flags.unsynced && a.flags.partialSync:
		return PartialSync(a.partialSyncHeight)
	default:
		return FullSync{}
	}
}

// PubKey returns the hex encoded pubkey for the address. Implementing PubKeyAddress.
func (a *btcAddress) PubKey() *ec.PublicKey {
	return a.pubKey
}
func (a *btcAddress) pubKeyBytes() []byte {
	if a.Compressed() {
		return a.pubKey.SerializeCompressed()
	}
	return a.pubKey.SerializeUncompressed()
}

// ExportPubKey returns the public key associated with the address serialised as a hex encoded string. Implements
// PubKeyAddress
func (a *btcAddress) ExportPubKey() string {
	return hex.EncodeToString(a.pubKeyBytes())
}

// PrivKey implements PubKeyAddress by returning the private key, or an error if the key store is locked, watching only
// or the private key is missing.
func (a *btcAddress) PrivKey() (*ec.PrivateKey, error) {
	if a.store.flags.watchingOnly {
		return nil, ErrWatchingOnly
	}
	if !a.flags.hasPrivKey {
		return nil, errors.New("no private key for address")
	}
	// Key store must be unlocked to decrypt the private key.
	if a.store.isLocked() {
		return nil, ErrLocked
	}
	// Unlock address with key store secret. unlock returns a copy of the clear text private key, and may be used safely
	// even during an address lock.
	privKeyCT, e := a.unlock(a.store.secret)
	if e != nil {
		return nil, e
	}
	return &ec.PrivateKey{
		PublicKey: *a.pubKey.ToECDSA(),
		D:         new(big.Int).SetBytes(privKeyCT),
	}, nil
}

// ExportPrivKey exports the private key as a WIF for encoding as a string in the Wallet Import Format.
func (a *btcAddress) ExportPrivKey() (*util.WIF, error) {
	pk, e := a.PrivKey()
	if e != nil {
		return nil, e
	}
	// NewWIF only errors if the network is nil. In this case, panic, as our program's assumptions are so broken that
	// this needs to be caught immediately, and a stack trace here is more useful than elsewhere.
	wif, e := util.NewWIF(
		pk, a.store.netParams(),
		a.Compressed(),
	)
	if e != nil {
		panic(e)
	}
	return wif, nil
}

// watchingCopy creates a copy of an address without a private key. This is used to fill a watching a key store with
// addresses from a normal key store.
func (a *btcAddress) watchingCopy(s *Store) walletAddress {
	return &btcAddress{
		store:   s,
		address: a.address,
		flags: addrFlags{
			hasPrivKey:              false,
			hasPubKey:               true,
			encrypted:               false,
			createPrivKeyNextUnlock: false,
			compressed:              a.flags.compressed,
			change:                  a.flags.change,
			unsynced:                a.flags.unsynced,
		},
		chaincode:         a.chaincode,
		chainIndex:        a.chainIndex,
		chainDepth:        a.chainDepth,
		pubKey:            a.pubKey,
		firstSeen:         a.firstSeen,
		lastSeen:          a.lastSeen,
		firstBlock:        a.firstBlock,
		partialSyncHeight: a.partialSyncHeight,
	}
}

// setSyncStatus sets the address flags and possibly the partial sync height depending on the type of s.
func (a *btcAddress) setSyncStatus(s SyncStatus) {
	switch e := s.(type) {
	case Unsynced:
		a.flags.unsynced = true
		a.flags.partialSync = false
		a.partialSyncHeight = 0
	case PartialSync:
		a.flags.unsynced = true
		a.flags.partialSync = true
		a.partialSyncHeight = int32(e)
	case FullSync:
		a.flags.unsynced = false
		a.flags.partialSync = false
		a.partialSyncHeight = 0
	}
}

// note that there is no encrypted bit here since if we had a script encrypted and then used it on the blockchain this
// provides a simple known plaintext in the key store file.
//
// It was determined that the script in a p2sh transaction is not a secret and any sane situation would also require a
// signature (which does have a secret).
type scriptFlags struct {
	hasScript   bool
	change      bool
	unsynced    bool
	partialSync bool
}

// ReadFrom implements the io.ReaderFrom interface by reading from r into sf.
func (sf *scriptFlags) ReadFrom(r io.Reader) (int64, error) {
	var b [8]byte
	n, e := io.ReadFull(r, b[:])
	if e != nil {
		return int64(n), e
	}
	// We match bits from addrFlags for similar fields. hence hasScript uses the same bit as hasPubKey and the change
	// bit is the same for both.
	sf.hasScript = b[0]&(1<<1) != 0
	sf.change = b[0]&(1<<5) != 0
	sf.unsynced = b[0]&(1<<6) != 0
	sf.partialSync = b[0]&(1<<7) != 0
	return int64(n), nil
}

// WriteTo implements the io.WriteTo interface by writing sf into w.
func (sf *scriptFlags) WriteTo(w io.Writer) (int64, error) {
	var b [8]byte
	if sf.hasScript {
		b[0] |= 1 << 1
	}
	if sf.change {
		b[0] |= 1 << 5
	}
	if sf.unsynced {
		b[0] |= 1 << 6
	}
	if sf.partialSync {
		b[0] |= 1 << 7
	}
	n, e := w.Write(b[:])
	return int64(n), e
}

// p2SHScript represents the variable length script entry in a key store.
type p2SHScript []byte

// ReadFrom implements the ReaderFrom interface by reading the P2SH script from r in the format <4 bytes little endian
// length><script bytes>
func (a *p2SHScript) ReadFrom(r io.Reader) (n int64, e error) {
	// read length
	var lenBytes [4]byte
	read, e := io.ReadFull(r, lenBytes[:])
	n += int64(read)
	if e != nil {
		return n, e
	}
	length := binary.LittleEndian.Uint32(lenBytes[:])
	script := make([]byte, length)
	read, e = io.ReadFull(r, script)
	n += int64(read)
	if e != nil {
		return n, e
	}
	*a = script
	return n, nil
}

// WriteTo implements the WriterTo interface by writing the P2SH script to w in the format <4 bytes little endian
// length><script bytes>
func (a *p2SHScript) WriteTo(w io.Writer) (n int64, e error) {
	// Prepare and write 32-bit little-endian length header
	var lenBytes [4]byte
	binary.LittleEndian.PutUint32(lenBytes[:], uint32(len(*a)))
	written, e := w.Write(lenBytes[:])
	n += int64(written)
	if e != nil {
		return n, e
	}
	// Now write the bytes themselves.
	written, e = w.Write(*a)
	return n + int64(written), e
}

type scriptAddress struct {
	store             *Store
	address           btcaddr.Address
	class             txscript.ScriptClass
	addresses         []btcaddr.Address
	reqSigs           int
	flags             scriptFlags
	script            p2SHScript // variable length
	firstSeen         int64
	lastSeen          int64
	firstBlock        int32
	partialSyncHeight int32
}

// ScriptAddress is an interface representing a Pay-to-Script-Hash style of bitcoind address.
type ScriptAddress interface {
	WalletAddress
	// Script Returns the script associated with the address.
	Script() []byte
	// ScriptClass Returns the class of the script associated with the address.
	ScriptClass() txscript.ScriptClass
	// Addresses Returns the addresses that are required to sign transactions from the script address.
	Addresses() []btcaddr.Address
	// RequiredSigs Returns the number of signatures required by the script address.
	RequiredSigs() int
}

// newScriptAddress initializes and returns a new P2SH address. iv must be 16 bytes, or nil (in which case it is
// randomly generated).
func newScriptAddress(s *Store, script []byte, bs *BlockStamp) (addr *scriptAddress, e error) {
	class, addresses, reqSigs, e := txscript.ExtractPkScriptAddrs(script, s.netParams())
	if e != nil {
		return nil, e
	}
	scriptHash := btcaddr.Hash160(script)
	var a *btcaddr.ScriptHash
	a, e = btcaddr.NewScriptHashFromHash(scriptHash, s.netParams())
	if e != nil {
		return nil, e
	}
	addr = &scriptAddress{
		store:     s,
		address:   a,
		addresses: addresses,
		class:     class,
		reqSigs:   reqSigs,
		flags: scriptFlags{
			hasScript: true,
			change:    false,
		},
		script:     script,
		firstSeen:  time.Now().Unix(),
		firstBlock: bs.Height,
	}
	return addr, nil
}

// ReadFrom reads an script address from an io.Reader.
func (sa *scriptAddress) ReadFrom(r io.Reader) (n int64, e error) {
	var read int64
	// Checksums
	var chkScriptHash uint32
	var chkScript uint32
	var scriptHash [ripemd160.Size]byte
	// Read serialized key store into addr fields and checksums.
	datas := []interface{}{
		&scriptHash,
		&chkScriptHash,
		make([]byte, 4), // ksVersion
		&sa.flags,
		&sa.script,
		&chkScript,
		&sa.firstSeen,
		&sa.lastSeen,
		&sa.firstBlock,
		&sa.partialSyncHeight,
	}
	for _, data := range datas {
		if rf, ok := data.(io.ReaderFrom); ok {
			read, e = rf.ReadFrom(r)
		} else {
			read, e = binaryRead(r, binary.LittleEndian, data)
		}
		if e != nil {
			return n + read, e
		}
		n += read
	}
	// Verify checksums, correct errors where possible.
	checks := []struct {
		data []byte
		chk  uint32
	}{
		{scriptHash[:], chkScriptHash},
		{sa.script, chkScript},
	}
	for i := range checks {
		if e = verifyAndFix(checks[i].data, checks[i].chk); E.Chk(e) {
			return n, e
		}
	}
	address, e := btcaddr.NewScriptHashFromHash(
		scriptHash[:],
		sa.store.netParams(),
	)
	if e != nil {
		return n, e
	}
	sa.address = address
	if !sa.flags.hasScript {
		return n, errors.New("read in an addresss with no script")
	}
	class, addresses, reqSigs, err :=
		txscript.ExtractPkScriptAddrs(sa.script, sa.store.netParams())
	if e != nil {
		return n, err
	}
	sa.class = class
	sa.addresses = addresses
	sa.reqSigs = reqSigs
	return n, nil
}

// WriteTo implements io.WriterTo by writing the scriptAddress to w.
func (sa *scriptAddress) WriteTo(w io.Writer) (n int64, e error) {
	var written int64
	hash := sa.address.ScriptAddress()
	datas := []interface{}{
		&hash,
		walletHash(hash),
		make([]byte, 4), // ksVersion
		&sa.flags,
		&sa.script,
		walletHash(sa.script),
		&sa.firstSeen,
		&sa.lastSeen,
		&sa.firstBlock,
		&sa.partialSyncHeight,
	}
	for _, data := range datas {
		if wt, ok := data.(io.WriterTo); ok {
			written, e = wt.WriteTo(w)
		} else {
			written, e = binaryWrite(w, binary.LittleEndian, data)
		}
		if e != nil {
			return n + written, e
		}
		n += written
	}
	return n, nil
}

// Address returns a util.ScriptHash for a btcAddress.
func (sa *scriptAddress) Address() btcaddr.Address {
	return sa.address
}

// AddrHash returns the script hash, implementing AddressInfo.
func (sa *scriptAddress) AddrHash() string {
	return string(sa.address.ScriptAddress())
}

// FirstBlock returns the first blockheight the address is known at.
func (sa *scriptAddress) FirstBlock() int32 {
	return sa.firstBlock
}

// Imported currently always returns true since script addresses are always imported addressed and not part of any
// chain.
func (sa *scriptAddress) Imported() bool {
	return true
}

// Change returns true if the address was created as a change address.
func (sa *scriptAddress) Change() bool {
	return sa.flags.change
}

// Compressed returns false since script addresses are never compressed. Implements WalletAddress.
func (sa *scriptAddress) Compressed() bool {
	return false
}

// Script returns the script that is represented by the address. It should not be modified.
func (sa *scriptAddress) Script() []byte {
	return sa.script
}

// Addresses returns the list of addresses that must sign the script.
func (sa *scriptAddress) Addresses() []btcaddr.Address {
	return sa.addresses
}

// ScriptClass returns the type of script the address is.
func (sa *scriptAddress) ScriptClass() txscript.ScriptClass {
	return sa.class
}

// RequiredSigs returns the number of signatures required by the script.
func (sa *scriptAddress) RequiredSigs() int {
	return sa.reqSigs
}

// SyncStatus returns a SyncStatus type for how the address is currently synced. For an Unsynced type, the value is the
// recorded first seen block height of the address.
//
// Implements WalletAddress.
func (sa *scriptAddress) SyncStatus() SyncStatus {
	switch {
	case sa.flags.unsynced && !sa.flags.partialSync:
		return Unsynced(sa.firstBlock)
	case sa.flags.unsynced && sa.flags.partialSync:
		return PartialSync(sa.partialSyncHeight)
	default:
		return FullSync{}
	}
}

// setSyncStatus sets the address flags and possibly the partial sync height depending on the type of s.
func (sa *scriptAddress) setSyncStatus(s SyncStatus) {
	switch e := s.(type) {
	case Unsynced:
		sa.flags.unsynced = true
		sa.flags.partialSync = false
		sa.partialSyncHeight = 0
	case PartialSync:
		sa.flags.unsynced = true
		sa.flags.partialSync = true
		sa.partialSyncHeight = int32(e)
	case FullSync:
		sa.flags.unsynced = false
		sa.flags.partialSync = false
		sa.partialSyncHeight = 0
	}
}

// watchingCopy creates a copy of an address without a private key.
//
// This is used to fill a watching key store with addresses from a normal key store.
func (sa *scriptAddress) watchingCopy(s *Store) walletAddress {
	return &scriptAddress{
		store:     s,
		address:   sa.address,
		addresses: sa.addresses,
		class:     sa.class,
		reqSigs:   sa.reqSigs,
		flags: scriptFlags{
			change:   sa.flags.change,
			unsynced: sa.flags.unsynced,
		},
		script:            sa.script,
		firstSeen:         sa.firstSeen,
		lastSeen:          sa.lastSeen,
		firstBlock:        sa.firstBlock,
		partialSyncHeight: sa.partialSyncHeight,
	}
}
func walletHash(b []byte) uint32 {
	sum := chainhash.DoubleHashB(b)
	return binary.LittleEndian.Uint32(sum)
}

// TODO(jrick) add error correction.
func verifyAndFix(b []byte, chk uint32) (e error) {
	if walletHash(b) != chk {
		return ErrChecksumMismatch
	}
	return nil
}

type kdfParameters struct {
	mem   uint64
	nIter uint32
	salt  [32]byte
}

// computeKdfParameters returns best guess parameters to the memory-hard key derivation function to make the computation
// last targetSec seconds, while using no more than maxMem bytes of memory.
func computeKdfParameters(targetSec float64, maxMem uint64) (params *kdfParameters, e error) {
	params = &kdfParameters{}
	if _, e = rand.Read(params.salt[:]); E.Chk(e) {
		return nil, e
	}
	testKey := []byte("This is an example key to test KDF iteration speed")
	memoryReqtBytes := uint64(1024)
	approxSec := float64(0)
	for approxSec <= targetSec/4 && memoryReqtBytes < maxMem {
		memoryReqtBytes *= 2
		before := time.Now()
		_ = keyOneIter(testKey, params.salt[:], memoryReqtBytes)
		approxSec = time.Since(before).Seconds()
	}
	allItersSec := float64(0)
	nIter := uint32(1)
	for allItersSec < 0.02 { // This is a magic number straight from armory's source.
		nIter *= 2
		before := time.Now()
		for i := uint32(0); i < nIter; i++ {
			_ = keyOneIter(testKey, params.salt[:], memoryReqtBytes)
		}
		allItersSec = time.Since(before).Seconds()
	}
	params.mem = memoryReqtBytes
	params.nIter = nIter
	return params, nil
}
func (params *kdfParameters) WriteTo(w io.Writer) (n int64, e error) {
	var written int64
	memBytes := make([]byte, 8)
	nIterBytes := make([]byte, 4)
	binary.LittleEndian.PutUint64(memBytes, params.mem)
	binary.LittleEndian.PutUint32(nIterBytes, params.nIter)
	chkedBytes := append(memBytes, nIterBytes...)
	chkedBytes = append(chkedBytes, params.salt[:]...)
	datas := []interface{}{
		&params.mem,
		&params.nIter,
		&params.salt,
		walletHash(chkedBytes),
		make([]byte, 256-(binary.Size(params)+4)), // padding
	}
	for _, data := range datas {
		if written, e = binaryWrite(w, binary.LittleEndian, data); E.Chk(e) {
			return n + written, e
		}
		n += written
	}
	return n, nil
}
func (params *kdfParameters) ReadFrom(r io.Reader) (n int64, e error) {
	var read int64
	// These must be read in but are not saved directly to chaincfg.
	chkedBytes := make([]byte, 44)
	var chk uint32
	padding := make([]byte, 256-(binary.Size(params)+4))
	datas := []interface{}{
		chkedBytes,
		&chk,
		padding,
	}
	for _, data := range datas {
		if read, e = binaryRead(r, binary.LittleEndian, data); E.Chk(e) {
			return n + read, e
		}
		n += read
	}
	// Verify checksum
	if e = verifyAndFix(chkedBytes, chk); E.Chk(e) {
		return n, e
	}
	// Read netparams
	buf := bytes.NewBuffer(chkedBytes)
	datas = []interface{}{
		&params.mem,
		&params.nIter,
		&params.salt,
	}
	for _, data := range datas {
		if e = binary.Read(buf, binary.LittleEndian, data); E.Chk(e) {
			return n, e
		}
	}
	return n, nil
}

type addrEntry struct {
	pubKeyHash160 [ripemd160.Size]byte
	addr          btcAddress
}

func (ae *addrEntry) WriteTo(w io.Writer) (n int64, e error) {
	var written int64
	// Write header
	if written, e = binaryWrite(w, binary.LittleEndian, addrHeader); E.Chk(e) {
		return n + written, e
	}
	n += written
	// Write hash
	if written, e = binaryWrite(w, binary.LittleEndian, &ae.pubKeyHash160); E.Chk(e) {
		return n + written, e
	}
	n += written
	// Write btcAddress
	written, e = ae.addr.WriteTo(w)
	n += written
	return n, e
}
func (ae *addrEntry) ReadFrom(r io.Reader) (n int64, e error) {
	var read int64
	if read, e = binaryRead(r, binary.LittleEndian, &ae.pubKeyHash160); E.Chk(e) {
		return n + read, e
	}
	n += read
	read, e = ae.addr.ReadFrom(r)
	return n + read, e
}

// scriptEntry is the entry type for a P2SH script.
type scriptEntry struct {
	scriptHash160 [ripemd160.Size]byte
	script        scriptAddress
}

// WriteTo implements io.WriterTo by writing the entry to w.
func (se *scriptEntry) WriteTo(w io.Writer) (n int64, e error) {
	var written int64
	// Write header
	if written, e = binaryWrite(w, binary.LittleEndian, scriptHeader); E.Chk(e) {
		return n + written, e
	}
	n += written
	// Write hash
	if written, e = binaryWrite(w, binary.LittleEndian, &se.scriptHash160); E.Chk(e) {
		return n + written, e
	}
	n += written
	// Write btcAddress
	written, e = se.script.WriteTo(w)
	n += written
	return n, e
}

// ReadFrom implements io.ReaderFrom by reading the entry from e.
func (se *scriptEntry) ReadFrom(r io.Reader) (n int64, e error) {
	var read int64
	if read, e = binaryRead(r, binary.LittleEndian, &se.scriptHash160); E.Chk(e) {
		return n + read, e
	}
	n += read
	read, e = se.script.ReadFrom(r)
	return n + read, e
}

// BlockStamp defines a block (by height and a unique hash) and is used to mark a point in the blockchain that a key
// store element is synced to.
type BlockStamp struct {
	Hash   *chainhash.Hash
	Height int32
}
