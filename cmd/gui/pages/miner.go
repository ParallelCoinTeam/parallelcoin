package pages

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gelook"
)

func Miner(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) *gelook.DuoUIpage {
	page := gelook.DuoUIpage{
		Title:         "MINER",
		TxColor:       "",
		Command:       func() {},
		Border:        4,
		Header:        func() {},
		HeaderBgColor: "",
		HeaderPadding: 0,
		Body:          DuoUIminer(rc, gtx, th),
		BodyBgColor:   "",
		BodyPadding:   0,
		Footer:        func() {},
		FooterBgColor: "",
		FooterPadding: 0,
	}
	return th.DuoUIpage(page)
}

func DuoUIminer(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	return func() {
		rc.GetDuoUIhashesPerSecList()
		layout.Flex{}.Layout(gtx,
			layout.Flexed(1, func() {
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
