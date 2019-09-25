// +build !nogui

package gui

import (
	"sync"
	"sync/atomic"

	"github.com/robfig/cron"

	"github.com/p9c/pod/cmd/gui/vue"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/util/interrupt"

	"github.com/p9c/pod/pkg/log"
)

func Main(cx *conte.Xt, wg *sync.WaitGroup) {
	cr := cron.New()
	log.WARN("starting gui")
	dV := vue.GetDuoVUE(cx, cr)
	cleaned := &atomic.Value{}
	cleaned.Store(false)
	cleanup := func() {
		if !cleaned.Load().(bool) {
			cleaned.Store(true)
			log.DEBUG("terminating webview")
			dV.Web.Terminate()
			interrupt.Request()
			log.DEBUG("waiting for waitgroup")
			wg.Wait()
			log.DEBUG("exiting webview")
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
