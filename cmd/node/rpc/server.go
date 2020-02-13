package rpc

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"net"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	
	"github.com/p9c/pod/cmd/node/mempool"
	"github.com/p9c/pod/cmd/node/state"
	"github.com/p9c/pod/cmd/node/upnp"
	"github.com/p9c/pod/cmd/node/version"
	blockchain "github.com/p9c/pod/pkg/chain"
	chaincfg "github.com/p9c/pod/pkg/chain/config"
	"github.com/p9c/pod/pkg/chain/config/netparams"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	indexers "github.com/p9c/pod/pkg/chain/index"
	"github.com/p9c/pod/pkg/chain/mining"
	cpuminer "github.com/p9c/pod/pkg/chain/mining/cpu"
	netsync "github.com/p9c/pod/pkg/chain/sync"
	txscript "github.com/p9c/pod/pkg/chain/tx/script"
	"github.com/p9c/pod/pkg/chain/wire"
	database "github.com/p9c/pod/pkg/db"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/peer"
	"github.com/p9c/pod/pkg/peer/addrmgr"
	"github.com/p9c/pod/pkg/peer/connmgr"
	"github.com/p9c/pod/pkg/pod"
	"github.com/p9c/pod/pkg/util"
	"github.com/p9c/pod/pkg/util/bloom"
)

const DefaultMaxOrphanTxSize = 100000

type (
	// BroadcastInventoryAdd is a type used to declare that the InvVect it
	// contains needs to be added to the rebroadcast map
	BroadcastInventoryAdd RelayMsg
	// BroadcastInventoryDel is a type used to declare that the InvVect it
	// contains needs to be removed from the rebroadcast map
	BroadcastInventoryDel *wire.InvVect
	// BroadcastMsg provides the ability to house a bitcoin message to be
	// broadcast to all connected peers except specified excluded peers.
	BroadcastMsg struct {
		Message      wire.Message
		ExcludePeers []*NodePeer
	}
	// CFHeaderKV is a tuple of a filter header and its associated block hash.
	// The struct is used to cache cfcheckpt responses.
	CFHeaderKV struct {
		BlockHash    chainhash.Hash
		FilterHeader chainhash.Hash
	}
	// CheckpointSorter implements sort.Interface to allow a slice of
	// checkpoints to be sorted.
	CheckpointSorter []chaincfg.Checkpoint
	ConnectNodeMsg   struct {
		Addr      string
		Permanent bool
		Reply     chan error
	}
	DisconnectNodeMsg struct {
		Cmp   func(*NodePeer) bool
		Reply chan error
	}
	GetAddedNodesMsg struct {
		Reply chan []*NodePeer
	}
	GetConnCountMsg struct {
		Reply chan int32
	}
	GetOutboundGroup struct {
		Key   string
		Reply chan int
	}
	GetPeersMsg struct {
		Reply chan []*NodePeer
	}
	// OnionAddr implements the net.Addr interface and represents a tor address.
	OnionAddr struct {
		Addr string
	}
	// PeerState maintains state of inbound, persistent,
	// outbound peers as well as banned peers and outbound groups.
	PeerState struct {
		InboundPeers    map[int32]*NodePeer
		OutboundPeers   map[int32]*NodePeer
		PersistentPeers map[int32]*NodePeer
		Banned          map[string]time.Time
		OutboundGroups  map[string]int
	}
	// RelayMsg packages an inventory vector along with the newly discovered
	// inventory so the relay has access to that information.
	RelayMsg struct {
		InvVect *wire.InvVect
		Data    interface{}
	}
	RemoveNodeMsg struct {
		Cmp   func(*NodePeer) bool
		Reply chan error
	}
	// Node provides a bitcoin Node for handling communications to and from
	// bitcoin peers.
	Node struct {
		// The following variables must only be used atomically. Putting the
		// uint64s first makes them 64-bit aligned for 32-bit systems.
		BytesReceived        uint64 // Total bytes received from all peers since start.
		BytesSent            uint64 // Total bytes sent by all peers since start.
		StartupTime          int64
		ChainParams          *netparams.Params
		AddrManager          *addrmgr.AddrManager
		ConnManager          *connmgr.ConnManager
		SigCache             *txscript.SigCache
		HashCache            *txscript.HashCache
		RPCServers           []*Server
		SyncManager          *netsync.SyncManager
		Chain                *blockchain.BlockChain
		TxMemPool            *mempool.TxPool
		CPUMiner             *cpuminer.CPUMiner
		ModifyRebroadcastInv chan interface{}
		NewPeers             chan *NodePeer
		DonePeers            chan *NodePeer
		BanPeers             chan *NodePeer
		Query                chan interface{}
		RelayInv             chan RelayMsg
		Broadcast            chan BroadcastMsg
		PeerHeightsUpdate    chan UpdatePeerHeightsMsg
		WG                   sync.WaitGroup
		Quit                 chan struct{}
		NAT                  upnp.NAT
		DB                   database.DB
		TimeSource           blockchain.MedianTimeSource
		Services             wire.ServiceFlag
		// The following fields are used for optional indexes.  They will be nil
		// if the associated index is not enabled.  These fields are set during
		// initial creation of the server and never changed afterwards, so they
		// do not need to be protected for concurrent access.
		TxIndex   *indexers.TxIndex
		AddrIndex *indexers.AddrIndex
		CFIndex   *indexers.CFIndex
		// The fee estimator keeps track of how long transactions are left in
		// the mempool before they are mined into blocks.
		FeeEstimator *mempool.FeeEstimator
		// CFCheckptCaches stores a cached slice of filter headers for
		// cfcheckpt messages for each filter type.
		CFCheckptCaches    map[wire.FilterType][]CFHeaderKV
		CFCheckptCachesMtx sync.RWMutex
		Algo               string
		Config             *pod.Config
		ActiveNet          *netparams.Params
		StateCfg           *state.Config
		GenThreads         uint32
		Started            int32
		Shutdown           int32
		ShutdownSched      int32
	}
	// NodePeer extends the peer to maintain state shared by the server and
	// the blockmanager.
	NodePeer struct {
		*peer.Peer
		// The following variables must only be used atomically
		FeeFilter      int64
		ConnReq        *connmgr.ConnReq
		Server         *Node
		ContinueHash   *chainhash.Hash
		RelayMtx       sync.Mutex
		Filter         *bloom.Filter
		KnownAddresses map[string]struct{}
		BanScore       connmgr.DynamicBanScore
		Quit           chan struct{}
		// The following chans are used to sync blockmanager and server.
		TxProcessed    chan struct{}
		BlockProcessed chan struct{}
		SentAddrs      bool
		IsWhitelisted  bool
		Persistent     bool
		DisableRelayTx bool
	}
	// SimpleAddr implements the net.Addr interface with two struct fields
	SimpleAddr struct {
		Net, Addr string
	}
	// UpdatePeerHeightsMsg is a message sent from the blockmanager to the
	// server after a new block has been accepted. The purpose of the message is
	// to update the heights of peers that were known to announce the block
	// before we connected it to the main chain or recognized it as an orphan.
	// With these updates, peer heights will be kept up to date, allowing for
	// fresh data when selecting sync peer candidacy.
	UpdatePeerHeightsMsg struct {
		NewHash    *chainhash.Hash
		NewHeight  int32
		OriginPeer *peer.Peer
	}
)

const (
	// DefaultServices describes the default services that are supported by
	// the server.
	DefaultServices = wire.SFNodeNetwork | wire.SFNodeBloom |
		wire.SFNodeWitness | wire.SFNodeCF
	// DefaultRequiredServices describes the default services that are required
	// to be supported by outbound peers.
	DefaultRequiredServices = wire.SFNodeNetwork
	// DefaultTargetOutbound is the default number of outbound peers to target.
	DefaultTargetOutbound = 125
	// ConnectionRetryInterval is the base amount of time to wait in between
	// retries when connecting to persistent peers.  It is adjusted by the
	// number of retries such that there is a retry backoff.
	ConnectionRetryInterval = time.Second
)

var (
	// Ensure simpleAddr implements the net.Addr interface.
	_ net.Addr = SimpleAddr{}
	// UserAgentName is the user agent name and is used to help identify
	// ourselves to peers.
	// nolint
	UserAgentName = "pod"
	// UserAgentVersion is the user agent version and is used to help identify
	// ourselves to peers.
	// nolint
	UserAgentVersion = fmt.Sprintf("%d.%d.%d", version.AppMajor,
		version.AppMinor, version.AppPatch)
	// // zeroHash is the zero value hash (all zeros).  It is defined as a
	// convenience.
	// zeroHash chainhash.Hash
)

// Network returns "onion". This is part of the net.Addr interface.
func (oa *OnionAddr) Network() string {
	return "onion"
}

// String returns the onion address. This is part of the net.Addr interface.
func (oa *OnionAddr) String() string {
	return oa.Addr
}

// Count returns the count of all known peers.
func (ps *PeerState) Count() int {
	return len(ps.InboundPeers) + len(ps.OutboundPeers) +
		len(ps.PersistentPeers)
}

// ForAllOutboundPeers is a helper function that runs closure on all outbound
// peers known to peerState.
func (ps *PeerState) ForAllOutboundPeers(closure func(sp *NodePeer)) {
	for _, e := range ps.OutboundPeers {
		closure(e)
	}
	for _, e := range ps.PersistentPeers {
		closure(e)
	}
}

// ForAllPeers is a helper function that runs closure on all peers known to
// peerState.
func (ps *PeerState) ForAllPeers(closure func(sp *NodePeer)) {
	for _, e := range ps.InboundPeers {
		closure(e)
	}
	ps.ForAllOutboundPeers(closure)
}

// AddBytesReceived adds the passed number of bytes to the total bytes received
// counter for the server.  It is safe for concurrent access.
func (s *Node) AddBytesReceived(bytesReceived uint64) {
	atomic.AddUint64(&s.BytesReceived, bytesReceived)
}

// AddBytesSent adds the passed number of bytes to the total bytes sent counter
// for the server.  It is safe for concurrent access.
func (s *Node) AddBytesSent(bytesSent uint64) {
	atomic.AddUint64(&s.BytesSent, bytesSent)
}

// AddPeer adds a new peer that has already been connected to the server.
func (s *Node) AddPeer(sp *NodePeer) {
	s.NewPeers <- sp
}

// AddRebroadcastInventory adds 'iv' to the list of inventories to be
// rebroadcasted at random intervals until they show up in a block.
func (s *Node) AddRebroadcastInventory(iv *wire.InvVect, data interface{}) {
	// Ignore if shutting down.
	if atomic.LoadInt32(&s.Shutdown) != 0 {
		return
	}
	s.ModifyRebroadcastInv <- BroadcastInventoryAdd{InvVect: iv, Data: data}
}

// AnnounceNewTransactions generates and relays inventory vectors and notifies
// both websocket and getblocktemplate long poll clients of the passed
// transactions.  This function should be called whenever new transactions are
// added to the mempool.
func (s *Node) AnnounceNewTransactions(txns []*mempool.TxDesc) {
	// Generate and relay inventory vectors for all newly accepted transactions.
	s.RelayTransactions(txns)
	// Notify both websocket and getblocktemplate long poll clients of all newly
	// accepted transactions.
	for i := range s.RPCServers {
		if s.RPCServers[i] != nil {
			s.RPCServers[i].NotifyNewTransactions(txns)
		}
	}
}

// BanPeer bans a peer that has already been connected to the server by ip.
func (s *Node) BanPeer(sp *NodePeer) {
	s.BanPeers <- sp
}

// BroadcastMessage sends msg to all peers currently connected to the server
// except those in the passed peers to exclude.
func (s *Node) BroadcastMessage(msg wire.Message, exclPeers ...*NodePeer) {
	// XXX: Need to determine if this is an alert that has already been
	// broadcast and refrain from broadcasting again.
	bmsg := BroadcastMsg{Message: msg, ExcludePeers: exclPeers}
	s.Broadcast <- bmsg
}

// ConnectedCount returns the number of currently connected peers.
func (s *Node) ConnectedCount() int32 {
	replyChan := make(chan int32)
	s.Query <- GetConnCountMsg{Reply: replyChan}
	return <-replyChan
}

// NetTotals returns the sum of all bytes received and sent across the network
// for all peers.  It is safe for concurrent access.
func (s *Node) NetTotals() (uint64, uint64) {
	return atomic.LoadUint64(&s.BytesReceived),
		atomic.LoadUint64(&s.BytesSent)
}

// OutboundGroupCount returns the number of peers connected to the given
// outbound group key.
func (s *Node) OutboundGroupCount(
	key string) int {
	replyChan := make(chan int)
	s.Query <- GetOutboundGroup{Key: key, Reply: replyChan}
	return <-replyChan
}

// RelayInventory relays the passed inventory vector to all connected peers
// that are not already known to have it.
func (s *Node) RelayInventory(invVect *wire.InvVect, data interface{}) {
	s.RelayInv <- RelayMsg{InvVect: invVect, Data: data}
}

// RemoveRebroadcastInventory removes 'iv' from the list of items to be
// rebroadcasted if present.
func (s *Node) RemoveRebroadcastInventory(iv *wire.InvVect) {
	// Log<-cl.Debug{emoveBroadcastInventory"
	// Ignore if shutting down.
	if atomic.LoadInt32(&s.Shutdown) != 0 {
		// Log<-cl.Debug{gnoring due to shutdown"
		return
	}
	s.ModifyRebroadcastInv <- BroadcastInventoryDel(iv)
}

// ScheduleShutdown schedules a server shutdown after the specified duration.
// It also dynamically adjusts how often to warn the server is going down based
// on remaining duration.
func (s *Node) ScheduleShutdown(duration time.Duration) {
	// Don't schedule shutdown more than once.
	if atomic.AddInt32(&s.ShutdownSched, 1) != 1 {
		return
	}
	log.WARNF("server shutdown in %v", duration)
	go func() {
		remaining := duration
		tickDuration := DynamicTickDuration(remaining)
		done := time.After(remaining)
		ticker := time.NewTicker(tickDuration)
	out:
		for {
			select {
			case <-done:
				ticker.Stop()
				err := s.Stop()
				if err != nil {
					log.ERROR(err)
				}
				break out
			case <-ticker.C:
				remaining -= -tickDuration
				if remaining < time.Second {
					continue
				}
				// Change tick duration dynamically based on remaining time.
				newDuration := DynamicTickDuration(remaining)
				if tickDuration != newDuration {
					tickDuration = newDuration
					ticker.Stop()
					ticker = time.NewTicker(tickDuration)
				}
				log.WARNF("server shutdown in %v", remaining)
			}
		}
	}()
}

