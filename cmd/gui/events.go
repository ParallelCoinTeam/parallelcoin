package gui

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/p9c/matrjoska/pkg/amt"
	"github.com/p9c/matrjoska/pkg/chainrpc/p2padvt"
	"github.com/p9c/matrjoska/pkg/transport"
	"github.com/p9c/matrjoska/pkg/wire"

	"github.com/p9c/matrjoska/pkg/btcjson"
	"github.com/p9c/matrjoska/pkg/chainhash"
	"github.com/p9c/matrjoska/pkg/rpcclient"
	"github.com/p9c/matrjoska/pkg/util"
)

func (wg *WalletGUI) WalletAndClientRunning() bool {
	running := wg.wallet.Running() && wg.WalletClient != nil && !wg.WalletClient.Disconnected()
	// D.Ln("wallet and wallet rpc client are running?", running)
	return running
}

func (wg *WalletGUI) Advertise() (e error) {
	if wg.node.Running() && wg.cx.Config.Discovery.True() {
		// I.Ln("sending out p2p advertisment")
		if e = wg.multiConn.SendMany(
			p2padvt.Magic,
			transport.GetShards(p2padvt.Get(uint64(wg.cx.Config.UUID.V()), (wg.cx.Config.P2PListeners.S())[0])),
		); E.Chk(e) {
		}
	}
	return
}

// func (wg *WalletGUI) Tickers() {
// 	first := true
// 	D.Ln("updating best block")
// 	var e error
// 	var height int32
// 	var h *chainhash.Hash
// 	if h, height, e = wg.ChainClient.GetBestBlock(); E.Chk(e) {
// 		// interrupt.Request()
// 		return
// 	}
// 	D.Ln(h, height)
// 	wg.State.SetBestBlockHeight(height)
// 	wg.State.SetBestBlockHash(h)
// 	go func() {
// 		var e error
// 		seconds := time.Tick(time.Second * 3)
// 		fiveSeconds := time.Tick(time.Second * 5)
// 	totalOut:
// 		for {
// 		preconnect:
// 			for {
// 				select {
// 				case <-seconds:
// 					D.Ln("preconnect loop")
// 					if e = wg.Advertise(); D.Chk(e) {
// 					}
// 					if wg.ChainClient != nil {
// 						wg.ChainClient.Disconnect()
// 						wg.ChainClient.Shutdown()
// 						wg.ChainClient = nil
// 					}
// 					if wg.WalletClient != nil {
// 						wg.WalletClient.Disconnect()
// 						wg.WalletClient.Shutdown()
// 						wg.WalletClient = nil
// 					}
// 					if !wg.node.Running() {
// 						break
// 					}
// 					break preconnect
// 				case <-fiveSeconds:
// 					continue
// 				case <-wg.quit.Wait():
// 					break totalOut
// 				}
// 			}
// 		out:
// 			for {
// 				select {
// 				case <-seconds:
// 					if e = wg.Advertise(); D.Chk(e) {
// 					}
// 					if !wg.cx.IsCurrent() {
// 						continue
// 					} else {
// 						wg.cx.Syncing.Store(false)
// 					}
// 					D.Ln("---------------------- ready", wg.ready.Load())
// 					D.Ln("---------------------- WalletAndClientRunning", wg.WalletAndClientRunning())
// 					D.Ln("---------------------- stateLoaded", wg.stateLoaded.Load())
// 					wg.node.Start()
// 					if e = wg.writeWalletCookie(); E.Chk(e) {
// 					}
// 					wg.wallet.Start()
// 					D.Ln("connecting to chain")
// 					if e = wg.chainClient(); E.Chk(e) {
// 						break
// 					}
// 					if wg.wallet.Running() { // && wg.WalletClient == nil {
// 						D.Ln("connecting to wallet")
// 						if e = wg.walletClient(); E.Chk(e) {
// 							break
// 						}
// 					}
// 					if !wg.node.Running() {
// 						D.Ln("breaking out node not running")
// 						break out
// 					}
// 					if wg.ChainClient == nil {
// 						D.Ln("breaking out chainclient is nil")
// 						break out
// 					}
// 					// if  wg.WalletClient == nil {
// 					// 	D.Ln("breaking out walletclient is nil")
// 					// 	break out
// 					// }
// 					if wg.ChainClient.Disconnected() {
// 						D.Ln("breaking out chainclient disconnected")
// 						break out
// 					}
// 					// if wg.WalletClient.Disconnected() {
// 					// 	D.Ln("breaking out walletclient disconnected")
// 					// 	break out
// 					// }
// 					// var e error
// 					if first {
// 						wg.updateChainBlock()
// 						wg.invalidate <- struct{}{}
// 					}
//
// 					if wg.WalletAndClientRunning() {
// 						if first {
// 							wg.processWalletBlockNotification()
// 						}
// 						// if wg.stateLoaded.Load() { // || wg.currentReceiveGetNew.Load() {
// 						// 	wg.ReceiveAddressbook = func(gtx l.Context) l.Dimensions {
// 						// 		var widgets []l.Widget
// 						// 		widgets = append(widgets, wg.ReceivePage.GetAddressbookHistoryCards("DocBg")...)
// 						// 		le := func(gtx l.Context, index int) l.Dimensions {
// 						// 			return widgets[index](gtx)
// 						// 		}
// 						// 		return wg.Flex().Rigid(
// 						// 			wg.lists["receiveAddresses"].Length(len(widgets)).Vertical().
// 						// 				ListElement(le).Fn,
// 						// 		).Fn(gtx)
// 						// 	}
// 						// }
// 						if wg.stateLoaded.Load() && !wg.State.IsReceivingAddress() { // || wg.currentReceiveGetNew.Load() {
// 							wg.GetNewReceivingAddress()
// 							if wg.currentReceiveQRCode == nil || wg.currentReceiveRegenerate.Load() { // || wg.currentReceiveGetNew.Load() {
// 								wg.GetNewReceivingQRCode(wg.ReceivePage.urn)
// 							}
// 						}
// 					}
// 					wg.invalidate <- struct{}{}
// 					first = false
// 				case <-fiveSeconds:
// 				case <-wg.quit.Wait():
// 					break totalOut
// 				}
// 			}
// 		}
// 	}()
// }

