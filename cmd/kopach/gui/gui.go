package gui

import (
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"

	"github.com/p9c/pod/pkg/gui/fonts/p9fonts"
	"github.com/p9c/pod/pkg/util/interrupt"
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
	th := material.NewTheme(p9fonts.Collection())
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
			h := material.Body1(th, "Kopach Miner")
			maroon := color.RGBA{127, 0, 0, 255}
			h.Color = maroon
			h.Alignment = text.Middle
			h.Font = text.Font{Typeface: "bariol", Weight: text.Bold}
			h.TextSize = unit.Dp(20)
			h.Layout(gtx)
			e.Frame(gtx.Ops)
		}
		select {
		case <-quit:
			return nil
		default:
		}
	}
}
