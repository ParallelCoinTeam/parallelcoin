package peer_test

import (
	"fmt"
	"github.com/p9c/log"
	"github.com/p9c/matrjoska/version"
	"net"
	"time"
	
	"github.com/p9c/qu"
	
	"github.com/p9c/matrjoska/pkg/chaincfg"
	"github.com/p9c/matrjoska/pkg/peer"
	"github.com/p9c/matrjoska/pkg/wire"
)

var subsystem = log.AddLoggerSubsystem(version.PathBase)
var F, E, W, I, D, T log.LevelPrinter = log.GetLogPrinterSet(subsystem)

// mockRemotePeer creates a basic inbound peer listening on the simnet port for use with Example_peerConnection. It does
// not return until the listner is active.
func mockRemotePeer() (e error) {
	// Configure peer to act as a simnet node that offers no services.
	peerCfg := &peer.Config{
		UserAgentName:    "peer",  // User agent name to advertise.
		UserAgentVersion: "1.0.0", // User agent version to advertise.
		ChainParams:      &chaincfg.SimNetParams,
		TrickleInterval:  time.Second * 10,
	}
	// Accept connections on the simnet port.
	listener, e := net.Listen("tcp", "127.0.0.1:18555")
	if e != nil {
		return e
	}
	go func() {
		conn, e := listener.Accept()
		if e != nil {
			E.F("Accept: error %v", e)
			return
		}
		// Create and start the inbound peer.
		p := peer.NewInboundPeer(peerCfg)
		p.AssociateConnection(conn)
	}()
	return nil
}

// This example demonstrates the basic process for initializing and creating an outbound peer. Peers negotiate by
// exchanging version and verack messages. For demonstration, a simple handler for version message is attached to the
// peer.
func Example_newOutboundPeer() {
	// Ordinarily this will not be needed since the outbound peer will be connecting to a remote peer, however, since this example is executed and tested, a mock remote peer is needed to listen for the outbound peer.
	if e := mockRemotePeer(); E.Chk(e) {
		E.F("mockRemotePeer: unexpected error %v", e)
		return
	}
	// Create an outbound peer that is configured to act as a simnet node that offers no services and has listeners for the version and verack messages.  The verack listener is used here to signal the code below when the handshake has been finished by signalling a channel.
	verack := qu.T()
	peerCfg := &peer.Config{
		UserAgentName:    "peer",  // User agent name to advertise.
		UserAgentVersion: "1.0.0", // User agent version to advertise.
		ChainParams:      &chaincfg.SimNetParams,
		Services:         0,
		TrickleInterval:  time.Second * 10,
		Listeners: peer.MessageListeners{
			OnVersion: func(p *peer.Peer, msg *wire.MsgVersion) *wire.MsgReject {
				fmt.Println("outbound: received version")
				return nil
			},
			OnVerAck: func(p *peer.Peer, msg *wire.MsgVerAck) {
				verack <- struct{}{}
			},
		},
	}
	p, e := peer.NewOutboundPeer(peerCfg, "127.0.0.1:18555")
	if e != nil {
		E.F("NewOutboundPeer: error %v", e)
		return
	}
	// Establish the connection to the peer address and mark it connected.
	conn, e := net.Dial("tcp", p.Addr())
	if e != nil {
		E.F("net.Dial: error %v", e)
		return
	}
	p.AssociateConnection(conn)
	// Wait for the verack message or timeout in case of failure.
	select {
	case <-verack.Wait():
	case <-time.After(time.Second * 1):
		E.F("Example_peerConnection: verack timeout")
	}
	// Disconnect the peer.
	p.Disconnect()
	p.WaitForDisconnect()
	// Output:
	// outbound: received version
}
