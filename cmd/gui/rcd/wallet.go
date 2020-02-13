package rcd

import (
	"fmt"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/node/rpc"
	"github.com/p9c/pod/pkg/chain/config/netparams"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	"github.com/p9c/pod/pkg/rpc/legacy"
	"github.com/p9c/pod/pkg/util"
	"time"
)

func
(rc *RcVar) GetDuoUIbalance(cx *conte.Xt) {
	acct := "default"
	minconf := 0
	getBalance, err := legacy.GetBalance(&btcjson.GetBalanceCmd{Account: &acct,
		MinConf: &minconf}, cx.WalletServer)
	if err != nil {
		rc.PushDuoUIalert("Error", err.Error(), "error")
	}
	gb, ok := getBalance.(float64)
	if ok {
		bb := fmt.Sprintf("%0.8f", gb)
		rc.Balance = bb
	}
	return
}
func
(rc *RcVar) GetDuoUIunconfirmedBalance(cx *conte.Xt) {
	acct := "default"
	getUnconfirmedBalance, err := legacy.GetUnconfirmedBalance(&btcjson.GetUnconfirmedBalanceCmd{Account: &acct}, cx.WalletServer)
	if err != nil {
		rc.PushDuoUIalert("Error", err.Error(), "error")
	}
	ub, ok := getUnconfirmedBalance.(float64)
	if ok {
		ubb := fmt.Sprintf("%0.8f", ub)
		rc.Unconfirmed = ubb
	}
	return
}

func
(rc *RcVar) GetDuoUItransactionsNumber(cx *conte.Xt) {
	// account, txcount, startnum, watchonly := "*", n, f, false
	// listTransactions, err := legacy.ListTransactions(&json.ListTransactionsCmd{Account: &account, Count: &txcount, From: &startnum, IncludeWatchOnly: &watchonly}, v.ws)
	lt, err := cx.WalletServer.ListTransactions(0, 999999999)
	if err != nil {
		rc.PushDuoUIalert("Error", err.Error(), "error")
	}
	rc.TxsNumber = len(lt)
}

func
(rc *RcVar) GetDuoUItransactions(cx *conte.Xt, sfrom, count int, cat string) *models.DuoUItransactions {
	// account, txcount, startnum, watchonly := "*", n, f, false
	// listTransactions, err := legacy.ListTransactions(&json.ListTransactionsCmd{Account: &account, Count: &txcount, From: &startnum, IncludeWatchOnly: &watchonly}, v.ws)
	lt, err := cx.WalletServer.ListTransactions(sfrom, count)
	if err != nil {
		rc.PushDuoUIalert("Error", err.Error(), "error")
	}
	rc.Transactions.TxsNumber = len(lt)
	txsArray := *new([]models.DuoUItx)
	// lt := listTransactions.([]json.ListTransactionsResult)
	switch c := cat; c {
	case "received":
		for _, tx := range lt {
			if tx.Category == "received" {
				txsArray = append(txsArray, txs(tx))
			}
		}
	case "sent":
		for _, tx := range lt {
			if tx.Category == "sent" {
				txsArray = append(txsArray, txs(tx))
			}
		}
	case "immature":
		for _, tx := range lt {
			if tx.Category == "immature" {
				txsArray = append(txsArray, txs(tx))
			}
		}
	case "generate":
		for _, tx := range lt {
			if tx.Category == "generate" {
				txsArray = append(txsArray, txs(tx))
			}
		}
	default:
		for _, tx := range lt {
			txsArray = append(txsArray, txs(tx))
		}
	}
	rc.Transactions.Txs = txsArray
	return &rc.Transactions
}
func
txs(t btcjson.ListTransactionsResult) models.DuoUItx {
	return models.DuoUItx{
		TxID:     t.TxID,
		Amount:   t.Amount,
		Category: t.Category,
		Time:     time.Unix(t.Time, 0),
	}

}
func
(rc *RcVar) GetDuoUIlastTxs(cx *conte.Xt) {
	rc.LastTxs = *rc.GetDuoUItransactions(cx, 0, 10, "")
	return
}
func
(rc *RcVar) GetDuoUITransactionsExcertps( cx *conte.Xt) {
	lt, err := cx.WalletServer.ListTransactions(0, rc.Txs.TxsListNumber)
	if err != nil {
		rc.PushDuoUIalert("Error", err.Error(), "error")
	}
	rc.Txs.TxsNumber = len(lt)
	// for i, j := 0, len(lt)-1; i < j; i, j = i+1, j-1 {
	//	lt[i], lt[j] = lt[j], lt[i]
	// }
	balanceHeight := 0.0
	txseRaw := []models.DuoUItransactionExcerpt{}
	for _, txRaw := range lt {
		unixTimeUTC := time.Unix(txRaw.Time, 0) // gives unix time stamp in utc
		txseRaw = append(txseRaw, models.DuoUItransactionExcerpt{
			// Balance:       txse.Balance + txRaw.Amount,
			Comment:       txRaw.Comment,
			Amount:        txRaw.Amount,
			Category:      txRaw.Category,
			Confirmations: txRaw.Confirmations,
			Time:          unixTimeUTC.Format(time.RFC3339),
			TxID:          txRaw.TxID,
		})
	}
	var balance float64
	txs := *new([]models.DuoUItransactionExcerpt)
	for _, tx := range txseRaw {
		balance = balance + tx.Amount
		tx.Balance = balance
		txs = append(txs, tx)
		if rc.Txs.Balance > balanceHeight {
			balanceHeight = rc.Txs.Balance
		}

	}
	rc.Txs.Txs = txs
	rc.Txs.BalanceHeight = balanceHeight
	return
}

