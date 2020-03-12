package pages

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/p9c/pod/pkg/gelook"
	"github.com/p9c/pod/cmd/gui/rcd"
)

func Miner(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) *gelook.DuoUIpage {
	return th.DuoUIpage("MINER", 0, func() {}, func() {}, DuoUIminer(rc, gtx, th), func() {})
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
							consoleOutputList.Layout(gtx, rc.Status.Kopach.Hps.Len(), func(i int) {
								t := rc.Status.Kopach.Hps.Get(i)
								layout.Flex{
									Axis:      layout.Vertical,
									Alignment: layout.End,
								}.Layout(gtx,
									layout.Rigid(func() {
										sat := th.Body1(fmt.Sprint(t))
										sat.Font.Typeface = th.Fonts["Mono"]
										sat.Color = gelook.HexARGB(th.Colors["Dark"])
										sat.Layout(gtx)
									}),
								)
							})
						}),
					)
				})
			}),
		)
	}
}
