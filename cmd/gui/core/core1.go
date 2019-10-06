package core

import (
	"github.com/p9c/pod/cmd/gui/db"
	"github.com/p9c/pod/cmd/gui/mod"
)

func CoreJs(d db.DuOSdb) string {
	return `
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

const core = new Vue({ 
	data () { return { 
	duoSystem }},
});`
}

func CompJs(cs mod.DuOScomps) (s string) {

	for _, c := range cs {

		cc := `var` + c.ID + ` = new Vue({
	el: '#` + c.ID + `',
	name: '` + c.ID + `',
	` + c.Js + `});`
		s = s + cc
	}
	return
}