func (wg *WalletGUI) updateThingies() (e error) {
	// update the configuration
	var b []byte
	if b, e = ioutil.ReadFile(wg.cx.Config.ConfigFile.V()); !E.Chk(e) {
		if e = json.Unmarshal(b, wg.cx.Config); !E.Chk(e) {
			return
		}
	}
	return
}
func (wg *WalletGUI) updateChainBlock() {
	D.Ln("processChainBlockNotification")
	var e error
	if wg.ChainClient != nil && wg.ChainClient.Disconnected() || wg.ChainClient.Disconnected() {
		D.Ln("connecting ChainClient")
		if e = wg.chainClient(); E.Chk(e) {
			return
		}
	}
	var h *chainhash.Hash
	var height int32
	D.Ln("updating best block")
	if h, height, e = wg.ChainClient.GetBestBlock(); E.Chk(e) {
		// interrupt.Request()
		return
	}
	D.Ln(h, height)
	wg.State.SetBestBlockHeight(height)
	wg.State.SetBestBlockHash(h)
}

func (wg *WalletGUI) processChainBlockNotification(hash *chainhash.Hash, height int32, t time.Time) {
	D.Ln("processChainBlockNotification")
	wg.State.SetBestBlockHeight(height)
	wg.State.SetBestBlockHash(hash)
	// if wg.WalletAndClientRunning() {
	// 	wg.processWalletBlockNotification()
	// }
}

func (wg *WalletGUI) processWalletBlockNotification() bool {
	D.Ln("processWalletBlockNotification")
	if !wg.WalletAndClientRunning() {
		D.Ln("wallet and client not running")
		return false
	}
	// check account balance
	var unconfirmed amt.Amount
	var e error
	if unconfirmed, e = wg.WalletClient.GetUnconfirmedBalance("default"); E.Chk(e) {
		return false
	}
	wg.State.SetBalanceUnconfirmed(unconfirmed.ToDUO())
	var confirmed amt.Amount
	if confirmed, e = wg.WalletClient.GetBalance("default"); E.Chk(e) {
		return false
	}
	wg.State.SetBalance(confirmed.ToDUO())
	var atr []btcjson.ListTransactionsResult
	// str := wg.State.allTxs.Load()
	if atr, e = wg.WalletClient.ListTransactionsCount("default", 2<<32); E.Chk(e) {
		return false
	}
	// D.Ln(len(atr))
	// wg.State.SetAllTxs(append(str, atr...))
	wg.State.SetAllTxs(atr)
	wg.txMx.Lock()
	wg.txHistoryList = wg.State.filteredTxs.Load()
	atrl := 10
	if len(atr) < atrl {
		atrl = len(atr)
	}
	wg.txMx.Unlock()
	wg.RecentTransactions(10, "recent")
	wg.RecentTransactions(-1, "history")
	return true
}

