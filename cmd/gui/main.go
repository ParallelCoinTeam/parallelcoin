package gui

import (
	"crypto/rand"
	"encoding/hex"
	"runtime"
	"sync"
	"time"

	"github.com/urfave/cli"
	uberatomic "go.uber.org/atomic"

	l "gioui.org/layout"

	"github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/pkg/gui/dialog"
	"github.com/p9c/pod/pkg/gui/toast"
	"github.com/p9c/pod/pkg/util/hdkeychain"

	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/pkg/comm/stdconn/worker"
	"github.com/p9c/pod/pkg/gui/cfg"
	"github.com/p9c/pod/pkg/gui/f"
	"github.com/p9c/pod/pkg/gui/fonts/p9fonts"
	"github.com/p9c/pod/pkg/gui/p9"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	rpcclient "github.com/p9c/pod/pkg/rpc/client"
	"github.com/p9c/pod/pkg/util/interrupt"
)

func Main(cx *conte.Xt, c *cli.Context) (err error) {
	var size int
	var noWallet, walletLocked bool
	wg := &WalletGUI{
		cx:         cx,
		c:          c,
		invalidate: make(chan struct{}),
		quit:       cx.KillAll,
		// runnerQuit: make(chan struct{}),
		size:         &size,
		noWallet:     &noWallet,
		walletLocked: uberatomic.NewBool(walletLocked),
	}
	return wg.Run()
}

type WalletGUI struct {
	cx               *conte.Xt
	c                *cli.Context
	w                map[string]*f.Window
	th               *p9.Theme
	size             *int
	App              *p9.App
	unlockPage       *p9.App
	unlockPassword   *p9.Password
	sidebarButtons   []*p9.Clickable
	buttonBarButtons []*p9.Clickable
	statusBarButtons []*p9.Clickable
	quitClickable    *p9.Clickable
	bools            map[string]*p9.Bool
	lists            map[string]*p9.List
	checkables       map[string]*p9.Checkable
	clickables       map[string]*p9.Clickable
	inputs           map[string]*p9.Input
	passwords        map[string]*p9.Password
	incdecs          map[string]*p9.IncDec
	// intSliders                map[string]*p9.IntSlider
	configs                            cfg.GroupsMap
	config                             *cfg.Config
	runningNode, runningWallet, mining uberatomic.Bool
	invalidate                         chan struct{}
	quit                               chan struct{}
	runnerQuit                         chan struct{}
	minerQuit                          chan struct{}
	sendAddresses                      []SendAddress
	NodeRunCommandChan                 chan string
	WalletRunCommandChan               chan string
	MinerRunCommandChan                chan string
	State                              State
	Node, Wallet, Miner                *worker.Worker
	ChainMutex, WalletMutex            sync.Mutex
	ChainClient, WalletClient          *rpcclient.Client
	txs                                []btcjson.ListTransactionsResult
	historyCurPage                     int
	console                            *Console
	toasts                             *toast.Toasts
	dialog                             *dialog.Dialog
	noWallet                           *bool
	walletLocked                       *uberatomic.Bool
	walletToLock                       time.Time
	walletLockTime                     int
	Size                               *int
	historyTable                       *p9.TextTable
}

