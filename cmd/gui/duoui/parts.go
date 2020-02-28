package duoui

import (
	"gioui.org/layout"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
)

func line(gtx *layout.Context, color string) func() {
	return func() {
		cs := gtx.Constraints
		theme.DuoUIdrawRectangle(gtx, cs.Width.Max, 1, color, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
	}
}
