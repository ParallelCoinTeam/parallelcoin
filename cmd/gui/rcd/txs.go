package rcd

import (
	"fmt"
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/model"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	"time"
)

func (r *RcVar) GetDuoUItransactionsNumber() {
	log.DEBUG("getting transaction count")
	// account, txcount, startnum, watchonly := "*", n, f, false
	// listTransactions, err := legacy.ListTransactions(&json.ListTransactionsCmd{Account: &account, Count: &txcount, From: &startnum, IncludeWatchOnly: &watchonly}, v.ws)
	lt, err := r.cx.WalletServer.ListTransactions(0, 999999999)
	if err != nil {
		// r.PushDuoUIalert("Error", err.Error(), "error")
	}
	r.Status.Wallet.TxsNumber.Store(uint64(len(lt)))
	return
}

func (r *RcVar) GetDuoUItransactions(sfrom, count int, cat string) *model.DuoUItransactions {
	log.DEBUG("getting transactions")
	// account, txcount, startnum, watchonly := "*", n, f, false
	// listTransactions, err := legacy.ListTransactions(&json.ListTransactionsCmd{Account: &account, Count: &txcount, From: &startnum, IncludeWatchOnly: &watchonly}, v.ws)
	lt, err := r.cx.WalletServer.ListTransactions(0, 11)
	if err != nil {
		log.INFO(err)
	}
	r.Status.Wallet.Transactions.TxsNumber = len(lt)
	log.INFO("ETO:" + fmt.Sprint(lt))
	txsArray := *new([]model.DuoUItx)
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
	r.Status.Wallet.Transactions.Txs = txsArray
	return r.Status.Wallet.Transactions
}
func txs(t btcjson.ListTransactionsResult) model.DuoUItx {
	return model.DuoUItx{
		TxID:     t.TxID,
		Amount:   t.Amount,
		Category: t.Category,
		Time:     helpers.FormatTime(time.Unix(t.Time, 0)),
	}

}
func (r *RcVar) GetLatestTransactions() {
	log.DEBUG("getting latest transactions")
	//r.Status.Wallet.LastTxs = r.GetDuoUItransactions(0, 10, "")
	return
}

func (r *RcVar) GetTransactions() func() {
	return func() {
		log.DEBUG("getting transactions")
		lt, err := r.cx.WalletServer.ListTransactions(0, r.Status.Wallet.Txs.TxsListNumber)
		if err != nil {
			// //r.PushDuoUIalert("Error", err.Error(), "error")
		}
		r.Status.Wallet.Txs.TxsNumber = len(lt)
		// for i, j := 0, len(lt)-1; i < j; i, j = i+1, j-1 {
		//	lt[i], lt[j] = lt[j], lt[i]
		// }
		balanceHeight := 0.0
		txseRaw := []model.DuoUItransactionExcerpt{}
		for _, txRaw := range lt {
			unixTimeUTC := time.Unix(txRaw.Time, 0) // gives unix time stamp in utc
			txseRaw = append(txseRaw, model.DuoUItransactionExcerpt{
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
		txs := *new([]model.DuoUItransactionExcerpt)
		for _, tx := range txseRaw {
			balance = balance + tx.Amount
			tx.Balance = balance
			txs = append(txs, tx)
			if r.Status.Wallet.Txs.Balance > balanceHeight {
				balanceHeight = r.Status.Wallet.Txs.Balance
			}

		}
		r.Status.Wallet.Txs.Txs = txs
		r.Status.Wallet.Txs.BalanceHeight = balanceHeight
		return
	}
}
