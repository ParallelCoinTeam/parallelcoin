package duoui

import (
	"fmt"
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/unit"
	"github.com/p9c/pod/cmd/gui/helpers"
	"image/color"
)

func DuoUIoverview(duo *DuoUI) layout.FlexChild {
	return duo.comp.Content.Layout.Flex(duo.gc, 1, func() {

		duo.GetDuOSbalance()
		//duo.GetDuOSblockHeight()
		//duo.GetDuOStatus()
		//duo.GetDuOSlocalLost()
		//duo.GetDuOSdifficulty()

		duo.comp.Overview.Inset.Layout(duo.gc, func() {
			helpers.DuoUIdrawRect(duo.gc, duo.cs.Width.Max, duo.cs.Height.Max, color.RGBA{A: 0xff, R: 0x30, G: 0xcf, B: 0xcf}, 0, 0, 0, 0)
			// Overview <<<
			overviewTop := duo.comp.Overview.Layout.Rigid(duo.gc, func() {
				//duo.comp.content.i.Layout(duo.gc, func() {
				helpers.DuoUIdrawRect(duo.gc, duo.cs.Width.Max, 180, color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}, 0, 0, 0, 0)
				// OverviewTop <<<
				balance := duo.comp.OverviewTop.Layout.Flex(duo.gc, 0.4, func() {
					helpers.DuoUIdrawRect(duo.gc, duo.cs.Width.Max-30, 180, color.RGBA{A: 0xff, R: 0x30, G: 0x30, B: 0x30}, 9.9, 9.9, 9.9, 9.9)
					in := layout.UniformInset(unit.Dp(60))

					in.Layout(duo.gc, func() {
						bal := duo.th.H3("Balance :" + duo.rc.Balance + " DUO")

						bal.Color = color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}
						bal.Layout(duo.gc)
					})

				})

				duo.comp.OverviewTop.Layout.Layout(duo.gc, balance, DuoUIsendreceive(duo))
				// OverviewTop >>>
				//})
			})
			overviewBottom := duo.comp.Overview.Layout.Flex(duo.gc, 1, func() {
				//duo.comp.content.i.Layout(duo.gc, func() {
				//DuoUIdrawRect(duo.gc, duo.cs.Width.Max, duo.cs.Height.Max, color.RGBA{A: 0xff, R: 0xcf, G: 0x30, B: 0x30})
				// OverviewBottom <<<
				transactions := duo.comp.OverviewBottom.Layout.Flex(duo.gc, 0.7, func() {
					helpers.DuoUIdrawRect(duo.gc, duo.cs.Width.Max, duo.cs.Height.Max, color.RGBA{A: 0xff, R: 0xcf, G: 0x30, B: 0xcf}, 0, 0, 0, 0)

					//duo.gc.Reset(e.Config, e.Size)

				})
				status := duo.comp.OverviewBottom.Layout.Flex(duo.gc, 0.3, func() {
					helpers.DuoUIdrawRect(duo.gc, duo.cs.Width.Max, duo.cs.Height.Max, color.RGBA{A: 0xff, R: 0x30, G: 0x30, B: 0xcf}, 0, 0, 0, 0)

					duo.comp.Status.Inset.Layout(duo.gc, func() {
						helpers.DuoUIdrawRect(duo.gc, duo.cs.Width.Max, duo.cs.Height.Max, color.RGBA{A: 0xff, R: 0xcf, G: 0x30, B: 0xcf}, 0, 0, 0, 0)

						balance := duo.comp.Status.Layout.Rigid(duo.gc, func() {
							duo.th.H5("balance :" + duo.rc.Balance).Layout(duo.gc)
						})
						blockheight := duo.comp.Status.Layout.Rigid(duo.gc, func() {
							duo.th.H5("blockheight :" + fmt.Sprint(duo.rc.BlockHeight)).Layout(duo.gc)
						})
						difficulty := duo.comp.Status.Layout.Rigid(duo.gc, func() {
							duo.th.H5("difficulty :" + fmt.Sprintf("%f", duo.rc.Difficulty)).Layout(duo.gc)
						})
						connections := duo.comp.Status.Layout.Rigid(duo.gc, func() {
							duo.th.H5("connections :" + fmt.Sprint(duo.rc.Connections)).Layout(duo.gc)
						})

						duo.comp.Status.Layout.Layout(duo.gc, balance, blockheight, difficulty, connections)
					})

				})
				duo.comp.OverviewBottom.Layout.Layout(duo.gc, transactions, status)
				// OverviewBottom >>>
				//})
			})
			duo.comp.Overview.Layout.Layout(duo.gc, overviewTop, overviewBottom)
			// Overview >>>
		})
	})
}
