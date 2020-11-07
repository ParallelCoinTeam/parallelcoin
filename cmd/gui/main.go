package gui

import (
	"github.com/urfave/cli"
	"runtime"
	"time"

	"github.com/p9c/pod/pkg/rpc/btcjson"

	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/pkg/comm/stdconn/worker"
	"github.com/p9c/pod/pkg/gui/cfg"
	"github.com/p9c/pod/pkg/gui/f"
	"github.com/p9c/pod/pkg/gui/fonts/p9fonts"
	"github.com/p9c/pod/pkg/gui/p9"
	rpcclient "github.com/p9c/pod/pkg/rpc/client"
	"github.com/p9c/pod/pkg/util/interrupt"
)

func Main(cx *conte.Xt, c *cli.Context) (err error) {
	var size int
	wg := &WalletGUI{
		cx:         cx,
		c:          c,
		invalidate: make(chan struct{}),
		quit:       cx.KillAll,
		// runnerQuit: make(chan struct{}),
		size: &size,
	}
	return wg.Run()
}

type WalletGUI struct {
	cx   *conte.Xt
	c    *cli.Context
	w    map[string]*f.Window
	th   *p9.Theme
	size *int
	*p9.App
	sidebarButtons            []*p9.Clickable
	buttonBarButtons          []*p9.Clickable
	statusBarButtons          []*p9.Clickable
	bools                     map[string]*p9.Bool
	quitClickable             *p9.Clickable
	lists                     map[string]*p9.List
	checkables                map[string]*p9.Checkable
	clickables                map[string]*p9.Clickable
	inputs                    map[string]*p9.Input
	passwords                 map[string]*p9.Password
	incdecs                   map[string]*p9.IncDec
	configs                   cfg.GroupsMap
	config                    *cfg.Config
	running, mining           bool
	invalidate                chan struct{}
	quit                      chan struct{}
	runnerQuit                chan struct{}
	sendAddresses             []SendAddress
	Worker                    *worker.Worker
	RunCommandChan            chan string
	State                     State
	Shell                     *worker.Worker
	ChainClient, WalletClient *rpcclient.Client
	txs                       []btcjson.ListTransactionsResult
	console                   *Console
}

func (wg *WalletGUI) Run() (err error) {
	wg.th = p9.NewTheme(p9fonts.Collection(), wg.quit)
	wg.th.Dark = wg.cx.Config.DarkTheme
	wg.th.Colors.SetTheme(*wg.th.Dark)
	wg.sidebarButtons = make([]*p9.Clickable, 10)
	for i := range wg.sidebarButtons {
		wg.sidebarButtons[i] = wg.th.Clickable()
	}
	wg.buttonBarButtons = make([]*p9.Clickable, 4)
	for i := range wg.buttonBarButtons {
		wg.buttonBarButtons[i] = wg.th.Clickable()
	}
	wg.statusBarButtons = make([]*p9.Clickable, 4)
	for i := range wg.statusBarButtons {
		wg.statusBarButtons[i] = wg.th.Clickable()
	}
	wg.lists = map[string]*p9.List{
		"createWallet": wg.th.List(),
		"overview":     wg.th.List(),
		"recent":       wg.th.List(),
		"send":         wg.th.List(),
		"transactions": wg.th.List(),
		"settings":     wg.th.List(),
		"received":     wg.th.List(),
		"console":      wg.th.List(),
	}
	wg.clickables = map[string]*p9.Clickable{
		"createWallet":            wg.th.Clickable(),
		"quit":                    wg.th.Clickable(),
		"sendSend":                wg.th.Clickable(),
		"sendClearAll":            wg.th.Clickable(),
		"sendAddRecipient":        wg.th.Clickable(),
		"receiveCreateNewAddress": wg.th.Clickable(),
		"receiveClear":            wg.th.Clickable(),
		"receiveShow":             wg.th.Clickable(),
		"receiveRemove":           wg.th.Clickable(),
		"transactions10":          wg.th.Clickable(),
		"transactions30":          wg.th.Clickable(),
		"transactions50":          wg.th.Clickable(),
	}
	wg.bools = map[string]*p9.Bool{
		"runstate":   wg.th.Bool(wg.running),
		"encryption": wg.th.Bool(false),
		"seed":       wg.th.Bool(false),
		"testnet":    wg.th.Bool(false),
	}
	pass := "password"
	wg.inputs = map[string]*p9.Input{
		"receiveLabel":   wg.th.Input("", "Label", "Primary", "DocText", 25, func(pass string) {}),
		"receiveAmount":  wg.th.Input("", "Amount", "Primary", "DocText", 25, func(pass string) {}),
		"receiveMessage": wg.th.Input("", "Message", "Primary", "DocText", 25, func(pass string) {}),
		"console":        wg.th.Input("", "ParallelCoin console", "Primary", "DocText", 25, func(pass string) {}),
	}
	wg.passwords = map[string]*p9.Password{
		"passEditor":        wg.th.Password(&pass, "Primary", "DocText", 25, func(pass string) {}),
		"confirmPassEditor": wg.th.Password(&pass, "Primary", "DocText", 25, func(pass string) {}),
	}
	wg.console = &Console{
		Commands: []ConsoleCommand{
			{
				ComID:    "input",
				Category: "input",
				Time:     time.Now(),
				// Out: input(duo),
			},
		},
		CommandsNumber: 1,
	}

	wg.w = make(map[string]*f.Window)
	if err = wg.Runner(); Check(err) {
	}
	// wg.RunCommandChan <- "run"
	wg.quitClickable = wg.th.Clickable()
	wg.w = map[string]*f.Window{
		"splash": f.NewWindow(),
		"main":   f.NewWindow(),
	}
	wg.incdecs = map[string]*p9.IncDec{
		"generatethreads": wg.th.IncDec(2, 0, runtime.NumCPU(), *wg.cx.Config.GenThreads,
			func(n int) {
				Debug("threads value now", n)
			},
		),
	}
	wg.App = wg.GetAppWidget()
	wg.Tickers()
	wg.CreateSendAddressItem()
	go func() {
		if err := wg.w["main"].
			Size(800, 480).
			Title("ParallelCoin Wallet").
			Open().
			Run(
				wg.Fn(),
				// wg.InitWallet(),
				func() {
					Debug("quitting wallet gui")
					wg.RunCommandChan <- "stop"
					close(wg.quit)
				}, wg.quit); Check(err) {
		}
	}()
	interrupt.AddHandler(func() {
		Debug("quitting wallet gui")
		wg.RunCommandChan <- "stop"
		close(wg.quit)
	})
out:
	for {
		select {
		case <-wg.invalidate:
			Debug("invalidating render queue")
			wg.w["main"].Window.Invalidate()
		case <-wg.quit:
			Debug("closing GUI on quit signal")
			break out
		}
	}
	// app.Main is just a synonym for select{} so don't do it, we want to be able to shut down
	// app.Main()
	return
}
