package gui

import (
	l "gioui.org/layout"
	"gioui.org/text"
	
	"github.com/p9c/pod/pkg/gui"
)

func (wg *WalletGUI) ReceivePage() l.Widget {
	return func(gtx l.Context) l.Dimensions {
		if wg.State != nil {
			// Debug(wg.State.isAddress)
			// Debug(wg.State.isAddress.Load())
			if wg.State.isAddress.Load() {
				ad := wg.State.currentReceivingAddress.Load()
				wg.currentReceiveAddress = ad.EncodeAddress()
				// var err error
				// // Debug(ad.ScriptAddress())
				// var conv []byte
				// if conv, err = bech32.ConvertBits(ad.ScriptAddress(), 8, 5, true); Check(err) {
				// }
				// if bech, err = bech32.Encode("pc", conv); Check(err) {
				// }
			}
		}
		return wg.Fill("PanelBg", l.W, wg.TextSize.V, 0, wg.VFlex().AlignMiddle().
			// Flexed(0, gui.EmptyMaxWidth()).
			Rigid(
				wg.H5("Receive").Alignment(text.Middle).Fn,
			).
			Rigid(
				wg.Body1("Scan to send or click to copy").Alignment(text.Middle).Fn,
			).
			Rigid(
				// wg.Fill("DocBg",
				wg.Inset(0.25,
					wg.VFlex().AlignMiddle().
						Rigid(
							wg.currentReceiveQR,
						).
						Rigid(
							wg.Body1(wg.currentReceiveAddress).Fn,
						).
						Rigid(
							wg.Inset(0.25,
								func(gtx l.Context) l.
								Dimensions {
									gtx.Constraints.Max.X = int(wg.TextSize.V * 29)
									return wg.inputs["receiveAmount"].Fn(gtx)
								},
							).Fn,
						).
						Rigid(
							wg.Inset(0.25,
								func(gtx l.Context) l.Dimensions {
									gtx.Constraints.Max.X = int(wg.TextSize.V * 29)
									return wg.inputs["receiveMessage"].Fn(gtx)
								},
							).Fn,
						).
						Rigid(
							wg.Inset(0.25,
								func(gtx l.Context) l.Dimensions {
									gtx.Constraints.Max.X = int(wg.TextSize.V * 29)
									return wg.ButtonLayout(wg.currentReceiveRegenClickable.SetClick(func() {
										Debug("clicked regenerate button")
										wg.currentReceiveGetNew.Store(true)
									})).Background("Primary").
										Embed(
											wg.Inset(0.5,
												wg.H6("regenerate").Fn,
											).Fn,
										).
										Fn(gtx)
								}).
								Fn,
						).
						Fn,
				).Fn,
				// l.W, wg.TextSize.V).Fn,
			).
			Flexed(0, gui.EmptyMaxWidth()).
			Fn).Fn(gtx)
	}
}
