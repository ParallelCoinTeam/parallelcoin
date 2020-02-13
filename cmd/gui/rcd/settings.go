package rcd

import (
	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/cmd/gui/mvc/model"
	"github.com/p9c/pod/pkg/pod"
)

func (r *RcVar) SaveDaemonCfg() {
	save.Pod(r.Settings.Daemon.Config)
}

func (r *RcVar) ComSettings() func() {
	return func() {
		r.Settings.Daemon = model.DaemonConfig{
			Config: r.Cx.Config,
			Schema: pod.GetConfigSchema(),
		}
	}
}
