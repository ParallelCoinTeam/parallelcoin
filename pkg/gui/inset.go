package gui

import (
	l "gioui.org/layout"
)

type Inset struct {
	th *Theme
	in l.Inset
	w  l.Widget
}

// Inset creates a padded empty space around a widget
func (th *Theme) Inset(pad float32, w l.Widget) (out *Inset) {
	out = &Inset{
		th: th,
		in: l.UniformInset(th.TextSize.Scale(pad)),
		w:  w,
	}
	return
}

// Embed sets the widget that will be inside the inset
func (in *Inset) Embed(w l.Widget) *Inset {
	in.w = w
	return in
}

// Fn lays out the given widget with the configured context and padding
func (in *Inset) Fn(c l.Context) l.Dimensions {
	return in.in.Layout(c, in.w)
}
