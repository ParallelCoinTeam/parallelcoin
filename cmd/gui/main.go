package gui

import (
	"crypto/rand"
	"github.com/tyler-smith/go-bip39"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
	
	"gioui.org/op/paint"
	uberatomic "go.uber.org/atomic"
	
	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/pkg/gui"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	"github.com/p9c/pod/pkg/util/interrupt"
	log "github.com/p9c/pod/pkg/util/logi"
	"github.com/p9c/pod/pkg/util/qu"
	
	"github.com/urfave/cli"
	
	l "gioui.org/layout"
	
	"github.com/p9c/pod/pkg/pipe/consume"
	"github.com/p9c/pod/pkg/util/rununit"
	
	"github.com/p9c/pod/app/apputil"
	
	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/pkg/gui/cfg"
	"github.com/p9c/pod/pkg/rpc/rpcclient"
)

func Main(cx *conte.Xt, c *cli.Context) (e error) {
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

type BoolMap map[string]*gui.Bool
type ListMap map[string]*gui.List
type CheckableMap map[string]*gui.Checkable
type ClickableMap map[string]*gui.Clickable
type Inputs map[string]*gui.Input
type Passwords map[string]*gui.Password
type IncDecMap map[string]*gui.IncDec

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
	WalletWatcher             qu.C
	*gui.Window
	Size                                     *int
	MainApp                                  *gui.App
	invalidate                               qu.C
	unlockPage                               *gui.App
	loadingPage                              *gui.App
	config                                   *cfg.Config
	configs                                  cfg.GroupsMap
	unlockPassword                           *gui.Password
	sidebarButtons                           []*gui.Clickable
	buttonBarButtons                         []*gui.Clickable
	statusBarButtons                         []*gui.Clickable
	receiveAddressbookClickables             []*gui.Clickable
	sendAddressbookClickables                []*gui.Clickable
	quitClickable                            *gui.Clickable
	bools                                    BoolMap
	lists                                    ListMap
	checkables                               CheckableMap
	clickables                               ClickableMap
	inputs                                   Inputs
	passwords                                Passwords
	incdecs                                  IncDecMap
	console                                  *Console
	RecentTxsWidget, TxHistoryWidget         l.Widget
	recentTxsClickables, txHistoryClickables []*gui.Clickable
	txRecentList, txHistoryList              []btcjson.ListTransactionsResult
	openTxID, prevOpenTxID                   *uberatomic.String
	originTxDetail                           string
	txMx                                     sync.Mutex
	Syncing                                  *uberatomic.Bool
	stateLoaded                              *uberatomic.Bool
	currentReceiveQRCode                     *paint.ImageOp
	currentReceiveAddress                    string
	currentReceiveQR                         l.Widget
	currentReceiveRegenClickable             *gui.Clickable
	currentReceiveCopyClickable              *gui.Clickable
	currentReceiveRegenerate                 *uberatomic.Bool
	// currentReceiveGetNew         *uberatomic.Bool
	sendClickable *gui.Clickable
	ready         *uberatomic.Bool
	mainDirection l.Direction
	preRendering  bool
	// ReceiveAddressbook l.Widget
	// SendAddressbook    l.Widget
	ReceivePage *ReceivePage
	SendPage    *SendPage
	// toasts                    *toast.Toasts
	// dialog                    *dialog.Dialog
	createSeed                          []byte
	createWords, showWords, createMatch string
	createVerifying                     bool
	restoring                           bool
}

