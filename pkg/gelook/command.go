// SPDX-License-Identifier: Unlicense OR MIT

package gelook

import (
	"time"

	"gioui.org/layout"
	"gioui.org/unit"
)

var ()

type Command struct {
	Com      interface{}
	ComID    string
	Category string
	Out      func()
	Time     time.Time
}

func (t *DuoUItheme) Command(name string) *Command {
	return &Command{
		ComID: name,
	}
}

func (p Command) Layout(gtx *layout.Context, f func()) {
	layout.Flex{}.Layout(gtx,
		layout.Flexed(1, func() {
			in := layout.UniformInset(unit.Dp(0))
			in.Layout(gtx, func() {
				cs := gtx.Constraints
				DuoUIdrawRectangle(gtx, cs.Width.Max, cs.Height.Max, "ffacacac", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
				layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
					cs := gtx.Constraints
					DuoUIdrawRectangle(gtx, cs.Width.Max, cs.Height.Max, "ffcfcfcf", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
					f()
				})
			})
		}),
	)
}
