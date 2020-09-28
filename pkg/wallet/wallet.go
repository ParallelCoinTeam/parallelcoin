package wallet

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/p9c/pkg/app/slog"

	"github.com/davecgh/go-spew/spew"

	blockchain "github.com/p9c/pod/pkg/chain"
	"github.com/p9c/pod/pkg/chain/config/netparams"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	txauthor "github.com/p9c/pod/pkg/chain/tx/author"
	wtxmgr "github.com/p9c/pod/pkg/chain/tx/mgr"
	txrules "github.com/p9c/pod/pkg/chain/tx/rules"
	txscript "github.com/p9c/pod/pkg/chain/tx/script"
	"github.com/p9c/pod/pkg/chain/wire"
	ec "github.com/p9c/pod/pkg/coding/elliptic"
	"github.com/p9c/pod/pkg/db/walletdb"
	"github.com/p9c/pod/pkg/pod"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	rpcclient "github.com/p9c/pod/pkg/rpc/client"
	"github.com/p9c/pod/pkg/util"
	"github.com/p9c/pod/pkg/util/hdkeychain"
	waddrmgr "github.com/p9c/pod/pkg/wallet/addrmgr"
	"github.com/p9c/pod/pkg/wallet/chain"
)

const (
	// InsecurePubPassphrase is the default outer encryption passphrase used
	// for public data (everything but private keys).  Using a non-default
	// public passphrase can prevent an attacker without the public
	// passphrase from discovering all past and future wallet addresses if
	// they gain access to the wallet database.
	//
	// NOTE: at time of writing, public encryption only applies to public
	// data in the waddrmgr namespace.  Transactions are not yet encrypted.
	InsecurePubPassphrase = ""
	// walletDbWatchingOnlyName = "wowallet.db"
	// recoveryBatchSize is the default number of blocks that will be
	// scanned successively by the recovery manager, in the event that the
	// wallet is started in recovery mode.
	recoveryBatchSize = 2000
)

// ErrNotSynced describes an error where an operation cannot complete
// due wallet being out of sync (and perhaps currently syncing with)
// the remote chain server.
var ErrNotSynced = errors.New("wallet is not synchronized with the chain server")

// Namespace bucket keys.
var (
	waddrmgrNamespaceKey = []byte("waddrmgr")
	wtxmgrNamespaceKey   = []byte("wtxmgr")
)

// Wallet is a structure containing all the components for a
// complete wallet.  It contains the Armory-style key store
// addresses and keys),
type Wallet struct {
	publicPassphrase []byte
	// Data stores
	db                 walletdb.DB
	Manager            *waddrmgr.Manager
	TxStore            *wtxmgr.Store
	chainClient        chain.Interface
	chainClientLock    sync.Mutex
	chainClientSynced  bool
	chainClientSyncMtx sync.Mutex
	lockedOutpoints    map[wire.OutPoint]struct{}
	recoveryWindow     uint32
	// Channels for rescan processing.  Requests are added and merged with
	// any waiting requests, before being sent to another goroutine to
	// call the rescan RPC.
	rescanAddJob        chan *RescanJob
	rescanBatch         chan *rescanBatch
	rescanNotifications chan interface{} // From chain server
	rescanProgress      chan *RescanProgressMsg
	rescanFinished      chan *RescanFinishedMsg
	// Channel for transaction creation requests.
	createTxRequests chan createTxRequest
	// Channels for the manager locker.
	unlockRequests     chan unlockRequest
	lockRequests       chan struct{}
	holdUnlockRequests chan chan heldUnlock
	lockState          chan bool
	changePassphrase   chan changePassphraseRequest
	changePassphrases  chan changePassphrasesRequest
	// Information for reorganization handling.
	// reorganizingLock sync.Mutex
	// reorganizeToHash chainhash.Hash
	// reorganizing     bool
	NtfnServer  *NotificationServer
	PodConfig   *pod.Config
	chainParams *netparams.Params
	wg          sync.WaitGroup
	started     bool
	quit        chan struct{}
	quitMu      sync.Mutex
	Update      chan struct{}
}

// Start starts the goroutines necessary to manage a wallet.
func (w *Wallet) Start() {
	slog.Trace("starting wallet")
	w.quitMu.Lock()
	select {
	case <-w.quit:
		slog.Trace("waiting for wallet shutdown")
		// Restart the wallet goroutines after shutdown finishes.
		w.WaitForShutdown()
		w.quit = make(chan struct{})
	default:
		if w.started {
			// Ignore when the wallet is still running.
			slog.Info("wallet already started")
			w.quitMu.Unlock()
			return
		}
		w.started = true
	}
	w.quitMu.Unlock()
	w.wg.Add(2)
	go w.txCreator()
	go w.walletLocker()
}

// SynchronizeRPC associates the wallet with the consensus RPC client,
// synchronizes the wallet with the latest changes to the blockchain, and
// continuously updates the wallet through RPC notifications.
//
// This method is unstable and will be removed when all syncing logic is moved
// outside of the wallet package.
func (w *Wallet) SynchronizeRPC(chainClient chain.Interface) {
	w.quitMu.Lock()
	select {
	case <-w.quit:
		w.quitMu.Unlock()
		return
	default:
	}
	w.quitMu.Unlock()
	// TODO: Ignoring the new client when one is already set breaks callers
	// who are replacing the client, perhaps after a disconnect.
	w.chainClientLock.Lock()
	if w.chainClient != nil {
		w.chainClientLock.Unlock()
		return
	}
	w.chainClient = chainClient
	// If the chain client is a NeutrinoClient instance, set a birthday so
	// we don't download all the filters as we go.
	switch cc := chainClient.(type) {
	case *chain.NeutrinoClient:
		cc.SetStartTime(w.Manager.Birthday())
	case *chain.BitcoindClient:
		cc.SetBirthday(w.Manager.Birthday())
	}
	w.chainClientLock.Unlock()
	// TODO: It would be preferable to either run these goroutines separately from the wallet (use wallet mutator
	//  functions to make changes from the RPC client) and not have to stop and restart them each time the client
	//  disconnects and reconnets.
	w.wg.Add(4)
	go w.handleChainNotifications()
	go w.rescanBatchHandler()
	go w.rescanProgressHandler()
	go w.rescanRPCHandler()
}

// requireChainClient marks that a wallet method can only be completed when the consensus RPC server is set. This
// function and all functions that call it are unstable and will need to be moved when the syncing code is moved out of
// the wallet.
func (w *Wallet) requireChainClient() (chainClient chain.Interface, err error) {
	w.chainClientLock.Lock()
	chainClient = w.chainClient
	w.chainClientLock.Unlock()
	if chainClient == nil {
		err = errors.New("blockchain RPC is inactive")
		slog.Debug(err)
		return
	}
	return
}

// ChainClientChan returns the optional consensus RPC client associated with the wallet.
//
// This function is unstable and will be removed once sync logic is moved out of the wallet.
func (w *Wallet) ChainClient() chain.Interface {
	w.chainClientLock.Lock()
	chainClient := w.chainClient
	w.chainClientLock.Unlock()
	return chainClient
}

// quitChan atomically reads the quit channel.
func (w *Wallet) quitChan() <-chan struct{} {
	w.quitMu.Lock()
	c := w.quit
	w.quitMu.Unlock()
	return c
}

// Stop signals all wallet goroutines to shutdown.
func (w *Wallet) Stop() {
	w.quitMu.Lock()
	quit := w.quit
	w.quitMu.Unlock()
	select {
	case <-quit:
	default:
		close(quit)
		w.chainClientLock.Lock()
		if w.chainClient != nil {
			w.chainClient.Stop()
			w.chainClient = nil
		}
		w.chainClientLock.Unlock()
	}
}

// ShuttingDown returns whether the wallet is currently in the process of shutting down or not.
func (w *Wallet) ShuttingDown() bool {
	select {
	case <-w.quitChan():
		return true
	default:
		return false
	}
}

// WaitForShutdown blocks until all wallet goroutines have finished executing.
func (w *Wallet) WaitForShutdown() {
	w.chainClientLock.Lock()
	if w.chainClient != nil {
		w.chainClient.WaitForShutdown()
	}
	w.chainClientLock.Unlock()
	w.wg.Wait()
}

// SynchronizingToNetwork returns whether the wallet is currently synchronizing with the Bitcoin network.
func (w *Wallet) SynchronizingToNetwork() bool {
	// At the moment, RPC is the only synchronization method. In the future, when SPV is added, a separate check will
	// also be needed, or SPV could always be enabled if RPC was not explicitly specified when creating the wallet.
	w.chainClientSyncMtx.Lock()
	syncing := w.chainClient != nil
	w.chainClientSyncMtx.Unlock()
	return syncing
}

// ChainSynced returns whether the wallet has been attached to a chain server and synced up to the best block on the
// main chain.
func (w *Wallet) ChainSynced() bool {
	w.chainClientSyncMtx.Lock()
	synced := w.chainClientSynced
	w.chainClientSyncMtx.Unlock()
	return synced
}

// SetChainSynced marks whether the wallet is connected to and currently in sync with the latest block notified by the
// chain server.
//
// NOTE: Due to an API limitation with rpcclient, this may return true after the client disconnected (and is attempting
// a reconnect). This will be unknown until the reconnect notification is received, at which point the wallet can be
// marked out of sync again until after the next rescan completes.
func (w *Wallet) SetChainSynced(synced bool) {
	w.chainClientSyncMtx.Lock()
	w.chainClientSynced = synced
	w.chainClientSyncMtx.Unlock()
}

// activeData returns the currently-active receiving addresses and all unspent outputs. This is primarily intended to
// provide the parameters for a rescan request.
func (w *Wallet) activeData(dbtx walletdb.ReadTx) (addrs []util.Address, unspent []wtxmgr.Credit, err error) {
	addrmgrNs := dbtx.ReadBucket(waddrmgrNamespaceKey)
	txmgrNs := dbtx.ReadBucket(wtxmgrNamespaceKey)
	if err = w.Manager.ForEachActiveAddress(addrmgrNs, func(addr util.Address) (err error) {
		addrs = append(addrs, addr)
		return
	}); slog.Check(err) {
		return
	}
	if unspent, err = w.TxStore.UnspentOutputs(txmgrNs); slog.Check(err) {
	}
	return
}

