package cfg

import (
	"github.com/p9c/gel"
	"github.com/p9c/pod/pkg/pod"
	"github.com/p9c/qu"
)

func New(cx *pod.State, w *gel.Window) *Config {
	cfg := &Config{
		Window: w,
		cx:     cx,
		quit:   cx.KillAll,
	}
	cfg.Theme = cx.App
	return cfg.Init()
}

type Config struct {
	cx *pod.State
	*gel.Window
	Bools      map[string]*gel.Bool
	lists      map[string]*gel.List
	enums      map[string]*gel.Enum
	checkables map[string]*gel.Checkable
	clickables map[string]*gel.Clickable
	editors    map[string]*gel.Editor
	inputs     map[string]*gel.Input
	multis     map[string]*gel.Multi
	configs    GroupsMap
	passwords  map[string]*gel.Password
	quit       qu.C
}

func (c *Config) Init() *Config {
	c.Theme.Colors.SetTheme(*c.Theme.Dark)
	c.enums = map[string]*gel.Enum{
		// "runmode": ng.th.Enum().SetValue(ng.runMode),
	}
	c.Bools = map[string]*gel.Bool{
		// "runstate": ng.th.Bool(false).SetOnChange(func(b bool) {
		// 	D.Ln("run state is now", b)
		// }),
	}
	c.lists = map[string]*gel.List{
		// "overview": ng.th.List(),
		"settings": c.List(),
	}
	c.clickables = map[string]*gel.Clickable{
		// "quit": ng.th.Clickable(),
	}
	c.checkables = map[string]*gel.Checkable{
		// "runmodenode":   ng.th.Checkable(),
		// "runmodewallet": ng.th.Checkable(),
		// "runmodeshell":  ng.th.Checkable(),
	}
	c.editors = make(map[string]*gel.Editor)
	c.inputs = make(map[string]*gel.Input)
	c.multis = make(map[string]*gel.Multi)
	c.passwords = make(map[string]*gel.Password)
	return c
}
