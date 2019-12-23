package duoui

import (
	"github.com/p9c/pod/pkg/gio/layout"
)

func DuoUImenu(duo *DuoUI) {
	duo.comp.Menu.Layout.Layout(duo.gc,
		layout.Rigid(func() {
			layout.Flex{}.Layout(duo.gc,
				layout.Rigid(func() {
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
				}),
			)
		}),
		layout.Rigid(func() {
			layout.Flex{}.Layout(duo.gc,
				layout.Rigid(func() {
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
				}),

			)
		}),
		layout.Rigid(func() {
			layout.Flex{}.Layout(duo.gc,
				layout.Rigid(func() {
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
				}),
			)
		}),
		layout.Rigid(func() {
			layout.Flex{}.Layout(duo.gc,
				layout.Rigid(func() {
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
				}),
			)
		}),
		layout.Rigid(func() {
			layout.Flex{}.Layout(duo.gc,
				layout.Rigid(func() {
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
				}),
			)
		}),
		layout.Rigid(func() {
			layout.Flex{}.Layout(duo.gc,
				layout.Rigid(func() {
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
				}),
			)
		}),
	)
}

//func DuoUImenu(duo *DuoUI) {
//	duo.comp.Menu.Layout.Layout(duo.gc,
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
//	layout.Flex{}.Layout(duo.gc,
//		layout.Rigid(func() {
//			in.Layout(duo.gc, func() {
//				var menuIco *material.Icon
//				var menuItem *widget.Button
//				switch icoName {
//				case "overview":
//					menuIco = duo.ico.Overview
//					menuItem = &duo.menu.Overview
//				case "history":
//					menuIco = duo.ico.History
//					menuItem = &duo.menu.History
//				case "addressbook":
//					menuIco = duo.ico.AddressBook
//					menuItem = &duo.menu.AddressBook
//				case "explorer":
//					menuIco = duo.ico.Explorer
//					menuItem = &duo.menu.Explorer
//				case "network":
//					menuIco = duo.ico.Network
//					//menuItem = duo.menu.Network
//				case "console":
//					menuIco = duo.ico.Console
//					menuItem = &duo.menu.Console
//				case "settings":
//					menuIco = duo.ico.Settings
//					menuItem = &duo.menu.Settings
//				}
//
//				for menuItem.Clicked(duo.gc) {
//					duo.menu.Current = icoName
//				}
//
//				b := duo.th.IconButton(menuIco)
//				b.Background = duo.menu.IcoBackground
//				b.Color = duo.menu.IcoColor
//				b.Padding = duo.menu.IcoPadding
//				b.Size = duo.menu.IcoSize
//				b.Layout(duo.gc, menuItem)
//			})
//		}),
//	)
//}
