package loader

import (
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/unit"
	"image/color"
)

func DuoUIloaderIntro(ldr *DuoUIload) {
	layout.Flex{}.Layout(ldr.gc,
		layout.Flexed(1, func() {
			helpers.DuoUIdrawRect(ldr.gc, ldr.cs.Width.Max, ldr.cs.Height.Max, color.RGBA{A: 0xff, R: 0x30, G: 0x30, B: 0x30}, 0, 0, 0, 0)
			// START View <<<

			layout.Flex{}.Layout(ldr.gc,
				layout.Rigid(func() {
				ldr.ico.Logo.Color = color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}
				ldr.ico.Logo.Layout(ldr.gc, unit.Dp(256))
			}))

			layout.Flex{}.Layout(ldr.gc,
				layout.Rigid(func() {

				}))

			//ldr.comp.View.Layout.Layout(ldr.gc, DuoUIloaderLogo(ldr))
		}))
}

//
//func DuoUIloaderLogo(ldr *DuoUIload) layout.FlexChild {
//	return ldr.comp.Intro.Layout.Rigid(ldr.gc, func() {
//		logo := ldr.comp.Intro.Layout.Rigid(ldr.gc, func() {
//			ldr.ico.Logo.Color = color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}
//			ldr.ico.Logo.Layout(ldr.gc, unit.Dp(256))
//		})
//
//		ldr.comp.Intro.Layout.Layout(ldr.gc, logo)
//	})
//}
