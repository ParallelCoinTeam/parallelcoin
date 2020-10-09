package gui

import (
	"fmt"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"

	mico "golang.org/x/exp/shiny/materialdesign/icons"

	"github.com/p9c/pod/pkg/gui/ico/svg"

	"github.com/p9c/pod/pkg/gui/fonts/p9fonts"
	w "github.com/p9c/pod/pkg/gui/widget"

	"github.com/p9c/pod/pkg/gui/f"
	"github.com/p9c/pod/pkg/gui/p9"
)

var (
	button0     = w.NewClickable()
	button1     = w.NewClickable()
	button2     = w.NewClickable()
	bool1       bool
	boolbutton1 = w.NewBool(&bool1)
	iconbutton  = w.NewClickable()
	iconbutton1 = w.NewClickable()
	quit        = make(chan struct{})
	th          *p9.Theme
	progressbar *p9.ProgressBar
	progress    int
	slider      = &widget.Float{}
	radio       = new(widget.Enum)
)

func Run(quit chan struct{}) {
	th = p9.NewTheme(p9fonts.Collection(), quit)
	progressbar = th.ProgressBar().SetProgress(0)
	go func() {
		fw := f.Window().Size(640, 480)
		fw.Run(func(ctx *layout.Context) {
			testLabels(th, *ctx)
		}, func() {
			close(quit)
			os.Exit(0)
		})
	}()
	app.Main()
}

func testLabels(th *p9.Theme, gtx layout.Context) {
	th.Flex().Flexed(1,
		th.Flex().Rigid(
			th.Flex().Flexed(0.5,
				th.Fill("PanelBg").Widget(
					th.Inset(0.25).Widget(
						blocks(th),
					).Fn,
				).Fn,
			).Flexed(0.5,
				th.Fill("DocBg").Widget(
					th.Inset(0.25).Widget(
						buttons(th),
					).Fn,
				).Fn,
			).Fn,
		).Fn,
	).Fn(gtx)
	progress++
	if progress == 100 {
		progress = 0
	}
}

func blocks(th *p9.Theme) layout.Widget {
	return th.Flex().Vertical().Rigid(
		th.Inset(0.25).Widget(
			th.Flex().Rigid(
				th.H1("this is a H1").
					Color("PanelText").
					Fn,
			).Fn,
		).Fn,
	).Rigid(
		th.Inset(0.25).Widget(
			th.H2("this is a H2").
				Font("bariol regular").
				Color("PanelText").Fn,
		).Fn,
	).Rigid(
		th.Inset(0.25).Widget(
			th.H3("this is a H3").
				Alignment(text.End).
				Color("PanelText").Fn,
		).Fn,
	).Rigid(
		th.Fill("DocBg").Widget(
			th.Inset(0.25).Widget(
				th.H4("this is a H4").
					Alignment(text.Middle).
					Color("DocText").Fn,
			).Fn,
		).Fn,
	).Rigid(
		th.Fill("PanelBg").Widget(
			th.Inset(0.25).Widget(
				th.H5("this is a H5").
					Color("PanelText").
					Fn,
			).Fn,
		).Fn,
	).Rigid(
		th.Inset(0.25).Widget(
			th.H6("this is a H6").
				Color("PanelText").Fn,
		).Fn,
	).Rigid(
		th.Inset(0.25).Widget(
			th.Body1("this is a Body1").Color("PanelText").Fn,
		).Fn,
	).Rigid(
		th.Inset(0.25).Widget(
			th.Body2("this is a Body2").Color("PanelText").Fn,
		).Fn,
	).Rigid(
		th.Inset(0.25).Widget(
			th.Caption("this is a Caption").Color("PanelText").Fn,
		).Fn,
	).Fn
}

func buttons(th *p9.Theme) layout.Widget {
	return th.Flex().Vertical().Rigid(
		th.Inset(0.25).Widget(
			th.Flex().Rigid(
				th.Button(
					button0.SetClick(func() {
						Info("clicked customised button")
					})).
					CornerRadius(1).
					Background("Secondary").
					Color("Dark").
					Font("bariol bold").
					TextScale(2).
					Text("customised button").
					Inset(1).
					Fn,
			).Fn,
		).Fn,
	).Rigid(
		th.Flex().Rigid(
			th.Inset(0.25).Widget(
				th.Button(
					button2.SetClick(func() {
						Info("clicked default style button")
					})).
					Text("default style").
					Fn,
			).Fn,
		).Rigid(
			th.Inset(0.25).Widget(
				th.Flex().Rigid(
					th.Indefinite().Color("Primary").Fn,
				).Fn,
			).Fn,
		).Rigid(
			th.Inset(0.25).Widget(
				th.Icon().Scale(2).Color("DocText").Src(icons.ParallelCoinRound).Fn,
			).Fn,
		).Rigid(
			th.Inset(0.25).Widget(
				th.IconButton(iconbutton.SetClick(
					func() {
						Debug("clicked parallelcoin button")
					})).
					Scale(50).
					Icon(
						th.Icon().
							Color("Light").
							Src(icons.ParallelCoin)).
					Fn,
			).Fn,
		).Rigid(
			th.Inset(0.25).Widget(
				th.IconButton(iconbutton1.SetClick(
					func() {
						Debug("clicked android button")
					})).
					Scale(50).
					Background("Secondary").
					Icon(
						th.Icon().
							Color("Light").
							Src(mico.ActionAndroid)).
					Fn,
			).Fn,
		).Rigid(
			th.Inset(0.25).Widget(
				th.CheckBox(boolbutton1.SetHook(func(b bool) {
					Debug("change state to", b)
				})).
					IconColor("Primary").
					TextColor("DocText").
					// IconScale(0.1).
					Text("checkbox").
					Fn,
			).Fn,
		).Fn,
	).Rigid(
		th.Flex().Rigid(
			th.Inset(0.25).Widget(
				th.Flex().Rigid(
					th.ProgressBar().Color("Primary").SetProgress(progress).Fn,
				).Fn,
			).Fn,
		).Fn,
	).Rigid(
		// th.Flex().Rigid(
			th.Inset(0.25).Widget(
				th.Flex().Rigid(
					th.ProgressBar().Color("Primary").SetProgress(100 - progress).Fn,
				).Fn,
			).Fn,
		// ).Fn,
	).Rigid(
		th.Inset(0.25).Widget(
			th.Flex().
				Flexed(1,
					th.Slider(slider, 0, 1).Fn,
				).
				Rigid(
					th.Body1(fmt.Sprintf("%3v", int(slider.Value*100))).
						Font("go regular").Color("DocText").
						Fn,
				).
				Fn,
		).Fn,
	).Rigid(
		th.Flex().Rigid(
			th.RadioButton(radio, "r1", "first").Fn,
		).Rigid(
			th.RadioButton(radio, "r2", "second").Fn,
		).Rigid(
			th.RadioButton(radio, "r3", "third").Fn,
		).Fn,
	).Fn
}
