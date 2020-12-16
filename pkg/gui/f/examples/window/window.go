package main

import (
	"gioui.org/layout"
	"github.com/p9c/pod/pkg/gui/fonts/p9fonts"
	"github.com/p9c/pod/pkg/gui/p9"

	"github.com/p9c/pod/pkg/gui/f"
)

func main() {
	quit := make(qu.C)
	if err := f.NewWindow(p9.NewTheme(p9fonts.Collection(), quit)).Title("Parallelcoin").Size(10, 20).
		Run(func(layout.Context) layout.Dimensions {
			Info("frame")
			return layout.Dimensions{}
		}, func(layout.Context) {
			Info("overlay")
		}, func() {
			Info("destroy")
		}, quit); Check(err) {
	}
	<-quit
	// app.Main()
}
