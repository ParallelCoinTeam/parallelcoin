package routeable

import (
	"github.com/p9c/pod/pkg/log"
	"net"
)

// GetInterface returns the address and interface of multicast capable
// interfaces
func GetInterface() (lanInterface []*net.Interface) {
	var err error
	var interfaces []net.Interface
	interfaces, err = net.Interfaces()
	if err != nil {
		log.ERROR("error:", err)
	}
	//log.SPEW(interfaces)
	for ifi := range interfaces {
		if interfaces[ifi].Flags&net.FlagMulticast != 0 {
			lanInterface=append(lanInterface, &interfaces[ifi])
		}
	}
	//log.SPEW(lanInterface)
	return
}
