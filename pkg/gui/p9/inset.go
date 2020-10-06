package p9

import (
	l "gioui.org/layout"
	"gioui.org/unit"
)

type _inset struct {
	in l.Inset
	w  l.Widget
}

// Inset creates a padded empty space around a widget
func (th *Theme) Inset(pad int, w l.Widget) (out *_inset) {
	out = &_inset{
		in: l.UniformInset(unit.Dp(float32(pad))),
		w:  w,
	}
	return
}

// Fn the given widget with the configured context and padding
func (in *_inset) Fn(gtx l.Context) l.Dimensions {
	return in.in.Layout(gtx, in.w)
}
