// +build !nogui

package gui

import (
	"github.com/parallelcointeam/parallelcoin/cmd/gui/vue"
	"github.com/parallelcointeam/parallelcoin/pkg/conte"
	"github.com/parallelcointeam/parallelcoin/pkg/util/cl"
	"github.com/parallelcointeam/parallelcoin/pkg/util/interrupt"
	"github.com/robfig/cron"
	"sync"
	"sync/atomic"
)

func Main(cx *conte.Xt, wg *sync.WaitGroup) {
	cr := cron.New()
	log <- cl.Warn{"starting gui", cl.Ine()}
	dV := vue.GetDuoVUE(cx, cr)
	cleaned := &atomic.Value{}
	cleaned.Store(false)
	cleanup := func() {
		if !cleaned.Load().(bool) {
			cleaned.Store(true)
			log <- cl.Debug{"terminating webview", cl.Ine()}
			dV.Web.Terminate()
			interrupt.Request()
			log <- cl.Debug{"waiting for waitgroup", cl.Ine()}
			wg.Wait()
			log <- cl.Debug{"exiting webview", cl.Ine()}
			dV.Web.Exit()
		}
	}
	interrupt.AddHandler(func() {
		cleanup()
	})
	defer cleanup()
	dV.Web.Dispatch(func() {
		cr.Start()
		vue.RunVue(*dV)

	})
	dV.Web.Run()
}
