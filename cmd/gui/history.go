package gui

import (
	l "gioui.org/layout"
)

func (wg *WalletGUI) HistoryPage() l.Widget {
	return func(gtx l.Context) l.Dimensions {
		return wg.th.VFlex().Rigid(
			wg.th.Fill("DocBg",
				wg.th.Inset(0.25,
					wg.th.Flex().
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
								func(gtx l.Context) l.Dimensions {
									return wg.th.CheckBox(wg.checkables["showGenerate"]).Text("generate").Fn(gtx)
								},
								// wg.th.Body1("generated").Fn,
							).Fn,
						).
						Rigid(
							wg.th.Inset(0.25,
								wg.th.Body1("sent").Fn,
							).Fn,
						).
						Rigid(
							wg.th.Inset(0.25,
								wg.th.Body1("received").Fn,
							).Fn,
						).
						Fn,
				).Fn,
			).Fn,
		).Fn(gtx)
	}
}
