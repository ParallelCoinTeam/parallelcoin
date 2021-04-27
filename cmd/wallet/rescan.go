package wallet

import (
	"github.com/p9c/log"
	"github.com/p9c/matrjoska/pkg/btcaddr"
	"github.com/p9c/matrjoska/pkg/chainclient"
	"github.com/p9c/matrjoska/pkg/txscript"
	"github.com/p9c/matrjoska/pkg/waddrmgr"
	"github.com/p9c/matrjoska/pkg/wire"
	"github.com/p9c/matrjoska/pkg/wtxmgr"
)

// RescanProgressMsg reports the current progress made by a rescan for a set of
// wallet addresses.
type RescanProgressMsg struct {
	Addresses    []btcaddr.Address
	Notification *chainclient.RescanProgress
}

// RescanFinishedMsg reports the addresses that were rescanned when a
// rescanfinished message was received rescanning a batch of addresses.
type RescanFinishedMsg struct {
	Addresses    []btcaddr.Address
	Notification *chainclient.RescanFinished
}

// RescanJob is a job to be processed by the RescanManager. The job includes a
// set of wallet addresses, a starting height to begin the rescan, and outpoints
// spendable by the addresses thought to be unspent. After the rescan completes,
// the error result of the rescan RPC is sent on the Err channel.
type RescanJob struct {
	InitialSync bool
	Addrs       []btcaddr.Address
	OutPoints   map[wire.OutPoint]btcaddr.Address
	BlockStamp  waddrmgr.BlockStamp
	err         chan error
}

// rescanBatch is a collection of one or more RescanJobs that were merged
// together before a rescan is performed.
type rescanBatch struct {
	initialSync bool
	addrs       []btcaddr.Address
	outpoints   map[wire.OutPoint]btcaddr.Address
	bs          waddrmgr.BlockStamp
	errChans    []chan error
}

// SubmitRescan submits a RescanJob to the RescanManager. A channel is returned
// with the final error of the rescan. The channel is buffered and does not need
// to be read to prevent a deadlock.
func (w *Wallet) SubmitRescan(job *RescanJob) <-chan error {
	errChan := make(chan error, 1)
	job.err = errChan
	w.rescanAddJob <- job
	return errChan
}

// batch creates the rescanBatch for a single rescan job.
func (job *RescanJob) batch() *rescanBatch {
	return &rescanBatch{
		initialSync: job.InitialSync,
		addrs:       job.Addrs,
		outpoints:   job.OutPoints,
		bs:          job.BlockStamp,
		errChans:    []chan error{job.err},
	}
}

// merge merges the work from k into j, setting the starting height to the
// minimum of the two jobs. This method does not check for duplicate addresses
// or outpoints.
func (b *rescanBatch) merge(job *RescanJob) {
	if job.InitialSync {
		b.initialSync = true
	}
	b.addrs = append(b.addrs, job.Addrs...)
	for op, addr := range job.OutPoints {
		b.outpoints[op] = addr
	}
	if job.BlockStamp.Height < b.bs.Height {
		b.bs = job.BlockStamp
	}
	b.errChans = append(b.errChans, job.err)
}

// done iterates through all error channels, duplicating sending the error to
// inform callers that the rescan finished (or could not complete due to an
// error).
func (b *rescanBatch) done(e error) {
	for _, c := range b.errChans {
		c <- e
	}
}

// rescanBatchHandler handles incoming rescan request, serializing rescan
// submissions, and possibly batching many waiting requests together so they can
// be handled by a single rescan after the current one completes.
func (w *Wallet) rescanBatchHandler() {
	var curBatch, nextBatch *rescanBatch
	quit := w.quitChan()
out:
	for {
		select {
		case job := <-w.rescanAddJob:
			if curBatch == nil {
				// Set current batch as this job and send request.
				curBatch = job.batch()
				w.rescanBatch <- curBatch
			} else {
				// Create next batch if it doesn't exist, or merge the job.
				if nextBatch == nil {
					nextBatch = job.batch()
				} else {
					nextBatch.merge(job)
				}
			}
		case n := <-w.rescanNotifications:
			switch n := n.(type) {
			case *chainclient.RescanProgress:
				if curBatch == nil {
					W.Ln(
						"received rescan progress notification but no rescan currently running",
					)
					continue
				}
				w.rescanProgress <- &RescanProgressMsg{
					Addresses:    curBatch.addrs,
					Notification: n,
				}
			case *chainclient.RescanFinished:
				if curBatch == nil {
					W.Ln(
						"received rescan finished notification but no rescan currently running",
					)
					continue
				}
				w.rescanFinished <- &RescanFinishedMsg{
					Addresses:    curBatch.addrs,
					Notification: n,
				}
				curBatch, nextBatch = nextBatch, nil
				if curBatch != nil {
					w.rescanBatch <- curBatch
				}
			default:
				// Unexpected message
				panic(n)
			}
		case <-quit.Wait():
			break out
		}
	}
	w.wg.Done()
}

