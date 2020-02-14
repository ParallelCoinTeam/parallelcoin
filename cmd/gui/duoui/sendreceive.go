package duoui

import (
	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/text"
	"github.com/p9c/pod/pkg/gui/unit"
)

var (
	topLabel          = "testtopLabel"
	addressLineEditor = &controller.Editor{
		SingleLine: true,
		Submit:     true,
	}
	amountLineEditor = &controller.Editor{
		SingleLine: true,
		Submit:     true,
	}
)

func (ui *DuoUI) DuoUIsend() func() {
	return func() {
		layout.Flex{}.Layout(ui.ly.Context,
			layout.Rigid(func() {
				cs := ui.ly.Context.Constraints
				theme.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, 180, ui.ly.Theme.Color.Bg, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
				layout.Flex{
					Axis: layout.Vertical,
				}.Layout(ui.ly.Context,
					layout.Rigid(func() {
						layout.UniformInset(unit.Dp(0)).Layout(ui.ly.Context, func() {
							cs := ui.ly.Context.Constraints
							theme.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, 32, "fff4f4f4", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
							layout.UniformInset(unit.Dp(0)).Layout(ui.ly.Context, func() {
								theme.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, 30, "ffffffff", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
								e := ui.ly.Theme.DuoUIeditor("DUO address", "DUO dva")
								e.Font.Style = text.Italic
								e.Font.Size = unit.Dp(24)
								e.Layout(ui.ly.Context, addressLineEditor)
								for _, e := range addressLineEditor.Events(ui.ly.Context) {
									if e, ok := e.(controller.SubmitEvent); ok {
										topLabel = e.Text
										addressLineEditor.SetText("")
									}
								}
							})
						})
					}),
					layout.Rigid(func() {
						layout.UniformInset(unit.Dp(0)).Layout(ui.ly.Context, func() {
							cs := ui.ly.Context.Constraints
							theme.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, 32, "fff4f4f4", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
							layout.UniformInset(unit.Dp(0)).Layout(ui.ly.Context, func() {
								theme.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, 30, "ffffffff", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
								e := ui.ly.Theme.DuoUIeditor("DUO Amount", "DUO dva")
								e.Font.Style = text.Italic
								e.Font.Size = unit.Dp(24)
								e.Layout(ui.ly.Context, amountLineEditor)
								for _, e := range amountLineEditor.Events(ui.ly.Context) {
									if e, ok := e.(controller.SubmitEvent); ok {
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
