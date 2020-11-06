package toast

import (
	"gioui.org/f32"
	l "gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"github.com/egonelbre/expgio/surface/f32color"
	"github.com/gioapp/gel/helper"
	"github.com/p9c/pod/pkg/gui/p9"
	icons2 "golang.org/x/exp/shiny/materialdesign/icons"
	"image"
	"image/color"
	"time"
)

type Toasts struct {
	toasts             []toast
	layout             *p9.List
	theme              *p9.Theme
	offset             image.Point
	duration           int
	singleSize         image.Point
	singleCornerRadius unit.Value
	singleElevation    unit.Value
}

type toast struct {
	title, content, level string
	background            color.RGBA
	time                  time.Time
	close                 p9.Clickable
	cornerRadius          unit.Value
	elevation             unit.Value
}

func New(th *p9.Theme) *Toasts {
	return &Toasts{
		layout:             th.List(),
		theme:              th,
		duration:           3000,
		singleSize:         image.Pt(250, 50),
		singleCornerRadius: unit.Dp(5),
		singleElevation:    unit.Dp(5),
	}
}
func (t *Toasts) AddToast(title, content string) {
	t.toasts = append(t.toasts, toast{
		title:        title,
		content:      content,
		time:         time.Now().Add(time.Duration(t.duration) * time.Millisecond),
		background:   helper.HexARGB("ffffffff"),
		cornerRadius: t.singleCornerRadius,
		elevation:    t.singleElevation,
	})
}

func (t *Toasts) DrawToasts(gtx l.Context) {
	t.layout.Vertical().Length(len(t.toasts)).ListElement(t.singleToast).Fn(gtx)
}

func (t *Toasts) singleToast(gtx l.Context, index int) l.Dimensions {
	if t.toasts[index].time != time.Now() {
		gtx.Constraints.Min = image.Pt(t.singleSize.X, t.singleSize.Y)
		gtx.Constraints.Max = image.Pt(t.singleSize.X, t.singleSize.Y)
		sz := gtx.Constraints.Min
		rr := float32(gtx.Px(t.singleCornerRadius))

		r := f32.Rect(0, 0, float32(sz.X), float32(sz.Y))
		t.toasts[index].singleToastLayoutShadow(gtx, r, rr)
		clip.UniformRRect(r, rr).Add(gtx.Ops)

		paint.Fill(gtx.Ops, t.toasts[index].background)

		return l.Flex{Axis: l.Vertical, Alignment: l.Middle}.Layout(gtx,
			l.Flexed(1,
				func(gtx l.Context) l.Dimensions {
					return t.theme.Inset(0.25,
						t.theme.VFlex().
							Rigid(
								t.theme.Fill("PanelBg",
									t.theme.Flex().
										Flexed(1,
											func(gtx l.Context) l.Dimensions {
												switch t.toasts[index].level {
												case "Warning":
													return t.theme.Icon().Color("DocText").Scale(1).Src(&icons2.AlertWarning).Fn(gtx)
												case "Success":
													return t.theme.Icon().Color("DocText").Scale(1).Src(&icons2.NavigationCheck).Fn(gtx)
												case "Danger":
													return t.theme.Icon().Color("DocText").Scale(1).Src(&icons2.AlertError).Fn(gtx)
												case "Info":
													return t.theme.Icon().Color("DocText").Scale(1).Src(&icons2.ActionInfo).Fn(gtx)
												}
												return l.Dimensions{}
											},
										).
										Rigid(
											t.theme.H6(t.toasts[index].title).Color("PanelText").Fn,
										).Fn,
								).Fn,
							).
							Rigid(
								t.theme.Body1(t.toasts[index].content).Color("PanelText").Fn,
							).Fn).Fn(gtx)
				}))
	} else {
		t.toasts = remove(t.toasts, index)
		return p9.EmptySpace(0, 0)(gtx)
	}
}

func (t *toast) singleToastLayoutShadow(gtx l.Context, r f32.Rectangle, rr float32) {
	if t.elevation.V <= 0 {
		return
	}

	offset := pxf(gtx.Metric, t.elevation)

	d := int(offset + 1)
	if d > 4 {
		d = 4
	}

	a := float32(t.background.A) / 0xff
	background := (f32color.RGBA{A: a * 0.4 / float32(d*d)}).SRGB()
	for x := 0; x <= d; x++ {
		for y := 0; y <= d; y++ {
			px, py := float32(x)/float32(d)-0.5, float32(y)/float32(d)-0.15
			stack := op.Push(gtx.Ops)
			op.Offset(f32.Pt(px*offset, py*offset)).Add(gtx.Ops)
			clip.UniformRRect(r, rr).Add(gtx.Ops)
			paint.Fill(gtx.Ops, background)
			stack.Pop()
		}
	}
}

func outset(r f32.Rectangle, y, s float32) f32.Rectangle {
	r.Min.X += s
	r.Min.Y += s + y
	r.Max.X += -s
	r.Max.Y += -s + y
	return r
}

func pxf(c unit.Metric, v unit.Value) float32 {
	switch v.U {
	case unit.UnitPx:
		return v.V
	case unit.UnitDp:
		s := c.PxPerDp
		if s == 0 {
			s = 1
		}
		return s * v.V
	case unit.UnitSp:
		s := c.PxPerSp
		if s == 0 {
			s = 1
		}
		return s * v.V
	default:
		panic("unknown unit")
	}
}

func remove(slice []toast, s int) []toast {
	return append(slice[:s], slice[s+1:]...)
}
