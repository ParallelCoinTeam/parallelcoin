package rcd

import (
	"fmt"
	"time"

	"github.com/p9c/pod/pkg/rpc/btcjson"
	"github.com/p9c/pod/pkg/rpc/legacy"
)

func (r *RcVar) GetDuoUIbalance() {
	Trace("getting balance")
	acct := "default"
	minconf := 0
	getBalance, err := legacy.GetBalance(&btcjson.GetBalanceCmd{Account: &acct,
		MinConf: &minconf}, r.cx.WalletServer)
	if err != nil {
		// r.PushDuoUIalert("Error", err.Error(), "error")
	}
	gb, ok := getBalance.(float64)
	if ok {
		bb := fmt.Sprintf("%0.8f", gb)
		r.Status.Wallet.Balance.Store(bb)
	}
	return
}

func (r *RcVar) GetDuoUIunconfirmedBalance() {
	Trace("getting unconfirmed balance")
	acct := "default"
	getUnconfirmedBalance, err := legacy.GetUnconfirmedBalance(&btcjson.GetUnconfirmedBalanceCmd{Account: &acct}, r.cx.WalletServer)
	if err != nil {
		// r.PushDuoUIalert("Error", err.Error(), "error")
	}
	ub, ok := getUnconfirmedBalance.(float64)
	if ok {
		ubb := fmt.Sprintf("%0.8f", ub)
		r.Status.Wallet.Unconfirmed.Store(ubb)
	}
	return
}

func (r *RcVar) DuoSend(wp string, ad string, am float64) func() {
	return func() {
		pass := legacy.RPCHandlers["walletpassphrase"].Result()
		pass.WalletPassphrase(&btcjson.WalletPassphraseCmd{
			Passphrase: "aaa",
			Timeout:    int64(time.Second),
		})
		pass.WalletPassphraseWait(nil)
		send := legacy.RPCHandlers["sendtoaddress"].Result()
		send.SendToAddress(&btcjson.SendToAddressCmd{
			Address:   ad,
			Amount:    am,
			Comment:   nil,
			CommentTo: nil,
		})
		send.SendToAddressWait(nil)
	}
}
