package duoui

import (
	"github.com/p9c/gio-parallel/layout"
	"github.com/p9c/gio-parallel/unit"
)

func DuoUImenu(duo *DuoUI) layout.FlexChild {
	return duo.comp.sidebar.l.Rigid(duo.gc, func() {

		in := layout.UniformInset(unit.Dp(0))

		overview := duo.comp.menu.l.Rigid(duo.gc, func() {
			in.Layout(duo.gc, func() {
				for duo.menu.overview.Clicked(duo.gc) {
					duo.menu.current = "overview"
				}
				b := duo.th.IconButton(duo.ico.Overview)
				b.Background = duo.menu.icoBackground
				b.Color = duo.menu.icoColor
				b.Padding = duo.menu.icoPadding
				b.Size = duo.menu.icoSize
				b.Layout(duo.gc, &duo.menu.overview)
			})
		})
		history := duo.comp.menu.l.Rigid(duo.gc, func() {
			in.Layout(duo.gc, func() {
				for duo.menu.history.Clicked(duo.gc) {
					duo.menu.current = "history"
				}
				b := duo.th.IconButton(duo.ico.History)
				b.Background = duo.menu.icoBackground
				b.Color = duo.menu.icoColor
				b.Padding = duo.menu.icoPadding
				b.Size = duo.menu.icoSize
				b.Layout(duo.gc, &duo.menu.history)
			})
		})
		addressbook := duo.comp.menu.l.Rigid(duo.gc, func() {
			in.Layout(duo.gc, func() {
				for duo.menu.addressbook.Clicked(duo.gc) {
					duo.menu.current = "addressbook"
				}
				b := duo.th.IconButton(duo.ico.AddressBook)
				b.Background = duo.menu.icoBackground
				b.Color = duo.menu.icoColor
				b.Padding = duo.menu.icoPadding
				b.Size = duo.menu.icoSize
				b.Layout(duo.gc, &duo.menu.addressbook)
			})
		})
		explorer := duo.comp.menu.l.Rigid(duo.gc, func() {
			in.Layout(duo.gc, func() {
				for duo.menu.explorer.Clicked(duo.gc) {
					duo.menu.current = "explorer"
				}
				b := duo.th.IconButton(duo.ico.Explorer)
				b.Background = duo.menu.icoBackground
				b.Color = duo.menu.icoColor
				b.Padding = duo.menu.icoPadding
				b.Size = duo.menu.icoSize
				b.Layout(duo.gc, &duo.menu.explorer)
			})
		})
		settings := duo.comp.menu.l.Rigid(duo.gc, func() {
			in.Layout(duo.gc, func() {
				for duo.menu.settings.Clicked(duo.gc) {
					duo.menu.current = "settings"
				}
				b := duo.th.IconButton(duo.ico.Settings)
				b.Background = duo.menu.icoBackground
				b.Color = duo.menu.icoColor
				b.Padding = duo.menu.icoPadding
				b.Size = duo.menu.icoSize
				b.Layout(duo.gc, &duo.menu.settings)
			})
		})
		duo.comp.menu.l.Layout(duo.gc, overview, history, addressbook, explorer, settings, )
	})
}
