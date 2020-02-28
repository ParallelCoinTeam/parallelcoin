package duoui

import (
	"gioui.org/layout"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
	"github.com/p9c/pod/cmd/gui/rcd"
)

func LoadPages(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme) (p map[string]*theme.DuoUIpage) {
	p = make(map[string]*theme.DuoUIpage)

	p["OVERVIEW"] = th.DuoUIpage("OVERVIEW", 0, func() {}, func() {}, overviewBody(rc, gtx, th), func() {})
	p["SEND"] = th.DuoUIpage("SEND", 10, func() {}, func() {}, send(rc, gtx, th), func() {})
	p["RECEIVE"] = th.DuoUIpage("RECEIVE", 10, func() {}, func() {}, func() { th.H5("receive :").Layout(gtx) }, func() {})
	p["ADDRESSBOOK"] = th.DuoUIpage("ADDRESSBOOK", 10, func() {}, func() {}, addressBook(rc, gtx, th), func() {})
	p["SETTINGS"] = th.DuoUIpage("SETTINGS", 0, func() {}, contentHeader(gtx, th, headerSettings(rc, gtx, th)), settingsBody(rc, gtx, th), func() {})
	p["NETWORK"] = th.DuoUIpage("NETWORK", 0, func() {}, func() {}, func() { th.H5("network :").Layout(gtx) }, func() {})
	p["BLOCK"] = th.DuoUIpage("BLOCK", 0, func() {}, func() {}, func() { th.H5("block :").Layout(gtx) }, func() {})
	p["HISTORY"] = th.DuoUIpage("HISTORY", 0, func() {}, contentHeader(gtx, th, headerTransactions(rc, gtx, th)), txsBody(rc, gtx, th), func() {})
	p["EXPLORER"] = th.DuoUIpage("EXPLORER", 0, rc.GetBlocksExcerpts(page.Value, perPage.Value), contentHeader(gtx, th, headerExplorer(gtx, th)), bodyExplorer(rc, gtx, th), func() {})
	p["MINER"] = th.DuoUIpage("MINER", 0, func() {}, func() {}, DuoUIminer(rc, gtx, th), func() {})
	p["CONSOLE"] = th.DuoUIpage("CONSOLE", 0, func() {}, func() {}, console(rc, gtx, th), func() {})
	p["LOG"] = th.DuoUIpage("LOG", 0, func() {}, func() {}, DuoUIlogger(rc, gtx, th), func() {})
	return
}

func setPage(rc *rcd.RcVar, page *theme.DuoUIpage) {
	rc.CurrentPage = page
}
