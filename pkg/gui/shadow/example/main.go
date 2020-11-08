package main

import (
	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"github.com/gioapp/gel/helper"
	"github.com/p9c/pod/pkg/gui/fonts/p9fonts"
	"github.com/p9c/pod/pkg/gui/p9"
	"github.com/p9c/pod/pkg/gui/shadow"
	"log"
	"os"
)

var (
	th = p9.NewTheme(p9fonts.Collection(), nil)
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
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			paint.Fill(gtx.Ops, helper.HexARGB("e5e5e5FF"))
			op.InvalidateOp{}.Add(gtx.Ops)

			th.Inset(0.25,
				th.VFlex().
					Rigid(
						shadow.Shadow(th).Drop(th.H6("Success").Color("Success").Fn(gtx)),
					).Fn).Fn(gtx)

			e.Frame(gtx.Ops)
			w.Invalidate()
		}
	}
}
