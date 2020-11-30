package gui

import (
	"fmt"

	l "gioui.org/layout"
)

func (wg *WalletGUI) ReceivePage() l.Widget {
	le := func(gtx l.Context, index int) l.Dimensions {
		return wg.th.Caption("BalaaaaaaaaaaaaaaaO_" + fmt.Sprint(index)).Color("DocBg").Fn(gtx)
	}
	return func(gtx l.Context) l.Dimensions {
		return wg.th.VFlex().
			Rigid(
				wg.receiveTop(),
			).
			Flexed(1,
				wg.th.Inset(0.25, wg.th.Fill("DocBg", wg.th.Inset(0.25,
					wg.lists["received"].Vertical().Length(len(wg.sendAddresses)).ListElement(le).Fn,
				).Fn).Fn).Fn,
			).Fn(gtx)
	}
}

func (wg *WalletGUI) receiveTop() l.Widget {
	return wg.th.Inset(0.25,
		wg.th.Fill("DocBg",
			wg.th.Inset(0.25,
				wg.th.VFlex().
					Rigid(
						wg.th.Inset(0.25,
							wg.th.Flex().
								SpaceBetween().
								Rigid(
									wg.th.Inset(0.0, wg.th.Fill("DocBg", wg.th.Inset(0.1, wg.th.Caption("Use this form to request payments. All fields are optional.").Color("DocText").Fn).Fn).Fn).Fn,
								).
								Rigid(
									wg.th.Inset(0.0, wg.th.Fill("DocBg", wg.th.Inset(0.1, wg.th.Caption("Label:").Color("DocText").Fn).Fn).Fn).Fn,
								).Fn,
						).Fn,
					).Rigid(
					wg.th.Inset(0.25,
						wg.th.Flex().
							SpaceBetween().
							Rigid(
								wg.th.Inset(0.0, wg.th.Fill("DocBg", wg.th.Inset(0.1, wg.th.Caption("Label:").Color("DocText").Fn).Fn).Fn).Fn,
							).
							Rigid(
								wg.th.Inset(0.0, wg.th.Fill("DocBg", wg.th.Inset(0.1,
									// wg.Caption("0.00000 DUO/kb").Color("DocText")
									wg.inputs["receiveAmount"].Fn).Fn).Fn).Fn,
							).Fn,
					).Fn,
				).Rigid(
					wg.th.Inset(0.25,
						wg.th.Flex().
							SpaceBetween().
							Rigid(
								wg.th.Inset(0.0, wg.th.Fill("DocBg", wg.th.Inset(0.1, wg.th.Caption("Amount:").Color("DocText").Fn).Fn).Fn).Fn,
							).
							Rigid(
								wg.th.Inset(0.0, wg.th.Fill("DocBg", wg.th.Inset(0.1,
									// wg.Caption("0.00000 DUO/kb").Color("DocText")
									wg.inputs["receiveLabel"].Fn).Fn).Fn).Fn,
							).Fn,
					).Fn,
				).Rigid(
					wg.th.Inset(0.25,
						wg.th.Flex().
							SpaceBetween().
							Rigid(
								wg.th.Inset(0.0, wg.th.Fill("DocBg", wg.th.Inset(0.1, wg.th.Caption("Message:").Color("DocText").Fn).Fn).Fn).Fn,
							).
							Rigid(
								wg.th.Inset(0.0, wg.th.Fill("DocBg", wg.th.Inset(0.1,
									// wg.Caption("0.00000 DUO/kb").Color("DocText")
									wg.inputs["receiveMessage"].Fn).Fn).Fn).Fn,
							).Fn,
					).Fn,
				).Rigid(
					wg.th.Inset(0.25,
						wg.th.Flex().
							SpaceBetween().
							Rigid(
								wg.th.Inset(0.25,
									wg.buttonText(wg.clickables["receiveCreateNewAddress"], "Create new receiving address", wg.Send),
								).Fn,
							).
							Rigid(
								wg.th.Inset(0.25,
									wg.buttonText(wg.clickables["receiveClear"], "Clear", wg.ClearAllAddresses),
								).Fn,
							).
						Fn,
					).Fn,
				).Fn,
			).Fn,
		).Fn,
	).Fn
}
