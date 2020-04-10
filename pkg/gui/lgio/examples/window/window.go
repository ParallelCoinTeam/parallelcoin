package main

import (
	"gioui.org/app"
	"gioui.org/layout"
	"github.com/p9c/pod/pkg/gui/lgio"
)

func main() {
	lgio.Window().Title("Parallelcoin").Size(640, 480).
		Run(func(*layout.Context){}, func() {})
	app.Main()
}