// Start begins accepting connections from peers.
func (s *Node) Start() {
	// Already started?
	if atomic.AddInt32(&s.Started, 1) != 1 {
		return
	}
	log.TRACE("starting server")
	// Server startup time. Used for the uptime command for uptime calculation.
	s.StartupTime = time.Now().Unix()
	// Start the peer handler which in turn starts the address and block
	// managers.
	s.WG.Add(1)
	go s.PeerHandler()
	if s.NAT != nil {
		s.WG.Add(1)
		go s.UPNPUpdateThread()
	}
	if !*s.Config.DisableRPC {
		s.WG.Add(1)
		// Start the rebroadcastHandler, which ensures user tx received by the
		// RPC server are rebroadcast until being included in a block.
		go s.RebroadcastHandler()
		for i := range s.RPCServers {
			s.RPCServers[i].Start()
		}
	}
	// Start the CPU miner if generation is enabled.
	if *s.Config.Generate {
		s.CPUMiner.Start()
	}
}

// Stop gracefully shuts down the server by stopping and disconnecting all
// peers and the main listener.
func (s *Node) Stop() error {
	// Make sure this only happens once.
	if atomic.AddInt32(&s.Shutdown, 1) != 1 {
		log.DEBUG("server is already in the process of shutting down")
		return nil
	}
	log.TRACE("node shutting down")
	// Stop the CPU miner if needed
	s.CPUMiner.Stop()
	// Shutdown the RPC server if it's not disabled.
	if !*s.Config.DisableRPC {
		for i := range s.RPCServers {
			err := s.RPCServers[i].Stop()
			if err != nil {
				log.ERROR(err)
			}
		}
	}
	// Save fee estimator state in the database.
	err := s.DB.Update(func(tx database.Tx) error {
		metadata := tx.Metadata()
		err := metadata.Put(mempool.EstimateFeeDatabaseKey, s.FeeEstimator.Save())
		if err != nil {
			log.ERROR(err)
		}
		return nil
	})
	if err != nil {
		log.ERROR(err)
	}
	// Signal the remaining goroutines to quit.
	close(s.Quit)
	return nil
}

// Transaction has one confirmation on the main chain. Now we can mark it as no
// longer needing rebroadcasting.
func (s *Node) TransactionConfirmed(tx *util.Tx) {
	// Log<-cl.Debug{ransactionConfirmed"
	// Rebroadcasting is only necessary when the RPC server is active.
	for i := range s.RPCServers {
		// Log<-cl.Debug{ending to RPC servers"
		if s.RPCServers[i] == nil {
			return
		}
	}
	// Log<-cl.Debug{etting new inventory vector"
	iv := wire.NewInvVect(wire.InvTypeTx, tx.Hash())
	// Log<-cl.Debug{emoving broadcast inventory"
	s.RemoveRebroadcastInventory(iv)
	// Log<-cl.Debug{one TransactionConfirmed"
}

// UpdatePeerHeights updates the heights of all peers who have have announced
// the latest connected main chain block, or a recognized orphan. These height
// updates allow us to dynamically refresh peer heights, ensuring sync peer
// selection has access to the latest block heights for each peer.
func (s *Node) UpdatePeerHeights(latestBlkHash *chainhash.Hash,
	latestHeight int32, updateSource *peer.Peer) {
	s.PeerHeightsUpdate <- UpdatePeerHeightsMsg{
		NewHash:    latestBlkHash,
		NewHeight:  latestHeight,
		OriginPeer: updateSource,
	}
}

// WaitForShutdown blocks until the main listener and peer handlers are stopped.
func (s *Node) WaitForShutdown() {
	s.WG.Wait()
}

// HandleAddPeerMsg deals with adding new peers.  It is invoked from the
// peerHandler goroutine.
func (s *Node) HandleAddPeerMsg(state *PeerState, sp *NodePeer) bool {
	if sp == nil {
		return false
	}
	// Ignore new peers if we're shutting down.
	if atomic.LoadInt32(&s.Shutdown) != 0 {
		log.INFOF("new peer %s ignored - server is shutting down", sp)
		sp.Disconnect()
		return false
	}
	// Disconnect banned peers.
	host, _, err := net.SplitHostPort(sp.Addr())
	if err != nil {
		log.ERROR("can't split host/port", err)
		sp.Disconnect()
		return false
	}
	if banEnd, ok := state.Banned[host]; ok {
		if time.Now().Before(banEnd) {
			log.DEBUGF("peer %s is banned for another %v - disconnecting %s",
				host, time.Until(banEnd))
			sp.Disconnect()
			return false
		}
		log.INFOF("peer %s is no longer banned", host)
		delete(state.Banned, host)
	}
	// TODO: Check for max peers from a single IP.
	//  Limit max number of total peers.
	if state.Count() >= *s.Config.MaxPeers {
		log.INFOF("max peers reached [%d] - disconnecting peer %s",
			s.Config.MaxPeers, sp)
		sp.Disconnect()
		// TODO: how to handle permanent peers here? they should be rescheduled.
		return false
	}
	// Add the new peer and start it.
	log.TRACE("new peer ", sp)
	if sp.Inbound() {
		state.InboundPeers[sp.ID()] = sp
	} else {
		state.OutboundGroups[addrmgr.GroupKey(sp.NA())]++
		if sp.Persistent {
			state.PersistentPeers[sp.ID()] = sp
		} else {
			state.OutboundPeers[sp.ID()] = sp
		}
	}
	return true
}

// HandleBanPeerMsg deals with banning peers.
// It is invoked from the peerHandler goroutine.
func (s *Node) HandleBanPeerMsg(state *PeerState, sp *NodePeer) {
	host, _, err := net.SplitHostPort(sp.Addr())
	if err != nil {
		log.ERRORF("can't split ban peer %s %v %s", sp.Addr(), err)
		return
	}
	direction := log.DirectionString(sp.Inbound())
	log.INFOF("banned peer %s (%s) for %v", host, direction, *s.Config.BanDuration)
	state.Banned[host] = time.Now().Add(*s.Config.BanDuration)
}

// HandleBroadcastMsg deals with broadcasting messages to peers.
// It is invoked from the peerHandler goroutine.
func (s *Node) HandleBroadcastMsg(state *PeerState, bmsg *BroadcastMsg) {
	state.ForAllPeers(func(sp *NodePeer) {
		if !sp.Connected() {
			return
		}
		for _, ep := range bmsg.ExcludePeers {
			if sp == ep {
				return
			}
		}
		sp.QueueMessage(bmsg.Message, nil)
	})
}

// HandleDonePeerMsg deals with peers that have signalled they are done.  It is
// invoked from the peerHandler goroutine.
func (s *Node) HandleDonePeerMsg(state *PeerState, sp *NodePeer) {
	var list map[int32]*NodePeer
	switch {
	case sp.Persistent:
		list = state.PersistentPeers
	case sp.Inbound():
		list = state.InboundPeers
	default:
		list = state.OutboundPeers
	}
	if _, ok := list[sp.ID()]; ok {
		if !sp.Inbound() && sp.VersionKnown() {
			state.OutboundGroups[addrmgr.GroupKey(sp.NA())]--
		}
		if !sp.Inbound() && sp.ConnReq != nil {
			s.ConnManager.Disconnect(sp.ConnReq.ID())
		}
		delete(list, sp.ID())
		log.TRACE("removed peer ", sp)
		return
	}
	if sp.ConnReq != nil {
		s.ConnManager.Disconnect(sp.ConnReq.ID())
	}
	// Update the address' last seen time if the peer has acknowledged our
	// version and has sent us its version as well.
	if sp.VerAckReceived() && sp.VersionKnown() && sp.NA() != nil {
		s.AddrManager.Connected(sp.NA())
	}
	// If we get here it means that either we didn't know about the peer or we
	// purposefully deleted it.
}

// HandleQuery is the central handler for all queries and commands from other
// goroutines related to peer state.
func (s *Node) HandleQuery(state *PeerState, querymsg interface{}) {
	switch msg := querymsg.(type) {
	case GetConnCountMsg:
		nconnected := int32(0)
		state.ForAllPeers(func(sp *NodePeer) {
			if sp.Connected() {
				nconnected++
			}
		})
		msg.Reply <- nconnected
	case GetPeersMsg:
		peers := make([]*NodePeer, 0, state.Count())
		state.ForAllPeers(func(sp *NodePeer) {
			if !sp.Connected() {
				return
			}
			peers = append(peers, sp)
		})
		msg.Reply <- peers
	case ConnectNodeMsg:
		// TODO: duplicate oneshots? Limit max number of total peers.
		if state.Count() >= *s.Config.MaxPeers {
			msg.Reply <- errors.New("max peers reached")
			return
		}
		for _, nodePeer := range state.PersistentPeers {
			if nodePeer.Addr() == msg.Addr {
				if msg.Permanent {
					msg.Reply <- errors.New("nodePeer already connected")
				} else {
					msg.Reply <- errors.New("nodePeer exists as a permanent nodePeer")
				}
				return
			}
		}
		netAddr, err := AddrStringToNetAddr(s.Config, s.StateCfg, msg.Addr)
		if err != nil {
			log.ERROR(err)
			msg.Reply <- err
			return
		}
		// TODO: if too many, nuke a non-perm nodePeer.
		go s.ConnManager.Connect(&connmgr.ConnReq{
			Addr:      netAddr,
			Permanent: msg.Permanent,
		})
		msg.Reply <- nil
	case RemoveNodeMsg:
		found := DisconnectPeer(state.PersistentPeers, msg.Cmp, func(sp *NodePeer) {
			// Keep group counts ok since we remove from the list now.
			state.OutboundGroups[addrmgr.GroupKey(sp.NA())]--
		})
		if found {
			msg.Reply <- nil
		} else {
			msg.Reply <- errors.New("nodePeer not found")
		}
	case GetOutboundGroup:
		count, ok := state.OutboundGroups[msg.Key]
		if ok {
			msg.Reply <- count
		} else {
			msg.Reply <- 0
		}
	// Request a list of the persistent (added) peers.
	case GetAddedNodesMsg:
		// Respond with a slice of the relevant peers.
		peers := make([]*NodePeer, 0, len(state.PersistentPeers))
		for _, sp := range state.PersistentPeers {
			peers = append(peers, sp)
		}
		msg.Reply <- peers
	case DisconnectNodeMsg:
		// Check inbound peers. We pass a nil callback since we don't require any
		// additional actions on disconnect for inbound peers.
		found := DisconnectPeer(state.InboundPeers, msg.Cmp, nil)
		if found {
			msg.Reply <- nil
			return
		}
		// Check outbound peers.
		found = DisconnectPeer(state.OutboundPeers, msg.Cmp, func(sp *NodePeer) {
			// Keep group counts ok since we remove from the list now.
			state.OutboundGroups[addrmgr.GroupKey(sp.NA())]--
		})
		if found {
			// If there are multiple outbound connections to the same ip:port,
			// continue disconnecting them all until no such peers are found.
			for found {
				found = DisconnectPeer(state.OutboundPeers, msg.Cmp,
					func(sp *NodePeer) {
						state.OutboundGroups[addrmgr.GroupKey(sp.NA())]--
					})
			}
			msg.Reply <- nil
			return
		}
		msg.Reply <- errors.New("nodePeer not found")
	}
}

// HandleRelayInvMsg deals with relaying inventory to peers that are not
// already known to have it.  It is invoked from the peerHandler goroutine.
func (s *Node) HandleRelayInvMsg(state *PeerState, msg RelayMsg) {
	state.ForAllPeers(func(sp *NodePeer) {
		if !sp.Connected() {
			return
		}
		// If the inventory is a block and the peer prefers headers, generate and
		// send a headers message instead of an inventory message.
		if msg.InvVect.Type == wire.InvTypeBlock && sp.WantsHeaders() {
			blockHeader, ok := msg.Data.(wire.BlockHeader)
			if !ok {
				log.WARN("underlying data for headers is not a block header")
				return
			}
			msgHeaders := wire.NewMsgHeaders()
			if err := msgHeaders.AddBlockHeader(&blockHeader); err != nil {
				log.ERROR("failed to add block header:", err)
				return
			}
			sp.QueueMessage(msgHeaders, nil)
			return
		}
		if msg.InvVect.Type == wire.InvTypeTx {
			// Don't relay the transaction to the peer when it has transaction
			// relaying disabled.
			if sp.IsRelayTxDisabled() {
				return
			}
			txD, ok := msg.Data.(*mempool.TxDesc)
			if !ok {
				log.WARNF("underlying data for tx inv relay is not a *mempool.TxDesc: %T",
					msg.Data)
				return
			}
			// Don't relay the transaction if the transaction fee-per-kb is less
			// than the peer's feefilter.
			feeFilter := atomic.LoadInt64(&sp.FeeFilter)
			if feeFilter > 0 && txD.FeePerKB < feeFilter {
				return
			}
			// Don't relay the transaction if there is a bloom filter loaded and
			// the transaction doesn't match it.
			if sp.Filter.IsLoaded() {
				if !sp.Filter.MatchTxAndUpdate(txD.Tx) {
					return
				}
			}
		}
		// Queue the inventory to be relayed with the next batch. It will be
		// ignored if the peer is already known to have the inventory.
		sp.QueueInventory(msg.InvVect)
	})
}

