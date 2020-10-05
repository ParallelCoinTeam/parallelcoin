package fui

import (
	"gioui.org/layout"
	"gioui.org/unit"
)

type inset struct {
	in layout.Inset
	w  layout.Widget
}

// Inset creates a padded empty space around a widget
func Inset(pad int, w layout.Widget) (out *inset) {
	out = &inset{
		in: layout.UniformInset(unit.Dp(float32(pad))),
		w:  w,
	}
	return
}

// Layout the given widget with the configured context and padding
func (in *inset) Layout(gtx layout.Context) layout.Dimensions {
	return in.in.Layout(gtx, in.w)
}

// Child the given widget with the configured context and padding
func (in *inset) Child(w layout.Widget) func(gtx layout.Context) layout.Dimensions {
	return func(gtx layout.Context) layout.Dimensions {
		return in.in.Layout(gtx, w)
	}
}
