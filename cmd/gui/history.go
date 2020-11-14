package gui

import (
	l "gioui.org/layout"
	"golang.org/x/exp/shiny/materialdesign/icons"

	"github.com/p9c/pod/pkg/gui/p9"
)

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
				wg.th.Body1("page 000").Fn,
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
				Color("DocText").Background("DocBg").Fn,
		).
		Rigid(
			wg.th.Inset(0.25,
				wg.th.Body1("tx/page").Fn,
			).Fn,
		).Fn
}

func (wg *WalletGUI) HistoryPageStatusFilter() l.Widget {
	return wg.th.Flex().AlignMiddle().
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
		Fn
}

func (wg *WalletGUI) HistoryPage() l.Widget {
	return func(gtx l.Context) l.Dimensions {
		return wg.th.VFlex().Rigid(
			wg.th.Inset(0.25,
				wg.th.Fill("DocBg",
					wg.th.Responsive(*wg.Size, p9.Widgets{
						p9.WidgetSize{
							Widget: wg.th.VFlex().Rigid(
								wg.th.Flex().AlignMiddle().SpaceBetween().
									Rigid(
										wg.HistoryPager(),
									).
									Rigid(
										wg.HistoryPagePerPageCount(),
									).
									Fn,
							).Rigid(
								wg.th.Flex().AlignMiddle().SpaceBetween().
									Rigid(
										wg.HistoryPageStatusFilter(),
									).
									Fn,
							).Fn,
						},
						p9.WidgetSize{
							Size: 1024,
							Widget: wg.th.Flex().AlignMiddle().SpaceBetween().
								Rigid(
									wg.HistoryPager(),
								).
								Rigid(
									wg.HistoryPageStatusFilter(),
								).
								Rigid(
									wg.HistoryPagePerPageCount(),
								).
								Fn,
						},
					}).Fn,
				).Fn,
			).Fn,
		).Fn(gtx)
	}
}
