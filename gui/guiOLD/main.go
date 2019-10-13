// +build !nogui

package guiOLD

import (
	"github.com/golang-ui/nuklear/nk"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/util/interrupt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/p9c/pod/pkg/log"
)

func init() {
	runtime.LockOSThread()
}

func Main(cx *conte.Xt, wg *sync.WaitGroup) {
	//cr := cron.New()
	//d := core.MountDuOS(cx, cr)
	log.WARN("starting guiOLD")
	cleaned := &atomic.Value{}
	cleaned.Store(false)
	cleanup := func() {
		if !cleaned.Load().(bool) {
			cleaned.Store(true)
			log.DEBUG("terminating webview")
			//d.Wv.Terminate()
			interrupt.Request()
			log.DEBUG("waiting for waitgroup")
			wg.Wait()
			log.DEBUG("exiting webview")
			//d.Wv.Exit()
		}
	}
	interrupt.AddHandler(func() {
		cleanup()
	})
	defer cleanup()

	// initialize glfw
	glfw.Init()
	// create main windows
	win, _ := glfw.CreateWindow(120, 200, "GoBkm", nil, nil)
	// make the new window current so that its context is used in the following code
	win.MakeContextCurrent()
	// initialize opengl context
	gl.Init()
	// create context
	ctx := nk.NkPlatformInit(win, nk.PlatformInstallCallbacks)
	// set default font
	atlas := nk.NewFontAtlas()
	nk.NkFontStashBegin(&atlas)
	font := nk.NkFontAtlasAddDefault(atlas, 18, nil)
	nk.NkFontStashEnd()
	nk.NkStyleSetFont(ctx, font.Handle())

	quit := make(chan struct{}, 1)
	ticker := time.NewTicker(time.Second / 30)
	// loop that handles GUI refreshing and event management
	for {
		select {
		case <-quit:
			nk.NkPlatformShutdown()
			glfw.Terminate()
			ticker.Stop()
			return
		case <-ticker.C:
			if win.ShouldClose() {
				close(quit)
				continue
			}
			glfw.PollEvents()
			// GUI definition
			draw(win, ctx)
		}
	}
}
