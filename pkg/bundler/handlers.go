package bnd

import (
	"fmt"
	"net/http"
)

func (b *DuOSsveBundle) DuOSassetsHandler(w http.ResponseWriter, r *http.Request) {

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
