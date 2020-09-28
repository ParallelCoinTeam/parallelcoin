// Package app is a multi-function universal binary that does all the things.
//
// Parallelcoin Pod
//
// This is the heart of configuration and coordination of
// the parts that compose the parallelcoin Pod - Ctl, Node and Wallet, and
// the extended, combined Shell and the Gio GUI.
package app

import (
	"github.com/stalker-loki/app/slog"
	"os"

	"github.com/p9c/pod/app/conte"
)

const (
	appName           = "pod"
	confExt           = ".json"
	appLanguage       = "en"
	podConfigFilename = appName + confExt
	PARSER            = "json"
)

// Main is the entrypoint for the pod AiO suite
func Main() int {
	cx := conte.GetNewContext(appName, appLanguage, "main")
	cx.App = GetApp(cx)
	if e := cx.App.Run(os.Args); slog.Check(e) {
		return 1
	}
	return 0
}
