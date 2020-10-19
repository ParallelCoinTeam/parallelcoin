package gui

import (
	"gioui.org/app"
	"github.com/urfave/cli"

	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/pkg/gui/f"
	"github.com/p9c/pod/pkg/gui/fonts/p9fonts"
	"github.com/p9c/pod/pkg/gui/p9"
	"github.com/p9c/pod/pkg/util/interrupt"
)

func Main(cx *conte.Xt, c *cli.Context) (err error) {
	ng := &NodeGUI{
		cx:         cx,
		c:          c,
		invalidate: make(chan struct{}),
		quit:       cx.KillAll,
	}
	return ng.Run()
}

type NodeGUI struct {
	cx         *conte.Xt
	c          *cli.Context
	w          *f.Window
	th         *p9.Theme
	appWidget  *p9.App
	sidebarButtons []*p9.Clickable
	invalidate chan struct{}
	quit       chan struct{}
}

func (ng *NodeGUI) Run() (err error) {
	ng.th = p9.NewTheme(p9fonts.Collection(), ng.quit)
	ng.th.Colors.SetTheme(ng.th.Dark)
	ng.sidebarButtons = make([]*p9.Clickable, 6)
	for i := range ng.sidebarButtons {
		ng.sidebarButtons[i] = ng.th.Clickable()
	}
	ng.appWidget = ng.GetAppWidget()
	ng.w = f.NewWindow()
	go func() {
		if err := ng.w.
			Size(640, 480).
			Title("parallelcoin node control panel").
			Open().
			Run(
				ng.appWidget.Fn,
				func() {
					Debug("quitting node gui")
					interrupt.Request()
				}); Check(err) {
		}
	}()
	go func() {
	out:
		for {
			select {
			case <-ng.invalidate:
				ng.w.Window.Invalidate()
			case <-ng.quit:
				break out
			}
		}
	}()
	app.Main()
	return
}
