package rpc

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/p9c/pod/cmd/node/state"
	"github.com/p9c/pod/pkg/util/cl"
)

// DefaultConnectTimeout is a reasonable 30 seconds
var DefaultConnectTimeout = time.Second * 30

// Dial connects to the address on the named network using the appropriate
// dial function depending on the address and configuration options.
// For example .onion addresses will be dialed using the onion specific proxy
// if one was specified, but will otherwise use the normal dial function (
// which could itself use a proxy or not).
var Dial = func(statecfg *state.Config) func(addr net.Addr) (net.Conn, error) {
	return func(addr net.Addr) (net.Conn, error) {
		if strings.Contains(addr.String(), ".onion:") {
			return statecfg.Oniondial(addr.Network(), addr.String(),
				DefaultConnectTimeout)
		}
		log <- cl.Trace{"StateCfg.Dial", addr.Network(), addr.String(), DefaultConnectTimeout}
		con, er := statecfg.Dial(addr.Network(), addr.String(), DefaultConnectTimeout)
		if er != nil {
			log <- cl.Trace{con, er}
		}
		return con, er
	}
}

// Lookup resolves the IP of the given host using the correct DNS lookup
// function depending on the configuration options.  For example,
// addresses will be resolved using tor when the --proxy flag was specified
// unless --noonion was also specified in which case the normal system DNS
// resolver will be used. Any attempt to resolve a tor address (.
// onion) will return an error since they are not intended to be resolved
// outside of the tor proxy.
var Lookup = func(statecfg *state.Config) func(host string) ([]net.IP, error) {
	return func(host string) ([]net.IP, error) {
		if strings.HasSuffix(host, ".onion") {
			return nil, fmt.Errorf("attempt to resolve tor address %s", host)
		}
		return statecfg.Lookup(host)
	}
}
