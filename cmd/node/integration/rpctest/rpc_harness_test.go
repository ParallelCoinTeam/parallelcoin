package rpctest

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/p9c/matrjoska/pkg/amt"
	"github.com/p9c/matrjoska/pkg/btcaddr"

	"github.com/p9c/qu"

	"github.com/p9c/matrjoska/pkg/chaincfg"
	"github.com/p9c/matrjoska/pkg/chainhash"
	"github.com/p9c/matrjoska/pkg/txscript"
	"github.com/p9c/matrjoska/pkg/util"
	"github.com/p9c/matrjoska/pkg/wire"
)

func testSendOutputs(r *Harness, t *testing.T) {
	genSpend := func(amt amt.Amount) *chainhash.Hash {
		// Grab a fresh address from the wallet.
		addr, e := r.NewAddress()
		if e != nil {
			t.Fatalf("unable to get new address: %v", e)
		}
		// Next, send amt DUO to this address,
		// spending from one of our mature coinbase outputs.
		addrScript, e := txscript.PayToAddrScript(addr)
		if e != nil {
			t.Fatalf("unable to generate pkscript to addr: %v", e)
		}
		output := wire.NewTxOut(int64(amt), addrScript)
		txid, e := r.SendOutputs([]*wire.TxOut{output}, 10)
		if e != nil {
			t.Fatalf("coinbase spend failed: %v", e)
		}
		return txid
	}
	assertTxMined := func(txid *chainhash.Hash, blockHash *chainhash.Hash) {
		block, e := r.Node.GetBlock(blockHash)
		if e != nil {
			t.Fatalf("unable to get block: %v", e)
		}
		numBlockTxns := len(block.Transactions)
		if numBlockTxns < 2 {
			t.Fatalf(
				"crafted transaction wasn't mined, block should have "+
					"at least %v transactions instead has %v", 2, numBlockTxns,
			)
		}
		minedTx := block.Transactions[1]
		txHash := minedTx.TxHash()
		if txHash != *txid {
			t.Fatalf("txid's don't match, %v vs %v", txHash, txid)
		}
	}
	// First, generate a small spend which will require only a single input.
	txid := genSpend(5 * amt.SatoshiPerBitcoin)
	// Generate a single block, the transaction the wallet created should be found in this block.
	blockHashes, e := r.Node.Generate(1)
	if e != nil {
		t.Fatalf("unable to generate single block: %v", e)
	}
	assertTxMined(txid, blockHashes[0])
	// Next, generate a spend much greater than the block reward. This transaction should also have been mined properly.
	txid = genSpend(500 * amt.SatoshiPerBitcoin)
	blockHashes, e = r.Node.Generate(1)
	if e != nil {
		t.Fatalf("unable to generate single block: %v", e)
	}
	assertTxMined(txid, blockHashes[0])
}

func assertConnectedTo(t *testing.T, nodeA *Harness, nodeB *Harness) {
	nodeAPeers, e := nodeA.Node.GetPeerInfo()
	if e != nil {
		t.Fatalf("unable to get nodeA's peer info")
	}
	nodeAddr := nodeB.node.config.listen
	addrFound := false
	for _, peerInfo := range nodeAPeers {
		if peerInfo.Addr == nodeAddr {
			addrFound = true
			break
		}
	}
	if !addrFound {
		t.Fatal("nodeA not connected to nodeB")
	}
}

func testConnectNode(r *Harness, t *testing.T) {
	// Create a fresh test harness.
	harness, e := New(&chaincfg.SimNetParams, nil, nil)
	if e != nil {
		t.Fatal(e)
	}
	if e := harness.SetUp(false, 0); E.Chk(e) {
		t.Fatalf("unable to complete rpctest setup: %v", e)
	}
	defer func() {
		if e := harness.TearDown(); E.Chk(e) {
		}
	}()
	// Establish a p2p connection from our new local harness to the main harness.
	if e := ConnectNode(harness, r); E.Chk(e) {
		t.Fatalf("unable to connect local to main harness: %v", e)
	}
	// The main harness should show up in our local harness' peer's list, and vice verse.
	assertConnectedTo(t, harness, r)
}

