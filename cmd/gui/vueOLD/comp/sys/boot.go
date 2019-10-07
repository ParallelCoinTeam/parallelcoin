package sys

import "github.com/p9c/pod/cmd/gui/vue/mod"

func Boot() mod.DuoVUEcomp {
	return mod.DuoVUEcomp{
		IsApp:    true,
		Name:     "BOOT",
		ID:       "boot",
		Version:  "0.0.1",
		CompType: "core",
		SubType:  "boot",
		Js:       `
	data () { return { 
	duoSystem }},
mounted: function(){
this.btnSClick();
},
	methods:{
        btnSClick: function(event) {
			this.duoSystem.bios.isBoot = false;
		}
	}
		`,
		Template: `<boot class="swrap boot" v-show="this.duoSystem.bios.isBoot">
<h1>Boot</h1>
 				<ejs-button cssClass='e-link' v-on:click.native='btnSClick'>Start!</ejs-button>


</boot>`,
		Css: `
		.boot{
			background:red;
		}
		`,
	}
}
