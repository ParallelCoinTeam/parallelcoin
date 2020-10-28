package cfg

import (
	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/pkg/gui/p9"
)

func New(cx *conte.Xt, th *p9.Theme) *Config {
	cfg := &Config{
		th:   th,
		cx:   cx,
		quit: cx.KillAll,
	}
	return cfg.Init()
}

type Config struct {
	cx         *conte.Xt
	th         *p9.Theme
	bools      map[string]*p9.Bool
	lists      map[string]*p9.List
	enums      map[string]*p9.Enum
	checkables map[string]*p9.Checkable
	clickables map[string]*p9.Clickable
	editors    map[string]*p9.Editor
	inputs     map[string]*p9.Input
	multis     map[string]*p9.Multi
	configs    GroupsMap
	passwords  map[string]*p9.Password
	quit       chan struct{}
}

func (c *Config) Init() *Config {
	// c.th = p9.NewTheme(p9fonts.Collection(), c.cx.KillAll)
	c.th.Colors.SetTheme(c.th.Dark)
	c.enums = map[string]*p9.Enum{
		// "runmode": ng.th.Enum().SetValue(ng.runMode),
	}
	c.bools = map[string]*p9.Bool{
		// "runstate": ng.th.Bool(false).SetOnChange(func(b bool) {
		// 	Debug("run state is now", b)
		// }),
	}
	c.lists = map[string]*p9.List{
		// "overview": ng.th.List(),
		"settings": c.th.List(),
	}
	c.clickables = map[string]*p9.Clickable{
		// "quit": ng.th.Clickable(),
	}
	c.checkables = map[string]*p9.Checkable{
		// "runmodenode":   ng.th.Checkable(),
		// "runmodewallet": ng.th.Checkable(),
		// "runmodeshell":  ng.th.Checkable(),
	}
	c.editors = make(map[string]*p9.Editor)
	c.inputs = make(map[string]*p9.Input)
	c.multis = make(map[string]*p9.Multi)
	c.passwords = make(map[string]*p9.Password)
	return c
}
