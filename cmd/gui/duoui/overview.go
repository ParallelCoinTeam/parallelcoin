package duoui

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/unit"
	"image/color"
)

func DuoUIoverview(duo *DuoUI) layout.FlexChild {
	return duo.comp.content.l.Flex(duo.gc, 1, func() {

		duo.GetDuOSbalance()
		//duo.GetDuOSblockHeight()
		//duo.GetDuOStatus()
		//duo.GetDuOSlocalLost()
		//duo.GetDuOSdifficulty()

		duo.comp.overview.i.Layout(duo.gc, func() {
			DuoUIdrawRect(duo.gc, duo.cs.Width.Max, duo.cs.Height.Max, color.RGBA{A: 0xff, R: 0x30, G: 0xcf, B: 0xcf}, 0, 0, 0, 0)
			// Overview <<<
			overviewTop := duo.comp.overview.l.Rigid(duo.gc, func() {
				//duo.comp.content.i.Layout(duo.gc, func() {
				DuoUIdrawRect(duo.gc, duo.cs.Width.Max, 180, color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}, 0, 0, 0, 0)
				// OverviewTop <<<
				balance := duo.comp.overviewTop.l.Flex(duo.gc, 0.4, func() {
					DuoUIdrawRect(duo.gc, duo.cs.Width.Max-30, 180, color.RGBA{A: 0xff, R: 0x30, G: 0x30, B: 0x30}, 9.9, 9.9, 9.9, 9.9)
					in := layout.UniformInset(unit.Dp(60))

					in.Layout(duo.gc, func() {
						bal := duo.th.H3("Balance :" + duo.rc.Balance + " DUO")

						bal.Color = color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}
						bal.Layout(duo.gc)
					})

				})

				duo.comp.overviewTop.l.Layout(duo.gc, balance, DuoUIsendreceive(duo))
				// OverviewTop >>>
				//})
			})
			overviewBottom := duo.comp.overview.l.Flex(duo.gc, 1, func() {
				//duo.comp.content.i.Layout(duo.gc, func() {
				//DuoUIdrawRect(duo.gc, duo.cs.Width.Max, duo.cs.Height.Max, color.RGBA{A: 0xff, R: 0xcf, G: 0x30, B: 0x30})
				// OverviewBottom <<<
				transactions := duo.comp.overviewBottom.l.Flex(duo.gc, 0.7, func() {
					DuoUIdrawRect(duo.gc, duo.cs.Width.Max, duo.cs.Height.Max, color.RGBA{A: 0xff, R: 0xcf, G: 0x30, B: 0xcf}, 0, 0, 0, 0)

					//duo.gc.Reset(e.Config, e.Size)

				})
				status := duo.comp.overviewBottom.l.Flex(duo.gc, 0.3, func() {
					DuoUIdrawRect(duo.gc, duo.cs.Width.Max, duo.cs.Height.Max, color.RGBA{A: 0xff, R: 0x30, G: 0x30, B: 0xcf}, 0, 0, 0, 0)

					duo.comp.status.i.Layout(duo.gc, func() {
						DuoUIdrawRect(duo.gc, duo.cs.Width.Max, duo.cs.Height.Max, color.RGBA{A: 0xff, R: 0xcf, G: 0x30, B: 0xcf}, 0, 0, 0, 0)

						balance := duo.comp.status.l.Rigid(duo.gc, func() {
							duo.th.H5("balance :" + duo.rc.Balance).Layout(duo.gc)
						})
						blockheight := duo.comp.status.l.Rigid(duo.gc, func() {
							duo.th.H5("blockheight :" + fmt.Sprint(duo.rc.BlockHeight)).Layout(duo.gc)
						})
						difficulty := duo.comp.status.l.Rigid(duo.gc, func() {
							duo.th.H5("difficulty :" + fmt.Sprintf("%f", duo.rc.Difficulty)).Layout(duo.gc)
						})
						connections := duo.comp.status.l.Rigid(duo.gc, func() {
							duo.th.H5("connections :" + fmt.Sprint(duo.rc.Connections)).Layout(duo.gc)
						})

						duo.comp.status.l.Layout(duo.gc, balance, blockheight, difficulty, connections)
					})

				})
				duo.comp.overviewBottom.l.Layout(duo.gc, transactions, status)
				// OverviewBottom >>>
				//})
			})
			duo.comp.overview.l.Layout(duo.gc, overviewTop, overviewBottom)
			// Overview >>>
		})
	})
}
