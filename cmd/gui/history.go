package gui

import (
	l "gioui.org/layout"
	"golang.org/x/exp/shiny/materialdesign/icons"

	"github.com/p9c/pod/pkg/gui/p9"
)

func (wg *WalletGUI) HistoryPage() l.Widget {
	return func(gtx l.Context) l.Dimensions {
		return wg.th.VFlex().Rigid(
			wg.th.Fill("DocBg",
				wg.th.Inset(0.25,
					wg.th.Flex().AlignMiddle().
						Rigid(
							// p9.If(wg.incdecs["transactionsPerPage"].GetCurrent() > 10,
							wg.th.IconButton(wg.clickables["txPageBack"]).
								Background("Transparent").
								Color("DocText").
								Scale(1).
								Icon(
									wg.th.Icon().Color("DocText").
										Scale(1).
										Src(&icons.NavigationArrowBack),
								).
								Fn,
							// p9.EmptySpace(0, 0),
							// ),
						).
						Rigid(
							wg.incdecs["transactionsPerPage"].
								Color("DocText").Background("DocBg").Fn,
						).
						Rigid(
							wg.th.Inset(0.25,
								wg.th.Body1("tx/page").Fn,
							).Fn,
						).
						Rigid(
							wg.th.Inset(0.25,
								wg.th.Body1("show").Font("bariol bold").Fn,
							).Fn,
						).
						Rigid(
							wg.th.Inset(0.25,
								func(gtx l.Context) l.Dimensions {
									return wg.th.CheckBox(wg.bools["showGenerate"]).
										TextColor("DocText").
										TextScale(1).
										Text("generate").
										IconScale(1).
										Fn(gtx)
								},
								// wg.th.Body1("generated").Fn,
							).Fn,
						).
						Rigid(
							wg.th.Inset(0.25,
								func(gtx l.Context) l.Dimensions {
									return wg.th.CheckBox(wg.bools["showSent"]).
										TextColor("DocText").
										TextScale(1).
										Text("sent").
										IconScale(1).
										Fn(gtx)
								},
							).Fn,
						).
						Rigid(
							wg.th.Inset(0.25,
								func(gtx l.Context) l.Dimensions {
									return wg.th.CheckBox(wg.bools["showReceived"]).
										TextColor("DocText").
										TextScale(1).
										Text("received").
										IconScale(1).
										Fn(gtx)
								},
							).Fn,
						).
						Flexed(1, p9.EmptyMaxWidth()).
						Rigid(
							// p9.If(wg.incdecs["transactionsPerPage"].GetCurrent() < 100,
							wg.th.IconButton(wg.clickables["txPageForward"]).
								Background("Transparent").
								Color("DocText").
								Scale(1).
								Icon(
									wg.th.Icon().Color("DocText").
										Scale(1).
										Src(&icons.NavigationArrowForward),
								).
								Fn,
							// p9.EmptySpace(0, 0),
							// ),
						).
						Fn,
				).Fn,
			).Fn,
		).Fn(gtx)
	}
}
