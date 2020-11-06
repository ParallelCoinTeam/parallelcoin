package main

import (
	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/p9c/pod/pkg/gui/toast/f32color"
	"image"
	"log"
	"math"
	"os"
)

type Toast struct {
	ticker   float64
	duration int
	offset   *offset
}
type offset struct {
	x, y float32
}

var (
	btn = new(widget.Clickable)
	th  = material.NewTheme(gofont.Collection())
)

func main() {
	go func() {
		w := app.NewWindow(app.Size(unit.Px(150*6+50), unit.Px(150*6-50)))
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window) error {
	var ops op.Ops

	t := &Toast{
		ticker:   0.0,
		duration: 111,
		offset:   &offset{-50, -50},
	}

	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			paint.Fill(gtx.Ops, f32color.RGBAHex(0xe5e5e5FF))
			op.InvalidateOp{}.Add(gtx.Ops)

			material.Button(th, btn, "toast").Layout(gtx)
			for btn.Clicked() {

			}
			t.drawSurface(gtx)
			t.ticker += 0.1
			e.Frame(gtx.Ops)
			w.Invalidate()
		}
	}
}

func (t *Toast) drawSurface(gtx layout.Context) {
	defer op.Push(gtx.Ops).Pop()
	op.Offset(f32.Pt(0, (float32(math.Cos(t.ticker))+1)*50)).Add(gtx.Ops)
	gtx.Constraints.Min = image.Pt(100, 100)
	gtx.Constraints.Max = image.Pt(100, 100)
	style := SurfaceLayoutStyle{
		Background:   f32color.RGBAHex(0xffffffff),
		CornerRadius: unit.Dp(5),
		Elevation:    unit.Dp(5),
	}
	style.Layout(gtx)
}