// syncWithChain brings the wallet up to date with the current chain server connection. It creates a rescan request and
// blocks until the rescan has finished.
func (w *Wallet) syncWithChain() (err error) {
	var chainClient chain.Interface
	if chainClient, err = w.requireChainClient(); slog.Check(err) {
		return
	}
	// Request notifications for transactions sending to all wallet addresses.
	var (
		addrs   []util.Address
		unspent []wtxmgr.Credit
	)
	if err = walletdb.View(w.db, func(dbtx walletdb.ReadTx) (err error) {
		addrs, unspent, err = w.activeData(dbtx)
		return
	}); slog.Check(err) {
		slog.Warn("error starting sync", err)
		return
	}
	startHeight := w.Manager.SyncedTo().Height
	// We'll mark this as our first sync if we don't have any unspent outputs as known by the wallet. This'll allow us
	// to skip a full rescan at this height, and instead wait for the backend to catch up.
	isInitialSync := len(unspent) == 0
	isRecovery := w.recoveryWindow > 0
	birthday := w.Manager.Birthday()
	// If an initial sync is attempted, we will try and find the block stamp of the first block past our birthday. This
	// will be fed into the rescan to ensure we catch transactions that are sent while performing the initial sync.
	var birthdayStamp *waddrmgr.BlockStamp
	// TODO(jrick): How should this handle a synced height earlier than the chain server best block?
	// When no addresses have been generated for the wallet, the rescan can be skipped.
	//
	// TODO: This is only correct because activeData above returns all
	//  addresses ever created, including those that don't need to be watched
	//  anymore.  This code should be updated when this assumption is no
	//  longer true, but worst case would result in an unnecessary rescan.
	if isInitialSync || isRecovery {
		// Find the latest checkpoint's height. This lets us catch up to at least that checkpoint, since we're
		// synchronizing from scratch, and lets us avoid a bunch of costly DB transactions in the case when we're using
		// BDB for the walletdb backend and Neutrino for the chain.Interface backend, and the chain backend starts
		// synchronizing at the same time as the wallet.
		var bestHeight int32
		if _, bestHeight, err = chainClient.GetBestBlock(); slog.Check(err) {
			return
		}
		slog.Debug("bestHeight", bestHeight)
		checkHeight := bestHeight
		if len(w.chainParams.Checkpoints) > 0 {
			checkHeight = w.chainParams.Checkpoints[len(w.chainParams.Checkpoints)-1].Height
		}
		logHeight := checkHeight
		if bestHeight > logHeight {
			logHeight = bestHeight
		}
		slog.Infof("catching up block hashes to height %d, this will take a while", logHeight)
		// Initialize the first database transaction.
		var tx walletdb.ReadWriteTx
		if tx, err = w.db.BeginReadWriteTx(); slog.Check(err) {
			return
		}
		ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
		// Only allocate the recoveryMgr if we are actually in recovery mode.
		var recoveryMgr *RecoveryManager
		if isRecovery {
			slog.Info("RECOVERY MODE ENABLED -- rescanning for used addresses with recovery_window =", w.recoveryWindow)
			// Initialize the recovery manager with a default batch size of 2000.
			recoveryMgr = NewRecoveryManager(w.recoveryWindow, recoveryBatchSize, w.chainParams)
			// In the event that this recovery is being resumed, we will need to repopulate all found addresses from the
			// database. For basic recovery, we will only do so for the default scopes.
			var scopedMgrs map[waddrmgr.KeyScope]*waddrmgr.ScopedKeyManager
			if scopedMgrs, err = w.defaultScopeManagers(); slog.Check(err) {
				return
			}
			txmgrNs := tx.ReadBucket(wtxmgrNamespaceKey)
			var credits []wtxmgr.Credit
			if credits, err = w.TxStore.UnspentOutputs(txmgrNs); slog.Check(err) {
				return
			}
			if err = recoveryMgr.Resurrect(ns, scopedMgrs, credits); slog.Check(err) {
				return
			}
		}
		var hash *chainhash.Hash
		var header *wire.BlockHeader
		for height := startHeight; height <= bestHeight; height++ {
			if hash, err = chainClient.GetBlockHash(int64(height)); slog.Check(err) {
				if err = tx.Rollback(); slog.Check(err) {
				}
				return
			}
			// If we're using the Neutrino backend, we can check if it's current or not. For other backends we'll assume
			// it is current if the best height has reached the last checkpoint.
			isCurrent := func(bestHeight int32) bool {
				switch c := chainClient.(type) {
				case *chain.NeutrinoClient:
					return c.CS.IsCurrent()
				}
				return bestHeight >= checkHeight
			}
			// If we've found the best height the backend knows about, and the backend is still synchronizing, we'll
			// wait. We can give it a little bit of time to synchronize further before updating the best height based on
			// the backend. Once we see that the backend has advanced, we can catch up to it.
			for height == bestHeight && !isCurrent(bestHeight) {
				time.Sleep(100 * time.Millisecond)
				if _, bestHeight, err = chainClient.GetBestBlock(); slog.Check(err) {
					if err = tx.Rollback(); slog.Check(err) {
					}
					return
				}
			}
			if header, err = chainClient.GetBlockHeader(hash); slog.Check(err) {
				return
			}
			// Check to see if this header's timestamp has surpassed our birthday or if we've surpassed one previously.
			timestamp := header.Timestamp
			if timestamp.After(birthday) || birthdayStamp != nil {
				// If this is the first block past our birthday, record the block stamp so that we can use this as the
				// starting point for the rescan. This will ensure we don't miss transactions that are sent to the
				// wallet during an initial sync.
				//
				// NOTE: The birthday persisted by the wallet is two days before the actual wallet birthday, to deal
				// with potentially inaccurate header timestamps.
				if birthdayStamp == nil {
					birthdayStamp = &waddrmgr.BlockStamp{
						Height:    height,
						Hash:      *hash,
						Timestamp: timestamp,
					}
				}
				// If we are in recovery mode and the check passes, we will add this block to our list of blocks to scan
				// for recovered addresses.
				if isRecovery {
					recoveryMgr.AddToBlockBatch(
						hash, height, timestamp,
					)
				}
			}
			if err = w.Manager.SetSyncedTo(ns, &waddrmgr.BlockStamp{
				Hash:      *hash,
				Height:    height,
				Timestamp: timestamp,
			}); slog.Check(err) {
				if err = tx.Rollback(); slog.Check(err) {
				}
				return
			}
			// If we are in recovery mode, attempt a recovery on blocks that have been added to the recovery manager's
			// block batch thus far. If block batch is empty, this will be a NOP.
			if isRecovery && height%recoveryBatchSize == 0 {
				if err = w.recoverDefaultScopes(
					chainClient, tx, ns,
					recoveryMgr.BlockBatch(),
					recoveryMgr.State(),
				); slog.Check(err) {
					if err = tx.Rollback(); slog.Check(err) {
					}
					return
				}
				// Clear the batch of all processed blocks.
				recoveryMgr.ResetBlockBatch()
			}
			// Every 10K blocks, commit and start a new database TX.
			if height%10000 == 0 {
				if err = tx.Commit(); slog.Check(err) {
					if err = tx.Rollback(); slog.Check(err) {
					}
					return
				}
				slog.Info("caught up to height", height)
				if tx, err = w.db.BeginReadWriteTx(); slog.Check(err) {
					return
				}
				ns = tx.ReadWriteBucket(waddrmgrNamespaceKey)
			}
		}
		// Perform one last recovery attempt for all blocks that were not batched at the default granularity of 2000
		// blocks.
		if isRecovery {
			if err = w.recoverDefaultScopes(chainClient, tx, ns, recoveryMgr.BlockBatch(), recoveryMgr.State()); slog.Check(err) {
				if err = tx.Rollback(); slog.Check(err) {
				}
				return
			}
		}
		// Commit (or roll back) the final database transaction.
		if err = tx.Commit(); slog.Check(err) {
			if err = tx.Rollback(); slog.Check(err) {
			}
			return
		}
		slog.Info("done catching up block hashes")
		// Since we've spent some time catching up block hashes, we might have new addresses waiting for us that were
		// requested during initial sync. Make sure we have those before we request a rescan later on.
		if err = walletdb.View(w.db, func(dbtx walletdb.ReadTx) (err error) {
			addrs, unspent, err = w.activeData(dbtx)
			return
		}); slog.Check(err) {
			return
		}
	}
	// Compare previously-seen blocks against the chain server. If any of these blocks no longer exist, rollback all of
	// the missing blocks before catching up with the rescan.
	var rollback bool
	rollbackStamp := w.Manager.SyncedTo()
	if err = walletdb.Update(w.db, func(tx walletdb.ReadWriteTx) (err error) {
		addrmgrNs := tx.ReadWriteBucket(waddrmgrNamespaceKey)
		txmgrNs := tx.ReadWriteBucket(wtxmgrNamespaceKey)
		var hash *chainhash.Hash
		var chainHash *chainhash.Hash
		var header *wire.BlockHeader
		for height := rollbackStamp.Height; true; height-- {
			if hash, err = w.Manager.BlockHash(addrmgrNs, height); slog.Check(err) {
				return
			}
			if chainHash, err = chainClient.GetBlockHash(int64(height)); slog.Check(err) {
				return
			}
			if header, err = chainClient.GetBlockHeader(chainHash); slog.Check(err) {
				return
			}
			rollbackStamp.Hash = *chainHash
			rollbackStamp.Height = height
			rollbackStamp.Timestamp = header.Timestamp
			if bytes.Equal(hash[:], chainHash[:]) {
				break
			}
			rollback = true
		}
		if rollback {
			if err = w.Manager.SetSyncedTo(addrmgrNs, &rollbackStamp); slog.Check(err) {
				return
			}
			// Rollback unconfirms transactions at and beyond the passed height, so add one to the new synced-to height
			// to prevent unconfirming txs from the synced-to block.
			if err = w.TxStore.Rollback(txmgrNs, rollbackStamp.Height+1); slog.Check(err) {
				return
			}
		}
		return
	}); slog.Check(err) {
		return
	}
	// If a birthday stamp was found during the initial sync and the rollback causes us to revert it, update the
	// birthday stamp so that it points at the new tip.
	if birthdayStamp != nil && rollbackStamp.Height <= birthdayStamp.Height {
		birthdayStamp = &rollbackStamp
	}
	// Request notifications for connected and disconnected blocks.
	//
	// TODO(jrick): Either request this notification only once, or when
	//  rpcclient is modified to allow some notification request to not
	//  automatically resent on reconnect, include the notifyblocks request
	//  as well.  I am leaning towards allowing off all rpcclient
	//  notification re-registrations, in which case the code here should be
	//  left as is.
	err = chainClient.NotifyBlocks()
	if err != nil {
		slog.Error(err)
		return err
	}
	return w.rescanWithTarget(addrs, unspent, birthdayStamp)
}

// defaultScopeManagers fetches the ScopedKeyManagers from the wallet using the default set of key scopes.
func (w *Wallet) defaultScopeManagers() (scopedMgrs map[waddrmgr.KeyScope]*waddrmgr.ScopedKeyManager, err error) {
	scopedMgrs = make(map[waddrmgr.KeyScope]*waddrmgr.ScopedKeyManager)
	var scopedMgr *waddrmgr.ScopedKeyManager
	for _, scope := range waddrmgr.DefaultKeyScopes {
		if scopedMgr, err = w.Manager.FetchScopedKeyManager(scope); slog.Check(err) {
			return
		}
		scopedMgrs[scope] = scopedMgr
	}
	return
}

// recoverDefaultScopes attempts to recover any addresses belonging to any active scoped key managers known to the
// wallet. Recovery of each scope's default account will be done iteratively against the same batch of blocks.
// TODO(conner): parallelize/pipeline/cache intermediate network requests
func (w *Wallet) recoverDefaultScopes(
	chainClient chain.Interface,
	tx walletdb.ReadWriteTx,
	ns walletdb.ReadWriteBucket,
	batch []wtxmgr.BlockMeta,
	recoveryState *RecoveryState,
) (err error) {
	var scopedMgrs map[waddrmgr.KeyScope]*waddrmgr.ScopedKeyManager
	if scopedMgrs, err = w.defaultScopeManagers(); slog.Check(err) {
		return
	}
	return w.recoverScopedAddresses(chainClient, tx, ns, batch, recoveryState, scopedMgrs)
}

// recoverAccountAddresses scans a range of blocks in attempts to recover any previously used addresses for a particular
// account derivation path. At a high level, the algorithm works as follows:
//
//  1) Ensure internal and external branch horizons are fully expanded.
//  2) Filter the entire range of blocks, stopping if a non-zero number of address are contained in a particular block.
//  3) Record all internal and external addresses found in the block.
//  4) Record any outpoints found in the block that should be watched for spends
//  5) Trim the range of blocks up to and including the one reporting the addrs.
//  6) Repeat from (1) if there are still more blocks in the range.
func (w *Wallet) recoverScopedAddresses(
	chainClient chain.Interface,
	tx walletdb.ReadWriteTx,
	ns walletdb.ReadWriteBucket,
	batch []wtxmgr.BlockMeta,
	recoveryState *RecoveryState,
	scopedMgrs map[waddrmgr.KeyScope]*waddrmgr.ScopedKeyManager,
) (err error) {
	// If there are no blocks in the batch, we are done.
	if len(batch) == 0 {
		return
	}
	slog.Infof("scanning %d blocks for recoverable addresses", len(batch))
expandHorizons:
	for scope, scopedMgr := range scopedMgrs {
		scopeState := recoveryState.StateForScope(scope)
		if err = expandScopeHorizons(ns, scopedMgr, scopeState); slog.Check(err) {
			return
		}
	}
	// With the internal and external horizons properly expanded, we now construct the filter blocks request. The
	// request includes the range of blocks we intend to scan, in addition to the scope-index -> addr map for all
	// internal and external branches.
	filterReq := newFilterBlocksRequest(batch, scopedMgrs, recoveryState)
	// Initiate the filter blocks request using our chain backend. If an error occurs, we are unable to proceed with the
	// recovery.
	var filterResp *chain.FilterBlocksResponse
	if filterResp, err = chainClient.FilterBlocks(filterReq); slog.Check(err) {
		return
	}
	// If the filter response is empty, this signals that the rest of the batch was completed, and no other addresses
	// were discovered. As a result, no further modifications to our recovery state are required and we can proceed to
	// the next batch.
	if filterResp == nil {
		return
	}
	// Otherwise, retrieve the block info for the block that detected a non-zero number of address matches.
	block := batch[filterResp.BatchIndex]
	// Log any non-trivial findings of addresses or outpoints.
	logFilterBlocksResp(block, filterResp)
	// Report any external or internal addresses found as a result of the appropriate branch recovery state. Adding
	// indexes above the last-found index of either will result in the horizons being expanded upon the next iteration.
	// Any found addresses are also marked used using the scoped key manager.
	if err = extendFoundAddresses(ns, filterResp, scopedMgrs, recoveryState); slog.Check(err) {
		return
	}
	// Update the global set of watched outpoints with any that were found in the block.
	for outPoint, addr := range filterResp.FoundOutPoints {
		recoveryState.AddWatchedOutPoint(&outPoint, addr)
	}
	// Finally, record all of the relevant transactions that were returned in the filter blocks response. This ensures
	// that these transactions and their outputs are tracked when the final rescan is performed.
	var txRecord *wtxmgr.TxRecord
	for _, txn := range filterResp.RelevantTxns {
		if txRecord, err = wtxmgr.NewTxRecordFromMsgTx(txn, filterResp.BlockMeta.Time); slog.Check(err) {
			return
		}
		if err = w.addRelevantTx(tx, txRecord, &filterResp.BlockMeta); slog.Check(err) {
			return
		}
	}
	// Update the batch to indicate that we've processed all block through the one that returned found addresses.
	batch = batch[filterResp.BatchIndex+1:]
	// If this was not the last block in the batch, we will repeat the filtering process again after expanding our
	// horizons.
	if len(batch) > 0 {
		goto expandHorizons
	}
	return
}

// expandScopeHorizons ensures that the ScopeRecoveryState has an adequately sized look ahead for both its internal and
// external branches. The keys derived here are added to the scope's recovery state, but do not affect the persistent
// state of the wallet. If any invalid child keys are detected, the horizon will be properly extended such that our
// lookahead always includes the proper number of valid child keys.
func expandScopeHorizons(ns walletdb.ReadWriteBucket, scopedMgr *waddrmgr.ScopedKeyManager,
	scopeState *ScopeRecoveryState) (err error) {
	// Compute the current external horizon and the number of addresses we must derive to ensure we maintain a
	// sufficient recovery window for the external branch.
	exHorizon, exWindow := scopeState.ExternalBranch.ExtendHorizon()
	count, childIndex := uint32(0), exHorizon
	for count < exWindow {
		keyPath := externalKeyPath(childIndex)
		var addr waddrmgr.ManagedAddress
		addr, err = scopedMgr.DeriveFromKeyPath(ns, keyPath)
		switch {
		case err == hdkeychain.ErrInvalidChild:
			// Record the existence of an invalid child with the external branch's recovery state. This also increments
			// the branch's horizon so that it accounts for this skipped child index.
			scopeState.ExternalBranch.MarkInvalidChild(childIndex)
			childIndex++
			continue
		case err != nil:
			return
		}
		// Register the newly generated external address and child index with the external branch recovery state.
		scopeState.ExternalBranch.AddAddr(childIndex, addr.Address())
		childIndex++
		count++
	}
	// Compute the current internal horizon and the number of addresses we must derive to ensure we maintain a
	// sufficient recovery window for the internal branch.
	inHorizon, inWindow := scopeState.InternalBranch.ExtendHorizon()
	count, childIndex = 0, inHorizon
	var addr waddrmgr.ManagedAddress
	for count < inWindow {
		keyPath := internalKeyPath(childIndex)
		addr, err = scopedMgr.DeriveFromKeyPath(ns, keyPath)
		switch {
		case err == hdkeychain.ErrInvalidChild:
			// Record the existence of an invalid child with the internal branch's recovery state. This also increments
			// the branch's horizon so that it accounts for this skipped child index.
			scopeState.InternalBranch.MarkInvalidChild(childIndex)
			childIndex++
			continue
		case err != nil:
			return
		}
		// Register the newly generated internal address and child index with the internal branch recovery state.
		scopeState.InternalBranch.AddAddr(childIndex, addr.Address())
		childIndex++
		count++
	}
	return
}

