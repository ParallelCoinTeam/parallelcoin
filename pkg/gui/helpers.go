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
	Theme *gelook.DuoUITheme
}

func (s *State) FlexV(children ...layout.FlexChild) {
	layout.Flex{Axis: layout.Vertical}.Layout(s.Gtx, children...)
}

func (s *State) FlexH(children ...layout.FlexChild) {
	layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(s.Gtx,
		children...)
}

func (s *State) FlexHStart(children ...layout.FlexChild) {
	layout.Flex{Axis: layout.Horizontal, Alignment: layout.Start}.Layout(s.Gtx,
		children...)
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

func (s *State) Spacer(bg string) layout.FlexChild {
	return Flexed(1, func() {
		// cs := s.Gtx.Constraints
		// s.Rectangle(cs.Width.Max, cs.Height.Max, bg, "FF")
	})
}

func (s *State) Rectangle(width, height int, color string, radius ...float32) {
	col := s.Theme.Colors[color]
	if col == "" {
		return
	}
	var r float32
	if len(radius) > 0 {
		r = radius[0]
	}
	gelook.DuoUIDrawRectangle(s.Gtx,
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
		bg := s.Theme.Colors[bg]
		if len(bg) == 0 {
			bg = "00000000"
		}
		s.Rectangle(cs.Height.Max, cs.Width.Max, bg)
		// s.Inset(0, func() {
		i := s.Theme.Icons[icon]
		// Debug(fg)
		// _ = fg
		i.Color = gelook.HexARGB(s.Theme.Colors[fg])
		i.Layout(s.Gtx, unit.Dp(float32(size)))
		// })
	}),
	)
}

func (s *State) IconButton(icon, fg, bg string, button *gel.Button, size ...int) {
	sz := 48
	if len(size) > 1 {
		sz = size[0]
	}
	s.Rectangle(sz, sz, bg)
	s.ButtonArea(func() {
		s.Inset(8, func() {
			s.Icon(icon, fg, "Transparent", sz-16)
		})
	}, button)
}

func (s *State) TextButton(label, fontFace string, fontSize int, fg, bg string,
	button *gel.Button) {
	s.Theme.DuoUIbutton(gelook.ButtonParams{
		TxtFont:       s.Theme.Fonts[fontFace],
		Txt:           label,
		TxtColor:      s.Theme.Colors[fg],
		BgColor:       s.Theme.Colors[bg],
		TxtHoverColor: s.Theme.Colors[bg],
		BgHoverColor:  s.Theme.Colors[fg],
		Icon:          "settingsIcon",
		IconColor:     s.Theme.Colors["Light"],
		TextSize:      fontSize,
		Width:         32,
		Height:        32,
		PaddingTop:    10,
		PaddingRight:  8,
		PaddingBottom: 7,
		PaddingLeft:   10,
	}).Layout(s.Gtx, button)
}

func (s *State) ButtonArea(content func(), button *gel.Button) {
	b := s.Theme.DuoUIbutton(gelook.ButtonParams{})
	b.InsideLayout(s.Gtx, button, content)
}

func (s *State) Label(txt, fg, bg string) {
	s.Gtx.Constraints.Height.Max = 48
	cs := s.Gtx.Constraints
	s.Rectangle(cs.Width.Max, cs.Height.Max, bg)
	s.Inset(10, func() {
		t := s.Theme.DuoUILabel(unit.Dp(float32(36)), txt)
		t.Color = s.Theme.Colors[fg]
		t.Font.Typeface = s.Theme.Fonts["Secondary"]
		// t.TextSize = unit.Dp(32)
		t.Layout(s.Gtx)
	})
}

func (s *State) Text(txt, fg, bg, face, tag string) func() {
	return func() {
		var desc gelook.DuoUILabel
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
			s.Rectangle(cs.Width.Max, cs.Height.Max, bg)
			desc.Layout(s.Gtx)
		})

	}
}

func Toggle(b *bool) bool {
	*b = !*b
	return *b
}
