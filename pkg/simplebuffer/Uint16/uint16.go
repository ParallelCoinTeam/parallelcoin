package Uint16

import (
	"encoding/binary"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/simplebuffer"
	"net"
	"strconv"
)

type Uint16 struct {
	Bytes [2]byte
}

func New() *Uint16 {
	return &Uint16{}
}

func (p *Uint16) DecodeOne(b []byte) *Uint16 {
	p.Decode(b)
	return p
}

func (p *Uint16) Decode(b []byte) (out []byte) {
	if len(b) >= 2 {
		p.Bytes = [2]byte{b[0], b[1]}
		if len(b) > 2 {
			out = b[2:]
		}
	}
	return
}

func (p *Uint16) Encode() []byte {
	return p.Bytes[:]
}

func (p *Uint16) Get() uint16 {
	return binary.BigEndian.Uint16(p.Bytes[:2])
}

func (p *Uint16) String() string {
	return strconv.FormatUint(uint64(binary.BigEndian.Uint16(p.Bytes[:2])), 10)
}

func (p *Uint16) Put(i uint16) *Uint16 {
	binary.BigEndian.PutUint16(p.Bytes[:], i)
	return p
}

func GetPort(listener string) simplebuffer.Serializer {
	//log.DEBUG(listener)
	_, p, err := net.SplitHostPort(listener)
	if err != nil {
		log.ERROR(err)
		return nil
	}
	oI, err := strconv.ParseUint(p, 10, 16)
	if err != nil {
		log.ERROR(err)
		return nil
	}
	port := &Uint16{}
	port.Put(uint16(oI))
	return port
}
