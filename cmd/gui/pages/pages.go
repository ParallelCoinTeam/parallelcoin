package pages

import (
	"encoding/json"
	"fmt"
	"gioui.org/op"
	"gioui.org/text"
	"github.com/stalker-loki/app/slog"
	"github.com/stalker-loki/pod/pkg/chain/fork"
	"github.com/stalker-loki/pod/pkg/pod"
	"github.com/stalker-loki/pod/pkg/rpc/btcjson"
	"runtime/debug"
	"strconv"
	"time"

	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/stalker-loki/pod/cmd/gui/component"
	"github.com/stalker-loki/pod/cmd/gui/model"
	"github.com/stalker-loki/pod/pkg/gui/clipboard"
	"github.com/stalker-loki/pod/pkg/gui/gel"
	"github.com/stalker-loki/pod/pkg/gui/gelook"
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
	txsPanelElement   = gel.NewPanel()
	consoleInputField = &gel.Editor{
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

func addressBookBody(c *component.State) func() {
	return func() {
		layout.Flex{}.Layout(c.Gtx,
			layout.Flexed(1, func() {
				layout.UniformInset(unit.Dp(0)).Layout(c.Gtx, func() {
					layout.Flex{
						Axis:    layout.Vertical,
						Spacing: layout.SpaceAround,
					}.Layout(c.Gtx,
						layout.Flexed(1, addressBookContent(c)))
				})
			}))
	}
}

func addressBookContent(c *component.State) func() {
	return func() {
		addressBookPanelElement.PanelObject = c.Rc.AddressBook.Addresses
		addressBookPanelElement.PanelObjectsNumber = len(c.Rc.AddressBook.Addresses)
		addressBookPanel := c.Thm.DuoUIpanel()
		addressBookPanel.ScrollBar = c.Thm.ScrollBar(16)
		addressBookPanel.Layout(c.Gtx, addressBookPanelElement, func(i int, in interface{}) {
			//if in != nil {
			//addresses := in.([]model.DuoUIaddress)
			t := c.Rc.AddressBook.Addresses[i]
			layout.Flex{Axis: layout.Vertical}.Layout(c.Gtx,
				layout.Rigid(func() {
					layout.Flex{
						Alignment: layout.Middle,
					}.Layout(c.Gtx,
						layout.Flexed(0.2,
							c.Label(c.Thm.Fonts["Primary"],
								12, c.Thm.Colors["Dark"], fmt.Sprint(t.Index)),
						),
						layout.Flexed(0.2,
							c.Label(c.Thm.Fonts["Primary"], 12,
								c.Thm.Colors["Dark"], t.Account),
						),
						layout.Rigid(c.MonoButton(
							t.Copy, 12, "", "",
							"Mono", t.Address, func() {
								clipboard.Set(t.Address)
							}),
						),
						layout.Flexed(0.4,
							c.Label(c.Thm.Fonts["Primary"],
								14, c.Thm.Colors["Dark"], t.Label),
						),
						layout.Flexed(0.2,
							c.Label(c.Thm.Fonts["Primary"],
								12, c.Thm.Colors["Dark"], fmt.Sprint(t.Amount)),
						),
						layout.Rigid(c.MonoButton(
							t.QrCode, 12, "", "",
							"Secondary", "QR",
							component.QrDialog(c.Rc, c.Gtx, t.Address)),
						),
					)
				}),
				layout.Rigid(c.Thm.DuoUIline(c.Gtx, 1, 0,
					1, c.Thm.Colors["Gray"])),
			)
			//}
		})
		// }).Layout(c.Gtx, addressBookPanel)
	}
}

func addressBookHeader(c *component.State, pageFunc func()) func() {
	return func() {
		layout.Flex{
			Spacing:   layout.SpaceBetween,
			Axis:      layout.Horizontal,
			Alignment: layout.Middle,
		}.Layout(c.Gtx,
			layout.Rigid(func() {
				if showMiningAddresses.Checked(c.Gtx) {
					c.Rc.AddressBook.ShowMiningAddresses = true
					// c.Rc.GetAddressBook()()
				} else {
					// c.Rc.GetAddressBook()()
				}
				c.Thm.DuoUIcheckBox("SHOW MINING ADDRESSES",
					c.Thm.Colors["Light"], c.Thm.Colors["Light"]).
					Layout(c.Gtx, showMiningAddresses)
			}),
			layout.Rigid(func() {
				// c.Thm.DuoUIcounter(c.Rc.GetBlocksExcerpts()).Layout(c.Gtx,
				//c.Rc.Explorer.Page, "PAGE", fmt.Sprint(c.Rc.Explorer.Page.Value))
			}),
			// layout.Rigid(component.Button(c.Gtx, th, buttonNewAddress,
			//c.Thm.Fonts["Secondary"], 12, c.Thm.Colors["ButtonText"], c.Thm.Colors["Dark"],
			//"NEW ADDRESS", component.QrDialog(rc, c.Gtx, c.Rc.CreateNewAddress("")))))
			layout.Rigid(c.MonoButton(
				buttonNewAddress, 12,
				"Primary", "Light", "Secondary",
				"NEW ADDRESS", func() {
					c.Rc.Dialog.Show = true
					c.Rc.Dialog = &model.DuoUIdialog{
						Show: true,
						Orange: func() {
							c.Rc.Dialog.Show = false
						},
						CustomField: component.DuoUIqrCode(c.Gtx, address, 256),
						Title:       "Copy address",
						Text:        c.Rc.CreateNewAddress(""),
					}
					pageFunc()
				}),
			),
		)
	}
}

func blockPage(c *component.State, block string) *gelook.DuoUIpage {
	page := gelook.DuoUIpage{
		Title:         "BLOCK",
		TxColor:       "",
		Command:       c.Rc.GetSingleBlock(block),
		Border:        4,
		BorderColor:   c.Thm.Colors["Light"],
		Header:        func() {},
		HeaderBgColor: "",
		HeaderPadding: 0,
		Body:          singleBlockBody(c, c.Rc.Explorer.SingleBlock),
		BodyBgColor:   c.Thm.Colors["Light"],
		BodyPadding:   0,
		Footer:        func() {},
		FooterBgColor: "",
		FooterPadding: 0,
	}
	return c.Thm.DuoUIpage(page)
}

func blockRow(c *component.State, block *model.DuoUIblock) {
	if block == nil || block.Link == nil {
		slog.Debug("blockRow empty result")
		debug.PrintStack()
		return
	}
	for block.Link.Clicked(c.Gtx) {
		c.Rc.ShowPage = fmt.Sprintf("BLOCK %s", block.BlockHash)
		c.Rc.GetSingleBlock(block.BlockHash)()
		component.SetPage(c.Rc, blockPage(c, block.BlockHash))
	}
	width := c.Gtx.Constraints.Width.Max
	button := c.Thm.DuoUIbutton(gelook.ButtonParams{})
	button.InsideLayout(c.Gtx, block.Link, func() {
		layout.Flex{
			Axis: layout.Vertical,
		}.Layout(c.Gtx,
			layout.Rigid(func() {
				c.Gtx.Constraints.Width.Min = width
				layout.Flex{
					Spacing: layout.SpaceBetween,
				}.Layout(c.Gtx,
					layout.Rigid(func() {
						var linkButton gelook.DuoUIbutton
						linkButton = c.Thm.DuoUIbutton(gelook.ButtonParams{
							TxtFont:       c.Thm.Fonts["Mono"],
							Txt:           fmt.Sprint(block.Height),
							TxtColor:      c.Thm.Colors["Light"],
							BgColor:       c.Thm.Colors["Info"],
							TxtHoverColor: c.Thm.Colors["Info"],
							BgHoverColor:  c.Thm.Colors["Dark"],
							IconColor:     c.Thm.Colors["Dark"],
							TextSize:      14,
							Width:         60,
							Height:        24,
							PaddingTop:    5,
							PaddingRight:  8,
							PaddingBottom: 6,
							PaddingLeft:   8,
						})
						linkButton.Layout(c.Gtx, block.Link)
					}),
					layout.Rigid(func() {
						layout.UniformInset(unit.Dp(8)).Layout(c.Gtx, func() {
							l := c.Thm.Body2(
								fmt.Sprint(time.Unix(block.Time, 0).
									Format("2006-01-02 15:04:05")))
							l.Font.Typeface = c.Thm.Fonts["Mono"]
							l.Alignment = text.Middle
							l.Color = c.Thm.Colors["Dark"]
							l.Layout(c.Gtx)
						})
					}),
					layout.Rigid(func() {
						layout.UniformInset(unit.Dp(8)).Layout(c.Gtx, func() {
							l := c.Thm.Body2(fmt.Sprint(block.Confirmations))
							l.Font.Typeface = c.Thm.Fonts["Mono"]
							l.Color = c.Thm.Colors["Dark"]
							l.Layout(c.Gtx)
						})
					}),
					layout.Rigid(func() {
						gelook.DuoUIcontainer{}.Layout(c.Gtx, layout.Center, func() {
							l := c.Thm.Body2(fmt.Sprint(block.TxNum))
							l.Font.Typeface = c.Thm.Fonts["Mono"]
							l.Color = c.Thm.Colors["Dark"]
							l.Layout(c.Gtx)
						})
					}),
					layout.Rigid(func() {
						layout.UniformInset(unit.Dp(8)).Layout(c.Gtx, func() {
							l := c.Thm.Body2(block.BlockHash)
							l.Font.Typeface = c.Thm.Fonts["Mono"]
							l.Color = c.Thm.Colors["Dark"]
							l.Layout(c.Gtx)
							txwidth = c.Gtx.Dimensions.Size.X
						})
					}),
				)
			}),
			layout.Rigid(func() {
				c.Thm.DuoUIline(c.Gtx, 0, 0, 1,
					c.Thm.Colors["Gray"])()
			}),
		)
	})
}

func blockRowCellLabels(c *component.State) {
	cs := c.Gtx.Constraints
	labels := c.Thm.DuoUIcontainer(0, c.Thm.Colors["Gray"])
	labels.FullWidth = true
	labels.Layout(c.Gtx, layout.W, func() {
		// component.HorizontalLine(c.Gtx, 1, c.Thm.Colors["Dark"])()
		c.Gtx.Constraints.Width.Min = cs.Width.Max
		layout.Flex{
			Spacing: layout.SpaceBetween,
		}.Layout(c.Gtx,
			layout.Rigid(func() {
				layout.UniformInset(unit.Dp(8)).Layout(c.Gtx, func() {
					l := c.Thm.Body2("Height")
					l.Font.Typeface = c.Thm.Fonts["Mono"]
					l.Alignment = text.Middle
					l.Color = c.Thm.Colors["Dark"]
					l.Layout(c.Gtx)
				})
			}),
			layout.Rigid(func() {
				layout.UniformInset(unit.Dp(8)).Layout(c.Gtx, func() {
					l := c.Thm.Body2("Time")
					l.Font.Typeface = c.Thm.Fonts["Mono"]
					l.Alignment = text.Middle
					l.Color = c.Thm.Colors["Dark"]
					l.Layout(c.Gtx)
				})
			}),
			layout.Rigid(func() {
				layout.UniformInset(unit.Dp(8)).Layout(c.Gtx, func() {
					l := c.Thm.Body2("Confirmations")
					l.Font.Typeface = c.Thm.Fonts["Mono"]
					l.Color = c.Thm.Colors["Dark"]
					l.Layout(c.Gtx)
				})
			}),
			layout.Rigid(func() {
				layout.UniformInset(unit.Dp(8)).Layout(c.Gtx, func() {
					l := c.Thm.Body2("TxNum")
					l.Font.Typeface = c.Thm.Fonts["Mono"]
					l.Color = c.Thm.Colors["Dark"]
					l.Layout(c.Gtx)
				})
			}),
			layout.Rigid(func() {
				layout.Inset{
					Top:   unit.Dp(8),
					Right: unit.Dp(float32(txwidth - 64)),
				}.Layout(c.Gtx, func() {
					l := c.Thm.Body2("BlockHash")
					l.Font.Typeface = c.Thm.Fonts["Mono"]
					l.Color = c.Thm.Colors["Dark"]
					l.Layout(c.Gtx)
				})
			}),
		)
	})
}

func bodyExplorer(c *component.State) func() {
	return func() {
		c.Rc.GetBlocksExcerpts()
		layout.Flex{Axis: layout.Vertical}.Layout(c.Gtx,
			layout.Rigid(func() {
				blockRowCellLabels(c)
			}),
			layout.Flexed(1, func() {
				explorerPanel := c.Thm.DuoUIpanel()
				explorerPanel.PanelObject = c.Rc.Explorer.Blocks
				explorerPanel.ScrollBar = c.Thm.ScrollBar(16)
				explorerPanelElement.PanelObjectsNumber =
					len(c.Rc.Explorer.Blocks)
				explorerPanel.Layout(c.Gtx, explorerPanelElement,
					func(i int, in interface{}) {
						blocks := in.([]model.DuoUIblock)
						b := blocks[i]
						//blocksList.Layout(c.Gtx, len(c.Rc.Explorer.Blocks), func(i int) {
						//	b := c.Rc.Explorer.Blocks[i]
						blockRow(c, &b)
						//})
					},
				)
			}),
		)
	}
}

func Console(c *component.State) *gelook.DuoUIpage {
	page := gelook.DuoUIpage{
		Title:         "CONSOLE",
		Command:       func() {},
		Border:        4,
		BorderColor:   "ff000000",
		Header:        func() {},
		HeaderPadding: 0,
		Body:          consoleBody(c),
		BodyBgColor:   "ff000000",
		BodyPadding:   4,
		Footer:        func() {},
		FooterPadding: 0,
	}
	return c.Thm.DuoUIpage(page)
}

func consoleBody(c *component.State) func() {
	return func() {
		layout.UniformInset(unit.Dp(8)).Layout(c.Gtx, func() {
			c.Thm.DuoUIcontainer(0, "ff000000").
				Layout(c.Gtx, layout.N, func() {
					layout.Flex{}.Layout(c.Gtx,
						layout.Flexed(1, func() {
							layout.Flex{
								Axis:    layout.Vertical,
								Spacing: layout.SpaceAround,
							}.Layout(c.Gtx, layout.Flexed(1, func() {
								consoleOutputList.Layout(c.Gtx,
									len(c.Rc.ConsoleHistory.Commands), func(i int) {
										t := c.Rc.ConsoleHistory.Commands[i]
										layout.Flex{
											Axis:      layout.Vertical,
											Alignment: layout.End,
										}.Layout(c.Gtx,
											layout.Rigid(
												c.Label(c.Thm.Fonts["Mono"],
													12, c.Thm.Colors["Light"],
													"ds://"+t.ComID),
											),
											layout.Rigid(
												c.Label(c.Thm.Fonts["Mono"],
													12, c.Thm.Colors["Light"],
													t.Out),
											),
										)
									})
							}),
								layout.Rigid(
									c.ConsoleInput(
										consoleInputField, "Run command",
										func(e gel.SubmitEvent) {
											outC := c.Rc.ConsoleCmd(e.Text)
											go func() {
												select {
												case out := <-outC:
													c.Rc.ConsoleHistory.Commands = append(
														c.Rc.ConsoleHistory.Commands,
														model.DuoUIconsoleCommand{
															ComID: e.Text,
															Time:  time.Time{},
															Out:   out,
														})
												}
											}()
										},
									),
								),
							)
						}),
					)
				})
		})
	}
}

func DuoUIaddressBook(c *component.State) *gelook.DuoUIpage {
	page := gelook.DuoUIpage{
		Title:         "ADDRESSBOOK",
		Command:       c.Rc.GetAddressBook(),
		Border:        4,
		BorderColor:   c.Thm.Colors["Light"],
		Header:        addressBookHeader(c, c.Rc.GetAddressBook()),
		HeaderPadding: 4,
		Body:          addressBookBody(c),
		BodyBgColor:   c.Thm.Colors["Light"],
		BodyPadding:   4,
		Footer:        func() {},
		FooterPadding: 0,
	}
	return c.Thm.DuoUIpage(page)
}

func DuoUIexplorer(c *component.State) *gelook.DuoUIpage {
	page := gelook.DuoUIpage{
		Title:         "EXPLORER",
		TxColor:       "",
		Command:       c.Rc.GetBlocksExcerpts(),
		Border:        4,
		BorderColor:   c.Thm.Colors["Light"],
		Header:        explorerHeader(c),
		HeaderBgColor: "",
		HeaderPadding: 4,
		Body:          bodyExplorer(c),
		BodyBgColor:   c.Thm.Colors["Light"],
		BodyPadding:   4,
		Footer:        func() {},
		FooterBgColor: "",
		FooterPadding: 0,
	}
	return c.Thm.DuoUIpage(page)
}

func DuoUIminer(c *component.State) func() {
	return func() {
		c.Rc.GetDuoUIhashesPerSecList()
		layout.Flex{}.Layout(c.Gtx, layout.Flexed(1, func() {
			layout.UniformInset(unit.Dp(0)).Layout(c.Gtx, func() {
				layout.Flex{
					Axis:    layout.Vertical,
					Spacing: layout.SpaceAround,
				}.Layout(c.Gtx,
					layout.Flexed(1, func() {
						//	consoleOutputList.Layout(c.Gtx, c.Rc.Status.Kopach.Hps.Len(), func(i int) {
						//		t := c.Rc.Status.Kopach.Hps.Get(i)
						//		layout.Flex{
						//			Axis:      layout.Vertical,
						//			Alignment: layout.End,
						//		}.Layout(c.Gtx,
						//			layout.Rigid(func() {
						//				sat := c.Thm.Body1(fmt.Sprint(t))
						//				sat.Font.Typeface = c.Thm.Fonts["Mono"]
						//				sat.Color = gelook.HexARGB(c.Thm.Colors["Dark"])
						//				sat.Layout(c.Gtx)
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

func explorerHeader(c *component.State) func() {
	return func() {
		layout.Flex{
			Spacing:   layout.SpaceBetween,
			Axis:      layout.Horizontal,
			Alignment: layout.Middle,
		}.Layout(c.Gtx,
			layout.Rigid(func() {
				c.Thm.DuoUIcounter(c.Rc.GetBlocksExcerpts()).Layout(c.Gtx,
					c.Rc.Explorer.Page, "PAGE",
					fmt.Sprint(c.Rc.Explorer.Page.Value))
			}),
			layout.Rigid(func() {
				c.Thm.DuoUIcounter(c.Rc.GetBlocksExcerpts()).Layout(c.Gtx,
					c.Rc.Explorer.PerPage, "PER PAGE",
					fmt.Sprint(c.Rc.Explorer.PerPage.Value))
			}),
		)
	}
}

func History(c *component.State) *gelook.DuoUIpage {
	page := gelook.DuoUIpage{
		Title:         "HISTORY",
		Command:       c.Rc.GetDuoUItransactions(),
		Border:        4,
		BorderColor:   c.Thm.Colors["Light"],
		Header:        historyHeader(c),
		HeaderPadding: 4,
		Body:          historyBody(c),
		BodyBgColor:   c.Thm.Colors["Light"],
		BodyPadding:   0,
		Footer:        func() {},
		FooterPadding: 0,
	}
	return c.Thm.DuoUIpage(page)
}

func historyBody(c *component.State) func() {
	return func() {
		layout.UniformInset(unit.Dp(16)).Layout(c.Gtx, func() {
			layout.Flex{
				Axis: layout.Vertical,
			}.Layout(c.Gtx,
				layout.Rigid(func() {
					txsPanel := c.Thm.DuoUIpanel()
					txsPanel.PanelObject = c.Rc.History.Txs.Txs
					txsPanel.ScrollBar = c.Thm.ScrollBar(16)
					txsPanelElement.PanelObjectsNumber = len(c.Rc.History.Txs.Txs)
					txsPanel.Layout(c.Gtx, txsPanelElement, func(i int, in interface{}) {
						txs := in.([]model.DuoUItransactionExcerpt)
						t := txs[i]
						c.Thm.DuoUIline(c.Gtx, 0, 0, 1, c.Thm.Colors["Hint"])()
						for t.Link.Clicked(c.Gtx) {
							c.Rc.ShowPage = fmt.Sprintf("TRANSACTION %s", t.TxID)
							c.Rc.GetSingleTx(t.TxID)()
							component.SetPage(c.Rc, txPage(c, t.TxID))
						}
						width := c.Gtx.Constraints.Width.Max
						button := c.Thm.DuoUIbutton(gelook.ButtonParams{})
						button.InsideLayout(c.Gtx, t.Link, func() {
							c.Gtx.Constraints.Width.Min = width
							layout.Flex{
								Spacing: layout.SpaceBetween,
							}.Layout(c.Gtx,
								layout.Rigid(c.TxsDetails(i, &t)),
								layout.Rigid(c.Label(c.Thm.Fonts["Mono"], 12,
									c.Thm.Colors["Secondary"],
									fmt.Sprintf("%0.8f", t.Amount))))
						})
					},
					)
				}),
			)
		})
	}
}

func historyHeader(c *component.State) func() {
	return func() {
		layout.Flex{
			Spacing: layout.SpaceBetween,
		}.Layout(c.Gtx,
			layout.Rigid(c.TransactionsFilter()),
			layout.Rigid(func() {
				c.Thm.DuoUIcounter(c.Rc.GetDuoUItransactions()).
					Layout(c.Gtx, c.Rc.History.PerPage, "TxNum per page: ",
						fmt.Sprint(c.Rc.History.PerPage.Value))
			}),
			layout.Rigid(func() {
				c.Thm.DuoUIcounter(c.Rc.GetDuoUItransactions()).
					Layout(c.Gtx, c.Rc.History.Page, "TxNum page: ",
						fmt.Sprint(c.Rc.History.Page.Value))
			}),
		)
	}
}

func LoadPages(c *component.State) (p map[string]*gelook.DuoUIpage) {
	p = make(map[string]*gelook.DuoUIpage)
	p["OVERVIEW"] = Overview(c)
	p["SEND"] = Send(c)
	// p["RECEIVE"] = c.Thm.DuoUIpage("RECEIVE", 10, func() {},
	// func() {}, func() { c.Thm.H5("receive :").Layout(c.Gtx) }, func() {})
	p["ADDRESSBOOK"] = DuoUIaddressBook(c)
	p["SETTINGS"] = Settings(c)
	p["NETWORK"] = Network(c)
	// p["BLOCK"] = c.Thm.DuoUIpage("BLOCK", 0, func() {},
	// func() {}, func() { c.Thm.H5("block :").Layout(c.Gtx) }, func() {})
	p["HISTORY"] = History(c)
	p["EXPLORER"] = DuoUIexplorer(c)
	p["MINER"] = Miner(c)
	p["CONSOLE"] = Console(c)
	p["LOG"] = Logger(c)
	return
}

func Logger(c *component.State) *gelook.DuoUIpage {
	page := gelook.DuoUIpage{
		Title:         "LOG",
		Command:       func() {},
		Border:        4,
		BorderColor:   c.Thm.Colors["Dark"],
		Header:        func() {},
		HeaderBgColor: "",
		Body:          c.DuoUIlogger(),
		BodyBgColor:   c.Thm.Colors["Dark"],
		Footer:        func() {},
		FooterBgColor: "",
	}
	return c.Thm.DuoUIpage(page)
}

func Miner(c *component.State) *gelook.DuoUIpage {
	page := gelook.DuoUIpage{
		Title:         "MINER",
		TxColor:       "",
		Command:       func() {},
		Border:        4,
		BorderColor:   c.Thm.Colors["Light"],
		Header:        func() {},
		HeaderBgColor: "",
		HeaderPadding: 0,
		Body:          DuoUIminer(c),
		BodyBgColor:   c.Thm.Colors["Dark"],
		BodyPadding:   0,
		Footer:        func() {},
		FooterBgColor: "",
		FooterPadding: 0,
	}
	return c.Thm.DuoUIpage(page)
}

func Network(c *component.State) *gelook.DuoUIpage {
	page := gelook.DuoUIpage{
		Title:         "NETWORK",
		TxColor:       "",
		Command:       c.Rc.GetPeerInfo(),
		Border:        4,
		BorderColor:   c.Thm.Colors["Light"],
		Header:        networkHeader(c),
		HeaderBgColor: "",
		HeaderPadding: 0,
		Body:          networkBody(c),
		BodyBgColor:   c.Thm.Colors["Light"],
		BodyPadding:   0,
		Footer:        func() {},
		FooterBgColor: "",
		FooterPadding: 0,
	}
	return c.Thm.DuoUIpage(page)
}

func networkBody(c *component.State) func() {
	return func() {
		layout.UniformInset(unit.Dp(16)).Layout(c.Gtx, func() {
			layout.Flex{
				Axis: layout.Vertical,
			}.Layout(c.Gtx,
				layout.Rigid(c.PeersList()))
		})
	}
}

func networkHeader(c *component.State) func() {
	return func() {
		layout.Flex{
			Spacing: layout.SpaceBetween,
		}.Layout(c.Gtx,
			// layout.Rigid(component.TransactionsFilter(c)),
			layout.Rigid(func() {
				c.Thm.DuoUIcounter(c.Rc.GetPeerInfo()).Layout(c.Gtx,
					c.Rc.Network.PerPage, "Peers per page: ",
					fmt.Sprint(c.Rc.Network.PerPage.Value))
			}),
			layout.Rigid(func() {
				c.Thm.DuoUIcounter(c.Rc.GetPeerInfo()).Layout(c.Gtx,
					c.Rc.Network.Page, "Peers page: ",
					fmt.Sprint(c.Rc.Network.Page.Value))
			}),
		)
	}
}

func Overview(c *component.State) *gelook.DuoUIpage {
	page := gelook.DuoUIpage{
		Title:         "OVERVIEW",
		Border:        0,
		Command:       func() {},
		Header:        func() {},
		HeaderBgColor: "",
		HeaderPadding: 0,
		Body:          overviewBody(c),
		BodyBgColor:   c.Thm.Colors["Light"],
		BodyPadding:   0,
		Footer:        func() {},
		FooterBgColor: "",
		FooterPadding: 0,
	}
	return c.Thm.DuoUIpage(page)
}

func overviewBody(c *component.State) func() {
	return func() {
		viewport := layout.Flex{Axis: layout.Horizontal}
		if c.Gtx.Constraints.Width.Max < 780 {
			viewport = layout.Flex{Axis: layout.Vertical}
		}
		viewport.Layout(c.Gtx,
			layout.Flexed(0.5, c.DuoUIstatus()),
			layout.Flexed(0.5, c.DuoUIlatestTransactions()),
		)
		op.InvalidateOp{}.Add(c.Gtx.Ops)
	}
}

func Send(c *component.State) *gelook.DuoUIpage {
	page := gelook.DuoUIpage{
		Title:         "SEND",
		TxColor:       "",
		Command:       func() {},
		Border:        4,
		BorderColor:   c.Thm.Colors["Light"],
		Header:        func() {},
		HeaderBgColor: "",
		HeaderPadding: 0,
		Body:          sendBody(c),
		BodyBgColor:   c.Thm.Colors["Light"],
		BodyPadding:   0,
		Footer:        func() {},
		FooterBgColor: "",
		FooterPadding: 0,
	}
	return c.Thm.DuoUIpage(page)
}

func getSendBodyWidgets(c *component.State) []func() {
	return []func(){
		func() {
			c.Thm.DuoUIcontainer(1,
				c.Thm.Colors["Gray"]).Layout(c.Gtx, layout.Center, func() {
				layout.Flex{}.Layout(c.Gtx,
					layout.Flexed(1,
						c.Editor(addressLineEditor, "DUO address",
							func(e gel.EditorEvent) {
								sendStruct.address = addressLineEditor.Text()
							}),
					),
					layout.Rigid(
						c.Button(
							buttonPasteAddress, c.Thm.Fonts["Primary"],
							10, 13, 8, 12, 8,
							c.Thm.Colors["ButtonText"], c.Thm.Colors["ButtonBg"],
							"PASTE ADDRESS", func() {
								addressLineEditor.SetText(clipboard.Get())
							}),
					),
				)
			})
		},
		func() {
			c.Thm.DuoUIcontainer(1, c.Thm.Colors["Gray"]).Layout(c.Gtx,
				layout.Center, func() {
					layout.Flex{}.Layout(c.Gtx,
						layout.Flexed(1, c.Editor(
							amountLineEditor, "DUO Amount",
							func(e gel.EditorEvent) {
								f, err := strconv.ParseFloat(
									amountLineEditor.Text(), 64)
								if err != nil {
								}
								sendStruct.amount = f
							}),
						),
						layout.Rigid(c.Button(
							buttonPasteAmount, c.Thm.Fonts["Primary"],
							10, 13, 8, 12, 8,
							c.Thm.Colors["ButtonText"],
							c.Thm.Colors["ButtonBg"],
							"PASTE AMOUNT",
							func() {
								amountLineEditor.SetText(clipboard.Get())
							}),
						),
					)
				})
		},
		func() {
			layout.Flex{}.Layout(c.Gtx,
				layout.Rigid(c.Button(
					buttonSend, c.Thm.Fonts["Primary"],
					14, 10, 10, 9, 10,
					c.Thm.Colors["ButtonText"], c.Thm.Colors["ButtonBg"],
					"SEND", func() {
						c.Rc.Dialog.Show = true
						c.Rc.Dialog = &model.DuoUIdialog{
							Show: true,
							Green: c.Rc.DuoSend(sendStruct.passPhrase,
								sendStruct.address, 11),
							GreenLabel: "SEND",
							CustomField: func() {
								layout.Flex{}.Layout(c.Gtx,
									layout.Flexed(1,
										c.Editor(
											passLineEditor,
											"Enter your password",
											func(e gel.EditorEvent) {
												sendStruct.passPhrase =
													passLineEditor.Text()
											},
										),
									),
								)
							},
							Red:      func() { c.Rc.Dialog.Show = false },
							RedLabel: "CANCEL",
							Title:    "Are you sure?",
							Text:     "Confirm ParallelCoin send",
						}
					}),
				),
			)
		},
	}
}

func sendBody(c *component.State) func() {
	return func() {
		layout.Flex{}.Layout(c.Gtx,
			layout.Rigid(func() {
				widgets := getSendBodyWidgets(c)
				layautList.Layout(c.Gtx, len(widgets), func(i int) {
					layout.UniformInset(unit.Dp(8)).Layout(c.Gtx, widgets[i])
				})
			}))
		//Info("passPhrase:" + sendStruct.passPhrase)
	}
}

func Settings(c *component.State) *gelook.DuoUIpage {
	page := gelook.DuoUIpage{
		Title:         "SETTINGS",
		TxColor:       "",
		Command:       func() {},
		Border:        4,
		BorderColor:   c.Thm.Colors["Light"],
		Header:        SettingsHeader(c),
		HeaderBgColor: "",
		HeaderPadding: 4,
		Body:          SettingsBody(c),
		BodyBgColor:   c.Thm.Colors["Light"],
		BodyPadding:   0,
		Footer:        func() {},
		FooterBgColor: "",
		FooterPadding: 0,
	}
	return c.Thm.DuoUIpage(page)

	// return c.Thm.DuoUIpage("SETTINGS", 0, func() {},
	//component.ContentHeader(c.Gtx, th, SettingsHeader(c)), SettingsBody(c),
	//func() {
	// var msg string
	// if c.Rc.Settings.Daemon.Config["DisableBanning"].(*bool) != true{
	//	msg = "ima"
	// }else{
	//	msg = "nema"
	// //}
	// ttt := c.Thm.H6(fmt.Sprint(c.Rc.Settings.Daemon.Config))
	// ttt.Color = gelook.HexARGB("ffcfcfcf")
	// ttt.Layout(c.Gtx)
	// })
}

func SettingsBody(c *component.State) func() {
	return func() {
		c.Thm.DuoUIcontainer(16,
			c.Thm.Colors["Light"]).Layout(c.Gtx, layout.N, func() {
			for _, fields := range c.Rc.Settings.Daemon.Schema.Groups {
				if fmt.Sprint(fields.Legend) == c.Rc.Settings.Tabs.Current {
					settingsPanel := c.Thm.DuoUIpanel()
					settingsPanel.PanelObject = fields.Fields
					settingsPanel.ScrollBar = c.Thm.ScrollBar(16)
					settingsPanelElement.PanelObjectsNumber = len(fields.Fields)
					settingsPanel.Layout(c.Gtx, settingsPanelElement,
						func(i int, in interface{}) {
							settings := in.(pod.Fields)
							//t := settings[i]
							//fieldsList.Layout(c.Gtx, len(fields.Fields), func(il int) {
							i = settingsPanelElement.PanelObjectsNumber - 1 - i
							tl := component.Field{
								Field: &settings[i],
							}
							layout.Flex{
								Axis: layout.Vertical,
							}.Layout(c.Gtx,
								layout.Rigid(SettingsItemRow(c, &tl)),
								layout.Rigid(c.Thm.DuoUIline(c.Gtx,
									4, 0, 1, c.Thm.Colors["LightGray"])))
						},
					)
				}
			}
		})
	}
}

func SettingsHeader(c *component.State) func() {
	return func() {
		layout.Flex{Spacing: layout.SpaceBetween}.Layout(c.Gtx,
			layout.Rigid(c.SettingsTabs()),
			layout.Rigid(func() {
				// var settingsRestartButton gelook.DuoUIbutton
				// settingsRestartButton = c.Thm.DuoUIbutton(c.Thm.Fonts["Secondary"],
				// 	"restart",
				// 	c.Thm.Colors["Light"],
				// 	c.Thm.Colors["Dark"],
				// 	c.Thm.Colors["Dark"],
				// 	c.Thm.Colors["Light"],
				// 	"",
				// 	c.Thm.Colors["Light"],
				// 	23, 0, 80, 48, 4, 4)
				// for buttonSettingsRestart.Clicked(c.Gtx) {
				// 	c.Rc.SaveDaemonCfg()
				// }
				// settingsRestartButton.Layout(c.Gtx, buttonSettingsRestart)
			}),
		)
	}
}

func SettingsItemRow(c *component.State, f *component.Field) func() {
	return func() {
		layout.Flex{
			Axis:      layout.Horizontal,
			Alignment: layout.Middle,
		}.Layout(c.Gtx,
			// layout.Rigid(func() {
			//	gelook.DuoUIdrawRectangle(c.Gtx, 30, 3, c.Thm.Colors["Light"],
			//	[4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
			// }),
			layout.Flexed(0.62, func() {
				layout.Flex{
					Axis:    layout.Vertical,
					Spacing: 4,
				}.Layout(c.Gtx,
					layout.Rigid(c.SettingsFieldLabel(f)),
					layout.Rigid(c.SettingsFieldDescription(f)),
				)
			}),
			layout.Flexed(1, c.DuoUIinputField(f)),
		)
	}
}

func singleBlockBody(c *component.State,
	block btcjson.GetBlockVerboseResult) func() {
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
		if c.Gtx.Constraints.Width.Max < 1280 {
			duo = layout.Vertical
		}
		trio := layout.Horizontal
		if c.Gtx.Constraints.Width.Max < 780 {
			trio = layout.Vertical
		}
		blockJSON, _ := json.MarshalIndent(block, "", "  ")
		blockText := string(blockJSON)
		widgets := []func(){
			component.UnoField(c.Gtx, c.ContentLabeledField(
				layout.Vertical, 4, 12, 14,
				"Hash", "Dark", "LightGrayII",
				"Dark", "LightGrayI", fmt.Sprint(block.Hash))),
			component.DuoFields(c.Gtx, duo,
				c.TrioFields(trio, 12, 16,
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
				c.TrioFields(trio, 12, 16,
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
			component.DuoFields(c.Gtx, duo,
				c.ContentLabeledField(layout.Vertical,
					4, 12, 12,
					"MerkleRoot", "Dark", "LightGrayII",
					"Dark", "LightGrayI", block.MerkleRoot),
				c.ContentLabeledField(layout.Vertical,
					4, 12, 12,
					"PowHash", "Dark", "LightGrayII",
					"Dark", "LightGrayI", fmt.Sprint(block.PowHash)),
			),
			component.DuoFields(c.Gtx, duo,
				c.TrioFields(trio, 12, 16,
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
				c.TrioFields(trio, 12, 16,
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
			component.UnoField(c.Gtx, c.ContentLabeledField(
				layout.Vertical, 4, 12, 12,
				"Tx", "Dark", "LightGrayII",
				"Dark", "LightGrayI", fmt.Sprint(block.Tx))),
			component.UnoField(c.Gtx, c.ContentLabeledField(
				layout.Vertical, 4, 12, 12,
				"RawTx", "Dark", "LightGrayII",
				"Dark", "LightGrayI", fmt.Sprint(blockText))),
			c.PageNavButtons(block.PreviousHash,
				block.NextHash, blockPage(c, block.PreviousHash),
				blockPage(c, block.NextHash)),
		}
		layautList.Layout(c.Gtx, len(widgets), func(i int) {
			layout.UniformInset(unit.Dp(0)).Layout(c.Gtx, widgets[i])
		})
	}
}

func singleTxBody(c *component.State,
	tx btcjson.GetTransactionResult) func() {
	return func() {

		//duo := layout.Horizontal
		//if c.Gtx.Constraints.Width.Max < 1280 {
		//	duo = layout.Vertical
		//}
		//trio := layout.Horizontal
		//if c.Gtx.Constraints.Width.Max < 780 {
		//	trio = layout.Vertical
		//}

		//blockJSON, _ := json.MarshalIndent(block, "", "  ")
		//blockText := string(blockJSON)
		widgets := []func(){

			func() {
				c.Thm.H6(tx.TxID).Layout(c.Gtx)
			},
			//component.UnoField(c.Gtx, component.ContentLabeledField(c.Gtx, th,
			//layout.Vertical, 4, 12, 14, "Hash", "Dark", "LightGrayII", "Dark", "LightGrayI", fmt.Sprint(block.Hash))),
			//component.DuoFields(c.Gtx, duo,
			//	component.TrioFields(c.Gtx, th, trio, 12, 16,
			//		"Height", fmt.Sprint(block.Height), "Dark", "LightGrayII",
			//		"Dark", "LightGrayI",
			//		"Confirmations", fmt.Sprint(block.Confirmations), "Dark",
			//		"LightGrayII", "Dark", "LightGrayI",
			//		"Time", fmt.Sprint(time.Unix(block.Time, 0).Format("2006-01-02 15:04:05")),
			//		"LightGrayII", "Dark", "Dark", "LightGrayI",
			//	),
			//	component.TrioFields(c.Gtx, th, trio, 12, 16,
			//		"PowAlgo", fmt.Sprint(block.PowAlgo), algoHeadColor,
			//		algoHeadBgColor, algoColor, algoBgColor,
			//		"Difficulty", fmt.Sprint(block.Difficulty), "Dark",
			//		"LightGrayII", "Dark", "LightGrayI",
			//		"Nonce", fmt.Sprint(block.Nonce), "LightGrayII", "Dark",
			//		"Dark", "LightGrayI",
			//	),
			//),
			//component.DuoFields(c.Gtx, duo,
			//	component.ContentLabeledField(c.Gtx, th, layout.Vertical,
			//	4, 12, 12, "MerkleRoot", "Dark", "LightGrayII", "Dark",
			//	"LightGrayI", block.MerkleRoot),
			//	component.ContentLabeledField(c.Gtx, th, layout.Vertical,
			//	4, 12, 12, "PowHash", "Dark", "LightGrayII", "Dark",
			//	"LightGrayI", fmt.Sprint(block.PowHash)),
			//),
			//component.DuoFields(c.Gtx, duo,
			//	component.TrioFields(c.Gtx, th, trio, 12, 16,
			//		"Size", fmt.Sprint(block.Size), "Dark", "LightGrayII",
			//		"Dark", "LightGrayI",
			//		"Weight", fmt.Sprint(block.Weight), "Dark", "LightGrayII",
			//		"Dark", "LightGrayI",
			//		"Bits", fmt.Sprint(block.Bits), "LightGrayII", "Dark",
			//		"Dark", "LightGrayI",
			//	),
			//	component.TrioFields(c.Gtx, th, trio, 12, 16,
			//		"TxNum", fmt.Sprint(block.TxNum), "Dark", "LightGrayII",
			//		"Dark", "LightGrayI",
			//		"StrippedSize", fmt.Sprint(block.StrippedSize), "Dark",
			//		"LightGrayII", "Dark", "LightGrayI",
			//		"Version", fmt.Sprint(block.Version), "LightGrayII",
			//		"Dark", "Dark", "LightGrayI",
			//	),
			//),
			//component.UnoField(c.Gtx, component.ContentLabeledField(c.Gtx, th,
			//layout.Vertical, 4, 12, 12, "Tx", "Dark", "LightGrayII", "Dark",
			//"LightGrayI", fmt.Sprint(block.Tx))),
			//component.UnoField(c.Gtx, component.ContentLabeledField(c.Gtx, th,
			//layout.Vertical, 4, 12, 12, "RawTx", "Dark", "LightGrayII",
			//"Dark", "LightGrayI", fmt.Sprint(blockText))),
			//component.PageNavButtons(c, block.PreviousHash,
			//block.NextHash, blockPage(c, block.PreviousHash),
			//blockPage(c, block.NextHash)),
		}
		layautList.Layout(c.Gtx, len(widgets), func(i int) {
			layout.UniformInset(unit.Dp(0)).Layout(c.Gtx, widgets[i])
		})
	}
}
func txPage(c *component.State, tx string) *gelook.DuoUIpage {
	page := gelook.DuoUIpage{
		Title:         "BLOCK",
		TxColor:       "",
		Command:       c.Rc.GetSingleTx(tx),
		Border:        4,
		BorderColor:   c.Thm.Colors["Light"],
		Header:        func() {},
		HeaderBgColor: "",
		HeaderPadding: 0,
		Body:          singleTxBody(c, c.Rc.History.SingleTx),
		BodyBgColor:   c.Thm.Colors["Light"],
		BodyPadding:   0,
		Footer:        func() {},
		FooterBgColor: "",
		FooterPadding: 0,
	}
	return c.Thm.DuoUIpage(page)
}
