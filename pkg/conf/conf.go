package conf

import (
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/pod"
)

type DuOSconfig struct {
	//db db.DuOSdb
	//Display mod.DisplayConfig `json:"display"`
	Pod PodConfig `json:"pod"`
}

type PodConfig struct {
	Config *pod.Config `json:"config"`
	Schema pod.Schema  `json:"schema"`
}

func (d *DuOSconfig) SaveDaemonCfg(c pod.Config) {
	*d.Pod.Config = c
	//save.Pod(d.Daemon.Config)
}

func (c *DuOSconfig) GetCoreCofig(cx *conte.Xt) *DuOSconfig {
	c.Pod = PodConfig{
		Config: cx.Config,
		Schema: pod.GetConfigSchema(),
	}
	//c.Display = mod.DisplayConfig{
	//	//Screens: conf.GetPanels(),
	//}
	return c
}
