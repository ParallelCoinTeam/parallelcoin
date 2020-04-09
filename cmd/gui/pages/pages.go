package pages

import (
	"encoding/json"
	"fmt"
	"gioui.org/op"
	"gioui.org/text"
	"github.com/p9c/pod/pkg/chain/fork"
	"github.com/p9c/pod/pkg/pod"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	"strconv"
	"time"

	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/p9c/pod/cmd/gui/component"
	"github.com/p9c/pod/cmd/gui/model"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gui/clipboard"
	"github.com/p9c/pod/pkg/gui/gel"
	"github.com/p9c/pod/pkg/gui/gelook"
)

type send struct {
	address    string
	amount     float64
	passPhrase string
}

var (
	addressBookPanelElement = gel.NewPanel()
	showMiningAddresses     = &gel.CheckBox{}
	buttonNewAddress        = new(gel.Button)
	address                 string
	previousBlockHashButton = new(gel.Button)
	nextBlockHashButton     = new(gel.Button)
	algoHeadColor, algoHeadBgColor,
	algoColor, algoBgColor string
	// itemValue = &gel.DuoUIcounter{
	//	Value:        11,
	//	OperateValue: 1,
	//	From:         0,
	//	To:           15,
	//	CounterInput: &gel.Editor{
	//		Alignment:  text.Middle,
	//		SingleLine: true,
	//	},
	//	CounterIncrease: new(gel.Button),
	//	//CounterDecrease: new(controller.Button),
	//	CounterReset: new(gel.Button),
	// }
	transactionsPanelElement = gel.NewPanel()
	consoleInputField        = &gel.Editor{
		SingleLine: true,
		Submit:     true,
	}
	consoleOutputList = &layout.List{
		Axis:        layout.Vertical,
		ScrollToEnd: true,
	}
	explorerPanelElement = gel.NewPanel()
	txwidth              int
	logOutputList        = &layout.List{
		Axis:        layout.Vertical,
		ScrollToEnd: true,
	}
	// itemValue = &gel.DuoUIcounter{
	//	Value:        11,
	//	OperateValue: 1,
	//	From:         0,
	//	To:           15,
	//	CounterInput: &gel.Editor{
	//		Alignment:  text.Middle,
	//		SingleLine: true,
	//	},
	//	CounterIncrease: new(gel.Button),
	//	//CounterDecrease: new(controller.Button),
	//	CounterReset: new(gel.Button),
	// }
	StartupTime = time.Now()
	layautList  = &layout.List{
		Axis: layout.Vertical,
	}
	addressLineEditor = &gel.Editor{
		SingleLine: true,
	}
	amountLineEditor = &gel.Editor{
		SingleLine: true,
	}
	passLineEditor = &gel.Editor{
		SingleLine: true,
	}
	buttonPasteAddress    = new(gel.Button)
	buttonPasteAmount     = new(gel.Button)
	buttonSend            = new(gel.Button)
	sendStruct            = new(send)
	settingsPanelElement  = gel.NewPanel()
	buttonSettingsRestart = new(gel.Button)
)

func addressBookBody(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	return func() {
		layout.Flex{}.Layout(gtx,
			layout.Flexed(1, func() {
				layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
					layout.Flex{
						Axis:    layout.Vertical,
						Spacing: layout.SpaceAround,
					}.Layout(gtx,
						layout.Flexed(1, addressBookContent(rc, gtx, th)))
				})
			}))
	}
}

func addressBookContent(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	return func() {
		addressBookPanelElement.PanelObject = rc.AddressBook.Addresses
		addressBookPanelElement.PanelObjectsNumber = len(rc.AddressBook.Addresses)
		addressBookPanel := th.DuoUIpanel()
		addressBookPanel.ScrollBar = th.ScrollBar(16)
		addressBookPanel.Layout(gtx, addressBookPanelElement, func(i int, in interface{}) {
			//if in != nil {
			//addresses := in.([]model.DuoUIaddress)
			t := rc.AddressBook.Addresses[i]
			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func() {
					layout.Flex{
						Alignment: layout.Middle,
					}.Layout(gtx,
						layout.Flexed(0.2,
							component.Label(gtx, th, th.Fonts["Primary"],
								12, th.Colors["Dark"], fmt.Sprint(t.Index)),
						),
						layout.Flexed(0.2,
							component.Label(gtx, th, th.Fonts["Primary"],
								12, th.Colors["Dark"], t.Account),
						),
						layout.Rigid(component.MonoButton(gtx, th,
							t.Copy, 12, "", "",
							"Mono", t.Address, func() {
								clipboard.Set(t.Address)
							}),
						),
						layout.Flexed(0.4,
							component.Label(gtx, th, th.Fonts["Primary"],
								14, th.Colors["Dark"], t.Label),
						),
						layout.Flexed(0.2,
							component.Label(gtx, th, th.Fonts["Primary"],
								12, th.Colors["Dark"], fmt.Sprint(t.Amount)),
						),
						layout.Rigid(component.MonoButton(gtx, th,
							t.QrCode, 12, "", "",
							"Secondary", "QR",
							component.QrDialog(rc, gtx, t.Address)),
						),
					)
				}),
				layout.Rigid(th.DuoUIline(gtx, 1, 0,
					1, th.Colors["Gray"])),
			)
			//}
		})
		// }).Layout(gtx, addressBookPanel)
	}
}