func testTearDownAll(t *testing.T) {
	// Grab a local copy of the currently active harnesses before attempting to tear them all down.
	initialActiveHarnesses := ActiveHarnesses()
	// Tear down all currently active harnesses.
	if e := TearDownAll(); E.Chk(e) {
		t.Fatalf("unable to teardown all harnesses: %v", e)
	}
	// The global testInstances map should now be fully purged with no active test harnesses remaining.
	if len(ActiveHarnesses()) != 0 {
		t.Fatalf("test harnesses still active after TearDownAll")
	}
	for _, harness := range initialActiveHarnesses {
		// Ensure all test directories have been deleted.
		var e error
		if _, e = os.Stat(harness.testNodeDir); e == nil {
			t.Errorf("created test datadir was not deleted.")
		}
	}
}
func testActiveHarnesses(r *Harness, t *testing.T) {
	numInitialHarnesses := len(ActiveHarnesses())
	// Create a single test harness.
	harness1, e := New(&chaincfg.SimNetParams, nil, nil)
	if e != nil {
		t.Fatal(e)
	}
	defer func() {
		if e := harness1.TearDown(); E.Chk(e) {
		}
	}()
	// With the harness created above, a single harness should be detected as active.
	numActiveHarnesses := len(ActiveHarnesses())
	if !(numActiveHarnesses > numInitialHarnesses) {
		t.Fatalf(
			"ActiveHarnesses not updated, should have an " +
				"additional test harness listed.",
		)
	}
}
func testJoinMempools(r *Harness, t *testing.T) {
	// Assert main test harness has no transactions in its mempool.
	pooledHashes, e := r.Node.GetRawMempool()
	if e != nil {
		t.Fatalf("unable to get mempool for main test harness: %v", e)
	}
	if len(pooledHashes) != 0 {
		t.Fatal("main test harness mempool not empty")
	}
	// Create a local test harness with only the genesis block. The nodes will be synced below so the same transaction
	// can be sent to both nodes without it being an orphan.
	var harness *Harness
	harness, e = New(&chaincfg.SimNetParams, nil, nil)
	if e != nil {
		t.Fatal(e)
	}
	if e = harness.SetUp(false, 0); E.Chk(e) {
		t.Fatalf("unable to complete rpctest setup: %v", e)
	}
	defer func() {
		if e = harness.TearDown(); E.Chk(e) {
		}
	}()
	nodeSlice := []*Harness{r, harness}
	// Both mempools should be considered synced as they are empty. Therefore, this should return instantly.
	if e = JoinNodes(nodeSlice, Mempools); E.Chk(e) {
		t.Fatalf("unable to join node on mempools: %v", e)
	}
	// Generate a coinbase spend to a new address within the main harness' mempool.
	var addr btcaddr.Address
	if addr, e = r.NewAddress(); E.Chk(e) {
	}
	var addrScript []byte
	addrScript, e = txscript.PayToAddrScript(addr)
	if e != nil {
		t.Fatalf("unable to generate pkscript to addr: %v", e)
	}
	output := wire.NewTxOut(5e8, addrScript)
	var testTx *wire.MsgTx
	testTx, e = r.CreateTransaction([]*wire.TxOut{output}, 10, true)
	if e != nil {
		t.Fatalf("coinbase spend failed: %v", e)
	}
	if _, e = r.Node.SendRawTransaction(testTx, true); E.Chk(e) {
		t.Fatalf("send transaction failed: %v", e)
	}
	// Wait until the transaction shows up to ensure the two mempools are not the same.
	harnessSynced := qu.T()
	go func() {
		for {
			var poolHashes []*chainhash.Hash
			poolHashes, e = r.Node.GetRawMempool()
			if e != nil {
				t.Fatalf("failed to retrieve harness mempool: %v", e)
			}
			if len(poolHashes) > 0 {
				break
			}
			time.Sleep(time.Millisecond * 100)
		}
		harnessSynced <- struct{}{}
	}()
	select {
	case <-harnessSynced.Wait():
	case <-time.After(time.Minute):
		t.Fatalf("harness node never received transaction")
	}
	// This select case should fall through to the default as the goroutine should be blocked on the JoinNodes call.
	poolsSynced := qu.T()
	go func() {
		if e = JoinNodes(nodeSlice, Mempools); E.Chk(e) {
			t.Fatalf("unable to join node on mempools: %v", e)
		}
		poolsSynced <- struct{}{}
	}()
	select {
	case <-poolsSynced.Wait():
		t.Fatalf("mempools detected as synced yet harness has a new tx")
	default:
	}
	// Establish an outbound connection from the local harness to the main harness and wait for the chains to be synced.
	if e = ConnectNode(harness, r); E.Chk(e) {
		t.Fatalf("unable to connect harnesses: %v", e)
	}
	if e = JoinNodes(nodeSlice, Blocks); E.Chk(e) {
		t.Fatalf("unable to join node on blocks: %v", e)
	}
	// Send the transaction to the local harness which will result in synced mempools.
	if _, e = harness.Node.SendRawTransaction(testTx, true); E.Chk(e) {
		t.Fatalf("send transaction failed: %v", e)
	}
	// Select once again with a special timeout case after 1 minute. The goroutine above should now be blocked on
	// sending into the unbuffered channel. The send should immediately succeed. In order to avoid the test hanging
	// indefinitely, a 1 minute timeout is in place.
	select {
	case <-poolsSynced.Wait():
	case <-time.After(time.Minute):
		t.Fatalf("mempools never detected as synced")
	}
}
func testJoinBlocks(r *Harness, t *testing.T) {
	// Create a second harness with only the genesis block so it is behind the main harness.
	harness, e := New(&chaincfg.SimNetParams, nil, nil)
	if e != nil {
		t.Fatal(e)
	}
	if e := harness.SetUp(false, 0); E.Chk(e) {
		t.Fatalf("unable to complete rpctest setup: %v", e)
	}
	defer func() {
		if e := harness.TearDown(); E.Chk(e) {
		}
	}()
	nodeSlice := []*Harness{r, harness}
	blocksSynced := qu.T()
	go func() {
		if e := JoinNodes(nodeSlice, Blocks); E.Chk(e) {
			t.Fatalf("unable to join node on blocks: %v", e)
		}
		blocksSynced <- struct{}{}
	}()
	// This select case should fall through to the default as the goroutine should be blocked on the JoinNodes calls.
	select {
	case <-blocksSynced.Wait():
		t.Fatalf("blocks detected as synced yet local harness is behind")
	default:
	}
	// Connect the local harness to the main harness which will sync the chains.
	if e := ConnectNode(harness, r); E.Chk(e) {
		t.Fatalf("unable to connect harnesses: %v", e)
	}
	// Select once again with a special timeout case after 1 minute. The goroutine above should now be blocked on
	// sending into the unbuffered channel. The send should immediately succeed. In order to avoid the test hanging
	// indefinitely, a 1 minute timeout is in place.
	select {
	case <-blocksSynced.Wait():
	case <-time.After(time.Minute):
		t.Fatalf("blocks never detected as synced")
	}
}
func testGenerateAndSubmitBlock(r *Harness, t *testing.T) {
	// Generate a few test spend transactions.
	addr, e := r.NewAddress()
	if e != nil {
		t.Fatalf("unable to generate new address: %v", e)
	}
	pkScript, e := txscript.PayToAddrScript(addr)
	if e != nil {
		t.Fatalf("unable to create script: %v", e)
	}
	output := wire.NewTxOut(amt.SatoshiPerBitcoin.Int64(), pkScript)
	const numTxns = 5
	txns := make([]*util.Tx, 0, numTxns)
	var tx *wire.MsgTx
	for i := 0; i < numTxns; i++ {
		tx, e = r.CreateTransaction([]*wire.TxOut{output}, 10, true)
		if e != nil {
			t.Fatalf("unable to create tx: %v", e)
		}
		txns = append(txns, util.NewTx(tx))
	}
	// Now generate a block with the default block version, and a zeroed out time.
	block, e := r.GenerateAndSubmitBlock(txns, ^uint32(0), time.Time{})
	if e != nil {
		t.Fatalf("unable to generate block: %v", e)
	}
	// Ensure that all created transactions were included, and that the block version was properly set to the default.
	numBlocksTxns := len(block.Transactions())
	if numBlocksTxns != numTxns+1 {
		t.Fatalf(
			"block did not include all transactions: "+
				"expected %v, got %v", numTxns+1, numBlocksTxns,
		)
	}
	blockVersion := block.WireBlock().Header.Version
	if blockVersion != BlockVersion {
		t.Fatalf(
			"block version is not default: expected %v, got %v",
			BlockVersion, blockVersion,
		)
	}
	// Next generate a block with a "non-standard" block version along with time stamp a minute after the previous
	// block's timestamp.
	timestamp := block.WireBlock().Header.Timestamp.Add(time.Minute)
	targetBlockVersion := uint32(1337)
	block, e = r.GenerateAndSubmitBlock(nil, targetBlockVersion, timestamp)
	if e != nil {
		t.Fatalf("unable to generate block: %v", e)
	}
	// Finally ensure that the desired block version and timestamp were set properly.
	header := block.WireBlock().Header
	blockVersion = header.Version
	if blockVersion != int32(targetBlockVersion) {
		t.Fatalf(
			"block version mismatch: expected %v, got %v",
			targetBlockVersion, blockVersion,
		)
	}
	if !timestamp.Equal(header.Timestamp) {
		t.Fatalf(
			"header time stamp mismatch: expected %v, got %v",
			timestamp, header.Timestamp,
		)
	}
}
func testGenerateAndSubmitBlockWithCustomCoinbaseOutputs(
	r *Harness,
	t *testing.T,
) {
	// Generate a few test spend transactions.
	addr, e := r.NewAddress()
	if e != nil {
		t.Fatalf("unable to generate new address: %v", e)
	}
	pkScript, e := txscript.PayToAddrScript(addr)
	if e != nil {
		t.Fatalf("unable to create script: %v", e)
	}
	output := wire.NewTxOut(amt.SatoshiPerBitcoin.Int64(), pkScript)
	const numTxns = 5
	txns := make([]*util.Tx, 0, numTxns)
	for i := 0; i < numTxns; i++ {
		var tx *wire.MsgTx
		tx, e = r.CreateTransaction([]*wire.TxOut{output}, 10, true)
		if e != nil {
			t.Fatalf("unable to create tx: %v", e)
		}
		txns = append(txns, util.NewTx(tx))
	}
	// Now generate a block with the default block version, a zero'd out time, and a burn output.
	block, e := r.GenerateAndSubmitBlockWithCustomCoinbaseOutputs(
		txns,
		^uint32(0), time.Time{}, []wire.TxOut{
			{
				Value:    0,
				PkScript: []byte{},
			},
		},
	)
	if e != nil {
		t.Fatalf("unable to generate block: %v", e)
	}
	// Ensure that all created transactions were included, and that the block version was properly set to the default.
	numBlocksTxns := len(block.Transactions())
	if numBlocksTxns != numTxns+1 {
		t.Fatalf(
			"block did not include all transactions: "+
				"expected %v, got %v", numTxns+1, numBlocksTxns,
		)
	}
	blockVersion := block.WireBlock().Header.Version
	if blockVersion != BlockVersion {
		t.Fatalf(
			"block version is not default: expected %v, got %v",
			BlockVersion, blockVersion,
		)
	}
	// Next generate a block with a "non-standard" block version along with time stamp a minute after the previous
	// block's timestamp.
	timestamp := block.WireBlock().Header.Timestamp.Add(time.Minute)
	targetBlockVersion := uint32(1337)
	block, e = r.GenerateAndSubmitBlockWithCustomCoinbaseOutputs(
		nil,
		targetBlockVersion, timestamp, []wire.TxOut{
			{
				Value:    0,
				PkScript: []byte{},
			},
		},
	)
	if e != nil {
		t.Fatalf("unable to generate block: %v", e)
	}
	// Finally ensure that the desired block version and timestamp were set properly.
	header := block.WireBlock().Header
	blockVersion = header.Version
	if blockVersion != int32(targetBlockVersion) {
		t.Fatalf(
			"block version mismatch: expected %v, got %v",
			targetBlockVersion, blockVersion,
		)
	}
	if !timestamp.Equal(header.Timestamp) {
		t.Fatalf(
			"header time stamp mismatch: expected %v, got %v",
			timestamp, header.Timestamp,
		)
	}
}
func testMemWalletReorg(r *Harness, t *testing.T) {
	// Create a fresh harness, we'll be using the main harness to force a re-org on this local harness.
	harness, e := New(&chaincfg.SimNetParams, nil, nil)
	if e != nil {
		t.Fatal(e)
	}
	if e := harness.SetUp(true, 5); E.Chk(e) {
		t.Fatalf("unable to complete rpctest setup: %v", e)
	}
	defer func() {
		if e := harness.TearDown(); E.Chk(e) {
		}
	}()
	// The internal wallet of this harness should now have 250 DUO.
	expectedBalance := 250 * amt.SatoshiPerBitcoin
	walletBalance := harness.ConfirmedBalance()
	if expectedBalance != walletBalance {
		t.Fatalf(
			"wallet balance incorrect: expected %v, got %v",
			expectedBalance, walletBalance,
		)
	}
	// Now connect this local harness to the main harness then wait for their chains to synchronize.
	if e := ConnectNode(harness, r); E.Chk(e) {
		t.Fatalf("unable to connect harnesses: %v", e)
	}
	nodeSlice := []*Harness{r, harness}
	if e := JoinNodes(nodeSlice, Blocks); E.Chk(e) {
		t.Fatalf("unable to join node on blocks: %v", e)
	}
	// The original wallet should now have a balance of 0 DUO as its entire chain should have been decimated in favor of
	// the main harness' chain.
	expectedBalance = amt.Amount(0)
	walletBalance = harness.ConfirmedBalance()
	if expectedBalance != walletBalance {
		t.Fatalf(
			"wallet balance incorrect: expected %v, got %v",
			expectedBalance, walletBalance,
		)
	}
}
func testMemWalletLockedOutputs(r *Harness, t *testing.T) {
	// Obtain the initial balance of the wallet at this point.
	startingBalance := r.ConfirmedBalance()
	// First, create a signed transaction spending some outputs.
	addr, e := r.NewAddress()
	if e != nil {
		t.Fatalf("unable to generate new address: %v", e)
	}
	pkScript, e := txscript.PayToAddrScript(addr)
	if e != nil {
		t.Fatalf("unable to create script: %v", e)
	}
	outputAmt := 50 * amt.SatoshiPerBitcoin
	output := wire.NewTxOut(int64(outputAmt), pkScript)
	tx, e := r.CreateTransaction([]*wire.TxOut{output}, 10, true)
	if e != nil {
		t.Fatalf("unable to create transaction: %v", e)
	}
	// The current wallet balance should now be at least 50 DUO less (accounting for fees) than the period balance
	currentBalance := r.ConfirmedBalance()
	if !(currentBalance <= startingBalance-outputAmt) {
		t.Fatalf(
			"spent outputs not locked: previous balance %v, "+
				"current balance %v", startingBalance, currentBalance,
		)
	}
	// Now unlocked all the spent inputs within the unbroadcast signed transaction. The current balance should now be
	// exactly that of the starting balance.
	r.UnlockOutputs(tx.TxIn)
	currentBalance = r.ConfirmedBalance()
	if currentBalance != startingBalance {
		t.Fatalf(
			"current and starting balance should now match: "+
				"expected %v, got %v", startingBalance, currentBalance,
		)
	}
}

