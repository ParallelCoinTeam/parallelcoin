package rcd

import (
	"fmt"
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/model"
	log "github.com/p9c/logi"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	"time"
)

func (r *RcVar) GetDuoUItransactionsNumber() {
	log.L.Debug("getting transaction count")
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
	log.L.Debug("getting transactions")
	// account, txcount, startnum, watchonly := "*", n, f, false
	// listTransactions, err := legacy.ListTransactions(&json.ListTransactionsCmd{Account: &account, Count: &txcount, From: &startnum, IncludeWatchOnly: &watchonly}, v.ws)
	lt, err := r.cx.WalletServer.ListTransactions(0, 11)
	if err != nil {
		log.L.Info(err)
	}
	log.L.Debug("TRZNZA:")
	log.L.Debug(lt)

	r.Status.Wallet.Transactions.TxsNumber = len(lt)
	log.L.Info("ETO:" + fmt.Sprint(lt))
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
	log.L.Debug("getting latest transactions")
		lt, err := r.cx.WalletServer.ListTransactions(0, 10)
		if err != nil {
			// //r.PushDuoUIalert("Error", err.Error(), "error")
		}
		r.History.Txs.TxsNumber = len(lt)
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
			if r.History.Txs.Balance > balanceHeight {
				balanceHeight = r.History.Txs.Balance
			}

		}
		r.History.Txs.Txs = txs
		r.History.Txs.BalanceHeight = balanceHeight

}

func (r *RcVar) GetTransactions() func() {
	return func() {
		log.L.Debug("getting transactions")
		lt, err := r.cx.WalletServer.ListTransactions(0, r.History.Txs.TxsListNumber)
		if err != nil {
			// //r.PushDuoUIalert("Error", err.Error(), "error")
		}
		r.History.Txs.TxsNumber = len(lt)
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
			if r.History.Txs.Balance > balanceHeight {
				balanceHeight = r.History.Txs.Balance
			}

		}
		r.History.Txs.Txs = txs
		r.History.Txs.BalanceHeight = balanceHeight
	}
}
