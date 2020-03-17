// +build !headless

package app

import (
	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/urfave/cli"

	"github.com/p9c/pod/cmd/gui/pages"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gel"
	"github.com/p9c/pod/pkg/gelook"
)

var (
	theme    = gelook.NewDuoUItheme()
	mainList = &layout.List{
		Axis: layout.Vertical,
	}
	logoButton                = new(gel.Button)
	runMenuButton             = new(gel.Button)
	stopMenuButton            = new(gel.Button)
	pauseMenuButton           = new(gel.Button)
	restartMenuButton         = new(gel.Button)
	settingsFoldButton        = new(gel.Button)
	runmodeFoldButton         = new(gel.Button)
	settingsOpen, runmodeOpen bool
	modesList                 = &layout.List{
		Axis:      layout.Horizontal,
		Alignment: layout.Start,
	}
	modesButtons = map[string]*gel.Button{
		"node":   new(gel.Button),
		"wallet": new(gel.Button),
		"shell":  new(gel.Button),
		"gui":    new(gel.Button),
	}
	runMode = "node"
	running = false
	pausing = false
)

var monitorHandle = func(cx *conte.Xt) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		Configure(cx, c)
		rc := rcd.RcInit(cx)
		L.Warn("starting monitor GUI")
		w := app.NewWindow(
			app.Size(unit.Dp(1600), unit.Dp(900)),
			app.Title("ParallelCoin"),
		)
		gtx := layout.NewContext(w.Queue())
		for e := range w.Events() {
			switch e := e.(type) {
			case system.DestroyEvent:
				L.Debug("destroy event received")
				close(cx.KillAll)
				return e.Err
			case system.FrameEvent:
				gtx.Reset(e.Config, e.Size)
				layout.Flex{
					Axis: layout.Vertical,
				}.Layout(gtx, layout.Rigid(func() {
					cs := gtx.Constraints
					gelook.DuoUIdrawRectangle(gtx, cs.Width.Max,
						cs.Height.Max, theme.Colors["Dark"],
						[4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
					DuoUIheader(gtx, theme)()
				}),
					Body(gtx, rc),
					BottomBar(gtx, rc),
				)
				e.Frame(gtx.Ops)
			}
			// w.Invalidate()
		}
		go app.Main()
		select {
		case <-cx.KillAll:
		}
		return
	}
}

func Body(gtx *layout.Context, rc *rcd.RcVar) layout.FlexChild {
	return layout.Flexed(1, func() {
		layout.Flex{
			Axis: layout.Vertical,
		}.Layout(gtx, layout.Flexed(1, func() {
			cs := gtx.Constraints
			gelook.DuoUIdrawRectangle(gtx, cs.Width.Max,
				cs.Height.Max, theme.Colors["Light"],
				[4]float32{0, 0, 0, 0},
				[4]float32{0, 0, 0, 0},
			)
		}),
		)
	})
}