func (wg *WalletGUI) Run() (e error) {
	wg.Syncing = uberatomic.NewBool(false)
	wg.openTxID = uberatomic.NewString("")
	wg.prevOpenTxID = uberatomic.NewString("")
	wg.stateLoaded = uberatomic.NewBool(false)
	wg.currentReceiveRegenerate = uberatomic.NewBool(true)
	// wg.currentReceiveGetNew = uberatomic.NewBool(false)
	wg.ready = uberatomic.NewBool(false)
	// wg.th = gui.NewTheme(p9fonts.Collection(), wg.quit)
	// wg.Window = gui.NewWindow(wg.th)
	wg.Window = gui.NewWindowP9(wg.quit)
	wg.Dark = wg.cx.Config.DarkTheme
	wg.Colors.SetTheme(*wg.Dark)
	*wg.noWallet = true
	wg.GetButtons()
	wg.lists = wg.GetLists()
	wg.clickables = wg.GetClickables()
	wg.checkables = wg.GetCheckables()
	before := func() { dbg.Ln("running before") }
	after := func() { dbg.Ln("running after") }
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
	wg.inputs = wg.GetInputs()
	wg.passwords = wg.GetPasswords()
	// wg.toasts = toast.New(wg.th)
	// wg.dialog = dialog.New(wg.th)
	wg.console = wg.ConsolePage()
	wg.quitClickable = wg.Clickable()
	wg.incdecs = wg.GetIncDecs()
	wg.Size = &wg.Window.Width
	wg.currentReceiveCopyClickable = wg.WidgetPool.GetClickable()
	wg.currentReceiveRegenClickable = wg.WidgetPool.GetClickable()
	wg.currentReceiveQR = func(gtx l.Context) l.Dimensions {
		return l.Dimensions{}
	}
	wg.ReceivePage = wg.GetReceivePage()
	wg.SendPage = wg.GetSendPage()
	wg.MainApp = wg.GetAppWidget()
	wg.State = GetNewState(wg.cx.ActiveNet, wg.MainApp.ActivePageGetAtomic())
	wg.unlockPage = wg.getWalletUnlockAppWidget()
	wg.loadingPage = wg.getLoadingPage()
	// wg.Watcher()
	if !apputil.FileExists(*wg.cx.Config.WalletFile) {
		inf.Ln("wallet file does not exist", *wg.cx.Config.WalletFile)
	} else {
		*wg.noWallet = false
		// if !*wg.cx.Config.NodeOff {
		// 	// wg.startNode()
		// 	wg.node.Start()
		// }
		if *wg.cx.Config.Generate && *wg.cx.Config.GenThreads != 0 {
			// wg.startMiner()
			wg.miner.Start()
		}
		wg.unlockPassword.Focus()
	}
	interrupt.AddHandler(
		func() {
			dbg.Ln("quitting wallet gui")
			// consume.Kill(wg.Node)
			// consume.Kill(wg.Miner)
			// wg.gracefulShutdown()
			wg.quit.Q()
		},
	)
	go func() {
	out:
		for {
			select {
			case <-wg.invalidate.Wait():
				trc.Ln("invalidating render queue")
				wg.Window.Window.Invalidate()
				// TODO: make a more appropriate trigger for this - ie, when state actually changes.
				// if wg.wallet.Running() && wg.stateLoaded.Load() {
				// 	filename := filepath.Join(wg.cx.DataDir, "state.json")
				// 	if e := wg.State.Save(filename, wg.cx.Config.WalletPass); err.Chk(e) {
				// 	}
				// }
			case <-wg.cx.KillAll.Wait():
				break out
			case <-wg.quit.Wait():
				break out
			}
		}
	}()
	if e := wg.Window.
		Size(56, 32).
		Title("ParallelCoin Wallet").
		Open().
		Run(
			func(gtx l.Context) l.Dimensions {
				return wg.Fill(
					"DocBg", l.Center, 0, 0, func(gtx l.Context) l.Dimensions {
						return gui.If(
							*wg.noWallet,
							wg.CreateWalletPage,
							func(gtx l.Context) l.Dimensions {
								switch {
								case wg.ready.Load() && wg.stateLoaded.Load():
									return wg.MainApp.Fn()(gtx)
								case wg.ready.Load() || wg.stateLoaded.Load():
									return wg.loadingPage.Fn()(gtx)
								default:
									return wg.unlockPage.Fn()(gtx)
								}
							},
							// gui.If(
							// 	wg.ready.Load(),
							// 	gui.If(
							// 		wg.WalletAndClientRunning(),
							// 		gui.If(
							// 			wg.stateLoaded.Load(),
							// 			wg.MainApp.Fn(),
							// 			wg.loadingPage.Fn(),
							// 		),
							// 		wg.loadingPage.Fn(),
							// 	),
							// 	gui.If(
							// 		wg.WalletAndClientRunning(),
							// 		wg.loadingPage.Fn(),
							// 		wg.unlockPage.Fn(),
							// 	),
							// ),
						)(gtx)
					},
				).Fn(gtx)
			},
			wg.MainApp.Overlay,
			interrupt.Request,
			wg.quit,
		); err.Chk(e) {
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

func (wg *WalletGUI) ShuffleSeed() {
	wg.createSeed = make([]byte, 32)
	_, _ = rand.Read(wg.createSeed)
	var e error
	var wk string
	if wk, e = bip39.NewMnemonic(wg.createSeed); err.Chk(e) {
		panic(e)
	}
	wg.createWords = wk
	// wg.createMatch = wk
	wks := strings.Split(wk, " ")
	var out string
	for i := 0; i < 24; i += 4 {
		out += strings.Join(wks[i:i+4], " ")
		if i+4 < 24 {
			out += "\n"
		}
	}
	wg.showWords = out
}

func (wg *WalletGUI) GetInputs() Inputs {
	wg.ShuffleSeed()
	return Inputs{
		"receiveAmount": wg.Input("", "Amount", "DocText", "PanelBg", "DocBg", func(amt string) {}, func(string) {}),
		"receiveMessage": wg.Input(
			"",
			"Description",
			"DocText",
			"PanelBg",
			"DocBg",
			func(pass string) {},
			func(string) {},
		),
		
		"sendAddress": wg.Input(
			"",
			"Parallelcoin Address",
			"DocText",
			"PanelBg",
			"DocBg",
			func(amt string) {},
			func(string) {},
		),
		"sendAmount": wg.Input("", "Amount", "DocText", "PanelBg", "DocBg", func(amt string) {}, func(string) {}),
		"sendMessage": wg.Input(
			"",
			"Description",
			"DocText",
			"PanelBg",
			"DocBg",
			func(pass string) {},
			func(string) {},
		),
		
		"console": wg.Input(
			"",
			"enter rpc command",
			"DocText",
			"Transparent",
			"PanelBg",
			func(pass string) {},
			func(string) {},
		),
		"walletWords": wg.Input(
			/*wg.createWords*/ "", "wallet word seed", "DocText", "DocBg", "PanelBg", func(string) {},
			func(seedWords string) {
				wg.createMatch = seedWords
				wg.Invalidate()
			},
		),
		"walletRestore": wg.Input(
			/*wg.createWords*/ "", "enter seed to restore", "DocText", "DocBg", "PanelBg", func(string) {},
			func(seedWords string) {
				var e error
				wg.createMatch = seedWords
				if wg.createSeed, e = bip39.EntropyFromMnemonic(seedWords); err.Chk(e) {
					return
				}
				wg.createWords = seedWords
				wg.Invalidate()
			},
		),
		// "walletSeed": wg.Input(
		// 	seedString, "wallet seed", "DocText", "DocBg", "PanelBg", func(seedHex string) {
		// 		var e error
		// 		if wg.createSeed, e = hex.DecodeString(seedHex); err.Chk(e) {
		// 			return
		// 		}
		// 		var wk string
		// 		if wk, e = bip39.NewMnemonic(wg.createSeed); err.Chk(e) {
		// 			panic(e)
		// 		}
		// 		wg.createWords=wk
		// 		wks := strings.Split(wk, " ")
		// 		var out string
		// 		for i := 0; i < 24; i += 4 {
		// 			out += strings.Join(wks[i:i+4], " ") + "\n"
		// 		}
		// 		wg.showWords = out
		// 	}, nil,
		// ),
	}
}

// GetPasswords returns the passwords used in the wallet GUI
func (wg *WalletGUI) GetPasswords() (passwords Passwords) {
	pass := ""
	passConfirm := ""
	passwords = Passwords{
		"passEditor": wg.Password(
			"password (minimum 8 characters length)",
			&pass,
			"DocText",
			"DocBg",
			"PanelBg",
			func(pass string) {},
		),
		"confirmPassEditor": wg.Password("confirm", &passConfirm, "DocText", "DocBg", "PanelBg", func(pass string) {}),
		"publicPassEditor": wg.Password(
			"public password (optional)",
			wg.cx.Config.WalletPass,
			"Primary",
			"DocText",
			"PanelBg",
			func(pass string) {},
		),
	}
	return
}

func (wg *WalletGUI) GetIncDecs() IncDecMap {
	return IncDecMap{
		"generatethreads": wg.IncDec().
			NDigits(2).
			Min(0).
			Max(runtime.NumCPU()).
			SetCurrent(*wg.cx.Config.GenThreads).
			ChangeHook(
				func(n int) {
					dbg.Ln("threads value now", n)
					go func() {
						dbg.Ln("setting thread count")
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
						// 	dbg.Ln("restarting miner")
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
					dbg.Ln("idle timeout", time.Duration(n)*time.Second)
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

func (wg *WalletGUI) GetLists() (o ListMap) {
	return ListMap{
		"createWallet":     wg.List(),
		"overview":         wg.List(),
		"balances":         wg.List(),
		"recent":           wg.List(),
		"send":             wg.List(),
		"sendMedium":       wg.List(),
		"sendAddresses":    wg.List(),
		"receive":          wg.List(),
		"receiveMedium":    wg.List(),
		"receiveAddresses": wg.List(),
		"transactions":     wg.List(),
		"settings":         wg.List(),
		"received":         wg.List(),
		"history":          wg.List(),
		"txdetail":         wg.List(),
	}
}

func (wg *WalletGUI) GetClickables() ClickableMap {
	return ClickableMap{
		"createWallet":            wg.Clickable(),
		"createVerify":            wg.Clickable(),
		"createShuffle":           wg.Clickable(),
		"createRestore":           wg.Clickable(),
		"genesis":                 wg.Clickable(),
		"autofill":                wg.Clickable(),
		"quit":                    wg.Clickable(),
		"sendSend":                wg.Clickable(),
		"sendSave":                wg.Clickable(),
		"sendFromRequest":         wg.Clickable(),
		"receiveCreateNewAddress": wg.Clickable(),
		"receiveClear":            wg.Clickable(),
		"receiveShow":             wg.Clickable(),
		"receiveRemove":           wg.Clickable(),
		"transactions10":          wg.Clickable(),
		"transactions30":          wg.Clickable(),
		"transactions50":          wg.Clickable(),
		"txPageForward":           wg.Clickable(),
		"txPageBack":              wg.Clickable(),
		"theme":                   wg.Clickable(),
	}
}

func (wg *WalletGUI) GetCheckables() CheckableMap {
	return CheckableMap{}
}

func (wg *WalletGUI) GetBools() BoolMap {
	return BoolMap{
		"runstate":     wg.Bool(wg.node.Running()),
		"encryption":   wg.Bool(false),
		"seed":         wg.Bool(false),
		"testnet":      wg.Bool(false),
		"lan":          wg.Bool(false),
		"solo":         wg.Bool(false),
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
		dbg.Ln(log.Caller("already called gracefulShutdown", 1))
		return
	} else {
		shuttingDown = true
	}
	dbg.Ln("\nquitting wallet gui\n")
	if wg.miner.Running() {
		dbg.Ln("stopping miner")
		wg.miner.Stop()
		wg.miner.Shutdown()
	}
	if wg.wallet.Running() {
		dbg.Ln("stopping wallet")
		wg.wallet.Stop()
		wg.wallet.Shutdown()
		wg.unlockPassword.Wipe()
		// wg.walletLocked.Store(true)
	}
	if wg.node.Running() {
		dbg.Ln("stopping node")
		wg.node.Stop()
		wg.node.Shutdown()
	}
	// wg.ChainMutex.Lock()
	if wg.ChainClient != nil {
		dbg.Ln("stopping chain client")
		wg.ChainClient.Shutdown()
		wg.ChainClient = nil
	}
	// wg.ChainMutex.Unlock()
	// wg.WalletMutex.Lock()
	if wg.WalletClient != nil {
		dbg.Ln("stopping wallet client")
		wg.WalletClient.Shutdown()
		wg.WalletClient = nil
	}
	// wg.WalletMutex.Unlock()
	// interrupt.Request()
	// time.Sleep(time.Second)
	wg.quit.Q()
}
