package p9

import (
	l "gioui.org/layout"
)

type _inset struct {
	th *Theme
	in l.Inset
	w  l.Widget
}

// Inset creates a padded empty space around a widget
func (th *Theme) Inset(pad float32) (out *_inset) {
	out = &_inset{
		th: th,
		in: l.UniformInset(
			th.
				textSize.
				Scale(pad)),
	}
	return
}

// Widget sets the widget that will be inside the inset
func (in *_inset) Widget(w l.Widget) *_inset {
	in.w = w
	return in
}

// Fn lays out the given widget with the configured context and padding
func (in *_inset) Fn(c l.Context) l.Dimensions {
	return in.in.Layout(c, in.w)
}
