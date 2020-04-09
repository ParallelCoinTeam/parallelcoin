package component

import (
	"fmt"
	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/nfnt/resize"
	"github.com/p9c/pod/cmd/gui/model"
	"github.com/p9c/pod/cmd/gui/rcd"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/coding/qrcode"
	"github.com/p9c/pod/pkg/gui/gel"
	"github.com/p9c/pod/pkg/gui/gelook"
	"github.com/p9c/pod/pkg/pod"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	"github.com/p9c/pod/pkg/util/interrupt"
	log "github.com/p9c/pod/pkg/util/logi"
	"image"
	"image/color"
	"strconv"
	"strings"
	"time"
)

type DuoUIcomponent struct {
	Name    string
	Version string
	Theme   *gelook.DuoUItheme
	M       interface{}
	V       func()
	C       func()
}

type Field struct {
	Field *pod.Field
}

var (
	previousBlockHashButton = new(gel.Button)
	nextBlockHashButton     = new(gel.Button)
	list                    = &layout.List{
		Axis: layout.Vertical,
	}
	buttonDialogCancel = new(gel.Button)
	buttonDialogOK     = new(gel.Button)
	buttonDialogClose  = new(gel.Button)
	buttonLog          = new(gel.Button)
	buttonSettings     = new(gel.Button)
	buttonNetwork      = new(gel.Button)
	buttonBlocks       = new(gel.Button)
	buttonConsole      = new(gel.Button)
	buttonHelp         = new(gel.Button)
	navItemIconColor   = "ffacacac"
	cornerNav          = &layout.List{
		Axis: layout.Horizontal,
	}
	footerNav = &layout.List{
		Axis: layout.Horizontal,
	}
	footerMenuItemWidth             = 48
	footerMenuItemHeight            = 48
	footerMenuItemTextSize          = 16
	footerMenuItemIconSize          = 32
	footerMenuItemPaddingVertical   = 0
	footerMenuItemPaddingHorizontal = 0
	latestTxsPanelElement           = gel.NewPanel()
	navButtonOverview               = new(gel.Button)
	navButtonSend                   = new(gel.Button)
	navButtonReceive                = new(gel.Button)
	navButtonAddressBook            = new(gel.Button)
	navButtonHistory                = new(gel.Button)
	mainNav                         = &layout.List{
		Axis: layout.Vertical,
	}
	groupsList = &layout.List{
		Axis:      layout.Horizontal,
		Alignment: layout.Start,
	}
	buttonHeader  = new(gel.Button)
	logOutputList = &layout.List{
		Axis:        layout.Vertical,
		ScrollToEnd: true,
	}
	StartupTime = time.Now()
	buttonQuit  = new(gel.Button)
	itemsList   = &layout.List{
		Axis: layout.Vertical,
	}
	singleItem = &layout.Flex{
		Axis:    layout.Horizontal,
		Spacing: layout.SpaceBetween,
	}
)

// rc, gtx, th

func DuoUIdialog(
	rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) {
	// cs := gtx.Constraints
	th.DuoUIcontainer(0, "ee000000").
		Layout(gtx, layout.Center, func() {
			cs := gtx.Constraints
			layout.Stack{
				Alignment: layout.Center,
			}.Layout(gtx,
				layout.Expanded(func() {
					rr := float32(gtx.Px(unit.Dp(0)))
					clip.Rect{
						Rect: f32.Rectangle{Max: f32.Point{
							X: float32(cs.Width.Max),
							Y: float32(cs.Height.Max),
						}},
						NE: rr, NW: rr, SE: rr, SW: rr,
					}.Op(gtx.Ops).Add(gtx.Ops)
					fill(gtx, gelook.HexARGB("ff888888"))
					pointer.Rect(image.Rectangle{
						Max: gtx.Dimensions.Size,
					}).Add(gtx.Ops)
				}),
				layout.Stacked(func() {
					if gtx.Constraints.Width.Max > 500 {
						gtx.Constraints.Width.Max = 500
					}
					layout.Center.Layout(gtx, func() {
						layout.Inset{
							Top: unit.Dp(16), Bottom: unit.Dp(16),
							Left: unit.Dp(8), Right: unit.Dp(8),
						}.Layout(gtx, func() {
							layout.Flex{
								Axis:      layout.Vertical,
								Alignment: layout.Middle,
							}.Layout(gtx, layout.Rigid(func() {
								layout.Flex{
									Axis:      layout.Horizontal,
									Alignment: layout.Middle,
								}.Layout(gtx, layout.Rigid(func() {
									layout.Inset{
										Top:    unit.Dp(0),
										Bottom: unit.Dp(8),
										Left:   unit.Dp(4),
										Right:  unit.Dp(4)}.
										Layout(gtx, func() {
											cur := th.DuoUIlabel(unit.Dp(14), rc.Dialog.Text)
											cur.Font.Typeface = th.Fonts["Primary"]
											cur.Color = th.Colors["Dark"]
											cur.Alignment = text.Start
											cur.Layout(gtx)
										})
								}),
								)
							}),
								layout.Rigid(func() {
									rc.Dialog.CustomField()
								}),
								layout.Rigid(func() {
									layout.Flex{
										Axis:      layout.Horizontal,
										Alignment: layout.Middle,
									}.Layout(gtx,
										layout.Rigid(
											dialogButton(gtx, th,
												rc.Dialog.Red, rc.Dialog.RedLabel,
												"ffcf3030", "iconCancel", "ffcf8080",
												buttonDialogCancel),
										),
										layout.Rigid(
											dialogButton(gtx, th,
												rc.Dialog.Green, rc.Dialog.GreenLabel,
												"ff30cf30", "iconOK",
												"ff80cf80", buttonDialogOK),
										),
										layout.Rigid(
											dialogButton(gtx, th,
												rc.Dialog.Orange, rc.Dialog.OrangeLabel,
												"ffcf8030", "iconClose",
												"ffcfa880",
												buttonDialogClose),
										),
									)
								}),
							)
						})
					})
				}))
		})

}

