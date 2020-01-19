package loader

import (
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/pkg/gui/layout"
)

func DuoUIloaderIntro(duo *models.DuoUI) {
	layout.Flex{}.Layout(duo.DuoUIcontext,
		layout.Flexed(1, func() {
			cs := duo.DuoUIcontext.Constraints
			//helpers.DuoUIdrawRect(duo.DuoUIcontext, cs.Width.Max, cs.Height.Max, color.RGBA{A: 0xff, R: 0x30, G: 0x30, B: 0x30}, 0, 0, 0, 0, unit.Dp(0))
			helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, 48, helpers.HexARGB("ff3030cf"), [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})

			// START View <<<

			layout.Flex{}.Layout(duo.DuoUIcontext,
				layout.Rigid(func() {
					//ldr.ico.Logo.Color = color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}
					//ldr.ico.Logo.Layout(duo.DuoUIcontext, unit.Dp(256))
				}))

			layout.Flex{}.Layout(duo.DuoUIcontext,
				layout.Rigid(func() {

				}))

			//ldr.comp.View.Layout.Layout(duo.DuoUIcontext, DuoUIloaderLogo(ldr))
		}))
}

//
//func DuoUIloaderLogo(ldr *DuoUIload) layout.FlexChild {
//	return ldr.comp.Intro.Layout.Rigid(duo.DuoUIcontext, func() {
//		logo := ldr.comp.Intro.Layout.Rigid(duo.DuoUIcontext, func() {
//			ldr.ico.Logo.Color = color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}
//			ldr.ico.Logo.Layout(duo.DuoUIcontext, unit.Dp(256))
//		})
//
//		ldr.comp.Intro.Layout.Layout(duo.DuoUIcontext, logo)
//	})
//}
