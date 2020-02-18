package duoui

import (
	"fmt"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/text"
	"github.com/p9c/pod/pkg/gui/unit"
)

var (
	latestTransList = &layout.List{
		Axis: layout.Vertical,
	}
)

func (ui *DuoUI) DuoUIlatestTransactions() func() {
	return func() {
		cs := ui.ly.Context.Constraints
		theme.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, cs.Height.Max, "ff424242", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
		layout.Flex{
			Axis: layout.Vertical,
		}.Layout(ui.ly.Context,
			layout.Rigid(func() {
				cs := ui.ly.Context.Constraints
				theme.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, 48, ui.ly.Theme.Color.Primary, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
				layout.UniformInset(unit.Dp(16)).Layout(ui.ly.Context, func() {
					latestx := ui.ly.Theme.H5("LATEST TRANSACTIONS")
					latestx.Color = theme.HexARGB(ui.ly.Theme.Color.Light)
					latestx.Alignment = text.Start
					latestx.Layout(ui.ly.Context)
				})
			}),
			layout.Flexed(1, func() {
				layout.UniformInset(unit.Dp(8)).Layout(ui.ly.Context, func() {
					layout.Flex{Axis: layout.Vertical}.Layout(ui.ly.Context,
						layout.Rigid(func() {
							cs := ui.ly.Context.Constraints
							latestTransList.Layout(ui.ly.Context, len(ui.rc.Status.Wallet.LastTxs.Txs), func(i int) {
								t := ui.rc.Status.Wallet.LastTxs.Txs[i]
								theme.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, cs.Height.Max, ui.ly.Theme.Color.Text, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
								layout.Inset{
									Top:    unit.Dp(8),
									Right:  unit.Dp(16),
									Bottom: unit.Dp(8),
									Left:   unit.Dp(16),
								}.Layout(ui.ly.Context, func() {
									layout.Flex{Axis: layout.Vertical}.Layout(ui.ly.Context,
										layout.Rigid(lTtxid(ui.ly.Context, ui.ly.Theme, t.TxID)),
										layout.Rigid(func() {
											layout.Flex{
												Spacing: layout.SpaceBetween,
											}.Layout(ui.ly.Context,
												layout.Rigid(func() {
													layout.Flex{
														Axis: layout.Vertical,
													}.Layout(ui.ly.Context,
														layout.Rigid(lTcategory(ui.ly.Context, ui.ly.Theme, t.Category)),
														layout.Rigid(lTtime(ui.ly.Context, ui.ly.Theme, t.Time)),
													)
												}),
												layout.Rigid(lTamount(ui.ly.Context, ui.ly.Theme, t.Amount)),
											)
										}),
										layout.Rigid(func() {
											cs := ui.ly.Context.Constraints
											theme.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, 1, ui.ly.Theme.Color.Hint, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
										}),
									)
								})
							})
						}))
				})
			}),
		)
	}
}

func lTtxid(gtx *layout.Context, th *theme.DuoUItheme, v string) func() {
	return func() {
		tim := th.Caption(v)
		tim.Font.Typeface = th.Font.Mono
		tim.Color = theme.HexARGB(th.Color.Light)
		tim.Layout(gtx)
	}
}

func lTcategory(gtx *layout.Context, th *theme.DuoUItheme, v string) func() {
	return func() {
		sat := th.Body1(v)
		sat.Color = theme.HexARGB(th.Color.Light)
		sat.Font.Typeface = "bariol"
		sat.Layout(gtx)
	}
}

func lTtime(gtx *layout.Context, th *theme.DuoUItheme, v string) func() {
	return func() {
		l := th.Body1(v)
		l.Font.Typeface = "bariol"
		l.Color = theme.HexARGB(th.Color.Light)
		l.Color = theme.HexARGB(th.Color.Hint)
		l.Layout(gtx)
	}
}

func lTamount(gtx *layout.Context, th *theme.DuoUItheme, v float64) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
			sat := th.Body1(fmt.Sprintf("%0.8f", v))
			sat.Font.Typeface = "bariol"
			sat.Color = theme.HexARGB(th.Color.Light)
			sat.Layout(gtx)
		})
	}
}
