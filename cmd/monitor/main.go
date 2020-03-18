// +build !headless

package monitor

import (
	"fmt"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/unit"

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
	closeButton                          = new(gel.Button)
	logoButton                           = new(gel.Button)
	runMenuButton                        = new(gel.Button)
	stopMenuButton                       = new(gel.Button)
	pauseMenuButton                      = new(gel.Button)
	restartMenuButton                    = new(gel.Button)
	settingsFoldButton                   = new(gel.Button)
	runmodeFoldButton                    = new(gel.Button)
	buildFoldButton                      = new(gel.Button)
	settingsOpen, runmodeOpen, buildOpen bool
	modesList                            = &layout.List{
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
					TopLevelLayout(gtx, rc, cx)()
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
