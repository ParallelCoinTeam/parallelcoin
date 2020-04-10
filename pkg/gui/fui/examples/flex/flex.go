package main

import (
	"gioui.org/app"
	"gioui.org/layout"
	"github.com/p9c/pod/pkg/gui/fui"
)

func main() {
	fui.Window().Title("Parallelcoin").Size(640, 480).Run(func(c *layout.Context) {
		fui.Flex().Vertical().Flexed(0.5,
			fui.Flex().Flexed(0.5,
				fui.Inset(32).Prep(c,
					fui.Widget().Fill(255, 255, 0, 255).Prep(c),
				),
			).Flexed(0.25,
				fui.Widget().Fill(0, 255, 255, 255).Prep(c),
			).Flexed(0.25,
				fui.Widget().Fill(255, 0, 255, 255).Prep(c),
			).Prep(c),
		).Flexed(0.25,
			fui.Widget().Fill(0, 255, 0, 255).Prep(c),
		).Flexed(0.25,
			fui.Widget().Fill(0, 0, 255, 255).Prep(c),
		).Layout(c)
	}, func() {
	})
	app.Main()
}
