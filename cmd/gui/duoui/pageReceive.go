package duoui

import (
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/unit"
)

func DuoUIreceive(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) {
	layout.Flex{}.Layout(duo.DuoUIcontext,
		layout.Flexed(1, func() {
			duo.DuoUIcomponents.AddressBook.Inset.Layout(duo.DuoUIcontext, func() {
				helpers.DuoUIdrawRectangle(duo.DuoUIcontext, duo.DuoUIconstraints.Width.Max, duo.DuoUIconstraints.Height.Max, helpers.HexARGB("ff888888"), [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
				// Overview <<<
				in := layout.UniformInset(unit.Dp(1))
				in.Layout(duo.DuoUIcontext, func() {
					helpers.DuoUIdrawRectangle(duo.DuoUIcontext, duo.DuoUIconstraints.Width.Max, duo.DuoUIconstraints.Height.Max, helpers.HexARGB("ffcfcfcf"), [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})

					duo.DuoUItheme.H5("receive :").Layout(duo.DuoUIcontext)
				})
				// Overview >>>
			})
		}),
	)
}