// handleUpdatePeerHeight updates the heights of all peers who were known to
// announce a block we recently accepted.
func (s *Node) HandleUpdatePeerHeights(state *PeerState,
	umsg UpdatePeerHeightsMsg) {
	state.ForAllPeers(func(sp *NodePeer) {
		// The origin peer should already have the updated height.
		if sp.Peer == umsg.OriginPeer {
			return
		}
		// This is a pointer to the underlying memory which doesn't change.
		latestBlkHash := sp.LastAnnouncedBlock()
		// Skip this peer if it hasn't recently announced any new blocks.
		if latestBlkHash == nil {
			return
		}
		// If the peer has recently announced a block, and this block matches our
		// newly accepted block, then update their block height.
		if *latestBlkHash == *umsg.NewHash {
			sp.UpdateLastBlockHeight(umsg.NewHeight)
			sp.UpdateLastAnnouncedBlock(nil)
		}
	})
}

// InboundPeerConnected is invoked by the connection manager when a new inbound
// connection is established.  It initializes a new inbound server peer
// instance, associates it with the connection, and starts a goroutine to wait
// for disconnection.
func (s *Node) InboundPeerConnected(conn net.Conn) {
	sp := NewServerPeer(s, false)
	sp.IsWhitelisted = GetIsWhitelisted(s.StateCfg, conn.RemoteAddr())
	sp.Peer = peer.NewInboundPeer(NewPeerConfig(sp))
	sp.AssociateConnection(conn)
	go s.PeerDoneHandler(sp)
}

// OutboundPeerConnected is invoked by the connection manager when a new
// outbound connection is established.  It initializes a new outbound server
// peer instance, associates it with the relevant state such as the connection
// request instance and the connection itself, and finally notifies the address
// manager of the attempt.
func (s *Node) OutboundPeerConnected(c *connmgr.ConnReq, conn net.Conn) {
	sp := NewServerPeer(s, c.Permanent)
	p, err := peer.NewOutboundPeer(NewPeerConfig(sp), c.Addr.String())
	if err != nil {
		log.ERRORF("cannot create outbound peer %s: %v %s", c.Addr, err)
		s.ConnManager.Disconnect(c.ID())
	}
	sp.Peer = p
	sp.ConnReq = c
	sp.IsWhitelisted = GetIsWhitelisted(s.StateCfg, conn.RemoteAddr())
	sp.AssociateConnection(conn)
	go s.PeerDoneHandler(sp)
	s.AddrManager.Attempt(sp.NA())
}

// PeerDoneHandler handles peer disconnects by notifiying the server that it's
// done along with other performing other desirable cleanup.
func (s *Node) PeerDoneHandler(sp *NodePeer) {
	sp.WaitForDisconnect()
	s.DonePeers <- sp
	// Only tell sync manager we are gone if we ever told it we existed.
	if sp.VersionKnown() {
		s.SyncManager.DonePeer(sp.Peer)
		// Evict any remaining orphans that were sent by the peer.
		numEvicted := s.TxMemPool.RemoveOrphansByTag(mempool.Tag(sp.ID()))
		if numEvicted > 0 {
			log.DEBUGF("Evicted %d %s from peer %v (id %d)",
				numEvicted, log.PickNoun(int(numEvicted), "orphan", "orphans"),
				sp, sp.ID())
		}
	}
	close(sp.Quit)
}

// PeerHandler is used to handle peer operations such as adding and removing
// peers to and from the server, banning peers, and broadcasting messages to
// peers.  It must be run in a goroutine.
func (s *Node) PeerHandler() {
	// Start the address manager and sync manager, both of which are needed by
	// peers.  This is done here since their lifecycle is closely tied to this
	// handler and rather than adding more channels to sychronize things, it's
	// easier and slightly faster to simply start and stop them in this handler.
	s.AddrManager.Start()
	s.SyncManager.Start()
	log.TRACE("starting peer handler")
	peerState := &PeerState{
		InboundPeers:    make(map[int32]*NodePeer),
		PersistentPeers: make(map[int32]*NodePeer),
		OutboundPeers:   make(map[int32]*NodePeer),
		Banned:          make(map[string]time.Time),
		OutboundGroups:  make(map[string]int),
	}
	if !*s.Config.DisableDNSSeed || len(*s.Config.ConnectPeers) < 0 {
		// Add peers discovered through DNS to the address manager.
		connmgr.SeedFromDNS(s.ActiveNet, DefaultRequiredServices,
			Lookup(s.StateCfg), func(addrs []*wire.NetAddress) {
				// Bitcoind uses a lookup of the dns seeder here. This is rather
				// strange since the values looked up by the DNS seed lookups will
				// vary quite a lot. to replicate this behaviour we put all
				// addresses as having come from the first one.
				log.DEBUG("adding addresses")
				s.AddrManager.AddAddresses(addrs, addrs[0])
			})
	}
	log.TRACE("starting connmgr")
	go s.ConnManager.Start()
out:
	for {
		select {
		// New peers connected to the server.
		case p := <-s.NewPeers:
			s.HandleAddPeerMsg(peerState, p)
		// Disconnected peers.
		case p := <-s.DonePeers:
			s.HandleDonePeerMsg(peerState, p)
		// Block accepted in mainchain or orphan, update peer height.
		case umsg := <-s.PeerHeightsUpdate:
			s.HandleUpdatePeerHeights(peerState, umsg)
		// Peer to ban.
		case p := <-s.BanPeers:
			s.HandleBanPeerMsg(peerState, p)
		// New inventory to potentially be relayed to other peers.
		case invMsg := <-s.RelayInv:
			s.HandleRelayInvMsg(peerState, invMsg)
		// Message to broadcast to all connected peers except those which are
		// excluded by the message.
		case bmsg := <-s.Broadcast:
			s.HandleBroadcastMsg(peerState, &bmsg)
		case qmsg := <-s.Query:
			s.HandleQuery(peerState, qmsg)
		case <-s.Quit:
			// Disconnect all peers on server shutdown.
			peerState.ForAllPeers(func(sp *NodePeer) {
				log.TRACEF("shutdown peer %s", sp)
				sp.Disconnect()
			})
			break out
		}
	}
	s.ConnManager.Stop()
	err := s.SyncManager.Stop()
	if err != nil {
		log.ERROR(err)
	}
	err = s.AddrManager.Stop()
	if err != nil {
		log.ERROR(err)
	}
	// Drain channels before exiting so nothing is left waiting around to send.
cleanup:
	for {
		select {
		case <-s.NewPeers:
		case <-s.DonePeers:
		case <-s.PeerHeightsUpdate:
		case <-s.RelayInv:
		case <-s.Broadcast:
		case <-s.Query:
		default:
			break cleanup
		}
	}
	s.WG.Done()
	log.TRACEF("peer handler done")
}

// PushBlockMsg sends a block message for the provided block hash to the
// connected peer.  An error is returned if the block hash is not known.
func (s *Node) PushBlockMsg(sp *NodePeer, hash *chainhash.Hash,
	doneChan chan<- struct{}, waitChan <-chan struct{},
	encoding wire.MessageEncoding) error {
	// Fetch the raw block bytes from the database.
	var blockBytes []byte
	err := sp.Server.DB.View(func(dbTx database.Tx) error {
		var err error
		blockBytes, err = dbTx.FetchBlock(hash)
		return err
	})
	if err != nil {
		log.ERRORF("unable to fetch requested block hash %v: %v",
			hash, err)
		if doneChan != nil {
			doneChan <- struct{}{}
		}
		return err
	}
	// Deserialize the block.
	var msgBlock wire.MsgBlock
	err = msgBlock.Deserialize(bytes.NewReader(blockBytes))
	if err != nil {
		log.ERRORF("unable to deserialize requested block hash %v: %v",
			hash, err)
		if doneChan != nil {
			doneChan <- struct{}{}
		}
		return err
	}
	// Once we have fetched data wait for any previous operation to finish.
	if waitChan != nil {
		<-waitChan
	}
	// We only send the channel for this message if we aren't sending an inv
	// straight after.
	var dc chan<- struct{}
	continueHash := sp.ContinueHash
	sendInv := continueHash != nil && continueHash.IsEqual(hash)
	if !sendInv {
		dc = doneChan
	}
	sp.QueueMessageWithEncoding(&msgBlock, dc, encoding)
	// When the peer requests the final block that was advertised in response to
	// a getblocks message which requested more blocks than would fit into a
	// single message, send it a new inventory message to trigger it to issue
	// another getblocks message for the next batch of inventory.
	if sendInv {
		best := sp.Server.Chain.BestSnapshot()
		invMsg := wire.NewMsgInvSizeHint(1)
		iv := wire.NewInvVect(wire.InvTypeBlock, &best.Hash)
		err := invMsg.AddInvVect(iv)
		if err != nil {
			log.ERROR(err)
		}
		sp.QueueMessage(invMsg, doneChan)
		sp.ContinueHash = nil
	}
	return nil
}

// PushMerkleBlockMsg sends a merkleblock message for the provided block hash
// to the connected peer.  Since a merkle block requires the peer to have a
// filter loaded, this call will simply be ignored if there is no filter
// loaded.  An error is returned if the block hash is not known.
func (s *Node) PushMerkleBlockMsg(sp *NodePeer, hash *chainhash.Hash,
	doneChan chan<- struct{}, waitChan <-chan struct{},
	encoding wire.MessageEncoding) error {
	// Do not send a response if the peer doesn't have a filter loaded.
	if !sp.Filter.IsLoaded() {
		if doneChan != nil {
			doneChan <- struct{}{}
		}
		return nil
	}
	// Fetch the raw block bytes from the database.
	blk, err := sp.Server.Chain.BlockByHash(hash)
	if err != nil {
		log.ERRORF("unable to fetch requested block hash %v: %v",
			hash, err)
		if doneChan != nil {
			doneChan <- struct{}{}
		}
		return err
	}
	// Generate a merkle block by filtering the requested block according to the
	// filter for the peer.
	merkle, matchedTxIndices := bloom.NewMerkleBlock(blk, sp.Filter)
	// Once we have fetched data wait for any previous operation to finish.
	if waitChan != nil {
		<-waitChan
	}
	// Send the merkleblock.  Only send the done channel with this message if no
	// transactions will be sent afterwards.
	var dc chan<- struct{}
	if len(matchedTxIndices) == 0 {
		dc = doneChan
	}
	sp.QueueMessage(merkle, dc)
	// Finally, send any matched transactions.
	blkTransactions := blk.MsgBlock().Transactions
	for i, txIndex := range matchedTxIndices {
		// Only send the done channel on the final transaction.
		var dc chan<- struct{}
		if i == len(matchedTxIndices)-1 {
			dc = doneChan
		}
		if txIndex < uint32(len(blkTransactions)) {
			sp.QueueMessageWithEncoding(blkTransactions[txIndex], dc,
				encoding)
		}
	}
	return nil
}

// PushTxMsg sends a tx message for the provided transaction hash to the
// connected peer.  An error is returned if the transaction hash is not known.
func (s *Node) PushTxMsg(sp *NodePeer, hash *chainhash.Hash,
	doneChan chan<- struct{}, waitChan <-chan struct{},
	encoding wire.MessageEncoding) error {
	// Attempt to fetch the requested transaction from the pool.  A call could
	// be made to check for existence first, but simply trying to fetch a
	// missing transaction results in the same behavior.
	tx, err := s.TxMemPool.FetchTransaction(hash)
	if err != nil {
		log.ERRORF("unable to fetch tx %v from transaction pool: %v", hash, err)
		if doneChan != nil {
			doneChan <- struct{}{}
		}
		return err
	}
	// Once we have fetched data wait for any previous operation to finish.
	if waitChan != nil {
		<-waitChan
	}
	sp.QueueMessageWithEncoding(tx.MsgTx(), doneChan, encoding)
	return nil
}

// RebroadcastHandler keeps track of user submitted inventories that we have
// sent out but have not yet made it into a block. We periodically rebroadcast
// them in case our peers restarted or otherwise lost track of them.
func (s *Node) RebroadcastHandler() {
	// Log<-cl.Debug{tarting rebroadcastHandler"
	// Wait 5 min before first tx rebroadcast.
	timer := time.NewTimer(5 * time.Minute)
	pendingInvs := make(map[wire.InvVect]interface{})
out:
	for {
		select {
		case riv := <-s.ModifyRebroadcastInv:
			switch msg := riv.(type) {
			// Incoming InvVects are added to our map of RPC txs.
			case BroadcastInventoryAdd:
				pendingInvs[*msg.InvVect] = msg.Data
			// When an InvVect has been added to a block,
			// we can now remove it, if it was present.
			case BroadcastInventoryDel:
				if _, ok := pendingInvs[*msg]; ok {
					delete(pendingInvs, *msg)
				}
			}
		case <-timer.C:
			// Any inventory we have has not made it into a block yet.
			// We periodically resubmit them until they have.
			for iv, data := range pendingInvs {
				ivCopy := iv
				s.RelayInventory(&ivCopy, data)
			}
			// Process at a random time up to 30mins (in seconds) in the future.
			timer.Reset(time.Second *
				time.Duration(RandomUint16Number(1800)))
		case <-s.Quit:
			break out
			// default:
		}
	}
	timer.Stop()
	// Drain channels before exiting so nothing is left waiting around to send.
cleanup:
	for {
		select {
		case <-s.ModifyRebroadcastInv:
		default:
			break cleanup
		}
	}
	s.WG.Done()
}

