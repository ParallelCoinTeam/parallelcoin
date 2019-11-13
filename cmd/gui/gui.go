package gui

import (
	"encoding/json"
	"github.com/p9c/pod/pkg/conte"
	"net/http"
	"time"
)

func GUI(cx *conte.Xt) {

	rc := rcvar{
		cx:     cx,
		alert:  DuOSalert{},
		status: DuOStatus{},
		txs:    DuOStransactionsExcerpts{},
		lastxs: DuOStransactions{},
	}

	go func() {
		for _ = range time.NewTicker(time.Second * 1).C {
			rc.GetDuOStatus()
			rc.GetTransactions(0, 5, "")
		}
	}()

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		js, err := json.Marshal(rc.status)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(js)
	})

	http.HandleFunc("/lastxs", func(w http.ResponseWriter, r *http.Request) {
		js, err := json.Marshal(rc.lastxs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(js)
	})

	go http.ListenAndServe(":3999", nil)

	//var fs http.FileSystem = http.Dir("./pkg/gui/svelte/assets")
	//err := vfsgen.Generate(fs, vfsgen.Options{})
	//if err != nil {
	//	log.FATAL("Shuttingdown GUI", err)
	//	os.Exit(1)
	//}
	//
	//ln, err := net.Listen("tcp", "127.0.0.1:0")
	//if err != nil {
	//	log.FATAL("Shuttingdown GUI", err)
	//	os.Exit(1)
	//}
	//defer ln.Close()
	//go http.Serve(ln, http.FileServer(fs))

}
