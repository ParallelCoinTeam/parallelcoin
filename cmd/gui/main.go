package gui

import (
	"crypto/rand"
	"fmt"
	"net"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/niubaoshu/gotiny"
	"github.com/tyler-smith/go-bip39"

	"github.com/p9c/log"
	"github.com/p9c/matrjoska/pkg/chainrpc/p2padvt"
	"github.com/p9c/matrjoska/pkg/pipe"
	"github.com/p9c/matrjoska/pkg/transport"
	"github.com/p9c/matrjoska/pod/state"
	"github.com/p9c/opts/meta"
	"github.com/p9c/opts/text"

	uberatomic "go.uber.org/atomic"

	"github.com/p9c/gio/op/paint"

	"github.com/p9c/qu"

	"github.com/p9c/interrupt"

	"github.com/p9c/gel"
	"github.com/p9c/matrjoska/pkg/btcjson"

	l "github.com/p9c/gio/layout"

	"github.com/p9c/matrjoska/cmd/gui/cfg"
	"github.com/p9c/matrjoska/pkg/apputil"
	"github.com/p9c/matrjoska/pkg/rpcclient"
	"github.com/p9c/matrjoska/pkg/util/rununit"
)

// Main is the entrypoint for the wallet GUI
func Main(cx *state.State) (e error) {
	size := uberatomic.NewInt32(0)
	noWallet := true
	wg := &WalletGUI{
		cx:         cx,
		invalidate: qu.Ts(16),
		quit:       cx.KillAll,
		Size:       size,
		noWallet:   &noWallet,
		otherNodes: make(map[uint64]*nodeSpec),
		certs:      cx.Config.ReadCAFile(),

	}
	return wg.Run()
}

type BoolMap map[string]*gel.Bool
type ListMap map[string]*gel.List
type CheckableMap map[string]*gel.Checkable
type ClickableMap map[string]*gel.Clickable
type Inputs map[string]*gel.Input
type Passwords map[string]*gel.Password
type IncDecMap map[string]*gel.IncDec

type WalletGUI struct {
	wg                        sync.WaitGroup
	cx                        *state.State
	quit                      qu.C
	State                     *State
	noWallet                  *bool
	node, wallet, miner       *rununit.RunUnit
	walletToLock              time.Time
	walletLockTime            int
	ChainMutex, WalletMutex   sync.Mutex
	ChainClient, WalletClient *rpcclient.Client
	WalletWatcher             qu.C
	*gel.Window
	Size                                     *uberatomic.Int32
	MainApp                                  *gel.App
	invalidate                               qu.C
	unlockPage                               *gel.App
	loadingPage                              *gel.App
	config                                   *cfg.Config
	configs                                  cfg.GroupsMap
	unlockPassword                           *gel.Password
	sidebarButtons                           []*gel.Clickable
	buttonBarButtons                         []*gel.Clickable
	statusBarButtons                         []*gel.Clickable
	receiveAddressbookClickables             []*gel.Clickable
	sendAddressbookClickables                []*gel.Clickable
	quitClickable                            *gel.Clickable
	bools                                    BoolMap
	lists                                    ListMap
	checkables                               CheckableMap
	clickables                               ClickableMap
	inputs                                   Inputs
	passwords                                Passwords
	incdecs                                  IncDecMap
	console                                  *Console
	RecentTxsWidget, TxHistoryWidget         l.Widget
	recentTxsClickables, txHistoryClickables []*gel.Clickable
	txHistoryList                            []btcjson.ListTransactionsResult
	openTxID, prevOpenTxID                   *uberatomic.String
	originTxDetail                           string
	txMx                                     sync.Mutex
	stateLoaded                              *uberatomic.Bool
	currentReceiveQRCode                     *paint.ImageOp
	currentReceiveAddress                    string
	currentReceiveQR                         l.Widget
	currentReceiveRegenClickable             *gel.Clickable
	currentReceiveCopyClickable              *gel.Clickable
	currentReceiveRegenerate                 *uberatomic.Bool
	// currentReceiveGetNew         *uberatomic.Bool
	sendClickable *gel.Clickable
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
	lastUpdated                         uberatomic.Int64
	multiConn                           *transport.Channel
	otherNodes                          map[uint64]*nodeSpec
	uuid                                uint64
	peerCount                           *uberatomic.Int32
	certs                               []byte
}

