package rpctest

import (
	"reflect"
	"time"
	
	"github.com/p9c/matrjoska/pkg/chainhash"
	"github.com/p9c/matrjoska/pkg/rpcclient"
)

// JoinType is an enum representing a particular type of "node join". A node
// join is a synchronization tool used to wait until a subset of nodes have a
// consistent state with respect to an attribute.
type JoinType uint8

const (
	// Blocks is a JoinType which waits until all nodes share the same block
	// height.
	Blocks JoinType = iota
	// Mempools is a JoinType which blocks until all nodes have identical mempool.
	Mempools
)

// JoinNodes is a synchronization tool used to block until all passed nodes are
// fully synced with respect to an attribute. This function will block for a
// period of time, finally returning once all nodes are synced according to the
// passed JoinType. This function be used to to ensure all active test harnesses
// are at a consistent state before proceeding to an assertion or check within
// rpc tests.
func JoinNodes(nodes []*Harness, joinType JoinType) (e error) {
	switch joinType {
	case Blocks:
		return syncBlocks(nodes)
	case Mempools:
		return syncMempools(nodes)
	}
	return nil
}

// syncMempools blocks until all nodes have identical mempools.
func syncMempools(nodes []*Harness) (e error) {
	poolsMatch := false
retry:
	for !poolsMatch {
		firstPool, e := nodes[0].Node.GetRawMempool()
		if e != nil {
			return e
		}
		// If all nodes have an identical mempool with respect to the first node,
		// then we're done. Otherwise drop back to the top of the loop and retry
		// after a short wait period.
		for _, node := range nodes[1:] {
			nodePool, e := node.Node.GetRawMempool()
			if e != nil {
				return e
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
func syncBlocks(nodes []*Harness) (e error) {
	blocksMatch := false
retry:
	for !blocksMatch {
		var prevHash *chainhash.Hash
		var prevHeight int32
		
		for _, node := range nodes {
			blockHash, blockHeight, e := node.Node.GetBestBlock()
			if e != nil {
				return e
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
	return nil
}

// ConnectNode establishes a new peer-to-peer connection between the "from" harness and the "to" harness. The connection
// made is flagged as persistent therefore in the case of disconnects, "from" will attempt to reestablish a connection
// to the "to" harness.
func ConnectNode(from *Harness, to *Harness) (e error) {
	peerInfo, e := from.Node.GetPeerInfo()
	if e != nil {
		return e
	}
	numPeers := len(peerInfo)
	targetAddr := to.node.config.listen
	if e = from.Node.AddNode(targetAddr, rpcclient.ANAdd); E.Chk(e) {
		return e
	}
	// Block until a new connection has been established.
	peerInfo, e = from.Node.GetPeerInfo()
	if e != nil {
		return e
	}
	for len(peerInfo) <= numPeers {
		peerInfo, e = from.Node.GetPeerInfo()
		if e != nil {
			return e
		}
	}
	return nil
}

// TearDownAll tears down all active test harnesses.
func TearDownAll() (e error) {
	harnessStateMtx.Lock()
	defer harnessStateMtx.Unlock()
	for _, harness := range testInstances {
		if e := harness.tearDown(); E.Chk(e) {
			return e
		}
	}
	return nil
}

// ActiveHarnesses returns a slice of all currently active test harnesses. A test harness if considered "active" if it
// has been created, but not yet torn down.
func ActiveHarnesses() []*Harness {
	harnessStateMtx.RLock()
	defer harnessStateMtx.RUnlock()
	activeNodes := make([]*Harness, 0, len(testInstances))
	for _, harness := range testInstances {
		activeNodes = append(activeNodes, harness)
	}
	return activeNodes
}
