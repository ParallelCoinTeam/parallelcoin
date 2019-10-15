package bnd

import (
	"fmt"
	"net/http"
)

func path(a DuOSasset) (path string) {
	if a.Sub != "" {
		path = a.Sub + "/" + a.Name
	} else {
		path = a.Name
	}
	return
}

func DuOSassetsHandler() {
	for _, a := range Assets() {
		http.HandleFunc(path(a), func(w http.ResponseWriter, r *http.Request) {
			headers := w.Header()
			headers["Content-Type"] = []string{a.ContentType}
			fmt.Fprint(w, string(a.DataRaw))
		})
	}

}

// Getting the headers so we can set the correct mime type
func PipeJsHandler(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()
	headers["Content-Type"] = []string{"application/javascript"}
	fmt.Fprint(w, string(DuOSsveBundler()["pipe.js"]))
}

// Getting the headers so we can set the correct mime type
func BndJsHandler(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()
	headers["Content-Type"] = []string{"application/javascript"}
	fmt.Fprint(w, string(DuOSsveBundler()["svelte.js"]))
}

// Getting the headers so we can set the correct mime type
func BndCssHandler(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()
	headers["Content-Type"] = []string{"text/css"}
	fmt.Fprint(w, string(DuOSsveBundler()["svelte.css"]))
}
