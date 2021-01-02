package gui

import (
	"gioui.org/app"
	l "gioui.org/layout"
	"github.com/urfave/cli"
	
	"github.com/p9c/pod/pkg/gui"
	qu "github.com/p9c/pod/pkg/util/quit"
	
	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/pkg/gui/cfg"
	"github.com/p9c/pod/pkg/gui/fonts/p9fonts"
	"github.com/p9c/pod/pkg/util/interrupt"
)

func Main(cx *conte.Xt, c *cli.Context) (err error) {
	var size int
	ng := &NodeGUI{
		cx:         cx,
		c:          c,
		invalidate: qu.T(),
		quit:       cx.KillAll,
		size:       &size,
	}
	return ng.Run()
}

type NodeGUI struct {
	cx               *conte.Xt
	c                *cli.Context
	w                *gui.Window
	th               *gui.Theme
	size             *int
	runMode          string
	app              *gui.App
	sidebarButtons   []*gui.Clickable
	buttonBarButtons []*gui.Clickable
	statusBarButtons []*gui.Clickable
	bools            map[string]*gui.Bool
	lists            map[string]*gui.List
	enums            map[string]*gui.Enum
	checkables       map[string]*gui.Checkable
	clickables       map[string]*gui.Clickable
	editors          map[string]*gui.Editor
	inputs           map[string]*gui.Input
	multis           map[string]*gui.Multi
	configs          cfg.GroupsMap
	config           *cfg.Config
	passwords        map[string]*gui.Password
	invalidate       qu.C
	quit             qu.C
}

func (ng *NodeGUI) Run() (err error) {
	ng.th = gui.NewTheme(p9fonts.Collection(), ng.quit)
	ng.th.Colors.SetTheme(*ng.th.Dark)
	ng.runMode = "node"
	ng.sidebarButtons = make([]*gui.Clickable, 9)
	for i := range ng.sidebarButtons {
		ng.sidebarButtons[i] = ng.th.Clickable()
	}
	ng.buttonBarButtons = make([]*gui.Clickable, 4)
	for i := range ng.buttonBarButtons {
		ng.buttonBarButtons[i] = ng.th.Clickable()
	}
	ng.statusBarButtons = make([]*gui.Clickable, 3)
	for i := range ng.statusBarButtons {
		ng.statusBarButtons[i] = ng.th.Clickable()
	}
	ng.enums = map[string]*gui.Enum{
		"runmode": ng.th.Enum().SetValue(ng.runMode),
	}
	ng.bools = map[string]*gui.Bool{
		"runstate": ng.th.Bool(false).SetOnChange(func(b bool) {
			Debug("run state is now", b)
		}),
	}
	ng.lists = map[string]*gui.List{
		"overview": ng.th.List(),
		"settings": ng.th.List(),
	}
	ng.clickables = map[string]*gui.Clickable{
		"quit": ng.th.Clickable(),
	}
	ng.checkables = map[string]*gui.Checkable{
		"runmodenode":   ng.th.Checkable(),
		"runmodewallet": ng.th.Checkable(),
		"runmodeshell":  ng.th.Checkable(),
	}
	ng.editors = make(map[string]*gui.Editor)
	ng.inputs = make(map[string]*gui.Input)
	ng.multis = make(map[string]*gui.Multi)
	ng.passwords = make(map[string]*gui.Password)
	ng.w = gui.NewWindow(ng.th)
	ng.app = ng.GetAppWidget()
	go func() {
		if err := ng.w.
			Size(64, 32).
			Title("parallelcoin node control panel").
			Open().
			Run(
				ng.app.Fn(),
				func(gtx l.Context) {},
				func() {
					Debug("quitting node gui")
					interrupt.Request()
				}, ng.quit); Check(err) {
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
