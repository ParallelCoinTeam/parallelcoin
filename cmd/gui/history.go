package gui

import (
	l "gioui.org/layout"
	"golang.org/x/exp/shiny/materialdesign/icons"

	"github.com/p9c/pod/pkg/gui/p9"
)

func (wg *WalletGUI) HistoryPage() l.Widget {

	return func(gtx l.Context) l.Dimensions {
		return wg.th.VFlex().
			Rigid(
				wg.th.Inset(0.25,
					wg.th.Fill("PanelBg",
						wg.th.Responsive(*wg.Size, p9.Widgets{
							{
								Widget: wg.th.VFlex().
									Flexed(1, wg.HistoryPageView()).
									Rigid(
										wg.th.Fill("DocBg",
											wg.th.Flex().AlignMiddle().SpaceBetween().
												Flexed(0.5, p9.EmptyMaxWidth()).
												Rigid(wg.HistoryPageStatusFilter()).
												Flexed(0.5, p9.EmptyMaxWidth()).
												Fn,
										).Fn,
									).
									Rigid(
										wg.th.Fill("DocBg",
											wg.th.Flex().AlignMiddle().SpaceBetween().
												Rigid(wg.HistoryPager()).
												Rigid(wg.HistoryPagePerPageCount()).
												Fn,
										).Fn,
									).
									Fn,
							},
							{
								Size: 800,
								Widget: wg.th.VFlex().
									Flexed(1, wg.HistoryPageView()).
									Rigid(
										wg.th.Fill("DocBg",
											wg.th.Flex().AlignMiddle().SpaceBetween().
												Rigid(wg.HistoryPager()).
												Rigid(wg.HistoryPageStatusFilter()).
												Rigid(wg.HistoryPagePerPageCount()).
												Fn,
										).Fn,
									).
									Fn,
							},
						}).Fn,
						// ).Fn,
					).Fn,
				).Fn,
			).Fn(gtx)
	}
}

func (wg *WalletGUI) HistoryPageView() l.Widget {
	return wg.th.Fill("DocBg",

		p9.EmptySpace(0, 0),
	).Fn
}

func (wg *WalletGUI) HistoryPager() l.Widget {
	return wg.th.Flex().AlignMiddle().
		Rigid(
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
		).
		Rigid(
			wg.th.Inset(0.25,
				wg.th.Caption("page 1/20").Fn,
			).Fn,
		).
		Rigid(
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
		).Fn
}

func (wg *WalletGUI) HistoryPagePerPageCount() l.Widget {
	return wg.th.Flex().AlignMiddle().
		Rigid(
			wg.incdecs["transactionsPerPage"].
				Color("DocText").Background("DocBg").Scale(p9.Scales["Caption"]).Fn,
		).
		Rigid(
			wg.th.Inset(0.25,
				wg.th.Caption("tx/page").Fn,
			).Fn,
		).Fn
}

func (wg *WalletGUI) HistoryPageStatusFilter() l.Widget {
	return wg.th.Flex().AlignMiddle().
		Rigid(
			wg.th.Inset(0.25,
				wg.th.Caption("show").Fn,
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
		Fn
}
