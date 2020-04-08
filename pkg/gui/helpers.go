package gui

import (
	"bytes"
	"gioui.org/app"
	"gioui.org/app/headless"
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gui/gel"
	"github.com/p9c/pod/pkg/gui/gelook"
	"golang.org/x/exp/shiny/iconvg"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"math"
	"os/exec"
	"time"
)

type IconFunc func(ins int) func(hl bool) layout.FlexChild

type WidgetFunc func(hl bool) layout.FlexChild

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
	if s.HW, err = headless.NewWindow(s.WindowWidth,
		s.WindowHeight); Check(err) {
		return
	}
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

func (s *State) FlexV(children ...layout.FlexChild) func(hl bool) {
	return func(hl bool) {
		gtx := s.Gtx
		if hl {
			gtx = s.Htx
		}
		layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
	}
}

func (s *State) FlexH(children ...layout.FlexChild) func(hl bool) {
	return func(hl bool) {
		gtx := s.Gtx
		if hl {
			gtx = s.Htx
		}
		layout.Flex{Axis: layout.Horizontal}.Layout(gtx, children...)
	}
}

func (s *State) Inset(hl bool, size int, fn func()) {
	gtx := s.Gtx
	if hl {
		gtx = s.Htx
	}
	layout.UniformInset(unit.Dp(float32(size))).Layout(gtx, fn)
}

func Rigid(widget func()) layout.FlexChild {
	return layout.Rigid(widget)
}

func Flexed(weight float32, widget func()) layout.FlexChild {
	return layout.Flexed(weight, widget)
}

func (s *State) Spacer() layout.FlexChild {
	return Flexed(1, func() {})
}

func (s *State) Rectangle(width, height int, color string,
	radius ...float32) func(hl bool) func() {
	return func(hl bool) func() {
		return func() {
			gtx := s.Gtx
			if hl {
				gtx = s.Htx
			}
			col := s.Theme.Colors[color]
			var r float32
			if len(radius) > 0 {
				r = radius[0]
			}
			gelook.DuoUIdrawRectangle(gtx,
				width, height, col,
				[4]float32{r, r, r, r},
				[4]float32{0, 0, 0, 0},
			)
		}
	}
}
func (s *State) IconSVGtoImage(icon []byte, fg string, size int) (render IconFunc, err error) {
	var m iconvg.Metadata
	var ico iconvg.Rasterizer
	if m, err = iconvg.DecodeMetadata(icon); Check(err) {
		return
	}
	dx, dy := m.ViewBox.AspectRatio()
	img := image.NewRGBA(
		image.Rectangle{
			Max: image.Point{
				X: size,
				Y: int(float32(size) * dy / dx),
			},
		},
	)
	ico.SetDstImage(img, img.Bounds(), draw.Src)
	m.Palette[0] = gelook.HexARGB(s.Theme.Colors[fg])
	_ = iconvg.Decode(&ico, icon, &iconvg.DecodeOptions{
		Palette: &m.Palette,
	})
	//return
	return func(ins int) func(hl bool) layout.FlexChild {
		return func(hl bool) layout.FlexChild {
			return Rigid(func() {
				gtx := s.Gtx
				if hl {
					gtx = s.Htx
				}
				gtx.Constraints.Width.Min = size
				gtx.Constraints.Height.Min = size
				s.Inset(hl, ins, func() {
					op := paint.NewImageOp(img)
					sz := op.Size()
					op.Add(gtx.Ops)
					paint.PaintOp{
						Rect: f32.Rectangle{
							Max: toPointF(sz),
						},
					}.Add(gtx.Ops)
				})
			})
		}
	}, nil
}

func toPointF(p image.Point) f32.Point {
	return f32.Point{X: float32(p.X), Y: float32(p.Y)}
}

func (s *State) ButtonArea(content func(hl bool) func(),
	hook func(), button *gel.Button) func(hl bool) layout.FlexChild {
	return func(hl bool) layout.FlexChild {
		return Rigid(func() {
			gtx := s.Gtx
			if hl {
				gtx = s.Htx
			}
			b := s.Theme.DuoUIbutton("", "", "",
				"", "", "", "", "",
				0, 0, 0, 0, 0, 0,
				0, 0)
			b.InsideLayout(gtx, button, content(hl))
			for button.Clicked(gtx) {
				hook()
			}
		})
	}
}

func (s *State) Label(hl bool, txt, fg, bg string) {
	gtx := s.Gtx
	if hl {
		gtx = s.Htx
	}
	s.Inset(hl, 10, func() {
		t := s.Theme.DuoUIlabel(unit.Dp(float32(36)), txt)
		t.Color = s.Theme.Colors[fg]
		t.Font.Typeface = s.Theme.Fonts["Secondary"]
		//t.TextSize = unit.Dp(32)
		t.Layout(gtx)
	})
}

func (s *State) Text(txt, fg, face, tag string,
	height int) func(hl bool) layout.FlexChild {
	return func(hl bool) layout.FlexChild {
		gtx := s.Gtx
		if hl {
			gtx = s.Htx
		}
		return Rigid(func() {
			gtx.Constraints.Height.Min = height
			gtx.Constraints.Height.Max = height
			layout.SW.Layout(gtx, func() {
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
				desc.Layout(gtx)
			})
		})
	}
}

func Toggle(b *bool) bool {
	*b = !*b
	return *b
}