// RelayTransactions generates and relays inventory vectors for all of the
// passed transactions to all connected peers.
func (s *Node) RelayTransactions(txns []*mempool.TxDesc) {
	for _, txD := range txns {
		iv := wire.NewInvVect(wire.InvTypeTx, txD.Tx.Hash())
		s.RelayInventory(iv, txD)
	}
}
func (s *Node) UPNPUpdateThread() {
	// Go off immediately to prevent code duplication, thereafter we renew lease
	// every 15 minutes.
	timer := time.NewTimer(0 * time.Second)
	lport, _ := strconv.ParseInt(s.ActiveNet.DefaultPort, 10, 16)
	first := true
out:
	for {
		select {
		case <-timer.C:
			// TODO: pick external port  more cleverly
			// TODO: know which ports we are listening to on an external net.
			// TODO: if specific listen port doesn't work then ask for wildcard
			//  listen port?
			// XXX this assumes timeout is in seconds.
			listenPort, err := s.NAT.AddPortMapping("tcp", int(lport),
				int(lport),
				"pod listen port", 20*60)
			if err != nil {
				log.ERRORF("can't add UPnP port mapping: %v %s", err)
			}
			if first && err == nil {
				// TODO: look this up periodically to see if upnp domain changed
				//  and so did ip.
				externalip, err := s.NAT.GetExternalAddress()
				if err != nil {
					log.ERROR(err)
					log.ERRORF("UPnP can't get external address: %v", err)
					continue out
				}
				na := wire.NewNetAddressIPPort(externalip, uint16(listenPort),
					s.Services)
				err = s.AddrManager.AddLocalAddress(na, addrmgr.UpnpPrio)
				if err != nil {
					log.ERROR(err)
					_ = err
					// XXX DeletePortMapping?
				}
				log.WARNF("successfully bound via UPnP to %s",
					addrmgr.NetAddressKey(na))
				first = false
			}
			timer.Reset(time.Minute * 15)
		case <-s.Quit:
			break out
		}
	}
	timer.Stop()
	if err := s.NAT.DeletePortMapping("tcp", int(lport),
		int(lport)); err != nil {
		log.DEBUGF("unable to remove UPnP port mapping: %v %s", err)
	} else {
		log.DEBUG("successfully cleared UPnP port mapping")
	}
	s.WG.Done()
}

// OnAddr is invoked when a peer receives an addr bitcoin message and is used
// to notify the server about advertised addresses.
func (sp *NodePeer) OnAddr(_ *peer.Peer,
	msg *wire.MsgAddr) {
	// Ignore addresses when running on the simulation test network.  This helps
	// prevent the network from becoming another public test network since it
	// will not be able to learn about other peers that have not specifically
	// been provided.
	if (*sp.Server.Config.Network)[0] == 's' {
		return
	}
	// Ignore old style addresses which don't include a timestamp.
	if sp.ProtocolVersion() < wire.NetAddressTimeVersion {
		return
	}
	// A message that has no addresses is invalid.
	if len(msg.AddrList) == 0 {
		log.ERRORF("command [%s] from %s does not contain any addresses",
			msg.Command(), sp.Peer)
		sp.Disconnect()
		return
	}
	for _, na := range msg.AddrList {
		// Don't add more address if we're disconnecting.
		if !sp.Connected() {
			return
		}
		// Set the timestamp to 5 days ago if it's more than 24 hours in the
		// future so this address is one of the first to be removed when space is
		// needed.
		now := time.Now()
		if na.Timestamp.After(now.Add(time.Minute * 10)) {
			na.Timestamp = now.Add(-1 * time.Hour * 24 * 5)
		}
		// Add address to known addresses for this peer.
		sp.AddKnownAddresses([]*wire.NetAddress{na})
	}
	// Add addresses to server address manager.  The address manager handles the
	// details of things such as preventing duplicate addresses, max addresses,
	// and last seen updates. XXX bitcoind gives a 2 hour time penalty here, do
	// we want to do the same?
	sp.Server.AddrManager.AddAddresses(msg.AddrList, sp.NA())
}

// OnBlock is invoked when a peer receives a block bitcoin message.  It blocks
// until the bitcoin block has been fully processed.
func (sp *NodePeer) OnBlock(_ *peer.Peer, msg *wire.MsgBlock, buf []byte) {
	// Convert the raw MsgBlock to a util.Block which provides some convenience
	// methods and things such as hash caching.
	block := util.NewBlockFromBlockAndBytes(msg, buf)
	// Add the block to the known inventory for the peer.
	iv := wire.NewInvVect(wire.InvTypeBlock, block.Hash())
	sp.AddKnownInventory(iv)
	// Queue the block up to be handled by the block manager and intentionally
	// block further receives until the bitcoin block is fully processed and
	// known good or bad.  This helps prevent a malicious peer from queuing up a
	// bunch of bad blocks before disconnecting (or being disconnected) and
	// wasting memory.  Additionally, this behavior is depended on by at least
	// the block acceptance test tool as the reference implementation processes
	// blocks in the same thread and therefore blocks further messages until the
	// bitcoin block has been fully processed.
	sp.Server.SyncManager.QueueBlock(block, sp.Peer, sp.BlockProcessed)
	<-sp.BlockProcessed
}

// OnFeeFilter is invoked when a peer receives a feefilter bitcoin message and
// is used by remote peers to request that no transactions which have a fee
// rate lower than provided value are inventoried to them.  The peer will be
// disconnected if an invalid fee filter value is provided.
func (sp *NodePeer) OnFeeFilter(_ *peer.Peer,
	msg *wire.MsgFeeFilter) {
	// Check that the passed minimum fee is a valid amount.
	if msg.MinFee < 0 || msg.MinFee > int64(util.MaxSatoshi) {
		log.DEBUGF("peer %v sent an invalid feefilter '%v' -- disconnecting %s",
			sp, util.Amount(msg.MinFee),
		)
		sp.Disconnect()
		return
	}
	atomic.StoreInt64(&sp.FeeFilter, msg.MinFee)
}

// OnFilterAdd is invoked when a peer receives a filteradd bitcoin message and
// is used by remote peers to add data to an already loaded bloom filter.  The
// peer will be disconnected if a filter is not loaded when this message is
// received or the server is not configured to allow bloom filters.
func (sp *NodePeer) OnFilterAdd(_ *peer.Peer,
	msg *wire.MsgFilterAdd) {
	// Disconnect and/or ban depending on the node bloom services flag and
	// negotiated protocol version.
	if !sp.EnforceNodeBloomFlag(msg.Command()) {
		return
	}
	if !sp.Filter.IsLoaded() {
		log.DEBUGF("%s sent a filteradd request with no filter loaded"+
			" -- disconnecting %s", sp)
		sp.Disconnect()
		return
	}
	sp.Filter.Add(msg.Data)
}

// OnFilterClear is invoked when a peer receives a filterclear bitcoin message
// and is used by remote peers to clear an already loaded bloom filter. The
// peer will be disconnected if a filter is not loaded when this message is
// received  or the server is not configured to allow bloom filters.
func (sp *NodePeer) OnFilterClear(_ *peer.Peer,
	msg *wire.MsgFilterClear) {
	// Disconnect and/or ban depending on the node bloom services flag and
	// negotiated protocol version.
	if !sp.EnforceNodeBloomFlag(msg.Command()) {
		return
	}
	if !sp.Filter.IsLoaded() {
		log.DEBUGF("%s sent a filterclear request with no filter loaded"+
			" -- disconnecting %s", sp)
		sp.Disconnect()
		return
	}
	sp.Filter.Unload()
}

// OnFilterLoad is invoked when a peer receives a filterload bitcoin message
// and it used to load a bloom filter that should be used for delivering merkle
// blocks and associated transactions that match the filter. The peer will be
// disconnected if the server is not configured to allow bloom filters.
func (sp *NodePeer) OnFilterLoad(_ *peer.Peer,
	msg *wire.MsgFilterLoad) {
	// Disconnect and/or ban depending on the node bloom services flag and
	// negotiated protocol version.
	if !sp.EnforceNodeBloomFlag(msg.Command()) {
		return
	}
	sp.SetDisableRelayTx(false)
	sp.Filter.Reload(msg)
}

// OnGetAddr is invoked when a peer receives a getaddr bitcoin message and is
// used to provide the peer with known addresses from the address manager.
func (sp *NodePeer) OnGetAddr(_ *peer.Peer,
	msg *wire.MsgGetAddr) {
	// Don't return any addresses when running on the simulation test network.
	// This helps prevent the network from becoming another public test network
	// since it will not be able to learn about other peers that have not
	// specifically been provided.
	if (*sp.Server.Config.Network)[0] == 's' {
		return
	}
	// Do not accept getaddr requests from outbound peers.  This reduces
	// fingerprinting attacks.
	if !sp.Inbound() {
		log.DEBUG("ignoring getaddr request from outbound peer", sp)
		return
	}
	// Only allow one getaddr request per connection to discourage address
	// stamping of inv announcements.
	if sp.SentAddrs {
		log.DEBUGF("ignoring repeated getaddr request from peer %s %s", sp)
		return
	}
	sp.SentAddrs = true
	// Get the current known addresses from the address manager.
	addrCache := sp.Server.AddrManager.AddressCache()
	// Push the addresses.
	sp.PreparePushAddrMsg(addrCache)
}

// OnGetBlocks is invoked when a peer receives a getblocks bitcoin message.
func (sp *NodePeer) OnGetBlocks(_ *peer.Peer,
	msg *wire.MsgGetBlocks) {
	// Find the most recent known block in the best chain based on the block
	// locator and fetch all of the block hashes after it until either
	// wire.MaxBlocksPerMsg have been fetched or the provided stop hash is
	// encountered. Use the block after the genesis block if no other blocks in
	// the provided locator are known.  This does mean the client will start
	// over with the genesis block if unknown block locators are provided. This
	// mirrors the behavior in the reference implementation.
	chain := sp.Server.Chain
	hashList := chain.LocateBlocks(msg.BlockLocatorHashes, &msg.HashStop,
		wire.MaxBlocksPerMsg)
	// Generate inventory message.
	invMsg := wire.NewMsgInv()
	for i := range hashList {
		iv := wire.NewInvVect(wire.InvTypeBlock, &hashList[i])
		err := invMsg.AddInvVect(iv)
		if err != nil {
			log.ERROR(err)
		}
	}
	// Send the inventory message if there is anything to send.
	if len(invMsg.InvList) > 0 {
		invListLen := len(invMsg.InvList)
		if invListLen == wire.MaxBlocksPerMsg {
			// Intentionally use a copy of the final hash so there is not a
			// reference into the inventory slice which would prevent the entire
			// slice from being eligible for GC as soon as it's sent.
			continueHash := invMsg.InvList[invListLen-1].Hash
			sp.ContinueHash = &continueHash
		}
		sp.QueueMessage(invMsg, nil)
	}
}

// OnGetCFCheckpt is invoked when a peer receives a getcfcheckpt bitcoin
// message.
func (sp *NodePeer) OnGetCFCheckpt(_ *peer.Peer,
	msg *wire.MsgGetCFCheckpt) {
	// Ignore getcfcheckpt requests if not in sync.
	if !sp.Server.SyncManager.IsCurrent() {
		return
	}
	// We'll also ensure that the remote party is requesting a set of
	// checkpoints for filters that we actually currently maintain.
	switch msg.FilterType {
	case wire.GCSFilterRegular:
		break
	default:
		log.DEBUG("filter request for unknown checkpoints for filter:",
			msg.FilterType)
		return
	}
	// Now that we know the client is fetching a filter that we know of, we'll
	// fetch the block hashes et each check point interval so we can compare
	// against our cache, and create new check points if necessary.
	blockHashes, err := sp.Server.Chain.IntervalBlockHashes(
		&msg.StopHash, wire.CFCheckptInterval,
	)
	if err != nil {
		log.ERROR("invalid getcfilters request:", err)
		return
	}
	checkptMsg := wire.NewMsgCFCheckpt(
		msg.FilterType, &msg.StopHash, len(blockHashes),
	)
	// Fetch the current existing cache so we can decide if we need to extend it
	// or if its adequate as is.
	sp.Server.CFCheckptCachesMtx.RLock()
	checkptCache := sp.Server.CFCheckptCaches[msg.FilterType]
	// If the set of block hashes is beyond the current size of the cache, then
	// we'll expand the size of the cache and also retain the write lock.
	var updateCache bool
	if len(blockHashes) > len(checkptCache) {
		// Now that we know we'll need to modify the size of the cache, we'll
		// release the read lock and grab the write lock to possibly expand the
		// cache size.
		sp.Server.CFCheckptCachesMtx.RUnlock()
		sp.Server.CFCheckptCachesMtx.Lock()
		defer sp.Server.CFCheckptCachesMtx.Unlock()
		// Now that we have the write lock, we'll check again as it's possible
		// that the cache has already been expanded.
		checkptCache = sp.Server.CFCheckptCaches[msg.FilterType]
		// If we still need to expand the cache, then We'll mark that we need to
		// update the cache for below and also expand the size of the cache in
		// place.
		if len(blockHashes) > len(checkptCache) {
			updateCache = true
			additionalLength := len(blockHashes) - len(checkptCache)
			newEntries := make([]CFHeaderKV, additionalLength)
			log.INFOF("growing size of checkpoint cache from %v to %v block hashes",
				len(checkptCache), len(blockHashes))
			// nolint
			checkptCache = append(
				sp.Server.CFCheckptCaches[msg.FilterType],
				newEntries...,
			)
		}
	} else {
		// Otherwise, we'll hold onto the read lock for the remainder of this
		// method.
		defer sp.Server.CFCheckptCachesMtx.RUnlock()
		log.TRACEF("serving stale cache of size %v", len(checkptCache))
	}
	// Now that we know the cache is of an appropriate size, we'll iterate
	// backwards until the find the block hash. We do this as it's possible a
	// re-org has occurred so items in the db are now in the main china while
	// the cache has been partially invalidated.
	var forkIdx int
	for forkIdx = len(blockHashes); forkIdx > 0; forkIdx-- {
		if checkptCache[forkIdx-1].BlockHash == blockHashes[forkIdx-1] {
			break
		}
	}
	// Now that we know the how much of the cache is relevant for this query,
	// we'll populate our check point message with the cache as is. Shortly
	// below, we'll populate the new elements of the cache.
	for i := 0; i < forkIdx; i++ {
		err := checkptMsg.AddCFHeader(&checkptCache[i].FilterHeader)
		if err != nil {
			log.ERROR(err)
		}
	}
	// We'll now collect the set of hashes that are beyond our cache so we can
	// look up the filter headers to populate the final cache.
	blockHashPtrs := make([]*chainhash.Hash, 0, len(blockHashes)-forkIdx)
	for i := forkIdx; i < len(blockHashes); i++ {
		blockHashPtrs = append(blockHashPtrs, &blockHashes[i])
	}
	filterHeaders, err := sp.Server.CFIndex.FilterHeadersByBlockHashes(
		blockHashPtrs, msg.FilterType,
	)
	log.ERROR("error retrieving cfilter headers:", err)
	if err != nil {
		log.ERROR(err)
		return
	}
	// Now that we have the full set of filter headers, we'll add them to the
	// checkpoint message, and also update our cache in line.
	for i, filterHeaderBytes := range filterHeaders {
		if len(filterHeaderBytes) == 0 {
			log.WARN("could not obtain CF header for", blockHashPtrs[i])
			return
		}
		filterHeader, err := chainhash.NewHash(filterHeaderBytes)
		if err != nil {
			log.ERROR("committed filter header deserialize failed:", err)
			return
		}
		err = checkptMsg.AddCFHeader(filterHeader)
		if err != nil {
			log.ERROR(err)
		}
		// If the new main chain is longer than what's in the cache, then we'll
		// override it beyond the fork point.
		if updateCache {
			checkptCache[forkIdx+i] = CFHeaderKV{
				BlockHash:    blockHashes[forkIdx+i],
				FilterHeader: *filterHeader,
			}
		}
	}
	// Finally, we'll update the cache if we need to, and send the final message
	// back to the requesting peer.
	if updateCache {
		sp.Server.CFCheckptCaches[msg.FilterType] = checkptCache
	}
	sp.QueueMessage(checkptMsg, nil)
}

