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
	cx      *conte.Xt
	c       *cli.Context
	w       *f.Window
	th      *p9.Theme
	size    *int
	runMode string
	app *p9.App
	sidebarButtons   []*p9.Clickable
	buttonBarButtons []*p9.Clickable
	statusBarButtons []*p9.Clickable
	bools            map[string]*p9.Bool
	lists            map[string]*p9.List
	enums            map[string]*p9.Enum
	checkables       map[string]*p9.Checkable
	clickables       map[string]*p9.Clickable
	editors          map[string]*p9.Editor
	inputs           map[string]*p9.Input
	multis           map[string]*p9.Multi
	configs          GroupsMap
	passwords        map[string]*p9.Password
	invalidate       chan struct{}
	quit             chan struct{}
}

func (ng *NodeGUI) Run() (err error) {
	ng.th = p9.NewTheme(p9fonts.Collection(), ng.quit)
	ng.th.Colors.SetTheme(ng.th.Dark)
	ng.runMode = "node"
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
	ng.enums = map[string]*p9.Enum{
		"runmode": ng.th.Enum().SetValue(ng.runMode),
	}
	ng.bools = map[string]*p9.Bool{
		"runstate": ng.th.Bool(false).SetOnChange(func(b bool) {
			Debug("run state is now", b)
		}),
	}
	ng.lists = map[string]*p9.List{
		"overview": ng.th.List(),
		"settings": ng.th.List(),
	}
	ng.clickables = map[string]*p9.Clickable{
		"quit": ng.th.Clickable(),
	}
	ng.checkables = map[string]*p9.Checkable{
		"runmodenode":   ng.th.Checkable(),
		"runmodewallet": ng.th.Checkable(),
		"runmodeshell":  ng.th.Checkable(),
	}
	ng.editors = make(map[string]*p9.Editor)
	ng.inputs = make(map[string]*p9.Input)
	ng.multis = make(map[string]*p9.Multi)
	ng.passwords = make(map[string]*p9.Password)
	ng.w = f.NewWindow()
	ng.app = ng.GetAppWidget()
	go func() {
		if err := ng.w.
			Size(640, 480).
			Title("parallelcoin node control panel").
			Open().
			Run(
				ng.app.Fn(),
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
