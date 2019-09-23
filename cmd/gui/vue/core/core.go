package core

import (
	"bytes"
	"fmt"
	"github.com/p9c/pod/cmd/gui/vue/comp"
	"text/template"

	"github.com/p9c/pod/cmd/gui/vue/db"
)

func CoreJs(d db.DuoVUEdb) string {
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
	err := jsRaw.ExecuteTemplate(&js, "vueapp", comp.Components(d))
	if err != nil {
		fmt.Println("error binding to webview:", err)
	}

	return js.String()
}
