package gel

import (
	l "github.com/p9c/gio/layout"
)

type Inset struct {
	*Window
	in    l.Inset
	w     l.Widget
}

// Inset creates a padded empty space around a widget
func (w *Window) Inset(pad float32, embed l.Widget) (out *Inset) {
	out = &Inset{
		Window: w,
		in:     l.UniformInset(w.TextSize.Scale(pad)),
		w:      embed,
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
