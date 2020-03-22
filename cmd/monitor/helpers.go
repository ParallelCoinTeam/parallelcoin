package monitor

import (
	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/p9c/pod/pkg/gel"
	"github.com/p9c/pod/pkg/gelook"
)

func (s *State) FlexV(children ...layout.FlexChild) () {
	layout.Flex{Axis: layout.Vertical}.Layout(s.Gtx, children...)
}

func (s *State) FlexH(children ...layout.FlexChild) {
	layout.Flex{Axis: layout.Horizontal}.Layout(s.Gtx, children...)
}

func (s *State) Inset(size int, fn func()) {
	layout.UniformInset(unit.Dp(float32(size))).Layout(s.Gtx, fn)
}

func Rigid(widget func()) layout.FlexChild {
	return layout.Rigid(widget)
}

func Flexed(weight float32, widget func()) layout.FlexChild {
	return layout.Flexed(weight, widget)
}

func Spacer() layout.FlexChild {
	return Flexed(1, func() {})
}

func (s *State) Rectangle(width, height int, color, opacity string, radius ...float32) {
	col := s.Theme.Colors[color]
	col = opacity + col[2:]
	var r float32
	if len(radius) > 0 {
		r = radius[0]
	}
	gelook.DuoUIdrawRectangle(s.Gtx,
		width, height, col,
		[4]float32{r, r, r, r},
		[4]float32{0, 0, 0, 0},
	)
}

func (s *State) IconButton(icon, fg, bg string, button *gel.Button, size ...int) {
	sz := 32
	if len(size)>1{
		sz = size[0]
	}
	s.Theme.DuoUIbutton("", "", "",
		s.Theme.Colors[bg], "", s.Theme.Colors[fg], icon,
		s.Theme.Colors[fg], 0, sz, sz+9, sz+9,
		0, 0).IconLayout(s.Gtx, button)
}

func (s *State) TextButton(label, fontFace string, fontSize int, fg, bg string,
	button *gel.Button) {
	s.Theme.DuoUIbutton(
		s.Theme.Fonts[fontFace],
		label,
		s.Theme.Colors[fg],
		s.Theme.Colors[bg],
		s.Theme.Colors[bg],
		s.Theme.Colors[fg],
		"settingsIcon",
		s.Theme.Colors["Light"],
		fontSize, 0, fontSize, fontSize, 0, 0).
		Layout(s.Gtx, button)
}

func (s *State) Text(txt, color, face, tag string) func() {
	return func() {
		var desc gelook.DuoUIlabel
		switch tag {
		case "body1":
			desc = s.Theme.Body1(txt)
		case "body2":
			desc = s.Theme.Body2(txt)
		case "h2":
			desc = s.Theme.H2(txt)
		case "h3":
			desc = s.Theme.H3(txt)
		case "h4":
			desc = s.Theme.H4(txt)
		case "h5":
			desc = s.Theme.H5(txt)
		case "h6":
			desc = s.Theme.H6(txt)
		}
		desc.Font.Typeface = s.Theme.Fonts[face]
		desc.Color = s.Theme.Colors[color]
		desc.Layout(s.Gtx)
	}
}

func Toggle(b *bool) bool {
	*b = !*b
	return *b
}
