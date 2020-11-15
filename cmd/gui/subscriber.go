package gui

import (
	"encoding/json"
	"time"

	"github.com/p9c/pod/cmd/walletmain"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	rpcclient "github.com/p9c/pod/pkg/rpc/client"
	"github.com/p9c/pod/pkg/util"
)

func (wg *WalletGUI) Subscriber() {
	ntfns := &rpcclient.NotificationHandlers{
		OnClientConnected: func() {},
		OnBlockConnected: func(hash *chainhash.Hash, height int32, t time.Time) {
			// update best block height

			// check account balance

			// pop up new block toast

		},
		OnFilteredBlockConnected:    func(height int32, header *wire.BlockHeader, txs []*util.Tx) {},
		OnBlockDisconnected:         func(hash *chainhash.Hash, height int32, t time.Time) {},
		OnFilteredBlockDisconnected: func(height int32, header *wire.BlockHeader) {},
		OnRecvTx:                    func(transaction *util.Tx, details *btcjson.BlockDetails) {},
		OnRedeemingTx:               func(transaction *util.Tx, details *btcjson.BlockDetails) {},
		OnRelevantTxAccepted:        func(transaction []byte) {},
		OnRescanFinished: func(hash *chainhash.Hash, height int32, blkTime time.Time) {
			// update best block height

			// stop showing syncing indicator

		},
		OnRescanProgress: func(hash *chainhash.Hash, height int32, blkTime time.Time) {
			// update best block height

			// set to show syncing indicator

		},
		OnTxAccepted:        func(hash *chainhash.Hash, amount util.Amount) {},
		OnTxAcceptedVerbose: func(txDetails *btcjson.TxRawResult) {},
		OnPodConnected:      func(connected bool) {},
		OnAccountBalance: func(account string, balance util.Amount, confirmed bool) {
			// what does this actually do
			Debug(account, balance, confirmed)
		},
		OnWalletLockState: func(locked bool) {
			// switch interface to unlock page

			// TODO: lock when idle... how to get trigger for idleness in UI?
		},
		OnUnknownNotification: func(method string, params []json.RawMessage) {},
	}
	_ = ntfns
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
					// close clients if they are open
					if wg.ChainSocket != nil {
						wg.ChainSocket.Disconnect()
						if wg.ChainSocket.Disconnected() {
							wg.ChainSocket = nil
						}
					}
					if wg.WalletSocket != nil {
						wg.WalletSocket.Disconnect()
						if wg.WalletSocket.Disconnected() {
							wg.WalletSocket = nil
						}
					}
					if err = wg.chainSocket(ntfns); Check(err) {
						break
					}
					if err = wg.walletSocket(ntfns); Check(err) {
						break
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
					// check if disconnected
					if wg.ChainSocket.Disconnected() {
						wg.ChainSocket = nil
					}
					if wg.WalletSocket.Disconnected() {
						wg.WalletSocket = nil
					}
					// if we were disconnected flip to connection mode
					if wg.ChainSocket == nil || wg.WalletSocket == nil {
						break out
					}
				case <-wg.quit:
					// close clients if they are open
					if wg.ChainSocket != nil {
						wg.ChainSocket.Disconnect()
						if wg.ChainSocket.Disconnected() {
							wg.ChainSocket = nil
						}
					}
					if wg.WalletSocket != nil {
						wg.WalletSocket.Disconnect()
						if wg.WalletSocket.Disconnected() {
							wg.WalletSocket = nil
						}
					}
					break totalOut
				}
			}
		}
		// Debug("*** Sending shutdown signal")
		// close(wg.quit)
	}()
}

func (wg *WalletGUI) chainSocket(ntfns *rpcclient.NotificationHandlers) (err error) {
	certs := walletmain.ReadCAFile(wg.cx.Config)
	wg.ChainSocket, err = rpcclient.New(&rpcclient.ConnConfig{
		Host:                 *wg.cx.Config.RPCConnect,
		Endpoint:             "ws",
		User:                 *wg.cx.Config.Username,
		Pass:                 *wg.cx.Config.Password,
		TLS:                  true,
		Certificates:         certs,
		Proxy:                "",
		ProxyUser:            "",
		ProxyPass:            "",
		DisableAutoReconnect: false,
		DisableConnectOnNew:  false,
		HTTPPostMode:         false,
		EnableBCInfoHacks:    false,
	}, ntfns)
	return
}

func (wg *WalletGUI) walletSocket(ntfns *rpcclient.NotificationHandlers) (err error) {
	certs := walletmain.ReadCAFile(wg.cx.Config)
	walletRPC := (*wg.cx.Config.WalletRPCListeners)[0]
	wg.WalletSocket, err = rpcclient.New(&rpcclient.ConnConfig{
		Host:                 walletRPC,
		Endpoint:             "ws",
		User:                 *wg.cx.Config.Username,
		Pass:                 *wg.cx.Config.Password,
		TLS:                  true,
		Certificates:         certs,
		Proxy:                "",
		ProxyUser:            "",
		ProxyPass:            "",
		DisableAutoReconnect: false,
		DisableConnectOnNew:  false,
		HTTPPostMode:         false,
		EnableBCInfoHacks:    false,
	}, ntfns)
	return
}
