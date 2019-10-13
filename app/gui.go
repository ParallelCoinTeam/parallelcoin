package app

import (
	"fmt"
	"github.com/p9c/pod/pkg/duos/core"
)

func gui(d *core.DuOS) {

	defer d.GuI.Close()

	// Wait until UI window is closed
	// Load HTML after Go functions are bound to JS

	d.GuI.Eval(`
		console.log("Hello, world!");
		console.log('Multiple values:', [1, false, {"x":5}]);
	`)

	// Css
	//injectCss(d)

	// Js

	//d.GetPeerInfo()

	//d.EvalJs()

	//vueLib, err := base64.StdEncoding.DecodeString(lib.VUE)
	//if err != nil {
	//	fmt.Printf("Error decoding string: %s ", err.Error())
	//	return
	//}
	//d.GuI.Eval(string(vueLib))

	//ej2
	//getEj2Vue, err := base64.StdEncoding.DecodeString(lib.EJS)
	//if err != nil {
	//	fmt.Printf("Error decoding string: %s ", err.Error())
	//	return
	//}
	//d.GuI.Eval(string(getEj2Vue))
	//
	//

	fmt.Println("teeeeeee")
	<-d.GuI.Done()

	// Start ticker goroutine
	//go func() {
	//	t := time.NewTicker(100 * time.Millisecond)
	//	for {
	//		select {
	//		case <-t.C: // Every 100ms increate number of ticks and update UI
	//			ui.Eval("document.querySelector('.timer').innerText =" + fmt.Sprint(getBlockCount))
	//		case <-togglec: // If paused - wait for another toggle event to unpause
	//			<-togglec
	//		}
	//	}
	//}()
	//<-ui.Done()
}
