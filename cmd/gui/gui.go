package gui

import (
	"fmt"
	"github.com/p9c/pod/pkg/gui/webview"
	"github.com/shurcooL/vfsgen"
	"log"
	"net"
	"net/http"
)

func GUI() {
	var fs http.FileSystem = http.Dir("./pkg/gui/vue/dist")
	err := vfsgen.Generate(fs, vfsgen.Options{})
	if err != nil {
		log.Fatalln(err)
	}
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	go http.Serve(ln, http.FileServer(fs))

	w := webview.New(webview.Settings{
		Width:  800,
		Height: 600,
		Title:  "DUo",
		URL:    fmt.Sprintf("http://%s", ln.Addr()),
		//ExternalInvokeCallback: handleRPC,
	})
	defer w.Exit()
	w.Run()
}

/*
func GUsssI() {

	var fs http.FileSystem = http.Dir("./pkg/gui/vue/dist")
	err := vfsgen.Generate(fs, vfsgen.Options{})
	if err != nil {
		log.Fatalln(err)
	}

	args := []string{}
	if runtime.GOOS == "linux" {
		args = append(args, "--class=ParallelCoin")
	}
	ui, err := engine.New("", "", 480, 320, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()

	// A simple way to know when UI is ready (uses body.onload event in JS)
	ui.Bind("start", func() {
		log.Println("UI is ready")
	})

	// Create and bind Go object to the UI
	c := &counter{}
	ui.Bind("counterAdd", c.Add)
	ui.Bind("counterValue", c.Value)

	// Load HTML.
	// You may also use `data:text/html,<base64>` approach to load initial HTML,
	// e.g: ui.Load("data:text/html," + url.PathEscape(html))

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	go http.Serve(ln, http.FileServer(fs))

	//ui.Load(fmt.Sprintf("http://%s", ln.Addr()))

	ui.Load("data:text/html," + url.PathEscape(indexHTML))

	// You may use console.log to debug your JS code, it will be printed via
	// log.Println(). Also exceptions are printed in a similar manner.
	ui.Eval(`
		console.log("Hello, world!");
		console.log('Multiple values:', [1, false, {"x":5}]);
	`)

	// Wait until the interrupt signal arrives or browser window is closed
	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt)
	select {
	case <-sigc:
	case <-ui.Done():
	}

	log.Println("exiting...")
}
*/
