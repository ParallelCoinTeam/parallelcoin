package pages

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/component"
	"github.com/p9c/pod/cmd/gui/model"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gui/clipboard"
	"github.com/p9c/pod/pkg/gui/controller"
	"github.com/p9c/pod/pkg/gui/theme"
)

var (
	addressBookList = &layout.List{
		Axis: layout.Vertical,
		//ScrollToEnd: true,
	}
	addressBookPanel = &controller.Panel{
		Name: "",
		PanelContentLayout: &layout.List{
			Axis:        layout.Vertical,
			ScrollToEnd: false,
		},
	}
	showMiningAddresses = &controller.CheckBox{}
	buttonNewAddress    = new(controller.Button)
)

func DuoUIaddressBook(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme) *theme.DuoUIpage {
	return th.DuoUIpage("ADDRESSBOOK", 0, rc.GetAddressBook(), component.ContentHeader(gtx, th, addressBookHeader(rc, gtx, th)), addressBookBody(rc, gtx, th), func() {})
}

func addressBookBody(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme) func() {
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
			}),
		)
	}
}

func addressBookHeader(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme) func() {
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
			layout.Rigid(component.Button(gtx, th, buttonNewAddress, th.Fonts["Secondary"], 12, th.Colors["ButtonText"], th.Colors["Dark"], "NEW ADDRESS", func() {
				rc.Dialog.Show = true

				rc.Dialog = &model.DuoUIdialog{
					Show: true,
					Ok:   nil,
					Close: func() {

					},
					CustomField: func() {
						layout.Flex{}.Layout(gtx,
							layout.Flexed(1, component.Editor(gtx, th, passLineEditor, "Enter your password", func(e controller.SubmitEvent) {
								passPharse = e.Text
							})))
					},
					Cancel: func() { rc.Dialog.Show = false },
					Title:  "Copy address",
					Text:   rc.CreateNewAddress(""),
				}
			})))
	}
}

func addressBookContent(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme) func() {
	return func() {
		addressBookList.Layout(gtx, len(rc.AddressBook.Addresses), func(i int) {
			t := rc.AddressBook.Addresses[i]
			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func() {
					layout.Flex{
						Alignment: layout.End,
					}.Layout(gtx,
						layout.Flexed(0.1, component.Label(gtx, th, th.Fonts["Primary"], 12, th.Colors["Dark"], fmt.Sprint(t.Index))),
						layout.Flexed(0.2, component.Label(gtx, th, th.Fonts["Primary"], 12, th.Colors["Dark"], t.Account)),
						layout.Rigid(component.Button(gtx, th, t.Copy, th.Fonts["Mono"], 12, th.Colors["ButtonText"], th.Colors["ButtonBg"], t.Address, func() { clipboard.Set(t.Address) })),
						layout.Flexed(0.4, component.Label(gtx, th, th.Fonts["Primary"], 14, th.Colors["Dark"], t.Label)),
						layout.Flexed(0.3, component.Label(gtx, th, th.Fonts["Primary"], 12, th.Colors["Dark"], fmt.Sprint(t.Amount))),
						layout.Rigid(component.Button(gtx, th, t.QrCode, th.Fonts["Mono"], 12, th.Colors["ButtonText"], th.Colors["ButtonBg"], "QR", component.QrDialog(rc,gtx))),
					)
				}),
				layout.Rigid(component.HorizontalLine(gtx, 1, th.Colors["Hint"])),
			)
		})
	}
}
