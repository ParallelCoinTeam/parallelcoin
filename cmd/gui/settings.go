//+build !headless

package gui

import (
	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/pod"
)

type DuOSConfig struct {
	db DuOSdb
	//Display mod.DisplayConfig `json:"display"`
	Daemon DaemonConfig `json:"daemon"`
}

type DaemonConfig struct {
	Config *pod.Config `json:"config"`
	Schema pod.Schema  `json:"schema"`
}

func (d *DuOSConfig) SaveDaemonCfg(c pod.Config) {
	*d.Daemon.Config = c
	save.Pod(d.Daemon.Config)
}

func GetCoreCofig(cx *conte.Xt) (c DuOSConfig) {
	c.Daemon = DaemonConfig{
		Config: cx.Config,
		Schema: pod.GetConfigSchema(),
	}
	//c.Display = mod.DisplayConfig{
	//	Screens: conf.GetPanels(),
	//}
	return c
}
