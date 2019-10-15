package core

import (
	"github.com/therecipe/qt/webkit"
	"github.com/therecipe/qt/widgets"
)

type GuI struct {
	//*widgets.QApplication
	QT      *widgets.QApplication `json:"qt"`
	WebView *webkit.QWebView      `json:"web"`
	Window  *widgets.QMainWindow  `json:"win"`
	Screen  *widgets.QWidget      `json:"scr"`
}

//func (d *DuOS) NewWindow(parent widgets.QWidget_ITF) *widgets.QMainWindow {
//	var window = new(widgets.QMainWindow)
//	window.QWidget = *widgets.NewQWidget(parent, 0)
//	window.SetWindowTitle("ParallelCoin")
//	window.Resize2(800, 600)
//	d.GuI.WebView = webkit.NewQWebView(window)
//	url := `https://parallelcoin.info/`
//	qurl := core.NewQUrl3(url, core.QUrl__TolerantMode)
//	d.GuI.WebView.Load(qurl)
//	return window
//}
