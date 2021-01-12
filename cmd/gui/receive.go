package gui

import (
	"fmt"
	
	l "gioui.org/layout"
	"gioui.org/text"
	"github.com/atotto/clipboard"
	
	"github.com/p9c/pod/pkg/gui"
)

const inputWidth float32 = 20
const Break1 = 48

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
		if wg.ReceiveAddressbook == nil {
			wg.ReceiveAddressbook = wg.Inset(0.25, wg.H1("addressbook").Alignment(text.End).Fn).Fn
		}
		var widgets []l.Widget
		if *wg.Size < int(wg.TextSize.Scale(Break1).V) {
			// assemble the list for the small, scrolling list view
			widgets = []l.Widget{
				wg.Inset(
					0.25,
					wg.Body2("Scan to send or click to copy").Alignment(text.Middle).Fn,
				).Fn,
				wg.Flex().AlignMiddle().
					Flexed(0.5, gui.EmptyMaxWidth()).
					Rigid(
						wg.ButtonLayout(
							wg.currentReceiveCopyClickable.SetClick(
								func() {
									qrText := fmt.Sprintf(
										"parallelcoin:%s?amount=%s&message=%s",
										wg.State.currentReceivingAddress.Load().EncodeAddress(),
										wg.inputs["receiveAmount"].GetText(),
										wg.inputs["receiveMessage"].GetText(),
									)
									Debug("clicked qr code copy clicker")
									if err := clipboard.WriteAll(qrText); Check(err) {
									}
								},
							),
						).
							// CornerRadius(0.5).
							// Corners(gui.NW | gui.SW | gui.NE).
							Background("white").
							Embed(
								wg.Inset(
									0.125,
									wg.Image().Src(*wg.currentReceiveQRCode).Scale(1).Fn,
								).Fn,
							).Fn,
					).
					Flexed(0.5, gui.EmptyMaxWidth()).
					Fn,
				// wg.Inset(
				// 	0.25,
				// 	wg.Caption(wg.currentReceiveAddress).Alignment(text.Middle).Font("go regular").Fn,
				// ).Fn,
				func(gtx l.Context) l.
				Dimensions {
					// gtx.Constraints.Max.X, gtx.Constraints.Min.X = int(wg.TextSize.V * 17),  int(wg.TextSize.V * 17)
					return wg.inputs["receiveSmallAmount"].Fn(gtx)
				},
				func(gtx l.Context) l.Dimensions {
					// gtx.Constraints.Max.X, gtx.Constraints.Min.X = int(wg.TextSize.V * 17),  int(wg.TextSize.V * 17)
					return wg.inputs["receiveSmallMessage"].Fn(gtx)
				},
				wg.ButtonLayout(
					wg.currentReceiveRegenClickable.SetClick(
						func() {
							Debug("clicked regenerate button")
							wg.currentReceiveGetNew.Store(true)
						},
					),
				).
					Background("Primary").
					Embed(
						wg.Inset(
							0.25,
							wg.H6("regenerate").Color("Light").Fn,
						).Fn,
					).
					Fn,
				wg.ReceiveAddressbook,
			}
		}
		le := func(gtx l.Context, index int) l.Dimensions {
			return widgets[index](gtx)
		}
		return wg.Responsive(
			*wg.Size, gui.Widgets{
				{Size: 0,
					Widget:
					wg.Flex().Flexed(
						1,
						wg.Fill(
							"DocBg", l.W, 0, 0,
							wg.Inset(
								0.25,
								wg.lists["receive"].
									Vertical().
									Length(len(widgets)).
									ListElement(le).Fn,
							).Fn,
						).Fn,
					).
						Fn,
				},
				{
					Size: Break1,
					Widget:
					wg.Fill(
						"PanelBg", l.W, wg.TextSize.V, 0,
						wg.Flex().AlignMiddle().Rigid(
							wg.VFlex().AlignMiddle().
								// Flexed(1, gui.EmptyMaxWidth()).
								Rigid(
									wg.VFlex().AlignMiddle().
										Rigid(
											wg.Inset(
												0.25,
												wg.Body2("Scan to send or click to copy").Alignment(text.Middle).Fn,
											).Fn,
										).
										Rigid(
											wg.currentReceiveQR,
										).
										// Rigid(
										// 	wg.Inset(
										// 		0.25,
										// 		wg.Caption(wg.currentReceiveAddress).Font("go regular").Fn,
										// 	).Fn,
										// ).
										Rigid(
											wg.Inset(
												0.25,
												func(gtx l.Context) l.
												Dimensions {
													gtx.Constraints.Max.X, gtx.Constraints.Min.X = int(wg.TextSize.V*inputWidth), int(wg.TextSize.V*inputWidth)
													return wg.inputs["receiveAmount"].Fn(gtx)
												},
											).Fn,
										).
										Rigid(
											wg.Inset(
												0.25,
												func(gtx l.Context) l.Dimensions {
													gtx.Constraints.Max.X, gtx.Constraints.Min.X = int(wg.TextSize.V*inputWidth), int(wg.TextSize.V*inputWidth)
													return wg.inputs["receiveMessage"].Fn(gtx)
												},
											).Fn,
										).
										Fn,
								).
								
								
								Rigid(
									wg.Inset(
										0.25,
										func(gtx l.Context) l.Dimensions {
											gtx.Constraints.Max.X, gtx.Constraints.Min.X = int(wg.TextSize.V*inputWidth), int(wg.TextSize.V*inputWidth)
											return wg.ButtonLayout(
												wg.currentReceiveRegenClickable.SetClick(
													func() {
														Debug("clicked regenerate button")
														wg.currentReceiveGetNew.Store(true)
													},
												),
											).
												// CornerRadius(0.5).Corners(gui.NW | gui.SW | gui.NE).
												Background("Primary").
												Embed(
													wg.Inset(
														0.5,
														wg.H6("regenerate").Color("Light").Fn,
													).Fn,
												).
												Fn(gtx)
										},
									).
										Fn,
								).
								Fn,
						).
							Flexed(
								1, wg.Flex().Rigid(
									wg.Fill(
										"DocBg", l.Center, wg.TextSize.V, 0,
										wg.Inset(
											0.25,
											wg.Flex().Flexed(1,
												wg.VFlex().Flexed(1,
													wg.ReceiveAddressbook,
												).Fn,
											).Fn,
										).Fn,
									).Fn,
								).
									Fn,
							).
							Fn,
					).
						Fn,
				},
				
				{
					Size: 64,
					Widget:
					wg.Fill(
						"PanelBg", l.W, wg.TextSize.V, 0,
						wg.Flex().AlignMiddle().Rigid(
							wg.VFlex().AlignMiddle().
								// Flexed(1, gui.EmptyMaxWidth()).
								Rigid(
									wg.VFlex().AlignMiddle().
										Rigid(
											wg.Inset(
												0.25,
												wg.Body2("Scan to send or click to copy").Alignment(text.Middle).Fn,
											).Fn,
										).
										Rigid(
											wg.currentReceiveQR,
										).
										// Rigid(
										// 	wg.Inset(
										// 		0.25,
										// 		wg.Caption(wg.currentReceiveAddress).Font("go regular").Fn,
										// 	).Fn,
										// ).
										Rigid(
											wg.Inset(
												0.25,
												func(gtx l.Context) l.
												Dimensions {
													gtx.Constraints.Max.X, gtx.Constraints.Min.X = int(wg.TextSize.V*inputWidth), int(wg.TextSize.V*inputWidth)
													return wg.inputs["receiveAmount"].Fn(gtx)
												},
											).Fn,
										).
										Rigid(
											wg.Inset(
												0.25,
												func(gtx l.Context) l.Dimensions {
													gtx.Constraints.Max.X, gtx.Constraints.Min.X = int(wg.TextSize.V*inputWidth), int(wg.TextSize.V*inputWidth)
													return wg.inputs["receiveMessage"].Fn(gtx)
												},
											).Fn,
										).
										Fn,
								).
								
								
								Rigid(
									wg.Inset(
										0.25,
										func(gtx l.Context) l.Dimensions {
											gtx.Constraints.Max.X, gtx.Constraints.Min.X = int(wg.TextSize.V*inputWidth), int(wg.TextSize.V*inputWidth)
											return wg.ButtonLayout(
												wg.currentReceiveRegenClickable.SetClick(
													func() {
														Debug("clicked regenerate button")
														wg.currentReceiveGetNew.Store(true)
													},
												),
											).
												// CornerRadius(0.5).Corners(gui.NW | gui.SW | gui.NE).
												Background("Primary").
												Embed(
													wg.Inset(
														0.5,
														wg.H6("regenerate").Color("Light").Fn,
													).Fn,
												).
												Fn(gtx)
										},
									).
										Fn,
								).
								Fn,
						).
							Flexed(
								1, wg.Flex().Rigid(
									wg.Fill(
										"DocBg", l.Center, wg.TextSize.V, 0,
										wg.Inset(
											0.25,
											wg.Flex().Flexed(1,
												wg.VFlex().Flexed(
													1,
													wg.ReceiveAddressbook,
												).Fn,
											).Fn,
										).Fn,
									).Fn,
								).
									Fn,
							).
							Fn,
					).
						Fn,
				},
				{
					Size: 96,
					Widget:
					wg.Fill(
						"PanelBg", l.W, wg.TextSize.V, 0,
						wg.Flex().AlignMiddle().Rigid(
							wg.VFlex().AlignMiddle().
								// Flexed(1, gui.EmptyMaxWidth()).
								Rigid(
									wg.Flex().AlignMiddle().
										Rigid(
											wg.VFlex().AlignMiddle().
												Rigid(
													wg.Inset(
														0.25,
														wg.Body2("Scan to send or click to copy").Alignment(text.Middle).Fn,
													).Fn,
												).
												Rigid(
													wg.currentReceiveQR,
												).
												// Rigid(
												// 	wg.Inset(
												// 		0.25,
												// 		wg.Caption(wg.currentReceiveAddress).Font("go regular").Fn,
												// 	).Fn,
												// ).
												Fn,
										).
										Rigid(
											wg.VFlex().AlignMiddle().
												Rigid(
													wg.Inset(
														0.25,
														func(gtx l.Context) l.
														Dimensions {
															gtx.Constraints.Max.X, gtx.Constraints.Min.X = int(wg.TextSize.V*inputWidth), int(wg.TextSize.V*inputWidth)
															return wg.inputs["receiveAmount"].Fn(gtx)
														},
													).Fn,
												).
												Rigid(
													wg.Inset(
														0.25,
														func(gtx l.Context) l.Dimensions {
															gtx.Constraints.Max.X, gtx.Constraints.Min.X = int(wg.TextSize.V*inputWidth), int(wg.TextSize.V*inputWidth)
															return wg.inputs["receiveMessage"].Fn(gtx)
														},
													).Fn,
												).
												Rigid(
													wg.Inset(
														0.25,
														func(gtx l.Context) l.Dimensions {
															gtx.Constraints.Max.X, gtx.Constraints.Min.X = int(wg.TextSize.V*inputWidth), int(wg.TextSize.V*inputWidth)
															return wg.ButtonLayout(
																wg.currentReceiveRegenClickable.SetClick(
																	func() {
																		Debug("clicked regenerate button")
																		wg.currentReceiveGetNew.Store(true)
																	},
																),
															).
																// CornerRadius(0.5).Corners(gui.NW | gui.SW | gui.NE).
																Background("Primary").
																Embed(
																	wg.Inset(
																		0.5,
																		wg.H6("regenerate").Color("Light").Fn,
																	).Fn,
																).
																Fn(gtx)
														},
													).
														Fn,
												).Fn,
										).
										Fn,
								).
								Fn,
						).
							Flexed(
								1, wg.Flex().Rigid(
									wg.Fill(
										"DocBg", l.Center, wg.TextSize.V, 0,
										wg.Inset(
											0.25,
											wg.Flex().Flexed(1,
												wg.VFlex().Flexed(1,
													wg.ReceiveAddressbook,
												).Fn,
											).Fn,
										).Fn,
									).Fn,
								).
									Fn,
							).
							Fn,
					).
						Fn,
				},
			},
		).
			Fn(gtx)
	}
}
