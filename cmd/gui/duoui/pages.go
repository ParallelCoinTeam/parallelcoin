package duoui

import (
	"github.com/p9c/pod/cmd/gui/mvc/theme"
)

func (ui *DuoUI) LoadPages() (p map[string]*theme.DuoUIpage) {
	p = make(map[string]*theme.DuoUIpage)
	//p := *new(*parallel.DuoUIpage)

	p["OVERVIEW"] = ui.ly.Theme.DuoUIpage("OVERVIEW", func() {}, ui.overviewBody(), func() {})
	p["SEND"] = ui.ly.Theme.DuoUIpage("SEND", func() {}, ui.DuoUIsend(), func() {})
	p["RECEIVE"] = ui.ly.Theme.DuoUIpage("RECEIVE", func() {}, func() { ui.ly.Theme.H5("receive :").Layout(ui.ly.Context) }, func() {})
	p["ADDRESSBOOK"] = ui.ly.Theme.DuoUIpage("ADDRESSBOOK", func() {}, ui.DuoUIaddressBook(), func() {})
	p["SETTINGS"] = ui.ly.Theme.DuoUIpage("SETTINGS", ui.contentHeader(ui.headerSettings()), ui.settingsBody(), func() {})
	p["NETWORK"] = ui.ly.Theme.DuoUIpage("NETWORK", func() {}, func() { ui.ly.Theme.H5("network :").Layout(ui.ly.Context) }, func() {})
	p["HISTORY"] = ui.ly.Theme.DuoUIpage("HISTORY", ui.contentHeader(ui.headerTransactions()), ui.txsBody(), func() {})
	p["EXPLORER"] = ui.ly.Theme.DuoUIpage("EXPLORER", ui.contentHeader(ui.headerExplorer()), ui.bodyExplorer(), func() {})
	p["MINER"] = ui.ly.Theme.DuoUIpage("MINER", func() {}, ui.DuoUIminer(), func() {})
	p["CONSOLE"] = ui.ly.Theme.DuoUIpage("CONSOLE", func() {}, ui.DuoUIconsole(), func() {})
	p["LOG"] = ui.ly.Theme.DuoUIpage("LOG", func() {}, ui.DuoUIlogger(), func() {})
	return
}