// externalKeyPath returns the relative external derivation path /0/0/index.
func externalKeyPath(index uint32) waddrmgr.DerivationPath {
	return waddrmgr.DerivationPath{
		Account: waddrmgr.DefaultAccountNum,
		Branch:  waddrmgr.ExternalBranch,
		Index:   index,
	}
}

// internalKeyPath returns the relative internal derivation path /0/1/index.
func internalKeyPath(index uint32) waddrmgr.DerivationPath {
	return waddrmgr.DerivationPath{
		Account: waddrmgr.DefaultAccountNum,
		Branch:  waddrmgr.InternalBranch,
		Index:   index,
	}
}

// newFilterBlocksRequest constructs FilterBlocksRequests using our current block range, scoped managers, and recovery
// state.
func newFilterBlocksRequest(
	batch []wtxmgr.BlockMeta,
	scopedMgrs map[waddrmgr.KeyScope]*waddrmgr.ScopedKeyManager,
	recoveryState *RecoveryState,
) (filterReq *chain.FilterBlocksRequest) {
	filterReq = &chain.FilterBlocksRequest{
		Blocks:           batch,
		ExternalAddrs:    make(map[waddrmgr.ScopedIndex]util.Address),
		InternalAddrs:    make(map[waddrmgr.ScopedIndex]util.Address),
		WatchedOutPoints: recoveryState.WatchedOutPoints(),
	}
	// Populate the external and internal addresses by merging the addresses sets belong to all currently tracked
	// scopes.
	for scope := range scopedMgrs {
		scopeState := recoveryState.StateForScope(scope)
		for index, addr := range scopeState.ExternalBranch.Addrs() {
			scopedIndex := waddrmgr.ScopedIndex{
				Scope: scope,
				Index: index,
			}
			filterReq.ExternalAddrs[scopedIndex] = addr
		}
		for index, addr := range scopeState.InternalBranch.Addrs() {
			scopedIndex := waddrmgr.ScopedIndex{
				Scope: scope,
				Index: index,
			}
			filterReq.InternalAddrs[scopedIndex] = addr
		}
	}
	return
}

// extendFoundAddresses accepts a filter blocks response that contains addresses found on chain, and advances the state
// of all relevant derivation paths to match the highest found child index for each branch.
func extendFoundAddresses(
	ns walletdb.ReadWriteBucket,
	filterResp *chain.FilterBlocksResponse,
	scopedMgrs map[waddrmgr.KeyScope]*waddrmgr.ScopedKeyManager,
	recoveryState *RecoveryState,
) (err error) {
	// Mark all recovered external addresses as used. This will be done only for scopes that reported a non-zero number
	// of external addresses in this block.
	for scope, indexes := range filterResp.FoundExternalAddrs {
		// First, report all external child indexes found for this scope. This ensures that the external last-found
		// index will be updated to include the maximum child index seen thus far.
		scopeState := recoveryState.StateForScope(scope)
		for index := range indexes {
			scopeState.ExternalBranch.ReportFound(index)
		}
		scopedMgr := scopedMgrs[scope]
		// Now, with all found addresses reported, derive and extend all external addresses up to and including the
		// current last found index for this scope.
		exNextUnfound := scopeState.ExternalBranch.NextUnfound()
		exLastFound := exNextUnfound
		if exLastFound > 0 {
			exLastFound--
		}
		if err = scopedMgr.ExtendExternalAddresses(ns, waddrmgr.DefaultAccountNum, exLastFound); slog.Check(err) {
			return
		}
		// Finally, with the scope's addresses extended, we mark used the external addresses that were found in the
		// block and belong to this scope.
		for index := range indexes {
			addr := scopeState.ExternalBranch.GetAddr(index)
			if err = scopedMgr.MarkUsed(ns, addr); slog.Check(err) {
				return
			}
		}
	}
	// Mark all recovered internal addresses as used. This will be done only for scopes that reported a non-zero number
	// of internal addresses in this block.
	for scope, indexes := range filterResp.FoundInternalAddrs {
		// First, report all internal child indexes found for this scope. This ensures that the internal last-found
		// index will be updated to include the maximum child index seen thus far.
		scopeState := recoveryState.StateForScope(scope)
		for index := range indexes {
			scopeState.InternalBranch.ReportFound(index)
		}
		scopedMgr := scopedMgrs[scope]
		// Now, with all found addresses reported, derive and extend all internal addresses up to and including the
		// current last found index for this scope.
		inNextUnfound := scopeState.InternalBranch.NextUnfound()
		inLastFound := inNextUnfound
		if inLastFound > 0 {
			inLastFound--
		}
		if err = scopedMgr.ExtendInternalAddresses(ns, waddrmgr.DefaultAccountNum, inLastFound); slog.Check(err) {
			return
		}
		// Finally, with the scope's addresses extended, we mark used the internal addresses that were found in the
		// blockand belong to this scope.
		for index := range indexes {
			addr := scopeState.InternalBranch.GetAddr(index)
			if err = scopedMgr.MarkUsed(ns, addr); slog.Check(err) {
				return
			}
		}
	}
	return
}

// logFilterBlocksResp provides useful logging information when filtering succeeded in finding relevant transactions.
func logFilterBlocksResp(block wtxmgr.BlockMeta,
	resp *chain.FilterBlocksResponse) {
	// Log the number of external addresses found in this block.
	var nFoundExternal int
	for _, indexes := range resp.FoundExternalAddrs {
		nFoundExternal += len(indexes)
	}
	if nFoundExternal > 0 {
		slog.Tracef("recovered %d external addrs at height=%d hash=%v", nFoundExternal, block.Height, block.Hash)
	}
	// Log the number of internal addresses found in this block.
	var nFoundInternal int
	for _, indexes := range resp.FoundInternalAddrs {
		nFoundInternal += len(indexes)
	}
	if nFoundInternal > 0 {
		slog.Tracef("recovered %d internal addrs at height=%d hash=%v", nFoundInternal, block.Height, block.Hash)
	}
	// Log the number of outpoints found in this block.
	nFoundOutPoints := len(resp.FoundOutPoints)
	if nFoundOutPoints > 0 {
		slog.Tracef("found %d spends from watched outpoints at height=%d hash=%v",
			nFoundOutPoints, block.Height, block.Hash)
	}
}

type (
	createTxRequest struct {
		account     uint32
		outputs     []*wire.TxOut
		minconf     int32
		feeSatPerKB util.Amount
		resp        chan createTxResponse
	}
	createTxResponse struct {
		tx  *txauthor.AuthoredTx
		err error
	}
)

// txCreator is responsible for the input selection and creation of transactions. These functions are the responsibility
// of this method (designed to be run as its own goroutine) since input selection must be serialized, or else it is
// possible to create double spends by choosing the same inputs for multiple transactions. Along with input selection,
// this method is also responsible for the signing of transactions, since we don't want to end up in a situation where
// we run out of inputs as multiple transactions are being created. In this situation, it would then be possible for
// both requests, rather than just one, to fail due to not enough available inputs.
func (w *Wallet) txCreator() {
	quit := w.quitChan()
out:
	for {
		select {
		case txr := <-w.createTxRequests:
			heldUnlock, err := w.holdUnlock()
			if err != nil {
				slog.Error(err)
				txr.resp <- createTxResponse{nil, err}
				continue
			}
			tx, err := w.txToOutputs(txr.outputs, txr.account,
				txr.minconf, txr.feeSatPerKB)
			heldUnlock.release()
			txr.resp <- createTxResponse{tx, err}
		case <-quit:
			break out
		}
	}
	w.wg.Done()
}

// CreateSimpleTx creates a new signed transaction spending unspent P2PKH outputs with at laest minconf confirmations
// spending to any number of address/amount pairs. Change and an appropriate transaction fee are automatically included,
// if necessary. All transaction creation through this function is serialized to prevent the creation of many
// transactions which spend the same outputs.
func (w *Wallet) CreateSimpleTx(account uint32, outputs []*wire.TxOut,
	minconf int32, satPerKb util.Amount) (tx *txauthor.AuthoredTx, err error) {
	req := createTxRequest{
		account:     account,
		outputs:     outputs,
		minconf:     minconf,
		feeSatPerKB: satPerKb,
		resp:        make(chan createTxResponse),
	}
	w.createTxRequests <- req
	resp := <-req.resp
	return resp.tx, resp.err
}

type (
	unlockRequest struct {
		passphrase []byte
		lockAfter  <-chan time.Time // nil prevents the timeout.
		err        chan error
	}
	changePassphraseRequest struct {
		old, new []byte
		private  bool
		err      chan error
	}
	changePassphrasesRequest struct {
		publicOld, publicNew   []byte
		privateOld, privateNew []byte
		err                    chan error
	}
	// heldUnlock is a tool to prevent the wallet from automatically locking after some timeout before an operation
	// which needed the unlocked wallet has finished. Any aquired heldUnlock *must* be released (preferably with a
	// defer) or the wallet will forever remain unlocked.
	heldUnlock chan struct{}
)

// walletLocker manages the locked/unlocked state of a wallet.
func (w *Wallet) walletLocker() {
	var timeout <-chan time.Time
	holdChan := make(heldUnlock)
	quit := w.quitChan()
	var err error
out:
	for {
		select {
		case req := <-w.unlockRequests:
			if err = walletdb.View(w.db, func(tx walletdb.ReadTx) (err error) {
				addrmgrNs := tx.ReadBucket(waddrmgrNamespaceKey)
				return w.Manager.Unlock(addrmgrNs, req.passphrase)
			}); slog.Check(err) {
				req.err <- err
				continue
			}
			timeout = req.lockAfter
			if timeout == nil {
				slog.Info("the wallet has been unlocked without a time limit")
			} else {
				slog.Info("the wallet has been temporarily unlocked")
			}
			req.err <- nil
			continue
		case req := <-w.changePassphrase:
			if err = walletdb.Update(w.db, func(tx walletdb.ReadWriteTx) (err error) {
				addrmgrNs := tx.ReadWriteBucket(waddrmgrNamespaceKey)
				return w.Manager.ChangePassphrase(
					addrmgrNs, req.old, req.new, req.private,
					&waddrmgr.DefaultScryptOptions,
				)
			}); slog.Check(err) {
			}
			req.err <- err
			continue
		case req := <-w.changePassphrases:
			if err = walletdb.Update(w.db, func(tx walletdb.ReadWriteTx) (err error) {
				addrmgrNs := tx.ReadWriteBucket(waddrmgrNamespaceKey)
				if err = w.Manager.ChangePassphrase(addrmgrNs, req.publicOld, req.publicNew, false,
					&waddrmgr.DefaultScryptOptions); slog.Check(err) {
					return
				}
				return w.Manager.ChangePassphrase(
					addrmgrNs, req.privateOld, req.privateNew,
					true, &waddrmgr.DefaultScryptOptions,
				)
			}); slog.Check(<-req.err) {
			}
			continue
		case req := <-w.holdUnlockRequests:
			if w.Manager.IsLocked() {
				close(req)
				continue
			}
			req <- holdChan
			<-holdChan // Block until the lock is released.
			// If, after holding onto the unlocked wallet for some time, the timeout has expired, lock it now instead of
			// hoping it gets unlocked next time the top level select runs.
			select {
			case <-timeout:
				// Let the top level select fallthrough so the wallet is locked.
			default:
				continue
			}
		case w.lockState <- w.Manager.IsLocked():
			continue
		case <-quit:
			break out
		case <-w.lockRequests:
		case <-timeout:
		}
		// Select statement fell through by an explicit lock or the timer expiring. Lock the manager here.
		timeout = nil
		if err = w.Manager.Lock(); slog.Check(err) && !waddrmgr.IsError(err, waddrmgr.ErrLocked) {
			slog.Error("could not lock wallet:", err)
		} else {
			slog.Info("the wallet has been locked")
		}
	}
	w.wg.Done()
}

// Unlock unlocks the wallet's address manager and relocks it after timeout has expired. If the wallet is already
// unlocked and the new passphrase is correct, the current timeout is replaced with the new one. The wallet will be
// locked if the passphrase is incorrect or any other error occurs during the unlock.
func (w *Wallet) Unlock(passphrase []byte, lock <-chan time.Time) (err error) {
	var cerr chan error
	cerr = make(chan error, 1)
	w.unlockRequests <- unlockRequest{
		passphrase: passphrase,
		lockAfter:  lock,
		err:        cerr,
	}
	return <-cerr
}

// Lock locks the wallet's address manager.
func (w *Wallet) Lock() {
	w.lockRequests <- struct{}{}
}

// Locked returns whether the account manager for a wallet is locked.
func (w *Wallet) Locked() bool {
	return <-w.lockState
}

// holdUnlock prevents the wallet from being locked. The heldUnlock object *must* be released, or the wallet will
// forever remain unlocked.
//
// TODO: To prevent the above scenario, perhaps closures should be passed to the walletLocker goroutine and disallow
//  callers from explicitly handling the locking mechanism.
func (w *Wallet) holdUnlock() (hl heldUnlock, err error) {
	req := make(chan heldUnlock)
	w.holdUnlockRequests <- req
	var ok bool
	hl, ok = <-req
	if !ok {
		// TODO(davec): This should be defined and exported from waddrmgr.
		err = waddrmgr.ManagerError{
			ErrorCode:   waddrmgr.ErrLocked,
			Description: "address manager is locked",
		}
	}
	return
}

// release releases the hold on the unlocked-state of the wallet and allows the wallet to be locked again. If a lock
// timeout has already expired, the wallet is locked again as soon as release is called.
func (c heldUnlock) release() {
	c <- struct{}{}
}

// ChangePrivatePassphrase attempts to change the passphrase for a wallet from old to new. Changing the passphrase is
// synchronized with all other address manager locking and unlocking. The lock state will be the same as it was before
// the password change.
func (w *Wallet) ChangePrivatePassphrase(old, new []byte) (err error) {
	er := make(chan error, 1)
	w.changePassphrase <- changePassphraseRequest{
		old:     old,
		new:     new,
		private: true,
		err:     er,
	}
	return <-er
}

