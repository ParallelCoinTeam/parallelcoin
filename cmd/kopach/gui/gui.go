package gui

import (
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/text"
	"golang.org/x/exp/shiny/materialdesign/icons"

	w "github.com/p9c/pod/pkg/gui/widget"

	"github.com/p9c/pod/pkg/gui/f"
	"github.com/p9c/pod/pkg/gui/fonts/p9fonts"
	"github.com/p9c/pod/pkg/gui/p9"
)

var (
	button0     = w.NewClickable()
	button1     = w.NewClickable()
	button2     = w.NewClickable()
	bool1       bool
	boolbutton1 = w.NewBool(&bool1)
	iconbutton  = w.NewClickable()
)

func Run(quit chan struct{}) {
	go func() {
		th := p9.NewTheme(p9fonts.Collection(), quit)
		fw := f.Window().Size(640, 480)
		fw.Run(func(ctx *layout.Context) {
			testLabels(th, *ctx)
		}, func() {
			os.Exit(0)
		})
	}()
	app.Main()
}

func testLabels(th *p9.Theme, gtx layout.Context) {
	th.Flex().Flexed(1,
		th.Flex().Vertical().Flexed(0.5,
			th.Flex().Rigid(
				th.Flex().Flexed(0.5,
					th.Fill("PanelBg").Widget(
						th.Inset(0.5).Widget(
							blocks(th),
						).Fn,
					).Fn,
				).Flexed(0.5,
					th.Fill("DocBg").Widget(
						buttons(th),
					).Fn,
				).Fn,
			).Fn,
		).Fn,
	).Fn(gtx)
}

func blocks(th *p9.Theme) layout.Widget {
	return th.Flex().Vertical().Rigid(
		th.Inset(0.5).Widget(
			th.Flex().Rigid(
				th.H1("this is a H1").
					Color("PanelText").
					Fn,
			).Fn,
		).Fn,
	).Rigid(
		th.Inset(0.5).Widget(
			th.H2("this is a H2").
				Font("bariol regular").
				Color("PanelText").Fn,
		).Fn,
	).Rigid(
		th.Inset(0.5).Widget(
			th.H3("this is a H3").
				Alignment(text.End).
				Color("PanelText").Fn,
		).Fn,
	).Rigid(
		th.Fill("DocBg").Widget(
			th.Inset(0.5).Widget(
				th.H4("this is a H4").
					Alignment(text.Middle).
					Color("DocText").Fn,
			).Fn,
		).Fn,
	).Rigid(
		th.Fill("PanelBg").Widget(
			th.Inset(0.5).Widget(
				th.H5("this is a H5").
					Color("PanelText").
					Fn,
			).Fn,
		).Fn,
	).Rigid(
		th.Inset(0.5).Widget(
			th.H6("this is a H6").
				Color("PanelText").Fn,
		).Fn,
	).Rigid(
		th.Inset(0.5).Widget(
			th.Body1("this is a Body1").Color("PanelText").Fn,
		).Fn,
	).Rigid(
		th.Inset(0.5).Widget(
			th.Body2("this is a Body2").Color("PanelText").Fn,
		).Fn,
	).Rigid(
		th.Inset(0.5).Widget(
			th.Caption("this is a Caption").Color("PanelText").Fn,
		).Fn,
	).Fn
}

func buttons(th *p9.Theme) layout.Widget {
	return th.Flex().Vertical().Rigid(
		th.Inset(0.5).Widget(
			th.Flex().Rigid(
				th.Button(
					button0.SetClick(func() {
						Info("clicked first button")
					})).
					CornerRadius(2).
					Background("Secondary").
					Color("Dark").
					Font("bariol bold").
					TextScale(2).
					Text("customised button").
					Inset(8).
					Fn,
			).Fn,
		).Fn,
	).Rigid(
		th.Flex().Vertical().Rigid(
			th.Flex().Rigid(
				th.Inset(0.5).Widget(
					th.Button(
						button2.SetClick(func() {
							Info("clicked third button")
						})).
						Text("default style button").
						Fn,
				).Fn,
			).Fn,
		).Rigid(
			th.Flex().Rigid(
				th.Inset(0.5).Widget(
					th.Button(
						button1.SetClick(func() {
							Info("clicked second button")
						})).
						TextScale(0.5).
						Text("button").
						Inset(4).
						Fn,
				).Fn,
			).Fn,
		).Rigid(
			th.CheckBox(boolbutton1.SetHook(func(b bool) {
				Debug("change state to", b)
			})).
				IconColor("Primary").
				TextColor("DocText").
				// Scale(0.1).
				Label("checkbox").
				Fn,
		).Rigid(
			th.Flex().Rigid(
				th.IconButton(iconbutton).
					// Inset(th.Inset(0.5)).
					Icon(th.Icon().Color("Light").Src(icons.ActionAndroid)).
					Fn,
			).Fn,
		).Fn,
	).Fn
}
