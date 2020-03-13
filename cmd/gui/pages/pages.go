package pages

import (
	"gioui.org/layout"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gelook"
)

func LoadPages(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) (p map[string]*gelook.DuoUIpage) {
	p = make(map[string]*gelook.DuoUIpage)

	p["OVERVIEW"] = Overview(rc, gtx, th)
	p["SEND"] = Send(rc, gtx, th)
	p["RECEIVE"] = th.DuoUIpage("RECEIVE", 10, func() {}, func() {}, func() { th.H5("receive :").Layout(gtx) }, func() {})
	p["ADDRESSBOOK"] = DuoUIaddressBook(rc, gtx, th)
	p["SETTINGS"] = Settings(rc, gtx, th)
	p["NETWORK"] = Network(rc, gtx, th)
	p["BLOCK"] = th.DuoUIpage("BLOCK", 0, func() {}, func() {}, func() { th.H5("block :").Layout(gtx) }, func() {})
	p["HISTORY"] = History(rc, gtx, th)
	p["EXPLORER"] = DuoUIexplorer(rc, gtx, th)
	p["MINER"] = Miner(rc, gtx, th)
	p["CONSOLE"] = Console(rc, gtx, th)
	p["LOG"] = Logger(rc, gtx, th)
	return
}
