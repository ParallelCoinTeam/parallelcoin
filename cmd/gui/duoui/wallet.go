package duoui

import (
	"fmt"
	"github.com/p9c/pod/cmd/node/rpc"
	"github.com/p9c/pod/pkg/chain/config/netparams"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	"github.com/p9c/pod/pkg/rpc/legacy"
	"github.com/p9c/pod/pkg/util"
	"time"
)

type
	DuOSbalance struct {
		Balance string `json:"balance"`
	}
type
	DuOSunconfirmed struct {
		Unconfirmed string `json:"unconfirmed"`
	}
type
	DuOStransactions struct {
		Txs       []btcjson.ListTransactionsResult `json:"txs"`
		TxsNumber int                              `json:"txsnumber"`
	}
type
	DuOStransactionsNumber struct {
		TxsNumber int `json:"txsnumber"`
	}
type
	DuOStransactionsExcerpts struct {
		Txs           []DuOStransactionExcerpt `json:"txs"`
		TxsNumber     int                      `json:"txsnumber"`
		Balance       float64                  `json:"balance"`
		BalanceHeight float64                  `json:"balanceheight"`
	}
type
	DuOStransactionExcerpt struct {
		Balance       float64 `json:"balance"`
		Amount        float64 `json:"amount"`
		Category      string  `json:"category"`
		Confirmations int64   `json:"confirmations"`
		Time          string  `json:"time"`
		TxID          string  `json:"txid"`
		Comment       string  `json:"comment,omitempty"`
	}

func
(duo *DuoUI) GetDuOSbalance() {
	acct := "default"
	minconf := 0
	getBalance, err := legacy.GetBalance(&btcjson.GetBalanceCmd{Account: &acct,
		MinConf: &minconf}, duo.cx.WalletServer)
	if err != nil {
		duo.rc.PushDuOSalert("Error", err.Error(), "error")
	}
	gb, ok := getBalance.(float64)
	if ok {
		bb := fmt.Sprintf("%0.8f", gb)
		duo.rc.Balance = bb
	}
	return
}
func
(duo *DuoUI) GetDuOSunconfirmedBalance() {
	acct := "default"
	getUnconfirmedBalance, err := legacy.GetUnconfirmedBalance(&btcjson.GetUnconfirmedBalanceCmd{Account: &acct}, duo.cx.WalletServer)
	if err != nil {
		duo.rc.PushDuOSalert("Error", err.Error(), "error")
	}
	ub, ok := getUnconfirmedBalance.(float64)
	if ok {
		ubb := fmt.Sprintf("%0.8f", ub)
		duo.rc.Unconfirmed = ubb
	}
	return
}

func
(duo *DuoUI) GetDuOStransactionsNumber() {
	// account, txcount, startnum, watchonly := "*", n, f, false
	// listTransactions, err := legacy.ListTransactions(&json.ListTransactionsCmd{Account: &account, Count: &txcount, From: &startnum, IncludeWatchOnly: &watchonly}, v.ws)
	lt, err := duo.cx.WalletServer.ListTransactions(0, 999999999)
	if err != nil {
		duo.rc.PushDuOSalert("Error", err.Error(), "error")
	}
	duo.rc.TxsNumber = len(lt)
}

