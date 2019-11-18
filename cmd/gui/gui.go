package gui

import (
	"encoding/json"
	"fmt"
	"github.com/p9c/pod/pkg/conte"
	"github.com/therecipe/qt/webengine"
	"time"
)

func GUI(cx *conte.Xt, view webengine.QWebEngineView) {

	rc := rcvar{
		cx:     cx,
		alert:  DuOSalert{},
		status: DuOStatus{},
		txs:    DuOStransactionsExcerpts{},
		lastxs: DuOStransactions{},
	}

	go func() {
		for _ = range time.NewTicker(time.Second * 1).C {



			//runJs("status", rc.GetDuOStatus(), view)
			//
			//runJs("txs", rc.GetTransactionsExcertps(), view)
			//
			//runJs("lastxs", rc.GetTransactions(0, 5, ""), view)

			view.Page().RunJavaScript(`
window.duoSystem = {
`+ evalVal("status", rc.GetDuOStatus()) +`
};


`)

			fmt.Println("Istos---->>>>>",rc.GetTransactionsExcertps())

			fmt.Println("Drugost---->>>>>", evalVal("txs", rc.GetTransactionsExcertps()))

			fmt.Println("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
			fmt.Println("aaaaaIstos---->>>>>", rc.GetTransactions(0, 5, ""))

			fmt.Println("aaaaaDrugost---->>>>>", evalVal("txs", rc.GetTransactions(0, 5, "")))
		}
	}()

}

func evalVal(n string, v interface{})string {
	vm, err := json.Marshal(v)
	if err != nil {
	}
	return n + `:` + string(vm) + `,`
}
