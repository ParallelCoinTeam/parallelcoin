package gui

import (
	"github.com/p9c/pod/pkg/rpc/btcjson"
	"time"
	
	"github.com/p9c/pod/pkg/util/qu"
)

// Watcher keeps the chain and wallet and rpc clients connected
func (wg *WalletGUI) Watcher() qu.C {
	quit := qu.T()
	// start things up first
	if !wg.node.Running() {
		D.Ln("watcher starting node")
		wg.node.Start()
	}
	if wg.ChainClient == nil {
		D.Ln("chain client is not initialized")
		var e error
		if e = wg.chainClient(); E.Chk(e) {
		}
	}
	var e error
	if !wg.wallet.Running() {
		D.Ln("watcher starting wallet")
		wg.wallet.Start()
		D.Ln("now we can open the wallet")
		if e = wg.writeWalletCookie(); E.Chk(e) {
		}
	}
	if wg.WalletClient == nil || wg.WalletClient.Disconnected() {
	allOut:
		for {
			if e = wg.walletClient(); !E.Chk(e) {
			out:
				for {
					// keep trying until shutdown or the wallet client connects
					var bci *btcjson.GetBlockChainInfoResult
					if bci, e = wg.WalletClient.GetBlockChainInfo(); E.Chk(e) {
						select {
						case <-time.After(time.Second):
							continue
						case <-wg.quit:
							return nil
						}
					}
					D.S(bci)
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
		var e error
	totalOut:
		for {
		disconnected:
			for {
				// D.Ln("top of watcher loop")
				select {
				case <-watchTick.C:
					if !wg.node.Running() {
						D.Ln("watcher starting node")
						wg.node.Start()
					}
					if wg.ChainClient.Disconnected() {
						if e = wg.chainClient(); E.Chk(e) {
							continue
						}
					}
					if !wg.wallet.Running() {
						D.Ln("watcher starting wallet")
						wg.wallet.Start()
					}
					if wg.WalletClient == nil {
						D.Ln("wallet client is not initialized")
						if e = wg.walletClient(); E.Chk(e) {
							continue
							// } else {
							// 	break disconnected
						}
					}
					if wg.WalletClient.Disconnected() {
						if e = wg.WalletClient.Connect(1); D.Chk(e) {
							continue
							// } else {
							// 	break disconnected
						}
					} else {
						D.Ln(
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
						D.Ln(">>>>>>>>>>>>>>>>>>>>>>>>>>>>> wallet not running, breaking out")
						break connected
					}
					if wg.WalletClient == nil || wg.WalletClient.Disconnected() {
						D.Ln(">>>>>>>>>>>>>>>>>>>>>>>>>>>>> wallet client disconnected, breaking out")
						break connected
					}
				case <-quit.Wait():
					break totalOut
				case <-wg.quit.Wait():
					break totalOut
				}
			}
		}
		D.Ln("shutting down watcher")
		if wg.WalletClient != nil {
			wg.WalletClient.Disconnect()
			wg.WalletClient.Shutdown()
		}
		wg.wallet.Stop()
	}()
	return quit
}
