package rcd

import (
	"encoding/hex"
	"fmt"
	"github.com/minio/highwayhash"
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"github.com/p9c/pod/cmd/gui/mvc/model"
	"github.com/p9c/pod/cmd/node/rpc"
	"github.com/p9c/pod/pkg/chain/config/netparams"
	wtxmgr "github.com/p9c/pod/pkg/chain/tx/mgr"
	txscript "github.com/p9c/pod/pkg/chain/tx/script"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	"github.com/p9c/pod/pkg/rpc/legacy"
	"github.com/p9c/pod/pkg/util"
	"github.com/p9c/pod/pkg/wallet"
	waddrmgr "github.com/p9c/pod/pkg/wallet/addrmgr"
	"sort"
	"time"
)

func
(r *RcVar) GetDuoUIbalance() {
	acct := "default"
	minconf := 0
	getBalance, err := legacy.GetBalance(&btcjson.GetBalanceCmd{Account: &acct,
		MinConf: &minconf}, r.cx.WalletServer)
	if err != nil {
		//r.PushDuoUIalert("Error", err.Error(), "error")
	}
	gb, ok := getBalance.(float64)
	if ok {
		bb := fmt.Sprintf("%0.8f", gb)
		r.Status.Wallet.Balance = bb
	}
	return
}
func
(r *RcVar) GetDuoUIunconfirmedBalance() {
	acct := "default"
	getUnconfirmedBalance, err := legacy.GetUnconfirmedBalance(&btcjson.GetUnconfirmedBalanceCmd{Account: &acct}, r.cx.WalletServer)
	if err != nil {
		//r.PushDuoUIalert("Error", err.Error(), "error")
	}
	ub, ok := getUnconfirmedBalance.(float64)
	if ok {
		ubb := fmt.Sprintf("%0.8f", ub)
		r.Status.Wallet.Unconfirmed = ubb
	}
	return
}

func
(r *RcVar) GetDuoUItransactionsNumber() {
	// account, txcount, startnum, watchonly := "*", n, f, false
	// listTransactions, err := legacy.ListTransactions(&json.ListTransactionsCmd{Account: &account, Count: &txcount, From: &startnum, IncludeWatchOnly: &watchonly}, v.ws)
	lt, err := r.cx.WalletServer.ListTransactions(0, 999999999)
	if err != nil {
		//r.PushDuoUIalert("Error", err.Error(), "error")
	}
	r.Status.Wallet.TxsNumber = len(lt)
	return
}