// ChangePublicPassphrase modifies the public passphrase of the wallet.
func (w *Wallet) ChangePublicPassphrase(old, new []byte) (err error) {
	er := make(chan error, 1)
	w.changePassphrase <- changePassphraseRequest{
		old:     old,
		new:     new,
		private: false,
		err:     er,
	}
	return <-er
}

// ChangePassphrases modifies the public and private passphrase of the wallet atomically.
func (w *Wallet) ChangePassphrases(publicOld, publicNew, privateOld,
	privateNew []byte) (err error) {
	er := make(chan error, 1)
	w.changePassphrases <- changePassphrasesRequest{
		publicOld:  publicOld,
		publicNew:  publicNew,
		privateOld: privateOld,
		privateNew: privateNew,
		err:        er,
	}
	return <-er
}

// // accountUsed returns whether there are any recorded transactions spending to
// // a given account. It returns true if atleast one address in the account was
// // used and false if no address in the account was used.
// func (w *Wallet) accountUsed(addrmgrNs walletdb.ReadWriteBucket, account uint32) (bool, err error) {
// 	var used bool
// 	err := w.Manager.ForEachAccountAddress(addrmgrNs, account,
// 		func(maddr waddrmgr.ManagedAddress) (err error) {
// 			used = maddr.Used(addrmgrNs)
// 			if used {
// 				return waddrmgr.Break
// 			}
// 			return nil
// 		})
// 	if err == waddrmgr.Break {
// 		err = nil
// 	}
// 	return used, err
// }

// AccountAddresses returns the addresses for every created address for an account.
func (w *Wallet) AccountAddresses(account uint32) (addrs []util.Address, err error) {
	if err = walletdb.View(w.db, func(tx walletdb.ReadTx) (err error) {
		addrmgrNs := tx.ReadBucket(waddrmgrNamespaceKey)
		return w.Manager.ForEachAccountAddress(addrmgrNs, account, func(maddr waddrmgr.ManagedAddress) (err error) {
			addrs = append(addrs, maddr.Address())
			return
		})
	}); slog.Check(err) {
	}
	return
}

// CalculateBalance sums the amounts of all unspent transaction outputs to addresses of a wallet and returns the
// balance.
//
// If confirmations is 0, all UTXOs, even those not present in a block (height -1), will be used to get the balance.
// Otherwise, a UTXO must be in a block. If confirmations is 1 or greater, the balance will be calculated based on how
// many how many blocks include a UTXO.
func (w *Wallet) CalculateBalance(confirms int32) (balance util.Amount, err error) {
	if err = walletdb.View(w.db, func(tx walletdb.ReadTx) (err error) {
		txmgrNs := tx.ReadBucket(wtxmgrNamespaceKey)
		blk := w.Manager.SyncedTo()
		balance, err = w.TxStore.Balance(txmgrNs, confirms, blk.Height)
		return
	}); slog.Check(err) {
	}
	return
}

// Balances records total, spendable (by policy), and immature coinbase reward balance amounts.
type Balances struct {
	Total          util.Amount
	Spendable      util.Amount
	ImmatureReward util.Amount
}

// CalculateAccountBalances sums the amounts of all unspent transaction outputs to the given account of a wallet and
// returns the balance.
//
// This function is much slower than it needs to be since transactions outputs are not indexed by the accounts they
// credit to, and all unspent transaction outputs must be iterated.
func (w *Wallet) CalculateAccountBalances(account uint32, confirms int32) (bals Balances, err error) {
	if err = walletdb.View(w.db, func(tx walletdb.ReadTx) (err error) {
		addrmgrNs := tx.ReadBucket(waddrmgrNamespaceKey)
		txmgrNs := tx.ReadBucket(wtxmgrNamespaceKey)
		// Get current block.  The block height used for calculating
		// the number of tx confirmations.
		syncBlock := w.Manager.SyncedTo()
		var unspent []wtxmgr.Credit
		if unspent, err = w.TxStore.UnspentOutputs(txmgrNs); slog.Check(err) {
			return
		}
		for i := range unspent {
			output := &unspent[i]
			var outputAcct uint32
			var addrs []util.Address
			if _, addrs, _, err = txscript.ExtractPkScriptAddrs(output.PkScript, w.chainParams,
			); !slog.Check(err) && len(addrs) > 0 {
				_, outputAcct, err = w.Manager.AddrAccount(addrmgrNs, addrs[0])
			} else if outputAcct != account {
				continue
			}
			bals.Total += output.Amount
			if output.FromCoinBase && !confirmed(int32(w.chainParams.CoinbaseMaturity), output.Height, syncBlock.Height) {
				bals.ImmatureReward += output.Amount
			} else if confirmed(confirms, output.Height, syncBlock.Height) {
				bals.Spendable += output.Amount
			}
		}
		return
	}); slog.Check(err) {
	}
	return
}

// CurrentAddress gets the most recently requested Bitcoin payment address from a wallet for a particular key-chain
// scope. If the address has already been used (there is at least one transaction spending to it in the blockchain or
// pod mempool), the next chained address is returned.
func (w *Wallet) CurrentAddress(account uint32, scope waddrmgr.KeyScope) (addr util.Address, err error) {
	var chainClient chain.Interface
	if chainClient, err = w.requireChainClient(); slog.Check(err) {
		return
	}
	var manager *waddrmgr.ScopedKeyManager
	if manager, err = w.Manager.FetchScopedKeyManager(scope); slog.Check(err) {
		return
	}
	var (
		props *waddrmgr.AccountProperties
	)
	if err = walletdb.Update(w.db, func(tx walletdb.ReadWriteTx) (err error) {
		addrmgrNs := tx.ReadWriteBucket(waddrmgrNamespaceKey)
		var maddr waddrmgr.ManagedAddress
		if maddr, err = manager.LastExternalAddress(addrmgrNs, account); slog.Check(err) {
			// If no address exists yet, create the first external
			// address.
			if waddrmgr.IsError(err, waddrmgr.ErrAddressNotFound) {
				addr, props, err = w.newAddress(
					addrmgrNs, account, scope,
				)
			}
			return
		}
		// Get next chained address if the last one has already been used.
		if maddr.Used(addrmgrNs) {
			addr, props, err = w.newAddress(addrmgrNs, account, scope)
			return
		}
		addr = maddr.Address()
		return
	}); slog.Check(err) {
		return
	}
	// If the props have been initially, then we had to create a new address to satisfy the query. Notify the rpc server
	// about the new address.
	if props != nil {
		if err = chainClient.NotifyReceived([]util.Address{addr}); slog.Check(err) {
			return
		}
		w.NtfnServer.notifyAccountProperties(props)
	}
	return
}

// PubKeyForAddress looks up the associated public key for a P2PKH address.
func (w *Wallet) PubKeyForAddress(a util.Address) (pubKey *ec.PublicKey, err error) {
	if err = walletdb.View(w.db, func(tx walletdb.ReadTx) (err error) {
		addrmgrNs := tx.ReadBucket(waddrmgrNamespaceKey)
		var managedAddr waddrmgr.ManagedAddress
		if managedAddr, err = w.Manager.Address(addrmgrNs, a); slog.Check(err) {
			return
		}
		managedPubKeyAddr, ok := managedAddr.(waddrmgr.ManagedPubKeyAddress)
		if !ok {
			err = errors.New("address does not have an associated public key")
			slog.Debug(err)
			return
		}
		pubKey = managedPubKeyAddr.PubKey()
		return
	}); slog.Check(err) {
	}
	return
}

// PrivKeyForAddress looks up the associated private key for a P2PKH or P2PK address.
func (w *Wallet) PrivKeyForAddress(a util.Address) (privKey *ec.PrivateKey, err error) {
	if err = walletdb.View(w.db, func(tx walletdb.ReadTx) (err error) {
		addrmgrNs := tx.ReadBucket(waddrmgrNamespaceKey)
		var managedAddr waddrmgr.ManagedAddress
		if managedAddr, err = w.Manager.Address(addrmgrNs, a); slog.Check(err) {
			return
		}
		managedPubKeyAddr, ok := managedAddr.(waddrmgr.ManagedPubKeyAddress)
		if !ok {
			err = errors.New("address does not have an associated private key")
			slog.Debug(err)
			return
		}
		privKey, err = managedPubKeyAddr.PrivKey()
		return
	}); slog.Check(err) {
	}
	return
}

// HaveAddress returns whether the wallet is the owner of the address a.
func (w *Wallet) HaveAddress(a util.Address) (have bool, err error) {
	if err = walletdb.View(w.db, func(tx walletdb.ReadTx) (err error) {
		addrmgrNs := tx.ReadBucket(waddrmgrNamespaceKey)
		_, err = w.Manager.Address(addrmgrNs, a)
		return
	}); slog.Check(err) {
		have = true
		return
	}
	if waddrmgr.IsError(err, waddrmgr.ErrAddressNotFound) {
		return
	}
	return
}

// AccountOfAddress finds the account that an address is associated with.
func (w *Wallet) AccountOfAddress(a util.Address) (account uint32, err error) {
	if err = walletdb.View(w.db, func(tx walletdb.ReadTx) (err error) {
		addrmgrNs := tx.ReadBucket(waddrmgrNamespaceKey)
		_, account, err = w.Manager.AddrAccount(addrmgrNs, a)
		return err
	}); slog.Check(err) {
	}
	return
}

// AddressInfo returns detailed information regarding a wallet address.
func (w *Wallet) AddressInfo(a util.Address) (managedAddress waddrmgr.ManagedAddress, err error) {
	if err = walletdb.View(w.db, func(tx walletdb.ReadTx) (err error) {
		addrmgrNs := tx.ReadBucket(waddrmgrNamespaceKey)
		managedAddress, err = w.Manager.Address(addrmgrNs, a)
		return
	}); slog.Check(err) {
	}
	return
}

// AccountNumber returns the account number for an account name under a particular key scope.
func (w *Wallet) AccountNumber(scope waddrmgr.KeyScope, accountName string) (account uint32, err error) {
	var manager *waddrmgr.ScopedKeyManager
	if manager, err = w.Manager.FetchScopedKeyManager(scope); slog.Check(err) {
		return
	}
	if err = walletdb.View(w.db, func(tx walletdb.ReadTx) (err error) {
		addrmgrNs := tx.ReadBucket(waddrmgrNamespaceKey)
		account, err = manager.LookupAccount(addrmgrNs, accountName)
		return
	}); slog.Check(err) {
	}
	return
}

// AccountName returns the name of an account.
func (w *Wallet) AccountName(scope waddrmgr.KeyScope, accountNumber uint32) (accountName string, err error) {
	var manager *waddrmgr.ScopedKeyManager
	if manager, err = w.Manager.FetchScopedKeyManager(scope); slog.Check(err) {
		return
	}
	if err = walletdb.View(w.db, func(tx walletdb.ReadTx) (err error) {
		addrmgrNs := tx.ReadBucket(waddrmgrNamespaceKey)
		accountName, err = manager.AccountName(addrmgrNs, accountNumber)
		return
	}); slog.Check(err) {
	}
	return
}

// AccountProperties returns the properties of an account, including address indexes and name. It first fetches the
// desynced information from the address manager, then updates the indexes based on the address pools.
func (w *Wallet) AccountProperties(scope waddrmgr.KeyScope, acct uint32) (props *waddrmgr.AccountProperties, err error) {
	var manager *waddrmgr.ScopedKeyManager
	if manager, err = w.Manager.FetchScopedKeyManager(scope); slog.Check(err) {
		return
	}
	err = walletdb.View(w.db, func(tx walletdb.ReadTx) (err error) {
		waddrmgrNs := tx.ReadBucket(waddrmgrNamespaceKey)
		props, err = manager.AccountProperties(waddrmgrNs, acct)
		return
	})
	return
}

// RenameAccount sets the name for an account number to newName.
func (w *Wallet) RenameAccount(scope waddrmgr.KeyScope, account uint32, newName string) (err error) {
	var manager *waddrmgr.ScopedKeyManager
	if manager, err = w.Manager.FetchScopedKeyManager(scope); slog.Check(err) {
		return
	}
	var props *waddrmgr.AccountProperties
	err = walletdb.Update(w.db, func(tx walletdb.ReadWriteTx) (err error) {
		addrmgrNs := tx.ReadWriteBucket(waddrmgrNamespaceKey)
		if err = manager.RenameAccount(addrmgrNs, account, newName); slog.Check(err) {
			return
		}
		props, err = manager.AccountProperties(addrmgrNs, account)
		return
	})
	if err == nil {
		w.NtfnServer.notifyAccountProperties(props)
	}
	return
}

// const maxEmptyAccounts = 100

// NextAccount creates the next account and returns its account number. The name must be unique to the account. In order
// to support automatic seed restoring, new accounts may not be created when all of the previous 100 accounts have no
// transaction history (this is a deviation from the BIP0044 spec, which allows no unused account gaps).
func (w *Wallet) NextAccount(scope waddrmgr.KeyScope, name string) (account uint32, err error) {
	var manager *waddrmgr.ScopedKeyManager
	if manager, err = w.Manager.FetchScopedKeyManager(scope); slog.Check(err) {
		return
	}
	var props *waddrmgr.AccountProperties
	if err = walletdb.Update(w.db, func(tx walletdb.ReadWriteTx) (err error) {
		addrmgrNs := tx.ReadWriteBucket(waddrmgrNamespaceKey)
		if account, err = manager.NewAccount(addrmgrNs, name); slog.Check(err) {
			return
		}
		props, err = manager.AccountProperties(addrmgrNs, account)
		return
	}); slog.Check(err) {
		slog.Error("cannot fetch new account properties for notification after account creation:", err)
	}
	w.NtfnServer.notifyAccountProperties(props)
	return
}

// CreditCategory describes the type of wallet transaction output. The category of "sent transactions" (debits) is
// always "send", and is not expressed by this type.
//
// TODO: This is a requirement of the RPC server and should be moved.
type CreditCategory byte

// These constants define the possible credit categories.
const (
	CreditReceive CreditCategory = iota
	CreditGenerate
	CreditImmature
)

// String returns the category as a string. This string may be used as the JSON string for categories as part of
// listtransactions and gettransaction RPC responses.
func (c CreditCategory) String() string {
	switch c {
	case CreditReceive:
		return "receive"
	case CreditGenerate:
		return "generate"
	case CreditImmature:
		return "immature"
	default:
		return "unknown"
	}
}

