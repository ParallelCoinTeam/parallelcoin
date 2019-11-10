package gui

import (
	"github.com/p9c/pod/pkg/conte"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/webengine"
	"github.com/therecipe/qt/widgets"
	"os"
)

func GUI(cx *conte.Xt) {
	//r := rcvar{
	//	cx:     cx,
	//	alert:  DuOSalert{},
	//	status: DuOStatus{},
	//}

	widgets.NewQApplication(len(os.Args), os.Args)
	var view = webengine.NewQWebEngineView(nil)
	//view.SetUrl(QUrl("qrc:/index.html"))
	view.Load(core.NewQUrl3("qrc:/index.html", 0))

	//view.setUrl(QUrl((QStringLiteral("http://localhost:8080"))));
	view.Show()
	widgets.QApplication_Exec()
}
