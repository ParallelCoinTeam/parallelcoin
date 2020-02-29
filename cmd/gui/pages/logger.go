package pages

import (
	"github.com/p9c/pod/cmd/gui/mvc/component"
	"github.com/p9c/pod/cmd/gui/rcd"
	"time"

	"gioui.org/layout"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
)

var (
	logOutputList = &layout.List{
		Axis:        layout.Vertical,
		ScrollToEnd: true,
	}
)

var StartupTime = time.Now()

func Logger(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme) *theme.DuoUIpage {
	return th.DuoUIpage("LOG", 0, func() {}, func() {}, component.DuoUIlogger(rc, gtx, th), func() {})
}
