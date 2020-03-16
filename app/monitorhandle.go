// +build !headless

package app

import (
	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/urfave/cli"

	"github.com/p9c/pod/cmd/gui/pages"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gelook"
	log "github.com/p9c/pod/pkg/logi"
)

var monitorHandle = func(cx *conte.Xt) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		var (
			theme    = gelook.NewDuoUItheme()
			mainList = &layout.List{
				Axis: layout.Vertical,
			}
		)
		Configure(cx, c)
		rc := rcd.RcInit(cx)
		log.L.Warn("starting monitor GUI")
		w := app.NewWindow(
			app.Size(unit.Dp(1024), unit.Dp(600)),
			app.Title("ParallelCoin"),
		)
		gtx := layout.NewContext(w.Queue())
		for e := range w.Events() {
			switch e := e.(type) {
			case system.DestroyEvent:
				log.L.Debug("destroy event received")
				close(cx.KillAll)
				return e.Err
			case system.FrameEvent:
				gtx.Reset(e.Config, e.Size)
				layout.Flex{
					Axis: layout.Vertical,
				}.Layout(gtx,
					layout.Rigid(func() {
						cs := gtx.Constraints
						gelook.DuoUIdrawRectangle(gtx, cs.Width.Max,
							cs.Height.Max, theme.Colors["Dark"],
							[4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
						pages.SettingsHeader(rc, gtx, theme)()
					}),
					layout.Flexed(1, func() {
						cs := gtx.Constraints
						gelook.DuoUIdrawRectangle(gtx, cs.Width.Max,
							cs.Height.Max, theme.Colors["Light"],
							[4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
						controllers := []func(){
							pages.SettingsBody(rc, gtx, theme),
						}
						mainList.Layout(gtx, len(controllers), func(i int) {
							layout.UniformInset(unit.Dp(10)).Layout(gtx,
								controllers[i])
						})
					}),
				)
				e.Frame(gtx.Ops)
			}
			// w.Invalidate()
		}
		go app.Main()
		select {
		case <-cx.KillAll:
		}
		return
	}
}
