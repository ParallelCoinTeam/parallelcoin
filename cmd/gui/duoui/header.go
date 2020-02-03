package duoui

import (
	"fmt"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/text"
	"github.com/p9c/pod/pkg/gui/unit"
	"github.com/p9c/pod/pkg/gui/widget"
	"github.com/p9c/pod/pkg/gui/widget/parallel"
	"image/color"
)

var (
	inLogo     = layout.Stack{Alignment: layout.Center}
	logoButton = new(widget.Button)
)

func DuoUIheader(duo *models.DuoUI, rc *rcd.RcVar) func(){
	return func() {
		// Header <<<
		duo.DuoUIcomponents.Header.Layout.Layout(duo.DuoUIcontext,
			layout.Rigid(func() {
				layout.Align(layout.Center).Layout(duo.DuoUIcontext, func() {
					layout.Inset{Top: unit.Dp(0), Bottom: unit.Dp(0), Left: unit.Dp(0), Right: unit.Dp(0)}.Layout(duo.DuoUIcontext, func() {

						logo := parallel.DuoUIlogo{
							Background: color.RGBA{},
							Color:      color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf},
							Icon:       duo.DuoUIico["Logo"],
							Size:       unit.Dp(96),
							Padding:    unit.Dp(8),
						}
						logo.Layout(duo.DuoUIcontext, logoButton)
						//
						//duo.Ico.Logo.Color = color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}
						//duo.Ico.Logo.Layout(duo.DuoUIcontext, unit.Dp(48))
					})
				})

			}),
			layout.Flexed(1, func() {
				layout.Align(layout.Start).Layout(duo.DuoUIcontext, func() {
					layout.Inset{Top: unit.Dp(24), Bottom: unit.Dp(8), Left: unit.Dp(0), Right: unit.Dp(4)}.Layout(duo.DuoUIcontext, func() {
						currentPage := duo.DuoUItheme.H4(duo.CurrentPage)
						currentPage.Color = color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}
						currentPage.Alignment = text.Start
						currentPage.Layout(duo.DuoUIcontext)
					})
				})

			}),
			layout.Rigid(func() {
				layout.Align(layout.Center).Layout(duo.DuoUIcontext, func() {
					layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(16), Left: unit.Dp(16), Right: unit.Dp(4)}.Layout(duo.DuoUIcontext, func() {
						balance := duo.DuoUItheme.Body2(rc.Balance + " " + duo.DuoUIconfiguration.Abbrevation)
						balance.Color = color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}
						balance.Alignment = text.End
						balance.Layout(duo.DuoUIcontext)
					})
				})

			}),
			layout.Rigid(func() {
				layout.Align(layout.Center).Layout(duo.DuoUIcontext, func() {
					layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(16), Left: unit.Dp(16), Right: unit.Dp(4)}.Layout(duo.DuoUIcontext, func() {
						balance := duo.DuoUItheme.Body2("dimenzion: " + fmt.Sprint(duo.DuoUIcontext.Constraints.Width.Max))
						balance.Color = color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}
						balance.Alignment = text.End
						balance.Layout(duo.DuoUIcontext)
					})
				})

			}))

		// Header >>>

	}
}
