package gui

import (
	"encoding/json"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/pod"
	"io/ioutil"
	"path/filepath"
	"time"
	
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	rpcclient "github.com/p9c/pod/pkg/rpc/client"
	"github.com/p9c/pod/pkg/util"
)

func (wg *WalletGUI) WalletAndClientRunning() bool {
	running := wg.wallet.Running() && wg.WalletClient != nil && !wg.WalletClient.Disconnected()
	// Debug("wallet and wallet rpc client are running?", running)
	return running
}

func (wg *WalletGUI) Tickers() {
	first := true
	Debug("updating best block")
	var err error
	var height int32
	var h *chainhash.Hash
	if h, height, err = wg.ChainClient.GetBestBlock(); Check(err) {
		// interrupt.Request()
		return
	}
	Debug(h, height)
	wg.State.SetBestBlockHeight(height)
	wg.State.SetBestBlockHash(h)
	
	go func() {
		var err error
		seconds := time.Tick(time.Second * 3)
		fiveSeconds := time.Tick(time.Second * 5)
	totalOut:
		for {
		preconnect:
			for {
				select {
				case <-seconds:
					Debug("---------------------- ready", wg.ready.Load())
					Debug("---------------------- WalletAndClientRunning", wg.WalletAndClientRunning())
					Debug("---------------------- stateLoaded", wg.stateLoaded.Load())
					Debug("preconnect loop")
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
					Debug("---------------------- ready", wg.ready.Load())
					Debug("---------------------- WalletAndClientRunning", wg.WalletAndClientRunning())
					Debug("---------------------- stateLoaded", wg.stateLoaded.Load())
					wg.node.Start()
					if err = wg.writeWalletCookie(); Check(err) {
					}
					wg.wallet.Start()
					Debug("connecting to chain")
					if err = wg.chainClient(); err != nil {
						break
					}
					if wg.wallet.Running() { // && wg.WalletClient == nil {
						Debug("connecting to wallet")
						if err = wg.walletClient(); Check(err) {
							break
						}
					}
					if !wg.node.Running() {
						Debug("breaking out node not running")
						break out
					}
					if wg.ChainClient == nil {
						Debug("breaking out chainclient is nil")
						break out
					}
					// if  wg.WalletClient == nil {
					// 	Debug("breaking out walletclient is nil")
					// 	break out
					// }
					if wg.ChainClient.Disconnected() {
						Debug("breaking out chainclient disconnected")
						break out
					}
					// if wg.WalletClient.Disconnected() {
					// 	Debug("breaking out walletclient disconnected")
					// 	break out
					// }
					// var err error
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

func (wg *WalletGUI) updateThingies() (err error) {
	// update the configuration
	var b []byte
	if b, err = ioutil.ReadFile(*wg.cx.Config.ConfigFile); !Check(err) {
		if err = json.Unmarshal(b, wg.cx.Config); !Check(err) {
			return
		}
	}
	return
}
func (wg *WalletGUI) updateChainBlock() {
	Debug("processChainBlockNotification")
	var err error
	if wg.ChainClient != nil && wg.ChainClient.Disconnected() || wg.ChainClient.Disconnected() {
		Debug("connecting ChainClient")
		if err = wg.chainClient(); Check(err) {
			return
		}
	}
	var h *chainhash.Hash
	var height int32
	Debug("updating best block")
	if h, height, err = wg.ChainClient.GetBestBlock(); Check(err) {
		// interrupt.Request()
		return
	}
	Debug(h, height)
	wg.State.SetBestBlockHeight(height)
	wg.State.SetBestBlockHash(h)
}

func (wg *WalletGUI) processChainBlockNotification(hash *chainhash.Hash, height int32, t time.Time) {
	Debug("processChainBlockNotification")
	wg.State.SetBestBlockHeight(height)
	wg.State.SetBestBlockHash(hash)
	if wg.WalletAndClientRunning() {
		wg.processWalletBlockNotification()
	}
}

func (wg *WalletGUI) processWalletBlockNotification() bool {
	Debug("processWalletBlockNotification")
	if !wg.WalletAndClientRunning() {
		Debug("wallet and client not running")
		return false
	}
	// check account balance
	var unconfirmed util.Amount
	var err error
	if unconfirmed, err = wg.WalletClient.GetUnconfirmedBalance("default"); Check(err) {
		return false
	}
	wg.State.SetBalanceUnconfirmed(unconfirmed.ToDUO())
	var confirmed util.Amount
	if confirmed, err = wg.WalletClient.GetBalance("default"); Check(err) {
		return false
	}
	wg.State.SetBalance(confirmed.ToDUO())
	var atr []btcjson.ListTransactionsResult
	// TODO: for some reason this function returns half as many as requested
	if atr, err = wg.WalletClient.ListTransactionsCountFrom("default", 2<<16, 0); Check(err) {
		return false
	}
	// Debug(len(atr))
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
	var err error
	var height int32
	var tip *chainhash.Hash
	if tip, height, err = wg.ChainClient.GetBestBlock(); Check(err) {
		return
	}
	var block *wire.MsgBlock
	if block, err = wg.ChainClient.GetBlock(tip); Check(err) {
	}
	t := block.Header.Timestamp
	wg.processChainBlockNotification(tip, height, t)
}

func (wg *WalletGUI) ChainNotifications() *rpcclient.NotificationHandlers {
	return &rpcclient.NotificationHandlers{
		OnClientConnected: func() {
			// go func() {
			Debug("(((NOTIFICATION))) CHAIN CLIENT CONNECTED!")
			wg.forceUpdateChain()
			wg.processWalletBlockNotification()
			wg.invalidate <- struct{}{}
		},
		OnBlockConnected: func(hash *chainhash.Hash, height int32, t time.Time) {
			Debug("(((NOTIFICATION))) chain OnBlockConnected", hash, height, t)
			wg.processChainBlockNotification(hash, height, t)
			wg.processWalletBlockNotification()
			// todo: send system notification of new block, set configuration to disable also
			wg.invalidate <- struct{}{}
			
		},
		OnFilteredBlockConnected: func(height int32, header *wire.BlockHeader, txs []*util.Tx) {
			Debug(
				"(((NOTIFICATION))) wallet OnFilteredBlockConnected hash", header.BlockHash(), "POW hash:",
				header.BlockHashWithAlgos(height), "height", height,
			)
			// Debugs(txs)
			nbh := header.BlockHash()
			wg.processChainBlockNotification(&nbh, height, header.Timestamp)
			if wg.processWalletBlockNotification() {
			}
			filename := filepath.Join(wg.cx.DataDir, "state.json")
			if err := wg.State.Save(filename, wg.cx.Config.WalletPass); Check(err) {
			}
			wg.invalidate <- struct{}{}
		},
		OnBlockDisconnected: func(hash *chainhash.Hash, height int32, t time.Time) {
			Debug("(((NOTIFICATION))) OnBlockDisconnected", hash, height, t)
			wg.forceUpdateChain()
			if wg.processWalletBlockNotification() {
			}
		},
		OnFilteredBlockDisconnected: func(height int32, header *wire.BlockHeader) {
			Debug("(((NOTIFICATION))) OnFilteredBlockDisconnected", height, header)
			wg.forceUpdateChain()
			if wg.processWalletBlockNotification() {
			}
		},
		OnRecvTx: func(transaction *util.Tx, details *btcjson.BlockDetails) {
			Debug("(((NOTIFICATION))) OnRecvTx", transaction, details)
			wg.forceUpdateChain()
			if wg.processWalletBlockNotification() {
			}
		},
		OnRedeemingTx: func(transaction *util.Tx, details *btcjson.BlockDetails) {
			Debug("(((NOTIFICATION))) OnRedeemingTx", transaction, details)
			wg.forceUpdateChain()
			if wg.processWalletBlockNotification() {
			}
		},
		OnRelevantTxAccepted: func(transaction []byte) {
			Debug("(((NOTIFICATION))) OnRelevantTxAccepted", transaction)
			wg.forceUpdateChain()
			if wg.processWalletBlockNotification() {
			}
		},
		OnRescanFinished: func(hash *chainhash.Hash, height int32, blkTime time.Time) {
			Debug("(((NOTIFICATION))) OnRescanFinished", hash, height, blkTime)
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
			Debug("(((NOTIFICATION))) OnRescanProgress", hash, height, blkTime)
			// update best block height
			// wg.processWalletBlockNotification()
			// set to show syncing indicator
			if wg.processWalletBlockNotification() {
			}
			wg.Syncing.Store(true)
			wg.invalidate <- struct{}{}
		},
		OnTxAccepted: func(hash *chainhash.Hash, amount util.Amount) {
			Debug("(((NOTIFICATION))) OnTxAccepted")
			Debug(hash, amount)
			if wg.processWalletBlockNotification() {
			}
		},
		OnTxAcceptedVerbose: func(txDetails *btcjson.TxRawResult) {
			Debug("(((NOTIFICATION))) OnTxAcceptedVerbose")
			Debugs(txDetails)
			if wg.processWalletBlockNotification() {
			}
		},
		OnPodConnected: func(connected bool) {
			Debug("(((NOTIFICATION))) OnPodConnected", connected)
			wg.forceUpdateChain()
			if wg.processWalletBlockNotification() {
			}
		},
		OnAccountBalance: func(account string, balance util.Amount, confirmed bool) {
			Debug("OnAccountBalance")
			// what does this actually do
			Debug(account, balance, confirmed)
		},
		OnWalletLockState: func(locked bool) {
			Debug("OnWalletLockState", locked)
			// switch interface to unlock page
			wg.forceUpdateChain()
			if wg.processWalletBlockNotification() {
			}
			// TODO: lock when idle... how to get trigger for idleness in UI?
		},
		OnUnknownNotification: func(method string, params []json.RawMessage) {
			Debug("(((NOTIFICATION))) OnUnknownNotification", method, params)
			wg.forceUpdateChain()
			if wg.processWalletBlockNotification() {
			}
		},
	}
	
}

func (wg *WalletGUI) WalletNotifications() *rpcclient.NotificationHandlers {
	// if !wg.wallet.Running() || wg.WalletClient == nil || wg.WalletClient.Disconnected() {
	// 	return nil
	// }
	return &rpcclient.NotificationHandlers{
		OnClientConnected: func() {
			Debug("(((NOTIFICATION))) wallet client connected, running initial processes")
			for !wg.processWalletBlockNotification() {
				time.Sleep(time.Second)
				Debug("(((NOTIFICATION))) retry attempting to update wallet transactions")
			}
			filename := filepath.Join(wg.cx.DataDir, "state.json")
			if err := wg.State.Save(filename, wg.cx.Config.WalletPass); Check(err) {
			}
			wg.invalidate <- struct{}{}
		},
		OnBlockConnected: func(hash *chainhash.Hash, height int32, t time.Time) {
			Debug("(((NOTIFICATION))) wallet OnBlockConnected", hash, height, t)
			wg.processWalletBlockNotification()
			filename := filepath.Join(wg.cx.DataDir, "state.json")
			if err := wg.State.Save(filename, wg.cx.Config.WalletPass); Check(err) {
			}
			wg.invalidate <- struct{}{}
		},
		OnFilteredBlockConnected: func(height int32, header *wire.BlockHeader, txs []*util.Tx) {
			Debug(
				"(((NOTIFICATION))) wallet OnFilteredBlockConnected hash", header.BlockHash(), "POW hash:",
				header.BlockHashWithAlgos(height), "height", height,
			)
			// Debugs(txs)
			nbh := header.BlockHash()
			wg.processChainBlockNotification(&nbh, height, header.Timestamp)
			if wg.processWalletBlockNotification() {
			}
			filename := filepath.Join(wg.cx.DataDir, "state.json")
			if err := wg.State.Save(filename, wg.cx.Config.WalletPass); Check(err) {
			}
			wg.invalidate <- struct{}{}
		},
		OnBlockDisconnected: func(hash *chainhash.Hash, height int32, t time.Time) {
			Debug("(((NOTIFICATION))) OnBlockDisconnected", hash, height, t)
			wg.forceUpdateChain()
			if wg.processWalletBlockNotification() {
			}
		},
		OnFilteredBlockDisconnected: func(height int32, header *wire.BlockHeader) {
			Debug("(((NOTIFICATION))) OnFilteredBlockDisconnected", height, header)
			wg.forceUpdateChain()
			if wg.processWalletBlockNotification() {
			}
		},
		OnRecvTx: func(transaction *util.Tx, details *btcjson.BlockDetails) {
			Debug("(((NOTIFICATION))) OnRecvTx", transaction, details)
			wg.forceUpdateChain()
			if wg.processWalletBlockNotification() {
			}
		},
		OnRedeemingTx: func(transaction *util.Tx, details *btcjson.BlockDetails) {
			Debug("(((NOTIFICATION))) OnRedeemingTx", transaction, details)
			wg.forceUpdateChain()
			if wg.processWalletBlockNotification() {
			}
		},
		OnRelevantTxAccepted: func(transaction []byte) {
			Debug("(((NOTIFICATION))) OnRelevantTxAccepted", transaction)
			wg.forceUpdateChain()
			if wg.processWalletBlockNotification() {
			}
		},
		OnRescanFinished: func(hash *chainhash.Hash, height int32, blkTime time.Time) {
			Debug("(((NOTIFICATION))) OnRescanFinished", hash, height, blkTime)
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
			Debug("(((NOTIFICATION))) OnRescanProgress", hash, height, blkTime)
			// update best block height
			// wg.processWalletBlockNotification()
			// set to show syncing indicator
			if wg.processWalletBlockNotification() {
			}
			wg.Syncing.Store(true)
			wg.invalidate <- struct{}{}
		},
		OnTxAccepted: func(hash *chainhash.Hash, amount util.Amount) {
			Debug("(((NOTIFICATION))) OnTxAccepted")
			Debug(hash, amount)
			if wg.processWalletBlockNotification() {
			}
		},
		OnTxAcceptedVerbose: func(txDetails *btcjson.TxRawResult) {
			Debug("(((NOTIFICATION))) OnTxAcceptedVerbose")
			Debugs(txDetails)
			if wg.processWalletBlockNotification() {
			}
		},
		OnPodConnected: func(connected bool) {
			Debug("(((NOTIFICATION))) OnPodConnected", connected)
			wg.forceUpdateChain()
			if wg.processWalletBlockNotification() {
			}
		},
		OnAccountBalance: func(account string, balance util.Amount, confirmed bool) {
			Debug("OnAccountBalance")
			// what does this actually do
			Debug(account, balance, confirmed)
		},
		OnWalletLockState: func(locked bool) {
			Debug("OnWalletLockState", locked)
			// switch interface to unlock page
			wg.forceUpdateChain()
			if wg.processWalletBlockNotification() {
			}
			// TODO: lock when idle... how to get trigger for idleness in UI?
		},
		OnUnknownNotification: func(method string, params []json.RawMessage) {
			Debug("(((NOTIFICATION))) OnUnknownNotification", method, params)
			wg.forceUpdateChain()
			if wg.processWalletBlockNotification() {
			}
		},
	}
	
}

func (wg *WalletGUI) chainClient() (err error) {
	Debug("starting up chain client")
	if *wg.cx.Config.NodeOff {
		Warn("node is disabled")
		return nil
	}
	
	if wg.ChainClient == nil { // || wg.ChainClient.Disconnected() {
		certs := pod.ReadCAFile(wg.cx.Config)
		Debug(*wg.cx.Config.RPCConnect)
		// wg.ChainMutex.Lock()
		// defer wg.ChainMutex.Unlock()
		if wg.ChainClient, err = rpcclient.New(
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
		); Check(err) {
			return
		}
	}
	if wg.ChainClient.Disconnected() {
		Debug("connecting chain client")
		if err = wg.ChainClient.Connect(1); Check(err) {
			return
		}
	}
	if err = wg.ChainClient.NotifyBlocks(); !Check(err) {
		Debug("subscribed to new blocks")
		// wg.WalletNotifications()
		wg.invalidate <- struct{}{}
	}
	return
}

func (wg *WalletGUI) walletClient() (err error) {
	Debug("connecting to wallet")
	if *wg.cx.Config.WalletOff {
		Warn("wallet is disabled")
		return nil
	}
	// walletRPC := (*wg.cx.Config.WalletRPCListeners)[0]
	certs := pod.ReadCAFile(wg.cx.Config)
	Info("config.tls", *wg.cx.Config.TLS)
	wg.WalletMutex.Lock()
	if wg.WalletClient, err = rpcclient.New(
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
	); Check(err) {
		wg.WalletMutex.Unlock()
		return
	}
	wg.WalletMutex.Unlock()
	// if err = wg.WalletClient.Connect(1); Check(err) {
	// 	return
	// }
	if err = wg.WalletClient.NotifyNewTransactions(true); !Check(err) {
		Debug("subscribed to new transactions")
	} else {
		// return
	}
	if err = wg.WalletClient.NotifyBlocks(); Check(err) {
		// return
	} else {
		Debug("subscribed to wallet client notify blocks")
	}
	Debug("wallet connected")
	return
}
