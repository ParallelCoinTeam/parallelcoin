package p9

import (
	"gioui.org/layout"
	"gioui.org/unit"
)

type inset struct {
	in layout.Inset
	w  layout.Widget
}

// Inset creates a padded empty space around a widget
func (th *Theme) Inset(pad int, w layout.Widget) (out *inset) {
	out = &inset{
		in: layout.UniformInset(unit.Dp(float32(pad))),
		w:  w,
	}
	return
}

// Fn the given widget with the configured context and padding
func (in *inset) Fn(gtx layout.Context) layout.Dimensions {
	return in.in.Layout(gtx, in.w)
}
