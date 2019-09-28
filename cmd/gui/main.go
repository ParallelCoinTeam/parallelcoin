// +build !nogui

package gui

import (
	"github.com/p9c/pod/cmd/gui/core"
	"github.com/robfig/cron"
	"sync"
	"sync/atomic"

	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/util/interrupt"

	"github.com/p9c/pod/pkg/log"
)

func Main(cx *conte.Xt, wg *sync.WaitGroup) {
	cr := cron.New()
	d := core.MountDuOS(cx, cr)
	log.WARN("starting gui")
	cleaned := &atomic.Value{}
	cleaned.Store(false)
	cleanup := func() {
		if !cleaned.Load().(bool) {
			cleaned.Store(true)
			log.DEBUG("terminating webview")
			d.Wv.Terminate()
			interrupt.Request()
			log.DEBUG("waiting for waitgroup")
			wg.Wait()
			log.DEBUG("exiting webview")
			d.Wv.Exit()
		}
	}
	interrupt.AddHandler(func() {
		cleanup()
	})
	defer cleanup()
	d.Wv.Dispatch(func() {
		cr.Start()
		core.RunDuOS(*d)
	})
	d.Wv.Run()
}
