// +build !nogui
// +build !headless

package core

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/minio/highwayhash"

	"github.com/p9c/pod/cmd/gui/mod"
	"github.com/p9c/pod/cmd/node/rpc"
	"github.com/p9c/pod/pkg/chain/config/netparams"
	wtxmgr "github.com/p9c/pod/pkg/chain/tx/mgr"
	txscript "github.com/p9c/pod/pkg/chain/tx/script"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	"github.com/p9c/pod/pkg/rpc/legacy"
	"github.com/p9c/pod/pkg/util"
	btcutil "github.com/p9c/pod/pkg/util"
	"github.com/p9c/pod/pkg/wallet"
	waddrmgr "github.com/p9c/pod/pkg/wallet/addrmgr"
)

func (d *DuOS) GetBalance() mod.DuOSbalance {
	acct := "default"
	minconf := 0
	getBalance, err := legacy.GetBalance(&btcjson.GetBalanceCmd{Account: &acct,
		MinConf: &minconf}, d.Cx.WalletServer)
	if err != nil {
		d.PushDuOSalert("Error", err.Error(), "error")
	}
	gb, ok := getBalance.(float64)
	if ok {
		bb := fmt.Sprintf("%0.8f", gb)
		d.Data.Status.Balance.Balance = bb
	}

	getUnconfirmedBalance, err := legacy.GetUnconfirmedBalance(&btcjson.
		GetUnconfirmedBalanceCmd{Account: &acct}, d.Cx.WalletServer)
	if err != nil {
		d.PushDuOSalert("Error", err.Error(), "error")
	}
	ub, ok := getUnconfirmedBalance.(float64)
	if ok {
		ubb := fmt.Sprintf("%0.8f", ub)
		d.Data.Status.Balance.Unconfirmed = ubb
	}
	return d.Data.Status.Balance
}
func (d *DuOS) GetTransactions(from, count int, cat string) (txs mod.DuOStransactions) {
	// account, txcount, startnum, watchonly := "*", n, f, false
	// listTransactions, err := legacy.ListTransactions(&json.ListTransactionsCmd{Account: &account, Count: &txcount, From: &startnum, IncludeWatchOnly: &watchonly}, v.ws)
	lt, err := d.Cx.WalletServer.ListTransactions(0, 10)
	if err != nil {
		d.PushDuOSalert("Error", err.Error(), "error")
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

func (d *DuOS) GetTransactionsExcertps() (txse mod.DuOStransactionsExcerpts) {
	lt, err := d.Cx.WalletServer.ListTransactions(0, 99999)
	if err != nil {
		d.PushDuOSalert("Error", err.Error(), "error")
	}
	txse.TxsNumber = len(lt)

	// for i, j := 0, len(lt)-1; i < j; i, j = i+1, j-1 {
	//	lt[i], lt[j] = lt[j], lt[i]
	// }

	balanceHeight := 0.0
	txseRaw := []mod.TransactionExcerpt{}
	for _, txRaw := range lt {
		unixTimeUTC := time.Unix(txRaw.Time, 0) // gives unix time stamp in utc
		txseRaw = append(txseRaw, mod.TransactionExcerpt{
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

func (d *DuOS) GetAddressBook() mod.DuOSaddressBook {
	addressbook := new(mod.DuOSaddressBook)
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
	syncBlock := d.Cx.WalletServer.Manager.SyncedTo()
	// Intermediate data for all addresses.
	allAddrData := make(map[string]AddrData)

	// Create an AddrData entry for each active address in the account.
	// Otherwise we'll just get addresses from transactions later.
	sortedAddrs, err := d.Cx.WalletServer.SortedActivePaymentAddresses()
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
	err = wallet.ExposeUnstableAPI(d.Cx.WalletServer).RangeTransactions(0, endHeight, func(details []wtxmgr.TxDetails) (bool, error) {
		for _, tx := range details {
			for _, cred := range tx.Credits {
				pkScript := tx.MsgTx.TxOut[cred.Index].PkScript
				_, addrs, _, err := txscript.ExtractPkScriptAddrs(
					pkScript, d.Cx.WalletServer.ChainParams())
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
	var addrs []mod.Address
	// Massage address data into output format.
	addressbook.Num = len(allAddrData)
	for address, addrData := range allAddrData {
		addr := btcjson.ListReceivedByAddressResult{
			Address: address,
			Amount:  addrData.amount.ToDUO(),
		}
		addrs = append(addrs, mod.Address{
			Index:   addrData.index,
			Account: addr.Account,
			Address: addr.Address,
			Amount:  addr.Amount,
		})
	}
	addressbook.Addresses = addrs
	return *addressbook
}

func (d *DuOS) DuoSend(wp string, ad string, am float64) string {
	if am > 0 {
		getBlockChain, err := rpc.HandleGetBlockChainInfo(d.Cx.RPCServer, nil, nil)
		if err != nil {
			d.PushDuOSalert("Error", err.Error(), "error")

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

		amount, _ := btcutil.NewAmount(am)
		addr, err := btcutil.DecodeAddress(ad, defaultNet)
		if err != nil {
			d.PushDuOSalert("Error", err.Error(), "error")

		}
		var validateAddr *btcjson.ValidateAddressWalletResult
		if err == nil {
			var va interface{}
			va, err = legacy.ValidateAddress(&btcjson.
				ValidateAddressCmd{Address: addr.String()}, d.Cx.WalletServer)
			if err != nil {
				d.PushDuOSalert("Error", err.Error(), "error")

			}
			vva := va.(btcjson.ValidateAddressWalletResult)
			validateAddr = &vva
			if validateAddr.IsValid {
				legacy.WalletPassphrase(btcjson.NewWalletPassphraseCmd(wp, 5),
					d.Cx.WalletServer)
				if err != nil {
					d.PushDuOSalert("Error", err.Error(), "error")
				}
				_, err = legacy.SendToAddress(
					&btcjson.SendToAddressCmd{
						Address: addr.EncodeAddress(), Amount: amount.ToDUO(),
					}, d.Cx.WalletServer)
				if err != nil {
					d.PushDuOSalert("error sending to address:", err.Error(), "error")

				} else {
					d.PushDuOSalert("Address OK", "OK", "success")
				}
			} else {
				if err != nil {
					d.PushDuOSalert("Invalid address", "INVALID", "error")
				}
			}
			d.PushDuOSalert("Payment sent", "PAYMENT", "success")
		}
	} else {
		// fmt.Println("low")

	}
	return "sent"
}

func (d *DuOS) CreateNewAddress(acctName string) string {
	account, err := d.Cx.WalletServer.AccountNumber(waddrmgr.KeyScopeBIP0044, acctName)
	if err != nil {
	}
	addr, err := d.Cx.WalletServer.NewAddress(account,
		waddrmgr.KeyScopeBIP0044, true)
	if err != nil {
	}
	d.PushDuOSalert("New address created:", addr.EncodeAddress(), "success")
	fmt.Println("low", addr.EncodeAddress())
	return addr.EncodeAddress()
}

func (d *DuOS) SaveAddressLabel(address, label string) {
	hf, err := highwayhash.New64(make([]byte, 32))
	if err != nil {
		panic(err)
	}
	addressHash := hex.EncodeToString(hf.Sum([]byte(address)))
	d.db.DbWrite("addressbook", addressHash, mod.AddBook{
		Address: addressHash,
		Label:   label,
	})

}
