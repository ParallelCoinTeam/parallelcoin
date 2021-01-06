package gui

import (
	l "gioui.org/layout"
	"gioui.org/text"
	
	"github.com/p9c/pod/pkg/gui"
)

func (wg *WalletGUI) ReceivePage() l.Widget {
	return func(gtx l.Context) l.Dimensions {
		var ra string
		// Debug(wg.State)
		if wg.State != nil {
			// Debug(wg.State.isAddress)
			// Debug(wg.State.isAddress.Load())
			if wg.State.isAddress.Load() {
				ra = wg.State.currentReceivingAddress.Load().EncodeAddress()
			}
		}
		return wg.Fill("PanelBg",
			wg.VFlex().AlignMiddle().
				// Flexed(0, gui.EmptyMaxWidth()).
				Rigid(
					wg.H5("ParallelCoin Pod Gio Wallet").Alignment(text.Middle).Fn,
				).
				Rigid(
					wg.Fill("DocBg",
						wg.Inset(0.5,
							wg.Body1(ra).Fn,
							// gui.EmptyMaxWidth(),
							// wg.VFlex().
							// 	AlignMiddle().
							// 	Rigid(
							//
							// 		wg.VFlex().AlignMiddle().
							// 			Rigid(
							// 				wg.Inset(0.25,
							// 					wg.Caption("Built from git repository:").
							// 						Font("bariol bold").Fn,
							// 				).Fn,
							// 			).
							// 			Rigid(
							// 				wg.Caption(version.URL).Fn,
							// 			).
							// 			Fn,
							//
							// 	).
							// 	Rigid(
							//
							// 		wg.VFlex().AlignMiddle().
							// 			Rigid(
							// 				wg.Inset(0.25,
							// 					wg.Caption("GitRef:").
							// 						Font("bariol bold").Fn,
							// 				).Fn,
							// 			).
							// 			Rigid(
							// 				wg.Caption(version.GitRef).Fn,
							// 			).
							// 			Fn,
							//
							// 	).
							// 	Rigid(
							//
							// 		wg.VFlex().AlignMiddle().
							// 			Rigid(
							// 				wg.Inset(0.25,
							// 					wg.Caption("GitCommit:").
							// 						Font("bariol bold").Fn,
							// 				).Fn,
							// 			).
							// 			Rigid(
							// 				wg.Caption(version.GitCommit).Fn,
							// 			).
							// 			Fn,
							//
							// 	).
							// 	Rigid(
							//
							// 		wg.VFlex().AlignMiddle().
							// 			Rigid(
							// 				wg.Inset(0.25,
							// 					wg.Caption("BuildTime:").
							// 						Font("bariol bold").Fn,
							// 				).Fn,
							// 			).
							// 			Rigid(
							// 				wg.Caption(version.BuildTime).Fn,
							// 			).
							// 			Fn,
							//
							// 	).
							// 	Rigid(
							//
							// 		wg.VFlex().AlignMiddle().
							// 			Rigid(
							// 				wg.Inset(0.25,
							// 					wg.Caption("Tag:").
							// 						Font("bariol bold").Fn,
							// 				).Fn,
							// 			).
							// 			Rigid(
							// 				wg.Caption(version.Tag).Fn,
							// 			).
							// 			Fn,
							//
							// 	).
							// 	Rigid(
							// 		wg.Icon().Scale(gui.Scales["H6"]).
							// 			Color("DocText").
							// 			Src(&p9icons.Gio).
							// 			Fn,
							// 	).
							// 	Rigid(
							// 		wg.Caption("powered by Gio").Fn,
							// 	).
							// 	Fn,
						).Fn,
						l.W, wg.TextSize.V,
					).Fn,
				).
				Flexed(0, gui.EmptyMaxWidth()).
				Fn,
			l.W, wg.TextSize.V).Fn(gtx)
	}
}
