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
	var size int
	ng := &NodeGUI{
		cx:         cx,
		c:          c,
		invalidate: make(chan struct{}),
		quit:       cx.KillAll,
		size:       &size,
	}
	return ng.Run()
}

type NodeGUI struct {
	cx   *conte.Xt
	c    *cli.Context
	w    *f.Window
	th   *p9.Theme
	size *int
	*p9.App
	sidebarButtons   []*p9.Clickable
	buttonBarButtons []*p9.Clickable
	statusBarButtons []*p9.Clickable
	bools            map[string]*p9.Bool
	lists            map[string]*p9.List
	quitClickable    *p9.Clickable
	invalidate       chan struct{}
	quit             chan struct{}
}

func (ng *NodeGUI) Run() (err error) {
	ng.th = p9.NewTheme(p9fonts.Collection(), ng.quit)
	ng.th.Colors.SetTheme(ng.th.Dark)
	ng.sidebarButtons = make([]*p9.Clickable, 9)
	for i := range ng.sidebarButtons {
		ng.sidebarButtons[i] = ng.th.Clickable()
	}
	ng.buttonBarButtons = make([]*p9.Clickable, 4)
	for i := range ng.buttonBarButtons {
		ng.buttonBarButtons[i] = ng.th.Clickable()
	}
	ng.statusBarButtons = make([]*p9.Clickable, 3)
	for i := range ng.statusBarButtons {
		ng.statusBarButtons[i] = ng.th.Clickable()
	}
	ng.bools = map[string]*p9.Bool{
		"runstate": ng.th.Bool(false).SetOnChange(func(b bool) {
			Debug("run state is now", b)
		}),
	}
	ng.lists = map[string]*p9.List{
		"overview": ng.th.List(),
	}
	ng.quitClickable = ng.th.Clickable()
	ng.w = f.NewWindow()
	ng.App = ng.GetAppWidget()
	go func() {
		if err := ng.w.
			Size(640, 480).
			Title("parallelcoin node control panel").
			Open().
			Run(
				ng.Fn(),
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
