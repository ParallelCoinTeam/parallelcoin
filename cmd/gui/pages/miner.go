package pages

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
	"github.com/p9c/pod/cmd/gui/rcd"
)

func Miner(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme) *theme.DuoUIpage {
	return th.DuoUIpage("MINER", 0, func() {}, func() {}, DuoUIminer(rc, gtx, th), func() {})
}

func DuoUIminer(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme) func() {
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
							consoleOutputList.Layout(gtx, len(rc.Status.Kopach.Hps), func(i int) {
								t := rc.Status.Kopach.Hps[i]
								layout.Flex{
									Axis:      layout.Vertical,
									Alignment: layout.End,
								}.Layout(gtx,
									layout.Rigid(func() {
										sat := th.Body1(fmt.Sprint(t))
										sat.Font.Typeface = th.Font.Mono
										sat.Color = theme.HexARGB(th.Color.Dark)
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
