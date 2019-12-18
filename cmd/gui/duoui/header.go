package duoui


import (
	"gioui.org/layout"
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/ico"
	"golang.org/x/exp/shiny/iconvg"
	"image"
	"image/color"
	"image/draw"
)

func DuoUIheader(duo *DuoUI) layout.FlexChild {
	return duo.comp.View.Layout.Rigid(duo.gc, func() {
		helpers.DuoUIdrawRect(duo.gc, duo.cs.Width.Max, 64, color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}, 0, 0, 0, 0)
		// Header <<<
		logo := duo.comp.Header.Layout.Rigid(duo.gc, func() {
			helpers.DuoUIdrawRect(duo.gc, 64, 64, color.RGBA{A: 0xff, R: 0x30, B: 0x30, G: 0x30}, 0, 0, 0, 0)

			sz := 48
			m, _ := iconvg.DecodeMetadata(ico.ParallelCoin)
			dx, dy := m.ViewBox.AspectRatio()
			img := image.NewRGBA(image.Rectangle{Max: image.Point{X: sz, Y: int(float32(sz) * dy / dx)}})
			var parallelcoinLogo iconvg.Rasterizer
			parallelcoinLogo.SetDstImage(img, img.Bounds(), draw.Src)
			// Use white for icons.
			m.Palette[0] = color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}
			iconvg.Decode(&parallelcoinLogo, ico.ParallelCoin, &iconvg.DecodeOptions{
				Palette: &m.Palette,
			})

		})
		balance := duo.comp.Header.Layout.Rigid(duo.gc, func() {
			duo.th.H5(duo.rc.Balance + " DUO").Layout(duo.gc)
		})
		duo.comp.Header.Layout.Layout(duo.gc, logo, balance)
		// Header >>>
	})
}
