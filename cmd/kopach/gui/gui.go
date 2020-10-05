package gui

import (
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/p9c/pod/pkg/gui/f"
	"github.com/p9c/pod/pkg/gui/fonts/p9fonts"
	"github.com/p9c/pod/pkg/gui/p9"
)

var (
	button = new(widget.Clickable)
)

func Run() {
	go func() {
		th := p9.NewTheme(p9fonts.Collection())
		f := f.Window().Size(640, 480)
		f.Run(func(ctx *layout.Context) {
			testLabels(th, *ctx)
		}, func() {
			os.Exit(0)
		})
	}()
	app.Main()
}

func testLabels(th *p9.Theme, gtx layout.Context) {
	f.Flex().Flexed(1,
		th.Inset(10,
			f.Flex().Vertical().Rigid(
				th.H1("this is a H1").Fn,
			).Rigid(
				th.H2("this is a H2").Fn,
			).Rigid(
				th.H3("this is a H3").Fn,
			).Rigid(
				th.H4("this is a H4").Fn,
			).Rigid(
				th.H5("this is a H5").Fn,
			).Rigid(
				th.H6("this is a H6").Fn,
			).Rigid(
				th.Body1("this is a Body1").Fn,
			).Rigid(
				th.Inset(10,
					th.Body2("this is a Body2").Fn,
				).Fn,
			).Rigid(
				th.Inset(10,
					th.Caption("this is a Caption").Fn,
				).Fn,
			).Rigid(
				th.Button(button, "plan9", "Click me!").Fn,
			).Fn,
		).Fn,
	).Fn(gtx)
}