func (wg *WalletGUI) forceUpdateChain() {
	wg.updateChainBlock()
	var e error
	var height int32
	var tip *chainhash.Hash
	if tip, height, e = wg.ChainClient.GetBestBlock(); E.Chk(e) {
		return
	}
	var block *wire.Block
	if block, e = wg.ChainClient.GetBlock(tip); E.Chk(e) {
	}
	t := block.Header.Timestamp
	wg.processChainBlockNotification(tip, height, t)
}

func (wg *WalletGUI) ChainNotifications() *rpcclient.NotificationHandlers {
	return &rpcclient.NotificationHandlers{
		// OnClientConnected: func() {
		// 	// go func() {
		// 	D.Ln("(((NOTIFICATION))) CHAIN CLIENT CONNECTED!")
		// 	wg.cx.Syncing.Store(true)
		// 	wg.forceUpdateChain()
		// 	wg.processWalletBlockNotification()
		// 	wg.RecentTransactions(10, "recent")
		// 	wg.RecentTransactions(-1, "history")
		// 	wg.invalidate <- struct{}{}
		// 	wg.cx.Syncing.Store(false)
		// },
		// OnBlockConnected: func(hash *chainhash.Hash, height int32, t time.Time) {
		// 	if wg.cx.Syncing.Load() {
		// 		return
		// 	}
		// 	D.Ln("(((NOTIFICATION))) chain OnBlockConnected", hash, height, t)
		// 	wg.processChainBlockNotification(hash, height, t)
		// 	// wg.processWalletBlockNotification()
		// 	// todo: send system notification of new block, set configuration to disable also
		// 	// if wg.WalletAndClientRunning() {
		// 	// 	var e error
		// 	// 	if _, e = wg.WalletClient.RescanBlocks([]chainhash.Hash{*hash}); E.Chk(e) {
		// 	// 	}
		// 	// }
		// 	wg.RecentTransactions(10, "recent")
		// 	wg.RecentTransactions(-1, "history")
		// 	wg.invalidate <- struct{}{}
		// },
		OnFilteredBlockConnected: func(height int32, header *wire.BlockHeader, txs []*util.Tx) {
			nbh := header.BlockHash()
			wg.processChainBlockNotification(&nbh, height, header.Timestamp)
			// if time.Now().Sub(time.Unix(wg.lastUpdated.Load(), 0)) < time.Second {
			// 	return
			// }
			wg.lastUpdated.Store(time.Now().Unix())
			hash := header.BlockHash()
			D.Ln(
				"(((NOTIFICATION))) OnFilteredBlockConnected hash", hash, "POW hash:",
				header.BlockHashWithAlgos(height), "height", height,
			)
			// D.S(txs)
			if wg.processWalletBlockNotification() {
			}
			// filename := filepath.Join(*wg.cx.Config.DataDir, "state.json")
			// if e := wg.State.Save(filename, wg.cx.Config.WalletPass); E.Chk(e) {
			// }
			// if wg.WalletAndClientRunning() {
			// 	var e error
			// 	if _, e = wg.WalletClient.RescanBlocks([]chainhash.Hash{hash}); E.Chk(e) {
			// 	}
			// }
			wg.RecentTransactions(10, "recent")
			wg.RecentTransactions(-1, "history")
			wg.invalidate <- struct{}{}
		},
		// OnBlockDisconnected: func(hash *chainhash.Hash, height int32, t time.Time) {
		// 	if wg.cx.Syncing.Load() {
		// 		return
		// 	}
		// 	D.Ln("(((NOTIFICATION))) OnBlockDisconnected", hash, height, t)
		// 	wg.forceUpdateChain()
		// 	if wg.processWalletBlockNotification() {
		// 	}
		// 	wg.RecentTransactions(10, "recent")
		// 	wg.RecentTransactions(-1, "history")
		// 	wg.invalidate <- struct{}{}
		// },
		// OnFilteredBlockDisconnected: func(height int32, header *wire.BlockHeader) {
		// 	if wg.cx.Syncing.Load() {
		// 		return
		// 	}
		// 	D.Ln("(((NOTIFICATION))) OnFilteredBlockDisconnected", height, header)
		// 	wg.forceUpdateChain()
		// 	if wg.processWalletBlockNotification() {
		// 	}
		// 	wg.RecentTransactions(10, "recent")
		// 	wg.RecentTransactions(-1, "history")
		// 	wg.invalidate <- struct{}{}
		// },
		// OnRecvTx: func(transaction *util.Tx, details *btcjson.BlockDetails) {
		// 	if wg.cx.Syncing.Load() {
		// 		return
		// 	}
		// 	D.Ln("(((NOTIFICATION))) OnRecvTx", transaction, details)
		// 	wg.forceUpdateChain()
		// 	if wg.processWalletBlockNotification() {
		// 	}
		// 	wg.RecentTransactions(10, "recent")
		// 	wg.RecentTransactions(-1, "history")
		// 	wg.invalidate <- struct{}{}
		// },
		// OnRedeemingTx: func(transaction *util.Tx, details *btcjson.BlockDetails) {
		// 	if wg.cx.Syncing.Load() {
		// 		return
		// 	}
		// 	D.Ln("(((NOTIFICATION))) OnRedeemingTx", transaction, details)
		// 	wg.forceUpdateChain()
		// 	if wg.processWalletBlockNotification() {
		// 	}
		// 	wg.RecentTransactions(10, "recent")
		// 	wg.RecentTransactions(-1, "history")
		// 	wg.invalidate <- struct{}{}
		// },
		// OnRelevantTxAccepted: func(transaction []byte) {
		// 	if wg.cx.Syncing.Load() {
		// 		return
		// 	}
		// 	D.Ln("(((NOTIFICATION))) OnRelevantTxAccepted", transaction)
		// 	wg.forceUpdateChain()
		// 	if wg.processWalletBlockNotification() {
		// 	}
		// 	wg.RecentTransactions(10, "recent")
		// 	wg.RecentTransactions(-1, "history")
		// 	wg.invalidate <- struct{}{}
		// },
		// OnRescanFinished: func(hash *chainhash.Hash, height int32, blkTime time.Time) {
		// 	if wg.cx.Syncing.Load() {
		// 		return
		// 	}
		// 	D.Ln("(((NOTIFICATION))) OnRescanFinished", hash, height, blkTime)
		// 	wg.processChainBlockNotification(hash, height, blkTime)
		// 	// update best block height
		// 	// wg.processWalletBlockNotification()
		// 	// stop showing syncing indicator
		// 	if wg.processWalletBlockNotification() {
		// 	}
		// 	wg.RecentTransactions(10, "recent")
		// 	wg.RecentTransactions(-1, "history")
		// 	wg.invalidate <- struct{}{}
		// },
		// OnRescanProgress: func(hash *chainhash.Hash, height int32, blkTime time.Time) {
		// 	D.Ln("(((NOTIFICATION))) OnRescanProgress", hash, height, blkTime)
		// 	// update best block height
		// 	// wg.processWalletBlockNotification()
		// 	// set to show syncing indicator
		// 	if wg.processWalletBlockNotification() {
		// 	}
		// 	wg.Syncing.Store(true)
		// 	wg.RecentTransactions(10, "recent")
		// 	wg.RecentTransactions(-1, "history")
		// 	wg.invalidate <- struct{}{}
		// },
		OnTxAccepted: func(hash *chainhash.Hash, amount amt.Amount) {
			// if wg.syncing.Load() {
			// 	D.Ln("OnTxAccepted but we are syncing")
			// 	return
			// }
			D.Ln("(((NOTIFICATION))) OnTxAccepted")
			D.Ln(hash, amount)
			// if wg.processWalletBlockNotification() {
			// }
			wg.RecentTransactions(10, "recent")
			wg.RecentTransactions(-1, "history")
			wg.invalidate <- struct{}{}
		},
		// OnTxAcceptedVerbose: func(txDetails *btcjson.TxRawResult) {
		// 	if wg.cx.Syncing.Load() {
		// 		return
		// 	}
		// 	D.Ln("(((NOTIFICATION))) OnTxAcceptedVerbose")
		// 	D.S(txDetails)
		// 	if wg.processWalletBlockNotification() {
		// 	}
		// 	wg.RecentTransactions(10, "recent")
		// 	wg.RecentTransactions(-1, "history")
		// 	wg.invalidate <- struct{}{}
		// },
		// OnPodConnected: func(connected bool) {
		// 	if wg.cx.Syncing.Load() {
		// 		return
		// 	}
		// 	D.Ln("(((NOTIFICATION))) OnPodConnected", connected)
		// 	wg.forceUpdateChain()
		// 	if wg.processWalletBlockNotification() {
		// 	}
		// 	wg.RecentTransactions(10, "recent")
		// 	wg.RecentTransactions(-1, "history")
		// 	wg.invalidate <- struct{}{}
		// },
		// OnAccountBalance: func(account string, balance util.Amount, confirmed bool) {
		// 	if wg.cx.Syncing.Load() {
		// 		return
		// 	}
		// 	D.Ln("OnAccountBalance")
		// 	// what does this actually do
		// 	D.Ln(account, balance, confirmed)
		// },
		// OnWalletLockState: func(locked bool) {
		// 	if wg.cx.Syncing.Load() {
		// 		return
		// 	}
		// 	D.Ln("OnWalletLockState", locked)
		// 	// switch interface to unlock page
		// 	wg.forceUpdateChain()
		// 	if wg.processWalletBlockNotification() {
		// 	}
		// 	// TODO: lock when idle... how to get trigger for idleness in UI?
		// },
		// OnUnknownNotification: func(method string, params []json.RawMessage) {
		// 	if wg.cx.Syncing.Load() {
		// 		return
		// 	}
		// 	D.Ln("(((NOTIFICATION))) OnUnknownNotification", method, params)
		// 	wg.forceUpdateChain()
		// 	if wg.processWalletBlockNotification() {
		// 	}
		// },
	}
	
}