func DuoUIinputField(
	rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme, f *Field) func() {
	rcs := rc.Settings
	rdw, rdCfg := rcs.Daemon.Widgets, rcs.Daemon.Config
	return func() {
		switch f.Field.Type {
		case "stringSlice":
			switch f.Field.InputType {
			case "text":
				if f.Field.Model != "MiningAddrs" {
					StringsArrayEditor(gtx, th,
						(rdw[f.Field.Model]).(*gel.Editor),
						(rdw[f.Field.Model]).(*gel.Editor).Text(),
						func(e gel.EditorEvent) {
							rdCfg[f.Field.Model] =
								strings.Fields((rdw[f.Field.Model]).(*gel.Editor).Text())
							Debug()
							if e != nil {
								rc.SaveDaemonCfg()
							}
						},
					)()
				}
			default:
			}
		case "input":
			switch f.Field.InputType {
			case "text":
				Editor(gtx, th,
					(rdw[f.Field.Model]).(*gel.Editor),
					(rdw[f.Field.Model]).(*gel.Editor).Text(),
					func(e gel.EditorEvent) {
						txt := rdw[f.Field.Model].(*gel.Editor).Text()
						rdCfg[f.Field.Model] = txt
						if e != nil {
							rc.SaveDaemonCfg()
						}
					},
				)()
			case "number":
				Editor(gtx, th,
					(rdw[f.Field.Model]).(*gel.Editor),
					(rdw[f.Field.Model]).(*gel.Editor).Text(),
					func(e gel.EditorEvent) {
						number, err :=
							strconv.Atoi((rdw[f.Field.Model]).(*gel.Editor).Text())
						if err == nil {
						}
						rdCfg[f.Field.Model] = number
						if e != nil {
							rc.SaveDaemonCfg()
						}
					},
				)()
			case "decimal":
				Editor(gtx, th,
					(rdw[f.Field.Model]).(*gel.Editor),
					(rdw[f.Field.Model]).(*gel.Editor).Text(),
					func(e gel.EditorEvent) {
						decimal, err :=
							strconv.ParseFloat(
								(rdw[f.Field.Model]).(*gel.Editor).Text(), 64)
						if err != nil {
						}
						rdCfg[f.Field.Model] = decimal
						if e != nil {
							rc.SaveDaemonCfg()
						}
					})()
			case "password":
				e := th.DuoUIeditor(f.Field.Label, "DocText", "DocBg", 32)
				e.Font.Typeface = th.Fonts["Primary"]
				e.Font.Style = text.Italic
				e.Layout(gtx, rdw[f.Field.Model].(*gel.Editor))
			default:
			}
		case "switch":
			sw := th.DuoUIcheckBox(
				f.Field.Label, th.Colors["Primary"], th.Colors["Primary"])
			sw.PillColor = th.Colors["LightGray"]
			sw.PillColorChecked = th.Colors["LightGrayI"]
			sw.CircleColor = th.Colors["LightGrayII"]
			sw.CircleColorChecked = th.Colors["Primary"]
			sw.DrawLayout(gtx, rdw[f.Field.Model].(*gel.CheckBox))
			if (rdw[f.Field.Model]).(*gel.CheckBox).Checked(gtx) {
				if !*rdCfg[f.Field.Model].(*bool) {
					tt := true
					rdCfg[f.Field.Model] = &tt
					rc.SaveDaemonCfg()
				}
			} else {
				if *rdCfg[f.Field.Model].(*bool) {
					ff := false
					rdCfg[f.Field.Model] = &ff
					rc.SaveDaemonCfg()
				}
			}
		case "radio":
			// radioButtonsGroup := (duo.Configuration.Settings.Daemon.Widgets[fieldName]).(*widget.Enum)
			// layout.Flex{}.Layout(gtx,
			//	layout.Rigid(func() {
			//		duo.Theme.RadioButton("r1", "RadioButton1").Layout(gtx,  radioButtonsGroup)
			//
			//	}),
			//	layout.Rigid(func() {
			//		duo.Theme.RadioButton("r2", "RadioButton2").Layout(gtx, radioButtonsGroup)
			//
			//	}),
			//	layout.Rigid(func() {
			//		duo.Theme.RadioButton("r3", "RadioButton3").Layout(gtx, radioButtonsGroup)
			//
			//	}))
		default:
			// duo.Theme.CheckBox("Checkbox").Layout(gtx, (duo.Configuration.Settings.Daemon.Widgets[fieldName]).(*widget.CheckBox))
		}
	}
}

func DuoUIlatestTransactions(
	rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	return func() {
		width := gtx.Constraints.Width.Max
		//cs := gtx.Constraints
		th.DuoUIcontainer(0, th.Colors["DarkGray"]).
			Layout(gtx, layout.NW, func() {
				layout.Flex{
					Axis: layout.Vertical,
				}.Layout(gtx,
					layout.Rigid(func() {
						th.DuoUIcontainer(8, th.Colors["Primary"]).
							Layout(gtx, layout.N,
								func() {
									gtx.Constraints.Width.Min = width
									latestx := th.H5("LATEST TRANSACTIONS")
									latestx.Color = th.Colors["Light"]
									latestx.Alignment = text.Start
									latestx.Layout(gtx)
								},
							)
					}),
					layout.Flexed(1, func() {
						layout.UniformInset(unit.Dp(8)).Layout(gtx,
							func() {
								layout.Flex{
									Axis: layout.Vertical,
								}.Layout(gtx,
									layout.Rigid(func() {
										latestTxsBookPanel := th.DuoUIpanel()
										latestTxsBookPanel.PanelObject =
											rc.Status.Wallet.LastTxs.Txs
										latestTxsBookPanel.ScrollBar =
											th.ScrollBar(0)
										latestTxsPanelElement.PanelObjectsNumber =
											len(rc.Status.Wallet.LastTxs.Txs)
										latestTxsBookPanel.Layout(gtx,
											latestTxsPanelElement, func(i int, in interface{}) {
												txs := in.([]model.DuoUItransactionExcerpt)
												t := txs[i]
												th.DuoUIcontainer(16, th.Colors["Dark"]).
													Layout(gtx, layout.NW, func() {
														width := gtx.Constraints.Width.Max
														layout.Flex{
															Axis: layout.
																Vertical,
														}.Layout(gtx,
															layout.Rigid(lTtxid(gtx, th, t.TxID)),
															layout.Rigid(func() {
																gtx.Constraints.Width.Min = width
																layout.Flex{
																	Spacing: layout.SpaceBetween,
																}.Layout(gtx,
																	layout.Rigid(func() {
																		layout.Flex{
																			Axis: layout.Vertical,
																		}.Layout(gtx,
																			layout.Rigid(lTcategory(gtx, th, t.Category)),
																			layout.Rigid(lTtime(gtx, th, t.Time)),
																		)
																	}),
																	layout.Rigid(lTamount(gtx, th, t.Amount)),
																)
															}),
															layout.Rigid(th.DuoUIline(gtx,
																0, 0, 1,
																th.Colors["Hint"]),
															),
														)
													})
											})
									}))
							})
					}),
				)
			})
	}
}

