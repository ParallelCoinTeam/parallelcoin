package duoui

import (
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/op"
	"github.com/p9c/pod/pkg/gui/text"
	"github.com/p9c/pod/pkg/gui/unit"
	"github.com/p9c/pod/pkg/gui/widget"
	"github.com/p9c/pod/pkg/gui/widget/parallel"
	"image/color"
	"time"
)

var (
	buttonToastOK = new(widget.Button)
	listToasts    = &layout.List{
		Axis:        layout.Vertical,
		ScrollToEnd: false,
		Alignment:   0,
		Position:    layout.Position{},
	}
)

func (duo *DuoUI) DuoUItoastSys(rc *rcd.RcVar) {
	layout.Align(layout.NE).Layout(duo.m.DuoUIcontext, func() {
		listToasts.Layout(duo.m.DuoUIcontext, len(rc.Toasts), func(i int) {
			layout.UniformInset(unit.Dp(16)).Layout(duo.m.DuoUIcontext, rc.Toasts[i])
		})
	})
}

func toastButton(text, txtColor, bgColor, iconColor string, duo *DuoUI, rc *rcd.RcVar, button *widget.Button, icon *parallel.DuoUIicon) func() {
	var b parallel.DuoUIbutton
	return func() {
		layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8), Left: unit.Dp(8), Right: unit.Dp(8)}.Layout(duo.m.DuoUIcontext, func() {
			b = duo.m.DuoUItheme.DuoUIbutton(text, txtColor, bgColor, iconColor, 24, 120, 60, 0, 0, icon)
			for button.Clicked(duo.m.DuoUIcontext) {
				//rc.ShowToast = false
			}
			b.Layout(duo.m.DuoUIcontext, button)
		})
	}
}

func toastAdd(duo *DuoUI, rc *rcd.RcVar) {
	rc.Toasts = append(rc.Toasts, func() {
		//iconOK, _ := parallel.NewDuoUIicon(icons.NavigationCheck)
		helpers.DuoUIdrawRectangle(duo.m.DuoUIcontext, 418, 160, helpers.HexARGB("aa000000"), [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
		helpers.DuoUIdrawRectangle(duo.m.DuoUIcontext, 408, 150, duo.m.DuoUItheme.Color.Primary, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
		layout.Flex{
			Axis:      layout.Vertical,
			Alignment: layout.Middle,
		}.Layout(duo.m.DuoUIcontext,
			layout.Rigid(func() {
				layout.Flex{
					Axis:      layout.Horizontal,
					Alignment: layout.Middle,
				}.Layout(duo.m.DuoUIcontext,
					layout.Rigid(func() {
						layout.Align(layout.Center).Layout(duo.m.DuoUIcontext, func() {
							layout.Inset{Top: unit.Dp(24), Bottom: unit.Dp(8), Left: unit.Dp(0), Right: unit.Dp(4)}.Layout(duo.m.DuoUIcontext, func() {
								cur := duo.m.DuoUItheme.H4("TOAST MESSAGE!!!")
								cur.Color = color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}
								cur.Alignment = text.Start
								cur.Layout(duo.m.DuoUIcontext)
							})
						})
					}),
				)
			}),
			layout.Rigid(func() {
				layout.Flex{
					Axis:      layout.Horizontal,
					Alignment: layout.Middle,
				}.Layout(duo.m.DuoUIcontext,
					//layout.Rigid(toastButton("OK", "ffcfcfcf", "ff308030", "ffcfcfcf", duo, rc, buttonToastOK, iconOK)),
				)
			}),
		)
	})
	go func(duo *DuoUI, ops *op.Ops) {
		time.Sleep(3 * time.Second)
		rc.Toasts[len(rc.Toasts)-1] = nil // or the zero value of T
		rc.Toasts = rc.Toasts[:len(rc.Toasts)-1]
		op.InvalidateOp{}.Add(ops)
	}(duo, duo.m.DuoUIcontext.Ops)
}