// rescanProgressHandler handles notifications for partially and fully completed
// rescans by marking each rescanned address as partially or fully synced.
func (w *Wallet) rescanProgressHandler() {
	quit := w.quitChan()
out:
	for {
		// These can't be processed out of order since both chans are unbuffured and are
		// sent from same context (the batch handler).
		select {
		case msg := <-w.rescanProgress:
			n := msg.Notification
			I.F(
				"rescanned through block %v (height %d)",
				n.Hash, n.Height,
			)
		case msg := <-w.rescanFinished:
			n := msg.Notification
			addrs := msg.Addresses
			noun := log.PickNoun(len(addrs), "address", "addresses")
			I.F(
				"finished rescan for %d %s (synced to block %s, height %d)",
				len(addrs), noun, n.Hash, n.Height,
			)
			go w.resendUnminedTxs()
		case <-quit.Wait():
			break out
		}
	}
	w.wg.Done()
}

// rescanRPCHandler reads batch jobs sent by rescanBatchHandler and sends the
// RPC requests to perform a rescan. New jobs are not read until a rescan
// finishes.
func (w *Wallet) rescanRPCHandler() {
	chainClient, e := w.requireChainClient()
	if e != nil {
		E.Ln("rescanRPCHandler called without an RPC client", e)
		w.wg.Done()
		return
	}
	quit := w.quitChan()
out:
	for {
		select {
		case batch := <-w.rescanBatch:
			// Log the newly-started rescan.
			numAddrs := len(batch.addrs)
			noun := log.PickNoun(numAddrs, "address", "addresses")
			I.F(
				"started rescan from block %v (height %d) for %d %s",
				batch.bs.Hash, batch.bs.Height, numAddrs, noun,
			)
			e := chainClient.Rescan(
				&batch.bs.Hash, batch.addrs,
				batch.outpoints,
			)
			if e != nil {
				E.F(
					"rescan for %d %s failed: %v", numAddrs, noun, e,
				)
			}
			batch.done(e)
		case <-quit.Wait():
			break out
		}
	}
	w.wg.Done()
}

// Rescan begins a rescan for all active addresses and unspent outputs of a
// wallet. This is intended to be used to sync a wallet back up to the current
// best block in the main chain, and is considered an initial sync rescan.
func (w *Wallet) Rescan(addrs []btcaddr.Address, unspent []wtxmgr.Credit) (e error) {
	return w.rescanWithTarget(addrs, unspent, nil)
}

// rescanWithTarget performs a rescan starting at the optional startStamp. If
// none is provided, the rescan will begin from the manager's sync tip.
func (w *Wallet) rescanWithTarget(
	addrs []btcaddr.Address,
	unspent []wtxmgr.Credit, startStamp *waddrmgr.BlockStamp,
) (e error) {
	outpoints := make(map[wire.OutPoint]btcaddr.Address, len(unspent))
	for _, output := range unspent {
		var outputAddrs []btcaddr.Address
		_, outputAddrs, _, e = txscript.ExtractPkScriptAddrs(
			output.PkScript, w.chainParams,
		)
		if e != nil {
			return e
		}
		outpoints[output.OutPoint] = outputAddrs[0]
	}
	// If a start block stamp was provided, we will use that as the initial starting
	// point for the rescan.
	if startStamp == nil {
		startStamp = &waddrmgr.BlockStamp{}
		*startStamp = w.Manager.SyncedTo()
	}
	job := &RescanJob{
		InitialSync: true,
		Addrs:       addrs,
		OutPoints:   outpoints,
		BlockStamp:  *startStamp,
	}
	// Submit merged job and block until rescan completes.
	return <-w.SubmitRescan(job)
}