func DuoUIlogger(
	rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	return func() {
		// const buflen = 9
		layout.UniformInset(unit.Dp(10)).Layout(gtx, func() {
			// const n = 1e6
			cs := gtx.Constraints
			gelook.DuoUIdrawRectangle(gtx,
				cs.Width.Max, cs.Height.Max, th.Colors["Dark"],
				[4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
			lm := rc.Log.LogMessages.Load().([]log.Entry)
			logOutputList.Layout(gtx, len(lm), func(i int) {
				t := lm[i]
				logText := th.Caption(
					fmt.Sprintf("%-12s",
						t.Time.Sub(StartupTime)/time.Second*time.Second) +
						" " + fmt.Sprint(t.Text))
				logText.Font.Typeface = th.Fonts["Mono"]
				logText.Color = th.Colors["Primary"]
				switch t.Level {
				case "TRC":
					logText.Color = th.Colors["Success"]
				case "DBG":
					logText.Color = th.Colors["Secondary"]
				case "INF":
					logText.Color = th.Colors["Info"]
				case "WRN":
					logText.Color = th.Colors["Warning"]
				case "ERROR":
					logText.Color = th.Colors["Danger"]
				case "FTL":
					logText.Color = th.Colors["Primary"]
				}
				logText.Layout(gtx)
			})
			op.InvalidateOp{}.Add(gtx.Ops)
		})
	}
}

func DuoUIstatus(
	rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	wall := rc.Status.Wallet
	nod := rc.Status.Node
	return func() {
		th.DuoUIcontainer(8, th.Colors["Light"]).Layout(gtx, layout.NW, func() {
			bigStatus := []func(){
				listItem(gtx, th,
					22, 6, "EditorMonetizationOn",
					"BALANCE :", wall.Balance.Load()+" "+rc.Settings.Abbrevation),
				th.DuoUIline(gtx,
					8, 0, 1, th.Colors["LightGray"]),
				listItem(gtx, th,
					22, 6, "MapsLayersClear", "UNCONFIRMED :",
					wall.Unconfirmed.Load()+" "+rc.Settings.Abbrevation),
				th.DuoUIline(gtx,
					8, 0, 1, th.Colors["LightGray"]),
				listItem(gtx, th,
					22, 6, "CommunicationImportExport",
					"TRANSACTIONS :", fmt.Sprint(wall.TxsNumber.Load())),
				th.DuoUIline(gtx,
					8, 0, 1, th.Colors["LightGray"]),
				listItem(gtx, th,
					16, 4, "DeviceWidgets",
					"Block Count :", fmt.Sprint(nod.BlockCount.Load())),
				th.DuoUIline(gtx,
					4, 0, 1, th.Colors["LightGray"]),
				listItem(gtx, th,
					16, 4, "ImageTimer",
					"Difficulty :", fmt.Sprint(nod.Difficulty.Load())),
				th.DuoUIline(gtx,
					4, 0, 1, th.Colors["LightGray"]),
				listItem(gtx, th,
					16, 4, "NotificationVPNLock",
					"Connections :", fmt.Sprint(nod.ConnectionCount.Load())),
			}
			itemsList.Layout(gtx, len(bigStatus), func(i int) {
				layout.UniformInset(unit.Dp(0)).Layout(gtx, bigStatus[i])
			})
		})
	}
}

func FooterLeftMenu(
	rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme, allPages *model.DuoUIpages) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
			cornerButtons := []func(){
				QuitButton(rc, gtx, th),
				// footerMenuButton(rc, gtx, th, allPages.Theme["EXPLORER"],
				//"BLOCKS: "+fmt.Sprint(rc.Status.Node.BlockCount), "", buttonBlocks),
				footerMenuButton(rc, gtx, th, allPages.Theme["LOG"], "LOG",
					"traceIcon", buttonLog),
			}
			cornerNav.Layout(gtx, len(cornerButtons), func(i int) {
				layout.UniformInset(unit.Dp(0)).Layout(gtx, cornerButtons[i])
			})
		})
	}
}

func footerMenuButton(
	rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme, page *gelook.DuoUIpage, text, icon string, footerButton *gel.Button) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
			var footerMenuItem gelook.DuoUIbutton
			if icon != "" {
				footerMenuItem = th.DuoUIbutton("", "",
					"", "", "", th.Colors["Dark"], icon,
					CurrentCurrentPageColor(rc.ShowPage,
						page.Title, navItemIconColor, th.Colors["Primary"]),
					footerMenuItemTextSize, footerMenuItemIconSize,
					footerMenuItemWidth, footerMenuItemHeight,
					0, 0, 0, 0)
				for footerButton.Clicked(gtx) {
					rc.ShowPage = page.Title
					SetPage(rc, page)
				}
			} else {
				footerMenuItem = th.DuoUIbutton(th.Fonts["Primary"], text,
					CurrentCurrentPageColor(rc.ShowPage, page.Title,
						th.Colors["Light"], th.Colors["Primary"]),
					"", "", "", "", "",
					footerMenuItemTextSize, footerMenuItemIconSize,
					0, footerMenuItemHeight,
					13, 16, 14, 16)
				footerMenuItem.Height = 48
				for footerButton.Clicked(gtx) {
					rc.ShowPage = page.Title
					SetPage(rc, page)
				}
			}
			footerMenuItem.Layout(gtx, footerButton)
		})
	}
}

