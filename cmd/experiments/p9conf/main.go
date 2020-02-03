// SPDX-License-Identifier: Unlicense OR MIT

package main

// A simple Gio program. See https://gioui.org for more information.

import (
	"image/color"
	"log"
	
	app "github.com/p9c/pod/pkg/gui/app"
	system "github.com/p9c/pod/pkg/gui/io/system"
	layout "github.com/p9c/pod/pkg/gui/layout"
	text "github.com/p9c/pod/pkg/gui/text"
	material "github.com/p9c/pod/pkg/gui/widget/material"
	
	gofont "github.com/p9c/pod/pkg/gui/font/gofont"
)

func main() {
	go func() {
		w := app.NewWindow()
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
	}()
	app.Main()
}

func loop(w *app.Window) error {
	gofont.Register()
	th := material.NewTheme()
	gtx := layout.NewContext(w.Queue())
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx.Reset(e.Config, e.Size)
			l := th.H1("Hello, Gio")
			maroon := color.RGBA{127, 0, 0, 255}
			l.Color = maroon
			l.Alignment = text.Middle
			l.Layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
}