// RecvCategory returns the category of received credit outputs from a transaction record. The passed block chain height
// is used to distinguish immature from mature coinbase outputs.
//
// TODO: This is intended for use by the RPC server and should be moved out of this package at a later time.
func RecvCategory(details *wtxmgr.TxDetails, syncHeight int32, net *netparams.Params) CreditCategory {
	if blockchain.IsCoinBaseTx(&details.MsgTx) {
		if confirmed(int32(net.CoinbaseMaturity), details.Block.Height,
			syncHeight) {
			return CreditGenerate
		}
		return CreditImmature
	}
	return CreditReceive
}

// listTransactions creates a object that may be marshalled to a response result for a listtransactions RPC.
//
// TODO: This should be moved to the legacyrpc package.
func listTransactions(tx walletdb.ReadTx, details *wtxmgr.TxDetails, addrMgr *waddrmgr.Manager,
	syncHeight int32, net *netparams.Params) []btcjson.ListTransactionsResult {
	addrmgrNs := tx.ReadBucket(waddrmgrNamespaceKey)
	var (
		blockHashStr  string
		blockTime     int64
		confirmations int64
	)
	if details.Block.Height != -1 {
		blockHashStr = details.Block.Hash.String()
		blockTime = details.Block.Time.Unix()
		confirmations = int64(confirms(details.Block.Height, syncHeight))
	}
	results := []btcjson.ListTransactionsResult{}
	txHashStr := details.Hash.String()
	received := details.Received.Unix()
	generated := blockchain.IsCoinBaseTx(&details.MsgTx)
	recvCat := RecvCategory(details, syncHeight, net).String()
	send := len(details.Debits) != 0
	// Fee can only be determined if every input is a debit.
	var feeF64 float64
	if len(details.Debits) == len(details.MsgTx.TxIn) {
		var debitTotal util.Amount
		for _, deb := range details.Debits {
			debitTotal += deb.Amount
		}
		var outputTotal util.Amount
		for _, output := range details.MsgTx.TxOut {
			outputTotal += util.Amount(output.Value)
		}
		// Note: The actual fee is debitTotal - outputTotal. However, this RPC reports negative numbers for fees, so the
		// inverse is calculated.
		feeF64 = (outputTotal - debitTotal).ToDUO()
	}
outputs:
	for i, output := range details.MsgTx.TxOut {
		// Determine if this output is a credit, and if so, determine its spentness.
		var isCredit bool
		var spentCredit bool
		for _, cred := range details.Credits {
			if cred.Index == uint32(i) {
				// Change outputs are ignored.
				if cred.Change {
					continue outputs
				}
				isCredit = true
				spentCredit = cred.Spent
				break
			}
		}
		var address string
		var accountName string
		_, addrs, _, _ := txscript.ExtractPkScriptAddrs(output.PkScript, net)
		if len(addrs) == 1 {
			addr := addrs[0]
			address = addr.EncodeAddress()
			mgr, account, err := addrMgr.AddrAccount(addrmgrNs, addrs[0])
			if err == nil {
				accountName, err = mgr.AccountName(addrmgrNs, account)
				if err != nil {
					slog.Error(err)
					accountName = ""
				}
			}
		}
		amountF64 := util.Amount(output.Value).ToDUO()
		result := btcjson.ListTransactionsResult{
			// Fields left zeroed:
			//   InvolvesWatchOnly
			//   BlockIndex
			//
			// Fields set below:
			//   Account (only for non-"send" categories)
			//   Category
			//   Amount
			//   Fee
			Address:         address,
			Vout:            uint32(i),
			Confirmations:   confirmations,
			Generated:       generated,
			BlockHash:       blockHashStr,
			BlockTime:       blockTime,
			TxID:            txHashStr,
			WalletConflicts: []string{},
			Time:            received,
			TimeReceived:    received,
		}
		// Add a received/generated/immature result if this is a credit. If the output was spent, create a second result
		// under the send category with the inverse of the output amount. It is therefore possible that a single output
		// may be included in the results set zero, one, or two times.
		//
		// Since credits are not saved for outputs that are not controlled by this wallet, all non-credits from
		// transactions with debits are grouped under the send category.
		if send || spentCredit {
			result.Category = "send"
			result.Amount = -amountF64
			result.Fee = &feeF64
			results = append(results, result)
		}
		if isCredit {
			result.Account = accountName
			result.Category = recvCat
			result.Amount = amountF64
			result.Fee = nil
			results = append(results, result)
		}
	}
	return results
}

// ListSinceBlock returns a slice of objects with details about transactions since the given block. If the block is -1
// then all transactions are included. This is intended to be used for listsinceblock RPC replies.
func (w *Wallet) ListSinceBlock(start, end, syncHeight int32) (txList []btcjson.ListTransactionsResult, err error) {
	txList = []btcjson.ListTransactionsResult{}
	if err = walletdb.View(w.db, func(tx walletdb.ReadTx) (err error) {
		txmgrNs := tx.ReadBucket(wtxmgrNamespaceKey)
		rangeFn := func(details []wtxmgr.TxDetails) (b bool, err error) {
			for _, detail := range details {
				jsonResults := listTransactions(tx, &detail, w.Manager, syncHeight, w.chainParams)
				txList = append(txList, jsonResults...)
			}
			return
		}
		return w.TxStore.RangeTransactions(txmgrNs, start, end, rangeFn)
	}); slog.Check(err) {
	}
	return
}

// ListTransactions returns a slice of objects with details about a recorded transaction. This is intended to be used
// for listtransactions RPC replies.
func (w *Wallet) ListTransactions(from, count int) (txList []btcjson.ListTransactionsResult, err error) {
	// txList := []btcjson.ListTransactionsResult{}
	slog.Trace("ListTransactions", from, count)
	if err = walletdb.View(w.db, func(tx walletdb.ReadTx) (err error) {
		txmgrNs := tx.ReadBucket(wtxmgrNamespaceKey)
		// Get current block. The block height used for calculating the number of tx confirmations.
		syncBlock := w.Manager.SyncedTo()
		slog.Trace("synced to", syncBlock)
		// Need to skip the first from transactions, and after those, only include the next count transactions.
		skipped := 0
		n := 0
		rangeFn := func(details []wtxmgr.TxDetails) (b bool, err error) {
			// Iterate over transactions at this height in reverse order. This does nothing for unmined transactions,
			// which are unsorted, but it will process mined transactions in the reverse order they were marked mined.
			for i := len(details) - 1; i >= 0; i-- {
				if from > skipped {
					skipped++
					continue
				}
				n++
				if n > count {
					return true, nil
				}
				jsonResults := listTransactions(tx, &details[i], w.Manager, syncBlock.Height, w.chainParams)
				txList = append(txList, jsonResults...)
				if len(jsonResults) > 0 {
					n++
				}
			}
			return
		}
		// Return newer results first by starting at mempool height and working
		// down to the genesis block.
		return w.TxStore.RangeTransactions(txmgrNs, -1, 0, rangeFn)
	}); slog.Check(err) {
	}
	return
}

// ListAddressTransactions returns a slice of objects with details about
// recorded transactions to or from any address belonging to a set.  This is
// intended to be used for listaddresstransactions RPC replies.
func (w *Wallet) ListAddressTransactions(pkHashes map[string]struct{}) (txList []btcjson.ListTransactionsResult, err error) {
	err = walletdb.View(w.db, func(tx walletdb.ReadTx) (err error) {
		txmgrNs := tx.ReadBucket(wtxmgrNamespaceKey)
		// Get current block.  The block height used for calculating
		// the number of tx confirmations.
		syncBlock := w.Manager.SyncedTo()
		rangeFn := func(details []wtxmgr.TxDetails) (b bool, err error) {
		loopDetails:
			for i := range details {
				detail := &details[i]
				for _, cred := range detail.Credits {
					pkScript := detail.MsgTx.TxOut[cred.Index].PkScript
					var addrs []util.Address
					if _, addrs, _, err = txscript.ExtractPkScriptAddrs(pkScript, w.chainParams); slog.Check(err) ||
						len(addrs) != 1 {
						continue
					}
					apkh, ok := addrs[0].(*util.AddressPubKeyHash)
					if !ok {
						continue
					}
					_, ok = pkHashes[string(apkh.ScriptAddress())]
					if !ok {
						continue
					}
					jsonResults := listTransactions(tx, detail, w.Manager, syncBlock.Height, w.chainParams)
					// if err != nil {
					// 	return false, err
					// }
					txList = append(txList, jsonResults...)
					continue loopDetails
				}
			}
			return false, nil
		}
		return w.TxStore.RangeTransactions(txmgrNs, 0, -1, rangeFn)
	})
	return
}

// ListAllTransactions returns a slice of objects with details about a recorded
// transaction.  This is intended to be used for listalltransactions RPC
// replies.
func (w *Wallet) ListAllTransactions() (txList []btcjson.ListTransactionsResult, err error) {
	err = walletdb.View(w.db, func(tx walletdb.ReadTx) (err error) {
		txmgrNs := tx.ReadBucket(wtxmgrNamespaceKey)
		// Get current block.  The block height used for calculating
		// the number of tx confirmations.
		syncBlock := w.Manager.SyncedTo()
		rangeFn := func(details []wtxmgr.TxDetails) (b bool, err error) {
			// Iterate over transactions at this height in reverse order. This does nothing for unmined transactions,
			// which are unsorted, but it will process mined transactions in the reverse order they were marked mined.
			for i := len(details) - 1; i >= 0; i-- {
				jsonResults := listTransactions(tx, &details[i], w.Manager,
					syncBlock.Height, w.chainParams)
				txList = append(txList, jsonResults...)
			}
			return false, nil
		}
		// Return newer results first by starting at mempool height and working down to the genesis block.
		return w.TxStore.RangeTransactions(txmgrNs, -1, 0, rangeFn)
	})
	return
}

// BlockIdentifier identifies a block by either a height or a hash.
type BlockIdentifier struct {
	height int32
	hash   *chainhash.Hash
}

// NewBlockIdentifierFromHeight constructs a BlockIdentifier for a block height.
func NewBlockIdentifierFromHeight(height int32) *BlockIdentifier {
	return &BlockIdentifier{height: height}
}

// NewBlockIdentifierFromHash constructs a BlockIdentifier for a block hash.
func NewBlockIdentifierFromHash(hash *chainhash.Hash) *BlockIdentifier {
	return &BlockIdentifier{hash: hash}
}

// GetTransactionsResult is the result of the wallet's GetTransactions method.
// See GetTransactions for more details.
type GetTransactionsResult struct {
	MinedTransactions   []Block
	UnminedTransactions []TransactionSummary
}

// GetTransactions returns transaction results between a starting and ending block. BlockC in the block range may be
// specified by either a height or a hash.
//
// Because this is a possibly lenghtly operation, a cancel channel is provided to cancel the task. If this channel
// unblocks, the results created thus far will be returned.
//
// Transaction results are organized by blocks in ascending order and unmined transactions in an unspecified order.
// Mined transactions are saved in a Block structure which records properties about the block.
func (w *Wallet) GetTransactions(startBlock, endBlock *BlockIdentifier, cancel <-chan struct{}) (
	res *GetTransactionsResult, err error) {
	var start, end int32 = 0, -1
	w.chainClientLock.Lock()
	chainClient := w.chainClient
	w.chainClientLock.Unlock()
	// TODO: Fetching block heights by their hashes is inherently racy because not all block headers are saved but when
	//  they are for SPV the db can be queried directly without this.
	var startResp, endResp rpcclient.FutureGetBlockVerboseResult
	if startBlock != nil {
		if startBlock.hash == nil {
			start = startBlock.height
		} else {
			if chainClient == nil {
				err = errors.New("no chain server client")
				slog.Debug(err)
				return
			}
			switch client := chainClient.(type) {
			case *chain.RPCClient:
				startResp = client.GetBlockVerboseTxAsync(startBlock.hash)
			case *chain.BitcoindClient:
				if start, err = client.GetBlockHeight(startBlock.hash); slog.Check(err) {
					return
				}
			case *chain.NeutrinoClient:
				if start, err = client.GetBlockHeight(startBlock.hash); slog.Check(err) {
					return
				}
			}
		}
	}
	if endBlock != nil {
		if endBlock.hash == nil {
			end = endBlock.height
		} else {
			if chainClient == nil {
				err = errors.New("no chain server client")
				slog.Debug(err)
				return
			}
			switch client := chainClient.(type) {
			case *chain.RPCClient:
				endResp = client.GetBlockVerboseTxAsync(endBlock.hash)
			case *chain.NeutrinoClient:
				end, err = client.GetBlockHeight(endBlock.hash)
				if err != nil {
					slog.Error(err)
					return nil, err
				}
			}
		}
	}
	if startResp != nil {
		resp, err := startResp.Receive()
		if err != nil {
			slog.Error(err)
			return nil, err
		}
		start = int32(resp.Height)
	}
	if endResp != nil {
		resp, err := endResp.Receive()
		if err != nil {
			slog.Error(err)
			return nil, err
		}
		end = int32(resp.Height)
	}
	if err = walletdb.View(w.db, func(dbtx walletdb.ReadTx) (err error) {
		txmgrNs := dbtx.ReadBucket(wtxmgrNamespaceKey)
		rangeFn := func(details []wtxmgr.TxDetails) (b bool, err error) {
			// TODO: probably should make RangeTransactions not reuse the
			// details backing array memory.
			dets := make([]wtxmgr.TxDetails, len(details))
			copy(dets, details)
			details = dets
			txs := make([]TransactionSummary, 0, len(details))
			for i := range details {
				txs = append(txs, makeTxSummary(dbtx, w, &details[i]))
			}
			if details[0].Block.Height != -1 {
				blockHash := details[0].Block.Hash
				res.MinedTransactions = append(res.MinedTransactions, Block{
					Hash:         &blockHash,
					Height:       details[0].Block.Height,
					Timestamp:    details[0].Block.Time.Unix(),
					Transactions: txs,
				})
			} else {
				res.UnminedTransactions = txs
			}
			select {
			case <-cancel:
				return true, nil
			default:
				return false, nil
			}
		}
		return w.TxStore.RangeTransactions(txmgrNs, start, end, rangeFn)
	}); slog.Check(err) {
	}
	return
}

