package gui

import (
	"encoding/json"
	"io/ioutil"
	"time"

	chainhash "github.com/p9c/pod/pkg/chain/hash"
	rpcclient "github.com/p9c/pod/pkg/rpc/client"
)

func (wg *WalletGUI) ConnectChainRPC() {
	go func() {
		ticker := time.Tick(time.Second)
	out:
		for {
			select {
			case <-ticker:
				Debug("connectChainRPC ticker")
				// update the configuration
				b, err := ioutil.ReadFile(*wg.cx.Config.ConfigFile)
				if err == nil {
					err = json.Unmarshal(b, wg.cx.Config)
					if err != nil {
					}
				}
				var client *rpcclient.Client
				connConfig := &rpcclient.ConnConfig{
					Host:                 *wg.cx.Config.RPCConnect,
					User:                 *wg.cx.Config.Username,
					Pass:                 *wg.cx.Config.Password,
					HTTPPostMode:         true,
				}
				if client, err = rpcclient.New(connConfig, nil); Check(err) {
					break
				}
				var height int32
				var h *chainhash.Hash
				if h, height, err = client.GetBestBlock(); Check(err){
					break
				}
				wg.State.bestBlockHeight = int(height)
				wg.State.bestBlockHash = h
			case <-wg.quit:
				break out
			}
		}
	}()
}
