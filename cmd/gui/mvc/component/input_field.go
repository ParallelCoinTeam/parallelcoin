package component

import (
	"gioui.org/layout"
	"gioui.org/text"
	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
	"github.com/p9c/pod/cmd/gui/rcd"
)

func DuoUIinputField(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme, f *Field) func() {
	return func() {
		e := th.DuoUIeditor(f.Field.Label)
		e.Font.Typeface = th.Font.Primary
		e.Font.Style = text.Italic
		e.Layout(gtx, (rc.Settings.Daemon.Widgets[f.Field.Label]).(*controller.Editor))
		//(ui.rc.Settings.Daemon.Widgets[f.Field.Label]).(*controller.Editor).SetText(f.Field.Value.(reflect.Value).String())
		//log.INFO(f.Field.Value.(reflect.Value).String())
		//for _, e := range lineEditor.Events(ui.ly.Context) {
		//	if _, ok := e.(controller.SubmitEvent); ok {
		//		//topLabel = e.Text
		//		lineEditor.SetText(f.Field.Value.(reflect.Value).String())
		//		log.INFO(f.Field.Value.(reflect.Value).String())
		//	}
		//}
	}
}