// AccountResult is a single account result for the AccountsResult type.
type AccountResult struct {
	waddrmgr.AccountProperties
	TotalBalance util.Amount
}

// AccountsResult is the resutl of the wallet's Accounts method.  See that
// method for more details.
type AccountsResult struct {
	Accounts           []AccountResult
	CurrentBlockHash   *chainhash.Hash
	CurrentBlockHeight int32
}

// Accounts returns the current names, numbers, and total balances of all
// accounts in the wallet restricted to a particular key scope.  The current
// chain tip is included in the result for atomicity reasons.
//
// TODO(jrick): Is the chain tip really needed, since only the total balances
// are included?
func (w *Wallet) Accounts(scope waddrmgr.KeyScope) (result *AccountsResult, err error) {
	var manager *waddrmgr.ScopedKeyManager
	if manager, err = w.Manager.FetchScopedKeyManager(scope); slog.Check(err) {
		return
	}
	var (
		accounts        []AccountResult
		syncBlockHash   *chainhash.Hash
		syncBlockHeight int32
	)
	if err = walletdb.View(w.db, func(tx walletdb.ReadTx) (err error) {
		addrmgrNs := tx.ReadBucket(waddrmgrNamespaceKey)
		txmgrNs := tx.ReadBucket(wtxmgrNamespaceKey)
		syncBlock := w.Manager.SyncedTo()
		syncBlockHash = &syncBlock.Hash
		syncBlockHeight = syncBlock.Height
		var unspent []wtxmgr.Credit
		if unspent, err = w.TxStore.UnspentOutputs(txmgrNs); slog.Check(err) {
			return
		}
		if err = manager.ForEachAccount(addrmgrNs, func(acct uint32) (err error) {
			var props *waddrmgr.AccountProperties
			if props, err = manager.AccountProperties(addrmgrNs, acct); slog.Check(err) {
				return
			}
			accounts = append(accounts, AccountResult{
				AccountProperties: *props,
				// TotalBalance set below
			})
			return
		}); slog.Check(err) {
			return
		}
		m := make(map[uint32]*util.Amount)
		for i := range accounts {
			a := &accounts[i]
			m[a.AccountNumber] = &a.TotalBalance
		}
		var addrs []util.Address
		var outputAcct uint32
		for i := range unspent {
			output := unspent[i]
			if _, addrs, _, err = txscript.ExtractPkScriptAddrs(output.PkScript, w.chainParams); !slog.Check(err) && len(addrs) > 0 {
				_, outputAcct, err = w.Manager.AddrAccount(addrmgrNs, addrs[0])
			}
			if err == nil {
				amt, ok := m[outputAcct]
				if ok {
					*amt += output.Amount
				}
			}
		}
		return
	}); slog.Check(err) {
	}
	return &AccountsResult{
			Accounts:           accounts,
			CurrentBlockHash:   syncBlockHash,
			CurrentBlockHeight: syncBlockHeight,
		},
		err
}

// AccountBalanceResult is a single result for the Wallet.AccountBalances method.
type AccountBalanceResult struct {
	AccountNumber  uint32
	AccountName    string
	AccountBalance util.Amount
}

// AccountBalances returns all accounts in the wallet and their balances. Balances are determined by excluding
// transactions that have not met requiredConfs confirmations.
func (w *Wallet) AccountBalances(scope waddrmgr.KeyScope, requiredConfs int32) (results []AccountBalanceResult, err error) {
	manager, err := w.Manager.FetchScopedKeyManager(scope)
	if err != nil {
		slog.Error(err)
		return nil, err
	}
	if err = walletdb.View(w.db, func(tx walletdb.ReadTx) (err error) {
		addrmgrNs := tx.ReadBucket(waddrmgrNamespaceKey)
		txmgrNs := tx.ReadBucket(wtxmgrNamespaceKey)
		syncBlock := w.Manager.SyncedTo()
		// Fill out all account info except for the balances.
		lastAcct, err := manager.LastAccount(addrmgrNs)
		if err != nil {
			slog.Error(err)
			return err
		}
		results = make([]AccountBalanceResult, lastAcct+2)
		for i := range results[:len(results)-1] {
			accountName, err := manager.AccountName(addrmgrNs, uint32(i))
			if err != nil {
				slog.Error(err)
				return err
			}
			results[i].AccountNumber = uint32(i)
			results[i].AccountName = accountName
		}
		results[len(results)-1].AccountNumber = waddrmgr.ImportedAddrAccount
		results[len(results)-1].AccountName = waddrmgr.ImportedAddrAccountName
		// Fetch all unspent outputs, and iterate over them tallying each account's balance where the output script pays
		// to an account address and the required number of confirmations is met.
		var unspentOutputs []wtxmgr.Credit
		if unspentOutputs, err = w.TxStore.UnspentOutputs(txmgrNs); slog.Check(err) {
			return
		}
		var outputAcct uint32
		var addrs []util.Address
		for i := range unspentOutputs {
			output := &unspentOutputs[i]
			if !confirmed(requiredConfs, output.Height, syncBlock.Height) {
				continue
			}
			if output.FromCoinBase && !confirmed(int32(w.ChainParams().CoinbaseMaturity),
				output.Height, syncBlock.Height) {
				continue
			}
			if _, addrs, _, err = txscript.ExtractPkScriptAddrs(output.PkScript, w.chainParams); slog.Check(err) || len(addrs) == 0 {
				continue
			}
			if outputAcct, err = manager.AddrAccount(addrmgrNs, addrs[0]); slog.Check(err) {
				continue
			}
			switch {
			case outputAcct == waddrmgr.ImportedAddrAccount:
				results[len(results)-1].AccountBalance += output.Amount
			case outputAcct > lastAcct:
				return errors.New("waddrmgr.Manager.AddrAccount returned account beyond recorded last account")
			default:
				results[outputAcct].AccountBalance += output.Amount
			}
		}
		return
	}); slog.Check(err) {
	}
	return
}

// creditSlice satisifies the sort.Interface interface to provide sorting transaction credits from oldest to newest.
// Credits with the same receive time and mined in the same block are not guaranteed to be sorted by the order they
// appear in the block. Credits from the same transaction are sorted by output index.
type creditSlice []wtxmgr.Credit

