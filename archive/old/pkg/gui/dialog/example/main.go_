package main

import (
	"log"
	"os"
	
	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	
	"github.com/p9c/pod/pkg/gui"
	"github.com/p9c/pod/pkg/gui/dialog"
	"github.com/p9c/pod/pkg/gui/fonts/p9fonts"
)

var (
	th         = gui.NewTheme(p9fonts.Collection(), nil)
	btnDanger  = th.Clickable()
	btnWarning = th.Clickable()
	btnSuccess = th.Clickable()
)

func main() {
	go func() {
		w := app.NewWindow(app.Size(unit.Px(150*6+50), unit.Px(150*6-50)))
		if e := loop(w); E.Chk(e) {
			log.F.Ln(e)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window) (e error) {
	var ops op.Ops
	d := dialog.New(th)
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			paint.Fill(gtx.Ops, gui.HexNRGB("e5e5e5FF"))
			op.InvalidateOp{}.Add(gtx.Ops)
			
			th.Inset(
				0.25,
				th.VFlex().
					Rigid(
						// th.Button(btnDanger).Text("Danger").Color("Danger").Fn,
						// ).
						// Rigid(
						//	th.Button(btnWarning).Text("Warning").Color("Warning").Fn,
						// ).
						// Rigid(
						th.Button(btnSuccess).Text("Success").Color("Success").SetClick(
							d.ShowDialog(
								"Success",
								"Success content",
								"Success",
							),
						).Fn,
					).Fn,
			).Fn(gtx)
			
			// for btnDanger.Clicked() {
			//	d.DrawDialog("Danger", "Danger content", "Danger")
			// }
			
			// for btnWarning.Clicked() {
			//	d.DrawDialog("Warning", "Warning content", "Warning")
			// }
			
			d.DrawDialog()(gtx)
			e.Frame(gtx.Ops)
			w.Invalidate()
		}
	}
}
