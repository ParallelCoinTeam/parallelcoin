package component

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/p9c/gelook"

	"github.com/p9c/pod/cmd/gui/rcd"
)

var (
	latestTransList = &layout.List{
		Axis: layout.Vertical,
	}
)

func DuoUIlatestTransactions(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	return func() {

		cs := gtx.Constraints
		gelook.DuoUIdrawRectangle(gtx, cs.Width.Max, cs.Height.Max, "ff424242", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
		layout.Flex{
			Axis: layout.Vertical,
		}.Layout(gtx,
			layout.Rigid(func() {
				cs := gtx.Constraints
				gelook.DuoUIdrawRectangle(gtx, cs.Width.Max, 48, th.Colors["Primary"], [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
				layout.UniformInset(unit.Dp(16)).Layout(gtx, func() {
					latestx := th.H5("LATEST TRANSACTIONS")
					latestx.Color = gelook.HexARGB(th.Colors["Light"])
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
								gelook.DuoUIdrawRectangle(gtx, cs.Width.Max, cs.Height.Max, th.Colors["Dark"], [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
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
										layout.Rigid(HorizontalLine(gtx, 1, th.Colors["Hint"])),
									)
								})
							})
						}))
				})
			}),
		)
	}
}

func lTtxid(gtx *layout.Context, th *gelook.DuoUItheme, v string) func() {
	return func() {
		tim := th.Caption(v)
		tim.Font.Typeface = th.Fonts["Mono"]
		tim.Color = gelook.HexARGB(th.Colors["Light"])
		tim.Layout(gtx)
	}
}

func lTcategory(gtx *layout.Context, th *gelook.DuoUItheme, v string) func() {
	return func() {
		sat := th.Body1(v)
		sat.Color = gelook.HexARGB(th.Colors["Light"])
		sat.Font.Typeface = "bariol"
		sat.Layout(gtx)
	}
}

func lTtime(gtx *layout.Context, th *gelook.DuoUItheme, v string) func() {
	return func() {
		l := th.Body1(v)
		l.Font.Typeface = "bariol"
		l.Color = gelook.HexARGB(th.Colors["Light"])
		l.Color = gelook.HexARGB(th.Colors["Hint"])
		l.Layout(gtx)
	}
}

func lTamount(gtx *layout.Context, th *gelook.DuoUItheme, v float64) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
			sat := th.Body1(fmt.Sprintf("%0.8f", v))
			sat.Font.Typeface = "bariol"
			sat.Color = gelook.HexARGB(th.Colors["Light"])
			sat.Layout(gtx)
		})
	}
}
