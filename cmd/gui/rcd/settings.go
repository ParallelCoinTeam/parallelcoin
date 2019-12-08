package rcd

import (
	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/pod"
)

type DuOSsettings struct {
	cx *conte.Xt
	db DuoUIdb
	//Display mod.DisplayConfig `json:"display"`
	Daemon DaemonConfig `json:"daemon"`
}

type DaemonConfig struct {
	Config *pod.Config `json:"config"`
	Schema pod.Schema  `json:"schema"`
}

func (d *DuOSsettings) SaveDaemonCfg(c pod.Config) {
	*d.Daemon.Config = c
	save.Pod(d.Daemon.Config)
}

func (d *DuOSsettings)GetCoreSettings() {
	d.Daemon = DaemonConfig{
		Config: d.cx.Config,
		Schema: pod.GetConfigSchema(),
	}
	//c.Display = mod.DisplayConfig{
	//	Screens: conf.GetPanels(),
	//}
	return
}
