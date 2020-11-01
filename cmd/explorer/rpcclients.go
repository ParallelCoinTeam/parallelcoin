package explorer

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"time"

	chainhash "github.com/p9c/pod/pkg/chain/hash"
	rpcclient "github.com/p9c/pod/pkg/rpc/client"
	"github.com/p9c/pod/pkg/util"
)

func (ex *Explorer) chainClient() (*rpcclient.Client, error) {
	return rpcclient.New(&rpcclient.ConnConfig{
		Host:         *ex.cx.Config.RPCConnect,
		User:         *ex.cx.Config.Username,
		Pass:         *ex.cx.Config.Password,
		HTTPPostMode: true,
	}, nil)
}

func (ex *Explorer) ConnectChainRPC() {
	go func() {
		ticker := time.Tick(time.Second)
	out:
		for {
			select {
			case <-ticker:
				// Debug("connectChainRPC ticker")
				// update the configuration
				b, err := ioutil.ReadFile(*ex.cx.Config.ConfigFile)
				if err == nil {
					err = json.Unmarshal(b, ex.cx.Config)
					if err != nil {
					}
				}
				// update chain data
				var chainClient *rpcclient.Client
				//chainConnConfig := &rpcclient.ConnConfig{
				//	Host:         *ex.cx.Config.RPCConnect,
				//	User:         *ex.cx.Config.Username,
				//	Pass:         *ex.cx.Config.Password,
				//	HTTPPostMode: true,
				//}
				if chainClient, err = ex.chainClient(); Check(err) {
					break
				}
				var height int32
				var h *chainhash.Hash
				if h, height, err = chainClient.GetBestBlock(); Check(err) {
					break
				}
				ex.State.SetBestBlockHeight(int(height))
				ex.State.SetBestBlockHash(h)
				// update wallet data
				walletRPC := (*ex.cx.Config.WalletRPCListeners)[0]
				var walletClient *rpcclient.Client
				var walletServer, port string
				if _, port, err = net.SplitHostPort(walletRPC); !Check(err) {
					walletServer = net.JoinHostPort("127.0.0.1", port)
				}
				walletConnConfig := &rpcclient.ConnConfig{
					Host:         walletServer,
					User:         *ex.cx.Config.Username,
					Pass:         *ex.cx.Config.Password,
					HTTPPostMode: true,
				}
				if walletClient, err = rpcclient.New(walletConnConfig, nil); Check(err) {
					break
				}
				var unconfirmed util.Amount
				if unconfirmed, err = walletClient.GetUnconfirmedBalance("default"); Check(err) {
					break
				}
				ex.State.SetBalanceUnconfirmed(unconfirmed.ToDUO())
				var confirmed util.Amount
				if confirmed, err = walletClient.GetBalance("default"); Check(err) {
					break
				}
				ex.State.SetBalance(confirmed.ToDUO())
			case <-ex.quit:
				break out
			}
		}
	}()
}
