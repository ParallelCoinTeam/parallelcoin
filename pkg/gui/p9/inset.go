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
func (th *Theme) Inset(pad float32) (out *_inset) {
	out = &_inset{
		in: l.UniformInset(unit.Sp(pad)),
	}
	return
}

func (in *_inset) Widget(w l.Widget) *_inset {
	in.w = w
	return in
}

// Fn the given widget with the configured context and padding
func (in *_inset) Fn(c l.Context) l.Dimensions {
	return in.in.Layout(c, in.w)
}
