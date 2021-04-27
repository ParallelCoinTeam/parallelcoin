// +build !headless

package launchers

import (
	"fmt"

	"github.com/p9c/matrjoska/cmd/gui"
	"github.com/p9c/matrjoska/pod/state"
)

// GUIHandle starts up the GUI wallet
func GUIHandle(ifc interface{}) (e error) {
	var cx *state.State
	var ok bool
	if cx, ok = ifc.(*state.State); !ok {
		return fmt.Errorf("cannot run without a state")
	}
	// // log.AppColorizer = color.Bit24(128, 255, 255, false).Sprint
	// // log.App = "   gui"
	D.Ln("starting up parallelcoin pod gui...")/**/
	// // fork.ForkCalc()
	// // podconfig.Configure(cx, true)
	// // D.Ln(os.Args)
	// // interrupt.AddHandler(func() {
	// // 	D.Ln("wallet gui is shut down")
	// // })
	if e = gui.Main(cx); E.Chk(e) {
	}
	D.Ln("pod gui finished")
	return
}
