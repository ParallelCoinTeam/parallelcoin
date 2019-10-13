package srv

import (
	"github.com/p9c/pod/pkg/alert"
	"github.com/p9c/pod/pkg/duos/mod"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	"github.com/p9c/pod/pkg/stat"
)

type DuOSservices struct {
	Status stat.DuOSstatus `json:"status"`
	Alert  alert.DuOSalert `json:"alert"`
	Data   *DuOSdata       `json:"data"`
}

type DuOSdata struct {
	Status               stat.DuOSstatus              `json:"status"`
	Peers                []*btcjson.GetPeerInfoResult `json:"peers"`
	TransactionsExcerpts mod.DuOStransactionsExcerpts `json:"txsex"`
	Blocks               mod.DuOSblocks               `json:"blocks"`
	Send                 mod.Send                     `json:"send"`
	Screens              map[string]string            `json:"screens"`
	Icons                map[string]string            `json:"ico"`
	//Addressbook          DuOSaddressBook              `json:"addressbook"`
}
