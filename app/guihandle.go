package app

import (
	"github.com/p9c/pod/cmd/gui"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/log"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/webengine"
	"github.com/therecipe/qt/widgets"
	"github.com/urfave/cli"
	"os"
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

		//qtgui.NewQGuiApplication(len(os.Args), os.Args)
		//var view = qml.NewQQmlApplicationEngine(nil)

		b := Bios{
			Theme:      false,
			IsBoot:     true,
			IsBootMenu: true,
			IsBootLogo: true,
			IsLoading:  false,
			IsDev:      true,
			IsScreen:   "overview",
		}
		log.INFO("starting GUI")

		var view = webengine.NewQWebEngineView(nil)
		//view.SetUrl(QUrl("qrc:/index.html"))
		view.Load(core.NewQUrl3("qrc:/index.html", 0))
		//view.Show()
		//utils.GetBiosMessage(view, "starting GUI")
		//view := quick.NewQQuickView(nil)
		//view.SetTitle("ctxproperty Example")
		//
		//view.SetResizeMode(quick.QQuickView__SizeRootObjectToView)
		//view.RootContext().SetContextProperty("ctxObject", NewCtxObject(nil))
		//view.Load(core.NewQUrl3("qrc:/main.qml", 0))
		//view.SetSource(core.NewQUrl3("qrc:/main.qml", 0))
		view.Show()

		Configure(cx, c)
		err := gui.Services(cx)
		gui.GUI(cx, *view)

		b.IsBootLogo = false
		b.IsBoot = false

		if !cx.Node.Load().(bool) {
			close(cx.WalletKill)
		}
		if !cx.Wallet.Load().(bool) {
			close(cx.NodeKill)
		}
		//qtgui.QGuiApplication_Exec()
		widgets.QApplication_Exec()
		return err
	}
}
