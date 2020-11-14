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

func (wg *WalletGUI) Subscriber() *rpcclient.NotificationHandlers {
	out := &rpcclient.NotificationHandlers{
		OnClientConnected: func() {},
		OnBlockConnected: func(hash *chainhash.Hash, height int32, t time.Time) {
			// check account balance

			// pop up new block toast

		},
		OnFilteredBlockConnected: func(height int32, header *wire.BlockHeader, txs []*util.Tx) {},
		OnBlockDisconnected: func(hash *chainhash.Hash, height int32, t time.Time) {},
		OnFilteredBlockDisconnected: func(height int32, header *wire.BlockHeader) {},
		OnRecvTx: func(transaction *util.Tx, details *btcjson.BlockDetails) {},
		OnRedeemingTx: func(transaction *util.Tx, details *btcjson.BlockDetails) {},
		OnRelevantTxAccepted: func(transaction []byte) {},
		OnRescanFinished: func(hash *chainhash.Hash, height int32, blkTime time.Time) {
			// update best block height

			// stop showing syncing indicator
		},
		OnRescanProgress: func(hash *chainhash.Hash, height int32, blkTime time.Time) {
			// update best block height

			// set to show syncing indicator
		},
		OnTxAccepted: func(hash *chainhash.Hash, amount util.Amount) {},
		OnTxAcceptedVerbose: func(txDetails *btcjson.TxRawResult) {},
		OnPodConnected: func(connected bool) {},
		OnAccountBalance: func(account string, balance util.Amount, confirmed bool) {
			// what does this actually do
			Debug(account, balance, confirmed)
		},
		OnWalletLockState: func(locked bool) {
			// switch interface to unlock page
		},
		OnUnknownNotification: func(method string, params []json.RawMessage) {},
	}
	return out
}
