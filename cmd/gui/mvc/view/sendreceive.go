package view

import (
	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"github.com/p9c/pod/cmd/gui/mvc/model"
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
	list = &layout.List{
		Axis: layout.Vertical,
	}
)

type DuoUIsend struct {
	*model.DuOScomponent
}

func DuoCOMsend() *DuoUIsend {
	send := *new(DuoUIsend)

	cs := &model.DuOScomponent{
		Name:    "logger",
		Version: "0.1",
		//Model:      ,
		//Controller: c,
	}

	*send.DuOScomponent = *cs

	return &send
}

func (c *DuoUIsend) View(gtx *layout.Context, th *theme.DuoUItheme) func() {
	return func() {
		layout.Flex{}.Layout(gtx,
			layout.Rigid(func() {
				cs := gtx.Constraints
				theme.DuoUIdrawRectangle(gtx, cs.Width.Max, 180, th.Color.Bg, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})

				layout.Flex{
					Axis: layout.Vertical,
				}.Layout(gtx,
					layout.Rigid(func() {
						layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
							cs := gtx.Constraints
							theme.DuoUIdrawRectangle(gtx, cs.Width.Max, 32, "fff4f4f4", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
							layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
								theme.DuoUIdrawRectangle(gtx, cs.Width.Max, 30, "ffffffff", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
								e := th.DuoUIeditor("DUO address", "DUO dva")
								e.Font.Style = text.Italic
								e.Font.Size = unit.Dp(24)
								e.Layout(gtx, addressLineEditor)
								for _, e := range addressLineEditor.Events(gtx) {
									if e, ok := e.(controller.SubmitEvent); ok {
										topLabel = e.Text
										addressLineEditor.SetText("")
									}
								}
							})
						})
					}),
					layout.Rigid(func() {
						layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
							cs := gtx.Constraints
							theme.DuoUIdrawRectangle(gtx, cs.Width.Max, 32, "fff4f4f4", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
							layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
								theme.DuoUIdrawRectangle(gtx, cs.Width.Max, 30, "ffffffff", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
								e := th.DuoUIeditor("DUO Amount", "DUO dva")
								e.Font.Style = text.Italic
								e.Font.Size = unit.Dp(24)
								e.Layout(gtx, amountLineEditor)
								for _, e := range amountLineEditor.Events(gtx) {
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