func (wg *WalletGUI) WalletNotifications() *rpcclient.NotificationHandlers {
	// if !wg.wallet.Running() || wg.WalletClient == nil || wg.WalletClient.Disconnected() {
	// 	return nil
	// }
	// var updating bool
	return &rpcclient.NotificationHandlers{
		// OnClientConnected: func() {
		// 	if wg.cx.Syncing.Load() {
		// 		return
		// 	}
		// 	if updating {
		// 		return
		// 	}
		// 	D.Ln("(((NOTIFICATION))) wallet client connected, running initial processes")
		// 	for !wg.processWalletBlockNotification() {
		// 		time.Sleep(time.Second)
		// 		D.Ln("(((NOTIFICATION))) retry attempting to update wallet transactions")
		// 	}
		// 	filename := filepath.Join(wg.cx.DataDir, "state.json")
		// 	if e := wg.State.Save(filename, wg.cx.Config.WalletPass); E.Chk(e) {
		// 	}
		// 	wg.invalidate <- struct{}{}
		// 	updating = false
		// },
		// OnBlockConnected: func(hash *chainhash.Hash, height int32, t time.Time) {
		// 	if wg.Syncing.Load() {
		// 		return
		// 	}
		// 	D.Ln("(((NOTIFICATION))) wallet OnBlockConnected", hash, height, t)
		// 	wg.processWalletBlockNotification()
		// 	filename := filepath.Join(wg.cx.DataDir, "state.json")
		// 	if e := wg.State.Save(filename, wg.cx.Config.WalletPass); E.Chk(e) {
		// 	}
		// 	wg.invalidate <- struct{}{}
		// },
		// OnFilteredBlockConnected: func(height int32, header *wire.BlockHeader, txs []*util.Tx) {
		// 	if wg.Syncing.Load() {
		// 		return
		// 	}
		// 	D.Ln(
		// 		"(((NOTIFICATION))) wallet OnFilteredBlockConnected hash", header.BlockHash(), "POW hash:",
		// 		header.BlockHashWithAlgos(height), "height", height,
		// 	)
		// 	// D.S(txs)
		// 	nbh := header.BlockHash()
		// 	wg.processChainBlockNotification(&nbh, height, header.Timestamp)
		// 	if wg.processWalletBlockNotification() {
		// 	}
		// 	filename := filepath.Join(wg.cx.DataDir, "state.json")
		// 	if e := wg.State.Save(filename, wg.cx.Config.WalletPass); E.Chk(e) {
		// 	}
		// 	wg.invalidate <- struct{}{}
		// },
		// OnBlockDisconnected: func(hash *chainhash.Hash, height int32, t time.Time) {
		// 	if wg.Syncing.Load() {
		// 		return
		// 	}
		// 	D.Ln("(((NOTIFICATION))) OnBlockDisconnected", hash, height, t)
		// 	wg.forceUpdateChain()
		// 	if wg.processWalletBlockNotification() {
		// 	}
		// },
		// OnFilteredBlockDisconnected: func(height int32, header *wire.BlockHeader) {
		// 	if wg.Syncing.Load() {
		// 		return
		// 	}
		// 	D.Ln("(((NOTIFICATION))) OnFilteredBlockDisconnected", height, header)
		// 	wg.forceUpdateChain()
		// 	if wg.processWalletBlockNotification() {
		// 	}
		// },
		// OnRecvTx: func(transaction *util.Tx, details *btcjson.BlockDetails) {
		// 	if wg.cx.Syncing.Load() {
		// 		return
		// 	}
		// 	D.Ln("(((NOTIFICATION))) OnRecvTx", transaction, details)
		// 	wg.forceUpdateChain()
		// 	if wg.processWalletBlockNotification() {
		// 	}
		// },
		// OnRedeemingTx: func(transaction *util.Tx, details *btcjson.BlockDetails) {
		// 	if wg.cx.Syncing.Load() {
		// 		return
		// 	}
		// 	D.Ln("(((NOTIFICATION))) OnRedeemingTx", transaction, details)
		// 	wg.forceUpdateChain()
		// 	if wg.processWalletBlockNotification() {
		// 	}
		// },
		// OnRelevantTxAccepted: func(transaction []byte) {
		// 	if wg.cx.Syncing.Load() {
		// 		return
		// 	}
		// 	D.Ln("(((NOTIFICATION))) OnRelevantTxAccepted", transaction)
		// 	wg.forceUpdateChain()
		// 	if wg.processWalletBlockNotification() {
		// 	}
		// },
		// OnRescanFinished: func(hash *chainhash.Hash, height int32, blkTime time.Time) {
		// 	if wg.cx.Syncing.Load() {
		// 		return
		// 	}
		// 	D.Ln("(((NOTIFICATION))) OnRescanFinished", hash, height, blkTime)
		// 	wg.processChainBlockNotification(hash, height, blkTime)
		// 	// update best block height
		// 	// wg.processWalletBlockNotification()
		// 	// stop showing syncing indicator
		// 	if wg.processWalletBlockNotification() {
		// 	}
		// 	wg.cx.Syncing.Store(false)
		// 	wg.invalidate <- struct{}{}
		// },
		// OnRescanProgress: func(hash *chainhash.Hash, height int32, blkTime time.Time) {
		// 	D.Ln("(((NOTIFICATION))) OnRescanProgress", hash, height, blkTime)
		// 	// // update best block height
		// 	// // wg.processWalletBlockNotification()
		// 	// // set to show syncing indicator
		// 	// if wg.processWalletBlockNotification() {
		// 	// }
		// 	// wg.Syncing.Store(true)
		// 	wg.invalidate <- struct{}{}
		// },
		// OnTxAccepted: func(hash *chainhash.Hash, amount util.Amount) {
		// 	if wg.cx.Syncing.Load() {
		// 		return
		// 	}
		// 	D.Ln("(((NOTIFICATION))) OnTxAccepted")
		// 	D.Ln(hash, amount)
		// 	if wg.processWalletBlockNotification() {
		// 	}
		// },
		// OnTxAcceptedVerbose: func(txDetails *btcjson.TxRawResult) {
		// 	if wg.cx.Syncing.Load() {
		// 		return
		// 	}
		// 	D.Ln("(((NOTIFICATION))) OnTxAcceptedVerbose")
		// 	D.S(txDetails)
		// 	if wg.processWalletBlockNotification() {
		// 	}
		// },
		// OnPodConnected: func(connected bool) {
		// 	D.Ln("(((NOTIFICATION))) OnPodConnected", connected)
		// 	wg.forceUpdateChain()
		// 	if wg.processWalletBlockNotification() {
		// 	}
		// },
		// OnAccountBalance: func(account string, balance util.Amount, confirmed bool) {
		// 	D.Ln("OnAccountBalance")
		// 	// what does this actually do
		// 	D.Ln(account, balance, confirmed)
		// },
		// OnWalletLockState: func(locked bool) {
		// 	D.Ln("OnWalletLockState", locked)
		// 	// switch interface to unlock page
		// 	wg.forceUpdateChain()
		// 	if wg.processWalletBlockNotification() {
		// 	}
		// 	// TODO: lock when idle... how to get trigger for idleness in UI?
		// },
		// OnUnknownNotification: func(method string, params []json.RawMessage) {
		// 	D.Ln("(((NOTIFICATION))) OnUnknownNotification", method, params)
		// 	wg.forceUpdateChain()
		// 	if wg.processWalletBlockNotification() {
		// 	}
		// },
	}
	
}

