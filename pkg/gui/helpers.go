package gui

import (
	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gui/gel"
	"github.com/p9c/pod/pkg/gui/gelook"
)

type State struct {
	Gtx   *layout.Context
	W     *app.Window
	Rc    *rcd.RcVar
	Theme *gelook.DuoUItheme
}

func (s *State) FlexV(children ...layout.FlexChild) {
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

func (s *State) Spacer() layout.FlexChild {
	return Flexed(1, func() { s.Inset(0, func() {}) })
}

func (s *State) Rectangle(width, height int, color, opacity string, radius ...float32) {
	col := s.Theme.Colors[color]
	if col == "" || col[:2] == "00" {
		return
	}
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

func (s *State) Icon(icon, fg, bg string, size int) {
	s.Gtx.Constraints.Width.Max = size
	s.Gtx.Constraints.Height.Max = size
	s.Gtx.Constraints.Width.Min = size
	s.Gtx.Constraints.Height.Min = size
	s.FlexH(Rigid(func() {
		cs := s.Gtx.Constraints
		tp := "ff"
		bg := s.Theme.Colors[bg]
		if len(bg) == 0 {
			bg = "00000000"
			tp = "00"
		}
		s.Rectangle(cs.Height.Max, cs.Width.Max, bg, tp)
		s.Inset(0, func() {
			i := s.Theme.Icons[icon]
			i.Color = gelook.HexARGB(s.Theme.Colors[fg])
			i.Layout(s.Gtx, unit.Dp(float32(size)))
		})
	}),
	)
}

func (s *State) IconButton(icon, fg, bg string, button *gel.Button, size ...int) {
	sz := 48
	if len(size) > 1 {
		sz = size[0]
	}
	s.Rectangle(sz, sz, bg, "ff")
	s.ButtonArea(func() {
		s.Inset(8, func() {
			s.Icon(icon, fg, "Transparent", sz-16)
		})
	}, button)
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
		fontSize, 0, 32, 32, 10, 8, 7, 10).
		Layout(s.Gtx, button)
}

func (s *State) ButtonArea(content func(), button *gel.Button) {
	b := s.Theme.DuoUIbutton("", "", "", "", "", "", "", "", 0, 0, 0, 0, 0, 0,
		0, 0)
	b.InsideLayout(s.Gtx, button, content)
}

func (s *State) Label(txt string) {
	s.Inset(8, func() {
		t := s.Theme.DuoUIlabel(unit.Dp(float32(39)), txt)
		t.Color = s.Theme.Colors["PanelText"]
		t.Font.Typeface = s.Theme.Fonts["Secondary"]
		//t.TextSize = unit.Dp(32)
		t.Layout(s.Gtx)
	})
}

func (s *State) Text(txt, fg, bg, face, tag string) func() {
	return func() {
		var desc gelook.DuoUIlabel
		switch tag {
		case "h1":
			desc = s.Theme.H1(txt)
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
		case "body1":
			desc = s.Theme.Body1(txt)
		case "body2":
			desc = s.Theme.Body2(txt)
		}
		desc.Font.Typeface = s.Theme.Fonts[face]
		desc.Color = s.Theme.Colors[fg]
		s.Inset(8, func() {
			cs := s.Gtx.Constraints
			s.Rectangle(cs.Width.Max, cs.Height.Max, bg, "ff")
			desc.Layout(s.Gtx)
		})

	}
}

func Toggle(b *bool) bool {
	*b = !*b
	return *b
}