// OnGetCFHeaders is invoked when a peer receives a getcfheader bitcoin message.
func (sp *NodePeer) OnGetCFHeaders(_ *peer.Peer,
	msg *wire.MsgGetCFHeaders) {
	// Ignore getcfilterheader requests if not in sync.
	if !sp.Server.SyncManager.IsCurrent() {
		return
	}
	// We'll also ensure that the remote party is requesting a set of headers
	// for filters that we actually currently maintain.
	switch msg.FilterType {
	case wire.GCSFilterRegular:
		break
	default:
		log.DEBUG(
			"filter request for unknown headers for filter:",
			msg.FilterType)
		return
	}
	startHeight := int32(msg.StartHeight)
	maxResults := wire.MaxCFHeadersPerMsg
	// If StartHeight is positive, fetch the predecessor block hash so we can
	// populate the PrevFilterHeader field.
	if msg.StartHeight > 0 {
		startHeight--
		maxResults++
	}
	// Fetch the hashes from the block index.
	hashList, err := sp.Server.Chain.HeightToHashRange(
		startHeight, &msg.StopHash, maxResults,
	)
	if err != nil {
		log.ERROR("invalid getcfheaders request:", err)
	}
	// This is possible if StartHeight is one greater that the height of
	// StopHash, and we pull a valid range of hashes including the previous
	// filter header.
	if len(hashList) == 0 || (msg.StartHeight > 0 && len(hashList) == 1) {
		log.DEBUG("no results for getcfheaders request")
		return
	}
	// Create []*chainhash.Hash from []chainhash.Hash to pass to
	// FilterHeadersByBlockHashes.
	hashPtrs := make([]*chainhash.Hash, len(hashList))
	for i := range hashList {
		hashPtrs[i] = &hashList[i]
	}
	// Fetch the raw filter hash bytes from the database for all blocks.
	filterHashes, err := sp.Server.CFIndex.FilterHashesByBlockHashes(
		hashPtrs, msg.FilterType,
	)
	if err != nil {
		log.ERROR("error retrieving cfilter hashes:", err)
		return
	}
	// Generate cfheaders message and send it.
	headersMsg := wire.NewMsgCFHeaders()
	// Populate the PrevFilterHeader field.
	if msg.StartHeight > 0 {
		prevBlockHash := &hashList[0]
		// Fetch the raw committed filter header bytes from the database.
		headerBytes, err := sp.Server.CFIndex.FilterHeaderByBlockHash(
			prevBlockHash, msg.FilterType)
		if err != nil {
			log.ERROR("error retrieving CF header:", err)
			return
		}
		if len(headerBytes) == 0 {
			log.WARN("could not obtain CF header for", prevBlockHash)
			return
		}
		// Deserialize the hash into PrevFilterHeader.
		err = headersMsg.PrevFilterHeader.SetBytes(headerBytes)
		if err != nil {
			log.ERROR("committed filter header deserialize failed:", err)
			return
		}
		hashList = hashList[1:]
		filterHashes = filterHashes[1:]
	}
	// Populate HeaderHashes.
	for i, hashBytes := range filterHashes {
		if len(hashBytes) == 0 {
			log.WARN("could not obtain CF hash for", hashList[i])
			return
		}
		// Deserialize the hash.
		filterHash, err := chainhash.NewHash(hashBytes)
		if err != nil {
			log.ERROR("committed filter hash deserialize failed:", err)
			return
		}
		err = headersMsg.AddCFHash(filterHash)
		if err != nil {
			log.ERROR(err)
		}
	}
	headersMsg.FilterType = msg.FilterType
	headersMsg.StopHash = msg.StopHash
	sp.QueueMessage(headersMsg, nil)
}

// OnGetCFilters is invoked when a peer receives a getcfilters bitcoin message.
func (sp *NodePeer) OnGetCFilters(_ *peer.Peer,
	msg *wire.MsgGetCFilters) {
	// Ignore getcfilters requests if not in sync.
	if !sp.Server.SyncManager.IsCurrent() {
		return
	}
	// We'll also ensure that the remote party is requesting a set of filters
	// that we actually currently maintain.
	switch msg.FilterType {
	case wire.GCSFilterRegular:
		break
	default:
		log.DEBUG("filter request for unknown filter:", msg.FilterType)
		return
	}
	hashes, err := sp.Server.Chain.HeightToHashRange(
		int32(msg.StartHeight), &msg.StopHash, wire.MaxGetCFiltersReqRange,
	)
	if err != nil {
		log.ERROR(err)
		log.ERROR("invalid getcfilters request:", err)
		return
	}
	// Create []*chainhash.Hash from []chainhash.Hash to pass to
	// FiltersByBlockHashes.
	hashPtrs := make([]*chainhash.Hash, len(hashes))
	for i := range hashes {
		hashPtrs[i] = &hashes[i]
	}
	filters, err := sp.Server.CFIndex.FiltersByBlockHashes(
		hashPtrs, msg.FilterType,
	)
	if err != nil {
		log.ERROR(err)
		log.ERROR("error retrieving cfilters:", err)
		return
	}
	for i, filterBytes := range filters {
		if len(filterBytes) == 0 {
			log.WARN("could not obtain cfilter for", hashes[i])
			return
		}
		filterMsg := wire.NewMsgCFilter(
			msg.FilterType, &hashes[i], filterBytes,
		)
		sp.QueueMessage(filterMsg, nil)
	}
}

// handleGetData is invoked when a peer receives a getdata bitcoin message and
// is used to deliver block and transaction information.
func (sp *NodePeer) OnGetData(_ *peer.Peer,
	msg *wire.MsgGetData) {
	numAdded := 0
	notFound := wire.NewMsgNotFound()
	length := len(msg.InvList)
	// A decaying ban score increase is applied to prevent exhausting resources
	// with unusually large inventory queries. Requesting more than the maximum
	// inventory vector length within a short period of time yields a score
	// above the default ban threshold. Sustained bursts of small requests are
	// not penalized as that would potentially ban peers performing IBD. This
	// incremental score decays each minute to half of its value.
	sp.AddBanScore(0, uint32(length)*99/wire.MaxInvPerMsg, "getdata")
	// We wait on this wait channel periodically to prevent queuing far more
	// data than we can send in a reasonable time, wasting memory. The waiting
	// occurs after the database fetch for the next one to provide a little
	// pipelining.
	var waitChan chan struct{}
	doneChan := make(chan struct{}, 1)
	for i, iv := range msg.InvList {
		var c chan struct{}
		// If this will be the last message we send.
		if i == length-1 && len(notFound.InvList) == 0 {
			c = doneChan
		} else if (i+1)%3 == 0 {
			// Buffered so as to not make the send goroutine block.
			c = make(chan struct{}, 1)
		}
		var err error
		switch iv.Type {
		case wire.InvTypeWitnessTx:
			err = sp.Server.PushTxMsg(sp, &iv.Hash, c, waitChan,
				wire.WitnessEncoding)
		case wire.InvTypeTx:
			err = sp.Server.PushTxMsg(sp, &iv.Hash, c, waitChan,
				wire.BaseEncoding)
		case wire.InvTypeWitnessBlock:
			err = sp.Server.PushBlockMsg(sp, &iv.Hash, c, waitChan,
				wire.WitnessEncoding)
		case wire.InvTypeBlock:
			err = sp.Server.PushBlockMsg(sp, &iv.Hash, c, waitChan,
				wire.BaseEncoding)
		case wire.InvTypeFilteredWitnessBlock:
			err = sp.Server.PushMerkleBlockMsg(sp, &iv.Hash, c, waitChan,
				wire.WitnessEncoding)
		case wire.InvTypeFilteredBlock:
			err = sp.Server.PushMerkleBlockMsg(sp, &iv.Hash, c, waitChan,
				wire.BaseEncoding)
		default:
			log.WARN("unknown type in inventory request", iv.Type)
			continue
		}
		if err != nil {
			log.ERROR(err)
			err := notFound.AddInvVect(iv)
			if err != nil {
				log.ERROR(err)
			}
			// When there is a failure fetching the final entry and the done
			// channel was sent in due to there being no outstanding not found
			// inventory, consume it here because there is now not found inventory
			// that will use the channel momentarily.
			if i == len(msg.InvList)-1 && c != nil {
				<-c
			}
		}
		numAdded++
		waitChan = c
	}
	if len(notFound.InvList) != 0 {
		sp.QueueMessage(notFound, doneChan)
	}
	// Wait for messages to be sent. We can send quite a lot of data at this
	// point and this will keep the peer busy for a decent amount of time. We
	// don't process anything else by them in this time so that we have an idea
	// of when we should hear back from them - else the idle timeout could fire
	// when we were only half done sending the blocks.
	if numAdded > 0 {
		<-doneChan
	}
}

// OnGetHeaders is invoked when a peer receives a getheaders bitcoin message.
func (sp *NodePeer) OnGetHeaders(_ *peer.Peer,
	msg *wire.MsgGetHeaders) {
	// Ignore getheaders requests if not in sync.
	if !sp.Server.SyncManager.IsCurrent() {
		return
	}
	// Find the most recent known block in the best chain based on the block
	// locator and fetch all of the headers after it until either
	// wire.MaxBlockHeadersPerMsg have been fetched or the provided stop hash is
	// encountered. Use the block after the genesis block if no other blocks in
	// the provided locator are known.  This does mean the client will start
	// over with the genesis block if unknown block locators are provided. This
	// mirrors the behavior in the reference implementation.
	chain := sp.Server.Chain
	headers := chain.LocateHeaders(msg.BlockLocatorHashes, &msg.HashStop)
	// Send found headers to the requesting peer.
	blockHeaders := make([]*wire.BlockHeader, len(headers))
	for i := range headers {
		blockHeaders[i] = &headers[i]
	}
	sp.QueueMessage(&wire.MsgHeaders{Headers: blockHeaders}, nil)
}

// OnHeaders is invoked when a peer receives a headers bitcoin message.
// The message is passed down to the sync manager.
func (sp *NodePeer) OnHeaders(_ *peer.Peer,
	msg *wire.MsgHeaders) {
	sp.Server.SyncManager.QueueHeaders(msg, sp.Peer)
}

// OnInv is invoked when a peer receives an inv bitcoin message and is used to
// examine the inventory being advertised by the remote peer and react
// accordingly.  We pass the message down to blockmanager which will call
// QueueMessage with any appropriate responses.
func (sp *NodePeer) OnInv(
	_ *peer.Peer,
	msg *wire.MsgInv) {
	if !*sp.Server.Config.BlocksOnly {
		if len(msg.InvList) > 0 {
			sp.Server.SyncManager.QueueInv(msg, sp.Peer)
		}
		return
	}
	newInv := wire.NewMsgInvSizeHint(uint(len(msg.InvList)))
	for _, invVect := range msg.InvList {
		if invVect.Type == wire.InvTypeTx {
			log.TRACEF("ignoring tx %v in inv from %v -- blocksonly enabled",
				invVect.Hash, sp)
			if sp.ProtocolVersion() >= wire.BIP0037Version {
				log.INFOF("peer %v is announcing transactions"+
					" -- disconnecting", sp)
				sp.Disconnect()
				return
			}
			continue
		}
		err := newInv.AddInvVect(invVect)
		if err != nil {
			log.ERROR("failed to add inventory vector:", err)
			break
		}
	}
	if len(newInv.InvList) > 0 {
		sp.Server.SyncManager.QueueInv(newInv, sp.Peer)
	}
}

// OnMemPool is invoked when a peer receives a mempool bitcoin message. It
// creates and sends an inventory message with the contents of the memory pool
// up to the maximum inventory allowed per message.  When the peer has a bloom
// filter loaded, the contents are filtered accordingly.
func (sp *NodePeer) OnMemPool(_ *peer.Peer,
	msg *wire.MsgMemPool) {
	// Only allow mempool requests if the server has bloom filtering enabled.
	if sp.Server.Services&wire.SFNodeBloom != wire.SFNodeBloom {
		log.DEBUG("peer", sp, "sent mempool request with bloom filtering disabled"+
			" -- disconnecting")
		sp.Disconnect()
		return
	}
	// A decaying ban score increase is applied to prevent flooding. The ban
	// score accumulates and passes the ban threshold if a burst of mempool
	// messages comes from a peer. The score decays each minute to half of its
	// value.
	sp.AddBanScore(0, 33, "mempool")
	// Generate inventory message with the available transactions in the
	// transaction memory pool.  Limit it to the max allowed inventory per
	// message.  The NewMsgInvSizeHint function automatically limits the passed
	// hint to the maximum allowed, so it's safe to pass it without double
	// checking it here.
	txMemPool := sp.Server.TxMemPool
	txDescs := txMemPool.TxDescs()
	invMsg := wire.NewMsgInvSizeHint(uint(len(txDescs)))
	for _, txDesc := range txDescs {
		// Either add all transactions when there is no bloom filter, or only the
		// transactions that match the filter when there is one.
		if !sp.Filter.IsLoaded() || sp.Filter.MatchTxAndUpdate(txDesc.Tx) {
			iv := wire.NewInvVect(wire.InvTypeTx, txDesc.Tx.Hash())
			err := invMsg.AddInvVect(iv)
			if err != nil {
				log.ERROR(err)
			}
			if len(invMsg.InvList)+1 > wire.MaxInvPerMsg {
				break
			}
		}
	}
	// Send the inventory message if there is anything to send.
	if len(invMsg.InvList) > 0 {
		sp.QueueMessage(invMsg, nil)
	}
}

