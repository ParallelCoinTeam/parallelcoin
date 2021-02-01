package gui

import (
	"time"
	
	qu "github.com/p9c/pod/pkg/util/quit"
)

func (wg *WalletGUI) Watcher() qu.C {
	quit := qu.T()
	go func() {
		watchTick := time.NewTicker(time.Second)
		var err error
	totalOut:
		for {
		disconnected:
			for {
				Debug("top of watcher loop")
				select {
				case <-watchTick.C:
					if !wg.wallet.Running() {
						Debug("watcher starting wallet")
						wg.wallet.Start()
					}
					if wg.WalletClient == nil {
						Debug("wallet client is not initialized")
						if err = wg.walletClient(); Check(err) {
							continue
						} else {
							break disconnected
						}
					}
					if wg.WalletClient.Disconnected() {
						if err = wg.WalletClient.Connect(1); Check(err) {
						} else {
							break disconnected
						}
					}
				case <-quit.Wait():
					break totalOut
				case <-wg.quit.Wait():
					break totalOut
				}
			}
			Debug("wallet and client are now connected")
		connected:
			for {
				select {
				case <-watchTick.C:
					if !wg.wallet.Running() || wg.WalletClient.Disconnected() {
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
