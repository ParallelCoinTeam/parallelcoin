// +build !headless

package conte

import (
	"go.uber.org/atomic"

	"github.com/urfave/cli"

	"github.com/p9c/pod/app/appdata"
	"github.com/p9c/pod/cmd/node/state"
	"github.com/p9c/pod/pkg/pod"
	"github.com/p9c/pod/pkg/util/lang"
)

// GetNewContext returns a fresh new context
func GetNewContext(appName, appLang, subtext string) *Xt {
	hr := &atomic.Value{}
	hr.Store(int(0))
	config, configMap := pod.EmptyConfig()
	chainClientReady := make(chan struct{})
	return &Xt{
		ChainClientReady: chainClientReady,
		KillAll:          make(chan struct{}),
		App:              cli.NewApp(),
		Config:           config,
		ConfigMap:        configMap,
		StateCfg:         new(state.Config),
		Language:         lang.ExportLanguage(appLang),
		DataDir:          appdata.Dir(appName, false),
	}
}
