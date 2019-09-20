// +build !nogui
// +build !headless

package vue

import (
	"git.parallelcoin.io/dev/pod/app/save"
	"git.parallelcoin.io/dev/pod/cmd/gui/vue/comp/conf"
	"git.parallelcoin.io/dev/pod/cmd/gui/vue/db"
	"git.parallelcoin.io/dev/pod/cmd/gui/vue/mod"
	"git.parallelcoin.io/dev/pod/pkg/conte"
	"git.parallelcoin.io/dev/pod/pkg/pod"
)

type DuoVUEConfig struct {
	db      db.DuoVUEdb
	Display mod.DisplayConfig `json:"display"`
	Daemon  DaemonConfig      `json:"daemon"`
}

type DaemonConfig struct {
	Config *pod.Config `json:"config"`
	Schema pod.Schema  `json:"schema"`
}

func (d *DuoVUEConfig) SaveDaemonCfg(c pod.Config) {
	*d.Daemon.Config = c
	save.Pod(d.Daemon.Config)
}

func GetCoreCofig(cx *conte.Xt) (c DuoVUEConfig) {
	c.Daemon = DaemonConfig{
		Config: cx.Config,
		Schema: pod.GetConfigSchema(),
	}
	c.Display = mod.DisplayConfig{
		Screens: conf.GetPanels(),
	}
	return c
}