func
(r *RcVar) GetDuoUItransactions(sfrom, count int, cat string) *model.DuoUItransactions {
	// account, txcount, startnum, watchonly := "*", n, f, false
	// listTransactions, err := legacy.ListTransactions(&json.ListTransactionsCmd{Account: &account, Count: &txcount, From: &startnum, IncludeWatchOnly: &watchonly}, v.ws)
	lt, err := r.cx.WalletServer.ListTransactions(sfrom, count)
	if err != nil {
		//r.PushDuoUIalert("Error", err.Error(), "error")
	}
	r.Status.Wallet.Transactions.TxsNumber = len(lt)
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
func
txs(t btcjson.ListTransactionsResult) model.DuoUItx {
	return model.DuoUItx{
		TxID:     t.TxID,
		Amount:   t.Amount,
		Category: t.Category,
		Time:     helpers.FormatTime(time.Unix(t.Time, 0)),
	}

}
func
(r *RcVar) GetLatestTransactions() {
	r.Status.Wallet.LastTxs = r.GetDuoUItransactions(0, 10, "")
	return
}
func
(r *RcVar) GetTransactions() {
	lt, err := r.cx.WalletServer.ListTransactions(0, r.Status.Wallet.Txs.TxsListNumber)
	if err != nil {
		////r.PushDuoUIalert("Error", err.Error(), "error")
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

func
(r *RcVar) DuoSend(wp string, ad string, am float64) {
	if am > 0 {
		getBlockChain, err := rpc.HandleGetBlockChainInfo(r.cx.RPCServer, nil, nil)
		if err != nil {
			//r.PushDuoUIalert("Error", err.Error(), "error")

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
			//r.PushDuoUIalert("Error", err.Error(), "error")
		}
		var validateAddr *btcjson.ValidateAddressWalletResult
		if err == nil {
			var va interface{}
			va, err = legacy.ValidateAddress(&btcjson.
			ValidateAddressCmd{Address: addr.String()}, r.cx.WalletServer)
			if err != nil {
				//r.PushDuoUIalert("Error", err.Error(), "error")
			}
			vva := va.(btcjson.ValidateAddressWalletResult)
			validateAddr = &vva
			if validateAddr.IsValid {
				legacy.WalletPassphrase(btcjson.NewWalletPassphraseCmd(wp, 5),
					r.cx.WalletServer)
				if err != nil {
					//r.PushDuoUIalert("Error", err.Error(), "error")
				}
				_, err = legacy.SendToAddress(
					&btcjson.SendToAddressCmd{
						Address: addr.EncodeAddress(), Amount: amount.ToDUO(),
					}, r.cx.WalletServer)
				if err != nil {
					//r.PushDuoUIalert("error sending to address:", err.Error(), "error")

				} else {
					//r.PushDuoUIalert("Address OK", "OK", "success")
				}
			} else {
				if err != nil {
					//r.PushDuoUIalert("Invalid address", "INVALID", "error")
				}
			}
			//r.PushDuoUIalert("Payment sent", "PAYMENT", "success")
		}
	} else {
		log.Println(am)
	}
	r.Sent = true
	return
}


func (r *RcVar) CreateNewAddress(acctName string) string {
	account, err := r.cx.WalletServer.AccountNumber(waddrmgr.KeyScopeBIP0044, acctName)
	if err != nil {
	}
	addr, err := r.cx.WalletServer.NewAddress(account,
		waddrmgr.KeyScopeBIP0044, true)
	if err != nil {
	}
	//dv.PushDuoVUEalert("New address created:", addr.EncodeAddress(), "success")
	fmt.Println("low", addr.EncodeAddress())
	return addr.EncodeAddress()
}

func (r *RcVar) SaveAddressLabel(address, label string) {
	hf, err := highwayhash.New64(make([]byte, 32))
	if err != nil {
		panic(err)
	}
	addressHash := hex.EncodeToString(hf.Sum([]byte(address)))
	r.db.DbWrite("addressbook", addressHash, AddBook{
		Address: addressHash,
		Label:   label,
	})

}



type AddressSlice []model.Address

func (a AddressSlice) Len() int {
	return len(a)
}

func (a AddressSlice) Less(i, j int) bool {
	return a[i].Index < a[j].Index
}

func (a AddressSlice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}


func (r *RcVar) GetAddressBook() {
	addressbook := new(model.DuoUIaddressBook)
	minConf := 1
	// Intermediate data for each address.
	type AddrData struct {
		// Total amount received.
		amount util.Amount
		// tx     []string
		// Account which the address belongs to
		// account string
		index int
	}
	syncBlock := r.cx.WalletServer.Manager.SyncedTo()
	// Intermediate data for all addresses.
	allAddrData := make(map[string]AddrData)

	// Create an AddrData entry for each active address in the account.
	// Otherwise we'll just get addresses from transactions later.
	sortedAddrs, err := r.cx.WalletServer.SortedActivePaymentAddresses()
	if err != nil {
	}
	idx := 0
	for _, address := range sortedAddrs {
		// There might be duplicates, just overwrite them.
		allAddrData[address] = AddrData{
			index: idx,
		}
		idx++
	}
	var endHeight int32
	if minConf == 0 {
		endHeight = -1
	} else {
		endHeight = syncBlock.Height - int32(minConf) + 1
	}
	err = wallet.ExposeUnstableAPI(r.cx.WalletServer).RangeTransactions(0, endHeight, func(details []wtxmgr.TxDetails) (bool, error) {
		for _, tx := range details {
			for _, cred := range tx.Credits {
				pkScript := tx.MsgTx.TxOut[cred.Index].PkScript
				_, addrs, _, err := txscript.ExtractPkScriptAddrs(
					pkScript, r.cx.WalletServer.ChainParams())
				if err != nil {
					// Non standard script, skip.
					continue
				}
				for _, addr := range addrs {
					addrStr := addr.EncodeAddress()
					addrData, ok := allAddrData[addrStr]
					if ok {
						addrData.amount += cred.Amount
					} else {
						addrData = AddrData{
							amount: cred.Amount,
						}
					}
					allAddrData[addrStr] = addrData
				}
			}
		}
		return false, nil
	})
	if err != nil {
	}
	var addrs AddressSlice
	// Massage address data into output format.
	addressbook.Num = len(allAddrData)
	for address, addrData := range allAddrData {
		addr := btcjson.ListReceivedByAddressResult{
			Address: address,
			Amount:  addrData.amount.ToDUO(),
		}
		addrs = append(addrs, model.Address{
			Index:   addrData.index,
			Account: addr.Account,
			Address: addr.Address,
			Amount:  addr.Amount,
			Copy: new(controller.Button),
		})
	}
	sort.Sort(addrs)
	addressbook.Addresses = addrs
	r.AddressBook = *addressbook
	return
}
