// SPDX-License-Identifier: Unlicense OR MIT

package win

import (
	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"

	"github.com/p9c/pod/pkg/gui/wallet/theme"
)

type Window struct {
	W  *app.Window
	th *theme.Theme
	L  func(gtx layout.Context) layout.Dimensions
}

type Windows struct {
	WindowCount int32
	W           map[string]*Window
}

func (w *Window) Loop(events <-chan event.Event) error {
	var ops op.Ops
	for {
		e := <-events
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			// for w.more.Clicked() {
			//	newWindow()
			// }
			// for w.close.Clicked() {
			//	w.Close()
			// }
			gtx := layout.NewContext(&ops, e)

			w.L(gtx)

			e.Frame(gtx.Ops)
		}
	}
}
