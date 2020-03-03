package component

import (
	"gioui.org/layout"
	"github.com/p9c/pod/pkg/gui/theme"
)

func HorizontalLine(gtx *layout.Context, height int, color string) func() {
	return func() {
		cs := gtx.Constraints
		theme.DuoUIdrawRectangle(gtx, cs.Width.Max, height, color, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
	}
}
