package duoui

import (
	svg "github.com/p9c/pod/pkg/gui/ico/svg"
	"image"

	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/p9c/pod/cmd/gui/component"
	"github.com/p9c/pod/pkg/gui/gelook"
)

// Main wallet screen
func (ui *DuoUI) DuoUIsplashScreen() {
	ctx := ui.ly.Context
	th := ui.ly.Theme
	th.DuoUIcontainer(0, th.Colors["Dark"]).
		Layout(ctx, layout.Center, func() {
			logo, _ := gelook.NewDuoUIicon(svg.ParallelCoin)
			layout.Flex{Axis: layout.Vertical}.Layout(ctx,
				layout.Rigid(func() {
					layout.Flex{Axis: layout.Horizontal}.Layout(ctx,
						layout.Rigid(func() {
							layout.UniformInset(unit.Dp(8)).
								Layout(ctx, func() {
									size := ctx.Px(unit.Dp(256)) - 2*ctx.
										Px(unit.Dp(8))
									if logo != nil {
										logo.Color = gelook.HexARGB(th.
											Colors["Light"])
										logo.Layout(ctx, unit.Px(float32(size)))
									}
									ctx.Dimensions = layout.Dimensions{
										Size: image.Point{X: size, Y: size},
									}
								})
						}),
						layout.Flexed(1, func() {
							layout.UniformInset(unit.Dp(60)).Layout(ctx, func() {
								txt := th.H1("PLAN NINE FROM FAR, " +
									"FAR AWAY SPACE")
								txt.Font.Typeface = th.Fonts["Secondary"]
								txt.Color = th.Colors["Light"]
								txt.Layout(ctx)
							})
						}),
					)
				}),
				layout.Flexed(1, component.DuoUIlogger(ui.rc, ctx, th)),
			)
		})
}

// Main wallet screen
func (ui *DuoUI) DuoUImainScreen() {
	ctx := ui.ly.Context
	th := ui.ly.Theme
	th.DuoUIcontainer(0, th.Colors["Dark"]).Layout(ctx,
		layout.Center, func() {
			layout.Flex{Axis: layout.Vertical}.Layout(ctx,
				layout.Rigid(ui.DuoUIheader()),
				layout.Flexed(1, ui.DuoUIbody()),
				layout.Rigid(ui.DuoUIfooter()),
			)
		})
}
