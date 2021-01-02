package cfg

import (
	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/pkg/gui"
	qu "github.com/p9c/pod/pkg/util/quit"
)

func New(cx *conte.Xt, w *gui.Window) *Config {
	cfg := &Config{
		Window: w,
		cx:    cx,
		quit:  cx.KillAll,
	}
	return cfg.Init()
}

type Config struct {
	cx         *conte.Xt
	*gui.Window
	Bools      map[string]*gui.Bool
	lists      map[string]*gui.List
	enums      map[string]*gui.Enum
	checkables map[string]*gui.Checkable
	clickables map[string]*gui.Clickable
	editors    map[string]*gui.Editor
	inputs     map[string]*gui.Input
	multis     map[string]*gui.Multi
	configs    GroupsMap
	passwords  map[string]*gui.Password
	quit       qu.C
}

func (c *Config) Init() *Config {
	// c.th = p9.NewTheme(p9fonts.Collection(), c.cx.KillAll)
	c.Theme.Colors.SetTheme(*c.Theme.Dark)
	c.enums = map[string]*gui.Enum{
		// "runmode": ng.th.Enum().SetValue(ng.runMode),
	}
	c.Bools = map[string]*gui.Bool{
		// "runstate": ng.th.Bool(false).SetOnChange(func(b bool) {
		// 	Debug("run state is now", b)
		// }),
	}
	c.lists = map[string]*gui.List{
		// "overview": ng.th.List(),
		"settings": c.List(),
	}
	c.clickables = map[string]*gui.Clickable{
		// "quit": ng.th.Clickable(),
	}
	c.checkables = map[string]*gui.Checkable{
		// "runmodenode":   ng.th.Checkable(),
		// "runmodewallet": ng.th.Checkable(),
		// "runmodeshell":  ng.th.Checkable(),
	}
	c.editors = make(map[string]*gui.Editor)
	c.inputs = make(map[string]*gui.Input)
	c.multis = make(map[string]*gui.Multi)
	c.passwords = make(map[string]*gui.Password)
	return c
}