// OnRead is invoked when a peer receives a message and it is used to update
// the bytes received by the server.
func (sp *NodePeer) OnRead(_ *peer.Peer,
	bytesRead int, msg wire.Message, err error) {
	sp.Server.AddBytesReceived(uint64(bytesRead))
}

// OnTx is invoked when a peer receives a tx bitcoin message.  It blocks until
// the bitcoin transaction has been fully processed.  Unlock the block handler
// this does not serialize all transactions through a single thread
// transactions don't rely on the previous one in a linear fashion like blocks.
func (sp *NodePeer) OnTx(
	_ *peer.Peer,
	msg *wire.MsgTx) {
	if *sp.Server.Config.BlocksOnly {
		log.TRACEF("ignoring tx %v from %v - blocksonly enabled", msg.TxHash(), sp)
		return
	}
	// Add the transaction to the known inventory for the peer. Convert the raw
	// MsgTx to a util.Tx which provides some convenience methods and things
	// such as hash caching.
	tx := util.NewTx(msg)
	iv := wire.NewInvVect(wire.InvTypeTx, tx.Hash())
	sp.AddKnownInventory(iv)
	// Queue the transaction up to be handled by the sync manager and
	// intentionally block further receives until the transaction is fully
	// processed and known good or bad.  This helps prevent a malicious peer
	// from queuing up a bunch of bad transactions before disconnecting (or
	// being disconnected) and wasting memory.
	sp.Server.SyncManager.QueueTx(tx, sp.Peer, sp.TxProcessed)
	<-sp.TxProcessed
}

// OnVersion is invoked when a peer receives a version bitcoin message and is
// used to negotiate the protocol version details as well as kick start the
// communications.
func (sp *NodePeer) OnVersion(
	_ *peer.Peer,
	msg *wire.MsgVersion) *wire.MsgReject {
	// Update the address manager with the advertised services for outbound
	// connections in case they have changed.  This is not done for inbound
	// connections to help prevent malicious behavior and is skipped when
	// running on the simulation test network since it is only intended to
	// connect to specified peers and actively avoids advertising and connecting
	// to discovered peers. NOTE: This is done before rejecting peers that are
	// too old to ensure it is updated regardless in the case a new minimum
	// protocol version is enforced and the remote node has not upgraded yet.
	isInbound := sp.Inbound()
	remoteAddr := sp.NA()
	addrManager := sp.Server.AddrManager
	if !((*sp.Server.Config.Network)[0] == 's') && !isInbound {
		addrManager.SetServices(remoteAddr, msg.Services)
	}
	// Ignore peers that have a protcol version that is too old.  The peer
	// negotiation logic will disconnect it after this callback returns.
	if msg.ProtocolVersion < int32(peer.MinAcceptableProtocolVersion) {
		return nil
	}
	// Reject outbound peers that are not full nodes.
	wantServices := wire.SFNodeNetwork
	if !isInbound && !GetHasServices(msg.Services, wantServices) {
		missingServices := wantServices & ^msg.Services
		log.DEBUGF("rejecting peer %s with services %v due to not providing"+
			" desired services %v %s", sp.Peer, msg.Services, missingServices)
		reason := fmt.Sprintf("required services %#x not offered",
			uint64(missingServices))
		return wire.NewMsgReject(msg.Command(), wire.RejectNonstandard, reason)
	}
	// Update the address manager and request known addresses from the remote
	// peer for outbound connections.  This is skipped when running on the
	// simulation test network since it is only intended to connect to specified
	// peers and actively avoids advertising and connecting to discovered peers.
	if !((*sp.Server.Config.Network)[0] == 't') && !isInbound {
		// After soft-fork activation, only make outbound connection to peers if
		// they flag that they're segwit enabled.
		chain := sp.Server.Chain
		segwitActive, err := chain.IsDeploymentActive(chaincfg.DeploymentSegwit)
		if err != nil {
			log.ERROR(err)
			log.ERROR("unable to query for segwit soft-fork state:", err)
			return nil
		}
		if segwitActive && !sp.IsWitnessEnabled() {
			log.INFO("disconnecting non-segwit peer", sp,
				"as it isn't segwit enabled and we need more segwit enabled peers")
			sp.Disconnect()
			return nil
		}
		// Advertise the local address when the server accepts incoming
		// connections and it believes itself to be close to the best known tip.
		if !*sp.Server.Config.DisableListen && sp.Server.SyncManager.IsCurrent() {
			// Get address that best matches.
			lna := addrManager.GetBestLocalAddress(remoteAddr)
			if addrmgr.IsRoutable(lna) {
				// Filter addresses the peer already knows about.
				addresses := []*wire.NetAddress{lna}
				sp.PreparePushAddrMsg(addresses)
			}
		}
		// Request known addresses if the server address manager needs more and
		// the peer has a protocol version new enough to include a timestamp with
		// addresses.
		hasTimestamp := sp.ProtocolVersion() >= wire.NetAddressTimeVersion
		if addrManager.NeedMoreAddresses() && hasTimestamp {
			sp.QueueMessage(wire.NewMsgGetAddr(), nil)
		}
		// Mark the address as a known good address.
		addrManager.Good(remoteAddr)
	}
	// Add the remote peer time as a sample for creating an offset against the
	// local clock to keep the network time in sync.
	sp.Server.TimeSource.AddTimeSample(sp.Addr(), msg.Timestamp)
	// Signal the sync manager this peer is a new sync candidate.
	sp.Server.SyncManager.NewPeer(sp.Peer)
	// Choose whether or not to relay transactions before a filter command is
	// received.
	sp.SetDisableRelayTx(msg.DisableRelayTx)
	// Add valid peer to the server.
	sp.Server.AddPeer(sp)
	return nil
}

// OnWrite is invoked when a peer sends a message and it is used to update the
// bytes sent by the server.
func (sp *NodePeer) OnWrite(_ *peer.Peer, bytesWritten int,
	msg wire.Message, err error) {
	sp.Server.AddBytesSent(uint64(bytesWritten))
}

// AddBanScore increases the persistent and decaying ban score fields by the
// values passed as parameters. If the resulting score exceeds half of the ban
// threshold, a warning is logged including the reason provided. Further, if
// the score is above the ban threshold, the peer will be banned and
// disconnected.
func (sp *NodePeer) AddBanScore(persistent, transient uint32, reason string) {
	// No warning is logged and no score is calculated if banning is disabled.
	if *sp.Server.Config.DisableBanning {
		return
	}
	if sp.IsWhitelisted {
		log.DEBUGF("misbehaving whitelisted peer %s: %s %s", sp, reason)
		return
	}
	warnThreshold := *sp.Server.Config.BanThreshold >> 1
	if transient == 0 && persistent == 0 {
		// The score is not being increased, but a warning message is still
		// logged if the score is above the warn threshold.
		score := sp.BanScore.Int()
		if int(score) > warnThreshold {
			log.WARNF("misbehaving peer %s: %s -- ban score is %d, "+
				"it was not increased this time", sp, reason, score)
		}
		return
	}
	score := sp.BanScore.Increase(persistent, transient)
	if int(score) > warnThreshold {
		log.WARNF("misbehaving peer %s: %s -- ban score increased to %d",
			sp, reason, score)
		if int(score) > *sp.Server.Config.BanThreshold {
			log.WARNF("misbehaving peer %s -- banning and disconnecting", sp)
			sp.Server.BanPeer(sp)
			sp.Disconnect()
		}
	}
}

// AddKnownAddresses adds the given addresses to the set of known addresses to
// the peer to prevent sending duplicate addresses.
func (sp *NodePeer) AddKnownAddresses(addresses []*wire.NetAddress) {
	for _, na := range addresses {
		sp.KnownAddresses[addrmgr.NetAddressKey(na)] = struct{}{}
	}
}

// IsAddressKnown true if the given address is already known to the peer.
func (sp *NodePeer) IsAddressKnown(na *wire.NetAddress) bool {
	_, exists := sp.KnownAddresses[addrmgr.NetAddressKey(na)]
	return exists
}

// EnforceNodeBloomFlag disconnects the peer if the server is not configured to
// allow bloom filters.  Additionally, if the peer has negotiated to a protocol
// version  that is high enough to observe the bloom filter service support bit,
// it will be banned since it is intentionally violating the protocol.
func (sp *NodePeer) EnforceNodeBloomFlag(cmd string) bool {
	if sp.Server.Services&wire.SFNodeBloom != wire.SFNodeBloom {
		// Ban the peer if the protocol version is high enough that the peer is
		// knowingly violating the protocol and banning is enabled. NOTE: Even
		// though the addBanScore function already examines whether or not
		// banning is enabled, it is checked here as well to ensure the violation
		// is logged and the peer is disconnected regardless.
		if sp.ProtocolVersion() >= wire.BIP0111Version &&
			!*sp.Server.Config.DisableBanning {
			// Disconnect the peer regardless of whether it was banned.
			sp.AddBanScore(100, 0, cmd)
			sp.Disconnect()
			return false
		}
		// Disconnect the peer regardless of protocol version or banning state.
		log.DEBUGF("%s sent an unsupported %s request -- disconnecting %s", sp, cmd)
		sp.Disconnect()
		return false
	}
	return true
}

// GetNewestBlock returns the current best block hash and height using the
// format required by the configuration for the peer package.
func (sp *NodePeer) GetNewestBlock() (*chainhash.Hash, int32, error) {
	best := sp.Server.Chain.BestSnapshot()
	return &best.Hash, best.Height, nil
}

// PreparePushAddrMsg sends an addr message to the connected peer using the
// provided addresses.
func (sp *NodePeer) PreparePushAddrMsg(addresses []*wire.NetAddress) {
	// Filter addresses already known to the peer.
	addrs := make([]*wire.NetAddress, 0, len(addresses))
	for _, addr := range addresses {
		if !sp.IsAddressKnown(addr) {
			addrs = append(addrs, addr)
		}
	}
	known, err := sp.PushAddrMsg(addrs)
	if err != nil {
		log.ERRORF("can't push address message to %s: %v", sp.Peer, err)
		sp.Disconnect()
		return
	}
	sp.AddKnownAddresses(known)
}

// IsRelayTxDisabled returns whether or not relaying of transactions for the
// given peer is disabled. It is safe for concurrent access.
func (sp *NodePeer) IsRelayTxDisabled() bool {
	sp.RelayMtx.Lock()
	isDisabled := sp.DisableRelayTx
	sp.RelayMtx.Unlock()
	return isDisabled
}

// SetDisableRelayTx toggles relaying of transactions for the given peer. It is
// safe for concurrent access.
func (sp *NodePeer) SetDisableRelayTx(disable bool) {
	sp.RelayMtx.Lock()
	sp.DisableRelayTx = disable
	sp.RelayMtx.Unlock()
}

// Len returns the number of checkpoints in the slice.  It is part of the
// sort.Interface implementation.
func (s CheckpointSorter) Len() int { return len(s) }

//	Less returns whether the checkpoint with index i should sort before the
// checkpoint with index j.  It is part of the sort.Interface implementation.
func (s CheckpointSorter) Less(i, j int) bool {
	return s[i].Height < s[j].
		Height
}

// Swap swaps the checkpoints at the passed indices.  It is part of the
// sort.Interface implementation.
func (s CheckpointSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Network returns the network. This is part of the net.Addr interface.
func (a SimpleAddr) Network() string {
	return a.Net
}

// String returns the address. This is part of the net.Addr interface.
func (a SimpleAddr) String() string {
	return a.Addr
}

//	AddLocalAddress adds an address that this node is listening on to the
// address manager so that it may be relayed to peers.
func AddLocalAddress(addrMgr *addrmgr.AddrManager, addr string, services wire.ServiceFlag) error {
	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		log.ERROR(err)
		return err
	}
	port, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		log.ERROR(err)
		return err
	}
	if ip := net.ParseIP(host); ip != nil && ip.IsUnspecified() {
		// If bound to unspecified address, advertise all local interfaces
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			log.ERROR(err)
			return err
		}
		for _, addr := range addrs {
			ifaceIP, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				log.ERROR(err)
				continue
			}
			//	If bound to 0.0.0.0, do not add IPv6 interfaces and if bound to
			// ::, do not add IPv4 interfaces.
			if (ip.To4() == nil) != (ifaceIP.To4() == nil) {
				continue
			}
			netAddr := wire.NewNetAddressIPPort(ifaceIP, uint16(port), services)
			err = addrMgr.AddLocalAddress(netAddr, addrmgr.BoundPrio)
			if err != nil {
				log.TRACE(err)
			}
		}
	} else {
		netAddr, err := addrMgr.HostToNetAddress(host, uint16(port), services)
		if err != nil {
			log.ERROR(err)
			return err
		}
		err = addrMgr.AddLocalAddress(netAddr, addrmgr.BoundPrio)
		if err != nil {
			log.ERROR(err)
		}
	}
	return nil
}

