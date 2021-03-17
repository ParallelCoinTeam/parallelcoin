package gui

import (
	"encoding/json"
	"github.com/p9c/pod/pkg/blockchain/wire"
	"github.com/p9c/pod/pkg/pod"
	"io/ioutil"
	"path/filepath"
	"time"
	
	"github.com/p9c/pod/pkg/blockchain/chainhash"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	"github.com/p9c/pod/pkg/rpc/rpcclient"
	"github.com/p9c/pod/pkg/util"
)

func (wg *WalletGUI) WalletAndClientRunning() bool {
	running := wg.wallet.Running() && wg.WalletClient != nil && !wg.WalletClient.Disconnected()
	// dbg.Ln("wallet and wallet rpc client are running?", running)
	return running
}

func (wg *WalletGUI) Tickers() {
	first := true
	dbg.Ln("updating best block")
	var e error
	var height int32
	var h *chainhash.Hash
	if h, height, e = wg.ChainClient.GetBestBlock(); err.Chk(e) {
		// interrupt.Request()
		return
	}
	dbg.Ln(h, height)
	wg.State.SetBestBlockHeight(height)
	wg.State.SetBestBlockHash(h)
	
	go func() {
		var e error
		seconds := time.Tick(time.Second * 3)
		fiveSeconds := time.Tick(time.Second * 5)
	totalOut:
		for {
		preconnect:
			for {
				select {
				case <-seconds:
					dbg.Ln("---------------------- ready", wg.ready.Load())
					dbg.Ln("---------------------- WalletAndClientRunning", wg.WalletAndClientRunning())
					dbg.Ln("---------------------- stateLoaded", wg.stateLoaded.Load())
					dbg.Ln("preconnect loop")
					if wg.ChainClient != nil {
						wg.ChainClient.Disconnect()
						wg.ChainClient.Shutdown()
						wg.ChainClient = nil
					}
					if wg.WalletClient != nil {
						wg.WalletClient.Disconnect()
						wg.WalletClient.Shutdown()
						wg.WalletClient = nil
					}
					if !wg.node.Running() {
						break
					}
					break preconnect
				case <-fiveSeconds:
					continue
				case <-wg.quit.Wait():
					break totalOut
				}
			}
		out:
			for {
				select {
				case <-seconds:
					dbg.Ln("---------------------- ready", wg.ready.Load())
					dbg.Ln("---------------------- WalletAndClientRunning", wg.WalletAndClientRunning())
					dbg.Ln("---------------------- stateLoaded", wg.stateLoaded.Load())
					wg.node.Start()
					if e = wg.writeWalletCookie(); err.Chk(e) {
					}
					wg.wallet.Start()
					dbg.Ln("connecting to chain")
					if e = wg.chainClient(); err.Chk(e) {
						break
					}
					if wg.wallet.Running() { // && wg.WalletClient == nil {
						dbg.Ln("connecting to wallet")
						if e = wg.walletClient(); err.Chk(e) {
							break
						}
					}
					if !wg.node.Running() {
						dbg.Ln("breaking out node not running")
						break out
					}
					if wg.ChainClient == nil {
						dbg.Ln("breaking out chainclient is nil")
						break out
					}
					// if  wg.WalletClient == nil {
					// 	dbg.Ln("breaking out walletclient is nil")
					// 	break out
					// }
					if wg.ChainClient.Disconnected() {
						dbg.Ln("breaking out chainclient disconnected")
						break out
					}
					// if wg.WalletClient.Disconnected() {
					// 	dbg.Ln("breaking out walletclient disconnected")
					// 	break out
					// }
					// var e error
					if first {
						wg.updateChainBlock()
						wg.invalidate <- struct{}{}
					}
					
					if wg.WalletAndClientRunning() {
						if first {
							wg.processWalletBlockNotification()
						}
						// if wg.stateLoaded.Load() { // || wg.currentReceiveGetNew.Load() {
						// 	wg.ReceiveAddressbook = func(gtx l.Context) l.Dimensions {
						// 		var widgets []l.Widget
						// 		widgets = append(widgets, wg.ReceivePage.GetAddressbookHistoryCards("DocBg")...)
						// 		le := func(gtx l.Context, index int) l.Dimensions {
						// 			return widgets[index](gtx)
						// 		}
						// 		return wg.Flex().Rigid(
						// 			wg.lists["receiveAddresses"].Length(len(widgets)).Vertical().
						// 				ListElement(le).Fn,
						// 		).Fn(gtx)
						// 	}
						// }
						if wg.stateLoaded.Load() && !wg.State.IsReceivingAddress() { // || wg.currentReceiveGetNew.Load() {
							wg.GetNewReceivingAddress()
							if wg.currentReceiveQRCode == nil || wg.currentReceiveRegenerate.Load() { // || wg.currentReceiveGetNew.Load() {
								wg.GetNewReceivingQRCode(wg.ReceivePage.urn)
							}
						}
					}
					wg.invalidate <- struct{}{}
					first = false
				case <-fiveSeconds:
				case <-wg.quit.Wait():
					break totalOut
				}
			}
		}
	}()
}