func addressBookHeader(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme, pageFunc func()) func() {
	return func() {
		layout.Flex{
			Spacing:   layout.SpaceBetween,
			Axis:      layout.Horizontal,
			Alignment: layout.Middle,
		}.Layout(gtx,
			layout.Rigid(func() {
				if showMiningAddresses.Checked(gtx) {
					rc.AddressBook.ShowMiningAddresses = true
					// rc.GetAddressBook()()
				} else {
					// rc.GetAddressBook()()
				}
				th.DuoUIcheckBox("SHOW MINING ADDRESSES",
					th.Colors["Light"], th.Colors["Light"]).
					Layout(gtx, showMiningAddresses)
			}),
			layout.Rigid(func() {
				// th.DuoUIcounter(rc.GetBlocksExcerpts()).Layout(gtx,
				//rc.Explorer.Page, "PAGE", fmt.Sprint(rc.Explorer.Page.Value))
			}),
			// layout.Rigid(component.Button(gtx, th, buttonNewAddress,
			//th.Fonts["Secondary"], 12, th.Colors["ButtonText"], th.Colors["Dark"],
			//"NEW ADDRESS", component.QrDialog(rc, gtx, rc.CreateNewAddress("")))))
			layout.Rigid(component.MonoButton(gtx, th,
				buttonNewAddress, 12,
				"Primary", "Light", "Secondary",
				"NEW ADDRESS", func() {
					rc.Dialog.Show = true
					rc.Dialog = &model.DuoUIdialog{
						Show: true,
						Orange: func() {
							rc.Dialog.Show = false
						},
						CustomField: component.DuoUIqrCode(gtx, address, 256),
						Title:       "Copy address",
						Text:        rc.CreateNewAddress(""),
					}
					pageFunc()
				}),
			),
		)
	}
}

func blockPage(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme, block string) *gelook.DuoUIpage {
	page := gelook.DuoUIpage{
		Title:         "BLOCK",
		TxColor:       "",
		Command:       rc.GetSingleBlock(block),
		Border:        4,
		BorderColor:   th.Colors["Light"],
		Header:        func() {},
		HeaderBgColor: "",
		HeaderPadding: 0,
		Body:          singleBlockBody(rc, gtx, th, rc.Explorer.SingleBlock),
		BodyBgColor:   th.Colors["Light"],
		BodyPadding:   0,
		Footer:        func() {},
		FooterBgColor: "",
		FooterPadding: 0,
	}
	return th.DuoUIpage(page)
}

func blockRow(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme, block *model.DuoUIblock) {
	for block.Link.Clicked(gtx) {
		rc.ShowPage = fmt.Sprintf("BLOCK %s", block.BlockHash)
		rc.GetSingleBlock(block.BlockHash)()
		component.SetPage(rc, blockPage(rc, gtx, th, block.BlockHash))
	}
	width := gtx.Constraints.Width.Max
	button := th.DuoUIbutton("", "",
		"", "",
		"", "",
		"", "",
		0, 0, 0, 0,
		0, 0, 0, 0)
	button.InsideLayout(gtx, block.Link, func() {
		layout.Flex{
			Axis: layout.Vertical,
		}.Layout(gtx,
			layout.Rigid(func() {
				gtx.Constraints.Width.Min = width
				layout.Flex{
					Spacing: layout.SpaceBetween,
				}.Layout(gtx,
					layout.Rigid(func() {
						var linkButton gelook.DuoUIbutton
						linkButton = th.DuoUIbutton(th.Fonts["Mono"],
							fmt.Sprint(block.Height), th.Colors["Light"],
							th.Colors["Info"], th.Colors["Info"],
							th.Colors["Dark"], "", th.Colors["Dark"],
							14, 0, 60, 24,
							5, 8, 6, 8)
						linkButton.Layout(gtx, block.Link)
					}),
					layout.Rigid(func() {
						layout.UniformInset(unit.Dp(8)).Layout(gtx, func() {
							l := th.Body2(
								fmt.Sprint(time.Unix(block.Time, 0).
									Format("2006-01-02 15:04:05")))
							l.Font.Typeface = th.Fonts["Mono"]
							l.Alignment = text.Middle
							l.Color = th.Colors["Dark"]
							l.Layout(gtx)
						})
					}),
					layout.Rigid(func() {
						layout.UniformInset(unit.Dp(8)).Layout(gtx, func() {
							l := th.Body2(fmt.Sprint(block.Confirmations))
							l.Font.Typeface = th.Fonts["Mono"]
							l.Color = th.Colors["Dark"]
							l.Layout(gtx)
						})
					}),
					layout.Rigid(func() {
						gelook.DuoUIcontainer{}.Layout(gtx, layout.Center, func() {
							l := th.Body2(fmt.Sprint(block.TxNum))
							l.Font.Typeface = th.Fonts["Mono"]
							l.Color = th.Colors["Dark"]
							l.Layout(gtx)
						})
					}),
					layout.Rigid(func() {
						layout.UniformInset(unit.Dp(8)).Layout(gtx, func() {
							l := th.Body2(block.BlockHash)
							l.Font.Typeface = th.Fonts["Mono"]
							l.Color = th.Colors["Dark"]
							l.Layout(gtx)
							txwidth = gtx.Dimensions.Size.X
						})
					}),
				)
			}),
			layout.Rigid(func() {
				th.DuoUIline(gtx, 0, 0, 1,
					th.Colors["Gray"])()
			}),
		)
	})
}