type nodeSpec struct {
	time.Time
	addr string
}

func (wg *WalletGUI) Run() (e error) {
	wg.openTxID = uberatomic.NewString("")
	var mc *transport.Channel
	quit := qu.T()
	// I.Ln(wg.cx.Config.MulticastPass.V(), string(wg.cx.Config.MulticastPass.
	// 	Bytes()))
	if mc, e = transport.NewBroadcastChannel(
		"controller",
		wg,
		wg.cx.Config.MulticastPass.Bytes(),
		transport.DefaultPort,
		16384,
		handlersMulticast,
		quit,
	); E.Chk(e) {
		return
	}
	wg.multiConn = mc
	wg.peerCount = uberatomic.NewInt32(0)
	wg.prevOpenTxID = uberatomic.NewString("")
	wg.stateLoaded = uberatomic.NewBool(false)
	wg.currentReceiveRegenerate = uberatomic.NewBool(true)
	wg.ready = uberatomic.NewBool(false)
	wg.Window = gel.NewWindowP9(wg.quit)
	wg.Dark = wg.cx.Config.DarkTheme
	wg.Colors.SetDarkTheme(wg.Dark.True())
	*wg.noWallet = true
	wg.GetButtons()
	wg.lists = wg.GetLists()
	wg.clickables = wg.GetClickables()
	wg.checkables = wg.GetCheckables()
	before := func() { D.Ln("running before") }
	after := func() { D.Ln("running after") }
	I.Ln(os.Args[1:])
	options := []string{os.Args[0]}
	// options = append(options, wg.cx.Config.FoundArgs...)
	// options = append(options, "pipelog")
	wg.node = wg.GetRunUnit(
		"NODE", before, after,
		append(options, "node")...,
		// "node",
	)
	wg.wallet = wg.GetRunUnit(
		"WLLT", before, after,
		append(options, "wallet")...,
		// "wallet",
	)
	wg.miner = wg.GetRunUnit(
		"MINE", before, after,
		append(options, "kopach")...,
		// "wallet",
	)
	// I.S(wg.node, wg.wallet, wg.miner)
	wg.bools = wg.GetBools()
	wg.inputs = wg.GetInputs()
	wg.passwords = wg.GetPasswords()
	// wg.toasts = toast.New(wg.th)
	// wg.dialog = dialog.New(wg.th)
	wg.console = wg.ConsolePage()
	wg.quitClickable = wg.Clickable()
	wg.incdecs = wg.GetIncDecs()
	wg.Size = wg.Window.Width
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
	if !apputil.FileExists(wg.cx.Config.WalletFile.V()) {
		I.Ln("wallet file does not exist", wg.cx.Config.WalletFile.V())
	} else {
		*wg.noWallet = false
		// if !*wg.cx.Config.NodeOff {
		// 	// wg.startNode()
		// 	wg.node.Start()
		// }
		if wg.cx.Config.Generate.True() && wg.cx.Config.GenThreads.V() != 0 {
			// wg.startMiner()
			wg.miner.Start()
		}
		wg.unlockPassword.Focus()
	}
	interrupt.AddHandler(
		func() {
			D.Ln("quitting wallet gui")
			// consume.Kill(wg.Node)
			// consume.Kill(wg.Miner)
			// wg.gracefulShutdown()
			wg.quit.Q()
		},
	)
	go func() {
		ticker := time.NewTicker(time.Second)
	out:
		for {
			select {
			case <-ticker.C:
				go func() {
					if e = wg.Advertise(); E.Chk(e) {
					}
					if wg.node.Running() {
						if wg.ChainClient != nil {
							if !wg.ChainClient.Disconnected() {
								var pi []btcjson.GetPeerInfoResult
								if pi, e = wg.ChainClient.GetPeerInfo(); E.Chk(e) {
									return
								}
								wg.peerCount.Store(int32(len(pi)))
								wg.Invalidate()
							}
						}
					}
				}()
			case <-wg.invalidate.Wait():
				T.Ln("invalidating render queue")
				wg.Window.Window.Invalidate()
				// TODO: make a more appropriate trigger for this - ie, when state actually changes.
				// if wg.wallet.Running() && wg.stateLoaded.Load() {
				// 	filename := filepath.Join(wg.cx.DataDir, "state.json")
				// 	if e := wg.State.Save(filename, wg.cx.Config.WalletPass); E.Chk(e) {
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
						return gel.If(
							*wg.noWallet,
							wg.CreateWalletPage,
							func(gtx l.Context) l.Dimensions {
								switch {
								case wg.stateLoaded.Load():
									return wg.MainApp.Fn()(gtx)
								default:
									return wg.unlockPage.Fn()(gtx)
								}
							},
						)(gtx)
					},
				).Fn(gtx)
			},
			wg.quit.Q,
			wg.quit,
		); E.Chk(e) {
	}
	wg.gracefulShutdown()
	wg.quit.Q()
	return
}

