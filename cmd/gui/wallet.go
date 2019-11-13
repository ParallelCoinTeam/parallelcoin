package gui

import (
	"fmt"
	"github.com/p9c/pod/cmd/node/rpc"
	"github.com/p9c/pod/pkg/chain/config/netparams"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	"github.com/p9c/pod/pkg/rpc/legacy"
	"github.com/p9c/pod/pkg/util"
	"time"
)

type DuOSbalance struct {
	Balance     string `json:"balance"`
	Unconfirmed string `json:"unconfirmed"`
}

type DuOStransactions struct {
	Txs       []btcjson.ListTransactionsResult `json:"txs"`
	TxsNumber int                              `json:"txsnumber"`
}

type DuOStransactionsExcerpts struct {
	Txs           []TransactionExcerpt `json:"txs"`
	TxsNumber     int                  `json:"txsnumber"`
	Balance       float64              `json:"balance"`
	BalanceHeight float64              `json:"balanceheight"`
}

type TransactionExcerpt struct {
	Balance       float64 `json:"balance"`
	Amount        float64 `json:"amount"`
	Category      string  `json:"category"`
	Confirmations int64   `json:"confirmations"`
	Time          string  `json:"time"`
	TxID          string  `json:"txid"`
	Comment       string  `json:"comment,omitempty"`
}

func (r *rcvar) GetBalance() (b DuOSbalance) {
	acct := "default"
	minconf := 0
	getBalance, err := legacy.GetBalance(&btcjson.GetBalanceCmd{Account: &acct,
		MinConf: &minconf}, r.cx.WalletServer)
	if err != nil {
		r.PushDuOSalert("Error", err.Error(), "error")
	}
	gb, ok := getBalance.(float64)
	if ok {
		bb := fmt.Sprintf("%0.8f", gb)
		b.Balance = bb
	}
	getUnconfirmedBalance, err := legacy.GetUnconfirmedBalance(&btcjson.
		GetUnconfirmedBalanceCmd{Account: &acct}, r.cx.WalletServer)
	if err != nil {
		r.PushDuOSalert("Error", err.Error(), "error")
	}
	ub, ok := getUnconfirmedBalance.(float64)
	if ok {
		ubb := fmt.Sprintf("%0.8f", ub)
		b.Unconfirmed = ubb
	}
	return
}

func (r *rcvar) GetTransactions(sfrom, count int, cat string) (txs DuOStransactions) {
	// account, txcount, startnum, watchonly := "*", n, f, false
	// listTransactions, err := legacy.ListTransactions(&json.ListTransactionsCmd{Account: &account, Count: &txcount, From: &startnum, IncludeWatchOnly: &watchonly}, v.ws)
	lt, err := r.cx.WalletServer.ListTransactions(0, 10)
	if err != nil {
		r.PushDuOSalert("Error", err.Error(), "error")
	}
	txs.TxsNumber = len(lt)
	// lt := listTransactions.([]json.ListTransactionsResult)
	switch c := cat; c {
	case "received":
		for _, tx := range lt {
			if tx.Category == "received" {
				txs.Txs = append(txs.Txs, tx)
			}
		}
	case "sent":
		for _, tx := range lt {
			if tx.Category == "sent" {
				txs.Txs = append(txs.Txs, tx)
			}
		}
	case "immature":
		for _, tx := range lt {
			if tx.Category == "immature" {
				txs.Txs = append(txs.Txs, tx)
			}
		}
	case "generate":
		for _, tx := range lt {
			if tx.Category == "generate" {
				txs.Txs = append(txs.Txs, tx)
			}
		}
	default:
		txs.Txs = lt
	}
	return
}

func (r *rcvar) GetTransactionsExcertps() (txse DuOStransactionsExcerpts) {
	lt, err := r.cx.WalletServer.ListTransactions(0, 99999)
	if err != nil {
		r.PushDuOSalert("Error", err.Error(), "error")
	}
	txse.TxsNumber = len(lt)
	// for i, j := 0, len(lt)-1; i < j; i, j = i+1, j-1 {
	//	lt[i], lt[j] = lt[j], lt[i]
	// }
	balanceHeight := 0.0
	txseRaw := []TransactionExcerpt{}
	for _, txRaw := range lt {
		unixTimeUTC := time.Unix(txRaw.Time, 0) // gives unix time stamp in utc
		txseRaw = append(txseRaw, TransactionExcerpt{
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
		txse.Txs = append(txse.Txs, tx)
		if txse.Balance > balanceHeight {
			balanceHeight = txse.Balance
		}
		fmt.Println("btititititmt", tx.Time)
		fmt.Println("bbbbbbbbb", tx.Amount)
	}
	fmt.Println("cccccccccccccccccccccccccccccccccccccccccccccc")
	fmt.Println("bbiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiii")
	fmt.Println("bbiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiii")
	fmt.Println("balanceHeightbalanceHeight", balanceHeight)
	fmt.Println("bbbbbbbbbbbbbbbbbbbbbbbbbbbb", txse.Balance)
	txse.BalanceHeight = balanceHeight
	return
}

func (r *rcvar) DuoSend(wp string, ad string, am float64) string {
	if am > 0 {
		getBlockChain, err := rpc.HandleGetBlockChainInfo(r.cx.RPCServer, nil, nil)
		if err != nil {
			r.PushDuOSalert("Error", err.Error(), "error")

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
			r.PushDuOSalert("Error", err.Error(), "error")
		}
		var validateAddr *btcjson.ValidateAddressWalletResult
		if err == nil {
			var va interface{}
			va, err = legacy.ValidateAddress(&btcjson.
				ValidateAddressCmd{Address: addr.String()}, r.cx.WalletServer)
			if err != nil {
				r.PushDuOSalert("Error", err.Error(), "error")
			}
			vva := va.(btcjson.ValidateAddressWalletResult)
			validateAddr = &vva
			if validateAddr.IsValid {
				legacy.WalletPassphrase(btcjson.NewWalletPassphraseCmd(wp, 5),
					r.cx.WalletServer)
				if err != nil {
					r.PushDuOSalert("Error", err.Error(), "error")
				}
				_, err = legacy.SendToAddress(
					&btcjson.SendToAddressCmd{
						Address: addr.EncodeAddress(), Amount: amount.ToDUO(),
					}, r.cx.WalletServer)
				if err != nil {
					r.PushDuOSalert("error sending to address:", err.Error(), "error")

				} else {
					r.PushDuOSalert("Address OK", "OK", "success")
				}
			} else {
				if err != nil {
					r.PushDuOSalert("Invalid address", "INVALID", "error")
				}
			}
			r.PushDuOSalert("Payment sent", "PAYMENT", "success")
		}
	} else {
		// fmt.Println("low")
	}
	return "sent"
}
