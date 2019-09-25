package core

import (
	"bytes"
	"fmt"
	"github.com/parallelcointeam/parallelcoin/cmd/gui/vue/comp"
	"text/template"

	"github.com/parallelcointeam/parallelcoin/cmd/gui/vue/db"
)

var CoreJs = `
const core = new Vue({ 
	el: '#x', 
	data () { return { 
	duoSystem }},
});
`

var CoreHeadJs = `
Vue.config.devtools = true;
Vue.use(VueFormGenerator);

Vue.prototype.$eventHub = new Vue(); 

const duoSystem = {
	alert:system.data.d.alert,
	config:system.data.conf,
	status: system.data.d.status,
	addressbook:system.data.d.addressbook,
	createAddress:'',
	txsEx:system.data.d.txsex,
	peers:system.data.d.peers,
	blocks:[],
	theme:false,
	logo:system.data.ico,
	bios:{
		isBoot:true,
		isDev:false,
	},
	isLoading:false,
	isScreen: 'overview',
	timer: '',
};
`


func AppsLoopJs(d db.DuoVUEdb) string {
	vueapp := `
{{define "vueapp"}}
	{{ range $key, $value := . }}
		var {{ .ID }} = new Vue({
		el: '#{{ .ID }}',
		name: '{{ .Name }}',
		{{ if .Template }}template: ` + "`" + `{{ .Template }}` + "`" + `,{{end}}
		{{ if .Js }}{{ .Js }}{{ end }}
		});
	{{ end }}
{{ end }}`
	var js bytes.Buffer
	jsRaw := template.Must(template.New("").Parse(string(vueapp)))
	err := jsRaw.ExecuteTemplate(&js, "vueapp", comp.Apps(d))
	if err != nil {
		fmt.Println("error binding to webview:", err)
	}
	return js.String()
}

func CompLoopJs(d db.DuoVUEdb) string {
	vuecomp := `
{{define "vuecomp"}}
	{{ range $key, $value := . }}
		var {{ .ID }} = {
		name: '{{ .Name }}',
		{{ if .Template }}template: ` + "`" + `{{ .Template }}` + "`" + `,{{end}}
		{{ if .Js }}{{ .Js }}{{ end }}
		}
	{{ end }}
{{ end }}`
	var js bytes.Buffer
	jsRaw := template.Must(template.New("").Parse(string(vuecomp)))
	err := jsRaw.ExecuteTemplate(&js, "vuecomp", comp.Components(d))
	if err != nil {
		fmt.Println("error binding to webview:", err)
	}
	return js.String()
}