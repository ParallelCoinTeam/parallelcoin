package gui

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"runtime"
	"sync"
	"time"
	
	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/pkg/util/interrupt"
	log "github.com/p9c/pod/pkg/util/logi"
	qu "github.com/p9c/pod/pkg/util/quit"
	
	"github.com/urfave/cli"
	
	l "gioui.org/layout"
	
	"github.com/p9c/pod/pkg/util/logi/pipe/consume"
	"github.com/p9c/pod/pkg/util/rununit"
	
	"github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/pkg/util/hdkeychain"
	
	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/pkg/gui/cfg"
	"github.com/p9c/pod/pkg/gui/f"
	"github.com/p9c/pod/pkg/gui/fonts/p9fonts"
	"github.com/p9c/pod/pkg/gui/p9"
	rpcclient "github.com/p9c/pod/pkg/rpc/client"
)

func Main(cx *conte.Xt, c *cli.Context) (err error) {
	var size int
	noWallet := true
	wg := &WalletGUI{
		cx:         cx,
		c:          c,
		invalidate: qu.Ts(16),
		quit:       cx.KillAll,
		Size:       &size,
		noWallet:   &noWallet,
	}
	return wg.Run()
}

type WalletGUI struct {
	wg                        sync.WaitGroup
	cx                        *conte.Xt
	c                         *cli.Context
	quit                      qu.C
	State                     State
	noWallet                  *bool
	node, wallet, miner       *rununit.RunUnit
	walletToLock              time.Time
	walletLockTime            int
	ChainMutex, WalletMutex   sync.Mutex
	ChainClient, WalletClient *rpcclient.Client
	w                         map[string]*f.Window
	Size                      *int
	th                        *p9.Theme
	App                       *p9.App
	invalidate                qu.C
	unlockPage                *p9.App
	config                    *cfg.Config
	configs                   cfg.GroupsMap
	unlockPassword            *p9.Password
	sidebarButtons            []*p9.Clickable
	buttonBarButtons          []*p9.Clickable
	statusBarButtons          []*p9.Clickable
	quitClickable             *p9.Clickable
	bools                     map[string]*p9.Bool
	lists                     map[string]*p9.List
	checkables                map[string]*p9.Checkable
	clickables                map[string]*p9.Clickable
	inputs                    map[string]*p9.Input
	passwords                 map[string]*p9.Password
	incdecs                   map[string]*p9.IncDec
	historyTable              *p9.TextTable
	sendAddresses             []SendAddress
	console                   *Console
	// toasts                    *toast.Toasts
	// dialog                    *dialog.Dialog
}

