package rcd

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	log "github.com/p9c/pod/pkg/logi"

	"github.com/p9c/pod/cmd/node/rpc"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	"github.com/p9c/pod/pkg/rpc/legacy"
)

func (r *RcVar) ConsoleCmd(com string) (o string) {
	split := strings.Split(com, " ")
	method := split[0]
	args := split[1:]
	var cmd, res interface{}
	var err error
	var errString, prev string
	if method == "help" {
		if len(args) < 1 {
			method = ""
			cmd = &btcjson.HelpCmd{Command: &method}
			if res, err = rpc.RPCHandlers["help"].Fn(r.cx.RPCServer, cmd, nil); log.L.Check(err) {
				errString += fmt.Sprintln(err)
			}
			o += fmt.Sprintln(res)
			if res, err = legacy.RPCHandlers["help"].
				Handler(cmd, r.cx.WalletServer, r.cx.ChainClient); log.L.Check(err) {
				errString += fmt.Sprintln(err)
			}
			o += fmt.Sprintln(res)
			splitted := strings.Split(o, "\n")
			sort.Strings(splitted)
			var dedup []string
			for i := range splitted {
				if i > 0 {
					if splitted[i] != prev {
						dedup = append(dedup, splitted[i])
					}
				}
				prev = splitted[i]
			}
			o = strings.Join(dedup, "\n")
			if errString != "" {
				o += "Error:\n"
				o += errString
			}
		} else {
			method = args[0]
			log.L.Debug("finding help for command", method)
			if help, err := r.cx.RPCServer.HelpCacher.RPCMethodHelp(
				method); log.L.Check(err) {
				o += err.Error() + "\n"
				o += fmt.Sprintln(res)
				cmd = &btcjson.HelpCmd{Command: &method}
				if res, err = legacy.RPCHandlers["help"].
					Handler(cmd, r.cx.WalletServer, r.cx.ChainClient); log.L.Check(err) {
					errString += fmt.Sprintln(err)
				}
				o += fmt.Sprintln(res)
			} else {
				o += help
			}
			// if _, ok := legacy.RPCHandlers[method]; ok {
			// 	o += "wallet server:\n"
			// 	o += legacy.HelpDescsEnUS()[method]
			// }
			// if _, ok := rpc.RPCHandlers[method]; ok {
			// 	o += "chain server:\n"
			// 	o += rpc.HelpDescsEnUS[method]
			// }
		}
		return
	}
	params := make([]interface{}, 0, len(split[1:]))
	for _, arg := range args {
		params = append(params, arg)
	}
	if cmd, err = btcjson.NewCmd(method, params...); log.L.Check(err) {
		o += fmt.Sprintln(err)
	}
	if x, ok := rpc.RPCHandlers[method]; !ok {
		if x, ok := legacy.RPCHandlers[method]; ok {
			if res, err = x.Handler(cmd, r.cx.WalletServer,
				r.cx.ChainClient); log.L.Check(err) {
				o += err.Error()
			}
			// o += fmt.Sprintln(res)
		}
	} else {
		if res, err = x.Fn(r.cx.RPCServer, cmd, nil); log.L.Check(err) {
			o += err.Error()
		}
		// o += fmt.Sprintln(res)
	}
	if res != nil {
		if j, err := json.MarshalIndent(res, "",
			"  "); !log.L.Check(err) {
			o += string(j)
		}
	}
	return

}
