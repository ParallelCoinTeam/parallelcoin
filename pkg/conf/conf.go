package conf

import (
	conte2 "github.com/p9c/pod/gui/____BEZI/test/pkg/conte"
	pod2 "github.com/p9c/pod/gui/____BEZI/test/pkg/pod"
)

type DuOSconfig struct {
	//db db.DuOSdb
	//Display mod.DisplayConfig `json:"display"`
	Pod PodConfig `json:"pod"`
}

type PodConfig struct {
	Config *pod2.Config `json:"config"`
	Schema pod2.Schema  `json:"schema"`
}

func (d *DuOSconfig) SaveDaemonCfg(c pod2.Config) {
	*d.Pod.Config = c
	//save.Pod(d.Daemon.Config)
}

func (c *DuOSconfig) GetCoreCofig(cx *conte2.Xt) *DuOSconfig {
	c.Pod = PodConfig{
		Config: cx.Config,
		Schema: pod2.GetConfigSchema(),
	}
	//c.Display = mod.DisplayConfig{
	//	//Screens: conf.GetPanels(),
	//}
	return c
}
