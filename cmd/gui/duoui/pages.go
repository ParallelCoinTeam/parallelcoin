package duoui

import (
	"github.com/p9c/pod/cmd/gui/mvc/theme"
	"github.com/p9c/pod/pkg/gui/layout"
)

func (ui *DuoUI) LoadPages() (p map[string]*theme.DuoUIpage) {
	p = make(map[string]*theme.DuoUIpage)
	//p := *new(*parallel.DuoUIpage)

	p["SEND"] = ui.ly.Theme.DuoUIpage("DuoUIsend", "ffcf30cf", "ffcf3030", 10, 10, func() { ui.ly.Theme.H5("send :").Layout(ui.ly.Context) })
	p["RECEIVE"] = ui.ly.Theme.DuoUIpage("RECEIVE", "ffcf30cf", "ffcf3030", 10, 10, func() { ui.ly.Theme.H5("receive :").Layout(ui.ly.Context) })
	p["OVERVIEW"] = ui.ly.Theme.DuoUIpage("OVERVIEW", "ffcf30cf", "ffcf3030", 10, 10, func() {
		viewport := layout.Flex{Axis: layout.Horizontal}
		if ui.ly.Context.Constraints.Width.Max < 780 {
			viewport = layout.Flex{Axis: layout.Vertical}
		}
		viewport.Layout(ui.ly.Context,
			layout.Flexed(0.5, func() {
				cs := ui.ly.Context.Constraints
				theme.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, cs.Height.Max, ui.ly.Theme.Color.Light, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
				ui.DuoUIbalance()
			}),
			layout.Flexed(0.5, func() {
				cs := ui.ly.Context.Constraints
				theme.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, cs.Height.Max, "ff424242", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
				ui.DuoUIlatestTransactions()
			}),
		)
	})
	p["ADDRESSBOOK"] = ui.ly.Theme.DuoUIpage("ADDRESSBOOK", "ffcf30cf", "ffcf3030", 10, 10, func() { ui.ly.Theme.H5("addressbook :").Layout(ui.ly.Context) })
	p["SETTINGS"] = ui.ly.Theme.DuoUIpage("SETTINGS", "ffcf30cf", "ffcf3030", 10, 10, func() { ui.DuoUIsettings() })
	p["NETWORK"] = ui.ly.Theme.DuoUIpage("NETWORK", "ffcf30cf", "ffcf3030", 10, 10, func() { ui.ly.Theme.H5("network :").Layout(ui.ly.Context) })
	p["HISTORY"] = ui.ly.Theme.DuoUIpage("HISTORY", "ffcf30cf", "ffcf3030", 10, 10, func() { ui.DuoUItransactions() })
	p["CONSOLE"] = ui.ly.Theme.DuoUIpage("CONSOLE", "ffcf30cf", "ffcf3030", 10, 10, func() { ui.DuoUIconsole() })
	p["LOG"] = ui.ly.Theme.DuoUIpage("LOG", "ffcf30cf", "ffcf3030", 10, 10, func() { ui.DuoUIlogger() })
	return
}
