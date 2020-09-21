package main

import (
	"gioui.org/app"
	"gioui.org/layout"
	"github.com/stalker-loki/pod/pkg/gui/fui"
)

func main() {
	fui.Window().Title("Parallelcoin").Size(640, 480).
		Run(func(*layout.Context){}, func() {})
	app.Main()
}
