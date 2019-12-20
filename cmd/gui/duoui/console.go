package duoui

import (
	"github.com/p9c/pod/cmd/gui/helpers"
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
	//ln = layout.UniformInset(unit.Dp(1))
	//in = layout.UniformInset(unit.Dp(8))
)

func DuoUIconsole(duo *DuoUI) layout.FlexChild {
	return duo.comp.Content.Layout.Flex(duo.gc, 0.9, func() {

		duo.comp.Console.Inset.Layout(duo.gc, func() {
			//helpers.DuoUIdrawRect(duo.gc, duo.cs.Width.Max, duo.cs.Height.Max, color.RGBA{A: 0xff, R: 0x30, G: 0xcf, B: 0xcf}, 0, 0, 0, 0)
			// Overview <<<
			consoleOut := duo.comp.Console.Layout.Flex(duo.gc, 1, func() {

				//duo.comp.content.i.Layout(duo.gc, func() {
				//helpers.DuoUIdrawRect(duo.gc, duo.cs.Width.Max, 180, color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}, 0, 0, 0, 0)
				// OverviewTop <<<
				//balance := duo.comp.OverviewTop.Layout.Flex(duo.gc, 0.4, func() {
				//	helpers.DuoUIdrawRect(duo.gc, duo.cs.Width.Max-30, 180, color.RGBA{A: 0xff, R: 0x30, G: 0x30, B: 0x30}, 9.9, 9.9, 9.9, 9.9)
				//	in := layout.UniformInset(unit.Dp(60))
				//
				//	in.Layout(duo.gc, func() {
				//		bal := duo.th.H3("Balance :" + duo.rc.Balance + " DUO")
				//
				//		bal.Color = color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}
				//		bal.Layout(duo.gc)
				//	})
				//
				//})

				//duo.comp.OverviewTop.Layout.Layout(duo.gc, balance, DuoUIsendreceive(duo))
				// OverviewTop >>>
				//})
			})
			consoleIn := duo.comp.Console.Layout.Rigid(duo.gc, func() {

				helpers.DuoUIdrawRect(duo.gc, duo.cs.Width.Max, 60, color.RGBA{A: 0xff, R: 0x30, G: 0x30, B: 0x30}, 9.9, 9.9, 9.9, 9.9)
				ln.Layout(duo.gc, func() {
					helpers.DuoUIdrawRect(duo.gc, duo.cs.Width.Max, 50, color.RGBA{A: 0xff, R: 0xf4, G: 0xf4, B: 0xf4}, 9.9, 9.9, 9.9, 9.9)
					in.Layout(duo.gc, func() {
						e := duo.th.Editor("Run command")
						e.Font.Style = text.Regular
						e.Font.Size = unit.Dp(24)
						e.Layout(duo.gc, consoleInputField)
						for _, e := range consoleInputField.Events(duo.gc) {
							if e, ok := e.(widget.SubmitEvent); ok {
								testLabel = e.Text
								consoleInputField.SetText("")
							}
						}
					})
				})

				//duo.comp.OverviewBottom.Layout.Layout(duo.gc, transactions, status)
				// OverviewBottom >>>
				//})
			})
			duo.comp.Console.Layout.Layout(duo.gc, consoleOut, consoleIn)
			// Overview >>>
		})
	})
	//return duo.comp.Content.Layout.Rigid(duo.gc, func() {
	//	//helpers.DuoUIdrawRect(duo.gc, duo.cs.Width.Max, 64, color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}, 0, 0, 0, 0)
	//	// Header <<<
	//	consoleOut := duo.comp.ConsoleOutput.Layout.Rigid(duo.gc, func() {
	//		//helpers.DuoUIdrawRect(duo.gc, 64, 64, color.RGBA{A: 0xff, R: 0x30, B: 0x30, G: 0x30}, 0, 0, 0, 0)
	//
	//
	//
	//	})
	//
	//	// Header >>>
	//})
}
