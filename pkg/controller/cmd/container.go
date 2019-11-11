package main

import (
	"github.com/p9c/pod/pkg/controller"
	"github.com/p9c/pod/pkg/log"
	"net"
)

func main() {
	log.L.SetLevel("trace", true)
	var ipa1 = net.ParseIP("127.0.0.1")
	var ipa2 = net.ParseIP("fe80::6382:2df5:7014:e156")
	ips := controller.NewIPs()
	ips.Put([]*net.IP{&ipa1, &ipa2})
	ips2 := controller.NewIPs()
	ips2.Decode(ips.Encode())
	dec := ips.Get()
	dec2 := ips2.Get()
	for i := range dec {
		if !dec[i].Equal(*dec2[i]) {

		}
	}
}
