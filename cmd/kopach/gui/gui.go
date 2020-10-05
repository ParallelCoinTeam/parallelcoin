package gui

import (
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/p9c/pod/pkg/gui/fonts/p9fonts"
	"github.com/p9c/pod/pkg/gui/fui"
	"github.com/p9c/pod/pkg/gui/plan9"
)

var (
	button = new(widget.Clickable)
)

func Run() {
	go func() {
		th := plan9.NewTheme(p9fonts.Collection())
		f := fui.Window().Size(640, 480)
		f.Run(func(ctx *layout.Context) {
			testLabels(th, *ctx)
		}, func() {
			os.Exit(0)
		})
	}()
	app.Main()
}

func testLabels(th *plan9.Theme, gtx layout.Context) {
	x := fui.Flex().Vertical().
		Rigid(plan9.H1(th, "this is a H1").Layout).
		Rigid(plan9.H2(th, "this is a H2").Layout).
		Rigid(plan9.H3(th, "this is a H3").Layout).
		Rigid(plan9.H4(th, "this is a H4").Layout).
		Rigid(plan9.H5(th, "this is a H5").Layout).
		Rigid(plan9.H6(th, "this is a H6").Layout).
		Rigid(plan9.Body1(th, "this is a Body1").Layout).
		Rigid(fui.Inset(10, plan9.Body2(th, "this is a Body2").Layout).Layout).
		Rigid(fui.Inset(10, plan9.Caption(th, "this is a Caption").Layout).Layout).
		Rigid(plan9.Button(th, button, "plan9", "Click me!").Layout).
		Layout
	fui.Flex().Flexed(1, fui.Inset(10, x).Layout).Layout(gtx)
}
