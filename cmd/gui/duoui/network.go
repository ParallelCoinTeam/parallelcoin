package duoui

import (
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/unit"
)

func DuoUInetwork(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) {
	layout.Flex{}.Layout(duo.Gc,
		layout.Flexed(1, func() {
			duo.Comp.Network.Inset.Layout(duo.Gc, func() {
				helpers.DuoUIdrawRectangle(duo.Gc, duo.Cs.Width.Max, duo.Cs.Height.Max, helpers.HexARGB("ff30cfcf"), [4]float32{0, 0, 0, 0}, unit.Dp(0))
				// Overview <<<
				in := layout.UniformInset(unit.Dp(60))
				in.Layout(duo.Gc, func() {
					helpers.DuoUIdrawRectangle(duo.Gc, duo.Cs.Width.Max, duo.Cs.Height.Max, helpers.HexARGB("ffcf30cf"), [4]float32{0, 0, 0, 0}, unit.Dp(0))

					duo.Th.H5("network :").Layout(duo.Gc)
				})
				// Overview >>>
			})
		}),
	)
}
