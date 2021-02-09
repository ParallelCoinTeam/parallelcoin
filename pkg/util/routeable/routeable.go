package routeable

import (
	"net"
	"strings"
)

// Address is a structure for storing a network address that is both routeable
// to the internet and has a local area network with multicast available
type Address struct {
	Interface net.Interface
	Address   string
}

// Addresses are the collection of address/interface combinations that have been
// found to be routeable and have multicast
var Addresses []Address

// Discover scans the OS network configuration and populates the foregoing fields
func Discover() {

}

// GetInterface returns the address and interface of multicast capable interfaces
func GetInterface() (interfaces []net.Interface, addresses []string) {
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
		addrs, _ := nif[i].Addrs()
		// Debug(addrs)
		for j := range addrs {
			// Debug(addresses[i].String())
			if !strings.ContainsAny(addrs[j].String(), ":") {
				routeableAddress = strings.Split(addrs[j].String(), "/")[0]
				// all addresses except localhost can exit potentially to the internet, on linux often these show first
				if strings.HasPrefix(routeableAddress, "127") {
					continue
				}
				if routeableAddress != "" {
					addresses = append(addresses, routeableAddress)
					
				}
				break
			}
		}
		// Debug(addresses)
		if len(addresses) > 0 {
			interfaces = append(interfaces, nif[i])
		}
	}
	if routeableAddress == "" {
		panic("no network available")
	}
	// Traces(lanInterface)
	return
}

func GetListenable() []net.TCPAddr {
	// first add the interface addresses
	rI, _ := GetInterface()
	var lA []net.TCPAddr
	for i := range rI {
		l, err := rI[i].Addrs()
		if err != nil {
			Error(err)
			return nil
		}
		for j := range l {
			ljs := l[j].String()
			ip := net.ParseIP(ljs)
			lA = append(
				lA, net.TCPAddr{IP: ip},
			)
		}
	}
	return lA
}