var harnessTestCases = []HarnessTestCase{
	testSendOutputs,
	testConnectNode,
	testActiveHarnesses,
	testJoinBlocks,
	testJoinMempools, // Depends on results of testJoinBlocks
	testGenerateAndSubmitBlock,
	testGenerateAndSubmitBlockWithCustomCoinbaseOutputs,
	testMemWalletReorg,
	testMemWalletLockedOutputs,
}

var mainHarness *Harness

const (
	numMatureOutputs = 25
)

func TestMain(m *testing.M) {
	var e error
	mainHarness, e = New(&chaincfg.SimNetParams, nil, nil)
	if e != nil {
		fmt.Println("unable to create main harness: ", e)
		os.Exit(1)
	}
	// Initialize the main mining node with a chain of length 125, providing 25 mature coinbases to allow spending from
	// for testing purposes.
	if e = mainHarness.SetUp(true, numMatureOutputs); E.Chk(e) {
		fmt.Println("unable to setup test chain: ", e)
		// Even though the harness was not fully setup, it still needs to be torn down to ensure all resources such as
		// temp directories are cleaned up. The error is intentionally ignored since this is already an error path and
		// nothing else could be done about it anyways.
		_ = mainHarness.TearDown()
		os.Exit(1)
	}
	exitCode := m.Run()
	// Clean up any active harnesses that are still currently running.
	if len(ActiveHarnesses()) > 0 {
		if e := TearDownAll(); E.Chk(e) {
			fmt.Println("unable to tear down chain: ", e)
			os.Exit(1)
		}
	}
	os.Exit(exitCode)
}
func TestHarness(t *testing.T) {
	// We should have (numMatureOutputs * 50 DUO) of mature unspendable
	// outputs.
	expectedBalance := numMatureOutputs * 50 * amt.SatoshiPerBitcoin
	harnessBalance := mainHarness.ConfirmedBalance()
	if harnessBalance != expectedBalance {
		t.Fatalf(
			"expected wallet balance of %v instead have %v",
			expectedBalance, harnessBalance,
		)
	}
	// Current tip should be at a height of numMatureOutputs plus the required number of blocks for coinbase maturity.
	nodeInfo, e := mainHarness.Node.GetInfo()
	if e != nil {
		t.Fatalf("unable to execute getinfo on node: %v", e)
	}
	expectedChainHeight := numMatureOutputs + uint32(mainHarness.ActiveNet.CoinbaseMaturity)
	if uint32(nodeInfo.Blocks) != expectedChainHeight {
		t.Errorf(
			"Chain height is %v, should be %v",
			nodeInfo.Blocks, expectedChainHeight,
		)
	}
	for _, testCase := range harnessTestCases {
		testCase(mainHarness, t)
	}
	testTearDownAll(t)
}
