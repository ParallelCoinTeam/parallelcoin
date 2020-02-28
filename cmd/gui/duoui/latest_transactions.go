package duoui

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
	"github.com/p9c/pod/cmd/gui/rcd"
)

var (
	latestTransList = &layout.List{
		Axis: layout.Vertical,
	}
)

func DuoUIlatestTransactions(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme) func() {
	return func() {
		cs := gtx.Constraints
		theme.DuoUIdrawRectangle(gtx, cs.Width.Max, cs.Height.Max, "ff424242", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
		layout.Flex{
			Axis: layout.Vertical,
		}.Layout(gtx,
			layout.Rigid(func() {
				cs := gtx.Constraints
				theme.DuoUIdrawRectangle(gtx, cs.Width.Max, 48, th.Color.Primary, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
				layout.UniformInset(unit.Dp(16)).Layout(gtx, func() {
					latestx := th.H5("LATEST TRANSACTIONS")
					latestx.Color = theme.HexARGB(th.Color.Light)
					latestx.Alignment = text.Start
					latestx.Layout(gtx)
				})
			}),
			layout.Flexed(1, func() {
				layout.UniformInset(unit.Dp(8)).Layout(gtx, func() {
					layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func() {
							cs := gtx.Constraints
							latestTransList.Layout(gtx, len(rc.Status.Wallet.LastTxs.Txs), func(i int) {
								t := rc.Status.Wallet.LastTxs.Txs[i]
								theme.DuoUIdrawRectangle(gtx, cs.Width.Max, cs.Height.Max, th.Color.Dark, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
								layout.Inset{
									Top:    unit.Dp(8),
									Right:  unit.Dp(16),
									Bottom: unit.Dp(8),
									Left:   unit.Dp(16),
								}.Layout(gtx, func() {
									layout.Flex{Axis: layout.Vertical}.Layout(gtx,
										layout.Rigid(lTtxid(gtx, th, t.TxID)),
										layout.Rigid(func() {
											layout.Flex{
												Spacing: layout.SpaceBetween,
											}.Layout(gtx,
												layout.Rigid(func() {
													layout.Flex{
														Axis: layout.Vertical,
													}.Layout(gtx,
														layout.Rigid(lTcategory(gtx, th, t.Category)),
														layout.Rigid(lTtime(gtx, th, t.Time)),
													)
												}),
												layout.Rigid(lTamount(gtx, th, t.Amount)),
											)
										}),
										layout.Rigid(line(gtx, th.Color.Hint)),
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
