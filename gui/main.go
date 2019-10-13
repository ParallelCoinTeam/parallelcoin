package gui

import (
	"fmt"
	"github.com/p9c/lorca"
	"github.com/p9c/pod/pkg/duos/core"
	"github.com/p9c/pod/pkg/vue/lib"
	"github.com/p9c/pod/pkg/vue/lib/css"
	"github.com/p9c/pod/pkg/vue/lib/html"
	"log"
	"net/url"
)

func Main(d *core.DuOS) {

	// Create UI with basic HTML passed via data URI
	// string(html.VUEHTML(html.VUEx(lib.VUElogo(), html.VUEheader(), html.VUEnav(lib.ICO()), html.ScreenOverview()), css.CSS(css.ROOT(), css.GRID(), css.COLORS(), css.HELPERS(), css.NAV())))))
	ui, err := lorca.New("data:text/html,"+url.PathEscape(string(html.VUEHTML(html.VUEx(lib.VUElogo(), html.VUEheader(), html.VUEnav(lib.ICO()), html.ScreenOverview()), css.CSS(css.ROOT(), css.GRID(), css.COLORS(), css.HELPERS(), css.NAV())))), "", 800, 600)
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()
	// Wait until UI window is closed
	<-ui.Done()
	// Load HTML after Go functions are bound to JS

	fmt.Println("teeeeeee")

	// Start ticker goroutine
	//go func() {
	//	t := time.NewTicker(100 * time.Millisecond)
	//	for {
	//		select {
	//		case <-t.C: // Every 100ms increate number of ticks and update UI
	//			ui.Eval("document.querySelector('.timer').innerText =" + fmt.Sprint(getBlockCount))
	//		case <-togglec: // If paused - wait for another toggle event to unpause
	//			<-togglec
	//		}
	//	}
	//}()
	//<-ui.Done()
}
