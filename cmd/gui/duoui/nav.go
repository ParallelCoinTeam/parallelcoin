package duoui

import (
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/unit"
)

func DuoUImenu(duo *DuoUI) layout.FlexChild {
	return duo.comp.Sidebar.Layout.Rigid(duo.gc, func() {

		in := layout.UniformInset(unit.Dp(0))

		overview := duo.comp.Menu.Layout.Rigid(duo.gc, func() {
			in.Layout(duo.gc, func() {
				for duo.menu.Overview.Clicked(duo.gc) {
					duo.menu.Current = "overview"
				}
				b := duo.th.IconButton(duo.ico.Overview)
				b.Background = duo.menu.IcoBackground
				b.Color = duo.menu.IcoColor
				b.Padding = duo.menu.IcoPadding
				b.Size = duo.menu.IcoSize
				b.Layout(duo.gc, &duo.menu.Overview)
			})
		})
		history := duo.comp.Menu.Layout.Rigid(duo.gc, func() {
			in.Layout(duo.gc, func() {
				for duo.menu.History.Clicked(duo.gc) {
					duo.menu.Current = "history"
				}
				b := duo.th.IconButton(duo.ico.History)
				b.Background = duo.menu.IcoBackground
				b.Color = duo.menu.IcoColor
				b.Padding = duo.menu.IcoPadding
				b.Size = duo.menu.IcoSize
				b.Layout(duo.gc, &duo.menu.History)
			})
		})
		addressbook := duo.comp.Menu.Layout.Rigid(duo.gc, func() {
			in.Layout(duo.gc, func() {
				for duo.menu.AddressBook.Clicked(duo.gc) {
					duo.menu.Current = "addressbook"
				}
				b := duo.th.IconButton(duo.ico.AddressBook)
				b.Background = duo.menu.IcoBackground
				b.Color = duo.menu.IcoColor
				b.Padding = duo.menu.IcoPadding
				b.Size = duo.menu.IcoSize
				b.Layout(duo.gc, &duo.menu.AddressBook)
			})
		})
		explorer := duo.comp.Menu.Layout.Rigid(duo.gc, func() {
			in.Layout(duo.gc, func() {
				for duo.menu.Explorer.Clicked(duo.gc) {
					duo.menu.Current = "explorer"
				}
				b := duo.th.IconButton(duo.ico.Explorer)
				b.Background = duo.menu.IcoBackground
				b.Color = duo.menu.IcoColor
				b.Padding = duo.menu.IcoPadding
				b.Size = duo.menu.IcoSize
				b.Layout(duo.gc, &duo.menu.Explorer)
			})
		})
		console := duo.comp.Menu.Layout.Rigid(duo.gc, func() {
			in.Layout(duo.gc, func() {
				for duo.menu.Console.Clicked(duo.gc) {
					duo.menu.Current = "console"
				}
				b := duo.th.IconButton(duo.ico.Console)
				b.Background = duo.menu.IcoBackground
				b.Color = duo.menu.IcoColor
				b.Padding = duo.menu.IcoPadding
				b.Size = duo.menu.IcoSize
				b.Layout(duo.gc, &duo.menu.Console)
			})
		})
		settings := duo.comp.Menu.Layout.Rigid(duo.gc, func() {
			in.Layout(duo.gc, func() {
				for duo.menu.Settings.Clicked(duo.gc) {
					duo.menu.Current = "settings"
				}
				b := duo.th.IconButton(duo.ico.Settings)
				b.Background = duo.menu.IcoBackground
				b.Color = duo.menu.IcoColor
				b.Padding = duo.menu.IcoPadding
				b.Size = duo.menu.IcoSize
				b.Layout(duo.gc, &duo.menu.Settings)
			})
		})
		duo.comp.Menu.Layout.Layout(duo.gc, overview, history, addressbook, explorer, console, settings )
	})
}
