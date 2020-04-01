package rcd

// )
//
// type
// DuOShistory struct {
//	cx *conte.Xt
//	db DuoUIdb
//	txs model.DuOStransactionsExcerpts
// }
//
// func (d *DuOShistory)GetDuOShistory() {
//	lt, err := d.cx.WalletServer.ListTransactions(0, 99999)
//	if err != nil {
//		//d.PushDuOSalert("Error", err.Error(), "error")
//	}
//	d.txs.TxsNumber = len(lt)
//	// for i, j := 0, len(lt)-1; i < j; i, j = i+1, j-1 {
//	//	lt[i], lt[j] = lt[j], lt[i]
//	// }
//	balanceHeight := 0.0
//	txseRaw := []DuOStransactionExcerpt{}
//	for _, txRaw := range lt {
//		unixTimeUTC := time.Unix(txRaw.Time, 0) // gives unix time stamp in utc
//		txseRaw = append(txseRaw, DuOStransactionExcerpt{
//			// Balance:       txse.Balance + txRaw.Amount,
//			Comment:       txRaw.Comment,
//			Amount:        txRaw.Amount,
//			Category:      txRaw.Category,
//			Confirmations: txRaw.Confirmations,
//			Time:          unixTimeUTC.Format(time.RFC3339),
//			TxID:          txRaw.TxID,
//		})
//	}
//	var balance float64
//	for _, tx := range txseRaw {
//		balance = balance + tx.Amount
//		tx.Balance = balance
//		d.txs.Txs = append(d.txs.Txs, tx)
//		if d.txs.Balance > balanceHeight {
//			balanceHeight = d.txs.Balance
//		}
//	}
//	d.txs.BalanceHeight = balanceHeight
//	Info("HISTORY-test:VAR->", d.txs)
//
//	return
// }