func blockRowCellLabels(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) {
	cs := gtx.Constraints
	labels := th.DuoUIcontainer(0, th.Colors["Gray"])
	labels.FullWidth = true
	labels.Layout(gtx, layout.W, func() {
		// component.HorizontalLine(gtx, 1, th.Colors["Dark"])()
		gtx.Constraints.Width.Min = cs.Width.Max
		layout.Flex{
			Spacing: layout.SpaceBetween,
		}.Layout(gtx,
			layout.Rigid(func() {
				layout.UniformInset(unit.Dp(8)).Layout(gtx, func() {
					l := th.Body2("Height")
					l.Font.Typeface = th.Fonts["Mono"]
					l.Alignment = text.Middle
					l.Color = th.Colors["Dark"]
					l.Layout(gtx)
				})
			}),
			layout.Rigid(func() {
				layout.UniformInset(unit.Dp(8)).Layout(gtx, func() {
					l := th.Body2("Time")
					l.Font.Typeface = th.Fonts["Mono"]
					l.Alignment = text.Middle
					l.Color = th.Colors["Dark"]
					l.Layout(gtx)
				})
			}),
			layout.Rigid(func() {
				layout.UniformInset(unit.Dp(8)).Layout(gtx, func() {
					l := th.Body2("Confirmations")
					l.Font.Typeface = th.Fonts["Mono"]
					l.Color = th.Colors["Dark"]
					l.Layout(gtx)
				})
			}),
			layout.Rigid(func() {
				layout.UniformInset(unit.Dp(8)).Layout(gtx, func() {
					l := th.Body2("TxNum")
					l.Font.Typeface = th.Fonts["Mono"]
					l.Color = th.Colors["Dark"]
					l.Layout(gtx)
				})
			}),
			layout.Rigid(func() {
				layout.Inset{
					Top:   unit.Dp(8),
					Right: unit.Dp(float32(txwidth - 64)),
				}.Layout(gtx, func() {
					l := th.Body2("BlockHash")
					l.Font.Typeface = th.Fonts["Mono"]
					l.Color = th.Colors["Dark"]
					l.Layout(gtx)
				})
			}),
		)
	})
}

func bodyExplorer(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	return func() {
		rc.GetBlocksExcerpts()
		layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func() {
				blockRowCellLabels(rc, gtx, th)
			}),
			layout.Flexed(1, func() {
				explorerPanel := th.DuoUIpanel()
				explorerPanel.PanelObject = rc.Explorer.Blocks
				explorerPanel.ScrollBar = th.ScrollBar(16)
				explorerPanelElement.PanelObjectsNumber =
					len(rc.Explorer.Blocks)
				explorerPanel.Layout(gtx, explorerPanelElement,
					func(i int, in interface{}) {
						blocks := in.([]model.DuoUIblock)
						b := blocks[i]
						//blocksList.Layout(gtx, len(rc.Explorer.Blocks), func(i int) {
						//	b := rc.Explorer.Blocks[i]
						blockRow(rc, gtx, th, &b)
						//})
					},
				)
			}),
		)
	}
}

func Console(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) *gelook.DuoUIpage {
	page := gelook.DuoUIpage{
		Title:         "CONSOLE",
		Command:       func() {},
		Border:        4,
		BorderColor:   "ff000000",
		Header:        func() {},
		HeaderPadding: 0,
		Body:          consoleBody(rc, gtx, th),
		BodyBgColor:   "ff000000",
		BodyPadding:   4,
		Footer:        func() {},
		FooterPadding: 0,
	}
	return th.DuoUIpage(page)
}

func consoleBody(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	return func() {
		layout.UniformInset(unit.Dp(8)).Layout(gtx, func() {
			th.DuoUIcontainer(0, "ff000000").
				Layout(gtx, layout.N, func() {
					layout.Flex{}.Layout(gtx, layout.Flexed(1, func() {
						layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
							layout.Flex{
								Axis:    layout.Vertical,
								Spacing: layout.SpaceAround,
							}.Layout(gtx, layout.Flexed(1, func() {
								consoleOutputList.Layout(gtx,
									len(rc.ConsoleHistory.Commands), func(i int) {
										t := rc.ConsoleHistory.Commands[i]
										layout.Flex{
											Axis:      layout.Vertical,
											Alignment: layout.End,
										}.Layout(gtx,
											layout.Rigid(
												component.Label(
													gtx, th, th.Fonts["Mono"],
													12, th.Colors["Light"],
													"ds://"+t.ComID),
											),
											layout.Rigid(
												component.Label(
													gtx, th, th.Fonts["Mono"],
													12, th.Colors["Light"],
													t.Out),
											),
										)
									})
							}),
								layout.Rigid(
									component.ConsoleInput(gtx, th,
										consoleInputField, "Run command",
										func(e gel.SubmitEvent) {
											rc.ConsoleHistory.Commands = append(
												rc.ConsoleHistory.Commands,
												model.DuoUIconsoleCommand{
													ComID: e.Text,
													Time:  time.Time{},
													Out:   rc.ConsoleCmd(e.Text),
												})
										},
									),
								),
							)
						})
					}),
					)
				})
		})
	}
}

func DuoUIaddressBook(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) *gelook.DuoUIpage {
	page := gelook.DuoUIpage{
		Title:         "ADDRESSBOOK",
		Command:       rc.GetAddressBook(),
		Border:        4,
		BorderColor:   th.Colors["Light"],
		Header:        addressBookHeader(rc, gtx, th, rc.GetAddressBook()),
		HeaderPadding: 4,
		Body:          addressBookBody(rc, gtx, th),
		BodyBgColor:   th.Colors["Light"],
		BodyPadding:   4,
		Footer:        func() {},
		FooterPadding: 0,
	}
	return th.DuoUIpage(page)
}

func DuoUIexplorer(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) *gelook.DuoUIpage {
	page := gelook.DuoUIpage{
		Title:         "EXPLORER",
		TxColor:       "",
		Command:       rc.GetBlocksExcerpts(),
		Border:        4,
		BorderColor:   th.Colors["Light"],
		Header:        explorerHeader(rc, gtx, th),
		HeaderBgColor: "",
		HeaderPadding: 4,
		Body:          bodyExplorer(rc, gtx, th),
		BodyBgColor:   th.Colors["Light"],
		BodyPadding:   4,
		Footer:        func() {},
		FooterBgColor: "",
		FooterPadding: 0,
	}
	return th.DuoUIpage(page)
}

