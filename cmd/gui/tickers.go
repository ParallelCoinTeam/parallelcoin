package gui

import (
	"encoding/json"
	"fmt"
	"image"
	"io/ioutil"
	"time"
	
	"gioui.org/op/paint"
	"github.com/atotto/clipboard"
	
	"github.com/p9c/pod/cmd/walletmain"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/coding/qrcode"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	rpcclient "github.com/p9c/pod/pkg/rpc/client"
	"github.com/p9c/pod/pkg/util"
	"github.com/p9c/pod/pkg/util/interrupt"
)

func (wg *WalletGUI) WalletAndClientRunning() bool {
	return wg.wallet.Running() && wg.WalletClient != nil && !wg.WalletClient.Disconnected()
}

func (wg *WalletGUI) Tickers() {
	first := true
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
					// Debug("preconnect loop")
					// update goroutines data
					// wg.goRoutines()
					// close clients if they are open
					// wg.ChainMutex.Lock()
					if wg.ChainClient != nil {
						wg.ChainClient.Disconnect()
						// if wg.ChainClient.Disconnected() {
						wg.ChainClient.Shutdown()
						wg.ChainClient = nil
						// }
					}
					// wg.ChainMutex.Unlock()
					// wg.WalletMutex.Lock()
					if wg.WalletClient != nil {
						wg.WalletClient.Disconnect()
						wg.WalletClient.Shutdown()
						wg.WalletClient = nil
					}
					// wg.WalletMutex.Unlock()
					if !wg.node.Running() {
						break
					}
					Debug("connecting to chain")
					if err = wg.chainClient(); Check(err) {
						break
					}
					Debug("connecting to wallet")
					if err = wg.walletClient(); Check(err) {
						// break
					}
					break preconnect
				case <-fiveSeconds:
					continue
				case <-wg.quit:
					break totalOut
				}
			}
		out:
			for {
				select {
				case <-seconds:
					// Debug("connected loop")
					// wg.goRoutines()
					// the remaining actions require a running shell, if it has been stopped we need to stop
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
					if wg.wallet.Running() && wg.WalletClient == nil {
						Debug("connecting to wallet")
						if err = wg.walletClient(); Check(err) {
							break
						}
					}
					if wg.WalletAndClientRunning() {
						if first {
							wg.processWalletBlockNotification()
						}
						if wg.stateLoaded.Load() && !wg.State.IsReceivingAddress() || wg.currentReceiveGetNew.Load() {
							wg.currentReceiveGetNew.Store(false)
							var addr util.Address
							if addr, err = wg.WalletClient.GetNewAddress("default"); !Check(err) {
								Debug("storing new receiving address", addr.EncodeAddress(),
									"as prior was empty", wg.State.currentReceivingAddress.String.Load())
								wg.State.SetReceivingAddress(addr)
								wg.State.isAddress.Store(true)
								Debugs(wg.State.Marshal())
								wg.invalidate <- struct{}{}
							}
						}
					}
					if (wg.currentReceiveQRCode == nil && wg.currentReceiveAddress != "" && wg.stateLoaded.Load()) ||
						wg.currentReceiveRegenerate.Load() {
						var qrc image.Image
						wg.currentReceiveRegenerate.Store(false)
						Debug("generating QR code")
						var err error
						qrText := fmt.Sprintf("parallelcoin:%s?amount=%s&message=%s",
							wg.State.currentReceivingAddress.Load().EncodeAddress(),
							wg.inputs["receiveAmount"].GetText(),
							wg.inputs["receiveMessage"].GetText(),
						)
						if qrc, err = qrcode.Encode(qrText, 0, qrcode.ECLevelQ, 6); !Check(err) {
							iop := paint.NewImageOp(qrc)
							wg.currentReceiveQRCode = &iop
							wg.currentReceiveQR = wg.ButtonLayout(wg.currentReceiveCopyClickable.SetClick(func() {
								Debug("clicked qr code copy clicker")
								if err := clipboard.WriteAll(qrText); Check(err) {
								}
							})).
								CornerRadius(0.5).
								Background("white").
								Embed(
									wg.Inset(0.5,
										wg.Image().Src(*wg.currentReceiveQRCode).Scale(1).Fn,
									).Fn,
								).Fn
							// *wg.currentReceiveQRCode = iop
						}
						
					}
					wg.invalidate <- struct{}{}
					first = false
				// }
				case <-fiveSeconds:
				case <-wg.quit:
					break totalOut
				}
			}
			// wg.runningNode.Store(false)
		}
		// Debug("*** Sending shutdown signal")
		// close(wg.quit)
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
	if wg.ChainClient == nil || wg.ChainClient.Disconnected() {
		return
	}
	var err error
	var h *chainhash.Hash
	var height int32
	if h, height, err = wg.ChainClient.GetBestBlock(); Check(err) {
		interrupt.Request()
		return
	}
	wg.State.SetBestBlockHeight(height)
	wg.State.SetBestBlockHash(h)
}

