package routeable

import (
	"net"

	log "github.com/p9c/logi"
)

// GetInterface returns the address and interface of multicast capable
// interfaces
func GetInterface() (lanInterface []*net.Interface) {
	var err error
	var interfaces []net.Interface
	interfaces, err = net.Interfaces()
	if err != nil {
		log.L.Error("error:", err)
	}
	// log.L.Traces(interfaces)
	for ifi := range interfaces {
		if interfaces[ifi].Flags&net.FlagLoopback == 0 && interfaces[ifi].
			HardwareAddr != nil {
			// iads, _ := interfaces[ifi].Addrs()
			// for i := range iads {
			//	//log.L.Traces(iads[i].Network())
			// }
			// log.L.Debug(interfaces[ifi].MulticastAddrs())
			lanInterface = append(lanInterface, &interfaces[ifi])
		}
	}
	// log.L.Traces(lanInterface)
	return
}
