package rcd

import (
	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/pod"
)

func SaveDaemonCfg(c *pod.Config) {
	save.Pod(c)
}

func GetCoreSettings(cx *conte.Xt) models.DaemonConfig {
	return models.DaemonConfig{
		Config: cx.Config,
		Schema: pod.GetConfigSchema(),
	}
	//c.Display = mod.DisplayConfig{
	//	Screens: conf.GetPanels(),
	//}

}
