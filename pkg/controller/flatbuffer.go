package controller

import (
	"bytes"
	"encoding/binary"
	"fmt"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/routeable"
	"net"
	"strconv"
	"strings"
	"sync"
)

type Serializer interface {
	// Encode returns the wire/storage form of the data
	Encode() []byte
	// Decode stores the decoded data from the head of the slice and returns
	// the remainder
	Decode(b []byte) []byte
}

type Serializers []Serializer

type Container struct {
	Data []byte
}

// CreateContainer takes an array of serializer interface objects and renders
// the data into bytes
func (srs Serializers) CreateContainer(magic [4]byte) (out *Container) {
	out = &Container{}
	var offset uint32
	var length uint16
	var nodes []uint32
	for i := range srs {
		b := srs[i].Encode()
		//log.DEBUG(i, len(b), hex.EncodeToString(b))
		length++
		nodes = append(nodes, offset)
		offset += uint32(len(b))
		out.Data = append(out.Data, b...)
	}
	//log.SPEW(out.Data)
	//log.DEBUG(offset, len(out.Data))
	nodeB := make([]byte, len(nodes)*4+2)
	start := uint32(len(nodeB) + 8)
	binary.BigEndian.PutUint16(nodeB[:2], length)
	for i := range nodes {
		b := nodeB[i*4+2 : i*4+4+2]
		binary.BigEndian.PutUint32(b, nodes[i]+start)
		//log.DEBUG(i, len(b), hex.EncodeToString(b))
	}
	//log.SPEW(nodeB)
	out.Data = append(nodeB, out.Data...)
	size := offset + uint32(len(nodeB)) + 8
	//log.DEBUG("size", size, len(out.Data))
	sB := make([]byte, 4)
	binary.BigEndian.PutUint32(sB, size)
	out.Data = append(append(magic[:], sB...), out.Data...)
	return
}

func (c *Container) Count() uint16 {
	size := binary.BigEndian.Uint32(c.Data[4:8])
	//log.DEBUG("size", size)
	if len(c.Data) >= int(size) {
		// we won't touch it if it's not at least as big so we don't get
		// bounds errors
		return binary.BigEndian.Uint16(c.Data[8:10])
	}
	return 0
}

func (c *Container) GetMagic() (out [4]byte) {
	copy(out[:], c.Data[:4])
	return
}

// Get returns the bytes that can be imported into an interface assuming the
// types are correct - field ordering is hard coded by the creation and
// identified by the magic. This is all read only and subslices so it should
// generate very little garbage or copy operations except as required for the
// output (we aren't going to go unsafe here,
// it isn't really necessary since already this library enables avoiding the
// decoding of values not being used from a message (or not used yet)
func (c *Container) Get(idx uint16) (out []byte) {
	length := c.Count()
	size := len(c.Data)
	if length > idx {
		//log.DEBUG("length", length, "idx", idx)
		if idx < length {
			offset := binary.BigEndian.Uint32(c.
				Data[10+idx*4 : 10+idx*4+4])
			//log.DEBUG("offset", offset)
			if idx < length-1 {
				nextOffset := binary.BigEndian.Uint32(c.
					Data[10+((idx+1)*4) : 10+((idx+1)*4)+4])
				//log.DEBUG("nextOffset", nextOffset)
				out = c.Data[offset:nextOffset]
			} else {
				nextOffset := len(c.Data)
				//log.DEBUG("last nextOffset", nextOffset)
				out = c.Data[offset:nextOffset]
			}
		}
	} else {
		panic(fmt.Sprintln("size mismatch", len(c.Data), size))
	}
	return
}

type Port struct {
	Bytes [2]byte
}

func NewPort() *Port {
	return &Port{}
}

func (p *Port) DecodeOne(b []byte) *Port {
	p.Decode(b)
	return p
}

func (p *Port) Decode(b []byte) (out []byte) {
	if len(b) >= 2 {
		p.Bytes = [2]byte{b[0], b[1]}
		if len(b) > 2 {
			out = b[2:]
		}
	}
	return
}

func (p *Port) Encode() []byte {
	return p.Bytes[:]
}

func (p *Port) Get() uint16 {
	return binary.BigEndian.Uint16(p.Bytes[:2])
}

func (p *Port) Put(i uint16) *Port {
	binary.BigEndian.PutUint16(p.Bytes[:], i)
	return p
}

type IP struct {
	Length byte
	Bytes  []byte
}

func NewIP() *IP {
	return &IP{}
}

