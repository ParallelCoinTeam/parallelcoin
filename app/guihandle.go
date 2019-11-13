package app

import (
	"encoding/json"
	"fmt"
	"github.com/p9c/pod/cmd/gui"
	"github.com/p9c/pod/cmd/node"
	"github.com/p9c/pod/cmd/node/rpc"
	"github.com/p9c/pod/cmd/walletmain"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gui/utils"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/util/interrupt"
	"github.com/p9c/pod/pkg/wallet"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/webengine"
	"github.com/therecipe/qt/widgets"
	"github.com/urfave/cli"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
)

type Bios struct {
	Theme      bool   `json:"theme"`
	IsBoot     bool   `json:"boot"`
	IsBootMenu bool   `json:"menu"`
	IsBootLogo bool   `json:"logo"`
	IsLoading  bool   `json:"loading"`
	IsDev      bool   `json:"dev"`
	IsScreen   string `json:"screen"`
}


var guiHandle = func(cx *conte.Xt) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		widgets.NewQApplication(len(os.Args), os.Args)

		b := Bios{
			Theme:      false,
			IsBoot:     true,
			IsBootMenu: true,
			IsBootLogo: true,
			IsLoading:  false,
			IsDev:      true,
			IsScreen:   "overview",
		}
		http.HandleFunc("/bios", func(w http.ResponseWriter, r *http.Request) {
			js, err := json.Marshal(&b)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Write(js)
		})
		b.IsBootLogo = true
		log.INFO("starting GUI")
		var view = webengine.NewQWebEngineView(nil)
		//view.SetUrl(QUrl("qrc:/index.html"))
		view.Load(core.NewQUrl3("qrc:/index.html", 0))
		view.Show()
		utils.GetBiosMessage(view, "starting GUI")

		Configure(cx)
		shutdownChan := make(chan struct{})
		walletChan := make(chan *wallet.Wallet)
		nodeChan := make(chan *rpc.Server)
		cx.WalletKill = make(chan struct{})
		cx.NodeKill = make(chan struct{})
		cx.Wallet = &atomic.Value{}
		cx.Wallet.Store(false)
		cx.Node = &atomic.Value{}
		cx.Node.Store(false)
		var err error
		var wg sync.WaitGroup
		if !*cx.Config.NodeOff {
			go func() {
				log.INFO(cx.Language.RenderText("goApp_STARTINGNODE"))
				utils.GetBiosMessage(view, cx.Language.RenderText("goApp_STARTINGNODE"))

				err = node.Main(cx, shutdownChan, cx.NodeKill, nodeChan, &wg)
				if err != nil {
					fmt.Println("error running node:", err)
					os.Exit(1)
				}
			}()
			log.DEBUG("waiting for nodeChan")
			cx.RPCServer = <-nodeChan
			log.DEBUG("nodeChan sent")
			cx.Node.Store(true)
		}
		if !*cx.Config.WalletOff {
			go func() {
				log.INFO("starting wallet")
				utils.GetBiosMessage(view, "starting wallet")
				err = walletmain.Main(cx.Config, cx.StateCfg,
					cx.ActiveNet, walletChan, cx.WalletKill, &wg)
				if err != nil {
					fmt.Println("error running wallet:", err)
					os.Exit(1)
				}
			}()
			log.DEBUG("waiting for walletChan")
			cx.WalletServer = <-walletChan
			log.DEBUG("walletChan sent")
			cx.Wallet.Store(true)
		}
		interrupt.AddHandler(func() {
			log.WARN("interrupt received, " +
				"shutting down shell modules")
			close(cx.WalletKill)
			close(cx.NodeKill)
		})
		gui.GUI(cx)
		b.IsBootLogo = false
		b.IsBoot = false
		widgets.QApplication_Exec()

		if !cx.Node.Load().(bool) {
			close(cx.WalletKill)
		}
		if !cx.Wallet.Load().(bool) {
			close(cx.NodeKill)
		}
		return err
	}
}
