package duoui

import (
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gui/layout"
)

func (duo *DuoUI) DuoUIbody(cx *conte.Xt, rc *rcd.RcVar) func() {
	return func() {
		layout.Flex{Axis: layout.Horizontal}.Layout(duo.Model.DuoUIcontext,
			layout.Rigid(duo.DuoUIsidebar(cx, rc)),
			layout.Flexed(1, DuoUIcontent(duo, cx, rc)),
		)
	}
}
