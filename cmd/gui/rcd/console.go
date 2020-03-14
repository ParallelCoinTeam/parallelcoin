package rcd

import (
	"fmt"
	"strings"

	log "github.com/p9c/logi"

	"github.com/p9c/pod/cmd/node/rpc"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	"github.com/p9c/pod/pkg/rpc/legacy"
)

func (r *RcVar) ConsoleCmd(com string) (o string) {
	split := strings.Split(com, " ")
	params := make([]interface{}, 0, len(split[1:]))
	cmd, err := btcjson.NewCmd(split[0], params...)
	if err != nil {
		o = fmt.Sprint(err)
	}
	log.L.Info(split)
	log.L.Infos(cmd)
	if x, ok := rpc.RPCHandlers[split[0]]; ok {
		if res, err := x.Fn(r.cx.RPCServer, cmd, nil); log.L.Check(err) {

		} else {
			return fmt.Sprint(res)
		}
	} else if x, ok := legacy.RPCHandlers[split[0]]; ok {
		_ = x
		// if res, err := x.Handler(cmd, r.cx.RPCServer,, nil); log.L.Check(err) {
		//
		// } else {
		// 	return fmt.Sprint(res)
		// }
	}

	return
}
