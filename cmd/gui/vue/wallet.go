//+build !nogui
// +build !headless

package vue

import (
	"encoding/hex"
	"fmt"
	"github.com/p9c/pod/cmd/gui/vue/alert"
	"github.com/p9c/pod/cmd/gui/vue/mod"
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
	"github.com/minio/highwayhash"
	"time"
)

func (d *DuoVUE) GetBalance() DuoVUEbalance {
	acct := "default"
	minconf := 0
	getBalance, err := legacy.GetBalance(&btcjson.GetBalanceCmd{Account: &acct, MinConf: &minconf}, d.cx.WalletServer)
	if err != nil {
		alert.Alert.Time = time.Now()
		alert.Alert.Alert = err.Error()
		alert.Alert.AlertType = "error"
	}
	gb, ok := getBalance.(float64)
	if ok {
		bb := fmt.Sprintf("%0.8f", gb)
		d.Status.Balance.Balance = bb
	}

	getUnconfirmedBalance, err := legacy.GetUnconfirmedBalance(&btcjson.GetUnconfirmedBalanceCmd{Account: &acct}, d.cx.WalletServer)
	if err != nil {
		alert.Alert.Time = time.Now()
		alert.Alert.Alert = err.Error()
		alert.Alert.AlertType = "error"
	}
	ub, ok := getUnconfirmedBalance.(float64)
	if ok {
		ubb := fmt.Sprintf("%0.8f", ub)
		d.Status.Balance.Unconfirmed = ubb
	}
	return d.Status.Balance
}
func (d *DuoVUE) GetTransactions(from, count int, cat string) (txs DuoVUEtransactions) {
	// account, txcount, startnum, watchonly := "*", n, f, false
	// listTransactions, err := legacy.ListTransactions(&json.ListTransactionsCmd{Account: &account, Count: &txcount, From: &startnum, IncludeWatchOnly: &watchonly}, v.ws)
	lt, err := d.cx.WalletServer.ListTransactions(0, 10)
	if err != nil {
		alert.Alert.Time = time.Now()
		alert.Alert.Alert = err.Error()
		alert.Alert.AlertType = "error"
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

func (d *DuoVUE) GetAddressBook() (addressbook DuoVUEAddressBook) {
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
	syncBlock := d.cx.WalletServer.Manager.SyncedTo()
	// Intermediate data for all addresses.
	allAddrData := make(map[string]AddrData)

	// Create an AddrData entry for each active address in the account.
	// Otherwise we'll just get addresses from transactions later.
	sortedAddrs, err := d.cx.WalletServer.SortedActivePaymentAddresses()
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
	err = wallet.ExposeUnstableAPI(d.cx.WalletServer).RangeTransactions(0, endHeight, func(details []wtxmgr.TxDetails) (bool, error) {
		for _, tx := range details {
			for _, cred := range tx.Credits {
				pkScript := tx.MsgTx.TxOut[cred.Index].PkScript
				_, addrs, _, err := txscript.ExtractPkScriptAddrs(
					pkScript, d.cx.WalletServer.ChainParams())
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
	return
}

func (d *DuoVUE) DuoSend(wp string, ad string, am float64) string {
	if am > 0 {
		getBlockChain, err := rpc.HandleGetBlockChainInfo(d.cx.RPCServer, nil, nil)
		if err != nil {
			alert.Alert.Time = time.Now()
			alert.Alert.Alert = err.Error()
			alert.Alert.AlertType = "error"
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
			alert.Alert.Time = time.Now()
			alert.Alert.Alert = err.Error()
			alert.Alert.AlertType = "error"
		}
		var validateAddr *btcjson.ValidateAddressWalletResult
		if err == nil {
			var va interface{}
			va, err = legacy.ValidateAddress(&btcjson.ValidateAddressCmd{Address: addr.String()}, d.cx.WalletServer)
			if err != nil {
				alert.Alert.Time = time.Now()
				alert.Alert.Alert = err.Error()
				alert.Alert.AlertType = "error"
			}
			vva := va.(btcjson.ValidateAddressWalletResult)
			validateAddr = &vva
			if validateAddr.IsValid {
				legacy.WalletPassphrase(btcjson.NewWalletPassphraseCmd(wp, 5), d.cx.WalletServer)
				if err != nil {
					alert.Alert.Time = time.Now()
					alert.Alert.Alert = err.Error()
					alert.Alert.AlertType = "error"
				}

				_, err = legacy.SendToAddress(
					&btcjson.SendToAddressCmd{
						Address: addr.EncodeAddress(), Amount: amount.ToDUO(),
					}, d.cx.WalletServer)
				if err != nil {
					alert.Alert.Time = time.Now()
					alert.Alert.Alert = "error sending to address:" + err.Error()
					alert.Alert.AlertType = "error"
				} else {
					alert.Alert.Time = time.Now()
					alert.Alert.Alert = "Address OK"
					alert.Alert.AlertType = "success"
				}
			} else {
				if err != nil {
					alert.Alert.Time = time.Now()
					alert.Alert.Alert = "Invalid address"
					alert.Alert.AlertType = "error"
				}
			}
			alert.Alert.Time = time.Now()
			alert.Alert.Alert = "Payment sent"
			alert.Alert.AlertType = "success"
		}
	} else {
		// fmt.Println("low")

	}
	return "sent"
}

func (d *DuoVUE) CreateNewAddress(acctName, label string) string {
	account, err := d.cx.WalletServer.AccountNumber(waddrmgr.KeyScopeBIP0044, acctName)
	if err != nil {
	}
	addr, err := d.cx.WalletServer.NewAddress(account,
		waddrmgr.KeyScopeBIP0044, false)
	if err != nil {
	}
	hf, err := highwayhash.New64(make([]byte, 32))
	if err != nil {
		panic(err)
	}
	addressHash := hex.EncodeToString(hf.Sum(addr.ScriptAddress()))
	d.db.DbWrite("addressbook", addressHash, mod.AddBook{
		Address: addressHash,
		Label:   label,
	})
	//fmt.Println("amo", w.Receive.Addr)

	fmt.Println("dddddddddddddddddddddddddddddddddd", addr.EncodeAddress())
	return addr.EncodeAddress()
}
