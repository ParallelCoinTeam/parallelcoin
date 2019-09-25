package discovery

import (
	"net"
	"strings"

	externalip "github.com/GlenDC/go-external-ip"
	"github.com/jackpal/gateway"

	"github.com/p9c/pod/pkg/util/cl"
)

// GetRouteableInterface returns the address and interface of the internet
// -facing network interface
func GetRouteableInterface() (lanInterface *net.Interface) {
	var gw net.IP
	var err error
	if gw, err = gateway.DiscoverGateway(); err != nil {
		log <- cl.Error{"gateway error: ", err, cl.Ine()}
		return nil
	}
	gwMasked := gw.Mask(gw.DefaultMask())
	var ifAddrs []net.Addr
	ifAddrs, err = net.InterfaceAddrs()
	if err != nil {
		log <- cl.Error{"gateway mask error: ", err, cl.Ine()}
		return nil
	}
	var ad net.IP
	for _, x := range ifAddrs {
		address := strings.Split(x.String(), "/")[0]
		a := net.ParseIP(address)
		masked := a.Mask(gw.DefaultMask())
		if masked.String() == gwMasked.String() {
			ad = a
		}
	}
	if ad == nil {
		log <- cl.Error{"somehow didn't find a LAN interface even though we" +
			" have a gateway", cl.Ine()}
		return nil
	}
	consensus := externalip.DefaultConsensus(nil, nil)
	ip, err := consensus.ExternalIP()
	nat := false
	if err != nil {
		log <- cl.Error{"could not get external IP, " +
			"probably no network connection", cl.Ine()}
		return nil
	} else {
		if ip.String() != ad.String() {
			nat = true
		}
	}
	if !nat {
		log <- cl.Warn{"we are directly on the internet", cl.Ine()}
	}
	var interfaces []net.Interface
	interfaces, err = net.Interfaces()
	if err != nil {
		log <- cl.Error{"error:", err}
	}
	for i := range interfaces {
		if ifs, err := interfaces[i].Addrs(); err == nil {
			for j := range ifs {
				ss := strings.Split(ifs[j].String(), "/")
				if ss[0] == ad.String() {
					lanInterface = &interfaces[i]
				}
			}
		}
	}
	return
}
