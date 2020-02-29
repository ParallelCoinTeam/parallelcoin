package pages

import (
	"gioui.org/layout"
	"gioui.org/op"
	"github.com/p9c/pod/cmd/gui/mvc/component"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
	"github.com/p9c/pod/cmd/gui/rcd"
)

func Overview(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme) *theme.DuoUIpage {
	return th.DuoUIpage("OVERVIEW", 0, func() {}, func() {}, overviewBody(rc, gtx, th), func() {})
}

func overviewBody(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme) func() {
	return func() {
		viewport := layout.Flex{Axis: layout.Horizontal}
		if gtx.Constraints.Width.Max < 780 {
			viewport = layout.Flex{Axis: layout.Vertical}
		}
		viewport.Layout(gtx,
			layout.Flexed(0.5, component.DuoUIbalance(rc, gtx, th)),
			layout.Flexed(0.5, component.DuoUIlatestTransactions(rc, gtx, th)),
		)
		op.InvalidateOp{}.Add(gtx.Ops)
	}
}
