package duoui

import (
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/cmd/gui/widget"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/unit"
)

var (
	testLabel         = "testtopLabel"
	consoleInputField = &widget.DuoUIeditor{
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
	layout.Flex{}.Layout(duo.DuoUIcontext,
		layout.Flexed(0.9, func() {

			duo.DuoUIcomponents.Console.Inset.Layout(duo.DuoUIcontext, func() {
				//helpers.DuoUIdrawRect(duo.DuoUIcontext, duo.Cs.Width.Max, duo.Cs.Height.Max, helpers.HexARGB("ff30cfcf"), [4]float32{0, 0, 0, 0}, unit.Dp(0))
				// Overview <<<
				layout.Flex{}.Layout(duo.DuoUIcontext,
					layout.Flexed(1, func() {

						//duo.comp.content.i.Layout(duo.DuoUIcontext, func() {
						//helpers.DuoUIdrawRect(duo.DuoUIcontext, duo.Cs.Width.Max, 180, helpers.HexARGB("ffcfcfcf"), [4]float32{0, 0, 0, 0}, unit.Dp(0))
						// OverviewTop <<<
						//balance := duo.comp.OverviewTop.Layout.Flex(duo.DuoUIcontext, 0.4, func() {
						//	helpers.DuoUIdrawRect(duo.DuoUIcontext, duo.Cs.Width.Max-30, 180, color.RGBA{A: 0xff, R: 0x30, G: 0x30, B: 0x30}, 9.9, 9.9, 9.9, 9.9)
						//	in := layout.UniformInset(unit.Dp(60))
						//
						//	in.Layout(duo.DuoUIcontext, func() {
						//		bal := duo.th.H3("Balance :" + duo.rc.Balance + " DUO")
						//
						//		bal.Color = color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}
						//		bal.Layout(duo.DuoUIcontext)
						//	})
						//
						//})

						//duo.comp.OverviewTop.Layout.Layout(duo.DuoUIcontext, balance, DuoUIsendreceive(duo))
						// OverviewTop >>>
						//})
					}),
				)

				layout.Flex{}.Layout(duo.DuoUIcontext,
					layout.Rigid(func() {

						//helpers.DuoUIdrawRectangle(duo.DuoUIcontext, duo.Cs.Width.Max, 60, helpers.HexARGB("ff303030"), [4]float32{0, 0, 0, 0}, unit.Dp(0))
						//ln.Layout(duo.DuoUIcontext, func() {
						//	helpers.DuoUIdrawRectangle(duo.DuoUIcontext, duo.Cs.Width.Max, 50, helpers.HexARGB("fff4f4f4"), [4]float32{0, 0, 0, 0}, unit.Dp(0))
						//	in.Layout(duo.DuoUIcontext, func() {
						//		e := duo.Th.DuoUIeditor("Run command", "Run txt")
						//		e.Font.Style = text.Regular
						//		e.Font.Size = unit.Dp(24)
						//		e.Layout(duo.DuoUIcontext, consoleInputField)
						//		for _, e := range consoleInputField.Events(duo.DuoUIcontext) {
						//			if e, ok := e.(widget.SubmitEvent); ok {
						//				testLabel = e.Text
						//				consoleInputField.SetText("")
						//			}
						//		}
						//	})
						//})

						//duo.comp.OverviewBottom.Layout.Layout(duo.DuoUIcontext, transactions, status)
						// OverviewBottom >>>
						//})

					}),
				)
				// Overview >>>
			})
		}),
	)
	//return duo.comp.Content.Layout.Rigid(duo.DuoUIcontext, func() {
	//	//helpers.DuoUIdrawRect(duo.DuoUIcontext, duo.Cs.Width.Max, 64, helpers.HexARGB("ffcfcfcf"), [4]float32{0, 0, 0, 0}, unit.Dp(0))
	//	// Header <<<
	//	consoleOut := duo.comp.ConsoleOutput.Layout.Rigid(duo.DuoUIcontext, func() {
	//		//helpers.DuoUIdrawRect(duo.DuoUIcontext, 64, 64, helpers.HexARGB("ff303030"), [4]float32{0, 0, 0, 0}, unit.Dp(0))
	//
	//
	//
	//	})
	//
	//	// Header >>>
	//})
}
