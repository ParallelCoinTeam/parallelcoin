package gui

import (
	"gioui.org/app"
	l "gioui.org/layout"
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
	th         *p9.Theme
	invalidate chan struct{}
	quit       chan struct{}
}

func (ng *NodeGUI) Run() (err error) {
	ng.th = p9.NewTheme(p9fonts.Collection(), ng.quit)
	win := f.Window()
	go func() {
		if err := win.
			Size(640, 480).
			Title("parallelcoin node control panel").
			Open().
			Run(
				func(gtx l.Context) l.Dimensions {
					return l.Dimensions{}
				},
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
				win.Window.Invalidate()
			case <-ng.quit:
				break out
			}
		}
	}()
	app.Main()
	return
}
