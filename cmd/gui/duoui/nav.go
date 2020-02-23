package duoui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
)

var (
	navButtonOverview    = new(controller.Button)
	navButtonSend        = new(controller.Button)
	navButtonReceive     = new(controller.Button)
	navButtonAddressBook = new(controller.Button)
	navButtonHistory     = new(controller.Button)
	mainNav              = &layout.List{
		Axis: layout.Vertical,
	}
)

func (ui *DuoUI) DuoUImenu() func() {
	return func() {
		var (
			width             int = 96
			height            int = 72
			textSize          int = 48
			iconSize          int = 36
			paddingVertical   int = 8
			paddingHorizontal int = 0
			bgColor               = "ff9a9a9a"
			textColor             = "ff303030"
			iconColor             = "ffacacac"
		)
		//overviewButton :=
		//historyButton :=

		in := layout.UniformInset(unit.Dp(0))

		layout.Flex{
			Axis:      layout.Vertical,
			Alignment: layout.Middle,
			Spacing:   layout.SpaceEvenly}.Layout(ui.ly.Context,
			layout.Rigid(func() {
				layout.Flex{}.Layout(ui.ly.Context,
					layout.Rigid(func() {
						in.Layout(ui.ly.Context, func() {

							navButtons := []func(){
								func() {
									in.Layout(ui.ly.Context, func() {
										var overviewMenuItem theme.DuoUIbutton
										overviewMenuItem = ui.ly.Theme.DuoUIbutton(ui.ly.Theme.Font.Secondary, "OVERVIEW", textColor, bgColor, "overviewIcon", iconColor, textSize, iconSize, width, height, paddingVertical, paddingHorizontal)
										for navButtonOverview.Clicked(ui.ly.Context) {
											ui.rc.ShowPage = "OVERVIEW"
										}
										overviewMenuItem.MenuLayout(ui.ly.Context, navButtonOverview)
									})
								},
								func() {
									theme.DuoUIdrawRectangle(ui.ly.Context, int(width), 1, "ff888888", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
								},
								func() {
									var sendMenuItem theme.DuoUIbutton
									sendMenuItem = ui.ly.Theme.DuoUIbutton(ui.ly.Theme.Font.Secondary, "SEND", textColor, bgColor, "sendIcon", iconColor, textSize, iconSize, width, height, paddingVertical, paddingHorizontal)
									for navButtonSend.Clicked(ui.ly.Context) {
										ui.rc.ShowPage = "SEND"
									}
									sendMenuItem.MenuLayout(ui.ly.Context, navButtonSend)
								},
								func() {
									theme.DuoUIdrawRectangle(ui.ly.Context, int(width), 1, "ff888888", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
								},
								func() {
									in.Layout(ui.ly.Context, func() {
										var receiveMenuItem theme.DuoUIbutton
										receiveMenuItem = ui.ly.Theme.DuoUIbutton(ui.ly.Theme.Font.Secondary, "RECEIVE", textColor, bgColor, "receiveIcon", iconColor, textSize, iconSize, width, height, paddingVertical, paddingHorizontal)
										for navButtonReceive.Clicked(ui.ly.Context) {
											ui.rc.ShowPage = "RECEIVE"
										}
										receiveMenuItem.MenuLayout(ui.ly.Context, navButtonReceive)
									})
								},
								func() {
									theme.DuoUIdrawRectangle(ui.ly.Context, int(width), 1, "ff888888", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
								},
								func() {
									var addressBookMenuItem theme.DuoUIbutton
									addressBookMenuItem = ui.ly.Theme.DuoUIbutton(ui.ly.Theme.Font.Secondary, "ADDRESS BOOK", textColor, bgColor, "addressBookIcon", iconColor, textSize, iconSize, width, height, paddingVertical, paddingHorizontal)
									for navButtonAddressBook.Clicked(ui.ly.Context) {
										ui.rc.ShowPage = "ADDRESSBOOK"
									}
									addressBookMenuItem.MenuLayout(ui.ly.Context, navButtonAddressBook)
								},
								func() {
									theme.DuoUIdrawRectangle(ui.ly.Context, int(width), 1, "ff888888", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
								},
								func() {
									var historyMenuItem theme.DuoUIbutton
									historyMenuItem = ui.ly.Theme.DuoUIbutton(ui.ly.Theme.Font.Secondary, "HISTORY", textColor, bgColor, "historyIcon", iconColor, textSize, iconSize, width, height, paddingVertical, paddingHorizontal)
									for navButtonHistory.Clicked(ui.ly.Context) {
										ui.rc.ShowPage = "HISTORY"
									}
									historyMenuItem.MenuLayout(ui.ly.Context, navButtonHistory)
								},
							}
							mainNav.Layout(ui.ly.Context, len(navButtons), func(i int) {
								layout.UniformInset(unit.Dp(0)).Layout(ui.ly.Context, navButtons[i])
							})
						})
					}),
				)
			}),
		)
	}
}
