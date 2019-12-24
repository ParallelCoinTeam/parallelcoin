package duoui

import (
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/unit"
)

func DuoUImenu(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) {

	in := layout.UniformInset(unit.Dp(8))

	duo.Comp.Menu.Layout.Layout(duo.Gc,
		layout.Rigid(func() {
			layout.Flex{}.Layout(duo.Gc,
				layout.Rigid(func() {
					in.Layout(duo.Gc, func() {
						for duo.Menu.Overview.Clicked(duo.Gc) {
							duo.Menu.Current = "overview"
						}
						b := duo.Th.IconButton(duo.Ico.Overview)
						b.Background = duo.Menu.IcoBackground
						b.Color = duo.Menu.IcoColor
						b.Padding = duo.Menu.IcoPadding
						b.Size = duo.Menu.IcoSize
						b.Layout(duo.Gc, &duo.Menu.Overview)
					})
				}),
			)
		}),
		layout.Rigid(func() {
			layout.Flex{}.Layout(duo.Gc,
				layout.Rigid(func() {
					in.Layout(duo.Gc, func() {
						for duo.Menu.History.Clicked(duo.Gc) {
							duo.Menu.Current = "history"
						}
						b := duo.Th.IconButton(duo.Ico.History)
						b.Background = duo.Menu.IcoBackground
						b.Color = duo.Menu.IcoColor
						b.Padding = duo.Menu.IcoPadding
						b.Size = duo.Menu.IcoSize
						b.Layout(duo.Gc, &duo.Menu.History)
					})
				}),

			)
		}),
		layout.Rigid(func() {
			layout.Flex{}.Layout(duo.Gc,
				layout.Rigid(func() {
					in.Layout(duo.Gc, func() {
						for duo.Menu.AddressBook.Clicked(duo.Gc) {
							duo.Menu.Current = "addressbook"
						}
						b := duo.Th.IconButton(duo.Ico.AddressBook)
						b.Background = duo.Menu.IcoBackground
						b.Color = duo.Menu.IcoColor
						b.Padding = duo.Menu.IcoPadding
						b.Size = duo.Menu.IcoSize
						b.Layout(duo.Gc, &duo.Menu.AddressBook)
					})
				}),
			)
		}),
		layout.Rigid(func() {
			layout.Flex{}.Layout(duo.Gc,
				layout.Rigid(func() {
					in.Layout(duo.Gc, func() {
						for duo.Menu.Explorer.Clicked(duo.Gc) {
							duo.Menu.Current = "explorer"
						}
						b := duo.Th.IconButton(duo.Ico.Explorer)
						b.Background = duo.Menu.IcoBackground
						b.Color = duo.Menu.IcoColor
						b.Padding = duo.Menu.IcoPadding
						b.Size = duo.Menu.IcoSize
						b.Layout(duo.Gc, &duo.Menu.Explorer)
					})
				}),
			)
		}),
		layout.Rigid(func() {
			layout.Flex{}.Layout(duo.Gc,
				layout.Rigid(func() {
					in.Layout(duo.Gc, func() {
						for duo.Menu.Console.Clicked(duo.Gc) {
							duo.Menu.Current = "console"
						}
						b := duo.Th.IconButton(duo.Ico.Console)
						b.Background = duo.Menu.IcoBackground
						b.Color = duo.Menu.IcoColor
						b.Padding = duo.Menu.IcoPadding
						b.Size = duo.Menu.IcoSize
						b.Layout(duo.Gc, &duo.Menu.Console)
					})
				}),
			)
		}),
		layout.Rigid(func() {
			layout.Flex{}.Layout(duo.Gc,
				layout.Rigid(func() {
					in.Layout(duo.Gc, func() {
						for duo.Menu.Settings.Clicked(duo.Gc) {
							duo.Menu.Current = "settings"
							rc.Settings.Daemon = rcd.GetCoreSettings(cx)
						}
						b := duo.Th.IconButton(duo.Ico.Settings)
						b.Background = duo.Menu.IcoBackground
						b.Color = duo.Menu.IcoColor
						b.Padding = duo.Menu.IcoPadding
						b.Size = duo.Menu.IcoSize
						b.Layout(duo.Gc, &duo.Menu.Settings)
					})
				}),
			)
		}),
	)
}

//func DuoUImenu(duo *DuoUI) {
//	duo.comp.Menu.Layout.Layout(duo.Gc,
//		layout.Rigid(func() {
//			layout.Rigid(func() {
//				menuItem(duo, "overview")
//				menuItem(duo, "history")
//				menuItem(duo, "addressbook")
//				menuItem(duo, "explorer")
//				menuItem(duo, "network")
//				menuItem(duo, "console")
//				menuItem(duo, "settings")
//			})
//		}))
//}
//
//func menuItem(duo *DuoUI, icoName string) {
//	layout.Flex{}.Layout(duo.Gc,
//		layout.Rigid(func() {
//			in.Layout(duo.Gc, func() {
//				var menuIco *material.Icon
//				var menuItem *widget.Button
//				switch icoName {
//				case "overview":
//					menuIco = duo.Ico.Overview
//					menuItem = &duo.Menu.Overview
//				case "history":
//					menuIco = duo.Ico.History
//					menuItem = &duo.Menu.History
//				case "addressbook":
//					menuIco = duo.Ico.AddressBook
//					menuItem = &duo.Menu.AddressBook
//				case "explorer":
//					menuIco = duo.Ico.Explorer
//					menuItem = &duo.Menu.Explorer
//				case "network":
//					menuIco = duo.Ico.Network
//					//menuItem = duo.Menu.Network
//				case "console":
//					menuIco = duo.Ico.Console
//					menuItem = &duo.Menu.Console
//				case "settings":
//					menuIco = duo.Ico.Settings
//					menuItem = &duo.Menu.Settings
//				}
//
//				for menuItem.Clicked(duo.Gc) {
//					duo.Menu.Current = icoName
//				}
//
//				b := duo.Th.IconButton(menuIco)
//				b.Background = duo.Menu.IcoBackground
//				b.Color = duo.Menu.IcoColor
//				b.Padding = duo.Menu.IcoPadding
//				b.Size = duo.Menu.IcoSize
//				b.Layout(duo.Gc, menuItem)
//			})
//		}),
//	)
//}
