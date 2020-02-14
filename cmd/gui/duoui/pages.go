package duoui

import (
	"github.com/p9c/pod/cmd/gui/mvc/theme"
)

func (ui *DuoUI) LoadPages() (p map[string]*theme.DuoUIpage) {
	p = make(map[string]*theme.DuoUIpage)
	//p := *new(*parallel.DuoUIpage)

	p["SEND"] = ui.ly.Theme.DuoUIpage("DuoUIsend", "ffcf30cf", "ffcf3030", 10, 10, func() { ui.ly.Theme.H5("send :").Layout(ui.ly.Context) })
	p["RECEIVE"] = ui.ly.Theme.DuoUIpage("RECEIVE", "ffcf30cf", "ffcf3030", 10, 10, func() { ui.ly.Theme.H5("receive :").Layout(ui.ly.Context) })
	p["OVERVIEW"] = ui.ly.Theme.DuoUIpage("OVERVIEW", "ffcf30cf", "ffcf3030", 10, 10, ui.DuoUIoverview())
	p["ADDRESSBOOK"] = ui.ly.Theme.DuoUIpage("ADDRESSBOOK", "ffcf30cf", "ffcf3030", 10, 10, func() { ui.ly.Theme.H5("addressbook :").Layout(ui.ly.Context) })
	p["SETTINGS"] = ui.ly.Theme.DuoUIpage("SETTINGS", "ffcf30cf", "ffcf3030", 10, 10, func() { ui.DuoUIsettings() })
	p["NETWORK"] = ui.ly.Theme.DuoUIpage("NETWORK", "ffcf30cf", "ffcf3030", 10, 10, func() { ui.ly.Theme.H5("network :").Layout(ui.ly.Context) })
	p["HISTORY"] = ui.ly.Theme.DuoUIpage("HISTORY", "ffcf30cf", "ffcf3030", 10, 10, func() { ui.DuoUItransactions() })
	p["CONSOLE"] = ui.ly.Theme.DuoUIpage("CONSOLE", "ffcf30cf", "ffcf3030", 10, 10, ui.DuoUIconsole())
	p["LOG"] = ui.ly.Theme.DuoUIpage("LOG", "ffcf30cf", "ffcf3030", 10, 10, ui.DuoUIlogger())
	return
}
