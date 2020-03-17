// +build !headless

package monitor

import (
	"fmt"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/p9c/pod/cmd/gui/pages"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gel"
	"github.com/p9c/pod/pkg/gelook"
	"github.com/p9c/pod/pkg/util/interrupt"
)

var (
	theme    = gelook.NewDuoUItheme()
	mainList = &layout.List{
		Axis: layout.Vertical,
	}
	closeButton               = new(gel.Button)
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
	runMode                   = "node"
	running                   = false
	pausing                   = false
	windowWidth, windowHeight int
)

func Run(cx *conte.Xt, rc *rcd.RcVar) (err error) {
	w := app.NewWindow(
		app.Size(unit.Dp(1600), unit.Dp(900)),
		app.Title("ParallelCoin Pod Monitor"),
	)
	gtx := layout.NewContext(w.Queue())
	go func() {
		L.Debug("starting up GUI event loop")
	out:
		for {
			select {
			case <-cx.KillAll:
				L.Debug("kill signal received")
				break out
			case e := <-w.Events():
				switch e := e.(type) {
				case system.DestroyEvent:
					L.Debug("destroy event received")
					close(cx.KillAll)
				case system.FrameEvent:
					gtx.Reset(e.Config, e.Size)
					cs := gtx.Constraints
					windowWidth, windowHeight = cs.Width.Max, cs.Height.Max
					TopLevelLayout(gtx,rc,cx)()
					e.Frame(gtx.Ops)
				}
			}
		}
		L.Debug("gui shut down")
		os.Exit(0)
	}()
	// w.Invalidate()
	interrupt.AddHandler(func() {
		close(cx.KillAll)
	})
	app.Main()
	return
}

func TopLevelLayout(gtx *layout.Context, rc *rcd.RcVar, cx *conte.Xt) func() {
	return func() {
		layout.Flex{
			Axis: layout.Vertical,
		}.Layout(gtx,
			DuoUIheader(gtx, cx),
			Body(gtx, rc),
			BottomBar(gtx, rc),
		)

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
		}.Layout(gtx,
			layout.Rigid(func() {
				SettingsAndRunmodeButtons(gtx, rc)
			}),
			SettingsPage(gtx, rc),
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
					case runmodeOpen:
						settingsOpen = !settingsOpen
					case !settingsOpen:
						settingsOpen = true
						runmodeOpen = false
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

func SettingsAndRunmodeButtons(gtx *layout.Context, rc *rcd.RcVar) {
	cs := gtx.Constraints
	gelook.DuoUIdrawRectangle(gtx, cs.Width.Max,
		cs.Height.Max, theme.Colors["Primary"],
		[4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
	layout.Flex{
		Axis: layout.Horizontal,
	}.Layout(gtx,
		RunControls(gtx),
		RunmodeButtons(gtx, rc),
		SettingsButtons(gtx, rc),
	)
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

func DuoUIheader(gtx *layout.Context, cx *conte.Xt) layout.FlexChild {
	return layout.Rigid(func() {
		layout.Flex{
			Axis:      layout.Horizontal,
			Spacing:   layout.SpaceBetween,
			Alignment: layout.Middle,
		}.Layout(gtx,
			layout.Rigid(func() {
				cs := gtx.Constraints
				gelook.DuoUIdrawRectangle(gtx, cs.Width.Max,
					cs.Height.Max, theme.Colors["Dark"],
					[4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
				var (
					textSize, iconSize               = 64, 64
					width, height                    = 72, 72
					paddingV, paddingH               = 8, 8
					insetSize, textInsetSize float32 = 16, 24
					closeInsetSize           float32 = 4
				)
				if windowWidth < 1024 || windowHeight < 1280 {
					textSize, iconSize = 24, 24
					width, height = 32, 32
					paddingV, paddingH = 8, 8
					insetSize = 10
					textInsetSize = 16
					closeInsetSize = 4
				}
				layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func() {
						layout.UniformInset(unit.Dp(insetSize)).Layout(gtx,
							func() {
								var logoMeniItem gelook.DuoUIbutton
								logoMeniItem = theme.DuoUIbutton(
									"", "",
									"", theme.Colors["Dark"],
									"", "",
									"logo", theme.Colors["Light"],
									textSize, iconSize,
									width, height,
									paddingV, paddingH)
								for logoButton.Clicked(gtx) {
									changeLightDark(theme)
								}
								logoMeniItem.IconLayout(gtx, logoButton)
							},
						)
					}),
					layout.Flexed(1, func() {
						layout.UniformInset(unit.Dp(textInsetSize)).Layout(gtx,
							func() {
								t := theme.DuoUIlabel(unit.Dp(float32(
									textSize)),
									"monitor")
								t.Color = theme.Colors["Light"]
								t.Layout(gtx)
							},
						)
					}),
					layout.Rigid(func() {
						layout.UniformInset(unit.Dp(closeInsetSize*2)).Layout(
							gtx,
							func() {
								t := theme.DuoUIlabel(unit.Dp(float32(
									24)),
									fmt.Sprintf("%dx%d",
										gtx.Constraints.Width.Max,
										gtx.Constraints.Height.Max))
								t.Color = theme.Colors["Light"]
								t.Font.Typeface = theme.Fonts["Primary"]
								t.Layout(gtx)
							})
					}),
					layout.Rigid(func() {
						layout.UniformInset(unit.Dp(closeInsetSize)).Layout(gtx,
							func() {
								theme.DuoUIbutton("", "settings", theme.Colors["Light"],
									"", "",
									theme.Colors["Dark"], "closeIcon",
									theme.Colors["Light"], 0, 32, 41, 41,
									0, 0).IconLayout(gtx, closeButton)
								for closeButton.Clicked(gtx) {
									L.Debug("close button clicked")
									close(cx.KillAll)
								}
							})
					}),
				)
			}),
		)
	})
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
