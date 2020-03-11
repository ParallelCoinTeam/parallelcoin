package duoui

import (
	"gioui.org/layout"
	"github.com/p9c/gelook"

	"github.com/p9c/pod/cmd/gui/pages"
	"github.com/p9c/pod/cmd/gui/rcd"
)

func LoadPages(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) (p map[string]*gelook.DuoUIpage) {
	p = make(map[string]*gelook.DuoUIpage)

	p["OVERVIEW"] = pages.Overview(rc, gtx, th)
	p["SEND"] = pages.Send(rc, gtx, th)
	p["RECEIVE"] = th.DuoUIpage("RECEIVE", 10, func() {}, func() {}, func() { th.H5("receive :").Layout(gtx) }, func() {})
	p["ADDRESSBOOK"] = pages.DuoUIaddressBook(rc, gtx, th)
	p["SETTINGS"] = pages.Settings(rc, gtx, th)
	p["NETWORK"] = th.DuoUIpage("NETWORK", 0, func() {}, func() {}, func() { th.H5("network :").Layout(gtx) }, func() {})
	p["BLOCK"] = th.DuoUIpage("BLOCK", 0, func() {}, func() {}, func() { th.H5("block :").Layout(gtx) }, func() {})
	p["HISTORY"] = pages.History(rc, gtx, th)
	p["EXPLORER"] = pages.DuoUIexplorer(rc, gtx, th)
	p["MINER"] = pages.Miner(rc, gtx, th)
	p["CONSOLE"] = pages.Console(rc, gtx, th)
	p["LOG"] = pages.Logger(rc, gtx, th)
	return
}
