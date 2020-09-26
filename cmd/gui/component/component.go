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
	"github.com/stalker-loki/app/slog"
	"image"
	"image/color"
	"strconv"
	"strings"
	"time"
)

type (
	Context struct {
		Gtx *layout.Context
		Thm *gelook.DuoUITheme
	}
	State struct {
		Context
		Rc *rcd.RcVar
	}
	DuoUIcomponent struct {
		Name    string
		Version string
		Theme   *gelook.DuoUITheme
		M       interface{}
		V       func()
		C       func()
	}
	Field struct {
		Field *pod.Field
	}
)

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

func NewContext(gtx *layout.Context, th *gelook.DuoUITheme) (cx *Context) {
	return &Context{
		Gtx: gtx,
		Thm: th,
	}
}

func NewState(rc *rcd.RcVar, gtx *layout.Context,
	th *gelook.DuoUITheme) (cx *State) {
	return &State{
		Context: *NewContext(gtx, th),
		Rc:      rc,
	}
}

// rc, gtx, th

func (s *State) DuoUIdialog() {
	t, g := s.Thm, s.Gtx
	// cs := g.Constraints
	t.DuoUIContainer(0, "ee000000").
		Layout(g, layout.Center, func() {
			cs := g.Constraints
			layout.Stack{
				Alignment: layout.Center,
			}.Layout(g,
				layout.Expanded(func() {
					rr := float32(g.Px(unit.Dp(0)))
					clip.Rect{
						Rect: f32.Rectangle{
							Max: f32.Point{
								X: float32(cs.Width.Max),
								Y: float32(cs.Height.Max),
							},
						},
						NE: rr, NW: rr, SE: rr, SW: rr,
					}.Op(g.Ops).Add(g.Ops)
					fill(g, gelook.HexARGB("ff888888"))
					pointer.Rect(image.Rectangle{
						Max: g.Dimensions.Size,
					}).Add(g.Ops)
				}),
				layout.Stacked(func() {
					if g.Constraints.Width.Max > 500 {
						g.Constraints.Width.Max = 500
					}
					layout.Center.Layout(g, func() {
						layout.Inset{
							Top: unit.Dp(16), Bottom: unit.Dp(16),
							Left: unit.Dp(8), Right: unit.Dp(8),
						}.Layout(g, func() {
							layout.Flex{
								Axis:      layout.Vertical,
								Alignment: layout.Middle,
							}.Layout(g,
								layout.Rigid(func() {
									layout.Flex{
										Axis:      layout.Horizontal,
										Alignment: layout.Middle,
									}.Layout(g, layout.Rigid(func() {
										layout.Inset{
											Top:    unit.Dp(0),
											Bottom: unit.Dp(8),
											Left:   unit.Dp(4),
											Right:  unit.Dp(4)}.
											Layout(g, func() {
												cur := t.DuoUILabel(unit.Dp(
													14), s.Rc.Dialog.Text)
												cur.Font.Typeface = t.
													Fonts["Primary"]
												cur.Color = t.Colors["Dark"]
												cur.Alignment = text.Start
												cur.Layout(g)
											})
									}),
									)
								}),
								layout.Rigid(func() {
									s.Rc.Dialog.CustomField()
								}),
								layout.Rigid(func() {
									layout.Flex{
										Axis:      layout.Horizontal,
										Alignment: layout.Middle,
									}.Layout(g,
										layout.Rigid(
											s.Context.dialogButton(
												s.Rc.Dialog.Red,
												s.Rc.Dialog.RedLabel,
												"ffcf3030", "iconCancel", "ffcf8080",
												buttonDialogCancel),
										),
										layout.Rigid(
											s.dialogButton(
												s.Rc.Dialog.Green,
												s.Rc.Dialog.GreenLabel,
												"ff30cf30", "iconOK",
												"ff80cf80", buttonDialogOK),
										),
										layout.Rigid(
											s.dialogButton(
												s.Rc.Dialog.Orange,
												s.Rc.Dialog.OrangeLabel,
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

func (s *State) DuoUIinputField(f *Field) func() {
	t, g := s.Thm, s.Gtx
	rcs := s.Rc.Settings
	rdw, rdCfg := rcs.Daemon.Widgets, rcs.Daemon.Config
	return func() {
		switch f.Field.Type {
		case "stringSlice":
			switch f.Field.InputType {
			case "text":
				if f.Field.Model != "MiningAddrs" {
					s.StringsArrayEditor(
						(rdw[f.Field.Model]).(*gel.Editor),
						(rdw[f.Field.Model]).(*gel.Editor).Text(),
						func(e gel.EditorEvent) {
							rdCfg[f.Field.Model] =
								strings.Fields((rdw[f.Field.Model]).(*gel.Editor).Text())
							slog.Debug()
							if e != nil {
								s.Rc.SaveDaemonCfg()
							}
						},
					)()
				}
			default:
			}
		case "input":
			switch f.Field.InputType {
			case "text":
				s.Editor(
					(rdw[f.Field.Model]).(*gel.Editor),
					(rdw[f.Field.Model]).(*gel.Editor).Text(),
					func(e gel.EditorEvent) {
						txt := rdw[f.Field.Model].(*gel.Editor).Text()
						rdCfg[f.Field.Model] = txt
						if e != nil {
							s.Rc.SaveDaemonCfg()
						}
					},
				)()
			case "number":
				s.Editor(
					(rdw[f.Field.Model]).(*gel.Editor),
					(rdw[f.Field.Model]).(*gel.Editor).Text(),
					func(e gel.EditorEvent) {
						number, err :=
							strconv.Atoi((rdw[f.Field.Model]).(*gel.Editor).Text())
						if err == nil {
						}
						rdCfg[f.Field.Model] = number
						if e != nil {
							s.Rc.SaveDaemonCfg()
						}
					},
				)()
			case "decimal":
				s.Editor(
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
							s.Rc.SaveDaemonCfg()
						}
					})()
			case "password":
				s.Editor(
					(rdw[f.Field.Model]).(*gel.Editor),
					(rdw[f.Field.Model]).(*gel.Editor).Text(),
					func(e gel.EditorEvent) {
						txt := rdw[f.Field.Model].(*gel.Editor).Text()
						rdCfg[f.Field.Model] = txt
						if e != nil {
							s.Rc.SaveDaemonCfg()
						}
					},
				)()
			case "time":
				//Debug("rendering duration")
				s.Editor(
					(rdw[f.Field.Model]).(*gel.Editor),
					(rdw[f.Field.Model]).(*gel.Editor).Text(),
					func(e gel.EditorEvent) {
						txt := rdw[f.Field.Model].(*gel.Editor).Text()
						var err error
						if rdCfg[f.Field.Model], err = time.ParseDuration(
							txt); slog.Check(err) {
						}
						if e != nil {
							s.Rc.SaveDaemonCfg()
						}
					},
				)()
			default:
			}
		case "switch":
			sw := t.DuoUICheckBox(
				f.Field.Label, t.Colors["Primary"], t.Colors["Primary"])
			sw.PillColor = t.Colors["LightGray"]
			sw.PillColorChecked = t.Colors["PrimaryDim"]
			sw.CircleColor = t.Colors["LightGrayII"]
			sw.CircleColorChecked = t.Colors["Primary"]
			sw.DrawLayout(g, rdw[f.Field.Model].(*gel.CheckBox))
			if (rdw[f.Field.Model]).(*gel.CheckBox).Checked(g) {
				if !*rdCfg[f.Field.Model].(*bool) {
					tt := true
					rdCfg[f.Field.Model] = &tt
					s.Rc.SaveDaemonCfg()
				}
			} else {
				if *rdCfg[f.Field.Model].(*bool) {
					ff := false
					rdCfg[f.Field.Model] = &ff
					s.Rc.SaveDaemonCfg()
				}
			}
		case "radio":
			// radioButtonsGroup := (duo.Configuration.Settings.Daemon.Widgets[fieldName]).(*widget.Enum)
			// layout.Flex{}.Layout(g,
			//	layout.Rigid(func() {
			//		duo.Theme.RadioButton("r1", "RadioButton1").Layout(g,
			//		radioButtonsGroup)
			//
			//	}),
			//	layout.Rigid(func() {
			//		duo.Theme.RadioButton("r2", "RadioButton2").Layout(g,
			//		radioButtonsGroup)
			//
			//	}),
			//	layout.Rigid(func() {
			//		duo.Theme.RadioButton("r3", "RadioButton3").Layout(g,
			//		radioButtonsGroup)
			//
			//	}))
		default:
			// duo.Theme.CheckBox("Checkbox").Layout(g,
			//(duo.Configuration.Settings.Daemon.Widgets[fieldName]).(*widget.CheckBox))
		}
	}
}

func (s *State) DuoUIlatestTransactions() func() {
	t, g := s.Thm, s.Gtx
	return func() {
		width := g.Constraints.Width.Max
		//cs := g.Constraints
		t.DuoUIContainer(0, t.Colors["DarkGray"]).
			Layout(g, layout.NW, func() {
				layout.Flex{
					Axis: layout.Vertical,
				}.Layout(g,
					layout.Rigid(func() {
						t.DuoUIContainer(8, t.Colors["Primary"]).
							Layout(g, layout.N,
								func() {
									g.Constraints.Width.Min = width
									latestx := t.H5("LATEST TRANSACTIONS")
									latestx.Color = t.Colors["Light"]
									latestx.Alignment = text.Start
									latestx.Layout(g)
								},
							)
					}),
					layout.Flexed(1, func() {
						layout.UniformInset(unit.Dp(8)).Layout(g,
							func() {
								layout.Flex{
									Axis: layout.Vertical,
								}.Layout(g,
									layout.Rigid(func() {
										latestTxsBookPanel := t.DuoUIPanel()
										latestTxsBookPanel.PanelObject =
											s.Rc.Status.Wallet.LastTxs.Txs
										latestTxsBookPanel.ScrollBar =
											t.ScrollBar(0)
										latestTxsPanelElement.PanelObjectsNumber =
											len(s.Rc.Status.Wallet.LastTxs.Txs)
										latestTxsBookPanel.Layout(g,
											latestTxsPanelElement, func(i int, in interface{}) {
												txs := in.([]model.DuoUItransactionExcerpt)
												tx := txs[i]
												t.DuoUIContainer(16,
													t.Colors["Dark"]).
													Layout(g, layout.NW,
														func() {
															width := g.Constraints.Width.Max
															layout.Flex{
																Axis: layout.
																	Vertical,
															}.Layout(g,
																layout.Rigid(
																	s.lTtxid(
																		tx.TxID)),
																layout.Rigid(func() {
																	g.
																		Constraints.Width.Min = width
																	layout.Flex{
																		Spacing: layout.SpaceBetween,
																	}.Layout(g,
																		layout.Rigid(func() {
																			layout.Flex{
																				Axis: layout.Vertical,
																			}.
																				Layout(g,
																					layout.Rigid(s.lTcategory(tx.Category)),
																					layout.Rigid(s.lTtime(tx.Time)),
																				)
																		}),
																		layout.
																			Rigid(
																				s.lTamount(tx.Amount)),
																	)
																}),
																layout.Rigid(s.
																	Thm.
																	DuoUILine(g,
																		0, 0, 1,
																		t.
																			Colors["Hint"]),
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

func (s *State) DuoUIlogger() func() {
	t, g := s.Thm, s.Gtx
	return func() {
		// const buflen = 9
		layout.UniformInset(unit.Dp(10)).Layout(g, func() {
			// const n = 1e6
			cs := g.Constraints
			gelook.DuoUIDrawRectangle(g,
				cs.Width.Max, cs.Height.Max, t.Colors["Dark"],
				[4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
			lm, ok := s.Rc.Log.LogMessages.Load().([]log.Entry)
			if !ok {
				return
			}
			logOutputList.Layout(g, len(lm), func(i int) {
				tt := lm[i]
				logText := t.Caption(
					fmt.Sprintf("%-12s",
						tt.Time.Sub(StartupTime)/time.Second*time.Second) +
						" " + fmt.Sprint(tt.Text))
				logText.Font.Typeface = t.Fonts["Mono"]
				logText.Color = t.Colors["Primary"]
				switch tt.Level {
				case "TRC":
					logText.Color = t.Colors["Success"]
				case "DBG":
					logText.Color = t.Colors["Secondary"]
				case "INF":
					logText.Color = t.Colors["Info"]
				case "WRN":
					logText.Color = t.Colors["Warning"]
				case "ERROR":
					logText.Color = t.Colors["Danger"]
				case "FTL":
					logText.Color = t.Colors["Primary"]
				}
				logText.Layout(g)
			})
			op.InvalidateOp{}.Add(g.Ops)
		})
	}
}

func (s *State) DuoUIstatus() func() {
	t, g := s.Thm, s.Gtx
	wall := s.Rc.Status.Wallet
	nod := s.Rc.Status.Node
	return func() {
		t.DuoUIContainer(8, t.Colors["Light"]).Layout(g, layout.NW,
			func() {
				bigStatus := []func(){
					s.listItem(22, 6, "EditorMonetizationOn",
						"BALANCE :", wall.Balance.Load()+" "+s.Rc.Settings.
							Abbrevation),
					t.DuoUILine(g,
						8, 0, 1, t.Colors["LightGray"]),
					s.listItem(22, 6, "MapsLayersClear", "UNCONFIRMED :",
						wall.Unconfirmed.Load()+" "+s.Rc.Settings.Abbrevation),
					t.DuoUILine(g,
						8, 0, 1, t.Colors["LightGray"]),
					s.listItem(22, 6, "CommunicationImportExport",
						"TRANSACTIONS :", fmt.Sprint(wall.TxsNumber.Load())),
					t.DuoUILine(g,
						8, 0, 1, t.Colors["LightGray"]),
					s.listItem(16, 4, "DeviceWidgets",
						"Block Count :", fmt.Sprint(nod.BlockCount.Load())),
					t.DuoUILine(g,
						4, 0, 1, t.Colors["LightGray"]),
					s.listItem(16, 4, "ImageTimer",
						"Difficulty :", fmt.Sprint(nod.Difficulty.Load())),
					t.DuoUILine(g,
						4, 0, 1, t.Colors["LightGray"]),
					s.listItem(16, 4, "NotificationVPNLock",
						"Connections :", fmt.Sprint(nod.ConnectionCount.Load())),
				}
				itemsList.Layout(g, len(bigStatus), func(i int) {
					layout.UniformInset(unit.Dp(0)).Layout(g, bigStatus[i])
				})
			})
	}
}

func (s *State) FooterLeftMenu(allPages *model.DuoUIpages) func() {
	g := s.Gtx
	return func() {
		cornerButtons := []func(){
			s.QuitButton(),
			// s.footerMenuButton(allPages.Theme["EXPLORER"],
			//"BLOCKS: "+fmt.Sprint(rc.Status.Node.BlockCount), "", buttonBlocks),
			s.footerMenuButton(allPages.Theme["LOG"], "LOG", "traceIcon",
				buttonLog),
		}
		cornerNav.Layout(g, len(cornerButtons), func(i int) {
			layout.UniformInset(unit.Dp(0)).Layout(g, cornerButtons[i])
		})
	}
}

func (s *State) footerMenuButton(page *gelook.DuoUIPage, text, icon string, footerButton *gel.Button) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(s.Gtx, func() {
			var footerMenuItem gelook.DuoUIbutton
			if icon != "" {
				footerMenuItem = s.Thm.DuoUIbutton(gelook.ButtonParams{
					BgHoverColor: s.Thm.Colors["Dark"],
					Icon:         icon,
					IconColor: CurrentCurrentPageColor(s.Rc.ShowPage,
						page.Title, navItemIconColor, s.Thm.Colors["Primary"]),
					TextSize: footerMenuItemTextSize, IconSize: footerMenuItemIconSize,
					Width: footerMenuItemWidth, Height: footerMenuItemHeight})
				for footerButton.Clicked(s.Gtx) {
					s.Rc.ShowPage = page.Title
					SetPage(s.Rc, page)
				}
				footerMenuItem.IconLayout(s.Gtx, footerButton)
			} else {
				footerMenuItem = s.Thm.DuoUIbutton(gelook.ButtonParams{
					TxtFont: s.Thm.Fonts["Primary"],
					Txt:     text,
					TxtColor: CurrentCurrentPageColor(s.Rc.ShowPage, page.Title,
						s.Thm.Colors["Light"], s.Thm.Colors["Primary"]),
					TextSize:      footerMenuItemTextSize,
					IconSize:      footerMenuItemIconSize,
					Height:        footerMenuItemHeight,
					PaddingTop:    13,
					PaddingRight:  16,
					PaddingBottom: 14,
					PaddingLeft:   16,
				})
				footerMenuItem.Height = 48
				for footerButton.Clicked(s.Gtx) {
					s.Rc.ShowPage = page.Title
					SetPage(s.Rc, page)
				}
				footerMenuItem.Layout(s.Gtx, footerButton)
			}
		})
	}
}

func (s *State) FooterRightMenu(allPages *model.DuoUIpages) func() {
	g := s.Gtx
	return func() {
		navButtons := []func(){
			// s.footerMenuButton(allPages.Theme["NETWORK"],
			// "", "networkIcon", buttonNetwork),
			// s.footerMenuButton(allPages.Theme["NETWORK"],
			// "CONNECTIONS: "+fmt.Sprint(rc.Status.Node.ConnectionCount),
			//"", buttonNetwork),
			// s.footerMenuButton(allPages.Theme["EXPLORER"],
			// "", "DeviceWidgets", buttonBlocks),
			s.footerMenuButton(allPages.Theme["NETWORK"],
				"", "network", buttonNetwork),
			s.footerMenuButton(allPages.Theme["NETWORK"],
				"CONNECTIONS: "+
					fmt.Sprint(s.Rc.Status.Node.ConnectionCount.Load()),
				"", buttonNetwork),
			s.footerMenuButton(allPages.Theme["EXPLORER"],
				"", "DeviceWidgets", buttonBlocks),
			s.footerMenuButton(allPages.Theme["EXPLORER"],
				"BLOCKS: "+fmt.Sprint(s.Rc.Status.Node.BlockCount.Load()), "",
				buttonBlocks),
			s.footerMenuButton(allPages.Theme["MINER"],
				"", "helpIcon", buttonHelp),
			s.footerMenuButton(allPages.Theme["CONSOLE"],
				"", "consoleIcon", buttonConsole),
			s.footerMenuButton(allPages.Theme["SETTINGS"],
				"", "settingsIcon", buttonSettings),
		}
		footerNav.Layout(g, len(navButtons), func(i int) {
			layout.UniformInset(unit.Dp(0)).Layout(g, navButtons[i])
		})
	}
}

func (s *State) HeaderMenu(allPages *model.DuoUIpages) func() {
	g := s.Gtx
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(g, func() {
			headerNav := []func(){
				s.headerMenuButton("",
					"CommunicationImportExport", buttonHeader),
				s.headerMenuButton("",
					"NotificationNetworkCheck", buttonHeader),
				s.headerMenuButton("",
					"NotificationSync", buttonHeader),
				s.headerMenuButton("",
					"NotificationSyncDisabled", buttonHeader),
				s.headerMenuButton("",
					"NotificationSyncProblem", buttonHeader),
				s.headerMenuButton("",
					"NotificationVPNLock", buttonHeader),
				s.headerMenuButton("",
					"MapsLayers", buttonHeader),
				s.headerMenuButton("",
					"MapsLayersClear", buttonHeader),
				s.headerMenuButton("",
					"ImageTimer", buttonHeader),
				s.headerMenuButton("",
					"ImageRemoveRedEye", buttonHeader),
				s.headerMenuButton("",
					"DeviceSignalCellular0Bar", buttonHeader),
				s.headerMenuButton("",
					"ActionTimeline", buttonHeader),
				s.headerMenuButton("",
					"HardwareWatch", buttonHeader),
			}
			footerNav.Layout(g, len(headerNav), func(i int) {
				layout.UniformInset(unit.Dp(0)).Layout(g, headerNav[i])
			})
		})
	}
}

func (s *State) headerMenuButton(text, icon string, headerButton *gel.Button) func() {
	t, g := s.Thm, s.Gtx
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(g, func() {
			var footerMenuItem gelook.DuoUIbutton
			footerMenuItem = t.DuoUIbutton(gelook.ButtonParams{
				BgHoverColor: t.Colors["Dark"],
				Icon:         icon,
				IconColor: CurrentCurrentPageColor(
					s.Rc.ShowPage, text, navItemIconColor, t.Colors["Primary"]),
				TextSize:      footerMenuItemTextSize,
				IconSize:      footerMenuItemIconSize,
				Width:         footerMenuItemWidth,
				Height:        footerMenuItemHeight,
				PaddingTop:    footerMenuItemPaddingVertical,
				PaddingRight:  footerMenuItemPaddingHorizontal,
				PaddingBottom: footerMenuItemPaddingVertical,
				PaddingLeft:   footerMenuItemPaddingHorizontal})
			for headerButton.Clicked(g) {
				s.Rc.ShowPage = text
			}
			footerMenuItem.IconLayout(g, headerButton)
		})
	}
}

func (s *State) iconButton(page *gelook.DuoUIPage) func() {
	t, g := s.Thm, s.Gtx
	return func() {
		var logMenuItem gelook.DuoUIbutton
		logMenuItem = t.DuoUIbutton(gelook.ButtonParams{
			BgColor: t.Colors["Dark"],
			Icon:    "traceIcon",
			IconColor: CurrentCurrentPageColor(
				s.Rc.ShowPage, "LOG", t.Colors["Light"], t.Colors["Primary"]),
			TextSize:      footerMenuItemTextSize,
			IconSize:      footerMenuItemIconSize,
			Width:         footerMenuItemWidth,
			Height:        footerMenuItemHeight,
			PaddingTop:    footerMenuItemPaddingVertical,
			PaddingRight:  footerMenuItemPaddingHorizontal,
			PaddingBottom: footerMenuItemPaddingVertical,
			PaddingLeft:   footerMenuItemPaddingHorizontal,
		})
		for buttonLog.Clicked(g) {
			SetPage(s.Rc, page)
			s.Rc.ShowPage = "LOG"
		}
		logMenuItem.IconLayout(g, buttonLog)
	}
}

func (s *State) MainNavigation(allPages *model.DuoUIpages, nav *model.DuoUInav) func() {
	t, g := s.Thm, s.Gtx
	return func() {
		navButtons := s.navButtons(allPages, nav)
		g.Constraints.Width.Max = nav.Width
		t.DuoUIContainer(0, t.Colors["Dark"]).Layout(g, layout.NW,
			func() {
				mainNav.Layout(g, len(navButtons), func(i int) {
					layout.UniformInset(unit.Dp(0)).Layout(g, navButtons[i])
				})
			})
	}
}

func (s *State) navButtons(allPages *model.DuoUIpages, nav *model.DuoUInav) []func() {
	t, g := s.Thm, s.Gtx
	return []func(){
		s.navMenuButton(allPages.Theme["OVERVIEW"], nav,
			"OVERVIEW", "overviewIcon", navButtonOverview),
		t.DuoUILine(g, 0, 0, 1,
			t.Colors["LightGrayIII"]),
		s.navMenuButton(allPages.Theme["SEND"], nav, "SEND", "sendIcon",
			navButtonSend),
		// navMenuLine(g, th),
		// navMenuButton(rc, g, t, allPages.Theme["RECEIVE"], "RECEIVE",
		//"receiveIcon", navButtonReceive),
		t.DuoUILine(g, 0, 0, 1,
			t.Colors["LightGrayIII"]),
		s.navMenuButton(allPages.Theme["ADDRESSBOOK"], nav,
			"ADDRESSBOOK", "addressBookIcon", navButtonAddressBook),
		t.DuoUILine(g, 0, 0, 1,
			t.Colors["LightGrayIII"]),
		s.navMenuButton(allPages.Theme["HISTORY"], nav,
			"HISTORY", "historyIcon", navButtonHistory),
	}
}

func (s *State) navMenuButton(page *gelook.DuoUIPage, nav *model.DuoUInav, title, icon string, navButton *gel.Button) func() {
	t, g := s.Thm, s.Gtx
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(g, func() {
			var menuItem gelook.DuoUIbutton
			menuItem = t.DuoUIbutton(gelook.ButtonParams{TxtFont: t.Fonts["Secondary"],
				Txt:           title,
				TxtColor:      t.Colors["Dark"],
				BgColor:       t.Colors["LightGrayII"],
				TxtHoverColor: t.Colors["LightGrayII"],
				BgHoverColor:  t.Colors["Dark"],
				Icon:          icon,
				IconColor: CurrentCurrentPageColor(
					s.Rc.ShowPage, title, navItemIconColor,
					t.Colors["Primary"]),
				TextSize:      nav.TextSize,
				IconSize:      nav.IconSize,
				Width:         nav.Width,
				Height:        nav.Height,
				PaddingTop:    nav.PaddingVertical,
				PaddingRight:  nav.PaddingHorizontal,
				PaddingBottom: nav.PaddingVertical,
				PaddingLeft:   nav.PaddingHorizontal,
			})
			for navButton.Clicked(g) {
				s.Rc.ShowPage = title
				SetPage(s.Rc, page)
			}
			menuItem.MenuLayout(g, navButton)
		})
	}
}

func (s *State) pageNavButton(page *gelook.DuoUIPage, b *gel.Button, label, hash string) {
	t, g := s.Thm, s.Gtx
	layout.UniformInset(unit.Dp(4)).Layout(g, func() {
		var blockButton gelook.DuoUIbutton
		blockButton = t.DuoUIbutton(gelook.ButtonParams{TxtFont: t.Fonts["Mono"],
			Txt:           label + " " + hash,
			TxtColor:      t.Colors["Light"],
			BgColor:       t.Colors["Info"],
			TxtHoverColor: t.Colors["Info"],
			BgHoverColor:  t.Colors["Light"],
			IconColor:     t.Colors["Light"],
			TextSize:      16,
			Width:         60,
			Height:        24,
		})
		for b.Clicked(g) {
			s.Rc.ShowPage = fmt.Sprintf("BLOCK %s", hash)
			s.Rc.GetSingleBlock(hash)()
			SetPage(s.Rc, page)
		}
		blockButton.Layout(g, b)
	})
}

func (s *State) PageNavButtons(previousBlockHash, nextBlockHash string, prevPage, nextPage *gelook.DuoUIPage) func() {
	g := s.Gtx
	return func() {
		layout.Flex{}.Layout(g,
			layout.Flexed(0.5, func() {
				eh := chainhash.Hash{}
				if previousBlockHash != eh.String() {
					s.pageNavButton(nextPage, previousBlockHashButton,
						"Previous Block", previousBlockHash)
				}
			}),
			layout.Flexed(0.5, func() {
				if nextBlockHash != "" {
					s.pageNavButton(nextPage, nextBlockHashButton,
						"Next Block", nextBlockHash)
				}
			}),
		)
	}
}

func (s *State) PeersList() func() {
	t, g := s.Thm, s.Gtx
	c := t.Colors["Dark"]
	return func() {
		s.Rc.Network.PeersList.Layout(g, len(s.Rc.Network.Peers), func(i int) {
			np := s.Rc.Network.Peers[i]
			t.DuoUILine(g, 0, 0, 1,
				t.Colors["Hint"])()
			layout.Flex{
				Spacing: layout.SpaceBetween,
			}.Layout(g,
				layout.Rigid(s.Context.Label(t.Fonts["Mono"], 14, c,
					fmt.Sprint(np.ID))),
				layout.Flexed(1, s.peerDetails(i, np)),
				layout.Rigid(s.Label(t.Fonts["Mono"], 14, c, np.Addr)))
		})
	}
}

func (s *State) QuitButton() func() {
	t, g := s.Thm, s.Gtx
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(g, func() {
			var closeMenuItem gelook.DuoUIbutton
			closeMenuItem = t.DuoUIbutton(gelook.ButtonParams{
				BgColor: t.Colors["Dark"],
				Icon:    "closeIcon",
				IconColor: CurrentCurrentPageColor(
					s.Rc.ShowPage, "CLOSE", t.Colors["Light"],
					t.Colors["Primary"]),
				TextSize: footerMenuItemTextSize,
				IconSize: footerMenuItemIconSize,
				Width:    footerMenuItemWidth,
				Height:   footerMenuItemHeight,
			})
			for buttonQuit.Clicked(g) {
				s.Rc.Dialog.Show = true
				s.Rc.Dialog = &model.DuoUIdialog{
					Show: true,
					Green: func() {
						interrupt.Request()
						// TODO make this close the window or at least switch to a shutdown screen
						s.Rc.Dialog.Show = false
					},
					GreenLabel: "QUIT",
					Orange: func() {
						interrupt.RequestRestart()
						// TODO make this close the window or at least switch to a shutdown screen
						s.Rc.Dialog.Show = false
					},
					OrangeLabel: "RESTART",
					Red:         func() { s.Rc.Dialog.Show = false },
					RedLabel:    "CANCEL",
					CustomField: func() {},
					Title:       "Are you sure?",
					Text:        "Confirm ParallelCoin close",
				}
			}
			closeMenuItem.IconLayout(g, buttonQuit)
		})
	}
}

func (s *State) SettingsTabs() func() {
	t, g := s.Thm, s.Gtx
	rcs := s.Rc.Settings
	return func() {
		groupsNumber := len(rcs.Daemon.Schema.Groups)
		groupsList.Layout(g, groupsNumber, func(i int) {
			layout.UniformInset(unit.Dp(0)).Layout(g, func() {
				col := t.Colors["Light"]
				bgColor := t.Colors["Dark"]
				i = groupsNumber - 1 - i
				tt := rcs.Daemon.Schema.Groups[i]
				txt := fmt.Sprint(tt.Legend)
				for rcs.Tabs.TabsList[txt].Clicked(g) {
					rcs.Tabs.Current = txt
				}
				if rcs.Tabs.Current == txt {
					col = t.Colors["Dark"]
					bgColor = t.Colors["Light"]
				}
				t.DuoUIbutton(gelook.ButtonParams{
					TxtFont:       t.Fonts["Primary"],
					Txt:           txt,
					TxtColor:      col,
					BgColor:       bgColor,
					TextSize:      16,
					Width:         80,
					Height:        32,
					PaddingTop:    4,
					PaddingRight:  4,
					PaddingBottom: 4,
					PaddingLeft:   4,
				}).Layout(g, rcs.Tabs.TabsList[txt])
			})
		})
	}
}

func (s *State) TransactionsFilter() func() {
	g := s.Gtx
	cats := s.Rc.History.Categories
	return func() {
		layout.Flex{}.Layout(g,
			layout.Rigid(s.txsFilterItem("ALL", cats.AllTxs)),
			layout.Rigid(s.txsFilterItem("MINTED", cats.MintedTxs)),
			layout.Rigid(s.txsFilterItem("IMATURE", cats.ImmatureTxs)),
			layout.Rigid(s.txsFilterItem("SENT", cats.SentTxs)),
			layout.Rigid(s.txsFilterItem("RECEIVED", cats.ReceivedTxs)))
		switch {
		case cats.AllTxs.Checked(g):
			s.Rc.History.Category = "all"
		case cats.MintedTxs.Checked(g):
			s.Rc.History.Category = "generate"
		case cats.ImmatureTxs.Checked(g):
			s.Rc.History.Category = "immature"
		case cats.SentTxs.Checked(g):
			s.Rc.History.Category = "sent"
		case cats.ReceivedTxs.Checked(g):
			s.Rc.History.Category = "received"
		}
	}
}

// gtx, th

func (s *Context) Button(buttonController *gel.Button, font text.Typeface, textSize, padT, padR, padB, padL int, color, bgColor, label string, handler func()) func() {
	t, g := s.Thm, s.Gtx
	return func() {
		t.DuoUIContainer(0, t.Colors[""])
		button := t.DuoUIbutton(gelook.ButtonParams{
			TxtFont:       font,
			Txt:           label,
			TxtColor:      color,
			BgColor:       bgColor,
			TextSize:      textSize,
			Width:         128,
			Height:        48,
			PaddingTop:    padT,
			PaddingRight:  padR,
			PaddingBottom: padB,
			PaddingLeft:   padL,
		})
		for buttonController.Clicked(g) {
			handler()
		}
		button.Layout(g, buttonController)
	}
}

func (s *Context) ConsoleInput(editorController *gel.Editor, label string, handler func(gel.SubmitEvent)) func() {
	t, g := s.Thm, s.Gtx
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(g, func() {
			e := t.DuoUIEditor(label, "Dark", "Light", 120)
			e.Font.Typeface = t.Fonts["Primary"]
			e.Color = gelook.HexARGB(t.Colors["Light"])
			e.Font.Style = text.Italic
			e.Layout(g, editorController)
			for _, e := range editorController.Events(g) {
				if e, ok := e.(gel.SubmitEvent); ok {
					handler(e)
					editorController.SetText("")
				}
			}
		})
	}
}

func (s *Context) contentField(text, color, bgColor string, font text.Typeface, padding, textSize float32) func() {
	t, g := s.Thm, s.Gtx
	return func() {
		hmin := g.Constraints.Width.Min
		vmin := g.Constraints.Height.Min
		layout.Stack{Alignment: layout.W}.Layout(g,
			layout.Expanded(func() {
				rr := float32(g.Px(unit.Dp(0)))
				clip.Rect{
					Rect: f32.Rectangle{Max: f32.Point{
						X: float32(g.Constraints.Width.Min),
						Y: float32(g.Constraints.Height.Min),
					}},
					NE: rr, NW: rr, SE: rr, SW: rr,
				}.Op(g.Ops).Add(g.Ops)
				fill(g, gelook.HexARGB(bgColor))
			}),
			layout.Stacked(func() {
				g.Constraints.Width.Min = hmin
				g.Constraints.Height.Min = vmin
				layout.Center.Layout(g, func() {
					layout.UniformInset(unit.Dp(padding)).Layout(g, func() {
						l := t.DuoUILabel(unit.Dp(textSize), text)
						l.Font.Typeface = font
						l.Color = color
						l.Layout(g)
					})
				})
			}),
		)
	}
}

func (s *Context) ContentLabeledField(axis layout.Axis, margin, labelTextSize, valueTextSize float32, label, headcolor, headbgColor, color, bgColor, value string) func() {
	t, g := s.Thm, s.Gtx
	return func() {
		layout.UniformInset(unit.Dp(margin)).Layout(g, func() {
			layout.Flex{
				Axis: axis,
			}.Layout(g,
				layout.Rigid(s.contentField(label,
					t.Colors[headcolor], t.Colors[headbgColor],
					t.Fonts["Primary"], 4, labelTextSize),
				),
				layout.Rigid(s.contentField(value,
					t.Colors[color], t.Colors[bgColor],
					t.Fonts["Mono"], 4, valueTextSize),
				),
			)
		})
	}
}

func (s *Context) dialogButton(f func(), txt, bgColor, icon, iconColor string,
	button *gel.Button) func() {
	t, g := s.Thm, s.Gtx
	return func() {
		if f != nil {
			var b gelook.DuoUIbutton
			layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8),
				Left: unit.Dp(8), Right: unit.Dp(8)}.Layout(g, func() {
				b = t.DuoUIbutton(gelook.ButtonParams{
					TxtFont:       t.Fonts["Primary"],
					Txt:           txt,
					TxtColor:      t.Colors["Dark"],
					BgColor:       bgColor,
					TxtHoverColor: t.Colors["Info"],
					BgHoverColor:  bgColor,
					Icon:          icon,
					IconColor:     iconColor,
					TextSize:      16,
					IconSize:      32,
					Width:         120,
					Height:        64,
				})
				for button.Clicked(g) {
					f()
				}
				b.MenuLayout(g, button)
			})
		}
	}
}

func (s *Context) Editor(editorController *gel.Editor, label string, handler func(gel.EditorEvent)) func() {
	t, g := s.Thm, s.Gtx
	return func() {
		t.DuoUIContainer(8, "ffffffff").
			Layout(g, layout.NW, func() {
				width := g.Constraints.Width.Max
				e := t.DuoUIEditor(label, t.Colors["Black"],
					t.Colors["White"], width)
				e.Font.Typeface = t.Fonts["Mono"]
				e.TextSize = unit.Dp(12)
				layout.UniformInset(unit.Dp(4)).Layout(g, func() {
					e.Layout(g, editorController)
				})
				for _, e := range editorController.Events(g) {
					switch e.(type) {
					case gel.ChangeEvent:
						handler(e)
					}
				}
			})
	}
}

func (s *Context) Label(font text.Typeface, size float32, color, label string) func() {
	t, g := s.Thm, s.Gtx
	return func() {
		l := t.DuoUILabel(unit.Dp(size), label)
		l.Font.Typeface = font
		l.Color = color
		l.Layout(g)
	}
}

func (s *Context) listItem(size, top int, iconName, name, value string) func() {
	t, g := s.Thm, s.Gtx
	return func() {
		icon := t.Icons[iconName]
		layout.Flex{
			Axis:    layout.Horizontal,
			Spacing: layout.SpaceBetween,
		}.Layout(g,
			layout.Rigid(func() {
				layout.Flex{}.Layout(g,
					layout.Rigid(func() {
						layout.Inset{
							Top:    unit.Dp(float32(top)),
							Bottom: unit.Dp(0), Left: unit.Dp(0),
							Right: unit.Dp(0)}.Layout(g, func() {
							if icon != nil {
								icon.Color = gelook.HexARGB(t.Colors["Dark"])
								icon.Layout(g, unit.Px(float32(size)))
							}
							g.Dimensions = layout.Dimensions{
								Size: image.Point{X: size, Y: size},
							}
						})
					}),
					layout.Rigid(func() {
						txt := t.DuoUILabel(unit.Dp(float32(size)), name)
						txt.Font.Typeface = t.Fonts["Primary"]
						txt.Color = t.Colors["Primary"]
						txt.Layout(g)
					}),
				)
			}),
			layout.Rigid(func() {
				v := t.H5(value)
				v.TextSize = unit.Dp(float32(size))
				v.Font.Typeface = t.Fonts["Primary"]
				v.Color = t.Colors["Dark"]
				v.Alignment = text.End
				v.Layout(g)
			}),
		)
	}
}

func (s *Context) lTamount(v float64) func() {
	t, g := s.Thm, s.Gtx
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(g, func() {
			sat := t.Body1(fmt.Sprintf("%0.8f", v))
			sat.Font.Typeface = "bariol"
			sat.Color = t.Colors["Light"]
			sat.Layout(g)
		})
	}
}

func (s *Context) lTcategory(v string) func() {
	t, g := s.Thm, s.Gtx
	return func() {
		sat := t.Body1(v)
		sat.Color = t.Colors["Light"]
		sat.Font.Typeface = "bariol"
		sat.Layout(g)
	}
}

func (s *Context) lTtime(v string) func() {
	t, g := s.Thm, s.Gtx
	return func() {
		l := t.Body1(v)
		l.Font.Typeface = "bariol"
		l.Color = t.Colors["Light"]
		l.Color = t.Colors["Hint"]
		l.Layout(g)
	}
}

func (s *Context) lTtxid(v string) func() {
	t, g := s.Thm, s.Gtx
	return func() {
		tim := t.Caption(v)
		tim.Font.Typeface = t.Fonts["Mono"]
		tim.Color = t.Colors["Light"]
		tim.Layout(g)
	}
}

func (s *Context) MonoButton(buttonController *gel.Button, textSize int, color, bgColor, font, label string, handler func()) func() {
	t, g := s.Thm, s.Gtx
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(g, func() {
			//var button gelook.Button
			button := t.Button(label)
			switch {
			case font != "":
				button.Font.Typeface = t.Fonts[font]
			case color != "":
				button.Color = gelook.HexARGB(t.Colors[color])
			case textSize != 0:
				button.TextSize = unit.Dp(float32(textSize))
			case bgColor != "":
				button.Background = gelook.HexARGB(t.Colors[bgColor])
			}
			for buttonController.Clicked(g) {
				handler()
			}
			button.Layout(g, buttonController)
		})
	}
}

func (s *Context) peerDetails(i int, pi *btcjson.GetPeerInfoResult) func() {
	t, g := s.Thm, s.Gtx
	prim := t.Fonts["Primary"]
	c := t.Colors["Dark"]
	return func() {
		layout.Flex{
			Axis:    layout.Horizontal,
			Spacing: layout.SpaceAround,
		}.Layout(g,
			layout.Rigid(s.Label(prim, 12, c, pi.AddrLocal)),
			layout.Rigid(s.Label(prim, 12, c, pi.Services)),
			layout.Rigid(s.Label(prim, 12, c, fmt.Sprint(pi.RelayTxes))),
			// layout.Rigid(s.Label(prim, 12, c
			//fmt.Sprint(t.LastSend))),
			// layout.Rigid(s.Label(prim, 12, c
			//fmt.Sprint(t.LastRecv))),
			layout.Rigid(s.Label(prim, 12, c, fmt.Sprint(pi.BytesSent))),
			layout.Rigid(s.Label(prim, 12, c, fmt.Sprint(pi.BytesRecv))),
			layout.Rigid(s.Label(prim, 12, c, fmt.Sprint(pi.ConnTime))),
			layout.Rigid(s.Label(prim, 12, c, fmt.Sprint(pi.TimeOffset))),
			layout.Rigid(s.Label(prim, 12, c, fmt.Sprint(pi.PingTime))),
			layout.Rigid(s.Label(prim, 12, c, fmt.Sprint(pi.PingWait))),
			layout.Rigid(s.Label(prim, 12, c, fmt.Sprint(pi.Version))),
			layout.Rigid(s.Label(prim, 12, c, pi.SubVer)),
			layout.Rigid(s.Label(prim, 12, c, fmt.Sprint(pi.Inbound))),
			layout.Rigid(s.Label(prim, 12, c, fmt.Sprint(pi.StartingHeight))),
			layout.Rigid(s.Label(prim, 12, c, fmt.Sprint(pi.CurrentHeight))),
			layout.Rigid(s.Label(prim, 12, c, fmt.Sprint(pi.BanScore))),
			layout.Rigid(s.Label(prim, 12, c, fmt.Sprint(pi.FeeFilter))),
			layout.Rigid(s.Label(prim, 12, c, fmt.Sprint(pi.SyncNode))))
	}
}

func (s *Context) SettingsFieldDescription(f *Field) func() {
	t, g := s.Thm, s.Gtx
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(g, func() {
			desc := t.Body2(fmt.Sprint(f.Field.Description))
			desc.Font.Typeface = t.Fonts["Primary"]
			desc.Color = t.Colors["Dark"]
			desc.Layout(g)
		})
	}
}

func (s *Context) SettingsFieldLabel(f *Field) func() {
	t, g := s.Thm, s.Gtx
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(g, func() {
			name := t.H6(fmt.Sprint(f.Field.Label))
			name.Color = t.Colors["Dark"]
			name.Font.Typeface = t.Fonts["Primary"]
			name.Layout(g)
		})
	}
}

func (s *Context) StringsArrayEditor(editorController *gel.Editor, label string, handler func(gel.EditorEvent)) func() {
	t, g := s.Thm, s.Gtx
	return func() {
		t.DuoUIContainer(8, t.Colors["White"]).Layout(g, layout.NW,
			func() {
				e := t.DuoUIEditor(label, t.Colors["Black"],
					t.Colors["White"], 16)
				e.Font.Typeface = t.Fonts["Mono"]
				layout.UniformInset(unit.Dp(4)).Layout(g, func() {
					e.Layout(g, editorController)
				})
				for _, e := range editorController.Events(g) {
					switch e.(type) {
					case gel.ChangeEvent:
						handler(e)
					}
				}
			})
	}
}

func (s *Context) TrioFields(axis layout.Axis, labelTextSize, valueTextSize float32, unoLabel, unoValue, unoHeadcolor, unoHeadbgColor, unoColor, unoBgColor, duoLabel, duoValue, duoHeadcolor, duoHeadbgColor, duoColor, duoBgColor, treLabel, treValue, treHeadcolor, treHeadbgColor, treColor, treBgColor string) func() {
	g := s.Gtx
	return func() {
		layout.Flex{
			Axis:    axis,
			Spacing: layout.SpaceAround,
		}.Layout(g,
			fieldAxis(axis, s.ContentLabeledField(
				layout.Vertical, 4,
				labelTextSize, valueTextSize,
				unoLabel, unoHeadcolor, unoHeadbgColor,
				unoColor, unoBgColor,
				fmt.Sprint(unoValue)), 0.3333),
			fieldAxis(axis, s.ContentLabeledField(
				layout.Vertical, 4,
				labelTextSize, valueTextSize,
				duoLabel, duoHeadcolor, duoHeadbgColor,
				duoColor, duoBgColor,
				fmt.Sprint(duoValue)), 0.3333),
			fieldAxis(axis, s.ContentLabeledField(
				layout.Vertical, 4,
				labelTextSize, valueTextSize,
				treLabel, treHeadbgColor, treHeadcolor,
				treColor, treBgColor,
				fmt.Sprint(treValue)), 0.3333),
		)
	}
}

func (s *Context) TxsDetails(i int, te *model.DuoUItransactionExcerpt) func() {
	t, g := s.Thm, s.Gtx
	return func() {
		layout.Flex{
			Axis: layout.Vertical,
		}.Layout(g,
			layout.Rigid(s.Label(t.Fonts["Primary"], 12,
				t.Colors["Dark"], fmt.Sprint(i))),
			layout.Rigid(s.Label(t.Fonts["Primary"], 12,
				t.Colors["Dark"], te.TxID)),
			layout.Rigid(s.Label(t.Fonts["Primary"], 12,
				t.Colors["Dark"], fmt.Sprintf("%0.8f", te.Amount))),
			layout.Rigid(s.Label(t.Fonts["Primary"], 12,
				t.Colors["Dark"], te.Category)),
			layout.Rigid(s.Label(t.Fonts["Primary"], 12,
				t.Colors["Dark"], te.Time)),
		)
	}
}

func (s *Context) txsFilterItem(id string, c *gel.CheckBox) func() {
	t, g := s.Thm, s.Gtx
	return func() {
		t.DuoUICheckBox(id, t.Colors["Light"],
			t.Colors["Light"]).Layout(g, c)
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

func SetPage(rc *rcd.RcVar, page *gelook.DuoUIPage) {
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
