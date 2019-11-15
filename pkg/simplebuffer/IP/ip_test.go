package IP

import (
	"net"
	"testing"
)

func TestIP(t *testing.T) {
	var ipa = net.ParseIP("127.0.0.1")
	ip := New()
	ip.Put(&ipa)
	ip2 := New()
	ip2.Decode(ip.Encode())
	if !ip.Get().Equal(*ip2.Get()) {
		t.Fail()
	}
}
