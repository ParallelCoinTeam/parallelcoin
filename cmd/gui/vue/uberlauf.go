//+build !nogui
// +build !headless

package vue

import (
	enjs "encoding/json"
	"github.com/zserge/webview"
	"log"
	"strings"
)

func (dv *DuoVUE) Render(cmd string, data interface{}) {
	b, err := enjs.Marshal(data)
	if err == nil {
		dv.Web.Eval("duoSystem." + cmd + "=" + string(b) + ";")
	}
}

func (dv *DuoVUE) HandleRPC(w webview.WebView, vc string) {
	switch {
	case vc == "close":
		dv.Web.Terminate()
	case vc == "fullscreen":
		dv.Web.SetFullscreen(true)
	case vc == "unfullscreen":
		dv.Web.SetFullscreen(false)
	case strings.HasPrefix(vc, "changeTitle:"):
		dv.Web.SetTitle(strings.TrimPrefix(vc, "changeTitle:"))
	case vc == "addressbook":
		dv.Render(vc, dv.GetAddressBook())
	case strings.HasPrefix(vc, "transactions:"):
		t := strings.TrimPrefix(vc, "transactions:")
		cmd := struct {
			From  int    `json:"from"`
			Count int    `json:"count"`
			C     string `json:"c"`
		}{}
		if err := enjs.Unmarshal([]byte(t), &cmd); err != nil {
			log.Println(err)
		}
		dv.Render("transactions", dv.GetTransactions(cmd.From, cmd.Count, cmd.C))
	case strings.HasPrefix(vc, "send:"):
		s := strings.TrimPrefix(vc, "send:")
		cmd := struct {
			Wp string  `json:"wp"`
			Ad string  `json:"ad"`
			Am float64 `json:"am"`
		}{}
		if err := enjs.Unmarshal([]byte(s), &cmd); err != nil {
			log.Println(err)
		}

		dv.Render("send", dv.DuoSend(cmd.Wp, cmd.Ad, cmd.Am))
	case strings.HasPrefix(vc, "createAddress:"):
		s := strings.TrimPrefix(vc, "createAddress:")
		cmd := struct {
			Account string `json:"account"`
			Label   string `json:"label"`
		}{}
		if err := enjs.Unmarshal([]byte(s), &cmd); err != nil {
			log.Println(err)
		}
		dv.Render("send", dv.CreateNewAddress(cmd.Account, cmd.Label))

	}

}
