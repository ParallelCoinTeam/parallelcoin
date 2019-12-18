package loader

import (
	"github.com/p9c/pod/pkg/gio/app"
	"github.com/p9c/pod/cmd/gui/duoui"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/log"
)

func DuoUIloader(duo *duoui.DuoUI, cx *conte.Xt, firstRun bool) (err error) {
	go func() {
		if err := DuoUIloaderLoop(cx); err != nil {
			log.FATAL(err)
		}
	}()
	app.Main()
	return

	//defer cx.Gui.Wv.Exit()
	//cx.Gui.Wv.Dispatch(func() {
	//
	//	_, err = cx.Gui.Wv.Bind("duos", &rcvar{
	//		cx:         cx,
	//		IsFirstRun: firstRun,
	//	})
	//
	//
	//	cx.Gui.Wv.InjectCSS(getFile("css/theme/root.css", *cx.Gui.Fs))
	//	cx.Gui.Wv.InjectCSS(getFile("css/theme/colors.css", *cx.Gui.Fs))
	//	cx.Gui.Wv.InjectCSS(getFile("css/theme/helpers.css", *cx.Gui.Fs))
	//	cx.Gui.Wv.InjectCSS(getFile("css/loader.css", *cx.Gui.Fs))
	//	cx.Gui.Wv.InjectCSS(getFile("css/svelte.css", *cx.Gui.Fs))

	// Load CSS
	//})
	//cx.Gui.Wv.Run()

	//
	//go func() {
	//	for _ = range time.NewTicker(time.Second * 1).C {
	//
	//
	//		//status, err := json.Marshal(rc.GetDuOStatus())
	//		//if err != nil {
	//		//}
	//		//transactions, err := json.Marshal(rc.GetTransactions(0, 555, ""))
	//		//if err != nil {
	//		//}
	//}
	//}()
	return
}
