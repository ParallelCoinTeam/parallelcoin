package duoui

import (
	"github.com/p9c/pod/cmd/gui/mvc/theme"
)

func (ui *DuoUI) LoadPages() (p map[string]*theme.DuoUIpage) {
	p = make(map[string]*theme.DuoUIpage)
	//p := *new(*parallel.DuoUIpage)

	p["OVERVIEW"] = ui.ly.Theme.DuoUIpage("OVERVIEW", "ffcf30cf", "ffcf3030", 10, 10, ui.DuoUIoverview())
	p["SEND"] = ui.ly.Theme.DuoUIpage("SEND", "ffcf30cf", "ffcf3030", 10, 10, ui.DuoUIsend())
	p["RECEIVE"] = ui.ly.Theme.DuoUIpage("RECEIVE", "ffcf30cf", "ffcf3030", 10, 10, func() { ui.ly.Theme.H5("receive :").Layout(ui.ly.Context) })
	p["ADDRESSBOOK"] = ui.ly.Theme.DuoUIpage("ADDRESSBOOK", "ffcf30cf", "ffcf3030", 10, 10, ui.DuoUIaddressBook())
	p["SETTINGS"] = ui.ly.Theme.DuoUIpage("SETTINGS", "ffcf30cf", "ffcf3030", 10, 10, ui.DuoUIsettings())
	p["NETWORK"] = ui.ly.Theme.DuoUIpage("NETWORK", "ffcf30cf", "ffcf3030", 10, 10, func() { ui.ly.Theme.H5("network :").Layout(ui.ly.Context) })
	p["HISTORY"] = ui.ly.Theme.DuoUIpage("HISTORY", "ffcf30cf", "ffcf3030", 10, 10, ui.DuoUItransactions())
	p["EXPLORER"] = ui.ly.Theme.DuoUIpage("EXPLORER", "ffcf30cf", "ffcf3030", 10, 10, ui.DuoUIexplorer())
	p["CONSOLE"] = ui.ly.Theme.DuoUIpage("CONSOLE", "ffcf30cf", "ffcf3030", 10, 10, ui.DuoUIconsole())
	p["LOG"] = ui.ly.Theme.DuoUIpage("LOG", "ffcf30cf", "ffcf3030", 10, 10, ui.DuoUIlogger())
	return
}
