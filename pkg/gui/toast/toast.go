package toast

import (
	"gioui.org/f32"
	l "gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"github.com/gioapp/gel/helper"
	"github.com/p9c/pod/pkg/gui/p9"
	icons2 "golang.org/x/exp/shiny/materialdesign/icons"
	"image"
	"image/color"
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
		duration:           100,
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

func (t *Toasts) DrawToasts() func(gtx l.Context) {
	return func(gtx l.Context) {
		defer op.Push(gtx.Ops).Pop()
		op.Offset(f32.Pt(float32(gtx.Constraints.Max.X)-250, 0)).Add(gtx.Ops)
		gtx.Constraints.Min = image.Pt(250, gtx.Constraints.Min.Y)
		gtx.Constraints.Max.X = 250
		//paint.Fill(gtx.Ops,  helper.HexARGB("ff559988"))
		t.theme.Inset(0,
			t.layout.Vertical().ScrollToEnd().Length(len(t.toasts)).ListElement(t.singleToast).Fn).Fn(gtx)
	}
}
func (t *Toasts) singleToast(gtx l.Context, index int) l.Dimensions {
	if t.toasts[index].ticker < float32(t.duration) {
		t.toasts[index].ticker += 1
		gtx.Constraints.Min = t.singleSize
		gtx.Constraints.Max = t.singleSize
		//gtx.Constraints.Max.X = t.singleSize.X
		//sz := gtx.Constraints.Min
		//rr := float32(gtx.Px(t.singleCornerRadius))

		//r := f32.Rect(0, 0, float32(sz.X), float32(sz.Y))
		//clip.UniformRRect(r, rr).Add(gtx.Ops)

		paint.Fill(gtx.Ops, t.toasts[index].bodyBackground)

		return t.theme.Inset(1,
			t.theme.Flex().Flexed(1,
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
				}).Fn).Fn(gtx)
	} else {
		t.toasts = remove(t.toasts, index)
		return p9.EmptySpace(0, 0)(gtx)
	}
}

func remove(slice []toast, s int) []toast {
	return append(slice[:s], slice[s+1:]...)
}
