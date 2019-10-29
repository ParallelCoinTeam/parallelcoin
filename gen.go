//+build generate

package main

import "github.com/p9c/gui"

func main() {
	// You can also run "npm build" or webpack here, or compress assets, or
	// generate manifests, or do other preparations for your assets.
	gui.Embed("ini", "pkg/duos/ini/assets.go", "pkg/svelte/frontend/public")
}
