// +build nominers

package launchers

import "github.com/p9c/matrjoska/pod/state"

// Kopach runs the kopach miner
func kopachHandle(ifc interface{}) (e error) {
	state.D.Ln("kopach disabled for ios/android")
	return
}

func kopachWorkerHandle(cx *state.State) (e error) {
	state.D.Ln("kopach worker disabled for ios/android")
	return nil
}