func
(rc *RcVar) DuoSend(cx *conte.Xt, wp string, ad string, am float64) {
	if am > 0 {
		getBlockChain, err := rpc.HandleGetBlockChainInfo(cx.RPCServer, nil, nil)
		if err != nil {
			rc.PushDuoUIalert("Error", err.Error(), "error")

		}
		result, ok := getBlockChain.(*btcjson.GetBlockChainInfoResult)
		if !ok {
			result = &btcjson.GetBlockChainInfoResult{}
		}
		var defaultNet *netparams.Params
		switch result.Chain {
		case "mainnet":
			defaultNet = &netparams.MainNetParams
		case "testnet":
			defaultNet = &netparams.TestNet3Params
		case "regtest":
			defaultNet = &netparams.RegressionTestParams
		default:
			defaultNet = &netparams.MainNetParams
		}
		amount, _ := util.NewAmount(am)
		addr, err := util.DecodeAddress(ad, defaultNet)
		if err != nil {
			rc.PushDuoUIalert("Error", err.Error(), "error")
		}
		var validateAddr *btcjson.ValidateAddressWalletResult
		if err == nil {
			var va interface{}
			va, err = legacy.ValidateAddress(&btcjson.
			ValidateAddressCmd{Address: addr.String()}, cx.WalletServer)
			if err != nil {
				rc.PushDuoUIalert("Error", err.Error(), "error")
			}
			vva := va.(btcjson.ValidateAddressWalletResult)
			validateAddr = &vva
			if validateAddr.IsValid {
				legacy.WalletPassphrase(btcjson.NewWalletPassphraseCmd(wp, 5),
					cx.WalletServer)
				if err != nil {
					rc.PushDuoUIalert("Error", err.Error(), "error")
				}
				_, err = legacy.SendToAddress(
					&btcjson.SendToAddressCmd{
						Address: addr.EncodeAddress(), Amount: amount.ToDUO(),
					}, cx.WalletServer)
				if err != nil {
					rc.PushDuoUIalert("error sending to address:", err.Error(), "error")

				} else {
					rc.PushDuoUIalert("Address OK", "OK", "success")
				}
			} else {
				if err != nil {
					rc.PushDuoUIalert("Invalid address", "INVALID", "error")
				}
			}
			rc.PushDuoUIalert("Payment sent", "PAYMENT", "success")
		}
	} else {
		log.Println(am)
	}
	rc.Sent = true
	return
}
