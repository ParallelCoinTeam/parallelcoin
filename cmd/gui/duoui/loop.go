package duoui

import (
	"errors"
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/ico"
	"github.com/p9c/pod/cmd/gui/loader"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gui/io/system"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/unit"
	"github.com/p9c/pod/pkg/gui/widget/parallel"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/util/interrupt"
	"image"
)

func DuoUImainLoop(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) error {
	for {
		select {
		case <-duo.Ready:
			duo.IsReady = true
		case <-duo.Quit:
			log.DEBUG("quit signal received")
			interrupt.Request()
			// This case is for handling when some external application is controlling the GUI and to gracefully
			// handle the back-end servers being shut down by the interrupt library receiving an interrupt signal
			// Probably nothing needs to be run between starting it and shutting down
			<-interrupt.HandlersDone
			log.DEBUG("closing GUI from interrupt/quit signal")
			return errors.New("shutdown triggered from back end")
		case e := <-duo.DuoUIwindow.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				log.DEBUG("destroy event received")
				interrupt.Request()
				// Here do cleanup like are you sure (optional) modal or shutting down indefinite spinner
				<-interrupt.HandlersDone
				return e.Err
			case system.FrameEvent:
				if rc.Boot.IsBoot {
					duo.DuoUIcontext.Reset(e.Config, e.Size)
					DuoUImainScreen(duo, rc)
					e.Frame(duo.DuoUIcontext.Ops)
				} else {
					duo.DuoUIcontext.Reset(e.Config, e.Size)
					if rc.Boot.IsFirstRun {
						loader.DuoUIloaderCreateWallet(duo, cx)
					} else {
						DuoUIgrid(duo, cx, rc)
						if rc.IsNotificationRun {
							DuoUIdialog(duo, cx, rc)
						}
					}
					e.Frame(duo.DuoUIcontext.Ops)
					duo.DuoUIcontext.Reset(e.Config, e.Size)
				}
			}
		}
	}
}

// Main wallet screen
func DuoUImainScreen(duo *models.DuoUI, rc *rcd.RcVar) {
	helpers.DuoUIdrawRectangle(duo.DuoUIcontext, duo.DuoUIcontext.Constraints.Width.Max, duo.DuoUIcontext.Constraints.Height.Max, duo.DuoUItheme.Color.Bg, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
	// START View <<<
	logo, _ := parallel.NewDuoUIicon(ico.ParallelCoin)
	layout.Flex{Axis: layout.Vertical}.Layout(duo.DuoUIcontext,
		layout.Flexed(0.6, func() {
		layout.Flex{Axis:layout.Horizontal}.Layout(duo.DuoUIcontext,

			layout.Rigid(func(){
				layout.UniformInset(unit.Dp(8)).Layout(duo.DuoUIcontext, func() {
					size := duo.DuoUIcontext.Px(unit.Dp(256)) - 2*duo.DuoUIcontext.Px(unit.Dp(8))
					if logo != nil {
						logo.Color = duo.DuoUItheme.Color.Dark
						logo.Layout(duo.DuoUIcontext, unit.Px(float32(size)))
					}
					duo.DuoUIcontext.Dimensions = layout.Dimensions{
						Size: image.Point{X: size, Y: size},
					}
				})
			}),
			layout.Flexed(1, func(){
				layout.UniformInset(unit.Dp(60)).Layout(duo.DuoUIcontext, func() {
					duo.DuoUItheme.H1("PLAN NINE FROM FAR, FAR AWAY SPACE").Layout(duo.DuoUIcontext)
				})
			}),
			)
		}),
		layout.Flexed(0.4, func() {
			loader.DuoUIloader(duo)
		}),
	)
}

// Main wallet screen
func DuoUIgrid(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) {
	// START View <<<
	cs := duo.DuoUIcontext.Constraints
	helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, cs.Height.Max, duo.DuoUItheme.Color.Dark, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})

	layout.Flex{Axis: layout.Vertical}.Layout(duo.DuoUIcontext,
		layout.Rigid(DuoUIheader(duo, rc)),
		layout.Flexed(1, DuoUIbody(duo, cx, rc)),
		layout.Rigid(DuoUIfooter(duo, rc)),
	)
}
