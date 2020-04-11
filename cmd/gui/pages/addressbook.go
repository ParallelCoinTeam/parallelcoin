package pages

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/p9c/pod/cmd/gui/component"
	"github.com/p9c/pod/cmd/gui/model"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gui/clipboard"
	"github.com/p9c/pod/pkg/gui/gel"
	"github.com/p9c/pod/pkg/gui/gelook"
)

var (
	addressBookPanelElement = gel.NewPanel()
	showMiningAddresses     = &gel.CheckBox{}
	buttonNewAddress        = new(gel.Button)
	address                 string
)

func DuoUIaddressBook(rc *rcd.RcVar, gtx *layout.Context,
	th *gelook.DuoUItheme) *gelook.DuoUIpage {
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

func addressBookBody(rc *rcd.RcVar, gtx *layout.Context,
	th *gelook.DuoUItheme) func() {
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

func addressBookHeader(rc *rcd.RcVar, gtx *layout.Context,
	th *gelook.DuoUItheme, pageFunc func()) func() {
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
				})))
	}
}

func addressBookContent(rc *rcd.RcVar, gtx *layout.Context,
	th *gelook.DuoUItheme) func() {
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
								12, th.Colors["Dark"], fmt.Sprint(t.Index))),
						layout.Flexed(0.2,
							component.Label(gtx, th, th.Fonts["Primary"],
								12, th.Colors["Dark"], t.Account)),
						layout.Rigid(component.MonoButton(gtx, th,
							t.Copy, 12, "", "",
							"Mono", t.Address, func() {
								clipboard.Set(t.Address)
							})),
						layout.Flexed(0.4,
							component.Label(gtx, th, th.Fonts["Primary"],
								14, th.Colors["Dark"], t.Label)),
						layout.Flexed(0.2,
							component.Label(gtx, th, th.Fonts["Primary"],
								12, th.Colors["Dark"], fmt.Sprint(t.Amount))),
						layout.Rigid(component.MonoButton(gtx, th,
							t.QrCode, 12, "", "",
							"Secondary", "QR",
							component.QrDialog(rc, gtx, t.Address))),
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
