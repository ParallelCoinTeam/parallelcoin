package gui

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/p9c/pod/cmd/walletmain"
)

func (wg *WalletGUI) ConnectChainRPC() {
	go func() {
		ticker := time.Tick(time.Second)
	out:
		for {
			select {
			case <-ticker:
				Debug("connectChainRPC ticker")
				if wg.ChainClient == nil {
					// update the configuration
					b, err := ioutil.ReadFile(*wg.cx.Config.ConfigFile)
					if err == nil {
						err = json.Unmarshal(b, wg.cx.Config)
						if err != nil {
						}
					} else {
					}
					Debug("connecting to", *wg.cx.Config.RPCConnect)
					if client, err := walletmain.StartChainRPC(wg.cx.Config, wg.cx.ActiveNet,
						walletmain.ReadCAFile(wg.cx.Config)); !Check(err) {
						wg.ChainClient = client
						if err := client.Start(); Check(err) {
							break
						}
						// if err := wg.ChainClient.Start(); Check(err) {
						// 	break
						// }
						// Debug("chain RPC connection succeeded")
						if h, height, err := wg.ChainClient.GetBestBlock(); !Check(err) {
							Debug("updating best block hash and height", h, height)
							wg.State.SetBestBlockHash(h)
							wg.State.SetBestBlockHeight(int(height))
						}
					} else {
						Debug("chain RPC connection failed")
						break
					}
				} else {
					Debug("connected, updating data")
					if h, height, err := wg.ChainClient.GetBestBlock(); !Check(err) {
						Debug("updating best block hash and height", h, height)
						wg.State.SetBestBlockHash(h)
						wg.State.SetBestBlockHeight(int(height))
					}
				}
				if wg.WalletClient == nil {

				}
			case <-wg.quit:
				break out
			}
		}
	}()
}
