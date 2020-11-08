package dialog

import (
	"gioui.org/io/pointer"
	l "gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"github.com/gioapp/gel/helper"
	"github.com/p9c/pod/pkg/gui/p9"
	"image"
	"image/color"
)

type Dialog struct {
	theme              *p9.Theme
	duration           int
	singleCornerRadius unit.Value
	singleElevation    unit.Value

	content          *content
	headerBackground color.RGBA
	bodyBackground   color.RGBA
	icon             *[]byte
	ticker           float32
	hideTitle        bool
	close            *p9.Clickable
	cornerRadius     unit.Value
	elevation        unit.Value
}

type content struct {
	title, level string
	content      interface{}
}

func New(th *p9.Theme) *Dialog {
	return &Dialog{
		theme:          th,
		duration:       100,
		close:          th.Clickable(),
		bodyBackground: helper.HexARGB("ee000000"),
		//singleSize:         image.Pt(300, 80),
		singleCornerRadius: unit.Dp(5),
		singleElevation:    unit.Dp(5),
	}
}
func (d *Dialog) ShowDialog(title, level string, contentInterface interface{}) func() {
	return func() {
		c := &content{
			title:   title,
			content: contentInterface,
			level:   level,
		}
		d.content = c
	}
}

func (d *Dialog) DrawDialog() func(gtx l.Context) {
	//switch d.content.level {
	//case "Warning":
	//	//ic = &icons2.AlertWarning
	//case "Success":
	//	//ic = &icons2.NavigationCheck
	//case "Danger":
	//	//ic = &icons2.AlertError
	//case "Info":
	//	//ic = &icons2.ActionInfo
	//}

	return func(gtx l.Context) {
		if d.content != nil {
			var content func(gtx l.Context) l.Dimensions
			switch c := d.content.content.(type) {
			case string:
				content = d.theme.Body1(c).Color("PanelText").Fn
			case func(gtx l.Context) l.Dimensions:
				content = c
			}

			defer op.Push(gtx.Ops).Pop()
			gtx.Constraints.Min = gtx.Constraints.Max
			d.theme.Stack().Alignment(l.Center).Expanded(
				func(gtx l.Context) l.Dimensions {
					paint.Fill(gtx.Ops, d.bodyBackground)
					pointer.Rect(
						image.Rectangle{Max: gtx.Constraints.Max}).Add(gtx.Ops)
					return l.Dimensions{Size: gtx.Constraints.Max}
				}).Stacked(
				d.theme.Fill("DocBg",
					d.theme.Inset(0.25,
						d.theme.Fill("PanelBg",
							d.theme.Inset(1,
								d.theme.VFlex().
									Rigid(
										d.theme.Body1(d.content.title).Color("PanelText").Fn,
									).
									Rigid(content).
									Rigid(
										d.theme.Button(d.close).Text("CLOSE").Color("Warning").SetClick(d.Close).Fn,
									).Fn).Fn).Fn).Fn).Fn).Fn(gtx)
		}
	}
}

func (d *Dialog) Close() {
	d.content = nil
}
