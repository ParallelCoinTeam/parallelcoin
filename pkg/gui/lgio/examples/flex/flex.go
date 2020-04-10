package main

import (
	"gioui.org/app"
	"gioui.org/layout"
	"github.com/p9c/pod/pkg/gui/lgio"
)

func main() {
	lgio.Window().Title("Parallelcoin").Size(640, 480).
		Run(func(c *layout.Context) {
			lgio.Flex().Vertical().Flexed(0.5,
				lgio.Flex().Flexed(0.5,
					lgio.Inset(32).Prepare(c,
						lgio.Widget().Fill(255, 128, 0, 255).Prepare(c),
					),
				).Flexed(0.25,
					lgio.Widget().Fill(0, 255, 128, 255).Prepare(c),
				).Flexed(0.25,
					lgio.Widget().Fill(128, 0, 255, 255).Prepare(c),
				).Prepare(c),
			).Flexed(0.25,
				lgio.Widget().Fill(0, 255, 0, 255).Prepare(c),
			).Flexed(0.25,
				lgio.Widget().Fill(0, 0, 255, 255).Prepare(c),
			).Layout(c)
		}, func() {

		})
	app.Main()
}