func (wg *WalletGUI) updateThingies() (e error) {
	// update the configuration
	var b []byte
	if b, e = ioutil.ReadFile(*wg.cx.Config.ConfigFile); !err.Chk(e) {
		if e = json.Unmarshal(b, wg.cx.Config); !err.Chk(e) {
			return
		}
	}
	return
}
func (wg *WalletGUI) updateChainBlock() {
	dbg.Ln("processChainBlockNotification")
	var e error
	if wg.ChainClient != nil && wg.ChainClient.Disconnected() || wg.ChainClient.Disconnected() {
		dbg.Ln("connecting ChainClient")
		if e = wg.chainClient(); err.Chk(e) {
			return
		}
	}
	var h *chainhash.Hash
	var height int32
	dbg.Ln("updating best block")
	if h, height, e = wg.ChainClient.GetBestBlock(); err.Chk(e) {
		// interrupt.Request()
		return
	}
	dbg.Ln(h, height)
	wg.State.SetBestBlockHeight(height)
	wg.State.SetBestBlockHash(h)
}

func (wg *WalletGUI) processChainBlockNotification(hash *chainhash.Hash, height int32, t time.Time) {
	dbg.Ln("processChainBlockNotification")
	wg.State.SetBestBlockHeight(height)
	wg.State.SetBestBlockHash(hash)
	if wg.WalletAndClientRunning() {
		wg.processWalletBlockNotification()
	}
}

func (wg *WalletGUI) processWalletBlockNotification() bool {
	dbg.Ln("processWalletBlockNotification")
	if !wg.WalletAndClientRunning() {
		dbg.Ln("wallet and client not running")
		return false
	}
	// check account balance
	var unconfirmed util.Amount
	var e error
	if unconfirmed, e = wg.WalletClient.GetUnconfirmedBalance("default"); err.Chk(e) {
		return false
	}
	wg.State.SetBalanceUnconfirmed(unconfirmed.ToDUO())
	var confirmed util.Amount
	if confirmed, e = wg.WalletClient.GetBalance("default"); err.Chk(e) {
		return false
	}
	wg.State.SetBalance(confirmed.ToDUO())
	var atr []btcjson.ListTransactionsResult
	// str := wg.State.allTxs.Load()
	if atr, e = wg.WalletClient.ListTransactionsCountFrom("default", 2<<16, /*len(str)*/0); err.Chk(e) {
		return false
	}
	// dbg.Ln(len(atr))
	// wg.State.SetAllTxs(append(str, atr...))
	wg.State.SetAllTxs(atr)
	wg.txMx.Lock()
	wg.txHistoryList = wg.State.filteredTxs.Load()
	atrl := 10
	if len(atr) < atrl {
		atrl = len(atr)
	}
	wg.txRecentList = atr[:atrl]
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
	if tip, height, e = wg.ChainClient.GetBestBlock(); err.Chk(e) {
		return
	}
	var block *wire.MsgBlock
	if block, e = wg.ChainClient.GetBlock(tip); err.Chk(e) {
	}
	t := block.Header.Timestamp
	wg.processChainBlockNotification(tip, height, t)
}

