package rcd

import (
	"encoding/hex"
	"fmt"
	"github.com/minio/highwayhash"
	"github.com/p9c/pod/cmd/gui/controller"
	"github.com/p9c/pod/cmd/gui/model"
	wtxmgr "github.com/p9c/pod/pkg/chain/tx/mgr"
	txscript "github.com/p9c/pod/pkg/chain/tx/script"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	"github.com/p9c/pod/pkg/util"
	"github.com/p9c/pod/pkg/wallet"
	waddrmgr "github.com/p9c/pod/pkg/wallet/addrmgr"
	"sort"
)

type DuoUItemplates struct {
	App  map[string][]byte            `json:"app"`
	Data map[string]map[string][]byte `json:"data"`
}

type DbAddress string

type Address struct {
	Index   int     `json:"num"`
	Label   string  `json:"label"`
	Account string  `json:"account"`
	Address string  `json:"address"`
	Amount  float64 `json:"amount"`
}
type Send struct {
	// Phrase string  `json:"phrase"`
	// Addr   string  `json:"addr"`
	// Amount float64 `json:"amount"`
	//exit
}

type AddBook struct {
	Address string `json:"address"`
	Label   string `json:"label"`
}

////////////////////////

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
			Copy:    new(controller.Button),
		})
	}
	sort.Sort(addrs)
	addressbook.Addresses = addrs
	r.AddressBook = *addressbook
	return
}