func DuoUIminer(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	return func() {
		rc.GetDuoUIhashesPerSecList()
		layout.Flex{}.Layout(gtx, layout.Flexed(1, func() {
			layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
				layout.Flex{
					Axis:    layout.Vertical,
					Spacing: layout.SpaceAround,
				}.Layout(gtx,
					layout.Flexed(1, func() {
						//	consoleOutputList.Layout(gtx, rc.Status.Kopach.Hps.Len(), func(i int) {
						//		t := rc.Status.Kopach.Hps.Get(i)
						//		layout.Flex{
						//			Axis:      layout.Vertical,
						//			Alignment: layout.End,
						//		}.Layout(gtx,
						//			layout.Rigid(func() {
						//				sat := th.Body1(fmt.Sprint(t))
						//				sat.Font.Typeface = th.Fonts["Mono"]
						//				sat.Color = gelook.HexARGB(th.Colors["Dark"])
						//				sat.Layout(gtx)
						//			}),
						//		)
						//	})
					}),
				)
			})
		}),
		)
	}
}

func explorerHeader(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	return func() {
		layout.Flex{
			Spacing:   layout.SpaceBetween,
			Axis:      layout.Horizontal,
			Alignment: layout.Middle,
		}.Layout(gtx,
			layout.Rigid(func() {
				th.DuoUIcounter(rc.GetBlocksExcerpts()).Layout(gtx,
					rc.Explorer.Page, "PAGE",
					fmt.Sprint(rc.Explorer.Page.Value))
			}),
			layout.Rigid(func() {
				th.DuoUIcounter(rc.GetBlocksExcerpts()).Layout(gtx,
					rc.Explorer.PerPage, "PER PAGE",
					fmt.Sprint(rc.Explorer.PerPage.Value))
			}),
		)
	}
}

func History(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) *gelook.DuoUIpage {
	page := gelook.DuoUIpage{
		Title:         "HISTORY",
		Command:       rc.GetDuoUItransactions(),
		Border:        4,
		BorderColor:   th.Colors["Light"],
		Header:        historyHeader(rc, gtx, th),
		HeaderPadding: 4,
		Body:          historyBody(rc, gtx, th),
		BodyBgColor:   th.Colors["Light"],
		BodyPadding:   0,
		Footer:        func() {},
		FooterPadding: 0,
	}
	return th.DuoUIpage(page)
}

func historyBody(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	return func() {
		layout.UniformInset(unit.Dp(16)).Layout(gtx, func() {
			layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				layout.Rigid(func() {
					transactionsPanel := th.DuoUIpanel()
					transactionsPanel.PanelObject = rc.History.Txs.Txs
					transactionsPanel.ScrollBar = th.ScrollBar(16)
					transactionsPanelElement.PanelObjectsNumber = len(rc.History.Txs.Txs)
					transactionsPanel.Layout(gtx, transactionsPanelElement,
						func(i int, in interface{}) {
							txs := in.([]model.DuoUItransactionExcerpt)
							t := txs[i]
							th.DuoUIline(gtx, 0, 0, 1, th.Colors["Hint"])()
							for t.Link.Clicked(gtx) {
								rc.ShowPage = fmt.Sprintf("TRANSACTION %s", t.TxID)
								rc.GetSingleTx(t.TxID)()
								component.SetPage(rc, txPage(rc, gtx, th, t.TxID))
							}
							width := gtx.Constraints.Width.Max
							button := th.DuoUIbutton("", "",
								"", "", "", "",
								"", "", 0, 0,
								0, 0, 0, 0,
								0, 0)
							button.InsideLayout(gtx, t.Link, func() {
								gtx.Constraints.Width.Min = width
								layout.Flex{
									Spacing: layout.SpaceBetween,
								}.Layout(gtx,
									layout.Rigid(component.TxsDetails(gtx,
										th, i, &t)),
									layout.Rigid(component.Label(gtx,
										th, th.Fonts["Mono"], 12,
										th.Colors["Secondary"],
										fmt.Sprintf("%0.8f", t.Amount))))
							})
						},
					)
				}),
			)
		})
	}
}

func historyHeader(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	return func() {
		layout.Flex{
			Spacing: layout.SpaceBetween,
		}.Layout(gtx,
			layout.Rigid(component.TransactionsFilter(rc, gtx, th)),
			layout.Rigid(func() {
				th.DuoUIcounter(rc.GetDuoUItransactions()).
					Layout(gtx, rc.History.PerPage, "TxNum per page: ",
						fmt.Sprint(rc.History.PerPage.Value))
			}),
			layout.Rigid(func() {
				th.DuoUIcounter(rc.GetDuoUItransactions()).
					Layout(gtx, rc.History.Page, "TxNum page: ",
						fmt.Sprint(rc.History.Page.Value))
			}),
		)
	}
}

func LoadPages(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) (p map[string]*gelook.DuoUIpage) {
	p = make(map[string]*gelook.DuoUIpage)
	p["OVERVIEW"] = Overview(rc, gtx, th)
	p["SEND"] = Send(rc, gtx, th)
	// p["RECEIVE"] = th.DuoUIpage("RECEIVE", 10, func() {},
	// func() {}, func() { th.H5("receive :").Layout(gtx) }, func() {})
	p["ADDRESSBOOK"] = DuoUIaddressBook(rc, gtx, th)
	p["SETTINGS"] = Settings(rc, gtx, th)
	p["NETWORK"] = Network(rc, gtx, th)
	// p["BLOCK"] = th.DuoUIpage("BLOCK", 0, func() {},
	// func() {}, func() { th.H5("block :").Layout(gtx) }, func() {})
	p["HISTORY"] = History(rc, gtx, th)
	p["EXPLORER"] = DuoUIexplorer(rc, gtx, th)
	p["MINER"] = Miner(rc, gtx, th)
	p["CONSOLE"] = Console(rc, gtx, th)
	p["LOG"] = Logger(rc, gtx, th)
	return
}

