package gui

import (
	"encoding/json"
	"fmt"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gui/engine"
	"github.com/p9c/pod/pkg/log"
	"github.com/shurcooL/vfsgen"
	"net"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"time"
)

// Go types that are bound to the UI must be thread-safe, because each binding
// is executed in its own goroutine. In this simple case we may use atomic
// operations, but for more complex cases one should use proper synchronization.
//type counter struct {
//	sync.Mutex
//	count int
//}
//
//func (c *counter) Add(n int) {
//	c.Lock()
//	defer c.Unlock()
//	c.count = c.count + n
//}
//
//func (c *counter) Value() int {
//	c.Lock()
//	defer c.Unlock()
//	return c.count
//}

func HTML() string {
	return `<!DOCTYPE html>
<html lang=en>
<head>
    <meta charset=utf-8>
    <meta http-equiv=X-UA-Compatible content="IE=edge">
    <meta name=viewport content="width=device-width,initial-scale=1">
    <link rel=icon href=favicon.ico>
    <title>duo</title>
</head>
<body></body>
</html>`
}

func GUI(cx *conte.Xt) {
	r := rcvar{
		cx:     cx,
		alert:  DuOSalert{},
		status: DuOStatus{},
	}

	var fs http.FileSystem = http.Dir("./pkg/gui/vue/dist")
	err := vfsgen.Generate(fs, vfsgen.Options{})
	if err != nil {
		log.FATAL("Shuttingdown GUI", err)
		os.Exit(1)
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.FATAL("Shuttingdown GUI", err)
		os.Exit(1)
	}
	defer ln.Close()
	go http.Serve(ln, http.FileServer(fs))

	args := []string{}
	if runtime.GOOS == "linux" {
		args = append(args, "--class=ParallelCoin")
	}
	//ui.Load("data:text/html," + url.PathEscape(indexHTML))
	//ui, err := engine.New(HTML(fmt.Sprintf("http://%s", ln.Addr())), "", 800, 600, args...)
	//ui, err := engine.New("", "", 800, 600, args...)
	ui, err := engine.New("data:text/html,"+url.PathEscape(HTML()), "", 800, 600, args...)
	if err != nil {
		log.FATAL("Shuttingdown GUI", err)
		os.Exit(1)
	}
	defer ui.Close()

	//ui.Load(url.PathEscape(HTML(fmt.Sprintf("http://%s", ln.Addr()))))

	// A simple way to know when UI is ready (uses body.onload event in JS)
	ui.Bind("start", func() {
		log.DEBUG("UI is ready", err)
		fmt.Println("dadadadadadaddadadadadadaddadadadadadaddadadadadadaddadadadadadad")
		fmt.Println("dadadadadadaddadadadadadaddadadadadadaddadadadadadaddadadadadadad")
		fmt.Println("dadadadadadaddadadadadadaddadadadadadaddadadadadadaddadadadadadad")
		fmt.Println("dadadadadadaddadadadadadaddadadadadadaddadadadadadaddadadadadadad")
		fmt.Println("dadadadadadaddadadadadadaddadadadadadaddadadadadadaddadadadadadad")
		fmt.Println("dadadadadadaddadadadadadaddadadadadadaddadadadadadaddadadadadadad")
	})

	ui.Eval(`
		console.log("Hello, world!");
		console.log('Multiple values:', [1, false, {"x":5}]);
	`)

	go func() {
		//some other thread
		for _ = range time.NewTicker(time.Second * 1).C {

			stat, err := json.Marshal(r.GetDuOStatus())
			if err != nil {
			}
			//ui.Eval(`console.log('Multiple values:', ` + string(stat) + `;`)

			ui.Eval(`console.log('Multiple values:', [ ` + string(stat) + ` ]);`)
			//fmt.Println("dadaddadadadad", string(stat) )

		}
	}()

	//// Create and bind Go object to the UI
	//c := &counter{}
	//ui.Bind("counterAdd", c.Add)
	//ui.Bind("counterValue", c.Value)

	<-ui.Done()
}