func (wg *WalletGUI) Run() (err error) {
	wg.th = p9.NewTheme(p9fonts.Collection(), wg.quit)
	wg.th.Dark = wg.cx.Config.DarkTheme
	wg.th.Colors.SetTheme(*wg.th.Dark)
	wg.sidebarButtons = make([]*p9.Clickable, 12)
	wl := true
	wg.walletLocked.Store(wl)
	for i := range wg.sidebarButtons {
		wg.sidebarButtons[i] = wg.th.Clickable()
	}
	wg.buttonBarButtons = make([]*p9.Clickable, 5)
	for i := range wg.buttonBarButtons {
		wg.buttonBarButtons[i] = wg.th.Clickable()
	}
	wg.statusBarButtons = make([]*p9.Clickable, 6)
	for i := range wg.statusBarButtons {
		wg.statusBarButtons[i] = wg.th.Clickable()
	}
	// wg.intSliders = map[string]*p9.IntSlider{
	// 	"lockTimeout": wg.th.IntSlider().
	// 		// TextColor("Danger").
	// 		// SliderColor("Danger").
	// 		Min(30).Max(3600).Value(300).
	// 		Hook(func(v int) {
	// 			wg.walletLockTime = v
	// 			wg.walletToLock = time.Now().Add(time.Duration(v) * time.Second)
	// 			// wg.intSliders["lockTimeout"].Value(v)
	// 		}),
	// }
	wg.State.AllTimeStrings.Store([]string{})
	wg.lists = map[string]*p9.List{
		"createWallet": wg.th.List(),
		"overview":     wg.th.List(),
		"recent":       wg.th.List(),
		"send":         wg.th.List(),
		"transactions": wg.th.List(),
		"settings":     wg.th.List(),
		"received":     wg.th.List(),
		"history":      wg.th.List(),
	}
	wg.historyTable = (&p9.TextTable{
		Theme:            wg.th,
		HeaderColor:      "DocText",
		HeaderBackground: "DocBg",
		HeaderFont:       "bariol bold",
		HeaderFontScale:  1,
		CellColor:        "PanelText",
		CellBackground:   "PanelBg",
		CellFont:         "go regular",
		CellFontScale:    p9.Scales["Caption"],
		Inset:            0.25,
		List:             wg.lists["history"],
	}).
		SetDefaults()
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
		"txPageForward":           wg.th.Clickable(),
		"txPageBack":              wg.th.Clickable(),
	}
	wg.checkables = map[string]*p9.Checkable{
	}
	wg.bools = map[string]*p9.Bool{
		"runstate":     wg.th.Bool(wg.runningNode.Load()),
		"encryption":   wg.th.Bool(false),
		"seed":         wg.th.Bool(false),
		"testnet":      wg.th.Bool(false),
		"ihaveread":    wg.th.Bool(false),
		"showGenerate": wg.th.Bool(true),
		"showSent":     wg.th.Bool(true),
		"showReceived": wg.th.Bool(true),
		"showImmature": wg.th.Bool(true),
	}
	pass := ""
	passConfirm := ""
	seed := make([]byte, hdkeychain.MaxSeedBytes)
	_, _ = rand.Read(seed)
	seedString := hex.EncodeToString(seed)
	wg.inputs = map[string]*p9.Input{
		"receiveLabel":   wg.th.Input("", "Label", "Primary", "DocText", 32, func(pass string) {}),
		"receiveAmount":  wg.th.Input("", "Amount", "Primary", "DocText", 32, func(pass string) {}),
		"receiveMessage": wg.th.Input("", "Message", "Primary", "DocText", 32, func(pass string) {}),
		"console":        wg.th.Input("", "enter rpc command", "Primary", "DocText", 32, func(pass string) {}),
		"walletSeed":     wg.th.Input(seedString, "wallet seed", "Primary", "DocText", 32, func(pass string) {}),
	}
	wg.passwords = map[string]*p9.Password{
		"passEditor":        wg.th.Password("password", &pass, "Primary", "DocText", 32, func(pass string) {}),
		"confirmPassEditor": wg.th.Password("confirm", &passConfirm, "Primary", "DocText", 32, func(pass string) {}),
		"publicPassEditor":  wg.th.Password("public password (optional)", wg.cx.Config.WalletPass, "Primary", "DocText", 32, func(pass string) {}),
	}
	wg.toasts = toast.New(wg.th)
	wg.dialog = dialog.New(wg.th)
	wg.console = wg.ConsolePage()
	wg.w = make(map[string]*f.Window)
	wg.quitClickable = wg.th.Clickable()
	wg.w = map[string]*f.Window{
		"splash": f.NewWindow(wg.th),
		"main":   f.NewWindow(wg.th),
	}
	wg.incdecs = map[string]*p9.IncDec{
		"generatethreads": wg.th.IncDec().
			NDigits(2).
			Min(0).
			Max(runtime.NumCPU()).
			SetCurrent(*wg.cx.Config.GenThreads).
			ChangeHook(
				func(n int) {
					Debug("threads value now", n)
					go func() {
						Debug("setting thread count")
						*wg.cx.Config.GenThreads = n
						save.Pod(wg.cx.Config)
						// wg.MinerThreadsChan <- n
						if wg.mining.Load() {
							Debug("restarting miner")
							wg.MinerRunCommandChan <- "stop"
							wg.MinerRunCommandChan <- "run"
						}
					}()
				},
			),
		"transactionsPerPage": wg.th.IncDec().
			Min(10).
			Max(100).
			NDigits(3).
			Amount(10).
			SetCurrent(10).
			ChangeHook(func(n int) {
				Debug("showing", n, "per page")
			}),
		"idleTimeout": wg.th.IncDec().
			Scale(4).
			Min(60).
			Max(3600).
			NDigits(4).
			Amount(60).
			SetCurrent(300).
			ChangeHook(func(n int) {
				Debug("idle timeout", time.Duration(n)*time.Second)
			}),
	}
	// wg.Subscriber()
	wg.App = wg.GetAppWidget()
	wg.unlockPage = wg.getWalletUnlockAppWidget()
	wg.Tickers()
	wg.CreateSendAddressItem()
	wg.runningNode.Store(!(*wg.cx.Config.NodeOff)) // || *wg.cx.Config.WalletOff)
	wg.mining.Store(*wg.cx.Config.Generate && *wg.cx.Config.GenThreads != 0)
	if !apputil.FileExists(*wg.cx.Config.WalletFile) {
		*wg.noWallet = true
		wg.runningNode.Store(false)
		wg.mining.Store(false)
		wg.inputs["walletseed"] = wg.th.Input("", "wallet seed", "Primary", "DocText", 25, func(pass string) {})
	} else {
		if err = wg.Runner(); Check(err) {
		}
	}
	if wg.runningNode.Load() {
		Debug("initial starting shell")
		wg.runningNode.Store(false)
		wg.NodeRunCommandChan <- "run"
	}
	if wg.mining.Load() {
		// wg.MinerThreadsChan <- *wg.cx.Config.GenThreads
		Debug("initial starting miner")
		wg.mining.Store(false)
		wg.MinerRunCommandChan <- "run"
	}
	wg.Size = wg.w["main"].Width
	go func() {
		if err := wg.w["main"].
			Size(64, 32).
			Title("ParallelCoin Wallet").
			Open().
			Run(
				func(gtx l.Context) l.Dimensions {
					return p9.If(*wg.noWallet,
						wg.CreateWalletPage,
						p9.If(wg.walletLocked.Load(),
							wg.unlockPage.Fn(),
							wg.App.Fn(),
						),
					)(gtx)
				},
				wg.App.Overlay(),
				// wg.InitWallet(),
				func() {
					Debug("quitting wallet gui")
					if wg.runningNode.Load() {
						// consume.Kill(wg.Node)
						close(wg.Node.Quit)
					}
					if wg.mining.Load() {
						// consume.Kill(wg.Miner)
						close(wg.Miner.Quit)
					}
					// close(wg.quit)
				}, wg.quit); Check(err) {
		}
	}()
	interrupt.AddHandler(func() {
		Debug("quitting wallet gui")
		// consume.Kill(wg.Node)
		// consume.Kill(wg.Miner)
		// close(wg.quit)
	})