func Logger(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) *gelook.DuoUIpage {
	page := gelook.DuoUIpage{
		Title:         "LOG",
		Command:       func() {},
		Border:        4,
		BorderColor:   th.Colors["Dark"],
		Header:        func() {},
		HeaderBgColor: "",
		Body:          component.DuoUIlogger(rc, gtx, th),
		BodyBgColor:   th.Colors["Dark"],
		Footer:        func() {},
		FooterBgColor: "",
	}
	return th.DuoUIpage(page)
}

func Miner(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) *gelook.DuoUIpage {
	page := gelook.DuoUIpage{
		Title:         "MINER",
		TxColor:       "",
		Command:       func() {},
		Border:        4,
		BorderColor:   th.Colors["Light"],
		Header:        func() {},
		HeaderBgColor: "",
		HeaderPadding: 0,
		Body:          DuoUIminer(rc, gtx, th),
		BodyBgColor:   th.Colors["Dark"],
		BodyPadding:   0,
		Footer:        func() {},
		FooterBgColor: "",
		FooterPadding: 0,
	}
	return th.DuoUIpage(page)
}

func Network(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) *gelook.DuoUIpage {
	page := gelook.DuoUIpage{
		Title:         "NETWORK",
		TxColor:       "",
		Command:       rc.GetPeerInfo(),
		Border:        4,
		BorderColor:   th.Colors["Light"],
		Header:        networkHeader(rc, gtx, th),
		HeaderBgColor: "",
		HeaderPadding: 0,
		Body:          networkBody(rc, gtx, th),
		BodyBgColor:   th.Colors["Light"],
		BodyPadding:   0,
		Footer:        func() {},
		FooterBgColor: "",
		FooterPadding: 0,
	}
	return th.DuoUIpage(page)
}

func networkBody(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	return func() {
		layout.UniformInset(unit.Dp(16)).Layout(gtx, func() {
			layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				layout.Rigid(component.PeersList(rc, gtx, th)))
		})
	}
}

func networkHeader(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	return func() {
		layout.Flex{
			Spacing: layout.SpaceBetween,
		}.Layout(gtx,
			// layout.Rigid(component.TransactionsFilter(rc, gtx, th)),
			layout.Rigid(func() {
				th.DuoUIcounter(rc.GetPeerInfo()).Layout(gtx,
					rc.Network.PerPage, "Peers per page: ",
					fmt.Sprint(rc.Network.PerPage.Value))
			}),
			layout.Rigid(func() {
				th.DuoUIcounter(rc.GetPeerInfo()).Layout(gtx,
					rc.Network.Page, "Peers page: ",
					fmt.Sprint(rc.Network.Page.Value))
			}),
		)
	}
}

func Overview(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) *gelook.DuoUIpage {
	page := gelook.DuoUIpage{
		Title:         "OVERVIEW",
		Border:        0,
		Command:       func() {},
		Header:        func() {},
		HeaderBgColor: "",
		HeaderPadding: 0,
		Body:          overviewBody(rc, gtx, th),
		BodyBgColor:   th.Colors["Light"],
		BodyPadding:   0,
		Footer:        func() {},
		FooterBgColor: "",
		FooterPadding: 0,
	}
	return th.DuoUIpage(page)
}

func overviewBody(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	return func() {
		viewport := layout.Flex{Axis: layout.Horizontal}
		if gtx.Constraints.Width.Max < 780 {
			viewport = layout.Flex{Axis: layout.Vertical}
		}
		viewport.Layout(gtx,
			layout.Flexed(0.5, component.DuoUIstatus(rc, gtx, th)),
			layout.Flexed(0.5, component.DuoUIlatestTransactions(rc, gtx, th)),
		)
		op.InvalidateOp{}.Add(gtx.Ops)
	}
}

func Send(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) *gelook.DuoUIpage {
	page := gelook.DuoUIpage{
		Title:         "SEND",
		TxColor:       "",
		Command:       func() {},
		Border:        4,
		BorderColor:   th.Colors["Light"],
		Header:        func() {},
		HeaderBgColor: "",
		HeaderPadding: 0,
		Body:          sendBody(rc, gtx, th),
		BodyBgColor:   th.Colors["Light"],
		BodyPadding:   0,
		Footer:        func() {},
		FooterBgColor: "",
		FooterPadding: 0,
	}
	return th.DuoUIpage(page)
}

