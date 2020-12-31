package main

import (
	"log"
	"os"
	
	"gioui.org/app"
	"gioui.org/io/system"
	l "gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	
	"github.com/p9c/pod/pkg/gui/fonts/p9fonts"
	"github.com/p9c/pod/pkg/gui/p9"
	"github.com/p9c/pod/pkg/gui/shadow"
)

var (
	th = p9.NewTheme(p9fonts.Collection(), nil)
)

func main() {
	go func() {
		w := app.NewWindow(app.Size(unit.Px(150*6+50), unit.Px(150*6-50)))
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window) error {
	var ops op.Ops
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := l.NewContext(&ops, e)
			paint.Fill(gtx.Ops, p9.HexNRGB("e5e5e5FF"))
			op.InvalidateOp{}.Add(gtx.Ops)
			
			th.Inset(
				5,
				th.VFlex().
					Flexed(
						1,
						th.VFlex().AlignMiddle().
							Rigid(
								th.Inset(
									1,
									func(gtx l.Context) l.Dimensions {
										return shadow.Shadow(
											gtx,
											unit.Dp(5),
											unit.Dp(3),
											p9.HexNRGB("ee000000"),
											th.Fill("DocBg", th.Inset(3,
												th.Body1("Shadow test 3").Color(
													"PanelText").Fn).Fn, l.Center).Fn,
										)
									},
								).Fn,
							).
							Rigid(
								th.Inset(
									1,
									func(gtx l.Context) l.Dimensions {
										return shadow.Shadow(
											gtx,
											unit.Dp(5),
											unit.Dp(5),
											p9.HexNRGB("ee000000"),
											th.Fill("DocBg", th.Inset(3, th.Body1("Shadow test 5").Color("PanelText").Fn).Fn, l.Center).Fn,
										)
									},
								).Fn,
							).
							Rigid(
								th.Inset(
									1,
									func(gtx l.Context) l.Dimensions {
										return shadow.Shadow(
											gtx,
											unit.Dp(5),
											unit.Dp(8),
											p9.HexNRGB("ee000000"),
											th.Fill("DocBg", th.Inset(3, th.Body1("Shadow test 8").Color("PanelText").Fn).Fn, l.Center).Fn,
										)
									},
								).Fn,
							).
							Rigid(
								th.Inset(
									1,
									func(gtx l.Context) l.Dimensions {
										return shadow.Shadow(
											gtx,
											unit.Dp(5),
											unit.Dp(12),
											p9.HexNRGB("ee000000"),
											th.Fill("DocBg", th.Inset(3, th.Body1("Shadow test 12").Color("PanelText").Fn).Fn, l.Center).Fn,
										)
									},
								).Fn,
							).Fn,
					).Fn,
			).Fn(gtx)
			e.Frame(gtx.Ops)
			w.Invalidate()
		}
	}
}
