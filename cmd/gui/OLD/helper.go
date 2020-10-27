package gui

import (
	"gioui.org/text"
	"github.com/p9c/pod/pkg/gui/wallet/lyt"
	"github.com/p9c/pod/pkg/gui/wallet/theme"
)

func labeledRow(th *theme.Theme, label string, content func(gtx C) D) func(gtx C) D {
	return func(gtx C) D {
		return lyt.Format(gtx, "hflex(start,r(inset(8dp0dp0dp8dp,_)),f(1,_))",
			func(gtx C) D {
				gtx.Constraints.Min.X = 80
				gtx.Constraints.Min.X = 80
				title := theme.Body(th, label)
				title.Alignment = text.Start
				return title.Layout(gtx)
			},
			content,
		)
	}
}
