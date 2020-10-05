package gui

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"

	"github.com/p9c/pod/pkg/gui/plan9"

	"github.com/p9c/pod/pkg/gui/fonts/p9fonts"
	"github.com/p9c/pod/pkg/util/interrupt"
)

type (
	D = layout.Dimensions
	C = layout.Context
	W = layout.Widget
)

func Run(quit chan struct{}) {
	quit = make(chan struct{})
	go func() {
		w := app.NewWindow(
			app.Size(unit.Dp(640), unit.Dp(480)),
			app.Title("kopach"),
		)
		if err := loop(w, quit); err != nil {
			log.Fatal(err)
		}
		Debug("exiting gui")
		os.Exit(0)
	}()
	go app.Main()
}

func loop(w *app.Window, quit chan struct{}) error {
	th := plan9.NewTheme(p9fonts.Collection())
	var ops op.Ops
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			// return e.Err
			interrupt.Request()
			close(quit)
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			testLabels(th, gtx)
			// h := plan9.H1(th, "Kopach Miner")
			// maroon := color.RGBA{127, 0, 0, 255}
			// h.Color = maroon
			// h.Alignment = text.Middle
			// h.Font = text.Font{Typeface: "bariol", Weight: text.Bold}
			// h.TextSize = unit.Dp(20)
			// h.Layout(gtx)
			e.Frame(gtx.Ops)
		}
		select {
		case <-quit:
			return nil
		default:
		}
	}
}

func testLabels(th *plan9.Theme, gtx C) D {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return plan9.H1(th, "this is a H1").Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return plan9.H2(th, "this is a H2").Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return plan9.H3(th, "this is a H3").Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return plan9.H4(th, "this is a H4").Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return plan9.H5(th, "this is a H5").Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return plan9.H6(th, "this is a H6").Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return plan9.Body1(th, "this is a Body1").Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return plan9.Body2(th, "this is a Body2").Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return plan9.Caption(th, "this is a Caption").Layout(gtx)
		}),
	)
}