func (i *IP) DecodeOne(b []byte) *IP {
	i.Decode(b)
	return i
}

func (i *IP) Decode(b []byte) (out []byte) {
	if len(b) >= 1 {
		i.Length = b[0]
		if len(b) > int(i.Length) {
			i.Bytes = b[1 : i.Length+1]
		}
	}
	total := int(i.Length) + 1
	if len(b) > total {
		out = b[total:]
	}
	return
}

func (i *IP) Encode() []byte {
	return append([]byte{i.Length}, i.Bytes...)
}

func (i *IP) Get() *net.IP {
	ip := make(net.IP, len(i.Bytes))
	copy(ip, i.Bytes)
	return &ip
}

func (i *IP) Put(ip *net.IP) *IP {
	i.Bytes = make([]byte, len(*ip))
	copy(i.Bytes, *ip)
	i.Length = byte(len(i.Bytes))
	return i
}

type IPs struct {
	Length byte
	IPs    []IP
}

func NewIPs() *IPs {
	return &IPs{}
}

func (ips *IPs) DecodeOne(b []byte) *IPs {
	ips.Decode(b)
	return ips
}

func (ips *IPs) Decode(b []byte) (out []byte) {
	if len(b) >= 1 {
		ips.Length = b[0]
		out = b[1:]
		count := ips.Length
		for ; count > 0; count-- {
			i := &IP{}
			out = i.Decode(out)
			ipa := make(net.IP, 16)
			copy(ipa, i.Bytes)
			nIP := NewIP()
			nIP.Decode(i.Encode())
			ips.IPs = append(ips.IPs, *nIP)
		}
	}
	return
}

func (ips *IPs) Encode() (out []byte) {
	out = []byte{ips.Length}
	for i := range ips.IPs {
		b := ips.IPs[i].Bytes
		out = append(out, append([]byte{byte(len(b))}, b...)...)
	}
	return
}

func (ips *IPs) Put(in []*net.IP) *IPs {
	ips.Length = byte(len(in))
	ips.IPs = make([]IP, len(in))
	for i := range in {
		ips.IPs[i].Put(in[i])
	}
	return ips
}

func (ips *IPs) Get() (out []*net.IP) {
	for i := range ips.IPs {
		out = append(out, ips.IPs[i].Get())
	}
	return
}

type Bits struct {
	Bytes [4]byte
}

func NewBits() *Bits {
	return &Bits{}
}

func (b *Bits) DecodeOne(by []byte) *Bits {
	b.Decode(by)
	return b
}

func (b *Bits) Decode(by []byte) (out []byte) {
	if len(by) >= 4 {
		b.Bytes = [4]byte{by[0], by[1], by[2], by[3]}
		if len(by) > 4 {
			out = by[4:]
		}
	}
	return
}

func (b *Bits) Encode() []byte {
	return b.Bytes[:]
}

func (b *Bits) Get() uint32 {
	return binary.BigEndian.Uint32(b.Bytes[:])
}

func (b *Bits) Put(bits uint32) *Bits {
	binary.BigEndian.PutUint32(b.Bytes[:], bits)
	return b
}

type Bitses struct {
	sync.Mutex
	Length  byte
	Byteses map[int32][]byte
}

func NewBitses() *Bitses {
	return &Bitses{Byteses: make(map[int32][]byte)}
}

func (b *Bitses) DecodeOne(by []byte) *Bitses {
	b.Decode(by)
	return b
}

func (b *Bitses) Decode(by []byte) (out []byte) {
	b.Lock()
	defer b.Unlock()
	//log.SPEW(by)
	if len(by) >= 7 {
		nB := by[0]
		if len(by) >= int(nB)*8 {
			for i := 0; i < int(nB); i++ {
				algoVer := int32(binary.BigEndian.Uint32(by[1+i*8 : 1+i*8+4]))
				//log.DEBUG("algoVer", algoVer, by[1+i*8+4:1+i*8+8], b.Byteses)
				b.Byteses[algoVer] = by[1+i*8+4 : 1+i*8+8]
			}
		}
		bL := int(nB)*8 + 1
		if len(by) > bL {
			out = by[bL:]
		}
	}
	//log.SPEW(b.Byteses)
	return
}

func (b *Bitses) Encode() (out []byte) {
	b.Lock()
	defer b.Unlock()
	out = []byte{b.Length}
	for algoVer := range b.Byteses {
		by := make([]byte, 4)
		binary.BigEndian.PutUint32(by, uint32(algoVer))
		out = append(out, append(by, b.Byteses[algoVer]...)...)
	}
	//log.SPEW(out)
	return
}

