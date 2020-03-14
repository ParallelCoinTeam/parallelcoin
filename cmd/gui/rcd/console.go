package rcd

import (
	"fmt"
	"strings"

	log "github.com/p9c/logi"

	"github.com/p9c/pod/pkg/rpc/btcjson"
)

func (r *RcVar) ConsoleCmd(com string) (o string) {
	split := strings.Split(com, " ")
	params := make([]interface{}, 0, len(split[1:]))
	cmd, err := btcjson.NewCmd(split[0], params...)
	if err != nil {
		o = fmt.Sprint(err)
	}
	log.L.Infos(cmd)
	return
}
