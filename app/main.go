// Package app is a multi-function universal binary that does all the things.
//
// Parallelcoin Pod
//
// This is the heart of configuration and coordination of
// the parts that compose the parallelcoin Pod - Ctl, Node and Wallet, and
// the extended, combined Shell and the Gio GUI.
package app

import (
	"github.com/p9c/pod/pkg/conte"
	log "github.com/p9c/logi"
	"os"
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
	log.L.Debug("running App")

	e := cx.App.Run(os.Args)
	if e != nil {
		log.Println("Pod ERROR:", e)
		return 1
	}

	return 0
}
