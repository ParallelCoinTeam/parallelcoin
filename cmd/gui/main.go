package gui

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
	
	uberatomic "go.uber.org/atomic"
	
	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/pkg/gui"
	"github.com/p9c/pod/pkg/rpc/btcjson"
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
	State                     *State
	noWallet                  *bool
	node, wallet, miner       *rununit.RunUnit
	walletToLock              time.Time
	walletLockTime            int
	ChainMutex, WalletMutex   sync.Mutex
	ChainClient, WalletClient *rpcclient.Client
	*gui.Window
	Size                        *int
	MainApp                     *gui.App
	invalidate                  qu.C
	unlockPage                  *gui.App
	config                      *cfg.Config
	configs                     cfg.GroupsMap
	unlockPassword              *gui.Password
	sidebarButtons              []*gui.Clickable
	buttonBarButtons            []*gui.Clickable
	statusBarButtons            []*gui.Clickable
	quitClickable               *gui.Clickable
	bools                       map[string]*gui.Bool
	lists                       map[string]*gui.List
	checkables                  map[string]*gui.Checkable
	clickables                  map[string]*gui.Clickable
	inputs                      map[string]*gui.Input
	passwords                   map[string]*gui.Password
	incdecs                     map[string]*gui.IncDec
	sendAddresses               []SendAddress
	console                     *Console
	RecentTransactionsWidget    l.Widget
	HistoryWidget               l.Widget
	txRecentList, txHistoryList []btcjson.ListTransactionsResult
	txMx                        sync.Mutex
	Syncing                     *uberatomic.Bool
	stateLoaded                 *uberatomic.Bool
	// toasts                    *toast.Toasts
	// dialog                    *dialog.Dialog
}

func (wg *WalletGUI) Run() (err error) {
	wg.Syncing = uberatomic.NewBool(false)
	wg.stateLoaded = uberatomic.NewBool(false)
	// wg.th = gui.NewTheme(p9fonts.Collection(), wg.quit)
	// wg.Window = gui.NewWindow(wg.th)
	wg.Window = gui.NewWindowP9(wg.quit)
	wg.Dark = wg.cx.Config.DarkTheme
	wg.Colors.SetTheme(*wg.Dark)
	*wg.noWallet = true
	wg.GetButtons()
	wg.lists = wg.GetLists()
	wg.clickables = wg.GetClickables()
	wg.checkables = map[string]*gui.Checkable{
	}
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
	wg.quitClickable = wg.Clickable()
	wg.GetIncDecs()
	wg.MainApp = wg.GetAppWidget()
	wg.State = GetNewState(wg.cx.ActiveNet, wg.MainApp.ActivePageGetAtomic())
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
	wg.Size = &wg.Window.Width
	interrupt.AddHandler(
		func() {
			Debug("quitting wallet gui")
			// consume.Kill(wg.Node)
			// consume.Kill(wg.Miner)
			wg.gracefulShutdown()
			wg.quit.Q()
		},
	)
	go func() {
	out:
		for {
			select {
			case <-wg.invalidate:
				Trace("invalidating render queue")
				wg.Window.Window.Invalidate()
				if wg.wallet.Running() && wg.stateLoaded.Load() {
					filename := filepath.Join(wg.cx.DataDir, "state.json")
					if err := wg.State.Save(filename, wg.cx.Config.WalletPass); Check(err) {
					}
				}
			case <-wg.cx.KillAll:
				break out
			case <-wg.quit:
				break out
			}
		}
	}()
	if err := wg.Window.
		Size(64, 32).
		Title("ParallelCoin Wallet").
		Open().
		Run(
			wg.Fill("DocBg",
				func(gtx l.Context) l.Dimensions {
					return gui.If(
						*wg.noWallet,
						wg.CreateWalletPage,
						gui.If(
							!wg.wallet.Running(),
							wg.unlockPage.Fn(),
							wg.MainApp.Fn(),
						),
					)(gtx)
				},
				l.W, 0).Fn,
			wg.MainApp.Overlay,
			// wg.InitWallet(),
			wg.gracefulShutdown,
			// func() { interrupt.Request() },
			wg.quit,
		); Check(err) {
	}
	wg.gracefulShutdown()
	wg.quit.Q()
	return
}

