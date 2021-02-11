package routeable

import (
	"errors"
	"github.com/davecgh/go-spew/spew"
	"net"
	"strings"
	
	"github.com/jackpal/gateway"
)

// TODO: android and ios need equivalent functions as gateway.DiscoverGateway

// Gateway stores the current network default gateway as discovered by
// github.com/jackpal/gateway
var Gateway net.IP

// Address is the network address that routes to the gateway and thus the
// internet
var Address net.IP

// Interface is the net.Interface of the Address above
var Interface *net.Interface

// SecondaryAddresses are all the other addresses that can be reached from
// somewhere (including localhost) but not necessarily the internet
var SecondaryAddresses []net.IP

// SecondaryInterfaces is the interfaces of the SecondaryAddresses stored in the
// corresponding slice index
var SecondaryInterfaces []*net.Interface

// Discover enumerates and evaluates all known network interfaces and addresses
// and filters it down to the ones that reach both a LAN and the internet
//
// We are only interested in IPv4 addresses because for the most part, domestic
// ISPs do not issue their customers with IPv6 routing, it's still a pain in the
// ass outside of large data centre connections
func Discover() (err error) {
	Info("discovering routeable interfaces and addresses...")
	var nif []net.Interface
	if nif, err = net.Interfaces(); Check(err) {
		return
	}
	// Debug("number of available network interfaces:", len(nif))
	// Debugs(nif)
	if Gateway, err = gateway.DiscoverGateway(); Check(err) {
		// todo: this error condition always happens on iOS and Android
		// return
		for i := range nif {
			Debugs(nif[i])
		}
	} else {
		if Gateway != nil {
			// Debug("default gateway:", Gateway)
			gws := Gateway.String()
			gw := net.ParseIP(gws)
			// out:
			for i := range nif {
				// Debugs(nif[i])
				var addrs []net.Addr
				if addrs, err = nif[i].Addrs(); Check(err) || addrs == nil {
					continue
				}
				for j := range addrs {
					var in *net.IPNet
					if _, in, err = net.ParseCIDR(addrs[j].String()); Check(err) {
						continue
					}
					if in.Contains(gw) {
						// Debug("network connection reachable from internet:", nif[i].Name, addrs[j].String())
						Address = net.ParseIP(strings.Split(addrs[j].String(), "/")[0])
						// Debugs(Address)
						Interface = &nif[i]
						// Debugs(Interface)
						// break out
					} else {
						// add them to the other address/interface lists
						// Debug(
						// 	"network connection not reachable (directly) from the internet",
						// 	nif[i].Name,
						// 	addrs[j].String(),
						// )
						ip, _, _ := net.ParseCIDR(addrs[j].String())
						SecondaryAddresses = append(SecondaryAddresses, ip)
						SecondaryInterfaces = append(SecondaryInterfaces, &nif[i])
					}
				}
			}
		}
	}
	Debug("Gateway", Gateway)
	Debug("Address", Address)
	Debug("Interface", spew.Sdump(Interface))
	Debug("SecondaryAddresses", spew.Sdump(SecondaryAddresses))
	Debug("SecondaryInterfaces", spew.Sdump(SecondaryInterfaces))
	return
}

// GetInterface returns the address and interface of multicast capable interfaces
func GetInterface() (ifc *net.Interface, address string, err error) {
	// var err error
	// var nif []net.Interface
	// nif, err = net.Interfaces()
	// if err != nil {
	// 	Error("error:", err)
	// }
	// // // Traces(interfaces)
	// // for ifi := range interfaces {
	// // 	if interfaces[ifi].Flags&net.FlagLoopback == 0 && interfaces[ifi].
	// // 		HardwareAddr != nil {
	// // 		// iads, _ := interfaces[ifi].Addrs()
	// // 		// for i := range iads {
	// // 		//	//Traces(iads[i].Network())
	// // 		// }
	// // 		// Debug(interfaces[ifi].MulticastAddrs())
	// // 		lanInterface = append(lanInterface, &interfaces[ifi])
	// // 	}
	// // }
	// var routeableAddress string
	// for i := range nif {
	// 	// Debug(nif[i].Addrs())
	// 	// Debug(nif[i].HardwareAddr)
	// 	// filter out known virtual devices
	// 	// microsoft hyper-v virtual interface
	// 	if strings.HasPrefix(nif[i].HardwareAddr.String(), "00:15:5d") {
	// 		continue
	// 	}
	// 	// todo: below here add discovered useful non-physical network interface tests like the one above
	// 	addrs, _ := nif[i].Addrs()
	// 	// Debug(addrs)
	// 	for j := range addrs {
	// 		// Debug(addresses[i].String())
	// 		if !strings.ContainsAny(addrs[j].String(), ":") {
	// 			routeableAddress = strings.Split(addrs[j].String(), "/")[0]
	// 			// all addresses except localhost can exit potentially to the internet, on linux often these show first
	// 			if strings.HasPrefix(routeableAddress, "127") {
	// 				continue
	// 			}
	// 			if routeableAddress != "" {
	// 				address = append(address, routeableAddress)
	//
	// 			}
	// 			break
	// 		}
	// 	}
	// 	// Debug(addresses)
	// 	if len(address) > 0 {
	// 		ifc = append(ifc, nif[i])
	// 	}
	// }
	// if routeableAddress == "" {
	// 	panic("no network available")
	// }
	// // Traces(lanInterface)
	if Address == nil || Interface == nil {
		if Discover() != nil {
			err = errors.New("no routeable address found")
			return
		}
	}
	address = Address.String()
	ifc = Interface
	return
}

func GetListenable() net.IP {
	// // first add the interface addresses
	// rI, _ := GetInterface()
	// var lA []net.TCPAddr
	// for i := range rI {
	// 	l, err := rI[i].Addrs()
	// 	if err != nil {
	// 		Error(err)
	// 		return nil
	// 	}
	// 	for j := range l {
	// 		ljs := l[j].String()
	// 		ip := net.ParseIP(ljs)
	// 		lA = append(
	// 			lA, net.TCPAddr{IP: ip},
	// 		)
	// 	}
	// }
	if Address == nil {
		if Discover() != nil {
			Error("no routeable address found")
		}
	}
	return Address
}
