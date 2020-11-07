package toast

import (
	"fmt"
	"image"
	"image/color"

	"gioui.org/f32"
	l "gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"github.com/egonelbre/expgio/surface/f32color"
	"github.com/gioapp/gel/helper"
	icons2 "golang.org/x/exp/shiny/materialdesign/icons"

	"github.com/p9c/pod/pkg/gui/p9"
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
	headerBackground      color.RGBA
	bodyBackground        color.RGBA
	icon                  *[]byte
	ticker                float32
	close                 p9.Clickable
	cornerRadius          unit.Value
	elevation             unit.Value
}

func New(th *p9.Theme) *Toasts {
	return &Toasts{
		layout:             th.List(),
		theme:              th,
		duration:           10,
		singleSize:         image.Pt(300, 80),
		singleCornerRadius: unit.Dp(5),
		singleElevation:    unit.Dp(5),
	}
}
func (t *Toasts) AddToast(title, content, level string) {
	ic := &icons2.ActionInfo
	switch level {
	case "Warning":
		ic = &icons2.AlertWarning
	case "Success":
		ic = &icons2.NavigationCheck
	case "Danger":
		ic = &icons2.AlertError
	case "Info":
		ic = &icons2.ActionInfo
	}
	t.toasts = append(t.toasts, toast{
		title:            title,
		content:          content,
		level:            level,
		ticker:           0,
		headerBackground: helper.HexARGB(t.theme.Colors[level]),
		bodyBackground:   helper.HexARGB(t.theme.Colors["PanelBg"]),
		cornerRadius:     t.singleCornerRadius,
		elevation:        t.singleElevation,
		icon:             ic,
	})
}

func (t *Toasts) DrawToasts(gtx l.Context) {
	defer op.Push(gtx.Ops).Pop()
	op.Offset(f32.Pt(float32(gtx.Constraints.Max.X)-250, 0)).Add(gtx.Ops)
	gtx.Constraints.Min = image.Pt(250, gtx.Constraints.Max.X)
	gtx.Constraints.Max.X = 250
	t.layout.Vertical().ScrollToEnd().Length(len(t.toasts)).ListElement(t.singleToast).Fn(gtx)
}

func (t *Toasts) singleToast(gtx l.Context, index int) l.Dimensions {
	fmt.Println("Tic:", t.toasts[index].ticker)
	fmt.Println("duration:", t.duration)

	if t.toasts[index].ticker < float32(t.duration) {
		t.toasts[index].ticker += 0.1
		gtx.Constraints.Min = t.singleSize
		gtx.Constraints.Max = t.singleSize
		sz := gtx.Constraints.Min
		rr := float32(gtx.Px(t.singleCornerRadius))

		r := f32.Rect(0, 0, float32(sz.X), float32(sz.Y))
		t.toasts[index].singleToastLayoutShadow(gtx, r, rr)
		clip.UniformRRect(r, rr).Add(gtx.Ops)

		paint.Fill(gtx.Ops, t.toasts[index].bodyBackground)

		return l.Flex{Axis: l.Vertical, Alignment: l.Middle}.Layout(gtx,
			l.Flexed(1,
				func(gtx l.Context) l.Dimensions {
					return t.theme.Inset(0.25,
						t.theme.VFlex().
							Rigid(
								t.theme.Inset(0.1,
									t.theme.Fill(t.toasts[index].level,
										t.theme.Flex().
											Rigid(
												func(gtx l.Context) l.Dimensions {
													return t.theme.Icon().Color("DocText").Scale(1).Src(t.toasts[index].icon).Fn(gtx)
												},
											).
											Flexed(1,
												t.theme.H6(t.toasts[index].title).Color("PanelBg").Fn,
											).Fn,
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

	a := float32(t.bodyBackground.A) / 0xff
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