func sendBody(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	return func() {
		layout.Flex{}.Layout(gtx,
			layout.Rigid(func() {
				widgets := []func(){
					func() {
						th.DuoUIcontainer(1,
							th.Colors["Gray"]).Layout(gtx, layout.Center, func() {
							layout.Flex{}.Layout(gtx,
								layout.Flexed(1, component.Editor(gtx, th,
									addressLineEditor, "DUO address",
									func(e gel.EditorEvent) {
										sendStruct.address = addressLineEditor.Text()
									})),
								layout.Rigid(component.Button(gtx, th,
									buttonPasteAddress, th.Fonts["Primary"],
									10, 13, 8, 12, 8,
									th.Colors["ButtonText"], th.Colors["ButtonBg"],
									"PASTE ADDRESS", func() {
										addressLineEditor.SetText(clipboard.Get())
									})))
						})
					},
					func() {
						th.DuoUIcontainer(1, th.Colors["Gray"]).Layout(gtx,
							layout.Center, func() {
								layout.Flex{}.Layout(gtx,
									layout.Flexed(1, component.Editor(gtx,
										th, amountLineEditor, "DUO Amount",
										func(e gel.EditorEvent) {
											f, err := strconv.ParseFloat(amountLineEditor.Text(), 64)
											if err != nil {
											}
											sendStruct.amount = f
										}),
									),
									layout.Rigid(component.Button(gtx, th,
										buttonPasteAmount, th.Fonts["Primary"],
										10, 13, 8, 12, 8,
										th.Colors["ButtonText"],
										th.Colors["ButtonBg"],
										"PASTE AMOUNT",
										func() {
											amountLineEditor.SetText(clipboard.Get())
										}),
									),
								)
							})
					},
					func() {
						layout.Flex{}.Layout(gtx,
							layout.Rigid(component.Button(gtx, th,
								buttonSend, th.Fonts["Primary"],
								14, 10, 10, 9, 10,
								th.Colors["ButtonText"], th.Colors["ButtonBg"],
								"SEND", func() {
									rc.Dialog.Show = true
									rc.Dialog = &model.DuoUIdialog{
										Show: true,
										Green: rc.DuoSend(sendStruct.passPhrase,
											sendStruct.address, 11),
										GreenLabel: "SEND",
										CustomField: func() {
											layout.Flex{}.Layout(gtx,
												layout.Flexed(1,
													component.Editor(gtx, th, passLineEditor,
														"Enter your password",
														func(e gel.EditorEvent) {
															sendStruct.passPhrase = passLineEditor.Text()
														},
													),
												),
											)
										},
										Red:      func() { rc.Dialog.Show = false },
										RedLabel: "CANCEL",
										Title:    "Are you sure?",
										Text:     "Confirm ParallelCoin send",
									}
								}),
							),
						)
					},
				}
				layautList.Layout(gtx, len(widgets), func(i int) {
					layout.UniformInset(unit.Dp(8)).Layout(gtx, widgets[i])
				})
			}))
		//Info("passPhrase:" + sendStruct.passPhrase)
	}
}

func Settings(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) *gelook.DuoUIpage {
	page := gelook.DuoUIpage{
		Title:         "SETTINGS",
		TxColor:       "",
		Command:       func() {},
		Border:        4,
		BorderColor:   th.Colors["Light"],
		Header:        SettingsHeader(rc, gtx, th),
		HeaderBgColor: "",
		HeaderPadding: 4,
		Body:          SettingsBody(rc, gtx, th),
		BodyBgColor:   th.Colors["Light"],
		BodyPadding:   0,
		Footer:        func() {},
		FooterBgColor: "",
		FooterPadding: 0,
	}
	return th.DuoUIpage(page)

	// return th.DuoUIpage("SETTINGS", 0, func() {}, component.ContentHeader(gtx, th, SettingsHeader(rc, gtx, th)), SettingsBody(rc, gtx, th), func() {
	// var msg string
	// if rc.Settings.Daemon.Config["DisableBanning"].(*bool) != true{
	//	msg = "ima"
	// }else{
	//	msg = "nema"
	// //}
	// ttt := th.H6(fmt.Sprint(rc.Settings.Daemon.Config))
	// ttt.Color = gelook.HexARGB("ffcfcfcf")
	// ttt.Layout(gtx)
	// })
}

func SettingsBody(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	return func() {
		th.DuoUIcontainer(16,
			th.Colors["Light"]).Layout(gtx, layout.N, func() {
			for _, fields := range rc.Settings.Daemon.Schema.Groups {
				if fmt.Sprint(fields.Legend) == rc.Settings.Tabs.Current {
					settingsPanel := th.DuoUIpanel()
					settingsPanel.PanelObject = fields.Fields
					settingsPanel.ScrollBar = th.ScrollBar(16)
					settingsPanelElement.PanelObjectsNumber = len(fields.Fields)
					settingsPanel.Layout(gtx, settingsPanelElement, func(i int, in interface{}) {
						settings := in.(pod.Fields)
						//t := settings[i]
						//fieldsList.Layout(gtx, len(fields.Fields), func(il int) {
						i = settingsPanelElement.PanelObjectsNumber - 1 - i
						tl := component.Field{
							Field: &settings[i],
						}
						layout.Flex{
							Axis: layout.Vertical,
						}.Layout(gtx,
							layout.Rigid(SettingsItemRow(rc, gtx, th, &tl)),
							layout.Rigid(th.DuoUIline(gtx,
								4, 0, 1, th.Colors["LightGray"])))
					})
				}
			}
		})
	}
}

func SettingsHeader(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	return func() {
		layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
			layout.Rigid(component.SettingsTabs(rc, gtx, th)),
			layout.Rigid(func() {
				// var settingsRestartButton gelook.DuoUIbutton
				// settingsRestartButton = th.DuoUIbutton(th.Fonts["Secondary"],
				// 	"restart",
				// 	th.Colors["Light"],
				// 	th.Colors["Dark"],
				// 	th.Colors["Dark"],
				// 	th.Colors["Light"],
				// 	"",
				// 	th.Colors["Light"],
				// 	23, 0, 80, 48, 4, 4)
				// for buttonSettingsRestart.Clicked(gtx) {
				// 	rc.SaveDaemonCfg()
				// }
				// settingsRestartButton.Layout(gtx, buttonSettingsRestart)
			}),
		)
	}
}

func SettingsItemRow(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme, f *component.Field) func() {
	return func() {
		layout.Flex{
			Axis:      layout.Horizontal,
			Alignment: layout.Middle,
		}.Layout(gtx,
			// layout.Rigid(func() {
			//	gelook.DuoUIdrawRectangle(gtx, 30, 3, th.Colors["Light"],
			//	[4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
			// }),
			layout.Flexed(0.62, func() {
				layout.Flex{
					Axis:    layout.Vertical,
					Spacing: 4,
				}.Layout(gtx,
					layout.Rigid(component.SettingsFieldLabel(gtx, th, f)),
					layout.Rigid(component.SettingsFieldDescription(gtx, th, f)),
				)
			}),
			layout.Flexed(1, component.DuoUIinputField(rc, gtx, th, f)),
		)
	}
}

