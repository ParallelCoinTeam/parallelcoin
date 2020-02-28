package duoui

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/unit"
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

func addressBook(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme) func() {
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
											layout.Flexed(0.1, func() {
												sat := th.Body1(fmt.Sprint(t.Index))
												sat.Font.Typeface = th.Font.Primary
												sat.Color = theme.HexARGB(th.Color.Dark)
												sat.Layout(gtx)
											}),
											layout.Rigid(func() {

												var copyButton theme.DuoUIbutton
												copyButton = th.DuoUIbutton(th.Font.Mono, t.Address, th.Color.Light, th.Color.Primary, "", "", 16, 0, 300, 24, 0, 0)

												for t.Copy.Clicked(gtx) {

													clipboard.Set(t.Address)
												}
												copyButton.Layout(gtx, t.Copy)

											}),
											layout.Flexed(0.2, func() {
												sat := th.Body1(t.Account)
												sat.Font.Typeface = th.Font.Primary
												sat.Color = theme.HexARGB(th.Color.Dark)
												sat.Layout(gtx)
											}),
											layout.Flexed(0.4, func() {
												sat := th.Body1(t.Label)
												sat.Font.Typeface = th.Font.Primary
												sat.Color = theme.HexARGB(th.Color.Dark)
												sat.Layout(gtx)
											}),
											layout.Flexed(0.3, func() {
												sat := th.Body1(fmt.Sprint(t.Amount))
												sat.Font.Typeface = th.Font.Primary
												sat.Color = theme.HexARGB(th.Color.Dark)
												sat.Layout(gtx)
											}),
										)
									}),
									layout.Rigid(line(gtx, th.Color.Hint)),
								)
							})
						}))
				})
			}),
		)
	}
}
