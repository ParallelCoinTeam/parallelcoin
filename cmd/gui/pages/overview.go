package pages

import (
	"gioui.org/layout"
	"gioui.org/op"

	"github.com/p9c/pod/cmd/gui/component"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gui/gelook"
)

func Overview(rc *rcd.RcVar, gtx *layout.Context,
	th *gelook.DuoUItheme) *gelook.DuoUIpage {
	page := gelook.DuoUIpage{
		Title:         "OVERVIEW",
		Border:        0,
		Command:       func() {},
		Header:        func() {},
		HeaderBgColor: "",
		HeaderPadding: 0,
		Body:          overviewBody(rc, gtx, th),
		BodyBgColor:   th.Colors["Light"],
		BodyPadding:   0,
		Footer:        func() {},
		FooterBgColor: "",
		FooterPadding: 0,
	}
	return th.DuoUIpage(page)
}

func overviewBody(rc *rcd.RcVar, gtx *layout.Context,
	th *gelook.DuoUItheme) func() {
	return func() {
		viewport := layout.Flex{Axis: layout.Horizontal}
		if gtx.Constraints.Width.Max < 780 {
			viewport = layout.Flex{Axis: layout.Vertical}
		}
		viewport.Layout(gtx,
			layout.Flexed(0.5,
				component.DuoUIstatus(rc, gtx, th)),
			layout.Flexed(0.5,
				component.DuoUIlatestTransactions(rc, gtx, th)),
		)
		op.InvalidateOp{}.Add(gtx.Ops)
	}
}
