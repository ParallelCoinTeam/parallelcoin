package fui

import (
	"gioui.org/layout"
	"gioui.org/unit"
)

type inset struct {
	in layout.Inset
}

// Inset creates a padded empty space around a widget
func Inset(pad int) (out *inset) {
	out = &inset{
		in: layout.UniformInset(unit.Dp(float32(pad))),
	}
	return
}

// Layout the given widget with the configured context and padding
func (in *inset) Layout(c *layout.Context, w layout.Widget) {
	in.in.Layout(c, w)
}

// Prep the given widget with the configured context and padding
func (in *inset) Prep(c *layout.Context, w layout.Widget) func() {
	return func() { in.in.Layout(c, w) }
}
