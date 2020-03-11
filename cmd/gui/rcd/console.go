package rcd

import (
	"encoding/json"
	"fmt"
	"github.com/p9c/wallet/walletmain"
	log "github.com/p9c/logi"
	"github.com/p9c/rpc/btcjson"
	"github.com/p9c/rpc/legacy"
	"github.com/p9c/wallet/chain"
	"strings"
)

func (r *RcVar) ConsoleCmd(com string) (o string) {
	split := strings.Split(com, " ")
	params := make([]interface{}, 0, len(split[1:]))
	log.L.Info(len(params))
	c, err := btcjson.NewCmd(split[0], params...)
	if err != nil {
		o = fmt.Sprint(err)
	}
	handler, ok := legacy.RPCHandlers[split[0]]
	if ok {
		var out interface{}
		if handler.HandlerWithChain != nil {
			rpcC, err := chain.NewRPCClient(r.cx.ActiveNet, *r.cx.Config.RPCConnect,
				*r.cx.Config.Username, *r.cx.Config.Password, walletmain.ReadCAFile(r.cx.Config), !*r.cx.Config.TLS, 0)
			if err != nil {
				log.L.Error(err)
			}
			err = rpcC.Start()
			if err != nil {
				log.L.Error(
					"unable to open connection to consensus RPC server:", err)
			}
			out, err = handler.HandlerWithChain(
				c,
				r.cx.WalletServer,
				rpcC)
			log.L.Debug("HandlerWithChain")
		}
		if handler.Handler != nil {
			out, err = handler.Handler(
				c,
				r.cx.WalletServer)
			if err != nil {
				log.L.Error(
					"unable to open connection to consensus RPC server:", err)
			}
			log.L.Debug("Handler")
		}
		if err != nil {
			o = fmt.Sprint(err)
		} else {
			if split[0] == "help" {
				o = out.(string)
			} else {
				j, _ := json.MarshalIndent(out, "", "  ")
				o = fmt.Sprint(string(j))
			}
		}
	} else {
		if split[0] == "" {
			o = ""
		} else {
			o = "Command does not exist"
		}
	}
	return
}