func (wg *WalletGUI) ChainNotifications() *rpcclient.NotificationHandlers {
	return &rpcclient.NotificationHandlers{
		OnClientConnected: func() {
			// go func() {
			dbg.Ln("(((NOTIFICATION))) CHAIN CLIENT CONNECTED!")
			wg.forceUpdateChain()
			wg.processWalletBlockNotification()
			wg.invalidate <- struct{}{}
		},
		OnBlockConnected: func(hash *chainhash.Hash, height int32, t time.Time) {
			dbg.Ln("(((NOTIFICATION))) chain OnBlockConnected", hash, height, t)
			wg.processChainBlockNotification(hash, height, t)
			// wg.processWalletBlockNotification()
			// todo: send system notification of new block, set configuration to disable also
			wg.invalidate <- struct{}{}
			
		},
		OnFilteredBlockConnected: func(height int32, header *wire.BlockHeader, txs []*util.Tx) {
			dbg.Ln(
				"(((NOTIFICATION))) wallet OnFilteredBlockConnected hash", header.BlockHash(), "POW hash:",
				header.BlockHashWithAlgos(height), "height", height,
			)
			// dbg.S(txs)
			nbh := header.BlockHash()
			wg.processChainBlockNotification(&nbh, height, header.Timestamp)
			if wg.processWalletBlockNotification() {
			}
			filename := filepath.Join(wg.cx.DataDir, "state.json")
			if e := wg.State.Save(filename, wg.cx.Config.WalletPass); err.Chk(e) {
			}
			wg.invalidate <- struct{}{}
		},
		// OnBlockDisconnected: func(hash *chainhash.Hash, height int32, t time.Time) {
		// 	dbg.Ln("(((NOTIFICATION))) OnBlockDisconnected", hash, height, t)
		// 	wg.forceUpdateChain()
		// 	if wg.processWalletBlockNotification() {
		// 	}
		// },
		// OnFilteredBlockDisconnected: func(height int32, header *wire.BlockHeader) {
		// 	dbg.Ln("(((NOTIFICATION))) OnFilteredBlockDisconnected", height, header)
		// 	wg.forceUpdateChain()
		// 	if wg.processWalletBlockNotification() {
		// 	}
		// },
		// OnRecvTx: func(transaction *util.Tx, details *btcjson.BlockDetails) {
		// 	dbg.Ln("(((NOTIFICATION))) OnRecvTx", transaction, details)
		// 	wg.forceUpdateChain()
		// 	if wg.processWalletBlockNotification() {
		// 	}
		// },
		// OnRedeemingTx: func(transaction *util.Tx, details *btcjson.BlockDetails) {
		// 	dbg.Ln("(((NOTIFICATION))) OnRedeemingTx", transaction, details)
		// 	wg.forceUpdateChain()
		// 	if wg.processWalletBlockNotification() {
		// 	}
		// },
		// OnRelevantTxAccepted: func(transaction []byte) {
		// 	dbg.Ln("(((NOTIFICATION))) OnRelevantTxAccepted", transaction)
		// 	wg.forceUpdateChain()
		// 	if wg.processWalletBlockNotification() {
		// 	}
		// },
		// OnRescanFinished: func(hash *chainhash.Hash, height int32, blkTime time.Time) {
		// 	dbg.Ln("(((NOTIFICATION))) OnRescanFinished", hash, height, blkTime)
		// 	wg.processChainBlockNotification(hash, height, blkTime)
		// 	// update best block height
		// 	// wg.processWalletBlockNotification()
		// 	// stop showing syncing indicator
		// 	if wg.processWalletBlockNotification() {
		// 	}
		// 	wg.Syncing.Store(false)
		// 	wg.invalidate <- struct{}{}
		// },
		// OnRescanProgress: func(hash *chainhash.Hash, height int32, blkTime time.Time) {
		// 	dbg.Ln("(((NOTIFICATION))) OnRescanProgress", hash, height, blkTime)
		// 	// update best block height
		// 	// wg.processWalletBlockNotification()
		// 	// set to show syncing indicator
		// 	if wg.processWalletBlockNotification() {
		// 	}
		// 	wg.Syncing.Store(true)
		// 	wg.invalidate <- struct{}{}
		// },
		// OnTxAccepted: func(hash *chainhash.Hash, amount util.Amount) {
		// 	dbg.Ln("(((NOTIFICATION))) OnTxAccepted")
		// 	dbg.Ln(hash, amount)
		// 	if wg.processWalletBlockNotification() {
		// 	}
		// },
		// OnTxAcceptedVerbose: func(txDetails *btcjson.TxRawResult) {
		// 	dbg.Ln("(((NOTIFICATION))) OnTxAcceptedVerbose")
		// 	dbg.S(txDetails)
		// 	if wg.processWalletBlockNotification() {
		// 	}
		// },
		// OnPodConnected: func(connected bool) {
		// 	dbg.Ln("(((NOTIFICATION))) OnPodConnected", connected)
		// 	wg.forceUpdateChain()
		// 	if wg.processWalletBlockNotification() {
		// 	}
		// },
		// OnAccountBalance: func(account string, balance util.Amount, confirmed bool) {
		// 	dbg.Ln("OnAccountBalance")
		// 	// what does this actually do
		// 	dbg.Ln(account, balance, confirmed)
		// },
		// OnWalletLockState: func(locked bool) {
		// 	dbg.Ln("OnWalletLockState", locked)
		// 	// switch interface to unlock page
		// 	wg.forceUpdateChain()
		// 	if wg.processWalletBlockNotification() {
		// 	}
		// 	// TODO: lock when idle... how to get trigger for idleness in UI?
		// },
		// OnUnknownNotification: func(method string, params []json.RawMessage) {
		// 	dbg.Ln("(((NOTIFICATION))) OnUnknownNotification", method, params)
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
	return &rpcclient.NotificationHandlers{
		OnClientConnected: func() {
			dbg.Ln("(((NOTIFICATION))) wallet client connected, running initial processes")
			for !wg.processWalletBlockNotification() {
				time.Sleep(time.Second)
				dbg.Ln("(((NOTIFICATION))) retry attempting to update wallet transactions")
			}
			filename := filepath.Join(wg.cx.DataDir, "state.json")
			if e := wg.State.Save(filename, wg.cx.Config.WalletPass); err.Chk(e) {
			}
			wg.invalidate <- struct{}{}
		},
		OnBlockConnected: func(hash *chainhash.Hash, height int32, t time.Time) {
			dbg.Ln("(((NOTIFICATION))) wallet OnBlockConnected", hash, height, t)
			wg.processWalletBlockNotification()
			filename := filepath.Join(wg.cx.DataDir, "state.json")
			if e := wg.State.Save(filename, wg.cx.Config.WalletPass); err.Chk(e) {
			}
			wg.invalidate <- struct{}{}
		},
		OnFilteredBlockConnected: func(height int32, header *wire.BlockHeader, txs []*util.Tx) {
			dbg.Ln(
				"(((NOTIFICATION))) wallet OnFilteredBlockConnected hash", header.BlockHash(), "POW hash:",
				header.BlockHashWithAlgos(height), "height", height,
			)
			// dbg.S(txs)
			nbh := header.BlockHash()
			wg.processChainBlockNotification(&nbh, height, header.Timestamp)
			if wg.processWalletBlockNotification() {
			}
			filename := filepath.Join(wg.cx.DataDir, "state.json")
			if e := wg.State.Save(filename, wg.cx.Config.WalletPass); err.Chk(e) {
			}
			wg.invalidate <- struct{}{}
		},
		OnBlockDisconnected: func(hash *chainhash.Hash, height int32, t time.Time) {
			dbg.Ln("(((NOTIFICATION))) OnBlockDisconnected", hash, height, t)
			wg.forceUpdateChain()
			if wg.processWalletBlockNotification() {
			}
		},
		OnFilteredBlockDisconnected: func(height int32, header *wire.BlockHeader) {
			dbg.Ln("(((NOTIFICATION))) OnFilteredBlockDisconnected", height, header)
			wg.forceUpdateChain()
			if wg.processWalletBlockNotification() {
			}
		},
		OnRecvTx: func(transaction *util.Tx, details *btcjson.BlockDetails) {
			dbg.Ln("(((NOTIFICATION))) OnRecvTx", transaction, details)
			wg.forceUpdateChain()
			if wg.processWalletBlockNotification() {
			}
		},
		OnRedeemingTx: func(transaction *util.Tx, details *btcjson.BlockDetails) {
			dbg.Ln("(((NOTIFICATION))) OnRedeemingTx", transaction, details)
			wg.forceUpdateChain()
			if wg.processWalletBlockNotification() {
			}
		},
		OnRelevantTxAccepted: func(transaction []byte) {
			dbg.Ln("(((NOTIFICATION))) OnRelevantTxAccepted", transaction)
			wg.forceUpdateChain()
			if wg.processWalletBlockNotification() {
			}
		},
		OnRescanFinished: func(hash *chainhash.Hash, height int32, blkTime time.Time) {
			dbg.Ln("(((NOTIFICATION))) OnRescanFinished", hash, height, blkTime)
			wg.processChainBlockNotification(hash, height, blkTime)
			// update best block height
			// wg.processWalletBlockNotification()
			// stop showing syncing indicator
			if wg.processWalletBlockNotification() {
			}
			wg.Syncing.Store(false)
			wg.invalidate <- struct{}{}
		},
		OnRescanProgress: func(hash *chainhash.Hash, height int32, blkTime time.Time) {
			dbg.Ln("(((NOTIFICATION))) OnRescanProgress", hash, height, blkTime)
			// update best block height
			// wg.processWalletBlockNotification()
			// set to show syncing indicator
			if wg.processWalletBlockNotification() {
			}
			wg.Syncing.Store(true)
			wg.invalidate <- struct{}{}
		},
		OnTxAccepted: func(hash *chainhash.Hash, amount util.Amount) {
			dbg.Ln("(((NOTIFICATION))) OnTxAccepted")
			dbg.Ln(hash, amount)
			if wg.processWalletBlockNotification() {
			}
		},
		OnTxAcceptedVerbose: func(txDetails *btcjson.TxRawResult) {
			dbg.Ln("(((NOTIFICATION))) OnTxAcceptedVerbose")
			dbg.S(txDetails)
			if wg.processWalletBlockNotification() {
			}
		},
		OnPodConnected: func(connected bool) {
			dbg.Ln("(((NOTIFICATION))) OnPodConnected", connected)
			wg.forceUpdateChain()
			if wg.processWalletBlockNotification() {
			}
		},
		OnAccountBalance: func(account string, balance util.Amount, confirmed bool) {
			dbg.Ln("OnAccountBalance")
			// what does this actually do
			dbg.Ln(account, balance, confirmed)
		},
		OnWalletLockState: func(locked bool) {
			dbg.Ln("OnWalletLockState", locked)
			// switch interface to unlock page
			wg.forceUpdateChain()
			if wg.processWalletBlockNotification() {
			}
			// TODO: lock when idle... how to get trigger for idleness in UI?
		},
		OnUnknownNotification: func(method string, params []json.RawMessage) {
			dbg.Ln("(((NOTIFICATION))) OnUnknownNotification", method, params)
			wg.forceUpdateChain()
			if wg.processWalletBlockNotification() {
			}
		},
	}
	
}

func (wg *WalletGUI) chainClient() (e error) {
	dbg.Ln("starting up chain client")
	if *wg.cx.Config.NodeOff {
		wrn.Ln("node is disabled")
		return nil
	}
	
	if wg.ChainClient == nil { // || wg.ChainClient.Disconnected() {
		certs := pod.ReadCAFile(wg.cx.Config)
		dbg.Ln(*wg.cx.Config.RPCConnect)
		// wg.ChainMutex.Lock()
		// defer wg.ChainMutex.Unlock()
		if wg.ChainClient, e = rpcclient.New(
			&rpcclient.ConnConfig{
				Host:                 *wg.cx.Config.RPCConnect,
				Endpoint:             "ws",
				User:                 *wg.cx.Config.Username,
				Pass:                 *wg.cx.Config.Password,
				TLS:                  *wg.cx.Config.TLS,
				Certificates:         certs,
				DisableAutoReconnect: false,
				DisableConnectOnNew:  false,
			}, wg.ChainNotifications(), wg.cx.KillAll,
		); err.Chk(e) {
			return
		}
	}
	if wg.ChainClient.Disconnected() {
		dbg.Ln("connecting chain client")
		if e = wg.ChainClient.Connect(1); err.Chk(e) {
			return
		}
	}
	if e = wg.ChainClient.NotifyBlocks(); !err.Chk(e) {
		dbg.Ln("subscribed to new blocks")
		// wg.WalletNotifications()
		wg.invalidate <- struct{}{}
	}
	return
}

func (wg *WalletGUI) walletClient() (e error) {
	dbg.Ln("connecting to wallet")
	if *wg.cx.Config.WalletOff {
		wrn.Ln("wallet is disabled")
		return nil
	}
	// walletRPC := (*wg.cx.Config.WalletRPCListeners)[0]
	certs := pod.ReadCAFile(wg.cx.Config)
	inf.Ln("config.tls", *wg.cx.Config.TLS)
	wg.WalletMutex.Lock()
	if wg.WalletClient, e = rpcclient.New(
		&rpcclient.ConnConfig{
			Host:                 *wg.cx.Config.WalletServer,
			Endpoint:             "ws",
			User:                 *wg.cx.Config.Username,
			Pass:                 *wg.cx.Config.Password,
			TLS:                  *wg.cx.Config.TLS,
			Certificates:         certs,
			DisableAutoReconnect: false,
			DisableConnectOnNew:  false,
		}, wg.WalletNotifications(), wg.cx.KillAll,
	); err.Chk(e) {
		wg.WalletMutex.Unlock()
		return
	}
	wg.WalletMutex.Unlock()
	// if e = wg.WalletClient.Connect(1); err.Chk(e) {
	// 	return
	// }
	if e = wg.WalletClient.NotifyNewTransactions(true); !err.Chk(e) {
		dbg.Ln("subscribed to new transactions")
	} else {
		// return
	}
	if e = wg.WalletClient.NotifyBlocks(); err.Chk(e) {
		// return
	} else {
		dbg.Ln("subscribed to wallet client notify blocks")
	}
	dbg.Ln("wallet connected")
	return
}
