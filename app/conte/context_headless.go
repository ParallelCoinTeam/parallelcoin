// +build headless

package conte

import (
	"github.com/p9c/pod/app/appdata"
	"github.com/p9c/pod/cmd/node/state"
	"github.com/p9c/pod/pkg/pod"
	"github.com/p9c/pod/pkg/util/lang"
	"github.com/urfave/cli"
)

// GetNewContext returns a fresh new context
func GetNewContext(appName, appLang, subtext string) *Xt {
	ec, _ := pod.EmptyConfig()
	return &Xt{
		App:      cli.NewApp(),
		Config:   ec,
		StateCfg: new(state.Config),
		Language: lang.ExportLanguage(appLang),
		DataDir:  appdata.Dir(appName, false),
	}
}