func (s creditSlice) Len() int {
	return len(s)
}
func (s creditSlice) Less(i, j int) bool {
	switch {
	// If both credits are from the same tx, sort by output index.
	case s[i].OutPoint.Hash == s[j].OutPoint.Hash:
		return s[i].OutPoint.Index < s[j].OutPoint.Index
	// If both transactions are unmined, sort by their received date.
	case s[i].Height == -1 && s[j].Height == -1:
		return s[i].Received.Before(s[j].Received)
	// Unmined (newer) txs always come last.
	case s[i].Height == -1:
		return false
	case s[j].Height == -1:
		return true
	// If both txs are mined in different blocks, sort by block height.
	default:
		return s[i].Height < s[j].Height
	}
}
func (s creditSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// ListUnspent returns a slice of objects representing the unspent wallet transactions fitting the given criteria. The
// confirmations will be more than minconf, less than maxconf and if addresses is populated only the addresses contained
// within it will be considered. If we know nothing about a transaction an empty array will be returned.
func (w *Wallet) ListUnspent(minconf, maxconf int32, addresses map[string]struct{}) (results []*btcjson.ListUnspentResult, err error) {
	if err = walletdb.View(w.db, func(tx walletdb.ReadTx) (err error) {
		addrmgrNs := tx.ReadBucket(waddrmgrNamespaceKey)
		txmgrNs := tx.ReadBucket(wtxmgrNamespaceKey)
		syncBlock := w.Manager.SyncedTo()
		filter := len(addresses) != 0
		var unspent []wtxmgr.Credit
		if unspent, err = w.TxStore.UnspentOutputs(txmgrNs); slog.Check(err) {
			return
		}
		sort.Sort(sort.Reverse(creditSlice(unspent)))
		defaultAccountName := "default"
		results = make([]*btcjson.ListUnspentResult, 0, len(unspent))
		var sc txscript.ScriptClass
		var addrs []util.Address
		for i := range unspent {
			output := unspent[i]
			// Outputs with fewer confirmations than the minimum or more confs than the maximum are excluded.
			confs := confirms(output.Height, syncBlock.Height)
			if confs < minconf || confs > maxconf {
				continue
			}
			// Only mature coinbase outputs are included.
			if output.FromCoinBase {
				target := int32(w.ChainParams().CoinbaseMaturity)
				if !confirmed(target, output.Height, syncBlock.Height) {
					continue
				}
			}
			// Exclude locked outputs from the result set.
			if w.LockedOutpoint(output.OutPoint) {
				continue
			}
			// Lookup the associated account for the output. Use the default account name in case there is no associated
			// account for some reason, although this should never happen.
			//
			// This will be unnecessary once transactions and outputs are grouped under the associated account in the
			// db.
			acctName := defaultAccountName
			if sc, addrs, _, err = txscript.ExtractPkScriptAddrs(output.PkScript, w.chainParams); slog.Check(err) {
				continue
			}
			if len(addrs) > 0 {
				var smgr *waddrmgr.ScopedKeyManager
				var acct uint32
				if smgr, acct, err = w.Manager.AddrAccount(addrmgrNs, addrs[0]); !slog.Check(err) {
					s, err := smgr.AccountName(addrmgrNs, acct)
					if err == nil {
						acctName = s
					}
				}
			}
			if filter {
				for _, addr := range addrs {
					_, ok := addresses[addr.EncodeAddress()]
					if ok {
						goto include
					}
				}
				continue
			}
		include:
			// At the moment watch-only addresses are not supported, so all recorded outputs that are not multisig are
			// "spendable". Multisig outputs are only "spendable" if all keys are controlled by this wallet.
			//
			// TODO: Each case will need updates when watch-only addrs is added. For P2PK, P2PKH, and P2SH, the address
			//  must be looked up and not be watching-only. For multisig, all pubkeys must belong to the manager with the
			//  associated private key (currently it only checks whether the pubkey exists, since the private key is
			//  required at the moment).
			var spendable bool
		scSwitch:
			switch sc {
			case txscript.PubKeyHashTy:
				spendable = true
			case txscript.PubKeyTy:
				spendable = true
			case txscript.WitnessV0ScriptHashTy:
				spendable = true
			case txscript.WitnessV0PubKeyHashTy:
				spendable = true
			case txscript.MultiSigTy:
				for _, a := range addrs {
					if _, err = w.Manager.Address(addrmgrNs, a); !slog.Check(err) {
						continue
					}
					if waddrmgr.IsError(err, waddrmgr.ErrAddressNotFound) {
						break scSwitch
					}
					return
				}
				spendable = true
			}
			result := &btcjson.ListUnspentResult{
				TxID:          output.OutPoint.Hash.String(),
				Vout:          output.OutPoint.Index,
				Account:       acctName,
				ScriptPubKey:  hex.EncodeToString(output.PkScript),
				Amount:        output.Amount.ToDUO(),
				Confirmations: int64(confs),
				Spendable:     spendable,
			}
			// BUG: this should be a JSON array so that all
			// addresses can be included, or removed (and the
			// caller extracts addresses from the pkScript).
			if len(addrs) > 0 {
				result.Address = addrs[0].EncodeAddress()
			}
			results = append(results, result)
		}
		return
	}); slog.Check(err) {
	}
	return
}

// DumpPrivKeys returns the WIF-encoded private keys for all addresses with private keys in a wallet.
func (w *Wallet) DumpPrivKeys() (privkeys []string, err error) {
	if err = walletdb.View(w.db, func(tx walletdb.ReadTx) (err error) {
		addrmgrNs := tx.ReadBucket(waddrmgrNamespaceKey)
		// Iterate over each active address, appending the private key to
		// privkeys.
		return w.Manager.ForEachActiveAddress(addrmgrNs, func(addr util.Address) (err error) {
			ma, err := w.Manager.Address(addrmgrNs, addr)
			if err != nil {
				slog.Error(err)
				return err
			}
			// Only those addresses with keys needed.
			pka, ok := ma.(waddrmgr.ManagedPubKeyAddress)
			if !ok {
				return
			}
			var wif *util.WIF
			if wif, err = pka.ExportPrivKey(); slog.Check(err) {
				// It would be nice to zero out the array here. However, since strings in go are immutable, and we have
				// no control over the caller I don't think we can. :(
				return
			}
			privkeys = append(privkeys, wif.String())
			return
		})
	}); slog.Check(err) {
	}
	return
}

// DumpWIFPrivateKey returns the WIF encoded private key for a
// single wallet address.
func (w *Wallet) DumpWIFPrivateKey(addr util.Address) (a string, err error) {
	var maddr waddrmgr.ManagedAddress
	if err = walletdb.View(w.db, func(tx walletdb.ReadTx) (err error) {
		waddrmgrNs := tx.ReadBucket(waddrmgrNamespaceKey)
		// Get private key from wallet if it exists.
		maddr, err = w.Manager.Address(waddrmgrNs, addr)
		return
	}); slog.Check(err) {
		return
	}
	pka, ok := maddr.(waddrmgr.ManagedPubKeyAddress)
	if !ok {
		err = fmt.Errorf("address %s is not a key type", addr)
		slog.Debug(err)
		return
	}
	var wif *util.WIF
	if wif, err = pka.ExportPrivKey(); slog.Check(err) {
		return
	}
	return wif.String(), nil
}

// ImportPrivateKey imports a private key to the wallet and writes the new wallet to disk.
func (w *Wallet) ImportPrivateKey(scope waddrmgr.KeyScope, wif *util.WIF,
	bs *waddrmgr.BlockStamp, rescan bool) (addrStr string, err error) {
	var manager *waddrmgr.ScopedKeyManager
	if manager, err = w.Manager.FetchScopedKeyManager(scope); slog.Check(err) {
		return
	}
	// The starting block for the key is the genesis block unless otherwise specified.
	var newBirthday time.Time
	if bs == nil {
		bs = &waddrmgr.BlockStamp{
			Hash:   *w.chainParams.GenesisHash,
			Height: 0,
		}
	} else {
		// Only update the new birthday time from default value if we actually have timestamp info in the header.
		header, err := w.chainClient.GetBlockHeader(&bs.Hash)
		if err == nil {
			newBirthday = header.Timestamp
		}
	}
	// Attempt to import private key into wallet.
	var addr util.Address
	var props *waddrmgr.AccountProperties
	if err = walletdb.Update(w.db, func(tx walletdb.ReadWriteTx) (err error) {
		addrmgrNs := tx.ReadWriteBucket(waddrmgrNamespaceKey)
		var maddr waddrmgr.ManagedPubKeyAddress
		if maddr, err = manager.ImportPrivateKey(addrmgrNs, wif, bs); slog.Check(err) {
			return
		}
		addr = maddr.Address()
		if props, err = manager.AccountProperties(addrmgrNs, waddrmgr.ImportedAddrAccount); slog.Check(err) {
			return
		}
		return w.Manager.SetBirthday(addrmgrNs, newBirthday)
	}); slog.Check(err) {
		return
	}
	// Rescan blockchain for transactions with txout scripts paying to the imported address.
	if rescan {
		job := &RescanJob{
			Addrs:      []util.Address{addr},
			OutPoints:  nil,
			BlockStamp: *bs,
		}
		// Submit rescan job and log when the import has completed. Do not block on finishing the rescan. The rescan
		// success or failure is logged elsewhere, and the channel is not required to be read, so discard the return
		// value.
		_ = w.SubmitRescan(job)
	} else {
		if err = w.chainClient.NotifyReceived([]util.Address{addr}); slog.Check(err) {
			err = fmt.Errorf("Failed to subscribe for address ntfns for address %s: %s",
				addr.EncodeAddress(), err)
			return
		}
	}
	addrStr = addr.EncodeAddress()
	slog.Info("imported payment address", addrStr)
	w.NtfnServer.notifyAccountProperties(props)
	// Return the payment address string of the imported private key.
	return
}

// LockedOutpoint returns whether an outpoint has been marked as locked and should not be used as an input for created
// transactions.
func (w *Wallet) LockedOutpoint(op wire.OutPoint) bool {
	_, locked := w.lockedOutpoints[op]
	return locked
}

// LockOutpoint marks an outpoint as locked, that is, it should not be used as an input for newly created transactions.
func (w *Wallet) LockOutpoint(op wire.OutPoint) {
	w.lockedOutpoints[op] = struct{}{}
}

// UnlockOutpoint marks an outpoint as unlocked, that is, it may be used as an input for newly created transactions.
func (w *Wallet) UnlockOutpoint(op wire.OutPoint) {
	delete(w.lockedOutpoints, op)
}

// ResetLockedOutpoints resets the set of locked outpoints so all may be used as inputs for new transactions.
func (w *Wallet) ResetLockedOutpoints() {
	w.lockedOutpoints = map[wire.OutPoint]struct{}{}
}

// LockedOutpoints returns a slice of currently locked outpoints. This is intended to be used by marshaling the result
// as a JSON array for listlockunspent RPC results.
func (w *Wallet) LockedOutpoints() []btcjson.TransactionInput {
	locked := make([]btcjson.TransactionInput, len(w.lockedOutpoints))
	i := 0
	for op := range w.lockedOutpoints {
		locked[i] = btcjson.TransactionInput{
			Txid: op.Hash.String(),
			Vout: op.Index,
		}
		i++
	}
	return locked
}

// resendUnminedTxs iterates through all transactions that spend from wallet credits that are not known to have been
// mined into a block, and attempts to send each to the chain server for relay.
func (w *Wallet) resendUnminedTxs() {
	var chainClient chain.Interface
	var err error
	if chainClient, err = w.requireChainClient(); slog.Check(err) {
		slog.Error("no chain server available to resend unmined transactions", err)
		return
	}
	var txs []*wire.MsgTx
	if err = walletdb.View(w.db, func(tx walletdb.ReadTx) (err error) {
		txmgrNs := tx.ReadBucket(wtxmgrNamespaceKey)
		if txs, err = w.TxStore.UnminedTxs(txmgrNs); slog.Check(err) {
		}
		return
	}); slog.Check(err) {
		slog.Error("cannot load unmined transactions for resending:", err)
		return
	}
	var resp *chainhash.Hash
	for _, tx := range txs {
		if resp, err = chainClient.SendRawTransaction(tx, false); slog.Check(err) {
			slog.Debugf("could not resend transaction %v: %v %s", tx.TxHash(), err)
			// We'll only stop broadcasting transactions if we detect that the output has already been fully spent, is
			// an orphan, or is conflicting with another transaction.
			//
			// TODO(roasbeef): SendRawTransaction needs to return concrete error types, no need for string matching
			switch {
			// The following are errors returned from pod's mempool.
			case strings.Contains(err.Error(), "spent"):
			case strings.Contains(err.Error(), "orphan"):
			case strings.Contains(err.Error(), "conflict"):
			case strings.Contains(err.Error(), "already exists"):
			case strings.Contains(err.Error(), "negative"):
			// The following errors are returned from bitcoind's mempool.
			case strings.Contains(err.Error(), "Missing inputs"):
			case strings.Contains(err.Error(), "already in block chain"):
			case strings.Contains(err.Error(), "fee not met"):
			default:
				continue
			}
			// As the transaction was rejected, we'll attempt to remove the unmined transaction all together. Otherwise,
			// we'll keep attempting to rebroadcast this, and we may be computing our balance incorrectly if this tx
			// credits or debits to us.
			tt := tx
			if err = walletdb.Update(w.db, func(dbTx walletdb.ReadWriteTx) (err error) {
				txmgrNs := dbTx.ReadWriteBucket(wtxmgrNamespaceKey)
				var txRec *wtxmgr.TxRecord
				if txRec, err = wtxmgr.NewTxRecordFromMsgTx(
					tt, time.Now(),
				); slog.Check(err) {
					return
				}
				return w.TxStore.RemoveUnminedTx(txmgrNs, txRec)
			}); slog.Check(err) {
				slog.Warnf("unable to remove conflicting tx %v: %v %s", tt.TxHash(), err)
				continue
			}
			slog.Info(func() string {
				return "removed conflicting tx:" + spew.Sdump(tt) + " "
			}())
			continue
		}
		slog.Debug("resent unmined transaction", resp)
	}
}

// SortedActivePaymentAddresses returns a slice of all active payment addresses in a wallet.
func (w *Wallet) SortedActivePaymentAddresses() (addrStrs []string, err error) {
	if err = walletdb.View(w.db, func(tx walletdb.ReadTx) (err error) {
		addrmgrNs := tx.ReadBucket(waddrmgrNamespaceKey)
		return w.Manager.ForEachActiveAddress(addrmgrNs, func(addr util.Address) (err error) {
			addrStrs = append(addrStrs, addr.EncodeAddress())
			return
		})
	}); slog.Check(err) {
		return
	}
	sort.Strings(addrStrs)
	return
}

// NewAddress returns the next external chained address for a wallet.
func (w *Wallet) NewAddress(account uint32, scope waddrmgr.KeyScope, nochain bool) (addr util.Address, err error) {
	var (
		chainClient chain.Interface
		props       *waddrmgr.AccountProperties
	)
	if !nochain {
		if chainClient, err = w.requireChainClient(); slog.Check(err) {
			return
		}
	}
	if err = walletdb.Update(w.db, func(tx walletdb.ReadWriteTx) (err error) {
		addrmgrNs := tx.ReadWriteBucket(waddrmgrNamespaceKey)
		addr, props, err = w.newAddress(addrmgrNs, account, scope)
		return
	}); slog.Check(err) {
		return
	}
	if !nochain {
		// Notify the rpc server about the newly created address.
		if err = chainClient.NotifyReceived([]util.Address{addr}); slog.Check(err) {
			return
		}
		w.NtfnServer.notifyAccountProperties(props)
	}
	return
}
func (w *Wallet) newAddress(addrmgrNs walletdb.ReadWriteBucket, account uint32,
	scope waddrmgr.KeyScope) (a util.Address, props *waddrmgr.AccountProperties, err error) {
	var manager *waddrmgr.ScopedKeyManager
	if manager, err = w.Manager.FetchScopedKeyManager(scope); slog.Check(err) {
		return
	}
	// Get next address from wallet.
	var addrs []waddrmgr.ManagedAddress
	if addrs, err = manager.NextExternalAddresses(addrmgrNs, account, 1); slog.Check(err) {
		return
	}
	if props, err = manager.AccountProperties(addrmgrNs, account); slog.Check(err) {
		slog.Error("cannot fetch account properties for notification after deriving next external address:", err)
		return
	}
	return addrs[0].Address(), props, nil
}

// NewChangeAddress returns a new change address for a wallet.
func (w *Wallet) NewChangeAddress(account uint32, scope waddrmgr.KeyScope) (addr util.Address, err error) {
	var chainClient chain.Interface
	if chainClient, err = w.requireChainClient(); slog.Check(err) {
		return
	}
	if err = walletdb.Update(w.db, func(tx walletdb.ReadWriteTx) (err error) {
		addrmgrNs := tx.ReadWriteBucket(waddrmgrNamespaceKey)
		addr, err = w.newChangeAddress(addrmgrNs, account)
		return
	}); slog.Check(err) {
		return
	}
	// Notify the rpc server about the newly created address.
	if err = chainClient.NotifyReceived([]util.Address{addr}); slog.Check(err) {
	}
	return
}

func (w *Wallet) newChangeAddress(addrmgrNs walletdb.ReadWriteBucket, account uint32) (a util.Address, err error) {
	// As we're making a change address, we'll fetch the type of manager that is able to make p2wkh output as they're
	// the most efficient.
	scopes := w.Manager.ScopesForExternalAddrType(
		waddrmgr.WitnessPubKey,
	)
	var manager *waddrmgr.ScopedKeyManager
	if manager, err = w.Manager.FetchScopedKeyManager(scopes[0]); slog.Check(err) {
		return
	}
	// Get next chained change address from wallet for account.
	var addrs []waddrmgr.ManagedAddress
	if addrs, err = manager.NextInternalAddresses(addrmgrNs, account, 1); slog.Check(err) {
		return
	}
	return addrs[0].Address(), nil
}

// confirmed checks whether a transaction at height txHeight has met minconf confirmations for a blockchain at height
// curHeight.
func confirmed(minconf, txHeight, curHeight int32) bool {
	return confirms(txHeight, curHeight) >= minconf
}

// confirms returns the number of confirmations for a transaction in a block at height txHeight (or -1 for an
// unconfirmed tx) given the chain height curHeight.
func confirms(txHeight, curHeight int32) int32 {
	switch {
	case txHeight == -1, txHeight > curHeight:
		return 0
	default:
		return curHeight - txHeight + 1
	}
}

// AccountTotalReceivedResult is a single result for the Wallet.TotalReceivedForAccounts method.
type AccountTotalReceivedResult struct {
	AccountNumber    uint32
	AccountName      string
	TotalReceived    util.Amount
	LastConfirmation int32
}

// TotalReceivedForAccounts iterates through a wallet's transaction history,
// returning the total amount of Bitcoin received for all accounts.
func (w *Wallet) TotalReceivedForAccounts(scope waddrmgr.KeyScope,
	minConf int32) (results []AccountTotalReceivedResult, err error) {
	manager, err := w.Manager.FetchScopedKeyManager(scope)
	if err != nil {
		slog.Error(err)
		return nil, err
	}
	if err = walletdb.View(w.db, func(tx walletdb.ReadTx) (err error) {
		addrmgrNs := tx.ReadBucket(waddrmgrNamespaceKey)
		txmgrNs := tx.ReadBucket(wtxmgrNamespaceKey)
		syncBlock := w.Manager.SyncedTo()
		if err = manager.ForEachAccount(addrmgrNs, func(account uint32) (err error) {
			var accountName string
			if accountName, err = manager.AccountName(addrmgrNs, account); slog.Check(err) {
				return
			}
			results = append(results, AccountTotalReceivedResult{
				AccountNumber: account,
				AccountName:   accountName,
			})
			return nil
		}); slog.Check(err) {
			return
		}
		var stopHeight int32
		if minConf > 0 {
			stopHeight = syncBlock.Height - minConf + 1
		} else {
			stopHeight = -1
		}
		rangeFn := func(details []wtxmgr.TxDetails) (b bool, err error) {
			for i := range details {
				detail := &details[i]
				var outputAcct uint32
				var addrs []util.Address
				for _, cred := range detail.Credits {
					pkScript := detail.MsgTx.TxOut[cred.Index].PkScript
					if _, addrs, _, err = txscript.ExtractPkScriptAddrs(pkScript, w.chainParams); !slog.Check(err) &&
						len(addrs) > 0 {
						_, outputAcct, err = w.Manager.AddrAccount(addrmgrNs, addrs[0])
					}
					if err == nil {
						acctIndex := int(outputAcct)
						if outputAcct == waddrmgr.ImportedAddrAccount {
							acctIndex = len(results) - 1
						}
						res := &results[acctIndex]
						res.TotalReceived += cred.Amount
						res.LastConfirmation = confirms(
							detail.Block.Height, syncBlock.Height)
					}
				}
			}
			return
		}
		return w.TxStore.RangeTransactions(txmgrNs, 0, stopHeight, rangeFn)
	}); slog.Check(err) {
	}
	return
}

// TotalReceivedForAddr iterates through a wallet's transaction history,
// returning the total amount of bitcoins received for a single wallet
// address.
func (w *Wallet) TotalReceivedForAddr(addr util.Address, minConf int32) (amount util.Amount, err error) {
	if err = walletdb.View(w.db, func(tx walletdb.ReadTx) (err error) {
		txmgrNs := tx.ReadBucket(wtxmgrNamespaceKey)
		syncBlock := w.Manager.SyncedTo()
		var (
			addrStr    = addr.EncodeAddress()
			stopHeight int32
		)
		if minConf > 0 {
			stopHeight = syncBlock.Height - minConf + 1
		} else {
			stopHeight = -1
		}
		rangeFn := func(details []wtxmgr.TxDetails) (b bool, err error) {
			for i := range details {
				detail := &details[i]
				for _, cred := range detail.Credits {
					pkScript := detail.MsgTx.TxOut[cred.Index].PkScript
					var addrs []util.Address
					if _, addrs, _, err = txscript.ExtractPkScriptAddrs(pkScript, w.chainParams); slog.Check(err) {
						// An error creating addresses from the output script only
						// indicates a non-standard script, so ignore this credit.
						slog.Error(err)
						continue
					}
					for _, a := range addrs {
						if addrStr == a.EncodeAddress() {
							amount += cred.Amount
							break
						}
					}
				}
			}
			return
		}
		return w.TxStore.RangeTransactions(txmgrNs, 0, stopHeight, rangeFn)
	}); slog.Check(err) {
	}
	return
}

// SendOutputs creates and sends payment transactions. It returns the transaction hash upon success.
func (w *Wallet) SendOutputs(outputs []*wire.TxOut, account uint32, minconf int32, satPerKb util.Amount,
) (h *chainhash.Hash, err error) {
	// Ensure the outputs to be created adhere to the network's consensus rules.
	for _, output := range outputs {
		if err := txrules.CheckOutput(output, satPerKb); err != nil {
			return nil, err
		}
	}
	// Create the transaction and broadcast it to the network. The transaction will be added to the database in order to
	// ensure that we continue to re-broadcast the transaction upon restarts until it has been confirmed.
	var createdTx *txauthor.AuthoredTx
	if createdTx, err = w.CreateSimpleTx(account, outputs, minconf, satPerKb); slog.Check(err) {
		return
	}
	return w.publishTransaction(createdTx.Tx)
}

// SignatureError records the underlying error when validating a transaction input signature.
type SignatureError struct {
	InputIndex uint32
	Error      error
}

// SignTransaction uses secrets of the wallet, as well as additional secrets passed in by the caller, to create and add
// input signatures to a transaction.
//
// Transaction input script validation is used to confirm that all signatures are valid. For any invalid input, a
// SignatureError is added to the returns. The final error return is reserved for unexpected or fatal errors, such as
// being unable to determine a previous output script to redeem.
//
// The transaction pointed to by tx is modified by this function.
func (w *Wallet) SignTransaction(
	tx *wire.MsgTx,
	hashType txscript.SigHashType,
	additionalPrevScripts map[wire.OutPoint][]byte,
	additionalKeysByAddress map[string]*util.WIF,
	p2shRedeemScriptsByAddress map[string][]byte,
) (signErrors []SignatureError, err error) {
	if err = walletdb.View(w.db, func(dbtx walletdb.ReadTx) (err error) {
		addrmgrNs := dbtx.ReadBucket(waddrmgrNamespaceKey)
		txmgrNs := dbtx.ReadBucket(wtxmgrNamespaceKey)
		var txDetails *wtxmgr.TxDetails
		for i, txIn := range tx.TxIn {
			prevOutScript, ok := additionalPrevScripts[txIn.PreviousOutPoint]
			if !ok {
				prevHash := &txIn.PreviousOutPoint.Hash
				prevIndex := txIn.PreviousOutPoint.Index
				if txDetails, err = w.TxStore.TxDetails(txmgrNs, prevHash); slog.Check(err) {
					err = fmt.Errorf("cannot query previous transaction "+
						"details for %v: %v", txIn.PreviousOutPoint, err)
					slog.Debug(err)
					return

				}
				if txDetails == nil {
					err = fmt.Errorf("%v not found", txIn.PreviousOutPoint)
					slog.Debug(err)
					return
				}
				prevOutScript = txDetails.MsgTx.TxOut[prevIndex].PkScript
			}
			// Set up our callbacks that we pass to txscript so it can look up the appropriate keys and scripts by
			// address.
			getKey := txscript.KeyClosure(func(addr util.Address) (key *ec.PrivateKey, compressed bool, err error) {
				if len(additionalKeysByAddress) != 0 {
					addrStr := addr.EncodeAddress()
					wif, ok := additionalKeysByAddress[addrStr]
					if !ok {
						err = errors.New("no key for address")
						slog.Debug(err)
						return
					}
					key = wif.PrivKey
					compressed = wif.CompressPubKey
					return
				}
				var address waddrmgr.ManagedAddress
				if address, err = w.Manager.Address(addrmgrNs, addr); slog.Check(err) {
					return
				}
				var pka waddrmgr.ManagedPubKeyAddress
				if pka, ok = address.(waddrmgr.ManagedPubKeyAddress); !ok {
					err = fmt.Errorf("address %v is not a pubkey address", address.Address().EncodeAddress())
					slog.Debug(err)
					return
				}
				if key, err = pka.PrivKey(); slog.Check(err) {
					return
				}
				return key, pka.Compressed(), nil
			})
			getScript := txscript.ScriptClosure(func(addr util.Address) (scr []byte, err error) {
				// If keys were provided then we can only use the redeem scripts provided with our inputs, too.
				if len(additionalKeysByAddress) != 0 {
					addrStr := addr.EncodeAddress()
					if script, ok := p2shRedeemScriptsByAddress[addrStr]; !ok {
						err = errors.New("no script for address")
						slog.Debug(err)
						return
					} else {
						scr = script
						return
					}
				}
				var address waddrmgr.ManagedAddress
				if address, err = w.Manager.Address(addrmgrNs, addr); slog.Check(err) {
					return
				}
				var sa waddrmgr.ManagedScriptAddress
				if sa, ok = address.(waddrmgr.ManagedScriptAddress); !ok {
					err = errors.New("address is not a script address")
					slog.Debug(err)
					return
				}
				return sa.Script()
			})
			// SigHashSingle inputs can only be signed if there's a corresponding output. However this could be already
			// signed, so we always verify the output.
			if (hashType&txscript.SigHashSingle) != txscript.SigHashSingle || i < len(tx.TxOut) {
				var mergedScript []byte
				if mergedScript, err = txscript.SignTxOutput(w.ChainParams(), tx, i, prevOutScript, hashType, getKey,
					getScript, txIn.SignatureScript); slog.Check(err) {
					// Failure to sign isn't an error, it just means that the tx isn't complete.
					signErrors = append(signErrors, SignatureError{
						InputIndex: uint32(i),
						Error:      err,
					})
					continue
				}
				txIn.SignatureScript = mergedScript
			}
			// Either it was already signed or we just signed it. Find out if it is completely satisfied or still needs
			// more.
			if vm, err := txscript.NewEngine(prevOutScript, tx, i, txscript.StandardVerifyFlags, nil,
				nil, 0); !slog.Check(err) {
				err = vm.Execute()
			} else {
				signErrors = append(signErrors, SignatureError{
					InputIndex: uint32(i),
					Error:      err,
				})
			}
		}
		return
	}); slog.Check(err) {
	}
	return
}

// PublishTransaction sends the transaction to the consensus RPC server so it can be propagated to other nodes and
// eventually mined.
//
// This function is unstable and will be removed once syncing code is moved out of the wallet.
func (w *Wallet) PublishTransaction(tx *wire.MsgTx) (err error) {
	_, err = w.publishTransaction(tx)
	slog.Debug(err)
	return err
}

// publishTransaction is the private version of PublishTransaction which contains the primary logic required for
// publishing a transaction, updating the relevant database state, and finally possible removing the transaction from
// the database (along with cleaning up all inputs used, and outputs created) if the transaction is rejected by the back
// end.
func (w *Wallet) publishTransaction(tx *wire.MsgTx) (txid *chainhash.Hash, err error) {
	var server chain.Interface
	if server, err = w.requireChainClient(); slog.Check(err) {
		return
	}
	// As we aim for this to be general reliable transaction broadcast API, we'll write this tx to disk as an
	// unconfirmed transaction. This way, upon restarts, we'll always rebroadcast it, and also add it to our set of
	// records.
	var txRec *wtxmgr.TxRecord
	if txRec, err = wtxmgr.NewTxRecordFromMsgTx(tx, time.Now()); slog.Check(err) {
		return
	}
	if err = walletdb.Update(w.db, func(dbTx walletdb.ReadWriteTx) (err error) {
		return w.addRelevantTx(dbTx, txRec, nil)
	}); slog.Check(err) {
		return
	}
	if txid, err = server.SendRawTransaction(tx, false); slog.Check(err) {
	}
	_ = txid
	switch {
	case err == nil:
		return
	// The following are errors returned from pod's mempool.
	case strings.Contains(err.Error(), "spent"):
		fallthrough
	case strings.Contains(err.Error(), "orphan"):
		fallthrough
	case strings.Contains(err.Error(), "conflict"):
		fallthrough
	// The following errors are returned from bitcoind's mempool.
	case strings.Contains(err.Error(), "fee not met"):
		fallthrough
	case strings.Contains(err.Error(), "Missing inputs"):
		fallthrough
	case strings.Contains(err.Error(), "already in block chain"):
		// If the transaction was rejected, then we'll remove it from the txstore, as otherwise, we'll attempt to
		// continually re-broadcast it, and the utxo state of the wallet won't be accurate.
		if dbErr := walletdb.Update(w.db, func(dbTx walletdb.ReadWriteTx) (err error) {
			txmgrNs := dbTx.ReadWriteBucket(wtxmgrNamespaceKey)
			return w.TxStore.RemoveUnminedTx(txmgrNs, txRec)
		}); slog.Check(err) {
			err = fmt.Errorf("unable to broadcast tx: %v, unable to remove invalid tx: %v", err, dbErr)
			slog.Debug(err)
			return
		}
		return
	default:
		return
	}
}

// ChainParams returns the network parameters for the blockchain the wallet belongs to.
func (w *Wallet) ChainParams() *netparams.Params {
	return w.chainParams
}

// Database returns the underlying walletdb database. This method is provided in order to allow applications wrapping
// btcwallet to store app-specific data with the wallet's database.
func (w *Wallet) Database() walletdb.DB {
	return w.db
}

// Create creates an new wallet, writing it to an empty database. If the passed seed is non-nil, it is used. Otherwise,
// a secure random seed of the recommended length is generated.
func Create(db walletdb.DB, pubPass, privPass, seed []byte, params *netparams.Params,
	birthday time.Time) (err error) {
	// If a seed was provided, ensure that it is of valid length. Otherwise, we generate a random seed for the wallet
	// with the recommended seed length.
	if seed == nil {
		var hdSeed []byte
		if hdSeed, err = hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen); slog.Check(err) {
			return
		}
		seed = hdSeed
	}
	if len(seed) < hdkeychain.MinSeedBytes ||
		len(seed) > hdkeychain.MaxSeedBytes {
		return hdkeychain.ErrInvalidSeedLen
	}
	return walletdb.Update(db, func(tx walletdb.ReadWriteTx) (err error) {
		var addrmgrNs walletdb.ReadWriteBucket
		if addrmgrNs, err = tx.CreateTopLevelBucket(waddrmgrNamespaceKey); slog.Check(err) {
			return
		}
		var txmgrNs walletdb.ReadWriteBucket
		if txmgrNs, err = tx.CreateTopLevelBucket(wtxmgrNamespaceKey); slog.Check(err) {
			return
		}
		if err = waddrmgr.Create(
			addrmgrNs, seed, pubPass, privPass, params, nil,
			birthday,
		); slog.Check(err) {
			return
		}
		return wtxmgr.Create(txmgrNs)
	})
}