func (wg *WalletGUI) chainClient() (e error) {
	D.Ln("starting up chain client")
	if wg.cx.Config.NodeOff.True() {
		W.Ln("node is disabled")
		return nil
	}
	if wg.ChainClient == nil { // || wg.ChainClient.Disconnected() {
		D.Ln(wg.cx.Config.RPCConnect.V())
		// wg.ChainMutex.Lock()
		// defer wg.ChainMutex.Unlock()
		// I.S(wg.certs)
		if wg.ChainClient, e = rpcclient.New(
			&rpcclient.ConnConfig{
				Host:                 wg.cx.Config.RPCConnect.V(),
				Endpoint:             "ws",
				User:                 wg.cx.Config.Username.V(),
				Pass:                 wg.cx.Config.Password.V(),
				TLS:                  wg.cx.Config.ClientTLS.True(),
				Certificates:         wg.certs,
				DisableAutoReconnect: false,
				DisableConnectOnNew:  false,
			}, wg.ChainNotifications(), wg.cx.KillAll,
		); E.Chk(e) {
			return
		}
	}
	if wg.ChainClient.Disconnected() {
		D.Ln("connecting chain client")
		if e = wg.ChainClient.Connect(1); E.Chk(e) {
			return
		}
	}
	if e = wg.ChainClient.NotifyBlocks(); !E.Chk(e) {
		D.Ln("subscribed to new blocks")
		// wg.WalletNotifications()
		wg.invalidate <- struct{}{}
	}
	return
}

