package duoui

import (
	"fmt"
	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"github.com/p9c/pod/cmd/gui/mvc/model"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/unit"
	"github.com/p9c/pod/pkg/util/interrupt"
)

var (
	buttonLog      = new(controller.Button)
	buttonQuit     = new(controller.Button)
	buttonSettings = new(controller.Button)
	buttonNetwork  = new(controller.Button)
	buttonBlocks   = new(controller.Button)
	buttonConsole  = new(controller.Button)
	buttonHelp     = new(controller.Button)
	cornerNav      = &layout.List{
		Axis: layout.Horizontal,
	}
	footerNav = &layout.List{
		Axis: layout.Horizontal,
	}
)

func (ui *DuoUI) DuoUIfooter() func() {
	return func() {
		cs := ui.ly.Context.Constraints
		theme.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, 64, ui.ly.Theme.Color.Dark, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
		var (
			width             float32 = 48
			height            float32 = 48
			iconSize          int     = 32
			paddingVertical   float32 = 8
			paddingHorizontal float32 = 8
		)

		layout.Flex{Spacing: layout.SpaceBetween}.Layout(ui.ly.Context,
			layout.Rigid(func() {
				layout.UniformInset(unit.Dp(0)).Layout(ui.ly.Context, func() {
					cornerButtons := []func(){
						func() {
							layout.UniformInset(unit.Dp(0)).Layout(ui.ly.Context, func() {
								var closeMeniItem theme.DuoUIbutton
								closeMeniItem = ui.ly.Theme.DuoUIbutton("", "", "", "ff303030", "closeIcon", "ffcfcfcf", iconSize, width, height, paddingVertical, paddingHorizontal)
								for buttonQuit.Clicked(ui.ly.Context) {
									ui.rc.Dialog.Show = true
									ui.rc.Dialog = &model.DuoUIdialog{
										Show: true,
										Ok: func() {
											interrupt.Request()
										},
										Close: func() {
											interrupt.RequestRestart()
										},
										Cancel: func() { ui.rc.Dialog.Show = false },
										Title:  "Are you sure?",
										Text:   "Confirm ParallelCoin close",
									}
								}
								closeMeniItem.Layout(ui.ly.Context, buttonQuit)
							})
						},

						func() {
							var logMenuItem theme.DuoUIbutton
							logMenuItem = ui.ly.Theme.DuoUIbutton("", "", "", "ff303030", "traceIcon", "ffcfcfcf", iconSize, width, height, paddingVertical, paddingHorizontal)

							for buttonLog.Clicked(ui.ly.Context) {
								ui.rc.ShowPage = "LOG"
							}
							logMenuItem.Layout(ui.ly.Context, buttonLog)
						},
					}
					cornerNav.Layout(ui.ly.Context, len(cornerButtons), func(i int) {
						layout.UniformInset(unit.Dp(0)).Layout(ui.ly.Context, cornerButtons[i])
					})
				})
			}),
			layout.Rigid(func() {
				layout.UniformInset(unit.Dp(0)).Layout(ui.ly.Context, func() {
					navButtons := []func(){

						//func() {
						//	layout.UniformInset(unit.Dp(14)).Layout(ui.ly.Context, func() {
						//		tim := ui.ly.Theme.Caption("Blocks: " + fmt.Sprint(ui.rc.Status.Node.BlockHeight))
						//		tim.Font.Typeface = "bariol"
						//		tim.Color = ui.ly.Theme.Color.Light
						//		tim.Layout(ui.ly.Context)
						//	})
						//},
						//func() {
						//	layout.UniformInset(unit.Dp(14)).Layout(ui.ly.Context, func() {
						//		tim := ui.ly.Theme.Caption("Connections: " + fmt.Sprint(ui.rc.Status.Node.ConnectionCount))
						//		tim.Font.Typeface = "bariol"
						//		tim.Color = ui.ly.Theme.Color.Light
						//		tim.Layout(ui.ly.Context)
						//	})
						//},

						func() {
							layout.UniformInset(unit.Dp(0)).Layout(ui.ly.Context, func() {
								var networkMeniItem theme.DuoUIbutton
								networkMeniItem = ui.ly.Theme.DuoUIbutton(ui.ly.Theme.Font.Primary, "Connections: "+fmt.Sprint(ui.rc.Status.Node.ConnectionCount), "ffcfcfcf", "", "", "", iconSize, 80, height, paddingVertical, 0)
								for buttonNetwork.Clicked(ui.ly.Context) {
									//ui.rc.ShowPage = "NETWORK"
								}
								networkMeniItem.Layout(ui.ly.Context, buttonNetwork)
							})
						},
						func() {
							layout.UniformInset(unit.Dp(0)).Layout(ui.ly.Context, func() {
								var blocksMenuItem theme.DuoUIbutton
								blocksMenuItem = ui.ly.Theme.DuoUIbutton(ui.ly.Theme.Font.Primary, "Blocks: "+fmt.Sprint(ui.rc.Status.Node.BlockHeight), "ffcfcfcf", "", "", "", iconSize, 80, height, paddingVertical, 0)
								for buttonBlocks.Clicked(ui.ly.Context) {
									//ui.rc.ShowPage = "EXPLORER"
									//ui.rc.ShowToast = true
									//ui.toastAdd()
								}
								blocksMenuItem.Layout(ui.ly.Context, buttonBlocks)
							})
						},

						//func() {
						//	layout.UniformInset(unit.Dp(0)).Layout(ui.ly.Context, func() {
						//		var networkMeniItem theme.DuoUIbutton
						//		networkMeniItem = ui.ly.Theme.DuoUIbutton("","", "", "ff303030", "networkIcon", "ffcfcfcf", iconSize, width, height, paddingVertical, paddingHorizontal)
						//		for buttonNetwork.Clicked(ui.ly.Context) {
						//			ui.rc.ShowPage = "NETWORK"
						//		}
						//		networkMeniItem.Layout(ui.ly.Context, buttonNetwork)
						//	})
						//},
						//func() {
						//	var blocksMenuItem theme.DuoUIbutton
						//	blocksMenuItem = ui.ly.Theme.DuoUIbutton("","", "", "ff303030", "blocksIcon", "ffcfcfcf", iconSize, width, height, paddingVertical, paddingHorizontal)
						//	for buttonBlocks.Clicked(ui.ly.Context) {
						//		//duo.CurrentPage = "EXPLORER"
						//		//ui.rc.ShowToast = true
						//		//toastAdd(duo, rc)
						//	}
						//	blocksMenuItem.Layout(ui.ly.Context, buttonBlocks)
						//},
						func() {
							var helpMenuItem theme.DuoUIbutton
							helpMenuItem = ui.ly.Theme.DuoUIbutton("", "", "", "ff303030", "helpIcon", "ffcfcfcf", iconSize, width, height, paddingVertical, paddingHorizontal)
							for buttonHelp.Clicked(ui.ly.Context) {
								//ui.rc.ShowDialog = true
							}
							helpMenuItem.Layout(ui.ly.Context, buttonHelp)
						},
						func() {
							layout.UniformInset(unit.Dp(0)).Layout(ui.ly.Context, func() {
								var consoleMenuItem theme.DuoUIbutton
								consoleMenuItem = ui.ly.Theme.DuoUIbutton("", "", "", "ff303030", "consoleIcon", "ffcfcfcf", iconSize, width, height, paddingVertical, paddingHorizontal)
								for buttonConsole.Clicked(ui.ly.Context) {
									ui.rc.ShowPage = "CONSOLE"
								}
								consoleMenuItem.Layout(ui.ly.Context, buttonConsole)
							})
						},
						func() {
							var settingsMenuItem theme.DuoUIbutton
							settingsMenuItem = ui.ly.Theme.DuoUIbutton("", "", "", "ff303030", "settingsIcon", "ffcfcfcf", iconSize, width, height, paddingVertical, paddingHorizontal)

							for buttonSettings.Clicked(ui.ly.Context) {
								ui.rc.ShowPage = "SETTINGS"
							}
							settingsMenuItem.Layout(ui.ly.Context, buttonSettings)
						},
					}
					footerNav.Layout(ui.ly.Context, len(navButtons), func(i int) {
						layout.UniformInset(unit.Dp(0)).Layout(ui.ly.Context, navButtons[i])
					})
				})
			}),
		)
	}
}
