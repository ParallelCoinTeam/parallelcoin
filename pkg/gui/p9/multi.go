package p9

import (
	l "gioui.org/layout"
)

type Multi struct {
	*Theme
}

func (th *Theme) Multiline(txt []string, borderColorFocused, borderColorUnfocused string,
	size int, handle func(txt []string)) (m *Multi) {

	return m
}

func (m *Multi) Fn(gtx l.Context) l.Dimensions {
	return l.Dimensions{}
}
