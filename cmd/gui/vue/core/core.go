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
	el: '#core', 
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
	config:null,
	status: system.data.d.status,
	addressbook:null,
	transactions:null,
	peers:null,
	blocks:[],
	theme:false,
	isBoot:false,
	isLoading:false,
	isDev:true,
	isScreen:'overview',
	timer: '',
};
`


func AppsLoopJs(d db.DuoVUEdb) string {
	vueapp := `
{{define "vueapp"}}
	{{ range $key, $value := . }}
		var {{ .ID }} = {{ if .IsApp }}new Vue({{ end }}{
		el: '#{{ .ID }}',
		name: '{{ .Name }}',
		{{ if .Template }}template: ` + "`" + `{{ .Template }}` + "`" + `,{{end}}
		{{ if .Js }}{{ .Js }}{{ end }}
		});
	{{ end }}
{{ end }}`
	var js bytes.Buffer
	jsRaw := template.Must(template.New("").Parse(string(vueapp)))
	err := jsRaw.ExecuteTemplate(&js, "vueapp", comp.Components(d))
	if err != nil {
		fmt.Println("error binding to webview:", err)
	}
	return js.String()
}