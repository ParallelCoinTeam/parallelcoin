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
	log "github.com/p9c/pod/pkg/logi"
)

var (
	theme    = gelook.NewDuoUItheme()
	mainList = &layout.List{
		Axis: layout.Vertical,
	}
	logoButton         = new(gel.Button)
	settingsFoldButton = new(gel.Button)
	settingsOpen       bool
)

var monitorHandle = func(cx *conte.Xt) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		Configure(cx, c)
		rc := rcd.RcInit(cx)
		log.L.Warn("starting monitor GUI")
		w := app.NewWindow(
			app.Size(unit.Dp(1600), unit.Dp(900)),
			app.Title("ParallelCoin"),
		)
		gtx := layout.NewContext(w.Queue())
		for e := range w.Events() {
			switch e := e.(type) {
			case system.DestroyEvent:
				log.L.Debug("destroy event received")
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
				}), layout.Flexed(1, func() {
					layout.Flex{
						Axis: layout.Vertical,
					}.Layout(gtx, layout.Flexed(1, func() {
						cs := gtx.Constraints
						gelook.DuoUIdrawRectangle(gtx, cs.Width.Max,
							cs.Height.Max, theme.Colors["Light"],
							[4]float32{0, 0, 0, 0},
							[4]float32{0, 0, 0, 0},
						)
					}), layout.Rigid(func() {
						cs := gtx.Constraints
						gelook.DuoUIdrawRectangle(gtx, cs.Width.Max,
							cs.Height.Max, theme.Colors["Light"],
							[4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
						layout.Flex{
							Axis: layout.Vertical,
						}.Layout(gtx, layout.Rigid(func() {
							cs := gtx.Constraints
							gelook.DuoUIdrawRectangle(gtx, cs.Width.Max,
								cs.Height.Max, theme.Colors["Dark"],
								[4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
							layout.Flex{
								Axis: layout.Horizontal,
							}.Layout(gtx, layout.Rigid(func() {
								s := theme.DuoUIbutton(
									theme.Fonts["Secondary"],
									"SETTINGS",
									theme.Colors["Light"],
									theme.Colors["Dark"],
									theme.Colors["Dark"],
									theme.Colors["Light"],
									"settingsIcon",
									theme.Colors["Light"],
									23, 0, 80, 48, 4, 4)
								for settingsFoldButton.Clicked(gtx) {
									log.L.Debug("settings folder clicked")
									settingsOpen = !settingsOpen
								}
								s.Layout(gtx, settingsFoldButton)
							}),
								layout.Flexed(1, func() {
									if settingsOpen {
										pages.SettingsHeader(rc, gtx, theme)()
									}
								}),
							)
						}),
							SettingsPage(gtx, rc),
						)
					}),
					)
				}),
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

func SettingsPage(gtx *layout.Context, rc *rcd.RcVar) layout.FlexChild {
	if !settingsOpen {
		return layout.Flexed(0, func() {})
	}
	return layout.Flexed(0.5, func() {
		// cs := gtx.Constraints
		// gelook.DuoUIdrawRectangle(gtx, cs.Width.Max,
		// 	cs.Height.Max, theme.Colors["Light"],
		// 	[4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
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
						layout.UniformInset(unit.Dp(10)).Layout(gtx,
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
