package gui

import (
	l "gioui.org/layout"
	
	"github.com/p9c/pod/pkg/gui"
)

func (wg *WalletGUI) HistoryPage() l.Widget {
	if wg.HistoryWidget == nil {
		wg.HistoryWidget = func(gtx l.Context) l.Dimensions {
			return l.Dimensions{Size: gtx.Constraints.Max}
		}
	}
	return func(gtx l.Context) l.Dimensions {
		return wg.VFlex().
			Rigid(
				// wg.Fill("DocBg", l.Center, 0, 0,
				// 	wg.Inset(0.25,
				wg.Responsive(*wg.Size, gui.Widgets{
					{
						Widget: wg.VFlex().
							Flexed(1, wg.HistoryPageView()).
							// Rigid(
							// 	// 	wg.Fill("DocBg",
							// 	wg.Flex().AlignMiddle().SpaceBetween().
							// 		Flexed(0.5, gui.EmptyMaxWidth()).
							// 		Rigid(wg.HistoryPageStatusFilter()).
							// 		Flexed(0.5, gui.EmptyMaxWidth()).
							// 		Fn,
							// 	// 	).Fn,
							// ).
							// Rigid(
							// 	wg.Fill("DocBg",
							// 		wg.Flex().AlignMiddle().SpaceBetween().
							// 			Rigid(wg.HistoryPager()).
							// 			Rigid(wg.HistoryPagePerPageCount()).
							// 			Fn,
							// 	).Fn,
							// ).
							Fn,
					},
					{
						Size: 64,
						Widget: wg.VFlex().
							Flexed(1, wg.HistoryPageView()).
							// Rigid(
							// 	// 	wg.Fill("DocBg",
							// 	wg.Flex().AlignMiddle().SpaceBetween().
							// 		// 			Rigid(wg.HistoryPager()).
							// 		Flexed(0.5, gui.EmptyMaxWidth()).
							// 		Rigid(wg.HistoryPageStatusFilter()).
							// 		Flexed(0.5, gui.EmptyMaxWidth()).
							// 		// 			Rigid(wg.HistoryPagePerPageCount()).
							// 		Fn,
							// 	// 	).Fn,
							// ).
							Fn,
					},
				}).Fn,
				// ).Fn,
				// ).Fn,
			).Fn(gtx)
	}
}

func (wg *WalletGUI) HistoryPageView() l.Widget {
	return wg.VFlex().
		Rigid(
			// wg.Fill("DocBg", l.Center, wg.TextSize.V, 0,
			// 	wg.Inset(0.25,
			wg.HistoryWidget,
			// ).Fn,
			// ).Fn,
		).Fn
}

func (wg *WalletGUI) HistoryPageStatusFilter() l.Widget {
	return wg.Flex().AlignMiddle().
		Rigid(
			wg.Inset(0.25,
				wg.Caption("show").Fn,
			).Fn,
		).
		Rigid(
			wg.Inset(0.25,
				func(gtx l.Context) l.Dimensions {
					return wg.CheckBox(wg.bools["showGenerate"]).
						TextColor("DocText").
						TextScale(1).
						Text("generate").
						IconScale(1).
						Fn(gtx)
				},
			).Fn,
		).
		Rigid(
			wg.Inset(0.25,
				func(gtx l.Context) l.Dimensions {
					return wg.CheckBox(wg.bools["showSent"]).
						TextColor("DocText").
						TextScale(1).
						Text("sent").
						IconScale(1).
						Fn(gtx)
				},
			).Fn,
		).
		Rigid(
			wg.Inset(0.25,
				func(gtx l.Context) l.Dimensions {
					return wg.CheckBox(wg.bools["showReceived"]).
						TextColor("DocText").
						TextScale(1).
						Text("received").
						IconScale(1).
						Fn(gtx)
				},
			).Fn,
		).
		Rigid(
			wg.Inset(0.25,
				func(gtx l.Context) l.Dimensions {
					return wg.CheckBox(wg.bools["showImmature"]).
						TextColor("DocText").
						TextScale(1).
						Text("immature").
						IconScale(1).
						Fn(gtx)
				},
			).Fn,
		).
		Fn
}
