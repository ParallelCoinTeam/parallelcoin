// +build !nogui
// +build !headless

package vue

import (
	"github.com/parallelcointeam/parallelcoin/app/save"
	"github.com/parallelcointeam/parallelcoin/cmd/gui/vue/comp/conf"
	"github.com/parallelcointeam/parallelcoin/cmd/gui/vue/db"
	"github.com/parallelcointeam/parallelcoin/cmd/gui/vue/mod"
	"github.com/parallelcointeam/parallelcoin/pkg/conte"
	"github.com/parallelcointeam/parallelcoin/pkg/pod"
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
