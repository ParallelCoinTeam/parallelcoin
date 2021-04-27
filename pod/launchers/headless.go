// +build headless

package launchers

import (
	"os"
)

func GUIHandle(ifc interface{}) (e error) {
	W.Ln("GUI was disabled for this podbuild (server only version)")
	os.Exit(1)
	return nil
}
