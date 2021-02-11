package routeable

import (
	"errors"
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
		var gw net.IP
		if Gateway != nil {
			gws := Gateway.String()
			gw = net.ParseIP(gws)
		}
		for i := range nif {
			var addrs []net.Addr
			if addrs, err = nif[i].Addrs(); Check(err) || addrs == nil {
				continue
			}
			for j := range addrs {
				var in *net.IPNet
				if _, in, err = net.ParseCIDR(addrs[j].String()); Check(err) {
					continue
				}
				if Gateway != nil && in.Contains(gw) {
					Address = net.ParseIP(strings.Split(addrs[j].String(), "/")[0])
					Interface = &nif[i]
					continue
				}
				ip, _, _ := net.ParseCIDR(addrs[j].String())
				SecondaryAddresses = append(SecondaryAddresses, ip)
				SecondaryInterfaces = append(SecondaryInterfaces, &nif[i])
			}
		}
	}
	Debug("Gateway", Gateway)
	Debug("Address", Address)
	Debug("Interface", Interface.Name)
	Debug("SecondaryAddresses")
	for i := range SecondaryInterfaces {
		Debug(SecondaryInterfaces[i].Name, SecondaryAddresses[i].String())
	}
	return
}

// GetInterface returns the address and interface of multicast-and-internet capable interfaces
func GetInterface() (ifc *net.Interface, address string, err error) {
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
	if Address == nil {
		if Discover() != nil {
			Error("no routeable address found")
		}
	}
	return Address
}
