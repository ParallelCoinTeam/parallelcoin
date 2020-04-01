package IP

import (
	"net"
)

type IP struct {
	Length byte
	Bytes  []byte
}

func New() *IP {
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

func (i *IP) String() string {
	return i.Get().String()
}

func (i *IP) Put(ip *net.IP) *IP {
	i.Bytes = make([]byte, len(*ip))
	copy(i.Bytes, *ip)
	i.Length = byte(len(i.Bytes))
	return i
}
