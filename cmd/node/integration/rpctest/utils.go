package rpctest

import (
	"github.com/p9c/pod/pkg/rpc/btcjson"
	"github.com/p9c/pkg/app/slog"
	"reflect"
	"time"

	chainhash "github.com/p9c/pod/pkg/chain/hash"
	client "github.com/p9c/pod/pkg/rpc/client"
)

// JoinType is an enum representing a particular type of "node join". A node
// join is a synchronization tool used to wait until a subset of nodes have a
// consistent state with respect to an attribute.
type JoinType uint8

const (
	// BlockC is a JoinType which waits until all nodes share the same block height.
	Blocks JoinType = iota
	// Mempools is a JoinType which blocks until all nodes have identical mempool.
	Mempools
)

// JoinNodes is a synchronization tool used to block until all passed nodes
// are fully synced with respect to an attribute.
// This function will block for a period of time,
// finally returning once all nodes are synced according to the passed
// JoinType. This function be used to to ensure all active test harnesses are
// at a consistent state before proceeding to an assertion or check within
// rpc tests.
func JoinNodes(nodes []*Harness, joinType JoinType) (err error) {
	switch joinType {
	case Blocks:
		return syncBlocks(nodes)
	case Mempools:
		return syncMempools(nodes)
	}
	return
}

// syncMempools blocks until all nodes have identical mempools.
func syncMempools(nodes []*Harness) (err error) {
	poolsMatch := false
retry:
	for !poolsMatch {
		var firstPool []*chainhash.Hash
		if firstPool, err = nodes[0].Node.GetRawMempool(); slog.Check(err) {
			return
		}
		// If all nodes have an identical mempool with respect to the first
		// node, then we're done. Otherwise drop back to the top of the loop
		// and retry after a short wait period.
		for _, node := range nodes[1:] {
			var nodePool []*chainhash.Hash
			if nodePool, err = node.Node.GetRawMempool(); slog.Check(err) {
				return
			}
			if !reflect.DeepEqual(firstPool, nodePool) {
				time.Sleep(time.Millisecond * 100)
				continue retry
			}
		}
		poolsMatch = true
	}
	return nil
}

// syncBlocks blocks until all nodes report the same best chain.
func syncBlocks(nodes []*Harness) (err error) {
	blocksMatch := false
retry:
	for !blocksMatch {
		var prevHash, blockHash *chainhash.Hash
		var prevHeight, blockHeight int32
		for _, node := range nodes {
			if blockHash, blockHeight, err = node.Node.GetBestBlock(); slog.Check(err) {
				return
			}
			if prevHash != nil && (*blockHash != *prevHash ||
				blockHeight != prevHeight) {
				time.Sleep(time.Millisecond * 100)
				continue retry
			}
			prevHash, prevHeight = blockHash, blockHeight
		}
		blocksMatch = true
	}
	return
}

// ConnectNode establishes a new peer-to-peer connection between the "from"
// harness and the "to" harness.
// The connection made is flagged as persistent therefore in the case of
// disconnects, "from" will attempt to reestablish a connection to the "to"
// harness.
func ConnectNode(from *Harness, to *Harness) (err error) {
	var peerInfo []btcjson.GetPeerInfoResult
	if peerInfo, err = from.Node.GetPeerInfo(); slog.Check(err) {
		return
	}
	numPeers := len(peerInfo)
	targetAddr := to.node.config.listen
	if err = from.Node.AddNode(targetAddr, client.ANAdd); slog.Check(err) {
		return
	}
	// Block until a new connection has been established.
	if peerInfo, err = from.Node.GetPeerInfo(); slog.Check(err) {
		return err
	}
	for len(peerInfo) <= numPeers {
		if peerInfo, err = from.Node.GetPeerInfo(); slog.Check(err) {
			return
		}
	}
	return
}

// TearDownAll tears down all active test harnesses.
func TearDownAll() (err error) {
	harnessStateMtx.Lock()
	defer harnessStateMtx.Unlock()
	for _, harness := range testInstances {
		if err = harness.tearDown(); slog.Check(err) {
			return
		}
	}
	return
}

// ActiveHarnesses returns a slice of all currently active test harnesses.
// A test harness if considered "active" if it has been created,
// but not yet torn down.
func ActiveHarnesses() []*Harness {
	harnessStateMtx.RLock()
	defer harnessStateMtx.RUnlock()
	activeNodes := make([]*Harness, 0, len(testInstances))
	for _, harness := range testInstances {
		activeNodes = append(activeNodes, harness)
	}
	return activeNodes
}