func (wg *WalletGUI) Run() (err error) {
	wg.th = p9.NewTheme(p9fonts.Collection(), wg.quit)
	wg.th.Dark = wg.cx.Config.DarkTheme
	wg.th.Colors.SetTheme(*wg.th.Dark)
	*wg.noWallet = true
	wg.GetButtons()
	wg.State.AllTimeStrings.Store([]string{})
	wg.lists = wg.GetLists()
	wg.clickables = wg.GetClickables()
	wg.checkables = map[string]*p9.Checkable{
	}
	wg.GetHistoryTable()
	before := func() { Debug("running before") }
	after := func() { Debug("running after") }
	wg.node = wg.GetRunUnit(
		"NODE", before, after,
		os.Args[0], "-D", *wg.cx.Config.DataDir, "--servertls=true", "--clienttls=true", "--pipelog", "node",
	)
	wg.wallet = wg.GetRunUnit(
		"WLLT", before, after,
		os.Args[0], "-D", *wg.cx.Config.DataDir, "--servertls=true", "--clienttls=true", "--pipelog", "wallet",
	)
	wg.miner = wg.GetRunUnit(
		"MINE", before, after,
		os.Args[0], "-D", *wg.cx.Config.DataDir, "--pipelog", "kopach",
	)
	wg.bools = wg.GetBools()
	wg.GetInputs()
	wg.GetPasswords()
	// wg.toasts = toast.New(wg.th)
	// wg.dialog = dialog.New(wg.th)
	wg.console = wg.ConsolePage()
	wg.w = make(map[string]*f.Window)
	wg.quitClickable = wg.th.Clickable()
	wg.w = map[string]*f.Window{
		"splash": f.NewWindow(wg.th),
		"main":   f.NewWindow(wg.th),
	}
	wg.GetIncDecs()
	wg.App = wg.GetAppWidget()
	wg.unlockPage = wg.getWalletUnlockAppWidget()
	wg.Tickers()
	if !apputil.FileExists(*wg.cx.Config.WalletFile) {
	} else {
		*wg.noWallet = false
		if !*wg.cx.Config.NodeOff {
			// wg.startNode()
			wg.node.Start()
		}
		if *wg.cx.Config.Generate && *wg.cx.Config.GenThreads != 0 {
			// wg.startMiner()
			wg.miner.Start()
		}
		wg.unlockPassword.Focus()
	}
	wg.Size = wg.w["main"].Width
	go func() {
		if err := wg.w["main"].
			Size(64, 32).
			Title("ParallelCoin Wallet").
			Open().
			Run(
				func(gtx l.Context) l.Dimensions {
					return p9.If(
						*wg.noWallet,
						wg.CreateWalletPage,
						p9.If(
							!wg.wallet.Running(),
							wg.unlockPage.Fn(),
							wg.App.Fn(),
						),
					)(gtx)
				},
				wg.App.Overlay(),
				// wg.InitWallet(),
				wg.gracefulShutdown,
				// func() { interrupt.Request() },
				wg.quit,
			); Check(err) {
		}
	}()
	// interrupt.AddHandler(
	// 	func() {
	// 		Debug("quitting wallet gui")
	// 		// consume.Kill(wg.Node)
	// 		// consume.Kill(wg.Miner)
	// 		// wg.gracefulShutdown()
	// 		// wg.quit.Q()
	// 	},
	// )
out:
	for {
		select {
		case <-wg.invalidate:
			Trace("invalidating render queue")
			wg.w["main"].Window.Invalidate()
		case <-wg.cx.KillAll:
			break out
		case <-wg.quit:
			break out
		}
	}
	wg.gracefulShutdown()
	wg.quit.Q()
	return
}

func (wg *WalletGUI) GetButtons() {
	wg.sidebarButtons = make([]*p9.Clickable, 12)
	// wg.walletLocked.Store(true)
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
}

func (wg *WalletGUI) GetHistoryTable() {
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
}

func (wg *WalletGUI) GetInputs() {
	seed := make([]byte, hdkeychain.MaxSeedBytes)
	_, _ = rand.Read(seed)
	seedString := hex.EncodeToString(seed)
	wg.inputs = map[string]*p9.Input{
		"receiveLabel":   wg.th.Input("", "Label", "Primary", "DocText", "DocBg", func(pass string) {}),
		"receiveAmount":  wg.th.Input("", "Amount", "Primary", "DocText", "DocBg", func(pass string) {}),
		"receiveMessage": wg.th.Input("", "Message", "Primary", "DocText", "DocBg", func(pass string) {}),
		"console":        wg.th.Input("", "enter rpc command", "Primary", "DocText", "DocBg", func(pass string) {}),
		"walletSeed":     wg.th.Input(seedString, "wallet seed", "Primary", "DocText", "DocBg", func(pass string) {}),
	}
}

func (wg *WalletGUI) GetPasswords() {
	pass := ""
	passConfirm := ""
	wg.passwords = map[string]*p9.Password{
		"passEditor":        wg.th.Password("password", &pass, "Primary", "DocText", "", func(pass string) {}),
		"confirmPassEditor": wg.th.Password("confirm", &passConfirm, "Primary", "DocText", "", func(pass string) {}),
		"publicPassEditor":  wg.th.Password("public password (optional)", wg.cx.Config.WalletPass, "Primary", "DocText", "", func(pass string) {}),
	}
}

