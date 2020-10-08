package p9

import (
	l "gioui.org/layout"
)

type _filler struct {
	th  *Theme
	col string
	w   l.Widget
}

// Fill fills underneath a widget you can put another widget over top
func (th *Theme) Fill(col string) *_filler {
	return &_filler{th: th, col: col}
}

// Widget sets the widget to draw over top
func (f *_filler) Widget(w l.Widget) *_filler {
	f.w = w
	return f
}

// Fn renders the fill and then the widget over top
func (f *_filler) Fn(gtx l.Context) l.Dimensions {
	Fill(gtx, f.th.Colors.Get(f.col))
	return f.w(gtx)
}
