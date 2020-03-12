// SPDX-License-Identifier: Unlicense OR MIT

package gelook

import (
	"gioui.org/layout"
	"gioui.org/unit"
)

var ()

type DuoUIcomponent struct {
	Name    string
	Version string
	M       interface{}
	V       func()
	C       func()
}

func (t *DuoUItheme) DuoUIcomponent(name string) *DuoUIcomponent {
	return &DuoUIcomponent{
		Name: name,
	}
}

func (p DuoUIcomponent) Layout(gtx *layout.Context, f func()) {
	layout.Flex{}.Layout(gtx,
		layout.Flexed(1, func() {
			in := layout.UniformInset(unit.Dp(0))
			in.Layout(gtx, func() {
				cs := gtx.Constraints
				DuoUIdrawRectangle(gtx, cs.Width.Max, cs.Height.Max, "ffacacac", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
				// Overview <<<
				layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
					cs := gtx.Constraints
					DuoUIdrawRectangle(gtx, cs.Width.Max, cs.Height.Max, "ffcfcfcf", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
					f()
				})
				// Overview >>>
			})
		}),
	)
}