func (wg *WalletGUI) GetIncDecs() {
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
						if wg.miner.Running() && n != 0 {
							wg.miner.Stop()
							wg.miner.Start()
						}
						if n == 0 {
							wg.miner.Stop()
						}
						*wg.cx.Config.GenThreads = n
						save.Pod(wg.cx.Config)
						// if wg.miner.Running() {
						// 	Debug("restarting miner")
						// 	wg.miner.Stop()
						// 	wg.miner.Start()
						// }
					}()
				},
			),
		"transactionsPerPage": wg.th.IncDec().
			Min(10).
			Max(100).
			NDigits(3).
			Amount(10).
			SetCurrent(10).
			ChangeHook(
				func(n int) {
					Debug("showing", n, "per page")
				},
			),
		"idleTimeout": wg.th.IncDec().
			Scale(4).
			Min(60).
			Max(3600).
			NDigits(4).
			Amount(60).
			SetCurrent(300).
			ChangeHook(
				func(n int) {
					Debug("idle timeout", time.Duration(n)*time.Second)
				},
			),
	}
}

func (wg *WalletGUI) GetRunUnit(name string, before, after func(), args ...string) *rununit.RunUnit {
	return rununit.New(
		before,
		after,
		consume.SimpleLog(name),
		consume.FilterNone,
		wg.quit,
		args...,
	)
}

func (wg *WalletGUI) GetLists() (o map[string]*p9.List) {
	return map[string]*p9.List{
		"createWallet": wg.th.List(),
		"overview":     wg.th.List(),
		"balances":     wg.th.List(),
		"recent":       wg.th.List(),
		"send":         wg.th.List(),
		"transactions": wg.th.List(),
		"settings":     wg.th.List(),
		"received":     wg.th.List(),
		"history":      wg.th.List(),
	}
}

func (wg *WalletGUI) GetClickables() map[string]*p9.Clickable {
	return map[string]*p9.Clickable{
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
}

func (wg *WalletGUI) GetBools() map[string]*p9.Bool {
	return map[string]*p9.Bool{
		"runstate":     wg.th.Bool(wg.node.Running()),
		"encryption":   wg.th.Bool(false),
		"seed":         wg.th.Bool(false),
		"testnet":      wg.th.Bool(false),
		"ihaveread":    wg.th.Bool(false),
		"showGenerate": wg.th.Bool(true),
		"showSent":     wg.th.Bool(true),
		"showReceived": wg.th.Bool(true),
		"showImmature": wg.th.Bool(true),
	}
}

var shuttingDown = false

func (wg *WalletGUI) gracefulShutdown() {
	if shuttingDown {
		Debug(log.Caller("already called gracefulShutdown", 1))
		return
	} else {
		shuttingDown = true
	}
	Debug("\n\nquitting wallet gui")
	if wg.miner.Running() {
		Debug("stopping miner")
		wg.miner.Stop()
		wg.miner.Shutdown()
	}
	if wg.wallet.Running() {
		Debug("stopping wallet")
		wg.wallet.Stop()
		wg.wallet.Shutdown()
		wg.unlockPassword.Wipe()
		// wg.walletLocked.Store(true)
	}
	if wg.node.Running() {
		Debug("stopping node")
		wg.node.Stop()
		wg.node.Shutdown()
	}
	// wg.ChainMutex.Lock()
	if wg.ChainClient != nil {
		Debug("stopping chain client")
		wg.ChainClient.Shutdown()
		wg.ChainClient = nil
	}
	// wg.ChainMutex.Unlock()
	// wg.WalletMutex.Lock()
	if wg.WalletClient != nil {
		Debug("stopping wallet client")
		wg.WalletClient.Shutdown()
		wg.WalletClient = nil
	}
	// wg.WalletMutex.Unlock()
	interrupt.Request()
	time.Sleep(time.Second)
	wg.quit.Q()
}
