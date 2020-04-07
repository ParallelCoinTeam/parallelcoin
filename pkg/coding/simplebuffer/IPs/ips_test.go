package IPs

import (
	"net"
	"testing"
)

func TestIPs(t *testing.T) {
	var ipa1 = net.ParseIP("127.0.0.1")
	var ipa2 = net.ParseIP("fe80::6382:2df5:7014:e156")
	ips := New()
	ips.Put([]*net.IP{&ipa1, &ipa2})
	ips2 := New()
	ips2.Decode(ips.Encode())
	dec := ips.Get()
	dec2 := ips2.Get()
	for i := range dec {
		if !dec[i].Equal(*dec2[i]) {
			t.Fail()
		}
	}
}
