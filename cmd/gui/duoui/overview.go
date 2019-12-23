package duoui

import (
	"fmt"
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/unit"
	"image/color"
)

func DuoUIoverview(duo *DuoUI) {
	duo.GetDuOSbalance()
	duo.GetDuOSunconfirmedBalance()
	duo.GetDuOSblockHeight()
	duo.GetDuOStatus()
	duo.GetDuOSlocalLost()
	duo.GetDuOSdifficulty()
	duo.comp.Overview.Layout.Layout(duo.gc,
		layout.Rigid(func() {
			// OverviewTop <<<
			duo.comp.OverviewTop.Layout.Layout(duo.gc,
				layout.Flexed(0.38, func() {
					helpers.DuoUIdrawRectangle(duo.gc, duo.cs.Width.Max-30, 180, color.RGBA{A: 0xff, R: 0x30, G: 0x30, B: 0x30}, 9.9, 9.9, 9.9, 9.9, unit.Dp(0))
					in := layout.Inset{
						Top:    unit.Dp(15),
						Right:  unit.Dp(30),
						Bottom: unit.Dp(15),
						Left:   unit.Dp(30),
					}
					in.Layout(duo.gc, func() {
						layout.Flex{
							Axis: layout.Vertical,
						}.Layout(duo.gc,
							layout.Rigid(func() {
								balanceTxt := duo.th.H6("Balance :")
								balanceTxt.Color = duo.conf.StatusTextColor
								balanceTxt.Layout(duo.gc)
							}),
							layout.Rigid(func() {
								al := layout.Align(layout.End)
								al.Layout(duo.gc, func() {
									balanceVal := duo.th.H4(duo.rc.Balance + " " + duo.conf.Abbrevation)
									balanceVal.Color = duo.conf.StatusTextColor
									balanceVal.Layout(duo.gc)
								})
							}),
							layout.Rigid(func() {
								al := layout.Align(layout.End)
								al.Layout(duo.gc, func() {
									balanceUnconfirmed := duo.th.H6("Unconfirmed :" + duo.rc.Unconfirmed)
									balanceUnconfirmed.Color = duo.conf.StatusTextColor
									balanceUnconfirmed.Layout(duo.gc)
								})
							}),
							layout.Rigid(func() {
								al := layout.Align(layout.End)
								al.Layout(duo.gc, func() {
									txsNumber := duo.th.H6("Transactions :" + fmt.Sprint(duo.rc.Transactions.TxsNumber))
									txsNumber.Color = duo.conf.StatusTextColor
									txsNumber.Layout(duo.gc)
								})
							}),

						)

					})
				}),
				layout.Flexed(0.62, func() {
					helpers.DuoUIdrawRectangle(duo.gc, duo.cs.Width.Max, 180, color.RGBA{A: 0xff, R: 0x30, G: 0x30, B: 0x30}, 9.9, 9.9, 9.9, 9.9, unit.Dp(0))
					in := layout.UniformInset(unit.Dp(60))
					in.Layout(duo.gc, func() {
						bal := duo.th.H3("Balance :" + duo.rc.Balance + " DUO")
						bal.Color = color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}
						bal.Layout(duo.gc)
					})
				}))
			// OverviewTop >>>
		}),
		layout.Flexed(1, func() {
			// OverviewBottom <<<
			in := layout.Inset{
				Top: unit.Dp(30),
			}
			in.Layout(duo.gc, func() {
				cs := duo.gc.Constraints
				duo.comp.OverviewBottom.Layout.Layout(duo.gc,
					layout.Flexed(0.76, func() {
						helpers.DuoUIdrawRectangle(duo.gc, duo.cs.Width.Max-30, cs.Height.Max, color.RGBA{A: 0xff, R: 0x30, G: 0x30, B: 0x30}, 9.9, 9.9, 9.9, 9.9, unit.Dp(0))
						in := layout.UniformInset(unit.Dp(60))
						in.Layout(duo.gc, func() {
							bal := duo.th.H3("Balance :" + duo.rc.Balance + " DUO")
							bal.Color = color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}
							bal.Layout(duo.gc)
						})
					}),
					layout.Flexed(0.24, func() {
						helpers.DuoUIdrawRectangle(duo.gc, duo.cs.Width.Max, cs.Height.Max, color.RGBA{A: 0xff, R: 0x30, G: 0x30, B: 0x30}, 9.9, 9.9, 9.9, 9.9, unit.Dp(0))
						in := layout.UniformInset(unit.Dp(30))
						in.Layout(duo.gc, func() {

							duo.comp.Status.Layout.Layout(duo.gc,
								// Balance status item
								layout.Rigid(func() {
									duo.comp.StatusItem.Layout.Layout(duo.gc,
										layout.Rigid(func() {
											balanceTxt := duo.th.H5("Balance :")
											balanceTxt.Color = duo.conf.StatusTextColor
											balanceTxt.Layout(duo.gc)
										}),
										layout.Rigid(func() {
											balance := duo.th.H5(duo.rc.Balance + " " + duo.conf.Abbrevation)
											balance.Color = duo.conf.StatusTextColor
											balance.Layout(duo.gc)
										}),

									)
								}),
								// Block height status item
								layout.Rigid(func() {
									duo.comp.StatusItem.Layout.Layout(duo.gc,
										layout.Rigid(func() {
											blockheightTxt := duo.th.H5("Block Height :")
											blockheightTxt.Color = duo.conf.StatusTextColor
											blockheightTxt.Layout(duo.gc)
										}),
										layout.Rigid(func() {
											blockheightVal := duo.th.H5(fmt.Sprint(duo.rc.BlockHeight))
											blockheightVal.Color = duo.conf.StatusTextColor
											blockheightVal.Layout(duo.gc)
										}),

									)
								}),

								// Difficulty height status item
								layout.Rigid(func() {
									duo.comp.StatusItem.Layout.Layout(duo.gc,
										layout.Rigid(func() {
											difficulty := duo.th.H5("Difficulty :")
											difficulty.Color = duo.conf.StatusTextColor
											difficulty.Layout(duo.gc)
										}),
										layout.Rigid(func() {
											difficulty := duo.th.H5(fmt.Sprintf("%f", duo.rc.Difficulty))
											difficulty.Color = duo.conf.StatusTextColor
											difficulty.Layout(duo.gc)
										}),

									)
								}),

								// Connections status item
								layout.Rigid(func() {
									duo.comp.StatusItem.Layout.Layout(duo.gc,
										layout.Rigid(func() {
											connections := duo.th.H5("Connections :")
											connections.Color = duo.conf.StatusTextColor
											connections.Layout(duo.gc)
										}),
										layout.Rigid(func() {
											connections := duo.th.H5(fmt.Sprint(duo.rc.Connections))
											connections.Color = duo.conf.StatusTextColor
											connections.Layout(duo.gc)
										}),

									)
								}),

							)
						})
					}))
				// OverviewBottom >>>
			})
		}),
	)
}
