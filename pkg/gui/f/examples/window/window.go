package main

import (
	"gioui.org/app"
	"gioui.org/layout"

	"github.com/p9c/pod/pkg/gui/f"
)

func main() {
	if err := f.Window().Title("Parallelcoin").Size(640, 480).
		Run(func(*layout.Context) {
			Info("frame")
		}, func() {
			Info("destroy")
		}); Check(err) {
	}
	app.Main()
}
