package gui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"runtime"
	"runtime/pprof"
	"strings"
	"time"
	
	"github.com/kofoworola/godate"
	
	l "gioui.org/layout"
	
	"github.com/p9c/pod/cmd/walletmain"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/gui/p9"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	rpcclient "github.com/p9c/pod/pkg/rpc/client"
	"github.com/p9c/pod/pkg/util"
)

func (wg *WalletGUI) WalletAndClientRunning() bool {
	return wg.wallet.Running() && wg.WalletClient != nil && !wg.WalletClient.Disconnected()
}

func (wg *WalletGUI) Tickers() {
	first := true
	go func() {
		var err error
		seconds := time.Tick(time.Second)
		// fiveSeconds := time.Tick(time.Second * 5)
	totalOut:
		for {
		preconnect:
			for {
				select {
				case <-seconds:
					Debug("preconnect loop")
					// update goroutines data
					wg.goRoutines()
					// close clients if they are open
					wg.ChainMutex.Lock()
					if wg.ChainClient != nil {
						wg.ChainClient.Disconnect()
						wg.ChainClient.Shutdown()
						wg.ChainClient = nil
					}
					wg.ChainMutex.Unlock()
					wg.WalletMutex.Lock()
					if wg.WalletClient != nil {
						wg.WalletClient.Disconnect()
						wg.WalletClient.Shutdown()
						wg.WalletClient = nil
					}
					wg.WalletMutex.Unlock()
					if !wg.node.Running() {
						break
					}
					if !*wg.cx.Config.NodeOff {
						Debug("connecting to chain")
						if err = wg.chainClient(); Check(err) {
							break
						}
					}
					if wg.wallet.Running() {
						Debug("connecting to wallet")
						if err = wg.walletClient(); Check(err) {
							break
						}
					}
					// if we got to here both are connected
					break preconnect
				case <-wg.quit:
					break totalOut
				}
			}
		out:
			for {
				select {
				case <-seconds:
					Debug("connected loop")
					wg.goRoutines()
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
					// if first {
					var height int32
					var h *chainhash.Hash
					if h, height, err = wg.ChainClient.GetBestBlock(); Check(err) {
						// break out
					}
					wg.State.SetBestBlockHeight(int(height))
					wg.State.SetBestBlockHash(h)
					Debug("wallet running:", wg.wallet.Running())
					Debug("wallet client running:", wg.WalletClient != nil)
					if wg.wallet.Running() && wg.WalletClient == nil {
						Debug("connecting to wallet")
						if err = wg.walletClient(); Check(err) {
							break
						}
					}
					if wg.WalletAndClientRunning() {
						Debug(!wg.WalletClient.Disconnected())
						Debug("wallet is unlocked")
						var unconfirmed util.Amount
						if unconfirmed, err = wg.WalletClient.GetUnconfirmedBalance("default"); Check(err) {
							break out
						}
						wg.State.SetBalanceUnconfirmed(unconfirmed.ToDUO())
						var confirmed util.Amount
						if confirmed, err = wg.WalletClient.GetBalance("default"); Check(err) {
							break out
						}
						wg.State.SetBalance(confirmed.ToDUO())
						Debug("updating recent transactions")
						var atr []btcjson.ListTransactionsResult
						// TODO: for some reason this function returns half as many as requested
						if atr, err = wg.WalletClient.ListTransactionsCountFrom("default", 2<<24, 0); Check(err) {
						}
						wg.State.SetAllTxs(atr)
						Debug("generate the widgets for the updated transactions")
						out := wg.State.FilteredTxs
						if wg.historyTable.Header == nil {
							// create the header
							wg.historyTable.Header = p9.TextTableHeader{
								"Amount",
								"Category",
								"Address",
								"Time",
								"Conf",
								"In Block",
								// "Transaction ID",
								// "Comment",
								// "Fee",
								// "BlockHash",
								// "BlockTime",
								// "Generated",
								// "Abandoned",
								// "Time Received",
								// "Trusted",
								// "Vout",
								// "Wallet Conflicts",
								// "Account",
								// "Other Account",
								// "Involves Watch Only",
							}
						}
						// startIndex := len(wg.historyTable.Body)
						// if wg.State.FilterChanged {
						// 	startIndex = 0
						// 	// all elements must be generated this time
						// 	wg.historyTable.Body = wg.historyTable.Body[:0]
						// }
						// // append all newly added items to the body. The caller can force regeneration by slicing the
						// // body to zero elements or nilling it, as then every element of out will be generated and added
						// // to a fresh empty slice.
						// if startIndex < len(out) {
						// 	// there is new elements appended to the end of the list
						o := out // [startIndex:]
						var bd p9.TextTableBody
						for x := range o {
							i := x
							oi := out[i]
							bd = append(
								bd, p9.TextTableRow{
									fmt.Sprintf("%6.8f", oi.Amount),
									oi.Category,
									oi.Address,
									fmt.Sprintf(
										"%v", godate.Now(time.Local).DifferenceForHumans(
											godate.Create(time.Unix(oi.Time, 0)),
										),
									),
									fmt.Sprintf("%v", oi.Confirmations),
									fmt.Sprintf("%v", *oi.BlockIndex),
									// wg.State.AllTxs[i].TxID,
									// wg.State.AllTxs[i].Comment,
									// fmt.Sprintf("%v", wg.State.AllTxs[i].Fee),
									// wg.State.AllTxs[i].BlockHash,
									// fmt.Sprintf("%v", wg.State.AllTxs[i].BlockTime),
									// fmt.Sprintf("%v", wg.State.AllTxs[i].Generated),
									// fmt.Sprintf("%v", wg.State.AllTxs[i].Abandoned),
									// fmt.Sprintf("%v", wg.State.AllTxs[i].Time),
									// fmt.Sprintf("%v", wg.State.AllTxs[i].Trusted),
									// fmt.Sprintf("%v", wg.State.AllTxs[i].Vout),
									// fmt.Sprintf("%v", wg.State.AllTxs[i].WalletConflicts),
									// wg.State.AllTxs[i].Account,
									// wg.State.AllTxs[i].OtherAccount,
									// fmt.Sprintf("%v", wg.State.AllTxs[i].InvolvesWatchOnly),
								},
							)
							// }
						}
						wg.historyTable.Body = bd // wg.historyTable.Body[:0]
						wg.historyTable.Regenerate(false)
					}
					wg.invalidate <- struct{}{}
					first = false
					// }
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

func (wg *WalletGUI) processChainBlockNotification(hash *chainhash.Hash, height int32, t time.Time) {
	// Debug("processChainBlockNotification")
	// update best block height
	wg.ChainMutex.Lock()
	defer wg.ChainMutex.Unlock()
	if wg.ChainClient == nil || wg.ChainClient.Disconnected() {
		return
	}
	wg.State.SetBestBlockHeight(int(height))
	wg.State.SetBestBlockHash(hash)
}

func (wg *WalletGUI) processWalletBlockNotification() {
	// Debug("processWalletBlockNotification", wg.WalletClient != nil)
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
			go func() {
				if !wg.WalletAndClientRunning() {
					return
				}
				Debug("WALLET CLIENT CONNECTED!")
				// // time.Sleep(time.Second * 3)
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
				// don't update this unless it's in view
				// if wg.ActivePageGet() == "main" {
				Debug("updating recent transactions")
				var atr []btcjson.ListTransactionsResult
				// TODO: for some reason this function returns half as many as requested
				if atr, err = wg.WalletClient.ListTransactionsCountFrom("default", 2<<16, 0); Check(err) {
					// break out
				}
				// Debug(len(atr))
				wg.State.SetAllTxs(atr)
				wg.invalidate <- struct{}{}
			}()
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
			
			wg.invalidate <- struct{}{}
		},
		OnRescanProgress: func(hash *chainhash.Hash, height int32, blkTime time.Time) {
			Debug("OnRescanProgress", hash, height, blkTime)
			// update best block height
			wg.processWalletBlockNotification()
			// set to show syncing indicator
			
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
			wg.stopWallet()
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
	wg.ChainMutex.Lock()
	defer wg.ChainMutex.Unlock()
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

func (wg *WalletGUI) goRoutines() {
	var err error
	if wg.App.ActivePageGet() == "goroutines" || wg.unlockPage.ActivePageGet() == "goroutines" {
		Debug("updating goroutines data")
		var b []byte
		buf := bytes.NewBuffer(b)
		if err = pprof.Lookup("goroutine").WriteTo(buf, 2); Check(err) {
		}
		lines := strings.Split(buf.String(), "\n")
		var out []l.Widget
		var clickables []*p9.Clickable
		for x := range lines {
			i := x
			clickables = append(clickables, wg.th.Clickable())
			var text string
			if strings.HasPrefix(lines[i], "goroutine") && i < len(lines)-2 {
				text = lines[i+2]
				text = strings.TrimSpace(strings.Split(text, " ")[0])
				// outString += text + "\n"
				out = append(
					out, func(gtx l.Context) l.Dimensions {
						return wg.th.ButtonLayout(clickables[i]).Embed(
							wg.th.Inset(
								0.25,
								wg.th.Caption(text).
									Color("DocText").Fn,
							).Fn,
						).Background("Transparent").SetClick(
							func() {
								go func() {
									out := make([]string, 2)
									split := strings.Split(text, ":")
									if len(split) > 2 {
										out[0] = strings.Join(split[:len(split)-1], ":")
										out[1] = split[len(split)-1]
									} else {
										out[0] = split[0]
										out[1] = split[1]
									}
									Debug("path", out[0], "line", out[1])
									goland := "goland64.exe"
									if runtime.GOOS != "windows" {
										goland = "goland"
									}
									launch := exec.Command(goland, "--line", out[1], out[0])
									if err = launch.Start(); Check(err) {
									}
								}()
							},
						).
							Fn(gtx)
					},
				)
			}
		}
		// Debug(outString)
		wg.State.SetGoroutines(out)
		wg.invalidate <- struct{}{}
	}
}