func FooterRightMenu(
	rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme, allPages *model.DuoUIpages) func() {
	return func() {
		navButtons := []func(){
			// footerMenuButton(rc, gtx, th, allPages.Theme["NETWORK"],
			// "", "networkIcon", buttonNetwork),
			// footerMenuButton(rc, gtx, th, allPages.Theme["NETWORK"],
			// "CONNECTIONS: "+fmt.Sprint(rc.Status.Node.ConnectionCount),
			//"", buttonNetwork),
			// footerMenuButton(rc, gtx, th, allPages.Theme["EXPLORER"],
			// "", "DeviceWidgets", buttonBlocks),
			footerMenuButton(rc, gtx, th, allPages.Theme["NETWORK"],
				"", "network", buttonNetwork),
			footerMenuButton(rc, gtx, th, allPages.Theme["NETWORK"],
				"CONNECTIONS: "+fmt.Sprint(rc.Status.Node.ConnectionCount.Load()), "", buttonNetwork),
			footerMenuButton(rc, gtx, th, allPages.Theme["EXPLORER"],
				"", "DeviceWidgets", buttonBlocks),
			footerMenuButton(rc, gtx, th, allPages.Theme["EXPLORER"],
				"BLOCKS: "+fmt.Sprint(rc.Status.Node.BlockCount.Load()), "", buttonBlocks),
			footerMenuButton(rc, gtx, th, allPages.Theme["MINER"],
				"", "helpIcon", buttonHelp),
			footerMenuButton(rc, gtx, th, allPages.Theme["CONSOLE"],
				"", "consoleIcon", buttonConsole),
			footerMenuButton(rc, gtx, th, allPages.Theme["SETTINGS"],
				"", "settingsIcon", buttonSettings),
		}
		footerNav.Layout(gtx, len(navButtons), func(i int) {
			layout.UniformInset(unit.Dp(0)).Layout(gtx, navButtons[i])
		})
	}
}

func HeaderMenu(
	rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme, allPages *model.DuoUIpages) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
			headerNav := []func(){
				headerMenuButton(rc, gtx, th, "",
					"CommunicationImportExport", buttonHeader),
				headerMenuButton(rc, gtx, th, "",
					"NotificationNetworkCheck", buttonHeader),
				headerMenuButton(rc, gtx, th, "",
					"NotificationSync", buttonHeader),
				headerMenuButton(rc, gtx, th, "",
					"NotificationSyncDisabled", buttonHeader),
				headerMenuButton(rc, gtx, th, "",
					"NotificationSyncProblem", buttonHeader),
				headerMenuButton(rc, gtx, th, "",
					"NotificationVPNLock", buttonHeader),
				headerMenuButton(rc, gtx, th, "",
					"MapsLayers", buttonHeader),
				headerMenuButton(rc, gtx, th, "",
					"MapsLayersClear", buttonHeader),
				headerMenuButton(rc, gtx, th, "",
					"ImageTimer", buttonHeader),
				headerMenuButton(rc, gtx, th, "",
					"ImageRemoveRedEye", buttonHeader),
				headerMenuButton(rc, gtx, th, "",
					"DeviceSignalCellular0Bar", buttonHeader),
				headerMenuButton(rc, gtx, th, "",
					"ActionTimeline", buttonHeader),
				headerMenuButton(rc, gtx, th, "",
					"HardwareWatch", buttonHeader),
			}
			footerNav.Layout(gtx, len(headerNav), func(i int) {
				layout.UniformInset(unit.Dp(0)).Layout(gtx, headerNav[i])
			})
		})
	}
}

func headerMenuButton(
	rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme, text, icon string, headerButton *gel.Button) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
			var footerMenuItem gelook.DuoUIbutton
			footerMenuItem = th.DuoUIbutton("", "",
				"", "", "",
				th.Colors["Dark"], icon,
				CurrentCurrentPageColor(rc.ShowPage, text, navItemIconColor,
					th.Colors["Primary"]),
				footerMenuItemTextSize, footerMenuItemIconSize,
				footerMenuItemWidth, footerMenuItemHeight,
				footerMenuItemPaddingVertical, footerMenuItemPaddingHorizontal,
				footerMenuItemPaddingVertical, footerMenuItemPaddingHorizontal)
			for headerButton.Clicked(gtx) {
				rc.ShowPage = text
			}
			footerMenuItem.IconLayout(gtx, headerButton)
		})
	}
}

func iconButton(
	rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme, page *gelook.DuoUIpage) func() {
	return func() {
		var logMenuItem gelook.DuoUIbutton
		logMenuItem = th.DuoUIbutton("", "", "",
			th.Colors["Dark"], "", "", "traceIcon",
			CurrentCurrentPageColor(rc.ShowPage, "LOG",
				th.Colors["Light"], th.Colors["Primary"]),
			footerMenuItemTextSize, footerMenuItemIconSize,
			footerMenuItemWidth, footerMenuItemHeight,
			footerMenuItemPaddingVertical, footerMenuItemPaddingHorizontal,
			footerMenuItemPaddingVertical, footerMenuItemPaddingHorizontal)
		for buttonLog.Clicked(gtx) {
			SetPage(rc, page)
			rc.ShowPage = "LOG"
		}
		logMenuItem.IconLayout(gtx, buttonLog)
	}
}

func MainNavigation(
	rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme, allPages *model.DuoUIpages, nav *model.DuoUInav) func() {
	return func() {
		navButtons := navButtons(rc, gtx, th, allPages, nav)
		gtx.Constraints.Width.Max = nav.Width
		th.DuoUIcontainer(0, th.Colors["Dark"]).Layout(gtx, layout.NW, func() {
			mainNav.Layout(gtx, len(navButtons), func(i int) {
				layout.UniformInset(unit.Dp(0)).Layout(gtx, navButtons[i])
			})
		})
	}
}

func navButtons(
	rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme, allPages *model.DuoUIpages, nav *model.DuoUInav) []func() {
	return []func(){
		navMenuButton(rc, gtx, th, allPages.Theme["OVERVIEW"], nav,
			"OVERVIEW", "overviewIcon", navButtonOverview),
		th.DuoUIline(gtx, 0, 0, 1,
			th.Colors["LightGrayIII"]),
		navMenuButton(rc, gtx, th, allPages.Theme["SEND"], nav,
			"SEND", "sendIcon", navButtonSend),
		// navMenuLine(gtx, th),
		// navMenuButton(rc, gtx, th, allPages.Theme["RECEIVE"], "RECEIVE",
		//"receiveIcon", navButtonReceive),
		th.DuoUIline(gtx, 0, 0, 1,
			th.Colors["LightGrayIII"]),
		navMenuButton(rc, gtx, th, allPages.Theme["ADDRESSBOOK"], nav,
			"ADDRESSBOOK", "addressBookIcon", navButtonAddressBook),
		th.DuoUIline(gtx, 0, 0, 1,
			th.Colors["LightGrayIII"]),
		navMenuButton(rc, gtx, th, allPages.Theme["HISTORY"], nav,
			"HISTORY", "historyIcon", navButtonHistory),
	}
}

