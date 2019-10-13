package core

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/p9c/pod/cmd/node/rpc"
	"github.com/p9c/pod/pkg/duos/db"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	"github.com/p9c/pod/pkg/stat"
	"github.com/p9c/pod/pkg/svelte/__OLDvue/lib"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"html/template"
	"time"
)

func (d *DuOS) EvalJs() {
	//svelte
	vueLib, err := base64.StdEncoding.DecodeString(lib.VUE)
	if err != nil {
		fmt.Printf("Error decoding string: %s ", err.Error())
		return
	}
	vue := d.GuI.Eval(string(vueLib))
	fmt.Println(vue)
	//ej2
	getEj2Vue, err := base64.StdEncoding.DecodeString(lib.EJS)
	if err != nil {
		fmt.Printf("Error decoding string: %s ", err.Error())
		return
	}
	ejs := d.GuI.Eval(string(getEj2Vue))
	fmt.Println(ejs)

	// libs
	//for _, lib := range lib.GetLibs() {
	//		return
	//	}
	//}
	//for _, js := range t.Data["js"] {
	//	err = w.Eval(string(js))
	//}
	// for _, js := range t.Data["js"] {
	// 	err = w.Eval(string(js))
	// }

	//err = d.GuI.Eval(CoreHeadJs)
	//if err != nil {
	//	fmt.Println("error binding to webview:", err)
	//}
	//
	//err = d.GuI.Eval(CompLoopJs(d.DbS))
	//if err != nil {
	//	fmt.Println("error binding to webview:", err)
	//}
	//
	//err = d.GuI.Eval(AppsLoopJs(d.DbS))
	//if err != nil {
	//	fmt.Println("error binding to webview:", err)
	//}

	d.GuI.Eval(CoreJs(d.DbS))
	//fmt.Println("MIkaaaaaaaaaa:", CoreJs(d))
}

func CoreJs(d db.DuOSdb) string {
	vueapp := `
	{{define "vueapp"}}

Vue.config.devtools = true;
Vue.use(VueFormGenerator);

Vue.prototype.$eventHub = new Vue(); 

 
const duoSystem = {
	config:null,
	node:null,
	wallet:null,
	status:null,
	balance:null,
	connectionCount:0,
	addressBook:null,
	transactions:null,
	peerInfo:null,
	blocks:[],
	theme:false,
	isBoot:false,
	isLoading:false,
	isDev:true,
	isScreen:'overview',
	timer: '',
};


{{ range $key, $value := . }}
var {{ .ID }} = {{ if .IsApp }}new Vue({{ end }}{
{{ if .IsApp }}el: '#{{ .ID }}',{{ end }}
  name: '{{ .Name }}',
{{ if .Template }}template: ` + "`" + `{{ .Template }}` + "`" + `,{{end}}
{{ if .Js }}{{ .Js }}{{ end }}
}{{ if .IsApp }});{{ end }}
{{ end }}


const core = new Vue({ 
	el: '#core', 
	data () { return { 
	duoSystem }},
});

{{ end }}`
	var js bytes.Buffer
	jsRaw := template.Must(template.New("").Parse(string(vueapp)))
	err := jsRaw.ExecuteTemplate(&js, "vueapp", Components(d))
	if err != nil {
		fmt.Println("error binding to webview:", err)
	}

	return js.String()
}

// GetMsg loads the message variable
func (d *DuOS) PushDuOSalert(t string, m interface{}, at string) {
	a := new(DuOSalert)
	a.Time = time.Now()
	a.Title = t
	a.Message = m
	a.AlertType = at
	//d.Render("alert", a)
}

func (d *DuOS) GetDuOSstatus() stat.DuOSstatus {
	status := *new(stat.DuOSstatus)
	sm, _ := mem.VirtualMemory()
	sc, _ := cpu.Info()
	sp, _ := cpu.Percent(0, true)
	sd, _ := disk.Usage("/")
	status.Cpu = sc
	status.CpuPercent = sp
	status.Memory = *sm
	status.Disk = *sd
	params := d.CtX.RPCServer.Cfg.ChainParams
	chain := d.CtX.RPCServer.Cfg.Chain
	chainSnapshot := chain.BestSnapshot()
	gnhpsCmd := btcjson.NewGetNetworkHashPSCmd(nil, nil)
	networkHashesPerSecIface, err := rpc.HandleGetNetworkHashPS(d.CtX.RPCServer, gnhpsCmd, nil)
	if err != nil {
	}
	networkHashesPerSec, ok := networkHashesPerSecIface.(int64)
	if !ok {
	}
	v, err := rpc.HandleVersion(d.CtX.RPCServer, nil, nil)
	if err != nil {
	}
	status.Version = "0.0.1"
	status.WalletVersion = v.(map[string]btcjson.VersionResult)
	status.UpTime = time.Now().Unix() - d.CtX.RPCServer.Cfg.StartupTime
	status.CurrentNet = d.CtX.RPCServer.Cfg.ChainParams.Net.String()
	status.NetworkHashPS = networkHashesPerSec
	status.HashesPerSec = int64(d.CtX.RPCServer.Cfg.CPUMiner.HashesPerSecond())
	status.Chain = params.Name
	status.Height = chainSnapshot.Height
	//s.Headers = chainSnapshot.Height
	status.BestBlockHash = chainSnapshot.Hash.String()
	status.Difficulty = rpc.GetDifficultyRatio(chainSnapshot.Bits, params, 2)
	status.Balance.Balance = d.GetBalance().Balance
	status.Balance.Unconfirmed = d.GetBalance().Unconfirmed
	status.BlockCount = d.GetBlockCount()
	status.ConnectionCount = d.GetConnectionCount()
	status.NetworkLastBlock = d.GetNetworkLastBlock()
	return status
}
