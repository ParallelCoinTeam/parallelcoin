package pages

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/p9c/pod/cmd/gui/component"
	"github.com/p9c/pod/cmd/gui/model"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gel"
	"github.com/p9c/pod/pkg/gelook"
	"github.com/p9c/pod/pkg/gui/clipboard"
)

var (
	addressBookPanelElement = &gel.Panel{
		PanelContentLayout: &layout.List{
			Axis:        layout.Vertical,
			ScrollToEnd: false,
		},
		ScrollBar: &gel.ScrollBar{
			//OperateValue: nil,
			Body: new(gel.ScrollBarBody),
			Up:   new(gel.Button),
			Down: new(gel.Button),
		},
	}
	showMiningAddresses = &gel.CheckBox{}
	buttonNewAddress    = new(gel.Button)
	address             string
)

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
				th.DuoUIcheckBox("SHOW MINING ADDRESSES", th.Colors["Light"], th.Colors["Light"]).Layout(gtx, showMiningAddresses)
			}),
			layout.Rigid(func() {
				// th.DuoUIcounter(rc.GetBlocksExcerpts()).Layout(gtx, rc.Explorer.Page, "PAGE", fmt.Sprint(rc.Explorer.Page.Value))
			}),
			// layout.Rigid(component.Button(gtx, th, buttonNewAddress, th.Fonts["Secondary"], 12, th.Colors["ButtonText"], th.Colors["Dark"], "NEW ADDRESS", component.QrDialog(rc, gtx, rc.CreateNewAddress("")))))
			layout.Rigid(component.Button(gtx, th, buttonNewAddress, th.Fonts["Secondary"], 12, th.Colors["Dark"], th.Colors["Light"], "NEW ADDRESS", func() {
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

func addressBookContent(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	return func() {
		addressBookPanel := gelook.Panel{}
		addressBookPanel.PanelObject = rc.AddressBook.Addresses
		addressBookPanel.ScrollBar = th.ScrollBar()
		addressBookPanelElement.PanelObjectsNumber = len(rc.AddressBook.Addresses)
		addressBookPanel.Layout(gtx, addressBookPanelElement, func(i int, in interface{}) {
			addresses := in.([]model.DuoUIaddress)
			t := addresses[i]
			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func() {
					layout.Flex{
						Alignment: layout.End,
					}.Layout(gtx,
						layout.Flexed(0.2, component.Label(gtx, th, th.Fonts["Primary"], 12, th.Colors["Dark"], fmt.Sprint(t.Index))),
						layout.Flexed(0.2, component.Label(gtx, th, th.Fonts["Primary"], 12, th.Colors["Dark"], t.Account)),
						layout.Rigid(component.Button(gtx, th, t.Copy, th.Fonts["Mono"], 12, th.Colors["ButtonText"], th.Colors["ButtonBg"], t.Address, func() { clipboard.Set(t.Address) })),
						layout.Flexed(0.4, component.Label(gtx, th, th.Fonts["Primary"], 14, th.Colors["Dark"], t.Label)),
						layout.Flexed(0.2, component.Label(gtx, th, th.Fonts["Primary"], 12, th.Colors["Dark"], fmt.Sprint(t.Amount))),
						layout.Rigid(component.Button(gtx, th, t.QrCode, th.Fonts["Mono"], 12, th.Colors["ButtonText"], th.Colors["Info"], "QR", component.QrDialog(rc, gtx, t.Address))),
					)
				}),
				layout.Rigid(th.DuoUIline(gtx, 0, 0, 1, th.Colors["Hint"])),
			)
		})
		// }).Layout(gtx, addressBookPanel)
	}
}