func navMenuButton(
	rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme, page *gelook.DuoUIpage, nav *model.DuoUInav, title, icon string, navButton *gel.Button) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
			var menuItem gelook.DuoUIbutton
			menuItem = th.DuoUIbutton(th.Fonts["Secondary"],
				title, th.Colors["Dark"],
				th.Colors["LightGrayII"],
				th.Colors["LightGrayII"],
				th.Colors["Dark"], icon,
				CurrentCurrentPageColor(rc.ShowPage, title, navItemIconColor,
					th.Colors["Primary"]), nav.TextSize, nav.IconSize,
				nav.Width, nav.Height, nav.PaddingVertical,
				nav.PaddingHorizontal, nav.PaddingVertical, nav.PaddingHorizontal)
			for navButton.Clicked(gtx) {
				rc.ShowPage = title
				SetPage(rc, page)
			}
			menuItem.MenuLayout(gtx, navButton)
		})
	}
}

func pageNavButton(
	rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme, page *gelook.DuoUIpage, b *gel.Button, label, hash string) {
	layout.UniformInset(unit.Dp(4)).Layout(gtx, func() {
		var blockButton gelook.DuoUIbutton
		blockButton = th.DuoUIbutton(th.Fonts["Mono"], label+" "+hash,
			th.Colors["Light"], th.Colors["Info"],
			th.Colors["Info"], th.Colors["Light"], "", th.Colors["Light"],
			16, 0, 60, 24,
			0, 0, 0, 0)
		for b.Clicked(gtx) {
			rc.ShowPage = fmt.Sprintf("BLOCK %s", hash)
			rc.GetSingleBlock(hash)()
			SetPage(rc, page)
		}
		blockButton.Layout(gtx, b)
	})
}

func PageNavButtons(
	rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme, previousBlockHash, nextBlockHash string, prevPage, nextPage *gelook.DuoUIpage) func() {
	return func() {
		layout.Flex{}.Layout(gtx,
			layout.Flexed(0.5, func() {
				eh := chainhash.Hash{}
				if previousBlockHash != eh.String() {
					pageNavButton(rc, gtx, th, nextPage,
						previousBlockHashButton, "Previous Block",
						previousBlockHash)
				}
			}),
			layout.Flexed(0.5, func() {
				if nextBlockHash != "" {
					pageNavButton(rc, gtx,
						th, nextPage, nextBlockHashButton,
						"Next Block", nextBlockHash)
				}
			}))
	}
}

func PeersList(
	rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	c := th.Colors["Dark"]
	return func() {
		rc.Network.PeersList.Layout(gtx, len(rc.Network.Peers), func(i int) {
			t := rc.Network.Peers[i]
			th.DuoUIline(gtx, 0, 0, 1,
				th.Colors["Hint"])()
			layout.Flex{
				Spacing: layout.SpaceBetween,
			}.Layout(gtx,
				layout.Rigid(Label(gtx, th, th.Fonts["Mono"], 14,
					c, fmt.Sprint(t.ID))),
				layout.Flexed(1, peerDetails(gtx, th, i, t)),
				layout.Rigid(Label(gtx, th, th.Fonts["Mono"], 14,
					c, t.Addr)))
		})
	}
}

func QuitButton(
	rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
			var closeMeniItem gelook.DuoUIbutton
			closeMeniItem = th.DuoUIbutton("", "", "",
				th.Colors["Dark"], "", "",
				"closeIcon",
				CurrentCurrentPageColor(
					rc.ShowPage, "CLOSE", th.Colors["Light"],
					th.Colors["Primary"]),
				footerMenuItemTextSize, footerMenuItemIconSize,
				footerMenuItemWidth, footerMenuItemHeight,
				0, 0, 0, 0)
			for buttonQuit.Clicked(gtx) {
				rc.Dialog.Show = true
				rc.Dialog = &model.DuoUIdialog{
					Show: true,
					Green: func() {
						interrupt.Request()
						// TODO make this close the window or at least switch to a shutdown screen
						rc.Dialog.Show = false
					},
					GreenLabel: "QUIT",
					Orange: func() {
						interrupt.RequestRestart()
						// TODO make this close the window or at least switch to a shutdown screen
						rc.Dialog.Show = false
					},
					OrangeLabel: "RESTART",
					Red:         func() { rc.Dialog.Show = false },
					RedLabel:    "CANCEL",
					CustomField: func() {},
					Title:       "Are you sure?",
					Text:        "Confirm ParallelCoin close",
				}
			}
			closeMeniItem.IconLayout(gtx, buttonQuit)
		})
	}
}

func SettingsTabs(
	rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	rcs := rc.Settings
	return func() {
		groupsNumber := len(rcs.Daemon.Schema.Groups)
		groupsList.Layout(gtx, groupsNumber, func(i int) {
			layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
				color := th.Colors["Light"]
				bgColor := th.Colors["Dark"]
				i = groupsNumber - 1 - i
				t := rcs.Daemon.Schema.Groups[i]
				txt := fmt.Sprint(t.Legend)
				for rcs.Tabs.TabsList[txt].Clicked(gtx) {
					rcs.Tabs.Current = txt
				}
				if rcs.Tabs.Current == txt {
					color = th.Colors["Dark"]
					bgColor = th.Colors["Light"]
				}
				th.DuoUIbutton(th.Fonts["Primary"],
					txt, color, bgColor, "", "",
					"", "", 16, 0, 80, 32,
					4, 4, 4, 4,
				).Layout(gtx, rcs.Tabs.TabsList[txt])
			})
		})
	}
}