func singleBlockBody(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme, block btcjson.GetBlockVerboseResult) func() {
	return func() {
		switch block.PowAlgo {
		case "scrypt":
			algoHeadColor = "Light"
			algoHeadBgColor = "Dark"
			algoColor = "DarkGrayI"
			algoBgColor = "Warning"
		case "sha256d":
			algoHeadColor = "Light"
			algoHeadBgColor = "Dark"
			algoColor = "LightGrayI"
			algoBgColor = "Info"
		case fork.List[1].AlgoVers[5]:
			algoHeadColor = "Light"
			algoHeadBgColor = "Dark"
			algoColor = "Light"
			algoBgColor = "DarkGray"
		case fork.List[1].AlgoVers[6]:
			algoHeadColor = "Light"
			algoHeadBgColor = "Dark"
			algoColor = "Dark"
			algoBgColor = "LightGrayI"
		case fork.List[1].AlgoVers[7]:
			algoHeadColor = "Light"
			algoHeadBgColor = "Dark"
			algoColor = "Light"
			algoBgColor = "DarkGrayI"
		case fork.List[1].AlgoVers[8]:
			algoHeadColor = "Light"
			algoHeadBgColor = "Dark"
			algoColor = "Light"
			algoBgColor = "Secondary"
		case fork.List[1].AlgoVers[9]:
			algoHeadColor = "Light"
			algoHeadBgColor = "Dark"
			algoColor = "Dark"
			algoBgColor = "Success"
		case fork.List[1].AlgoVers[10]:
			algoHeadColor = "Light"
			algoHeadBgColor = "Dark"
			algoColor = "Light"
			algoBgColor = "DarkGrayII"
		case fork.List[1].AlgoVers[11]:
			algoHeadColor = "Light"
			algoHeadBgColor = "Dark"
			algoColor = "Light"
			algoBgColor = "Danger"
		case fork.List[1].AlgoVers[12]:
			algoHeadColor = "Light"
			algoHeadBgColor = "Dark"
			algoColor = "Light"
			algoBgColor = "Hint"
		case fork.List[1].AlgoVers[13]:
			algoHeadColor = "Light"
			algoHeadBgColor = "Dark"
			algoColor = "Light"
			algoBgColor = "Fatal"
		}
		duo := layout.Horizontal
		if gtx.Constraints.Width.Max < 1280 {
			duo = layout.Vertical
		}
		trio := layout.Horizontal
		if gtx.Constraints.Width.Max < 780 {
			trio = layout.Vertical
		}
		blockJSON, _ := json.MarshalIndent(block, "", "  ")
		blockText := string(blockJSON)
		widgets := []func(){
			component.UnoField(gtx, component.ContentLabeledField(gtx, th,
				layout.Vertical, 4, 12, 14,
				"Hash", "Dark", "LightGrayII",
				"Dark", "LightGrayI", fmt.Sprint(block.Hash))),
			component.DuoFields(gtx, duo,
				component.TrioFields(gtx, th, trio, 12, 16,
					"Height", fmt.Sprint(block.Height),
					"Dark", "LightGrayII",
					"Dark", "LightGrayI",
					"Confirmations", fmt.Sprint(block.Confirmations),
					"Dark", "LightGrayII",
					"Dark", "LightGrayI",
					"Time", fmt.Sprint(time.Unix(block.Time, 0).
						Format("2006-01-02 15:04:05")),
					"LightGrayII", "Dark",
					"Dark", "LightGrayI",
				),
				component.TrioFields(gtx, th, trio, 12, 16,
					"PowAlgo", fmt.Sprint(block.PowAlgo),
					algoHeadColor, algoHeadBgColor, algoColor, algoBgColor,
					"Difficulty", fmt.Sprint(block.Difficulty),
					"Dark", "LightGrayII",
					"Dark", "LightGrayI",
					"Nonce", fmt.Sprint(block.Nonce),
					"LightGrayII", "Dark",
					"Dark", "LightGrayI",
				),
			),
			component.DuoFields(gtx, duo,
				component.ContentLabeledField(gtx, th, layout.Vertical,
					4, 12, 12,
					"MerkleRoot", "Dark", "LightGrayII",
					"Dark", "LightGrayI", block.MerkleRoot),
				component.ContentLabeledField(gtx, th, layout.Vertical,
					4, 12, 12,
					"PowHash", "Dark", "LightGrayII",
					"Dark", "LightGrayI", fmt.Sprint(block.PowHash)),
			),
			component.DuoFields(gtx, duo,
				component.TrioFields(gtx, th, trio, 12, 16,
					"Size", fmt.Sprint(block.Size),
					"Dark", "LightGrayII",
					"Dark", "LightGrayI",
					"Weight", fmt.Sprint(block.Weight),
					"Dark", "LightGrayII",
					"Dark", "LightGrayI",
					"Bits", fmt.Sprint(block.Bits),
					"LightGrayII", "Dark",
					"Dark", "LightGrayI",
				),
				component.TrioFields(gtx, th, trio, 12, 16,
					"TxNum", fmt.Sprint(block.TxNum),
					"Dark", "LightGrayII",
					"Dark", "LightGrayI",
					"StrippedSize", fmt.Sprint(block.StrippedSize),
					"Dark", "LightGrayII",
					"Dark", "LightGrayI",
					"Version", fmt.Sprint(block.Version),
					"LightGrayII", "Dark",
					"Dark", "LightGrayI",
				),
			),
			component.UnoField(gtx, component.ContentLabeledField(gtx, th,
				layout.Vertical, 4, 12, 12,
				"Tx", "Dark", "LightGrayII",
				"Dark", "LightGrayI", fmt.Sprint(block.Tx))),
			component.UnoField(gtx, component.ContentLabeledField(gtx, th,
				layout.Vertical, 4, 12, 12,
				"RawTx", "Dark", "LightGrayII",
				"Dark", "LightGrayI", fmt.Sprint(blockText))),
			component.PageNavButtons(rc, gtx, th, block.PreviousHash,
				block.NextHash, blockPage(rc, gtx, th, block.PreviousHash),
				blockPage(rc, gtx, th, block.NextHash)),
		}
		layautList.Layout(gtx, len(widgets), func(i int) {
			layout.UniformInset(unit.Dp(0)).Layout(gtx, widgets[i])
		})
	}
}

