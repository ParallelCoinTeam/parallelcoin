package gui

import (
	"github.com/p9c/pod/pkg/rpc/btcjson"
	"time"
	
	qu "github.com/p9c/pod/pkg/util/quit"
)

// Watcher keeps the chain and wallet and rpc clients connected
func (wg *WalletGUI) Watcher() qu.C {
	quit := qu.T()
	// start things up first
	if !wg.node.Running() {
		Debug("watcher starting node")
		wg.node.Start()
	}
	if wg.ChainClient == nil {
		Debug("chain client is not initialized")
		var err error
		if err = wg.chainClient(); Check(err) {
		}
	}
	var err error
	if !wg.wallet.Running() {
		Debug("watcher starting wallet")
		wg.wallet.Start()
		Debug("now we can open the wallet")
		if err = wg.writeWalletCookie(); Check(err) {
		}
	}
	if wg.WalletClient == nil || wg.WalletClient.Disconnected() {
	allOut:
		for {
			if err = wg.walletClient(); !Check(err) {
			out:
				for {
					// keep trying until shutdown or the wallet client connects
					var bci *btcjson.GetBlockChainInfoResult
					if bci, err = wg.WalletClient.GetBlockChainInfo(); Check(err) {
						select {
						case <-time.After(time.Second):
							continue
						case <-wg.quit:
							return nil
						}
					}
					Debugs(bci)
					break out
				}
			}
			wg.unlockPassword.Wipe()
			wg.ready.Store(true)
			wg.Invalidate()
			select {
			case <-time.After(time.Second):
				break allOut
			case <-wg.quit:
				return nil
			}
		}
	}
	go func() {
		watchTick := time.NewTicker(time.Second)
		var err error
	totalOut:
		for {
		disconnected:
			for {
				// Debug("top of watcher loop")
				select {
				case <-watchTick.C:
					if !wg.node.Running() {
						Debug("watcher starting node")
						wg.node.Start()
					}
					if wg.ChainClient.Disconnected() {
						if err = wg.chainClient(); Check(err) {
							continue
						}
					}
					if !wg.wallet.Running() {
						Debug("watcher starting wallet")
						wg.wallet.Start()
					}
					if wg.WalletClient == nil {
						Debug("wallet client is not initialized")
						if err = wg.walletClient(); Check(err) {
							continue
							// } else {
							// 	break disconnected
						}
					}
					if wg.WalletClient.Disconnected() {
						if err = wg.WalletClient.Connect(1); Check(err) {
							continue
							// } else {
							// 	break disconnected
						}
					} else {
						Debug(
							"chain, chainclient, wallet and client are now connected",
							wg.node.Running(),
							!wg.ChainClient.Disconnected(),
							wg.wallet.Running(),
							!wg.WalletClient.Disconnected(),
						)
						wg.updateChainBlock()
						wg.processWalletBlockNotification()
						break disconnected
					}
				case <-quit.Wait():
					break totalOut
				case <-wg.quit.Wait():
					break totalOut
				}
			}
		connected:
			for {
				select {
				case <-watchTick.C:
					if !wg.wallet.Running() {
						Debug(">>>>>>>>>>>>>>>>>>>>>>>>>>>>> wallet not running, breaking out")
						break connected
					}
					if wg.WalletClient == nil || wg.WalletClient.Disconnected() {
						Debug(">>>>>>>>>>>>>>>>>>>>>>>>>>>>> wallet client disconnected, breaking out")
						break connected
					}
				case <-quit.Wait():
					break totalOut
				case <-wg.quit.Wait():
					break totalOut
				}
			}
		}
		Debug("shutting down watcher")
		if wg.WalletClient != nil {
			wg.WalletClient.Disconnect()
			wg.WalletClient.Shutdown()
		}
		wg.wallet.Stop()
	}()
	return quit
}
