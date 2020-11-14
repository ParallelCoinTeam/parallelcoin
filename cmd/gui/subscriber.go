package gui

import (
	"encoding/json"
	"time"

	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	rpcclient "github.com/p9c/pod/pkg/rpc/client"
	"github.com/p9c/pod/pkg/util"
)

func (wg *WalletGUI) Subscriber() {
	out := &rpcclient.NotificationHandlers{
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
	_ = out
}

func (wg *WalletGUI) chainSocket(ntfns *rpcclient.NotificationHandlers) (err error) {
	wg.ChainSocket, err = rpcclient.New(&rpcclient.ConnConfig{
		Host:         *wg.cx.Config.RPCConnect,
		User:         *wg.cx.Config.Username,
		Pass:         *wg.cx.Config.Password,
		HTTPPostMode: false,
	}, ntfns)
	return
}

func (wg *WalletGUI) walletSocket(ntfns *rpcclient.NotificationHandlers) (err error) {
	walletRPC := (*wg.cx.Config.WalletRPCListeners)[0]
	wg.WalletSocket, err = rpcclient.New(&rpcclient.ConnConfig{
		Host:         walletRPC,
		User:         *wg.cx.Config.Username,
		Pass:         *wg.cx.Config.Password,
		HTTPPostMode: false,
	}, ntfns)
	return
}
