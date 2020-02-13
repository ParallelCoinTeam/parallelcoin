// SPDX-License-Identifier: Unlicense OR MIT

package theme

import (
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/unit"
	"time"
)

var ()

type
	Command struct {
		Com      interface{}
		ComID    string
		Category string
		Out      func()
		Time     time.Time
	}

func (t *DuoUItheme) Command(name string) *DuoUIcomponent {
	return &DuoUIcomponent{
		Name: name,
	}
}

func (p Command) Layout(gtx *layout.Context, f func()) {
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
