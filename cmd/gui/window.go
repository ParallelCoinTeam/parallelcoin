// SPDX-License-Identifier: Unlicense OR MIT

package gui

// Multiple windows in Gio.

import (
	qu "github.com/p9c/pod/pkg/util/quit"
	"log"
	"os"
	"sync/atomic"
	
	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/io/system"
	l "gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	
	"github.com/p9c/pod/pkg/gui/p9"
)

type window struct {
	*app.Window
	quit  qu.C
	more  *p9.Clickable
	close *p9.Clickable
}

//
// func _main() {
// 	newWindow()
// 	app.Main()
// }

var windowCount int32

func newWindow() {
	atomic.AddInt32(&windowCount, +1)
	go func() {
		w := new(window)
		if w.quit == nil {
			w.quit = qu.T()
		}
		w.Window = app.NewWindow()
		if err := w.loop(w.Events()); err != nil {
			log.Fatal(err)
		}
		if c := atomic.AddInt32(&windowCount, -1); c == 0 {
			os.Exit(0)
		}
	}()
}

func (w *window) loop(events <-chan event.Event) error {
	// th := p9.NewTheme(gofont.Collection(), w.quit)
	var ops op.Ops
	for {
		e := <-events
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			for w.more.Clicked() {
				newWindow()
			}
			for w.close.Clicked() {
				w.Close()
			}
			gtx := l.NewContext(&ops, e)
			
			l.Center.Layout(
				gtx, func(gtx l.Context) l.Dimensions {
					return l.Dimensions{}
					// return layout.Flex{
					// 	Alignment: layout.Middle,
					// }.Layout(gtx,
					// 	RigidInset(material.Button(th, &w.more, "More!").Layout),
					// 	RigidInset(material.Button(th, &w.close, "Close").Layout),
					// )
				},
			)
			e.Frame(gtx.Ops)
		}
	}
}

func RigidInset(w l.Widget) l.FlexChild {
	return l.Rigid(
		func(gtx l.Context) l.Dimensions {
			return l.UniformInset(unit.Sp(5)).Layout(gtx, w)
		},
	)
}
