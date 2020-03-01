package model

import (
	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"github.com/p9c/pod/pkg/rpc/btcjson"
)

type Explorer struct {
	Page        *controller.DuoUIcounter
	PerPage     *controller.DuoUIcounter
	Blocks      []DuoUIblock
	SingleBlock btcjson.GetBlockVerboseResult
}
