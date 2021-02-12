package p2padvt

import (
	"net"
	
	"github.com/niubaoshu/gotiny"
	
	"github.com/p9c/pod/pkg/util"
	"github.com/p9c/pod/pkg/util/routeable"
	
	"github.com/p9c/pod/app/conte"
)

var Magic = []byte{'a', 'd', 'v', 1}

type Advertisment struct {
	IPs  []net.IP
	P2P  uint16
	UUID uint64
}

// Get returns an advertisment message
func Get(cx *conte.Xt) []byte {
	P2P := util.GetActualPort((*cx.Config.P2PListeners)[0])
	_, ips := routeable.GetAddressesAndInterfaces()
	adv := Advertisment{
		IPs:  ips, // routeable.GetListenable(),
		P2P:  P2P,
		UUID: cx.UUID,
	}
	ad := gotiny.Marshal(&adv)
	return ad
}

// GetAdvt returns an advertisment struct
func GetAdvt(cx *conte.Xt) *Advertisment {
	_, ips := routeable.GetAddressesAndInterfaces()
	adv := &Advertisment{
		IPs:  ips,
		P2P:  util.GetActualPort((*cx.Config.P2PListeners)[0]),
		UUID: cx.UUID,
	}
	return adv
}