func
(duo *DuoUI) GetDuOStransactions(sfrom, count int, cat string) DuOStransactions {
	// account, txcount, startnum, watchonly := "*", n, f, false
	// listTransactions, err := legacy.ListTransactions(&json.ListTransactionsCmd{Account: &account, Count: &txcount, From: &startnum, IncludeWatchOnly: &watchonly}, v.ws)
	lt, err := duo.cx.WalletServer.ListTransactions(sfrom, count)
	if err != nil {
		duo.rc.PushDuOSalert("Error", err.Error(), "error")
	}
	duo.rc.Transactions.TxsNumber = len(lt)
	// lt := listTransactions.([]json.ListTransactionsResult)
	switch c := cat; c {
	case "received":
		for _, tx := range lt {
			if tx.Category == "received" {
				duo.rc.Transactions.Txs = append(duo.rc.Transactions.Txs, tx)
			}
		}
	case "sent":
		for _, tx := range lt {
			if tx.Category == "sent" {
				duo.rc.Transactions.Txs = append(duo.rc.Transactions.Txs, tx)
			}
		}
	case "immature":
		for _, tx := range lt {
			if tx.Category == "immature" {
				duo.rc.Transactions.Txs = append(duo.rc.Transactions.Txs, tx)
			}
		}
	case "generate":
		for _, tx := range lt {
			if tx.Category == "generate" {
				duo.rc.Transactions.Txs = append(duo.rc.Transactions.Txs, tx)
			}
		}
	default:
		duo.rc.Transactions.Txs = lt
	}
	return duo.rc.Transactions
}
func
(duo *DuoUI) GetDuOSlastTxs() {
	duo.rc.LastTxs = duo.GetDuOStransactions(0, 10, "")
	return
}
func
(duo *DuoUI) GetDuOSTransactionsExcertps() {
	lt, err := duo.cx.WalletServer.ListTransactions(0, 99999)
	if err != nil {
		duo.rc.PushDuOSalert("Error", err.Error(), "error")
	}
	duo.rc.Txs.TxsNumber = len(lt)
	// for i, j := 0, len(lt)-1; i < j; i, j = i+1, j-1 {
	//	lt[i], lt[j] = lt[j], lt[i]
	// }
	balanceHeight := 0.0
	txseRaw := []DuOStransactionExcerpt{}
	for _, txRaw := range lt {
		unixTimeUTC := time.Unix(txRaw.Time, 0) // gives unix time stamp in utc
		txseRaw = append(txseRaw, DuOStransactionExcerpt{
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
	for _, tx := range txseRaw {
		balance = balance + tx.Amount
		tx.Balance = balance
		duo.rc.Txs.Txs = append(duo.rc.Txs.Txs, tx)
		if duo.rc.Txs.Balance > balanceHeight {
			balanceHeight = duo.rc.Txs.Balance
		}
	}
	duo.rc.Txs.BalanceHeight = balanceHeight
	return
}

func
(duo *DuoUI) DuoSend(wp string, ad string, am float64) {
	if am > 0 {
		getBlockChain, err := rpc.HandleGetBlockChainInfo(duo.cx.RPCServer, nil, nil)
		if err != nil {
			duo.rc.PushDuOSalert("Error", err.Error(), "error")

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
			duo.rc.PushDuOSalert("Error", err.Error(), "error")
		}
		var validateAddr *btcjson.ValidateAddressWalletResult
		if err == nil {
			var va interface{}
			va, err = legacy.ValidateAddress(&btcjson.
			ValidateAddressCmd{Address: addr.String()}, duo.cx.WalletServer)
			if err != nil {
				duo.rc.PushDuOSalert("Error", err.Error(), "error")
			}
			vva := va.(btcjson.ValidateAddressWalletResult)
			validateAddr = &vva
			if validateAddr.IsValid {
				legacy.WalletPassphrase(btcjson.NewWalletPassphraseCmd(wp, 5),
					duo.cx.WalletServer)
				if err != nil {
					duo.rc.PushDuOSalert("Error", err.Error(), "error")
				}
				_, err = legacy.SendToAddress(
					&btcjson.SendToAddressCmd{
						Address: addr.EncodeAddress(), Amount: amount.ToDUO(),
					}, duo.cx.WalletServer)
				if err != nil {
					duo.rc.PushDuOSalert("error sending to address:", err.Error(), "error")

				} else {
					duo.rc.PushDuOSalert("Address OK", "OK", "success")
				}
			} else {
				if err != nil {
					duo.rc.PushDuOSalert("Invalid address", "INVALID", "error")
				}
			}
			duo.rc.PushDuOSalert("Payment sent", "PAYMENT", "success")
		}
	} else {
		log.Println(am)
	}
	duo.rc.Sent = true
	return
}