func (wg *WalletGUI) GetButtons() {
	wg.sidebarButtons = make([]*gel.Clickable, 12)
	// wg.walletLocked.Store(true)
	for i := range wg.sidebarButtons {
		wg.sidebarButtons[i] = wg.Clickable()
	}
	wg.buttonBarButtons = make([]*gel.Clickable, 5)
	for i := range wg.buttonBarButtons {
		wg.buttonBarButtons[i] = wg.Clickable()
	}
	wg.statusBarButtons = make([]*gel.Clickable, 8)
	for i := range wg.statusBarButtons {
		wg.statusBarButtons[i] = wg.Clickable()
	}
}

func (wg *WalletGUI) ShuffleSeed() {
	wg.createSeed = make([]byte, 32)
	_, _ = rand.Read(wg.createSeed)
	var e error
	var wk string
	if wk, e = bip39.NewMnemonic(wg.createSeed); E.Chk(e) {
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
			"Title",
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
			"Title",
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
				if wg.createSeed, e = bip39.EntropyFromMnemonic(seedWords); E.Chk(e) {
					return
				}
				wg.createWords = seedWords
				wg.Invalidate()
			},
		),
		// "walletSeed": wg.Input(
		// 	seedString, "wallet seed", "DocText", "DocBg", "PanelBg", func(seedHex string) {
		// 		var e error
		// 		if wg.createSeed, e = hex.DecodeString(seedHex); E.Chk(e) {
		// 			return
		// 		}
		// 		var wk string
		// 		if wk, e = bip39.NewMnemonic(wg.createSeed); E.Chk(e) {
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
	passwords = Passwords{
		"passEditor": wg.Password(
			"password (minimum 8 characters length)",
			text.New(meta.Data{}, ""),
			"DocText",
			"DocBg",
			"PanelBg",
			func(pass string) {},
		),
		"confirmPassEditor": wg.Password(
			"confirm",
			text.New(meta.Data{}, ""),
			"DocText",
			"DocBg",
			"PanelBg",
			func(pass string) {},
		),
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
			SetCurrent(wg.cx.Config.GenThreads.V()).
			ChangeHook(
				func(n int) {
					D.Ln("threads value now", n)
					go func() {
						D.Ln("setting thread count")
						if wg.miner.Running() && n != 0 {
							wg.miner.Stop()
							wg.miner.Start()
						}
						if n == 0 {
							wg.miner.Stop()
						}
						wg.cx.Config.GenThreads.Set(n)
						_ = wg.cx.Config.WriteToFile(wg.cx.Config.ConfigFile.V())
						// if wg.miner.Running() {
						// 	D.Ln("restarting miner")
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
					D.Ln("idle timeout", time.Duration(n)*time.Second)
				},
			),
	}
}

func (wg *WalletGUI) GetRunUnit(
	name string, before, after func(), args ...string,
) *rununit.RunUnit {
	I.Ln("getting rununit for", name, args)
	// we have to copy the args otherwise further mutations affect this one
	argsCopy := make([]string, len(args))
	copy(argsCopy, args)
	return rununit.New(name, before, after, pipe.SimpleLog(name),
		pipe.FilterNone, wg.quit, argsCopy...)
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
		"balanceConfirmed":        wg.Clickable(),
		"balanceUnconfirmed":      wg.Clickable(),
		"balanceTotal":            wg.Clickable(),
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
		D.Ln(log.Caller("already called gracefulShutdown", 1))
		return
	} else {
		shuttingDown = true
	}
	D.Ln("\nquitting wallet gui\n")
	if wg.miner.Running() {
		D.Ln("stopping miner")
		wg.miner.Stop()
		wg.miner.Shutdown()
	}
	if wg.wallet.Running() {
		D.Ln("stopping wallet")
		wg.wallet.Stop()
		wg.wallet.Shutdown()
		wg.unlockPassword.Wipe()
		// wg.walletLocked.Store(true)
	}
	if wg.node.Running() {
		D.Ln("stopping node")
		wg.node.Stop()
		wg.node.Shutdown()
	}
	// wg.ChainMutex.Lock()
	if wg.ChainClient != nil {
		D.Ln("stopping chain client")
		wg.ChainClient.Shutdown()
		wg.ChainClient = nil
	}
	// wg.ChainMutex.Unlock()
	// wg.WalletMutex.Lock()
	if wg.WalletClient != nil {
		D.Ln("stopping wallet client")
		wg.WalletClient.Shutdown()
		wg.WalletClient = nil
	}
	// wg.WalletMutex.Unlock()
	// interrupt.Request()
	// time.Sleep(time.Second)
	wg.quit.Q()
}

