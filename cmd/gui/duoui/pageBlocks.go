package duoui

import (
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/unit"
)

func DuoUIexplorer(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) {
	layout.Flex{}.Layout(duo.DuoUIcontext,
		layout.Flexed(1, func() {
			duo.DuoUIcomponents.Explorer.Inset.Layout(duo.DuoUIcontext, func() {
				helpers.DuoUIdrawRectangle(duo.DuoUIcontext, duo.DuoUIconstraints.Width.Max, duo.DuoUIconstraints.Height.Max, helpers.HexARGB("ff30cfcf"), [4]float32{0, 0, 0, 0}, unit.Dp(0))
				// Overview <<<
				in := layout.UniformInset(unit.Dp(60))
				in.Layout(duo.DuoUIcontext, func() {
					helpers.DuoUIdrawRectangle(duo.DuoUIcontext, duo.DuoUIconstraints.Width.Max, duo.DuoUIconstraints.Height.Max, helpers.HexARGB("ffcf30cf"), [4]float32{0, 0, 0, 0}, unit.Dp(0))

					duo.DuoUItheme.H5("explorer :").Layout(duo.DuoUIcontext)
				})
				// Overview >>>
			})
		}),
	)
}
