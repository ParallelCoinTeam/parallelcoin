// SPDX-License-Identifier: Unlicense OR MIT
package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"gioui.org/widget/material"
	//"github.com/btcsuite/btcd/btcjson"

	txscript "github.com/p9c/pod/pkg/chain/tx/script"
	"github.com/p9c/pod/pkg/util"

	chaincfg "github.com/p9c/pod/pkg/chain/config"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	rpcclient "github.com/p9c/pod/pkg/rpc/client"

	"golang.org/x/exp/shiny/materialdesign/icons"
)

type Wallet struct {
	addr    util.Address
	disc    chan struct{}
	done    chan struct{}
	events  chan Event
	balance util.Amount
}

type Transaction struct {
	New    bool
	Hash   string
	amount util.Amount
	time   time.Time
	added  time.Time
}

type Event interface {
	isWalletEvent()
}

type TransactionEvent struct {
	Trans Transaction
}

type ErrorEvent struct {
	Err error
}

var qrIcn *material.Icon

func init() {
	icn, err := material.NewIcon(icons.ContentAdd)
	if err != nil {
		log.Fatal(err)
	}
	qrIcn = icn
}

var (
	pubAddr = flag.String("addr", "tb1qy9cem33xwpttzdh7a3nsmqsyz8ytz2jz28w860", "bitcoin address")
	host    = flag.String("host", "localhost:18334", "btcd host")
)

const startHeight = 1582962

var bitcoinNet = &chaincfg.TestNet3Params

const (
	rpcUser = "Mk1Xfdws+n0OI6fARguxpLtZt48="
	rpcPass = "zZT4QgeqwADcEKcrkOJr/2x7760="
	rpcCert = `-----BEGIN CERTIFICATE-----
MIICmTCCAfugAwIBAgIRAJwxVLc1YHY3A35LaTAsS3MwCgYIKoZIzj0EAwQwNDEg
MB4GA1UEChMXYnRjZCBhdXRvZ2VuZXJhdGVkIGNlcnQxEDAOBgNVBAMTB3Rlc3Rt
YWMwHhcNMTkxMDE2MjAyMzU0WhcNMjkxMDE0MjAyMzU0WjA0MSAwHgYDVQQKExdi
dGNkIGF1dG9nZW5lcmF0ZWQgY2VydDEQMA4GA1UEAxMHdGVzdG1hYzCBmzAQBgcq
hkjOPQIBBgUrgQQAIwOBhgAEASBx1IWFOr6/y82v6nJHYB7tGHjHk7WpbEHxqxi2
raovxw4aM2d/gncuPNinMInP6JbRdvV30CYZ5/GrimZjuNRlARihFHbYQ6lkYLAy
wncw4Y7rFOWbmbj9YsFOgnUkhuTkBm3r56UHfeqaO7LbG+zVZo2/mBewkfkhFyEi
SXmpQPXmo4GqMIGnMA4GA1UdDwEB/wQEAwICpDAPBgNVHRMBAf8EBTADAQH/MIGD
BgNVHREEfDB6ggd0ZXN0bWFjgglsb2NhbGhvc3SHBH8AAAGHEAAAAAAAAAAAAAAA
AAAAAAGHEP6AAAAAAAAAAAAAAAAAAAGHEP6AAAAAAAAABMMVBQWbMICHBMCoVraH
EP6AAAAAAAAAhHcl//6k7zSHEP6AAAAAAAAA27SKa9kr6WkwCgYIKoZIzj0EAwQD
gYsAMIGHAkIBrl769uPROZS+KHnBS40kIIdGgwHcA88I8jv4udTphdGRrCOKNFsB
FAu7fv/4YmTAafgmUX0s66Jmg81tKGQRnbUCQQxvDfJYIdPG4nO3piAoihZY7g6U
+zUCmQjcegjk0jY0UpBN3ire+iNEkFMLm2CDtTGTgHZJH7DinyGMgqXb0Kc5
-----END CERTIFICATE-----`
)

func NewWallet(pubaddr string, host string) (*Wallet, error) {
	addr, err := util.DecodeAddress(pubaddr, bitcoinNet)
	if err != nil {
		return nil, err
	}
	w := &Wallet{
		addr: addr,
	}
	ccfg := &rpcclient.ConnConfig{
		Host:         host,
		Endpoint:     "ws",
		User:         rpcUser,
		Pass:         rpcPass,
		Certificates: []byte(rpcCert),
	}
	w.Connect(ccfg)
	return w, nil
}

func (w *Wallet) Balance() util.Amount {
	return util.Amount(atomic.LoadInt64((*int64)(&w.balance)))
}

func (w *Wallet) Connect(cfg *rpcclient.ConnConfig) {
	w.disc = make(chan struct{})
	w.done = make(chan struct{})
	w.events = make(chan Event)
	go w.run(cfg)
}

func (w *Wallet) run(cfg *rpcclient.ConnConfig) {
	var cl *rpcclient.Client
	onerr := func(err error) {
		if w.disc == nil {
			return
		}
		log.Printf("wallet error: %v", err)
		w.events <- ErrorEvent{err}
	}
	rescanned := false
	connected := make(chan struct{})
	handlers := &rpcclient.NotificationHandlers{
		OnClientConnected: func() {
			connected <- struct{}{}
		},
		OnRescanFinished: func(hash *chainhash.Hash, height int32, blkTime time.Time) {
			rescanned = true
		},
		OnRecvTx: func(trans *util.Tx, details *btcjson.BlockDetails) {
			for _, out := range trans.MsgTx().TxOut {
				// Extract and print details from the script.
				_, addresses, _, err := txscript.ExtractPkScriptAddrs(out.PkScript, bitcoinNet)
				if err != nil {
					onerr(err)
					return
				}
				if len(addresses) != 1 {
					continue
				}
				addr := addresses[0]
				if !bytes.Equal(addr.ScriptAddress(), w.addr.ScriptAddress()) {
					continue
				}
				atomic.AddInt64((*int64)(&w.balance), out.Value)
				var t time.Time
				if details != nil {
					t = time.Unix(details.Time, 0)
				} else {
					t = time.Now()
				}
				w.events <- TransactionEvent{
					Trans: Transaction{
						New:    rescanned,
						Hash:   trans.Hash().String(),
						amount: util.Amount(out.Value),
						time:   t,
						added:  time.Now(),
					},
				}
			}
		},
	}
	var err error
	cl, err = rpcclient.New(cfg, handlers)
	if err != nil {
		onerr(err)
	}
	go func() {
		<-connected
		block, err := cl.GetBlockHash(startHeight)
		if err != nil {
			onerr(err)
			return
		}
		addrs := []util.Address{w.addr}
		if err := cl.NotifyReceived(addrs); err != nil {
			log.Fatal(err)
		}
		if err := cl.Rescan(block, addrs, nil); err != nil {
			onerr(err)
			return
		}
	}()
	for {
		select {
		case <-w.disc:
			cl.Shutdown()
			w.done <- struct{}{}
			return
		}
	}
}

func (w *Wallet) disconnect() {
	w.disc <- struct{}{}
	<-w.done
}

func (t *Transaction) FormatAmount() string {
	return fmt.Sprintf("%d", t.amount)
	//return fmt.Sprintf("%d", t.amount)
}

func (t *Transaction) FormatTime() string {
	return t.time.Local().Format("2006-01-02 15:04:05")
}

func (e TransactionEvent) isWalletEvent() {}
func (e ErrorEvent) isWalletEvent()       {}
