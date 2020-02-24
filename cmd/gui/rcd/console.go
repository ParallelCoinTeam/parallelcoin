package rcd

import (
	"encoding/json"
	"fmt"
	"github.com/p9c/pod/cmd/node/rpc"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	"strings"
)

func (r *RcVar) ConsoleCmd(com string) (o string) {
	split := strings.Split(com, " ")
	params := split[1:]
	log.INFO(len(params))
	//if len(params) < 1 {
	//	params = nil
	//}
	c, err := btcjson.NewCmd(split[0], strings.Join(params, " "))
	out, err := rpc.RPCHandlers[split[0]](r.cx.RPCServer, c, make(chan struct{}))
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
	return
}
