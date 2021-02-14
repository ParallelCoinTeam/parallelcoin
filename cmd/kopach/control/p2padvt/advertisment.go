package p2padvt

import (
	"github.com/niubaoshu/gotiny"
	"github.com/p9c/pod/pkg/chain/wire"
	
	"github.com/p9c/pod/pkg/util"
	"github.com/p9c/pod/pkg/util/routeable"
	
	"github.com/p9c/pod/app/conte"
)

var Magic = []byte{'a', 'd', 'v', 1}

// Advertisment is the contact details, UUID and services available on a node
type Advertisment struct {
	// IPs is here stored as an empty struct map so shuffling is automatic
	IPs map[string]struct{}
	// P2P is the port on which the node can be contacted by other peers
	P2P uint16
	// UUID is a unique identifier randomly generated at each initialisation
	UUID uint64
	// Services reflects the services available, in time the controller will be
	// merged into wire, along with other key protocols implemented in pod
	Services wire.ServiceFlag
}

// Get returns an advertisment message
func Get(cx *conte.Xt) []byte {
	P2P := util.GetActualPort((*cx.Config.P2PListeners)[0])
	_, ips := routeable.GetAddressesAndInterfaces()
	adv := Advertisment{
		IPs:      ips, // routeable.GetListenable(),
		P2P:      P2P,
		UUID:     cx.UUID,
		Services: cx.RealNode.Services,
	}
	ad := gotiny.Marshal(&adv)
	return ad
}

// GetAdvt returns an advertisment struct
func GetAdvt(cx *conte.Xt) *Advertisment {
	_, ips := routeable.GetAddressesAndInterfaces()
	adv := &Advertisment{
		IPs:      ips,
		P2P:      util.GetActualPort((*cx.Config.P2PListeners)[0]),
		UUID:     cx.UUID,
		Services: cx.RealNode.Services,
	}
	return adv
}
