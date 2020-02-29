package pages

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/mvc/component"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gui/clipboard"
)

var (
	addressBookList = &layout.List{
		Axis: layout.Vertical,
		//ScrollToEnd: true,
	}
)

func AddressBook(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme) *theme.DuoUIpage {
	return th.DuoUIpage("ADDRESSBOOK", 10, func() {}, func() {}, addressBookBody(rc, gtx, th), func() {})
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
							addressBookList.Layout(gtx, len(rc.AddressBook.Addresses), func(i int) {
								t := rc.AddressBook.Addresses[i]
								layout.Flex{Axis: layout.Vertical}.Layout(gtx,
									layout.Rigid(func() {
										layout.Flex{
											Alignment: layout.End,
										}.Layout(gtx,
											layout.Flexed(0.1, component.Label(gtx, th, th.Font.Primary, fmt.Sprint(t.Index))),
											layout.Flexed(0.2, component.Label(gtx, th, th.Font.Primary, t.Account)),
											layout.Rigid(component.Button(gtx, th, t.Copy, th.Font.Mono, t.Address, func() { clipboard.Set(t.Address) })),
											layout.Flexed(0.4, component.Label(gtx, th, th.Font.Primary, t.Label)),
											layout.Flexed(0.3, component.Label(gtx, th, th.Font.Primary, fmt.Sprint(t.Amount))),
										)
									}),
									layout.Rigid(component.HorizontalLine(gtx, 1, th.Color.Hint)),
								)
							})
						}))
				})
			}),
		)
	}
}
