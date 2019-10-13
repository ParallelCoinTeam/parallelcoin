//+build !nogui
// +build !headless

package vueOLD

import (
	enjs "encoding/json"
	"github.com/zserge/webview"
	"log"
	"strings"
)

func (dV *DuoVUE) Render(cmd string, data interface{}) {
	b, err := enjs.Marshal(data)
	if err == nil {
		dV.Web.Eval("duoSystem." + cmd + "=" + string(b) + ";")
	}
}

func (dV *DuoVUE) HandleRPC(w webview.WebView, vc string) {
	switch {
	case vc == "close":
		dV.Web.Terminate()
	case vc == "fullscreen":
		dV.Web.SetFullscreen(true)
	case vc == "unfullscreen":
		dV.Web.SetFullscreen(false)
	case strings.HasPrefix(vc, "changeTitle:"):
		dV.Web.SetTitle(strings.TrimPrefix(vc, "changeTitle:"))
	//case vc == "status":
	//	dV.cr.AddFunc("@every 1s", func() {
	//		dV.Web.Dispatch(func() {
	//			//dV.Render("status", dV.GetDuoVUEstatus())
	//		})
	//	})
	//case vc == "peers":
	//	dV.cr.AddFunc("@every 3s", func() {
	//		dV.Web.Dispatch(func() {
	//			//dV.Render("status", dV.GetPeerInfo())
	//		})
	//	})
	case vc == "addressbook":
		dV.Render(vc, dV.GetAddressBook())
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
		dV.Render("transactions", dV.GetTransactions(cmd.From, cmd.Count, cmd.C))
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

		//dV.Render("send", dV.DuoSend(cmd.Wp, cmd.Ad, cmd.Am))
	case strings.HasPrefix(vc, "createAddress:"):
		s := strings.TrimPrefix(vc, "createAddress:")
		cmd := struct {
			Account string `json:"account"`
		}{}
		if err := enjs.Unmarshal([]byte(s), &cmd); err != nil {
			log.Println(err)
		}
		b, err := enjs.Marshal(dV.CreateNewAddress(cmd.Account))
		if err == nil {
			dV.Web.Eval("createAddress=" + string(b) + ";")
		}
		//dV.Render("createAddress", dV.CreateNewAddress(cmd.Account))
	case strings.HasPrefix(vc, "saveAddressLabel:"):
		s := strings.TrimPrefix(vc, "saveAddressLabel:")
		cmd := struct {
			Address string `json:"address"`
			Label   string `json:"label"`
		}{}
		if err := enjs.Unmarshal([]byte(s), &cmd); err != nil {
			log.Println(err)
		}
		//dV.Render("saveAddressLabel", dV.SaveAddressLabel(cmd.Address, cmd.Label))
		dV.SaveAddressLabel(cmd.Address, cmd.Label)

	}

}