// AddrStringToNetAddr takes an address in the form of 'host:port' and returns
// a net.Addr which maps to the original address with any host names resolved
// to IP addresses.  It also handles tor addresses properly by returning a
// net.Addr that encapsulates the address.
func AddrStringToNetAddr(config *pod.Config, stateCfg *state.Config, addr string) (net.Addr, error) {
	host, strPort, err := net.SplitHostPort(addr)
	if err != nil {
		log.ERROR(err)
		return nil, err
	}
	port, err := strconv.Atoi(strPort)
	if err != nil {
		log.ERROR(err)
		return nil, err
	}
	// Skip if host is already an IP address.
	if ip := net.ParseIP(host); ip != nil {
		return &net.TCPAddr{
				IP:   ip,
				Port: port,
			},
			nil
	}
	// Tor addresses cannot be resolved to an IP, so just return an onion
	// address instead.
	if strings.HasSuffix(host, ".onion") {
		if !*config.Onion {
			return nil, errors.New("tor has been disabled")
		}
		return &OnionAddr{Addr: addr}, nil
	}
	// Attempt to look up an IP address associated with the parsed host.
	ips, err := Lookup(stateCfg)(host)
	if err != nil {
		log.ERROR(err)
		return nil, err
	}
	if len(ips) == 0 {
		return nil, fmt.Errorf("no addresses found for %s", host)
	}
	return &net.TCPAddr{
			IP:   ips[0],
			Port: port,
		},
		nil
}

// DisconnectPeer attempts to drop the connection of a targeted peer in the
// passed peer list. Targets are identified via usage of the passed
// `compareFunc`, which should return `true` if the passed peer is the target
// peer. This function returns true on success and false if the peer is unable
// to be located. If the peer is found, and the passed callback: `whenFound'
// isn't nil, we call it with the peer as the argument before it is removed
// from the peerList, and is disconnected from the server.
func DisconnectPeer(peerList map[int32]*NodePeer,
	compareFunc func(*NodePeer) bool, whenFound func(*NodePeer)) bool {
	for addr, nodePeer := range peerList {
		if compareFunc(nodePeer) {
			if whenFound != nil {
				whenFound(nodePeer)
			}
			// This is ok because we are not continuing to iterate so won't
			// corrupt the loop.
			delete(peerList, addr)
			nodePeer.Disconnect()
			return true
		}
	}
	return false
}

//	DynamicTickDuration is a convenience function used to dynamically choose a
// tick duration based on remaining time.  It is primarily used during
// server shutdown to make shutdown warnings more frequent as the shutdown time
// approaches.
func DynamicTickDuration(remaining time.Duration) time.Duration {
	switch {
	case remaining <= time.Second*5:
		return time.Second
	case remaining <= time.Second*15:
		return time.Second * 5
	case remaining <= time.Minute:
		return time.Second * 15
	case remaining <= time.Minute*5:
		return time.Minute
	case remaining <= time.Minute*15:
		return time.Minute * 5
	case remaining <= time.Hour:
		return time.Minute * 15
	}
	return time.Hour
}

// GetHasServices returns whether or not the provided advertised service flags
// have all of the provided desired service flags set.
func GetHasServices(advertised, desired wire.ServiceFlag) bool {
	return advertised&desired == desired
}

// InitListeners initializes the configured net listeners and adds any bound
// addresses to the address manager. Returns the listeners and a upnp.NAT
// interface, which is non-nil if UPnP is in use.
func InitListeners(config *pod.Config, activeNet *netparams.Params,
	aMgr *addrmgr.AddrManager, listenAddrs []string, services wire.ServiceFlag) ([]net.Listener, upnp.NAT, error) {
	// Listen for TCP connections at the configured addresses
	log.TRACE("listenAddrs ", listenAddrs)
	netAddrs, err := ParseListeners(listenAddrs)
	if err != nil {
		log.ERROR(err)
		return nil, nil, err
	}
	log.TRACE("netAddrs ", netAddrs)
	listeners := make([]net.Listener, 0, len(netAddrs))
	for _, addr := range netAddrs {
		log.TRACE("addr ", addr, " ", addr.Network(), " ", addr.String())
		listener, err := net.Listen(addr.Network(), addr.String())
		if err != nil {
			log.WARNF("can't listen on %s: %v %s", addr, err)
			continue
		}
		listeners = append(listeners, listener)
	}
	var nat upnp.NAT
	if len(*config.ExternalIPs) != 0 {
		defaultPort, err := strconv.ParseUint(activeNet.DefaultPort, 10, 16)
		if err != nil {
			log.ERRORF("can not parse default port %s for active chain: %v",
				activeNet.DefaultPort, err)
			return nil, nil, err
		}
		for _, sip := range *config.ExternalIPs {
			eport := uint16(defaultPort)
			host, portstr, err := net.SplitHostPort(sip)
			if err != nil {
				// no port, use default.
				host = sip
			} else {
				port, err := strconv.ParseUint(portstr, 10, 16)
				if err != nil {
					log.ERRORF("can not parse port from %s for externalip: %v",
						sip, err)
					continue
				}
				eport = uint16(port)
			}
			na, err := aMgr.HostToNetAddress(host, eport, services)
			if err != nil {
				log.ERRORF("not adding %s as externalip: %v", sip, err)
				continue
			}
			err = aMgr.AddLocalAddress(na, addrmgr.ManualPrio)
			if err != nil {
				log.ERRORF("skipping specified external IP: %v", err)
			}
		}
	} else {
		if *config.UPNP {
			var err error
			nat, err = upnp.Discover()
			if err != nil {
				log.ERRORF("can't discover upnp: %v", err)
			}
			// nil upnp.nat here is fine, just means no upnp on network.
		}
		// Add bound addresses to address manager to be advertised to peers.
		for _, listener := range listeners {
			addr := listener.Addr().String()
			err := AddLocalAddress(aMgr, addr, services)
			if err != nil {
				log.ERRORF("skipping bound address %s: %v", addr, err)
			}
		}
	}
	return listeners, nat, nil
}

// GetIsWhitelisted returns whether the IP address is included in the
// whitelisted networks and IPs.
func GetIsWhitelisted(statecfg *state.Config, addr net.Addr) bool {
	if len(statecfg.ActiveWhitelists) == 0 {
		return false
	}
	host, _, err := net.SplitHostPort(addr.String())
	if err != nil {
		log.ERROR(err)
		log.ERRORF("unable to SplitHostPort on '%s': %v", addr, err)
		return false
	}
	ip := net.ParseIP(host)
	if ip == nil {
		log.WARNF("unable to parse IP '%s'", addr)
		return false
	}
	for _, ipnet := range statecfg.ActiveWhitelists {
		if ipnet.Contains(ip) {
			return true
		}
	}
	return false
}

func // MergeCheckpoints returns two slices of checkpoints merged into one slice
// such that the checkpoints are sorted by height.  In the case the additional
// checkpoints contain a checkpoint with the same height as a checkpoint in the
// default checkpoints, the additional checkpoint will take precedence and
// overwrite the default one.
MergeCheckpoints(defaultCheckpoints, additional []chaincfg.Checkpoint) []chaincfg.Checkpoint {
	//	Create a map of the additional checkpoints to remove duplicates while
	//	leaving the most recently-specified checkpoint.
	extra := make(map[int32]chaincfg.Checkpoint)
	for _, checkpoint := range additional {
		extra[checkpoint.Height] = checkpoint
	}
	// Add all default checkpoints that do not have an override in the
	// additional checkpoints.
	numDefault := len(defaultCheckpoints)
	checkpoints := make([]chaincfg.Checkpoint, 0, numDefault+len(extra))
	for _, checkpoint := range defaultCheckpoints {
		if _, exists := extra[checkpoint.Height]; !exists {
			checkpoints = append(checkpoints, checkpoint)
		}
	}
	// Append the additional checkpoints and return the sorted results.
	for _, checkpoint := range extra {
		checkpoints = append(checkpoints, checkpoint)
	}
	sort.Sort(CheckpointSorter(checkpoints))
	return checkpoints
}

func // NewPeerConfig returns the configuration for the given ServerPeer.
NewPeerConfig(sp *NodePeer) *peer.Config {
	return &peer.Config{
		Listeners: peer.MessageListeners{
			OnVersion:      sp.OnVersion,
			OnMemPool:      sp.OnMemPool,
			OnTx:           sp.OnTx,
			OnBlock:        sp.OnBlock,
			OnInv:          sp.OnInv,
			OnHeaders:      sp.OnHeaders,
			OnGetData:      sp.OnGetData,
			OnGetBlocks:    sp.OnGetBlocks,
			OnGetHeaders:   sp.OnGetHeaders,
			OnGetCFilters:  sp.OnGetCFilters,
			OnGetCFHeaders: sp.OnGetCFHeaders,
			OnGetCFCheckpt: sp.OnGetCFCheckpt,
			OnFeeFilter:    sp.OnFeeFilter,
			OnFilterAdd:    sp.OnFilterAdd,
			OnFilterClear:  sp.OnFilterClear,
			OnFilterLoad:   sp.OnFilterLoad,
			OnGetAddr:      sp.OnGetAddr,
			OnAddr:         sp.OnAddr,
			OnRead:         sp.OnRead,
			OnWrite:        sp.OnWrite,
			// Note: The reference client currently bans peers that send alerts
			// not signed with its key.  We could verify against their key, but
			// since the reference client is currently unwilling to support other
			// implementations' alert messages, we will not relay theirs.
			OnAlert: nil,
		},
		NewestBlock:       sp.GetNewestBlock,
		HostToNetAddress:  sp.Server.AddrManager.HostToNetAddress,
		Proxy:             *sp.Server.Config.Proxy,
		UserAgentName:     UserAgentName,
		UserAgentVersion:  UserAgentVersion,
		UserAgentComments: *sp.Server.Config.UserAgentComments,
		ChainParams:       sp.Server.ChainParams,
		Services:          sp.Server.Services,
		DisableRelayTx:    *sp.Server.Config.BlocksOnly,
		ProtocolVersion:   peer.MaxProtocolVersion,
		TrickleInterval:   *sp.Server.Config.TrickleInterval,
	}
}

