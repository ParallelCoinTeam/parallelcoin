package routeable

import (
	"net"
)

// GetInterface returns the address and interface of multicast capable
// interfaces
func GetInterface() (lanInterface []*net.Interface) {
	var err error
	var interfaces []net.Interface
	interfaces, err = net.Interfaces()
	if err != nil {
		L.Error("error:", err)
	}
	// L.Traces(interfaces)
	for ifi := range interfaces {
		if interfaces[ifi].Flags&net.FlagLoopback == 0 && interfaces[ifi].
			HardwareAddr != nil {
			// iads, _ := interfaces[ifi].Addrs()
			// for i := range iads {
			//	//L.Traces(iads[i].Network())
			// }
			// L.Debug(interfaces[ifi].MulticastAddrs())
			lanInterface = append(lanInterface, &interfaces[ifi])
		}
	}
	// L.Traces(lanInterface)
	return
}
