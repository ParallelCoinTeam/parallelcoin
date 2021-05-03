package peersummary

import "net"

// PeerSummary stores the salient network location information about a peer
type PeerSummary struct {
	IP      net.IP
	Inbound bool
}
type PeerSummaries []PeerSummary
