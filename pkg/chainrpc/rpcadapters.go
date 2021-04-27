package chainrpc

import (
	"sync/atomic"

	"github.com/p9c/matrjoska/pkg/block"

	"github.com/p9c/matrjoska/pkg/blockchain"
	"github.com/p9c/matrjoska/pkg/chainhash"
	"github.com/p9c/matrjoska/pkg/mempool"
	"github.com/p9c/matrjoska/pkg/netsync"
	"github.com/p9c/matrjoska/pkg/peer"
	"github.com/p9c/matrjoska/pkg/wire"
)

// Peer provides a peer for use with the RPC Server and implements the RPCServerPeer interface.
type Peer NodePeer

// Ensure rpcPeer implements the RPCServerPeer interface.
var _ ServerPeer = (*Peer)(nil)

// ToPeer returns the underlying peer instance.
//
// This function is safe for concurrent access and is part of the RPCServerPeer interface implementation.
func (p *Peer) ToPeer() *peer.Peer {
	if p == nil {
		return nil
	}
	return (*NodePeer)(p).Peer
}

// IsTxRelayDisabled returns whether or not the peer has disabled transaction relay.
//
// This function is safe for concurrent access and is part of the RPCServerPeer interface implementation.
func (p *Peer) IsTxRelayDisabled() bool {
	return (*NodePeer)(p).DisableRelayTx
}

// GetBanScore returns the current integer value that represents how close the peer is to being banned.
//
// This function is safe for concurrent access and is part of the RPCServerPeer interface implementation.
func (p *Peer) GetBanScore() uint32 {
	return (*NodePeer)(p).BanScore.Int()
}

// GetFeeFilter returns the requested current minimum fee rate for which transactions should be announced.
//
// This function is safe for concurrent access and is part of the RPCServerPeer interface implementation.
func (p *Peer) GetFeeFilter() int64 {
	return atomic.LoadInt64(&(*NodePeer)(p).FeeFilter)
}

// ConnManager provides a connection manager for use with the RPC Server and implements the rpcserver ConnManager
// interface.
type ConnManager struct {
	Server *Node
}

// Ensure rpcConnManager implements the RPCServerConnManager interface.
var _ ServerConnManager = &ConnManager{}

// Connect adds the provided address as a new outbound peer.
//
// The permanent flag indicates whether or not to make the peer persistent and reconnect if the connection is lost.
//
// Attempting to connect to an already existing peer will return an error.
//
// This function is safe for concurrent access and is part of the RPCServerConnManager interface implementation.
func (cm *ConnManager) Connect(addr string, permanent bool) (e error) {
	replyChan := make(chan error)
	cm.Server.Query <- ConnectNodeMsg{
		Addr:      addr,
		Permanent: permanent,
		Reply:     replyChan,
	}
	return <-replyChan
}

// RemoveByID removes the peer associated with the provided id from the list of persistent peers.
//
// Attempting to remove an id that does not exist will return an error.
//
// This function is safe for concurrent access and is part of the RPCServerConnManager interface implementation.
func (cm *ConnManager) RemoveByID(id int32) (e error) {
	replyChan := make(chan error)
	cm.Server.Query <- RemoveNodeMsg{
		Cmp:   func(sp *NodePeer) bool { return sp.ID() == id },
		Reply: replyChan,
	}
	return <-replyChan
}

// RemoveByAddr removes the peer associated with the provided address from the list of persistent peers.
//
// Attempting to remove an address that does not exist will return an error.
//
// This function is safe for concurrent access and is part of the RPCServerConnManager interface implementation.
func (cm *ConnManager) RemoveByAddr(addr string) (e error) {
	replyChan := make(chan error)
	cm.Server.Query <- RemoveNodeMsg{
		Cmp:   func(sp *NodePeer) bool { return sp.Addr() == addr },
		Reply: replyChan,
	}
	return <-replyChan
}

// DisconnectByID disconnects the peer associated with the provided id.
//
// This applies to both inbound and outbound peers.
//
// Attempting to remove an id that does not exist will return an error.
//
// This function is safe for concurrent access and is part of the RPCServerConnManager interface implementation.
func (cm *ConnManager) DisconnectByID(id int32) (e error) {
	replyChan := make(chan error)
	cm.Server.Query <- DisconnectNodeMsg{
		Cmp:   func(sp *NodePeer) bool { return sp.ID() == id },
		Reply: replyChan,
	}
	return <-replyChan
}

// DisconnectByAddr disconnects the peer associated with the provided address. This applies to both inbound and outbound
// peers.
//
// Attempting to remove an address that does not exist will return an error.
//
// This function is safe for concurrent access and is part of the RPCServerConnManager interface implementation.
func (cm *ConnManager) DisconnectByAddr(addr string) (e error) {
	replyChan := make(chan error)
	cm.Server.Query <- DisconnectNodeMsg{
		Cmp:   func(sp *NodePeer) bool { return sp.Addr() == addr },
		Reply: replyChan,
	}
	return <-replyChan
}

