package integration

import (
	"bytes"
	"fmt"
	"os"
	"runtime/debug"
	"testing"

	"github.com/p9c/matrjoska/cmd/node/integration/rpctest"
	"github.com/p9c/matrjoska/pkg/chaincfg"
)

func testGetBestBlock(r *rpctest.Harness, t *testing.T) {
	_, prevbestHeight, e := r.Node.GetBestBlock()
	if e != nil {
		t.Fatalf("Call to `getbestblock` failed: %v", e)
	}
	// Create a new block connecting to the current tip.
	generatedBlockHashes, e := r.Node.Generate(1)
	if e != nil {
		t.Fatalf("Unable to generate block: %v", e)
	}
	bestHash, bestHeight, e := r.Node.GetBestBlock()
	if e != nil {
		t.Fatalf("Call to `getbestblock` failed: %v", e)
	}
	// Hash should be the same as the newly submitted block.
	if !bytes.Equal(bestHash[:], generatedBlockHashes[0][:]) {
		t.Fatalf(
			"Block hashes do not match. Returned hash %v, wanted "+
				"hash %v", bestHash, generatedBlockHashes[0][:],
		)
	}
	// Block height should now reflect newest height.
	if bestHeight != prevbestHeight+1 {
		t.Fatalf(
			"Block heights do not match. Got %v, wanted %v",
			bestHeight, prevbestHeight+1,
		)
	}
}

func testGetBlockCount(r *rpctest.Harness, t *testing.T) {
	// Save the current count.
	currentCount, e := r.Node.GetBlockCount()
	if e != nil {
		t.Fatalf("Unable to get block count: %v", e)
	}
	if _, e = r.Node.Generate(1); E.Chk(e) {
		t.Fatalf("Unable to generate block: %v", e)
	}
	// Count should have increased by one.
	newCount, e := r.Node.GetBlockCount()
	if e != nil {
		t.Fatalf("Unable to get block count: %v", e)
	}
	if newCount != currentCount+1 {
		t.Fatalf(
			"Block count incorrect. Got %v should be %v",
			newCount, currentCount+1,
		)
	}
}

func testGetBlockHash(r *rpctest.Harness, t *testing.T) {
	// Create a new block connecting to the current tip.
	generatedBlockHashes, e := r.Node.Generate(1)
	if e != nil {
		t.Fatalf("Unable to generate block: %v", e)
	}
	info, e := r.Node.GetInfo()
	if e != nil {
		t.Fatalf("call to getinfo cailed: %v", e)
	}
	blockHash, e := r.Node.GetBlockHash(int64(info.Blocks))
	if e != nil {
		t.Fatalf("Call to `getblockhash` failed: %v", e)
	}
	// Block hashes should match newly created block.
	if !bytes.Equal(generatedBlockHashes[0][:], blockHash[:]) {
		
		t.Fatalf(
			"Block hashes do not match. Returned hash %v, wanted "+
				"hash %v", blockHash, generatedBlockHashes[0][:],
		)
	}
}

var rpcTestCases = []rpctest.HarnessTestCase{
	testGetBestBlock,
	testGetBlockCount,
	testGetBlockHash,
}
var primaryHarness *rpctest.Harness

func TestMain(m *testing.M) {
	var e error
	// In order to properly test scenarios on as if we were on mainnet, ensure that non-standard transactions aren't
	// accepted into the mempool or relayed.
	podCfg := []string{"--rejectnonstd"}
	primaryHarness, e = rpctest.New(&chaincfg.SimNetParams, nil, podCfg)
	if e != nil {
		fmt.Println("unable to create primary harness: ", e)
		os.Exit(1)
	}
	// Initialize the primary mining node with a chain of length 125, providing 25 mature coinbases to allow spending
	// from for testing purposes.
	if e := primaryHarness.SetUp(true, 25); E.Chk(e) {
		fmt.Println("unable to setup test chain: ", e)
		// Even though the harness was not fully setup, it still needs to be torn down to ensure all resources such as
		// temp directories are cleaned up. The error is intentionally ignored since this is already an error path and
		// nothing else could be done about it anyways.
		_ = primaryHarness.TearDown()
		os.Exit(1)
	}
	exitCode := m.Run()
	// Clean up any active harnesses that are still currently running. This includes removing all temporary directories,
	// and shutting down any created processes.
	if e := rpctest.TearDownAll(); E.Chk(e) {
		fmt.Println("unable to tear down all harnesses: ", e)
		os.Exit(1)
	}
	os.Exit(exitCode)
}

func TestRpcServer(t *testing.T) {
	var currentTestNum int
	defer func() {
		// If one of the integration tests caused a panic within the main goroutine, then tear down all the harnesses in
		// order to avoid any leaked pod processes.
		if r := recover(); r != nil {
			fmt.Println("recovering from test panic: ", r)
			if e := rpctest.TearDownAll(); E.Chk(e) {
				fmt.Println("unable to tear down all harnesses: ", e)
			}
			t.Fatalf("test #%v panicked: %s", currentTestNum, debug.Stack())
		}
	}()
	for _, testCase := range rpcTestCases {
		testCase(primaryHarness, t)
		currentTestNum++
	}
}
