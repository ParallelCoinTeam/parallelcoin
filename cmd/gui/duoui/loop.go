package duoui

import (
	"errors"
	"image"
	
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/ico"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gui/io/system"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/unit"
	"github.com/p9c/pod/pkg/gui/widget/parallel"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/util/interrupt"
)

func DuoUImainLoop(d *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) error {
	duo := DuoUI{Model: d}
	for {
		select {
		case <-duo.Model.Ready:
			updateTrigger := make(chan struct{}, 1)
			go func() {
			quitTrigger:
				for {
					select {
					case <-updateTrigger:
						log.DEBUG("repaint forced")
						duo.Model.DuoUIwindow.Invalidate()
					case <-duo.Model.Quit:
						break quitTrigger
					}
				}
			}()
			rcd.ListenInit(cx, rc, updateTrigger)
			duo.Model.IsReady = true
		case <-duo.Model.Quit:
			log.DEBUG("quit signal received")
			interrupt.Request()
			// This case is for handling when some external application is controlling the GUI and to gracefully
			// handle the back-end servers being shut down by the interrupt library receiving an interrupt signal
			// Probably nothing needs to be run between starting it and shutting down
			<-interrupt.HandlersDone
			log.DEBUG("closing GUI from interrupt/quit signal")
			return errors.New("shutdown triggered from back end")
		case e := <-duo.Model.DuoUIwindow.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				log.DEBUG("destroy event received")
				interrupt.Request()
				// Here do cleanup like are you sure (optional) modal or shutting down indefinite spinner
				<-interrupt.HandlersDone
				return e.Err
			case system.FrameEvent:
				if rc.Boot.IsBoot {
					duo.Model.DuoUIcontext.Reset(e.Config, e.Size)
					duo.DuoUImainScreen()
					e.Frame(duo.Model.DuoUIcontext.Ops)
				} else {
					duo.Model.DuoUIcontext.Reset(e.Config, e.Size)
					if rc.Boot.IsFirstRun {
						DuoUIloaderCreateWallet(duo.Model, cx)
					} else {
						duo.DuoUIgrid(cx, rc)
						if rc.ShowDialog {
							duo.DuoUIdialog(cx, rc)
						}
						duo.DuoUItoastSys(rc)
					}
					e.Frame(duo.Model.DuoUIcontext.Ops)
					duo.Model.DuoUIcontext.Reset(e.Config, e.Size)
				}
			}
		}
	}
}

// Main wallet screen
func (duo *DuoUI) DuoUImainScreen() {
	helpers.DuoUIdrawRectangle(duo.Model.DuoUIcontext, duo.Model.DuoUIcontext.Constraints.Width.Max, duo.Model.DuoUIcontext.Constraints.Height.Max, duo.Model.DuoUItheme.Color.Bg, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
	// START View <<<
	logo, _ := parallel.NewDuoUIicon(ico.ParallelCoin)
	layout.Flex{Axis: layout.Vertical}.Layout(duo.Model.DuoUIcontext,
		layout.Flexed(0.6, func() {
			layout.Flex{Axis: layout.Horizontal}.Layout(duo.Model.DuoUIcontext,
				
				layout.Rigid(func() {
					layout.UniformInset(unit.Dp(8)).Layout(duo.Model.DuoUIcontext, func() {
						size := duo.Model.DuoUIcontext.Px(unit.Dp(256)) - 2*duo.Model.DuoUIcontext.Px(unit.Dp(8))
						if logo != nil {
							logo.Color = duo.Model.DuoUItheme.Color.Dark
							logo.Layout(duo.Model.DuoUIcontext, unit.Px(float32(size)))
						}
						duo.Model.DuoUIcontext.Dimensions = layout.Dimensions{
							Size: image.Point{X: size, Y: size},
						}
					})
				}),
				layout.Flexed(1, func() {
					layout.UniformInset(unit.Dp(60)).Layout(duo.Model.DuoUIcontext, func() {
						duo.Model.DuoUItheme.H1("PLAN NINE FROM FAR, FAR AWAY SPACE").Layout(duo.Model.DuoUIcontext)
					})
				}),
			)
		}),
		layout.Flexed(0.4, func() {
			DuoUIloader(duo.Model)
		}),
	)
}

// Main wallet screen
func (duo *DuoUI) DuoUIgrid(cx *conte.Xt, rc *rcd.RcVar) {
	// START View <<<
	cs := duo.Model.DuoUIcontext.Constraints
	helpers.DuoUIdrawRectangle(duo.Model.DuoUIcontext, cs.Width.Max, cs.Height.Max, duo.Model.DuoUItheme.Color.Dark, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
	
	layout.Flex{Axis: layout.Vertical}.Layout(duo.Model.DuoUIcontext,
		layout.Rigid(duo.DuoUIheader(rc)),
		layout.Flexed(1, duo.DuoUIbody(cx, rc)),
		layout.Rigid(duo.DuoUIfooter(rc)),
	)
}
