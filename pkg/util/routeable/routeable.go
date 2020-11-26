package routeable

import (
	"net"
	"strings"
)

// GetInterface returns the address and interface of multicast capable interfaces
func GetInterface() (lanInterface []*net.Interface) {
	var err error
	var nif []net.Interface
	nif, err = net.Interfaces()
	if err != nil {
		Error("error:", err)
	}
	// // Traces(interfaces)
	// for ifi := range interfaces {
	// 	if interfaces[ifi].Flags&net.FlagLoopback == 0 && interfaces[ifi].
	// 		HardwareAddr != nil {
	// 		// iads, _ := interfaces[ifi].Addrs()
	// 		// for i := range iads {
	// 		//	//Traces(iads[i].Network())
	// 		// }
	// 		// Debug(interfaces[ifi].MulticastAddrs())
	// 		lanInterface = append(lanInterface, &interfaces[ifi])
	// 	}
	// }
	var routeableAddress string
	for i := range nif {
		// Debug(nif[i].Addrs())
		// Debug(nif[i].HardwareAddr)
		// filter out known virtual devices
		// microsoft hyper-v virtual interface
		if strings.HasPrefix(nif[i].HardwareAddr.String(), "00:15:5d") {
			continue
		}
		// todo: below here add discovered useful non-physical network interface tests like the one above
		addresses, _ := nif[i].Addrs()
		for i := range addresses {
			// Debug(addresses[i].String())
			if !strings.ContainsAny(addresses[i].String(), ":") {
				routeableAddress = strings.Split(addresses[i].String(), "/")[0]
				if routeableAddress != "" {
					lanInterface = []*net.Interface{&nif[i]}
					break
				}
			}
		}
	}
	if routeableAddress == "" {
		panic("no network available")
	}
	// Traces(lanInterface)
	return
}
