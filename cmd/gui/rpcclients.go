package gui

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net"
	"runtime/pprof"
	"strings"
	"time"

	l "gioui.org/layout"

	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	rpcclient "github.com/p9c/pod/pkg/rpc/client"
	"github.com/p9c/pod/pkg/util"
)

func (wg *WalletGUI) chainClient() (*rpcclient.Client, error) {
	// update the configuration
	b, err := ioutil.ReadFile(*wg.cx.Config.ConfigFile)
	if err == nil {
		err = json.Unmarshal(b, wg.cx.Config)
		if err != nil {
		}
	}
	return rpcclient.New(&rpcclient.ConnConfig{
		Host:         *wg.cx.Config.RPCConnect,
		User:         *wg.cx.Config.Username,
		Pass:         *wg.cx.Config.Password,
		HTTPPostMode: true,
	}, nil)
}

func (wg *WalletGUI) walletClient() (*rpcclient.Client, error) {
	// update wallet data
	walletRPC := (*wg.cx.Config.WalletRPCListeners)[0]
	var walletServer, port string
	var err error
	if _, port, err = net.SplitHostPort(walletRPC); !Check(err) {
		walletServer = net.JoinHostPort("127.0.0.1", port)
	}
	return rpcclient.New(&rpcclient.ConnConfig{
		Host:         walletServer,
		User:         *wg.cx.Config.Username,
		Pass:         *wg.cx.Config.Password,
		HTTPPostMode: true,
	}, nil)
}

func (wg *WalletGUI) Tickers() {
	go func() {
		seconds := time.Tick(time.Second)
		fiveSeconds := time.Tick(time.Second * 5)
	out:
		for {
			select {
			case <-seconds:
				// Debug("connectChainRPC ticker")
				var chainClient *rpcclient.Client
				var err error
				if chainClient, err = wg.chainClient(); Check(err) {
					break
				}
				var height int32
				var h *chainhash.Hash
				if h, height, err = chainClient.GetBestBlock(); Check(err) {
					break
				}
				wg.State.SetBestBlockHeight(int(height))
				wg.State.SetBestBlockHash(h)
				// // update wallet data
				// walletRPC := (*wg.cx.Config.WalletRPCListeners)[0]
				var walletClient *rpcclient.Client
				// var walletServer, port string
				// if _, port, err = net.SplitHostPort(walletRPC); !Check(err) {
				//	walletServer = net.JoinHostPort("127.0.0.1", port)
				// }
				// walletConnConfig := &rpcclient.ConnConfig{
				//	Host:         walletServer,
				//	User:         *wg.cx.Config.Username,
				//	Pass:         *wg.cx.Config.Password,
				//	HTTPPostMode: true,
				// }
				if walletClient, err = wg.walletClient(); Check(err) {
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
				var ltr []btcjson.ListTransactionsResult
				// TODO: for some reason this function returns half as many as requested
				if ltr, err = walletClient.ListTransactionsCount("default", 20); Check(err) {
					break
				}
				// Debugs(ltr)
				wg.State.SetLastTxs(ltr)
			case <-fiveSeconds:
				var b []byte
				buf := bytes.NewBuffer(b)
				pprof.Lookup("goroutine").WriteTo(buf, 2)
				lines := strings.Split(buf.String(), "\n")
				// Debugs(lines)
				var out []l.Widget
				for i := range lines {
					out = append(out, wg.th.Caption(lines[i]).Color("DocText").Fn)
				}
				wg.State.SetGoroutines(out)
			case <-wg.quit:
				break out
			}
		}
	}()
}