func (wg *WalletGUI) processChainBlockNotification(hash *chainhash.Hash, height int32, t time.Time) {
	Debug("processChainBlockNotification")
	wg.State.SetBestBlockHeight(height)
	wg.State.SetBestBlockHash(hash)
	// wg.invalidate <- struct{}{}
}

func (wg *WalletGUI) processWalletBlockNotification() {
	Debug("processWalletBlockNotification")
	if !wg.WalletAndClientRunning() {
		return
	}
	// check account balance
	var unconfirmed util.Amount
	var err error
	if unconfirmed, err = wg.WalletClient.GetUnconfirmedBalance("default"); Check(err) {
		// break out
	}
	wg.State.SetBalanceUnconfirmed(unconfirmed.ToDUO())
	var confirmed util.Amount
	if confirmed, err = wg.WalletClient.GetBalance("default"); Check(err) {
		// break out
	}
	wg.State.SetBalance(confirmed.ToDUO())
	var atr []btcjson.ListTransactionsResult
	// TODO: for some reason this function returns half as many as requested
	if atr, err = wg.WalletClient.ListTransactionsCountFrom("default", 2<<16, 0); Check(err) {
		// break out
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
}

func (wg *WalletGUI) ChainNotifications() *rpcclient.NotificationHandlers {
	return &rpcclient.NotificationHandlers{
		OnClientConnected: func() {
			// go func() {
			Debug("CHAIN CLIENT CONNECTED!")
			// var err error
			// var height int32
			// var h *chainhash.Hash
			// if h, height, err = wg.ChainClient.GetBestBlock(); Check(err) {
			// }
			// wg.State.SetBestBlockHeight(int(height))
			// wg.State.SetBestBlockHash(h)
			// wg.invalidate <- struct{}{}
			// }()
		},
		OnBlockConnected: func(hash *chainhash.Hash, height int32, t time.Time) {
			Trace("chain OnBlockConnected", hash, height, t)
			wg.processChainBlockNotification(hash, height, t)
			wg.processWalletBlockNotification()
			// pop up new block toast
			
			wg.invalidate <- struct{}{}
		},
		// OnFilteredBlockConnected:    func(height int32, header *wire.BlockHeader, txs []*util.Tx) {},
		// OnBlockDisconnected:         func(hash *chainhash.Hash, height int32, t time.Time) {},
		// OnFilteredBlockDisconnected: func(height int32, header *wire.BlockHeader) {},
		// OnRecvTx:                    func(transaction *util.Tx, details *btcjson.BlockDetails) {},
		// OnRedeemingTx:               func(transaction *util.Tx, details *btcjson.BlockDetails) {},
		// OnRelevantTxAccepted:        func(transaction []byte) {},
		// OnRescanFinished: func(hash *chainhash.Hash, height int32, blkTime time.Time) {
		// 	Debug("OnRescanFinished", hash, height, blkTime)
		// 	// update best block height
		//
		// 	// stop showing syncing indicator
		//
		// },
		// OnRescanProgress: func(hash *chainhash.Hash, height int32, blkTime time.Time) {
		// 	Debug("OnRescanProgress", hash, height, blkTime)
		// 	// update best block height
		//
		// 	// set to show syncing indicator
		//
		// },
		// OnTxAccepted:        func(hash *chainhash.Hash, amount util.Amount) {},
		// OnTxAcceptedVerbose: func(txDetails *btcjson.TxRawResult) {},
		// OnPodConnected:      func(connected bool) {},
		// OnAccountBalance: func(account string, balance util.Amount, confirmed bool) {
		// 	Debug("OnAccountBalance")
		// 	// what does this actually do
		// 	Debug(account, balance, confirmed)
		// },
		// OnWalletLockState: func(locked bool) {
		// 	Debug("OnWalletLockState", locked)
		// 	// switch interface to unlock page
		//
		// 	// TODO: lock when idle... how to get trigger for idleness in UI?
		// },
		// OnUnknownNotification: func(method string, params []json.RawMessage) {},
	}
	
}

func (wg *WalletGUI) WalletNotifications() *rpcclient.NotificationHandlers {
	if !wg.wallet.Running() || wg.WalletClient == nil || wg.WalletClient.Disconnected() {
		return nil
	}
	return &rpcclient.NotificationHandlers{
		OnClientConnected: func() {
			wg.processWalletBlockNotification()
			wg.invalidate <- struct{}{}
		},
		// OnBlockConnected: func(hash *chainhash.Hash, height int32, t time.Time) {
		// 	Debug("wallet OnBlockConnected", hash, height, t)
		// 	wg.processWalletBlockNotification()
		// 	wg.processChainBlockNotification(hash, height, t)
		// wg.invalidate <- struct{}{}
		// },
		// OnFilteredBlockConnected:    func(height int32, header *wire.BlockHeader, txs []*util.Tx) {},
		// OnBlockDisconnected:         func(hash *chainhash.Hash, height int32, t time.Time) {},
		// OnFilteredBlockDisconnected: func(height int32, header *wire.BlockHeader) {},
		// OnRecvTx:                    func(transaction *util.Tx, details *btcjson.BlockDetails) {},
		// OnRedeemingTx:               func(transaction *util.Tx, details *btcjson.BlockDetails) {},
		// OnRelevantTxAccepted:        func(transaction []byte) {},
		OnRescanFinished: func(hash *chainhash.Hash, height int32, blkTime time.Time) {
			Debug("OnRescanFinished", hash, height, blkTime)
			// update best block height
			wg.processWalletBlockNotification()
			// stop showing syncing indicator
			wg.Syncing.Store(false)
			wg.invalidate <- struct{}{}
		},
		OnRescanProgress: func(hash *chainhash.Hash, height int32, blkTime time.Time) {
			Debug("OnRescanProgress", hash, height, blkTime)
			// update best block height
			// wg.processWalletBlockNotification()
			// set to show syncing indicator
			wg.Syncing.Store(true)
			wg.invalidate <- struct{}{}
		},
		// OnTxAccepted:        func(hash *chainhash.Hash, amount util.Amount) {},
		// OnTxAcceptedVerbose: func(txDetails *btcjson.TxRawResult) {},
		// // OnPodConnected:      func(connected bool) {},
		OnAccountBalance: func(account string, balance util.Amount, confirmed bool) {
			Debug("OnAccountBalance")
			// what does this actually do
			Debug(account, balance, confirmed)
		},
		OnWalletLockState: func(locked bool) {
			Debug("OnWalletLockState", locked)
			// switch interface to unlock page
			wg.wallet.Stop()
			// TODO: lock when idle... how to get trigger for idleness in UI?
			wg.invalidate <- struct{}{}
		},
		// OnUnknownNotification: func(method string, params []json.RawMessage) {},
	}
	
}

func (wg *WalletGUI) chainClient() (err error) {
	Debug("starting up chain client")
	if *wg.cx.Config.NodeOff {
		Warn("node is disabled")
		return nil
	}
	certs := walletmain.ReadCAFile(wg.cx.Config)
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
	if err = wg.ChainClient.NotifyBlocks(); !Check(err) {
		Debug("subscribed to new transactions")
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
	certs := walletmain.ReadCAFile(wg.cx.Config)
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
	if err = wg.WalletClient.NotifyNewTransactions(true); !Check(err) {
		Debug("subscribed to new transactions")
		// wg.WalletNotifications()
		wg.invalidate <- struct{}{}
	}
	Debug("wallet connected")
	return
}

// func (wg *WalletGUI) goRoutines() {
// 	var err error
// 	if wg.App.ActivePageGet() == "goroutines" || wg.unlockPage.ActivePageGet() == "goroutines" {
// 		Debug("updating goroutines data")
// 		var b []byte
// 		buf := bytes.NewBuffer(b)
// 		if err = pprof.Lookup("goroutine").WriteTo(buf, 2); Check(err) {
// 		}
// 		lines := strings.Split(buf.String(), "\n")
// 		var out []l.Widget
// 		var clickables []*p9.Clickable
// 		for x := range lines {
// 			i := x
// 			clickables = append(clickables, wg.Clickable())
// 			var text string
// 			if strings.HasPrefix(lines[i], "goroutine") && i < len(lines)-2 {
// 				text = lines[i+2]
// 				text = strings.TrimSpace(strings.Split(text, " ")[0])
// 				// outString += text + "\n"
// 				out = append(
// 					out, func(gtx l.Context) l.Dimensions {
// 						return wg.ButtonLayout(clickables[i]).Embed(
// 							wg.ButtonInset(
// 								0.25,
// 								wg.Caption(text).
// 									Color("DocText").Fn,
// 							).Fn,
// 						).Background("Transparent").SetClick(
// 							func() {
// 								go func() {
// 									out := make([]string, 2)
// 									split := strings.Split(text, ":")
// 									if len(split) > 2 {
// 										out[0] = strings.Join(split[:len(split)-1], ":")
// 										out[1] = split[len(split)-1]
// 									} else {
// 										out[0] = split[0]
// 										out[1] = split[1]
// 									}
// 									Debug("path", out[0], "line", out[1])
// 									goland := "goland64.exe"
// 									if runtime.GOOS != "windows" {
// 										goland = "goland"
// 									}
// 									launch := exec.Command(goland, "--line", out[1], out[0])
// 									if err = launch.Start(); Check(err) {
// 									}
// 								}()
// 							},
// 						).
// 							Fn(gtx)
// 					},
// 				)
// 			}
// 		}
// 		// Debug(outString)
// 		wg.State.SetGoroutines(out)
// 		wg.invalidate <- struct{}{}
// 	}
// }