// Open loads an already-created wallet from the passed database and namespaces.
func Open(db walletdb.DB, pubPass []byte, cbs *waddrmgr.OpenCallbacks,
	params *netparams.Params, recoveryWindow uint32, podConfig *pod.Config) (w *Wallet, err error) {
	slog.Warn("opening wallet")
	if err = walletdb.View(db, func(tx walletdb.ReadTx) (err error) {
		waddrmgrBucket := tx.ReadBucket(waddrmgrNamespaceKey)
		if waddrmgrBucket == nil {
			err = errors.New("missing address manager namespace")
			slog.Debug(err)
			return
		}
		wtxmgrBucket := tx.ReadBucket(wtxmgrNamespaceKey)
		if wtxmgrBucket == nil {
			err = errors.New("missing transaction manager namespace")
			slog.Debug(err)
		}
		return
	}); slog.Check(err) {
		return
	}
	// Perform upgrades as necessary. Each upgrade is done under its own transaction, which is managed by each package
	// itself, so the entire DB is passed instead of passing already opened write transaction.
	//
	// This will need to change later when upgrades in one package depend on data in another (such as removing chain
	// synchronization from address manager).
	if err = waddrmgr.DoUpgrades(db, waddrmgrNamespaceKey, pubPass, params, cbs); slog.Check(err) {
		return
	}
	if err = wtxmgr.DoUpgrades(db, wtxmgrNamespaceKey); slog.Check(err) {
		return
	}
	// Open database abstraction instances
	var (
		addrMgr *waddrmgr.Manager
		txMgr   *wtxmgr.Store
	)
	if err = walletdb.View(db, func(tx walletdb.ReadTx) (err error) {
		addrmgrNs := tx.ReadBucket(waddrmgrNamespaceKey)
		txmgrNs := tx.ReadBucket(wtxmgrNamespaceKey)
		if addrMgr, err = waddrmgr.Open(addrmgrNs, pubPass, params); slog.Check(err) {
			return
		}
		txMgr, err = wtxmgr.Open(txmgrNs, params)
		return
	}); slog.Check(err) {
		return
	}
	slog.Trace("opened wallet") // TODO: log balance? last sync height?
	w = &Wallet{
		publicPassphrase:    pubPass,
		db:                  db,
		Manager:             addrMgr,
		TxStore:             txMgr,
		lockedOutpoints:     map[wire.OutPoint]struct{}{},
		recoveryWindow:      recoveryWindow,
		rescanAddJob:        make(chan *RescanJob),
		rescanBatch:         make(chan *rescanBatch),
		rescanNotifications: make(chan interface{}),
		rescanProgress:      make(chan *RescanProgressMsg),
		rescanFinished:      make(chan *RescanFinishedMsg),
		createTxRequests:    make(chan createTxRequest),
		unlockRequests:      make(chan unlockRequest),
		lockRequests:        make(chan struct{}),
		holdUnlockRequests:  make(chan chan heldUnlock),
		lockState:           make(chan bool),
		changePassphrase:    make(chan changePassphraseRequest),
		changePassphrases:   make(chan changePassphrasesRequest),
		chainParams:         params,
		PodConfig:           podConfig,
		quit:                make(chan struct{}),
	}
	w.NtfnServer = newNotificationServer(w)
	w.TxStore.NotifyUnspent = func(hash *chainhash.Hash, index uint32) {
		w.NtfnServer.notifyUnspentOutput(0, hash, index)
	}
	return
}