func singleTxBody(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme, tx btcjson.GetTransactionResult) func() {
	return func() {

		//duo := layout.Horizontal
		//if gtx.Constraints.Width.Max < 1280 {
		//	duo = layout.Vertical
		//}
		//trio := layout.Horizontal
		//if gtx.Constraints.Width.Max < 780 {
		//	trio = layout.Vertical
		//}

		//blockJSON, _ := json.MarshalIndent(block, "", "  ")
		//blockText := string(blockJSON)
		widgets := []func(){

			func() {
				th.H6(tx.TxID).Layout(gtx)
			},
			//component.UnoField(gtx, component.ContentLabeledField(gtx, th,
			//layout.Vertical, 4, 12, 14, "Hash", "Dark", "LightGrayII", "Dark", "LightGrayI", fmt.Sprint(block.Hash))),
			//component.DuoFields(gtx, duo,
			//	component.TrioFields(gtx, th, trio, 12, 16,
			//		"Height", fmt.Sprint(block.Height), "Dark", "LightGrayII",
			//		"Dark", "LightGrayI",
			//		"Confirmations", fmt.Sprint(block.Confirmations), "Dark",
			//		"LightGrayII", "Dark", "LightGrayI",
			//		"Time", fmt.Sprint(time.Unix(block.Time, 0).Format("2006-01-02 15:04:05")),
			//		"LightGrayII", "Dark", "Dark", "LightGrayI",
			//	),
			//	component.TrioFields(gtx, th, trio, 12, 16,
			//		"PowAlgo", fmt.Sprint(block.PowAlgo), algoHeadColor,
			//		algoHeadBgColor, algoColor, algoBgColor,
			//		"Difficulty", fmt.Sprint(block.Difficulty), "Dark",
			//		"LightGrayII", "Dark", "LightGrayI",
			//		"Nonce", fmt.Sprint(block.Nonce), "LightGrayII", "Dark",
			//		"Dark", "LightGrayI",
			//	),
			//),
			//component.DuoFields(gtx, duo,
			//	component.ContentLabeledField(gtx, th, layout.Vertical,
			//	4, 12, 12, "MerkleRoot", "Dark", "LightGrayII", "Dark",
			//	"LightGrayI", block.MerkleRoot),
			//	component.ContentLabeledField(gtx, th, layout.Vertical,
			//	4, 12, 12, "PowHash", "Dark", "LightGrayII", "Dark",
			//	"LightGrayI", fmt.Sprint(block.PowHash)),
			//),
			//component.DuoFields(gtx, duo,
			//	component.TrioFields(gtx, th, trio, 12, 16,
			//		"Size", fmt.Sprint(block.Size), "Dark", "LightGrayII",
			//		"Dark", "LightGrayI",
			//		"Weight", fmt.Sprint(block.Weight), "Dark", "LightGrayII",
			//		"Dark", "LightGrayI",
			//		"Bits", fmt.Sprint(block.Bits), "LightGrayII", "Dark",
			//		"Dark", "LightGrayI",
			//	),
			//	component.TrioFields(gtx, th, trio, 12, 16,
			//		"TxNum", fmt.Sprint(block.TxNum), "Dark", "LightGrayII",
			//		"Dark", "LightGrayI",
			//		"StrippedSize", fmt.Sprint(block.StrippedSize), "Dark",
			//		"LightGrayII", "Dark", "LightGrayI",
			//		"Version", fmt.Sprint(block.Version), "LightGrayII",
			//		"Dark", "Dark", "LightGrayI",
			//	),
			//),
			//component.UnoField(gtx, component.ContentLabeledField(gtx, th,
			//layout.Vertical, 4, 12, 12, "Tx", "Dark", "LightGrayII", "Dark",
			//"LightGrayI", fmt.Sprint(block.Tx))),
			//component.UnoField(gtx, component.ContentLabeledField(gtx, th,
			//layout.Vertical, 4, 12, 12, "RawTx", "Dark", "LightGrayII",
			//"Dark", "LightGrayI", fmt.Sprint(blockText))),
			//component.PageNavButtons(rc, gtx, th, block.PreviousHash,
			//block.NextHash, blockPage(rc, gtx, th, block.PreviousHash),
			//blockPage(rc, gtx, th, block.NextHash)),
		}
		layautList.Layout(gtx, len(widgets), func(i int) {
			layout.UniformInset(unit.Dp(0)).Layout(gtx, widgets[i])
		})
	}
}
func txPage(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme, tx string) *gelook.DuoUIpage {
	page := gelook.DuoUIpage{
		Title:         "BLOCK",
		TxColor:       "",
		Command:       rc.GetSingleTx(tx),
		Border:        4,
		BorderColor:   th.Colors["Light"],
		Header:        func() {},
		HeaderBgColor: "",
		HeaderPadding: 0,
		Body:          singleTxBody(rc, gtx, th, rc.History.SingleTx),
		BodyBgColor:   th.Colors["Light"],
		BodyPadding:   0,
		Footer:        func() {},
		FooterBgColor: "",
		FooterPadding: 0,
	}
	return th.DuoUIpage(page)
}
