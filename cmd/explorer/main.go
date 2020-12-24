package explorer

import (
	"gioui.org/app"
	l "gioui.org/layout"
	qu "github.com/p9c/pod/pkg/util/quit"
	"github.com/urfave/cli"

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
	ex := &Explorer{
		cx:         cx,
		c:          c,
		invalidate: qu.T(),
		quit:       cx.KillAll,
		size:       &size,
	}
	return ex.Run()
}

type Explorer struct {
	cx   *conte.Xt
	c    *cli.Context
	w    *f.Window
	th   *p9.Theme
	size *int
	*p9.App
	buttonBarButtons          []*p9.Clickable
	statusBarButtons          []*p9.Clickable
	bools                     map[string]*p9.Bool
	quitClickable             *p9.Clickable
	lists                     map[string]*p9.List
	checkables                map[string]*p9.Checkable
	clickables                map[string]*p9.Clickable
	inputs                    map[string]*p9.Input
	configs                   cfg.GroupsMap
	config                    *cfg.Config
	running                   bool
	invalidate                qu.C
	quit                      qu.C
	Worker                    *worker.Worker
	RunCommandChan            chan string
	State                     State
	Shell                     *worker.Worker
	blocks                    []btcjson.BlockDetails
	ChainClient, WalletClient *rpcclient.Client
}

func (ex *Explorer) Run() (err error) {
	ex.th = p9.NewTheme(p9fonts.Collection(), ex.quit)
	ex.th.Dark = ex.cx.Config.DarkTheme
	ex.th.Colors.SetTheme(*ex.th.Dark)
	ex.buttonBarButtons = make([]*p9.Clickable, 4)
	for i := range ex.buttonBarButtons {
		ex.buttonBarButtons[i] = ex.th.Clickable()
	}
	ex.statusBarButtons = make([]*p9.Clickable, 3)
	for i := range ex.statusBarButtons {
		ex.statusBarButtons[i] = ex.th.Clickable()
	}
	ex.lists = map[string]*p9.List{
		"blocks": ex.th.List(),
	}
	ex.clickables = map[string]*p9.Clickable{
		"quit": ex.th.Clickable(),
	}
	ex.bools = map[string]*p9.Bool{
		"runstate":   ex.th.Bool(ex.running),
		"encryption": ex.th.Bool(false),
		"seed":       ex.th.Bool(false),
		"testnet":    ex.th.Bool(false),
	}

	ex.inputs = map[string]*p9.Input{
		"receiveLabel":   ex.th.Input("", "Label", "Primary", "DocText", "DocBg", func(pass string) {}),
		"receiveAmount":  ex.th.Input("", "Amount", "Primary", "DocText", "DocBg", func(pass string) {}),
		"receiveMessage": ex.th.Input("", "Message", "Primary", "DocText", "DocBg", func(pass string) {}),
	}

	ex.RunCommandChan = make(chan string)
	if err = ex.Runner(); Check(err) {
	}
	ex.RunCommandChan <- "run"
	ex.quitClickable = ex.th.Clickable()
	ex.w = f.NewWindow(ex.th)

	ex.App = ex.GetAppWidget()
	go func() {
		if err := ex.w.
			Size(64, 32).
			Title("ParallelCoin Wallet").
			Open().
			Run(
				ex.Fn(),
				func(gtx l.Context) {},
				func() {
					Debug("quitting wallet gui")
					interrupt.Request()
				}, ex.quit); Check(err) {
		}
	}()
	// tickers and triggers
	go func() {
	out:
		for {
			select {
			case <-ex.invalidate:
				ex.w.Window.Invalidate()
			case <-ex.quit:
				break out
			}
		}
	}()
	app.Main()
	return
}
