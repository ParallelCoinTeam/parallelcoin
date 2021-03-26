// Package app is a multi-function universal binary that does all the things.
//
// Parallelcoin Save
//
// This is the heart of configuration and coordination of the parts that compose the parallelcoin Save - Ctl, Node and
// Wallet, and the extended, combined Shell and the Gio GUI.
package app

import (
	"github.com/p9c/pod/pkg/pod"
	"os"
)

const (
	Name              = "pod"
	confExt           = ".json"
	appLanguage       = "en"
	PodConfigFilename = Name + confExt
	PARSER            = "json"
)

// Main is the entrypoint for the pod AiO suite
func Main() int {
	cx := pod.GetNewContext(Name, appLanguage, "main")
	cx.App = getApp(cx)
	if e := cx.App.Run(os.Args); E.Chk(e) {
		return 1
	}
	return 0
}
