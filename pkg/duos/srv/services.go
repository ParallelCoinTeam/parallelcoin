package srv

import (
	alert2 "github.com/p9c/pod/gui/____BEZI/test/pkg/alert"
	mod2 "github.com/p9c/pod/gui/____BEZI/test/pkg/duos/mod"
	stat2 "github.com/p9c/pod/gui/____BEZI/test/pkg/stat"
	"github.com/p9c/pod/pkg/rpc/btcjson"
)

type DuOSservices struct {
	Status stat2.DuOSstatus `json:"status"`
	Alert  alert2.DuOSalert `json:"alert"`
	Data   *DuOSdata        `json:"data"`
}

type DuOSdata struct {
	Status               stat2.DuOSstatus              `json:"status"`
	Peers                []*btcjson.GetPeerInfoResult  `json:"peers"`
	TransactionsExcerpts mod2.DuOStransactionsExcerpts `json:"txsex"`
	Blocks               mod2.DuOSblocks               `json:"blocks"`
	Send                 mod2.Send                     `json:"send"`
	Screens              map[string]string             `json:"screens"`
	Icons                map[string]string             `json:"ico"`
	//Addressbook          DuOSaddressBook              `json:"addressbook"`
}
