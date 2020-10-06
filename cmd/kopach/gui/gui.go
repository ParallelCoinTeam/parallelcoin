package gui

import (
	"os"

	"gioui.org/app"
	"gioui.org/layout"

	w "github.com/p9c/pod/pkg/gui/widget"

	"github.com/p9c/pod/pkg/gui/f"
	"github.com/p9c/pod/pkg/gui/fonts/p9fonts"
	"github.com/p9c/pod/pkg/gui/p9"
)

var (
	button = w.NewClickable()
)

func Run(quit chan struct{}) {
	go func() {
		th := p9.NewTheme(p9fonts.Collection(), quit)
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
	th.Flex().Flexed(1,
		th.Inset(10,
			th.Flex().Vertical().Rigid(
				th.Inset(4,
					th.H1("this is a H1").Fn,
				).Fn,
			).Rigid(
				th.Inset(4,
					th.H2("this is a H2").Fn,
				).Fn,
			).Rigid(
				th.Inset(4,
					th.H3("this is a H3").Fn,
				).Fn,
			).Rigid(
				th.Inset(4,
					th.H4("this is a H4").Fn,
				).Fn,
			).Rigid(
				th.Inset(4,
					th.H5("this is a H5").Fn,
				).Fn,
			).Rigid(
				th.Inset(4,
					th.H6("this is a H6").Fn,
				).Fn,
			).Rigid(
				th.Inset(4,
					th.Body1("this is a Body1").Fn,
				).Fn,
			).Rigid(
				th.Inset(4,
					th.Body2("this is a Body2").Fn,
				).Fn,
			).Rigid(
				th.Inset(4,
					th.Caption("this is a Caption").Fn,
				).Fn,
			).Rigid(
				th.Button(button, "plan9", "Click me!", w.ClickEvents{
					Click: func() {
						Info("click event")
					}, Cancel: func() {
						Info("cancel event")
					}, Press: func() {
						Info("press event")
					},
				}).Fn,
			).Fn,
		).Fn,
	).Fn(gtx)
}
