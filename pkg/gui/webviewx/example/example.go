package main

import (
	"github.com/p9c/pod/pkg/gui/webview"
)

func main() {
	w := webviewx.New(true)
	w.Navigate("https://github.com")
	w.SetTitle("Hello")
	w.Dispatch(func() {
		println("Hello dispatch")
	})
	w.Run()
}
