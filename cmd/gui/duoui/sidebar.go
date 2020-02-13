package duoui

import (
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gui/layout"
)

func (duo *DuoUI)DuoUIsidebar(cx *conte.Xt, rc *rcd.RcVar) func() {
	return func() {
		layout.Flex{Axis: layout.Vertical}.Layout(duo.m.DuoUIcontext,
			layout.Rigid(DuoUImenu(duo, cx, rc)),
		)
	}
}
