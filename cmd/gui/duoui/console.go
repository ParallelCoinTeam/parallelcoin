package duoui

import (
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/text"
	"github.com/p9c/pod/pkg/gio/unit"
	"github.com/p9c/pod/pkg/gio/widget"
	"image/color"
)

var (
	testLabel         = "testtopLabel"
	consoleInputField = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
	consoleOutputList = &layout.List{
		Axis: layout.Vertical,
	}
	ln = layout.UniformInset(unit.Dp(1))
	in = layout.UniformInset(unit.Dp(8))
)

func DuoUIconsole(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) {
	layout.Flex{}.Layout(duo.Gc,
		layout.Flexed(0.9, func() {

			duo.Comp.Console.Inset.Layout(duo.Gc, func() {
				//helpers.DuoUIdrawRect(duo.Gc, duo.Cs.Width.Max, duo.Cs.Height.Max, color.RGBA{A: 0xff, R: 0x30, G: 0xcf, B: 0xcf}, 0, 0, 0, 0, unit.Dp(0))
				// Overview <<<
				layout.Flex{}.Layout(duo.Gc,
					layout.Flexed(1, func() {

						//duo.comp.content.i.Layout(duo.Gc, func() {
						//helpers.DuoUIdrawRect(duo.Gc, duo.Cs.Width.Max, 180, color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}, 0, 0, 0, 0, unit.Dp(0))
						// OverviewTop <<<
						//balance := duo.comp.OverviewTop.Layout.Flex(duo.Gc, 0.4, func() {
						//	helpers.DuoUIdrawRect(duo.Gc, duo.Cs.Width.Max-30, 180, color.RGBA{A: 0xff, R: 0x30, G: 0x30, B: 0x30}, 9.9, 9.9, 9.9, 9.9)
						//	in := layout.UniformInset(unit.Dp(60))
						//
						//	in.Layout(duo.Gc, func() {
						//		bal := duo.th.H3("Balance :" + duo.rc.Balance + " DUO")
						//
						//		bal.Color = color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}
						//		bal.Layout(duo.Gc)
						//	})
						//
						//})

						//duo.comp.OverviewTop.Layout.Layout(duo.Gc, balance, DuoUIsendreceive(duo))
						// OverviewTop >>>
						//})
					}),
				)

				layout.Flex{}.Layout(duo.Gc,
					layout.Rigid(func() {

						helpers.DuoUIdrawRectangle(duo.Gc, duo.Cs.Width.Max, 60, color.RGBA{A: 0xff, R: 0x30, G: 0x30, B: 0x30}, 9.9, 9.9, 9.9, 9.9, unit.Dp(0))
						ln.Layout(duo.Gc, func() {
							helpers.DuoUIdrawRectangle(duo.Gc, duo.Cs.Width.Max, 50, color.RGBA{A: 0xff, R: 0xf4, G: 0xf4, B: 0xf4}, 9.9, 9.9, 9.9, 9.9, unit.Dp(0))
							in.Layout(duo.Gc, func() {
								e := duo.Th.Editor("Run command")
								e.Font.Style = text.Regular
								e.Font.Size = unit.Dp(24)
								e.Layout(duo.Gc, consoleInputField)
								for _, e := range consoleInputField.Events(duo.Gc) {
									if e, ok := e.(widget.SubmitEvent); ok {
										testLabel = e.Text
										consoleInputField.SetText("")
									}
								}
							})
						})

						//duo.comp.OverviewBottom.Layout.Layout(duo.Gc, transactions, status)
						// OverviewBottom >>>
						//})

					}),
				)
				// Overview >>>
			})
		}),
	)
	//return duo.comp.Content.Layout.Rigid(duo.Gc, func() {
	//	//helpers.DuoUIdrawRect(duo.Gc, duo.Cs.Width.Max, 64, color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}, 0, 0, 0, 0, unit.Dp(0))
	//	// Header <<<
	//	consoleOut := duo.comp.ConsoleOutput.Layout.Rigid(duo.Gc, func() {
	//		//helpers.DuoUIdrawRect(duo.Gc, 64, 64, color.RGBA{A: 0xff, R: 0x30, B: 0x30, G: 0x30}, 0, 0, 0, 0, unit.Dp(0))
	//
	//
	//
	//	})
	//
	//	// Header >>>
	//})
}