func BottomBar(gtx *layout.Context, rc *rcd.RcVar) layout.FlexChild {
	return layout.Rigid(func() {
		cs := gtx.Constraints
		gelook.DuoUIdrawRectangle(gtx, cs.Width.Max,
			cs.Height.Max, theme.Colors["Primary"],
			[4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
		layout.Flex{
			Axis: layout.Vertical,
		}.Layout(gtx, layout.Rigid(func() {
			cs := gtx.Constraints
			gelook.DuoUIdrawRectangle(gtx, cs.Width.Max,
				cs.Height.Max, theme.Colors["Primary"],
				[4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
			layout.Flex{
				Axis: layout.Horizontal,
			}.Layout(gtx,
				layout.Rigid(func() {
					settingsAndRunmodeButtons(gtx, rc)
				}),
			)
		}),
			SettingsPage(gtx, rc),
			// RunModePage(gtx, rc),
		)
	})
}

func settingsAndRunmodeButtons(gtx *layout.Context, rc *rcd.RcVar) {
	layout.Flex{
		Axis: layout.Horizontal,
	}.Layout(gtx,
		layout.Rigid(func() {
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
		}),
		layout.Rigid(func() {
			bg := theme.Colors["Primary"]
			if runmodeOpen {
				bg = theme.Colors["Dark"]
			}
			s := theme.DuoUIbutton(
				theme.Fonts["Secondary"],
				"run mode",
				theme.Colors["Light"],
				bg,
				theme.Colors["Primary"],
				theme.Colors["Light"],
				"settingsIcon",
				theme.Colors["Light"],
				23, 23, 23, 23, 0, 0)
			for runmodeFoldButton.Clicked(gtx) {
				L.Debug("run mode folder clicked")
				if runmodeOpen && settingsOpen {
					settingsOpen = false
				}
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
		layout.Rigid(func() {
			bg := theme.Colors["Primary"]
			if settingsOpen {
				bg = theme.Colors["Dark"]
			}
			s := theme.DuoUIbutton(
				theme.Fonts["Secondary"],
				"SETTINGS",
				theme.Colors["Light"],
				bg,
				theme.Colors["Primary"],
				theme.Colors["Light"],
				"settingsIcon",
				theme.Colors["Light"],
				23, 0, 80, 48, 4, 4)
			for settingsFoldButton.Clicked(gtx) {
				L.Debug("settings folder clicked")
				switch {
				case runmodeOpen:
					settingsOpen = !settingsOpen
				case !settingsOpen:
					settingsOpen = true
					runmodeOpen = false
				case settingsOpen:
					settingsOpen = false
				}
			}
			s.Layout(gtx, settingsFoldButton)
		}),
		layout.Rigid(func() {
			if settingsOpen {
				pages.SettingsHeader(rc, gtx, theme)()
			}
		}),
		layout.Flexed(1, func() {
			// if settingsOpen {
			// 	pages.SettingsHeader(rc, gtx, theme)()
			// }
		}),
	)
}

func SettingsPage(gtx *layout.Context, rc *rcd.RcVar) layout.FlexChild {
	if !settingsOpen {
		return layout.Flexed(0, func() {})
	}
	return layout.Flexed(0.5, func() {
		cs := gtx.Constraints
		gelook.DuoUIdrawRectangle(gtx, cs.Width.Max,
			cs.Height.Max, theme.Colors["Dark"],
			[4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
		controllers := []func(){
			func() {
				pages.SettingsBody(rc, gtx, theme)()
			},
		}
		mainList.Layout(gtx, len(controllers), func(i int) {
			layout.UniformInset(unit.Dp(10)).Layout(gtx,
				controllers[i])
		})
	})

}

func DuoUIheader(gtx *layout.Context, theme *gelook.DuoUItheme) func() {
	return func() {
		layout.Flex{
			Axis:      layout.Horizontal,
			Spacing:   layout.SpaceBetween,
			Alignment: layout.Middle,
		}.Layout(gtx,
			layout.Rigid(func() {
				layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func() {
						var logoMeniItem gelook.DuoUIbutton
						logoMeniItem = theme.DuoUIbutton("", "", "",
							theme.Colors["Dark"], "", "", "logo",
							theme.Colors["Light"], 16, 64, 96, 96, 8, 8)
						for logoButton.Clicked(gtx) {
							changeLightDark(theme)
						}
						logoMeniItem.IconLayout(gtx, logoButton)
					}),
					layout.Flexed(1, func() {
						layout.UniformInset(unit.Dp(20)).Layout(gtx,
							func() {
								t := theme.H2("monitor")
								t.Color = theme.Colors["Light"]
								t.Layout(gtx)
							},
						)
					}),
				)
			}),
		)
	}
}

func changeLightDark(theme *gelook.DuoUItheme) {
	light := theme.Colors["Light"]
	dark := theme.Colors["Dark"]
	lightGray := theme.Colors["LightGrayIII"]
	darkGray := theme.Colors["DarkGrayII"]
	theme.Colors["Light"] = dark
	theme.Colors["Dark"] = light
	theme.Colors["LightGrayIII"] = darkGray
	theme.Colors["DarkGrayII"] = lightGray
}
