package p9

import (
	l "gioui.org/layout"
	"github.com/urfave/cli"
)

type Multi struct {
	*Theme
	lines      *cli.StringSlice
	clickables []*Clickable
	buttons    []*Button
}

func (th *Theme) Multiline(txt *cli.StringSlice, borderColorFocused, borderColorUnfocused string,
	size int, handle func(txt []string)) (m *Multi) {
	m = &Multi{Theme: th, lines: txt}
	return m
}

func (m *Multi) Fn(gtx l.Context) l.Dimensions {
	return l.Dimensions{}
}
