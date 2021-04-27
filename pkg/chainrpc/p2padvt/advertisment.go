package p2padvt

import (
	"github.com/niubaoshu/gotiny"
	
	"github.com/p9c/matrjoska/pkg/util"
	"github.com/p9c/matrjoska/pkg/util/routeable"
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
	// // Services reflects the services available, in time the controller will be
	// // merged into wire, along with other key protocols implemented in pod
	// Services wire.ServiceFlag
}

// Get returns an advertisment message
func Get(uuid uint64, listeners string) []byte {
	_, ips := routeable.GetAddressesAndInterfaces()
	adv := &Advertisment{
		IPs:  ips,
		P2P:  util.GetActualPort(listeners),
		UUID: uuid,
		// Services: node.Services,
	}
	ad := gotiny.Marshal(&adv)
	return ad
}
