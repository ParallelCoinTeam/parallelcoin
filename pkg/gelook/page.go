// SPDX-License-Identifier: Unlicense OR MIT

package gelook

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
)

var ()

type DuoUIpage struct {
	Title       string
	TxColor     string
	Font        text.Font
	BgColor     string
	BorderColor string
	Border      float32
	shaper      text.Shaper
	Command     func()
	header      func()
	body        func()
	footer      func()
}

func (t *DuoUItheme) DuoUIpage(txt string, border float32, command, header, body, footer func()) *DuoUIpage {
	return &DuoUIpage{
		Title: txt,
		Font: text.Font{
			// Size: t.TextSize.Scale(14.0 / 16.0),
		},
		TxColor:     t.Colors["Dark"],
		BgColor:     t.Colors["Light"],
		BorderColor: t.Colors["LightGrayI"],
		Border:      border,
		shaper:      t.Shaper,
		Command:     command,
		header:      header,
		body:        body,
		footer:      footer,
	}
}

func (p DuoUIpage) Layout(gtx *layout.Context) {
	layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(p.header),
		layout.Flexed(1, func() {
			cs := gtx.Constraints
			DuoUIdrawRectangle(gtx, cs.Width.Max, cs.Height.Max, p.BorderColor, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
			layout.UniformInset(unit.Dp(p.Border)).Layout(gtx, p.body)
		}),
		layout.Rigid(p.footer),
	)
}
