package rpc

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/parallelcointeam/parallelcoin/cmd/node/state"
)

// DefaultConnectTimeout is a reasonable 30 seconds
var DefaultConnectTimeout = time.Second * 30

// Dial connects to the address on the named network using the appropriate
// dial function depending on the address and configuration options.
// For example .onion addresses will be dialed using the onion specific proxy
// if one was specified, but will otherwise use the normal dial function (
// which could itself use a proxy or not).
var Dial = func(stateCfg *state.Config) func(addr net.Addr) (net.Conn, error) {
	return func(addr net.Addr) (net.Conn, error) {
		if strings.Contains(addr.String(), ".onion:") {
			return stateCfg.Oniondial(addr.Network(), addr.String(),
				DefaultConnectTimeout)
		}
		l.Trace("StateCfg.Dial", addr.Network(), addr.String(),
			DefaultConnectTimeout)
		conn, er := stateCfg.Dial(addr.Network(), addr.String(), DefaultConnectTimeout)
		if er != nil {
			l.Trace(conn, er)
		}
		return conn, er
	}
}

// Lookup resolves the IP of the given host using the correct DNS lookup
// function depending on the configuration options.  For example,
// addresses will be resolved using tor when the --proxy flag was specified
// unless --noonion was also specified in which case the normal system DNS
// resolver will be used. Any attempt to resolve a tor address (.
// onion) will return an error since they are not intended to be resolved
// outside of the tor proxy.
var Lookup = func(stateCfg *state.Config) func(host string) ([]net.IP, error) {
	return func(host string) ([]net.IP, error) {
		if strings.HasSuffix(host, ".onion") {
			return nil, fmt.Errorf("attempt to resolve tor address %s", host)
		}
		return stateCfg.Lookup(host)
	}
}
