package pages

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/component"
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
)

func AddressBook(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme) *theme.DuoUIpage {
	return th.DuoUIpage("ADDRESSBOOK", 0, func() {}, func() {}, addressBookBody(rc, gtx, th), func() {})
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
						layout.Flexed(1, func() {
							th.DuoUIpanel(addressBookContent(rc, gtx, th)).Layout(gtx, addressBookPanel)
						}))
				})
			}),
		)
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
						layout.Flexed(0.4, component.Label(gtx, th, th.Fonts["Primary"], 12, th.Colors["Dark"], t.Label)),
						layout.Flexed(0.3, component.Label(gtx, th, th.Fonts["Primary"], 12, th.Colors["Dark"], fmt.Sprint(t.Amount))),
					)
				}),
				layout.Rigid(component.HorizontalLine(gtx, 1, th.Colors["Hint"])),
			)
		})
	}
}
