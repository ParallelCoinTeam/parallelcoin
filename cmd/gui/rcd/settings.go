package rcd

import (
	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/cmd/gui/mvc/model"
	"github.com/p9c/pod/pkg/pod"
)

func (rc *RcVar) SaveDaemonCfg() {
	save.Pod(rc.Settings.Daemon.Config)
}

func (rc *RcVar) ComSettings() func() {
	return func() {
		rc.Settings.Daemon = model.DaemonConfig{
			Config: rc.Cx.Config,
			Schema: pod.GetConfigSchema(),
		}
	}
}