func (b *Bitses) Get() (out map[int32]uint32) {
	b.Lock()
	defer b.Unlock()
	out = make(map[int32]uint32)
	for algoVer := range b.Byteses {
		oB := binary.BigEndian.Uint32(b.Byteses[algoVer])
		out[algoVer] = oB
	}
	//log.SPEW(out)
	return
}

func (b *Bitses) Put(in map[int32]uint32) *Bitses {
	b.Lock()
	defer b.Unlock()
	b.Length = byte(len(in))
	b.Byteses = make(map[int32][]byte, b.Length)
	for algoVer := range in {
		bits := make([]byte, 4)
		binary.BigEndian.PutUint32(bits, in[algoVer])
		b.Byteses[algoVer] = bits
	}
	//log.SPEW(b.Byteses)
	return b
}

type Hash struct {
	Hash *chainhash.Hash
}

func NewHash() *Hash {
	return &Hash{Hash: new(chainhash.Hash)}
}

func (h *Hash) DecodeOne(b []byte) *Hash {
	h.Decode(b)
	return h
}

func (h *Hash) Decode(b []byte) (out []byte) {
	if len(b) >= 32 {
		err := h.Hash.SetBytes(b[:32])
		if err != nil {
			log.ERROR(err)
			return
		}
		if len(b) > 32 {
			out = b[32:]
		}
	}
	return
}

func (h *Hash) Encode() []byte {
	return h.Hash.CloneBytes()
}

func (h *Hash) Get() *chainhash.Hash {
	return h.Hash
}

func (h *Hash) Put(pH chainhash.Hash) *Hash {
	// this should avoid a copy
	h.Hash = &pH
	return h
}

type Transaction struct {
	Length uint32
	Bytes  []byte
}

func NewTransaction() *Transaction {
	return &Transaction{}
}

func (t *Transaction) DecodeOne(b []byte) *Transaction {
	t.Decode(b)
	return t
}

func (t *Transaction) Decode(b []byte) (out []byte) {
	if len(b) >= 4 {
		t.Length = binary.BigEndian.Uint32(b[:4])
		if len(b) >= 4+int(t.Length) {
			t.Bytes = b[4 : 4+t.Length]
			if len(b) > 4+int(t.Length) {
				out = b[4+t.Length:]
			}
		}
	}
	return
}

func (t *Transaction) Encode() (out []byte) {
	out = make([]byte, 4+len(t.Bytes))
	binary.BigEndian.PutUint32(out[:4], t.Length)
	copy(out[4:], t.Bytes)
	return
}

func (t *Transaction) Get() (txs *wire.MsgTx) {
	txs = new(wire.MsgTx)
	buffer := bytes.NewBuffer(t.Bytes)
	err := txs.Deserialize(buffer)
	if err != nil {
		log.ERROR(err)
	}
	return
}

func (t *Transaction) Put(txs *wire.MsgTx) *Transaction {
	var buffer bytes.Buffer
	err := txs.Serialize(&buffer)
	if err != nil {
		log.ERROR(err)
		return t
	}
	t.Bytes = buffer.Bytes()
	t.Length = uint32(len(t.Bytes))
	return t
}

func GetRouteableIPs() Serializer {
	// first add the interface addresses
	rI := routeable.GetInterface()
	//log.SPEW(rI)
	var lA []net.Addr
	for i := range rI {
		l, err := rI[i].Addrs()
		//log.SPEW(lA)
		if err != nil {
			log.ERROR(err)
			return nil
		}
		lA = append(lA, l...)
	}
	ips := NewIPs()
	var ipslice []*net.IP
	for i := range lA {
		//log.DEBUG(lA[i])
		addIP := net.ParseIP(strings.Split(lA[i].String(), "/")[0])
		ipslice = append(ipslice, &addIP)
	}
	ips.Put(ipslice)
	//log.SPEW(ipslice)
	//log.SPEW(ips)
	//log.SPEW(ips.Get())
	return ips
}

func GetPort(listener string) Serializer {
	//log.DEBUG(listener)
	_, p, err := net.SplitHostPort(listener)
	if err != nil {
		log.ERROR(err)
		return nil
	}
	oI, err := strconv.ParseInt(p, 10, 32)
	if err != nil {
		log.ERROR(err)
		return nil
	}
	port := &Port{}
	port.Put(uint16(oI))
	return port
}