out:
	for {
		select {
		case <-wg.invalidate:
			// Debug("invalidating render queue")
			wg.w["main"].Window.Invalidate()
		case <-wg.quit:
			Debug("closing GUI on quit signal")
			Debug("disconnecting chain client")
			wg.ChainMutex.Lock()
			if wg.ChainClient != nil {
				wg.ChainClient.Disconnect()
				if wg.ChainClient.Disconnected() {
					wg.ChainClient = nil
				}
			}
			wg.ChainMutex.Unlock()
			Debug("disconnecting wallet client")
			wg.WalletMutex.Lock()
			if wg.WalletClient != nil {
				wg.WalletClient.Disconnect()
				if wg.WalletClient.Disconnected() {
					wg.WalletClient = nil
				}
			}
			wg.WalletMutex.Unlock()
			// if wg.Node != nil {
			//	Debug("stopping shell")
			//	// wg.NodeRunCommandChan <- "stop"
			//	consume.Kill(wg.Node)
			// }
			// if wg.Miner != nil {
			//	Debug("stopping miner")
			//	consume.Kill(wg.Miner)
			//	// wg.MinerRunCommandChan <- "stop"
			// }
			break out
		}
	}
	// app.Main is just a synonym for select{} so don't do it, we want to be able to shut down
	// app.Main()
	return
}
