package monitor

import (
	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/p9c/pod/cmd/gui/pages"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gelook"
)

func BottomBar(gtx *layout.Context, rc *rcd.RcVar) layout.FlexChild {
	return layout.Rigid(func() {
		cs := gtx.Constraints
		gelook.DuoUIdrawRectangle(gtx, cs.Width.Max,
			cs.Height.Max, theme.Colors["Primary"],
			[4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
		layout.Flex{
			Axis: layout.Vertical,
		}.Layout(gtx,
			StatusBar(gtx, rc),
			SettingsPage(gtx, rc),
			BuildPage(gtx, rc),
		)
	})
}

func StatusBar(gtx *layout.Context,
	rc *rcd.RcVar) layout.FlexChild {
	return layout.Rigid(func() {
		cs := gtx.Constraints
		gelook.DuoUIdrawRectangle(gtx, cs.Width.Max,
			cs.Height.Max, theme.Colors["Primary"],
			[4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
		layout.Flex{
			Axis: layout.Horizontal,
		}.Layout(gtx,
			RunControls(gtx),
			RunmodeButtons(gtx, rc),
			BuildButtons(gtx, rc),
			SettingsButtons(gtx, rc),
		)
	})
}

func RunControls(gtx *layout.Context) layout.FlexChild {
	return layout.Rigid(func() {
		if !running {
			theme.DuoUIbutton("", "", "",
				theme.Colors["Primary"], "",
				theme.Colors["Dark"], "Run",
				theme.Colors["Light"], 0, 41, 41, 41,
				0, 0).IconLayout(gtx, runMenuButton)
			for runMenuButton.Clicked(gtx) {
				L.Debug("clicked run button")
				running = true
			}
		}
		if running {
			layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func() {
					theme.DuoUIbutton("", "", "",
						theme.Colors["Dark"], "",
						theme.Colors["Dark"], "Stop",
						theme.Colors["Light"], 0, 41, 41, 41,
						0, 0).IconLayout(gtx, stopMenuButton)
					for stopMenuButton.Clicked(gtx) {
						L.Debug("clicked stop button")
						running = false
						pausing = false
					}
				}),
				layout.Rigid(func() {
					ic := "Pause"
					rc := theme.Colors["Dark"]
					if pausing {
						ic = "Run"
						rc = theme.Colors["Primary"]
					}
					theme.DuoUIbutton("", "", "",
						rc, "",
						theme.Colors["Dark"], ic,
						theme.Colors["Light"], 0, 41, 41, 41,
						0, 0).IconLayout(gtx, pauseMenuButton)
					for pauseMenuButton.Clicked(gtx) {
						if pausing {
							L.Debug("clicked on resume button")
						} else {
							L.Debug("clicked pause button")
						}
						pausing = !pausing
					}
				}),
				layout.Rigid(func() {
					theme.DuoUIbutton("", "", "",
						theme.Colors["Dark"], "",
						theme.Colors["Dark"], "Restart",
						theme.Colors["Light"], 0, 41, 41, 41,
						0, 0).IconLayout(gtx, restartMenuButton)
					for restartMenuButton.Clicked(gtx) {
						L.Debug("clicked restart button")
					}
				}),
			)
		}
	})
}

func RunmodeButtons(gtx *layout.Context, rc *rcd.RcVar) layout.FlexChild {
	return layout.Rigid(func() {
		layout.Flex{
			Axis: layout.Horizontal,
		}.Layout(gtx,
			layout.Rigid(func() {
				bg := theme.Colors["Primary"]
				if runmodeOpen {
					bg = theme.Colors["Dark"]
				}
				s := theme.DuoUIbutton(
					theme.Fonts["Secondary"],
					"mode",
					theme.Colors["Light"],
					bg,
					theme.Colors["Primary"],
					theme.Colors["Light"],
					"settingsIcon",
					theme.Colors["Light"],
					23, 23, 23, 23, 0, 0)
				for runmodeFoldButton.Clicked(gtx) {
					L.Debug("run mode folder clicked")
					// if runmodeOpen && settingsOpen {
					// 	settingsOpen = false
					// }
					runmodeOpen = !runmodeOpen

				}
				s.Layout(gtx, runmodeFoldButton)
			}),
			layout.Rigid(func() {
				if runmodeOpen {
					modes := []string{
						"node", "wallet", "shell", "gui",
					}
					modesList.Layout(gtx, len(modes), func(i int) {
						fg, bg := theme.Colors["Light"], theme.Colors["Dark"]
						if runMode == modes[i] {
							fg, bg = theme.Colors["Dark"], theme.Colors["Light"]
						}
						theme.DuoUIbutton(theme.Fonts["Primary"],
							modes[i],
							fg,
							bg,
							"", "",
							"", "",
							16, 0, 80, 32, 4, 4).Layout(gtx, modesButtons[modes[i]])
						for modesButtons[modes[i]].Clicked(gtx) {
							L.Debug(modes[i], "clicked")
							if runmodeOpen {
								runMode = modes[i]
								runmodeOpen = false
							}
						}
					})
				} else {
					layout.UniformInset(unit.Dp(8)).Layout(gtx, func() {
						t := theme.DuoUIlabel(unit.Dp(18), runMode)
						t.Font.Typeface = theme.Fonts["Primary"]
						t.Layout(gtx)
					})
				}
			}),
		)
	})
}

func SettingsButtons(gtx *layout.Context, rc *rcd.RcVar) layout.FlexChild {
	return layout.Flexed(1, func() {
		layout.Flex{
			Axis: layout.Horizontal,
		}.Layout(gtx,
			layout.Rigid(func() {
				bg := theme.Colors["Primary"]
				if settingsOpen {
					bg = theme.Colors["Dark"]
				}
				theme.DuoUIbutton(theme.Fonts["Primary"], "settings",
					theme.Colors["Light"],
					bg, "",
					theme.Colors["Dark"], "settingsIcon",
					theme.Colors["Light"], 23, 32, 41, 41,
					0, 0).IconLayout(gtx, settingsFoldButton)
				for settingsFoldButton.Clicked(gtx) {
					L.Debug("settings folder clicked")
					switch {
					case !settingsOpen:
						buildOpen = false
						settingsOpen = true
					case settingsOpen:
						settingsOpen = false
					}
				}
				// s.Layout(gtx, settingsFoldButton)
			}),
			layout.Rigid(func() {
				if windowWidth > 1024 && settingsOpen {
					pages.SettingsHeader(rc, gtx, theme)()
				}
			}),
		)
	})
}

func SettingsPage(gtx *layout.Context, rc *rcd.RcVar) layout.FlexChild {
	if !settingsOpen {
		return layout.Flexed(0, func() {})
	}
	var weight float32 = 0.5
	var settingsInset float32 = 0
	switch {
	case windowWidth < 1024 && windowHeight > 1024:
		// weight = 0.333
	case windowHeight < 1024 && windowWidth < 1024:
		weight = 1
	case windowHeight < 600 && windowWidth > 1024:
		weight = 1
	}
	return layout.Flexed(weight, func() {
		cs := gtx.Constraints
		gelook.DuoUIdrawRectangle(gtx, cs.Width.Max,
			cs.Height.Max, theme.Colors["Primary"],
			[4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
		layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func() {
				if windowWidth < 1024 {
					pages.SettingsHeader(rc, gtx, theme)()
				}
			}),
			layout.Rigid(func() {
				cs := gtx.Constraints
				gelook.DuoUIdrawRectangle(gtx, cs.Width.Max,
					cs.Height.Max, theme.Colors["Dark"],
					[4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
				layout.UniformInset(unit.Dp(settingsInset)).Layout(gtx,
					pages.SettingsBody(rc, gtx, theme))
			}),
		)
	})
}

func BuildButtons(gtx *layout.Context, rc *rcd.RcVar) layout.FlexChild {
	return layout.Rigid(func() {
		layout.Flex{
			Axis: layout.Horizontal,
		}.Layout(gtx,
			layout.Rigid(func() {
				bg := theme.Colors["Primary"]
				if buildOpen {
					bg = theme.Colors["Dark"]
				}
				s := theme.DuoUIbutton(
					theme.Fonts["Secondary"],
					"build",
					theme.Colors["Light"],
					bg,
					theme.Colors["Primary"],
					theme.Colors["Light"],
					"settingsIcon",
					theme.Colors["Light"],
					23, 23, 23, 23, 0, 0)
				for buildFoldButton.Clicked(gtx) {
					L.Debug("run mode folder clicked")
					switch {
					case !buildOpen:
						buildOpen = true
						settingsOpen = false
					case buildOpen:
						buildOpen = false
					}

				}
				s.Layout(gtx, buildFoldButton)
			}),
		)
	})
}

func BuildPage(gtx *layout.Context, rc *rcd.RcVar) layout.FlexChild {
	if !buildOpen {
		return layout.Flexed(0, func() {})
	}
	var weight float32 = 0.5
	// var settingsInset float32 = 0
	switch {
	case windowWidth < 1024 && windowHeight > 1024:
		// weight = 0.333
	case windowHeight < 1024 && windowWidth < 1024:
		weight = 1
	case windowHeight < 600 && windowWidth > 1024:
		weight = 1
	}
	return layout.Flexed(weight, func() {
		cs := gtx.Constraints
		gelook.DuoUIdrawRectangle(gtx, cs.Width.Max,
			cs.Height.Max, theme.Colors["Dark"],
			[4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
		// layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		// 	layout.Rigid(func() {
		// 		if windowWidth < 1024 {
		// 			pages.SettingsHeader(rc, gtx, theme)()
		// 		}
		// 	}),
		// 	layout.Rigid(func() {
		// 		cs := gtx.Constraints
		// 		gelook.DuoUIdrawRectangle(gtx, cs.Width.Max,
		// 			cs.Height.Max, theme.Colors["Dark"],
		// 			[4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
		// 		layout.UniformInset(unit.Dp(settingsInset)).Layout(gtx,
		// 			pages.SettingsBody(rc, gtx, theme))
		// 	}),
		// )
	})
}
