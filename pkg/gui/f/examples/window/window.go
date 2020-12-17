package main

import (
	"gioui.org/layout"
	"github.com/p9c/pod/pkg/gui/fonts/p9fonts"
	"github.com/p9c/pod/pkg/gui/p9"
	qu "github.com/p9c/pod/pkg/util/quit"
	
	"github.com/p9c/pod/pkg/gui/f"
)

func main() {
	quit := qu.T()
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
