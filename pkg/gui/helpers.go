package gui

import (
	"bytes"
	"gioui.org/app"
	"gioui.org/app/headless"
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gui/gel"
	"github.com/p9c/pod/pkg/gui/gelook"
	"image"
	"image/png"
	"io/ioutil"
	"math"
	"os/exec"
	"time"
)

type ScaledConfig struct {
	Scale float32
}

func (s *ScaledConfig) Now() time.Time {
	return time.Now()
}

func (s *ScaledConfig) Px(v unit.Value) int {
	scale := s.Scale
	if v.U == unit.UnitPx {
		scale = 1
	}
	return int(math.Round(float64(scale * v.V)))
}

// State stores the state for a gui
type State struct {
	Gtx   *layout.Context
	Htx   *layout.Context
	W     *app.Window
	HW    *headless.Window
	Rc    *rcd.RcVar
	Theme *gelook.DuoUItheme
	// these two values need to be updated by the main render pipeline loop
	WindowWidth, WindowHeight int
	DarkTheme                 bool
	ScreenShooting            bool
}

func (s *State) Screenshot(widget func(),
	path string) (err error) {
	Debug("capturing screenshot")
	s.ScreenShooting = true
	sz := image.Point{X: s.WindowWidth, Y: s.WindowHeight}
	s.Htx.Reset(&ScaledConfig{1}, sz)
	widget()
	s.HW.Frame(s.Htx.Ops)
	var img *image.RGBA
	if img, err = s.HW.Screenshot(); Check(err) {
	}
	Debug("image captured", len(img.Pix))
	//Debugs(img)
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
	}
	Debug("png", buf.Len())
	b64 := buf.Bytes()
	if err := ioutil.WriteFile(path, b64, 0600); !Check(err) {
		cmd := exec.Command("chromium", path)
		err = cmd.Run()
	}
	//Debug("bytes", len(b64))
	//clip := make([]byte, len(b64)*2)
	//base64.StdEncoding.Encode(clip, b64)
	//Debug("clip", len(clip))
	//st := "data:image/png;base64," + string(clip)
	//Debug(st)
	//if cmdIn, err := cmd.StdinPipe(); !Check(err) {
	//	cmdIn.Write([]byte(st))
	//}
	//clipboard.Set(st)
	//time.Sleep(time.Second / 2)
	s.ScreenShooting = false
	return
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

func (s *State) Spacer(bg string) layout.FlexChild {
	return Flexed(1, func() {
		//cs := s.Gtx.Constraints
		//s.Rectangle(cs.Width.Max, cs.Height.Max, bg, "FF")
	})
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

func (s *State) Label(txt, fg, bg string) {
	s.Gtx.Constraints.Height.Max = 48
	cs := s.Gtx.Constraints
	s.Rectangle(cs.Width.Max, cs.Height.Max, bg, "ff")
	s.Inset(10, func() {
		t := s.Theme.DuoUIlabel(unit.Dp(float32(36)), txt)
		t.Color = s.Theme.Colors[fg]
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
