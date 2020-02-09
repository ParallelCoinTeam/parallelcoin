package components

import (
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/text"
	"github.com/p9c/pod/pkg/gui/unit"
	"github.com/p9c/pod/pkg/gui/widget"
)

var (
	topLabel          = "testtopLabel"
	addressLineEditor = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
	amountLineEditor = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
	list = &layout.List{
		Axis: layout.Vertical,
	}
)

func DuoUIsend(duo *models.DuoUI) func() {
	return func() {
		layout.Flex{}.Layout(duo.DuoUIcontext,
			layout.Rigid(func() {
				cs := duo.DuoUIcontext.Constraints
				helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, 180, duo.DuoUItheme.Color.Bg, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})

				layout.Flex{
					Axis: layout.Vertical,
				}.Layout(duo.DuoUIcontext,
					layout.Rigid(func() {
						layout.UniformInset(unit.Dp(0)).Layout(duo.DuoUIcontext, func() {
							cs := duo.DuoUIcontext.Constraints
							helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, 32, helpers.HexARGB("fff4f4f4"), [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
							layout.UniformInset(unit.Dp(0)).Layout(duo.DuoUIcontext, func() {
								helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, 30, helpers.HexARGB("ffffffff"), [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
								e := duo.DuoUItheme.DuoUIeditor("DUO address", "DUO dva")
								e.Font.Style = text.Italic
								e.Font.Size = unit.Dp(24)
								e.Layout(duo.DuoUIcontext, addressLineEditor)
								for _, e := range addressLineEditor.Events(duo.DuoUIcontext) {
									if e, ok := e.(widget.SubmitEvent); ok {
										topLabel = e.Text
										addressLineEditor.SetText("")
									}
								}
							})
						})
					}),
					layout.Rigid(func() {
						layout.UniformInset(unit.Dp(0)).Layout(duo.DuoUIcontext, func() {
							cs := duo.DuoUIcontext.Constraints
							helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, 32, helpers.HexARGB("fff4f4f4"), [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
							layout.UniformInset(unit.Dp(0)).Layout(duo.DuoUIcontext, func() {
								helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, 30, helpers.HexARGB("ffffffff"), [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
								e := duo.DuoUItheme.DuoUIeditor("DUO Amount", "DUO dva")
								e.Font.Style = text.Italic
								e.Font.Size = unit.Dp(24)
								e.Layout(duo.DuoUIcontext, amountLineEditor)
								for _, e := range amountLineEditor.Events(duo.DuoUIcontext) {
									if e, ok := e.(widget.SubmitEvent); ok {
										topLabel = e.Text
										amountLineEditor.SetText("")
									}
								}
							})
						})
					}))
			}),
		)
	}
}
