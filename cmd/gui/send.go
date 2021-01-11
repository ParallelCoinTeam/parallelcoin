package gui

import (
	l "gioui.org/layout"
	"gioui.org/text"
	"golang.org/x/exp/shiny/materialdesign/icons"
	
	"github.com/p9c/pod/pkg/gui"
)

type SendAddress struct {
	AddressInput      *gui.Input
	LabelInput        *gui.Input
	AddressBookBtn    *gui.Clickable
	PasteClipboardBtn *gui.Clickable
	ClearBtn          *gui.Clickable
	AmountInput       *gui.Input
	// AmountInput       *counter.Counter
	SubtractFee     *gui.Bool
	AllAvailableBtn *gui.Clickable
}

func (wg *WalletGUI) getInput(name string, width int) l.Widget {
	return func(gtx l.Context) l.Dimensions {
		return wg.Inset(0.25,
			func(gtx l.Context) l.Dimensions {
				if width > 0 {
					gtx.Constraints.Max.X = int(wg.TextSize.V) * width
				}
				return wg.inputs[name].Fn(gtx)
			},
		).Fn(gtx)
	}
}

func (wg *WalletGUI) SendPage() l.Widget {
	return func(gtx l.Context) l.Dimensions {
		return wg.Responsive(*wg.Size, gui.Widgets{
			{
				Size: 64,
				Widget:
				wg.Fill("PanelBg", l.NE, wg.TextSize.V, 0, wg.Flex().AlignMiddle().
					Rigid(
						wg.Fill("DocBg", l.NE, wg.TextSize.V, 0,
							wg.VFlex().AlignMiddle().
								Rigid(
									wg.VFlex().AlignMiddle().
										Rigid(
											wg.getInput("sendAddress", 26),
										).
										Rigid(
											wg.getInput("sendAmount", 26),
										).
										Rigid(
											wg.getInput("sendMessage", 26),
										).
										Rigid(
											wg.Flex().
												Rigid(
													wg.Inset(0.25,
														func(gtx l.Context) l.Dimensions {
															return wg.ButtonLayout(wg.clickables["sendSend"].SetClick(func() {
																Debug("clicked regenerate button")
																wg.currentReceiveGetNew.Store(true)
															})).
																// CornerRadius(0.5).Corners(0).
																Background("Primary").
																Embed(
																	wg.Inset(0.25,
																		wg.Flex().AlignMiddle().
																			Rigid(
																				wg.Icon().
																					Scale(gui.
																						Scales["H4"]).
																					Color("Light").
																					Src(
																						&icons.
																							ContentSend,
																					).Fn,
																			).
																			Rigid(
																				wg.Inset(0.5, gui.EmptySpace(0, 0)).Fn,
																			).
																			Rigid(
																				wg.H6("send").Color("Light").Fn,
																			).
																			Rigid(
																				wg.Inset(0.5, gui.EmptySpace(0, 0)).Fn,
																			).
																			Fn,
																	).Fn,
																).
																Fn(gtx)
														}).
														Fn,
												).
												Rigid(
													wg.Inset(0.25,
														func(gtx l.Context) l.Dimensions {
															q := gtx.Queue
															gtx.Queue = nil
															defer func() {
																gtx.Queue = q
															}()
															return wg.ButtonLayout(wg.clickables["sendSave"].SetClick(func() {
																Debug("clicked regenerate button")
																wg.currentReceiveGetNew.Store(true)
															})).
																CornerRadius(0.5).Corners(gui.NW | gui.SW | gui.NE).
																Background("Primary").
																Embed(
																	wg.Inset(0.25,
																		wg.Flex().AlignMiddle().
																			Rigid(
																				wg.Icon().
																					Scale(gui.
																						Scales["H4"]).
																					Color("Light").
																					Src(
																						&icons.
																							ContentSave,
																					).Fn,
																			).
																			Rigid(
																				wg.Inset(0.5, gui.EmptySpace(0, 0)).Fn,
																			).
																			Rigid(
																				wg.H6("save").Color("Light").Fn,
																			).
																			Rigid(
																				wg.Inset(0.5, gui.EmptySpace(0, 0)).Fn,
																			).
																			Fn,
																	).Fn,
																).
																Fn(gtx)
														}).
														Fn,
												).
												Fn,
										).
										Fn,
								).
								Fn,
						).Fn,
					).
					Rigid(
						wg.Fill("DocBg", l.Center, wg.TextSize.V, 0,
							wg.Inset(0.25,
								wg.Flex().Flexed(1,
									wg.VFlex().
										Flexed(1,
											wg.H1("addressbook").Alignment(text.End).Fn,
										).
										Fn,
								).Fn,
							).Fn,
						).
							Fn,
					).
					Fn).
					Fn,
			},
			{
				Size: 0,
				Widget:
				// wg.Fill("scrim", l.NE, wg.TextSize.V, 0,
				// 	gui.EmptyMaxWidth(),
				// ).Fn,
				
				// wg.Fill("PanelBg", l.NE, wg.TextSize.V, l.W,
				wg.VFlex().AlignMiddle().
					Rigid(
						wg.Fill("DocBg", l.NE, wg.TextSize.V, 0,
							wg.VFlex().AlignMiddle().
								Rigid(
									wg.VFlex().AlignMiddle().
										Rigid(
											wg.getInput("sendAddress", 0),
										).
										Rigid(
											wg.getInput("sendAmount", 0),
										).
										Rigid(
											wg.getInput("sendMessage", 0),
										).
										Rigid(
											wg.Flex().
												Rigid(
													wg.Inset(0.25,
														func(gtx l.Context) l.Dimensions {
															return wg.ButtonLayout(wg.clickables["sendSend"].SetClick(func() {
																Debug("clicked regenerate button")
																wg.currentReceiveGetNew.Store(true)
															})).Background("Primary").
																Embed(
																	wg.Inset(0.25,
																		wg.Flex().AlignMiddle().
																			Rigid(
																				wg.Icon().
																					Scale(gui.
																						Scales["H4"]).
																					Color("Light").
																					Src(
																						&icons.ContentSend,
																					).Fn,
																			).
																			Rigid(
																				wg.Inset(0.5, gui.EmptySpace(0, 0)).Fn,
																			).
																			Rigid(
																				wg.H6("send").Color("Light").Fn,
																			).
																			Rigid(
																				wg.Inset(0.5, gui.EmptySpace(0, 0)).Fn,
																			).
																			Fn,
																	).Fn,
																).
																Fn(gtx)
														}).
														Fn,
												).
												Rigid(
													wg.Inset(0.25,
														func(gtx l.Context) l.Dimensions {
															q := gtx.Queue
															gtx.Queue = nil
															defer func() {
																gtx.Queue = q
															}()
															return wg.ButtonLayout(wg.clickables["sendSave"].SetClick(func() {
																Debug("clicked regenerate button")
																wg.currentReceiveGetNew.Store(true)
															})).Background("Primary").
																Embed(
																	wg.Inset(0.25,
																		wg.Flex().AlignMiddle().
																			Rigid(
																				wg.Icon().
																					Scale(gui.
																						Scales["H4"]).
																					Color("Light").
																					Src(
																						&icons.ContentSave,
																					).Fn,
																			).
																			Rigid(
																				wg.Inset(0.5, gui.EmptySpace(0, 0)).Fn,
																			).
																			Rigid(
																				wg.H6("save").Color("Light").Fn,
																			).
																			Rigid(
																				wg.Inset(0.5, gui.EmptySpace(0, 0)).Fn,
																			).
																			Fn,
																	).Fn,
																).
																Fn(gtx)
														}).
														Fn,
												).Fn,
										).
										Fn,
								).
								Fn,
						).Fn,
					).
					Rigid(
						wg.Inset(0.125, gui.EmptySpace(0, 0)).Fn,
					).
					Rigid(
						wg.Fill("DocBg", l.Center, wg.TextSize.V, 0,
							wg.Inset(0.25,
								wg.Flex().Flexed(1,
									wg.VFlex().
										Flexed(1,
											wg.H1("addressbook").Alignment(text.End).Fn,
										).
										Fn,
								).Fn,
							).
								Fn,
						).
							Fn,
					).
					Fn,
			},
		}).
			Fn(gtx)
	}
}