var handlersMulticast = transport.Handlers{
	// string(sol.Magic):      processSolMsg,
	string(p2padvt.Magic): processAdvtMsg,
	// string(hashrate.Magic): processHashrateMsg,
}

func processAdvtMsg(
	ctx interface{}, src net.Addr, dst string, b []byte,
) (e error) {
	wg := ctx.(*WalletGUI)
	if wg.cx.Config.Discovery.False() {
		return
	}
	if wg.ChainClient == nil {
		T.Ln("no chain client to process advertisment")
		return
	}
	var j p2padvt.Advertisment
	gotiny.Unmarshal(b, &j)
	// I.S(j)
	var peerUUID uint64
	peerUUID = j.UUID
	// I.Ln("peerUUID of advertisment", peerUUID, wg.otherNodes)
	if int(peerUUID) == wg.cx.Config.UUID.V() {
		D.Ln("ignoring own advertisment message")
		return
	}
	if _, ok := wg.otherNodes[peerUUID]; !ok {
		var pi []btcjson.GetPeerInfoResult
		if pi, e = wg.ChainClient.GetPeerInfo(); E.Chk(e) {
		}
		for i := range pi {
			for k := range j.IPs {
				jpa := net.JoinHostPort(k, fmt.Sprint(j.P2P))
				if jpa == pi[i].Addr {
					I.Ln("not connecting to node already connected outbound")
					return
				}
				if jpa == pi[i].AddrLocal {
					I.Ln("not connecting to node already connected inbound")
					return
				}
			}
		}
		// if we haven't already added it to the permanent peer list, we can add it now
		I.Ln("connecting to lan peer with same PSK", j.IPs, peerUUID)
		wg.otherNodes[peerUUID] = &nodeSpec{}
		wg.otherNodes[peerUUID].Time = time.Now()
		for i := range j.IPs {
			addy := net.JoinHostPort(i, fmt.Sprint(j.P2P))
			for j := range pi {
				if addy == pi[j].Addr || addy == pi[j].AddrLocal {
					// not connecting to peer we already have connected to
					return
				}
			}
		}
		// try all IPs
		for addr := range j.IPs {
			peerIP := net.JoinHostPort(addr, fmt.Sprint(j.P2P))
			if e = wg.ChainClient.AddNode(peerIP, "add"); E.Chk(e) {
				continue
			}
			D.Ln("connected to peer via address", peerIP)
			wg.otherNodes[peerUUID].addr = peerIP
			break
		}
		I.Ln(peerUUID, "added", "otherNodes", wg.otherNodes)
	} else {
		// update last seen time for peerUUID for garbage collection of stale disconnected
		// nodes
		D.Ln("other node seen again", peerUUID, wg.otherNodes[peerUUID].addr)
		wg.otherNodes[peerUUID].Time = time.Now()
	}
	// I.S(wg.otherNodes)
	// If we lose connection for more than 9 seconds we delete and if the node
	// reappears it can be reconnected
	for i := range wg.otherNodes {
		D.Ln(i, wg.otherNodes[i])
		tn := time.Now()
		if tn.Sub(wg.otherNodes[i].Time) > time.Second*6 {
			// also remove from connection manager
			if e = wg.ChainClient.AddNode(wg.otherNodes[i].addr, "remove"); E.Chk(e) {
			}
			D.Ln("deleting", tn, wg.otherNodes[i], i)
			delete(wg.otherNodes, i)
		}
	}
	// on := int32(len(wg.otherNodes))
	// wg.otherNodeCount.Store(on)
	return
}