// ConnectedCount returns the number of currently connected peers.
//
// This function is safe for concurrent access and is part of the RPCServerConnManager interface implementation.
func (cm *ConnManager) ConnectedCount() int32 {
	return cm.Server.ConnectedCount()
}

// NetTotals returns the sum of all bytes received and sent across the network for all peers.
//
// This function is safe for concurrent access and is part of the RPCServerConnManager interface implementation.
func (cm *ConnManager) NetTotals() (uint64, uint64) {
	return cm.Server.NetTotals()
}

// ConnectedPeers returns an array consisting of all connected peers.
//
// This function is safe for concurrent access and is part of the RPCServerConnManager interface implementation.
func (cm *ConnManager) ConnectedPeers() []ServerPeer {
	replyChan := make(chan []*NodePeer)
	cm.Server.Query <- GetPeersMsg{Reply: replyChan}
	serverPeers := <-replyChan
	// Convert to RPC Server peers.
	peers := make([]ServerPeer, 0, len(serverPeers))
	for _, sp := range serverPeers {
		peers = append(peers, (*Peer)(sp))
	}
	return peers
}

// PersistentPeers returns an array consisting of all the added persistent peers.
//
// This function is safe for concurrent access and is part of the RPCServerConnManager interface implementation.
func (cm *ConnManager) PersistentPeers() []ServerPeer {
	replyChan := make(chan []*NodePeer)
	cm.Server.Query <- GetAddedNodesMsg{Reply: replyChan}
	serverPeers := <-replyChan
	// Convert to generic peers.
	peers := make([]ServerPeer, 0, len(serverPeers))
	for _, sp := range serverPeers {
		peers = append(peers, (*Peer)(sp))
	}
	return peers
}

// BroadcastMessage sends the provided message to all currently connected peers.
//
// This function is safe for concurrent access and is part of the RPCServerConnManager interface implementation.
func (cm *ConnManager) BroadcastMessage(msg wire.Message) {
	cm.Server.BroadcastMessage(msg)
}

// AddRebroadcastInventory adds the provided inventory to the list of inventories to be rebroadcast at random intervals
// until they show up in a block.
//
// This function is safe for concurrent access and is part of the RPCServerConnManager interface implementation.
func (cm *ConnManager) AddRebroadcastInventory(iv *wire.InvVect,
	data interface{},
) {
	cm.Server.AddRebroadcastInventory(iv, data)
}

// RelayTransactions generates and relays inventory vectors for all of the passed transactions to all connected peers.
func (cm *ConnManager) RelayTransactions(txns []*mempool.TxDesc) {
	cm.Server.RelayTransactions(txns)
}

// SyncManager provides a block manager for use with the RPC Server and implements the RPCServerSyncManager interface.
type SyncManager struct {
	Server  *Node
	SyncMgr *netsync.SyncManager
}

// Ensure rpcSyncMgr implements the RPCServerSyncManager interface.
var _ ServerSyncManager = (*SyncManager)(nil)

// IsCurrent returns whether or not the sync manager believes the chain is current as compared to the rest of the
// network.
//
// This function is safe for concurrent access and is part of the RPCServerSyncManager interface implementation.
func (b *SyncManager) IsCurrent() bool {
	return b.SyncMgr.IsCurrent()
}

// SubmitBlock submits the provided block to the network after processing it locally.
//
// This function is safe for concurrent access and is part of the RPCServerSyncManager interface implementation.
func (b *SyncManager) SubmitBlock(block *block.Block,
	flags blockchain.BehaviorFlags,
) (bool, error) {
	return b.SyncMgr.ProcessBlock(block, flags)
}

// Pause pauses the sync manager until the returned channel is closed.
//
// This function is safe for concurrent access and is part of the RPCServerSyncManager interface implementation.
func (b *SyncManager) Pause() chan<- struct{} {
	return b.SyncMgr.Pause()
}

// SyncPeerID returns the peer that is currently the peer being used to sync from.
//
// This function is safe for concurrent access and is part of the RPCServerSyncManager interface implementation.
func (b *SyncManager) SyncPeerID() int32 {
	return b.SyncMgr.SyncPeerID()
}

// LocateHeaders returns the hashes of the blocks after the first known block in
// the provided locators until the provided
// stop hash or the current tip is reached, up to a max of wire.MaxBlockHeadersPerMsg hashes.
//
// This function is safe for concurrent access and is part of the RPCServerSyncManager interface implementation.
func (b *SyncManager) LocateHeaders(locators []*chainhash.Hash,
	hashStop *chainhash.Hash,
) []wire.BlockHeader {
	return b.Server.Chain.LocateHeaders(locators, hashStop)
}