func (wg *WalletGUI) GetButtons() {
	wg.sidebarButtons = make([]*gui.Clickable, 12)
	// wg.walletLocked.Store(true)
	for i := range wg.sidebarButtons {
		wg.sidebarButtons[i] = wg.Clickable()
	}
	wg.buttonBarButtons = make([]*gui.Clickable, 5)
	for i := range wg.buttonBarButtons {
		wg.buttonBarButtons[i] = wg.Clickable()
	}
	wg.statusBarButtons = make([]*gui.Clickable, 6)
	for i := range wg.statusBarButtons {
		wg.statusBarButtons[i] = wg.Clickable()
	}
}

func (wg *WalletGUI) GetInputs() {
	seed := make([]byte, hdkeychain.MaxSeedBytes)
	_, _ = rand.Read(seed)
	seedString := hex.EncodeToString(seed)
	wg.inputs = map[string]*gui.Input{
		"receiveLabel":   wg.Input("", "Label", "Primary", "DocText", "DocBg", func(pass string) {}),
		"receiveAmount":  wg.Input("", "Amount", "Primary", "DocText", "DocBg", func(pass string) {}),
		"receiveMessage": wg.Input("", "Message", "Primary", "DocText", "DocBg", func(pass string) {}),
		"console":        wg.Input("", "enter rpc command", "Primary", "DocText", "DocBg", func(pass string) {}),
		"walletSeed":     wg.Input(seedString, "wallet seed", "Primary", "DocText", "DocBg", func(pass string) {}),
	}
}

func (wg *WalletGUI) GetPasswords() {
	pass := ""
	passConfirm := ""
	wg.passwords = map[string]*gui.Password{
		"passEditor":        wg.Password("password", &pass, "Primary", "DocText", "", func(pass string) {}),
		"confirmPassEditor": wg.Password("confirm", &passConfirm, "Primary", "DocText", "", func(pass string) {}),
		"publicPassEditor":  wg.Password("public password (optional)", wg.cx.Config.WalletPass, "Primary", "DocText", "", func(pass string) {}),
	}
}

func (wg *WalletGUI) GetIncDecs() {
	wg.incdecs = map[string]*gui.IncDec{
		"generatethreads": wg.IncDec().
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
		"idleTimeout": wg.IncDec().
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

func (wg *WalletGUI) GetLists() (o map[string]*gui.List) {
	return map[string]*gui.List{
		"createWallet": wg.List(),
		"overview":     wg.List(),
		"balances":     wg.List(),
		"recent":       wg.List(),
		"send":         wg.List(),
		"transactions": wg.List(),
		"settings":     wg.List(),
		"received":     wg.List(),
		"history":      wg.List(),
	}
}

func (wg *WalletGUI) GetClickables() map[string]*gui.Clickable {
	return map[string]*gui.Clickable{
		"createWallet":            wg.Clickable(),
		"quit":                    wg.Clickable(),
		"sendSend":                wg.Clickable(),
		"sendClearAll":            wg.Clickable(),
		"sendAddRecipient":        wg.Clickable(),
		"receiveCreateNewAddress": wg.Clickable(),
		"receiveClear":            wg.Clickable(),
		"receiveShow":             wg.Clickable(),
		"receiveRemove":           wg.Clickable(),
		"transactions10":          wg.Clickable(),
		"transactions30":          wg.Clickable(),
		"transactions50":          wg.Clickable(),
		"txPageForward":           wg.Clickable(),
		"txPageBack":              wg.Clickable(),
	}
}

func (wg *WalletGUI) GetBools() map[string]*gui.Bool {
	return map[string]*gui.Bool{
		"runstate":     wg.Bool(wg.node.Running()),
		"encryption":   wg.Bool(false),
		"seed":         wg.Bool(false),
		"testnet":      wg.Bool(false),
		"ihaveread":    wg.Bool(false),
		"showGenerate": wg.Bool(true),
		"showSent":     wg.Bool(true),
		"showReceived": wg.Bool(true),
		"showImmature": wg.Bool(true),
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