func // NewNode returns a new pod server configured to listen on addr for the
// bitcoin network type specified by chainParams.  Use start to begin accepting
// connections from peers.
// TODO: simplify/modularise this
NewNode(config *pod.Config, stateCfg *state.Config,
	activeNet *netparams.Params, listenAddrs []string, db database.DB,
	chainParams *netparams.Params, interruptChan <-chan struct{},
	algo string) (*Node, error) {
	log.TRACE("listenAddrs ", listenAddrs)
	services := DefaultServices
	if *config.NoPeerBloomFilters {
		services &^= wire.SFNodeBloom
	}
	if *config.NoCFilters {
		services &^= wire.SFNodeCF
	}
	aMgr := addrmgr.New(*config.DataDir+string(os.PathSeparator)+activeNet.
		Name, Lookup(stateCfg))
	var listeners []net.Listener
	var nat upnp.NAT
	if !*config.DisableListen {
		var err error
		listeners, nat, err = InitListeners(config, activeNet, aMgr,
			listenAddrs, services)
		if err != nil {
			log.ERROR(err)
			return nil, err
		}
		if len(listeners) == 0 {
			return nil, errors.New("no valid listen address")
		}
	}
	nThreads := runtime.NumCPU()
	var thr int
	if *config.GenThreads == -1 || thr > nThreads {
		thr = nThreads
	} else {
		thr = *config.GenThreads
	}
	log.TRACE("set genthreads to ", nThreads)
	s := Node{
		ChainParams:          chainParams,
		AddrManager:          aMgr,
		NewPeers:             make(chan *NodePeer, *config.MaxPeers),
		DonePeers:            make(chan *NodePeer, *config.MaxPeers),
		BanPeers:             make(chan *NodePeer, *config.MaxPeers),
		Query:                make(chan interface{}),
		RelayInv:             make(chan RelayMsg, *config.MaxPeers),
		Broadcast:            make(chan BroadcastMsg, *config.MaxPeers),
		Quit:                 make(chan struct{}),
		ModifyRebroadcastInv: make(chan interface{}),
		PeerHeightsUpdate:    make(chan UpdatePeerHeightsMsg),
		NAT:                  nat,
		DB:                   db,
		TimeSource:           blockchain.NewMedianTime(),
		Services:             services,
		SigCache:             txscript.NewSigCache(uint(*config.SigCacheMaxSize)),
		HashCache:            txscript.NewHashCache(uint(*config.SigCacheMaxSize)),
		CFCheckptCaches:      make(map[wire.FilterType][]CFHeaderKV),
		GenThreads:           uint32(thr),
		Algo:                 algo,
		Config:               config,
		StateCfg:             stateCfg,
		ActiveNet:            activeNet,
	}
	// Create the transaction and address indexes if needed.
	// CAUTION: the txindex needs to be first in the indexes array because the
	// addrindex uses data from the txindex during catchup.
	// If the addrindex is run first,
	// it may not have the transactions from the current block indexed.
	var indexes []indexers.Indexer
	log.DEBUG("txindex", *config.TxIndex, "addrindex", *config.AddrIndex)
	if *config.TxIndex || *config.AddrIndex {
		// Enable transaction index if address index is enabled since it
		// requires it.
		if !*config.TxIndex {
			log.INFO("transaction index enabled because it is required by the" +
				" address index")
			*config.TxIndex = true
		} else {
			log.INFO("transaction index is enabled")
		}
		s.TxIndex = indexers.NewTxIndex(db)
		indexes = append(indexes, s.TxIndex)
	}
	if *config.AddrIndex {
		log.INFO("address index is enabled")
		s.AddrIndex = indexers.NewAddrIndex(db, chainParams)
		indexes = append(indexes, s.AddrIndex)
	}
	if !*config.NoCFilters {
		log.TRACE("committed filter index is enabled")
		s.CFIndex = indexers.NewCfIndex(db, chainParams)
		indexes = append(indexes, s.CFIndex)
	}
	// Create an index manager if any of the optional indexes are enabled.
	var indexManager blockchain.IndexManager
	if len(indexes) > 0 {
		indexManager = indexers.NewManager(db, indexes)
	}
	// Merge given checkpoints with the default ones unless they are disabled.
	var checkpoints []chaincfg.Checkpoint
	if !*config.DisableCheckpoints {
		checkpoints = MergeCheckpoints(
			s.ChainParams.Checkpoints, stateCfg.AddedCheckpoints)
	}
	// Create a new block chain instance with the appropriate configuration.
	var err error
	s.Chain, err = blockchain.New(
		&blockchain.Config{
			DB:           s.DB,
			Interrupt:    interruptChan,
			ChainParams:  s.ChainParams,
			Checkpoints:  checkpoints,
			TimeSource:   s.TimeSource,
			SigCache:     s.SigCache,
			IndexManager: indexManager,
			HashCache:    s.HashCache,
		},
	)
	if err != nil {
		log.ERROR(err)
		return nil, err
	}
	s.Chain.DifficultyAdjustments = make(map[string]float64)
	// Search for a FeeEstimator state in the database.
	// If none can be found or if it cannot be loaded, create a new one.
	e := db.Update(func(tx database.Tx) error {
		metadata := tx.Metadata()
		feeEstimationData := metadata.Get(mempool.EstimateFeeDatabaseKey)
		if feeEstimationData != nil {
			// delete it from the database so that we don't try to restore
			// the same thing again somehow.
			e := metadata.Delete(mempool.EstimateFeeDatabaseKey)
			if e != nil {
				return e
			}
			// If there is an error, log it and make a new fee estimator.
			var err error
			s.FeeEstimator, err = mempool.RestoreFeeEstimator(feeEstimationData)
			if err != nil {
				log.ERROR(err)
				return fmt.Errorf("failed to restore fee estimator %v", err)
			}
		}
		return nil
	})
	if e != nil {
		log.ERROR(e)
	}
	// If no feeEstimator has been found, or if the one that has been found is
	// behind somehow, create a new one and start over.
	if s.FeeEstimator == nil || s.FeeEstimator.LastKnownHeight() != s.Chain.
		BestSnapshot().Height {
		s.FeeEstimator = mempool.NewFeeEstimator(
			mempool.DefaultEstimateFeeMaxRollback,
			mempool.DefaultEstimateFeeMinRegisteredBlocks,
		)
	}
	txC := mempool.Config{
		Policy: mempool.Policy{
			DisableRelayPriority: *config.NoRelayPriority,
			AcceptNonStd:         *config.RelayNonStd,
			FreeTxRelayLimit:     *config.FreeTxRelayLimit,
			MaxOrphanTxs:         *config.MaxOrphanTxs,
			MaxOrphanTxSize:      DefaultMaxOrphanTxSize,
			MaxSigOpCostPerTx:    blockchain.MaxBlockSigOpsCost / 4,
			MinRelayTxFee:        stateCfg.ActiveMinRelayTxFee,
			MaxTxVersion:         2,
		},
		ChainParams:   chainParams,
		FetchUtxoView: s.Chain.FetchUtxoView,
		BestHeight: func() int32 {
			return s.Chain.BestSnapshot().Height
		},
		MedianTimePast: func() time.Time {
			return s.Chain.BestSnapshot().MedianTime
		},
		CalcSequenceLock: func(tx *util.Tx, view *blockchain.UtxoViewpoint) (
			*blockchain.SequenceLock, error) {
			return s.Chain.CalcSequenceLock(tx, view, true)
		},
		IsDeploymentActive: s.Chain.IsDeploymentActive,
		SigCache:           s.SigCache,
		HashCache:          s.HashCache,
		AddrIndex:          s.AddrIndex,
		FeeEstimator:       s.FeeEstimator,
	}
	s.TxMemPool = mempool.New(&txC)
	s.SyncManager, err =
		netsync.New(
			&netsync.Config{
				PeerNotifier:       &s,
				Chain:              s.Chain,
				TxMemPool:          s.TxMemPool,
				ChainParams:        s.ChainParams,
				DisableCheckpoints: *config.DisableCheckpoints,
				MaxPeers:           *config.MaxPeers,
				FeeEstimator:       s.FeeEstimator,
			},
		)
	if err != nil {
		log.ERROR(err)
		return nil, err
	}
	// Create the mining policy and block template generator based on the
	// configuration options.
	// NOTE: The CPU miner relies on the mempool, so the mempool has to be
	// created before calling the function to create the CPU miner.
	policy := mining.Policy{
		BlockMinWeight:    uint32(*config.BlockMinWeight),
		BlockMaxWeight:    uint32(*config.BlockMaxWeight),
		BlockMinSize:      uint32(*config.BlockMinSize),
		BlockMaxSize:      uint32(*config.BlockMaxSize),
		BlockPrioritySize: uint32(*config.BlockPrioritySize),
		TxMinFreeFee:      stateCfg.ActiveMinRelayTxFee,
	}
	blockTemplateGenerator := mining.NewBlkTmplGenerator(&policy,
		s.ChainParams, s.TxMemPool, s.Chain, s.TimeSource,
		s.SigCache, s.HashCache, s.Algo)
	s.CPUMiner = cpuminer.New(&cpuminer.Config{
		Blockchain:             s.Chain,
		ChainParams:            chainParams,
		BlockTemplateGenerator: blockTemplateGenerator,
		MiningAddrs:            stateCfg.ActiveMiningAddrs,
		ProcessBlock:           s.SyncManager.ProcessBlock,
		ConnectedCount:         s.ConnectedCount,
		IsCurrent:              s.SyncManager.IsCurrent,
		NumThreads:             s.GenThreads,
		Algo:                   s.Algo,
		Solo:                   *config.Solo,
	})
	// Only setup a function to return new addresses to connect to when
	// not running in connect-only mode.  The simulation network is always
	// in connect-only mode since it is only intended to connect to
	// specified peers and actively avoid advertising and connecting to
	// discovered peers in order to prevent it from becoming a public test
	// network.
	var newAddressFunc func() (net.Addr, error)
	if !((*config.Network)[0] == 's') && len(*config.ConnectPeers) == 0 {
		newAddressFunc = func() (net.Addr, error) {
			for tries := 0; tries < 100; tries++ {
				addr := s.AddrManager.GetAddress()
				if addr == nil {
					break
				}
				// Address will not be invalid, local or unrouteable
				// because addrmanager rejects those on addition.
				// Just check that we don't already have an address
				// in the same group so that we are not connecting
				// to the same network segment at the expense of
				// others.
				key := addrmgr.GroupKey(addr.NetAddress())
				if s.OutboundGroupCount(key) != 0 {
					continue
				}
				// only allow recent nodes (10 min) after we failed 30 times
				if tries < 30 && time.Since(addr.LastAttempt()) < 10*time.Minute {
					continue
				}
				// allow non default ports after 50 failed tries.
				if tries < 50 && fmt.Sprintf("%d", addr.NetAddress().Port) !=
					activeNet.DefaultPort {
					continue
				}
				addrString := addrmgr.NetAddressKey(addr.NetAddress())
				return AddrStringToNetAddr(config, stateCfg, addrString)
			}
			return nil, errors.New("no valid connect address")
		}
	}
	// Create a connection manager.
	targetOutbound := DefaultTargetOutbound
	if *config.MaxPeers < targetOutbound {
		targetOutbound = *config.MaxPeers
	}
	cMgr, err :=
		connmgr.New(
			&connmgr.Config{
				Listeners:      listeners,
				OnAccept:       s.InboundPeerConnected,
				RetryDuration:  ConnectionRetryInterval,
				TargetOutbound: uint32(targetOutbound),
				Dial:           Dial(stateCfg),
				OnConnection:   s.OutboundPeerConnected,
				GetNewAddress:  newAddressFunc,
			},
		)
	if err != nil {
		log.ERROR(err)
		return nil, err
	}
	s.ConnManager = cMgr
	// Start up persistent peers.
	permanentPeers := *config.ConnectPeers
	if len(permanentPeers) == 0 {
		permanentPeers = *config.AddPeers
	}
	for _, addr := range permanentPeers {
		netAddr, err := AddrStringToNetAddr(config, stateCfg, addr)
		if err != nil {
			log.ERROR(err)
			return nil, err
		}
		go s.ConnManager.Connect(
			&connmgr.ConnReq{
				Addr:      netAddr,
				Permanent: true,
			},
		)
	}
	if !*config.DisableRPC {
		// Setup listeners for the configured RPC listen addresses and
		// TLS settings.
		listeners := map[string][]string{
			"sha256d": *config.RPCListeners,
		}
		for l := range listeners {
			rpcListeners, err := SetupRPCListeners(config, listeners[l])
			if err != nil {
				log.ERROR(err)
				return nil, err
			}
			if len(rpcListeners) == 0 {
				return nil, errors.New("RPCS: No valid listen address")
			}
			
			rp, err := NewRPCServer(&ServerConfig{
				Listeners:    rpcListeners,
				StartupTime:  s.StartupTime,
				ConnMgr:      &ConnManager{&s},
				SyncMgr:      &SyncManager{&s, s.SyncManager},
				TimeSource:   s.TimeSource,
				Chain:        s.Chain,
				ChainParams:  chainParams,
				DB:           db,
				TxMemPool:    s.TxMemPool,
				Generator:    blockTemplateGenerator,
				CPUMiner:     s.CPUMiner,
				TxIndex:      s.TxIndex,
				AddrIndex:    s.AddrIndex,
				CfIndex:      s.CFIndex,
				FeeEstimator: s.FeeEstimator,
				Algo:         l,
			}, stateCfg, config)
			if err != nil {
				log.ERROR(err)
				return nil, err
			}
			s.RPCServers = append(s.RPCServers, rp)
		}
		// Signal process shutdown when the RPC server requests it.
		go func() {
			for i := range s.RPCServers {
				<-s.RPCServers[i].RequestedProcessShutdown()
			}
			// interrupt.Request()
		}()
	}
	return &s, nil
}

// NewServerPeer returns a new ServerPeer instance. The peer needs to be set by
// the caller.
func NewServerPeer(s *Node, isPersistent bool) *NodePeer {
	return &NodePeer{
		Server:         s,
		Persistent:     isPersistent,
		Filter:         bloom.LoadFilter(nil),
		KnownAddresses: make(map[string]struct{}),
		Quit:           make(chan struct{}),
		TxProcessed:    make(chan struct{}, 1),
		BlockProcessed: make(chan struct{}, 1),
	}
}

// ParseListeners determines whether each listen address is IPv4 and IPv6 and
// returns a slice of appropriate net.Addrs to listen on with TCP. It also
// properly detects addresses which apply to "all interfaces" and adds the
// address as both IPv4 and IPv6.
func ParseListeners(addrs []string) ([]net.Addr, error) {
	netAddrs := make([]net.Addr, 0, len(addrs)*2)
	for _, addr := range addrs {
		host, _, err := net.SplitHostPort(addr)
		if err != nil {
			log.ERROR(err)
			// Shouldn't happen due to already being normalized.
			return nil, err
		}
		// Empty host or host of * on plan9 is both IPv4 and IPv6.
		if host == "" || (host == "*" && runtime.GOOS == "plan9") {
			netAddrs = append(netAddrs, SimpleAddr{Net: "tcp4", Addr: addr})
			netAddrs = append(netAddrs, SimpleAddr{Net: "tcp6", Addr: addr})
			continue
		}
		// Strip IPv6 zone id if present since net.ParseIP does not handle it.
		zoneIndex := strings.LastIndex(host, "%")
		if zoneIndex > 0 {
			host = host[:zoneIndex]
		}
		// Parse the IP.
		ip := net.ParseIP(host)
		if ip == nil {
			return nil, fmt.Errorf("'%s' is not a valid IP address", host)
		}
		// To4 returns nil when the IP is not an IPv4 address,
		// so use this determine the address type.
		if ip.To4() == nil {
			netAddrs = append(netAddrs, SimpleAddr{Net: "tcp6", Addr: addr})
		} else {
			netAddrs = append(netAddrs, SimpleAddr{Net: "tcp4", Addr: addr})
		}
	}
	return netAddrs, nil
}

// RandomUint16Number returns a random uint16 in a specified input range.  Note
// that the range is in zeroth ordering; if you pass it 1800, you will get
// values from 0 to 1800.
func RandomUint16Number(max uint16) uint16 {
	// In order to avoid modulo bias and ensure every possible outcome in [0,
	// max) has equal probability, the random number must be sampled from a
	// random source that has a range limited to a multiple of the modulus.
	var randomNumber uint16
	var limitRange = (math.MaxUint16 / max) * max
	for {
		err := binary.Read(rand.Reader, binary.LittleEndian, &randomNumber)
		if err != nil {
			log.ERROR(err)
		}
		if randomNumber < limitRange {
			return randomNumber % max
		}
	}
}

// SetupRPCListeners returns a slice of listeners that are configured for use
// with the RPC server depending on the configuration settings for listen
// addresses and TLS.
func SetupRPCListeners(config *pod.Config, urls []string) ([]net.Listener, error) {
	// Setup TLS if not disabled.
	listenFunc := net.Listen
	if *config.TLS {
		// Generate the TLS cert and key file if both don't already exist.
		if !FileExists(*config.RPCKey) && !FileExists(*config.RPCCert) {
			err := GenCertPair(*config.RPCCert, *config.RPCKey)
			if err != nil {
				log.ERROR(err)
				return nil, err
			}
		}
		keyPair, err := tls.LoadX509KeyPair(*config.RPCCert, *config.RPCKey)
		if err != nil {
			log.ERROR(err)
			return nil, err
		}
		tlsConfig := tls.Config{
			Certificates:       []tls.Certificate{keyPair},
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: *config.TLSSkipVerify,
		}
		// Change the standard net.Listen function to the tls one.
		listenFunc = func(net string, laddr string) (net.Listener, error) {
			return tls.Listen(net, laddr, &tlsConfig)
		}
	}
	netAddrs, err := ParseListeners(urls)
	if err != nil {
		log.ERROR(err)
		return nil, err
	}
	listeners := make([]net.Listener, 0, len(netAddrs))
	for _, addr := range netAddrs {
		listener, err := listenFunc(addr.Network(), addr.String())
		if err != nil {
			log.ERROR(err)
			log.ERRORF("can't listen on %s: %v", addr, err)
			continue
		}
		listeners = append(listeners, listener)
	}
	return listeners, nil
}

// FileExists reports whether the named file or directory exists.
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