func TransactionsFilter(
	rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	cats := rc.History.Categories
	return func() {
		layout.Flex{}.Layout(gtx,
			layout.Rigid(txsFilterItem(gtx, th, "ALL", cats.AllTxs)),
			layout.Rigid(txsFilterItem(gtx, th, "MINTED", cats.MintedTxs)),
			layout.Rigid(txsFilterItem(gtx, th, "IMATURE", cats.ImmatureTxs)),
			layout.Rigid(txsFilterItem(gtx, th, "SENT", cats.SentTxs)),
			layout.Rigid(txsFilterItem(gtx, th, "RECEIVED", cats.ReceivedTxs)))
		switch {
		case cats.AllTxs.Checked(gtx):
			rc.History.Category = "all"
		case cats.MintedTxs.Checked(gtx):
			rc.History.Category = "generate"
		case cats.ImmatureTxs.Checked(gtx):
			rc.History.Category = "immature"
		case cats.SentTxs.Checked(gtx):
			rc.History.Category = "sent"
		case cats.ReceivedTxs.Checked(gtx):
			rc.History.Category = "received"
		}
	}
}

// gtx, th

func Button(
	gtx *layout.Context, th *gelook.DuoUItheme, buttonController *gel.Button, font text.Typeface, textSize, padT, padR, padB, padL int, color, bgColor, label string, handler func()) func() {
	return func() {
		th.DuoUIcontainer(0, th.Colors[""])
		button := th.DuoUIbutton(font, label, color, bgColor,
			"", "", "", "",
			textSize, 0, 128, 48,
			padT, padR, padB, padL)
		for buttonController.Clicked(gtx) {
			handler()
		}
		button.Layout(gtx, buttonController)
	}
}

func ConsoleInput(
	gtx *layout.Context, th *gelook.DuoUItheme, editorController *gel.Editor, label string, handler func(gel.SubmitEvent)) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
			e := th.DuoUIeditor(label, "Dark", "Light", 120)
			e.Font.Typeface = th.Fonts["Primary"]
			e.Color = gelook.HexARGB(th.Colors["Light"])
			e.Font.Style = text.Italic
			e.Layout(gtx, editorController)
			for _, e := range editorController.Events(gtx) {
				if e, ok := e.(gel.SubmitEvent); ok {
					handler(e)
					editorController.SetText("")
				}
			}
		})
	}
}

func contentField(
	gtx *layout.Context, th *gelook.DuoUItheme, text, color, bgColor string, font text.Typeface, padding, textSize float32) func() {
	return func() {
		hmin := gtx.Constraints.Width.Min
		vmin := gtx.Constraints.Height.Min
		layout.Stack{Alignment: layout.W}.Layout(gtx,
			layout.Expanded(func() {
				rr := float32(gtx.Px(unit.Dp(0)))
				clip.Rect{
					Rect: f32.Rectangle{Max: f32.Point{
						X: float32(gtx.Constraints.Width.Min),
						Y: float32(gtx.Constraints.Height.Min),
					}},
					NE: rr, NW: rr, SE: rr, SW: rr,
				}.Op(gtx.Ops).Add(gtx.Ops)
				fill(gtx, gelook.HexARGB(bgColor))
			}),
			layout.Stacked(func() {
				gtx.Constraints.Width.Min = hmin
				gtx.Constraints.Height.Min = vmin
				layout.Center.Layout(gtx, func() {
					layout.UniformInset(unit.Dp(padding)).Layout(gtx, func() {
						l := th.DuoUIlabel(unit.Dp(textSize), text)
						l.Font.Typeface = font
						l.Color = color
						l.Layout(gtx)
					})
				})
			}),
		)
	}
}

func ContentLabeledField(
	gtx *layout.Context, th *gelook.DuoUItheme, axis layout.Axis, margin, labelTextSize, valueTextSize float32, label, headcolor, headbgColor, color, bgColor, value string) func() {
	return func() {
		layout.UniformInset(unit.Dp(margin)).Layout(gtx, func() {
			layout.Flex{
				Axis: axis,
			}.Layout(gtx,
				layout.Rigid(contentField(gtx, th, label,
					th.Colors[headcolor], th.Colors[headbgColor],
					th.Fonts["Primary"], 4, labelTextSize),
				),
				layout.Rigid(contentField(gtx, th, value,
					th.Colors[color], th.Colors[bgColor],
					th.Fonts["Mono"], 4, valueTextSize),
				),
			)
		})
	}
}

func dialogButton(
	gtx *layout.Context, th *gelook.DuoUItheme, f func(), t, bgColor, icon, iconColor string, button *gel.Button) func() {
	return func() {
		if f != nil {
			var b gelook.DuoUIbutton
			layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8),
				Left: unit.Dp(8), Right: unit.Dp(8)}.Layout(gtx, func() {
				b = th.DuoUIbutton(th.Fonts["Primary"], t,
					th.Colors["Dark"], bgColor, th.Colors["Info"], bgColor,
					icon, iconColor, 16, 32,
					120, 64,
					0, 0, 0, 0)
				for button.Clicked(gtx) {
					f()
				}
				b.MenuLayout(gtx, button)
			})
		}
	}
}

func Editor(
	gtx *layout.Context, th *gelook.DuoUItheme, editorController *gel.Editor, label string, handler func(gel.EditorEvent)) func() {
	return func() {
		th.DuoUIcontainer(8, "ffffffff").
			Layout(gtx, layout.NW, func() {
				width := gtx.Constraints.Width.Max
				e := th.DuoUIeditor(label, th.Colors["Black"], th.Colors["White"], width)
				e.Font.Typeface = th.Fonts["Mono"]
				e.TextSize = unit.Dp(12)
				layout.UniformInset(unit.Dp(4)).Layout(gtx, func() {
					e.Layout(gtx, editorController)
				})
				for _, e := range editorController.Events(gtx) {
					switch e.(type) {
					case gel.ChangeEvent:
						handler(e)
					}
				}
			})
	}
}

func Label(
	gtx *layout.Context, th *gelook.DuoUItheme, font text.Typeface, size float32, color, label string) func() {
	return func() {
		l := th.DuoUIlabel(unit.Dp(size), label)
		l.Font.Typeface = font
		l.Color = color
		l.Layout(gtx)
	}
}

