package main

import (
	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/p9c/pod/pkg/gel"
	"github.com/p9c/pod/pkg/gelook"
)

func main() {
	theme := gelook.NewDuoUItheme()
	var (
		mainList = &layout.List{
			Axis: layout.Vertical,
		}
		startButton = new(gel.Button)
	)
	go func() {
		w := app.NewWindow(
			app.Size(unit.Dp(800), unit.Dp(600)),
			app.Title("ParallelCoin"),
		)
		gtx := layout.NewContext(w.Queue())
		for e := range w.Events() {
			if e, ok := e.(system.FrameEvent); ok {
				gtx.Reset(e.Config, e.Size)

				layout.Flex{
					Axis: layout.Horizontal,
				}.Layout(gtx, layout.Flexed(1, func() {
					cs := gtx.Constraints
					gelook.DuoUIdrawRectangle(gtx, cs.Width.Max,
						cs.Height.Max, theme.Colors["Light"],
						[4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
					layout.Center.Layout(gtx, func() {
						controllers := []func(){
							func() {
								bal := theme.H5(
									"Enter the private passphrase for your new wallet:")
								bal.Font.Typeface = theme.Fonts["Primary"]
								bal.Color = theme.Colors["Dark"]
								bal.Layout(gtx)
							},
							func() {
								var createWalletbuttonComp gelook.DuoUIbutton
								createWalletbuttonComp = theme.DuoUIbutton(theme.
									Fonts["Secondary"], "start wallet",
									theme.Colors["Dark"], theme.Colors["Light"],
									theme.Colors["Light"],
									theme.Colors["Dark"], "",
									theme.Colors["Dark"], 16, 0, 125, 32, 4, 4)
								for startButton.Clicked(gtx) {
								}
								createWalletbuttonComp.Layout(gtx,
									startButton)
							},
						}
						mainList.Layout(gtx, len(controllers), func(i int) {
							layout.UniformInset(unit.Dp(10)).Layout(gtx,
								controllers[i])
						})
					})
				}))
				e.Frame(gtx.Ops)
				w.Invalidate()
			}

		}
	}()
	app.Main()
}
