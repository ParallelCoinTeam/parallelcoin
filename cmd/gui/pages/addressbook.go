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
	addressBookList = &layout.List{
		Axis: layout.Vertical,
		//ScrollToEnd: true,
	}
	addressBookPanel = &gel.Panel{
		Name: "",
		PanelContentLayout: &layout.List{
			Axis:        layout.Vertical,
			ScrollToEnd: false,
		},
	}
	showMiningAddresses = &gel.CheckBox{}
	buttonNewAddress    = new(gel.Button)
	address             string
)

func DuoUIaddressBook(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) *gelook.DuoUIpage {
	return th.DuoUIpage("ADDRESSBOOK", 0, rc.GetAddressBook(), component.ContentHeader(gtx, th, addressBookHeader(rc, gtx, th, rc.GetAddressBook())), addressBookBody(rc, gtx, th), func() {})
}

func addressBookBody(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	return func() {
		layout.UniformInset(unit.Dp(8)).Layout(gtx, func() {
			th.DuoUIitem(0, th.Colors["Dark"]).Layout(gtx, layout.N, func() {
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
			})
		})
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
					//rc.GetAddressBook()()
				} else {
					//rc.GetAddressBook()()
				}
				th.DuoUIcheckBox("SHOW MINING ADDRESSES", th.Colors["Dark"], th.Colors["Dark"]).Layout(gtx, showMiningAddresses)
			}),
			layout.Rigid(func() {
				//th.DuoUIcounter(rc.GetBlocksExcerpts()).Layout(gtx, rc.Explorer.Page, "PAGE", fmt.Sprint(rc.Explorer.Page.Value))
			}),
			//layout.Rigid(component.Button(gtx, th, buttonNewAddress, th.Fonts["Secondary"], 12, th.Colors["ButtonText"], th.Colors["Dark"], "NEW ADDRESS", component.QrDialog(rc, gtx, rc.CreateNewAddress("")))))
			layout.Rigid(component.Button(gtx, th, buttonNewAddress, th.Fonts["Secondary"], 12, th.Colors["ButtonText"], th.Colors["Dark"], "NEW ADDRESS", func() {
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
		//th.DuoUIpanel(func() {
		addressBookList.Layout(gtx, len(rc.AddressBook.Addresses), func(i int) {
			t := rc.AddressBook.Addresses[i]
			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func() {
					layout.Flex{
						Alignment: layout.End,
					}.Layout(gtx,
						layout.Flexed(0.2, component.Label(gtx, th, th.Fonts["Primary"], 12, th.Colors["Light"], fmt.Sprint(t.Index))),
						layout.Flexed(0.2, component.Label(gtx, th, th.Fonts["Primary"], 12, th.Colors["Light"], t.Account)),
						layout.Rigid(component.Button(gtx, th, t.Copy, th.Fonts["Mono"], 12, th.Colors["ButtonText"], th.Colors["ButtonBg"], t.Address, func() { clipboard.Set(t.Address) })),
						layout.Flexed(0.4, component.Label(gtx, th, th.Fonts["Primary"], 14, th.Colors["Light"], t.Label)),
						layout.Flexed(0.2, component.Label(gtx, th, th.Fonts["Primary"], 12, th.Colors["Light"], fmt.Sprint(t.Amount))),
						layout.Rigid(component.Button(gtx, th, t.QrCode, th.Fonts["Mono"], 12, th.Colors["ButtonText"], th.Colors["Info"], "QR", component.QrDialog(rc, gtx, t.Address))),
					)
				}),
				layout.Rigid(component.HorizontalLine(gtx, 1, th.Colors["Hint"])),
			)
		})
		//}).Layout(gtx, addressBookPanel)
	}
}
