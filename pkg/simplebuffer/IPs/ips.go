package IPs

import (
	log "github.com/p9c/logi"
	"github.com/p9c/pod/pkg/routeable"
	"github.com/p9c/pod/pkg/simplebuffer"
	"github.com/p9c/pod/pkg/simplebuffer/IP"
	"net"
	"strings"
)

type IPs struct {
	Length byte
	IPs    []IP.IP
}

func New() *IPs {
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
			i := &IP.IP{}
			out = i.Decode(out)
			ipa := make(net.IP, 16)
			copy(ipa, i.Bytes)
			nIP := IP.New()
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
	ips.IPs = make([]IP.IP, len(in))
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

func GetListenable() simplebuffer.Serializer {
	// first add the interface addresses
	rI := routeable.GetInterface()
	var lA []net.Addr
	for i := range rI {
		l, err := rI[i].Addrs()
		if err != nil {
			log.L.Error(err)
			return nil
		}
		lA = append(lA, l...)
	}
	ips := New()
	var ipslice []*net.IP
	for i := range lA {
		addIP := net.ParseIP(strings.Split(lA[i].String(), "/")[0])
		if addIP.To4() != nil {
			// we only want ipv4 addresses because even ipv6 on lans is still
			// rare
			ipslice = append(ipslice, &addIP)
		}
	}
	ips.Put(ipslice)
	return ips
}