func listItem(
	gtx *layout.Context, th *gelook.DuoUItheme, size, top int, iconName, name, value string) func() {
	return func() {
		icon := th.Icons[iconName]
		layout.Flex{
			Axis:    layout.Horizontal,
			Spacing: layout.SpaceBetween,
		}.Layout(gtx,
			layout.Rigid(func() {
				layout.Flex{}.Layout(gtx,
					layout.Rigid(func() {
						layout.Inset{
							Top:    unit.Dp(float32(top)),
							Bottom: unit.Dp(0), Left: unit.Dp(0),
							Right: unit.Dp(0)}.Layout(gtx, func() {
							if icon != nil {
								icon.Color = gelook.HexARGB(th.Colors["Dark"])
								icon.Layout(gtx, unit.Px(float32(size)))
							}
							gtx.Dimensions = layout.Dimensions{
								Size: image.Point{X: size, Y: size},
							}
						})
					}),
					layout.Rigid(func() {
						txt := th.DuoUIlabel(unit.Dp(float32(size)), name)
						txt.Font.Typeface = th.Fonts["Primary"]
						txt.Color = th.Colors["Primary"]
						txt.Layout(gtx)
					}),
				)
			}),
			layout.Rigid(func() {
				v := th.H5(value)
				v.TextSize = unit.Dp(float32(size))
				v.Font.Typeface = th.Fonts["Primary"]
				v.Color = th.Colors["Dark"]
				v.Alignment = text.End
				v.Layout(gtx)
			}),
		)
	}
}

func lTamount(
	gtx *layout.Context, th *gelook.DuoUItheme, v float64) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
			sat := th.Body1(fmt.Sprintf("%0.8f", v))
			sat.Font.Typeface = "bariol"
			sat.Color = th.Colors["Light"]
			sat.Layout(gtx)
		})
	}
}

func lTcategory(
	gtx *layout.Context, th *gelook.DuoUItheme, v string) func() {
	return func() {
		sat := th.Body1(v)
		sat.Color = th.Colors["Light"]
		sat.Font.Typeface = "bariol"
		sat.Layout(gtx)
	}
}

func lTtime(
	gtx *layout.Context, th *gelook.DuoUItheme, v string) func() {
	return func() {
		l := th.Body1(v)
		l.Font.Typeface = "bariol"
		l.Color = th.Colors["Light"]
		l.Color = th.Colors["Hint"]
		l.Layout(gtx)
	}
}

func lTtxid(
	gtx *layout.Context, th *gelook.DuoUItheme, v string) func() {
	return func() {
		tim := th.Caption(v)
		tim.Font.Typeface = th.Fonts["Mono"]
		tim.Color = th.Colors["Light"]
		tim.Layout(gtx)
	}
}

func MonoButton(
	gtx *layout.Context, th *gelook.DuoUItheme, buttonController *gel.Button, textSize int, color, bgColor, font, label string, handler func()) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
			//var button gelook.Button
			button := th.Button(label)
			switch {
			case font != "":
				button.Font.Typeface = th.Fonts[font]
			case color != "":
				button.Color = gelook.HexARGB(th.Colors[color])
			case textSize != 0:
				button.TextSize = unit.Dp(float32(textSize))
			case bgColor != "":
				button.Background = gelook.HexARGB(th.Colors[bgColor])
			}
			for buttonController.Clicked(gtx) {
				handler()
			}
			button.Layout(gtx, buttonController)
		})
	}
}

func peerDetails(
	gtx *layout.Context, th *gelook.DuoUItheme, i int, t *btcjson.GetPeerInfoResult) func() {
	prim := th.Fonts["Primary"]
	c := th.Colors["Dark"]
	return func() {
		layout.Flex{
			Axis:    layout.Horizontal,
			Spacing: layout.SpaceAround,
		}.Layout(gtx,
			layout.Rigid(Label(gtx, th, prim, 12, c, t.AddrLocal)),
			layout.Rigid(Label(gtx, th, prim, 12, c, t.Services)),
			layout.Rigid(Label(gtx, th, prim, 12, c, fmt.Sprint(t.RelayTxes))),
			// layout.Rigid(Label(gtx, th, prim, 12, c
			//fmt.Sprint(t.LastSend))),
			// layout.Rigid(Label(gtx, th, prim, 12, c
			//fmt.Sprint(t.LastRecv))),
			layout.Rigid(Label(gtx, th, prim, 12, c, fmt.Sprint(t.BytesSent))),
			layout.Rigid(Label(gtx, th, prim, 12, c,
				fmt.Sprint(t.BytesRecv))),
			layout.Rigid(Label(gtx, th, prim, 12, c, fmt.Sprint(t.ConnTime))),
			layout.Rigid(Label(gtx, th, prim, 12, c, fmt.Sprint(t.TimeOffset))),
			layout.Rigid(Label(gtx, th, prim, 12, c, fmt.Sprint(t.PingTime))),
			layout.Rigid(Label(gtx, th, prim, 12, c, fmt.Sprint(t.PingWait))),
			layout.Rigid(Label(gtx, th, prim, 12, c, fmt.Sprint(t.Version))),
			layout.Rigid(Label(gtx, th, prim, 12, c, t.SubVer)),
			layout.Rigid(Label(gtx, th, prim, 12, c, fmt.Sprint(t.Inbound))),
			layout.Rigid(Label(gtx, th, prim, 12, c, fmt.Sprint(t.StartingHeight))),
			layout.Rigid(Label(gtx, th, prim, 12, c, fmt.Sprint(t.CurrentHeight))),
			layout.Rigid(Label(gtx, th, prim, 12, c, fmt.Sprint(t.BanScore))),
			layout.Rigid(Label(gtx, th, prim, 12, c, fmt.Sprint(t.FeeFilter))),
			layout.Rigid(Label(gtx, th, prim, 12, c, fmt.Sprint(t.SyncNode))))
	}
}

func SettingsFieldDescription(
	gtx *layout.Context, th *gelook.DuoUItheme, f *Field) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
			desc := th.Body2(fmt.Sprint(f.Field.Description))
			desc.Font.Typeface = th.Fonts["Primary"]
			desc.Color = th.Colors["Dark"]
			desc.Layout(gtx)
		})
	}
}

func SettingsFieldLabel(
	gtx *layout.Context, th *gelook.DuoUItheme, f *Field) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
			name := th.H6(fmt.Sprint(f.Field.Label))
			name.Color = th.Colors["Dark"]
			name.Font.Typeface = th.Fonts["Primary"]
			name.Layout(gtx)
		})
	}
}

