package duoui

import (
	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/unit"
)

var (
	buttonOverview    = new(controller.Button)
	buttonSend        = new(controller.Button)
	buttonReceive     = new(controller.Button)
	buttonAddressBook = new(controller.Button)
	buttonHistory     = new(controller.Button)
	mainNav           = &layout.List{
		Axis: layout.Vertical,
	}
)

func (ui *DuoUI)DuoUImenu() func() {
	return func() {

		var (
			width             float32 = 96
			height            float32 = 72
			iconSize          int     = 72
			paddingVertical   float32 = 0
			paddingHorizontal float32 = 0
			bgColor = "ff9a9a9a"
			textColor = "ff303030"
			iconColor = "ffacacac"
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
										overviewMenuItem = ui.ly.Theme.DuoUIbutton(ui.ly.Theme.Font.Secondary,"OVERVIEW", textColor, bgColor, "overviewIcon", iconColor, iconSize, width, height, paddingVertical, paddingHorizontal)
										for buttonOverview.Clicked(ui.ly.Context) {
											ui.rc.ShowPage = "OVERVIEW"
										}
										overviewMenuItem.Layout(ui.ly.Context, buttonOverview)
									})
								},
								func() {
									theme.DuoUIdrawRectangle(ui.ly.Context, int(width), 1, "ff888888", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
								},
								func() {
									var sendMenuItem theme.DuoUIbutton
									sendMenuItem = ui.ly.Theme.DuoUIbutton(ui.ly.Theme.Font.Secondary,"SEND", textColor, bgColor, "sendIcon", iconColor, iconSize, width, height, paddingVertical, paddingHorizontal)
									for buttonSend.Clicked(ui.ly.Context) {
										ui.rc.ShowPage = "SEND"
									}
									sendMenuItem.Layout(ui.ly.Context, buttonSend)
								},
								func() {
									theme.DuoUIdrawRectangle(ui.ly.Context, int(width), 1, "ff888888", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
								},
								func() {
									in.Layout(ui.ly.Context, func() {
										var receiveMenuItem theme.DuoUIbutton
										receiveMenuItem = ui.ly.Theme.DuoUIbutton(ui.ly.Theme.Font.Secondary,"RECEIVE", textColor,bgColor,  "receiveIcon", iconColor, iconSize, width, height, paddingVertical, paddingHorizontal)
										for buttonReceive.Clicked(ui.ly.Context) {
											ui.rc.ShowPage = "RECEIVE"
										}
										receiveMenuItem.Layout(ui.ly.Context, buttonReceive)
									})
								},
								func() {
									theme.DuoUIdrawRectangle(ui.ly.Context, int(width), 1, "ff888888", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
								},
								func() {
									var addressBookMenuItem theme.DuoUIbutton
									addressBookMenuItem = ui.ly.Theme.DuoUIbutton(ui.ly.Theme.Font.Secondary,"ADDRESS BOOK", textColor, bgColor, "addressBookIcon", iconColor, iconSize, width, height, paddingVertical, paddingHorizontal)
									for buttonAddressBook.Clicked(ui.ly.Context) {
										ui.rc.ShowPage = "ADDRESSBOOK"
									}
									addressBookMenuItem.Layout(ui.ly.Context, buttonAddressBook)
								},
								func() {
									theme.DuoUIdrawRectangle(ui.ly.Context, int(width), 1, "ff888888", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
								},
								func() {
									var historyMenuItem theme.DuoUIbutton
									historyMenuItem = ui.ly.Theme.DuoUIbutton(ui.ly.Theme.Font.Secondary,"HISTORY", textColor,  bgColor, "historyIcon", iconColor, iconSize, width, height, paddingVertical, paddingHorizontal)
									for buttonHistory.Clicked(ui.ly.Context) {
										ui.rc.ShowPage = "HISTORY"
									}
									historyMenuItem.Layout(ui.ly.Context, buttonHistory)
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
