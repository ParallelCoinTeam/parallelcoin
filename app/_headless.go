// +build headless

package app

import (
	"github.com/p9c/pod/pkg/opts"
	"os"
)

func GUIHandle(ifc interface{}) (e error) {
	opts.W.Ln("GUI was disabled for this podbuild (server only version)")
	os.Exit(1)
	return nil
}
