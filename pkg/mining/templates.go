package mining

//
// import (
// 	"container/heap"
// 	"fmt"
// 	blockchain "github.com/p9c/matrjoska/pkg/chain"
// 	"github.com/p9c/matrjoska/pkg/chain/fork"
// 	chainhash "github.com/p9c/matrjoska/pkg/chain/hash"
// 	txscript "github.com/p9c/matrjoska/pkg/chain/tx/script"
// 	"github.com/p9c/matrjoska/pkg/chain/wire"
// 	"github.com/p9c/matrjoska/pkg/util"
// 	"math/rand"
// 	"time"
// )
//
// // GenBlockHeader generate a block given a version number to use for mining
// // (nonce is empty, date can be updated, version changes merkle and target bits.
// // All the data required for this is in the exported fields that are serialized
// // for over the wire
// func (m *MsgBlockTemplate) GenBlockHeader(vers int32) *wire.BlockHeader {
// 	return &wire.BlockHeader{
// 		Version:    vers,
// 		PrevBlock:  m.PrevBlock,
// 		MerkleRoot: m.Merkles[vers],
// 		Timestamp:  m.Timestamp,
// 		Bits:       m.Bits[vers],
// 	}
// }
//
// // Reconstruct takes a received block from the wire and reattaches the transactions
// func (m *MsgBlockTemplate) Reconstruct(hdr *wire.BlockHeader) *wire.WireBlock {
// 	msgBlock := &wire.WireBlock{Header: *hdr}
// 	// the coinbase is the last transaction
// 	txs := append(m.txs, m.coinbases[msgBlock.Header.Version])
// 	for _, tx := range txs {
// 		if e := msgBlock.AddTransaction(tx.MsgTx()); E.Chk(e) {
// 			return nil
// 		}
// 	}
// 	return msgBlock
// }
//
// // NewBlockTemplates returns a data structure which has methods to construct
// // block version specific block headers and reconstruct their transactions
// func (g *BlkTmplGenerator) NewBlockTemplates(
// 	workerNumber uint32,
// 	payToAddress util.Address,
// ) (*MsgBlockTemplate, error) {
// 	mbt := &MsgBlockTemplate{Bits: make(map[int32]uint32), Merkles: make(map[int32]chainhash.Hash)}
// 	// Extend the most recently known best block.
// 	best := g.Chain.BestSnapshot()
// 	mbt.PrevBlock = best.Hash
// 	mbt.Timestamp = g.Chain.BestChain.Tip().Header().Timestamp.Add(time.Second)
// 	mbt.Height = best.Height + 1
// 	// Create a standard coinbase transaction paying to the provided address.
// 	//
// 	// NOTE: The coinbase value will be updated to include the fees from the
// 	// selected transactions later after they have actually been selected. It is
// 	// created here to detect any errors early before potentially doing a lot of
// 	// work below. The extra nonce helps ensure the transaction is not a duplicate
// 	// transaction (paying the same value to the same public key address would
// 	// otherwise be an identical transaction for block version 1).
// 	rand.Seed(time.Now().UnixNano())
// 	extraNonce := rand.Uint64()
// 	var e error
// 	// numAlgos := fork.GetNumAlgos(mbt.Height)
// 	// coinbaseScripts := make(map[int32][]byte, numAlgos)
// 	// coinbaseTxs := make(map[int32]*util.Tx, numAlgos)
// 	var coinbaseSigOpCost int64
// 	// blockTemplates := make(map[int32]*BlockTemplate, numAlgos)
// 	var priorityQueues *txPriorityQueue
// 	// Get the current source transactions and create a priority queue to hold the
// 	// transactions which are ready for inclusion into a block along with some
// 	// priority related and fee metadata. Reserve the same number of items that are
// 	// available for the priority queue. Also, choose the initial txsort order for the
// 	// priority queue based on whether or not there is an area allocated for
// 	// high-priority transactions.
// 	sourceTxns := g.TxSource.MiningDescs()
// 	sortedByFee := g.Policy.BlockPrioritySize == 0
// 	blockUtxos := blockchain.NewUtxoViewpoint()
// 	// dependers is used to track transactions which depend on another transaction
// 	// in the source pool. This, in conjunction with the dependsOn map kept with
// 	// each dependent transaction helps quickly determine which dependent
// 	// transactions are now eligible for inclusion in the block once each
// 	// transaction has been included.
// 	dependers := make(map[chainhash.Hash]map[chainhash.Hash]*txPrioItem)
// mempoolLoop:
// 	for _, txDesc := range sourceTxns {
// 		// A block can't have more than one coinbase or contain non-finalized
// 		// transactions.
// 		tx := txDesc.Tx
// 		if blockchain.IsCoinBase(tx) {
// 			Tracec(
// 				func() string {
// 					return fmt.Sprintf("skipping coinbase tx %s", tx.Hash())
// 				},
// 			)
// 			continue
// 		}
// 		if !blockchain.IsFinalizedTransaction(
// 			tx, mbt.Height,
// 			g.TimeSource.AdjustedTime(),
// 		) {
// 			Tracec(
// 				func() string {
// 					return "skipping non-finalized tx " + tx.Hash().String()
// 				},
// 			)
// 			continue
// 		}
// 		// Fetch all of the utxos referenced by the this transaction.
// 		//
// 		// NOTE: This intentionally does not fetch inputs from the mempool since a
// 		// transaction which depends on other transactions in the mempool must come
// 		// after those dependencies in the final generated block.
// 		utxos, e := g.Chain.FetchUtxoView(tx)
// 		if e != nil  {
// 			Warnc(
// 				func() string {
// 					return "unable to fetch utxo view for tx " + tx.Hash().String() + ": " + err.Error()
// 				},
// 			)
// 			continue
// 		}
// 		// Setup dependencies for any transactions which reference other transactions in
// 		// the mempool so they can be properly ordered below.
// 		prioItem := &txPrioItem{tx: tx}
// 		for _, txIn := range tx.MsgTx().TxIn {
// 			originHash := &txIn.PreviousOutPoint.Hash
// 			entry := utxos.LookupEntry(txIn.PreviousOutPoint)
// 			if entry == nil || entry.IsSpent() {
// 				if !g.TxSource.HaveTransaction(originHash) {
// 					Tracec(
// 						func() string {
// 							return "skipping tx %s because it references unspent output %s which is not available" +
// 								tx.Hash().String() +
// 								txIn.PreviousOutPoint.String()
// 						},
// 					)
// 					continue mempoolLoop
// 				}
// 				// The transaction is referencing another transaction in the source pool, so
// 				// setup an ordering dependency.
// 				deps, exists := dependers[*originHash]
// 				if !exists {
// 					deps = make(map[chainhash.Hash]*txPrioItem)
// 					dependers[*originHash] = deps
// 				}
// 				deps[*prioItem.tx.Hash()] = prioItem
// 				if prioItem.dependsOn == nil {
// 					prioItem.dependsOn = make(
// 						map[chainhash.Hash]struct{},
// 					)
// 				}
// 				prioItem.dependsOn[*originHash] = struct{}{}
// 				// Skip the check below. We already know the referenced transaction is
// 				// available.
// 				continue
// 			}
// 		}
// 		// Calculate the final transaction priority using the input value age sum as
// 		// well as the adjusted transaction size. The formula is: sum (inputValue *
// 		// inputAge) / adjustedTxSize
// 		prioItem.priority = CalcPriority(
// 			tx.MsgTx(), utxos,
// 			mbt.Height,
// 		)
// 		// Calculate the fee in Satoshi/kB.
// 		prioItem.feePerKB = txDesc.FeePerKB
// 		prioItem.fee = txDesc.Fee
// 		// Add the transaction to the priority queue to mark it ready for inclusion in
// 		// the block unless it has dependencies.
// 		if prioItem.dependsOn == nil {
// 			heap.Push(priorityQueues, prioItem)
// 		}
// 		// Merge the referenced outputs from the input transactions to this transaction
// 		// into the block utxo view. This allows the code below to avoid a second
// 		// lookup.
// 		mergeUtxoView(blockUtxos, utxos)
// 	}
// 	priorityQueues = newTxPriorityQueue(len(sourceTxns), sortedByFee)
// 	var coinbaseScript []byte
// 	if coinbaseScript, e = standardCoinbaseScript(mbt.Height, extraNonce); E.Chk(e) {
// 		return nil, e
// 	}
// 	algos := fork.GetAlgos(mbt.Height)
// 	var alg int32
// 	mbt.coinbases = make(map[int32]*util.Tx)
// 	// Create a slice to hold the transactions to be included in the generated block
// 	var coinbaseTx *util.Tx
// 	for i := range algos {
// 		alg = algos[i].Version
// 		if coinbaseTx, e = createCoinbaseTx(
// 			g.ChainParams, coinbaseScript, mbt.Height, payToAddress, alg,
// 		); E.Chk(e) {
// 			return nil, e
// 		}
// 		mbt.coinbases[alg] = coinbaseTx
// 		// this should be the same for all anyhow, as they are all the same format just
// 		// diff amounts (note: this might be wrawwnnggrrr)
// 		coinbaseSigOpCost = int64(blockchain.CountSigOps(mbt.coinbases[alg]))
// 	}
// 	// Create slices to hold the fees and number of signature operations for each of
// 	// the selected transactions and add an entry for the coinbase. This allows the
// 	// code below to simply append details about a transaction as it is selected for
// 	// inclusion in the final block. However, since the total fees aren't known yet,
// 	// use a dummy value for the coinbase fee which will be updated later.
// 	txFees := make([]int64, 0, len(sourceTxns))
// 	txSigOpCosts := make([]int64, 0, len(sourceTxns))
// 	txFees = append(txFees, -1) // Updated once known
// 	txSigOpCosts = append(txSigOpCosts, coinbaseSigOpCost)
// 	Tracef("considering %d transactions for inclusion to new block", len(sourceTxns))
// 	// The starting block size is the size of the block header plus the max possible
// 	// transaction count size, plus the size of the coinbase transaction.
// 	// with reserved space. Also create a utxo view to house all of the input
// 	// transactions so multiple lookups can be avoided.
// 	blockWeight := uint32((blockHeaderOverhead) + blockchain.GetTransactionWeight(coinbaseTx))
// 	blockSigOpCost := coinbaseSigOpCost
// 	totalFees := int64(0)
// 	mbt.txs = make([]*util.Tx, 0, len(sourceTxns))
// 	// Choose which transactions make it into the block.
// 	for priorityQueues.Len() > 0 {
// 		// Grab the highest priority (or highest fee per kilobyte depending on the txsort
// 		// order) transaction.
// 		prioItem := heap.Pop(priorityQueues).(*txPrioItem)
// 		tx := prioItem.tx
// 		// Grab any transactions which depend on this one.
// 		deps := dependers[*tx.Hash()]
// 		// Enforce maximum block size.  Also check for overflow.
// 		txWeight := uint32(blockchain.GetTransactionWeight(tx))
// 		blockPlusTxWeight := blockWeight + txWeight
// 		if blockPlusTxWeight < blockWeight ||
// 			blockPlusTxWeight >= g.Policy.BlockMaxWeight {
// 			Tracef("skipping tx %s because it would exceed the max block weight", tx.Hash())
// 			logSkippedDeps(tx, deps)
// 			continue
// 		}
// 		// Enforce maximum signature operation cost per block. Also check for overflow.
// 		sigOpCost, e := blockchain.GetSigOpCost(tx, false, blockUtxos, true, false)
// 		if e != nil  {
// 			Tracec(
// 				func() string {
// 					return "skipping tx " + tx.Hash().String() +
// 						"due to error in GetSigOpCost: " + err.Error()
// 				},
// 			)
// 			logSkippedDeps(tx, deps)
// 			continue
// 		}
// 		if blockSigOpCost+int64(sigOpCost) < blockSigOpCost ||
// 			blockSigOpCost+int64(sigOpCost) > blockchain.MaxBlockSigOpsCost {
// 			Tracec(
// 				func() string {
// 					return "skipping tx " + tx.Hash().String() +
// 						" because it would exceed the maximum sigops per block"
// 				},
// 			)
// 			logSkippedDeps(tx, deps)
// 			continue
// 		}
// 		// Skip free transactions once the block is larger than the minimum block size.
// 		if sortedByFee &&
// 			prioItem.feePerKB < int64(g.Policy.TxMinFreeFee) &&
// 			blockPlusTxWeight >= g.Policy.BlockMinWeight {
// 			Tracec(
// 				func() string {
// 					return fmt.Sprintf(
// 						"skipping tx %v with feePerKB %v < TxMinFreeFee %v and block weight %v >= minBlockWeight %v",
// 						tx.Hash(),
// 						prioItem.feePerKB,
// 						g.Policy.TxMinFreeFee,
// 						blockPlusTxWeight,
// 						g.Policy.BlockMinWeight,
// 					)
// 				},
// 			)
// 			logSkippedDeps(tx, deps)
// 			continue
// 		}
// 		// Prioritize by fee per kilobyte once the block is larger than the priority
// 		// size or there are no more high-priority transactions.
// 		if !sortedByFee && (blockPlusTxWeight >= g.Policy.BlockPrioritySize ||
// 			prioItem.priority <= MinHighPriority.ToDUO()) {
// 			Tracef(
// 				"switching to txsort by fees per kilobyte blockSize %d"+
// 					" >= BlockPrioritySize %d || priority %.2f <= minHighPriority %.2f",
// 				blockPlusTxWeight,
// 				g.Policy.BlockPrioritySize,
// 				prioItem.priority,
// 				MinHighPriority,
// 			)
// 			sortedByFee = true
// 			priorityQueues.SetLessFunc(txPQByFee)
// 		}
// 		// Put the transaction back into the priority queue and skip it so it is
// 		// re-prioritized by fees if it won't fit into the high-priority section or the
// 		// priority is too low. Otherwise this transaction will be the final one in the
// 		// high-priority section, so just fall though to the code below so it is added
// 		// now.
// 		if blockPlusTxWeight > g.Policy.BlockPrioritySize ||
// 			prioItem.priority < MinHighPriority.ToDUO() {
// 			heap.Push(priorityQueues, prioItem)
// 			continue
// 		}
//
// 		// Ensure the transaction inputs pass all of the necessary preconditions before
// 		// allowing it to be added to the block.
// 		_, e = blockchain.CheckTransactionInputs(
// 			tx, mbt.Height,
// 			blockUtxos, g.ChainParams,
// 		)
// 		if e != nil  {
// 			Tracef(
// 				"skipping tx %s due to error in CheckTransactionInputs: %v",
// 				tx.Hash(), e,
// 			)
// 			logSkippedDeps(tx, deps)
// 			continue
// 		}
// 		if e = blockchain.ValidateTransactionScripts(
// 			g.Chain, tx, blockUtxos,
// 			txscript.StandardVerifyFlags, g.SigCache,
// 			g.HashCache,
// 		); E.Chk(e) {
// 			Tracef(
// 				"skipping tx %s due to error in ValidateTransactionScripts: %v",
// 				tx.Hash(), e,
// 			)
// 			logSkippedDeps(tx, deps)
// 			continue
// 		}
// 		// Spend the transaction inputs in the block utxo view and add an entry for it
// 		// to ensure any transactions which reference this one have it available as an
// 		// input and can ensure they aren't double spending.
// 		if e = spendTransaction(blockUtxos, tx, mbt.Height); E.Chk(e) {
// 		}
// 		// Add the transaction to the block, increment counters, and save the fees and
// 		// signature operation counts to the block template.
// 		mbt.txs = append(mbt.txs, tx)
// 		blockWeight += txWeight
// 		blockSigOpCost += int64(sigOpCost)
// 		totalFees += prioItem.fee
// 		txFees = append(txFees, prioItem.fee)
// 		txSigOpCosts = append(txSigOpCosts, int64(sigOpCost))
// 		Tracef(
// 			"adding tx %s (priority %.2f, feePerKB %.2f)",
// 			prioItem.tx.Hash(),
// 			prioItem.priority,
// 			prioItem.feePerKB,
// 		)
// 		// Add transactions which depend on this one (and also do not have any other
// 		// unsatisfied dependencies) to the priority queue.
// 		for _, item := range deps {
// 			// Add the transaction to the priority queue if there are no more dependencies
// 			// after this one.
// 			delete(item.dependsOn, *tx.Hash())
// 			if len(item.dependsOn) == 0 {
// 				heap.Push(priorityQueues, item)
// 			}
// 		}
// 	}
// 	if fork.GetCurrent(mbt.Height) < 1 {
// 		// for legacy chain this is the consensus timestamp to use, post hard fork there
// 		// is no allowance for less than 1 second between block timestamps of
// 		// sequential, linked blocks, which was filled earlier by default
// 		mbt.Timestamp = medianAdjustedTime(best, g.TimeSource)
// 	}
//
// 	for next, curr, more := fork.AlgoVerIterator(mbt.Height); more(); next() {
// 		tX := append(mbt.txs, mbt.coinbases[curr()])
// 		// Now that the actual transactions have been selected, update the block weight
// 		// for the real transaction count and coinbase value with the total fees
// 		// accordingly.
// 		blockWeight -= wire.MaxVarIntPayload -
// 			(uint32(wire.VarIntSerializeSize(uint64(len(mbt.txs)))))
// 		mbt.coinbases[curr()].MsgTx().TxOut[0].Value += totalFees
// 		txFees[0] = -totalFees
// 		// Calculate the required difficulty for the block. The timestamp is potentially
// 		// adjusted to ensure it comes after the median time of the last several blocks
// 		// per the chain consensus rules.
// 		algo := fork.GetAlgoName(mbt.Height, curr())
// 		D.Ln("algo", algo)
// 		if mbt.Bits[curr()], e = g.Chain.CalcNextRequiredDifficulty(algo); E.Chk(e) {
// 			return nil, e
// 		}
// 		D.F(
// 			"%s %d reqDifficulty %08x %064x", algo, curr(),
// 			mbt.Bits[curr()], fork.CompactToBig(mbt.Bits[curr()]),
// 		)
// 		// Create a new block ready to be solved.
// 		D.S(tX)
//
// 		merkles := blockchain.BuildMerkleTreeStore(tX, false)
// 		mbt.Merkles[curr()] = *merkles[len(merkles)-1]
// 		// TODO: can we do this once instead of 9 times?
// 		var msgBlock wire.WireBlock
// 		msgBlock.Header = wire.BlockHeader{
// 			Version:    curr(),
// 			PrevBlock:  mbt.PrevBlock,
// 			MerkleRoot: mbt.Merkles[curr()],
// 			Timestamp:  mbt.Timestamp,
// 			Bits:       mbt.Bits[curr()],
// 		}
// 		for _, tx := range tX {
// 			if e := msgBlock.AddTransaction(tx.MsgTx()); E.Chk(e) {
// 				return nil, e
// 			}
// 		}
// 		// Finally, perform a full check on the created block against the chain
// 		// consensus rules to ensure it properly connects to the current best chain with
// 		// no issues.
// 		block := util.NewBlock(&msgBlock)
// 		block.SetHeight(mbt.Height)
// 		e = g.Chain.CheckConnectBlockTemplate(workerNumber, block)
// 		if e != nil  {
// 			D.Ln("checkconnectblocktemplate err:", e)
// 			return nil, e
// 		}
// 		Tracec(
// 			func() string {
// 				bh := msgBlock.Header.BlockHash()
// 				return fmt.Sprintf(
// 					"created new block template (algo %s, %d transactions, "+
// 						"%d in fees, %d signature operations cost, %d weight, "+
// 						"target difficulty %064x prevblockhash %064x %064x subsidy %d)",
// 					algo,
// 					len(msgBlock.Transactions),
// 					totalFees,
// 					blockSigOpCost,
// 					blockWeight,
// 					fork.CompactToBig(msgBlock.Header.Bits),
// 					msgBlock.Header.PrevBlock.CloneBytes(),
// 					bh.CloneBytes(),
// 					msgBlock.Transactions[0].TxOut[0].Value,
// 				)
// 			},
// 		)
// 		// // Tracec(func() string { return spew.Sdump(msgBlock) })
// 		// blockTemplate := &BlockTemplate{
// 		// 	Block:           &msgBlock,
// 		// 	Fees:            txFees,
// 		// 	SigOpCosts:      txSigOpCosts,
// 		// 	Height:          mbt.Height,
// 		// 	ValidPayAddress: payToAddress != nil,
// 		// }
// 		// blockTemplates[curr()] = blockTemplate
// 	}
// 	return mbt, nil
// }
