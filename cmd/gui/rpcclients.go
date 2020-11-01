package gui

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"time"

	chainhash "github.com/p9c/pod/pkg/chain/hash"
	rpcclient "github.com/p9c/pod/pkg/rpc/client"
	"github.com/p9c/pod/pkg/util"
)

func (wg *WalletGUI) chainClient() (*rpcclient.Client, error) {
	return rpcclient.New(&rpcclient.ConnConfig{
		Host:         *wg.cx.Config.RPCConnect,
		User:         *wg.cx.Config.Username,
		Pass:         *wg.cx.Config.Password,
		HTTPPostMode: true,
	}, nil)
}

func (wg *WalletGUI) ConnectChainRPC() {
	go func() {
		ticker := time.Tick(time.Second)
		connectedOnce := false
	out:
		for {
			select {
			case <-ticker:
				// Debug("connectChainRPC ticker")
				// update the configuration
				var err error
				if !connectedOnce {
					var b []byte
					b, err = ioutil.ReadFile(*wg.cx.Config.ConfigFile)
					if err == nil {
						err = json.Unmarshal(b, wg.cx.Config)
						if err != nil {
						}
					}
				}
				// update chain data
				var chainClient *rpcclient.Client
				// chainConnConfig := &rpcclient.ConnConfig{
				//	Host:         *wg.cx.Config.RPCConnect,
				//	User:         *wg.cx.Config.Username,
				//	Pass:         *wg.cx.Config.Password,
				//	HTTPPostMode: true,
				// }
				if chainClient, err = wg.chainClient(); Check(err) {
					break
				}
				var height int32
				var h *chainhash.Hash
				if h, height, err = chainClient.GetBestBlock(); Check(err) {
					break
				}
				connectedOnce = true
				wg.State.SetBestBlockHeight(int(height))
				wg.State.SetBestBlockHash(h)
				// update wallet data
				walletRPC := (*wg.cx.Config.WalletRPCListeners)[0]
				var walletClient *rpcclient.Client
				var walletServer, port string
				if _, port, err = net.SplitHostPort(walletRPC); !Check(err) {
					walletServer = net.JoinHostPort("127.0.0.1", port)
				}
				walletConnConfig := &rpcclient.ConnConfig{
					Host:         walletServer,
					User:         *wg.cx.Config.Username,
					Pass:         *wg.cx.Config.Password,
					HTTPPostMode: true,
				}
				if walletClient, err = rpcclient.New(walletConnConfig, nil); Check(err) {
					break
				}
				var unconfirmed util.Amount
				if unconfirmed, err = walletClient.GetUnconfirmedBalance("default"); Check(err) {
					break
				}
				wg.State.SetBalanceUnconfirmed(unconfirmed.ToDUO())
				var confirmed util.Amount
				if confirmed, err = walletClient.GetBalance("default"); Check(err) {
					break
				}
				wg.State.SetBalance(confirmed.ToDUO())
			case <-wg.quit:
				break out
			}
		}
	}()
}