func StringsArrayEditor(
	gtx *layout.Context, th *gelook.DuoUItheme, editorController *gel.Editor, label string, handler func(gel.EditorEvent)) func() {
	return func() {
		th.DuoUIcontainer(8, th.Colors["White"]).Layout(gtx, layout.NW, func() {
			e := th.DuoUIeditor(label, th.Colors["Black"], th.Colors["White"], 16)
			e.Font.Typeface = th.Fonts["Mono"]
			layout.UniformInset(unit.Dp(4)).Layout(gtx, func() {
				e.Layout(gtx, editorController)
			})
			for _, e := range editorController.Events(gtx) {
				switch e.(type) {
				case gel.ChangeEvent:
					handler(e)
				}
			}
		})
	}
}

func TrioFields(
	gtx *layout.Context, th *gelook.DuoUItheme, axis layout.Axis, labelTextSize, valueTextSize float32, unoLabel, unoValue, unoHeadcolor, unoHeadbgColor, unoColor, unoBgColor, duoLabel, duoValue, duoHeadcolor, duoHeadbgColor, duoColor, duoBgColor, treLabel, treValue, treHeadcolor, treHeadbgColor, treColor, treBgColor string) func() {
	return func() {
		layout.Flex{
			Axis:    axis,
			Spacing: layout.SpaceAround,
		}.Layout(gtx,
			fieldAxis(axis, ContentLabeledField(gtx, th,
				layout.Vertical, 4,
				labelTextSize, valueTextSize,
				unoLabel, unoHeadcolor, unoHeadbgColor,
				unoColor, unoBgColor,
				fmt.Sprint(unoValue)), 0.3333),
			fieldAxis(axis, ContentLabeledField(gtx, th,
				layout.Vertical, 4,
				labelTextSize, valueTextSize,
				duoLabel, duoHeadcolor, duoHeadbgColor,
				duoColor, duoBgColor,
				fmt.Sprint(duoValue)), 0.3333),
			fieldAxis(axis, ContentLabeledField(gtx, th,
				layout.Vertical, 4,
				labelTextSize, valueTextSize,
				treLabel, treHeadbgColor, treHeadcolor,
				treColor, treBgColor,
				fmt.Sprint(treValue)), 0.3333),
		)
	}
}

func TxsDetails(
	gtx *layout.Context, th *gelook.DuoUItheme, i int, t *model.DuoUItransactionExcerpt) func() {
	return func() {
		layout.Flex{
			Axis: layout.Vertical,
		}.Layout(gtx,
			layout.Rigid(Label(gtx, th, th.Fonts["Primary"], 12,
				th.Colors["Dark"], fmt.Sprint(i))),
			layout.Rigid(Label(gtx, th, th.Fonts["Primary"], 12,
				th.Colors["Dark"], t.TxID)),
			layout.Rigid(Label(gtx, th, th.Fonts["Primary"], 12,
				th.Colors["Dark"], fmt.Sprintf("%0.8f", t.Amount))),
			layout.Rigid(Label(gtx, th, th.Fonts["Primary"], 12,
				th.Colors["Dark"], t.Category)),
			layout.Rigid(Label(gtx, th, th.Fonts["Primary"], 12,
				th.Colors["Dark"], t.Time)),
		)
	}
}

func txsFilterItem(
	gtx *layout.Context, th *gelook.DuoUItheme, id string, c *gel.CheckBox) func() {
	return func() {
		th.DuoUIcheckBox(id, th.Colors["Light"], th.Colors["Light"]).Layout(gtx, c)
	}
}

// misc

func CurrentCurrentPageColor(showPage, page, color, currentPageColor string) (c string) {
	if showPage == page {
		c = currentPageColor
	} else {
		c = color
	}
	return
}

func DuoFields(gtx *layout.Context, axis layout.Axis, left, right func()) func() {
	return func() {
		layout.Flex{
			Axis:    axis,
			Spacing: layout.SpaceAround,
		}.Layout(gtx,
			fieldAxis(axis, left, 0.5),
			fieldAxis(axis, right, 0.5),
		)
	}
}

func DuoUIqrCode(gtx *layout.Context, hash string, size uint) func() {
	return func() {
		qr, err := qrcode.Encode(hash, 3, qrcode.ECLevelM)
		if err != nil {
		}
		qrResize := resize.Resize(size, 0, qr, resize.NearestNeighbor)
		addrQR := paint.NewImageOp(qrResize)
		sz := gtx.Constraints.Width.Constrain(gtx.Px(unit.Dp(float32(size))))
		addrQR.Add(gtx.Ops)
		paint.PaintOp{
			Rect: f32.Rectangle{
				Max: f32.Point{
					X: float32(sz), Y: float32(sz),
				},
			},
		}.Add(gtx.Ops)
		gtx.Dimensions.Size = image.Point{X: sz, Y: sz}
	}
}

func fieldAxis(axis layout.Axis, field func(), size float32) layout.FlexChild {
	var f layout.FlexChild
	switch axis {
	case layout.Horizontal:
		f = layout.Flexed(size, field)
	case layout.Vertical:
		f = layout.Rigid(field)
	}
	return f
}

func fill(gtx *layout.Context, col color.RGBA) {
	cs := gtx.Constraints
	d := image.Point{X: cs.Width.Min, Y: cs.Height.Min}
	dr := f32.Rectangle{
		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
	}
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{Rect: dr}.Add(gtx.Ops)
	gtx.Dimensions = layout.Dimensions{Size: d}
}

func QrDialog(rc *rcd.RcVar, gtx *layout.Context, address string) func() {
	return func() {
		// clipboard.Set(t.Address)
		rc.Dialog.Show = true
		rc.Dialog = &model.DuoUIdialog{
			Show:        true,
			CustomField: DuoUIqrCode(gtx, address, 256),
			Orange:      func() { rc.Dialog.Show = false },
			OrangeLabel: "CLOSE",
			Title:       "ParallelCoin address",
			Text:        address,
		}
	}
}

func SetPage(rc *rcd.RcVar, page *gelook.DuoUIpage) {
	page.Command()
	rc.CurrentPage = page
}

func UnoField(gtx *layout.Context, field func()) func() {
	return func() {
		layout.Flex{
			Axis:    layout.Horizontal,
			Spacing: layout.SpaceAround,
		}.Layout(gtx,
			layout.Flexed(1, field),
		)

	}
}
