package pages

import (
	"github.com/p9c/pod/cmd/gui/component"
	"github.com/p9c/pod/cmd/gui/rcd"
	"time"

	"gioui.org/layout"
	"github.com/p9c/pod/pkg/gelook"
)

var (
	logOutputList = &layout.List{
		Axis:        layout.Vertical,
		ScrollToEnd: true,
	}
)

var StartupTime = time.Now()

func Logger(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) *gelook.DuoUIpage {
	page := gelook.DuoUIpage{
		Title:         "LOG",
		Command:       func() {},
		Border:        4,
		BorderColor:   th.Colors["Light"],
		Header:        func() {},
		HeaderBgColor: "",
		Body:          component.DuoUIlogger(rc, gtx, th),
		BodyBgColor:   th.Colors["Dark"],
		Footer:        func() {},
		FooterBgColor: "",
	}
	return th.DuoUIpage(page)
}