func (wg *WalletGUI) walletClient() (e error) {
	D.Ln("connecting to wallet")
	if wg.cx.Config.WalletOff.True() {
		W.Ln("wallet is disabled")
		return nil
	}
	// walletRPC := (*wg.cx.Config.WalletRPCListeners)[0]
	// certs := wg.cx.Config.ReadCAFile()
	// I.Ln("config.tls", wg.cx.Config.ClientTLS.True())
	wg.WalletMutex.Lock()
	if wg.WalletClient, e = rpcclient.New(
		&rpcclient.ConnConfig{
			Host:                 wg.cx.Config.WalletServer.V(),
			Endpoint:             "ws",
			User:                 wg.cx.Config.Username.V(),
			Pass:                 wg.cx.Config.Password.V(),
			TLS:                  wg.cx.Config.ClientTLS.True(),
			Certificates:         wg.certs,
			DisableAutoReconnect: false,
			DisableConnectOnNew:  false,
		}, wg.WalletNotifications(), wg.cx.KillAll,
	); E.Chk(e) {
		wg.WalletMutex.Unlock()
		return
	}
	wg.WalletMutex.Unlock()
	// if e = wg.WalletClient.Connect(1); E.Chk(e) {
	// 	return
	// }
	if e = wg.WalletClient.NotifyNewTransactions(true); !E.Chk(e) {
		D.Ln("subscribed to new transactions")
	} else {
		// return
	}
	// if e = wg.WalletClient.NotifyBlocks(); E.Chk(e) {
	// 	// return
	// } else {
	// 	D.Ln("subscribed to wallet client notify blocks")
	// }
	D.Ln("wallet connected")
	return
}
